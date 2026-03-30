package main

import (
	"net/http"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`[
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
		{"id": 3, "name": "Charlie"}
	]`))
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`[
		{"user_id": 1, "id": 101, "status": "completed", "amount": 250.75},
		{"user_id": 1, "id": 102, "status": "pending", "amount": 100.00},
		{"user_id": 2, "id": 103, "status": "completed", "amount": 300.50},
		{"user_id": 3, "id": 104, "status": "completed", "amount": 150.25},
		{"user_id": 3, "id": 105, "status": "pending", "amount": 200.00}
	]`))
}

func main() {
	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/orders", orderHandler)
	http.ListenAndServe(":8081", nil)
}
