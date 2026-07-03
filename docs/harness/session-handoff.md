# Session Handoff

## Verified Now

```text
$ cd backend && go build -o api.exe ./cmd/api
(exit 0)

$ cd backend && go test ./...
ok  	github.com/MaiconGambini/erpGolang/backend/internal/auth	(cached)
ok  	github.com/MaiconGambini/erpGolang/backend/internal/customers	(cached)
ok  	github.com/MaiconGambini/erpGolang/backend/internal/products	(cached)
ok  	github.com/MaiconGambini/erpGolang/backend/internal/suppliers	(cached)
ok  	github.com/MaiconGambini/erpGolang/backend/internal/sales	(cached)

$ git status
On branch main
Your branch is ahead of 'origin/main' by 5 commits.
nothing to commit, working tree clean

$ git log --oneline -5
a2365f4 feat(sales): sales module with stock confirm/cancel workflow
8ccf43c feat(suppliers): suppliers CRUD module
53eb926 feat(products): catalog module with low-stock endpoint
f25d1f3 feat: MVP phases 0-10 — auth, customers, deploy, and CI
8f0a5d2 chore(harness): bootstrap AGENTS, agent-os, and project docs
```

Prior session (not re-run this handoff pass):

```text
$ cd frontend && npm run typecheck && npm run test:unit
Test Files  4 passed (4), Tests  12 passed (12)

$ npx playwright test e2e/{auth-customers,products,suppliers,sales}.spec.ts
13 passed individually (5+3+3+2); full 14-test suite not run as one command
```

## Changed This Session

**Ordered git commits (local, not pushed)**

| SHA | Message |
|-----|---------|
| `8f0a5d2` | chore(harness): bootstrap AGENTS, agent-os, and project docs |
| `f25d1f3` | feat: MVP phases 0-10 — auth, customers, deploy, and CI |
| `53eb926` | feat(products): catalog module with low-stock endpoint |
| `8ccf43c` | feat(suppliers): suppliers CRUD module |
| `a2365f4` | feat(sales): sales module with stock confirm/cancel workflow |

**MVP phases 0–10** (`f25d1f3`)
- Backend: chi, pgx, Redis, JWT, middleware, httpx, sqlc, customers CRUD, audit, migrate/seed CLIs
- Frontend: FSD auth, customers UI, Playwright e2e
- Deploy: fly.toml, prod compose migrate job, Caddy `/api/*` fix, GitHub Actions workflows

**Post-MVP — Products** (`53eb926`)
- `backend/internal/products/`, migration `202601020001_products.sql`, `GET /products/low-stock`
- Frontend `/products`, `e2e/products.spec.ts`

**Post-MVP — Suppliers** (`8ccf43c`)
- `backend/internal/suppliers/`, migration `202601030001_suppliers.sql`
- Frontend `/suppliers`, `e2e/suppliers.spec.ts`, `feature_list.json` post-mvp entries

**Post-MVP — Sales** (`a2365f4`)
- `backend/internal/sales/` — draft → confirmed → cancelled; stock on confirm/cancel
- Migration `202601040001_sales.sql`, frontend `/sales`, `e2e/sales.spec.ts`
- `main.go` registers auth, users, tenants, customers, products, suppliers, sales, audit

**Harness**
- `AGENTS.md`, `agent-os/`, `docs/{ARCHITECTURE,PRODUCT,RELIABILITY}.md`, `docs/harness/*`
- `feature_list.json` — phases 0–10 + products/suppliers/sales marked `passing`

## Broken Or Unverified

- Full `npx playwright test` (14 tests) — not run as single suite; modules tested individually.
- `fly deploy` — config only, not executed against Fly.io.
- Prod `docker compose up --build` on VPS — not run.
- `go test -tags=integration` — blocked on Windows App Control locally; intended for CI/Linux.
- Stale `api.exe` on `:8080` causes 404 for new routes if not rebuilt before E2E (`reuseExistingServer: true` when `CI` unset).

## Decisions Made

- Post-MVP modules follow customers FSD template; sqlc output committed manually when Windows App Control blocks `sqlc.exe`.
- Sales prices snapshotted from product at draft creation; stock moves only on confirm/cancel.
- Money as JSON strings (`"19.90"`) with `NUMERIC(12,2)` in Postgres — no decimal library.
- Migrate via `cmd/migrate` + `schema_migrations` table.
- Low-stock default threshold: 5 (`GET /products/low-stock?threshold=`).
- Session work split into ordered commits: harness → MVP → products → suppliers → sales.

## Next Best Step

Push `main` (5 commits ahead of `origin/main`), rebuild `api.exe`, run full `npx playwright test`, then either **Fly/VPS deploy** or **wire dashboard KPIs** to live APIs.

## Commands

```bash
# Infra
docker compose -f docker-compose.dev.yml up -d

# Migrations (mark initial if DB seeded manually)
docker exec goerp-postgres psql -U goerp -d goerp -c \
  "INSERT INTO schema_migrations (filename) VALUES ('202601010001_initial.sql') ON CONFLICT DO NOTHING;"
cd backend && go build -o migrate.exe ./cmd/migrate && .\migrate.exe

# API (rebuild after every backend module change)
cd backend && go build -o api.exe ./cmd/api && .\api.exe

# Seed (tenants acme/beta, admin123)
cd backend && go build -o seed.exe ./cmd/seed && .\seed.exe

# Validation
cd backend && go test ./...
cd frontend && npm run typecheck && npm run test:unit
cd frontend && npx playwright test

# Push local commits
git push origin main

# Prod config check
GOERP_ENV_FILE="$(pwd)/deploy/env/production.env.example" GOERP_DOMAIN=example.com \
  docker compose -f deploy/compose/docker-compose.prod.yml config
```
