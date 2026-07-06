# Harness Eval Contract

## Capability Evals

- Dashboard KPI API returns tenant-scoped counts
- Module CRUD E2E per domain (customers, products, suppliers, sales)

## Regression Evals

- `go test ./...` backend
- `npm run typecheck` + `npm run test:unit` frontend
- `npx playwright test` E2E

## Command Evals

- `make validate` when Makefile present

## Judge Rubric

- See `agent-os/judges/project-judge.md`

## Evidence

- Exact command output in `docs/harness/session-handoff.md`
