.PHONY: build-authservice run-authservice

AUTH_DIR := services/auth
AUTH_MIGRATIONS := $(AUTH_DIR)/db/migrations
LOAD_AUTH_ENV := set -a && . ./$(AUTH_DIR)/.env && set +a


build-authservice:
	go build -o bin/authservice ./$(AUTH_DIR)/cmd/server

run-authservice: build-authservice
	set -a && . ./$(AUTH_DIR)/.env && set +a && ./bin/authservice

run-auth-migrate-new:
	@test -n "$(name)" || (echo "usage: make run-auth-migrate-new name=add_phone_to_users" && exit 1)
	migrate create -ext sql -dir ${AUTH_MIGRATIONS} -seq ${name}
	
auth-migrate-up:
	$(LOAD_AUTH_ENV) && migrate -path $(AUTH_MIGRATIONS) -database "$$AUTH_DATABASE_URL" up

auth-migrate-down:
	$(LOAD_AUTH_ENV) && migrate -path $(AUTH_MIGRATIONS) -database "$$AUTH_DATABASE_URL" down 1

