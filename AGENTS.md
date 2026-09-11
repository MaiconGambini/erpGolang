## Harness

Read this first. These rules apply to every session.

### Session Start

Invoke `harness-session-start` before feature work. It discovers instructions, progress state, feature state, handoff, startup path, recent commits, and declares exactly one active task.

### Session End

Invoke `harness-clean-handoff` before closing. Record verification, blockers, next action, and git status.

### Development Rules

- WIP=1.
- Plan before editing.
- Use standards from `agent-os/standards/`.
- Use `agent-os/specs/` for meaningful work.
- Completion requires evidence, not confidence.
- Secrets stay server-only.

### goERP Project Context

- Stack: Go 1.25+ (chi, pgx, sqlc, Atlas, Redis, gofpdf) + Vue 3 (FSD, PrimeVue, Pinia, Vue Query, Chart.js).
- Plan: `plan.md` (phases 0–10). Context: `context.md` (brainstorm), canonical docs: `docs/README.md`.
- Reference modules: `customers` (CRUD), `sales` (workflow), `dashboard` + `reports` (read-model).
- RBAC: `docs/ROLES.md`; middleware `RequireRole`; frontend `shared/lib/roles.ts`.
- Never create tenant-owned tables without `tenant_id`.
- Never access tenant data without filtering by `tenant_id`.
- API errors: `{ error: { code, message, details } }`.
- Paginated responses: `{ data, pagination }`.

### Core Commands

- `/harness-bootstrap` — install full harness with confirmation.
- `/harness-session-start` — start session.
- `/prevc` — plan, review, execute, validate, judge, confirm, handoff.
- `/harness-clean-handoff` — close session.

### Baseline Verification

```bash
make validate
# or:
docker compose -f docker-compose.dev.yml config
cd backend && go test ./... && go build ./...
go test -tags=integration ./internal/customers/... ./internal/sales/... ./internal/reports/... ./internal/dashboard/...
cd frontend && npm run typecheck && npm run build && npx playwright test
```
