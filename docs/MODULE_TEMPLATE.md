# Module Template

Use this template for every new business module after `customers` becomes the reference module.

## Module Overview

```md
# Module: Products

Purpose:
Fields:
Business rules:
Endpoints:
UI screens:
Permissions:
Audit events:
```

## Backend Checklist

- Migration created.
- SQL queries named `ListX`, `GetX`, `CreateX`, `UpdateX`, `SoftDeleteX`.
- Domain model created.
- DTOs created.
- Repository interface created.
- sqlc repository adapter created.
- Service created with tenant isolation and business rules.
- Handler created with validation and standard responses.
- Routes registered.
- Tests cover service, handler, and tenant isolation.

## Frontend Checklist

- Entity types created.
- API client functions created.
- List feature created.
- Create/edit features created.
- Delete feature created if applicable.
- Widget/page created.
- Sidebar route added.
- Loading, empty, error, and success states implemented.

## API Checklist

- Response format follows `CONVENTIONS.md`.
- Error codes documented.
- Pagination documented when list endpoints exist.
- OpenAPI updated once contract tooling is introduced.

## Definition Of Done

- Backend compiles.
- Frontend compiles.
- Migrations run.
- Tests pass.
- Tenant isolation verified.
- Audit log added for write endpoints.
- Docs updated if behavior changed.
