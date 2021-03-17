#!/bin/sh

# Do some sanity-checking to make sure we are running this from the right place
if ! [ -f _config.yml ]; then
	echo "Could not find _config.yml. We may not be in the right place. Bailing."
	exit 1
fi

# Helper function to deal with sed differences between osx and Linux
# See https://stackoverflow.com/a/38595160
sedi () {
	sed --version >/dev/null 2>&1 && sed -i -- "$@" || sed -i "" "$@"
}

# Parse latest_engine_api_version variables from _config.yml to replace the value
# in toc.yaml. This is brittle!
latest_engine_api_version="$(grep 'latest_engine_api_version:' ./_config.yml | grep -oh '"[0-9].*"$' | sed 's/"//g')"
sedi "s/{{ site.latest_engine_api_version }}/${latest_engine_api_version}/g" ./_data/toc.yaml
