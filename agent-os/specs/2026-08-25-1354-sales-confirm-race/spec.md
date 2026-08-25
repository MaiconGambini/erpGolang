# Feature Spec — Sales confirm race guard

Date: 2026-08-25 · Lane: P1 "Sales race — concurrent confirm guard" · Source: `backend/docs/SALES_TRANSACTIONS.md` §Known Concurrency Gap

## Objective

Two concurrent `POST /api/v1/sales/{id}/confirm` calls on the same draft sale must commit
**at most one** stock decrement set, with the loser answering `409 INVALID_STATUS` — and this
must be pinned by a regression test, not by reading code.

## Current state (verified 2026-08-25, file:line)

- `Confirm` (`backend/internal/sales/service.go:293-355`) reads status **outside** the tx via
  `s.Get` (L294-300), then decrements stock per item inside the tx (L320-331), then flips
  status.
- The status flip IS guarded: `UpdateSaleStatus`
  (`backend/queries/sales.sql:45-48`) carries `AND status = @from_status`, and the service maps
  its `ErrNoRows` to `errInvalidStatus` (service.go:333-340). Any error path triggers the
  deferred rollback (L311), so the loser's decrements are undone.
- `DecrementProductStock` (`backend/queries/sales.sql:74-78`) guards with
  `AND stock >= @quantity` → `INSUFFICIENT_STOCK`.
- Conclusion: the doc's stated hole ("UPDATE does not require `status='draft'`") is **already
  closed in code**. What remains is (a) zero proof, (b) a wasteful/deadlock-prone window where
  the losing tx performs all item work before discovering it lost, (c) stale documentation.

## Requirements

- REQ-001: A concurrency regression test drives ≥2 simultaneous confirms of one seeded draft
  sale against real Postgres and asserts exactly one `200` and N−1 `409 INVALID_STATUS`, with
  final per-product stock decremented exactly once. Skips without `DATABASE_URL` (same
  convention as `backend/internal/dashboard/integration_test.go`).
- REQ-002: `Confirm` (and `Cancel`) re-checks status on a locked row **inside** the tx —
  new sqlc query `GetSaleForUpdate` (`SELECT ... FOR UPDATE`) — so the loser exits with
  `INVALID_STATUS` before touching product rows. The guarded `UPDATE ... WHERE status=...`
  stays as belt-and-braces.
- REQ-003: `backend/docs/SALES_TRANSACTIONS.md` §Known Concurrency Gap is rewritten to match
  reality: guard exists; remaining caveat is retry semantics (doc §Client Contract unchanged).

## Acceptance Criteria

- [ ] `cd backend && go test -race ./internal/sales/...` passes including the new concurrency test.
- [ ] Losing confirm never reaches `DecrementProductStock` (observable via test ordering or lock-wait assertions; minimally: test passes deterministically across repeated runs).
- [ ] Doc gap section matches implemented behavior; no contradiction between `SALES_TRANSACTIONS.md` and queries.
- [ ] No change to the HTTP contract in `contract/openapi.yaml` (status codes and error codes unchanged).
