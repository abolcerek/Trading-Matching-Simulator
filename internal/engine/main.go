package engine

import (
	"github.com/google/uuid"
	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine/types"
	"github.com/igrmk/treemap/v2"
)

type OrderNode struct {
	Order types.Order
	Next *OrderNode
	Prev *OrderNode
}

type Price struct {
	Price int64
	Head *OrderNode
	Tail *OrderNode
}

type OrderBook struct {
	BidTree *treemap.TreeMap[int64, *Price]
	AskTree *treemap.TreeMap[int64, *Price]
	Hashmap map[uuid.UUID]*OrderNode
}

