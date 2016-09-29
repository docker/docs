
# Interesting strings to find in the logs

Here are things we could find in the logs and highlight automatically,
with suggested tickets we could link them to. It's not guaranteed the
linked issue is a duplicate but it might still be useful as a reference.

## transfused crash

In docker-system.log:

Transport endpoint is not connected

see #5

## qcow2 corruption

In docker-system.log:

Invalid_argument("Cstruct.sub

see #11

## /var/lib/docker corruption

In /moby/var/log/docker.log:

invalid checksum digest format
invalid mount id value
Error starting daemon:

see #20

## failure to talk to vmnetd

In docker-system.log:

Communication with networking components faile

see #61

## too many hyperkits

In ps-ax, more than one "... hyperkit -A -m ..."

see #71

## still has virtualbox

In docker-system.log:

Docker does not rely on Virtualbox but may not work properly on systems with VirtualBox versions prior to v4.3.30
        VirtualBox v4.3.28 is currently installed.
        Please upgrade or uninstall Virtualbox.

see #78

## moby failed to boot

Check the console, look for the login prompt.

see #84

## Moby kernel problems

In docker-system.log: "BUG" or "rcu_sched self-detected stall"

see #87

## Invariants

In docker-system.log: "INVARIANT VIOLATED"

see #89

