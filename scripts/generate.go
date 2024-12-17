package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kostenbl4/wb_task_L0/internal/store"
)

func generateFakeOrder() store.Order {
	gofakeit.Seed(time.Now().UnixNano())

	orderUID := gofakeit.UUID()
	requestId := strconv.Itoa(gofakeit.Number(1000000, 9999999))
	internalSignature := gofakeit.LetterN(10)

	items := []store.Item{}
	nItems := gofakeit.Number(1, 5)
	for i := 0; i < nItems; i++ {
		items = append(items, store.Item{
			ChrtID:      gofakeit.Number(1000000, 9999999),
			TrackNumber: strconv.Itoa(gofakeit.Number(100000000, 999999999)),
			Price:       gofakeit.Number(50, 1500),
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.Number(0, 50),
			Size:        strconv.Itoa(gofakeit.Number(0, 5)),
			TotalPrice:  gofakeit.Number(50, 1500),
			NmID:        gofakeit.Number(1000000, 9999999),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.Number(100, 300),
		})
	}

	return store.Order{
		OrderUID:    orderUID,
		TrackNumber: strconv.Itoa(gofakeit.Number(100000000, 999999999)),
		Entry:       "WBIL",
		Delivery: store.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Address,
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: store.Payment{
			Transaction:  orderUID,
			RequestID:    &requestId,
			Currency:     gofakeit.CurrencyShort(),
			Provider:     "wbpay",
			Amount:       gofakeit.Number(100, 15000),
			PaymentDt:    time.Now().Unix(),
			Bank:         gofakeit.Company(),
			DeliveryCost: gofakeit.Number(100, 500),
			GoodsTotal:   gofakeit.Number(50, 3000),
			CustomFee:    gofakeit.Number(0, 100),
		},
		Items: items,
		Locale:            gofakeit.Language(),
		InternalSignature: &internalSignature,
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   "meest",
		ShardKey:          fmt.Sprint(gofakeit.Number(1, 10)),
		SmID:              gofakeit.Number(1, 100),
		DateCreated:       gofakeit.Date(),
		OofShard:          fmt.Sprint(gofakeit.Number(1, 10)),
	}
}
