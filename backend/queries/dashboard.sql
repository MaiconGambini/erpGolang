-- Aggregate read model: depends on customers, sales, and products tables.
-- name: GetDashboardSummary :one
SELECT
  (SELECT count(*)::bigint FROM customers c
   WHERE c.tenant_id = @tenant_id AND c.deleted_at IS NULL AND c.active = true
  ) AS active_customers,
  (SELECT count(*)::bigint FROM customers c
   WHERE c.tenant_id = @tenant_id AND c.deleted_at IS NULL
     AND c.created_at >= now() - interval '30 days'
  ) AS new_customers_30d,
  (SELECT count(*)::bigint FROM sales s
   JOIN customers c ON c.id = s.customer_id AND c.tenant_id = s.tenant_id AND c.deleted_at IS NULL
   WHERE s.tenant_id = @tenant_id AND s.deleted_at IS NULL AND s.status = 'draft'
  ) AS draft_sales,
  (SELECT count(*)::bigint FROM products p
   WHERE p.tenant_id = @tenant_id AND p.deleted_at IS NULL AND p.active = true
     AND p.stock <= @threshold
  ) AS low_stock_alerts,
  (SELECT count(*)::bigint FROM sales s
   WHERE s.tenant_id = @tenant_id AND s.deleted_at IS NULL AND s.status = 'confirmed'
  ) AS confirmed_sales_count,
  (SELECT coalesce(sum(s.total), 0)::numeric FROM sales s
   WHERE s.tenant_id = @tenant_id AND s.deleted_at IS NULL AND s.status = 'confirmed'
  ) AS confirmed_sales_total;
