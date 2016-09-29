#!/bin/sh
# SUMMARY: Wordpress docker-compose example
# LABELS: release, master, nightly, apps
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# Based on: https://docs.docker.com/compose/wordpress/

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

WORDPRESS_TEST_DIR="${D4X_LOCAL_TMPDIR}/wordpress_test"
WORDPRESS="https://en-gb.wordpress.org/wordpress-4.5.3-en_GB.tar.gz"
TESTDIR=$(pwd)

clean_up() {
    cd "${WORDPRESS_TEST_DIR}"
    docker-compose down --rmi all || true
    rm -rf "${WORDPRESS_TEST_DIR}"
}
trap clean_up EXIT

rm -rf "${WORDPRESS_TEST_DIR}" | true
mkdir "${WORDPRESS_TEST_DIR}"
cd "${WORDPRESS_TEST_DIR}"

# The tutorial suggest downloading wordpress and copy the config file
curl --retry 5 --retry-delay 10 "${WORDPRESS}" -o wordpress.tar.gz
tar -xvf wordpress.tar.gz
rm wordpress.tar.gz

cp "${TESTDIR}/Dockerfile" "${TESTDIR}/docker-compose.yml" .
cp "${TESTDIR}/wp-config.php" ./wordpress

# start
docker-compose pull
docker-compose build
docker-compose --verbose up -d

"${RT_UTILS}/rt-urltest" -r 20 -s WordPress "http://${D4X_HOST_NAME}:8000"

exit 0
