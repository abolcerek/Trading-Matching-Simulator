package engine

import (
	"fmt"

	"github.com/abolcerek/Trading-Matching-Simulator/internal/engine/types"
	"github.com/google/uuid"
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

func NewOrderBook() *OrderBook{
	CreateAskTree := func() *treemap.TreeMap[int64, *Price]{
		AskTree := treemap.NewWithKeyCompare[int64, *Price](func(a, b int64) bool {
			return a < b
		})
		return AskTree
	}

	CreateBidTree := func() *treemap.TreeMap[int64, *Price]{
		BidTree := treemap.NewWithKeyCompare[int64, *Price](func(a, b int64) bool {
			return a > b
		})
		return BidTree
	}
	hashmap := make(map[uuid.UUID]*OrderNode)
	orderbook := OrderBook{
		BidTree: CreateBidTree(),
		AskTree: CreateAskTree(),
		Hashmap: hashmap,
	}
	return &orderbook
}


func (orderbook *OrderBook) AddOrder(order *OrderNode) {
	var tree *treemap.TreeMap[int64, *Price]
	switch order.Order.Side {
	case "bid", "buy":
		tree = orderbook.BidTree
	case "ask", "sell":
		tree = orderbook.AskTree
	default:
		fmt.Println("Incorrect order side")
		return
	}
	price_node, ok := tree.Get(order.Order.Price)
	if !ok {
		new_price_node := Price{
			Price: order.Order.Price,
			Head: order,
			Tail: order,
		}
		tree.Set(order.Order.Price, &new_price_node)
		orderbook.Hashmap[order.Order.Id] = order
		return
	}
	price_node.Tail.Next = order
	order.Prev = price_node.Tail
	price_node.Tail = order
	orderbook.Hashmap[order.Order.Id] = order
}


func (orderbook *OrderBook) RemoveNode(order *OrderNode) {
	var tree *treemap.TreeMap[int64, *Price]
	switch order.Order.Side {
	case "bid", "buy":
		tree = orderbook.BidTree
	case "ask", "sell":
		tree = orderbook.AskTree
	default:
		fmt.Println("Incorrect order side")
		return
	}
	price_node, ok := tree.Get(order.Order.Price)
	if !ok {
		fmt.Println("Node not found in the tree")
		return
	}
	// When there is only one order in the doubly linked list
	if price_node.Head == order && price_node.Tail == order{
		tree.Del(order.Order.Price)
		delete(orderbook.Hashmap, order.Order.Id)
		return
	}
	// When the order node is the head of the doubly linked list
	if price_node.Head == order && price_node.Tail != order {
		price_node.Head = order.Next
		order.Next.Prev = nil
		order.Next = nil
		delete(orderbook.Hashmap, order.Order.Id)
		return		
	}
	// When the order node is the tail of the doubly linked list
	if price_node.Tail == order && price_node.Head != order {
		price_node.Tail = order.Prev
		order.Prev.Next = nil
		order.Prev = nil
		delete(orderbook.Hashmap, order.Order.Id)
		return
	}
	// When the order node is somewhere in the middle of the doubly linked list
	if price_node.Head != order && price_node.Tail != order {
		order.Prev.Next = order.Next
		order.Next.Prev = order.Prev
		delete(orderbook.Hashmap, order.Order.Id)
		return
	}
}
