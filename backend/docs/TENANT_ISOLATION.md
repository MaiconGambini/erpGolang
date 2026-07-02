# Tenant Isolation

- Tenant-owned tables require `tenant_id`.
- Tenant-owned queries must filter by `tenant_id`.
- Tenant ID comes from authenticated context.
- Cross-tenant access must be tested.
