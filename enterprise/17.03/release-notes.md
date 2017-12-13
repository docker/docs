---
title: Docker EE 17.03 release notes
description: Docker Enterprise Edition 17.03 release notes
keywords: docker ee, release notes, upgrade
---

This document describes the latest changes, additions, known issues, and fixes
for Docker Enterprise Edition (Docker EE).

Docker EE is functionally equivalent to the corresponding Docker CE that
it references. However, Docker EE also includes back-ported fixes
(security-related and priority defects) from the open source. It incorporates
defect fixes that you can use in environments where new features cannot be
adopted as quickly for consistency and compatibility reasons.

## Docker EE 17.03.2-ee-8 (2017-12-13)

* Handle cleanup DNS for attachable container to prevent leak in name resolution [docker/libnetwork#1999](https://github.com/docker/libnetwork/pull/1999)
* When a node is removed, delete all of its attachment tasks so networks used by those tasks can be removed [docker/swarmkit#2417](https://github.com/docker/swarmkit/pull/2417)
* Increase gRPC request timeout to 20 seconds for sending snapshots to prevent `context deadline exceeded` errors [docker/swarmkit#2406](https://github.com/docker/swarmkit/pull/2406)
* Avoid using a map for log attributes to prevent panic [moby/moby#34174](https://github.com/moby/moby/pull/34174)
* Fix "raw" mode with the Splunk logging driver [moby/moby#34520](https://github.com/moby/moby/pull/34520)
* Don't unmount entire plugin manager tree on remove [moby/moby#33422](https://github.com/moby/moby/pull/33422)
* Redact secret data on secret creation [moby/moby#33884](https://github.com/moby/moby/pull/33884)
* Sort secrets and configs to ensure idempotence and prevent `docker stack deploy` from useless restart of services [docker/cli#509](https://github.com/docker/cli/pull/509)
* Automatically set `may_detach_mounts=1` on startup to prevent `device or resource busy` errors [moby/moby#34886](https://github.com/moby/moby/pull/34886)
* Don't abort when setting `may_detach_mounts` [moby/moby#35172](https://github.com/moby/moby/pull/35172)
* Protect health monitor channel to prevent engine panic [moby/moby#35482](https://github.com/moby/moby/pull/35482)

## Docker EE 17.03.2-ee-7 (2017-10-04)

* Fix logic in network resource reaping to prevent memory leak [docker/libnetwork#1944](https://github.com/docker/libnetwork/pull/1944) [docker/libnetwork#1960](https://github.com/docker/libnetwork/pull/1960)
* Increase max GRPC message size to 128MB for larger snapshots so newly added managers can successfully join [docker/swarmkit#2375](https://github.com/docker/swarmkit/pull/2375)

## Docker EE 17.03.2-ee-6 (2017-08-24)

* Fix daemon panic on docker image push [moby/moby#33105](https://github.com/moby/moby/pull/33105)
* Fix panic in concurrent network creation/deletion operations [docker/libnetwork#1861](https://github.com/docker/libnetwork/pull/1861)
* Improve network db stability under stressful situations [docker/libnetwork#1860](https://github.com/docker/libnetwork/pull/1860)
* Enable TCP Keep-Alive in Docker client [docker/cli#415](https://github.com/docker/cli/pull/415)
* Lock goroutine to OS thread while changing NS [docker/libnetwork#1911](https://github.com/docker/libnetwork/pull/1911)
* Ignore PullOptions for running tasks [docker/swarmkit#2351](https://github.com/docker/swarmkit/pull/2351)

## Docker EE 17.03.2-ee-5 (20 Jul 2017)

* Add more locking to storage drivers [#31136](https://github.com/moby/moby/pull/31136)
* Prevent data race on `docker network connect/disconnect` [#33456](https://github.com/moby/moby/pull/33456)
* Improve service discovery reliability [#1796](https://github.com/docker/libnetwork/pull/1796) [#18078](https://github.com/docker/libnetwork/pull/1808)
* Fix resource leak in swarm mode [#2215](https://github.com/docker/swarmkit/pull/2215)
* Optimize `docker system df` for volumes on NFS [#33620](https://github.com/moby/moby/pull/33620)
* Fix validation bug with host-mode ports in swarm mode [#2177](https://github.com/docker/swarmkit/pull/2177)
* Fix potential crash in swarm mode [#2268](https://github.com/docker/swarmkit/pull/2268)
* Improve network control-plane reliability [#1704](https://github.com/docker/libnetwork/pull/1704)
* Do not error out when selinux relabeling is not supported on volume filesystem [#33831](https://github.com/moby/moby/pull/33831)
* Remove debugging code for aufs ebusy errors [#31665](https://github.com/moby/moby/pull/31665)
* Prevent resource leak on healthchecks [#33781](https://github.com/moby/moby/pull/33781)
* Fix issue where containerd supervisor may exit prematurely [#32590](https://github.com/moby/moby/pull/32590)
* Fix potential containerd crash [#2](https://github.com/docker/containerd/pull/2)
* Ensure server details are set in client even when an error is returned [#33827](https://github.com/moby/moby/pull/33827)
* Fix issue where slow/dead `docker logs` clients can block the container [#33897](https://github.com/moby/moby/pull/33897)
* Fix potential panic on Windows when running as a service [#32244](https://github.com/moby/moby/pull/32244)

## Docker EE 17.03.2-ee-4 (01 Jun 2017)

Refer to the [detailed list](https://github.com/moby/moby/releases/tag/v17.03.2-ce) of all changes since the release of Docker EE 17.03.1-ee-3

**Note**: This release includes a fix for potential data loss under certain
circumstances with the local (built-in) volume driver.

## Docker EE 17.03.1-ee-3 (30 Mar 2017)

* Fix an issue with the SELinux policy for Oracle Linux [#31501](https://github.com/docker/docker/pull/31501)

## Docker EE 17.03.1-ee-2 (28 Mar 2017)

* Fix issue with swarm CA timeouts [#2063](https://github.com/docker/swarmkit/pull/2063) [#2064](https://github.com/docker/swarmkit/pull/2064/files)

Refer to the [detailed list](https://github.com/moby/moby/releases/tag/v17.03.1-ce) of all changes since the release of Docker EE 17.03.0-ee-1

## Docker EE 17.03.0-ee-1 (2 Mar 2017)

Initial Docker EE release, based on Docker CE 17.03.0

* Optimize size calculation for `docker system df` container size [#31159](https://github.com/docker/docker/pull/31159)
