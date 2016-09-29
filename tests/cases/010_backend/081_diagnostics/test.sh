#!/bin/sh
# SUMMARY: Verify that Moby diagnostics work
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm -rf foo.tar
}
trap clean_up EXIT

case "${RT_OS}" in
    osx)
        echo "00000003.0000f3a6" | nc -U ~/Library/Containers/com.docker.docker/Data/@connect > foo.tar
        ;;
    win)
        vmid=$(d4w_vmid)
        # XXX This does not work when run against a local build
        "/c/Program Files/Docker/Docker/resources/moby-diag-dl.exe" -o foo.tar -vmid "$vmid"
        ;;
    *)
        exit 1
        ;;
esac

# check the file is there and non-zero
[ -s foo.tar ] || exit 1

out=$(tar tvf foo.tar)

# check that the tarfile list contains these strings
STRINGS="Database date uname ifconfig route brctl docker.log messages daemon.json resolv dig wget vsockd vsudd"

for string in ${STRINGS}; do
    echo "Checking: ${string}"
    echo "${out}" | grep -q "${string}"
done

exit 0
