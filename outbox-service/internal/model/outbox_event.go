package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	ProcessedAt *time.Time
	EventType   string
	Data        []byte
}

func (o OutboxEvent) String() string {
	return fmt.Sprintf("OutboxEvent[ID=%s, EventType=%s]", o.ID.String(), o.EventType)
}
