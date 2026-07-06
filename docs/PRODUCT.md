# Product

## Purpose

goERP is a modular, tenant-aware ERP for small and medium businesses. MVP 1 delivers authentication, multi-tenant isolation, reference CRUD modules (customers, products, suppliers, sales), and a live dashboard summary.

## Core User Journeys

1. **Admin login** — tenant slug + email + password → dashboard with live KPIs
2. **Session persistence** — refresh on page reload without re-login
3. **Customer management** — list, search, paginate, create, edit, soft-delete customers
4. **Catalog & supply** — products (stock, low-stock alerts), suppliers CRUD
5. **Sales** — draft sales, confirm (stock decrement), cancel (stock restore)
6. **Tenant isolation** — tenant A never sees tenant B data (404 on cross-tenant access)

## Differentiators

- Strong tenant isolation from day one
- Modular monolith: new business module in ~1–2 days using customers template
- FSD frontend with consistent UX patterns (`docs/UX_PATTERNS.md`)

## MVP 1 Module List

| Module | Status |
|---|---|
| auth | Implemented |
| customers | Implemented |
| products | Implemented |
| suppliers | Implemented |
| sales | Implemented |
| dashboard | Implemented |
| users CRUD | Deferred |

See `docs/ARCHITECTURE.md`, `docs/CI_CD.md`, and `DEPLOYMENT.md` for technical detail.

## Non-Goals (current scope)

- Tenant self-registration
- Roles beyond admin on routes (schema supports all roles; MVP applies admin only)
- Invoicing, purchasing workflows beyond sales draft/confirm
- i18n, dark mode, PWA, offline
- Microservices, GraphQL, real-time
