package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/gustapinto/go-transactional-outbox/inventory-service/internal/model"
)

type InventoryRepository interface {
	Update(context.Context, string, uuid.UUID, int64) error

	OrderHasAlreadyBeenProcessed(context.Context, uuid.UUID) (bool, error)
}

type Inventory struct {
	InventoryRepository InventoryRepository
}

func (i Inventory) ProcessOrder(event model.OrderCreatedEvent) error {
	log.Printf("Processing event %s", event.String())

	exists, err := i.InventoryRepository.OrderHasAlreadyBeenProcessed(context.Background(), event.OrderID)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Order for event %s has already been processed, skipping", event.String())
		return nil
	}

	if err := i.InventoryRepository.Update(context.Background(), event.Product, event.OrderID, event.Quantity); err != nil {
		return err
	}

	log.Printf("Updated inventory for event %s", event.String())
	return nil
}
