run: build
	@. .envrc
	@./bin/main

build:
	@go build -o bin/main ./cmd

migrate-up:
	@migrate -path db/migrations -database $(DATABASE_URL) up

migrate-down:
	@migrate -path db/migrations -database $(DATABASE_URL) down

services:
	@docker compose up -d

.PHONY: build run migrate-up migrate-down services