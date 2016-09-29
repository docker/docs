#!/bin/sh
# SUMMARY: docker-compose elastic search test
# LABELS: master, release, nightly, apps
# AUTHOR: Ian Campbell <ian.campbell@docker.com>
# Partially covers #2017: File-sharing compatibility/correctness testing for
# common apps.

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up () {
    docker-compose down --rmi all || true
}

trap clean_up EXIT

docker-compose pull
docker-compose up -d

for retries in $(seq 30) ; do
    if curl -fsS 'http://localhost:9200' ; then
	break
    fi
    echo "Waiting 1s"
    sleep 1s
done
if [ "$retries" -eq 30 ] ; then
    echo "Service not up after $retries pings"
    exit 1
else
    echo "Service up after $retries pings"
fi

# Simple workflow from https://github.com/elastic/elasticsearch

# _version and created fields change over repeated runs, normalise
normalise_result()
{
    # OSX sed doesn't support alternatives ("|") without -E and on
    # Windows GNU sed doesn't support -E, so just spell out the
    # alternatives for "created".
    sed -e 's/"_version":[0-9][0-9]*/"_version":N/g; s/"created":true/"created":BOOL/g; s/"created":false/"created":BOOL/g'
}

check_result()
{
    if [ "x$expected" != "x$actual" ] ; then
	echo "Unexpected result"
	echo "Expected: $expected"
	echo "Actual:   $actual"
	exit 1
    fi
}

expected='{"_index":"twitter","_type":"user","_id":"kimchy","_version":N,"created":BOOL}'
actual=$(curl -fsS -XPUT 'http://localhost:9200/twitter/user/kimchy' -d '{ "name" : "Shay Banon" }' | normalise_result)
check_result

expected='{"_index":"twitter","_type":"tweet","_id":"1","_version":N,"created":BOOL}'
actual=$(curl -fsS -XPUT 'http://localhost:9200/twitter/tweet/1' -d '
{
    "user": "kimchy",
    "post_date": "2009-11-15T13:12:00",
    "message": "Trying out Elasticsearch, so far so good?"
}' | normalise_result)
check_result

expected='{"_index":"twitter","_type":"tweet","_id":"2","_version":N,"created":BOOL}'
actual=$(curl -fsS -XPUT 'http://localhost:9200/twitter/tweet/2' -d '
{
    "user": "kimchy",
    "post_date": "2009-11-15T14:12:12",
    "message": "Another tweet, will it be indexed?"
}' | normalise_result)
check_result

# This produces quite a bit of output, don't bother checking the content, just that it succeeded.
curl -fsS -XGET 'http://localhost:9200/twitter/tweet/_search?q=user:kimchy&pretty=true'


docker-compose stop
docker-compose down --rmi all

exit 0
