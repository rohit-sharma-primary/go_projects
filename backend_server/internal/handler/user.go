package handler

import (
	"errors"
	"net/http"
	"strconv"

	"server/internal/model"
	"server/internal/service"
	"server/internal/utils"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", h.getUsersHandler)
	mux.HandleFunc("GET /users/{id}", h.getUserByIDHandler)
	mux.HandleFunc("POST /users", h.postUserHandler)
}

func (h *UserHandler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, h.service.ListUsers())
}

func (h *UserHandler) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, svcErr := h.service.GetUser(id)
	if errors.Is(svcErr, service.ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	if svcErr != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) postUserHandler(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, svcErr := h.service.CreateUser(req)
	if errors.Is(svcErr, service.ErrInvalidInput) {
		utils.WriteError(w, http.StatusBadRequest, "name is required and email must contain '@'")
		return
	}
	if svcErr != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}
