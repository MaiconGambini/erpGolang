# Harness Progress

## P0 — Local docs + MVP polish — **complete** (2026-07-06)

- [x] README five-step local run aligned with `docker-compose.dev.yml` (Postgres **5434**, Redis **6381**)
- [x] `backend/docs/LOCAL_DEV.md` + `backend/.env.example` synced with compose ports and Windows `127.0.0.1` note
- [x] Sale draft edit UI (`EditSaleDialog` + backend PATCH)
- [x] Playwright port fix — defaults `DATABASE_URL` `:5434`, `REDIS_URL` `:6381`; E2E Vite on **5174** (`playwright.config.ts`)
- [x] E2E suite **15/15** passing

## Current Active Work

**P1** — Fly CD secret, dashboard drill-down, OpenAPI contract, ops hardening (see `session-handoff.md`).

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
```

## Next Best Action

Set `FLY_API_TOKEN` for Fly CD; start P1 dashboard drill-down or OpenAPI contract; commit when user approves.
