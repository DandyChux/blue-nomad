package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildHTMLList(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  string
	}{
		{
			name:  "empty list",
			items: []string{},
			want:  "",
		},
		{
			name:  "single item",
			items: []string{"Cleanser"},
			want:  "<li>Cleanser</li>",
		},
		{
			name:  "multiple items",
			items: []string{"Cleanser", "Moisturizer", "SPF 50"},
			want:  "<li>Cleanser</li><li>Moisturizer</li><li>SPF 50</li>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildHTMLList(tt.items)
			if got != tt.want {
				t.Errorf("buildHTMLList() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDiagnosticHandler_Submit_InvalidBody(t *testing.T) {
	h := &DiagnosticHandler{}

	req := httptest.NewRequest(http.MethodPost, "/diagnostic", strings.NewReader("invalid-json"))
	rr := httptest.NewRecorder()

	h.Submit(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rr.Code)
	}
}
