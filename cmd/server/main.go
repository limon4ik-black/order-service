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
	connDB, ctx := storage.InitConnDb(log) // ctx вместо _
	defer connDB.Close()
	//TODO: init kafka
	consumer := kafka.InitKafkaConsumer("localhost:9092", "orders", "order-service-group")
	defer consumer.Close()
	rdb := storage.InitConnCash()
	go kafka.StartConsumer(consumer, log, connDB, ctx, rdb)
	//TODO: init Redis

	http.Handle("/order", &handlers.GetOrderHandler{
		Conn: connDB,
		Ctx:  ctx,
		Log:  log,
		Rdb:  rdb,
	})

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Error("failed to listen", "error", err)
	}
}
