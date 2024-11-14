package main

import (
	"log/slog"
	"net/http"

	database "github.com/ntentasd/ecommerce-microservices/pkg/database"
	"github.com/ntentasd/ecommerce-microservices/pkg/kafka"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	db       *database.Database
	redis    *redis.Client
	port     string
	version  string
	producer *kafka.Producer
}

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}

	redis, err := database.NewRedisClient()
	if err != nil {
		slog.Error("Failed to connect to redis", "error", err)
		return
	}

	producer, err := kafka.SetupProducer()
	if err != nil {
		slog.Error("Failed to setup kafka producer", "error", err)
		return
	}

	app := &Application{
		db:       db,
		redis:    redis,
		producer: producer,
		port:     ":8000",
		version:  "1.0.0",
	}

	srv := http.Server{
		Addr:    app.port,
		Handler: app.registerRoutes(),
	}

	slog.Info("Starting server on port " + app.port)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
		app.db.Close()
		app.redis.Close()
	}
}
