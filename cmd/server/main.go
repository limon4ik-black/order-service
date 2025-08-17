package main

import (
	"net/http"
	"order-service/internal/handlers"
	"order-service/internal/kafka"
	"order-service/internal/logger"
	"order-service/internal/storage"
)

func main() {
	//TODO: logger
	log := logger.InitLoggerSlogger("local")
	//TODO: init DB
	connDB, _ := storage.InitConnDb(log) // ctx вместо _
	defer connDB.Close()
	//TODO: init kafka
	consumer := kafka.InitKafkaConsumer("localhost:9092", "orders", "order-service-group")
	defer consumer.Close()
	go kafka.StartConsumer(consumer, log)
	//TODO: init Redis

	http.Handle("/order", new(handlers.Order))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Error("failed to listen", "error", err)
	}
}
