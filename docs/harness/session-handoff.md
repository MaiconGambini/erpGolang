# Session Handoff

## Verified Now

```text
$ cd backend && go build -o api.exe ./cmd/api && go build -o migrate.exe ./cmd/migrate
(exit 0)

$ $env:DATABASE_URL="postgres://goerp:goerp@localhost:5432/goerp?sslmode=disable"; .\migrate.exe
skip 202601010001_initial.sql (already applied)
skip 202601020001_products.sql (already applied)
skip 202601030001_suppliers.sql (already applied)
apply 202601040001_sales.sql
migrations complete

$ cd backend && go test ./...
ok  internal/auth
ok  internal/customers
ok  internal/products
ok  internal/suppliers
ok  internal/sales

$ cd frontend && npm run typecheck
(exit 0)

$ cd frontend && npm run test:unit
Test Files  4 passed (4)
Tests  12 passed (12)

$ cd frontend && npx playwright test e2e/sales.spec.ts
2 passed (create+confirm, tenant isolation)

$ cd frontend && npx playwright test e2e/products.spec.ts
(3 passed in prior session)

$ cd frontend && npx playwright test e2e/suppliers.spec.ts
(3 passed in prior session)

$ cd frontend && npx playwright test e2e/auth-customers.spec.ts
(5 passed in prior session — full suite not re-run this session)
```

## Changed This Session

**Harness bootstrap**
- `AGENTS.md`, `feature_list.json`, `agent-os/`, `docs/{ARCHITECTURE,PRODUCT,RELIABILITY}.md`, `docs/harness/*`

**MVP phases 0–10 (prior work, uncommitted)**
- Backend: chi router, pgx, Redis, JWT auth, httpx, middleware, sqlc gen, customers CRUD, audit, seed/migrate CLIs
- Frontend: FSD auth (boot, guards, interceptors), customers CRUD UI, Playwright e2e
- Deploy: `fly.toml`, `Dockerfile.fly`, `Dockerfile.migrate`, prod compose migrate job, Caddy `/api/*` fix, `.github/workflows/*`

**Post-MVP — Products**
- `backend/schema.sql`, `migrations/202601020001_products.sql`, `queries/products.sql`, `gen/db/products.sql.go`
- `backend/internal/products/`, frontend `entities/product`, `features/product/*`, `pages/products`, `e2e/products.spec.ts`

**Post-MVP — Suppliers**
- `migrations/202601030001_suppliers.sql`, `queries/suppliers.sql`, `gen/db/suppliers.sql.go`
- `backend/internal/suppliers/`, frontend supplier FSD stack, `e2e/suppliers.spec.ts`

**Post-MVP — Sales**
- `migrations/202601040001_sales.sql`, `queries/sales.sql`, `gen/db/sales.sql.go`
- `backend/internal/sales/` — draft/confirmed/cancelled, stock decrement on confirm, restore on cancel
- Frontend `/sales`, create dialog (customer + product lines), confirm/cancel actions, `e2e/sales.spec.ts`

## Broken Or Unverified

- Full `npx playwright test` (14 tests) — not run as single suite this session; modules tested individually (5+3+3+2 PASS).
- `fly deploy` to real Fly.io — config only, not executed.
- Full prod `docker compose up --build` on VPS — not run.
- `go test -tags=integration ./internal/customers/...` — blocked on Windows App Control locally; intended for CI/Linux.
- Stale `api.exe` on `:8080` causes 404 for new routes if not rebuilt — Playwright webServer may reuse old binary (`reuseExistingServer: true` when `CI` unset).

## Decisions Made

- Post-MVP modules follow customers FSD template; sqlc output committed manually when Windows App Control blocks `sqlc.exe`.
- Sales prices snapshotted from product at draft creation; stock moves only on confirm/cancel (not on draft).
- Money exposed as JSON strings (`"19.90"`) with `NUMERIC(12,2)` in Postgres — no shopspring/decimal dependency.
- Migrate via `cmd/migrate` + `schema_migrations` table (not Atlas in local/CI path).
- Low-stock default threshold: 5 (`GET /products/low-stock?threshold=`).

## Next Best Step

Run full E2E suite once after `go build -o api.exe ./cmd/api`, then either **real Fly/VPS deploy** or **wire dashboard KPIs** to live APIs (low-stock count, draft sales count).

## Commands

```bash
# Infra
docker compose -f docker-compose.dev.yml up -d

# Migrations (mark initial if DB seeded manually)
docker exec goerp-postgres psql -U goerp -d goerp -c \
  "INSERT INTO schema_migrations (filename) VALUES ('202601010001_initial.sql') ON CONFLICT DO NOTHING;"
cd backend && go build -o migrate.exe ./cmd/migrate && .\migrate.exe

# API (rebuild after every backend module change)
cd backend && go build -o api.exe ./cmd/api && .\api.exe

# Seed (tenants acme/beta, admin123)
cd backend && go build -o seed.exe ./cmd/seed && .\seed.exe

# Validation
cd backend && go test ./...
cd frontend && npm run typecheck && npm run test:unit
cd frontend && npx playwright test

# Prod config check
GOERP_ENV_FILE="$(pwd)/deploy/env/production.env.example" GOERP_DOMAIN=example.com \
  docker compose -f deploy/compose/docker-compose.prod.yml config
```
