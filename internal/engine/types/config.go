package types

import (
	"time"
	"github.com/google/uuid"
)


type Order struct {
	Id uuid.UUID
	UserID uuid.UUID
	Sequence_num int64
	Side string
	Type string
	Price int64
	Quantity int64
	Remaining_quantity int64
	Created_at time.Time
}

type Envelope struct {
	Tag string
	Order Order
	Event_sequence_num int64
}
