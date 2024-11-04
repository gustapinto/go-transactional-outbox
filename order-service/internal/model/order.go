package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Title     string
	Value     float64
}
