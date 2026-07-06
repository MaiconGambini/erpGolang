# Reliability

## Startup And Verification

```bash
# Infra
docker compose -f docker-compose.dev.yml up -d

# Migrations + seed
cd backend && go run ./cmd/migrate
go run ./cmd/seed   # tenants acme/beta, password admin123

# Backend health (when running)
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz

# Baseline gate
make validate
```

| Layer | Command | Pass |
|---|---|---|
| Compose | `docker compose -f docker-compose.dev.yml config` | exit 0 |
| Backend | `cd backend && go test ./... && go build ./...` | exit 0 |
| Backend integration | `go test -tags=integration ./internal/customers/... ./internal/dashboard/...` | exit 0 (needs `DATABASE_URL` + seed) |
| Frontend | `cd frontend && npm run typecheck && npm run test:unit && npm run build` | exit 0 |
| E2E | `cd frontend && npx playwright test` | all green (dockerized API + DB) |

## Error Handling

- API errors: `{ error: { code, message, details } }` — never leak stack traces to clients
- 500 → `{ code: "INTERNAL_ERROR" }` with full trace in structured logs only
- Cross-tenant access → 404 (not 403) to avoid enumeration
- Validation failures → 422 with field details
- Auth login rate limit → 429 `{ code: "RATE_LIMITED" }`
- Sales: `INVALID_STATUS` (409), `INSUFFICIENT_STOCK` (409), `VALIDATION_ERROR` (400)

## Observability

- Structured JSON logs via `slog` with `request_id`, `tenant_id`, `user_id`, `method`, `path`, `status`, `duration_ms`
- `tenant_id` / `user_id` populated only after `AuthJWT` — absent on public auth routes
- `/healthz` — process alive
- `/readyz` — PostgreSQL + Redis reachable
- Graceful shutdown: 30s timeout on `SIGTERM`
- HTTP server timeouts (`cmd/api/main.go`): ReadHeader 5s, Read 10s, Write 30s, Idle 60s

## Rate Limiting

- Scope: `POST /api/v1/auth/login` only
- Limit: 5 attempts / 15 minutes per client IP
- Redis key: `ratelimit:login:{ip}`
- **Degraded mode:** if Redis is unavailable, login is not rate-limited (fail-open)
- `/readyz` requires Redis today — deploy blocked if Redis is down

See `backend/docs/RATE_LIMITING.md`.

## Frontend Auth And Cache

| Behavior | Notes |
|---|---|
| Access token | Pinia memory only |
| Refresh | HttpOnly cookie, path `/api/v1/auth` |
| 401 interceptor | Single-flight refresh; clears session on failure (no auto-redirect) |
| Logout | `queryClient.clear()` |
| Dashboard query key | `['dashboard', 'summary', userId]` |
| Cross-user risk | Login without prior logout may show stale lists until refetch; logout clears cache |

Low-stock threshold `5` is duplicated in Go and TypeScript — keep in sync or omit param and use server default.

## Transactional Writes (Sales)

| Failure | Outcome | Recovery |
|---|---|---|
| Insufficient stock on confirm | 409, tx rolled back | Reduce quantities |
| Invalid status transition | 409 | User action required |
| DB error mid-tx | 500, rolled back | Safe to retry draft create only |
| Audit insert fails post-commit | 200 today | Reconciliation / alert (known gap) |

See `backend/docs/SALES_TRANSACTIONS.md` for invariants and known concurrent-confirm race.

## Database

- Migrations run as separate job (`cmd/migrate`), not on app boot in production
- Stock `CHECK >= 0`; conditional decrement prevents negative stock per row
- Backup: daily `pg_dump`; test restore periodically (see `DEPLOYMENT.md`)

## Deployment Risks

- Secrets via env only (`.env.example` committed, `.env` gitignored)
- Refresh cookie must use `Secure=true` in production (see `backend/docs/AUTH_SESSION.md`)
- CORS via `ALLOWED_ORIGINS` must match frontend origin
- Fly CD gated on Backend CI success (see `docs/CI_CD.md`)
- `JWT_REFRESH_SECRET` is required in config but unused (refresh uses opaque token + SHA-256 hash)

See also: `backend/docs/DEPLOY_READINESS.md`, `DEPLOYMENT.md`, `docs/CI_CD.md`.
