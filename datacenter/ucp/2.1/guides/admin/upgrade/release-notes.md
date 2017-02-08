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
1.13.0. Before installing or upgrading this version, you need to install CS
Docker Engine 1.13.1 in the nodes that you plan to manage with UCP.

**New features**

* Core
  * You can now deploy an application stack composed of multiple services using
  a compose file v3
  * Support for managing secrets like passwords of private keys, and using them
  when deploying services. You can configure who has access to configure secrets
  and use them in their applications, without having to given them access to the
  sensitive information directly
  * Official support for routing hostnames to services. It now supports HTTPS
  passthrough where the TLS termination is performed by your services. It also
  supports Service Name Indication (SNI) extension of TLS
  * Early access Windows Server support for worker nodes
  * You can now see node metrics like disk and memory usage
  * Added access control for volumes

* UI/UX
  * You can now view and manage application stacks directly from the UI
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
  * Added more tooltips to guide users

**Bug Fixes**

* Core
  * HTTP routing mesh can now be enabled or reconfigured when UCP is configured
  to only run images signed by specific teams
  * Fixed an error in which `_ping` calls were causing multiple TCP connections
  to open up on the cluster
  * Fixed an issue in which UCP install occasionally failed with the error
  "failed to change temp password"
  * Defining multiple HRM networks with overlapping subnets now correctly causes
  the HTTP Routing Mesh `ucp-hrm` service to fail.
  * Fixed an issue where multiple rapid updates of HTTP Routing Mesh configuration
  would not register correctly
  * With HTTP Routing Mesh, using the "default" backend option with an empty
  external route now works correctly
  * Volumes label-based access control now correctly supports volumes created
  via the `mount` format flag

* UI/UX
  * When creating a user, pressing enter on keyboard no longer causes problems
  * Fixed assorted icon and text visibility glitches
  * Installing DTR no longer fails when "Enable scheduling on UCP controllers and
  DTR nodes" is unchecked.
  * Publishing a port to both TCP and UDP in a service via UI now works correctly
  * Nodes now stay sorted after clicking a parameter to sort by in the Nodes screen


**Version Compatibility**

UCP 2.1 requires minimum versions of the following Docker components:

* Docker Engine 1.13.0
* Docker Remote API 1.25
* Compose 1.9
