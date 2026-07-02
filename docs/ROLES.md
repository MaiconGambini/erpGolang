# Roles

## Personas

- Tenant Admin: owns tenant configuration and user management.
- Manager: manages daily business operations but cannot change critical configuration.
- Operator: performs routine operational work.
- Viewer: read-only access.

## Permission Matrix

| Capability | Admin | Manager | Operator | Viewer |
|---|---:|---:|---:|---:|
| View customers | Yes | Yes | Yes | Yes |
| Create customers | Yes | Yes | Yes | No |
| Edit customers | Yes | Yes | Yes | No |
| Delete customers | Yes | No | No | No |
| Manage users | Yes | No | No | No |
| View financial data | Yes | Yes | No | No |
| Change tenant settings | Yes | No | No | No |

MVP note: the first implementation may only enforce `admin`, but the data model should not block future roles.
