package models

type Order struct {
	ID          int     `json:"id"`
	CustomerID  int     `json:"customer_id"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
}

type OrderInput struct {
	CustomerID int         `json:"customer_id"`
	Products   map[int]int `json:"products"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Category  string  `json:"category"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderEvent struct {
	EventType string
	Order     Order
}
