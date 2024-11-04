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
	CREATE SCHEMA IF NOT EXISTS "order_service_schema";

	CREATE TABLE IF NOT EXISTS "order_service_schema"."orders" (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		title VARCHAR(100) NOT NULL,
		product VARCHAR(100) NOT NULL,
		quantity BIGINT NOT NULL,
		value DOUBLE PRECISION NOT NULL
	);

	CREATE TABLE IF NOT EXISTS "order_service_schema"."outbox" (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		processed_at TIMESTAMP,
		event_type VARCHAR(255) NOT NULL,
		data JSONB NOT NULL
	);
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
