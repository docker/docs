<!--[metadata]>
+++
aliases = ["/docker-trusted-registry/cse-release-notes/"]
title = "CS Engine release notes"
description = "Commercially supported Docker Engine release notes"
keywords = ["docker, documentation, about, technology, understanding, enterprise, hub, registry, Commercially Supported Docker Engine, release notes"]
[menu.main]
parent="menu_csengine_release_notes"
weight=0
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
CS Engine. For notes on older versions, see the [CS Engine prior release notes archive](prior-release-notes.md).

## CS Engine 1.10.2-cs1
(22 February 2016)

In this release the CS Engine is supported on SUSE Linux Enterprise 12 OS.

Refer to the [detailed list](https://github.com/docker/docker/releases) of all changes since the release of CS Engine 1.9.1.

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

Starting with this release, upgrading minor versions, for example, from 1.9.0 to 1.9.1, is faster and easier.

You can refer to the detailed list of all changes since the release of CS Engine
1.9.0
https://github.com/docker/docker/releases.
