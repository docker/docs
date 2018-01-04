#!/bin/sh

TARGET="$1"
VER="$2"

if [ -z "$TARGET" ]; then
  echo "Usage: $0 <target> <version>"
  echo "No target provided. Exiting."
  exit 1
fi

if [ -z "$VER" ]; then
  echo "Usage: $0 <target> <version>"
  echo "No version provided. Exiting."
  exit 1
else
  BASEURL="$VER/"
fi

if ! [ -d "$TARGET" ]; then
  echo "Target directory $TARGET does not exist. Exiting."
  exit 1
fi

# Create permalinks for archived versions
# first the HTML, then JS which is used to build the TOC in 17.09+
# Note: pattern '\(src\|href\)=\("\{0,1\}\)' matches:
# - src=
# - href=
# followed by an optional double quote
# the pattern for the JS regex matches exactly: "path":"<absolute path>"
# The goal is to change URLs like /blah to /v17.09/blah
printf "Creating permalinks for $VER"
printf ".HTML..."; find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#\(src\|href\)=\("\{0,1\}\)/#\1=\2/'"$BASEURL"'#g';
printf ".JS..."; find ${TARGET} -type f -name 'toc.js' -print0 | xargs -0 sed -i 's#"path":"/#"path":"/'"$BASEURL"'#g';
echo "done"