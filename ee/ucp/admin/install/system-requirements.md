---
title: UCP System requirements
description: Learn about the system requirements for installing Docker Universal Control Plane.
keywords: UCP, architecture, requirements, Docker Engine - Enterprise
redirect_from:
- /enterprise/admin/install/system-requirements/
---

>{% include enterprise_label_shortform.md %}

Docker Universal Control Plane can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and software requirements

You can install UCP on-premises or on a cloud provider. Common requirements:

* [Docker Engine - Enterprise](/ee/supported-platforms.md) version {{ site.docker_ee_version }}
* Linux kernel version 3.10 or higher. For debugging purposes, it is suggested to match the host OS kernel versions as close as possible.
* [A static IP address for each node in the cluster](/ee/ucp/admin/install/plan-installation/#static-ip-addresses)
* User namespaces should not be configured on any node. This function is not currently supported by UCP. See [Isolate containers with a user namespace](https://docs.docker.com/engine/security/userns-remap/) for more information.

### Minimum requirements

* 8GB of RAM for manager nodes
* 4GB of RAM for worker nodes
* 2 vCPUs for manager nodes
* 10GB of free disk space for the `/var` partition for manager nodes (A minimum of 6GB is recommended.)
* 500MB of free disk space for the `/var` partition for worker nodes

* Default install directories:
   - /var/lib/docker (Docker Data Root Directory)
   - /var/lib/kubelet (Kubelet Data Root Directory)
   - /var/lib/containerd (Containerd Data Root Directory)

> Note
>
> Increased storage is required for Kubernetes manager nodes in UCP 3.1. If you are upgrading to UCP 3.1, refer to [Kubelet restarting after upgrade to Universal Control Plane 3.1](https://success.docker.com/article/kublet-restarting-after-upgrade-to-universal-control-plane-31) for information on how to increase the size of the `/var/lib/kubelet` filesystem.

### Recommended production requirements

 * 16GB of RAM for manager nodes
 * 4 vCPUs for manager nodes
 * 25-100GB of free disk space

Note that Windows container images are typically larger than Linux container images. For
this reason, you should provision more local storage for Windows
nodes and for any DTR setups that store Windows container images.

Also, make sure the nodes are running an [operating system supported by Docker Enterprise](https://success.docker.com/Policies/Compatibility_Matrix).

For highly-available installations, you also need a way to transfer files
between hosts.

> Workloads on manager nodes
>
> Docker does not support workloads other than those required for UCP on UCP manager nodes.

## Ports used

When installing UCP on a host, a series of ports need to be opened to incoming
traffic. Each of these ports will expect incoming traffic from a set of hosts,
indicated as the "Scope" of that port. The three scopes are:
- External: Traffic arrives from outside the cluster through end-user
  interaction.
- Internal: Traffic arrives from other hosts in the same cluster.
- Self: Traffic arrives to that port only from processes on the same host.

> Note
>
> When installing UCP on Microsoft Azure, an overlay network is not used for
> Kubernetes; therefore, any containerized service deployed onto Kubernetes and
> exposed as a Kubernetes Service may need its corresponding port to be opened
> on the underlying Azure Network Security Group. For more information see
> [Installing on
> Azure](/ee/ucp/admin/install/cloudproviders/install-on-azure/#azure-prerequisites).

Make sure the following ports are open for incoming traffic on the respective
host types:

|       Hosts       |          Port           |       Scope        |                                    Purpose                                    |
| :---------------- | :---------------------- | :----------------- | :---------------------------------------------------------------------------- |
| managers, workers | TCP 179                 | Internal           | Port for BGP peers, used for Kubernetes networking                            |
| managers          | TCP 443  (configurable) | External, Internal | Port for the UCP web UI and API                                               |
| managers          | TCP 2376 (configurable) | Internal           | Port for the Docker Swarm manager. Used for backwards compatibility           |
| managers          | TCP 2377 (configurable) | Internal           | Port for control communication between swarm nodes                            |
| managers, workers | UDP 4789                | Internal           | Port for overlay networking                                                   |
| managers          | TCP 6443 (configurable) | External, Internal | Port for Kubernetes API server endpoint                                       |
| managers, workers | TCP 6444                | Self               | Port for Kubernetes API reverse proxy                                         |
| managers, workers | TCP, UDP 7946           | Internal           | Port for gossip-based clustering                                              |
| managers, workers | TCP 9099                | Self               | Port for calico health check
| managers, workers | TCP 10250               | Internal           | Port for Kubelet                                                              |
| managers, workers | TCP 12376               | Internal           | Port for a TLS authentication proxy that provides access to the Docker Engine |
| managers, workers | TCP 12378               | Self               | Port for Etcd reverse proxy                                                   |
| managers          | TCP 12379               | Internal           | Port for Etcd Control API                                                     |
| managers          | TCP 12380               | Internal           | Port for Etcd Peer API                                                        |
| managers          | TCP 12381               | Internal           | Port for the UCP cluster certificate authority                                |
| managers          | TCP 12382               | Internal           | Port for the UCP client certificate authority                                 |
| managers          | TCP 12383               | Internal           | Port for the authentication storage backend                                   |
| managers          | TCP 12384               | Internal           | Port for the authentication storage backend for replication across managers   |
| managers          | TCP 12385               | Internal           | Port for the authentication service API                                       |
| managers          | TCP 12386               | Internal           | Port for the authentication worker                                            |
| managers          | TCP 12388               | Internal           | Internal Port for the Kubernetes API Server                                   |

## Disable `CLOUD_NETCONFIG_MANAGE` for SLES 15
For SUSE Linux Enterprise Server 15 (SLES 15) installations, you must disable `CLOUD_NETCONFIG_MANAGE`
prior to installing UCP.

    1. In the network interface configuration file, `/etc/sysconfig/network/ifcfg-eth0`, set
    ```
    CLOUD_NETCONFIG_MANAGE="no"
    ```
    2. Run `service network restart`.

## Enable ESP traffic

For overlay networks with encryption to work, you need to ensure that
IP protocol 50 (Encapsulating Security Payload) traffic is allowed.

## Enable IP-in-IP traffic

The default networking plugin for UCP is Calico, which uses IP Protocol
Number 4 for IP-in-IP encapsulation.

If you're deploying to AWS or another cloud provider, enable IP-in-IP
traffic for your cloud provider's security group.

## Enable connection tracking on the loopback interface for SLES
Calico's Kubernetes controllers can't reach the Kubernetes API server
unless connection tracking is enabled on the loopback interface. SLES
disables connection tracking by default.

On each node in the cluster:

```
sudo mkdir -p /etc/sysconfig/SuSEfirewall2.d/defaults
echo FW_LO_NOTRACK=no | sudo tee /etc/sysconfig/SuSEfirewall2.d/defaults/99-docker.cfg
sudo SuSEfirewall2 start
```

## Timeout settings

Make sure the networks you're using allow the UCP components enough time
to communicate before they time out.

| Component                              | Timeout (ms) | Configurable |
|:---------------------------------------|:-------------|:-------------|
| Raft consensus between manager nodes   | 3000         | no           |
| Gossip protocol for overlay networking | 5000         | no           |
| etcd                                   | 500          | yes          |
| RethinkDB                              | 10000        | no           |
| Stand-alone cluster                    | 90000        | no           |

## Time Synchronization

In distributed systems like Docker UCP, time synchronization is critical
to ensure proper operation. As a best practice to ensure consistency between
the engines in a UCP cluster, all engines should regularly synchronize time
with a Network Time Protocol (NTP) server. If a server's clock is skewed,
unexpected behavior may cause poor performance or even failures.

## Compatibility and maintenance lifecycle

Docker Enterprise is a software subscription that includes three products:

* Docker Engine - Enterprise with enterprise-grade support
* Docker Trusted Registry
* Docker Universal Control Plane

Learn more about compatibility and the maintenance lifecycle for these products:

- [Compatibility Matrix](https://success.docker.com/Policies/Compatibility_Matrix)
- [Maintenance Lifecycle](https://success.docker.com/Policies/Maintenance_Lifecycle)

## Where to go next

- [Plan your installation](plan-installation.md)
- [UCP architecture](../../ucp-architecture.md)