#!/bin/sh
# Daily database backup for the docker-compose.prod.yml stack.
#   scripts/backup.sh [backup-dir]        (default: ./backups, keeps 14 days)
# Schedule it with cron, e.g.  15 2 * * *  cd /srv/alshuyukh && scripts/backup.sh /var/backups/alshuyukh
# Copy the files off the server too: a backup on the same disk is not a backup.
set -eu
DIR=${1:-./backups}
KEEP_DAYS=${KEEP_DAYS:-14}
COMPOSE="docker compose -f docker-compose.prod.yml --env-file ${ENV_FILE:-deploy/.env}"
mkdir -p "$DIR"
FILE="$DIR/alshuyukh-$(date -u +%Y%m%dT%H%M%SZ).dump"
# Custom format: compressed, keeps owners and grants, restorable with pg_restore.
$COMPOSE exec -T postgres pg_dump -U postgres -d alshuyukh --format=custom > "$FILE.partial"
mv "$FILE.partial" "$FILE"
# A dump that pg_restore cannot list is not a backup.
$COMPOSE exec -T postgres pg_restore --list < "$FILE" > /dev/null
find "$DIR" -name 'alshuyukh-*.dump' -mtime +"$KEEP_DAYS" -delete
echo "backup written: $FILE ($(du -h "$FILE" | cut -f1))"
