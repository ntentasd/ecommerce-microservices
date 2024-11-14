package main

import (
	"net/http"

	"github.com/ntentasd/ecommerce-microservices/utils"
)

func (app *Application) errorResponse(w http.ResponseWriter, status int, message any) {
	utils.FormatResponse(w, status, utils.Envelope{"error": message})
}

func (app *Application) serverErrorResponse(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (app *Application) notFoundResponse(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusNotFound, "The requested resource could not be found")
}

func (app *Application) badRequestResponse(w http.ResponseWriter, err error) {
	app.errorResponse(w, http.StatusBadRequest, err.Error())
}
