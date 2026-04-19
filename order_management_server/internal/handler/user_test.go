package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"server/internal/model"
	"server/internal/repository"
	"server/internal/service"
)

func newUserTestMux() *http.ServeMux {
	mux := http.NewServeMux()
	repo := repository.NewUserRepository()
	svc := service.NewUserService(repo)
	NewUserHandler(svc).Register(mux)
	return mux
}

func TestPostUserHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "creates user",
			body:       `{"name":"Alice","email":"alice@example.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "rejects invalid email",
			body:       `{"name":"Alice","email":"aliceexample.com"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			newUserTestMux().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantStatus == http.StatusCreated {
				var got model.User
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if got.ID != 1 {
					t.Fatalf("expected id 1, got %d", got.ID)
				}
			}
		})
	}
}

func TestGetUsersHandler(t *testing.T) {
	t.Parallel()

	repo := repository.NewUserRepository()
	svc := service.NewUserService(repo)
	handler := NewUserHandler(svc)
	mux := http.NewServeMux()
	handler.Register(mux)

	if _, err := svc.CreateUser(model.CreateUserRequest{Name: "Alice", Email: "alice@example.com"}); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got []model.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 user, got %d", len(got))
	}
}
