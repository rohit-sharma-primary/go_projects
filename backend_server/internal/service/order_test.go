package service

import (
	"errors"
	"testing"

	"server/internal/model"
	"server/internal/repository"
)

func TestOrderServiceCreateOrder(t *testing.T) {
	t.Parallel()

	userRepo := repository.NewUserRepository()
	createdUser := userRepo.Create(model.User{Name: "Alice", Email: "alice@example.com"})

	tests := []struct {
		name    string
		req     model.CreateOrderRequest
		wantErr error
	}{
		{
			name: "creates order",
			req: model.CreateOrderRequest{
				UserID: createdUser.ID,
				Amount: 125.50,
			},
		},
		{
			name: "rejects missing user",
			req: model.CreateOrderRequest{
				UserID: createdUser.ID + 99,
				Amount: 125.50,
			},
			wantErr: ErrInvalidInput,
		},
		{
			name: "rejects zero amount",
			req: model.CreateOrderRequest{
				UserID: createdUser.ID,
				Amount: 0,
			},
			wantErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewOrderService(repository.NewOrderRepository(), userRepo)
			order, err := svc.CreateOrder(tt.req)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if order.ID != 1 {
					t.Fatalf("expected generated id 1, got %d", order.ID)
				}
				if order.Status != "created" {
					t.Fatalf("expected status created, got %s", order.Status)
				}
			}
		})
	}
}

func TestOrderServiceListOrders(t *testing.T) {
	t.Parallel()

	userRepo := repository.NewUserRepository()
	userOne := userRepo.Create(model.User{Name: "Alice", Email: "alice@example.com"})
	userTwo := userRepo.Create(model.User{Name: "Bob", Email: "bob@example.com"})

	orderRepo := repository.NewOrderRepository()
	svc := NewOrderService(orderRepo, userRepo)

	if _, err := svc.CreateOrder(model.CreateOrderRequest{UserID: userOne.ID, Amount: 10}); err != nil {
		t.Fatalf("create order for userOne: %v", err)
	}
	if _, err := svc.CreateOrder(model.CreateOrderRequest{UserID: userTwo.ID, Amount: 20}); err != nil {
		t.Fatalf("create order for userTwo: %v", err)
	}

	allOrders := svc.ListOrders(nil)
	if len(allOrders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(allOrders))
	}

	filterID := userOne.ID
	filteredOrders := svc.ListOrders(&filterID)
	if len(filteredOrders) != 1 {
		t.Fatalf("expected 1 filtered order, got %d", len(filteredOrders))
	}
	if filteredOrders[0].UserID != userOne.ID {
		t.Fatalf("expected user id %d, got %d", userOne.ID, filteredOrders[0].UserID)
	}
}
