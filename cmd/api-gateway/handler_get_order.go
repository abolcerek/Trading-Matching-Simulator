package main

import (
	"encoding/json"
	"net/http"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/auth"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine/types"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandlerGetOrder(w http.ResponseWriter, r *http.Request) {
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
	order_id, err := uuid.Parse(r.PathValue("orderID"))
	if err != nil {
		err_params.Error = "Something went wrong"
		handleErrors(w, &err_params, 400)
		return
	}
	database_order, err := cfg.database.GetOrder(r.Context(), order_id)
	if err != nil {
		err_params.Error = "Order not found"
		handleErrors(w, &err_params, 404)
		return
	}
	if database_order.UserID != id {
		err_params.Error = "Invalid access"
		handleErrors(w, &err_params, 403)
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
	data, err := json.Marshal(&order)
	if err != nil {
		err_params.Error = "Error marshalling JSON"
		handleErrors(w, &err_params, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(data)
}