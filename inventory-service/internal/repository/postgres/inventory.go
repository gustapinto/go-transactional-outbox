package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	DB *sql.DB
}

func (Inventory) updateInventory(tx *sql.Tx, ctx context.Context, productCode string, quantity int64) error {
	query := `
	UPDATE
		"inventory_service_schema"."inventory"
	SET
		quantity_in_stock = (quantity_in_stock - $1::BIGINT)::BIGINT,
		updated_at = $2::TIMESTAMP
	WHERE
		product_code = $3::VARCHAR;
	`

	if _, err := tx.ExecContext(ctx, query, quantity, time.Now(), productCode); err != nil {
		return err
	}

	return nil
}

func (Inventory) insertInventoryTransaction(tx *sql.Tx, ctx context.Context, productCode string, orderId uuid.UUID, quantity int64) error {
	query := `
	INSERT INTO "inventory_service_schema"."inventory_transactions" (
		id,
		created_at,
		product_code,
		order_id,
		quantity
	) VALUES (
	 	$1::UUID,
		$2::TIMESTAMP,
		$3::VARCHAR,
		$4::UUID,
		$5::BIGINT
	);
	`

	id := uuid.New()
	if _, err := tx.ExecContext(ctx, query, id, time.Now(), productCode, orderId, quantity); err != nil {
		return err
	}

	return nil
}

func (i Inventory) Update(ctx context.Context, productCode string, orderId uuid.UUID, quantity int64) error {
	tx, err := i.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := i.insertInventoryTransaction(tx, ctx, productCode, orderId, quantity); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Join(err, errRollback)
		}

		return err
	}

	if err := i.updateInventory(tx, ctx, productCode, quantity); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Join(err, errRollback)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (i Inventory) OrderHasAlreadyBeenProcessed(ctx context.Context, orderId uuid.UUID) (bool, error) {
	query := `
	SELECT
		count(1)
	FROM
		"inventory_service_schema"."inventory_transactions"
	WHERE
		order_id = $1::UUID;
	`

	row := i.DB.QueryRowContext(ctx, query, orderId)
	if row.Err() != nil {
		if strings.Contains(row.Err().Error(), "no row") {
			return false, nil
		}

		return false, row.Err()
	}

	var count int64
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count != 0, nil
}
