# Backend Architecture

The backend is a modular monolith. Modules own HTTP routes, business rules, and sqlc data access. Shared packages stay small and infrastructure-neutral.

## Registered Modules

From `cmd/api/main.go`:

| Module | Type | Audit on writes |
|---|---|---|
| auth | Identity | No |
| users | Placeholder (CRUD deferred) | — |
| tenants | Current tenant | No |
| customers | Entity CRUD | Yes |
| dashboard | Read-model aggregate | No |
| products | Entity CRUD | Yes |
| suppliers | Entity CRUD | Yes |
| sales | Workflow | Yes |
| audit | Write-side recorder | N/A |

## Layering

```text
cmd/api → internal/app → modules → gen/db + shared/* + platform/*
```

- `handler.go` — HTTP only (decode, validate, `httpx` responses).
- `service.go` — DTOs, business rules, sqlc calls, transactions.
- sqlc generated queries are the persistence layer.

## Route Inventory

| Prefix | Methods | Auth |
|---|---|---|
| `/auth/login` | POST | Public + rate limit |
| `/auth/refresh`, `/auth/logout` | POST | Public (cookie) |
| `/auth/me` | GET | JWT + tenant |
| `/tenants/current` | GET | JWT + tenant |
| `/customers`, `/products`, `/suppliers` | GET/POST + `/{id}` GET/PATCH/DELETE | JWT + tenant |
| `/products/low-stock` | GET | JWT + tenant |
| `/sales` | GET/POST + `/{id}` GET/PATCH/DELETE | JWT + tenant |
| `/sales/{id}/confirm`, `/cancel` | POST | JWT + tenant |
| `/dashboard/summary` | GET (`?threshold=`) | JWT + tenant |

## Cross-Module Boundaries

- Modules do not import each other's Go packages.
- Sales and dashboard use sqlc queries across tables (customers, products, sales) — intentional MVP orchestration.
- Sales confirm/cancel mutates `products.stock` inside service transactions.

## Related Docs

- `SCHEMA.md` — table ownership
- `TENANT_ISOLATION.md` — isolation policy
- `SALES_TRANSACTIONS.md` — workflow invariants
- `DASHBOARD_AGGREGATES.md` — KPI definitions
- `RATE_LIMITING.md` — Redis login limits
- `AUTH_SESSION.md` — tokens and cookies
