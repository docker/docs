#!/bin/sh

set -xe

function slug() {
  echo $1 | sed -E s/[^a-zA-Z0-9]+/-/g | sed -E s/^-+\|-+$//g | tr A-Z a-z
}

function get_site_id() {
  netlify sites:list --json | jq --raw-output ".[] | select(.name==\"$1\") | .id"
}

function get_site_url() {
  netlify sites:list --json | jq --raw-output ".[] | select(.name==\"$1\") | .url"
}

echo "searching site ${NETLIFY_SITE_NAME}"

CLEAN_NAME=$(slug $NETLIFY_SITE_NAME)

SITE_ID=$(get_site_id $CLEAN_NAME)

echo "deleting site"

netlify sites:delete --force "${SITE_ID}"
