package model

type Order struct {
	UserID int     `json:"user_id"`
	ID     int     `json:"id"`
	Status string  `json:"status"`
	Amount float64 `json:"amount"`
}

type Report struct {
	UserID          int     `json:"user_id"`
	Name            string  `json:"name"`
	TotalSpend      float64 `json:"total_spend"`
	CompletedOrders int     `json:"completed_orders"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
