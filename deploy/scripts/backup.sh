#!/usr/bin/env bash
#
# goERP Postgres backup — streams a custom-format pg_dump out of the prod
# compose stack and verifies it before calling it a backup.
#
# Usage:
#   ./backup.sh <backup-dir> <retention-days>
#
# Environment:
#   GOERP_ENV_FILE   required by the prod compose file (e.g. /etc/goerp/goerp.env)
#   GOERP_COMPOSE_FILE  override compose file (default deploy/compose/docker-compose.prod.yml)
#
# Recommended cron tiers (7/14/30 retention, per DEPLOYMENT.md):
#   15 * * * *  /opt/goerp/deploy/scripts/backup.sh /var/backups/goerp/hourly 1
#   30 2 * * *  /opt/goerp/deploy/scripts/backup.sh /var/backups/goerp/daily 14
#   45 3 * * 0  /opt/goerp/deploy/scripts/backup.sh /var/backups/goerp/weekly 30
#
# Ship the directories off-server (rsync/restic) — this script does not upload.

set -euo pipefail

DIR="${1:?usage: backup.sh <backup-dir> <retention-days>}"
RETENTION_DAYS="${2:?usage: backup.sh <backup-dir> <retention-days>}"
COMPOSE_FILE="${GOERP_COMPOSE_FILE:-deploy/compose/docker-compose.prod.yml}"
PG_USER="${POSTGRES_USER:-goerp}"
PG_DB="${POSTGRES_DB:-goerp}"

compose() {
	docker compose -f "$COMPOSE_FILE" "$@"
}

command -v docker >/dev/null || { echo "docker CLI not found" >&2; exit 1; }
mkdir -p "$DIR"

stamp="$(date +%Y%m%d-%H%M%S)"
out="$DIR/goerp-$stamp.dump"
tmp="$out.tmp"

echo "[backup] dumping $PG_DB to $tmp"
compose exec -T postgres pg_dump -U "$PG_USER" -d "$PG_DB" -Fc > "$tmp"

[ -s "$tmp" ] || { echo "[backup] ERROR: empty dump" >&2; rm -f "$tmp"; exit 1; }

# Integrity gate: the archive must at least be listable by pg_restore.
if [ "${VERIFY:-1}" = "1" ]; then
	echo "[backup] verifying archive"
	compose cp "$tmp" postgres:/tmp/.goerp-verify.dump
	if ! compose exec -T postgres pg_restore --list /tmp/.goerp-verify.dump >/dev/null; then
		compose exec postgres rm -f /tmp/.goerp-verify.dump || true
		echo "[backup] ERROR: archive failed pg_restore --list" >&2
		rm -f "$tmp"
		exit 1
	fi
	compose exec postgres rm -f /tmp/.goerp-verify.dump
fi

mv "$tmp" "$out"
echo "[backup] wrote $out ($(du -h "$out" | cut -f1))"

pruned=$(find "$DIR" -name 'goerp-*.dump' -mtime "+$RETENTION_DAYS" -print -delete | wc -l)
echo "[backup] pruned $pruned dump(s) older than ${RETENTION_DAYS}d from $DIR"
