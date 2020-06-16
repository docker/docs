#!/bin/sh

set -xe

function usage() {
  echo "Usage: <netlify_auth_token> <site_name>"
  exit 1
}

function slug() {
  echo $1 | sed -E s/[^a-zA-Z0-9]+/-/g | sed -E s/^-+\|-+$//g | tr A-Z a-z
}

function get_site_id() {
  netlify sites:list --json | jq --raw-output ".[] | select(.name==\"$1\") | .id"
}

if [ -z "$1" ] || [ -z "$2" ]; then
  usage
fi

export NETLIFY_AUTH_TOKEN=$1
site_name=$2


echo "searching site ${site_name}"

clean_site_name=$(slug $site_name)
site_id=$(get_site_id $clean_site_name)

echo "deleting site"

netlify sites:delete --force "${site_id}"
