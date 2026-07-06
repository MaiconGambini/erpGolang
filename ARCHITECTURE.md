# Architecture

goERP is a modular monolith with a REST API and a Vue frontend. The architecture optimizes for fast MVP delivery, strong tenant isolation, and a clear path to deployment without microservices.

## System Shape

```text
Browser
  -> Frontend app
  -> REST API (/api/v1)
  -> Go backend
  -> PostgreSQL 16

Go backend
  -> Redis 7
```

## Backend Architecture

The backend is a modular monolith. Each module owns handlers, service, and sqlc queries. DTOs live in `service.go`. sqlc is the persistence layer (no separate repository interface files in MVP).

MVP modules:

- `tenants`: tenant identity, slug, status.
- `users`: users, roles, tenant membership (CRUD routes deferred).
- `auth`: login, refresh, logout, access token verification.
- `customers`: reference CRUD module.
- `products`, `suppliers`: catalog CRUD modules.
- `sales`: draft/confirm/cancel workflow with stock effects.
- `dashboard`: tenant KPI aggregate read model.
- `audit`: append-only write event log (no HTTP routes).

Dependency direction:

```text
cmd/api -> internal/app -> modules -> shared/platform
```

Rules:

- Domain files do not import HTTP, SQL, logging, or framework packages.
- Handlers translate HTTP and call services.
- Services enforce business rules and tenant isolation.
- Services map sqlc rows to DTOs in `service.go`.
- Modules do not import each other's packages; cross-table reads/writes happen via sqlc inside service transactions (sales, dashboard).
- Shared packages contain infrastructure-neutral helpers only.

## Frontend Architecture

The frontend follows Feature-Sliced Design:

```text
app -> processes -> pages -> widgets -> features -> entities -> shared
```

Rules:

- Imports only flow downward.
- `entities` mirror API DTOs and basic API calls.
- `features` own user actions such as create customer or login.
- `widgets` compose features into reusable page sections.
- `pages` compose widgets and features; they should contain little business logic.
- `shared` contains API client, UI primitives, formatting, and config.

## Tenant Isolation

Every tenant-owned table has `tenant_id`. Every tenant-owned query filters by `tenant_id`. Tenant IDs come from authenticated context, never from request bodies.

Cross-tenant access must return `404` for tenant-owned resources unless an explicit system-admin capability is introduced later.

## Persistence

PostgreSQL is the source of truth. Redis is allowed for cache, rate limiting, session coordination, and temporary state, but not as the only store for domain records.

## Deployment Path

Phase 1: local machine with Docker Compose for PostgreSQL and Redis.

Phase 2: VPS with Docker Compose and Caddy.

Phase 3: AWS EC2 + RDS. ECS/Fargate is deferred until operational needs justify it.
