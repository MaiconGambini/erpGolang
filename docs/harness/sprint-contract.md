# Sprint Contract

## Sprint Overview

- Objective: Complete goERP MVP per `plan.md` phases 3–10 (gap-fill from current scaffold)
- Feature ID: goerp-mvp-1
- Date: 2026-07-03

## Scope

**In scope:**

- Gap analysis and completion of phases 3–7 (backend foundation through customers E2E)
- Phase 8 tests (unit, integration, handler, Vitest, Playwright smoke)
- Phase 9 observability (logs, health, shutdown, rate limit, CORS)
- Phase 10 deploy artifacts (Dockerfiles, CI, Caddy)

**Out of scope:**

- New business modules beyond customers
- i18n, dark mode, PWA, GraphQL, microservices
- Tenant self-registration (seed/CLI only for MVP)

## Roles

- Planner: spec-lead + PREVC
- Generator: domain subagents (Go backend, Vue frontend)
- Evaluator: project-judge

## Acceptance Criteria

- [ ] `make validate` exits 0
- [ ] Login + refresh + logout work via API and UI
- [ ] Customers CRUD with tenant isolation and audit logs
- [ ] At least one backend service test and one handler test
- [ ] At least one Playwright E2E (login + create customer)
- [ ] README local quick-start verified in 5 steps

## Verification Plan

| Check | Command | Pass Condition |
|---|---|---|
| Infra | `docker compose -f docker-compose.dev.yml config` | exit 0 |
| Backend | `cd backend && go test ./...` | exit 0, tests exist |
| Backend build | `cd backend && go build ./...` | exit 0 |
| Frontend | `cd frontend && npm run typecheck && npm run build` | exit 0 |
| Full | `make validate` | exit 0 |

## Evidence Log

| Check | Output | Pass? |
|---|---|---|
| docker compose config | name: erpgolang | yes |
| go test | no test files, exit 0 | partial |

## Sprint Log

- 2026-07-03: Harness bootstrap applied. PREVC Plan pending user approval.
