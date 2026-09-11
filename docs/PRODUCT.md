# Product

## Purpose

goERP is a modular, tenant-aware ERP for small and medium businesses. MVP 1 delivers authentication, multi-tenant isolation, reference CRUD modules (customers, products, suppliers, sales), RBAC, reporting (CSV/PDF/charts), users admin, audit viewer, and a live dashboard summary.

## Core User Journeys

1. **Admin login** — tenant slug + email + password → dashboard with live KPIs and charts
2. **Session persistence** — refresh on page reload without re-login
3. **Customer management** — list, search, paginate, create, edit, soft-delete customers
4. **Catalog & supply** — products (stock, low-stock alerts), suppliers CRUD
5. **Sales** — draft sales, confirm (stock decrement), cancel (stock restore), PDF export
6. **Reporting** — CSV export on lists, sales date filters, dashboard charts, period PDF summary
7. **User admin** — list users, edit role/active (admin only)
8. **Audit viewer** — read-only audit log list (admin only)
9. **Viewer role** — read-only access; create/edit/delete hidden and blocked server-side
10. **Tenant isolation** — tenant A never sees tenant B data (404 on cross-tenant access)

## Differentiators

- Strong tenant isolation from day one
- RBAC enforced per `docs/ROLES.md` on API and UI
- Modular monolith: new business module in ~1–2 days using customers template
- FSD frontend with consistent UX patterns (`docs/UX_PATTERNS.md`)
- Portfolio-ready reporting without claiming fiscal NF-e scope

## MVP 1 Module List

| Module | Status |
|---|---|
| auth | Implemented |
| customers | Implemented |
| products | Implemented |
| suppliers | Implemented |
| sales | Implemented |
| dashboard | Implemented |
| reports | Implemented (CSV via lists, charts, PDF) |
| users | Implemented (list/get/update, admin) |
| audit | Implemented (write recorder + admin list) |

See `docs/ARCHITECTURE.md`, `docs/CI_CD.md`, and `DEPLOYMENT.md` for technical detail.

## Non-Goals (current scope)

- Tenant self-registration
- NF-e, SPED, fiscal invoicing
- Payments, AP/AR, purchase orders, stock ledger
- i18n, PWA, offline (dark mode added 2026-08-25 by operator decision)
- Microservices, GraphQL, real-time
- Live hosted demo (screenshots in README)
