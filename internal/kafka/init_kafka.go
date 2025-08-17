package kafka

import (
	"context"
	"log/slog"

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
		log.Info("Kafka message", "value", string(m.Value))

		// ПАРС
		// ЗАЛИВ В БД
		// ЗАЛИВ В КЕШ

	}
}
