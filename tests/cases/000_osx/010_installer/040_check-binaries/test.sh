#!/bin/sh
# SUMMARY:Validate binaries
# LABELS: !win
# AUTHOR: Ian Campbell <ian.campbell@docker.com>

# Source libraries. Uncomment if needed/defined
#. ${RT_ROOT}/lib/lib.sh
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

set -e # Exit on error

# Verify we aren't shipping artifacts which reference homebrew in /usr/local
# Verify our binaries will work on 10.10 (yosemite)

group=$(readlink /usr/local/bin/docker)

if echo "${group}" | grep -q '/var/root'; then
  echo "Hit bug #4880, sleeping for 10 secs to give vmnetd time to behave"
  sleep 10
fi

path=$(readlink "${group}")
bundle=$(dirname "$(dirname "$(dirname "$(dirname "${path}")")")")
echo bundle="${bundle}"

DIRS="Contents/MacOS Contents/Resources/bin Contents/Resources/lib"

# shellcheck disable=SC2045
# (SC2045: Iterating over ls output is fragile. Use globs. )

for dir in ${DIRS}; do
  for file in $(ls -1 "${bundle}/${dir}"); do
    path="${bundle}/${dir}/${file}"
    if otool -L "${path}" | grep /usr/local > /dev/null ; then
      echo "ERROR: ${path} is referencing homebrew libraries in /usr/local"
      res=1
    fi
    version=$(otool -l "${path}" | grep LC_VERSION_MIN_MACOS -2 | grep "version 10" | cut -d "." -f 2)
    if [ -n "${version}" ]; then
      if [ "${version}" -gt 10 ]; then
        echo "ERROR: ${path} is built for the wrong version of OSX: 10.${version} > 10.10"
        res=1
      else
        echo "INFO: Binary ${path} is built for 10.${version} <= 10.10"
      fi
    fi
  done
done

exit $res
