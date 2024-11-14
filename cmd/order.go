package main

import (
	"net/http"

	"github.com/ntentasd/ecommerce-microservices/models"
	"github.com/ntentasd/ecommerce-microservices/utils"
)

func (app *Application) createOrder(w http.ResponseWriter, r *http.Request) {
	var input models.OrderInput
	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	order, err := app.db.CreateOrder(input.CustomerID, input.Products)
	if err != nil {
		app.serverErrorResponse(w)
		return
	}

	event := models.OrderEvent{
		EventType: "OrderCreated",
		Order:     *order,
	}

	app.producer.SendOrderEvent(event, nil)

	utils.FormatResponse(w, http.StatusCreated, utils.Envelope{"order": order})
}
