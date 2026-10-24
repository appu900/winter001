-include .env
export

.PHONY: build-api run-api

build-api:
	go build -o bin/api ./cmd/api

run-api:
	go run ./cmd/api
