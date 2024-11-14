run: build
	@. .envrc
	@./bin/main

build:
	@go build -o bin/main ./cmd

.PHONY: build run