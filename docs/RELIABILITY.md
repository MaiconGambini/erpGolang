# Reliability

## Startup And Verification

```bash
# Infra
docker compose -f docker-compose.dev.yml up -d

# Backend health (when running)
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz

# Baseline CI gate
make validate
```

| Layer | Command | Pass |
|---|---|---|
| Compose | `docker compose -f docker-compose.dev.yml config` | exit 0 |
| Backend | `cd backend && go test ./... && go build ./...` | exit 0 |
| Frontend | `cd frontend && npm run typecheck && npm run build` | exit 0 |

## Error Handling

- API errors: `{ error: { code, message, details } }` — never leak stack traces to clients
- 500 → `{ code: "INTERNAL_ERROR" }` with full trace in structured logs only
- Cross-tenant access → 404 (not 403) to avoid enumeration
- Validation failures → 422 with field details

## Observability

- Structured JSON logs via `slog` with `request_id`, `tenant_id`, `user_id`, `method`, `path`, `status`, `duration_ms`
- `/healthz` — process alive
- `/readyz` — PostgreSQL + Redis reachable
- Graceful shutdown: 30s timeout on `SIGTERM`

## Deployment Risks

- Secrets via env only (`.env.example` committed, `.env` gitignored)
- Migrations run as separate deploy job (not on app boot)
- Rate limit on `/api/v1/auth/login` (Redis, 5 attempts / 15min per IP)
- CORS via `ALLOWED_ORIGINS` env

See also: `backend/docs/DEPLOY_READINESS.md`, `DEPLOYMENT.md`.
