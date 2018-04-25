---
title: UCP 3.0 release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: UCP, release notes
redirect_from:
  - /datacenter/ucp/3.0/guides/release-notes/
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for the latest UCP version.
You can then use [the upgrade instructions](admin/install/upgrade.md), to
upgrade your installation to the latest release.

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
  rectified in an upcoming patch release (#12894).
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

## Release notes for earlier versions

- [UCP 2.2 release notes](/datacenter/ucp/2.2/guides/release-notes.md)
- [UCP 2.1 release notes](/datacenter/ucp/2.1/guides/release-notes/index.md)
- [UCP 2.0 release notes](/datacenter/ucp/2.0/guides/release-notes.md)
- [UCP 1.1 release notes](/datacenter/ucp/1.1/release_notes.md)
