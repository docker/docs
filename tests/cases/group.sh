#!/bin/sh
# SUMMARY: Docker for Mac and Windows regression tests
# NAME: pinata
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

# OS X settings
OSX_INSTALLER_NAME="Docker.dmg"
OSX_LOCAL_INSTALLER="$(pwd)/../../v1/mac/build/${OSX_INSTALLER_NAME}"
OSX_LOCAL_DIR="$(pwd)/../../v1/mac/build/Docker.app"
OSX_SCRIPTS="$(pwd)/../../v1/mac/scripts"

# Windows settings
WIN_INSTALLER_NAME="InstallDocker.msi"
WIN_LOCAL_INSTALLER="$(pwd)/../../win/build/${WIN_INSTALLER_NAME}"
WIN_LOCAL_DIR="$(pwd)/../../win/build/win"

# This mostly sets up the env.sh file
# Adding the PATH is a bit trickier as we support different modes.
# - If the `installer` label is set, set the PATH to the location where the
#   installer will put the files
# - If the `nostart` label is *not* set we assume we run from the local
#   build tree and set the path to that.
# - if the `nostart` label is set, we assume the user has set everything up
_osx_group_init() {
    echo "D4X_INSTALLER_NAME=${OSX_INSTALLER_NAME}" >> "${D4X_TMPDIR}/env.sh"
    echo "D4X_LOCAL_INSTALLER=${OSX_LOCAL_INSTALLER}" >> "${D4X_TMPDIR}/env.sh"
    echo "D4X_LOCAL_DIR=${OSX_LOCAL_DIR}" >> "${D4X_TMPDIR}/env.sh"
    echo "OSX_SCRIPTS=${OSX_SCRIPTS}" >> "${D4X_TMPDIR}/env.sh"
    # If D4X_LOCAL_DIR exists, use Docker.app from there
    if [ -d "$(dirname "${OSX_LOCAL_DIR}")" ]; then
      OSX_APP_DIR=${OSX_LOCAL_DIR}
    else
      OSX_APP_DIR="/Applications/Docker.app"
    fi
    echo "OSX_APP_DIR=${OSX_APP_DIR}" >> "${D4X_TMPDIR}/env.sh"
    echo "" >> "${D4X_TMPDIR}/env.sh"

    if rt_label_set "installer"; then
        echo "PATH=/usr/local/bin:\$PATH"  >> "${D4X_TMPDIR}/env.sh"
    elif rt_label_set "nostart"; then
        echo "PATH=${OSX_LOCAL_DIR}/Contents/Resources/bin:\$PATH"  >> "${D4X_TMPDIR}/env.sh"
    fi
}

_win_group_init() {
    echo "D4X_INSTALLER_NAME=${WIN_INSTALLER_NAME}" >> "${D4X_TMPDIR}/env.sh"
    echo "D4X_LOCAL_INSTALLER=${WIN_LOCAL_INSTALLER}" >> "${D4X_TMPDIR}/env.sh"
    echo "D4X_LOCAL_DIR=${WIN_LOCAL_DIR}" >> "${D4X_TMPDIR}/env.sh"
    echo "" >> "${D4X_TMPDIR}/env.sh"

    if rt_label_set "installer"; then
        echo "PATH=/c/Program\ Files/Docker/Docker/Resources/bin:\$PATH"  >> "${D4X_TMPDIR}/env.sh"
    elif rt_label_set "nostart"; then
        echo "PATH=${WIN_LOCAL_DIR}/bin:\$PATH"  >> "${D4X_TMPDIR}/env.sh"
    fi
}

# Create a version file
# This must be called *after* a version is installed
_set_version() {
    [ -z "${RT_RESULTS}" ] && return

    # If version is passed in via the environment, use it.
    if [ -n "${D4X_VERSION}" ]; then
        echo "${D4X_VERSION}" > "${RT_RESULTS}/VERSION.txt"
        return
    fi

    if rt_label_set "nostart"; then
        # If we run against something the user started we have no
        # reliable way to determine the version.
        echo "UNKNOWN" > "${RT_RESULTS}/VERSION.txt"
    elif rt_label_set "installer"; then
        # XXX TODO: On windows the binary will get a `-Version` option.
        echo "UNKNOWN" > "${RT_RESULTS}/VERSION.txt"
    else
        set +e
        git rev-parse HEAD > "${RT_RESULTS}/VERSION.txt"
        [ $? -ne 0 ] &&  echo "UNKNOWN REV" > "${RT_RESULTS}/VERSION.txt"
        set -e
    fi
}

group_init() {
    [ -r "${D4X_TMPDIR}" ] && rm -rf "${D4X_TMPDIR}"
    mkdir "${D4X_TMPDIR}"

    # Create env.sh file
    case "${RT_OS}" in
        osx)
            _osx_group_init
            ;;
        win)
            _win_group_init
            ;;
        *)
            exit 1
            ;;
    esac

    echo "export D4X_LOCAL_INSTALLER" >> "${D4X_TMPDIR}/env.sh"
    echo "export D4X_LOCAL_DIR" >> "${D4X_TMPDIR}/env.sh"

    [ -z "${D4X_PASSWORD}" ] && D4X_PASSWORD="Pass30rd!"
    [ -z "${D4X_USERNAME}" ] && D4X_USERNAME=$(whoami)
    echo "D4X_PASSWORD=${D4X_PASSWORD}" >> "${D4X_TMPDIR}/env.sh"
    echo "D4X_USERNAME=${D4X_USERNAME}" >> "${D4X_TMPDIR}/env.sh"

    source "${D4X_TMPDIR}/env.sh"

    # If the installer label is set, either use the installer from the
    # local directory or download it (if D4X_INSTALLER_URL is
    # set). Either way, the installer will be in D4X_TMPDIR afterwards.
    # shellcheck disable=SC2153
    if rt_label_set "installer"; then
        if [ -z "${D4X_INSTALLER_URL}" ]; then
            # no installer URL, check local dir
            if [ -f "${D4X_LOCAL_INSTALLER}" ]; then
                cp "${D4X_LOCAL_INSTALLER}" "${D4X_TMPDIR}"
            else
                echoerr "No local installer in ${D4X_LOCAL_INSTALLER} and D4X_INSTALLER_URL not set"
                exit 1
            fi
        else
            (cd "${D4X_TMPDIR}" || exit 1; curl --retry 5 --retry-delay 10 "${D4X_INSTALLER_URL}" -o "${D4X_INSTALLER_NAME}")
        fi

        # ensure that no docker is running nor is installed, when
        # running installer tests.
        d4x_app_stop
        d4x_app_uninstall
    fi

    return 0
}

group_deinit() {
    source "${D4X_TMPDIR}/env.sh"
    _set_version
    # get the last 4 logs on windows
    # shellcheck disable=SC2012
    if [ "${RT_OS}" = "win" ]; then
        for f in $(ls -Art "${HOME}"/AppData/Local/Docker/log* | tail -4); do
            cp "$f" "${RT_RESULTS}"
        done
        for f in $(ls -Art "${HOME}"/AppData/Local/Docker/inst* | tail -4); do
            cp "$f" "${RT_RESULTS}"
        done
    fi
    return 0
}

CMD=$1
case $CMD in
init)
    group_init
    res=$?
    ;;
deinit)
    group_deinit
    res=$?
    ;;
*)
    res=1
    ;;
esac

exit $res
