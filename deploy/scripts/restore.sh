#!/usr/bin/env bash
#
# goERP Postgres restore / drill — loads a backup.sh dump into a scratch
# database inside the prod compose stack WITHOUT touching the live one.
#
# Usage:
#   CONFIRM=YES ./restore.sh <dump-file> [target-db]   # non-interactive
#   ./restore.sh <dump-file> [target-db]               # prompts first
#
# Defaults to target db `goerp_restore_check` so a drill can never clobber
# production by accident; pass the real db name ONLY for an actual incident.
#
# Environment:
#   GOERP_ENV_FILE      required by the prod compose file
#   GOERP_COMPOSE_FILE  override (default deploy/compose/docker-compose.prod.yml)


set -euo pipefail

DUMP="${1:?usage: restore.sh <dump-file> [target-db]}"
TARGET_DB="${2:-goerp_restore_check}"
COMPOSE_FILE="${GOERP_COMPOSE_FILE:-deploy/compose/docker-compose.prod.yml}"
PG_USER="${POSTGRES_USER:-goerp}"

compose() {
	docker compose -f "$COMPOSE_FILE" "$@"
}

# Windows-native form of a host path for the docker CLI (identity on Linux).
host_path() {
	cygpath -m "$1" 2>/dev/null || printf '%s' "$1"
}

[ -f "$DUMP" ] || { echo "no such dump: $DUMP" >&2; exit 1; }
live_db="$(compose exec -T postgres printenv POSTGRES_DB 2>/dev/null || true)"
if [ -n "$live_db" ] && [ "$TARGET_DB" = "$live_db" ] && [ "${OVERRIDE_LIVE_GUARD:-}" != "YES" ]; then
	echo "refusing: target looks like the live database; use a scratch name or OVERRIDE_LIVE_GUARD=YES" >&2
	exit 1
fi
if [ "${CONFIRM:-}" != "YES" ]; then
	read -r -p "Restore $DUMP into database '$TARGET_DB'? Type YES: " answer
	[ "$answer" = "YES" ] || { echo "aborted"; exit 1; }
fi

echo "[restore] copying archive into container"
compose cp "$(host_path "$DUMP")" postgres:/tmp/.goerp-restore.dump

echo "[restore] recreating database $TARGET_DB"
compose exec -T postgres psql -U "$PG_USER" -d postgres -v ON_ERROR_STOP=1 \
	-c "DROP DATABASE IF EXISTS \"$TARGET_DB\";" \
	-c "CREATE DATABASE \"$TARGET_DB\" OWNER \"$PG_USER\";"

echo "[restore] running pg_restore"
MSYS_NO_PATHCONV=1 compose exec -T postgres pg_restore -U "$PG_USER" --no-owner --role "$PG_USER" \
	-d "$TARGET_DB" /tmp/.goerp-restore.dump

MSYS_NO_PATHCONV=1 compose exec postgres rm -f /tmp/.goerp-restore.dump

echo "[restore] OK — sanity-check row counts:"
echo "  $COMPOSE_FILE exec postgres psql -U $PG_USER -d $TARGET_DB -c '\\dt'"
echo "Drill pass criterion: key tables present and counts plausible vs production (see DEPLOYMENT.md §Backup & Restore Drill)."
