#!/bin/bash

# Colour resources
ESC_SEQ="\x1b["
COL_RESET=$ESC_SEQ"39;49;00m"
COL_RED=$ESC_SEQ"31;01m"
COL_GREEN=$ESC_SEQ"32;01m"
COL_YELLOW=$ESC_SEQ"33;01m"
COL_BLUE=$ESC_SEQ"34;01m"
COL_MAGENTA=$ESC_SEQ"35;01m"
COL_CYAN=$ESC_SEQ"36;01m"


#
# Usage:
#       msg [options] <message_to_print>
#
#       Options:
#               -f      Fatal error; exit after print
#               -e      Red error
#               -w      Yellow alert
#               -nc     No color
#               -n      No newline after message
#
msg() {
    local col=$COL_GREEN
    local flags=""
    while test $# -gt 0; do
        case "$1" in
            -w) col=$COL_YELLOW ; shift ;;
            -e) col=$COL_RED ; shift ;;
            -f) col=$COL_RED ; fatal=yes ; shift ;;
            -nc) col=$COL_RESET ; shift ;;
            -n) flags="-n" ; shift ;;
            --) shift ; break ;;
            -*) break ;;
            *)  break ;;
        esac
    done
    echo -e $flags "${col}$*${COL_RESET}"
    test "$fatal" && exit 1 || :
}

#
# Print a message and ask for confirmation. Takes the same args as msg().
# Honours $AUTO_YES if set, which skips the user input step.
#
confirm() {
    msg "$@"
    msg -n '(hit enter to continue) '
    if [ ! "$AUTO_YES" ]; then
        read input
    fi
}

#
# Same as confirm(), but allow skip.
# Honours $AUTO_YES if set, which skips the user input step.
# Return codes:
#       0       Continue
#       1       Skip
#
skippable_confirm() {
    msg "$@"
    msg -n '("skip" to skip, or hit enter to continue) '
    if [ ! "$AUTO_YES" ]; then
        read input
        if [ "$input" == "skip" ]; then
            return 1
        fi
    fi
    return 0
}

#
# Run a step with confirmation. Honours $DRYRUN if set.
#
run_step() {
    local doc=
    while test $# -gt 0; do
        case "$1" in
            -d) shift; doc="$1"; shift ;;
            *)  break ;;
        esac
    done

    local cmd="$*"

    if [ "$DRYRUN" ]; then
        confirm -w "\n[DRY RUN] About to: $doc -- \`$cmd\`"
        echo [DRY RUN] "$cmd"
    else
        skippable_confirm -w "\nCONFIRM: $doc -- \`$cmd\`"
        if [ $? -eq 1 ]; then
            msg -w "Skipping step"
            return 0
        else
            msg -nc "running ..."
            eval "$cmd"
        fi
    fi

    if [ $? -ne 0 ]; then
        confirm -e "\`$cmd\` FAILED! Fix it please."
        return 1
    else
        msg "DONE with: $doc (\`$cmd\`)"
        return 0
    fi
}

#
# Preflight checks. Usage:
#
#       preflight [options]
#
#       Options:
#               -h url_to_curl
#
preflight() {
    local host=
    while test $# -gt 0; do
        case "$1" in
            -h) shift ; host="$1" ; shift ;;
            *)  break ;;
        esac
    done

    if [ -n "$host" ]; then
        run_step -d "Check VPN" curl -s -o /dev/null $host
    fi
    run_step -d "Check hub-boss" boss version
    run_step -d "Update pass" pass git pull
}

#
# Usage:
#       stop_containers [options] <regex>
#
#       Options:
#               -d description_message
#               -n number_of_expected_containers_to_stop
#               -h host_expression
#
# E.g. stop_containers -n 4 docker-index-1.2.3
#
stop_containers() {
    local n_expected=
    local doc=
    local hosts=
    while test $# -gt 0; do
        case "$1" in
            -n) shift ; n_expected="$1" ; shift ;;
            -d) shift ; doc="$1" ; shift ;;
            -h) shift ; hosts="$1" ; shift ;;
            *)  break ;;
        esac
    done

    if [ -z "$hosts" ]; then
        msg -f "INTERNAL ERROR: Please specify host expression (-h)"
        exit 1
    fi

    local regex="$1"
    if [ -z "$regex" ]; then
        msg -f "INTERNAL ERROR: Please specify regex for containers to stop"
        exit 1
    fi

    msg "\n$doc -- about to stop the following $n_expected container(s):"
    msg "boss ps $hosts --grep '$regex'"
    boss ps $hosts --grep "$regex"
    echo

    local container_ids=`boss ps $hosts --grep "$regex" | select_cid`
    IFS=',' read -a cid_array <<< $container_ids
    if [ -n "$n_expected" ]; then
        # Be lenient and allow the actual number being smaller than expected
        if [ "${#cid_array[@]}" -gt "$n_expected" ]; then
            msg -e "Expecting to stop $n_expected containers. Found ${#cid_array[@]}. Aborting step."
            msg -e "$matching"
            return 1
        fi
    fi

    run_step -d "$doc -- stop container IDs: $container_ids" boss stop $hosts $container_ids
}

# Notify the #store channel.
#
# Usage:
#       notify_slack <message>
#
#       Options:
#               -c channel (with #)
#
# E.g. notify_slack "deployment starting now"
# E.g  notify_slack -c #marcus-test "deployment starting 3, 2, 1..."
#
notify_slack() {
    local channel="#store-team"
    while test $# -gt 0; do
        case "$1" in
            -c) shift ; channel="$1" ; shift ;;
            *)  break ;;
        esac
    done

    local msg="$1"
    if [ -z "$msg" ]; then
        echo "You need to pass a parameter like notify_slack 'my message'"
        return 1
    fi

    local username="Store Bot"
    local avatar=":department_store:"
    local hook_url="https://hooks.slack.com/services/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"

    local cmd='curl -X POST --data-urlencode "payload={\"channel\": \"$channel\", \"username\": \"$username\", \"text\": \"$msg\", \"icon_emoji\": \"$avatar\"}" $hook_url'

    if [ "$DRYRUN" ]; then
        echo [DRY RUN] "$cmd"
    else
        eval "$cmd"
    fi
}

