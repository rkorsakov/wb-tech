package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func NewPostgresStorage(ctx context.Context, connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := createTables(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}
	return &Storage{db: db}, nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			order_uid          VARCHAR(255) PRIMARY KEY,
			track_number       VARCHAR(255),
			entry              VARCHAR(255),
			locale             VARCHAR(10),
			internal_signature VARCHAR(255),
			customer_id        VARCHAR(255),
			delivery_service   VARCHAR(255),
			shardkey           VARCHAR(255),
			sm_id              INTEGER,
			date_created       TIMESTAMP,
			oof_shard          VARCHAR(255))
	`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS deliveries (
			order_uid VARCHAR(255) REFERENCES orders(order_uid),
			name      VARCHAR(255),
			phone     VARCHAR(255),
			zip       VARCHAR(255),
			city      VARCHAR(255),
			address   VARCHAR(255),
			region    VARCHAR(255),
			email     VARCHAR(255),
			PRIMARY KEY (order_uid)
		)`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS payments (
			order_uid      VARCHAR(255) REFERENCES orders(order_uid),
			transaction    VARCHAR(255),
			request_id     VARCHAR(255),
			currency       VARCHAR(3),
			provider       VARCHAR(255),
			amount         INTEGER,
			payment_dt     BIGINT,
			bank           VARCHAR(255),
			delivery_cost  INTEGER,
			goods_total    INTEGER,
			custom_fee     INTEGER,
			PRIMARY KEY (order_uid)
		)`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS items (
			order_uid     VARCHAR(255) REFERENCES orders(order_uid),
			chrt_id       INTEGER,
			track_number  VARCHAR(255),
			price         INTEGER,
			rid           VARCHAR(255),
			name          VARCHAR(255),
			sale          INTEGER,
			size          VARCHAR(255),
			total_price   INTEGER,
			nm_id         INTEGER,
			brand         VARCHAR(255),
			status        INTEGER,
			PRIMARY KEY (order_uid, chrt_id)
		)`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) Close() error {
	return s.db.Close()
}
