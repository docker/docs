#!/bin/sh

TARGET="$1"

if [ -z "$TARGET" ]; then
  echo "Usage: $0 <target>"
  echo "No target was given. Exiting."
  exit 1
fi

if ! [ -d "$TARGET" ]; then
  echo "Target directory $TARGET does not exist. Exiting."
  exit 1
fi

# Pre-gzip files. note that the ngx_http_gzip_static_module requires  both the
# compressed, and uncompressed files to be present see:
# http://nginx.org/en/docs/http/ngx_http_gzip_static_module.html
#
# Compressed content is roughly 80% smaller than uncompressed but will make the
# final image 20% bigger (due to both uncompressed and compressed content being
# included in the image)
printf "Compressing assets in $TARGET"
printf "."; find ${TARGET} -type f -iname "*.html" -exec gzip -f -9 --keep {} +
printf "."; find ${TARGET} -type f -iname "*.js"   -exec gzip -f -9 --keep {} +
printf "."; find ${TARGET} -type f -iname "*.css"  -exec gzip -f -9 --keep {} +
printf "."; find ${TARGET} -type f -iname "*.json" -exec gzip -f -9 --keep {} +
printf "."; find ${TARGET} -type f -iname "*.svg"  -exec gzip -f -9 --keep {} +
printf "."; find ${TARGET} -type f -iname "*.txt"  -exec gzip -f -9 --keep {} +
echo "done."