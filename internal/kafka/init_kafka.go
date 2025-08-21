package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"order-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func InitKafkaConsumer(broker, topic, group string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  group,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func StartConsumer(reader *kafka.Reader, log *slog.Logger, conn *pgxpool.Pool, ctx context.Context, rdb *redis.Client) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error("Kafka read error", "error", err)
			continue
		}

		orderStruct, err := ParseMessage(m, log)

		if err != nil {
			continue
		}

		cashOrder, err := json.Marshal(orderStruct)
		if err != nil {
			log.Error("failed trans struct to json", "error", err)
			panic(err)
		}

		err = rdb.Set(ctx, orderStruct.Order_uid, cashOrder, 0).Err()
		if err != nil {
			log.Error("failed to insert into cash", "error", err)
			panic(err)
		}

		_, err = conn.Exec(ctx, "SELECT insert_general($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			orderStruct.Order_uid,
			orderStruct.Track_number,
			orderStruct.Entry,
			orderStruct.Locale,
			orderStruct.Internal_signature,
			orderStruct.Customer_id,
			orderStruct.Delivery_service,
			orderStruct.Shard_key,
			orderStruct.Sm_id,
			orderStruct.Date_created,
			orderStruct.Oof_shard)

		if err != nil {
			log.Error("failed to insert into general", "error", err)
			panic(err)
		}

		_, err = conn.Exec(ctx,
			"INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) "+
				"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			orderStruct.Order_uid,
			orderStruct.Delivery.Name,
			orderStruct.Delivery.Phone,
			orderStruct.Delivery.Zip,
			orderStruct.Delivery.City,
			orderStruct.Delivery.Address,
			orderStruct.Delivery.Region,
			orderStruct.Delivery.Email,
		)
		if err != nil {
			log.Error("failed to insert into delivery", "error", err)
			panic(err)
		}

		_, err = conn.Exec(ctx,
			"INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) "+
				"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			orderStruct.Order_uid,
			orderStruct.Payment.Transaction,
			orderStruct.Payment.Request_id,
			orderStruct.Payment.Currency,
			orderStruct.Payment.Provider,
			orderStruct.Payment.Amount,
			orderStruct.Payment.Payment_dt,
			orderStruct.Payment.Bank,
			orderStruct.Payment.Delivery_cost,
			orderStruct.Payment.Goods_total,
			orderStruct.Payment.Custom_fee,
		)
		if err != nil {
			log.Error("failed to insert into payment", "error", err)
			panic(err)
		}

		for _, item := range orderStruct.Items {
			_, err = conn.Exec(ctx,
				"INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) "+
					"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
				orderStruct.Order_uid,
				item.ChrtID,
				item.TrackNumber,
				item.Price,
				item.Rid,
				item.Name,
				item.Sale,
				item.Size,
				item.TotalPrice,
				item.NmID,
				item.Brand,
				item.Status,
			)
			if err != nil {
				log.Error("failed to insert into items", "error", err)
				panic(err)
			}
		}
		log.Info("вроде все")
	}
}

func ParseMessage(m kafka.Message, log *slog.Logger) (*models.Order, error) {

	var order models.Order
	err := json.Unmarshal(m.Value, &order)

	if err != nil {
		log.Error("failed to parse json", "error", err)
		return nil, err
	}

	return &order, nil
}
