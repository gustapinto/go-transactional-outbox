package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gustapinto/go-transactional-outbox/outbox-service/internal/model"
)

type Outbox struct {
	DB *sql.DB
}

func (o Outbox) GetNonProcessedOutboxEvents(ctx context.Context) ([]model.OutboxEvent, error) {
	query := `
	SELECT
		id,
		created_at,
		event_type,
		data
	FROM
		outbox
	WHERE
		processed_at IS NULL
	ORDER BY
		created_at ASC
	`

	rows, err := o.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]model.OutboxEvent, 0)
	for rows.Next() {
		var event model.OutboxEvent
		if err := rows.Scan(&event.ID, &event.CreatedAt, &event.EventType, &event.Data); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (o Outbox) SetOutboxEventAsProcessed(ctx context.Context, id uuid.UUID) error {
	query := `
	UPDATE
		outbox
	SET
		processed_at = $1
	WHERE
		id = $2
	`

	if _, err := o.DB.ExecContext(ctx, query, time.Now(), id); err != nil {
		return err
	}

	return nil
}
