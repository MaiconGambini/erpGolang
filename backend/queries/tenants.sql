-- name: GetTenantBySlug :one
SELECT id, slug, name, status, created_at, updated_at
FROM tenants
WHERE slug = $1 AND status = 'active';

-- name: GetTenantByID :one
SELECT id, slug, name, status, created_at, updated_at
FROM tenants
WHERE id = $1 AND status = 'active';

-- name: CreateTenant :one
INSERT INTO tenants (slug, name)
VALUES ($1, $2)
RETURNING id, slug, name, status, created_at, updated_at;
