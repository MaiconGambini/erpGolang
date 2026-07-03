# Product

## Purpose

goERP is a modular, tenant-aware ERP for small and medium businesses. MVP 1 delivers authentication, multi-tenant isolation, and a complete customers module as the reference pattern for future modules (products, orders, etc.).

## Core User Journeys

1. **Admin login** — tenant slug + email + password → dashboard
2. **Session persistence** — refresh on page reload without re-login
3. **Customer management** — list, search, paginate, create, edit, soft-delete customers
4. **Tenant isolation** — tenant A never sees tenant B data (404 on cross-tenant access)

## Differentiators

- Strong tenant isolation from day one
- Modular monolith: new business module in ~1–2 days using customers template
- FSD frontend with consistent UX patterns (`docs/UX_PATTERNS.md`)

## Non-Goals (MVP 1)

- Tenant self-registration
- Roles beyond admin on routes (schema supports all roles; MVP applies admin only)
- Products, orders, invoicing
- i18n, dark mode, PWA, offline
- Microservices, GraphQL, real-time
