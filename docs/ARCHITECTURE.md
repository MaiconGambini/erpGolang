# Architecture

> Canonical detail: see also [`ARCHITECTURE.md`](../ARCHITECTURE.md) at repo root and [`backend/docs/ARCHITECTURE.md`](../backend/docs/ARCHITECTURE.md).

## Subsystems

| Subsystem | Path | Responsibility |
|---|---|---|
| API | `backend/cmd/api` | HTTP entry, graceful shutdown |
| App composition | `backend/internal/app` | Module registration, deps wiring |
| Auth | `backend/internal/auth` | Login, refresh, logout, me |
| Users | `backend/internal/users` | Tenant user CRUD |
| Tenants | `backend/internal/tenants` | Current tenant info |
| Customers | `backend/internal/customers` | Reference business CRUD |
| Audit | `backend/internal/audit` | Append-only write log |
| Platform | `backend/internal/platform` | DB, Redis, logger, validation, HTTP server |
| Shared | `backend/internal/shared` | Errors, pagination, auth/tenant context |
| Frontend | `frontend/src` | Vue 3 FSD app |

## Data Flow

```text
Browser → Vue app → axios (/api/v1) → chi router → handler → service → sqlc repo → PostgreSQL
                                              ↓
                                           Redis (cache, rate limit)
```

## Integration Points

- PostgreSQL 16 — source of truth (`backend/schema.sql`, Atlas migrations)
- Redis 7 — rate limit, cache coordination (not session source of truth)
- JWT access (15min) + refresh cookie (HttpOnly, hash in `auth_sessions`)

## Key Decisions

| Date | Decision | Reason |
|---|---|---|
| MVP | Multi-tenant via `tenant_id` | Isolation without subdomain complexity |
| MVP | chi/v5 router | Middleware ecosystem, explicit routing |
| MVP | sqlc + Atlas | Typed queries, declarative migrations |
| MVP | FSD frontend | Scalable slice isolation for ERP modules |
| MVP | Customers as template module | All future modules copy this pattern |
