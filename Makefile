SHELL := /bin/bash
# include .env
# $(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

ARGS=$(filter-out $@,$(MAKECMDGOALS))

## Golang
GOCMD=go
GORUN=$(GOCMD) run

run-local:
	export GOSUMDB=off
	gofmt -s -w .
	set -o allexport; source configuration/local.env; set +o allexport && ${GORUN} main.go ${ARGS}

test:
	go test ./...

test-coverage:
	export GOSUMDB=off
	go test ./... -coverprofile=test_result/coverage-all.out && go tool cover -html=test_result/coverage-all.out

migration-setup:
	go get github.com/golang-migrate/migrate

migrate-up:
	migrate -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" -path ${PATH_MIGRATE} ${ARGS} up

migrate-new:
	migrate create -ext sql -dir ./db/migrations -seq ${ARGS}
