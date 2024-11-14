package kafka

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/IBM/sarama"
)

var (
	ProductTopic = "ProductEvents"
	OrderTopic   = "OrderEvents"
)

type Producer struct {
	sarama.SyncProducer
}

func SetupProducer() (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 5
	config.Admin.Timeout = 30 * time.Second

	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer}, nil
}

func (p *Producer) SendProductEvent(message any, key sarama.Encoder) error {
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: ProductTopic,
		Key:   key,
		Value: sarama.ByteEncoder(jsonMsg),
	}

	go func() {
		_, _, err = p.SendMessage(msg)
		if err != nil {
			slog.Error("Failed to send message", "error", err)
		}
	}()

	return nil
}

func (p *Producer) SendOrderEvent(message any, key sarama.Encoder) error {
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: OrderTopic,
		Key:   key,
		Value: sarama.ByteEncoder(jsonMsg),
	}

	go func() {
		_, _, err = p.SendMessage(msg)
		if err != nil {
			slog.Error("Failed to send message", "error", err)
		}
	}()

	return nil
}
