# Roles

## Personas

- Tenant Admin: owns tenant configuration and user management.
- Manager: manages daily business operations but cannot change critical configuration.
- Operator: performs routine operational work.
- Viewer: read-only access.

## Permission Matrix

Routes enforce `RequireRole` middleware on the API and role helpers in the Vue UI (`shared/lib/roles.ts`).

| Capability | Admin | Manager | Operator | Viewer |
|---|---:|---:|---:|---:|
| View dashboard | Yes | Yes | Yes | Yes |
| View customers | Yes | Yes | Yes | Yes |
| Create customers | Yes | Yes | Yes | No |
| Edit customers | Yes | Yes | Yes | No |
| Delete customers | Yes | No | No | No |
| View products | Yes | Yes | Yes | Yes |
| Create/edit products | Yes | Yes | Yes | No |
| Delete products | Yes | No | No | No |
| View suppliers | Yes | Yes | Yes | Yes |
| Create/edit suppliers | Yes | Yes | Yes | No |
| Delete suppliers | Yes | No | No | No |
| View sales | Yes | Yes | Yes | Yes |
| Create/confirm sales | Yes | Yes | Yes | No |
| Cancel sales | Yes | Yes | No | No |
| Delete draft sales | Yes | Yes | No | No |
| Manage users | Yes | No | No | No |
| View audit log | Yes | No | No | No |
| View financial data (revenue charts, period PDF, aggregate KPIs) | Yes | Yes | No | No |
| Export sale PDF (single order) | Yes | Yes | Yes | Yes |
| Change tenant settings | Yes | No | No | No |
