---
title: Docker EE 17.03 release notes
description: Docker Enterprise Edition release notes
keywords: docker, docker-ee, engine, install, release notes, enterprise
redirect_from:
- /enterprise/release-notes/
---

This document describes the latest changes, additions, known issues, and fixes
for Docker Enterprise Edition (Docker EE).

Docker EE is functionally equivalent to the corresponding Docker CE that
it references. However, Docker EE also includes back-ported fixes
(security-related and priority defects) from the open source. It incorporates
defect fixes that you can use in environments where new features cannot be
adopted as quickly for consistency and compatibility reasons.

## Docker EE 17.03.2-ee-4
(01 Jun 2017)

Refer to the [detailed list](https://github.com/moby/moby/releases/tag/v17.03.2-ce) of all changes since the release of Docker EE 17.03.1-ee-3

*Note*: This release includes a fix for potential data loss under certain
circumstances with the local (built-in) volume driver.

## Docker EE 17.03.1-ee-3
(30 Mar 2017)

* Fix an issue with the SELinux policy for Oracle Linux [#31501](https://github.com/docker/docker/pull/31501)

## Docker EE 17.03.1-ee-2
(28 March 2017)

* Fix issue with swarm CA timeouts [#2063](https://github.com/docker/swarmkit/pull/2063) [#2064](https://github.com/docker/swarmkit/pull/2064/files)

Refer to the [detailed list](https://github.com/moby/moby/releases/tag/v17.03.1-ce) of all changes since the release of Docker EE 17.03.0-ee-1

## Docker EE 17.03.0-ee-1

(2 Mar 2017)

Initial Docker EE release, based on Docker CE 17.03.0

* Optimize size calculation for `docker system df` container size [#31159](https://github.com/docker/docker/pull/31159)
