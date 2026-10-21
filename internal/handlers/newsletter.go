package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dandychux/blue-nomad/internal/services"
)

// ── Request types ──────────────────────────────────────────────────────────────

type SubscribeRequest struct {
	Email      string         `json:"email"`
	Properties map[string]any `json:"properties,omitempty"`
}

type SendMailRequest struct {
	Email   string `json:"email"`
	SendTo  string `json:"sendTo,omitempty"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html,omitempty"`
}

// ── Handler ────────────────────────────────────────────────────────────────────

// NewsletterHandler holds dependencies for newsletter-related HTTP handlers
type NewsletterHandler struct {
	hubspot *services.HTTPClient
	apiKey  string
	listID  string
}

func NewNewsletterHandler() *NewsletterHandler {
	apiKey := os.Getenv("HUBSPOT_API_KEY")
	listID := os.Getenv("HUBSPOT_LIST_ID")

	client := services.NewHTTPClientWithConfig(services.HTTPClientConfig{
		BaseURL:    "https://api.hubapi.com",
		Timeout:    15 * time.Second,
		RetryCount: 2,
		RetryDelay: 1 * time.Second,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if apiKey != "" {
		client.SetBearerToken(apiKey)
	}

	return &NewsletterHandler{
		hubspot: client,
		apiKey:  apiKey,
		listID:  listID,
	}
}

// ── HTTP handlers ──────────────────────────────────────────────────────────────

func (h *NewsletterHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var payload SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if h.apiKey == "" {
		http.Error(w, "HUBSPOT_API_KEY not set", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	result, err := h.addToHubSpot(ctx, payload.Email, payload.Properties)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error with HubSpot: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *NewsletterHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	var req SendMailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sendTo := req.SendTo
	if sendTo == "" {
		sendTo = os.Getenv("EMAIL_USERNAME")
	}

	msg := []*services.EmailMessage{
		{
			To:       []string{sendTo},
			Subject:  req.Subject,
			BodyHTML: req.HTML,
			BodyText: req.Text,
		},
	}

	if _, err := services.SendEmail(msg); err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email sent successfully",
		"to":      sendTo,
	})
}

// ── HubSpot helpers ────────────────────────────────────────────────────────────

func (h *NewsletterHandler) addToHubSpot(ctx context.Context, email string, properties map[string]any) (map[string]any, error) {
	contactProperties := map[string]any{
		"email":          email,
		"hs_lead_status": "NEW",
		"lifecyclestage": "subscriber",
	}
	for k, v := range properties {
		contactProperties[k] = v
	}

	body := map[string]any{"properties": contactProperties}

	resp, err := h.hubspot.Post(ctx, "/crm/v3/objects/contacts", body)

	// Creation succeeded
	if err == nil && resp.IsSuccess() {
		var contact map[string]any
		if err := resp.JSON(&contact); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}

		if h.listID != "" {
			if contactID, ok := contact["id"].(string); ok {
				h.addContactToList(ctx, contactID)
			}
		}

		return map[string]any{
			"success": true,
			"contact": contact,
			"created": true,
		}, nil
	}

	// Contact already exists — update instead
	if resp != nil && resp.StatusCode == 409 {
		return h.handleExistingContact(ctx, email, properties)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}
	return nil, fmt.Errorf("failed to create contact: %s", resp.String())
}

func (h *NewsletterHandler) handleExistingContact(ctx context.Context, email string, properties map[string]any) (map[string]any, error) {
	searchBody := map[string]any{
		"filterGroups": []map[string]any{
			{
				"filters": []map[string]any{
					{
						"propertyName": "email",
						"operator":     "EQ",
						"value":        email,
					},
				},
			},
		},
	}

	searchResp, err := h.hubspot.Post(ctx, "/crm/v3/objects/contacts/search", searchBody)
	if err != nil || !searchResp.IsSuccess() {
		return nil, fmt.Errorf("failed to search for existing contact: %w", err)
	}

	var searchData map[string]any
	if err := searchResp.JSON(&searchData); err != nil {
		return nil, fmt.Errorf("error decoding search response: %w", err)
	}

	results, ok := searchData["results"].([]any)
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("contact not found in search results")
	}

	firstResult, ok := results[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid search result format")
	}

	contactID, ok := firstResult["id"].(string)
	if !ok {
		return nil, fmt.Errorf("contact ID not found")
	}

	updateProperties := map[string]any{
		"lifecyclestage": "subscriber",
	}
	for k, v := range properties {
		updateProperties[k] = v
	}

	updateResp, err := h.hubspot.Patch(ctx, fmt.Sprintf("/crm/v3/objects/contacts/%s", contactID), map[string]any{"properties": updateProperties})
	if err != nil || !updateResp.IsSuccess() {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}

	var contact map[string]any
	if err := updateResp.JSON(&contact); err != nil {
		return nil, fmt.Errorf("error decoding update response: %w", err)
	}

	if h.listID != "" {
		h.addContactToList(ctx, contactID)
	}

	return map[string]any{
		"success": true,
		"contact": contact,
		"updated": true,
	}, nil
}

func (h *NewsletterHandler) addContactToList(ctx context.Context, contactID string) {
	if h.apiKey == "" {
		return
	}

	resp, err := h.hubspot.Request(
		ctx,
		"PUT",
		fmt.Sprintf("/crm/v3/lists/%s/memberships/add", h.listID),
		nil,
		[]string{contactID},
		nil,
	)

	if err != nil || !resp.IsSuccess() {
		fmt.Printf("Failed to add contact to list: %v\n", err)
		return
	}

	fmt.Printf("Contact %s added to list %s\n", contactID, h.listID)
}

// AddContact is a public wrapper around addToHubSpot so other handlers can reuse it.
func (h *NewsletterHandler) AddContact(ctx context.Context, email string, properties map[string]any) (map[string]any, error) {
	return h.addToHubSpot(ctx, email, properties)
}
