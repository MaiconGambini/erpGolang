# Harness Progress

## P0 — Local docs + MVP polish — **complete** (2026-07-06)

- [x] README five-step local run aligned with `docker-compose.dev.yml` (Postgres **5434**, Redis **6381**)
- [x] `backend/docs/LOCAL_DEV.md` + `backend/.env.example` synced with compose ports and Windows `127.0.0.1` note
- [x] Sale draft edit UI (`EditSaleDialog` + backend PATCH)
- [x] Playwright port fix — defaults `DATABASE_URL` `:5434`, `REDIS_URL` `:6381`; E2E Vite on **5174** (`playwright.config.ts`)
- [x] E2E suite **15/15** passing

## Current Active Work

**P1 backlog — CLEAR** (2026-08-25). UI Polish lane entregue via impeccable `polish`: dark mode com tokens (toggle no topbar, persistência, anti-FOUC), skeletons + empty states em 6 superfícies, shell responsivo ≤768px (sidebar virava conteúdo espremido), gráficos vivos (range 90d + `cmd/seed` com dataset demo idempotente), superfícies de browser tematizadas. Verificação: typecheck/build/unit(14)/E2E(21) verdes, detect.mjs (1 warning de fonte = exceção intencional), inspeção visual batched + rodada de confirmação. Screenshots do README re-capturados (dark, dados demo limpos). Commits b59fec0..510bb8a.

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

Backup/restore drill EXECUTED 2026-08-25 (dev compose): PASS — see DEPLOYMENT.md drill log
$ browser: login/dashboard/sales/customer-form captured from seeded stack (goerp-api :8080 + goerp-web :5173, hub-managed)
$ /audit page renders tenant action table — mutation audit trail confirmed present
```

## Session Start Failure — 2026-08-31

- Startup command discovered from `Makefile`: `make validate`.
- Exact failure: `error: command not found: make` (exit 127).

- Windows baseline equivalents passed: `docker compose ... config`; `go test ./...`; `go build ./...`; `npm run typecheck`; `npm run build`.

## Frontend polish extension — 2026-08-31

- Shared semantic theme tokens consumed by shell, login, tables, status badges, errors, and dialogs; hardcoded credential defaults removed from the client bundle.
- Added shared `AppDialog` with `aria-modal`, Escape close, focus trap/restore, scroll lock, responsive sizing, and reduced-motion transition; migrated all feature dialogs.
- Completed responsive table overflow, keyboard/focus states, stable sale item keys, Portuguese 404, and status/danger contrast fixes.
- Added `docs/images/login-light.png` and `docs/images/login-dark.png`; README now presents both themes.
- Verification: `npm run typecheck` PASS; `npm run test:unit` PASS (14 tests); `npm run build` PASS; browser light/dark login screenshots PASS; `admin123` absent from built assets.
- Detector run once: Inter warning is an intentional incumbent-font exception; the temporary top accent warning was removed after the scan per bounded QA rules.

## Verification Update — 2026-08-31

- Docker Desktop started; `npm run test:e2e` PASS — 21 tests.
- Authenticated browser smoke PASS: dashboard light/dark at desktop, dashboard dark at mobile, customer modal initial focus + Escape close + focus restoration, and mobile table horizontal overflow.

## Next Best Action

Review the final diff, commit the frontend polish and screenshot assets, then push to GitHub.
