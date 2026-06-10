package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/auth"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/database"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine/types"
	"github.com/google/uuid"
)

type OrderRequest struct {
	Side string `json:"side"`
	Type string `json:"type"`
	Price int64 `json:"price"`
	Quantity int64 `json:"quantity"`
}

func (cfg apiConfig) HandlerCreateOrder(w http.ResponseWriter, r *http.Request) {
	err_params := error_parameters{}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err_params.Error = "Invalid access token"
		handleErrors(w, &err_params, 401)
		return
	}
	id, err := auth.ValidateJWT(token, cfg.JWT_secret)
	if err != nil {
		err_params.Error = "Invalid access token"
		handleErrors(w, &err_params, 401)
		return
	}
	orderRequest := OrderRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&orderRequest)
	if err != nil {
		err_params.Error = "Error decoding JSON"
		handleErrors(w, &err_params, 400)
		return
	}
	if orderRequest.Quantity <= 0 {
		err_params.Error = "No quantity for order"
		handleErrors(w, &err_params, 400)
		return
	}
	if orderRequest.Type != "market" && orderRequest.Type != "limit" {
		err_params.Error = "Incorrect order type"
		handleErrors(w, &err_params, 400)
		return
	}
	if orderRequest.Side != "buy" && orderRequest.Side != "sell" {
		err_params.Error = "Incorrect order side"
		handleErrors(w, &err_params, 400)
		return
	}
	var price sql.NullInt64
	switch orderRequest.Type {
	case "market":
		price.Valid = false
	case "limit":
		if orderRequest.Price <= 0 {
			err_params.Error = "Price must be greater than 0"
			handleErrors(w, &err_params, 400)
			return
		} else {
			price.Int64 = orderRequest.Price
			price.Valid = true
		}
	}
	create_order_params := database.CreateOrderParams{
		OrderID: uuid.New(),
		UserID: id,
		Side: orderRequest.Side,
		Type: orderRequest.Type,
		Price: price,
		Quantity: orderRequest.Quantity,
		RemainingQuantity: orderRequest.Quantity,
		Status: "pending",
		CreatedAt: time.Now(),
	}
	database_order, err := cfg.database.CreateOrder(r.Context(), create_order_params)
	if err != nil {
		err_params.Error = "Error creating order"
		handleErrors(w, &err_params, 400)
		return
	}
	order := types.Order{
		Id: database_order.OrderID,
		UserID: database_order.UserID,
		Sequence_num: database_order.SequenceNum.Int64,
		Side: database_order.Side,
		Type: database_order.Type,
		Price: database_order.Price.Int64,
		Quantity: database_order.Quantity,
		Remaining_quantity: database_order.RemainingQuantity,
		Created_at: database_order.CreatedAt,
	}
	envelope := types.Envelope{
		Tag: place,
		Order: order,
		Event_sequence_num: 0,
		// Event sequence number will be implemented later
	}
	cfg.orderChannel <- envelope
	data, err := json.Marshal(&order)
	if err != nil {
		err_params.Error = "Error marshalling JSON"
		handleErrors(w, &err_params, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(202)
	w.Write(data)
}