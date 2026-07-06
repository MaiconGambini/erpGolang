# Harness Status

## Goal

- **Source of truth:** `~/.cursor/harness-goals/goal-state.json` → `goal-20260706-095500`
- **Status:** **complete** (all requirements achieved)
- **Branch:** `main` (dirty — thermo remediation uncommitted)

## Readiness

| Check | Status |
|-------|--------|
| Phases 0–10 + post-MVP | passing in `feature_list.json` |
| Dashboard remediation | threshold centralization, tests, Vue Query fixes |
| Backend tests | **exit 0** (`go test ./...`) |
| Frontend unit | **12 passed** (`npm run test:unit`) |
| CI golangci-lint + integration | **configured** in `backend.yml` |
| Fly CD workflow | **added** — needs `FLY_API_TOKEN` GitHub secret |
| Full Playwright (15 tests) | pending run with API + DB |

## Next Action

Run `npx playwright test`; add `FLY_API_TOKEN`; commit when user approves (see sprint-contract commit groups).
