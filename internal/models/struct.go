package models

type Order struct {
	Order_uid          string   `json:"order_uid"`
	Track_number       string   `json:"track_number"`
	Entry              string   `json:"entry"`
	Delivery           Delivery `json:"delivery"`
	Payment            Payment  `json:"payment"`
	Items              []Item   `json:"items"`
	Locale             string   `json:"locale"`
	Internal_signature string   `json:"internal_signature"`
	Customer_id        string   `json:"customer_id"`
	Delivery_service   string   `json:"delivery_string"`
	Shard_key          string   `json:"shard_key"`
	Sm_id              int      `json:"sm_key"`
	Date_created       string   `json:"date_created"`
	Oof_shard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        int    `json:"amount"`
	Payment_dt    int    `json:"payment_dt"`
	Bank          string `json:"bank"`
	Delivery_cost int    `json:"delivery_cost"`
	Goods_total   int    `json:"goods_total"`
	Custom_fee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
