# Tenant Isolation

## Enforcement Layers

1. JWT claim → `tenantctx` (never trust client-supplied `tenant_id`)
2. `TenantScope` middleware on protected routes
3. sqlc: every tenant-owned query includes `tenant_id` in `WHERE`
4. Integration tests per module (`//go:build integration`)

## Response Policy

- Cross-tenant resource access → `404 NOT_FOUND` (not 403)
- Missing/invalid token → `401 UNAUTHORIZED`

## Per-Module Pattern

```sql
-- Example: all queries use @tenant_id parameter
WHERE tenant_id = @tenant_id AND deleted_at IS NULL
```

Sales joins include `c.tenant_id = s.tenant_id`. Dashboard subqueries filter `tenant_id` on every table.

## Workflow Isolation

- `DecrementProductStock` includes `tenant_id`; insufficient stock → 409, not cross-tenant leak
- Confirm/cancel operate only on sales matching tenant from JWT

## Test Matrix

| Module | Integration test |
|---|---|
| customers | `internal/customers/integration_test.go` |
| dashboard | `internal/dashboard/integration_test.go` |
| products, suppliers, sales | Handler tests; add integration tests when extending |

Run: `go test -tags=integration ./internal/customers/... ./internal/dashboard/...` with `DATABASE_URL` and seed data.

## Known DB Gaps (MVP)

- FKs on `sales.customer_id` / `sale_items.product_id` do not include composite `tenant_id` — application validates parents with tenant filter
- `users.email` is globally unique
