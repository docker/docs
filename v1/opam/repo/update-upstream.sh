#!/usr/bin/env sh

set -ex

OPAM_SWITCH=4.02.3
OPAM_REPO=https://github.com/ocaml/opam-repository.git
REPO_DIR_NAME=opam-mini-repo

TARGET_DIR=$(pwd)
WORK_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'opam-mini-repo')

cleanup() {
    rm -r "${WORK_DIR:?}"
}

#trap cleanup EXIT

## Fetch the opam-repository

cd "${WORK_DIR}"
git clone --depth=1 ${OPAM_REPO} ${REPO_DIR_NAME}
cd ${REPO_DIR_NAME}

## copy the dev/ and local/ packages in the repo
cp -LR "${TARGET_DIR}/packages/dev" packages/dev
cp -LR "${TARGET_DIR}/packages/local" packages/local
git add packages/dev packages/local
git commit -a -m "Adding local and dev packages"

# Remove the upstream packages that are copied in /dev
for pkg in $(ls packages/dev); do
  upstream="packages/${pkg%%.*}/${pkg}"
  if [ -d "${upstream}" ]; then
    rm -rf "${upstream}"
  fi
done
git commit -a -m "Remove upstream source of dev packages" || echo "ok"

## Compute the list of packages needed

PACKAGES="$* $(ls packages/local | xargs) $(ls packages/dev | xargs)"
echo "PACKAGES=${PACKAGES}"

export OPAMROOT=${WORK_DIR}/opam

opam init --root=${OPAMROOT} -n .

export OPAMSWITCH=${OPAM_SWITCH}
export OPAMNO=1

# ugly hack to make opam think that the switch is already installed
# and to overwrite opam internal variables
echo "${OPAM_SWITCH} ${OPAM_SWITCH}" > ${OPAMROOT}/aliases
mkdir -p "${OPAMROOT}/${OPAM_SWITCH}/config"

config="${OPAMROOT}/${OPAM_SWITCH}/config/global-config.config"

function add {
    key=$1
    value=$2
    echo "${key}: \"${value}\"" > ${config}
}

add ocaml_version ${OPAM_SWITCH%%+*}
add compiler ${OPAM_SWITCH}
add preinstalled false
add os darwin

echo "ocaml-version=$(opam config var ocaml-version)"
echo "compiler=$(opam config var compiler)"
echo "preinstalled=$(opam config var preinstalled)"
echo "os=$(opam config var os)"

OUTPUT=${WORK_DIR}/pkgs.json
opam install --root=${OPAMROOT} ${PACKAGES} --dry-run --json=${OUTPUT}
ALL_PACKAGES=$(jq '.[] | map(select(.install)) | map( [.install.name, .install.version] | join(".")) | join(" ")' ${OUTPUT})

## Copy the package metadata that are needed in packages/upstream

rm -rf "${TARGET_DIR}/packages/upstream/"
mkdir -p "${TARGET_DIR}/packages/upstream/"

# Note: not sure why this is needed, but it is
BASE_PKGS="base-unix.base base-threads.base base-bigarray.base"

for pkg in ${BASE_PKGS} ${ALL_PACKAGES//\"}; do
    echo Adding ${pkg}
    if [ -d ${TARGET_DIR}/packages/dev/${pkg} ]; then
        echo "${pkg} is a dev package, skipping."
    elif [ -d ${TARGET_DIR}/packages/local/${pkg} ]; then
        echo "${pkg} is a local package, skipping."
    else
        cp -R packages/${pkg%%.*}/${pkg} "${TARGET_DIR}/packages/upstream/"
    fi
done

# Install the compiler
# rm -rf ${TARGET_DIR}/compilers
# mkdir ${TARGET_DIR}/compilers
# find compilers -name "${OPAM_SWITCH}.comp" -exec cp {} ${TARGET_DIR}/compilers/ \;
# find compilers -name "${OPAM_SWITCH}.descr" -exec cp {} ${TARGET_DIR}/compilers/ \;
