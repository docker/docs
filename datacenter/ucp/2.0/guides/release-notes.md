---
title: UCP release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords:
- Docker, UCP
- Release notes
- Versions
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for the latest UCP version.
You can then use [the upgrade instructions](installation/upgrade.md), to
upgrade your installation to the latest release.

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
