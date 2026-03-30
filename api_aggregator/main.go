package main

import (
	. "api_aggregator/internal/model"
	"api_aggregator/pkg"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to go aggregator server!")
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	// var userCh = make(chan []User)
	// var orderCh = make(chan []Order)
	users := []User{}
	orders := []Order{}
	userCh := make(chan []User)
	orderCh := make(chan []Order)
	errCh := make(chan error, 2)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go func() {
		users, err := getUsers(ctx)
		if err != nil {
			errCh <- err
			return
		}
		userCh <- users
	}()

	go func() {
		orders, err := getOrders(ctx)
		if err != nil {
			errCh <- err
			return
		}
		orderCh <- orders
	}()

	for i := 0; i < 2; i++ {
		select {
		case u := <-userCh:
			users = u
		case o := <-orderCh:
			orders = o
		case err := <-errCh:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	reports := utils.GenerateReport(users, orders)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func getUsers(ctx context.Context) ([]User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/users", nil)
	if err != nil {
		fmt.Println("Error fetching users:", err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching users:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var users []User
	if err = json.NewDecoder(resp.Body).Decode(&users); err != nil {
		fmt.Println("Error parsing users:", err)
		return nil, err
	}
	fmt.Printf("Fetched users: %+v\n", users)
	return users, nil
}

func getOrders(ctx context.Context) ([]Order, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081/orders", nil)
	if err != nil {
		fmt.Println("Error fetching orders", err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching users:", err)
		return nil, err
	}
	defer resp.Body.Close()
	var orders []Order
	if err = json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		fmt.Println("Error parsing orders:", err)
		return nil, err
	}
	fmt.Printf("Fetched orders: %+v\n", orders)
	return orders, nil
}

func main() {
	fmt.Println("Starting server on :8080")
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/report", reportHandler)
	http.ListenAndServe(":8080", nil)
}
