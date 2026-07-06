# Dashboard Aggregates

## Endpoint

`GET /api/v1/dashboard/summary?threshold=N`

- `threshold` optional; default `5`, max `1_000_000` via `inventory.NormalizeLowStockThreshold`

## Fields

| Field | Source | Filters |
|---|---|---|
| `active_customers` | `customers` | tenant, active, not deleted |
| `new_customers_30d` | `customers` | `created_at >= now() - 30 days` |
| `draft_sales` | `sales` JOIN `customers` | `status = draft`, not deleted |
| `low_stock_alerts` | `products` | active, `stock <= threshold` |

SQL: `backend/queries/dashboard.sql`

## Consistency

- Point-in-time snapshot across subqueries is not a single MVCC snapshot.
- Acceptable for MVP dashboard; not for financial reporting.
- No Redis cache — always reads PostgreSQL.

## Scaling Path (when needed)

1. Partial indexes: `sales(tenant_id, status) WHERE deleted_at IS NULL`
2. Materialized view or rollup table per tenant
3. Short TTL cache with explicit stale-read tolerance

## Frontend

- Composable: `use-dashboard-summary.ts`
- Invalidated by customer/product/sale mutations via `invalidateDashboardSummary()`
