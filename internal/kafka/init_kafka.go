package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"order-service/internal/models"

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

func StartConsumer(reader *kafka.Reader, log *slog.Logger) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error("Kafka read error", "error", err)
			continue
		}
		//log.Info("Kafka message", "value", string(m.Value))

		//valid
		// ПАРС

		orderStruct, err := ParseMessage(m, log)
		if err != nil {
			continue
		}
		log.Info("OMG"+orderStruct.Order_uid+orderStruct.Customer_id, orderStruct.Delivery.Name, orderStruct.Delivery.Zip)
		//...
		// ЗАЛИВ В БД
		// ЗАЛИВ В КЕШ

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
