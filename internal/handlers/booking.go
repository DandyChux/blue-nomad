package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/dandychux/blue-nomad/internal/services"
)

type BookingHandler struct {
	square *services.SquareClient
	flow   *services.BookingFlowService
}

func NewBookingHandler(square *services.SquareClient, flow *services.BookingFlowService) *BookingHandler {
	return &BookingHandler{
		square: square,
		flow:   flow,
	}
}

func (h *BookingHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	data, err := h.square.GetBookingServices(r.Context())
	if err != nil {
		slog.Error("failed to fetch booking services", "error", err)
		http.Error(w, "Failed to fetch services", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *BookingHandler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	var req services.SearchAvailabilityRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid availability request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.ServiceVariationID == "" || req.StartAt == "" {
		http.Error(w, "Missing service_variation_id or start_at", http.StatusBadRequest)
		return
	}

	data, err := h.square.SearchAvailability(r.Context(), req)
	if err != nil {
		slog.Error("failed to search square availability", "error", err)
		http.Error(w, "Failed to check availability", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *BookingHandler) GetPaymentConfig(w http.ResponseWriter, r *http.Request) {
	appID := os.Getenv("SQUARE_APPLICATION_ID")
	locationID := os.Getenv("SQUARE_LOCATION_ID")

	if appID == "" || locationID == "" {
		http.Error(w, "Square payment configuration is missing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"application_id": appID,
		"location_id":    locationID,
	})
}

func (h *BookingHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req services.CreateBookingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid booking request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.ServiceVariationID == "" || req.TeamMemberID == "" || req.StartAt == "" || req.EmailAddress == "" || req.PhoneNumber == "" {
		http.Error(w, "Missing required booking details", http.StatusBadRequest)
		return
	}
	if req.ServiceName == "" || req.PriceCents < 0 {
		http.Error(w, "Missing service name or price", http.StatusBadRequest)
		return
	}

	result, err := h.flow.CreateRequest(r.Context(), req)
	if err != nil {
		slog.Error("failed to create local booking request", "error", err)
		http.Error(w, "Could not create booking request", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"request_id": result.ID,
		"status":     result.Status,
	})
}

func (h *BookingHandler) StoreCardAndBook(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("id")
	if requestID == "" {
		http.Error(w, "Missing booking request id", http.StatusBadRequest)
		return
	}

	var req struct {
		SourceID          string `json:"source_id"`
		VerificationToken string `json:"verification_token,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid store-card request body", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.SourceID == "" {
		http.Error(w, "Missing card source_id", http.StatusBadRequest)
		return
	}

	result, err := h.flow.StoreCardAndBook(r.Context(), services.StoreCardAndBookInput{
		BookingRequestID:  requestID,
		SourceID:          req.SourceID,
		VerificationToken: req.VerificationToken,
	})
	if err != nil {
		switch {
		case err == services.ErrBookingRequestNotFound:
			http.Error(w, "Booking request not found", http.StatusNotFound)
		case err == services.ErrSlotNoLongerAvailable:
			http.Error(w, "That time slot is no longer available. Please pick another.", http.StatusConflict)
		default:
			slog.Error("failed to store card on file and create booking", "error", err, "request_id", requestID)
			http.Error(w, fmt.Sprintf("failed to store card on file and create booking: %v", err), http.StatusUnprocessableEntity)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"request_id":     result.ID,
		"booking_id":     result.SquareBookingID,
		"status":         result.Status,
		"booking_status": result.SquareBookingStatus,
	})
}

func (h *BookingHandler) BookFreeService(
	w http.ResponseWriter,
	r *http.Request,
) {
	requestID := r.PathValue("id")
	if requestID == "" {
		http.Error(w, "Missing booking request id", http.StatusBadRequest)
		return
	}

	result, err := h.flow.BookFreeService(r.Context(), requestID)
	if err != nil {
		switch {
		case err == services.ErrBookingRequestNotFound:
			http.Error(w, "Booking request not found", http.StatusNotFound)

		case err == services.ErrSlotNoLongerAvailable:
			http.Error(
				w,
				"That time slot is no longer available. Please pick another.",
				http.StatusConflict,
			)

		default:
			slog.Error(
				"failed to create free booking",
				"error", err,
				"request_id", requestID,
			)
			http.Error(
				w,
				fmt.Sprintf("failed to create free booking: %v", err),
				http.StatusUnprocessableEntity,
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"request_id":     result.ID,
		"booking_id":     result.SquareBookingID,
		"status":         result.Status,
		"booking_status": result.SquareBookingStatus,
	})
}
