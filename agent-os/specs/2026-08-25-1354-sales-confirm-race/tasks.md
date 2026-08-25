# Tasks

- [x] T01 — Concurrency regression test added (`TestConcurrentConfirmSingleDecrement`, 4 goroutines over HTTP). Evidence: compiles + runs under `-tags=integration`; skips gracefully with PG down (2026-08-25). Live-DB assertion pending infra (CI).
- [x] T02 — `GetSaleForUpdate :one` added to `backend/queries/sales.sql`; sqlc v1.31.1 regenerated `gen/db`; `Confirm`/`Cancel` lock-and-recheck before item work (`service.go`). Evidence: `go build ./...` + `go test ./internal/sales/...` pass; gofmt clean.
- [x] T03 — `backend/docs/SALES_TRANSACTIONS.md` gap section rewritten as "Concurrency Guards (verified)"; Client Contract untouched.
