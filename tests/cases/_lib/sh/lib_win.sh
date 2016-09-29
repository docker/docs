#!/bin/sh
# Copyright (C) 2016 Rolf Neugebauer <rolf.neugebauer@docker.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#

# utility function
_win_path() {
    PATH=$1
    PATH=${PATH/\/c\//C:\\} # replace /c/ with C:\
    PATH=${PATH//\//\\}     # replace / with \
    echo "${PATH}"
}

##
## Common functions which must be implemented on OSX and Windows
##

# Return true is Docker for Windows is installed
d4x_app_installed() {
    # shellcheck disable=SC2016
    installed=$(rt_ps_cmd 'Get-ItemProperty HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | ? {$_.DisplayName -eq "Docker"} | % {Write-Host $_.PSChildName}')
    if [ -z "${installed}" ]; then
        return 1
    else
        return 0
    fi
}

# Return true if the Docker for Windows is running
d4x_app_running() {
    out="$(tasklist.exe | grep Docker)"
    if [ -z "$out" ]; then
        return 1
    else
        return 0
    fi
}


# Install Docker for Windows
d4x_app_install() {
    _installer="${D4X_TMPDIR}/${D4X_INSTALLER_NAME}"
    echo "Install Docker from ${_installer}..."
    installer=$(_win_path "${_installer}")
    "${RT_UTILS}/rt-elevate.exe" -wait msiexec /i "${installer}" /quiet
}

# Uninstall Docker for Windows
d4x_app_uninstall() {
    # shellcheck disable=SC2016
    installed=$(rt_ps_cmd 'Get-ItemProperty HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | ? {$_.DisplayName -eq "Docker"} | % {Write-Host $_.PSChildName}')
    for product_code in ${installed}; do
        echo "Uninstall Docker for Windows ${product_code}..."
        "${RT_UTILS}/rt-elevate.exe" -wait msiexec /x "${product_code}" /quiet
    done
}


# Start the Docker for Windows application. Return true on success
d4x_app_start() {
    # Start the application
    # redirect stdout/stderr to properly start in background
    # Make sure it starts somewhere known, not just the current directory)
    pushd /
    /c/Program\ Files/Docker/Docker/Docker\ for\ Windows.exe \
        -DisableCheckForUpdates -DisableWelcomeWhale -DisableToolboxMigration \
        -Username="${D4X_USERNAME}" -Password="${D4X_PASSWORD}" \
        >/dev/null 2>/dev/null &
    popd
    # wait till it's up
    echo "Wait for docker daemon"
    PATH=$PATH:"/c/Program Files/Docker/Docker/Resources/bin"
    RUNNING=1
    for i in $(seq 1 60); do
        set +e
        docker version
        EXIT_CODE=$?
        set -e

        if [ "$EXIT_CODE" -ne 0 ]; then
            echo "$i: Docker is not yet available. Waiting"
            sleep 1
        else
            echo "Docker is up and running"
            RUNNING=0
            break
        fi
    done
    
    d4w_backend_cli -Mount=C

    return ${RUNNING}
}

# Start the backend only (return true on success)
d4x_backend_start() {
    # On windows we start the full application
    d4x_app_start
}

# Stop Docker (kills the frontend process)
d4x_app_stop() {
    echo "Stop Docker for Windows..."
    set +e
    rt_ps_cmd "Stop-Process -Name \"Docker for Windows\" -ea SilentlyContinue"
    set -e
}


##
## Windows only functions
##

# Connect to Docker for Windows Backend to trigger features that are
# otherwise only available through the GUI.
d4w_backend_cli() {
    SECRET_KEY="--testftw!928374kasljf039"
    DOCKER_FOR_WINDOWS_PATH="/c/Program Files/Docker/Docker/DockerCli.exe"
    DOCKER_FOR_WINDOWS_BUILD_PATH="$RT_PROJECT_ROOT/../../win/build/win/DockerCli.exe"

    if [ -e "$DOCKER_FOR_WINDOWS_PATH" ]; then
        echo "Using $DOCKER_FOR_WINDOWS_PATH as backend cli"

        "$DOCKER_FOR_WINDOWS_PATH" "$SECRET_KEY" -Username="${D4X_USERNAME}" -Password="${D4X_PASSWORD}" "$@"
        return $?
    elif [ -e "$DOCKER_FOR_WINDOWS_BUILD_PATH" ]; then
        echo "Using $DOCKER_FOR_WINDOWS_BUILD_PATH as backend cli"

        "$DOCKER_FOR_WINDOWS_BUILD_PATH" "$SECRET_KEY" -Username="${D4X_USERNAME}" -Password="${D4X_PASSWORD}" "$@"
        return $?
    fi

    echo "Unable to find Backend CLI"
    return "${RT_TEST_CANCEL}"
}

# Get the VM ID. This is a privileged operation and we can't get the
# output from the invocation. Redirect to a file and then read it.
d4w_vmid() {
    "${RT_UTILS}/rt-elevate.exe" -wait powershell.exe -NoLogo -WindowStyle Hidden -NoProfile -NonInteractive -Command "(get-vm mobylinuxvm).Id.Guid > ${D4X_TMPDIR}/vmid"
    dos2unix.exe -q "${D4X_TMPDIR}"/vmid
    cat "${D4X_TMPDIR}"/vmid
}

# Docker load the d4w/nsenter image
d4w_load_nsenter() {
    DOCKER_FOR_WINDOWS_PATH="/c/Program Files/Docker/Docker/src/resources/nsenter.tar"
    DOCKER_FOR_WINDOWS_BUILD_PATH="${RT_PROJECT_ROOT}/../../win/src/resources/nsenter.tar"

    if [ -e "$DOCKER_FOR_WINDOWS_PATH" ]; then
        docker load -i "${DOCKER_FOR_WINDOWS_PATH}"
    elif [ -e "$DOCKER_FOR_WINDOWS_BUILD_PATH" ]; then
        docker load -i "${DOCKER_FOR_WINDOWS_BUILD_PATH}"
    fi
}

# Run the d4w/nsenter image
d4w_nsenter() {
    d4w_load_nsenter
    docker run --rm --privileged --pid=host d4w/nsenter /bin/sh -c "$@"
}