---
title: UCP 2.2 release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: UCP, release notes
---

> **UCP 3.0 now available**. [You can check the release notes here](/ucp/dtr/release-notes.md).

Here you can learn about new features, bug fixes, breaking changes, and
known issues for the latest UCP version.
You can then use [the upgrade instructions](admin/install/upgrade.md), to
upgrade your installation to the latest release.

## Version 2.2.9 (2018-04-17)

**Bug fixes**

* Security
  * Fixed an issue that allows users to incorrectly interact with local volumes.
  * Fixed an issue where setting minimum TLS version that causes `ucp-agent` to
   keep restarting on worker nodes.

* Core
  * Fixed an issue that causes container fail to start with `container ID not found`
   during high concurrent API calls to create and start containers.

## Version 2.2.7 (2018-03-26)

**Bug fixes**

* Fixed an issue where the minimum TLS version setting is not correctly handled,
  leading to non-default values causing `ucp-controller` and `ucp-agent` to keep
  restarting.

## Version 2.2.6 (2018-03-19)

**New features**

* Security
  * Default TLS connections to TLS 1.2, and allow users to configure the minimum
  TLS version used by the UCP controller.
* Support and troubleshoot
  * The support dump now includes the output of `dmesg`.
  * Added more information to the telemetry data: kernel version, graph driver, and
  logging driver.
  * The `dsinfo` image used for support dumps is now smaller.

**Bug fixes**

* Core
  * The HRM service is no longer deployed with constraints that might prevent
  the service from ever getting scheduled.
  * Fixed a problem causing the HRM service to be restarted multiple times.
  * The `ucp-agent` service is now deployed without adding extra collection labels.
  This doesn't change the behavior of the service.
  * Fixed problem causing a healthy `ucp-auth-store` component to be reported as
  unhealthy.
  * Fixed a race condition causing the labels for the UCP controller container
  to be reset.
  * Fixed an issue causing the `ucp-agent` service to be deployed with the wrong
  architecture on Windows nodes.
* RBAC
  * Role-based access control can now be enforced for third-party volume plugins,
  fixing a known issue from UCP 2.2.5.
  * Admins can now clean up volumes and networks that had inconsistent collection
  labels across different nodes in the cluster. Previously, they would have had
  to go onto each node and clean up those resources directly.
  * When upgrading from UCP 2.1, inactive user accounts are no longer migrated
  to the new RBAC model.
  * Fixed an issue preventing users from seeing a collection when they have
  permissions to deploy services on a child collection.
  * Grants are now deleted when deleting an organization whose teams have grants.
* UI
  * Fixed a problem in the Settings page that would cause Docker to stop when
  you made changes to UCP settings and a new manager node is promoted to leader.
  * Fixed bug causing the Grants list page not to render after deleting an
  organization mentioned used on a grant.
  * Fixed a problem that would intermittently cause settings not to be persisted.
  * Fixed an issue that prevented users from being able to change LDAP settings.

**Known issues**

* RethinkDB can only run with up to 127 CPU cores.
* When integrating with LDAP and using multiple domain servers, if the
default server configuration is not chosen, then the last server configuration
is always used, regardless of which one is actually the best match.


## Version 2.2.5 (16 January 2018)

**Bug fixes**

* Role-based access control is now enforced for volumes managed by
third-party volume plugins. This is a critical security fix for customers that
use third-party volume drivers and rely on Docker Universal Control Plane for
tenant isolation of workloads and data.
**Caution** is advised when applying this update because users or automated
workflows may have come to rely on lack of access control enforcement when
manipulating volumes created with 3rd party volume plugins.

**Known issues**

* UCP role-based access control is not compatible with all third-party volume
plugins. If you’re using certain third-party volume plugins (such as Netapp)
and are planning on upgrading UCP, you can skip 2.2.5 and wait for the upcoming
2.2.6 release, which will provide an alternative way to turn on RBAC enforcement
for volumes.

## Version 2.2.4 (2 November 2017)

**News**

* Docker Universal Control Plane now supports running managers on IBM Z on RHEL, SLES and Ubuntu. Previously, only workers were supported on IBM Z.

**Bug fixes**

* Core
  * `ucp-etcd` system images are now hidden. Previously, these system images were erroneously displayed in the images list.
  * `disable_usageinfo` will now disable usage metrics. A regression caused this setting to not be respected.
  * UCP now outputs "Initializing..." log messages during setup so that administrators can establish that setup or install has begun.
  * Windows worker promotion is now blocked. Previously, Windows workers could be promoted using the CLI, which would fail.
  * Loading gzipped images with the Docker CLI is now supported. This would previously cause a panic.
  * Permissions are now checked when filtering nodes by container. Previously, permissions were not considered.
  * An LDAP sync is now triggered as soon as an LDAP user is lazy-provisioned. Previously, lazy-provisioned users would not immediately be added to teams and orgs.

* UI/UX
  * License page now shows all capabilities. Previously it was not clear if a license supported Docker image scanning or not.
  * Additional translations added for internationalization.
  * UI for adding users to teams simplified.
  * The grant list can now sorted and pagination in the grants view has been improved. The grants view previously had glitches on systems with many grants.
  * Fixed problem where UI would hang when pulling images.
  * "Max failure ratio" and "Failure action" re-introduced in service definitions. These settings were not available in UCP 2.2, but were available in previous UCP versions.
  * Collection labels are no longer applied to UCP system services. UCP previously auto-applied labels, which was confusing.

**Known issues**

 * Docker currently has limitations related to overlay networking and services using VIP-based endpoints. These limitations apply to use of the HTTP Routing Mesh (HRM). HRM users should familiarize themselves with these limitations. In particular, HRM may encounter virtual IP exhaustion (as evidenced by `failed to allocate network IP for task` Docker log messages). If this happens, and if the HRM service is restarted or rescheduled for any reason, HRM may fail to resume operation automatically. See the Docker EE 17.06-ee5 release notes for details.
 * The Swarm admin UI for UCP versions 2.2.0 and later contain a bug. If used with Docker Engine version 17.06.2-ee5 or earlier, attempting to update "Task History Limit", "Heartbeat Period" and "Node Certificate Expiry" settings using the UI will cause the cluster to crash on next restart. Using UCP 2.2.X and Docker Engine 17.06-ee6 and later, updating these settings will fail (but not cause the cluster to crash). Users are encouraged to update to Docker Engine version 17.06.2-ee6 and later, and to use the Docker CLI (instead of the UCP UI) to update these settings. Rotating join tokens works with any combination of Docker Engine and UCP versions. Docker Engine versions 17.03 and earlier (which use UCP version 2.1 and earlier) are not affected by this problem.

## Version 2.2.3 (13 September 2017)

**Bug fixes**

* Core
  * Node list will no longer show duplicated worker node entries.
  * Volume mount options are no longer dropped when creating volumes.
  * `docker stack deploy` with secrets specified in docker-compose file now works.
* UI/UX
  * Upgrade button is now greyed out and deacticated after initiating upgrade.
  * If an error is encountered while creating a service, the UI no longer freezes.
  * Upgrade notification fixed to have working link.
  * "Default Role For All Private Collections" can now be updated. Updating this
  role in the UI previously had no effect.
  * Added notification to UI to show that an upgrade is in progress.
  * Client bundle can now be downloaded with Safari browser.
  * Windows nodes are no longer displayed in the DTR install UI.
  * DTR settings state in UCP is now preserved when switching tabs. Previously,
  un-saved state was lost when switching tabs.
  * Fixed problem where first manager node may have IP address `0.0.0.0`,
  causing dashboard to not update.
  * UI for adding Windows nodes improved to include full join instructions.
  * Node Task UI fixed. Displaying tasks for a node previously did not work.
  * LDAP settings UI improved. Sync interval setting is now validated, a
  never-ending update spinner been fixed and it's UI action sequencing bugs have
  been fixed so that it's now possible to disable LDAP.
  * Uploading Docker images in the UI now has better error messages and improved
  validation.
  * Containers removed in UI are now force-removed. Previously removing
  containers would fail.
  * DTR install instructions `--ucp-url` parameter fixed to have valid value.
  * Deleting multiple users in succession fixed. Previously, an error would
  result when deleting more than one user at a time.
  * Added validation when adding DTR URL in UCP admin settings.
  * Left-nav now shows resource counts, addressing an UI regression from UCP 2.1.

**Known issues**

 * Upgrading heterogeneous swarms from CLI may fail because x86 images are used
 instead of the correct image for the worker architecture.
 * Agent container log is empty even though it's running correctly.
 * Rapid UI settings updates may cause unintended settings changes for logging
 settings and other admin settings.
 * Attempting to load an (unsupported) `tar.gz` image results in a poor error
 message.
 * Searching for images in the UCP images UI doesn't work.
 * Removing a stack may leave orphaned volumes.
 * Storage metrics are not available for Windows.
 * You can't create a bridge network from the web UI. As a workaround use
 `<node-name>/<network-name>`.


## version 2.2.2 (30 August 2017)

**Bug fixes**

* Core
  * Fixed an issue that caused timeouts during install, preventing UCP 2.2.1 from
  being released.
  * Fixed a number of issues in which access control labels and roles could not
  be upgraded to their new format, when upgrading UCP.
  [Learn more](https://success.docker.com/KBase/Auth_system_migration_errors).
  * Fixed an issue that caused an upgrade with multiple manager nodes to fail
  with RethinkDB startup errors.
  * Fixed an issue that caused upgrades to fail due to UCP being unable to
  remove and replace older UCP containers.
  * Fixed an issue in which upgrade timed out due to lack of available disk space.
  * Fixed an issue in which rescheduling of containers not belonging in services
  could fail due to a request for a duplicate IP address.
  * DTR containers are no longer omitted from `docker ps` commands.
* UI/UX
  * Fixed known issue from 2.2.0 where config changes (including LDAP/AD) take
  an extended period to update after making changes in the UI settings.
  * Fixed an issue where the `/apidocs` url redirected to the login page.
  * Fixed an issue in which the UI does not redirect to a bad URL immediately
  after an upgrade.
  * Config and API docs now show the correct LDAP sync cron schedule format.
* docker/ucp image
  * Support dump now contains information about access control migrations.
  * The `ucp-auth-store` and `ucp-auth-api` containers now report health checks.

**Known issues**

* When deploying compose files that use secrets, the secret definition must
include `external: true`, otherwise the deployment fails with the error
`unable to inspect secret`.

## Version 2.2.0 (16 August 2017)

**New features**

* The role-based access control system has been overhauled for additional
granularity and customization. Admins now define access control through Grants,
a 1:1:1 mapping of a Subject, a Role, and a Collection:
  * Subject: A user, team, or organization.
  * Role: A set of permissions. In addition to the existing predefined roles,
  admins can now create custom roles with their choice of permissions taken
  from the full Docker API.
  * Collection: A group of containers or container-based resources (e.g. volumes,
  networks, secrets, etc.). Collections have a hierarchical directory-like structure
  and replace the old access control labels from the previous system (though they
  still use labels in the CLI).
  * [Read the documentation](access-control/index.md#transition-from-ucp-21-access-control)
   for more information and examples of the new system and how your old access
   control settings are migrated during an upgrade.
* UCP now provides access control for nodes, where an admin can enforce
physical isolation between users on different nodes in the cluster. This means two
different teams can only view and deploy on the nodes to which they have access.
This is only available with an EE Advanced license.
* Enhancements to the user management system:
  * UCP now supports the user concept of organizations, which are groups of teams.
  * Users can now specify a default collection which automatically applies
  access control labels to all CLI deploy commands when no label is specified by
  the user.
* Support for UCP workers running Windows Server 2016, and the ability to deploy
Windows-based containerized applications on the cluster.
  * [Read the documentation](admin/configure/join-windows-worker-nodes/index.md)
  for instructions on how to join Windows nodes, and current limitations when
  deploying Windows applications.
* Support for UCP workers running on IBM Z systems with RHEL 7.3, Ubuntu 16.04,
and SLES 12.
* UCP now provides a public, stable API for cluster configuration and access control,
and the API is fully interactive within the UCP UI.
* Support for using services with macvlan networks and configuring network scope in UI.
* The UCP UI has been redesigned for ease-of-use and data management:
  * Redesigned dashboard with time-series historical graphs for usage metrics.
  * Compact layout to more easily view resource information at a glance.
  * Detail panels for resources no longer slide out and cover the main panel.
  * Filtering mechanism to display related items (e.g. resources in a collection or stack).

**Known issues**

* UI issues:
  * Cannot currently remove nodes using UCP UI. Workaround is to remove from CLI
  instead.
  * Search does not function correctly for images.
  * Cannot view label constraints from a collection's details pages. Workaround
  is to view by editing the collection.
  * Certain config changes to UCP make take several minutes to update after making
  changes in the UI. In particular this affects LDAP/AD configuration changes.
  * Turning `LDAP Enabled` from "Yes" to "No" disables the save button. Workaround
  is to do a page refresh which completes the configuration change.
  * Removing stacks from the UI may cause certain resources to not be deleted,
  including networks or volumes. Workaround is to delete the resources directly.
  * When you create a network and check 'Enable hostname based routing', the web
  UI doesn't apply the HRM labels to the network. As a workaround,
  [create the network using the CLI](https://docs.docker.com/datacenter/ucp/2.2/guides/user/services/use-domain-names-to-access-services/#service-labels).
  * The web UI does not currently persist changes to session timeout settings.
  As a workaround you can update the settings from the CLI, by [adapting these instructions for the
  session timeout](https://docs.docker.com/datacenter/ucp/2.2/guides/admin/configure/external-auth/enable-ldap-config-file/).
* docker/ucp
  * The `support` command does not currently produce a valid support dump. As a
  workaround you can download a support dumps from the web UI.
* Windows issues
  * Disk related metrics do not display for Windows worker nodes.
  * If upgrading from an existing deployment, ensure that HRM is using a non-encrypted
  network prior to attaching Windows services.
