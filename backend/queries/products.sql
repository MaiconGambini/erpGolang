-- name: ListProducts :many
SELECT id, tenant_id, name, sku, price, stock, unit, barcode, active, deleted_at, created_at, updated_at
FROM products
WHERE tenant_id = @tenant_id AND deleted_at IS NULL
  AND (@search::text = '' OR name ILIKE '%' || @search || '%' OR sku ILIKE '%' || @search || '%')
  AND (sqlc.narg('active')::boolean IS NULL OR active = sqlc.narg('active'))
ORDER BY created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: CountProducts :one
SELECT count(*)::bigint
FROM products
WHERE tenant_id = @tenant_id AND deleted_at IS NULL
  AND (@search::text = '' OR name ILIKE '%' || @search || '%' OR sku ILIKE '%' || @search || '%')
  AND (sqlc.narg('active')::boolean IS NULL OR active = sqlc.narg('active'));

-- name: ListLowStockProducts :many
SELECT id, tenant_id, name, sku, price, stock, unit, barcode, active, deleted_at, created_at, updated_at
FROM products
WHERE tenant_id = @tenant_id AND deleted_at IS NULL AND active = true
  AND stock <= @threshold
ORDER BY stock ASC, name ASC
LIMIT @limit_count;

-- name: GetProduct :one
SELECT id, tenant_id, name, sku, price, stock, unit, barcode, active, deleted_at, created_at, updated_at
FROM products
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL;

-- name: CreateProduct :one
INSERT INTO products (tenant_id, name, sku, price, stock, unit, barcode, active)
VALUES (@tenant_id, @name, @sku, @price, @stock, @unit, @barcode, @active)
RETURNING id, tenant_id, name, sku, price, stock, unit, barcode, active, deleted_at, created_at, updated_at;

-- name: UpdateProduct :one
UPDATE products
SET name = @name,
    sku = @sku,
    price = @price,
    stock = @stock,
    unit = @unit,
    barcode = @barcode,
    active = @active,
    updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL
RETURNING id, tenant_id, name, sku, price, stock, unit, barcode, active, deleted_at, created_at, updated_at;

-- name: SoftDeleteProduct :one
UPDATE products
SET deleted_at = now(), updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL
RETURNING id, tenant_id, name, sku, price, stock, unit, barcode, active, deleted_at, created_at, updated_at;
