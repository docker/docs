---
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: Docker, UCP, Release notes, Versions
title: UCP 2.1 release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for the latest UCP version.
You can then use [the upgrade instructions](index.md), to
upgrade your installation to the latest release.

## Version 2.1.0

(9 Feb 2017)

This version of UCP extends the functionality provided by CS Docker Engine
1.13. Before installing or upgrading this version, you need to install CS
Docker Engine 1.13 in the nodes that you plan to manage with UCP.

**New features**

* Core
  * Support for managing secrets (e.g. sensitive information such as passwords
  or private keys) and using them when deploying services. You can store secrets
  securely on the cluster and configure who has access to them, all without having
  to give users access to the sensitive information directly
  * Support for Compose yml 3.1 to deploy stacks of services, networks, volumes,
  and secrets.
  * HTTP Routing Mesh now generally available. It now supports HTTPS passthrough
  where the TLS termination is performed by your services, Service Name  Indication
  (SNI) extension of TLS, multiple networks for app isolation, and Sticky Sessions
  * Granular label-based access control for secrets and volumes
  (NOTE: unlike other resources controlled via label-based access control, a
  volume without a label is accessible by all UCP users with Restricted Control
  or higher default permissions)

* UI/UX
  * You can now view and manage application stacks directly from the UI
  * You can now view cluster and node level resource usage metrics
  * When updating a service, the UI now shows more information about the service status
  * Rolling update for services now have `failure-action` which you can use to
  * Several improvements to service lifecycle management
  specify rollback, pausing, or continuing if the update fails for a task
  * LDAP synching has more configuration options for extra flexibility
  * UCP now warns when the cluster has nodes with different Docker Engine versions
  * The HTTP routing mesh settings page now lists all services using the
  routing mesh, with details on parameters and health status
  * Admins can now view team membership in a user's details screen
  * You can now customize session timeouts in the authentication settings page
  * Can now mount `tmpfs` or existing local volumes to a service when deploying
  services from the UI
  * Added more tooltips to guide users on the above features

**Bug fixes**

* Core
    * HTTP routing mesh can now be enabled or reconfigured when UCP is configured
    to only run images signed by specific teams
    * Fixed an error in which `_ping` calls were causing multiple TCP connections
    to open up on the cluster
    * Fixed an issue in which UCP install occasionally failed with the error
    "failed to change temp password"
    * Fixed an issue where multiple rapid updates of HTTP Routing Mesh configuration
    would not register correctly
    * Demoting a manager while in HA configuration no longer causes the `ucp-auth-api`
     container to provide errors

* UI/UX
    * When creating a user, pressing enter on keyboard no longer causes problems
    * Fixed assorted icon and text visibility glitches
    * Installing DTR no longer fails when "Enable scheduling on UCP controllers and
    DTR nodes" is unchecked.
    * Publishing a port to both TCP and UDP in a service via UI now works correctly

**Known issues**


The `docker stats` command is sometimes wrongly reporting high CPU usage.
Use the `top` command to confirm the real CPU usage of your node.
[Learn more](https://github.com/docker/docker/issues/28941).


**Version compatibility**

UCP 2.1 requires minimum versions of the following Docker components:

* Docker Engine 1.13.0
* Docker Remote API 1.25
* Compose 1.9
