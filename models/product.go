package models

type Product struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
}

type ProductInput struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
}

type ProductEvent struct {
	EventType string
	Product   Product
}
