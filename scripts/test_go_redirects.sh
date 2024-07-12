#!/usr/bin/env sh

echo "Testing redirects in data/redirects.yml..."

# Get all redirects from data/redirects.yml
REDIRECTS=$(yq eval 'keys | .[]' data/redirects.yml)
LOCAL_REDIRECTS=$(echo $REDIRECTS | awk 'BEGIN { RS = " " } /^\// { print $1 }')

echo "Checking for duplicate redirects..."
DUPLICATES=$(echo $LOCAL_REDIRECTS | tr ' ' '\n' | sort | uniq -d)
if [ -n "$DUPLICATES" ]; then
  echo "Duplicate redirects found:"
  echo $DUPLICATES
  exit 1
fi

echo "Checking for redirects to nonexistent paths..."
for file in $(echo REDIRECTS | awk 'BEGIN { RS = " " } /^\// { print $1 }'); do
  if [ ! -e "./public/${file%/*}/index.html" ]; then
    echo "Redirect to nonexistent path: $file"
    exit 1
  fi
done

echo "All redirects are valid."
