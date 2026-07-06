# Architecture

> Canonical overview: [`ARCHITECTURE.md`](../ARCHITECTURE.md) at repo root. Backend implementation detail: [`backend/docs/ARCHITECTURE.md`](../backend/docs/ARCHITECTURE.md).

## Subsystems

| Subsystem | Path | Responsibility |
|---|---|---|
| API | `backend/cmd/api` | HTTP entry, graceful shutdown |
| App composition | `backend/internal/app` | Module registration, deps wiring |
| Auth | `backend/internal/auth` | Login, refresh, logout, me |
| Users | `backend/internal/users` | Admin list/get/update users |
| Tenants | `backend/internal/tenants` | Current tenant info |
| Customers | `backend/internal/customers` | Reference CRUD module |
| Products | `backend/internal/products` | Catalog CRUD, low-stock list |
| Suppliers | `backend/internal/suppliers` | Supplier CRUD |
| Sales | `backend/internal/sales` | Draft sales, confirm/cancel, stock effects |
| Dashboard | `backend/internal/dashboard` | Tenant KPI aggregate read model |
| Reports | `backend/internal/reports` | Charts, period PDF, sale PDF read model |
| Audit | `backend/internal/audit` | Append-only write log + admin list API |
| Shared | `backend/internal/shared` | Errors, pagination, auth/tenant context, export, inventory |
| Platform | `backend/internal/platform` | DB, Redis, logger, validation, middleware (incl. RBAC) |
| Frontend | `frontend/src` | Vue 3 FSD app |

## Module Archetypes

| Archetype | Reference | Backend shape | Frontend shape |
|---|---|---|---|
| Entity CRUD | `customers`, `products`, `suppliers` | `module.go`, `handler.go`, `service.go`, sqlc queries | `entities/*`, `features/*/list|create|edit|delete`, `pages/*` |
| Workflow | `sales` | CRUD + `POST /{id}/confirm`, `POST /{id}/cancel`, DB transactions | `features/sale/actions`, confirm/cancel in UI |
| Read-model | `dashboard`, `reports` | Aggregate SQL, optional query params, CSV via `?format=csv` on lists | `entities/dashboard`, `entities/reports`, `pages/dashboard` |
| Admin | `users`, `audit` | Admin-only routes, tenant-scoped list/update | `pages/users`, `pages/audit`, `requireAdmin` guard |

Services own business rules and map sqlc rows to DTOs. sqlc generated queries are the persistence layer (no separate repository interface files today).

## Data Flow

```text
Browser → Vue (FSD) → entities/*/api → shared/api/client (axios)
       → /api/v1 → chi → AuthJWT → TenantScope → RequireRole → handler → service → sqlc → PostgreSQL
       → Redis (login rate limit only; not response cache)
```

Pinia `session` holds access token and user. Vue Query holds server state. Tenant scope is derived from JWT claims server-side; the frontend never sends `tenantId`.

## Integration Points

- PostgreSQL 16 — source of truth (`backend/schema.sql`, Atlas migrations)
- Redis 7 — login rate limiting (`POST /auth/login`); `/readyz` checks Redis
- JWT access (~15 min) + refresh HttpOnly cookie (hash in `auth_sessions`, rotation on refresh)
- Vite dev proxy: `/api` → `http://localhost:8080`; `VITE_API_BASE_URL=/api/v1`

## API Envelope

| Case | Shape |
|---|---|
| Single resource | `{ data: T }` |
| Paginated list | `{ data: T[], pagination }` |
| Error | `{ error: { code, message, details? } }` |
| CSV export | `?format=csv` on list endpoints (UTF-8 BOM) |

## Cross-Module Rules

- Modules do not import each other's Go packages.
- Multi-entity workflows may use sqlc across tables inside one service transaction (sales confirm/cancel).
- Read-model modules (`dashboard`, `reports`) use aggregate SQL over customers, sales, and products.
- Low-stock threshold default `5` is shared: `backend/internal/shared/inventory` and `frontend/src/shared/config/inventory.ts`.
- RBAC is enforced per `docs/ROLES.md` via `RequireRole` middleware and frontend `shared/lib/roles.ts`.

## Frontend Routes

| Route | Page |
|---|---|
| `/login` | Login |
| `/` | Dashboard KPIs + charts (financial gated) |
| `/customers` | Customer CRUD |
| `/products` | Product CRUD |
| `/suppliers` | Supplier CRUD |
| `/sales` | Sales list + draft/confirm/cancel |
| `/users` | User admin (admin only) |
| `/audit` | Audit log viewer (admin only) |

## Key Decisions

| Date | Decision | Reason |
|---|---|---|
| MVP | Multi-tenant via `tenant_id` | Isolation without subdomain complexity |
| MVP | chi/v5 router | Middleware ecosystem, explicit routing |
| MVP | sqlc + Atlas | Typed queries, declarative migrations |
| MVP | FSD frontend | Scalable slice isolation for ERP modules |
| MVP 1 | Three module archetypes | CRUD, workflow, read-model cover current domains |
| MVP 1 | Customers as CRUD reference | Products/suppliers mirror customers; sales adds workflow |
| Portfolio 8.5 | Reports module + gofpdf | Portfolio reporting without NF-e scope |
| Portfolio 8.5 | Route-level RBAC | Matches ROLES.md matrix; viewer read-only |

See also: `docs/CI_CD.md`, `docs/RELIABILITY.md`, `DEPLOYMENT.md`.
