package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShopHandler_CreateCheckoutLink_Validation(t *testing.T) {
	h := &ShopHandler{}

	tests := []struct {
		name           string
		body           string
		wantStatusCode int
	}{
		{
			name:           "invalid json payload",
			body:           `{not a json}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "empty cart items",
			body:           `{"items": []}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			h.CreateCheckoutLink(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("CreateCheckoutLink() status = %d, want %d", rr.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestShopHandler_GetProductImage_MissingID(t *testing.T) {
	h := &ShopHandler{}

	req := httptest.NewRequest(http.MethodGet, "/shop/images/", nil)
	// path parameter imageID is omitted
	rr := httptest.NewRecorder()

	h.GetProductImage(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request when imageID is empty, got %d", rr.Code)
	}
}
