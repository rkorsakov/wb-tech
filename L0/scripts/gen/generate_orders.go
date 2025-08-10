package main

import (
	"L0/internal/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	baseOrder := `{
  "order_uid": "test123456789",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "test123456789",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`
	os.Mkdir("testdata", 0755)
	amountToGen, _ := strconv.Atoi(os.Args[1])
	for i := 0; i < amountToGen; i++ {
		var order models.Order
		json.Unmarshal([]byte(baseOrder), &order)
		order.OrderUID = generateRandomString(10)
		order.TrackNumber = fmt.Sprintf("%s", generateRandomString(6))
		order.Payment.Transaction = order.OrderUID
		order.DateCreated = time.Now().Format(time.RFC3339)
		jsonData, _ := json.MarshalIndent(order, "", "  ")
		filename := filepath.Join("testdata", fmt.Sprintf("order_%d.json", i+1))
		os.WriteFile(filename, jsonData, 0644)
	}
}

func generateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
