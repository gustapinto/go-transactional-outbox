package model

import (
	"time"

	"github.com/google/uuid"
)

type Outbox struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	PublishedAt *time.Time
	EventType   string
	Data        []byte
}
