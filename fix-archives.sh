#!/bin/sh

TARGET="$1"

if [ -z "$TARGET" ]; then
  echo "Usage: $0 <target> <version>"
  echo "No target provided. Exiting."
  exit 1
fi

VER="$2"

if [ -z "$VER" ]; then
  echo "Usage: $0 <target> <version>"
  echo "No version provided. Exiting."
  exit 1
fi


if ! [ -d "$TARGET/$VER" ]; then
  echo "Target directory $TARGET/$VER does not exist. Exiting."
  exit 1
fi


echo "Doing extra processing for archive in $TARGET/$VER:"

echo "  Fixing links..."

sh /usr/bin/normalize-links.sh "$TARGET" "$VER"

echo "  Minifying assets..."

sh /usr/bin/minify-assets.sh "$TARGET" "$VER"

echo "  Creating permalinks..."

sh /usr/bin/create-permalinks.sh "$TARGET" "$VER"

echo "  Compressing assets..."

sh /usr/bin/compress-assets.sh "$TARGET"

echo "Finished cleaning up $TARGET/$VER."