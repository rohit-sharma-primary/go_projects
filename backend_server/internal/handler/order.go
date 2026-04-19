package handler

import (
	"errors"
	"net/http"
	"strconv"

	"server/internal/model"
	"server/internal/service"
	"server/internal/utils"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /orders", h.getOrdersHandler)
	mux.HandleFunc("GET /orders/{id}", h.getOrderByIDHandler)
	mux.HandleFunc("POST /orders", h.postOrderHandler)
}

func (h *OrderHandler) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var userID *int

	if rawUserID := r.URL.Query().Get("user_id"); rawUserID != "" {
		parsedUserID, err := strconv.Atoi(rawUserID)
		if err != nil || parsedUserID <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &parsedUserID
	}

	utils.WriteJSON(w, http.StatusOK, h.service.ListOrders(userID))
}

func (h *OrderHandler) getOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, svcErr := h.service.GetOrder(id)
	if errors.Is(svcErr, service.ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, "order not found")
		return
	}
	if svcErr != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch order")
		return
	}

	utils.WriteJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) postOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req model.CreateOrderRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, svcErr := h.service.CreateOrder(req)
	if errors.Is(svcErr, service.ErrInvalidInput) {
		utils.WriteError(w, http.StatusBadRequest, "user must exist and amount must be greater than 0")
		return
	}
	if svcErr != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, order)
}
