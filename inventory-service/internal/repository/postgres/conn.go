package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func OpenDatabaseConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func InitializeDatabase(db *sql.DB) error {
	query := `
	CREATE SCHEMA IF NOT EXISTS "inventory_service_schema";

	CREATE TABLE IF NOT EXISTS "inventory_service_schema"."inventory" (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		product_code VARCHAR(100) UNIQUE NOT NULL,
		quantity_in_stock BIGINT NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS "inventory_service_schema"."inventory_transactions" (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		product_code VARCHAR(100) NOT NULL REFERENCES "inventory_service_schema"."inventory" (product_code),
		order_id UUID UNIQUE NOT NULL,
		quantity BIGINT
	);

	INSERT INTO "inventory_service_schema"."inventory" (
		id,
		created_at,
		product_code,
		quantity_in_stock
	) VALUES (
		GEN_RANDOM_UUID(),
		CURRENT_TIMESTAMP,
		'PRODUCT_1',
		50
	), (
		GEN_RANDOM_UUID(),
		CURRENT_TIMESTAMP,
		'PRODUCT_2',
		50
	), (
		GEN_RANDOM_UUID(),
		CURRENT_TIMESTAMP,
		'PRODUCT_3',
		50
	)
	ON CONFLICT
		(product_code)
	DO NOTHING;
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
