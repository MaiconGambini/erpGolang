# Backend Architecture

The backend is a modular monolith. Modules own HTTP routes, business rules, and sqlc data access. Shared packages stay small and infrastructure-neutral.

## Registered Modules

From `cmd/api/main.go`:

| Module | Type | RBAC | Audit on writes |
|---|---|---|---|
| auth | Identity | Public / JWT | No |
| users | Admin | Admin only | No |
| tenants | Identity | JWT + tenant | No |
| customers | Entity CRUD | Read/write/delete roles | Yes |
| products | Entity CRUD | Read/write/delete roles | Yes |
| suppliers | Entity CRUD | Read/write/delete roles | Yes |
| sales | Workflow | Read/write/cancel roles | Yes |
| dashboard | Read-model | All roles; financial KPIs redacted | No |
| reports | Read-model | Financial: admin/manager; sale PDF: all read roles | No |
| audit | Cross-cutting | Admin list; write via recorder | N/A |

## Layering

```text
cmd/api → internal/app → modules → gen/db + shared/* + platform/*
```

- `handler.go` — HTTP only (decode, validate, `httpx` responses).
- `service.go` — DTOs, business rules, sqlc calls, transactions.
- `middleware/role.go` — `RequireRole` RBAC guards on routes.
- sqlc generated queries are the persistence layer.

## Route Inventory

| Prefix | Methods | Auth / RBAC |
|---|---|---|
| `/auth/login` | POST | Public + rate limit |
| `/auth/refresh`, `/auth/logout` | POST | Public (cookie) |
| `/auth/me` | GET | JWT + tenant |
| `/tenants/current` | GET | JWT + tenant |
| `/users` | GET; `/{id}` GET/PATCH | JWT + tenant + **admin** |
| `/audit-logs` | GET | JWT + tenant + **admin** |
| `/customers`, `/products`, `/suppliers` | GET/POST + `/{id}` GET/PATCH/DELETE | JWT + tenant + role matrix |
| `/customers`, etc. | `?format=csv` | Same as list |
| `/products/low-stock` | GET | JWT + tenant + read roles |
| `/sales` | GET/POST + `/{id}` GET/PATCH/DELETE | JWT + tenant + role matrix |
| `/sales` | `?format=csv&from=&to=` | Same as list |
| `/sales/{id}/confirm`, `/cancel` | POST | Write / manager roles |
| `/dashboard/summary` | GET (`?threshold=`) | JWT + tenant + read roles |
| `/reports/sales-by-day`, `/top-products` | GET (`?from=&to=`) | JWT + tenant + **admin/manager** |
| `/reports/sales-summary.pdf` | GET | JWT + tenant + **admin/manager** |
| `/reports/sales/{id}/pdf` | GET | JWT + tenant + read roles |

## Cross-Module Boundaries

- Modules do not import each other's Go packages.
- `reports` owns all PDF and chart endpoints (including sale PDF).
- Sales and dashboard use sqlc across tables — intentional read-model / workflow orchestration.
- Sales confirm/cancel mutates `products.stock` inside service transactions.

## Related Docs

- `SCHEMA.md` — table ownership
- `TENANT_ISOLATION.md` — isolation policy
- `SALES_TRANSACTIONS.md` — workflow invariants
- `DASHBOARD_AGGREGATES.md` — KPI definitions
- `RATE_LIMITING.md` — Redis login limits
- `AUTH_SESSION.md` — tokens and cookies
- `../../docs/ROLES.md` — permission matrix
