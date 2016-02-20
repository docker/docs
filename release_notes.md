<!--[metadata]>
+++
title ="Release Notes"
description="Release notes for Docker Universal Control Plane. Learn more about the changes introduced in the latest versions."
keywords = ["Docker, UCP", "Release notes", "Versions"]
[menu.main]
parent="mn_ucp"
weight="99"
+++
<![end-metadata]-->

# UCP Release Notes

## Version 1.0

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

## Version 0.9

**Features**

* Allow editing user accounts
* Renamed 'role' to 'permission level'
* Improved the UI of the Container and Settings screens
* Added tooltips for contextual help
* The dashboard now shows the scheduling strategy being used

**Bug fixes**

* Fixed http 404 error when accessing UCP

**Other notes**

It's not possible to upgrade from a previous version to v0.9. If you've
already installed UCP, use the `--fresh-install` option with the `ucp install`
command, to do remove the old installation, and install v0.9.

## Version 0.8

**Features**

* LDAP/AD integration

    You can now choose between the built-in, LDAP, or Active Directory service
    for authentication. To change the authentication service, login into UCP
    with an administrator account, navigate to the Settings page, and click
    the Auth tab.

* DTR integration

    You can now configure UCP to connect to a Docker Trusted Registry version
    1.4.3 or higher.

* Teams and ACLs

    You can now apply labels to resources, and manage permissions based on
    labels. You can also create teams to apply the same permissions to
    multiple users.

* Multi-host networking

    The `ucp` install tool now lets you set up multi-host networking.
    For more information run `docker run --rm docker/ucp engine-discovery --help`.

* UI

    Overall changes to the UI to make UCP easier to use. The UI for managing
    teams and LDAP integration was improved.


**Other notes**

This version of UCP requires Docker Engine 1.10.0-rc1 or higher. It was also
changed to use Swarm 1.1.0-RC2, and Etcd 2.2.4 internally.

**Known issues**

If you've upgraded from a previous version, you might have
access control problems, when using non-admin users.
As a work around, recreate those users and delete the old ones.
