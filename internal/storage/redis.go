package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitConnCash() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func ReloadCash(rdb *redis.Client, conn *pgxpool.Pool, ctx context.Context) error {

	query := `
        SELECT order_uid, track_number, entry, locale, internal_signature,
               customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM general ORDER BY date_created DESC LIMIT 5
    `
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to load orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order

		// general
		err := rows.Scan(
			&order.Order_uid,
			&order.Track_number,
			&order.Entry,
			&order.Locale,
			&order.Internal_signature,
			&order.Customer_id,
			&order.Delivery_service,
			&order.Shard_key,
			&order.Sm_id,
			&order.Date_created,
			&order.Oof_shard,
		)
		if err != nil {
			return fmt.Errorf("failed to scan general: %w", err)
		}

		// delivery
		var delivery models.Delivery
		err = conn.QueryRow(ctx,
			`SELECT name, phone, zip, city, address, region, email 
             FROM delivery WHERE order_uid = $1`, order.Order_uid,
		).Scan(
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			return fmt.Errorf("failed to scan delivery: %w", err)
		}
		order.Delivery = delivery

		// payment
		var payment models.Payment
		err = conn.QueryRow(ctx,
			`SELECT transaction, request_id, currency, provider, amount,
                    payment_dt, bank, delivery_cost, goods_total, custom_fee
             FROM payment WHERE order_uid = $1`, order.Order_uid,
		).Scan(
			&payment.Transaction,
			&payment.Request_id,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.Payment_dt,
			&payment.Bank,
			&payment.Delivery_cost,
			&payment.Goods_total,
			&payment.Custom_fee,
		)
		if err != nil {
			return fmt.Errorf("failed to scan payment: %w", err)
		}
		order.Payment = payment

		// items
		order.Items = []models.Item{}
		itemRows, err := conn.Query(ctx,
			`SELECT chrt_id, track_number, price, rid, name, sale, size,
                    total_price, nm_id, brand, status
             FROM items WHERE order_uid = $1`, order.Order_uid)
		if err != nil {
			return fmt.Errorf("failed to load items: %w", err)
		}
		for itemRows.Next() {
			var item models.Item
			err := itemRows.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.Rid,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				itemRows.Close()
				return fmt.Errorf("failed to scan item: %w", err)
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()

		// сериализация
		data, err := json.Marshal(order)
		if err != nil {
			return fmt.Errorf("failed to marshal order: %w", err)
		}

		// сохранение в Redis
		err = rdb.Set(ctx, order.Order_uid, data, 0).Err()
		if err != nil {
			return fmt.Errorf("failed to set redis key: %w", err)
		}
	}

	return nil
}
