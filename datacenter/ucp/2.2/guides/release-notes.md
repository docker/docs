---
title: UCP 2.1 release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: Docker, UCP, release notes
---

Here you can learn about new features, bug fixes, breaking changes and
known issues for the latest UCP version.
You can then use [the upgrade instructions](admin/install/upgrade.md), to
upgrade your installation to the latest release.

## Version 2.2.0 Beta

(17 July 2017)

**New Features**

* The role-based access control system has been  overhauled for additional
granularity and customization. Admins now define access control through Grants,
a 1:1:1 mapping of a Subject, a Role, and a Collection:
  * Subject: A user, team, or organization
  * Role: A set of permissions. In addition to the existing predefined roles,
  admins can now create custom roles with their choice of permissions taken
  from the full Docker API
  * Collection: A group of containers or container-based resources (e.g. volumes,
  networks, secrets, etc.). Collections have a hierarchical directory-like structure
  and replace the old access control labels from the previous system (though they
  still use labels in the CLI.
  * Please read the documentation <here> for more information on the new system
  and how your old access control settings are migrated during an upgrade
* (EE Advanced only) UCP now provides RBAC for nodes, where an admin can enforce
physical isolation between users on different nodes in the cluster. This means two
different teams can only view and deploy on the nodes to which they have access.
* Enhancements to the user management system:
  * UCP now supports the user concept of organizations, which are groups of teams
  * Users can now specify a default collection which automatically applies
  access control labels to all CLI deploy commands when no label is specified
  by the user
* Support for UCP workers running Windows Server 2016, and the ability to deploy
Windows-based containerized applications on the cluster
  * Please read the documentation <here> for instructions on how to join
  Windows nodes, and current limitations when deploying Windows applications
* Support for UCP workers running Linux on IBM Z systems
* UCP now provides a public, stable API for cluster configuration and access control,
and the API is fully interactive within the UCP UI.
* The UCP UI has been redesigned for ease-of-use and data management:
  * Redesigned dashboard with time-series historical graphs for usage metrics
  * Compact layout to more easily view resource information at a glance
  * Detail panels for resources no longer slide out and cover the main panel

**Known Issues**
