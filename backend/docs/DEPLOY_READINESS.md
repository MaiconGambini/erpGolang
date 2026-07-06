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

## Migrations

- Atlas: `atlas.hcl`, files in `backend/migrations/`
- `Dockerfile.migrate` / Fly release command runs migrate before app

## Smoke Test Matrix (post-deploy)

1. `GET /healthz`, `GET /readyz`
2. `POST /auth/login` → 200
3. `GET /dashboard/summary` with Bearer token
4. `GET /customers?limit=1`

## Known Gaps

- No route-level RBAC beyond authentication
- `users` CRUD not exposed
- No OpenAPI contract published
- Fly deploy does not include frontend (see `docs/CI_CD.md` Profile B)

See also: `docs/RELIABILITY.md`, `DEPLOYMENT.md`, `docs/CI_CD.md`.
