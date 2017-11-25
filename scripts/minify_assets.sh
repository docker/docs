#!/bin/sh

TARGET="$1"
VER="$2"

if ! [ -d "$TARGET" ]; then
  echo "Target directory $TARGET does not exist. Exiting."
  exit 1
fi

# Minify assets. This benefits both the compressed, and uncompressed versions
printf "Optimizing "

printf "html..."; minify -r --type=html --match=\.html -o ${TARGET}/ ${TARGET} || true
printf "css..." ; minify -r --type=css  --match=\.css  -o ${TARGET}/ ${TARGET} || true
printf "json..."; minify -r --type=json --match=\.json -o ${TARGET}/ ${TARGET} || true

echo "done."
