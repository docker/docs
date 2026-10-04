#!/bin/sh
# Restores a backup made by scripts/backup.sh into the docker-compose.prod.yml stack.
#   scripts/restore.sh backups/alshuyukh-20260101T021500Z.dump
# This REPLACES the current database. The API is stopped during the restore.
set -eu
FILE=${1:?usage: scripts/restore.sh <dump-file>}
COMPOSE="docker compose -f docker-compose.prod.yml --env-file ${ENV_FILE:-deploy/.env}"
printf 'Replace the database with %s? Type "restore" to continue: ' "$FILE"
read -r answer
[ "$answer" = "restore" ] || { echo "cancelled"; exit 1; }
$COMPOSE stop api
$COMPOSE exec -T postgres psql -U postgres -v ON_ERROR_STOP=1 -d postgres <<'SQL'
SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'alshuyukh' AND pid <> pg_backend_pid();
DROP DATABASE alshuyukh;
CREATE DATABASE alshuyukh OWNER alshuyukh_owner;
REVOKE ALL ON DATABASE alshuyukh FROM PUBLIC;
GRANT CONNECT ON DATABASE alshuyukh TO alshuyukh_app;
SQL
# Ownership and grants come from the dump: tables stay owned by alshuyukh_owner,
# so row-level security still applies to the app role.
$COMPOSE exec -T postgres pg_restore -U postgres -d alshuyukh --exit-on-error < "$FILE"
# Applies any migrations newer than the backup.
$COMPOSE run --rm migrate
$COMPOSE start api
echo "restored from $FILE"
