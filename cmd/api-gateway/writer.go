package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/database"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine"
	"github.com/google/uuid"
)

func (cfg *apiConfig) Consumer(ch <-chan []engine.Event){
	for output := range ch {
		err := cfg.handleEvents(output)
		if err != nil {
			fmt.Printf("Transaction failed with error: %v", err)
		}
	}
}

func (cfg *apiConfig) handleEvents(event []engine.Event) error {
	tx, err := cfg.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := cfg.database.WithTx(tx)
	for _, event := range event {
		switch event.Type {
		case engine.Filled:
			trade_params :=	database.CreateTradeParams{
				TradeID: uuid.New(),
				MakerOrderID: event.Fill.Maker_order_id,
				TakerOrderID: event.Fill.Taker_order_id,
				MakerUserID: event.Fill.Maker_user_id,
				TakerUserID: event.Fill.Taker_user_id,
				Price: event.Fill.Price,
				Quantity: event.Fill.Quantity,
				CreatedAt: time.Now(),
			}
			_, err := qtx.CreateTrade(context.Background(), trade_params)
			if err != nil {
				return err
			}
			var maker_order_status string
			if event.Fill.Maker_remaining > 0 {
				maker_order_status = "partially_filled"
			} else {
				maker_order_status = "filled"
			}
			update_maker_order_params := database.UpdateOrderParams{
				RemainingQuantity: event.Fill.Maker_remaining,
				Status: maker_order_status,
				OrderID: event.Fill.Maker_order_id,
			}
			err = qtx.UpdateOrder(context.Background(), update_maker_order_params)
			if err != nil {
				return err
			}
			var taker_order_status string
			if event.Fill.Taker_remaining > 0 {
				taker_order_status = "partially_filled"
			} else {
				taker_order_status = "filled"
			}
			update_taker_order_params := database.UpdateOrderParams{
				RemainingQuantity: event.Fill.Taker_remaining,
				Status: taker_order_status,
				OrderID: event.Fill.Taker_order_id,
			}
			err = qtx.UpdateOrder(context.Background(), update_taker_order_params)
			if err != nil {
				return err
			}
			var maker_balance int64
			var taker_balance int64
			if event.Fill.Side == "buy" {
				maker_balance = event.Fill.Price * event.Fill.Quantity
				taker_balance = -event.Fill.Price * event.Fill.Quantity
			} else {
				maker_balance = -event.Fill.Price * event.Fill.Quantity
				taker_balance = event.Fill.Price * event.Fill.Quantity				
			}
			update_maker_params := database.UpdateBalanceParams{
				Balance: maker_balance,
				ID: event.Fill.Maker_user_id,
			}
			err = qtx.UpdateBalance(context.Background(), update_maker_params)
			if err != nil {
				return err
			}
			update_taker_params := database.UpdateBalanceParams{
				Balance: taker_balance,
				ID: event.Fill.Taker_user_id,
			}
			err = qtx.UpdateBalance(context.Background(), update_taker_params)
			if err != nil {
				return err
			}
		case engine.Rested:
			order_params := database.UpdateOrderParams{
				RemainingQuantity: event.Rest.Remaining_quantity,
				Status: event.Rest.Status,
				OrderID: event.Rest.Order_id,
			}
			err := qtx.UpdateOrder(context.Background(), order_params)
			if err != nil {
				return err
			}
		case engine.Canceled:
			order_params := database.UpdateOrderParams{
				RemainingQuantity: event.Cancel.Quantity_removed,
				Status: "canceled",
				OrderID: event.Cancel.Order_id,
			}
			err := qtx.UpdateOrder(context.Background(), order_params)
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}