package service

import (
	"server/internal/model"
	"server/internal/repository"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
	userRepo  *repository.UserRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, userRepo *repository.UserRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

func (s *OrderService) CreateOrder(req model.CreateOrderRequest) (model.Order, error) {
	if req.UserID <= 0 || req.Amount <= 0 {
		return model.Order{}, ErrInvalidInput
	}

	if _, ok := s.userRepo.GetByID(req.UserID); !ok {
		return model.Order{}, ErrInvalidInput
	}

	return s.orderRepo.Create(model.Order{
		UserID: req.UserID,
		Amount: req.Amount,
		Status: "created",
	}), nil
}

func (s *OrderService) GetOrder(id int) (model.Order, error) {
	order, ok := s.orderRepo.GetByID(id)
	if !ok {
		return model.Order{}, ErrNotFound
	}
	return order, nil
}

func (s *OrderService) ListOrders(userID *int) []model.Order {
	return s.orderRepo.List(userID)
}
