<!--[metadata]>
+++
title = "CS Engine release notes"
description = "Commercially supported Docker Engine release notes"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, Commercially Supported Docker Engine, release notes"]
[menu.main]
parent="smn_dhe"
weight=102
+++
<![end-metadata]-->

# Commercially supported Engine release notes

This document describes the latest changes, additions, known issues, and fixes
for the commercially supported Docker Engine (CS Engine).

The CS Engine is functionally equivalent to the corresponding Docker Engine that
it references. However, a commercially supported release also includes
back-ported fixes (security-related and priority defects) from the open source.
It incorporates defect fixes that you can use in environments where new features
cannot be adopted as quickly for consistency and compatibility reasons.

#### Prior versions

These notes refer to the current and immediately prior releases of the
CS Engine. For notes on older versions, see the [CS Engine prior release notes archive](cse-prior-release-notes.md).

## CS Engine 1.9.1-cs3
(6 January 2016)

This release addresses the following issues:

* The commercially supported Engine 1.9.1-cs3 now supports multi-host networking
for all the kernels that the base CS Engine is supported on.

>**Note**: Centos 7 has its firewall enabled by default and it prevents the VXLAN tunnel from communicating. If this applies to you, then after installing the CS Engine, execute the following command in the Linux host:

`sudo firewall-cmd --zone=public --permanent --add-port=4789/udp`


* Corrected an issue where Docker didn't remove the Masquerade NAT rule from `iptables` when the network was removed. This caused the gateway address to be
incorrectly propagated as the source address of a connection.

* Fixed an issue where if the daemon started multiple containers concurrently, then the `/etc/hosts` files were incompletely populated. This issue occurred randomly.

* Corrected an issue where the same IP address for different Docker containers resulted in network connection inconsistencies. Now each container has a separate IP address.

* Corrected an issue where the IPv6 gateway was not created when using custom networks although the network had a configured gateway.

* Fixed an issue where users might have experienced a panic error if the  daemon was started with the `—cluster-store` option, but without the `—cluster-advertise` option.

## CS Engine 1.9.1-cs2
(4 December 2015)

Starting with this release, upgrading minor versions, for example, from 1.9.0 to 1.9.1, is faster and easier. See the [upgrade](install/upgrade.md) documentation for details.

You can refer to the detailed list of all changes since the release of CS Engine
1.9.0
https://github.com/docker/docker/releases.

## CS Engine 1.9.0
(12 November 2015)

Highlighted feature summary:

* Network Management and Plugins. Networks are now first class objects that can be listed, created, deleted, inspected, and connected to or disconnected from a
container. They can be manipulated outside of the container themselves and are
fully manageable on its own lifecycle. You can also use plugins to extend
network functionality.

* Docker, Inc. now provides support for the in-box Overlay (for cross-host networking) and Bridge network plugins. You can find more information about how
to manage networks and using network plugins in the [documentation](https://docs.docker.com/engine/userguide/networking/dockernetworks/).

* Volume Management and Plugins. Volumes also become discrete, manageable objects in Docker. Volumes can be listed, created, deleted, and inspected.
Similar to networks, they have their own managed lifecycle outside of the
container. Plugins allow others to write and extend the functionality of volumes
or provide integration with other types of storage.

* The in-box volume driver is included and supported. You can find more information about how to manage volumes  and using  volume plugins in the
documentation.

* Docker Content Trust. Use Content Trust to both verify the integrity and the publisher of all the data received from a registry over any channel. Content Trust is currently only supported using Docker Hub notary servers.

* Updated the release cadence of the CS Docker Engine. Starting with this version, Docker supports **every** major release of Docker Engine from open
source with three releases under support at one time. This means you’ll be able
to take advantage of the latest and greatest features and you won’t have to wait
for a supported release to take advantage of a specific feature.

Refer to the [detailed list](https://github.com/docker/docker/releases) of all changes since the release of CS Engine 1.6.
