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

# Note: pattern '\(src\|href\)=\("\{0,1\}\)' matches:
# - src=
# - href=
# followed by an optional double quote

printf "Cleaning up $VER"

# Fix relative links for archive
printf "."; find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#\(src\|href\)=\("\{0,1\}\)\(http\|https\)://\(docs\|docs-stage\).docker.com/#\1=\2/#g'

# Substitute https:// for schema-less resources (src="//analytics.google.com")
# We're replacing them to prevent them being seen as absolute paths below
printf "."; find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#\(src\|href\)=\("\{0,1\}\)//#\1="https://#g'

# And some archive versions already have URLs starting with '/version/'
printf "."; find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#\(src\|href\)=\("\{0,1\}\)/'"$BASEURL"'#\1="/#g'

case "$VER" in v1.4|v1.5|v1.6|v1.7|v1.10)
	# Archived versions 1.7 and under use some absolute links, and v1.10 uses
	# "relative" links to sources (href="./css/"). Remove those to make them
	# work :)
	printf "."; find ${TARGET} -type f -name '*.html' -print0 | xargs -0 sed -i 's#\(src\|href\)=\("\{0,1\}\)\./#\1="/#g'
	;;
esac

echo "done"
