#!/bin/sh
# SUMMARY: Run the docker-compose integration tests
# LABELS: master, release, !win
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# tox fails on Windows. See https://github.com/docker/pinata/issues/4253

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"


# install tox if not already installed
if [ "${RT_OS}" = "osx" ]; then
    which tox || sudo easy_install tox
fi
if [ "${RT_OS}" = "win" ]; then
    # On windows tox is not on the default path. Add possible
    # locations to the path
    PATH="${PATH}":/c/ProgramData/Chocolatey/lib/python2/tools/Scripts
    PATH="${PATH}":/c/ProgramData/Chocolatey/lib/python3/tools/Scripts
    which tox || "${RT_UTILS}/rt-elevate.exe" -wait easy_install tox
fi

# download source code matching our version:
# the version is the third word, followed by a ','. strip it.
VERSION=$(docker-compose --version | awk '{print $3}' | sed 's/.$//')
cd "${D4X_TMPDIR}"
curl -sL https://github.com/docker/compose/archive/"${VERSION}".zip -o "${VERSION}".zip
unzip "${VERSION}".zip
cd compose-"${VERSION}"

# compose has a test for git pre-commit hooks which fails if not
# executed in a git repository, filter it out
TOX_ENV=
for e in $(tox -l | grep '^py'); do
    TOX_ENV="${TOX_ENV}",$e
done

# fd leak in test harness in 1.8.0-rc2 consumes 800-900 descriptors.
# See https://github.com/docker/compose/pull/3748
ulimit -n 1024

# run tox and return its return code
tox -e "${TOX_ENV}" --skip-missing-interpreters
