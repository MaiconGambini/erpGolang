# Harness Progress

## Current Active Work

**Thermo-nuclear remediation** — complete (pending full E2E run + Fly CD secret).

## Completed (this session)

- Centralized low-stock threshold (backend + frontend)
- Dashboard code simplification + Vue Query cache safety
- Dashboard handler/integration/E2E tests
- CI: Go 1.25, golangci-lint, integration job, frontend unit tests
- CD: `deploy.yml` for Fly.io

## Verification Status

```text
$ cd backend && go test ./...
exit 0

$ cd frontend && npm run typecheck && npm run test:unit
exit 0 (12 tests)
```

## Next Best Action

Run `npx playwright test`; configure `FLY_API_TOKEN`; commit when user approves.
