package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dandychux/blue-nomad/internal/services"
)

// ── Request type ───────────────────────────────────────────────────────────────

type DiagnosticRequest struct {
	Email            string   `json:"email"`
	FirstName        string   `json:"firstName"`
	SkinType         string   `json:"skinType"`
	SkinTypeOther    string   `json:"skinTypeOther,omitempty"`
	Products         []string `json:"products"`
	ProductsOther    string   `json:"productsOther,omitempty"`
	Ingredients      []string `json:"ingredients"`
	IngredientsOther string   `json:"ingredientsOther,omitempty"`
	SpecificProducts string   `json:"specificProducts,omitempty"`
	Concerns         []string `json:"concerns"`
	ConcernsOther    string   `json:"concernsOther,omitempty"`
}

// ── Handler ────────────────────────────────────────────────────────────────────

type DiagnosticHandler struct {
	newsletter *NewsletterHandler
}

func NewDiagnosticHandler(newsletter *NewsletterHandler) *DiagnosticHandler {
	return &DiagnosticHandler{newsletter: newsletter}
}

// Submit handles POST /diagnostic
func (h *DiagnosticHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req DiagnosticRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid diagnostic request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Persist to HubSpot synchronously — that's the durable record.
	properties := map[string]any{
		"firstname": req.FirstName,
	}
	if _, err := h.newsletter.AddContact(r.Context(), req.Email, properties); err != nil {
		// Soft-fail: HubSpot is the source of truth, but a transient failure
		// shouldn't block the user. Log loudly and move on.
		slog.Warn("failed to add diagnostic contact to HubSpot",
			"error", err, "email", req.Email)
	}

	// Fire the admin notification email asynchronously so SMTP latency
	// (or outright failure) doesn't block the user's response.
	go func(req DiagnosticRequest) {
		if _, err := h.sendNotificationEmail(&req); err != nil {
			slog.Error("failed to send diagnostic email",
				"error", err,
				"email", req.Email,
			)
		}
	}(req)

	slog.Info("diagnostic submitted", "email", req.Email, "name", req.FirstName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Diagnostic submitted successfully",
	})
}

// ── Email builder ──────────────────────────────────────────────────────────────

func (h *DiagnosticHandler) sendNotificationEmail(req *DiagnosticRequest) ([]error, error) {
	// Resolve skin type display value
	skinType := req.SkinType
	if skinType == "Other" && req.SkinTypeOther != "" {
		skinType = fmt.Sprintf("Other — %s", req.SkinTypeOther)
	}

	// Resolve products
	products := append([]string{}, req.Products...)
	if req.ProductsOther != "" {
		for i, p := range products {
			if p == "Other" {
				products[i] = fmt.Sprintf("Other — %s", req.ProductsOther)
			}
		}
	}

	// Resolve ingredients
	ingredients := append([]string{}, req.Ingredients...)
	if req.IngredientsOther != "" {
		for i, ing := range ingredients {
			if ing == "Other" {
				ingredients[i] = fmt.Sprintf("Other — %s", req.IngredientsOther)
			}
		}
	}

	// Resolve concerns
	concerns := append([]string{}, req.Concerns...)
	if req.ConcernsOther != "" {
		for i, c := range concerns {
			if c == "Other" {
				concerns[i] = fmt.Sprintf("Other — %s", req.ConcernsOther)
			}
		}
	}

	// Resolve specific products
	specificProducts := req.SpecificProducts
	if specificProducts == "" {
		specificProducts = "Not provided"
	}

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: 'Helvetica Neue', Arial, sans-serif; color: #1a1a1a; line-height: 1.6; margin: 0; padding: 0; background-color: #f7f5ef; }
				.container { max-width: 600px; margin: 40px auto; }
				.header { border-bottom: 2px solid #1a1a1a; padding-bottom: 16px; margin-bottom: 32px; }
				.header h1 { font-size: 18px; text-transform: uppercase; letter-spacing: 3px; margin: 0 0 4px 0; font-weight: 600; }
				.header p { font-size: 11px; text-transform: uppercase; letter-spacing: 2px; color: #777; margin: 0; }
				.section { margin-bottom: 28px; }
				.label { font-size: 10px; text-transform: uppercase; letter-spacing: 2px; color: #888; margin-bottom: 6px; }
				.value { font-size: 14px; }
				.list { padding-left: 0; list-style: none; }
				.list li { padding: 4px 0; font-size: 14px; }
				.list li::before { content: "— "; color: #aaa; }
				.divider { height: 1px; background-color: #ddd; margin: 24px 0; }
				.footer { font-size: 10px; text-transform: uppercase; letter-spacing: 2px; color: #aaa; margin-top: 40px; padding-top: 16px; border-top: 1px solid #ddd; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Skin Health Diagnostic</h1>
					<p>New submission received</p>
				</div>

				<div class="section">
					<div class="label">Name</div>
					<div class="value">%s</div>
				</div>

				<div class="section">
					<div class="label">Email</div>
					<div class="value"><a href="mailto:%s" style="color: #1a1a1a;">%s</a></div>
				</div>

				<div class="divider"></div>

				<div class="section">
					<div class="label">Skin Type</div>
					<div class="value">%s</div>
				</div>

				<div class="divider"></div>

				<div class="section">
					<div class="label">Current Products</div>
					<ul class="list">%s</ul>
				</div>

				<div class="section">
					<div class="label">Current Ingredients</div>
					<ul class="list">%s</ul>
				</div>

				<div class="section">
					<div class="label">Specific Products & Brands</div>
					<div class="value">%s</div>
				</div>

				<div class="divider"></div>

				<div class="section">
					<div class="label">Skin Concerns</div>
					<ul class="list">%s</ul>
				</div>

				<div class="footer">
					Blue Nomad · 1123 Broadway, #1014 · New York, NY 10010
				</div>
			</div>
		</body>
		</html>
	`,
		req.FirstName,
		req.Email, req.Email,
		skinType,
		buildHTMLList(products),
		buildHTMLList(ingredients),
		specificProducts,
		buildHTMLList(concerns),
	)

	return services.SendEmail([]*services.EmailMessage{
		{
			To:       []string{"hello@bluenomad.nyc"},
			Subject:  fmt.Sprintf("Express Skin Diagnostic — %s", req.FirstName),
			BodyHTML: htmlBody,
		},
	})
}

// buildHTMLList turns a string slice into <li> elements
func buildHTMLList(items []string) string {
	var b strings.Builder
	for _, item := range items {
		b.WriteString(fmt.Sprintf("<li>%s</li>", item))
	}
	return b.String()
}
