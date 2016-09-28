# Bootstrapper Spec

Excluding the initial `install` command, the hope is that customers will
rarely interact with the bootstrapper outside of troubleshooting/break-fix
scenarios, although it may take more than one release to fully achieve
that goal. (Even `install` might be something we can address in a future
release once our "proxy beachhead" matures.)

Installation Straw-man:

1. Generate Root CA materials
2. Generate local certs for the UCP services
    * Including server certs for the new swarm v2 manager, placed on the host filesystem in a well known location
3. Stand up Root CA containers
4. Stand up remaining containers
5. Poke the local daemon to start the swarm manager
    * Point to our CA
        * Most likely based on https://github.com/docker/swarm-v2/issues/521

* If the node is already part of the cluster, we reset it
    * In the future we can offer the ability to "take over" the cluster, but that's not critical for this release (mentioning here in case it turns out to be low cost)


## Transition Plan

Basic ideas, things might get shuffled around depending on the shape of the new
code, or if anything is blocking important workflows. Fuzziness increases with
further away milestones.

1. Sprint 1: Target the initial internal release from core group.
    * If core `cluster init` works, build on that to get a single node cluster
      functional.
    * Factor out as much as possible from `join` command for reuse within the
      proxy.
    * Fix any install breakage and otherwise preserve the capability to make a
      cluster.
    * Convert installation to be a service (if that is solidified here)
2. Sprint 2:
    * As soon as swarm v2 joining is possible, either implement the bare
      minimum join code in utility functions, or document the manual workflow.
    * Rev targetted release when core team provides a new milestone, fix any
      breakage to keep basic workflow intact
    * Finish migration of join out of bootstrapper.
3. Post-Dockercon 1:
    * Remove hacks
    * Remove bootstrapper commands that won't be staying around in the release
      or replace them with stubs that inform the user of changes (join,
      engine-discovery)
    * Backup / Restore
    * Other things we've missed
4. Post-Dockercon 2:
    * Upgrade
    * Other missing workflows, testing

## Other notable aspects

Like etcd, if quorum is lost, the system becomes wedged to prevent
split-brain.  Also like etcd, there isn't a super simple recovery model.
You must shut down all managers, then bounce one of them with a `--force`
flag and it will reset the cluster.

## Misc Bootstrapper Commands

* Backup
    * Should include cluster state (`/var/lib/docker/cluster...?`)
* Restore
    * Swarm v2 manager will have some sort of `--force` flag to reset to a cluster of 1 with the current configuration
* Join
    * once we have the above auto join logic wired up, this command should error out with a usage message telling users to run `docker cluster join ...` instead
* Uninstall
    * Keep this around, and make sure it cleanly unwinds a node at the swarm level
* `regen-certs`
    * Some certs go away (superseded by core), but still have this command to regenerate our certs
* `engine-discovery`
    * Removed, should be automatically set up by core
* Upgrade
    * TBD but probably needs some work
* `dump-certs`, `fingerprint`, `id`, `images`
    * No change
* Stop and Restart
    * No change
