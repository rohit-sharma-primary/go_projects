package service

import (
	"errors"
	"testing"

	"server/internal/model"
	"server/internal/repository"
)

func TestUserServiceCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     model.CreateUserRequest
		wantErr error
	}{
		{
			name: "creates user",
			req: model.CreateUserRequest{
				Name:  "Alice",
				Email: "alice@example.com",
			},
		},
		{
			name: "rejects empty name",
			req: model.CreateUserRequest{
				Name:  " ",
				Email: "alice@example.com",
			},
			wantErr: ErrInvalidInput,
		},
		{
			name: "rejects invalid email",
			req: model.CreateUserRequest{
				Name:  "Alice",
				Email: "aliceexample.com",
			},
			wantErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewUserService(repository.NewUserRepository())
			user, err := svc.CreateUser(tt.req)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if user.ID != 1 {
					t.Fatalf("expected generated id 1, got %d", user.ID)
				}
				if user.Name != tt.req.Name {
					t.Fatalf("expected name %q, got %q", tt.req.Name, user.Name)
				}
			}
		})
	}
}

func TestUserServiceGetUser(t *testing.T) {
	t.Parallel()

	repo := repository.NewUserRepository()
	created := repo.Create(model.User{Name: "Alice", Email: "alice@example.com"})
	svc := NewUserService(repo)

	user, err := svc.GetUser(created.ID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user.ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, user.ID)
	}

	_, err = svc.GetUser(created.ID + 1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
