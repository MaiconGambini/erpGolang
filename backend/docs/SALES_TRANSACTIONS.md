# Sales Transactions

## Invariants

1. Draft sales do not affect stock.
2. Confirm is atomic: status transition + all line stock decrements in one transaction.
3. Cancel is allowed only from `confirmed`; restores stock.
4. Soft-delete only for `draft` sales.

## Status Machine

| From | Action | To |
|---|---|---|
| `draft` | confirm | `confirmed` |
| `confirmed` | cancel | `cancelled` |

Editable/deletable only in `draft`.

## Error Codes

| Code | HTTP | When |
|---|---|---|
| `INVALID_STATUS` | 409 | Wrong status for action |
| `INSUFFICIENT_STOCK` | 409 | Confirm with insufficient product stock |
| `VALIDATION_ERROR` | 400 | Invalid input |

## Client Contract

- `POST /sales/{id}/confirm` is **not** safely retryable until idempotency is added.
- Do not blind-retry confirm after timeout.

## Concurrency Guards (verified 2026-08-25)

- Status transitions are conditional updates: `UpdateSaleStatus`
  (`backend/queries/sales.sql`) requires `AND status = @from_status`; a losing contender gets
  `ErrNoRows` -> `INVALID_STATUS` (409) and its whole transaction rolls back.
- `Confirm`/`Cancel` additionally take `SELECT ... FOR UPDATE` on the sale row at transaction
  start (`GetSaleForUpdate`) and re-check status inside the tx, so the loser exits before
  touching any product row.
- Stock guards: decrement requires `stock >= quantity` (`INSUFFICIENT_STOCK`); increment is unguarded by design.
- Regression test: `TestConcurrentConfirmSingleDecrement` in
  `backend/internal/sales/integration_test.go` — run with
  `go test -tags=integration ./internal/sales/...` against a live Postgres; skips otherwise.

The former "Known Concurrency Gap" (status flip without `status = 'draft'` predicate) is
closed by the combination above. Confirm/cancel remain **not idempotent** — §Client Contract
still applies.

## Audit

Audit recorded post-commit. Audit failure does not roll back the sale (best-effort).
