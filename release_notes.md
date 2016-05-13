<!--[metadata]>
+++
title ="Release Notes"
description="Release notes for Docker Universal Control Plane. Learn more about the changes introduced in the latest versions."
keywords = ["Docker, UCP", "Release notes", "Versions"]
[menu.main]
identifier="ucp_rnotes"
parent="mn_ucp"
weight=110
+++
<![end-metadata]-->

# UCP Release Notes

## Version 1.1.0

(4 May 2016)

**Features**

* Core
  * UCP and DTR are now using a unified authentication service,
  * Users and teams created in UCP are displayed in DTR under the 'Datacenter'
  organization,
  * All controllers joined to the cluster now have replicated CAs. For this,
  you need to copy the root key material to controllers joined to the cluster,
  * All UCP components were compiled with Go 1.5.4 and 1.6 to address a
  security vulnerability in Go,
  * When joining nodes to the cluster, UCP automatically runs
  'engine-discovery' to configure the Docker Engine for multi-host networking,
  * If you're using Docker Engine 1.11 with default configurations, when joining
  new nodes to the cluster multi-host networking is automatically configured
  without needing to restart the Docker daemon.

* docker/dtr image
  * Added the 'backup' command to create backups of controller nodes,
  * Added the 'restore' command, to restore a controller node from a backup,
  * Added the 'regen-certs' command, to regenerate keys and certificates used on
  a controller node. You can use this for changing the SANS on the certificates
  or in case a CA is compromised,
  * Added the 'stop' and 'restart' commands, to stop and start UCP containers.
​
* UI
  * Now you can deploy apps from the UI using a docker-compse.yml file,
  * There's a new setting to prevent users from deploying containers to the UCP
  controller nodes,
  * Improved usability of LDAP configuration settings,
  * Images page no longer shows the sha256 id of each image ID,
  * User profiles now display default permissions,
  * Improved feedback when creating users and teams with invalid characters,
  * Added horizontal scrollbar to wide pages.

**Bug Fixes**

* Improved messages when installing UCP on a host with firewall rules,
* Images page no longer shows images generated from intermediate builds,
* Images page no longer hangs when pulling an image,
* Scaling a container from the UI now preserves parameters like 'net' and
'privileged',
* Fixed `docker ps --filter` to filter containers correctly.


**Misc**

* All UCP containers now have the 'com.docker.ucp.version' label with their
upstream version or UCP version,
* When running docker/ucp in interactive mode, the parameters and environment
variables passed to the command are displayed,
* Renamed 'external-ucp-ca' flag to 'external-server-cert' for clarity.
The first is deprecated but still available.

**Known issues**

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

**Component Versions**

UCP 1.1.0 uses:

* cfssl 1.2.0
* Docker Compose 1.7.0
* Docker Swarm: 1.1.3
* etcd 2.2.5
* RethinkDB 2.3.0
​

**Upgrading**

[Learn how to upgrade to the latest version](installation/upgrade.md)

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

**Known issues**

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
    * Users can now launch "private" containers if default permission is Restricted Control or greater

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
* Switched from 'Full Control' to 'Restricted Control' for managing non-container resources

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
