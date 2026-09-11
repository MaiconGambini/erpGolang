-- name: SalesByDay :many
SELECT date_trunc('day', s.created_at)::date AS sale_date,
       count(*)::bigint AS sale_count,
       coalesce(sum(s.total), 0)::numeric(12, 2) AS sale_total
FROM sales s
WHERE s.tenant_id = @tenant_id
  AND s.deleted_at IS NULL
  AND s.status = 'confirmed'
  AND s.created_at >= @from_date
  AND s.created_at < @to_date
GROUP BY 1
ORDER BY 1 ASC;

-- name: TopProducts :many
SELECT p.name AS product_name,
       p.sku AS product_sku,
       sum(si.quantity)::bigint AS quantity_sold,
       coalesce(sum(si.line_total), 0)::numeric(12, 2) AS revenue
FROM sale_items si
JOIN sales s ON s.id = si.sale_id AND s.tenant_id = si.tenant_id
JOIN products p ON p.id = si.product_id AND p.tenant_id = si.tenant_id
WHERE si.tenant_id = @tenant_id
  AND s.deleted_at IS NULL
  AND s.status = 'confirmed'
  AND s.created_at >= @from_date
  AND s.created_at < @to_date
GROUP BY p.id, p.name, p.sku
ORDER BY revenue DESC
LIMIT @limit_count;

-- name: ListConfirmedSalesForReport :many
SELECT s.id, c.name AS customer_name, s.status, s.total, s.created_at
FROM sales s
JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id AND c.deleted_at IS NULL
WHERE s.tenant_id = @tenant_id
  AND s.deleted_at IS NULL
  AND s.status = 'confirmed'
  AND s.created_at >= @from_date
  AND s.created_at < @to_date
ORDER BY s.created_at DESC
LIMIT 500;
