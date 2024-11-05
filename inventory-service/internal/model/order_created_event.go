package model

import (
	"fmt"

	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	OrderID  uuid.UUID `json:"order_id,omitempty"`
	Title    string    `json:"title,omitempty"`
	Product  string    `json:"product_code,omitempty"`
	Quantity int64     `json:"quantity,omitempty"`
	Value    float64   `json:"value,omitempty"`
}

func (o OrderCreatedEvent) String() string {
	return fmt.Sprintf("OrderCreatedEvent[OrderID=%s]", o.OrderID.String())

}
