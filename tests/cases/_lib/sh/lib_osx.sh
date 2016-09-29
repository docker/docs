#!/bin/sh
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

##
## Common functions which must be implemented on OSX and Windows
##

# Return true if Docker for Mac is installed
d4x_app_installed() {
    if [ ! -d "${OSX_APP_DIR}" ]; then
        echo "Application not found in ${OSX_APP_DIR}"
        return 1
    fi
}

# Return true if Docker for Mac is running
d4x_app_running() {
    pgrep -xq com.docker.osx.hyperkit.linux && return 0 || return 1
}


# Install Docker for Mac
d4x_app_install() {
    echo "### Installing ${OSX_APP_DIR}..."
    sudo -n id || echo "WARNING: sudo requires password"

    sudo "${OSX_APP_DIR}/Contents/MacOS/Docker" --quit-after-install --unattended

    for f in $(seq 1 60); do
        d4x_app_running || break
        if [ "$f" == "60" ]; then
            echo "Timed out while waiting for Docker to exit after --quit-after-install"
            return 1
        fi
        echo "Waiting for Docker to exit, sleeping 1s..."
        sleep 1
    done

    echo "### Starting Docker.app"
    open "${OSX_APP_DIR}" --args --unattended
    d4x_wait_for_docker
}

# Uninstall Docker for Mac
d4x_app_uninstall() {
    if [ -d "${OSX_APP_DIR}" ]; then
      "${OSX_APP_DIR}/Contents/MacOS/Docker" --uninstall --unattended
      for f in $(seq 1 10); do
          echo "Waiting for docker symlink to disappear..."
          if [ ! -L /usr/local/bin/docker ]; then
              break
          fi
          if [ "$f" == "10" ]; then
              echo "timed out while waiting for docker symlinks to disappear"
              return 1
          fi
          sleep 2
      done;

      for f in $(seq 1 10); do
          echo "Waiting Docker processes to stop"
          d4x_app_running || break
          if [ "$f" == "10" ]; then
              echo "timed out while waiting for docker processes to stop"
              return 1
          fi
          sleep 2
      done
    fi
}


# Start the Docker application (with all its trimmings)
d4x_app_start() {
    open -a "$OSX_APP_DIR"
    d4x_wait_for_docker
}

# Start the backend only
d4x_backend_start() {
    # XXX TODO
    d4x_app_start
    return $?
}

# Stop Docker (generic, independent of backend or full app)
d4x_app_stop() {
    # XXX TODO
    local pid
    pid="$(pgrep Docker)" || true
    if [ -n "$pid" ]; then
        kill -2 "$pid"
        kill -9 "$pid"
    fi
    return 0
}


##
## OS X specific functions
##

# Wait for process to appear n seconds - return 0 on success 1 on failure
d4x_wait_for_process() {
    local N=${1}
    local PROC=${2}
    local i

    echo "waiting for process ${PROC}..."
    for i in $(seq 1 "${N}"); do
        pgrep -q "${PROC}" && break
        if [ "$i" -eq "$N" ]; then
            echo "timeout while waiting for ${PROC}"
            return 1
        fi
        echo "waiting for process ${PROC}, sleeping 1s"
        sleep 1
    done
    return 0
}

# Wait for file to appear in n seconds - return 0 on success 1 on failure
d4x_wait_for_file() {
    local N=${1}
    local F=${2}
    local i

    echo "waiting for file ${F}..."
    for i in $(seq 1 "${N}"); do
        [ -e "${F}" ] && break
        if [ "$i" -eq "$N" ]; then
            echo "timeout while waiting for file ${F}"
            return 1
        fi
        echo "waiting for file ${F}, sleeping 1s"
        sleep 1
    done
    return 0
}

# Run a docker command - returns 0 on success. Parameters: [timeout] [docker binary] {commands...}
_d4x_docker() {
    local N=${1}
    local CMD=${2}
    local PID
    local f

    if [ ! -x "${CMD}" ]; then
        echo "file ${CMD} does not exist or is not executable"
        return 0
    fi

    # skip first params
    shift 2

    # execute command in background
    echo "Running ${CMD} $*"
    ${CMD} "$@" &
    PID=$!
    for f in $(seq 1 "$N"); do
        ps -p $PID > /dev/null && break
        if [ "$f" == "$N" ];  then
            kill -9 $PID
            return 1
        fi
        sleep 1
    done
    wait $PID
    return $?
}

# Wait n seconds for docker to start - return 0 on success
d4x_wait_for_docker() {
    local h
    local f
    if [ "$HOME" != "" ]; then
        h=${HOME}
    else
        h="/Users/root"
    fi
    local base_dir="${h}/Library/Containers/com.docker.docker/Data"
    local driver_dir="${base_dir}/com.docker.driver.amd64-linux"

    d4x_wait_for_process 30 Docker || return 1
    d4x_wait_for_file 10 "/var/run/docker.sock" || return 1
    # former /var/tmp/com.docker.port.socket, now s51 in container directory
    d4x_wait_for_file 30 "${base_dir}/s51" || return 1
    d4x_wait_for_file 30 "${driver_dir}/Docker.qcow2" || return 1
    d4x_wait_for_file 30 "${driver_dir}/console-ring" || return 1

    for f in $(seq 1 10); do
        _d4x_docker 15 "${OSX_APP_DIR}/Contents/Resources/bin/docker" ps && break
        if [ "$f" == "10" ]; then
            echo "Unable to run docker ps successfuly"
            return 1
        fi
        sleep 1
    done

    return 0
}

d4m_extract_dmg() {
    local installer="${D4X_TMPDIR}/${D4X_INSTALLER_NAME}"
    local mountpoint="${D4X_TMPDIR}/dmg_mount"

    echo "Verifying dmg..."
    hdiutil verify "${installer}"

    mkdir "${mountpoint}"
    echo "Mounting in ${mountpoint}..."
    hdiutil attach "${installer}" -mountpoint "${mountpoint}" -nobrowse

    echo "Extracting to ${OSX_APP_DIR}"
    mkdir -p "${OSX_APP_DIR}"
    local appname
    appname=$(basename "${OSX_APP_DIR}")

    if [ ! -d "${mountpoint}/${appname}" ]; then
        echo "App dir not found in ${mountpoint}/${appname}..."
        return 1
    fi

    cp -av "${mountpoint}/${appname}/" "${OSX_APP_DIR}"

    hdiutil detach "${mountpoint}" || true
    rmdir "${mountpoint}" || true
}
