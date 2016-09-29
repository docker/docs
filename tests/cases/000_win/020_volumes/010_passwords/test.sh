#!/bin/sh
# SUMMARY: Check sharing with different passwords
# LABELS: win
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

ORIG_PASSWORD="${D4X_PASSWORD}"
ORIG_USERNAME="${D4X_USERNAME}"
TESTDIR="/c/Users/Public/password_test"

change_passwd() {
    # For some reason,
    # ([adsi]"WinNT://$(Hostname)/$env:username").ChangePassword("old", "new!")
    # only seems to work on some systems.
    # Also, escaping special characters for multiple shells is tedious
    # if not nay impossible. So, we just create a PS script on the fly
    # and execute it elevated.
    echo "([adsi]\"WinNT://${HOSTNAME}/${D4X_USERNAME}\").SetPassword('$1')" > tmp.ps1
    # convert the script to UTF-16
    iconv.exe -c -f utf-8 -t utf-16  < tmp.ps1 > passwd.ps1
    up=$(pwd)/passwd.ps1
    wp=$(cygpath -w "$up")
    "${RT_UTILS}/rt-elevate.exe" -wait powershell.exe -NoLogo -WindowStyle Hidden -NoProfile -NonInteractive -File "$wp"
    return 0
}

clean_up() {
    D4X_PASSWORD="${ORIG_PASSWORD}"
    D4X_USERNAME="${ORIG_USERNAME}"
    d4w_backend_cli -Mount=C
    rm -rf "${TESTDIR}"
    rm -rf passwd.ps1
}
trap clean_up EXIT

# assumes that the C drive is not shared
test_passwd () {
    _pass="$1"

    set +e
    out=$(docker run --rm -v "${TESTDIR}":/testdir alpine ls /testdir)
    EXIT_CODE=$?
    set -e

    [ "$EXIT_CODE" -eq 0 ] && exit 1
    [ "$out" != "${out/drive is not shared/}" ] && exit 1

    change_passwd "$_pass"
    D4X_PASSWORD="$_pass"
    d4w_backend_cli -Mount=C

    out=$(docker run --rm -v "${TESTDIR}":/testdir alpine ls /testdir)
    [ -z "$out" ] && exit 1

    d4w_backend_cli -Unmount=C

    set +e
    out=$(docker run --rm -v "${TESTDIR}":/testdir alpine ls /testdir)
    EXIT_CODE=$?
    set -e

    [ "$EXIT_CODE" -eq 0 ] && exit 1
    [ "$out" != "${out/drive is not shared/}" ] && exit 1

    return 0
}

for username in 'TestUser' 'Test User'; do
    "${RT_UTILS}/rt-elevate.exe" -wait powershell.exe -NoLogo -WindowStyle Hidden -NoProfile -NonInteractive -File "newuser.ps1" \""${username}"\"

    D4X_USERNAME="${username}"
    rm -rf "${TESTDIR}" || true
    mkdir -p "${TESTDIR}"
    touch "${TESTDIR}"/foo
    testdirpath=$(cygpath -w "${TESTDIR}")
    icacls "${testdirpath}" /grant "${D4X_USERNAME}":F /T

    d4w_backend_cli -Unmount=C

    for p in 'P@s;w," .%^*03rd$`\/r!' 'd$!ø§£µîüéàæøå4' 's=#~2))-+|{}[]<>'; do
        echo "Trying: '$p'"
        test_passwd "$p"
    done
done

exit 0
