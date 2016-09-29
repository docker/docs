#!/bin/sh
# SUMMARY: django docker-compose example
# LABELS: master, release, nightly, !win, apps
# The docker-compose run line errors on Windows
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# Based on https://docs.docker.com/compose/django/

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
# clean up
    docker-compose down --rmi all || true
    rm -rf composeexample
    rm -f manage.py
}
trap clean_up EXIT

# initialise
docker-compose run --rm web django-admin.py startproject composeexample .
sed -ie '/^DATABASES/{n;N;N;N;d;}' composeexample/settings.py
sed -ie '/^DATABASES/r db.txt'     composeexample/settings.py
rm composeexample/settings.pye

docker-compose up -d
sleep 5
"${RT_UTILS}/rt-urltest" -r 10 -s Congratulations "http://${D4X_HOST_NAME}:8000"

exit 0
