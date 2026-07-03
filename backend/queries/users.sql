-- name: CreateUser :one
INSERT INTO users (tenant_id, email, password_hash, name, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, tenant_id, email, name, role, active, created_at, updated_at;
