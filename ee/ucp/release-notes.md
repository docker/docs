---
title: UCP 3.0 Beta release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: UCP, release notes
ui_tabs:
- version: ucp-3.0
  orhigher: false
- version: ucp-2.2
  orlower: true
redirect_from:
  - /datacenter/ucp/3.0/guides/release-notes/
---
{% if include.version=="ucp-3.0" %}

Here you can learn about new features, bug fixes, breaking changes, and
known issues for the latest UCP version.
You can then use [the upgrade instructions](admin/install/upgrade.md), to
upgrade your installation to the latest release.

## Version 3.0.0-beta2

(18 January 2018)

**New features**

* This is the second beta version of Docker Universal Control Plane version
  3.0.0 (and the first publicly available one).
* Interlock-2/HRM now defaults to 2 replicas for improved durability.

**Bug fixes**

* Core
  * Use default service accounts for admin pods. This means that more Kubernetes
  apps will run without modifications.
  * `ucp-calico-node` now works correctly on hosts with local bridge networks.
  Previously, setup would fail.
  * Creating support dumps no longer fails when blocking containers are present.
  * Stopping and starting the Docker engine on a UCP master no longer fails due
  to mount propagation issues.
  * In case errors are encountered during backup and restore, the errors are
  now printed. Previously, no output was provided in case of backup or restore
  errors.
  * `kubectl exec`, `kubectl logs --follow` and other interactive `kubectl`
  commands now work correctly. Proxy limitations previously blocked these CLI actions.
  * Fixed quoting bug that prevented creation of Kubernetes secrets.
  * `docker images --filter` now works correctly when using the UCP API endpoint.
  * Kubernetes workloads are now evicted from nodes in case the node-mode is
  changed to not run Kubernetes workloads.
  * Fixed a problem that would cause bad Kubernetes configurations for stacks
  deployed using `docker-compose.yml` files.
  * Kubernetes networking components now install correctly on hosts with SELinux
  in "Enforcing" mode.
  * Fixed bug that caused nodes to revert to revert to `/Shared` collection after
  having been assigned to a different collection.
  * Interlock default backend is now correctly configured.
  * Fixed a bug that could cause cluster failure due to missing labels and taints
  when promoting workers to managers.
  * Fixed reconciliation bug caused by problems mounting `/var` on RHEL and CentOS.
  * The Docker Kubernetes flex-volume bridge doesn't work with volumes created
  with volume plugins. Only local volumes (and v1 plugin volumes) work.
  * Fixed Calico-related problem that caused Kubernetes HostPort publishing to
  not work.

* UI/UX
  * Fixed bug where container list view would load forever when there where no
  containers to show.
  * Clicking Kubernetes published port URLs now opens new browser window or tab.
  * Kubernetes published port URL now work. Previously, these urls where not
  reliable due to problems with `kube-proxy`.
  * Swarm service status indicator now works. A missing API call caused the the
  indicator to not work after service creation.
  * Fixed quoting issue in client bundle instructions.
  * Kubernetes resources with identical names in different namespaces are now
  visible in the UI. A bug previously meant that only one resource would be shown.
  * Fixed a problem in the Kubernetes parts of the UI where hidden system
  services would count towards the total count displayed for the list.

**Known Issues**

* Deploying Kubernetes Helm charts through Docker Universal Control plane does
not work due to compatibility issues with default service accounts.
* The product versions reported by components in the beta release are preliminary
and will change for the GA final release.
* Security hardening of Kubernetes managed by Docker EE is not fully complete.
Only use this release for testing and validation in controlled environments with
trusted users.
* Beta2 has been tested on RHEL 7.3, 7.4, and Ubuntu 16.04. There are confirmed
incompatibilities with SLES 12 and Ubuntu 14.04.
* Rotating UCP certificates can cause Kubernetes networking to stop working.
* Currently, the only tested and supported Kubernetes CNI networking plugin is
Calico (which is included).
* The Kubernetes overlay networking implementation (based on Calico) will not
currently default to relevant underlay networking on infrastructure where that
is supported.
* Certificate used for kubelet API do not include node IP in the SAN
(Subject Alternate Name) list. This may cause extensions and other software that
interacts with Kubernetes to fail.
* Deploying Docker Compose files as Kubernetes objects is currently only
supported with v3 Compose files and later. A bug means that pre-v3 files are not
rejected and deployment silently fails.
* In-product interactive Kubernetes docs are currently broken.

## Version 3.0 Beta1

(11 December 2017)

**New features**

* UCP now supports Kubernetes as an orchestrator, in addition to the existing
Swarmkit and "Classic" Swarm orchestrators. Kubernetes system components are
installed automatically on all manager and worker nodes in the cluster.
Kubernetes in UCP is not yet supported on Windows or IBM Z nodes.
* Worker nodes can be allocated to run only Swarm workloads, only Kubernetes
workloads, or mixed workloads. Mixed workloads are not recommended for use in
a production environment due to potential resource contention issues across
orchestrators. Manager nodes are by default Mixed in order to support Swarm
and Kubernetes system components.
* Hostname Routing Mesh (HRM) has been upgraded to use Interlock 2 backend for
layer 7 routing. This adds increased performance, stability, and new features
including SSL Termination, Contextual Path-based Routing, Websocket Support,
and Canary Application Instance deployments. Existing HRM labels (and newly
added labels with the old format) will automatically migrate to the new format.
It is recommended to use the new format for new applications in order to take
advantage of the added features.

**Known issues**

* UCP 3.0 Beta1 has been tested on Ubuntu 16.04 and RHEL 7.3.
* Installation on SLES 12 and Ubuntu 14.04 is not currently possible because of
an iptables incompatibility.
* UCP 3.0 requires more resources to run than UCP 2.2 and is unlikely to work
correctly on nodes with less than 4GB of total memory.
* The Kubernetes SPDY operations such as `kubectl logs` or `kubectl exec` are
not possible when using the client bundle feature. As a workaround, you may
change all references from `:443` to `:6443` in the `kube.yml` and `env.sh`
files of a user client bundle.
* The default service account of each namespace currently has no permissions,
while all other service accounts have admin-level permissions and are usable
only by admin users.  Admins should create custom service accounts for workloads
intended to use the service account feature.
* Security hardening of Kubernetes managed by Docker EE is not fully complete.
Only use this release for testing and validation in controlled environments with
trusted users. Apps that rely on the default Kubernetes service account may not
work because of access restrictions.
* The product versions reported by components in the beta release are
preliminary and will change for the GA final release.
* Interlock, the successor to the HTTP Routing Mesh (HRM), has not yet been
fully scale-tested or optimized and the default settings in the beta release
are not optimized for production use.
* HRM supports configuring a catch-all fallback service for requests that do
not match any routing directive. Interlock does not support this feature.
* When changing the mode for a worker node between Kubernetes, Swarm and Mixed,
when going back to "Swarm" some Kubernetes tasks may not be evicted from the
node.
* A Kubernetes bug causes workloads published using `NodePort` to only be
accessible on the particular nodes that are running pods for the workload.
* Problems with state-reconciliation may cause nodes that are changed from
managers to workers to get into a state where they're incorrectly running
Kubernetes master components.
* Installing on systems with SELinux in enforcing mode currently fails because
of a Calico installation problem.
* Deleting Kubernetes Pods may leave pods in "Terminating" state with no way
to delete them.
* Removing a node from Docker Swarm may not remove the node from the Kubernetes
node set.
* When promoting nodes from worker to master, not all required labels and taints
are correctly applied. This can cause cluster failure if master nodes are lost.
Reconciliation may fail on RHEL and CentOS because of problems with mounting the
`/var` folder.
* The Docker Kubernetes flex-volume bridge doesn't work with volumes created
with volume plugins. Only local volumes (and v1 plugin volumes) work.
* Rotating UCP certificates can cause Kubernetes networking to stop working.
* Currently, the only tested and supported Kubernetes CNI networking plugin is
Calico (which is included).
* Kubernetes resources with identical names in different namespaces are not
visible in the UI because of a problem with the middleware object store.
* If all the nodes in UCP are using IBM Z, Kubernetes functionality will show
up on the web UI but not be available for use.
* Deploying a Compose V2 or older file as a Kubernetes is not supported but
currently UCP doesn't present a meaningful error message to the user.

{% elsif include.version=="ucp-2.2" %}

- [UCP 2.2 release notes](/datacenter/ucp/2.2/guides/release-notes.md)
- [UCP 2.1 release notes](/datacenter/ucp/2.1/guides/release-notes/index.md)
- [UCP 2.0 release notes](/datacenter/ucp/2.0/guides/release-notes.md)
- [UCP 1.1 release notes](/datacenter/ucp/1.1/release_notes.md)

{% endif %}
