package main

import (
	"net/http"

	"github.com/ntentasd/ecommerce-microservices/utils"
)

func (app *Application) registerRoutes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/health", app.health)

	router.Handle("/products", utils.MethodHandler(app.getProducts, http.MethodGet))
	router.Handle("/products/create", utils.MethodHandler(app.createProduct, http.MethodPost))
	router.Handle("/products/", utils.MethodHandler(app.getProduct, http.MethodGet))

	router.Handle("/orders/create", utils.MethodHandler(app.createOrder, http.MethodPost))

	return router
}
