# Glossary

## Tenant

A tenant is an isolated company/account using goERP. Tenant-owned data must include `tenant_id`.

## User

A user is a person who signs in to a tenant. Users have roles such as admin, manager, operator, or viewer. MVP enforces admin only on routes; the schema supports all roles.

## Customer

A customer is an individual or company that buys from a tenant. Customers are tenant-scoped and can be active or inactive. Soft-deleted via `deleted_at`.

## Supplier

A supplier provides goods or services to a tenant. Suppliers are tenant-scoped CRUD entities (MVP 1). Document (CPF/CNPJ) is unique per tenant when present.

## Product

A product is an item sold or managed by a tenant. Products have SKU (unique per tenant), price, and stock quantity. Low-stock alerts use threshold default `5` (configurable on dashboard and low-stock API).

## Sale (Order)

A sale is a tenant-scoped commercial transaction with line items. Status lifecycle: `draft` → `confirmed` → `cancelled`. Draft sales do not affect stock; confirm decrements stock; cancel on a confirmed sale restores stock.

## Invoice

An invoice is a fiscal or billing document. **Out of MVP 1 scope.**

## Payment

A payment records money movement for an order or invoice. **Out of MVP 1 scope.**

## Stock

Stock is the quantity of a product available to a tenant. Mutated on sales confirm (decrement) and cancel (restore). Products CRUD can set stock directly. Schema enforces `stock >= 0`.

## Dashboard Summary

A read-only aggregate of tenant KPIs: active customers, new customers (30 days), draft sales count, low-stock product count.

## Audit Log

An audit log is an immutable record of a write action, including tenant, actor, action, entity, and timestamp. Auth login/logout is not audited in MVP.
