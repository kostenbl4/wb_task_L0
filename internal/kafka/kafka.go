package kafka

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	handleFunc func(msg []byte) error
	run         bool
}

func NewKafkaConsumer(topic, group string, address []string, handlerFunc func(msg []byte) error) (*KafkaConsumer, error) {
	slog.Info(fmt.Sprintf("Creating new Kafka consumer for topic: %s, group: %s", topic, group))

	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
		"group.id":          group,
		"auto.offset.reset": "earliest",
	}

	slog.Info("Kafka consumer configuration set")

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		slog.Error("Failed to create Kafka consumer: " + err.Error())
		return nil, err
	}

	slog.Info("Kafka consumer created successfully")

	err = c.Subscribe(topic, nil)
	if err != nil {
		slog.Error("Failed to subscribe to topic: " + err.Error())
		return nil, err
	}

	slog.Info(fmt.Sprintf("Subscribed to topic: %s", topic))

	return &KafkaConsumer{consumer: c, handleFunc: handlerFunc}, nil
}

func (c *KafkaConsumer) Start() {
	c.run = true
	slog.Info("Starting Kafka consumer loop")
	for c.run {
		slog.Info("Waiting for a message from Kafka")
		kafkaMsg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		if kafkaMsg == nil {
			slog.Info("No message received from Kafka, waiting for the next one")
			continue
		}

		slog.Info(fmt.Sprintf("Received message from kafka topic %v partition %v at offset %v", kafkaMsg.TopicPartition.Topic, kafkaMsg.TopicPartition.Partition, kafkaMsg.TopicPartition.Offset))
		if err := c.handleFunc(kafkaMsg.Value); err != nil {
			slog.Error(err.Error())
		}
	}
}

func (c *KafkaConsumer) Stop() error {
	c.run = false
	return c.consumer.Close()
}
