---
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: Docker, UCP, Release notes, Versions
title: UCP release notes
redirect_from:
- /ucp/release_notes/
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for the latest UCP version.
You can then use [the upgrade instructions](installation/upgrade.md), to
upgrade your installation to the latest release.

## Version 2.0.4

(17 July 2017)

**Security Update**

* Remediated a privilege escalation where an authenticated user could obtain
admin-level privileges

This issue affects UCP versions 2.0.0-2.0.3 and 2.1.0-2.1.4. The were discovered
by our development team during internal testing.

## Version 2.0.3

(8 Feb 2017)

**Bug Fixes**

* Core
  * Label-based access control now supported for volumes. Unlike other
    resources controlled via label-based access control, a volume without a
    label is accessible by all UCP users with Restricted Control or higher
    default permissions.)
  * Demoting a manager while in HA configuration no longer causes the `ucp-auth-api`
  container to provide errors
  * Improved the performance of joining new nodes or promoting nodes to UCP managers
  * LDAP integration now supports iterating through group members to avoid
  exceeding LDAP length limits
  * UCP now detects if Docker is accessible on ports 2375-2376, and warns you to
  reconfigure Docker Engine if that's the case

* UI/UX
  * The users page now allows filtering by active, inactive, and admin users
  * Settings page now shows a list with the latest UCP versions that you can upgrade to
  * Improved behavior of upgrading from the web UI to make it more consistent
  with the behavior of upgrading through the CLI
  * UI now correctly displays environment variables with values containing `=`

* docker/ucp image
  * Fixed a problem that caused the upgrade command to fail when upgrading a
cluster using the HTTP routing mesh



## Version 2.0.2

(18 Jan 2017)

**Security update**

This patch contains the following security-related updates:

* Fixed an issue by which a high number of ping requests could result in temporarily
unresponsive UCP services
* Fixed an issue by which non-admin users with "View-Only" permissions could use
the undocumented private API to restart/stop/delete containers.
* Only admins are now allowed to tag, save, and load images as UCP/DTR system images
* UCP will now warn admins during installation if there is an open TCP port which
could be used to perform unauthorized actions on the cluster

These issues affect UCP version 2.0.1 and 2.0.0. They were discovered by our
development team during internal testing.

We've revised our guidelines on access control permissions as well. Read the
[permissions levels section](user-management/permission-levels.md) for more details.

**Features**

* Core
	* Label-based access control now supported for volumes.
	(NOTE: unlike other resources controlled via label-based access control, a
	volume without a label is accessible by all UCP users with Restricted Control
	or higher default permissions.)
	* Authentication now supports LDAP servers that don't use `memberOf`. Instead
	can look within LDAP groups and sync members in that group.

* docker/ucp image
	* Can now add input the text of a license file at install time with the flag
	`--license "cat <license_name>.lic"`

* UI/UX
	* Can now configure labels, drives, and other options when mounting volumes
	while using Deploy Service wizard
	* Task errors are now shown in each service's details page in the GUI
	* Can now add individual container labels when using Deploy Service wizard

**Bug Fixes**

* Core
	* Setting `network_mode=host` or `--net=host` no longer causes container
	scheduling and websocket errors.
	* Inspecting networks now correctly shows attached containers across the entire
	cluster
	* UCP now prompts you to rename overlay networks which have illegal characters
	during an upgrade of the cluster
	* HTTP Routing Mesh now correctly upgrades to latest version when UCP is upgraded
	* UCP images now pulled correctly on worker nodes when using Docker for Azure template

	* UCP now correctly ensures that a non-admin user designated as an LDAP recovery
	admin is promoted to admin status within the platform


* UI/UX
	* Fixed an issue preventing non-admins from creating volumes in the GUI
	* Fixed hyperlink in banner for upgrading UCP
	* SANs added during installation with `--interactive` flag now appear correctly
	in each manager node's details page in the GUI
	* Clarified banner warning provided mid-upgrade when certain nodes are running
	different versions of UCP
	* Fixed an error where "Show System Services" toggle did not function correctly
	* Users now prevented from using "Drain" option on manager nodes

* docker/ucp image
	* Fixed an issue preventing `stop` and `restart` commands from working correctly
	* UCP install no longer occasionally stalls with error "failed to change temp password"

## Version 2.0.1

(22 Nov 2016)

**Features**

* UI/UX
  * The node details page now shows information about the node's CPU and RAM
  * Improved applications page to provide more guidance when there are no apps deployed

**Bug fixes**

* Core
  * Fixed an issue with rethinkDB sync that causes timeout failures during upgrades
  * HTTP routing mesh no longer crashes if a routed service fails. It also
  provides better error messages in the CLI

* docker/ucp image
  * Install and upgrade timeouts were increased to ensure the swarm managed by
  UCP is fully operational
  * Using `--external-server-cert` at install time now correctly preserves
  pre-existing server certificates
  * Joining a node that's using a non-default Docker runtime path now works as
  expected

* UI/UX
  * Several buttons and links are now hidden if the user has "Read-Only"
  permissions, e.g. deploy application, create volumes, create images
  * When importing users from LDAP, the change password fields are no longer
  displayed in the user profile
  * When integrating with LDAP, the LDAP reader password is no longer displayed
  on the UI or HTML
  * Clarified that service resource constraints use absolute instead of nano CPU
  shares. This is consistent with the Docker CLI
  * UI now prompts for confirmation when switching from LDAP to built-in
  authentication
  * Improved DTR integration screen, to provide more guidance on how to install
  DTR

**Known issues**

* When deploying applications from the UI or CLI with the compose file format
2.1, overlay networks will not be created in attachable mode, so containers
will fail to attach to that network. As a workaround you can create the networks
upfront and make them attachable, and change your compose file to use those
networks.

## Version 2.0.0

(10 Nov 2016)

**Features**

* Native support for Docker Engine 1.12 running in swarm mode
* Support for declaring the application desired state with `docker service`
* Backwards compatibility for using `docker run`, `docker-compose`,
`docker logs`, and `docker exec` commands on UCP
* Pause and drain nodes for putting nodes in maintenance
* Specify container health checks in a Dockerfile and view container status
in the web UI
* Granular label-based access control for services and networks
* Use Docker content trust to enforce only running images signed by members
of specific teams
* Built-in TCP load balancing and service discovery for services
* Added an HTTP routing mesh for enabling hostname routing for services
(experimental)
* The UCP web UI now lets you know when a new version is available, and upgrades
to the new version with a single click

**Installer**

* You can now install UCP on an existing Docker Engine 1.12 swarm, by running
the UCP install command on manager node. Requires CS Engine 1.12 or higher.
* It's now possible to install a new UCP cluster using a backup using the
`restore` command. This preserves the UCP cluster configurations but not any
application services or networks from the old cluster.

**Web UI**

* Web UI has been significantly redesigned for this release
* Add, remove, pause, and drain nodes from the web UI and see their status
in real-time
* One-click promotion of worker nodes to managers for high availability
* You can now reconfigure certificates, SANs, and key-value store parameters
from the UI
* Added getting started guides and tooltips for guiding new users on the
platform
* Added instructions to the web UI on how to install and integrate with DTR
* Added a wizard to make it simpler to deploy Docker services from the UI

**General improvements**

* UCP and DTR system resources are now hidden and protected from non-admin users
* Improved support dump with more system information, readable file formats,
and container names/IDs
* LDAP/AD now lists up to 1 million users

**Known issues**

* When deploying applications from the UI or CLI with the compose file format
2.1, overlay networks will not be created in attachable mode, so containers
will fail to attach to that network. As a workaround you can create the networks
upfront and make them attachable, and change your compose file to use those
networks.
