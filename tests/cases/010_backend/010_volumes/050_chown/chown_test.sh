#!/bin/sh -e

# root is the actor
# nobody is the chowned user/group
# postgres is the observer

NOBODY=65534
REAL_FILE=/tmp/chown_test
LINK_FILE=/tmp/chown_link_test

test_chown()
{
    FILE=$1
    CREATE=$2
    POST_CHECK=$3

    TEST="$FILE chown only uid"
    UID_EXPECTED=nobody:postgres
    rm -f "$FILE"
    $CREATE "$FILE"
    /chown_test uid "$NOBODY" "$FILE"
    UID_RESULT=$(su -c "stat -c %U:%G $FILE" postgres)
    if [ "$UID_EXPECTED" != "$UID_RESULT" ]; then
        echo "$TEST expected $UID_EXPECTED got $UID_RESULT"
        exit 1
    fi
    $POST_CHECK "$TEST"

    TEST="$FILE chown only uid root"
    ROOT_UID_EXPECTED=root:postgres
    /chown_test uid 0 "$FILE"
    ROOT_UID_RESULT=$(su -c "stat -c %U:%G $FILE" postgres)
    if [ "$ROOT_UID_EXPECTED" != "$ROOT_UID_RESULT" ]; then
        echo "$TEST expected $ROOT_UID_EXPECTED got $ROOT_UID_RESULT"
        exit 1
    fi
    $POST_CHECK "$TEST"

    TEST="$FILE chown only gid"
    GID_EXPECTED=postgres:nobody
    rm -f "$FILE"
    $CREATE "$FILE"
    /chown_test gid "$NOBODY" "$FILE"
    GID_RESULT="$(su -c "stat -c %U:%G $FILE" postgres)"
    if [ "$GID_EXPECTED" != "$GID_RESULT" ]; then
        echo "$TEST expected $GID_EXPECTED got $GID_RESULT"
        exit 1
    fi
    $POST_CHECK "$TEST"

    TEST="$FILE chown only gid root"
    ROOT_GID_EXPECTED=postgres:root
    /chown_test gid 0 "$FILE"
    ROOT_GID_RESULT=$(su -c "stat -c %U:%G $FILE" postgres)
    if [ "$ROOT_GID_EXPECTED" != "$ROOT_GID_RESULT" ]; then
        echo "$TEST expected $ROOT_GID_EXPECTED got $ROOT_GID_RESULT"
        exit 1
    fi
    $POST_CHECK "$TEST"

    TEST="$FILE chown uid and gid"
    UIDGID_EXPECTED=nobody:nobody
    chown -h nobody:nobody "$FILE"
    UIDGID_RESULT=$(su -c "stat -c %U:%G $FILE" postgres)
    if [ "$UIDGID_EXPECTED" != "$UIDGID_RESULT" ]; then
        echo "$TEST expected $UIDGID_EXPECTED got $UIDGID_RESULT"
        exit 1
    fi
    $POST_CHECK "$TEST"
}

noop()
{
    TEST=$1
    echo -e "$TEST\tPASS"
}

check_real_file()
{
    TEST=$1

    EXPECTED=postgres:postgres
    RESULT="$(su -c "stat -c %U:%G $REAL_FILE" postgres)"
    if [ "$EXPECTED" != "$RESULT" ]; then
        echo "after [$TEST] test expected $EXPECTED got $RESULT"
        exit 1
    fi
    echo -e "$TEST\tPASS"
}

test_chown "$REAL_FILE" "touch" noop
rm -f "$REAL_FILE"

touch "$REAL_FILE"
test_chown "$LINK_FILE" "ln -s $REAL_FILE" check_real_file
