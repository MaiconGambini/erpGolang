.PHONY: dev test build fmt vet sqlc migrate-apply migrate-status

dev:
	go run ./cmd/api

test:
	go test ./...

build:
	go build ./...

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
