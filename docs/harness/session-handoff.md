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

$ go test -tags=integration -count=3 ./internal/sales/   (live PG)
9/9 PASS — TestConcurrentConfirmSingleDecrement deterministic: 1×200 + 3×409, single decrement

$ backup.sh → restore.sh → row-count diff (dev compose)
DRILL PASS: users 3/3, customers 74/74; scratch db dropped

$ browser: 4 PNG screenshots from seeded stack; /audit page verified
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
5. **Portfolio polish** — DONE 2026-08-25: real screenshots captured from seeded stack replace SVG mockups (`docs/images/*.png`, README updated); "user mutation audit trail" verified already shipped (`/audit` admin page + service-level audit events on create/update/delete/confirm/cancel)

## Next Best Step

1. Pick first P1 item (Fly CD or dashboard drill-down)
2. Commit P0 doc + harness updates when user approves (suggested groups in `docs/harness/sprint-contract.md`)

## Verified Now — 2026-08-31

- `npm run typecheck` — PASS.
- `npm run test:unit` — PASS, 14 tests.
- `npm run build` — PASS.
- Browser smoke — PASS for the login surface at desktop size in light and dark themes; fields are empty, placeholders render, and theme persistence works.
- `docs/images/login-light.png` and `docs/images/login-dark.png` were captured and read back as valid PNG assets.
- `admin123` is absent from the built frontend assets.

## Changed This Session

- Continued the incumbent slate/blue ERP polish without redesigning the visual world.
- Added semantic light/dark tokens, stronger focus rings, themed controls, shell navigation icons, user avatar, refined login, table surfaces, responsive horizontal table scrolling, status badges, stable sale item keys, and Portuguese 404 styling.
- Added shared `frontend/src/shared/ui/AppDialog.vue` with modal semantics, Escape close, focus trap/restore, scroll lock, responsive sizing, and reduced-motion transition; migrated all feature dialogs.
- Updated `README.md` with light/dark login screenshots and synchronized `docs/DESIGN_SYSTEM.md`.

## Broken Or Unverified

- No runtime blocker remains after Docker Desktop started; E2E and authenticated browser smoke completed.
- No measured gate report exists under `docs/harness/quality/*.json`; coverage, mutation, regression, E2E, complexity, boundary, and security metrics remain unavailable.
- `make validate` is unavailable because `make` is not installed; Windows-equivalent baseline checks passed.
- Final review-agent retries failed with provider HTTP 429; the earlier design and code reviews were completed, and their critical findings were addressed.
- Impeccable detector ran once and reported the intentional Inter font warning plus a temporary accent warning; the accent was removed afterward without a second detector run per bounded QA rules.

- GitHub PR E2E run `33428231102` failed before the lifecycle correction in `CreateSaleDialog`; focused local regression passes and requires a rerun after the fix is pushed.

## Decisions Made

- Keep Inter and the existing slate/blue token world; the detector font warning is an intentional portfolio identity exception.
- Keep README theme screenshots focused on the login surface; authenticated dashboard light/dark and mobile states were inspected in-browser but are not added as static assets.
- Do not commit or push automatically; the worktree remains operator-owned.

## Next Best Step

Push the dialog lifecycle correction, wait for GitHub checks, then merge the PR into `main`.

## Commands

```text
npm run typecheck                         PASS
npm run test:unit                         PASS (14 tests)
npm run build                             PASS
node .../impeccable/scripts/detect.mjs    exit 2 (2 warnings; see Broken Or Unverified)
make validate                             exit 127 (make unavailable)
Windows baseline equivalents              PASS
npm run test:e2e                         PASS (21 tests)
browser authenticated smoke              PASS (dashboard themes/mobile; modal focus/Escape/restore; table overflow)
npx playwright test e2e/sales.spec.ts --grep tenant-isolation  PASS (1 test)
```

No PREVC lifecycle state, quality report, or commit trailer was produced.
