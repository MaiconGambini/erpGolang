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

## Backup & Restore Drill

Automated by `deploy/scripts/backup.sh` (custom-format `pg_dump` streamed from the
compose `postgres` service, verified via `pg_restore --list`, retention sweep built in):

```bash
export GOERP_ENV_FILE=/etc/goerp/goerp.env
/opt/goerp/deploy/scripts/backup.sh /var/backups/goerp/daily 14
```

Cron tiers (7/14/30 retention policy):

```text
15 * * * * /opt/goerp/deploy/scripts/backup.sh /var/backups/goerp/hourly 1
30 2 * * * /opt/goerp/deploy/scripts/backup.sh /var/backups/goerp/daily   14
45 3 * * 0 /opt/goerp/deploy/scripts/backup.sh /var/backups/goerp/weekly 30
```

Ship `/var/backups/goerp/` off-server (rsync/restic) — the script does not upload.

**Restore drill** (run quarterly, and before declaring production-ready) using
`deploy/scripts/restore.sh` — defaults to a scratch database so it can never clobber live data:

```bash
CONFIRM=YES /opt/goerp/deploy/scripts/restore.sh \
  /var/backups/goerp/daily/goerp-YYYYMMDD-HHMMSS.dump goerp_restore_check

# pass criterion: key tables present, counts plausible vs production
docker compose -f deploy/compose/docker-compose.prod.yml exec postgres \
  psql -U goerp -d goerp_restore_check -c '\dt'
```

Record drill date + result in this file below.

### Drill log

- (no drills yet — first drill pending a running Postgres)

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
