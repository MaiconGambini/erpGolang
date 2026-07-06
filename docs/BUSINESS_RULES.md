# Business Rules

## Tenant Isolation

- Tenant-owned data belongs to exactly one tenant.
- Cross-tenant access is forbidden.
- Tenant-owned queries must filter by `tenant_id`.
- Cross-tenant reads return `404` for tenant-owned resources.

## Users

- Passwords are never stored or logged in plain text.
- Deactivation is preferred over hard deletion.
- MVP enforces admin only on routes; schema preserves the target role model.
- Email is globally unique in schema (MVP simplification).

## Customers

- Customer documents (CPF/CNPJ) are unique per tenant when present.
- Customers are soft deleted by default.
- Customers with linked sales should not be deleted without a business rule (not enforced in MVP).

## Products

- SKU is unique per tenant when present.
- Price must be greater than zero.
- Stock must be greater than or equal to zero.
- Low-stock threshold default is `5` (shared constant in Go and frontend config).
- `GET /products/low-stock` and dashboard KPI use the same threshold semantics.

## Suppliers

- Supplier documents are unique per tenant when present.
- Suppliers are soft deleted by default.
- Mirror customer CRUD rules where applicable.

## Sales

| Status | Allowed operations |
|---|---|
| `draft` | get, update, delete, confirm |
| `confirmed` | get, cancel |
| `cancelled` | get only |

- Only draft sales are editable or deletable.
- Confirm transitions to `confirmed` and decrements stock for each line item.
- Cancel is allowed only from `confirmed` and restores stock.
- Insufficient stock on confirm returns `409 INSUFFICIENT_STOCK`.
- Invalid status transitions return `409 INVALID_STATUS`.
- Unit price is snapshotted from product at create/update time.
- Draft sales count on dashboard includes non-deleted drafts with live customers.

## Dashboard

- KPIs are point-in-time SQL aggregates (no cache).
- `threshold` query param defaults to `5`, capped at 1_000_000.
- Dashboard does not write data or emit audit events.

## Money

- Never use floating point for persisted money.
- Store money with fixed precision or decimal types.

## Audit

- Write operations on business modules should be auditable.
- Audit records include tenant, actor, action, entity type, entity ID, and timestamp.
- Do not log passwords, tokens, or secrets.
- Audit insert failures are currently best-effort (logged internally, not rolled back).
