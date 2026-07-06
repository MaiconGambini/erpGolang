-- name: ListSales :many
SELECT s.id, s.tenant_id, s.customer_id, c.name AS customer_name, s.status, s.total, s.notes,
       s.deleted_at, s.created_at, s.updated_at
FROM sales s
JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id AND c.deleted_at IS NULL
WHERE s.tenant_id = @tenant_id AND s.deleted_at IS NULL
  AND (@search::text = '' OR c.name ILIKE '%' || @search || '%')
  AND (sqlc.narg('status')::text IS NULL OR s.status = sqlc.narg('status'))
  AND (sqlc.narg('from_date')::timestamptz IS NULL OR s.created_at >= sqlc.narg('from_date'))
  AND (sqlc.narg('to_date')::timestamptz IS NULL OR s.created_at < sqlc.narg('to_date'))
ORDER BY s.created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: CountSales :one
SELECT count(*)::bigint
FROM sales s
JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id AND c.deleted_at IS NULL
WHERE s.tenant_id = @tenant_id AND s.deleted_at IS NULL
  AND (@search::text = '' OR c.name ILIKE '%' || @search || '%')
  AND (sqlc.narg('status')::text IS NULL OR s.status = sqlc.narg('status'))
  AND (sqlc.narg('from_date')::timestamptz IS NULL OR s.created_at >= sqlc.narg('from_date'))
  AND (sqlc.narg('to_date')::timestamptz IS NULL OR s.created_at < sqlc.narg('to_date'));

-- name: GetSale :one
SELECT s.id, s.tenant_id, s.customer_id, c.name AS customer_name, s.status, s.total, s.notes,
       s.deleted_at, s.created_at, s.updated_at
FROM sales s
JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id AND c.deleted_at IS NULL
WHERE s.id = @id AND s.tenant_id = @tenant_id AND s.deleted_at IS NULL;

-- name: CreateSale :one
INSERT INTO sales (tenant_id, customer_id, status, total, notes)
VALUES (@tenant_id, @customer_id, @status, @total, @notes)
RETURNING id, tenant_id, customer_id, status, total, notes, deleted_at, created_at, updated_at;

-- name: UpdateSaleDraft :one
UPDATE sales
SET customer_id = @customer_id, total = @total, notes = @notes, updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND status = 'draft' AND deleted_at IS NULL
RETURNING id, tenant_id, customer_id, status, total, notes, deleted_at, created_at, updated_at;

-- name: CountSalesByCustomer :one
SELECT count(*)::bigint FROM sales WHERE tenant_id = @tenant_id AND customer_id = @customer_id AND deleted_at IS NULL;

-- name: UpdateSaleStatus :one
UPDATE sales
SET status = @status, updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND status = @from_status AND deleted_at IS NULL
RETURNING id, tenant_id, customer_id, status, total, notes, deleted_at, created_at, updated_at;

-- name: SoftDeleteSale :one
UPDATE sales
SET deleted_at = now(), updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND status = 'draft' AND deleted_at IS NULL
RETURNING id, tenant_id, customer_id, status, total, notes, deleted_at, created_at, updated_at;

-- name: ListSaleItems :many
SELECT si.id, si.tenant_id, si.sale_id, si.product_id, p.name AS product_name, p.sku AS product_sku,
       si.quantity, si.unit_price, si.line_total, si.created_at
FROM sale_items si
JOIN products p ON p.id = si.product_id AND p.tenant_id = si.tenant_id AND p.deleted_at IS NULL
WHERE si.sale_id = @sale_id AND si.tenant_id = @tenant_id
ORDER BY si.created_at ASC;

-- name: DeleteSaleItems :exec
DELETE FROM sale_items
WHERE sale_id = @sale_id AND tenant_id = @tenant_id;

-- name: CreateSaleItem :one
INSERT INTO sale_items (tenant_id, sale_id, product_id, quantity, unit_price, line_total)
VALUES (@tenant_id, @sale_id, @product_id, @quantity, @unit_price, @line_total)
RETURNING id, tenant_id, sale_id, product_id, quantity, unit_price, line_total, created_at;

-- name: DecrementProductStock :one
UPDATE products
SET stock = stock - @quantity, updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL AND active = true
  AND stock >= @quantity
RETURNING id, stock;

-- name: IncrementProductStock :exec
UPDATE products
SET stock = stock + @quantity, updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL;
