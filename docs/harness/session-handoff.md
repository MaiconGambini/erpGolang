# Session Handoff

## P0 Status — **complete** (2026-07-06, historical evidence)

Local documentation and MVP polish items are done:

| Item | Evidence |
|------|----------|
| README 5-step local run | Ports **5434** / **6381**; link to `backend/docs/LOCAL_DEV.md` |
| Docs ↔ compose ↔ `.env.example` | `docker compose -f docker-compose.dev.yml config` — exit 0 |
| Sale draft edit UI | `EditSaleDialog.vue`, `use-edit-sale.ts`, table/detail edit actions |
| Playwright port fix | `playwright.config.ts` — `5434`/`6381` defaults, Vite E2E on **5174** |
| Full E2E | `npx playwright test` — **15/15** passed (historical, superseded by the later **21-test** result below) |

## Verified Now — historical baseline (prior to 2026-08-31)

```text
$ docker compose -f docker-compose.dev.yml config
exit 0 (postgres published 5434, redis published 6381)

$ cd backend && go test ./...
ok  auth, customers, dashboard, products, sales, suppliers, shared/inventory
exit 0

$ cd frontend && npm run typecheck && npm run test:unit
vue-tsc — exit 0
Historical baseline: vitest run — 12 passed (4 files)
exit 0

$ cd frontend && npx playwright test
Historical baseline: `npx playwright test` — 15 passed (superseded by the later **21-test** result below)
Historical baseline note: the earlier full E2E run was pending at **20 tests**; the later **21-test** result below supersedes that status.

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

## P1 — Next items (historical snapshot; superseded by later verified sections)

The dated facts and deployment caveats below are preserved as historical evidence. Later verified sections are authoritative for current status.

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

## Broken Or Unverified — 2026-08-31 historical snapshot

- Historical 2026-08-31 snapshot: no runtime blocker remained after Docker Desktop started; E2E and authenticated browser smoke completed.
- Historical 2026-08-31 snapshot: no measured gate report existed under `docs/harness/quality/*.json`; coverage, mutation, regression, E2E, complexity, boundary, and security metrics remained unavailable.
- Historical 2026-08-31 snapshot: `make validate` was unavailable because `make` was not installed; Windows-equivalent baseline checks passed.
- Historical 2026-08-31 snapshot: final review-agent retries failed with provider HTTP 429; the earlier design and code reviews were completed, and their critical findings were addressed.
- Historical 2026-08-31 snapshot: Impeccable detector ran once and reported the intentional Inter font warning plus a temporary accent warning; the accent was removed afterward without a second detector run per bounded QA rules.

- Historical 2026-08-31 snapshot: GitHub PR E2E run `33428231102` failed before the lifecycle correction in `CreateSaleDialog`; focused local regression passed and required a rerun after the fix was pushed.




## Decisions Made

- Keep Inter and the existing slate/blue token world; the detector font warning is an intentional portfolio identity exception.
- Keep README theme screenshots focused on the login surface; authenticated dashboard light/dark and mobile states were inspected in-browser but are not added as static assets.
- Do not commit or push automatically; the worktree remains operator-owned.

## Next Best Step

Operator-controlled next action: review the current 18-file public-readiness diff and commit only if approved. This session does not claim a push.

## Commands — 2026-08-31 historical snapshot

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

Historical 2026-08-31 snapshot: no PREVC lifecycle state, quality report, or commit trailer was produced.

## Verified Now — 2026-09-11

- Quality gate `node ~/.config/opencode/scripts/harness-quality-gate.mjs --mode local --label public-readiness-final10` — PASS. Report: `docs/harness/quality/2026-09-11T01-40-44-local-public-readiness-final10-10540.json`.
- Risk router: tier `full/high`, credibility `credible`, `filesTouched=18`, `netLines=225`, with no new dependencies or schema change.
- Redteam `node ~/.omp/agent/scripts/harness-redteam-scan.mjs --metric` — PASS; `scanned=917`, `FAIL=0`.
- All five workflow YAML files parse — `WORKFLOWS_YAML_OK 5`.
- Production Compose config — PASS: `GOERP_ENV_FILE=../../deploy/env/production.env.example GOERP_DOMAIN=example.com docker compose -f deploy/compose/docker-compose.prod.yml config --quiet`.
- Seed guard — `go build ./cmd/seed` PASS; a production invocation with required dummy environment stopped at the new guard before DB connection.
- `git diff --check` — PASS, with only CRLF normalization warnings.
- Tracked secret-like paths query — empty.
- Baseline evidence: backend Go test/build PASS; frontend typecheck/build/unit PASS (14 tests); authenticated E2E PASS (21 tests, 2026-08-31) is historical evidence, not current-run evidence.

## Changed This Session — 2026-09-11

- Current public-readiness scope is the 18-file diff measured by the risk router, with `netLines=225`.
- All five workflows pin third-party actions to immutable commit SHAs, use action version v2.13.2 where applicable, and request `contents: read`.
- CI seed jobs set `APP_ENV=development`; the seed guard fails closed when `APP_ENV` is omitted.
- VPS deployment checks out the exact `DEPLOY_SHA`.
- `backend/docs/LOCAL_DEV.md` and `docs/RELIABILITY.md` use root-safe runbook commands.

The gate is PASS, but only 1/14 metrics is measured: 6 are `not_configured` and 7 are `unavailable`. Those gaps are not green evidence.

```text
Metric                       Value  Threshold  Status
line_coverage                -      -          not_configured (no coverage artifact found — run the test suite with coverage reporting first)
branch_coverage              -      -          not_configured (no coverage artifact found — run the test suite with coverage reporting first)
cyclomatic_max               -      -          not_configured (no complexity tool detected for this stack)
module_lines_max              596    <= 300     observe  backend\\internal\\sales\\service.go
security_findings             -      -          not_configured (no static security scanner detected for this stack)
boundary_violations           -      -          not_configured (no dependency-boundary config found (dependency-cruiser / import-linter / Packwerk))
regression_suite              -      -          not_configured (no regression suite declared — set suites.regression.command in agent-os/quality-thresholds.json)
rule_violations_enforced     -      -          unavailable (no findings file in docs/harness/findings — adherence unmeasured; a gap to name, never a pass, never a zero)
rule_violations_prose        -      -          unavailable (no findings file in docs/harness/findings — adherence unmeasured; a gap to name, never a pass, never a zero)
adherence_per_changed_lines  -      -          unavailable (no findings file in docs/harness/findings — adherence unmeasured; a gap to name, never a pass, never a zero)
unciteable_findings_ratio    -      -          unavailable (no findings file in docs/harness/findings — adherence unmeasured; a gap to name, never a pass, never a zero)
enforced_fraction            -      -          unavailable (no findings file in docs/harness/findings — adherence unmeasured; a gap to name, never a pass, never a zero)
rules_active                 -      -          unavailable (no findings file in docs/harness/findings — adherence unmeasured; a gap to name, never a pass, never a zero)
rules_retired                -      -          unavailable (no findings file in docs/harness/findings — adherence unmeasured; a gap to name, never a pass, never a zero)
```

## Broken Or Unverified — 2026-09-11

- PREVC inspect reports no lifecycle state (`WARN`); no PREVC lifecycle transition is claimed.
- `make validate` is unavailable on Windows.
- The Docker daemon is unavailable, so current live integration and E2E were not rerun.
- The prior authenticated 21-test E2E result from 2026-08-31 is historical evidence, not current-run evidence.
- The project static security scanner is unavailable; deterministic redteam scanning passed.
- The first live VPS drill and required deployment secrets remain pending.

## Commands — 2026-09-11

```text
node ~/.config/opencode/scripts/harness-quality-gate.mjs --mode local --label public-readiness-final10
report: docs/harness/quality/2026-09-11T01-40-44-local-public-readiness-final10-10540.json
verdict: PASS
risk router: full/high; credibility: credible; filesTouched=18; netLines=225; no new dependencies/schema change

node ~/.omp/agent/scripts/harness-redteam-scan.mjs --metric
status: PASS; scanned=917; FAIL=0

workflow YAML parse
WORKFLOWS_YAML_OK 5

GOERP_ENV_FILE=../../deploy/env/production.env.example GOERP_DOMAIN=example.com docker compose -f deploy/compose/docker-compose.prod.yml config --quiet
PASS

go build ./cmd/seed
PASS
production invocation with required dummy environment
stopped at the new guard before DB connection

git diff --check
PASS (only CRLF normalization warnings)

tracked secret-like paths query
empty

backend Go test/build
PASS
frontend typecheck/build/unit
PASS (14 tests)
authenticated E2E (2026-08-31)
PASS (21 tests; historical, not current-run evidence)
make validate
unavailable on Windows
PREVC inspect
WARN (no lifecycle state)
current live integration and E2E
not rerun; Docker daemon unavailable
project static security scanner
unavailable
```

commit trailer:
1ab8373 Harden CI and deployment workflows
117a8ec Refresh public documentation and runbooks
e7f755d Record public release readiness evidence
push verification: `HEAD == origin/main == e7f755dcea322a4cc3089945169e3c1c07d094fe`; worktree clean
backend verification: `go test ./... && go build ./...` — PASS
