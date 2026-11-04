package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBookingHandler_GetPaymentConfig(t *testing.T) {
	h := &BookingHandler{}

	t.Run("missing environment variables returns 500", func(t *testing.T) {
		t.Setenv("SQUARE_APPLICATION_ID", "")
		t.Setenv("SQUARE_LOCATION_ID", "")

		req := httptest.NewRequest(http.MethodGet, "/booking/config", nil)
		rr := httptest.NewRecorder()

		h.GetPaymentConfig(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 InternalServerError, got %d", rr.Code)
		}
	})

	t.Run("configured environment variables return 200 JSON", func(t *testing.T) {
		t.Setenv("SQUARE_APPLICATION_ID", "sq-app-id-test")
		t.Setenv("SQUARE_LOCATION_ID", "sq-loc-id-test")

		req := httptest.NewRequest(http.MethodGet, "/booking/config", nil)
		rr := httptest.NewRecorder()

		h.GetPaymentConfig(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}

		var resp map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp["application_id"] != "sq-app-id-test" || resp["location_id"] != "sq-loc-id-test" {
			t.Errorf("unexpected config response: %+v", resp)
		}
	})
}

func TestBookingHandler_GetAvailability_Validation(t *testing.T) {
	h := &BookingHandler{}

	tests := []struct {
		name           string
		body           string
		wantStatusCode int
	}{
		{
			name:           "invalid json body",
			body:           `{bad json}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "missing service_variation_id",
			body:           `{"start_at":"2026-05-01T10:00:00Z"}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "missing start_at",
			body:           `{"service_variation_id":"var_123"}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/booking/availability", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			h.GetAvailability(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("GetAvailability() status = %d, want %d", rr.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestBookingHandler_CreateRequest_Validation(t *testing.T) {
	h := &BookingHandler{}

	tests := []struct {
		name           string
		body           string
		wantStatusCode int
	}{
		{
			name:           "malformed json",
			body:           `{"invalid"`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing required contact info",
			body: `{
				"service_variation_id": "var_1",
				"team_member_id": "tm_1",
				"start_at": "2026-05-01T10:00:00Z",
				"service_name": "Facial",
				"price_cents": 15000
			}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing service name",
			body: `{
				"service_variation_id": "var_1",
				"team_member_id": "tm_1",
				"start_at": "2026-05-01T10:00:00Z",
				"email_address": "test@example.com",
				"phone_number": "+1234567890",
				"service_name": "",
				"price_cents": 15000
			}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/booking/request", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			h.CreateRequest(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("CreateRequest() status = %d, want %d", rr.Code, tt.wantStatusCode)
			}
		})
	}
}
