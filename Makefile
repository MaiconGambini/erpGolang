.PHONY: infra-up infra-down infra-logs backend-dev backend-test backend-build frontend-dev frontend-build validate

infra-up:
	docker compose -f docker-compose.dev.yml up -d

infra-down:
	docker compose -f docker-compose.dev.yml down

infra-logs:
	docker compose -f docker-compose.dev.yml logs -f

backend-dev:
	cd backend && go run ./cmd/api

backend-test:
	cd backend && go test ./...

backend-build:
	cd backend && go build ./...

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

validate:
	docker compose -f docker-compose.dev.yml config
	cd backend && go test ./... && go build ./...
	cd frontend && npm run build
