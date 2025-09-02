package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"order-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type GetOrderHandler struct {
	Conn *pgxpool.Pool
	Rdb  *redis.Client
	Ctx  context.Context
	Log  *slog.Logger
}

func (h *GetOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "count is %d\n", h.n)
	order_uid := r.URL.Query().Get("order_uid")
	exists, err := h.Rdb.Exists(h.Ctx, order_uid).Result()
	if err != nil {
		h.Log.Error("failed to dostup k rdb", "error", err)
		panic(err)
	}

	if exists == 1 {

		valOrder, err := h.Rdb.Get(h.Ctx, order_uid).Result()
		if err != nil {
			panic(err) //
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(valOrder))
		h.Log.Info("ОТРАБОТАЛО")
	} else {
		var orderResponse models.Order

		//DELIVERY
		query := "SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1"
		err := h.Conn.QueryRow(h.Ctx, query, order_uid).Scan(
			&orderResponse.Delivery.Name,
			&orderResponse.Delivery.Phone,
			&orderResponse.Delivery.Zip,
			&orderResponse.Delivery.City,
			&orderResponse.Delivery.Address,
			&orderResponse.Delivery.Region,
			&orderResponse.Delivery.Email,
		)
		if err != nil {
			h.Log.Error("failed to get order from db", "error", err) // обработать нормально типо 500
			w.WriteHeader(http.StatusInternalServerError)
			//panic(err)
		}

		//PAYMENT
		query = "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1"
		err = h.Conn.QueryRow(h.Ctx, query, order_uid).Scan(
			&orderResponse.Payment.Transaction,
			&orderResponse.Payment.Request_id,
			&orderResponse.Payment.Currency,
			&orderResponse.Payment.Provider,
			&orderResponse.Payment.Amount,
			&orderResponse.Payment.Payment_dt,
			&orderResponse.Payment.Bank,
			&orderResponse.Payment.Delivery_cost,
			&orderResponse.Payment.Goods_total,
			&orderResponse.Payment.Custom_fee,
		)
		if err != nil {
			h.Log.Error("failed to get order from db", "error", err) // обработать нормально типо 500
			w.WriteHeader(http.StatusInternalServerError)
			//panic(err)
		}

		//ITEM NOT ITEMs
		var item models.Item
		query = "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1"
		err = h.Conn.QueryRow(h.Ctx, query, order_uid).Scan(
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
			h.Log.Error("failed to get order from db", "error", err) // обработать нормально типо 500
			w.WriteHeader(http.StatusInternalServerError)
			//panic(err)
		}

		//GENERAL
		query = "SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM general WHERE order_uid = $1"
		err = h.Conn.QueryRow(h.Ctx, query, order_uid).Scan(
			&orderResponse.Order_uid,
			&orderResponse.Track_number,
			&orderResponse.Entry,
			&orderResponse.Locale,
			&orderResponse.Internal_signature,
			&orderResponse.Customer_id,
			&orderResponse.Delivery_service,
			&orderResponse.Shard_key,
			&orderResponse.Sm_id,
			&orderResponse.Date_created,
			&orderResponse.Oof_shard,
		)
		if err != nil {
			h.Log.Error("failed to get order from db", "error", err) // обработать нормально типо 500
			w.WriteHeader(http.StatusInternalServerError)
			//panic(err)
		}

		//ITEMS
		orderResponse.Items = []models.Item{item}

		cashOrder, err := json.Marshal(orderResponse)
		if err != nil {
			h.Log.Error("failed trans struct to json", "error", err)
			panic(err)
		}

		err = h.Rdb.Set(h.Ctx, order_uid, cashOrder, 0).Err()
		if err != nil {
			h.Log.Error("failed to insert into cash", "error", err)
			panic(err)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderResponse)
		return
	}
}
