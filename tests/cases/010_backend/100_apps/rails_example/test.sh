#!/bin/sh
# SUMMARY: rails docker-compose example
# LABELS: master, release, nightly, !win, apps
# The docker-compose run line errors on Windows
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# Based on: https://docs.docker.com/compose/rails/

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rm railsexample_web_run_1 || true
    docker-compose down --rmi all || true
    docker rmi "$(docker images -f "dangling=true" -q)" || true

    rm -rf Gemfile.lock Gemfilee README.rdoc Rakefile app bin config config.ru
    rm -rf db lib log public test tmp vendor Gemfile
}
trap clean_up EXIT

# copy files in place
touch Gemfile.lock
cp Gemfile.new Gemfile

# Prepare
docker-compose run --rm web rails new . --force --database=postgresql --skip-bundle

sed -ie "s/# gem 'therubyracer/gem 'therubyracer/" Gemfile

docker-compose build --force-rm

cp database.yml config

docker-compose up -d
docker-compose run web rake db:create
sleep 5

"${RT_UTILS}/rt-urltest" -r 10 -s aboard "http://${D4X_HOST_NAME}:3000"
exit 0
