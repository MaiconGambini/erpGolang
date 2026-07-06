# Business Rules

## Tenant Isolation

- Tenant-owned data belongs to exactly one tenant.
- Cross-tenant access is forbidden.
- Tenant-owned queries must filter by `tenant_id`.
- Cross-tenant reads return `404` for tenant-owned resources.

## Users

- Passwords are never stored or logged in plain text.
- Deactivation is preferred over hard deletion.
- List/get/update routes are admin-only; role changes enforced server-side.
- Email is globally unique in schema (MVP simplification).
- User mutations are not yet written to the audit log (known gap).

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

## RBAC

- Route-level enforcement via `RequireRole` middleware per `docs/ROLES.md`.
- Viewer is read-only on all write endpoints.
- Delete on catalog entities (customers, products, suppliers) is admin-only.
- Cancel sales and delete draft sales: admin and manager only.
- Users list/update and audit log list: admin only.
- Financial data (revenue charts, period PDF, `confirmedSalesCount` / `confirmedSalesTotal` KPIs): admin and manager only.
- Single-sale PDF and CSV list export: per roles matrix in `docs/ROLES.md`.

## Dashboard

- KPIs are point-in-time SQL aggregates (no cache).
- `threshold` query param defaults to `5`, capped at 1_000_000.
- Dashboard does not write data or emit audit events.
- Financial KPI fields are zeroed in the API response for operator and viewer roles.

## Money

- Never use floating point for persisted money.
- Store money with fixed precision or decimal types.

## Reporting

- CSV export is read-only; respects list filters and tenant scope.
- Chart and period PDF queries require `from` and `to` date parameters (ISO dates).
- Period PDF caps at 500 confirmed sales per export (see `ListConfirmedSalesForReport`).
- Reports are not fiscal documents (no NF-e).

## Audit

- Write operations on business modules should be auditable.
- Audit records include tenant, actor, action, entity type, entity ID, and timestamp.
- Admins can list audit logs via `GET /audit-logs` (paginated).
- Do not log passwords, tokens, or secrets.
- Audit insert failures are currently best-effort (logged internally, not rolled back).
- User mutations are not audited in MVP (known gap).
