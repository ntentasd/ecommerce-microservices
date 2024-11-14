package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/ntentasd/ecommerce-microservices/models"
)

var (
	ProductTopic = "ProductEvents"
	OrderTopic   = "OrderEvents"

	consumerGroupID = "product-consumer-group"
)

type Consumer struct {
	store *EventStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var eventType struct {
			Type string `json:"type"`
		}

		// First, unmarshal to get the event type
		err := json.Unmarshal(msg.Value, &eventType)
		if err != nil {
			slog.Error("Failed to unmarshal event type", "error", err)
			continue
		}

		// Now, based on the event type, unmarshal into the specific event struct
		switch eventType.Type {
		case string(ProductCreated):
			var productEvent models.ProductEvent
			err = json.Unmarshal(msg.Value, &productEvent)
			if err != nil {
				slog.Error("Failed to unmarshal ProductEvent", "error", err)
				continue
			}
			slog.Info("Received ProductEvent", "event", productEvent)
			consumer.store.Add(ProductCreated, productEvent)

		case string(OrderCreated):
			var orderEvent models.OrderEvent
			err = json.Unmarshal(msg.Value, &orderEvent)
			if err != nil {
				slog.Error("Failed to unmarshal OrderEvent", "error", err)
				continue
			}
			slog.Info("Received OrderEvent", "event", orderEvent)
			consumer.store.Add(OrderCreated, orderEvent)

		default:
			slog.Error("Unknown event type", "type", eventType.Type)
			continue
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func SetupConsumer(ctx context.Context, store *EventStore) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Admin.Timeout = 30 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup([]string{os.Getenv("KAFKA_BROKER")}, consumerGroupID, config)
	if err != nil {
		slog.Error("Failed to create consumer group", "error", err)
		return
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		store: store,
	}

	for {
		err = consumerGroup.Consume(ctx, []string{ProductTopic, OrderTopic}, consumer)
		if err != nil {
			slog.Error("Failed to consume messages", "error", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}
