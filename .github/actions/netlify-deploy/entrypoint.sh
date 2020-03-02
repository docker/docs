#!/bin/sh

function usage() {
  echo "Usage: <netlify_token> <site_name>"
  exit 1
}

if [ -z "$1" ] || [ -z "$2" ]; then
  usage
fi

netlify_token=$1
site_name=$2
