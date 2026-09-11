.PHONY: dev test build fmt vet sqlc migrate-apply migrate-status migrate test-integration

dev:
	go run ./cmd/api

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

build:
	go build ./...

migrate:
	go run ./cmd/migrate

fmt:
	gofmt -w .

vet:
	go vet ./...

sqlc:
	sqlc generate

migrate-apply:
	atlas migrate apply --env local

migrate-status:
	atlas migrate status --env local
