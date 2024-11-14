run: build
	@. .envrc
	@./bin/main

run-consumer: build-consumer
	@./bin/consumer

build:
	@go build -o bin/main ./cmd

build-consumer:
	@go build -o bin/consumer ./consumer

daemon:
	CompileDaemon -directory=cmd/ -color=true -build="go build -o ../bin/cmd" -command="./bin/cmd"

migrate-up:
	@migrate -path db/migrations -database $(DATABASE_URL) up

migrate-down:
	@migrate -path db/migrations -database $(DATABASE_URL) down

services:
	@docker compose up -d

.PHONY: build run migrate-up migrate-down services