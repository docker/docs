---
description: Commercially supported Docker Engine release notes
keywords: docker, documentation, about, technology, understanding, enterprise, hub, registry, Commercially Supported Docker Engine, release notes
redirect_from:
- /docker-trusted-registry/cse-release-notes/
- /docker-trusted-registry/cs-engine/release-notes/release-notes/
- /cs-engine/release-notes/release-notes/
title: Commercially Supported Engine release notes
---

This document describes the latest changes, additions, known issues, and fixes
for the commercially supported Docker Engine (CS Engine).

The CS Engine is functionally equivalent to the corresponding Docker Engine that
it references. However, a commercially supported release also includes
back-ported fixes (security-related and priority defects) from the open source.
It incorporates defect fixes that you can use in environments where new features
cannot be adopted as quickly for consistency and compatibility reasons.

## Prior versions

These notes refer to the current and immediately prior releases of the
CS Engine. For notes on older versions, see the [CS Engine prior release notes archive](prior-release-notes.md).

## CS Engine 1.12.6-cs13
(28 Jul 2017)

* Fix packaging issue where packages were missing a `containerd` patch.
  This resolves an issue with a deadlock in containerd related to healtchecks.
* Fix a deadlock on cancelling healthcecks. [#28462](https://github.com/moby/moby/pull/28462)

## CS Engine 1.12.6-cs12
(01 Jun 2017)

* Fix an issue where if a volume using the local volume driver which has
mount options fails to unmount on container exit, the data in the mount may be
lost if the user attempts to subsequently remove the volume. [#32327](https://github.com/docker/docker/pulls/32327)

## CS Engine 1.12.6-cs11
(11 May 2017)

* Fix an issue with overlay networks L2 miss notifications not being handled in
some cases [#1642](https://github.com/docker/libnetwork/pull/1642)

## CS Engine 1.12.6-cs10
(6 Mar 2017)

* Fix concurrency issue in libnetwork

## CS Engine 1.12.6-cs9
(28 Feb 2017)

* Fixes an issue causing containerd to deadlock [#336](https://github.com/docker/containerd/pull/336)
* Fixes an issue where encrypted overlay networks stop working [#30727](https://github.com/docker/docker/issues/30727)

## CS Engine 1.12.6-cs8
(8 Feb 2017)

This release addresses the following issues:

* Addresses performance issues introduced by external KV-Store access with the
  `docker network ls` endpoint with large amounts of overlay networks and containers
  attached to those networks

* Addresses an inconsistent mac -> vtep binding issue when using overlay networks

* Adds a new repository for RHEL 7.2 users, to deal with issues
  users have encountered when installing the docker-engine-selinux package
  on systems pinned to 7.2 packages that are older than those available in the
  normal 7.2 install. This change relates to packaging changes for
  [1.12.6-cs7](#cs-engine-1126-cs7).

  Users experiencing issues installing the selinux package should switch to this
  repository. See [install instructions](/cs-engine/install.md) for more details.
  Only switch to this repository if you encounter problems installing the
  selinux packages from the centos/7 repo.

## CS Engine 1.12.6-cs7
(24 Jan 2017)

This release addresses the following issues:

* [#28406](https://github.com/docker/docker/issues/28406) Fix conflicts introduced
by the updated `selinux-policy` base package from RHEL/CentOS 7.3
* [#26639](https://github.com/docker/docker/issues/26639) Resolves hostnames passed
to the local volume driver for nfs mount options.
* [#26111](https://github.com/docker/docker/issues/26111) Fix issue with adding
iptables rules due to xtables lock message change.

## CS Engine 1.12.6-cs6
(10 Jan 2017)

Bumps RunC version to address CVE-2016-9962.

Refer to the [detailed list](https://github.com/docker/docker/releases/tag/v1.12.6) of all
changes since the release of CS Engine 1.12.5-cs5.

## CS Engine 1.12.5-cs5
(21 Dec 2016)

Refer to the [detailed list](https://github.com/docker/docker/releases/tag/v1.12.5) of all
changes since the release of CS Engine 1.12.3-cs4.

## CS Engine 1.12.3-cs4
(11 Nov 2016)

This releases addresses the following issues:

* [#27370](https://github.com/docker/docker/issues/27370) Fix `--net-alias` for
`--attachable` networks
* [#28051](https://github.com/docker/docker/issues/28051) Fix an issue removing
a `--attachable` network by ID.

## CS Engine 1.12.3-cs3
(27 Oct 2016)

Refer to the [detailed list](https://github.com/docker/docker/releases) of all
changes since the release of CS Engine 1.12.2-cs2.

## CS Engine 1.12.2-cs2
(13 Oct 2016)

Refer to the [detailed list](https://github.com/docker/docker/releases) of all
changes since the release of CS Engine 1.12.1-cs1.

## CS Engine 1.12.1-cs1
(20 Sep 2016)

Refer to the [detailed list](https://github.com/docker/docker/releases) of all
changes since the release of CS Engine 1.11.2-cs5.

This release addresses the following issues:

* [#25962](https://github.com/docker/docker/pull/25962) Allow normal containers
to connect to swarm-mode overlay network
* Various bug fixes in swarm mode networking

## CS Engine 1.11.2-cs8
(01 Jun 2017)

* Fix an issue where if a volume using the local volume driver which has
mount options fails to unmount on container exit, the data in the mount may be
lost if the user attempts to subsequently remove the volume. [#32327](https://github.com/docker/docker/pulls/32327)

## CS Engine 1.11.2-cs7
(24 Jan 2017)

This release addresses the following issues:

* [#26639](https://github.com/docker/docker/issues/26639) Resolves hostnames passed
to the local volume driver for nfs mount options.
* [#26111](https://github.com/docker/docker/issues/26111) Fix issue with adding
iptables rules due to xtables lock message change.
* [#1572](https://github.com/docker/libnetwork/issues/1572) Fix daemon panic
* [#1130](https://github.com/docker/libnetwork/pull/1130) Fix IPAM out of sync
issue on ungraceful shutdown.

## CS Engine 1.11.2-cs6
(12 Jan 2017)

Bumps RunC version to address CVE-2016-9962.

## CS Engine 1.11.2-cs5
(13 Sep 2016)

This release addresses the following issues:

* Make the docker daemon ignore the `SIGPIPE` signal
[#19728](https://github.com/docker/docker/issues/19728)
* Fix race in libdevicemapper symlink handling
[#24671](https://github.com/docker/docker/issues/24671)
* Generate additional logging when unmarshalling devicemapper metadata
[#23974](https://github.com/docker/docker/pull/23974)
* Drop queries in root domain when ndots is set
[#1441](https://github.com/docker/libnetwork/pull/1441)

## CS Engine 1.11.2-cs4
(16 Aug 2016)

This release addresses the following issues:

* Change systemd kill mode to `process` so systemd only stops the docker daemon
[#21933](https://github.com/docker/docker/issues/21933)
* Fix dropped external DNS responses when greater than 512 bytes
[#1373](https://github.com/docker/libnetwork/pull/1373)
* Remove UDP connection caching in embedded DNS server
[#1352](https://github.com/docker/libnetwork/pull/1352)
* Fix issue where truncated DNS replies were discarded by the embedded DNS server
[#1351](https://github.com/docker/libnetwork/pull/1351)

## CS Engine 1.11.2-cs3
(7 Jun 2016)

This release addresses the following issues:

* Fix potential panic when running `docker build`
[#23032](https://github.com/docker/docker/pull/23032)
* Fix interpretation of `--user` parameter
[#22998](https://github.com/docker/docker/pull/22998)
* Fix a bug preventing container statistics from being correctly reported
[#22955](https://github.com/docker/docker/pull/22955)
* Fix an issue preventing containers from being restarted after daemon restart
[#22947](https://github.com/docker/docker/pull/22947)
* Fix a possible deadlock on image deletion and container attach
[#22918](https://github.com/docker/docker/pull/22918)
* Fix an issue causing `docker ps` to hang when using devicemapper
[#22168](https://github.com/docker/docker/pull/22168)
* Fix a bug preventing to `docker exec` into a container when using
devicemapper [#22168](https://github.com/docker/docker/pull/22168)

## CS Engine 1.11.1-cs2
(17 May 2016)

This release fixes the following issue which prevented DTR containers to be automatically restarted on a docker daemon restart:

https://github.com/docker/docker/issues/22486

## CS Engine 1.11.1-cs1
(27 April 2016)

In this release the CS Engine is supported on RHEL 7.2 OS
