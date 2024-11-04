package service

import (
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(context.Context, string, float64) (uuid.UUID, error)
}

type Order struct {
	OrderRepository OrderRepository
}

func (o Order) Create(title string, value float64) (uuid.UUID, error) {
	orderId, err := o.OrderRepository.Create(context.Background(), title, value)
	if err != nil {
		return uuid.Nil, err
	}

	return orderId, nil
}
