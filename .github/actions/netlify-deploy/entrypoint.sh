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

CLEAN_NAME=$(slug $NETLIFY_SITE_NAME)

if [ -z "$(get_site_id $CLEAN_NAME)" ]; then
  echo "creating site"
  netlify sites:create \
    --account-slug ${NETLIFY_ACCOUNT_SLUG} \
    --manual \
    --name "${CLEAN_NAME}"
else
  echo "site already exists"
fi

echo "fetching site id"

SITE_ID=$(get_site_id ${CLEAN_NAME})
SITE_URL=$(get_site_url ${CLEAN_NAME})

echo
echo "site id $SITE_ID"
echo "site url $SITE_URL"
echo ::set-output name=site_id::$SITE_ID
echo ::set-output name=site_url::$SITE_URL
echo

echo "deploy site"

netlify deploy --prod --dir ${NETLIFY_DIRECTORY} --site $SITE_ID
