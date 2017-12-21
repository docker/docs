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
printf "Creating permalinks for $VER"
printf "."; find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#\(src\|href\)=\("\{0,1\}\)/#\1=\2/'"$BASEURL"'#g';
echo "done"