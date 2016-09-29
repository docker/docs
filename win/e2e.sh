#!/bin/bash

# Check it's running in elevated mode
net session >/dev/null 2>&1
if [[ $? -ne 0 ]]; then
    echo "Must run this script with elevated rights"
    exit 1
fi

# Check python is installed
command -v python >/dev/null 2>&1
if [[ $? -ne 0 ]]; then
    echo "Python has to be on the path"
    exit 1
fi

function psCommand {
    powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Unrestricted -Command "$1"
}

function psScript {
    powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Unrestricted -File "$1" $2
}

function winPath {
    PATH=$1
    PATH=${PATH/\/c\//C:\\} # replace /c/ with C:\
    PATH=${PATH//\//\\}     # replace / with \
    echo $PATH
}

function killDocker {
    echo "Stop Docker for Windows..."
    set +e
    psCommand "Stop-Process -ProcessName \"Docker for Windows\" -ea SilentlyContinue" >/dev/null 2>/dev/null
    set -e
}

function uninstallDocker {
    for installed in $(psCommand 'Get-ItemProperty HKLM:\\Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | ? {$_.DisplayName -eq "Docker"} | % {Write-Host $_.PSChildName}'); do
        echo "Uninstall Docker for Windows $installed..."
        psCommand "Start-Process msiexec -Wait -arg \"/x$installed /quiet\""
    done
}

set -e

export MSYS_NO_PATHCONV=1
killDocker
uninstallDocker

if [[ -z "$SKIPBUILD" ]]; then
  echo "Build the msi..."
  export SKIP_TESTS="1"
  BUILD_SCRIPT=$(winPath "$(pwd)/please.ps1")
  psScript "$BUILD_SCRIPT" package
fi

echo "Install Docker..."
INSTALLER=$(winPath "$(pwd)/build/InstallDocker.msi")
psCommand "Start-Process msiexec -Wait -arg \"/i ${INSTALLER} /quiet\""

D4X_USERNAME="$(whoami)"
D4X_PASSWORD="${USERPASSWORD-Docker4theTeam!}"

echo "Run Docker for Windows"
"/c/Program Files/Docker/Docker/Docker For Windows.exe" \
  -DisableCheckForUpdates \
  -DisableWelcomeWhale \
  -DisableToolboxMigration \
  -Username=${D4X_USERNAME} \
  -Password=${D4X_PASSWORD} \
  >/dev/null 2>/dev/null &

echo "Wait for docker daemon"
export PATH=$PATH:"/c/Program Files/Docker/Docker/Resources/bin"
RUNNING=0
for i in `seq 1 60`; do
    set +e
    docker run --rm busybox echo hello >/dev/null 2>/dev/null
    EXIT_CODE=$?
    set -e

    if [[ $EXIT_CODE -ne 0 ]]; then
        echo "Docker is not yet available. Waiting"
        sleep 2
    else
        echo "Docker is up and running"
        RUNNING=1
        break
    fi
done

if [[ $RUNNING -eq 1 ]]; then
    echo "Share C drive..."
    "/c/Program Files/Docker/Docker/DockerCli.exe" \
      "--testftw!928374kasljf039" \
      -Share=C -Username="${D4X_USERNAME}" -Password="${D4X_PASSWORD}"

    echo "Run tests..."
    pushd ../tests
    ./rt-local -vvx -l ${TARGET-nostart} run
    popd
fi

killDocker
uninstallDocker

if [[ $RUNNING -ne 1 ]]; then
    exit 1
fi
