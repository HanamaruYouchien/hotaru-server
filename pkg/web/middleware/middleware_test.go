// FILE: pkg/web/middleware/middleware_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMaxBodyLength(t *testing.T) {
	handler := MaxBodyLength(10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{"Body within limit", "12345", http.StatusOK},
		{"Body exceeds limit", "12345678901", http.StatusRequestEntityTooLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRateLimiter(t *testing.T) {
	handler := RateLimiter(1, time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// First request should pass
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, w.Code)
	}
	print(w.Body.String())

	// Second request should be rate limited
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status %v, got %v", http.StatusTooManyRequests, w.Code)
	}
	print(w.Body.String())

}
