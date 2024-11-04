package model

import "github.com/google/uuid"

type OrderCreatedEvent struct {
	OrderID uuid.UUID `json:"order_id,omitempty"`
	Title   string    `json:"title,omitempty"`
	Value   float64   `json:"value,omitempty"`
}
