# AI Context

goERP is a modular, tenant-aware ERP monorepo. **Scope authority:** `docs/PRODUCT.md`.

## Target Stack

- Backend: Go 1.25+, chi, pgx, sqlc, Atlas, Redis, JWT, bcrypt, slog.
- Frontend: Vite, Vue 3, TypeScript, Pinia, TanStack Vue Query, PrimeVue, Tailwind CSS.
- Data: PostgreSQL 16 and Redis 7.
- Local infra: Docker Compose for infrastructure only.

## Implemented Modules (reference paths)

| Module | Backend | Primary endpoints |
|---|---|---|
| auth | `internal/auth` | `POST /auth/login`, `/refresh`, `/logout`; `GET /auth/me` |
| tenants | `internal/tenants` | `GET /tenants/current` |
| customers | `internal/customers` | CRUD `/customers` |
| products | `internal/products` | CRUD `/products`, `GET /products/low-stock` |
| suppliers | `internal/suppliers` | CRUD `/suppliers` |
| sales | `internal/sales` | CRUD `/sales`, `POST /sales/{id}/confirm`, `/cancel` |
| dashboard | `internal/dashboard` | `GET /dashboard/summary?threshold=` |
| audit | `internal/audit` | Write-only via `audit.Recorder` |

**CRUD reference:** `customers`. **Workflow reference:** `sales`. **Read-model reference:** `dashboard`.

## Core Rules For Agents

- Never create a tenant-owned table without `tenant_id`.
- Never access tenant-owned data without filtering by `tenant_id`.
- Never trust tenant ID from a request body.
- Keep API errors in `{ error: { code, message, details } }` format.
- Keep paginated responses in `{ data, pagination }` format.
- Do not add new dependencies without a clear reason.
- When stock thresholds matter, keep `shared/inventory` (Go) and `shared/config/inventory.ts` (FE) in sync.

## How To Add A Module

1. Read `docs/GLOSSARY.md`, `docs/BUSINESS_RULES.md`, and `docs/MODULE_TEMPLATE.md`.
2. Choose archetype: Entity CRUD, Workflow, or Read-model.
3. Add schema/migration with tenant rules.
4. Add sqlc queries with explicit `tenant_id` parameters.
5. Add `module.go`, `handler.go`, `service.go` (DTOs in service; no separate domain/repository files unless testing requires it).
6. Mirror DTOs in `frontend/src/entities/<module>/model/types.ts`.
7. Add features, page, sidebar route.
8. Verify tenant isolation (integration test with `//go:build integration`).
9. Wire Vue Query invalidation; call `invalidateDashboardSummary()` when KPIs are affected.

## Validation Expectations

```bash
# Infra
docker compose -f docker-compose.dev.yml config

# Backend
cd backend && go test ./... && go build ./...
cd backend && go run ./cmd/seed   # tenants acme/beta, password admin123
go test -tags=integration ./internal/customers/... ./internal/dashboard/...

# Frontend
cd frontend && npm run typecheck && npm run build && npm run test:unit

# E2E (dockerized DB + API)
cd frontend && npx playwright test

# Full gate
make validate
```
