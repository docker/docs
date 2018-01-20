---
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: Docker, UCP, Release notes, Versions
title: UCP release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for each UCP version.
You can then use [the upgrade instructions](installation/upgrade.md), to
upgrade your installation to the latest release.

## Version 1.1.6

(18 Jan 2017)

Note: UCP 1.1.6 supports Docker Engine 1.12 but does not use the built-in
orchestration capabilities provided by the Docker Engine with swarm mode enabled.
When installing this UCP version on a Docker Engine 1.12 host, UCP creates a
cluster using the older Docker Swarm v1.2.

**Security Update**

This patch contains the following security-related updates:

* Fixed an issue by which a high number of ping requests could result in temporarily
unresponsive UCP services
* Fixed an issue by which non-admin users with "View-Only" permissions could use
the undocumented private API to restart/stop/delete containers.
* Only admins are now allowed to tag, save, and load images as UCP/DTR system images
* UCP will now warn admins during installation if there is an open TCP port which
could be used to perform unauthorized actions on the cluster

These issues affect UCP version 1.1.5 and below. They were discovered by our
development team during internal testing.

We've revised our guidelines on access control permissions as well. Read
the [permissions levels section](user-management/permission-levels.md) for more details.

## Version 1.1.5

(8 Dec 2016)

Note: UCP 1.1.5 supports Docker Engine 1.12 but does not use the built-in
orchestration capabilities provided by the Docker Engine with swarm mode enabled.
When installing this UCP version on a Docker Engine 1.12 host, UCP creates a
cluster using the older Docker Swarm v1.2.

**Bug fixes**

* Fixed an issue where containers created using a network alias would be
rescheduled to another node on a `on-node-failure` event, but have incorrect
DNS entries
* Errors during an LDAP synchronization are now displayed in the UI
* UCP is now correctly configured after the controller service is restarted
* Usernames from both LDAP and managed users are now correctly normalized during
the authentication process.


## Version 1.1.4

(29 Sept 2016)

Note: UCP 1.1.4 supports Docker Engine 1.12 but does not use the built-in
orchestration capabilities provided by the Docker Engine with swarm mode enabled.
When installing this UCP version on a Docker Engine 1.12 host, UCP creates a
cluster using Docker Swarm v1.2.5.

**Bug fixes**

* Fixed an issue that prevented authentication tokens to be generated for
organization accounts
* Improved performance when deploying to a cluster with thousands of users
* Increased timeouts for adding new nodes and other network operations


## Version 1.1.3

Note: UCP 1.1.3 supports Docker Engine 1.12 but does not use the built-in
orchestration capabilities provided by the Docker Engine with swarm mode enabled.
When installing this UCP version on a Docker Engine 1.12 host, UCP creates a
cluster using Docker Swarm v1.2.5.

**Security Update**

Fixes a security issue by which a malicious user with limited privileges can
escalate their privileges to perform unauthorized actions on the cluster via
the API.

This issue affects deployments of Universal Control Plane versions 1.1.2 or
prior, and can only be used to gain access to the system by someone who already
has a UCP account.

This issue was discovered by our development team during internal testing.

**Features**

* Core
  * Upgraded Docker Swarm to 1.2.5
  * Non-admin users no longer have the ability to edit or delete UCP/DTR volumes
	and networks.
  * The Pull Image, Delete Image, Create Volume, Delete Volume, Create Network
    and Delete Network operations are now inaccessible to users with View Only
	default permissions or lower.

**Bug Fixes**

* Improved system performance when large numbers of overlay networks are deployed
  on the cluster.
* Fixed an issue which affected container rescheduling on clusters with overlay
  networks.
* Fixed an issue which affected synchronizing organization owners (admins) in
  LDAP when migrating from DTR 1.4.3 to 2.0.x
* Fixed an issue where UCP/DTR integration config was not loaded when UCP
  controller was restarted.
* Fixed an issue in the GUI where the sidebar does not display when first
  logging into UCP.
* Fixed an issue where volumes created through the UCP GUI did not correctly
  populate the labels field.

**Known Issues**

* This version of UCP cannot be installed on Engine 1.12 host with swarm mode
enabled, and is not compatible with swarm-mode based APIs, e.g. `docker service`.

## Version 1.1.2

Note: UCP 1.1.2 supports Docker Engine 1.12 but doesn't use the new clustering
capabilities provided by the Docker swarm mode. When installing this UCP version
on a Docker Engine 1.12, UCP creates a "classic" Docker Swarm 1.2.3 cluster.

**Features**

* Core
  * Upgraded etcd to version 2.3.6.
  * Upgraded rethinkDB to version 2.3.4.
  * The support dump generated by `dsinfo` now provides more information about
  the UCP deployment, that can be used by Docker support.

* docker/ucp image
  * It's now possible to generate a support dump directly from the CLI using the
  `support` command.
  * It's now possible to tune how often UCP's key-value store takes snapshots
  using the `docker/ucp install --kv-snapshot-count` option. This can be used
  with the `--kv-timeout` flag to tune the performance of the key-value store.
  [Learn more about tuning the key-value store](https://github.com/coreos/etcd/blob/master/Documentation/v2/tuning.md#snapshots)

* UI
  * The dashboard now notifies admin users when an update for UCP is available.
  * It's now possible to see which specific controllers need to have root CAs
  inserted to achieve high-availability.
  * It's now possible to filter images on the `Images` tab.

**Bug Fixes**

* Fixed an issue in which UCP failed to install in machines where the hostname
has more than 41 characters.
* Fixed an issue in which `ping` requests caused a memory leak in the
`ucp-controller` and `ucp-kv` containers.
* When installing in the CLI, UCP now displays the specified `ADMIN_USERNAME`
variable rather than just "admin".
* Fixed an issue where container owner label permissions took priority over access
label permissions when displaying a list of containers.
* Fixed an issue in which upgrading to UCP caused a user to still see an older
version in the UI.

**Known Issues**

* This version of UCP can't be installed on Engine 1.12 swarm-mode based
clusters, and is not compatible with swarm-mode based APIs, e.g. `docker service`.

## Version 1.1.1

**Features**

* Core
  * Upgraded Docker Swarm to version 1.2.3.
  * An administrator can now reset their password. Use the `docker/ucp-auth
  passwd` command for this.

* docker/ucp image
  * It's now possible to configure the election timeout of the UCP key-value store
  with the `docker/ucp install --kv-timeout` option. This is useful when running
  UCP across multiple regions. The heartbeat interval will be 1/10th of
  the specified election timeout value.
  [Learn more](https://coreos.com/etcd/docs/latest/tuning.html)
  * It's now possible to skip TLS verification when joining new nodes to the
  cluster, using the `docker/ucp join --insecure-fingerprint` option. However, to
  ensure your cluster is secure, don't use this option for normal UCP deployments.
  * The restore operation now supports `--interactive, -i` flags, which require a
  backup file to be mounted in `/backup.tar` instead of streamed through `stdin`.

* UI
  * When pulling images on the UCP UI, you can now provide login credentials for
  a private registry.
  * It's now possible to disable a user account, to make it easier to switch
  from managed authentication to LDAP and vice-versa.
  * Added a setting to submit usage reports without anonymizing data.
  * When failing to pull an image on the UCP UI, a feedback message is now
  displayed.
  * The Containers page now allows showing and hiding columns.
  * The Containers page now allows filtering for running, stopped, and system
  containers.

**Bug Fixes**

* Fixed an issue that prevented new nodes to be joined to a cluster, after
upgrading UCP from an older version to 1.1.0.
* Fixed an issue that prevented UCP from integrating with DTR for single-sign-on
when pushing/pulling images.
* When upgrading, configurations for user, teams, and organizations are now
preserved.
* When upgrading, version labels are correctly added to the containers.
* Improved error logs generated by the UP key-value store.
* The restore command now ensures the backup is not corrupt, that the UCP
cluster is healthy and is running the same or later version of UCP before
restoring.
* The restore command now works correctly on a freshly installed instance of
UCP, assuming the same host IP and a correct backup file.
* LDAP domain names are now case-insensitive for easier syncing.
* Fixed an issue that caused LDAP syncs to run every minute, after upgrading
UCP from an older version to 1.1.0.
* Fixed error by which user could get an "access denied" message when deploying
a container from the UI due to cached permission labels.
* Fixed issue where environment variables were not being passed to new containers
when "Allow users to deploy containers on UCP controllers" setting was disabled.

**Misc**

* Since container rescheduling has reached GA on Docker Swarm, you can use it
without having to install UCP with the `--swarm-experimental-flag`.
* UCP now requires a minimum of 2 GB of RAM per node, instead of 1.5 GB.
* During installation, UCP now warns you to only restart the Docker Engine
after joining all controller nodes to the cluster.

**Known Issues**

* When using UCP with a Docker Engine prior to 1.11.1-cs2, containers with a
restart policy set to `restart=always` and using an overlay network, may not
resume properly when the Docker daemon is restarted. Upgrade the Docker Engine
on your nodes to version 1.11.1-cs2 to fix this. This is especially important
when running UCP and DTR on the same nodes, and with high-availability.
* When attempting to restore a v1.1.0 backup on a new cluster installed with
the `fresh-install` flag, the restore operation may fail due to engine-discovery
configuration issues. You should create new backups after upgrading to v1.1.1.
* UCP fails to install in machines where the hostname has more than 41
characters. This will be fixed in a future release. (Fixed in UCP 1.1.2)

## Version 1.1.0

**Features**

* Core
  * UCP and DTR are now using a unified authentication service.
  * Users and teams created in UCP are displayed in DTR under the 'Datacenter'
  organization.
  * All controllers joined to the cluster now have replicated CAs. For this,
  you need to copy the root key material to controllers joined to the cluster,
  * All UCP components were compiled with Go 1.5.4 and 1.6 to address a
  security vulnerability in Go.
  * When joining nodes to the cluster, UCP automatically runs
  'engine-discovery' to configure the Docker Engine for multi-host networking.
  * If you're using Docker Engine 1.11 with default configurations, when joining
  new nodes to the cluster multi-host networking is automatically configured
  without needing to restart the Docker daemon.

* docker/ucp image
  * Added the 'backup' command to create backups of controller nodes.
  * Added the 'restore' command, to restore a controller node from a backup.
  * Added the 'regen-certs' command, to regenerate keys and certificates used on
  a controller node. You can use this for changing the SANS on the certificates
  or in case a CA is compromised.
  * Added the 'stop' and 'restart' commands, to stop and start UCP containers.
​
* UI
  * Now you can deploy apps from the UI using a docker-compose.yml file.
  * There's a new setting to prevent users from deploying containers to the UCP
  controller nodes.
  * Improved usability of LDAP configuration settings.
  * Images page no longer shows the sha256 id of each image ID.
  * User profiles now display default permissions.
  * Improved feedback when creating users and teams with invalid characters.
  * Added horizontal scrollbar to wide pages.

**Bug Fixes**

* Improved messages when installing UCP on a host with firewall rules.
* Images page no longer shows images generated from intermediate builds.
* Images page no longer hangs when pulling an image.
* Scaling a container from the UI now preserves parameters like 'net' and
'privileged'.
* Fixed `docker ps --filter` to filter containers correctly.


**Misc**

* All UCP containers now have the 'com.docker.ucp.version' label with their
upstream version or UCP version.
* When running docker/ucp in interactive mode, the parameters and environment
variables passed to the command are displayed.
* Renamed 'external-ucp-ca' flag to 'external-server-cert' for clarity.
The former name is deprecated but still available.
* UCP is automatically configured to use overlay networking. Make sure ports
4789 and 7946 are open for this to work.
* The new authentication service requires ports 12383-12386 to be open.

**Known Issues**

* After upgrading to version `1.1.0`, if you join new nodes to the cluster,
a success message is displayed, but that node will not be part of the
cluster. As a workaround, join new controller nodes before upgrading, or
perform a fresh installation of UCP 1.1.0.
* If you have an active login session in UCP and do an upgrade, you should force
refresh the browser or you may run into UI errors.
* When joining replicas to the cluster, you may be prompted to restart the
Docker daemon on that node. For a faster installation, only restart the Docker
daemon after joining all replicas.
* When deploying applications from the UI, using the `host` network option
might cause errors. If this happens, deploy the application from the CLI.
* UCP 1.1.0 may not integrate correctly with DTR for purposes of single-sign-on
for pushing/pulling images. It is recommended to upgrade to UCP 1.1.1 for this.

**Component Versions**

UCP 1.1.0 uses:

* cfssl 1.2.0
* Docker Compose 1.7.0
* Docker Swarm: 1.1.3
* etcd 2.2.5
* RethinkDB 2.3.0

## Version 1.0.4

**Security update**

Fixes a security issue by which a user can can obtain unauthorized access to
UCP via LDAP authentication.

## Version 1.0.3

Fixes a bug introduced by version 1.0.2 that was causing problems when a user
navigated to their profile page.

## Version 1.0.2

**Security update**

Fixes a security issue by which a non-admin user account can gain admin-level
privileges via the UCP API.

**Known Issues**

Non-admin users might have an error when navigating to their profile page. This
happens when the user is part of a team that has a label applied to it.


## Version 1.0.1

**Features**

* Core
  * Upgraded Swarm to 1.1.3
  * Improved support for `docker cp`
  * System CA pool fallback for secure DTR connections
  * Added `--swarm-experimental` option during UCP install

* UI
  * Can provide one-time credentials to deploy a container from a private registry in UI
  * Added checkbox to select all containers in Containers screen
  * Removed click handlers from UI elements containing checkboxes
  * Usernames and team names now need to be url-compatible
  * Several usability improvements to Team screen
  * Messages now display team name, instead of Id
  * Added support for Growl style notifications
  * Improved usability of Applications page, when there are no applications
  deployed
  * Several improvements to form validations
  * Improved error messages displayed when users try to pull an image with
  no name
  * Don't allow creating teams with the same name
  * Non-admin users can no longer see cluster overview in Dashboard screen
  * Page size control is no longer displayed when the list has few elements
  * Renamed 'Roles' to 'Permissions'

**Bug fixes**

* Users that are on a team and have permission set to 'None', can no longer see
containers
* Volume driver options are now being correctly sent to Docker Engine
* Fix bug with visibility to User containers with the owner the same as a label


## Version 1.0.0

**Features**

* Core
    * License is now required to add nodes
    * Improved access control system
    * /\_ping endpoint now checks the state of datastore and Swarm
    * Use mutual TLS in CFSSL
    * Improved access control for Docker Engine proxy
    * Added support for custom server certificates and user bundles
    * Users can now launch "private" containers if default permission is
    Restricted Control or greater

* UI
    * Pages for Containers, Images, and Applications are now consistent
    * Improved usability of LDAP configuration page
    * Logs are displayed during LDAP configuration
    * Users can now see their permissions and teams on their profile page
    * Improved license configuration
    * Improved error messages for restricted operations
    * Support for enabling and disabling DTR integration

**Bug fixes**

* Users only see volumes, images, and networks if they have permissions
* User default role now setup properly with LDAP authentication
* Fixed container privilege escalation in access control
* Fixed UI issue that caused errors in Safari

**Misc**

* UCP now uses a vendored UCP Swarm image
* Removed timestamps from controller logs
* Switched from 'Full Control' to 'Restricted Control' for managing non-container
resources

**Known issues**

In version 1.0.0 it's not possible to create containers on user-defined
bridge networks, using the UCP web app.
This happens because the UCP web app is using the \<node\>/\<network_name\> syntax,
which is not supported.

As a workaround, create the containers using the CLI and:

* Use only \<network_name\>, and let Swarm find the node with that network, or
* Use the network ID instead.


**Upgrade notes**

It's not possible to upgrade from previous versions to v1.0. If you've
participated in the Docker UCP beta program, you need to uninstall the beta
version, before installing v1.0.

To ensure a smooth transition process, start by uninstalling UCP from
the regular nodes, followed by the controller nodes. Also, make sure you
use `ucp uninstall` command from version 1.0:

    docker run --rm -it --name ucp -v /var/run/docker.sock:/var/run/docker.sock docker/ucp:1.0.0 uninstall -i

After uninstalling, you can [Install UCP on a sandbox](install-sandbox.md),
or [Install UCP for production](installation/install-production.md).
