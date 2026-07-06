# Backend Local Development

## Prerequisites

```bash
cp backend/.env.example backend/.env
# Set JWT_ACCESS_SECRET, JWT_REFRESH_SECRET, DATABASE_URL, REDIS_URL
```

## Infrastructure

From repository root:

```bash
docker compose -f docker-compose.dev.yml up -d
```

## Database

```bash
cd backend
go run ./cmd/migrate
go run ./cmd/seed    # tenants acme/beta, password admin123
make sqlc            # after query changes
```

## Run API

```bash
cd backend
make dev             # air via .air.toml
# or: go run ./cmd/api
```

## Run Frontend

```bash
cd frontend
npm install
npm run dev          # http://localhost:5173, proxies /api → :8080
```

## Tests

```bash
cd backend
go test ./...
go test -tags=integration ./internal/customers/... ./internal/dashboard/...

cd frontend
npm run typecheck
npm run test:unit
npx playwright test  # requires API + DB (see e2e.yml)
```

## Manual Smoke

1. Login as `acme` / `admin@acme.com` / `admin123`
2. Dashboard shows four KPIs
3. CRUD flows: customers, products, suppliers, sales (draft → confirm)
