package main

import (
	"fmt"
	"log/slog"
	"strings"

	"errors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var errUnkownEventType = errors.New("unknown event type")

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(message []byte, topic string) error {
	slog.Info(fmt.Sprintf("Producing message to topic: %s", topic))

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}

	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		slog.Error("Failed to produce message: " + err.Error())
		return err
	}

	e := <-kafkaChan

	switch ev := e.(type) {
	case *kafka.Message:
		slog.Info("Message produced successfully")
		return nil
	case kafka.Error:
		slog.Error("Kafka error: " + ev.Error())
		return ev
	default:
		slog.Error("Encountered unknown event type")
		return errUnkownEventType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(1000)
	p.producer.Close()
}
