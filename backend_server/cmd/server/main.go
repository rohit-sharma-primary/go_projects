package main

import (
	"log"
	"net/http"

	"server/internal/handler"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"
)

func main() {
	mux := http.NewServeMux()

	userRepo := repository.NewUserRepository()
	orderRepo := repository.NewOrderRepository()

	userService := service.NewUserService(userRepo)
	orderService := service.NewOrderService(orderRepo, userRepo)

	handler.NewUserHandler(userService).Register(mux)
	handler.NewOrderHandler(orderService).Register(mux)

	server := middleware.RecoverPanics(
		middleware.HandleAuth(
			middleware.LogRequests(mux),
		),
	)

	log.Println("server listening on :3000")
	if err := http.ListenAndServe(":3000", server); err != nil {
		log.Fatal(err)
	}
}
