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
	Type    string  `json:"type"`
	Product Product `json:"product"`
}

func (e ProductEvent) GetType() string {
	return e.Type
}

func (e ProductEvent) GetObject() interface{} {
	return e.Product
}
