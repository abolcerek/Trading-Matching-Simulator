package main

import (
	"context"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine/types"
)

func (cfg *apiConfig) Replay() (error) {
	orders, err := cfg.database.Replay(context.Background())
	if err != nil {
		return err
	}
	for _, order := range orders {
		switch order.Status {
		case "open", "partially_filled":
			order_node := engine.OrderNode{
			Order: types.Order{
				Id: order.OrderID,
				UserID: order.UserID,
				Sequence_num: order.SequenceNum.Int64,
				Side: order.Side,
				Type: order.Type,
				Price: order.Price.Int64,
				Quantity: order.Quantity,
				Remaining_quantity: order.RemainingQuantity,
				Created_at: order.CreatedAt,
			},
			Next: nil,
			Prev: nil,
		}
		cfg.orderbook.AddOrder(&order_node)
		}
	}
	return nil
}