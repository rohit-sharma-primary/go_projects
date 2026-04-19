package repository

import (
	"sync"
	"time"

	"server/internal/model"
)

type OrderRepository struct {
	mu     sync.RWMutex
	nextID int
	orders map[int]model.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		nextID: 1,
		orders: make(map[int]model.Order),
	}
}

func (r *OrderRepository) Create(order model.Order) model.Order {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = r.nextID
	order.CreatedAt = time.Now().UTC()
	r.nextID++
	r.orders[order.ID] = order

	return order
}

func (r *OrderRepository) GetByID(id int) (model.Order, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	return order, ok
}

func (r *OrderRepository) List(userID *int) []model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]model.Order, 0, len(r.orders))
	for _, order := range r.orders {
		if userID != nil && order.UserID != *userID {
			continue
		}
		orders = append(orders, order)
	}

	return orders
}
