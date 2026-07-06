# goERP

Modular, tenant-aware ERP monorepo for small and medium businesses — built with Go, Vue 3, PostgreSQL, and Redis.

Inspired by learning-oriented READMEs like [person-crud](https://github.com/KozielGPC/person-crud), this document starts with **what the project teaches**, then describes **how the system works**.

---

## What I Learned Building This Project

### Backend (Go)

- **Modular monolith composition** — registering domain modules (`auth`, `customers`, `products`, `suppliers`, `sales`, `dashboard`) behind a single `app.Module` interface and chi router.
- **Multi-tenant isolation** — every tenant-owned table has `tenant_id`; JWT claims feed `tenantctx`; cross-tenant access returns `404` (not `403`).
- **Typed SQL with sqlc** — queries live in `.sql` files; generated Go code removes stringly-typed SQL in handlers.
- **Schema migrations with Atlas** — declarative `schema.sql` + versioned migrations; migrate job runs before app boot in production.
- **Auth with JWT + refresh rotation** — short-lived access token in memory; refresh token as HttpOnly cookie; SHA-256 hash stored in Postgres; session revoked on rotation.
- **Transactional workflows** — sales `confirm` / `cancel` use `pgx` transactions to update status and product stock atomically.
- **Read-model aggregates** — dashboard KPIs as a single tenant-scoped SQL query instead of N+1 API calls.
- **Structured logging** — `slog` JSON logs with `request_id`, `tenant_id`, `user_id`, latency.
- **Rate limiting** — Redis counter on login (5 / 15 min per IP); fail-open when Redis is down.
- **Audit trail** — append-only `audit_logs` on business writes via injected `audit.Recorder`.
- **Testing layers** — handler tests, integration tests with `//go:build integration`, golangci-lint in CI.

### Frontend (Vue 3)

- **Feature-Sliced Design (FSD)** — `entities` → `features` → `widgets` → `pages`; imports flow downward only.
- **TanStack Vue Query** — server state, cache invalidation, dashboard refresh after mutations.
- **Pinia session store** — access token in memory; axios interceptor with single-flight refresh on `401`.
- **PrimeVue + Tailwind** — CRUD tables, dialogs, forms with consistent UX patterns.
- **E2E with Playwright** — auth, customers, products, suppliers, sales, dashboard KPI smoke tests.

### DevOps & Reliability

- **Docker Compose** — local Postgres + Redis; production full stack with Caddy TLS.
- **CI pipelines** — backend lint/test/integration, frontend typecheck/unit/build, Playwright E2E.
- **CD to Fly.io** — deploy gated on green Backend CI (`workflow_run`).
- **Two deploy profiles** — VPS full stack (default) vs Fly API split (see `docs/CI_CD.md`).

### Architecture Patterns

| Pattern | Where |
|---|---|
| Entity CRUD | `customers`, `products`, `suppliers` |
| Workflow + transactions | `sales` (draft → confirm → cancel) |
| Read-model aggregate | `dashboard` (KPI summary) |
| Shared invariants | low-stock threshold in Go + TypeScript |

---

## The System

goERP MVP 1 delivers authentication, tenant isolation, catalog CRUD, sales with stock effects, and a live dashboard.

```text
Browser
  → Vue 3 (FSD) + Vue Query + Pinia
  → REST API (/api/v1)
  → Go modular monolith (chi)
  → PostgreSQL 16 (source of truth)
  → Redis 7 (login rate limit)
```

### Backend Modules

| Module | Responsibility |
|---|---|
| `auth` | Login, refresh, logout, me |
| `tenants` | Current tenant info |
| `customers` | Reference CRUD module |
| `products` | Catalog, stock, low-stock list |
| `suppliers` | Supplier CRUD |
| `sales` | Draft sales, confirm (stock −), cancel (stock +) |
| `dashboard` | KPI aggregate (`active_customers`, `new_customers_30d`, `draft_sales`, `low_stock_alerts`) |
| `audit` | Write-side event log (no HTTP routes) |
| `users` | Registered; CRUD routes deferred |

### API Conventions

- Single resource: `{ data: T }`
- Paginated list: `{ data: T[], pagination }`
- Error: `{ error: { code, message, details? } }`
- Tenant scope from JWT — never from request body

### Frontend Routes

| Route | Screen |
|---|---|
| `/login` | Tenant slug + email + password |
| `/` | Dashboard KPIs |
| `/customers` | Customer CRUD |
| `/products` | Product CRUD |
| `/suppliers` | Supplier CRUD |
| `/sales` | Sales list, confirm, cancel |

### Auth Flow

1. `POST /auth/login` with `tenantSlug`, `email`, `password` → access token + user JSON.
2. Refresh token set as HttpOnly cookie (`Secure` in production).
3. Axios sends `Authorization: Bearer`; on `401`, single-flight `POST /auth/refresh` then retry.
4. Logout clears cookie and `queryClient.clear()`.

### Database

PostgreSQL tables: `tenants`, `users`, `auth_sessions`, `customers`, `products`, `suppliers`, `sales`, `sale_items`, `audit_logs`.

Sales lifecycle: `draft` → `confirmed` → `cancelled`. Stock changes only on confirm/cancel.

Seed tenants: `acme` / `beta` — password `admin123` (`go run ./cmd/seed`).

---

## Usage

### Prerequisites

- Go 1.25+
- Node.js 20+
- Docker (for Postgres + Redis)

### Local run (5 steps)

`docker-compose.dev.yml` publishes **Postgres on host port `5434`** and **Redis on `6381`** (avoids conflicts with `5432`/`6379` or `goerp-prod`). Full backend runbook: [`backend/docs/LOCAL_DEV.md`](backend/docs/LOCAL_DEV.md).

**1. Environment**

```bash
cp backend/.env.example backend/.env
# Set JWT_ACCESS_SECRET and JWT_REFRESH_SECRET (non-default values)
```

Default connection strings in `.env.example` already target `5434` / `6381`. On Windows, prefer `127.0.0.1` over `localhost` in `DATABASE_URL` to avoid IPv6 hitting a different Postgres instance.

**2. Infrastructure**

```bash
docker compose -f docker-compose.dev.yml up -d
```

Verify compose file (read-only): `docker compose -f docker-compose.dev.yml config`

**3. Database**

```bash
cd backend
go run ./cmd/migrate
go run ./cmd/seed    # tenants acme/beta — password admin123
```

**4. Backend**

```bash
cd backend
make dev          # hot reload via air
# or: go run ./cmd/api
```

API: `http://localhost:8080` — health: `/healthz`, `/readyz`

**5. Frontend**

```bash
cd frontend
npm install
npm run dev
```

App: `http://localhost:5173` (proxies `/api` → backend). Login: `acme` / `admin@acme.com` / `admin123`.

### Tests

```bash
# Backend
cd backend && go test ./...
go test -tags=integration ./internal/customers/... ./internal/dashboard/...

# Frontend
cd frontend && npm run typecheck && npm run test:unit

# E2E (docker-compose.dev.yml up; Playwright starts API on :8080 and Vite on :5174)
cd frontend && npx playwright test   # 15 tests; defaults DATABASE_URL :5434, REDIS_URL :6381

# Full gate
make validate
```

---

## Target Stack

| Layer | Technologies |
|---|---|
| Backend | Go 1.25+, chi, pgx, sqlc, Atlas, Redis, JWT, bcrypt, slog |
| Frontend | Vite, Vue 3, TypeScript, Pinia, TanStack Vue Query, PrimeVue, Tailwind |
| Data | PostgreSQL 16, Redis 7 |
| Deploy | Docker Compose + Caddy (VPS), Fly.io (API CD) |

---

## Repository Layout

```text
backend/                 Go API, migrations, sqlc queries
frontend/                Vue 3 FSD application
docs/                    Product, architecture, UX, CI/CD
backend/docs/            Backend-specific runbooks
deploy/                  Production compose + Caddy
docker-compose.dev.yml   Local PostgreSQL and Redis
```

---

## Tasks

### Backend

- [x] Modular monolith with `app.Module` registration
- [x] Multi-tenant schema and middleware
- [x] Auth: login, refresh rotation, logout, JWT
- [x] Customers CRUD + audit + integration tests
- [x] Products CRUD + low-stock endpoint
- [x] Suppliers CRUD
- [x] Sales workflow (confirm/cancel + stock transactions)
- [x] Dashboard KPI aggregate
- [x] Login rate limiting (Redis)
- [x] Structured logging + health/readiness
- [x] golangci-lint + integration tests in CI
- [ ] Users CRUD (deferred)
- [ ] OpenAPI contract (`contract/`)

### Frontend

- [x] FSD structure with entity API clients
- [x] Auth boot + axios refresh interceptor
- [x] CRUD pages: customers, products, suppliers, sales
- [x] Dashboard with live KPIs + cache invalidation
- [x] Loading, empty, error states
- [x] Playwright E2E suite (15 tests)
- [x] Sale draft edit UI (`EditSaleDialog` + PATCH)
- [ ] Drill-down links from dashboard metrics

### DevOps

- [x] GitHub Actions: backend, frontend, E2E workflows
- [x] Fly deploy gated on Backend CI
- [x] VPS production compose + Caddy
- [ ] VPS deploy automation in CI
- [ ] Backup/restore drill documented and tested

---

## Documentation

| Doc | Purpose |
|---|---|
| [`ARCHITECTURE.md`](ARCHITECTURE.md) | System shape and module boundaries |
| [`DEPLOYMENT.md`](DEPLOYMENT.md) | VPS and production deploy |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Detailed subsystems and API envelope |
| [`docs/CI_CD.md`](docs/CI_CD.md) | Workflows and deploy profiles |
| [`docs/AI_CONTEXT.md`](docs/AI_CONTEXT.md) | Context for AI agents |
| [`docs/BUSINESS_RULES.md`](docs/BUSINESS_RULES.md) | Domain invariants |
| [`docs/MODULE_TEMPLATE.md`](docs/MODULE_TEMPLATE.md) | How to add a module |
| [`docs/UX_PATTERNS.md`](docs/UX_PATTERNS.md) | UI behavior patterns |
| [`backend/docs/`](backend/docs/) | Auth, schema, tenant isolation, sales transactions |

---

## Next Steps

- Run full Playwright suite in CI locally
- Set `FLY_API_TOKEN` for live Fly deploy
- Fix sales concurrent-confirm race (see `backend/docs/SALES_TRANSACTIONS.md`)
- Add `agent-os/specs/` per module with acceptance criteria
