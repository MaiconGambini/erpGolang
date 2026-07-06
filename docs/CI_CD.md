# CI/CD

## Workflows

| Workflow | Trigger paths | Jobs |
|---|---|---|
| `backend.yml` | `backend/**`, `fly.toml`, workflow file | migrate, golangci-lint, unit test, build, integration (customers, dashboard, sales, reports), coverage profile |
| `frontend.yml` | `frontend/**` | npm ci, typecheck, unit test, build |
| `e2e.yml` | `frontend/**`, `backend/**` | migrate, seed, Playwright (20 tests) |
| `deploy.yml` | After **Backend CI** succeeds on `main` | `flyctl deploy` to Fly.io |

## Deploy Gate

`deploy.yml` uses `workflow_run` and deploys only when **Backend CI** completes with `success` on `main`. This prevents broken backend code from reaching Fly production without passing lint, unit tests, and integration tests.

Required secret: `FLY_API_TOKEN` in GitHub repository settings.

## Local Verification (mirrors CI)

```bash
make validate
cd backend && go test -tags=integration ./internal/customers/... ./internal/dashboard/... ./internal/sales/... ./internal/reports/...
cd frontend && npm run test:unit && npx playwright test
```

Integration tests require `DATABASE_URL` pointing at a migrated database (CI uses Postgres service; local: `docker-compose.dev.yml` on port `5434`).

## Deploy Profiles

| Profile | Target | Automation |
|---|---|---|
| **A — VPS full stack** (default) | `deploy/compose/docker-compose.prod.yml` + Caddy | Manual / operator SSH |
| **B — Fly API split** | `fly.toml` + external Postgres/Redis | GitHub Actions on green Backend CI |

Profile B does not deploy the frontend. Host static assets separately and set `ALLOWED_ORIGINS` to the frontend URL.

See `DEPLOYMENT.md` for operator runbooks.

## Known Gaps

- E2E is not a hard gate for Fly deploy (Backend CI only).
- No GitHub Actions workflow for VPS compose deploy.
- Branch protection should require Backend CI + Frontend CI status checks on PRs.
