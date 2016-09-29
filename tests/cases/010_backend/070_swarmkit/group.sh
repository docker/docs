#!/bin/sh
# SUMMARY: SwarmKit regression tests
# LABELS: !benchmarks

# Source libraries. Uncomment if needed/defined
# . ${RT_ROOT}/lib/lib.sh
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

SWARM_DID_INIT_PATH="${D4X_TMPDIR}/010_070_001_swarm_did_init"

group_init()
{
    # Group initialisation code goes here

    # enable swarm mode if it isn't already
    docker swarm inspect
    if [[ $? -ne 0 ]]; then
        touch "${SWARM_DID_INIT_PATH}"
        [ $? -ne 0 ] && return $?
        docker swarm init
        # XXX There may be a race between swarm init and before
        # starting a service. Sleep here. Remove when fixed
        sleep 5
    fi
    return $?
}

group_deinit()
{
    # Group de-initialisation code goes here

    # disable swarm mode if we enabled it earlier in this test
    if [[ -e "${SWARM_DID_INIT_PATH}" ]]; then
        rm "${SWARM_DID_INIT_PATH}"
        [ $? -ne 0 ] && return $?
        docker swarm leave --force
        return $?
    fi
    return 0
}

CMD=$1
case $CMD in
init)
    group_init
    res=$?
    ;;
deinit)
    group_deinit
    res=$?
    ;;
*)
    res=1
    ;;
esac

exit $res
