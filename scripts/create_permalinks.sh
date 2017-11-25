#!/bin/sh

TARGET="$1"
VER="$2"
BASEURL="$VER/"

if [ -z "$VER" ]; then
  echo "No version provided. Exiting."
  exit 1
fi

if ! [ -d "$TARGET" ]; then
  echo "Target directory $TARGET does not exist. Exiting."
  exit 1
fi

# Create permalinks for archived versions
printf "Creating permalinks for $VER"
printf "."; find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#\(src\|href\)=\("\{0,1\}\)/#\1=\2/'"$BASEURL"'#g';
echo "done"
