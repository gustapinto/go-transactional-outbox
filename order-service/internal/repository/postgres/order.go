package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gustapinto/go-transactional-outbox/order-service/internal/model"
)

type Order struct {
	DB *sql.DB
}

func (Order) createOrderOutbox(tx *sql.Tx, ctx context.Context, eventType string, data []byte) error {
	query := `
	INSERT INTO "outbox" (
		id,
		created_at,
		event_type,
		data
	)
	VALUES (
		$1::UUID,
		$2::TIMESTAMP,
		$3::VARCHAR,
		$4::JSONB
	);
	`

	id := uuid.New()
	if _, err := tx.ExecContext(ctx, query, id, time.Now(), eventType, data); err != nil {
		return err
	}

	return nil
}

func (o Order) createOrder(tx *sql.Tx, ctx context.Context, title string, value float64) (uuid.UUID, error) {
	query := `
	INSERT INTO "orders" (
		id,
		created_at,
		title,
		value
	)
	VALUES (
		$1::UUID,
		$2::TIMESTAMP,
		$3::VARCHAR,
		$4::DOUBLE PRECISION
	);
	`

	id := uuid.New()
	if _, err := tx.ExecContext(ctx, query, id, time.Now(), title, value); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (o Order) Create(ctx context.Context, title string, value float64) (uuid.UUID, error) {
	tx, err := o.DB.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}

	orderId, err := o.createOrder(tx, ctx, title, value)
	if err != nil {
		return uuid.Nil, err
	}

	event, err := json.Marshal(&model.OrderCreatedEvent{
		OrderID: orderId,
		Title:   title,
		Value:   value,
	})
	if err != nil {
		_ = tx.Rollback()
		return uuid.Nil, err
	}

	if err := o.createOrderOutbox(tx, ctx, "ORDER_CREATED", event); err != nil {
		_ = tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return orderId, nil
}
