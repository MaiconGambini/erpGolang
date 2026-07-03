# Deployment

The first production path for goERP is intentionally simple: VPS + Docker Compose + Caddy. AWS should begin with EC2 + RDS only after the MVP works locally and on a VPS.

## Local Development

Local Docker Compose runs infrastructure only:

- PostgreSQL 16
- Redis 7

Backend and frontend run on the host for fast hot reload.

```bash
docker compose -f docker-compose.dev.yml up -d
```

## VPS MVP Topology

```text
Internet
  -> Caddy TLS reverse proxy
  -> frontend static assets
  -> backend container on :8080
  -> PostgreSQL container or managed database
  -> Redis container
```

Only Caddy should be public. PostgreSQL and Redis must not be exposed to the internet.

## Production Environment Variables

- `APP_ENV`
- `HTTP_ADDR`
- `DATABASE_URL`
- `REDIS_URL`
- `ALLOWED_ORIGINS`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`
- `JWT_ACCESS_TTL` (optional)
- `JWT_REFRESH_TTL` (optional)
- `BCRYPT_COST` (optional)
- `POSTGRES_DB` / `POSTGRES_USER` / `POSTGRES_PASSWORD` (compose postgres service)
- `GOERP_DOMAIN` (Caddy TLS)

Secrets must live outside Git. On a VPS, use a protected environment file such as `/etc/goerp/goerp.env` with restrictive permissions.

Production Compose requires `GOERP_ENV_FILE` to point to that real environment file and `GOERP_DOMAIN` to be exported for Caddy configuration. Do not run production Compose with the example file except for local configuration validation.

Example VPS command:

```bash
export GOERP_ENV_FILE=/etc/goerp/goerp.env
export GOERP_DOMAIN=app.example.com
docker compose -f deploy/compose/docker-compose.prod.yml config
docker compose -f deploy/compose/docker-compose.prod.yml up -d --build
```

Local validation with the committed example file:

```bash
GOERP_ENV_FILE=../../deploy/env/production.env.example GOERP_DOMAIN=example.com docker compose -f deploy/compose/docker-compose.prod.yml config
```

## Migration Strategy

Do not run migrations automatically on backend boot in production. The deploy process should run migrations as a separate step before rolling out the new backend.

Minimum sequence:

1. Build images.
2. Run migration job.
3. If migration succeeds, restart services.
4. Run health checks.
5. Abort rollout if migration or health checks fail.

## Backup

For VPS PostgreSQL:

- daily `pg_dump`
- off-server storage
- 7/14/30 day retention
- restore tested before production use

For AWS, prefer RDS automated backups and snapshots.

## AWS Path

Recommended first AWS setup:

- EC2 for the app host.
- RDS PostgreSQL.
- Redis container initially, ElastiCache later if Redis becomes critical.
- Caddy on EC2, or ALB + ACM later.
- SSM Parameter Store or Secrets Manager for secrets.

Avoid Kubernetes, ECS/Fargate, Terraform, blue/green deploys, and multi-region setups until the MVP has real operational demand.

## Fly.io Path

For a managed deploy without managing a VPS:

```bash
# One-time setup
fly apps create goerp-api
fly secrets set \
  DATABASE_URL=postgres://... \
  REDIS_URL=redis://... \
  JWT_ACCESS_SECRET=... \
  JWT_REFRESH_SECRET=... \
  ALLOWED_ORIGINS=https://goerp-api.fly.dev

# Deploy (runs /goerp-migrate as release_command)
fly deploy --config fly.toml
```

Use external managed Postgres (Neon) and Redis (Upstash) for the MVP. The `backend/Dockerfile.fly` image includes both the API and migrate binaries.

## Production Compose Validation

```bash
make prod-config
```

Deploy sequence with Compose:

```bash
export GOERP_ENV_FILE=/etc/goerp/goerp.env
export GOERP_DOMAIN=app.example.com
docker compose -f deploy/compose/docker-compose.prod.yml up -d --build
```

The `migrate` service runs before `backend` starts.
