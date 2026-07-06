# Rate Limiting

## Login Endpoint

| Property | Value |
|---|---|
| Route | `POST /api/v1/auth/login` |
| Limit | 5 attempts / 15 minutes |
| Key | `ratelimit:login:{client_ip}` |
| Response | `429 { code: "RATE_LIMITED" }` |

Implementation: `internal/platform/middleware/ratelimit.go`

## Redis Failure Behavior

- `redis == nil` or `Incr` error → **fail-open** (request proceeds)
- Documented tradeoff: availability over abuse resistance during Redis outage

## Gaps (MVP)

- No rate limit on `POST /auth/refresh`
- `INCR` + `EXPIRE` not atomic (key may lack TTL on crash between ops)
- Behind reverse proxy: ensure real client IP is used (Caddy `X-Forwarded-For` / trusted headers)

## ReadyZ Coupling

`/readyz` fails if Redis is unreachable, blocking deploy even though CRUD works without rate limiting.

Future option: split `redis_ready` warning from hard readyz failure.
