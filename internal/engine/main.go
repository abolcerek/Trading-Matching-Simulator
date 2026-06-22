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
	Next  *OrderNode
	Prev  *OrderNode
}

type Price struct {
	Price int64
	Head  *OrderNode
	Tail  *OrderNode
}

type PriceLevel struct {
	Price int64
	Quantity int64
}

type OrderBook struct {
	BidTree *treemap.TreeMap[int64, *Price]
	AskTree *treemap.TreeMap[int64, *Price]
	Hashmap map[uuid.UUID]*OrderNode
}

type Fill struct {
	Maker_order_id  uuid.UUID
	Taker_order_id  uuid.UUID
	Maker_user_id   uuid.UUID
	Taker_user_id   uuid.UUID
	Price           int64
	Quantity        int64
	Side            string
	Maker_remaining int64
	Taker_remaining int64
	Created_at      time.Time
}

type Rest struct {
	Order_id           uuid.UUID
	User_id            uuid.UUID
	Side               string
	Price              int64
	Status             string
	Remaining_quantity int64
}

type Cancel struct {
	Order_id         uuid.UUID
	User_id          uuid.UUID
	Side             string
	Price            int64
	Quantity_removed int64
}

type EventType string

const (
	Filled   EventType = "FILL"
	Canceled EventType = "CANCEL"
	Rested   EventType = "REST"
)

type Event struct {
	Type   EventType
	Fill   Fill
	Cancel Cancel
	Rest   Rest
}

type Snapshot struct {
	Bid []PriceLevel
	Ask []PriceLevel
}

type SnapshotRequest struct {
	Reply chan Snapshot
}

func NewOrderBook() *OrderBook {
	CreateAskTree := func() *treemap.TreeMap[int64, *Price] {
		AskTree := treemap.NewWithKeyCompare[int64, *Price](func(a, b int64) bool {
			return a < b
		})
		return AskTree
	}

	CreateBidTree := func() *treemap.TreeMap[int64, *Price] {
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
			Head:  order,
			Tail:  order,
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
	if price_node.Head == order && price_node.Tail == order {
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

func (orderbook *OrderBook) Match(order *OrderNode) ([]Event, error) {
	events := []Event{}
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
		return []Event{}, fmt.Errorf("Incorrect order side")
	}
	orderLoop:
	for order.Order.Remaining_quantity > 0 { // While the orders remaining quantity is > 0
		it := tree.Iterator()
		if !it.Valid() {
			fmt.Println("No opposing orders left")
			break orderLoop
		}
		best_price, price_node := it.Key(), it.Value() // Get the lowest ask or highest bid
		if order.Order.Type == "limit" {               // If its a limit buy or a limit ask
			if isBuying && best_price > order.Order.Price { // If the lowest ask is greater then the price of the buy order
				break orderLoop // break the loop
			} else if isBuying == false && best_price < order.Order.Price { // if the highest bid is less than the price of the ask order
				break orderLoop // break the loop
			}
		}
		fill_quantity := min(order.Order.Remaining_quantity, price_node.Head.Order.Remaining_quantity) // Compute the quantity that will be filled
		order.Order.Remaining_quantity -= fill_quantity                                                // Decrement by fill quantity
		price_node.Head.Order.Remaining_quantity -= fill_quantity                                      // Decrement by fill quantity
		fill := Fill{                                                                                  // Create a fill event
			Maker_order_id:  price_node.Head.Order.Id,
			Taker_order_id:  order.Order.Id,
			Maker_user_id:   price_node.Head.Order.UserID,
			Taker_user_id:   order.Order.UserID,
			Price:           price_node.Price,
			Quantity:        fill_quantity,
			Side:            order.Order.Side,
			Maker_remaining: price_node.Head.Order.Remaining_quantity,
			Taker_remaining: order.Order.Remaining_quantity,
			Created_at:      time.Now(),
		}
		events = append(events, Event{
			Type: Filled,
			Fill: fill,
		})
		if price_node.Head.Order.Remaining_quantity == 0 { // If the latest order at that price has been fufilled
			orderbook.RemoveNode(price_node.Head) // Remove it from the orderbook
		}
	}
	if order.Order.Type == "limit" && order.Order.Remaining_quantity > 0 { // If its a limit buy or ask and theres a remaining quantity
		orderbook.AddOrder(order) // Add the order to the orderbook
		var status string
		if order.Order.Quantity == order.Order.Remaining_quantity {
			status = "open"
		} else {
			status = "partially_filled"
		}
		events = append(events, Event{
			Type: Rested,
			Rest: Rest{
				Order_id:           order.Order.Id,
				User_id:            order.Order.UserID,
				Side:               order.Order.Side,
				Price:              order.Order.Price,
				Status:             status,
				Remaining_quantity: order.Order.Remaining_quantity,
			},
		})
	}
	return events, nil
}

func (orderbook *OrderBook) Cancel(order *OrderNode) ([]Event, error) {
	var events []Event
	orderNode, ok := orderbook.Hashmap[order.Order.Id]
	if !ok {
		return []Event{}, fmt.Errorf("Order not found in the orderbook")
	}
	orderbook.RemoveNode(orderNode)
	events = append(events, Event{
		Type: Canceled,
		Cancel: Cancel{
			Order_id:         order.Order.Id,
			User_id:          order.Order.UserID,
			Side:             order.Order.Side,
			Price:            order.Order.Price,
			Quantity_removed: orderNode.Order.Remaining_quantity,
		},
	})
	return events, nil
}

func (orderbook *OrderBook) Snapshot() Snapshot {
	snapshot := Snapshot{}
	BidPriceLevel := []PriceLevel{}
	AskPriceLevel := []PriceLevel{}
	for it := orderbook.BidTree.Iterator(); it.Valid(); it.Next() {
		priceNode := it.Value()
		price := priceNode.Price
		quantity := int64(0)
		for current := priceNode.Head; current != nil; current = current.Next {
			quantity = quantity + current.Order.Remaining_quantity
		}
		priceLevel := PriceLevel{
			Price: price,
			Quantity: quantity,
		}
		BidPriceLevel = append(BidPriceLevel, priceLevel)
	}
	for it := orderbook.AskTree.Iterator(); it.Valid(); it.Next() {
		priceNode := it.Value()
		price := priceNode.Price
		quantity := int64(0)
		for current := priceNode.Head; current != nil; current = current.Next {
			quantity = quantity + current.Order.Remaining_quantity
		}
		priceLevel := PriceLevel{
			Price: price,
			Quantity: quantity,
		}
		AskPriceLevel = append(AskPriceLevel, priceLevel)
	}
	snapshot.Bid = BidPriceLevel
	snapshot.Ask = AskPriceLevel
	return snapshot
}

func RunEngine(orderbook *OrderBook, in <-chan types.Envelope, eventChan chan<- []Event, snapshot <- chan SnapshotRequest) {
	for {
		select {
		case order := <- in:
			switch order.Tag {
			case "place":
				placed_order_node := OrderNode{
					Order: order.Order,
					Next:  nil,
					Prev:  nil,
				}
				events, err := orderbook.Match(&placed_order_node)
				if err != nil {
					fmt.Println("Error when matching the order")
				}
				eventChan <- events
				fmt.Printf("Here are the events : %v", events)
			case "cancel":
				canceled_order_node := OrderNode{
					Order: order.Order,
					Next:  nil,
					Prev:  nil,
				}
				canceled_event, err := orderbook.Cancel(&canceled_order_node)
				if err != nil {
					fmt.Println("Error when matching the order")
				}
				eventChan <- canceled_event
				fmt.Printf("Here is the canceled event : %v", canceled_event)
			default:
				fmt.Print("Incorrect order tag")
			}
		case request := <- snapshot:
			current_snapshot := orderbook.Snapshot()
			request.Reply <- current_snapshot
		}
	}
}