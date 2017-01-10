---
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: Docker, UCP, Release notes, Versions
title: UCP 2.1 release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for the latest UCP version.
You can then use [the upgrade instructions](install/upgrade.md), to
upgrade your installation to the latest release.

## Version 2.1.0 beta 1

(10 Jan 2017)

This version of UCP extends the functionality provided by Docker Engine 1.13.
Before installing or upgrading this version, you need to install
Docker Engine 1.13 in the nodes that you plan to manage with UCP.

**New features**

* Core
  * Support for managing secrets like passwords of private keys, and using them
  when deploying services. You can configure who has access to configure secrets
  and use them in their applications, without having to given them access to the
  sensitive information directly
  * You can now deploy an application stack composed of multiple services using
  Docker Compose
  * Official support for routing hostnames to services. It now supports HTTPS
  passthrough where the TLS termination is performed by your services. It also
  supports Service Name Indication (SNI) extension of TLS
  * Early access Windows Server support for worker nodes
  * You can now see node metrics like disk and memory usage
  * Added access control for volumes

* UI/UX
  * Rolling update for services now have `failure-action` which you can use to
  specify rollback, pausing, or continuing if the update fails for a task
  * LDAP synching has more configuration options for extra flexibility
  * Several improvements to service lifecycle management
  * When updating a service, the UI now shows more information about the service status
  * The service details page now displays information about task errors
  * UCP now warns when the cluster has nodes with different Docker Engine versions

**Bug Fixes**

* Core
  * HTTP routing mesh can now be enabled or reconfigured when UCP is configured
  to only run images signed by specific teams.

* UI/UX
  * When creating a user, pressing enter on keyboard no longer causes problems
	* Fixed assorted icon and text visibility glitches

**Version Compatibility**

UCP 2.1 requires minimum versions of the following Docker components:

* Docker Engine 1.13
* Compose 1.9
* Docker Remote API v. 1.25
