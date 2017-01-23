---
description: Commercially supported Docker Engine release notes
keywords: docker, documentation, about, technology, understanding, enterprise, hub, registry, Commercially Supported Docker Engine, release notes
redirect_from:
- /docker-trusted-registry/cse-release-notes/
- /docker-trusted-registry/cs-engine/release-notes/release-notes/
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

## CS Engine 1.12.6-cs6
(10 Jan 2017)

Bumps RunC version to address CVE-2016-9962.

Refer to the [detailed list](https://github.com/docker/docker/releases/tag/v1.12.6) of all
changes since the release of CS Engine 1.12.5-cs5.

## CS Engine 1.12.5-cs5
(21 Dec 2016)

Refer to the [detailed list](https://github.com/docker/docker/releases/tag/v1.12.5) of all
changes since the release of CS Engine 1.12.3-cs4

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
changes since the release of CS Engine 1.11.2-cs5

This release addresses the following issues:

* [#25962](https://github.com/docker/docker/pull/25962) Allow normal containers
to connect to swarm-mode overlay network
* Various bug fixes in swarm mode networking

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

This releases addresses the following issues:

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
