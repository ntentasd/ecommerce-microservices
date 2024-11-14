package database

import (
	"context"
	"time"

	"github.com/ntentasd/ecommerce-microservices/models"
)

func (db *Database) CreateOrder(customerID int, productIDs map[int]int) (*models.Order, error) {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return &models.Order{}, err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var totalAmount float64
	var products []models.Product
	for productID, quantity := range productIDs {
		product, err := db.GetProduct(productID)
		if err != nil {
			return &models.Order{}, err
		}

		products = append(products, product)

		totalAmount += product.Price * float64(quantity)
	}

	var order models.Order
	order.CustomerID = customerID
	err = tx.QueryRowContext(ctx, "INSERT INTO orders (customer_id, total_amount) VALUES ($1, $2) RETURNING id, status, total_amount", customerID, totalAmount).Scan(&order.ID, &order.Status, &order.TotalAmount)
	if err != nil {
		return &models.Order{}, err
	}

	for _, product := range products {
		quantity := productIDs[product.ID]
		_, err = tx.ExecContext(ctx, "INSERT INTO order_items (order_id, product_id, category, quantity, price) VALUES ($1, $2, $3, $4, $5)", order.ID, product.ID, product.Category, quantity, product.Price)
		if err != nil {
			return &models.Order{}, err
		}
	}

	return &order, tx.Commit()
}
