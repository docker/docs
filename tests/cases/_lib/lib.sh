#!/bin/sh
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

# Shell library functions which may be useful for writing tests

# Source the main regression test library if present
[ -f "${RT_ROOT}/lib/lib.sh" ] && . "${RT_ROOT}/lib/lib.sh"


# Temporary directory for tests to use. This is different to
# D4X_LOCAL_TMPDIR defined below as the latter *must* be on the local
# filesystem while D4X_TMPDIR might be on a shared drive.
D4X_TMPDIR="${RT_PROJECT_ROOT}/_tmp"


# The top-level group.sh of the project creates a env.sh file
# containing environment variables for tests. Source it if present.
[ -f "${D4X_TMPDIR}/env.sh" ] && . "${D4X_TMPDIR}/env.sh"

# /tmp is handled completely different to other file paths by
# Cygwin/MINGW on Windows. It's actually transparently mapped to
# <user>\AppData\Local\Temp, not \c\tmp or other places one might
# expect.  Further, it is not translated propoerly with/without
# MSYS_NO_PATHCONV set. It's a mess.
#
# Since we use /tmp a lot on unit tests, we define an environment
# variable (D4X_TMPDIR) which points to a good place for temporary
# files both on Windows and other platforms. The intention is that
# it's usable for volume mounts, host and VM FS access.
# shellcheck disable=SC2034
case "${RT_OS}" in
    osx)
        D4X_LOCAL_TMPDIR=/tmp
        ;;
    win)
        D4X_LOCAL_TMPDIR="${HOME}/AppData/Local/Temp"
        ;;
esac

# define the host name of the Docker VM
# shellcheck disable=SC2034
D4X_HOST_NAME=localhost

# Kill containers for a given image and remove that image
d4x_cleanup_image() {
    for IMAGE_NAME in "$@"; do
        docker ps -qa -f ancestor="${IMAGE_NAME}" | xargs docker rm -f || true
        docker ps -qa -f name="${IMAGE_NAME}" | xargs docker rm -f || true
        docker rmi "${IMAGE_NAME}" || true
    done

    # remove intermediate images
    docker rmi "$(docker images -f "dangling=true" -q)" || true
}

# Include Windows and OS X specific library files
case "${RT_OS}" in
    osx)
        [ -f "${RT_PROJECT_ROOT}/_lib/sh/lib_osx.sh" ] && . "${RT_PROJECT_ROOT}/_lib/sh/lib_osx.sh"
        ;;
    win)
        [ -f "${RT_PROJECT_ROOT}/_lib/sh/lib_win.sh" ] && . "${RT_PROJECT_ROOT}/_lib/sh/lib_win.sh"
        ;;
esac
