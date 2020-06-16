#!/bin/sh

set -xe

function usage() {
  echo "Usage: <netlify_auth_token> <netlify_account_slug> <site_name> <directory>"
  exit 1
}

function slug() {
  echo $1 | sed -E s/[^a-zA-Z0-9]+/-/g | sed -E s/^-+\|-+$//g | tr A-Z a-z
}

function get_site_id() {
  netlify sites:list --json | jq --raw-output ".[] | select(.name==\"$1\") | .id"
}

function get_site_url() {
  netlify sites:list --json | jq --raw-output ".[] | select(.name==\"$1\") | .url"
}

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ] || [ -z "$4" ]; then
  usage
fi

export NETLIFY_AUTH_TOKEN=$1
account_slug=$2
site_name=$3
directory=$4

clean_site_name=$(slug $site_name)

if [ -z "$(get_site_id $clean_site_name)" ]; then
  echo "creating site"
  netlify sites:create \
    --account-slug ${account_slug} \
    --manual \
    --name "${clean_site_name}"
else
  echo "site already exists"
fi

echo "fetching site id"

site_id=$(get_site_id ${clean_site_name})
site_url=$(get_site_url ${clean_site_name})

echo
echo "site id $site_id"
echo "site url $site_url"
echo ::set-output name=site_id::$site_id
echo ::set-output name=site_url::$site_url
echo

echo "deploy site"

netlify deploy --prod --dir ${directory} --site $site_id
