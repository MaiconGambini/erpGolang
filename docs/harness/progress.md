# Harness Progress

## P0 — Local docs + MVP polish — **complete** (2026-07-06)

- [x] README five-step local run aligned with `docker-compose.dev.yml` (Postgres **5434**, Redis **6381**)
- [x] `backend/docs/LOCAL_DEV.md` + `backend/.env.example` synced with compose ports and Windows `127.0.0.1` note
- [x] Sale draft edit UI (`EditSaleDialog` + backend PATCH)
- [x] Playwright port fix — defaults `DATABASE_URL` `:5434`, `REDIS_URL` `:6381`; E2E Vite on **5174** (`playwright.config.ts`)
- [x] E2E suite **15/15** passing

## Current Active Work

**P1** — remaining: Fly CD secret (blocked on operator). DONE 2026-08-25: OpenAPI contract (`contract/openapi.yaml`); sales confirm race guard (`agent-os/specs/2026-08-25-1354-sales-confirm-race/`); VPS deploy automation (`deploy-vps.yml`) + backup/restore scripts (`deploy/scripts/`) + drill runbook (`DEPLOYMENT.md §Backup & Restore Drill`).

## Completed (prior session)

- Centralized low-stock threshold (backend + frontend)
- Dashboard code simplification + Vue Query cache safety
- Dashboard handler/integration/E2E tests
- CI: Go 1.25, golangci-lint, integration job, frontend unit tests
- CD: `deploy.yml` for Fly.io

## Verification Status

```text
$ docker compose -f docker-compose.dev.yml config
exit 0

$ cd backend && go test ./...
exit 0

$ cd frontend && npm run typecheck && npm run test:unit
exit 0 (12 tests)

$ cd frontend && npx playwright test
15 passed
$ npx @redocly/cli lint contract/openapi.yaml
valid — 0 errors, 6 warnings (probe/logout ops have no 4XX response by design)

$ op extraction: contract/openapi.yaml vs backend/internal/**/module.go + app/routes.go
40/40 registered chi operations covered exactly
$ cd backend && go build ./... && go test ./internal/sales/...
exit 0 (unit); integration tag compiles and skips cleanly without Postgres

$ go test -tags=integration -run "TestConcurrentConfirmSingleDecrement|..." ./internal/sales/
3 SKIP (database unavailable — Docker daemon down); live assertion runs in CI
$ sh -n deploy/scripts/{backup,restore}.sh && python yaml parse deploy-vps.yml
SYNTAX_OK_SH / WORKFLOW_YAML_OK

$ GOERP_ENV_FILE=deploy/env/production.env.example docker compose -f deploy/compose/docker-compose.prod.yml config
COMPOSE_CONFIG_OK (offline validation)

NOT executed (Docker daemon down): first backup run, restore drill — logged as pending in DEPLOYMENT.md drill log
```

## Next Best Action

Commit `contract/openapi.yaml` when user approves; next P1 lane: sales confirm race guard (`backend/docs/SALES_TRANSACTIONS.md`) or ops hardening.
