#!/bin/sh
##
## A library of shell functions
##

# Special return code to indicate that a test was cancelled
RT_TEST_CANCEL=253

# Echo to stderr
echoerr() {
    echo "$@" 1>&2
}


# Check if a label is set. Returns 0 if not set, 1 otherwise
rt_label_set() {
    label=$1
    res=1

    if [ "x$RT_LABELS" = "x" ]; then
        return $res
    fi

    OLDIFS=$IFS
    IFS=:
    for l in $RT_LABELS; do
        if [ "$l" = "$label" ]; then
            res=0
            break
        fi
    done

    IFS=$OLDIFS
    return $res
}

# assert_equals [message] expected actual
assert_equals() {
    message=''
    if [ $# -eq 3 ]; then
        message=$1
        shift
    fi
    expected=${1:-}
    actual=${2:-}

    echo "Actual $actual - Expected $expected - Msg $message"

    returned=0
    if [ "${expected}" = "${actual}" ]; then
        returned=0
    else
        echo "${message} - expected: <${expected}> but was: <${actual}>"
        returned=1
    fi

    unset message expected actual
    return $returned
}

# assert_true [message] result
assert_true() {
    message=''
    if [ $# -eq 2 ]; then
        message=$1
        shift
    fi
    condition=${1:-}

    returned=0

    # see if condition is an integer, i.e. a return value
    matching=`expr "${condition}" : '\([0-9]*\)'`
    if [ -z "${condition}" ]; then
        # null condition
        returned=1
    elif [ "${condition}" = "${matching}" ]; then
        # possible return value. treating 0 as true, and non-zero as false.
        [ ${condition} -ne 0 ] && returned=1
    else
        # (hopefully) a condition
        ( eval ${condition} ) >/dev/null 2>&1
        [ $? -ne 0 ] && returned=1
    fi

    # record the test
    if [ ${returned} -eq 0 ]; then
        echo "Condition is true"
    else
        echo "Condition is false - ${message}"
    fi

    unset message condition matching
    return ${returned}
}

# assert_not_null [message] result
assert_not_null() {
    if [ $# -eq 2 ]; then
        assert_true "$1" "[ -n '$2' ]"
    else
        assert_true "[ -n '${1:-}' ]"
    fi
}

# assert_null [message] result
assert_null() {
    if [ $# -eq 2 ]; then
        assert_true "$1" "[ -z '$2' ]"
    else
        assert_true "[ -z '${1:-}' ]"
    fi
}

# Some helper functions for windows

# Run a Powershell command
rt_ps_cmd() {
    powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Unrestricted -Command "$@"
}

# Run a Powershell script
rt_ps_script() {
    _script=$1
    shift
    powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Unrestricted -File "$_script" "$@"
}
