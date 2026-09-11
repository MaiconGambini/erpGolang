-- name: InsertAuditLog :one
INSERT INTO audit_logs (tenant_id, actor_user_id, action, entity_type, entity_id, before, after)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, tenant_id, actor_user_id, action, entity_type, entity_id, before, after, created_at;

-- name: ListAuditLogs :many
SELECT id, tenant_id, actor_user_id, action, entity_type, entity_id, created_at
FROM audit_logs
WHERE tenant_id = @tenant_id
ORDER BY created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: CountAuditLogs :one
SELECT count(*)::bigint FROM audit_logs WHERE tenant_id = @tenant_id;
