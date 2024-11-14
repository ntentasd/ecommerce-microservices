package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ntentasd/ecommerce-microservices/models"
	"github.com/ntentasd/ecommerce-microservices/utils"
	"github.com/redis/go-redis/v9"
)

func (app *Application) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.db.GetProducts()
	if err != nil {
		app.serverErrorResponse(w)
		return
	}

	utils.FormatResponse(w, http.StatusOK, products)
}

func (app *Application) getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseInt(r.URL.Path)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	cached, err := app.redis.Get(r.Context(), strconv.Itoa(id)).Result()
	if err == redis.Nil {
		product, err := app.db.GetProduct(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				app.notFoundResponse(w)
				return
			}

			app.serverErrorResponse(w)
			return
		}

		productJSON, err := json.Marshal(product)
		if err != nil {
			app.serverErrorResponse(w)
			return
		}

		app.redis.Set(r.Context(), strconv.Itoa(id), productJSON, 15*time.Second)
		utils.FormatResponse(w, http.StatusOK, product)
		return
	}

	var product models.Product
	if err := json.Unmarshal([]byte(cached), &product); err != nil {
		app.serverErrorResponse(w)
		return
	}

	utils.FormatResponse(w, http.StatusOK, product)
}

func (app *Application) createProduct(w http.ResponseWriter, r *http.Request) {
	var product models.ProductInput
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		if errors.Is(err, io.EOF) {
			app.badRequestResponse(w, errors.New("request body is empty"))
			return
		}
		app.badRequestResponse(w, err)
		return
	}

	id, err := app.db.CreateProduct(product)
	if err != nil {
		slog.Error("failed to create product", "error", err)
		app.serverErrorResponse(w)
		return
	}

	outProduct := models.Product{
		ID:            id,
		Name:          product.Name,
		Description:   product.Description,
		Category:      product.Category,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
	}

	event := models.ProductEvent{
		EventType: "ProductCreated",
		Product:   outProduct,
	}

	app.producer.SendProductEvent(event, nil)
	utils.FormatResponse(w, http.StatusCreated, utils.Envelope{"product": outProduct})
}
