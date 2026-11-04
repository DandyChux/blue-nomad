package handlers

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowedImageFormat(t *testing.T) {
	tests := []struct {
		format string
		want   bool
	}{
		{"webp", true},
		{"avif", true},
		{"jpg", false},
		{"jpeg", false},
		{"png", false},
		{"gif", false},
		{"svg", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			if got := allowedImageFormat(tt.format); got != tt.want {
				t.Errorf("allowedImageFormat(%q) = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func TestAllowedImageWidth(t *testing.T) {
	tests := []struct {
		name  string
		width int
		want  bool
	}{
		{"lower bound", 1, true},
		{"standard width", 800, true},
		{"upper bound", 2400, true},
		{"zero width", 0, false},
		{"negative width", -100, false},
		{"exceeds max width", 2401, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowedImageWidth(tt.width); got != tt.want {
				t.Errorf("allowedImageWidth(%d) = %v, want %v", tt.width, got, tt.want)
			}
		})
	}
}

func TestImageHandler_Redirect(t *testing.T) {
	// Setup a test handler with valid dummy hex keys
	keyHex := hex.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	saltHex := hex.EncodeToString([]byte("fedcba9876543210fedcba9876543210"))

	t.Setenv("IMGPROXY_URL", "https://imgproxy.example.com")
	t.Setenv("IMGPROXY_KEY", keyHex)
	t.Setenv("IMGPROXY_SALT", saltHex)
	t.Setenv("IMGPROXY_ALLOWED_HOSTS", "cdn.sanity.io,example.com")

	handler, err := NewImageHandler()
	if err != nil {
		t.Fatalf("failed to create ImageHandler: %v", err)
	}

	tests := []struct {
		name           string
		widthParam     string
		rawQuery       string
		wantStatusCode int
	}{
		{
			name:           "valid image redirect",
			widthParam:     "600",
			rawQuery:       "src=https://cdn.sanity.io/images/photo.jpg&format=webp&quality=80",
			wantStatusCode: http.StatusFound,
		},
		{
			name:           "default format and quality",
			widthParam:     "1200",
			rawQuery:       "src=https://example.com/asset.png",
			wantStatusCode: http.StatusFound,
		},
		{
			name:           "invalid image width (out of range)",
			widthParam:     "5000",
			rawQuery:       "src=https://cdn.sanity.io/images/photo.jpg",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "non-numeric image width",
			widthParam:     "invalid",
			rawQuery:       "src=https://cdn.sanity.io/images/photo.jpg",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "unsupported format",
			widthParam:     "400",
			rawQuery:       "src=https://cdn.sanity.io/images/photo.jpg&format=bmp",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "quality out of bounds",
			widthParam:     "400",
			rawQuery:       "src=https://cdn.sanity.io/images/photo.jpg&quality=150",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "missing src parameter",
			widthParam:     "400",
			rawQuery:       "format=webp",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "disallowed host (SSRF protection)",
			widthParam:     "400",
			rawQuery:       "src=https://evil.com/malicious.jpg",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "non-https scheme disallowed",
			widthParam:     "400",
			rawQuery:       "src=http://cdn.sanity.io/images/photo.jpg",
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/images/"+tt.widthParam+"?"+tt.rawQuery, nil)
			req.SetPathValue("width", tt.widthParam)

			rr := httptest.NewRecorder()
			handler.Redirect(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("Redirect() status = %d, want %d (body: %s)", rr.Code, tt.wantStatusCode, rr.Body.String())
			}

			if tt.wantStatusCode == http.StatusFound {
				loc := rr.Header().Get("Location")
				if loc == "" {
					t.Error("expected Location header on 302 redirect, got empty")
				}
				cacheControl := rr.Header().Get("Cache-Control")
				if cacheControl != "public, max-age=3600" {
					t.Errorf("Cache-Control = %q, want public, max-age=3600", cacheControl)
				}
			}
		})
	}
}
