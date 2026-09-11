-- name: CreateUser :one
INSERT INTO users (tenant_id, email, password_hash, name, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, tenant_id, email, password_hash, name, role, active, created_at, updated_at;

-- name: ListUsers :many
SELECT id, tenant_id, email, name, role, active, created_at, updated_at
FROM users
WHERE tenant_id = @tenant_id
ORDER BY name ASC
LIMIT @limit_count OFFSET @offset_count;

-- name: CountUsers :one
SELECT count(*)::bigint FROM users WHERE tenant_id = @tenant_id;

-- name: GetUser :one
SELECT id, tenant_id, email, name, role, active, created_at, updated_at
FROM users
WHERE id = @id AND tenant_id = @tenant_id;

-- name: UpdateUser :one
UPDATE users
SET name = @name, role = @role, active = @active, updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id
RETURNING id, tenant_id, email, name, role, active, created_at, updated_at;
