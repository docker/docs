#!/bin/sh

VER="$2"

if [ -z "$VER" ]; then
  echo "No version provided. Exiting."
  exit 1
else
  BASEURL="$VER/"
fi

TARGET="$1"

if ! [ -d "$TARGET" ]; then
  echo "Target directory $TARGET does not exist. Exiting."
  exit 1
fi

echo "Cleaning up $VER..."
# Fix relative links for archive
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs-stage.docker.com/#href="/#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs-stage.docker.com/#src="/#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="https://docs.docker.com/#href="/#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="https://docs.docker.com/#src="/#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="http://docs.docker.com/#href="/#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="http://docs.docker.com/#src="/#g'

# Substitute https:// for schema-less resources (src="//analytics.google.com")
# We're replacing them to prevent them being seen as absolute paths below
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="//#href="https://#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="//#src="https://#g'

# And some archive versions already have URLs starting with '/version/'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/'"$BASEURL"'#href="/#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/'"$BASEURL"'#src="/#g'

# Archived versions 1.7 and under use some absolute links, and v1.10 uses
# "relative" links to sources (href="./css/"). Remove those to make them
# work :)
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="\./#href="/#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="\./#src="/#g'

# Create permalinks for archived versions
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#href="/#href="/'"$BASEURL"'#g'
find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#src="/#src="/'"$BASEURL"'#g'

echo "Finished cleaning up $VER."
