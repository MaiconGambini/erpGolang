# Sprint Contract

## Sprint Overview

- **Objective:** Thermo-nuclear review remediation — dashboard quality, tests, CI/CD
- **Feature ID:** `post-mvp-dashboard`, production checklist gaps
- **Date:** 2026-07-06

## Acceptance Criteria

- [x] Centralized low-stock threshold (backend `shared/inventory`, frontend `shared/config/inventory`)
- [x] Dashboard handler + integration tests; E2E `dashboard.spec.ts`
- [x] Vue Query cache cleared on logout; dashboard invalidation on mutations
- [x] CI: Go 1.25, golangci-lint, integration job, frontend `test:unit`
- [x] CD: `.github/workflows/deploy.yml` (requires `FLY_API_TOKEN` secret)
- [x] `go test ./...`, `npm run typecheck`, `npm run test:unit` pass locally

## Evidence Log (2026-07-06 remediation)

| Check | Output | Pass? |
|-------|--------|-------|
| Backend unit | `go test ./...` — exit 0 (dashboard tests skip without DATABASE_URL) | yes |
| Inventory unit | `go test ./internal/shared/inventory/...` — exit 0 | yes |
| Frontend typecheck | `vue-tsc --noEmit` — exit 0 | yes |
| Frontend unit | `vitest run` — 12 passed | yes |
| Dashboard E2E | `e2e/dashboard.spec.ts` added (2 tests) | pending full suite run |
| CI workflows | `backend.yml`, `frontend.yml`, `e2e.yml`, `deploy.yml` updated | yes |

## Project Judge Verdict: **Accept** (remediation scope)

Historical note from the 2026-07-06 remediation: the full `npx playwright test` run with API + DB was pending at that time for 20 tests; the later verification dated 2026-08-31 recorded Playwright E2E PASS — 21 tests. User sets `FLY_API_TOKEN` for live CD.

## Suggested commit groups

1. `fix(docker): align migrate Dockerfile to Go 1.25`
2. `feat(dashboard): live KPI summary API + page`
3. `refactor(inventory): centralize low-stock threshold`
4. `test(dashboard): handler, integration, e2e`
5. `ci: go 1.25, golangci-lint, integration, test:unit`
6. `ci: fly deploy on main`
7. `chore(harness): policy docs + handoff`
