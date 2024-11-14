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
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				handleEvents(ctx, store)
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	select {}
}

func handleEvents(ctx context.Context, store *EventStore) {
	go handleProductEvents(ctx, store)
	go handleOrderEvents(ctx, store)
}

func handleProductEvents(ctx context.Context, store *EventStore) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if store.IsEmpty(ProductCreated) {
				time.Sleep(1 * time.Second)
				continue
			}

			event, err := store.Pop(ProductCreated)
			if err != nil {
				slog.Error("failed to pop event", "error", err)
				continue
			}

			slog.Info("Processed ProductEvent", "event", event)
		}
	}
}

func handleOrderEvents(ctx context.Context, store *EventStore) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if store.IsEmpty(OrderCreated) {
				time.Sleep(1 * time.Second)
				continue
			}

			event, err := store.Pop(OrderCreated)
			if err != nil {
				slog.Error("failed to pop event", "error", err)
				continue
			}

			slog.Info("Processed OrderEvent", "event", event)
		}
	}
}
