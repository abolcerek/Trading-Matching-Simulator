package engine

import (
	"fmt"
	"time"

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

type Fill struct {
	Maker_order_id uuid.UUID
	Taker_order_id uuid.UUID
	Maker_user_id uuid.UUID
	Taker_user_id uuid.UUID
	Price int64
	Quantity int64
	Maker_remaining int64
	Taker_remaining int64
	Created_at time.Time
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


func (orderbook *OrderBook) Match(order *OrderNode) []Fill{
	fills := []Fill{}
	var tree *treemap.TreeMap[int64, *Price]
	var isBuying bool
	switch order.Order.Side {
	case "bid", "buy":
		tree = orderbook.AskTree
		isBuying = true
	case "ask", "sell":
		tree = orderbook.BidTree
		isBuying = false
	default:
		fmt.Println("Incorrect order side")
		return fills
	}
	orderLoop:
	for order.Order.Remaining_quantity > 0 { // While the orders remaining quantity is > 0
		it := tree.Iterator()
		if !it.Valid() {
			fmt.Println("No opposing orders left")
			break orderLoop
		}
		best_price, price_node := it.Key(), it.Value() // Get the lowest ask or highest bid 
		if order.Order.Type == "limit" { // If its a limit buy or a limit ask
			if isBuying && best_price > order.Order.Price { // If the lowest ask is greater then the price of the buy order
				break orderLoop // break the loop
			} else if isBuying == false && best_price < order.Order.Price{ // if the highest bid is less than the price of the ask order
				break orderLoop // break the loop
			} 
		}
		fill_quantity := min(order.Order.Remaining_quantity, price_node.Head.Order.Remaining_quantity) // Compute the quantity that will be filled
		order.Order.Remaining_quantity -= fill_quantity // Decrement by fill quantity
		price_node.Head.Order.Remaining_quantity -= fill_quantity // Decrement by fill quantity
		fill := Fill{ // Create a fill event
			Maker_order_id: price_node.Head.Order.Id,
			Taker_order_id: order.Order.Id,
			Maker_user_id: price_node.Head.Order.UserID,
			Taker_user_id: order.Order.UserID,
			Price: price_node.Price,
			Quantity: fill_quantity,
			Maker_remaining: price_node.Head.Order.Remaining_quantity,
			Taker_remaining: order.Order.Remaining_quantity,
			Created_at: time.Now(),
		}
		fills = append(fills, fill) // Append the fill to the slice of fills
		if price_node.Head.Order.Remaining_quantity == 0 { // If the latest order at that price has been fufilled
			orderbook.RemoveNode(price_node.Head) // Remove it from the orderbook
		}
	}
	if order.Order.Type == "limit" && order.Order.Remaining_quantity > 0 { // If its a limit buy or ask and theres a remaining quantity
		orderbook.AddOrder(order) // Add the order to the orderbook
	}
	return fills
}


func (orderbook *OrderBook) Cancel(order *OrderNode) {
	orderNode, ok := orderbook.Hashmap[order.Order.Id]
	if !ok {
		fmt.Println("Order not found in the orderbook")
		return
	}
	orderbook.RemoveNode(orderNode)
}