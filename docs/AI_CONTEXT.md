# AI Context

goERP is a modular, tenant-aware ERP monorepo.

## Target Stack

- Backend: Go 1.23+, chi, pgx, sqlc, Atlas, Redis, JWT, bcrypt, slog.
- Frontend: Vite, Vue 3, TypeScript, Pinia, TanStack Vue Query, PrimeVue, Tailwind CSS.
- Data: PostgreSQL 16 and Redis 7.
- Local infra: Docker Compose for infrastructure only.

## Core Rules For Agents

- Never create a tenant-owned table without `tenant_id`.
- Never access tenant-owned data without filtering by `tenant_id`.
- Never trust tenant ID from a request body.
- Keep API errors in `{ error: { code, message, details } }` format.
- Keep paginated responses in `{ data, pagination }` format.
- Do not add new dependencies without a clear reason.
- Follow the customers module as the reference module once implemented.

## How To Add A Module

1. Read `docs/GLOSSARY.md` and `docs/BUSINESS_RULES.md`.
2. Add schema/migration with tenant rules.
3. Add sqlc queries with explicit `tenant_id` parameters.
4. Add backend domain, DTO, service, handler, repository, and tests.
5. Mirror DTOs in frontend entity types.
6. Add list/create/edit/delete features as needed.
7. Update sidebar/routes.
8. Verify tenant isolation and API response shapes.

## Validation Expectations

- Backend: `go test ./...` and `go build ./...`.
- Frontend: `npm run typecheck` and `npm run build`.
- Infra: `docker compose -f docker-compose.dev.yml config`.
