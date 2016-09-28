# Upgrade support in UCP

There are a number of considerations involved during upgrade.
The following sections break down the parts of the system involved
during an upgrade, and pointers to where changes must be made when adding
support for new upgrade paths.

From a users perspective, they use the bootstrapper to upgrade each
individual node in the cluster, with the controller nodes being upgraded
first.  The bootstrapper is involved in preserving the data volumes,
deleting the old containers, and spinning up new continuers with new
images pointed at the old data volumes.


## Data Transformation/Update (aka, schema upgrade)

At present, all "schema" in UCP is stored in the KV store.  As a general
philosophy, we strive to migrate the data forward during controller
startup.  When performing the data upgrade, new fields should be added,
but old fields should be left intact.  The motivation for this is
to allow individual controllers to come up in an HA cluster during an
upgrade without having inconsistencies, provided the user is not editing
the configuration during the upgrade process.  If the schema must be
transformed in a non-backwards compatible way, a new versioned tree in
the KV store should be established.

* On startup, the controller should detect if old configuration is
  present **without** the corresponding new configuration.
* If detected, migrate/transform the old configuration into the new form
  while preserving the old configuration.
    * The first HA node that is ugpraded will perform this transformation
    * Not-yet upgraded nodes will continue to operate with the old data
      until they are upgraded.  Once upgraded, they will detect the data
      has already been transformed and will do nothing.

The implementation for the data transformation belongs in the relevant
configuration subsystems and should only be triggered on controller
startup.  Typically `controller/manager/XXX.go` in a routine called
something like `setupXXX`.

While it is possible to add data transformations within the bootstrapper
the goal is to avoid this unless deemed absolutely necessary.  Examples
where this may be required are volume "renames" or shifting of data
from one volume to another. (Not yet implemented as we haven't hit
that use-case.)  While the bootstrapper does have some knowledge of the
internal data mode of parts of the configuration, to the extent possible
we want to limit that.


### Cleaning up old/stale data

This is not yet implemented, and until we have a large dataset that has
been migrated, we can leave the cruft since users don't interact with the
KV store directly.  At some point we will implement the cleanup, which
must ensure all nodes in the cluster have been upgraded before removing
the old state.  Depending on the scenario, this could be as simple as
a "curl script" we explain to support that they can hep customers run
when/if they have questions about extra data in the KV store.

This may be a case where implementing this within the bootstrapper
makes sense, as it can be done explicitly by the user once all nodes
are upgraded and they're confident they don't need to support rolling
back to the prior version.


## Identifying Compatible Versions

During an upgrade, the bootstrapper checks image metadata to determine
if the users requested upgrade path is supported.  This is accomplished
by inspecting Labels on the image - `com.docker.ucp.version` and
`com.docker.ucp.upgrades_from` for compatibility.  At present, the
`upgrades_from` is a simple list so each upgrade vector must be explicitly
listed. **For each new external release, these fields must be updated.**
If no compatibility information is detected, the bootstrapper makes
the assumption that the container in question can be upgraded from any
prior version.  All containers in the system must be compatible with their
prior versions for the bootstrapper to proceed with the upgrade flow.

The algorithm for processing this is located in the bootstrapper in
`client/images.go` with `CheckUpgradeCompatible`


# TODO

* Add built-in end-to-end upgrade support in the UI where the bootstrapper
  is run by the system
* Add support for wildcard and/or range based upgrade compatibility identification
* Implement cleanup algorithm
