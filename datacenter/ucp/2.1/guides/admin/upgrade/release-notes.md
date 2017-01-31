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

## Version 2.1.0 beta 2

(25 Jan 2017)

This version of UCP extends the functionality provided by Docker Engine 1.13.
Before installing or upgrading this version, you need to install Docker Engine 1.13
in the nodes that you plan to manage with UCP.

**New features**

* Core
  * You can now deploy secrets using Compose for Services

* UI/UX
  * You can now view and manage Stacks (Compose v3 applications consisting of
  services, volumes, and networks) directly from the UI
  * The HTTP Routing Mesh settings page now contains a table of all services
  using HRM, with details on parameters and health status
  * Admins can now view team membership in a user's details screen
  * You can now customize session timeouts length in Authentication config screen
  * Can now mount `tmpfs` or existing local volumes to a service via the UI
  * Added numerous tooltips throughout the UI

**Bug fixes**

* Core
  * Promoting a worker to a manager no longer requires manually pulling
  down `ucp-metrics` image on the node (fixes known beta1 issue)
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
  * Installing DTR no longer fails when "Enable scheduling on UCP controllers and
  DTR nodes" is unchecked.
  * HTTP Routing Mesh configuration status messages are now temporary and no
  longer overlap with text on a service's details screen.
  * Publishing a port to both TCP and UDP in a service via UI now works correctly
  * Removed Internal Scheme option when publishing a hostname route
  * Nodes now stay sorted after clicking a parameter to sort by in the Nodes screen
  * Metrics on the Nodes screen no longer incorrectly show `-` instead of 0% or 100%
  * In service deploy/edit screens, secrets have been moved to Environments tab

**Known issues**

When deploying compose files that use secrets, you need to specify all the parameters
for the service like:

```none
secrets:
  - source: "foo"
    target: "foo"
    uid: "0"
    gid: "0"
    mode: 0400
```

UCP returns an error if you only use the name of the secret.


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

**Known Issues**

* Promoting a worker to manager causes the promotion process to hang or take a
very long time, because the node will not pull the `ucp-metrics` image correctly.
The workaround is to use `docker pull docker/ucp-metrics:2.1.0-beta1` on all
nodes that you plan to promote to manager.
* Dashboard metrics may show inaccurately high usage values when using an AWS template
based deployment

**Version Compatibility**

UCP 2.1 requires minimum versions of the following Docker components:

* Docker Engine 1.13
* Compose 1.9
* Docker Remote API v. 1.25
