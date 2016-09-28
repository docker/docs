# Node Administration Spec

Typical Sequence:
1. `docker cluster join ...` pointed to UCP controller -- swarm manager reaches out to CA
    * This request is async, and the daemon will poll until the task completes or is rejected
2. CA plugin adds node to pending list in UCP
3. Admin accepts join(s) in UCP UI
    * Or they run `docker cluster accept ...` via CLI pointed at UCP, which we hijack for our node management
4. CA generates cert, and updates the task as completed
5. Daemon now gets the cert, and "connects" to the cluster

* Switch all of our node management APIs over to map to the new 1.12 node management APIs
    * We'll interpose/modify/extend them as needed.
    * On node removal, make sure the relevant cert CNs are blacklisted (overlap with CA spec...)
        * For manager remove, do a proper clean-up of our metadata
        * Implement https://github.com/docker/orca/issues/1018
* Switch UI to start using the `/cluster/info` endpoint to get information about the local controller
    * Keep this info in the `/info` endpoint as it is today for legacy CLI/API support
* We'll need API for cluster membership for both v1 and v2 clusters
    * For Seattle, linux nodes should be 1:1 in the two clusters (unless we implement refinmenets below)
    * Support for Windows systems is still TBD
    * Don't need first-class UI for v1 membership, but at the API level make sure we keep the v1/v2 distinction and don't codify bad assumptions that they're identical

Possibilities/refinements: (optional although might be mandatory as we flesh it out)
* Update advertise address
    * Still somewhat murky, but we may want (need?) to offer the ability to change the advertise address of a node after it has joined (we may be able to punt this and require the user to simply re-join with the correct public/accessible address for the node)
* Opt-in backwards compatibility
    * Possibly for the next release after this one we may offer a toggle on a per-node basis to enable/disable the `ucp-swarm-join` container which will control what nodes swarmv1 sees.

