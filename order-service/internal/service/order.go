package service

import (
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(context.Context, string, string, int64, float64) (uuid.UUID, error)
}

type Order struct {
	OrderRepository OrderRepository
}

func (o Order) Create(title, product string, quantity int64, value float64) (uuid.UUID, error) {
	orderId, err := o.OrderRepository.Create(context.Background(), title, product, quantity, value)
	if err != nil {
		return uuid.Nil, err
	}

	return orderId, nil
}
