# Roles

## Personas

- Tenant Admin: owns tenant configuration and user management.
- Manager: manages daily business operations but cannot change critical configuration.
- Operator: performs routine operational work.
- Viewer: read-only access.

## Permission Matrix

MVP note: routes enforce authentication only (`requireAuth`). Role checks are deferred; the matrix documents target behavior.

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
| View financial data | Yes | Yes | No | No |
| Change tenant settings | Yes | No | No | No |
