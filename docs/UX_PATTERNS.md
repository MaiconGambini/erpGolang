# UX Patterns

## CRUD Page Layout

```text
Header: title, description, primary action
Toolbar: search, filters, export
Content: table or responsive card list
Footer: pagination
```

## Dashboard / Summary Page

### Layout

```text
Header: title, short description (no primary action required)
Metrics: responsive KPI grid (1 col mobile → 2–4 cols desktop)
```

- Use `AppShell` + `AppPageHeader` like CRUD pages.
- Metric grid: `repeat(auto-fit, minmax(220px, 1fr))` with `gap: 16px`.

### Metric Card

```text
┌─────────────────────────┐
│ Label (text-muted)      │
│ 1.234 (text-primary)    │
└─────────────────────────┘
```

- Values use `Intl.NumberFormat('pt-BR')`.
- Zero is valid — show `0`, not em dash.
- Per-card skeletons during load (not page-level “Carregando…”).
- Page-level error with retry when summary query fails.
- “Alertas” (`low_stock_alerts`): warning accent when `> 0` (future: link to `/products?low_stock=true`).

### Tokens

| Property | Token |
|---|---|
| Background | `--color-surface` |
| Border | `1px solid var(--color-border)` |
| Radius | `--radius-lg` |
| Label | `--color-text-muted` |
| Value | `--color-text-primary` |
| Warning metric | `--color-warning` |

## Async Query States (Vue Query)

| State | Target pattern |
|---|---|
| Loading | Skeleton for layout-stable UI |
| Error | Shared message + retry (`refetch`) |
| Empty | “Nenhum X encontrado” for lists; dashboard zeros are normal |
| Success | Formatted data |

- Gate queries with `enabled` when auth is required.
- Invalidate dashboard summary after customer/product/sale mutations.

## Sales Flow

- Draft sales: create dialog, edit (backend supports PATCH; UI may defer).
- Confirm: destructive-adjacent action with confirmation; show status badge after success.
- Cancel: only on confirmed sales; restores stock server-side.
- Stock warnings: highlight low-stock rows in product table (`LOW_STOCK_THRESHOLD`).

## Loading

- Use skeletons for table/page content.
- Use button-level loading for submit actions.
- Use spinner only for blocking operations.

## Empty State

- Short title.
- Helpful explanation.
- CTA when the user has permission.

## Form Validation

- Field errors appear inline below fields.
- Global failures can use toast plus inline summary.
- Toast must not be the only feedback for critical errors.

## Delete Confirmation

- Use a confirmation dialog.
- Include the entity name.
- Use destructive styling only on the confirm action.

## Create/Edit Success

- Show success toast.
- Close dialog.
- Invalidate/refetch the list query and dashboard summary when applicable.

## Pagination

- Server-side pagination.
- Default page size: 20.
- API uses `limit` and `offset`.

## Search

- Debounce search by 300ms.
- Show inline loading while results update.

## Accessibility

- Metric grid: `role="region"` + `aria-label="Resumo da operação"`.
- Error messages: `role="alert"`; retry button keyboard-focusable.
- Do not rely on color alone for alert metrics.

## Role-Aware UI

- Central helpers: `shared/lib/roles.ts` (`canCreate`, `canEdit`, `canDelete`, `canViewFinancial`, `requireAdmin`).
- Viewer: hide create/edit/delete actions; server still enforces RBAC.
- Admin-only routes: `/users`, `/audit` with `requireAdmin` router guard.
- Financial surfaces (revenue KPIs, charts, period PDF) hidden unless `canViewFinancial` (admin/manager).
- Sale PDF and CSV export remain available per `docs/ROLES.md`.

## Reporting & Export

### CSV export

- Toolbar button on list pages (customers, products, suppliers, sales).
- API: `?format=csv` on list endpoints; UTF-8 BOM for Excel (pt-BR).
- Sales export passes current table filters (`from`, `to`, `status`, search) via exposed `exportParams`.

### PDF export

- Sale detail: row action or sales table → `GET /reports/sales/{id}/pdf`.
- Period summary: dashboard toolbar → `GET /reports/sales-summary.pdf?from=&to=`.
- Download via blob helper in `shared/lib/export.ts`; show toast on failure.

### Dashboard charts

- Date range picker (default last 30 days) drives chart queries and period PDF.
- Sales-by-day line chart and top-products bar chart (Chart.js).
- KPI drill-downs: customers → `/customers`; draft sales → `/sales?status=draft`; low stock → `/products?lowStock=1`; confirmed → `/sales?status=confirmed`.
- Six KPI cards; operators/viewers see operational metrics only (financial totals redacted server-side).
