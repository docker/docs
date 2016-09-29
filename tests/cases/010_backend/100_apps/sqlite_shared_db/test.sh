#!/bin/sh
# SUMMARY: test sql3 on shared volume database
# LABELS:
# AUTHOR: Emmanuel Briney <emmanuel.briney@docker.com>
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=sqlite_shared_db

clean_up() {
    rm -f "${D4X_LOCAL_TMPDIR}/database.db" || true
    docker rmi "${IMAGE_NAME}" || true
}
trap clean_up EXIT

clean_up
docker build -t "${IMAGE_NAME}" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}":/data "${IMAGE_NAME}"
exit 0
