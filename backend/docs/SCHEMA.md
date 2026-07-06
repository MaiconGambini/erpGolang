# Schema

Source of truth: `backend/schema.sql`. Migrations in `backend/migrations/`.

## Table Ownership

| Module | Tables |
|---|---|
| tenants | `tenants` |
| users | `users` |
| auth | `auth_sessions` |
| customers | `customers` |
| products | `products` |
| suppliers | `suppliers` |
| sales | `sales`, `sale_items` |
| audit | `audit_logs` |

## Relationships

- `sales.customer_id` → `customers`
- `sale_items.sale_id` → `sales` (ON DELETE CASCADE on hard delete; sales use soft delete)
- `sale_items.product_id` → `products`
- All tenant-owned tables include `tenant_id`

## Constraints (high level)

- `sales.status`: `draft`, `confirmed`, `cancelled`
- `products.stock >= 0`, `products.price > 0`
- Partial unique indexes: `(tenant_id, sku)` on products, `(tenant_id, document)` on customers/suppliers when document present
- Soft delete: `deleted_at` on customers, products, suppliers, sales

## Inventory

- Stock mutated on sales confirm (decrement) and cancel on confirmed sale (increment).
- Products CRUD can set stock directly.
- Low-stock threshold default `5` — see `shared/inventory` and dashboard aggregate docs.

## Migrations (post-MVP modules)

- `202601020001_products.sql`
- `202601030001_suppliers.sql`
- `202601040001_sales.sql`

Dashboard has no tables; it reads aggregates from customers, sales, and products.
