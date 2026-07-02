# goERP

goERP is a modular, tenant-aware ERP monorepo for small and medium-sized businesses.

The MVP starts with authentication, tenant-aware users, a dashboard, and a complete customers module. The first implementation goal is to run locally with PostgreSQL and Redis, then deploy simply to a VPS with Docker Compose and Caddy. AWS can follow with EC2 and RDS once the MVP is proven.

## Target Stack

- Backend: Go 1.23+, chi, pgx, sqlc, Atlas, Redis, JWT, bcrypt, slog.
- Frontend: Vite, Vue 3, TypeScript, Pinia, TanStack Vue Query, PrimeVue, Tailwind CSS, SCSS.
- Data: PostgreSQL 16 as source of truth, Redis 7 for cache/session coordination/rate limits.
- Local infra: Docker Compose for PostgreSQL and Redis only.
- First deployment target: VPS + Docker Compose + Caddy.

## Repository Layout

```text
backend/                 Go API and database layer
frontend/                Vue application
docs/                    Product, engineering, UX, and AI context docs
contract/                Future OpenAPI contract boundary
docker-compose.dev.yml   Local PostgreSQL and Redis
```

## Local Quick Start

1. Copy environment variables:

```bash
cp .env.example .env
```

2. Start local infrastructure:

```bash
docker compose -f docker-compose.dev.yml up -d
```

3. Run the backend:

```bash
cd backend
go run ./cmd/api
```

4. Run the frontend:

```bash
cd frontend
npm install
npm run dev
```

## Common Commands

```bash
docker compose -f docker-compose.dev.yml ps
docker compose -f docker-compose.dev.yml logs -f
docker compose -f docker-compose.dev.yml down
docker compose -f docker-compose.dev.yml down -v
```

Backend:

```bash
cd backend
go test ./...
go build ./...
```

Frontend:

```bash
cd frontend
npm run typecheck
npm run build
```

## Documentation Index

- `ARCHITECTURE.md` - system architecture and module boundaries.
- `CONVENTIONS.md` - naming, API responses, pagination, tenant rules.
- `DEPLOYMENT.md` - VPS-first deployment path and AWS evolution.
- `docs/AI_CONTEXT.md` - primary context for AI agents.
- `docs/GLOSSARY.md` - domain terms.
- `docs/BUSINESS_RULES.md` - business invariants.
- `docs/ROLES.md` - personas and permissions.
- `docs/UX_PATTERNS.md` - reusable UI behavior.
- `docs/DESIGN_SYSTEM.md` - light goERP visual system.
- `docs/MODULE_TEMPLATE.md` - how to add a business module.
- `docs/PROMPT_TEMPLATES.md` - reusable prompts for future agent work.

## MVP Definition of Ready

- Backend compiles and tests pass.
- Frontend typecheck/build pass.
- Migrations and seeds can recreate a local environment once the database implementation phase is complete.
- Tenant isolation is tested for every tenant-owned module.
- New modules follow the customers module template once it exists.
- Docs are updated when architecture, APIs, or business rules change.
