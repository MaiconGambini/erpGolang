# Conventions

## Language

- Code comments and docs: English.
- User-facing UI copy: pt-BR for the Brazilian ERP audience.
- Commit messages: Conventional Commits in English.

## Naming

- Database: `snake_case`.
- Go exported types and methods: `PascalCase`.
- Go package names: short lowercase names.
- JSON fields: `camelCase`.
- TypeScript types/interfaces: `PascalCase`.
- Vue and frontend files: `kebab-case`.
- Feature slices: verb-first, for example `create-customer`, not `customer-create`.

## API Response Shape

Success:

```json
{
  "data": {}
}
```

Error:

```json
{
  "error": {
    "code": "CUSTOMER_NOT_FOUND",
    "message": "Customer not found.",
    "details": {}
  }
}
```

Pagination:

```json
{
  "data": [],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 100
  }
}
```

## SQL Query Names

- `ListX`
- `GetX`
- `CreateX`
- `UpdateX`
- `SoftDeleteX`

## Tenant Rules

- Every tenant-owned table must have `tenant_id NOT NULL`.
- Every tenant-owned query must receive and filter by `tenant_id`.
- Never trust tenant ID from request body or query string.
- Cross-tenant access must be covered by tests.

## Export

- CSV list export: `?format=csv` on paginated list endpoints.
- UTF-8 BOM prefix for Excel compatibility (pt-BR locale).
- PDF/binary exports live under `reports` module; use `Content-Disposition: attachment`.
- Export respects the same RBAC roles as the underlying list/read endpoint.

## Git

Use Conventional Commits:

- `feat:`
- `fix:`
- `docs:`
- `chore:`
- `test:`
- `refactor:`
