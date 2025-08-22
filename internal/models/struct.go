package models

import "time"

type Order struct {
	Order_uid          string    `json:"order_uid" validate:"required"`
	Track_number       string    `json:"track_number" validate:"required"`
	Entry              string    `json:"entry" validate:"required"`
	Delivery           Delivery  `json:"delivery" validate:"required"`
	Payment            Payment   `json:"payment" validate:"required"`
	Items              []Item    `json:"items" validate:"required,min=1,dive"`
	Locale             string    `json:"locale" validate:"required"`
	Internal_signature string    `json:"internal_signature"`
	Customer_id        string    `json:"customer_id" validate:"required"`
	Delivery_service   string    `json:"delivery_service" validate:"required"`
	Shard_key          string    `json:"shard_key"`
	Sm_id              int       `json:"sm_id" validate:"required,gt=0"`
	Date_created       time.Time `json:"date_created" validate:"required"`
	Oof_shard          string    `json:"oof_shard" validate:"required"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction   string `json:"transaction" validate:"required"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency" validate:"required,len=3"`
	Provider      string `json:"provider" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
	Payment_dt    int    `json:"payment_dt" validate:"required,gt=0"`
	Bank          string `json:"bank" validate:"required"`
	Delivery_cost int    `json:"delivery_cost" validate:"gte=0"`
	Goods_total   int    `json:"goods_total" validate:"gte=0"`
	Custom_fee    int    `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"required,gt=0"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required,gt=0"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required,gte=0"`
	NmID        int    `json:"nm_id" validate:"required,gt=0"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}
