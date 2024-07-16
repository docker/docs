#!/usr/bin/env sh

# Find all media files {svg,png,webp,mp4,jpg,jpeg} in {content,static}
MEDIA=$(fd . -e "svg" -e "png" -e "webp" -e "mp4" -e "jpg" -e "jpeg" ./content ./static)
TEMPFILE=$(mktemp)

for file in $MEDIA; do
  rg -q "$(basename $file)"
  if [ $? -ne 0 ]; then
    echo "$file" >> "$TEMPFILE"
  fi
done

UNUSED_FILES=$(< $TEMPFILE)
rm $TEMPFILE

if [ -z "$UNUSED_FILES" ]; then
  exit 0
else
  echo "$(echo "$UNUSED_FILES" | wc -l) unused media files. Please remove them."
  printf "%s\n" ${UNUSED_FILES[@]}
  exit 1
fi
