package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		authHeader   string
		wantStatus   int
		wantExecuted bool
	}{
		{
			name:         "allows valid bearer token",
			authHeader:   "Bearer " + SecretToken,
			wantStatus:   http.StatusNoContent,
			wantExecuted: true,
		},
		{
			name:       "rejects missing token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects malformed token",
			authHeader: SecretToken,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			executed := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				executed = true
				w.WriteHeader(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			HandleAuth(next).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if executed != tt.wantExecuted {
				t.Fatalf("expected executed=%t, got %t", tt.wantExecuted, executed)
			}
		})
	}
}

func TestRecoverPanics(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	RecoverPanics(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Fatalf("expected internal server error message, got %q", body["error"])
	}
}
