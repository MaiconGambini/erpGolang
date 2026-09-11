# Dashboard Aggregates

## Endpoint

`GET /api/v1/dashboard/summary?threshold=N`

- `threshold` optional; default `5`, max `1_000_000` via `inventory.NormalizeLowStockThreshold`
- **Not date-scoped** — KPIs use fixed windows or all-time counts (see table below)
- Operators and viewers receive `confirmedSalesCount = 0` and `confirmedSalesTotal = "0"` (financial redaction in handler)

## Fields

| JSON field | SQL source | Filters |
|---|---|---|
| `activeCustomers` | `customers` | tenant, active, not deleted |
| `newCustomers30d` | `customers` | `created_at >= now() - 30 days` |
| `draftSales` | `sales` JOIN `customers` | `status = draft`, not deleted |
| `lowStockAlerts` | `products` | active, `stock <= threshold` |
| `confirmedSalesCount` | `sales` | `status = confirmed`, all-time; **admin/manager only** |
| `confirmedSalesTotal` | `sales` | sum of confirmed totals, all-time; **admin/manager only** |

SQL: `backend/queries/dashboard.sql`

## Charts vs KPIs

| Surface | Date range | Module |
|---|---|---|
| KPI cards on `/` | Fixed / all-time (above) | `dashboard` |
| Sales-by-day chart | User-selected (`from`/`to`, default 30d) | `reports` |
| Top products chart | User-selected | `reports` |
| Period PDF export | User-selected | `reports` |

Do not assume KPI cards and charts share the same time window.

## Consistency

- Point-in-time snapshot across subqueries is not a single MVCC snapshot.
- Acceptable for MVP dashboard; not for fiscal reporting.
- No Redis cache — always reads PostgreSQL.

## Scaling Path (when needed)

1. Partial indexes: `sales(tenant_id, status, created_at) WHERE deleted_at IS NULL`
2. Materialized view or rollup table per tenant
3. Short TTL cache with explicit stale-read tolerance

## Frontend

- Composable: `features/dashboard/summary/model/use-dashboard-summary.ts`
- Invalidated by customer/product/sale mutations via `invalidateDashboardSummary()`
- Financial charts and PDF hidden for operator/viewer (`canViewFinancial` in `shared/lib/roles.ts`)
- KPI drill-downs: active customers → `/customers`; draft → `/sales?status=draft`; low stock → `/products?lowStock=1`; confirmed → `/sales?status=confirmed`
