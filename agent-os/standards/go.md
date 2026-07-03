# Go Standard

- Use Go 1.23+ with explicit error handling; no panic in production paths.
- Keep domain logic free of HTTP, SQL, and framework imports.
- Handlers translate HTTP; services enforce business rules; repositories wrap sqlc.
- Every DB query receives `ctx context.Context` as first parameter.
- Every tenant-owned query filters by `tenant_id` from JWT context, never from request body.
- Use `shopspring/decimal` for money; `google/uuid` for IDs.
- Structured errors via `httpx.AppError`; never leak internals in 500 responses.
- Tests: unit (mocked repo), integration (`//go:build integration` + testcontainers), handler (httptest).
