---
title: UCP release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: UCP, release notes
toc_min: 1
toc_max: 2
redirect_from:
  - /datacenter/ucp/2.2/guides/release-notes/
  - /datacenter/ucp/3.0/guides/release-notes/
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for the latest UCP version.
You can then use [the upgrade instructions](admin/install/upgrade.md) to
upgrade your installation to the latest release.

* [Version 3.0](#version-30)
* [Version 2.2](#version-22)

# Version 3.0

## 3.0.4 (2018-08-09)

**Bug fixes**

* Security
  * Fixed a critical security issue where the LDAP bind username and password 
    were stored in clear text on UCP hosts. Please refer to [this KB article](https://success.docker.com/article/upgrading-to-ucp-2-2-12-ucp-3-0-4/) for proper implementation of this fix.

**Known Issue**

* You must manually pull `docker/ucp-agent:3.0.4` in the images section of the web UI before upgrading. Alternately, you can just pull `docker/ucp-agent:3.0.4` on every manager node.

## 3.0.3 (2018-07-26)

**New platforms**
* UCP now supports running Windows Server 1803 workers
   * Server 1803 ingress routing, VIP service discovery, and named pipe mounting are not supported in this release.
   * Offline bundles `ucp_images_win_1803_3.0.3.tar.gz` have been added.
* UCP 3.0.3 now supports IBM Z (s390x) as worker nodes on 3.0.x for SLES 12 SP 3. Interlock is currently not supported for 3.0.x on Z.

**Bug Fixes**
* Core
   * Optimize swarm service read api calls through UCP
   * Fixes an issue where some UCP Controller API calls may hang indefinitely.
   * Default Calico MTU set to 1480
   * Calico is upgraded to 3.0.8
   * Compose for Kubernetes logging improvements
   * Fixes an issue where backups would fail if UCP was not licensed.
   * Fixes an issue where DTR admins are missing the Full Control Grant against /Shared Collection even though they have logged in at least once to the UI.
   * Add support for bind mount volumes to kubernetes stacks and fixes sporadic errors in kubernetes stack validator that would incorrectly reject stacks.


## 3.0.2 (2018-06-21)

**New Features**
* UCP now supports running Windows Server 1709 workers
   * Server 1709 provides smaller Windows base image sizes, as detailed [here](https://docs.microsoft.com/en-us/windows-server/get-started/whats-new-in-windows-server-1709)
   * Server 1709 provides relaxed image compatibility requirements, as detailed [here](https://docs.microsoft.com/en-us/virtualization/windowscontainers/deploy-containers/version-compatibility)
   * Server 1709 ingress routing, VIP service discovery, and named pipe mounting are not supported in this release.
   * Offline bundle names are changed from `ucp_images_win_3.0.1.tar.gz` to `ucp_images_win_2016_3.0.2.tar.gz` and `ucp_images_win_1709_3.0.2.tar.gz` based on Windows Server versions.
* UCP now supports running RHEL 7.5.  Please refer to the [compatibility matrix](https://success.docker.com/article/compatibility-matrix).
* Added support for dynamic volume provisioning in Kubernetes for AWS EBS and
Azure Disk when installing UCP with the `--cloud-provider` option.

**Bug Fixes**
* Core
   * Fixed an issue for anonymous volumes in Compose for Kubernetes.
   * Fixed an issue where a fresh install would have an initial per-user session
   limit of unlimited rather than the expected limit of 10 minutes.
   * Added separate resource types for Kubernetes subresources (e.g. pod/log)
   so that users can get separate permissions for those resources, as with the
   built-in Kubernetes RBAC authorizer. If you had a custom role with
   (for instance) Pod Get permissions, you may need to create a new custom
   role with permissions for each new subresource.
   * To deploy Pods with Privileged options, users now require a grant with the
   role `Full Control` for all namespaces.
   * The `/api/ucp/config` endpoint now includes default node orchestrator.
   * Added `cni_mtu` and `cni_ipinip_mtu` settings to UCP config for controlling
   MTU sizes in Calico.
   * When a route is not found, interlock now returns a 503 (the expected
     behavior) instead of a 404.

* UI/UX
   * Fixed an issue that causes LDAP configuration UI to not work properly.

## 3.0.1 (2018-05-17)

**Bug Fixes**
* Core
  * Bumped Kubernetes version to 1.8.11.
  * Compose for Kubernetes now respects the specified port services are exposed on.
  This port must be in the `NodePort` range.
  * Kubernetes API server port is now configurable via `--kube-apiserver-port`
  flag at install or `cluster_config.kube_apiserver_port` in UCP config.
  * Fixed an issue where upgrade fails due to missing `ucp-kv` snapshots.
  * Fixed an issue where upgrade fails due to layer 7 routing issues.
  * `ucp-interlock-proxy` no longer tries to schedule components on Windows nodes.
  * Fixed an issue where a Kubernetes networking failure would not stop UCP from
  installing successfully.
  * Fixed an issue where encrypted overlay networks could not communicate on
  firewalled hosts.
  * Fixed an issue where Pod CIDR and Node IP values could conflict at install
  Installation no longer fails if an empty `PodCIDR` value is set in the UCP
  config at install time. Instead, it falls back to default CIDR.

*  UI/UX
  * Fixed an issue where UCP banners redirected to older UCP 2.2 documentation.


**Known issues**

* Encrypted overlay networks may not work after upgrade from 3.0.0.  Apply the following to
  all the nodes after the upgrade.
    ```
    iptables -t nat -D KUBE-MARK-DROP -j MARK --set-xmark 0x8000/0x8000
    iptables -t filter -D KUBE-FIREWALL -m comment --comment "kubernetes firewall for dropping marked packets" -m mark --mark 0x8000/0x8000 -j DROP
    ```
 * `ucp-kube-controller-manager` emits a large number of container logs.
 * Excessive delay is seen when sending `docker service ls` via UCP client
   bundle on a cluster that is running thousands of services.
 *  Inter-node networking may break on Kubernetes pods while the `calico-node`
  pods are being upgraded on each node. This may cause up to a few minutes of
  networking disruption for pods on each node during the upgrade process,
  depending on how quickly `calico-node` gets upgraded on those nodes.
  * `ucp-interlock-proxy` may fail to start when two or more services are
 configured with two or more backend hosts.  [You can use this workaround](https://success.docker.com/article/how-do-i-ensure-the-ucp-routing-mesh-ucp-interlock-proxy-continues-running-in-the-event-of-a-failed-update).

## Version 3.0.0 (2018-04-17)

The UCP system requirements were updated with 3.0.0. Make sure to
[check the system](https://docs.docker.com/ee/ucp/admin/install/system-requirements/)
requirements before upgrading.

### Orchestration

* UCP now supports Kubernetes as an orchestrator, in addition to the existing
Swarmkit and "classic" Swarm orchestrators. Kubernetes system components are
automatically deployed on all manager and Linux worker nodes managed by UCP.
[Learn more about Kubernetes support](https://docs.docker.com/ee/ucp/ucp-architecture/).
* Worker nodes running Linux on amd64 architectures can be configured to run
only Swarm workloads, only Kubernetes workloads, or mixed workloads. Manager
nodes are by default Mixed in order to support Swarm and Kubernetes system
components. However, it is not recommended to run Worker nodes as Mixed due to
potential resource contention issues.
* Users can deploy Kubernetes workloads through the web UI, and the CLI using
a UCP client bundle and `kubectl`.
[Learn more](https://docs.docker.com/ee/ucp/kubernetes/).
* Users can now use Compose to deploy Kubernetes workloads from the web UI.
[Lean more](https://docs.docker.com/ee/ucp/kubernetes/deploy-with-compose/).

### Networking

* UCP includes Calico as the default CNI plugin for networking of Kubernetes
applications. [Learn more](https://docs.projectcalico.org/v3.1/introduction/).
The following Calico features are supported:
   * L3 IP-IP Overlay Data Path.
   * BGP control plane.
   * Calico IPAM.
   * Management of Calico CNI plugin lifecycle.
   * Kubernetes Network Policy. This is experimental in 3.0.0.
* You can now use layer 7 routing in your Kubernetes workloads by using an
NGINX-based ingress controller.
[Learn more](https://docs.docker.com/ee/ucp/kubernetes/deploy-ingress-controller/).
* Layer 7 routing for Swarmkit applications has been upgraded to use
Interlock backend.
This adds increased performance, stability, and new features including SSL Termination,
Contextual Path-based Routing, Websocket Support, and Canary Application Instance
deployments. Existing Hostname Routing Mesh (HRM) labels (and newly added labels
with the old format) will automatically migrate to the new format. It is strongly
recommended to use the new format for new applications in order to take advantage
of the new features. [Learn more](https://docs.docker.com/ee/ucp/interlock/).

### Storage

* Support for NFS-based Kubernetes persistent volumes. Additional volume plugins
will be available in future releases.

### Security

* Role-based access control now supports Kubernetes resources.
[Lean more](https://docs.docker.com/ee/ucp/authorization/migrate-kubernetes-roles/).
  * In addition to users, teams, organizations, and grants you can now use
  Kubernetes Service Accounts as a subject type.
  [Learn more](https://docs.docker.com/ee/ucp/kubernetes/create-service-account/).
  * You can now create custom roles with Kubernetes API permissions. Default
  roles include Kubernetes API permissions based on their access type.
  As an example, View-Only contains Swarm and Kubernetes read-only API permissions.
  * In addition to collections, grants can now use Kubernetes Namespaces as
  a resource set type.
  * Admins can now link a Kubernetes namespace to a collection of nodes in order
  to isolate users and workloads between different nodes.
* Administrators can now enforce only running trusted images for both swarm and
  Kubernetes applications. [Learn more](https://docs.docker.com/ee/ucp/admin/configure/run-only-the-images-you-trust/)
* API support for registering multiple UCP clusters to a single DTR for the
  purposes of signed image enforcement. [Learn more](https://docs.docker.com/ee/ucp/admin/configure/integrate-with-multiple-registries/).
* The `Restricted Control` role includes the `User Impersonation` Kubernetes action,
  which can allow a user to escalate to admin privileges if the role is granted against
  `All Kubernetes Namespaces`. For this version, we recommend that administrators do not
  grant the `Restricted Control` role against Kubernetes namespaces, and use custom roles
  instead. This issue does not affect any other roles in the system, or any of the grants
  using `Restricted Control` against collections.
* For increased security UCP now requires clients to use TLS version 1.2.

### Known issues

* Platform support
  * Kubernetes is not yet supported for Windows based workloads. Use Swarmkit for
  Windows based workloads instead.
  * EE 2.0 is not yet supported in IBM Z platforms.
* Upgrade
  * After upgrading to 3.0, UCP uses the `191.168.0.0/16` CIDR to allocate IPs
  for Kubernetes Pods. This option is customizable when doing a fresh installation,
  but not during an upgrade. As a workaround to configure this option, before
  upgrading, [update the UCP configuration](admin/configure/ucp-configuration-file.md)
  to include `pod_cidr = "<ip>/<mask>"` under the `[cluster_config]` option.
* CLI
  * Both Docker and `kubectl` CLIs report that UCP is running Kubernetes 1.8.2,
  when in fact it is running 1.8.9.
* Networking
  * Swarm encrypted overlay networks might not work as expected because default
  Kubernetes firewall rules are interfering with them. [Learn more](https://success.docker.com/article/KB000717).
  * Calico networking for Kubernetes is not supported for Microsoft Azure. UCP
  leverages Azure networking and IPAM for control-plane and connectivity.
  [Learn how to deploy EE 2.0 on Azure](https://docs.docker.com/ee/ucp/admin/install/install-on-azure/).
  * Azure IPAM will fail if nodes in the cluster are connected to different subnets.
  As a workaround ensure network setup avoids multiple subnets. This will be
  rectified in an upcoming patch release.
  * UCP Calico control-plane supports full-mesh BGP peering only at release-time.
  Calico control-plane may cause high CPU on nodes in clusters above 100 nodes.
  A route reflector based partial-mesh BGP control-plane will reduce CPU
  consumption when scaling past 100 nodes.
  Route-reflector configurations will be included in a future release.
  * In some deployments the `kube-dns` component won't be able to resolve external
  domain names. [Deploy a ConfigMap to work around this](https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/#configure-stub-domain-and-upstream-dns-servers).
  * If you upgrade from UCP 2.x to UCP 3.x on Azure, Kubernetes networking doesn't work.
  The cluster upgrade completes, and Swarm workloads work, but Kubernetes networking
  will be down.
* Management
  * If upgrading UCP through the UI, UCP will not check to ensure the manager node
  has the minimum memory required of 4 GB. Upgrading through the CLI does check for
  this requirement.
  * Putting a node in `drain` mode currently removes only Swarm workloads, and not
  Kubernetes workloads. This will be fixed in a future release.
  * Kubernetes base image layer uses Ubuntu 16.04 which contains some known CVE
  vulnerabilities. These will be removed when the base image layer is updated.
  * Running `docker system prune -a` directly on individual worker nodes in the cluster
  will potentially delete UCP system images. This behavior will not occur if the
  prune command is run using a UCP client bundle.
  * Compose for Kubernetes only supports v3 or higher YAML files. Any older
  version YAML files will silently fail without errors.
  * Linking Kubernetes namespace to a collection of nodes in order to isolate
  Kubernetes workloads between different nodes is not working as expected.
  [You can use this workaround](https://success.docker.com/article/workaround-for-link-nodes-in-collection-ucp-300).
  * Running `kubectl get cs` might show some internal UCP components as
  unhealthy when that's not the case.
* Storage
  * UCP does not yet support dynamic volume provisioning (NFS volumes do not
    support this). This will change in future releases when more volume types
    are available.

### Deprecation notice

The following functionality has been deprecated with UCP 3.0.0 and will be
unavailable in the next UCP feature release.

* The web UI is going to stop supporting users to deploy stacks with basic
containers. You should update your Compose files to version 3, and deploy your
stack as a Swarm service or Kubernetes workload.
* The option to integrate with a remote Syslog system is going to be removed
from the UCP web UI. You can configure Docker Engine for this.
* The option to configure a rescheduling policy for basic containers is
deprecated. Deploy your applications as Swarm services or Kubernetes workloads.

# Version 2.2

## Version 2.2.12 (2018-08-09)

**Bug fixes**

* Security
  * Fixed a critical security issue where the LDAP bind username and password 
    were stored in clear text on UCP hosts. Please refer to the following KB article 
    https://success.docker.com/article/upgrading-to-ucp-2-2-12-ucp-3-0-4/ 
    for proper implementation of this fix.

## Version 2.2.11 (2018-07-26)

**New platforms**
* UCP 2.2.11 is supported running on RHEL 7.5 and Ubuntu 18.04.

**Bug fixes**

* Security
  * Fixed an issue that causes some security headers to not be added to all API responses.

* Core
  * Optimized swarm service read API calls through UCP.
  * Upgraded `RethinkDB` image to address potential security vulnerabilities.
  * Fixee an issue where removing a worker node from the cluster would cause an etcd member to be removed on a manager node.
  * Upgraded `etcd` version to 2.3.8.
  * Fixed an issue that causes classic Swarm to provide outdated data.
  * Fixed an issue that raises `ucp-kv` collection error with un-named volumes.

* UI
  * Fixed an issue that causes UI to not parse volume options correctly.
  * Fixed an issue that prevents the user from deploying stacks via UI.

## Version 2.2.10 (2018-05-17)

**Bug fixes**

* Security
  * Security headers are added for PCI compliance to all API responses.

* UI
  * Users can now set log driver name and options on both create and update
   service screens.
  * Fixed an issue that causes legacy collections on services to break UI. Now
   legacy collections are properly prepended with "/Shared/Legacy/".
  * Fixed an issue that causes service counts in status summary to be shown
   incorrectly.

* Authentication/Authorization
  * Private collections are now only created for those users which have
   previously logged in.
  * The logic which reconciles collection labels is now skipped if the
   node already has an access label.
  * Fixed an issue where LDAP syncs would always search against the last server
   in the list of additional domains if the search base DN matched any of those
   domains.

* Core
  * UCP can now be displayed in an iframe for pages hosted on the same domain.
  * Fixed an issue that prevents non-admin users to do `docker build` on UCP.
  * Fixed an issue where a node's status may be reported incorrectly in node
   listings.
  * UCP can now be installed on a system with more than 127 logical CPU cores.
  * Improved the performance of UCP's local and global health checks.

**Known Issue**

* Excessive delay is seen when sending `docker service ls` via UCP client
 bundle on a cluster that is running thousands of services.


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

## Release notes for earlier versions

- [UCP 2.1 release notes](/datacenter/ucp/2.1/guides/release-notes/index.md)
- [UCP 2.0 release notes](/datacenter/ucp/2.0/guides/release-notes.md)
- [UCP 1.1 release notes](/datacenter/ucp/1.1/release_notes.md)
