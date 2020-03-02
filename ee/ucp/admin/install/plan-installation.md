---
title: Plan your installation
description: Learn about the Docker Universal Control Plane architecture, and the requirements to install it on production.
keywords: UCP, install, Docker EE
---

>{% include enterprise_label_shortform.md %}

Docker Universal Control Plane helps you manage your container cluster from a
centralized place. This article explains what you need to consider before
deploying Docker Universal Control Plane for production.

## System requirements

Before installing UCP, make sure that all nodes (physical or virtual
machines) that you'll manage with UCP:

* [Comply with the system requirements](system-requirements.md), and
* Are running the same version of Docker Engine.

## Hostname strategy

Docker UCP requires Docker Enterprise. Before installing Docker Enterprise on
your cluster nodes, you should plan for a common hostname strategy.

Decide if you want to use short hostnames, like `engine01`, or Fully Qualified
Domain Names (FQDN), like `node01.company.example.com`. Whichever you choose,
confirm your naming strategy is consistent across the cluster, because
Docker Engine and UCP use hostnames.

For example, if your cluster has three hosts, you can name them:

```none
node1.company.example.com
node2.company.example.com
node3.company.example.com
```

## Static IP addresses

Docker UCP requires each node on the cluster to have a static IPv4 address.
Before installing UCP, ensure your network and nodes are configured to support
this.

## Avoid IP range conflicts

The following table lists recommendations to avoid IP range conflicts.

| Component  | Subnet                     | Range                                        | Default IP address                            |
| ---------- | -------------------------- | -------------------------------------------- | --------------------------------------------- |
| Engine     | `default-address-pools`    | CIDR range for interface and bridge networks | 172.17.0.0/16 - 172.30.0.0/16, 192.168.0.0/16 |
| Swarm      | `default-addr-pool`        | CIDR range for Swarm overlay networks        | 10.0.0.0/8                                    |
| Kubernetes | `pod-cidr`                 | CIDR range for Kubernetes pods               | 192.168.0.0/16                                |
| Kubernetes | `service-cluster-ip-range` | CIDR range for Kubernetes services           | 10.96.0.0/16                                  |

### Engine

There are two IP ranges used by the engine for the `docker0` and `docker_gwbridge` interface:

#### default-address-pools

`default-address-pools` defines a pool of CIDR ranges that are used to allocated subnets for local bridge networks. By default the first available subnet(`172.17.0.0/16`) is assigned to `docker0` and the next available subnet(`172.18.0.0/16`) is assigned to `docker_gwbridge`. Both the `docker0` and `docker_gwbridge` subnet can be modified by changing the `default-address-pools` value or as described in their individual sections below.

The default value for `default-address-pools` is:

 ```json
 {
     "default-address-pools": [
          {"base":"172.17.0.0/16","size":16}, <-- docker0
          {"base":"172.18.0.0/16","size":16}, <-- docker_gwbridge
          {"base":"172.19.0.0/16","size":16},
          {"base":"172.20.0.0/16","size":16},
          {"base":"172.21.0.0/16","size":16},
          {"base":"172.22.0.0/16","size":16},
          {"base":"172.23.0.0/16","size":16},
          {"base":"172.24.0.0/16","size":16},
          {"base":"172.25.0.0/16","size":16},
          {"base":"172.26.0.0/16","size":16},
          {"base":"172.27.0.0/16","size":16},
          {"base":"172.28.0.0/16","size":16},
          {"base":"172.29.0.0/16","size":16},
          {"base":"172.30.0.0/16","size":16},
          {"base":"192.168.0.0/16","size":20}
     ]
 }
 ```

`default-address-pools`:  A list of IP address pools for local bridge networks. Each entry in the list contain the following:

`base`: CIDR range to be allocated for bridge networks.

`size`: CIDR netmask that determines the subnet size to allocate from the `base` pool

As an example, `{"base":"192.168.0.0/16","size":20}` will allocate `/20` subnets from `192.168.0.0/16` yielding the following subnets for bridge networks:\
`192.168.0.0/20` (`192.168.0.0` - `192.168.15.255`)\
`192.168.16.0/20` (`192.168.16.0` - `192.168.31.255`)\
`192.168.32.0/20` (`192.168.32.0` - `192.168.47.255`)\
`192.168.48.0/20` (`192.168.32.0` - `192.168.63.255`)\
`192.168.64.0/20` (`192.168.64.0` - `192.168.79.255`)\
...\
`192.168.240.0/20` (`192.168.240.0` - `192.168.255.255`)

> Note
> 
> If the `size` matches the netmask of the `base`, then that pool only containers one subnet.
> 
> For example, `{"base":"172.17.0.0/16","size":16}` will only yield one subnet `172.17.0.0/16` (`172.17.0.0` - `172.17.255.255`).

#### docker0

By default, the Docker engine creates and configures the host system with a virtual network interface called `docker0`, which is an ethernet bridge device. If you don't specify a different network when starting a container, the container is connected to the bridge and all traffic coming from and going to the container flows over the bridge to the Docker engine, which handles routing on behalf of the container.

Docker engine creates `docker0` with a configurable IP range. Containers which are connected to the default bridge are allocated IP addresses within this range. Certain default settings apply to `docker0` unless you specify otherwise. The default subnet for `docker0` is the first pool in `default-address-pools` which is `172.17.0.0/16`.

The recommended way to configure the `docker0` settings is to use the `daemon.json` file.

If only the subnet needs to be customized, it can be changed by modifying the first pool of `default-address-pools`  in the `daemon.json` file.

```json
 {
     "default-address-pools": [
          {"base":"172.17.0.0/16","size":16}, <-- Modify this value
          {"base":"172.18.0.0/16","size":16},
          {"base":"172.19.0.0/16","size":16},
          {"base":"172.20.0.0/16","size":16},
          {"base":"172.21.0.0/16","size":16},
          {"base":"172.22.0.0/16","size":16},
          {"base":"172.23.0.0/16","size":16},
          {"base":"172.24.0.0/16","size":16},
          {"base":"172.25.0.0/16","size":16},
          {"base":"172.26.0.0/16","size":16},
          {"base":"172.27.0.0/16","size":16},
          {"base":"172.28.0.0/16","size":16},
          {"base":"172.29.0.0/16","size":16},
          {"base":"172.30.0.0/16","size":16},
          {"base":"192.168.0.0/16","size":20}
     ]
 }
```

> Note
>
> Modifying this value can also affect the `docker_gwbridge` if the `size` doesn't match the netmask of the `base`.

To configure a CIDR range and not rely on `default-address-pools`, the `fixed-cidr` setting can used:

```json
{
  "fixed-cidr": "172.17.0.0/16",
}
```

`fixed-cidr`: Specify the subnet for `docker0`, using standard CIDR notation. Default is `172.17.0.0/16`, the network gateway will be `172.17.0.1` and IPs for your containers will be allocated from (`172.17.0.2` - `172.17.255.254`).

To configure a gateway IP and CIDR range while not relying on `default-address-pools`, the `bip` setting can used:

```json
{
  "bip": "172.17.0.1/16",
}
```

`bip`: Specific a gateway IP address and CIDR netmask of the `docker0` network. The notation is `<gateway IP>/<CIDR netmask>` and the default is `172.17.0.1/16` which will make the `docker0` network gateway `172.17.0.1` and subnet `172.17.0.0/16`.

#### docker_gwbridge

The `docker_gwbridge` is a virtual network interface that connects the overlay networks (including the `ingress` network) to an individual Docker engine's physical network. Docker creates it automatically when you initialize a swarm or join a Docker host to a swarm, but it is not a Docker device. It exists in the kernel of the Docker host. The default subnet for `docker_gwbridge` is the next available subnet in `default-address-pools` which with defaults is `172.18.0.0/16`.

> Note
>
> If you need to customize the `docker_gwbridge` settings, you must do so before joining the host to the swarm, or after temporarily removing the host from the swarm.

The recommended way to configure the `docker_gwbridge` settings is to use the `daemon.json` file.

For `docker_gwbridge`, the second available subnet will be allocated from `default-address-pools`. If any customizations where made to the `docker0` interface it could affect which subnet is allocated. With the default `default-address-pools` settings you would modify the second pool.

```json
 {
     "default-address-pools": [
          {"base":"172.17.0.0/16","size":16},
          {"base":"172.18.0.0/16","size":16}, <-- Modify this value
          {"base":"172.19.0.0/16","size":16},
          {"base":"172.20.0.0/16","size":16},
          {"base":"172.21.0.0/16","size":16},
          {"base":"172.22.0.0/16","size":16},
          {"base":"172.23.0.0/16","size":16},
          {"base":"172.24.0.0/16","size":16},
          {"base":"172.25.0.0/16","size":16},
          {"base":"172.26.0.0/16","size":16},
          {"base":"172.27.0.0/16","size":16},
          {"base":"172.28.0.0/16","size":16},
          {"base":"172.29.0.0/16","size":16},
          {"base":"172.30.0.0/16","size":16},
          {"base":"192.168.0.0/16","size":20}
     ]
 }
```

### Swarm

Swarm uses a default address pool of `10.0.0.0/8` for its overlay networks. If this conflicts with your current network implementation, please use a custom IP address pool. To specify a custom IP address pool, use the `--default-addr-pool` command line option during [Swarm initialization](../../../../engine/swarm/swarm-mode.md).

> Note
>
> The Swarm `default-addr-pool` setting is separate from the Docker engine `default-address-pools` setting. They are two separate ranges that are used for different purposes.

> Note
>
> Currently, the UCP installation process does not support this flag. To deploy with a custom IP pool, Swarm must first be initialized using this flag and UCP must be installed on top of it.

### Kubernetes

There are two internal IP ranges used within Kubernetes that may overlap and
conflict with the underlying infrastructure:

* The Pod Network -  Each Pod in Kubernetes is given an IP address from either
  the Calico or Azure IPAM services. In a default installation Pods are given
  IP addresses on the `192.168.0.0/16` range. This can be customized at install time by passing the `--pod-cidr` flag to the 
  [UCP install command](/reference/ucp/{{ site.ucp_version }}/cli/install/). 
* The Services Network - When a user exposes a Service in Kubernetes it is
  accessible via a VIP, this VIP comes from a Cluster IP Range. By default on UCP
  this range is `10.96.0.0/16`. Beginning with 3.1.8, this value can be
  changed at install time with the `--service-cluster-ip-range` flag.

## Avoid firewall conflicts

For SUSE Linux Enterprise Server 12 SP2 (SLES12), the `FW_LO_NOTRACK` flag is turned on by default in the openSUSE firewall. This speeds up packet processing on the loopback interface, and breaks certain firewall setups that need to redirect outgoing packets via custom rules on the local machine.

To turn off the FW_LO_NOTRACK option, edit the `/etc/sysconfig/SuSEfirewall2` file and set `FW_LO_NOTRACK="no"`. Save the file and restart the firewall or reboot.

For SUSE Linux Enterprise Server 12 SP3, the default value for `FW_LO_NOTRACK` was changed to `no`.

For Red Hat Enterprise Linux (RHEL) 8, if firewalld is running and `FirewallBackend=nftables` is set in `/etc/firewalld/firewalld.conf`, change this to `FirewallBackend=iptables`, or you can explicitly run the following commands to allow traffic to enter the default bridge (docker0) network:

```
firewall-cmd --permanent --zone=trusted --add-interface=docker0
firewall-cmd --reload
```
## Time synchronization

In distributed systems like Docker UCP, time synchronization is critical
to ensure proper operation. As a best practice to ensure consistency between
the engines in a UCP cluster, all engines should regularly synchronize time
with a Network Time Protocol (NTP) server. If a host node's clock is skewed,
unexpected behavior may cause poor performance or even failures.

## Load balancing strategy

Docker UCP doesn't include a load balancer. You can configure your own
load balancer to balance user requests across all manager nodes.

If you plan to use a load balancer, you need to decide whether you'll
add the nodes to the load balancer using their IP addresses or their FQDNs.
Whichever you choose, be consistent across nodes. When this is decided,
take note of all IPs or FQDNs before starting the installation.

[Learn how to set up your load balancer](../configure/join-nodes/use-a-load-balancer.md).

## Load balancing UCP and DTR

By default, UCP and DTR both use port 443. If you plan on deploying UCP and
DTR, your load balancer needs to distinguish traffic between the two by IP
address or port number.

* If you want to configure your load balancer to listen on port 443:
  * Use one load balancer for UCP and another for DTR.
  * Use the same load balancer with multiple virtual IPs.
* Configure your load balancer to expose UCP or DTR on a port other than 443.

If you want to install UCP in a high-availability configuration that uses
a load balancer in front of your UCP controllers, include the appropriate IP
address and FQDN of the load balancer's VIP by using
one or more `--san` flags in the 
[UCP install command](/reference/ucp/{{ site.ucp_version }}/cli/install/)
or when you're asked for additional SANs in interactive mode.
[Learn about high availability](../configure/set-up-high-availability.md).

## Use an external Certificate Authority

You can customize UCP to use certificates signed by an external Certificate
Authority. When using your own certificates, you need to have a certificate
bundle that has:

* A ca.pem file with the root CA public certificate,
* A cert.pem file with the server certificate and any intermediate CA public
certificates. This certificate should also have SANs for all addresses used to
reach the UCP manager,
* A key.pem file with server private key.

You can have a certificate for each manager, with a common SAN. For
example, on a three-node cluster, you can have:

* node1.company.example.org with SAN ucp.company.org
* node2.company.example.org with SAN ucp.company.org
* node3.company.example.org with SAN ucp.company.org

You can also install UCP with a single externally-signed certificate
for all managers, rather than one for each manager node. In this case,
the certificate files are copied automatically to any new
manager nodes joining the cluster or being promoted to a manager role.

## Where to go next

* [System requirements](system-requirements.md)
* [Install UCP](index.md)

