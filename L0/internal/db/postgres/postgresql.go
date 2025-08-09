package postgres

import (
	"L0/internal/models"
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

func (s *Storage) SaveOrder(ctx context.Context, order models.Order) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        )
        ON CONFLICT (order_uid) DO UPDATE SET
            track_number = EXCLUDED.track_number,
            entry = EXCLUDED.entry,
            locale = EXCLUDED.locale,
            internal_signature = EXCLUDED.internal_signature,
            customer_id = EXCLUDED.customer_id,
            delivery_service = EXCLUDED.delivery_service,
            shardkey = EXCLUDED.shardkey,
            sm_id = EXCLUDED.sm_id,
            date_created = EXCLUDED.date_created,
            oof_shard = EXCLUDED.oof_shard`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
        INSERT INTO deliveries (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
        ON CONFLICT (order_uid) DO UPDATE SET
            name = EXCLUDED.name,
            phone = EXCLUDED.phone,
            zip = EXCLUDED.zip,
            city = EXCLUDED.city,
            address = EXCLUDED.address,
            region = EXCLUDED.region,
            email = EXCLUDED.email`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
        INSERT INTO payments (
            order_uid, transaction, request_id, currency, provider,
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        )
        ON CONFLICT (order_uid) DO UPDATE SET
            transaction = EXCLUDED.transaction,
            request_id = EXCLUDED.request_id,
            currency = EXCLUDED.currency,
            provider = EXCLUDED.provider,
            amount = EXCLUDED.amount,
            payment_dt = EXCLUDED.payment_dt,
            bank = EXCLUDED.bank,
            delivery_cost = EXCLUDED.delivery_cost,
            goods_total = EXCLUDED.goods_total,
            custom_fee = EXCLUDED.custom_fee`,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to delete old items: %w", err)
	}
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid,
                name, sale, size, total_price, nm_id, brand, status
            ) VALUES (
                $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
            )`,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status)
		if err != nil {
			return fmt.Errorf("failed to save item: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *Storage) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	var order models.Order
	var delivery models.Delivery
	var payment models.Payment
	var items []models.Item
	err := s.db.QueryRowContext(ctx, `
        SELECT 
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders 
        WHERE order_uid = $1`, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	err = s.db.QueryRowContext(ctx, `
        SELECT 
            name, phone, zip, city, address, region, email
        FROM deliveries 
        WHERE order_uid = $1`, orderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	order.Delivery = delivery
	err = s.db.QueryRowContext(ctx, `
        SELECT 
            transaction, request_id, currency, provider, amount,
            payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments 
        WHERE order_uid = $1`, orderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	order.Payment = payment

	rows, err := s.db.QueryContext(ctx, `
        SELECT 
            chrt_id, track_number, price, rid, name,
            sale, size, total_price, nm_id, brand, status
        FROM items 
        WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}
	order.Items = items
	return &order, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
