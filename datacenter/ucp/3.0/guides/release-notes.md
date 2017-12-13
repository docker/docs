---
title: UCP 3.0 Beta1 release notes
description: Release notes for Docker Universal Control Plane. Learn more about the
  changes introduced in the latest versions.
keywords: UCP, release notes
---

Here you can learn about new features, bug fixes, breaking changes, and
known issues for the latest UCP version.
You can then use [the upgrade instructions](admin/install/upgrade.md), to
upgrade your installation to the latest release.

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
* The kubernetes SPDY operations such as `kubectl logs` or `kubectl exec` are
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