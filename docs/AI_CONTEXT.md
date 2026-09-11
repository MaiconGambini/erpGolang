# AI Context

goERP is a modular, tenant-aware ERP monorepo. **Scope authority:** `docs/PRODUCT.md`.

## Target Stack

- Backend: Go 1.25+, chi, pgx, sqlc, Atlas, Redis, JWT, bcrypt, slog, gofpdf.
- Frontend: Vite, Vue 3, TypeScript, Pinia, TanStack Vue Query, PrimeVue, Tailwind, Chart.js.
- Data: PostgreSQL 16 and Redis 7.
- Local infra: Docker Compose (`docker-compose.dev.yml`).

## Implemented Modules

| Module | Backend | Primary endpoints |
|---|---|---|
| auth | `internal/auth` | `POST /auth/login`, `/refresh`, `/logout`; `GET /auth/me` |
| tenants | `internal/tenants` | `GET /tenants/current` |
| users | `internal/users` | `GET/PATCH /users`, `GET /users/{id}` (admin) |
| customers | `internal/customers` | CRUD `/customers`; `?format=csv` |
| products | `internal/products` | CRUD `/products`; `GET /products/low-stock`; `?format=csv` |
| suppliers | `internal/suppliers` | CRUD `/suppliers`; `?format=csv` |
| sales | `internal/sales` | CRUD `/sales`; confirm/cancel; `?format=csv&from=&to=` |
| dashboard | `internal/dashboard` | `GET /dashboard/summary?threshold=` |
| reports | `internal/reports` | `GET /reports/sales-by-day`, `/top-products`, `/sales/{id}/pdf`, `/sales-summary.pdf` |
| audit | `internal/audit` | Write via `audit.Recorder`; `GET /audit-logs` (admin) |

**References:** `customers` (CRUD), `sales` (workflow), `dashboard` + `reports` (read-model).

## RBAC

- Middleware: `internal/platform/middleware/role.go` (`RequireRole`).
- Matrix: `docs/ROLES.md`.
- Frontend: `shared/lib/roles.ts`, `requireAdmin` on `/users` and `/audit`.
- Financial data (charts, period PDF, revenue KPIs): admin + manager only.
- Viewer seed for local development and intentional demo/test use only (never production): `viewer@acme.com` / `admin123`.

## Core Rules For Agents

- Never create a tenant-owned table without `tenant_id`.
- Never access tenant-owned data without filtering by `tenant_id`.
- Never trust tenant ID from a request body.
- API errors: `{ error: { code, message, details? } }`.
- Paginated lists: `{ data, pagination }`.
- CSV export: `?format=csv` with UTF-8 BOM (`internal/shared/export/csv.go`).
- Do not add dependencies without clear reason.
- Keep `shared/inventory` (Go) and `shared/config/inventory.ts` (FE) in sync.

## How To Add A Module

1. Read `docs/GLOSSARY.md`, `docs/BUSINESS_RULES.md`, `docs/MODULE_TEMPLATE.md`.
2. Choose archetype: Entity CRUD, Workflow, Read-model, or Admin.
3. Add schema/migration with tenant rules.
4. Add sqlc queries with explicit `tenant_id`.
5. Add `module.go`, `handler.go`, `service.go`.
6. Apply `RequireRole` guards per `docs/ROLES.md`.
7. Mirror DTOs in `frontend/src/entities/<module>/model/types.ts`.
8. Add features, page, sidebar route (with role guards if needed).
9. Integration test for tenant isolation (`//go:build integration`).
10. Wire Vue Query invalidation; call `invalidateDashboardSummary()` when KPIs change.

## Validation Expectations

```bash
docker compose -f docker-compose.dev.yml config

(
  cd backend
  go test ./... && go build ./...
)
# Bash/Git Bash: load DATABASE_URL from backend/.env before migration
(
  cd backend
  set -a && source .env && set +a
  go run ./cmd/migrate && go run ./cmd/seed
)
(
  cd backend
  go test -tags=integration ./internal/customers/... ./internal/sales/... ./internal/reports/... ./internal/dashboard/...
)

(
  cd frontend
  npm run typecheck && npm run build && npm run test:unit
)
(
  cd frontend
  npx playwright test   # 21 tests; global-setup migrates + seeds
)

make validate
```

Seed logins are intentional demo/test credentials for local development only and must never be used in production: `acme` / `admin@acme.com` / `admin123`; `viewer@acme.com` / `admin123`.
