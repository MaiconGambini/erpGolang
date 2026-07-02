# Business Rules

## Tenant Isolation

- Tenant-owned data belongs to exactly one tenant.
- Cross-tenant access is forbidden.
- Tenant-owned queries must filter by `tenant_id`.
- Cross-tenant reads should return `404` for tenant-owned resources.

## Users

- Passwords are never stored or logged in plain text.
- Deactivation is preferred over hard deletion.
- MVP may implement only tenant admin, but schema and docs must preserve the target role model.
- Email uniqueness must be decided before full auth implementation. Default recommendation: globally unique email for MVP simplicity.

## Customers

- Customer documents (CPF/CNPJ) are unique per tenant when present.
- A customer may be created without a document only if the tenant workflow allows it.
- Customers are soft deleted by default.
- Customers with linked business records must not be hard deleted.

## Money

- Never use floating point for persisted money.
- Store money with fixed precision or decimal types.

## Audit

- Write operations should be auditable.
- Audit records include tenant, actor, action, entity type, entity ID, and timestamp.
- Do not log passwords, tokens, or secrets.
