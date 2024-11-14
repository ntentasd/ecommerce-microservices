package database

import (
	"context"
	"time"

	"github.com/ntentasd/ecommerce-microservices/models"
)

func (db *Database) GetProducts() ([]models.Product, error) {
	rows, err := db.db.Query("SELECT id, name, description, category, price, stock_quantity FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.StockQuantity); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (db *Database) GetProduct(id int) (models.Product, error) {
	var product models.Product

	stmt, err := db.db.Prepare("SELECT id, name, description, category, price, stock_quantity FROM products WHERE id = $1")
	if err != nil {
		return models.Product{}, err
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := stmt.QueryRowContext(ctx, id).Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.StockQuantity); err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (db *Database) CreateProduct(product models.ProductInput) (int, error) {
	stmt, err := db.db.Prepare("INSERT INTO products (name, description, category, price, stock_quantity) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	if err := stmt.QueryRowContext(ctx, product.Name, product.Description, product.Category, product.Price, product.StockQuantity).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
