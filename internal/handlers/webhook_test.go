package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebhookHandler_Health(t *testing.T) {
	h := &WebhookHandler{}

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	rr := httptest.NewRecorder()

	h.Health(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", res["status"])
	}
}

func TestWebhookHandler_Revalidate_Unauthorized(t *testing.T) {
	h := &WebhookHandler{secret: "my-secret-key"}

	req := httptest.NewRequest(http.MethodPost, "/revalidate?secret=wrong-secret", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()

	h.Revalidate(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rr.Code)
	}
}

func TestWebhookHandler_VerifySquareSignature(t *testing.T) {
	signatureKey := "square-sig-key-123"
	webhookURL := "https://example.com/api/webhooks/square"
	body := `{"event_id":"evt_123","type":"booking.created"}`

	h := &WebhookHandler{
		squareSignatureKey: signatureKey,
		squareWebhookURL:   webhookURL,
	}

	// Compute valid signature: HMAC-SHA256(key, webhookURL + body)
	mac := hmac.New(sha256.New, []byte(signatureKey))
	mac.Write([]byte(webhookURL + body))
	validSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name      string
		body      string
		signature string
		want      bool
	}{
		{
			name:      "valid signature",
			body:      body,
			signature: validSignature,
			want:      true,
		},
		{
			name:      "tampered body",
			body:      `{"event_id":"tampered"}`,
			signature: validSignature,
			want:      false,
		},
		{
			name:      "invalid base64 signature",
			body:      body,
			signature: "not-base-64!!",
			want:      false,
		},
		{
			name:      "empty signature",
			body:      body,
			signature: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := h.verifySquareSignature(tt.body, tt.signature)
			if got != tt.want {
				t.Errorf("verifySquareSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
