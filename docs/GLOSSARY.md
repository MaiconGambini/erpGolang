# Glossary

## Tenant

A tenant is an isolated company/account using goERP. Tenant-owned data must include `tenant_id`.

## User

A user is a person who signs in to a tenant. Users have roles: admin, manager, operator, or viewer. RBAC is enforced on API routes (`RequireRole` middleware) and in the UI (`shared/lib/roles.ts`).

## Customer

A customer is an individual or company that buys from a tenant. Customers are tenant-scoped and can be active or inactive. Soft-deleted via `deleted_at`.

## Supplier

A supplier provides goods or services to a tenant. Suppliers are tenant-scoped CRUD entities. Document (CPF/CNPJ) is unique per tenant when present.

## Product

A product is an item sold or managed by a tenant. Products have SKU (unique per tenant), price, and stock quantity. Low-stock alerts use threshold default `5` (configurable on dashboard and low-stock API).

## Sale (Order)

A sale is a tenant-scoped commercial transaction with line items. Status lifecycle: `draft` → `confirmed` → `cancelled`. Draft sales do not affect stock; confirm decrements stock; cancel on a confirmed sale restores stock.

## Report

A read-only aggregation or export of business data. Financial reports (sales-by-day, top products, period summary PDF) are limited to admin and manager roles per `docs/ROLES.md`. CSV list export (`?format=csv`) and single-sale PDF are available per the roles matrix. Not fiscal NF-e.

## Export

Export means downloading tenant-scoped data outside the app UI. CSV exports use UTF-8 BOM for Excel (pt-BR). PDF exports use gofpdf server-side.

## Invoice

An invoice is a fiscal or billing document. **Out of scope.**

## Payment

A payment records money movement for an order or invoice. **Out of scope.**

## Stock

Stock is the quantity of a product available to a tenant. Mutated on sales confirm (decrement) and cancel (restore). Products CRUD can set stock directly. Schema enforces `stock >= 0`.

## Dashboard Summary

A read-only aggregate of tenant KPIs: active customers, new customers (30 days), draft sales count, low-stock product count, confirmed sales count (all-time), and confirmed revenue total (all-time). Operators and viewers see operational KPIs only; financial totals are redacted. Dashboard charts and period PDF use a selectable date range (default last 30 days).

## Audit Log

An audit log is an immutable record of a write action, including tenant, actor, action, entity, and timestamp. Admins can list audit logs via `/api/v1/audit-logs`. Auth login/logout is not audited in MVP.

## RBAC

Role-based access control. See `docs/ROLES.md` for the permission matrix. Viewer is read-only; delete is admin-only on catalog entities.
