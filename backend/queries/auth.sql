-- name: GetUserByEmail :one
SELECT u.id, u.tenant_id, u.email, u.password_hash, u.name, u.role, u.active
FROM users u
WHERE u.email = $1 AND u.active = true;

-- name: GetUserByID :one
SELECT u.id, u.tenant_id, u.email, u.name, u.role, u.active
FROM users u
WHERE u.id = $1 AND u.tenant_id = $2 AND u.active = true;

-- name: CreateAuthSession :one
INSERT INTO auth_sessions (tenant_id, user_id, refresh_token_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, tenant_id, user_id, refresh_token_hash, expires_at, revoked_at, rotated_at, created_at;

-- name: GetAuthSessionByHash :one
SELECT id, tenant_id, user_id, refresh_token_hash, expires_at, revoked_at, rotated_at, created_at
FROM auth_sessions
WHERE refresh_token_hash = $1 AND revoked_at IS NULL AND expires_at > now();

-- name: RevokeAuthSession :exec
UPDATE auth_sessions SET revoked_at = now() WHERE id = $1;

-- name: RevokeAuthSessionsByUser :exec
UPDATE auth_sessions SET revoked_at = now() WHERE user_id = $1 AND revoked_at IS NULL;
