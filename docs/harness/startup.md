# Harness Startup Path — goERP

## Detected Stack

- **Backend:** Go 1.25+ (`go.mod`), chi, pgx, sqlc, Atlas
- **Frontend:** Vite + Vue 3 + TypeScript (`frontend/package.json`)
- **Infra:** Docker Compose (`docker-compose.dev.yml`) — Postgres 16 + Redis 7

## Baseline Command

`make validate` is the canonical baseline (see root `Makefile`).

On Windows (no `make`), run equivalent steps:

```powershell
docker compose -f docker-compose.dev.yml config
cd backend; go test ./...; go build ./...
cd ..\frontend; npm run typecheck; npm run build
```

## Required Setup

1. Copy `.env.example` → `.env` at repo root and in `backend/` if needed.
2. `docker compose -f docker-compose.dev.yml up -d`
3. `cd backend && go run ./cmd/migrate` (or `make migrate`)
4. `cd backend && go run ./cmd/seed` (seed tenants + admin users)

## Dev Servers

```powershell
cd backend; go run ./cmd/api
cd frontend; npm run dev
```

## Focused Checks

| Check | Command |
|-------|---------|
| Backend unit | `cd backend && go test ./...` |
| Frontend typecheck | `cd frontend && npm run typecheck` |
| E2E (full) | `cd frontend && npx playwright test` (port **5174**, see below) |
| Prod compose | `make prod-config` or see `deploy/` |

## Windows E2E (App Control blocks Go binaries)

When Device Guard blocks `api.exe` / `go run`, start the API in Docker before Playwright:

```powershell
docker compose -f docker-compose.dev.yml up -d
cd backend; go run ./cmd/migrate; go run ./cmd/seed

docker build -t goerp-api-e2e ./backend
docker run --rm -p 8080:8080 `
  -e DATABASE_URL="postgres://goerp:goerp@host.docker.internal:5432/goerp?sslmode=disable" `
  -e REDIS_URL="redis://host.docker.internal:6379/0" `
  -e JWT_ACCESS_SECRET=test-access `
  -e JWT_REFRESH_SECRET=test-refresh `
  goerp-api-e2e

cd frontend
npx playwright install chromium   # first run only
npx playwright test               # uses port 5174 (--strictPort)
```

Playwright config uses port **5174** to avoid colliding with other Vite apps on `:5173`.

## Risks

- `make` not available on Windows — use PowerShell equivalents above.
- sqlc binary may be blocked by App Control; commit generated `gen/db/*.go` manually.
- E2E requires running API + frontend + seeded DB.

## Session Start

Invoke `/harness-session-start` at the beginning of each session or subagent run.
