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

## Known Concurrency Gap

`Confirm` reads status outside the transaction; `UPDATE sales SET status` does not require `status = 'draft'`. Two concurrent confirms could both decrement stock.

**Target fix:** `SELECT ... FOR UPDATE` inside tx; `UPDATE ... WHERE status = 'draft'`.

## Audit

Audit recorded post-commit. Audit failure does not roll back the sale (best-effort).
