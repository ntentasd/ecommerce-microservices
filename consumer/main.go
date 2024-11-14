package main

import (
	"context"
	"log/slog"
	"time"
)

func main() {
	store := &EventStore{
		data: make(Events),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go SetupConsumer(ctx, store)

	slog.Info("Consumer running")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			go handleProductEvents(store)
			go handleOrderEvents(store)
		}
	}
}

func handleProductEvents(store *EventStore) {
	if store.IsEmpty(ProductCreated) {
		time.Sleep(1 * time.Second)
		return
	}

	event, err := store.Pop(ProductCreated)
	if err != nil {
		slog.Error("failed to pop event", "error", err)
	}

	slog.Info("Processed ProductEvent", "event", event)
}

func handleOrderEvents(store *EventStore) {
	if store.IsEmpty(OrderCreated) {
		time.Sleep(1 * time.Second)
		return
	}

	event, err := store.Pop(OrderCreated)
	if err != nil {
		slog.Error("failed to pop event", "error", err)
	}

	slog.Info("Processed OrderEvent", "event", event)
}
