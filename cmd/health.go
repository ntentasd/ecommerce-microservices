package main

import (
	"net/http"

	"github.com/ntentasd/ecommerce-microservices/utils"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func (app *Application) health(w http.ResponseWriter, r *http.Request) {
	utils.FormatResponse(w, http.StatusOK, HealthResponse{
		Status:  "OK",
		Message: "Service is healthy",
		Version: "1.0.0",
	})
}
