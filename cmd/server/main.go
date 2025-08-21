package main

// валидация и тесты остались
import (
	"net/http"
	"order-service/internal/handlers"
	"order-service/internal/kafka"
	"order-service/internal/logger"
	"order-service/internal/storage"
)

func main() {

	log := logger.InitLoggerSlogger("local")

	connDB, ctx := storage.InitConnDb(log)
	defer connDB.Close()

	consumer := kafka.InitKafkaConsumer("localhost:9092", "orders", "order-service-group")
	defer consumer.Close()

	rdb := storage.InitConnCash()
	err := storage.ReloadCash(rdb, connDB, ctx)
	if err != nil {
		log.Error("failed to reload cash", "error", err)
	}

	go kafka.StartConsumer(consumer, log, connDB, ctx, rdb)

	http.Handle("/order", &handlers.GetOrderHandler{
		Conn: connDB,
		Ctx:  ctx,
		Log:  log,
		Rdb:  rdb,
	})

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Error("failed to listen", "error", err)
	}
}
