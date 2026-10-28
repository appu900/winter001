.PHONY: build-authservice run-authservice

AUTH_DIR := services/auth

build-authservice:
	go build -o bin/authservice ./$(AUTH_DIR)/cmd/server

run-authservice: build-authservice
	set -a && . ./$(AUTH_DIR)/.env && set +a && ./bin/authservice
