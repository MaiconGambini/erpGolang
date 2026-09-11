# Auth And Sessions

## Login

- Body: `{ tenantSlug, email, password }`
- Tenant matched by slug → user's `tenant_id`
- Response: `{ data: { accessToken, user } }` (envelope via handler map)
- Sets HttpOnly `refresh_token` cookie

## Access Token

- JWT HS256; claims: `user_id`, `tenant_id`, `role`
- TTL from `JWT_ACCESS_TTL` (default 15 minutes)
- Sent as `Authorization: Bearer` header

## Refresh Token

- 32-byte random hex; SHA-256 hash stored in `auth_sessions`
- Rotation on refresh (old session revoked)
- Cookie name: `refresh_token`; path `/api/v1/auth`
- `HttpOnly`, `SameSite=Lax`
- `Secure=true` when `APP_ENV` is not development

## Public vs Protected

| Route | Auth |
|---|---|
| `POST /auth/login` | Public (+ rate limit) |
| `POST /auth/refresh` | Cookie |
| `POST /auth/logout` | Cookie |
| `GET /auth/me` | Bearer JWT + tenant scope |

## Rate Limiting

5 attempts / 15 minutes per IP on login. See `RATE_LIMITING.md`.

## Config Note

`JWT_REFRESH_SECRET` is required in `config.Config` but **unused** — refresh uses opaque tokens hashed in Postgres. Remove or document as reserved for future HMAC refresh tokens.

## Production Checklist

- `Secure` cookie flag enabled (non-development)
- `ALLOWED_ORIGINS` matches frontend URL
- Frontend uses `withCredentials: true` on axios client
