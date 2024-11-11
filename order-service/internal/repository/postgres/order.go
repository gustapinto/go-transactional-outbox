package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gustapinto/go-transactional-outbox/order-service/internal/model"
)

type Order struct {
	DB *sql.DB
}

func (Order) createOrderOutboxEvent(tx *sql.Tx, ctx context.Context, eventType string, data []byte) error {
	query := `
	INSERT INTO "order_service_schema"."outbox" (
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

func (o Order) createOrder(tx *sql.Tx, ctx context.Context, title, product string, quantity int64, value float64) (uuid.UUID, error) {
	query := `
	INSERT INTO "order_service_schema"."orders" (
		id,
		created_at,
		title,
		product,
		quantity,
		value
	)
	VALUES (
		$1::UUID,
		$2::TIMESTAMP,
		$3::VARCHAR,
		$4::VARCHAR,
		$5::BIGINT,
		$6::DOUBLE PRECISION
	);
	`

	id := uuid.New()
	if _, err := tx.ExecContext(ctx, query, id, time.Now(), title, product, quantity, value); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (o Order) Create(ctx context.Context, title, product string, quantity int64, value float64) (uuid.UUID, error) {
	tx, err := o.DB.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}

	orderId, err := o.createOrder(tx, ctx, title, product, quantity, value)
	if err != nil {
		return uuid.Nil, err
	}

	event, err := json.Marshal(&model.OrderCreatedEvent{
		OrderID:  orderId,
		Title:    title,
		Product:  product,
		Quantity: quantity,
		Value:    value,
	})
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return uuid.Nil, errors.Join(err, errRollback)
		}

		return uuid.Nil, err
	}

	if err := o.createOrderOutboxEvent(tx, ctx, "ORDER_CREATED", event); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return uuid.Nil, errors.Join(err, errRollback)
		}

		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return orderId, nil
}
