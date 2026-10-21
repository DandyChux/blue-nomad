package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "time/tzdata"

	"github.com/dandychux/blue-nomad/internal/services"
)

type WebhookPayload struct {
	Type       string `json:"_type,omitempty"`
	Operation  string `json:"operation,omitempty"`
	Event      string `json:"event,omitempty"`
	Collection string `json:"collection,omitempty"`
}

// --- Square Webhook Structs ---
type SquareWebhookPayload struct {
	MerchantID string `json:"merchant_id"`
	Type       string `json:"type"`
	EventID    string `json:"event_id"`
	Data       struct {
		Type   string `json:"type"`
		ID     string `json:"id"`
		Object struct {
			Payment struct {
				ID          string `json:"id"`
				OrderID     string `json:"order_id"`
				ReferenceID string `json:"reference_id"`
				CustomerID  string `json:"customer_id"`
				Status      string `json:"status"`
				AmountMoney struct {
					Amount   int64  `json:"amount"`
					Currency string `json:"currency"`
				} `json:"amount_money"`
			} `json:"payment"`
			Booking struct {
				ID                  string `json:"id"`
				Version             int    `json:"version"`
				Status              string `json:"status"`
				CreatedAt           string `json:"created_at"`
				StartAt             string `json:"start_at"`
				LocationID          string `json:"location_id"`
				CustomerID          string `json:"customer_id"`
				AppointmentSegments []struct {
					DurationMinutes         int    `json:"duration_minutes"`
					ServiceVariationID      string `json:"service_variation_id"`
					TeamMemberID            string `json:"team_member_id"`
					ServiceVariationVersion int64  `json:"service_variation_version"`
				} `json:"appointment_segments"`
			} `json:"booking"`
		} `json:"object"`
	} `json:"data"`
}

type WebhookHandler struct {
	secret             string
	squareSignatureKey string
	squareWebhookURL   string
	sanity             *services.SanityClient
	square             *services.SquareClient
	queue              *services.WebhookQueue
	flow               *services.BookingFlowService

	mu            sync.RWMutex
	lastInvalidAt time.Time
	lastEvent     string
}

// Updated Constructor
func NewWebhookHandler(secret, squareSignatureKey, squareWebhookURL string, sanity *services.SanityClient, square *services.SquareClient, queue *services.WebhookQueue, flow *services.BookingFlowService) *WebhookHandler {
	return &WebhookHandler{
		secret:             secret,
		squareSignatureKey: squareSignatureKey,
		squareWebhookURL:   squareWebhookURL,
		sanity:             sanity,
		square:             square,
		queue:              queue,
		flow:               flow,
	}
}

func (h *WebhookHandler) Revalidate(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	if h.secret == "" || secret != h.secret {
		http.Error(w, "Invalid secret", http.StatusUnauthorized)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := describeEvent(payload)
	slog.Info("Webhook received — invalidating cache", "event", event)

	h.sanity.InvalidateCache()

	h.mu.Lock()
	h.lastInvalidAt = time.Now().UTC()
	h.lastEvent = event
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":        "Revalidation triggered",
		"revalidated":    true,
		"invalidated_at": h.lastInvalidAt.Format(time.RFC3339),
	})
}

func (h *WebhookHandler) LastInvalidation(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	t := h.lastInvalidAt
	evt := h.lastEvent
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, max-age=0")
	json.NewEncoder(w).Encode(map[string]any{
		"invalidated_at": t.Format(time.RFC3339),
		"event":          evt,
	})
}

func (h *WebhookHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HandleSquareWebhook receives payment events from Square
func (h *WebhookHandler) HandleSquareWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read Square webhook body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Square-Hmacsha256-Signature")
	if !h.verifySquareSignature(string(body), signature) {
		slog.Warn("Invalid Square webhook signature. Possible spoofing attempt.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload SquareWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		slog.Error("Failed to unmarshal Square webhook payload", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.EventID == "" || payload.Type == "" {
		slog.Warn("Square webhook missing event metadata",
			"event_id", payload.EventID,
			"event_type", payload.Type,
		)
		http.Error(w, "Invalid Square event", http.StatusBadRequest)
		return
	}

	inserted, err := h.queue.Enqueue(r.Context(), services.EnqueueWebhookEventParams{
		Provider:  services.WebhookProviderSquare,
		EventID:   payload.EventID,
		EventType: payload.Type,
		Payload:   body,
	})
	if err != nil {
		// Important: do NOT return 200 if we failed to persist the webhook.
		// Let Square retry because durability is the point of this queue.
		slog.Error("Failed to persist Square webhook event",
			"error", err,
			"event_id", payload.EventID,
			"event_type", payload.Type,
		)
		http.Error(w, "Webhook persistence unavailable", http.StatusServiceUnavailable)
		return
	}

	if !inserted {
		slog.Info("Duplicate Square webhook ignored",
			"event_id", payload.EventID,
			"event_type", payload.Type,
		)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) ProcessQueuedWebhook(ctx context.Context, provider string, rawPayload []byte) error {
	if provider != services.WebhookProviderSquare {
		return fmt.Errorf("unsupported webhook provider: %s", provider)
	}

	var payload SquareWebhookPayload
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return fmt.Errorf("unmarshal queued Square webhook payload: %w", err)
	}

	switch payload.Type {
	case "payment.created":
		return h.handlePaymentEvent(ctx, payload)
	case "booking.created":
		return h.handleBookingCreated(ctx, payload)
	case "booking.updated":
		return h.handleBookingUpdated(ctx, payload)
	case "catalog.version.updated":
		return h.handleCatalogChanged(ctx, payload)
	default:
		slog.Debug("Unhandled queued Square event ignored", "type", payload.Type, "event_id", payload.EventID)
		return nil
	}
}

func (h *WebhookHandler) handleCatalogChanged(
	_ context.Context,
	payload SquareWebhookPayload,
) error {
	slog.Info(
		"Square catalog changed",
		"event_type", payload.Type,
		"event_id", payload.EventID,
	)

	// h.sanity.InvalidateCache()

	h.mu.Lock()
	h.lastInvalidAt = time.Now().UTC()
	h.lastEvent = payload.Type
	h.mu.Unlock()

	return nil
}

func (h *WebhookHandler) handlePaymentEvent(_ context.Context, payload SquareWebhookPayload) error {
	payment := payload.Data.Object.Payment

	slog.Debug("Ignoring Square payment event for booking flow", "payment_id", payment.ID, "status", payment.Status, "event_type", payload.Type)
	return nil
}

func (h *WebhookHandler) sendBookingNotification(req *services.BookingRequest) ([]error, error) {
	adminEmail := notificationEmail()
	if adminEmail == "" {
		return []error{fmt.Errorf("ADMIN_EMAIL or EMAIL_USERNAME must be configured")}, nil
	}

	formattedDate, formattedTime := formatBookingTime(req.StartAt)
	formattedAmount := fmt.Sprintf("%.2f", float64(req.PriceCents)/100.0)

	msg := []*services.EmailMessage{
		{
			To:      []string{adminEmail},
			Subject: fmt.Sprintf("New Blue Nomad Booking: %s at %s", formattedDate, formattedTime),
			BodyText: fmt.Sprintf(
				"A new appointment has been created.\n\nBooking ID: %s\nStatus: %s\nDate: %s\nTime: %s\nService: %s\nPrice: $%s %s\nCustomer: %s %s\nEmail: %s\nPhone: %s\n\nThe customer's card has been saved on file and the Square booking has been created. No payment has been charged online. Final payment should be collected on site in Square.",
				req.SquareBookingID,
				req.SquareBookingStatus,
				formattedDate,
				formattedTime,
				req.ServiceName,
				formattedAmount,
				req.Currency,
				req.GivenName,
				req.FamilyName,
				req.EmailAddress,
				req.PhoneNumber,
			),
			BodyHTML: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #000; text-transform: uppercase;">New Appointment Request</h2>
				<p>A new appointment has been created via the Blue Nomad website.</p>
				<div style="background-color: #f9fafb; padding: 20px; border-radius: 4px; margin: 20px 0;">
					<p><strong>Date:</strong> %s</p>
					<p><strong>Time:</strong> %s</p>
					<p><strong>Service:</strong> %s</p>
					<p><strong>Booking ID:</strong> %s</p>
					<p><strong>Price:</strong> $%s %s</p>
					<p><strong>Client:</strong> %s %s</p>
					<p><strong>Email:</strong> %s</p>
				</div>
				<p>The customer's card was saved on file during booking.</p>
				<p><a href="https://squareup.com/dashboard/appointments">Open in Square Appointments</a></p>
			</div>
			`, formattedDate, formattedTime, req.ServiceName, req.SquareBookingID, formattedAmount, req.Currency, req.GivenName, req.FamilyName, req.EmailAddress),
		},
		{
			To:       []string{req.EmailAddress},
			Subject:  fmt.Sprintf("Your booking has been received - %s", req.ServiceName),
			BodyText: fmt.Sprintf("Thank you for booking with us. Your appointment is confirmed and we'll see you on %s at %s.", formattedDate, formattedTime),
			BodyHTML: fmt.Sprintf(`
				<div>
					<h2>Your appointment is confirmed.</h2>
					<p>Thank you for booking with us. Your appointment is confirmed and we'll see you soon!</p>
					<div style="background-color: #f9fafb; padding: 20px; border-radius: 4px; margin: 20px 0;">
						<p>Date: %s</p>
						<p>Time: %s</p>
						<p>Service: %s</p>
					</div>
				</div>
			`, formattedDate, formattedTime, req.ServiceName),
		},
	}

	if emailErrs, err := services.SendEmail(msg); err != nil {
		return []error{fmt.Errorf("send booking update email for booking %s: %w", req.ID, err)}, nil
	} else if emailErrs != nil {
		return []error{fmt.Errorf("send booking update email for booking %s: %w", req.ID, emailErrs[0])}, nil
	}

	return nil, nil
}

func (h *WebhookHandler) handleBookingCreated(ctx context.Context, payload SquareWebhookPayload) error {
	booking := payload.Data.Object.Booking

	slog.Info("📅 Booking CREATED",
		"booking_id", booking.ID,
		"status", booking.Status,
		"start_at", booking.StartAt,
		"customer_id", booking.CustomerID,
	)

	req, err := h.flow.GetBySquareBookingID(ctx, booking.ID)
	if err != nil {
		if err == services.ErrBookingRequestNotFound {
			slog.Debug("No local booking request found for booking.created", "booking_id", booking.ID)
			return nil
		}
		return err
	}

	if req.AdminNotifiedAt != nil {
		return nil
	}

	if errs, err := h.sendBookingNotification(req); err != nil {
		return errs[0]
	}

	// if len(errs) > 0 {
	// 	return errs[0]
	// }

	return h.flow.MarkAdminNotified(ctx, req.ID)
}

// handleBookingUpdated notifies the admin when a booking changes state
// (confirmation, cancellation, reschedule, decline, etc.). Square fires
// booking.updated for any of these transitions, so we translate the status
// into a human-readable subject line and include the full context in the body.
func (h *WebhookHandler) handleBookingUpdated(_ context.Context, payload SquareWebhookPayload) error {
	booking := payload.Data.Object.Booking

	formattedDate, formattedTime, err := parseAndFormatBookingTime(booking.StartAt)
	if err != nil {
		slog.Warn("failed to parse booking start time", "start_at", booking.StartAt, "error", err)
		formattedDate = booking.StartAt
	}

	var durationMin int
	var serviceVariationID string
	if len(booking.AppointmentSegments) > 0 {
		durationMin = booking.AppointmentSegments[0].DurationMinutes
		serviceVariationID = booking.AppointmentSegments[0].ServiceVariationID
	}

	switch booking.Status {
	case "ACCEPTED", "DECLINED", "CANCELLED_BY_CUSTOMER", "NO_SHOW":
		// Map Square's status to a friendly subject prefix
		statusLabel := bookingStatusLabel(booking.Status)

		slog.Info("📅 Booking UPDATED",
			"booking_id", booking.ID,
			"status", booking.Status,
			"version", booking.Version,
			"start_at", booking.StartAt,
			"customer_id", booking.CustomerID,
		)

		adminEmail := notificationEmail()
		if adminEmail == "" {
			return fmt.Errorf("ADMIN_EMAIL or EMAIL_USERNAME must be configured")
		}

		subject := fmt.Sprintf("Booking %s: %s at %s", statusLabel, formattedDate, formattedTime)

		msg := []*services.EmailMessage{
			{

				To:      []string{adminEmail},
				Subject: subject,
				BodyText: fmt.Sprintf(
					"Booking %s\n\nBooking ID: %s\nStatus: %s\nDate: %s\nTime: %s\nDuration: %d minutes\nCustomer ID: %s\nService Variation: %s\nVersion: %d\n\nReview in your Square Dashboard for full details.",
					statusLabel, booking.ID, booking.Status, formattedDate, formattedTime,
					durationMin, booking.CustomerID, serviceVariationID, booking.Version,
				),
				BodyHTML: fmt.Sprintf(`
					<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
						<h2 style="color: #000; text-transform: uppercase;">Booking %s</h2>
						<p>An existing appointment was updated on Square.</p>
						<div style="background-color: #f9fafb; padding: 20px; border-radius: 4px; margin: 20px 0;">
							<p><strong>Date:</strong> %s</p>
							<p><strong>Time:</strong> %s</p>
							<p><strong>Duration:</strong> %d minutes</p>
							<p><strong>New Status:</strong> %s</p>
							<p><strong>Booking ID:</strong> %s</p>
							<p><strong>Version:</strong> %d</p>
						</div>
						<p><a href="https://squareup.com/dashboard/appointments">Open in Square Dashboard</a></p>
					</div>
					`, statusLabel, formattedDate, formattedTime, durationMin, booking.Status, booking.ID, booking.Version),
			},
		}

		if _, err := services.SendEmail(msg); err != nil {
			return fmt.Errorf("send booking update email for booking %s: %w", booking.ID, err)
		}

		slog.Info("Booking update notification sent", "booking_id", booking.ID)
		return nil

	default:
		slog.Debug("Skipping non-final booking update", "booking_id", booking.ID, "status", booking.Status)
		return nil
	}
}

// verifySquareSignature mathematically proves the request came from Square
func (h *WebhookHandler) verifySquareSignature(body, signature string) bool {
	if h.squareSignatureKey == "" || signature == "" {
		return false
	}

	message := h.squareWebhookURL + body

	mac := hmac.New(sha256.New, []byte(h.squareSignatureKey))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	providedMAC, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedMAC, providedMAC)
}

const defaultBusinessTimezone = "America/New_York"

func bookingLocation() *time.Location {
	tz := os.Getenv("BUSINESS_TIMEZONE")
	if tz == "" {
		tz = defaultBusinessTimezone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		slog.Error("failed to load booking timezone", "timezone", tz, "error", err)
		return time.UTC
	}

	return loc
}

func formatBookingTime(t time.Time) (string, string) {
	local := t.In(bookingLocation())
	return local.Format("Monday, January 2, 2006"), local.Format("3:04 PM MST")
}

func parseAndFormatBookingTime(raw string) (string, string, error) {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw, "", err
	}

	date, clock := formatBookingTime(t)
	return date, clock, nil
}

// bookingStatusLabel maps Square booking statuses to human-friendly text.
func bookingStatusLabel(status string) string {
	switch status {
	case "ACCEPTED":
		return "Confirmed"
	case "PENDING":
		return "Pending"
	case "DECLINED":
		return "Declined"
	case "CANCELLED_BY_CUSTOMER":
		return "Cancelled (by customer)"
	case "CANCELLED_BY_SELLER":
		return "Cancelled (by seller)"
	case "NO_SHOW":
		return "No-show"
	default:
		if status == "" {
			return "Updated"
		}
		return status
	}
}

func describeEvent(p WebhookPayload) string {
	if p.Collection != "" && p.Event != "" {
		return fmt.Sprintf("directus:%s/%s", p.Collection, p.Event)
	}
	if p.Type != "" && p.Operation != "" {
		return fmt.Sprintf("sanity:%s/%s", p.Type, p.Operation)
	}
	return "unknown"
}

func bookingIDFromReference(reference string) (string, bool) {
	const prefix = "booking:"
	if !strings.HasPrefix(reference, prefix) {
		return "", false
	}

	bookingID := strings.TrimPrefix(reference, prefix)
	if bookingID == "" {
		return "", false
	}

	return bookingID, true
}

func notificationEmail() string {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail != "" {
		return adminEmail
	}
	return os.Getenv("EMAIL_USERNAME")
}
