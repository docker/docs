<!--[metadata]>
+++
title ="Release Notes"
description="Release notes for Docker Universal Control Plane. Learn more about the changes introduced in the latest versions."
keywords = ["Docker, UCP", "Release notes", "Versions"]
[menu.main]
identifier="ucp_rnotes"
parent="mn_ucp"
weight="99"
+++
<![end-metadata]-->

# UCP Release Notes

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

After uninstalling, you can [Install UCP for evaluation](evaluation-install.md),
or [Install UCP for production](production-install.md).
