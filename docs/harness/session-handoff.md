# Session Handoff

## P0 Status — **complete** (2026-07-06)

Local documentation and MVP polish items are done:

| Item | Evidence |
|------|----------|
| README 5-step local run | Ports **5434** / **6381**; link to `backend/docs/LOCAL_DEV.md` |
| Docs ↔ compose ↔ `.env.example` | `docker compose -f docker-compose.dev.yml config` — exit 0 |
| Sale draft edit UI | `EditSaleDialog.vue`, `use-edit-sale.ts`, table/detail edit actions |
| Playwright port fix | `playwright.config.ts` — `5434`/`6381` defaults, Vite E2E on **5174** |
| Full E2E | `npx playwright test` — **15/15** passed |

## Verified Now

```text
$ docker compose -f docker-compose.dev.yml config
exit 0 (postgres published 5434, redis published 6381)

$ cd backend && go test ./...
ok  auth, customers, dashboard, products, sales, suppliers, shared/inventory
exit 0

$ cd frontend && npm run typecheck && npm run test:unit
vue-tsc — exit 0
vitest run — 12 passed (4 files)
exit 0

$ cd frontend && npx playwright test
15 passed

$ npx @redocly/cli lint contract/openapi.yaml
valid — 0 errors, 6 warnings (probe/logout ops have no 4XX by design); coverage 40/40 chi routes

$ cd backend && sqlc generate && go build ./... && go test ./internal/sales/...
sqlc v1.31.1 OK; unit tests pass

$ go test -tags=integration -run TestConcurrent ./internal/sales/   (needs live PG)
SKIP — database unavailable (Docker daemon down locally); CI exercises it

$ sh -n deploy/scripts/*.sh; python -c yaml.safe_load(deploy-vps.yml); docker compose config (prod)
SYNTAX_OK_SH · WORKFLOW_YAML_OK · COMPOSE_CONFIG_OK
```

## Thermo-Nuclear Remediation (prior session)

**Code quality**
- `backend/internal/shared/inventory/threshold.go` — shared `DefaultLowStockThreshold` + `NormalizeLowStockThreshold`
- `frontend/src/shared/config/inventory.ts` — `LOW_STOCK_THRESHOLD`
- Removed `SummaryDTO` pass-through; `DashboardPage.vue` metrics via `v-for`
- Vue Query: `queryClient.clear()` on logout; `invalidateDashboardSummary` on KPI-affecting mutations

**Tests**
- `backend/internal/dashboard/handler_test.go` — 401 + 200 envelope (skips without `DATABASE_URL`)
- `backend/internal/dashboard/integration_test.go` — tenant-scoped summary counts
- `frontend/e2e/dashboard.spec.ts` — KPI smoke + active customers increment

**CI/CD**
- `.github/workflows/backend.yml` — Go 1.25, golangci-lint, integration job with seed
- `.github/workflows/frontend.yml` — `npm run test:unit`
- `.github/workflows/e2e.yml` — Go 1.25
- `.github/workflows/deploy.yml` — `flyctl deploy` on push to `main` (needs `FLY_API_TOKEN`)

## Blockers

- **Fly CD:** set GitHub secret `FLY_API_TOKEN` + Fly app secrets (`DATABASE_URL`, `REDIS_URL`, `JWT_*`, `ALLOWED_ORIGINS`)
- **Commits:** deferred until user approves

## P1 — Next items

1. **Fly CD** — configure `FLY_API_TOKEN` and verify deploy workflow on `main`
2. **OpenAPI contract** — DONE 2026-08-25: `contract/openapi.yaml` (OpenAPI 3.0.3, 40/40 chi routes, `redocly lint` clean, bare-vs-wrapped envelope inconsistency documented)
3. **Sales race** — DONE 2026-08-25: verified the documented gap was already half-closed (conditional `UpdateSaleStatus`), added `GetSaleForUpdate FOR UPDATE` recheck in `Confirm`/`Cancel`, regression test `TestConcurrentConfirmSingleDecrement`; doc `SALES_TRANSACTIONS.md` synced
4. **Ops** — DONE 2026-08-25 (modulo live drill): `.github/workflows/deploy-vps.yml` (workflow_dispatch; secrets `VPS_SSH_KEY`/`VPS_HOST`/`VPS_USER` needed before first run), `deploy/scripts/backup.sh` (+verify+retention), `deploy/scripts/restore.sh` (scratch-db default + live-guard), `DEPLOYMENT.md §Backup & Restore Drill` with cron tiers + pass criteria; compose config validated offline
5. **Portfolio polish** — replace `docs/images/*.svg` with real screenshots; user mutation audit trail

## Next Best Step

1. Pick first P1 item (Fly CD or dashboard drill-down)
2. Commit P0 doc + harness updates when user approves (suggested groups in `docs/harness/sprint-contract.md`)
