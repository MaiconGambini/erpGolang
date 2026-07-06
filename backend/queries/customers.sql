-- name: ListCustomers :many
SELECT id, tenant_id, name, document, document_type, email, phone, postal_code, street, street_number, city, state, active, deleted_at, created_at, updated_at
FROM customers
WHERE tenant_id = @tenant_id AND deleted_at IS NULL
  AND (@search::text = '' OR name ILIKE '%' || @search || '%' OR COALESCE(document, '') ILIKE '%' || @search || '%')
  AND (sqlc.narg('active')::boolean IS NULL OR active = sqlc.narg('active'))
ORDER BY created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: CountCustomers :one
SELECT count(*)::bigint
FROM customers
WHERE tenant_id = @tenant_id AND deleted_at IS NULL
  AND (@search::text = '' OR name ILIKE '%' || @search || '%' OR COALESCE(document, '') ILIKE '%' || @search || '%')
  AND (sqlc.narg('active')::boolean IS NULL OR active = sqlc.narg('active'));

-- name: GetCustomer :one
SELECT id, tenant_id, name, document, document_type, email, phone, postal_code, street, street_number, city, state, active, deleted_at, created_at, updated_at
FROM customers
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL;

-- name: CreateCustomer :one
INSERT INTO customers (tenant_id, name, document, document_type, email, phone, postal_code, street, street_number, city, state, active)
VALUES (@tenant_id, @name, @document, @document_type, @email, @phone, @postal_code, @street, @street_number, @city, @state, @active)
RETURNING id, tenant_id, name, document, document_type, email, phone, postal_code, street, street_number, city, state, active, deleted_at, created_at, updated_at;

-- name: UpdateCustomer :one
UPDATE customers
SET name = @name,
    document = @document,
    document_type = @document_type,
    email = @email,
    phone = @phone,
    postal_code = @postal_code,
    street = @street,
    street_number = @street_number,
    city = @city,
    state = @state,
    active = @active,
    updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL
RETURNING id, tenant_id, name, document, document_type, email, phone, postal_code, street, street_number, city, state, active, deleted_at, created_at, updated_at;

-- name: SoftDeleteCustomer :one
UPDATE customers
SET deleted_at = now(), updated_at = now()
WHERE id = @id AND tenant_id = @tenant_id AND deleted_at IS NULL
RETURNING id, tenant_id, name, document, document_type, email, phone, postal_code, street, street_number, city, state, active, deleted_at, created_at, updated_at;
