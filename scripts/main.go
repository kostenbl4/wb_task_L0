package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	p, err := NewProducer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}

	// Последовательное заполнение
	for i := 0; i < 1000; i++ {

		time.Sleep(time.Second)
		fakeOrder := generateFakeOrder()
		jsonData, err := json.Marshal(fakeOrder)
		if err != nil {
			panic(err)
		}
		if err := p.Produce(jsonData, "orders"); err != nil {
			panic(err)
		}
		slog.Info(fmt.Sprintf("Produced order_uid=%v", fakeOrder.OrderUID))
	}

	// if err := p.Produce([]byte("trash data"), "orders"); err != nil {
	// 	panic(err)
	// }
	p.Close()
}
