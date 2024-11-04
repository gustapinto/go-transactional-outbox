package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/gustapinto/go-transactional-outbox/message-relay-service/internal/model"
)

type OutboxRepository interface {
	GetNonProcessedOutboxEvents(context.Context) ([]model.OutboxEvent, error)

	SetOutboxEventAsProcessed(context.Context, uuid.UUID) error
}

type OutboxEventProcessor interface {
	Process(context.Context, model.OutboxEvent) error
}

type Outbox struct {
	OutboxRepository OutboxRepository
}

func (o Outbox) GetAndProcessNonProcessedOutboxEvents(processorMapping map[string]OutboxEventProcessor) error {
	events, err := o.OutboxRepository.GetNonProcessedOutboxEvents(context.Background())
	if err != nil {
		return err
	}

	eventsSize := len(events)
	if eventsSize == 0 {
		log.Printf("Found no event to process, skipping...")
		return nil
	}

	log.Printf("Found %d events to process, starting...", eventsSize)

	for _, event := range events {
		log.Printf("Processing %s", event.String())

		processor, ok := processorMapping[event.EventType]
		if !ok {
			log.Printf("There is no event processor mapped to event %s, skipping %s", event.String(), event.ID.String())
			continue
		}

		if err := processor.Process(context.Background(), event); err != nil {
			log.Printf("Failed to process event %s, got error %s", event.String(), err.Error())
			continue
		}

		if err := o.OutboxRepository.SetOutboxEventAsProcessed(context.Background(), event.ID); err != nil {
			log.Printf("Failed to set event %s as processed, got error %s", event.String(), err.Error())
			continue
		}
	}

	return nil
}
