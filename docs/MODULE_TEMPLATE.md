# Module Template

Use this template for every new business module. **Reference modules:** `customers` (CRUD), `sales` (workflow), `dashboard` (read-model).

## Choose Module Kind

- [ ] **Entity CRUD** — copy: `customers`, `products`, `suppliers`
- [ ] **Workflow** — copy: `sales`
- [ ] **Read-model** — copy: `dashboard`

## Module Overview

```md
# Module: <Name>

Purpose:
Fields:
Business rules:
Endpoints:
UI screens:
Permissions:
Audit events:
```

## Backend Checklist (all kinds)

- Migration created with `tenant_id` on tenant-owned tables.
- SQL queries in `queries/<module>.sql` with explicit `tenant_id` parameters.
- Files: `module.go`, `handler.go`, `service.go`.
- DTOs and `toDTO()` in `service.go`.
- sqlc generated code used directly in service (no separate repository interface unless unit-test mocks are required).
- Routes: `AuthJWT` + `TenantScope` on tenant-owned resources.
- `handler_test.go` for HTTP contract; `integration_test.go` for tenant isolation (`//go:build integration`).

### Entity CRUD extras

- Queries: `ListX`, `GetX`, `CreateX`, `UpdateX`, `SoftDeleteX`.
- `errNotFound` → 404 on cross-tenant miss.
- Audit on create/update/delete via `audit.Recorder`.
- Paginated lists via `httpx.Paginated`.

### Workflow extras

- Document status machine table in spec.
- Use `pool.Begin` + transaction for multi-row effects.
- Action routes: `POST /{id}/<action>`.
- Sentinel errors (`errInvalidStatus`, `errInsufficientStock`, etc.) mapped in handler.
- Integration tests for status conflicts and stock edge cases.

### Read-model extras

- Aggregate SQL with `tenant_id` in every subquery.
- No audit, no mutations.
- Prefer module-owned DTO with camelCase JSON (avoid leaking sqlc row shape).
- Tenant isolation integration test.
- Document Vue Query invalidation sources from other modules.

## Frontend Checklist

- `entities/<module>/model/types.ts` — API DTO types.
- `entities/<module>/api/*.api.ts` — thin axios wrappers.
- Features under `features/<module>/` (list, create, edit, delete, actions).
- Page under `pages/<module>/`.
- Sidebar route added.
- Loading, empty, error, and success states.
- Invalidate related query keys; call `invalidateDashboardSummary()` when KPIs change.

## Cache Invalidation Matrix

| Mutation | Invalidate |
|---|---|
| Customer CRUD | `['customers']`, dashboard summary |
| Product CRUD | `['products']`, dashboard summary |
| Supplier CRUD | `['suppliers']` |
| Sale create/confirm/cancel | `['sales']`, `['products']`, dashboard summary |
| Sale delete (draft) | `['sales']`, dashboard summary |
| Logout | `queryClient.clear()` |

## API Checklist

- Response format follows project conventions (`{ data }`, `{ data, pagination }`).
- Error codes documented.
- Pagination uses `limit` and `offset`.

## Definition Of Done

- Backend compiles; `go test ./...` passes.
- Integration test for tenant isolation (when module is tenant-owned).
- Frontend typecheck and build pass.
- Migrations apply cleanly.
- Audit on write endpoints (except read-models).
- Docs updated: `docs/BUSINESS_RULES.md`, `backend/docs/SCHEMA.md` if schema changed.
