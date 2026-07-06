# Backend Deploy Readiness

## Environment Variables

| Variable | Required | Used for |
|---|---|---|
| `APP_ENV` | Yes | development vs production behavior |
| `HTTP_ADDR` | Yes | Listen address |
| `DATABASE_URL` | Yes | PostgreSQL |
| `REDIS_URL` | Yes | Rate limit + readyz |
| `JWT_ACCESS_SECRET` | Yes | Access JWT signing |
| `JWT_REFRESH_SECRET` | Yes | Config validation only (unused in refresh flow) |
| `ALLOWED_ORIGINS` | Yes | CORS |
| `JWT_ACCESS_TTL` | No | Default 15m |
| `BCRYPT_COST` | No | Default 12 |

## Health Checks

- `GET /healthz` — liveness
- `GET /readyz` — PostgreSQL + Redis; returns `NOT_READY` on failure

## Security Checklist

- [ ] Refresh cookie `Secure=true` in production
- [ ] `ALLOWED_ORIGINS` set to real frontend origin(s)
- [ ] Secrets not in Git; use `GOERP_ENV_FILE` on VPS
- [ ] PostgreSQL and Redis not exposed to internet (VPS compose)
- [ ] Migrations run before traffic (`cmd/migrate`, Fly `release_command`)
- [ ] RBAC matrix in `docs/ROLES.md` matches deployed routes

## Migrations

- Atlas: `atlas.hcl`, files in `backend/migrations/`
- `Dockerfile.migrate` / Fly release command runs migrate before app

## Smoke Test Matrix (post-deploy)

1. `GET /healthz`, `GET /readyz`
2. `POST /auth/login` → 200
3. `GET /dashboard/summary` with Bearer token
4. `GET /customers?limit=1`
5. `GET /reports/sales-by-day?from=2026-01-01&to=2026-01-31` (admin/manager token)
6. `GET /audit-logs?limit=1` (admin token)

## Implemented (Portfolio 8.5+)

- Route-level RBAC via `RequireRole` on all business routes
- Users list/get/update API (`/users`, admin)
- Audit log list API (`/audit-logs`, admin)
- CSV export on list endpoints; PDF reports under `/reports/*`

## Known Gaps

- No OpenAPI contract published (`contract/` planned)
- User mutations not written to audit log
- No last-admin guard on user deactivation
- Fly deploy does not include frontend (see `docs/CI_CD.md` Profile B)

See also: `docs/RELIABILITY.md`, `DEPLOYMENT.md`, `docs/CI_CD.md`, `docs/ROLES.md`.
