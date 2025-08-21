package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"order-service/internal/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func randomName() string {
	names := []string{"Alex", "Ivan", "Sergey", "John", "Anna", "Maria", "Kate"}
	surnames := []string{"Petrov", "Ivanov", "Smith", "Johnson", "Brown"}
	return names[rand.Intn(len(names))] + " " + surnames[rand.Intn(len(surnames))]
}

func randomCity() string {
	cities := []string{"Moscow", "Berlin", "New York", "Tokyo", "Paris", "London"}
	return cities[rand.Intn(len(cities))]
}

func randomBrand() string {
	brands := []string{"Nike", "Adidas", "Puma", "Reebok", "Vivienne Sabo", "L'Oreal"}
	return brands[rand.Intn(len(brands))]
}

func randomBank() string {
	banks := []string{"alpha", "sber", "tinkoff", "vtb", "raiffeisen"}
	return banks[rand.Intn(len(banks))]
}

func generateOrder() models.Order {
	order := models.Order{
		Order_uid:    uuid.NewString(),
		Track_number: "TRACK-" + strconv.Itoa(rand.Intn(1000000)),
		Entry:        "WBIL",
		Delivery: models.Delivery{
			Name:    randomName(),
			Phone:   "+972" + strconv.Itoa(1000000+rand.Intn(8999999)),
			Zip:     strconv.Itoa(100000 + rand.Intn(900000)),
			City:    randomCity(),
			Address: "Street " + strconv.Itoa(rand.Intn(100)),
			Region:  "Region-" + strconv.Itoa(rand.Intn(50)),
			Email:   "user" + strconv.Itoa(rand.Intn(1000)) + "@mail.com",
		},
		Payment: models.Payment{
			Transaction:   uuid.NewString(),
			Request_id:    "",
			Currency:      "USD",
			Provider:      "wbpay",
			Amount:        100 + rand.Intn(2000),
			Payment_dt:    int(time.Now().Unix()),
			Bank:          randomBank(),
			Delivery_cost: 100 + rand.Intn(2000),
			Goods_total:   50 + rand.Intn(500),
			Custom_fee:    0,
		},
		Items:              []models.Item{},
		Locale:             "en",
		Internal_signature: "",
		Customer_id:        "cust" + strconv.Itoa(rand.Intn(1000)),
		Delivery_service:   "meest",
		Shard_key:          strconv.Itoa(rand.Intn(10)),
		Sm_id:              rand.Intn(200),
		Date_created:       time.Now().UTC(),
		Oof_shard:          strconv.Itoa(rand.Intn(5)),
	}

	// генерим от 1 до 5 товаров
	n := rand.Intn(5) + 1
	for i := 0; i < n; i++ {
		item := models.Item{
			ChrtID:      rand.Intn(9999999),
			TrackNumber: "ITEM-" + strconv.Itoa(rand.Intn(100000)),
			Price:       50 + rand.Intn(1000),
			Rid:         uuid.NewString(),
			Name:        "Product-" + strconv.Itoa(rand.Intn(100)),
			Sale:        rand.Intn(50),
			Size:        strconv.Itoa(rand.Intn(5)),
			TotalPrice:  50 + rand.Intn(1000),
			NmID:        rand.Intn(999999),
			Brand:       randomBrand(),
			Status:      100 + rand.Intn(200),
		}
		order.Items = append(order.Items, item)
	}

	return order
}

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
	})
	defer writer.Close()

	rand.Seed(time.Now().UnixNano())

	for {
		order := generateOrder()

		data, err := json.Marshal(order)
		if err != nil {
			log.Println("Ошибка сериализации:", err)
			continue
		}

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(order.Order_uid),
				Value: data,
			},
		)
		if err != nil {
			log.Println("Ошибка отправки:", err)
		} else {
			log.Println("Отправлен заказ:", order.Order_uid)
		}
		time.Sleep(1 * time.Second)
	}
}
