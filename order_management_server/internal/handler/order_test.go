package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"server/internal/model"
	"server/internal/repository"
	"server/internal/service"
)

func newOrderTestMux() (*http.ServeMux, *service.UserService, *service.OrderService) {
	mux := http.NewServeMux()
	userRepo := repository.NewUserRepository()
	orderRepo := repository.NewOrderRepository()
	userSvc := service.NewUserService(userRepo)
	orderSvc := service.NewOrderService(orderRepo, userRepo)
	NewOrderHandler(orderSvc).Register(mux)
	return mux, userSvc, orderSvc
}

func TestPostOrderHandler(t *testing.T) {
	t.Parallel()

	t.Run("creates order for existing user", func(t *testing.T) {
		t.Parallel()

		mux, userSvc, _ := newOrderTestMux()
		user, err := userSvc.CreateUser(model.CreateUserRequest{Name: "Alice", Email: "alice@example.com"})
		if err != nil {
			t.Fatalf("seed user: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"user_id":`+strconv.Itoa(user.ID)+`,"amount":99.99}`))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d", rec.Code)
		}

		var got model.Order
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if got.UserID != user.ID {
			t.Fatalf("expected user id %d, got %d", user.ID, got.UserID)
		}
	})

	t.Run("rejects invalid order", func(t *testing.T) {
		t.Parallel()

		mux, _, _ := newOrderTestMux()
		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"user_id":999,"amount":0}`))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})
}

func TestGetOrdersHandlerFiltersByUserID(t *testing.T) {
	t.Parallel()

	mux, userSvc, orderSvc := newOrderTestMux()
	userOne, err := userSvc.CreateUser(model.CreateUserRequest{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		t.Fatalf("seed userOne: %v", err)
	}
	userTwo, err := userSvc.CreateUser(model.CreateUserRequest{Name: "Bob", Email: "bob@example.com"})
	if err != nil {
		t.Fatalf("seed userTwo: %v", err)
	}

	if _, err := orderSvc.CreateOrder(model.CreateOrderRequest{UserID: userOne.ID, Amount: 10}); err != nil {
		t.Fatalf("seed order one: %v", err)
	}
	if _, err := orderSvc.CreateOrder(model.CreateOrderRequest{UserID: userTwo.ID, Amount: 20}); err != nil {
		t.Fatalf("seed order two: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/orders?user_id="+strconv.Itoa(userOne.ID), nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got []model.Order
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 order, got %d", len(got))
	}
	if got[0].UserID != userOne.ID {
		t.Fatalf("expected user id %d, got %d", userOne.ID, got[0].UserID)
	}
}
