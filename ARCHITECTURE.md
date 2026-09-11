# Architecture

goERP is a modular monolith with a REST API and a Vue frontend. The architecture optimizes for fast MVP delivery, strong tenant isolation, RBAC, portfolio-grade reporting, and a clear path to deployment without microservices.

## System Shape

```text
Browser
  -> Vue 3 (FSD) + Vue Query + Pinia
  -> REST API (/api/v1)
  -> Go modular monolith (chi)
  -> PostgreSQL 16

Go backend
  -> Redis 7 (login rate limit)
```

## Backend Architecture

The backend is a modular monolith. Each module owns handlers, service, and sqlc queries. DTOs live in `service.go`. sqlc is the persistence layer.

### Modules

| Module | Type | Notes |
|---|---|---|
| `auth` | Identity | Login, refresh, logout, JWT |
| `tenants` | Identity | Current tenant info |
| `users` | Admin | List/get/update users (admin only) |
| `customers` | Entity CRUD | Reference module |
| `products` | Entity CRUD | Stock, low-stock list |
| `suppliers` | Entity CRUD | Party fields (BR) |
| `sales` | Workflow | Draft → confirm → cancel; stock effects |
| `dashboard` | Read-model | 6 KPI aggregates |
| `reports` | Read-model | Charts, PDF exports |
| `audit` | Cross-cutting | Write recorder + admin list API |

Dependency direction:

```text
cmd/api -> internal/app -> modules -> shared/platform -> gen/db
```

### Rules

- Handlers translate HTTP and call services.
- Services enforce business rules and tenant isolation.
- Services map sqlc rows to DTOs in `service.go`.
- Modules do not import each other's Go packages.
- Cross-table reads/writes happen via sqlc inside service transactions (sales) or dedicated read-model modules (dashboard, reports).
- RBAC via `RequireRole` middleware per `docs/ROLES.md`.
- Shared packages contain infrastructure-neutral helpers (`export`, `inventory`, `httpx`).

## Frontend Architecture

Feature-Sliced Design:

```text
app -> processes -> pages -> widgets -> features -> entities -> shared
```

Rules:

- Imports only flow downward.
- `entities` mirror API DTOs and API clients.
- `features` own user actions (create customer, confirm sale).
- `widgets` compose tables and page sections.
- `pages` compose widgets; minimal business logic.
- `shared` contains API client, roles, export helpers, UI primitives.

Role helpers live in `shared/lib/roles.ts`; admin routes use `requireAdmin` guard.

## Tenant Isolation

Every tenant-owned table has `tenant_id`. Every tenant-owned query filters by `tenant_id`. Tenant IDs come from JWT claims, never from request bodies.

Cross-tenant access returns `404` for tenant-owned resources.

## Reporting

| Capability | Location |
|---|---|
| CSV list export | `?format=csv` on customers, products, suppliers, sales |
| Dashboard KPIs | `GET /dashboard/summary` |
| Charts | `GET /reports/sales-by-day`, `/top-products` |
| PDF | `GET /reports/sales/{id}/pdf`, `/reports/sales-summary.pdf` |

Financial aggregates (revenue charts, period PDF, confirmed-sales KPIs) are limited to admin and manager roles.

## Persistence

PostgreSQL is the source of truth. Redis is used for login rate limiting and `/readyz` — not as domain record storage.

## Deployment Path

1. **Local** — Docker Compose (`docker-compose.dev.yml`, Postgres `5434`, Redis `6381`).
2. **VPS** — Docker Compose + Caddy (`deploy/`).
3. **Fly.io** — API CD gated on Backend CI (`fly.toml`).

See `DEPLOYMENT.md` and `docs/CI_CD.md`.

## Related Docs

| Doc | Purpose |
|---|---|
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Detailed subsystems, routes, API envelope |
| [`docs/ROLES.md`](docs/ROLES.md) | RBAC matrix |
| [`backend/docs/ARCHITECTURE.md`](backend/docs/ARCHITECTURE.md) | Backend route inventory |
| [`docs/README.md`](docs/README.md) | Full documentation index |
