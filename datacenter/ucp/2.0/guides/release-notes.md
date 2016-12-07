---
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: Docker, UCP, Release notes, Versions
title: UCP release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for the latest UCP version.
You can then use [the upgrade instructions](installation/upgrade.md), to
upgrade your installation to the latest release.

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
