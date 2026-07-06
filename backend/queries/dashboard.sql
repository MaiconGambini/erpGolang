-- Aggregate read model: depends on customers, sales, and products tables.
-- name: GetDashboardSummary :one
SELECT
  (SELECT count(*)::bigint FROM customers
   WHERE tenant_id = @tenant_id AND deleted_at IS NULL AND active = true
  ) AS active_customers,
  (SELECT count(*)::bigint FROM customers
   WHERE tenant_id = @tenant_id AND deleted_at IS NULL
     AND created_at >= now() - interval '30 days'
  ) AS new_customers_30d,
  (SELECT count(*)::bigint FROM sales s
   JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id AND c.deleted_at IS NULL
   WHERE s.tenant_id = @tenant_id AND s.deleted_at IS NULL AND s.status = 'draft'
  ) AS draft_sales,
  (SELECT count(*)::bigint FROM products
   WHERE tenant_id = @tenant_id AND deleted_at IS NULL AND active = true
     AND stock <= @threshold
  ) AS low_stock_alerts;
