# Backend Local Development

Root README: five-step local run in [`README.md`](../../README.md#local-run-5-steps).

## 1. Environment

```bash
cp backend/.env.example backend/.env
# Set JWT_ACCESS_SECRET and JWT_REFRESH_SECRET (non-default values)
```

`backend/.env.example` defaults match `docker-compose.dev.yml` host ports:

| Service  | Container | Host port |
|----------|-----------|-----------|
| Postgres | `5432`    | **`5434`** |
| Redis    | `6379`    | **`6381`** |

```env
DATABASE_URL=postgres://goerp:goerp@127.0.0.1:5434/goerp?sslmode=disable
REDIS_URL=redis://127.0.0.1:6381/0
```

Use `127.0.0.1` instead of `localhost` on Windows to avoid IPv6 hitting a different Postgres instance.

CLI tools (`go run ./cmd/migrate`) read `DATABASE_URL` from the environment — load `.env` first or export vars before running migrate.

## 2. Infrastructure

From repository root:

```bash
docker compose -f docker-compose.dev.yml up -d
```

Verify compose (read-only): `docker compose -f docker-compose.dev.yml config`

## 3. Database

### Migrations

The canonical migration runner is `go run ./cmd/migrate` (used by CI and Docker). It applies versioned SQL files from `backend/migrations/` and tracks applied versions in the `schema_migrations` table.

Atlas (`atlas.hcl`) is optional and intended for schema diffing only — do not run both Atlas migrate and `cmd/migrate` against the same database, or you will get conflicting migration trackers.

```bash
cd backend
go run ./cmd/migrate
go run ./cmd/seed    # tenants acme/beta; admin@*.com + viewer@acme.com / admin123
make sqlc            # after query changes
```

## 4. Run API

```bash
cd backend
make dev             # air via .air.toml
# or: go run ./cmd/api
```

API: `http://localhost:8080` — `/healthz`, `/readyz`

## 5. Run Frontend

```bash
cd frontend
npm install
npm run dev          # http://localhost:5173, proxies /api → :8080
```

## Tests

```bash
cd backend
go test ./...
go test -tags=integration ./internal/customers/... ./internal/dashboard/... ./internal/sales/... ./internal/reports/...

cd frontend
npm run typecheck
npm run test:unit
npx playwright test  # 20 tests; Vite on :5174; defaults :5434/:6381 (see playwright.config.ts)
```

## Manual Smoke

1. Login as `acme` / `admin@acme.com` / `admin123`
2. Dashboard shows six KPIs (financial totals visible for admin)
3. CRUD flows: customers, products, suppliers, sales (draft → confirm)
4. CSV export on a list page; period PDF on dashboard (admin/manager)
5. Viewer smoke: `viewer@acme.com` / `admin123` — no create buttons; `/users` redirects
