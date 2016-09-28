#!/bin/bash

script_name=`basename $0`
DRYRUN=
AUTO_YES=
: ${ENVIRONMENT?"Must be set"}
: ${SECONDARY_ROLE:=store}

projects_with_cron_container=(billing accounts repos)
projects_with_static_container=(billing accounts repos)

source $(dirname $0)/funcs.sh

usage() {
    echo "Deploy projects on:"
    echo "  ENVIRONMENT:    $ENVIRONMENT"
    echo "  SECONDARY_ROLE: $SECONDARY_ROLE"
    echo "Usage: $script_name [options]"
    echo "  Options:
        -p project          : specify project (e.g. 'accounts')
        -from version       : specify current deployment version (e.g. '1.2.3')
        -to version         : specify new deployment version (e.g. '1.2.4')
        -pgbouncer version  : specify pgbouncer version (e.g. '0.8.0')
        -webs n             : specify num web containers
        -workers n          : specify num worker containers
        -yes                : skip all confirmation, and exit on error
        -n                  : dry run
        -h                  : print this help message
    "
}

preflight_addr () {
    if [[ "$ENVIRONMENT" == "aws_stage" ]] ; then
        echo "http://infra.stage-us-east-1.aws.dckr.io:8500"
    elif [[ "$ENVIRONMENT" == "aws_prod" ]] ; then
        echo "http://infra.us-east-1.aws.dckr.io:8500"
    fi
}

has_cron_container() {
    for p in ${projects_with_cron_container[@]} ; do
        if [ "$p" == "$PROJECT" ] ; then
            return 0
        fi
    done
    return 1
}

has_static_container() {
    for p in ${projects_with_static_container[@]} ; do
        if [ "$p" == "$PROJECT" ] ; then
            return 0
        fi
    done
    return 1
}

#
# Main
#
while test $# -gt 0; do
    case "$1" in
        -from)  shift ; OLD_VERSION=$1 ; shift ;;
        -to)    shift ; NEW_VERSION=$1 ; shift ;;
        -pgbouncer) shift ; PGB_VERSION=$1 ; shift ;;
        -webs)  shift ; NUM_WEBS=$1 ; shift ;;
        -workers) shift ; NUM_WORKERS=$1 ; shift ;;
        -p)     shift ; PROJECT=$1 ; shift ;;
        -yes)   AUTO_YES=yes ; set -e ; shift ;;
        -n)     DRYRUN=yes ; shift ;;
        -h)     usage ; exit 0 ;;
        -*)     usage ; exit 1 ;;
        --)     break ;;
    esac
done

# Test input
if [ -z "$PROJECT" ]; then
    usage
    msg -f "Project not set. Please specify '-p' (e.g. repos)."
fi
if [ -z "$OLD_VERSION" ]; then
    usage
    msg -f "From version not set. Please specify '-from' (e.g. 1.2.3)."
fi
if [ -z "$NEW_VERSION" ]; then
    usage
    msg -f "To version not set. Please specify '-to' (e.g. 1.2.4)."
fi
if [ -z "$NUM_WEBS" ]; then
    usage
    msg -f "Num web containers not set. Please specify '-webs' (e.g. 3). Pass in '0' to disable."
fi
if [ -z "$NUM_WORKERS" ]; then
    usage
    msg -f "Num worker containers not set. Please specify '-workers' (e.g. 8). Pass in '0' to disable."
fi

ENV_AND_ROLE=$ENVIRONMENT.$SECONDARY_ROLE

#
# Find the matching pgbouncer
# - PGB_PROJECT signifies which pgbouncer project to use. It should be set even if
#   we're not upgrading pgbouncer.
# - PGB_NEW_IMAGE, if present, will trigger an upgrade of pgbouncer.
#
PGB_PROJECT=
PGB_NEW_IMAGE=
case "$PROJECT" in
    accounts)           PGB_PROJECT=accounts-pgbouncer ;;
    repos)       PGB_PROJECT=repos-pgbouncer ;;
esac

if [ -n "$PGB_VERSION" ]; then
    if [ -n "$PGB_PROJECT" ]; then
        PGB_NEW_IMAGE=docker/${PGB_PROJECT}:${PGB_VERSION}
    else
        usage
        msg -f "I don't know which pgbouncer repo to use for $PROJECT."
    fi
fi

export REPO=docker/${PROJECT}:${NEW_VERSION}

confirm -w "
This script guides you through the installation of
${COL_RED}${REPO}${COL_RESET} on $ENV_AND_ROLE,
upgrading from $OLD_VERSION, running on each node:
- $NUM_WEBS web containers
- $NUM_WORKERS worker containers
If any step fails, pause, fix the issue, and continue.
"

if [ -n "$PGB_NEW_IMAGE" ]; then
    confirm -w "In addition, pgbouncer will be upgraded to ${COL_RED}${PGB_NEW_IMAGE}${COL_RESET}.
"
fi

preflight -h $( preflight_addr )

notify_slack "${SLACK_PREFIX}Deploying version $NEW_VERSION (upgrade from $OLD_VERSION) of $PROJECT to $ENV_AND_ROLE.\n"

# Pull repo
unset DOCKER_CERT_PATH
run_step -d "Pull $REPO on $ENV_AND_ROLE" \
    boss pull $ENV_AND_ROLE $REPO

if [ -n "$PGB_NEW_IMAGE" ]; then
    run_step -d "Pull $PGB_NEW_IMAGE on $ENV_AND_ROLE" \
        boss pull $ENV_AND_ROLE $PGB_NEW_IMAGE
fi

# Start the deployment
run_step -d "List old containers" boss ps $ENV_AND_ROLE --grep $PROJECT
old_container_pattern="$ENVIRONMENT-$PROJECT-$OLD_VERSION"
old_container_cids=`boss ps $ENV_AND_ROLE -a --grep $old_container_pattern | select_cid`
run_step -d "Test that the 'from' version is correct" test "x${old_container_cids}" != 'x'


HOSTS=( $( boss list_hosts -g $ENV_AND_ROLE | tail -n +2 | awk '{print $1}' ) )
FIRST_HOST="${HOSTS[0]}"

if has_static_container ; then
    run_step -d "Run static container" boss run $FIRST_HOST $REPO --type static --wait --skip_pass
fi

if has_cron_container ; then
    stop_containers -d "Stop cron container" -h $FIRST_HOST -n 1 ${old_container_pattern}_cron
    run_step -d "Run cron container" boss run $FIRST_HOST $REPO --type cron --skip_pass
fi

for host in "${HOSTS[@]}"; do
    # Stop all web & worker containers
    stop_containers -d "Stop old app containers on $host" -h $host "${old_container_pattern}_(web|worker)"

    # Upgrade pgbouncer if necessary
    if [ -n "$PGB_NEW_IMAGE" ]; then
        stop_containers -d "Stop pgbouncer container on $host" -h $host $ENVIRONMENT-$PGB_PROJECT
        run_step -d "Run pgbouncer container on $host" \
            boss run $host $PGB_NEW_IMAGE --type worker -n 1 --skip_pass
    fi

    # Check that pgbouncer is running
    if [ -n "$PGB_PROJECT" ]; then
        # It's important to make sure that pgbouncer is started and listening on the right port.
        # This is a hackish way to check it.
        #
        # NOTE: We assume pgbouncer to be listening on port 60?8.
        run_step -d "Check that pgbouncer started correctly on $host" \
            "boss ps $host --grep $PGB_PROJECT | grep -q 60.8:6543"
    fi

    # Start web & worker containers
    if [ $NUM_WEBS -gt 0 ]; then
        run_step -d "Run web containers on $host" \
            boss run $host $REPO --type web -n $NUM_WEBS --skip_pass
    fi

    if [ $NUM_WORKERS -gt 0 ]; then
        run_step -d "Run worker containers on $host" \
            boss run $host $REPO --type worker -n $NUM_WORKERS --skip_pass
    fi

    # List containers on the host
    if [ -z "$AUTO_YES" ]; then
        run_step -d "List app containers on $host" \
            boss ps $host --grep ${PROJECT}
    fi

    # Remove old containers
    run_step -d "Remove old containers on $host" "boss ps $host -a --grep $old_container_pattern | select_cid | xargs boss rm $host"
done

run_step -d "Check final set of containers" boss ps $ENV_AND_ROLE --grep ${PROJECT}

notify_slack "${SLACK_PREFIX}Finished deployment of $PROJECT to $ENV_AND_ROLE"

msg -w "
Congrats. You've finished. Now you should:
. Check the site is up
. Do a docker push/pull
. Check Bugsnag errors
. Check Kibana logs
. Check New Relic
Kill the old containers after you confirmed that things are working.
"
