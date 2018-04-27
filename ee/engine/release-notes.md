---
title: Docker EE Engine release notes
description: Learn about new features, improvements, and known issues in the
  Enterprise Edition of Docker engine.
keywords: ee, release notes, upgrade
redirect_from:
- /enterprise/release-notes/
- /enterprise/17.06/
- /enterprise/17.06/release-notes/
- /enterprise/17.03/release-notes/
---

This document describes the latest changes, additions, known issues, and fixes
for Docker Enterprise Edition (Docker EE).

Docker EE is functionally equivalent to the corresponding Docker CE that
it references. However, Docker EE also includes back-ported fixes
(security-related and priority defects) from the open source. It incorporates
defect fixes that you can use in environments where new features cannot be
adopted as quickly for consistency and compatibility reasons.

## 17.06.2-ee-10 (2018-04-27)

### Runtime

* Fix version output to not have `-dev`.

## 17.06.2-ee-9 (2018-04-26)

### Runtime

* Make Swarm manager Raft quorum parameters configurable in daemon config. [moby/moby#36726](https://github.com/moby/moby/pull/36726)
* Windows: Ignore missing tombstone files when closing an image.
* Windows: Fix directory deletes when a container sharing a base image is running.

### Swarm mode

- Increase raft ElectionTick to 10xHeartbeatTick. [docker/swarmkit#2564](https://github.com/docker/swarmkit/pull/2564)
- Adding logic to restore networks in order. [docker/swarmkit#2584](https://github.com/docker/swarmkit/pull/2584)

## 17.06.2-ee-8 (2018-04-17)

### Networking

- Update libnetwork to fix stale HNS endpoints on Windows. [moby/moby#36603](https://github.com/moby/moby/pull/36603)

### Packaging

* Ensure the graphdriver dir is a shared mount within docker systemd service.

### Known issues

* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.

## 17.06.2-ee-7 (2018-03-19)

### Important notes about this release

- The `overlay2` detection has been improved in this release. On
  Linux distributions where `devicemapper` was the default storage driver,
  `overlay2` is now used by default, if the kernel supports it.

### Logging

* Set timeout on splunk batch send [moby/moby#35496](https://github.com/moby/moby/pull/35496)
- AWS: Fix batch size calculation for large logs[moby/moby#35726](https://github.com/moby/moby/pull/35726)
* Support a proxy in splunk log driver [moby/moby#36220](https://github.com/moby/moby/pull/36220)

### Networking

- Fix NetworkDB node management race condition [docker/libnetwork#2033](https://github.com/docker/libnetwork/pull/2033)
* Update Netlink socket timeout [docker/libnetwork#2044](https://github.com/docker/libnetwork/pull/2044)
- Fix for duplicate IP issues [docker/libnetwork#2105](https://github.com/docker/libnetwork/pull/2105)

### Packaging

+ Add a platform version to `Docker --version` [docker/cli#427](https://github.com/docker/cli/pull/427) and [moby/moby#35705](https://github.com/moby/moby/pull/35705)
* SLES and Ubuntu set TasksMax in docker.service [docker/docker-ce-packaging#78](https://github.com/docker/docker-ce-packaging/pull/78)

### Runtime

* Adjust min TLS Version to v1.2 for PCI compliance [docker/go-connections#45](https://github.com/docker/go-connections/pull/45)
* Fix systemd cgroup after memory type changed [opencontainers/runc#1573](https://github.com/opencontainers/runc/pull/1573)
* Detect overlay2 support on pre-4.0 kernels [moby/moby#35527](https://github.com/moby/moby/pull/35527)
* Enables deferred device deletion/removal by default if the driver version in the kernel supports the feature [moby/moby#33698](https://github.com/moby/moby/pull/33698)
- Fix EBUSY errors under overlayfs and v4.13+ kernels [moby/moby#34914](https://github.com/moby/moby/pull/34914) and [moby/moby#34948](https://github.com/moby/moby/pull/34948)
- Fix TestMount under a selinux system [moby/moby#34965](https://github.com/moby/moby/pull/34965)
- Fix devicemapper error: cannot remove container filesystem, layer not retained [moby/moby#36160](https://github.com/moby/moby/pull/36160)
+ Golang bumped to 1.8.7
* Add timeouts for volume plugin ops [moby/moby#35441](https://github.com/moby/moby/pull/35441)
+ Add `REMOVE` and `ORPHANED` to `TaskState` [moby/moby#36146](https://github.com/moby/moby/pull/36146)
- Fix abort when setting `may_detach_mounts` [moby/moby#35172](https://github.com/moby/moby/pull/35172)
* Windows: Ensure Host Network Service exists [moby/moby#34928](https://github.com/moby/moby/pull/34928)
- Fix issue where network inspect does not show created time in swarm scope [moby/moby#36095](https://github.com/moby/moby/pull/36095)
* Windows: Daemon should respect `DOCKER_TMPDIR` [moby/moby#35077](https://github.com/moby/moby/pull/35077)
- Merge global storage options on create [moby/moby#34508](https://github.com/moby/moby/pull/34508)
- Remove support for overlay/overlay2 without d_type [moby/moby#35514](https://github.com/moby/moby/pull/35514)

### Swarm mode

* Add required call to allocate VIPs when endpoints are restored [docker/swarmkit#2468](https://github.com/docker/swarmkit/pull/2468)
- Synchronize Dispatcher.Stop() with incoming rpcs [docker/swarmkit#2524](https://github.com/docker/swarmkit/pull/2524)
- Fix IP overlap with empty EndpointSpec [docker/swarmkit#2511](https://github.com/docker/swarmkit/pull/2511)

## 17.06.2-ee-6 (2017-11-27)

### Runtime

* Create labels when volume exists only remotely [moby/moby#34896](https://github.com/moby/moby/pull/34896)
* Fix leaking container/exec state [moby/moby#35484](https://github.com/moby/moby/pull/35484)
* Protect health monitor channel to prevent panics [moby/moby#35482](https://github.com/moby/moby/pull/35482)
* Mask `/proc/scsi` path from use in container [moby/moby#35399](https://github.com/moby/moby/pull/35399)
* Fix memory exhaustion when a malformed image could cause the daemon to crash [moby/moby#35424](https://github.com/moby/moby/pull/35424)

### Swarm mode

* Fix deadlock on getting swarm info [moby/moby#35388](https://github.com/moby/moby/issues/35388)
* Only shut down old tasks on success [docker/swarmkit#2308](https://github.com/docker/swarmkit/pull/2308)
* Error on cluster spec name change [docker/swarmkit#2436](https://github.com/docker/swarmkit/pull/2436)

## 17.06.2-ee-5 (2017-11-02)

### Important notes about this release

- Starting with Docker EE 17.06.2-ee-5, Ubuntu, SLES, RHEL packages are also available
  for IBM Power using the ppc64le architecture.

- Docker EE 17.06.2-ee-5 now enables the [telemetry plugin](/enterprise/telemetry/)
  by default on all supported Linux distributions. For more details, including how to
  opt out, see [the documentation](/enterprise/telemetry/).

### Client

* Set APIVersion on the client, even when Ping fails [docker/cli#546](https://github.com/docker/cli/pull/546)

### Logging

* Fix "raw" mode with the Splunk logging driver [moby/moby#34520](https://github.com/moby/moby/pull/34520)

### Networking

* Disable hostname lookup to speed up check if chain chain exists [docker/libnetwork#1974](https://github.com/docker/libnetwork/pull/1974)
* Handle cleanup DNS for attachable container to prevent leak in name resolution [docker/libnetwork#1989](https://github.com/docker/libnetwork/pull/1989)

### Packaging

+ Add telemetry plugin for all linux distributions
+ Fix install of docker-ee on RHEL7 s390x by removing dependency on `container-selinux`

### Runtime

* Automatically set `may_detach_mounts=1` on startup [moby/moby#34886](https://github.com/moby/moby/pull/34886)
* Fallback to use naive diff driver if enable CONFIG_OVERLAY_FS_REDIRECT_DIR [moby/moby#34342](https://github.com/moby/moby/pull/34342)
* Set selinux label on local volumes from mounts API [moby/moby#34684](https://github.com/moby/moby/pull/34684)
* Close pipe in overlay2 graphdriver [moby/moby#34863](https://github.com/moby/moby/pull/34863)
* Relabel config files [moby/moby#34732](https://github.com/moby/moby/pull/34732)
+ Add support for Windows version filtering on pull of docker image [moby/moby#35090](https://github.com/moby/moby/pull/35090)

### Swarm mode

* Increase gRPC request timeout to 20 seconds for sending snapshots to prevent `context deadline exceeded` errors [docker/swarmkit#2391](https://github.com/docker/swarmkit/pull/2391)
* When a node is removed, delete all of its attachment tasks so networks used by those tasks can be removed [docker/swarmkit#2414](https://github.com/docker/swarmkit/pull/2414)

### Known issues

 * It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
 * Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.

## 17.06.2-ee-4 (2017-10-12)

### Client

* Fix idempotence of `docker stack deploy` when secrets or configs are used [docker/cli#509](https://github.com/docker/cli/pull/509)

### Logging

* Avoid using a map for log attributes to prevent panic [moby/moby#34174](https://github.com/moby/moby/pull/34174)

### Networking

* Fix for garbage collection logic in NetworkDB. Entries were not properly garbage collected and deleted within the expected time [docker/libnetwork#1944](https://github.com/docker/libnetwork/pull/1944) [docker/libnetwork#1960](https://github.com/docker/libnetwork/pull/1960)
* Allow configuration of max packet size in network DB to use the full available MTU. This requires a configuration in the docker daemon and need a dockerd restart [docker/libnetwork#1839](https://github.com/docker/libnetwork/pull/1839)
* Overlay fix for transient IP reuse [docker/libnetwork#1935](https://github.com/docker/libnetwork/pull/1935) [docker/libnetwork#1968](https://github.com/docker/libnetwork/pull/1968)
* Serialize IP allocation [docker/libnetwork#1788](https://github.com/docker/libnetwork/pull/1788)


## 17.06.2-ee-3 (2017-09-22)

### Swarm mode

- Increase max message size to allow larger snapshots [docker/swarmkit#131](https://github.com/docker/swarmkit/pull/131)

## 17.06.1-ee-2 (2017-08-24)

### Client

- Enable TCP Keep-Alive in Docker client [#415](https://github.com/docker/cli/pull/415)

### Networking

- Lock goroutine to OS thread while changing NS [#1911](https://github.com/docker/libnetwork/pull/1911)

### Runtime

- devmapper: ensure that UdevWait is called after calls to setCookie [#33732](https://github.com/moby/moby/pull/33732)
- aufs: ensure diff layers are correctly removed to prevent leftover files from using up storage [#34587](https://github.com/moby/moby/pull/34587)

### Swarm mode

- Ignore PullOptions for running tasks [#2351](https://github.com/docker/swarmkit/pull/2351)

## 17.06.1-ee (2017-08-16)

### Important notes about this release

- Starting with Docker EE 17.06.1, Ubuntu, SLES, RHEL packages are also available
  for IBM Z using the s390x architecture.

- Docker EE 17.06.1 includes a new [telemetry plugin](/enterprise/telemetry/)
  which is enabled by default on Ubuntu hosts. For more details, including how to
  opt out, see [the documentation(/enterprise/telemetry/).

- Docker 17.06 by default disables communication with legacy (v1)
  registries. If you require interaction with registries that have not yet
  migrated to the v2 protocol, set the `--disable-legacy-registry=false` daemon
  option.


### Builder

+ Add `--iidfile` option to docker build. It allows specifying a location where to save the resulting image ID
+ Allow specifying any remote ref in git checkout URLs [#32502](https://github.com/moby/moby/pull/32502)
+ Add multi-stage build support [#31257](https://github.com/moby/moby/pull/31257) [#32063](https://github.com/moby/moby/pull/32063)
+ Allow using build-time args (`ARG`) in `FROM` [#31352](https://github.com/moby/moby/pull/31352)
+ Add an option for specifying build target [#32496](https://github.com/moby/moby/pull/32496)
* Accept `-f -` to read Dockerfile from `stdin`, but use local context for building [#31236](https://github.com/moby/moby/pull/31236)
* The values of default build time arguments (e.g `HTTP_PROXY`) are no longer displayed in docker image history unless a corresponding `ARG` instruction is written in the Dockerfile. [#31584](https://github.com/moby/moby/pull/31584)
- Fix setting command if a custom shell is used in a parent image [#32236](https://github.com/moby/moby/pull/32236)
- Fix `docker build --label` when the label includes single quotes and a space [#31750](https://github.com/moby/moby/pull/31750)
* Disable container logging for build containers [#29552](https://github.com/moby/moby/pull/29552)
* Fix use of `**/` in `.dockerignore` [#29043](https://github.com/moby/moby/pull/29043)
* Fix a regression, where `ADD` from remote URL's extracted archives [#89](https://github.com/docker/docker-ce/pull/89)
* Fix handling of remote "git@" notation [#100](https://github.com/docker/docker-ce/pull/100)
* Fix copy `--from` conflict with force pull [#86](https://github.com/docker/docker-ce/pull/86)

### Client

+ Add `--format` option to `docker stack ls` [#31557](https://github.com/moby/moby/pull/31557)
+ Add support for labels in compose initiated builds [#32632](https://github.com/moby/moby/pull/32632) [#32972](https://github.com/moby/moby/pull/32972)
+ Add `--format` option to `docker history` [#30962](https://github.com/moby/moby/pull/30962)
+ Add `--format` option to `docker system df` [#31482](https://github.com/moby/moby/pull/31482)
+ Allow specifying Nameservers and Search Domains in stack files [#32059](https://github.com/moby/moby/pull/32059)
+ Add support for `read_only` service to `docker stack deploy` [#docker/cli/73](https://github.com/docker/cli/pull/73)
* Display Swarm cluster and node TLS information [#docker/cli/44](https://github.com/docker/cli/pull/44)
+ Add support for placement preference to `docker stack deploy` [#docker/cli/35](https://github.com/docker/cli/pull/35)
+ Add new `ca ` subcommand to `docker swarm` to allow managing a swarm CA [#docker/cli/48](https://github.com/docker/cli/pull/48)
+ Add credential-spec to compose [#docker/cli/71](https://github.com/docker/cli/pull/71)
+ Add support for csv format options to `--network` and `--network-add` [#docker/cli/62](https://github.com/docker/cli/pull/62) [#33130](https://github.com/moby/moby/pull/33130)
- Fix stack compose bind-mount volumes on Windows [#docker/cli/136](https://github.com/docker/cli/pull/136)
- Correctly handle a Docker daemon without registry info [#docker/cli/126](https://github.com/docker/cli/pull/126)
+ Allow `--detach` and `--quiet` flags when using --rollback [#docker/cli/144](https://github.com/docker/cli/pull/144)
+ Remove deprecated `--email` flag from `docker login` [#docker/cli/143](https://github.com/docker/cli/pull/143)
* Adjusted `docker stats` memory output [#docker/cli/80](https://github.com/docker/cli/pull/80)
* Add `--mount` flag to `docker run` and `docker create` [#32251](https://github.com/moby/moby/pull/32251)
* Add `--type=secret` to `docker inspect` [#32124](https://github.com/moby/moby/pull/32124)
* Add `--format` option to `docker secret ls` [#31552](https://github.com/moby/moby/pull/31552)
* Add `--filter` option to `docker secret ls` [#30810](https://github.com/moby/moby/pull/30810)
* Add `--filter scope=<swarm|local>` to `docker network ls` [#31529](https://github.com/moby/moby/pull/31529)
* Add `--cpus` support to `docker update` [#31148](https://github.com/moby/moby/pull/31148)
* Add label filter to `docker system prune` and other `prune` commands [#30740](https://github.com/moby/moby/pull/30740)
* `docker stack rm` now accepts multiple stacks as input [#32110](https://github.com/moby/moby/pull/32110)
* Improve `docker version --format` option when the client has downgraded the API version [#31022](https://github.com/moby/moby/pull/31022)
* Prompt when using an encrypted client certificate to connect to a docker daemon [#31364](https://github.com/moby/moby/pull/31364)
* Display created tags on successful `docker build` [#32077](https://github.com/moby/moby/pull/32077)
* Cleanup compose convert error messages [#32087](https://github.com/moby/moby/pull/32087)
+ Sort `docker stack ls` by name [#31085](https://github.com/moby/moby/pull/31085)
+ Flags for specifying bind mount consistency [#31047](https://github.com/moby/moby/pull/31047)
* Output of docker CLI --help is now wrapped to the terminal width [#28751](https://github.com/moby/moby/pull/28751)
* Suppress image digest in docker ps [#30848](https://github.com/moby/moby/pull/30848)
* Hide command options that are related to Windows [#30788](https://github.com/moby/moby/pull/30788)
* Fix `docker plugin install` prompt to accept "enter" for the "N" default [#30769](https://github.com/moby/moby/pull/30769)
+ Add `truncate` function for Go templates [#30484](https://github.com/moby/moby/pull/30484)
* Support expanded syntax of ports in `stack deploy` [#30476](https://github.com/moby/moby/pull/30476)
* Support expanded syntax of mounts in `stack deploy` [#30597](https://github.com/moby/moby/pull/30597) [#31795](https://github.com/moby/moby/pull/31795)
+ Add `--add-host` for docker build [#30383](https://github.com/moby/moby/pull/30383)
+ Add `.CreatedAt` placeholder for `docker network ls --format` [#29900](https://github.com/moby/moby/pull/29900)
* Update order of `--secret-rm` and `--secret-add` [#29802](https://github.com/moby/moby/pull/29802)
+ Add `--filter enabled=true` for `docker plugin ls` [#28627](https://github.com/moby/moby/pull/28627)
+ Add `--format` to `docker service ls` [#28199](https://github.com/moby/moby/pull/28199)
+ Add `publish` and `expose` filter for `docker ps --filter` [#27557](https://github.com/moby/moby/pull/27557)
* Support multiple service IDs on `docker service ps` [#25234](https://github.com/moby/moby/pull/25234)
+ Allow swarm join with `--availability=drain` [#24993](https://github.com/moby/moby/pull/24993)
* Docker inspect now shows "docker-default" when AppArmor is enabled and no other profile was defined [#27083](https://github.com/moby/moby/pull/27083)
* Make pruning volumes optional when running `docker system prune`, and add a `--volumes` flag [#109](https://github.com/docker/docker-ce/pull/109)
* Show progress of replicated tasks before they are assigned [#97](https://github.com/docker/docker-ce/pull/97)
* Fix `docker wait` hanging if the container does not exist [#106](https://github.com/docker/docker-ce/pull/106)
* If `docker swarm ca` is called without the `--rotate` flag, warn if other flags are passed [#110](https://github.com/docker/docker-ce/pull/110)
* Fix API version negotiation not working if the daemon returns an error [#115](https://github.com/docker/docker-ce/pull/115)
* Print an error if "until" filter is combined with "--volumes" on system prune [#154](https://github.com/docker/docker-ce/pull/154)


### Contrib

+ Add support for building docker debs for Ubuntu 17.04 Zesty on amd64 [#32435](https://github.com/moby/moby/pull/32435)

### Daemon

- Fix `--api-cors-header` being ignored if `--api-enable-cors` is not set [#32174](https://github.com/moby/moby/pull/32174)
- Cleanup docker tmp dir on start [#31741](https://github.com/moby/moby/pull/31741)
- Deprecate `--graph` flag in favor or `--data-root` [#28696](https://github.com/moby/moby/pull/28696)

### Distribution

* Select digest over tag when both are provided during a pull [#33214](https://github.com/moby/moby/pull/33214)

### Logging

+ Add monitored resource type metadata for GCP logging driver [#32930](https://github.com/moby/moby/pull/32930)
+ Add multiline processing to the AWS CloudWatch logs driver [#30891](https://github.com/moby/moby/pull/30891)
+ Add support for logging driver plugins [#28403](https://github.com/moby/moby/pull/28403)
* Add support for showing logs of individual tasks to `docker service logs`, and add `/task/{id}/logs` REST endpoint [#32015](https://github.com/moby/moby/pull/32015)
* Add `--log-opt env-regex` option to match environment variables using a regular expression [#27565](https://github.com/moby/moby/pull/27565)
+ Implement optional ring buffer for container logs [#28762](https://github.com/moby/moby/pull/28762)
+ Add `--log-opt awslogs-create-group=<true|false>` for awslogs (CloudWatch) to support creation of log groups as needed [#29504](https://github.com/moby/moby/pull/29504)
- Fix segfault when using the gcplogs logging driver with a "static" binary [#29478](https://github.com/moby/moby/pull/29478)
* Fix stderr logging for `journald` and `syslog` [#95](https://github.com/docker/docker-ce/pull/95)
* Fix log readers can block writes indefinitely [#98](https://github.com/docker/docker-ce/pull/98)
* Fix `awslogs` driver repeating last event [#151](https://github.com/docker/docker-ce/pull/151)

### Networking

+ Add Support swarm-mode services with node-local networks such as macvlan, ipvlan, bridge, host [#32981](https://github.com/moby/moby/pull/32981)
+ Pass driver-options to network drivers on service creation [#32981](https://github.com/moby/moby/pull/33130)
+ Isolate Swarm Control-plane traffic from Application data traffic using --data-path-addr [#32717](https://github.com/moby/moby/pull/32717)
* Several improvements to Service Discovery [#docker/libnetwork/1796](https://github.com/docker/libnetwork/pull/1796)
+ Allow user to replace, and customize the ingress network [#31714](https://github.com/moby/moby/pull/31714)
- Fix UDP traffic in containers not working after the container is restarted [#32505](https://github.com/moby/moby/pull/32505)
- Fix files being written to `/var/lib/docker` if a different data-root is set [#32505](https://github.com/moby/moby/pull/32505)
* Check parameter `--ip`, `--ip6` and `--link-local-ip` in `docker network connect` [#30807](https://github.com/moby/moby/pull/30807)
+ Added support for `dns-search` [#30117](https://github.com/moby/moby/pull/30117)
+ Added --verbose option for docker network inspect to show task details from all swarm nodes [#31710](https://github.com/moby/moby/pull/31710)
* Clear stale datapath encryption states when joining the cluster [docker/libnetwork#1354](https://github.com/docker/libnetwork/pull/1354)
+ Ensure iptables initialization only happens once [docker/libnetwork#1676](https://github.com/docker/libnetwork/pull/1676)
* Fix bad order of iptables filter rules [docker/libnetwork#961](https://github.com/docker/libnetwork/pull/961)
+ Add anonymous container alias to service record on attachable network [docker/libnetwork#1651](https://github.com/docker/libnetwork/pull/1651)
+ Support for `com.docker.network.container_interface_prefix` driver label [docker/libnetwork#1667](https://github.com/docker/libnetwork/pull/1667)
+ Improve network list performance by omitting network details that are not used [#30673](https://github.com/moby/moby/pull/30673)
* Fix issue with driver options not received by network drivers [#127](https://github.com/docker/docker-ce/pull/127)

### Packaging

+ Rely on `container-selinux` on Centos/Fedora/RHEL when available [#32437](https://github.com/moby/moby/pull/32437)

### Plugins

* Make plugin removes more resilient to failure [#91](https://github.com/docker/docker-ce/pull/91)

### Runtime

+ Add build & engine info prometheus metrics [#32792](https://github.com/moby/moby/pull/32792)
* Update containerd to d24f39e203aa6be4944f06dd0fe38a618a36c764 [#33007](https://github.com/moby/moby/pull/33007)
* Update runc to 992a5be178a62e026f4069f443c6164912adbf09 [#33007](https://github.com/moby/moby/pull/33007)
+ Add option to auto-configure blkdev for devmapper [#31104](https://github.com/moby/moby/pull/31104)
+ Add log driver list to `docker info` [#32540](https://github.com/moby/moby/pull/32540)
+ Add API endpoint to allow retrieving an image manifest [#32061](https://github.com/moby/moby/pull/32061)
* Do not remove container from memory on error with `forceremove` [#31012](https://github.com/moby/moby/pull/31012)
+ Add support for metric plugins [#32874](https://github.com/moby/moby/pull/32874)
* Return an error when an invalid filter is given to `prune` commands [#33023](https://github.com/moby/moby/pull/33023)
+ Add daemon option to allow pushing foreign layers [#33151](https://github.com/moby/moby/pull/33151)
- Fix an issue preventing containerd to be restarted after it died [#32986](https://github.com/moby/moby/pull/32986)
+ Add cluster events to Docker event stream. [#32421](https://github.com/moby/moby/pull/32421)
+ Add support for DNS search on windows [#33311](https://github.com/moby/moby/pull/33311)
* Upgrade to Go 1.8.3 [#33387](https://github.com/moby/moby/pull/33387)
- Prevent a containerd crash when journald is restarted [#33007](https://github.com/moby/moby/pull/33007)
- Fix healthcheck failures due to invalid environment variables [#33249](https://github.com/moby/moby/pull/33249)
- Prevent a directory to be created in lieu of the daemon socket when a container mounting it is to be restarted during a shutdown [#30348](https://github.com/moby/moby/pull/33330)
- Prevent a container to be restarted upon stop if its stop signal is set to `SIGKILL` [#33335](https://github.com/moby/moby/pull/33335)
- Ensure log drivers get passed the same filename to both StartLogging and StopLogging endpoints [#33583](https://github.com/moby/moby/pull/33583)
- Remove daemon data structure dump on `SIGUSR1` to avoid a panic [#33598](https://github.com/moby/moby/pull/33598)
- Ensure health probe is stopped when a container exits [#32274](https://github.com/moby/moby/pull/32274)
* Handle paused container when restoring without live-restore set [#31704](https://github.com/moby/moby/pull/31704)
- Do not allow sub second in healthcheck options in Dockerfile [#31177](https://github.com/moby/moby/pull/31177)
* Support name and id prefix in `secret update` [#30856](https://github.com/moby/moby/pull/30856)
* Use binary frame for websocket attach endpoint [#30460](https://github.com/moby/moby/pull/30460)
* Fix linux mount calls not applying propagation type changes [#30416](https://github.com/moby/moby/pull/30416)
* Fix ExecIds leak on failed `exec -i` [#30340](https://github.com/moby/moby/pull/30340)
* Prune named but untagged images if `danglingOnly=true` [#30330](https://github.com/moby/moby/pull/30330)
+ Add daemon flag to set `no_new_priv` as default for unprivileged containers [#29984](https://github.com/moby/moby/pull/29984)
+ Add daemon option `--default-shm-size` [#29692](https://github.com/moby/moby/pull/29692)
+ Support registry mirror config reload [#29650](https://github.com/moby/moby/pull/29650)
- Ignore the daemon log config when building images [#29552](https://github.com/moby/moby/pull/29552)
* Move secret name or ID prefix resolving from client to daemon [#29218](https://github.com/moby/moby/pull/29218)
+ Add the ability to specify extra rules for a container device `cgroup devices.allow` mechanism [#22563](https://github.com/moby/moby/pull/22563)
- Fix `cpu.cfs_quota_us` being reset when running `systemd daemon-reload` [#31736](https://github.com/moby/moby/pull/31736)
* Prevent a `goroutine` leak when `healthcheck` gets stopped [#90](https://github.com/docker/docker-ce/pull/90)
* Do not error on relabel when relabel not supported [#92](https://github.com/docker/docker-ce/pull/92)
* Limit max backoff delay to 2 seconds for GRPC connection [#94](https://github.com/docker/docker-ce/pull/94)
* Fix issue preventing containers to run when memory cgroup was specified due to bug in certain kernels [#102](https://github.com/docker/docker-ce/pull/102)
* Fix container not responding to SIGKILL when paused [#102](https://github.com/docker/docker-ce/pull/102)
* Improve error message if an image for an incompatible OS is loaded [#108](https://github.com/docker/docker-ce/pull/108)
* Fix a handle leak in `go-winio` [#112](https://github.com/docker/docker-ce/pull/112)
* Fix issue upon upgrade, preventing docker from showing running containers when `--live-restore` is enabled [#117](https://github.com/docker/docker-ce/pull/117)
* Fix bug where services using secrets would fail to start on daemons using the `userns-remap` feature [#121](https://github.com/docker/docker-ce/pull/121)
* Fix error handling with `not-exist` errors on remove [#142](https://github.com/docker/docker-ce/pull/142)
* Fix REST API Swagger representation cannot be loaded with SwaggerUI [#156](https://github.com/docker/docker-ce/pull/156)

### Security

+ Allow personality with UNAME26 bit set in default seccomp profile [#32965](https://github.com/moby/moby/pull/32965)
* Allow setting SELinux type or MCS labels when using `--ipc=container:` or `--ipc=host` [#30652](https://github.com/moby/moby/pull/30652)
* Redact secret data on secret creation [#99](https://github.com/docker/docker-ce/pull/99)

### Swarm mode

+ Add an option to allow specifying a different interface for the data traffic (as opposed to control traffic) [#32717](https://github.com/moby/moby/pull/32717)
* Allow specifying a secret location within the container [#32571](https://github.com/moby/moby/pull/32571)
+ Add support for secrets on Windows [#32208](https://github.com/moby/moby/pull/32208)
+ Add TLS Info to swarm info and node info endpoint [#32875](https://github.com/moby/moby/pull/32875)
+ Add support for services to carry arbitrary config objects [#32336](https://github.com/moby/moby/pull/32336), [#docker/cli/45](https://github.com/docker/cli/pull/45),[#33169](https://github.com/moby/moby/pull/33169)
+ Add API to rotate swarm CA certificate [#32993](https://github.com/moby/moby/pull/32993)
* Service digest pining is now handled client side [#32388](https://github.com/moby/moby/pull/32388), [#33239](https://github.com/moby/moby/pull/33239)
+ Placement now also take platform in account [#33144](https://github.com/moby/moby/pull/33144)
- Fix possible hang when joining fails [#docker-ce/19](https://github.com/docker/docker-ce/pull/19)
- Fix an issue preventing external CA to be accepted [#33341](https://github.com/moby/moby/pull/33341)
- Fix possible orchestration panic in mixed version clusters [#swarmkit/2233](https://github.com/docker/swarmkit/pull/2233)
- Avoid assigning duplicate IPs during initialization [#swarmkit/2237](https://github.com/docker/swarmkit/pull/2237)
+ Add update/rollback order for services (`--update-order` / `--rollback-order`) [#30261](https://github.com/moby/moby/pull/30261)
+ Add support for synchronous `service create` and `service update` [#31144](https://github.com/moby/moby/pull/31144)
+ Add support for "grace periods" on healthchecks through the `HEALTHCHECK --start-period` and `--health-start-period` flag to
  `docker service create`, `docker service update`, `docker create`, and `docker run` to support containers with an initial startup
  time [#28938](https://github.com/moby/moby/pull/28938)
* `docker service create` now omits fields that are not specified by the user, when possible. This allows defaults to be applied inside the manager [#32284](https://github.com/moby/moby/pull/32284)
* `docker service inspect` now shows default values for fields that are not specified by the user [#32284](https://github.com/moby/moby/pull/32284)
* Move `docker service logs` out of experimental [#32462](https://github.com/moby/moby/pull/32462)
* Add support for Credential Spec and SELinux to services to the API [#32339](https://github.com/moby/moby/pull/32339)
* Add `--entrypoint` flag to `docker service create` and `docker service update` [#29228](https://github.com/moby/moby/pull/29228)
* Add `--network-add` and `--network-rm` to `docker service update` [#32062](https://github.com/moby/moby/pull/32062)
* Add `--credential-spec` flag to `docker service create` and `docker service update` [#32339](https://github.com/moby/moby/pull/32339)
* Add `--filter mode=<global|replicated>` to `docker service ls` [#31538](https://github.com/moby/moby/pull/31538)
* Resolve network IDs on the client side, instead of in the daemon when creating services [#32062](https://github.com/moby/moby/pull/32062)
* Add `--format` option to `docker node ls` [#30424](https://github.com/moby/moby/pull/30424)
* Add `--prune` option to `docker stack deploy` to remove services that are no longer defined in the docker-compose file [#31302](https://github.com/moby/moby/pull/31302)
* Add `PORTS` column for `docker service ls` when using `ingress` mode [#30813](https://github.com/moby/moby/pull/30813)
- Fix unnescessary re-deploying of tasks when environment-variables are used [#32364](https://github.com/moby/moby/pull/32364)
- Fix `docker stack deploy` not supporting `endpoint_mode` when deploying from a docker compose file [#32333](https://github.com/moby/moby/pull/32333)
- Proceed with startup if cluster component cannot be created to allow recovering from a broken swarm setup [#31631](https://github.com/moby/moby/pull/31631)
+ Topology-aware scheduling [#30725](https://github.com/moby/moby/pull/30725)
+ Automatic service rollback on failure [#31108](https://github.com/moby/moby/pull/31108)
+ Worker and manager on the same node are now connected through a UNIX socket [docker/swarmkit#1828](https://github.com/docker/swarmkit/pull/1828), [docker/swarmkit#1850](https://github.com/docker/swarmkit/pull/1850), [docker/swarmkit#1851](https://github.com/docker/swarmkit/pull/1851)
* Improve raft transport package [docker/swarmkit#1748](https://github.com/docker/swarmkit/pull/1748)
* No automatic manager shutdown on demotion/removal [docker/swarmkit#1829](https://github.com/docker/swarmkit/pull/1829)
* Use TransferLeadership to make leader demotion safer [docker/swarmkit#1939](https://github.com/docker/swarmkit/pull/1939)
* Decrease default monitoring period [docker/swarmkit#1967](https://github.com/docker/swarmkit/pull/1967)
+ Add Service logs formatting [#31672](https://github.com/moby/moby/pull/31672)
* Fix service logs API to be able to specify stream [#31313](https://github.com/moby/moby/pull/31313)
+ Add `--stop-signal` for `service create` and `service update` [#30754](https://github.com/moby/moby/pull/30754)
+ Add `--read-only` for `service create` and `service update` [#30162](https://github.com/moby/moby/pull/30162)
+ Renew the context after communicating with the registry [#31586](https://github.com/moby/moby/pull/31586)
+ (experimental) Add `--tail` and `--since` options to `docker service logs` [#31500](https://github.com/moby/moby/pull/31500)
+ (experimental) Add `--no-task-ids` and `--no-trunc` options to `docker service logs` [#31672](https://github.com/moby/moby/pull/31672)
* Do not add duplicate platform information to service spec [#107](https://github.com/docker/docker-ce/pull/107)
* Cluster update and memory issue fixes [#114](https://github.com/docker/docker-ce/pull/114)
* Changing get network request to return predefined network in swarm [#150](https://github.com/docker/docker-ce/pull/150)

### Windows

* Block pulling Windows images on non-Windows daemons [#29001](https://github.com/moby/moby/pull/29001)

### Deprecation

* Disable legacy registry (v1) by default [#33629](https://github.com/moby/moby/pull/33629)
- Deprecate `--api-enable-cors` daemon flag. This flag was marked deprecated in Docker 1.6.0 but not listed in deprecated features [#32352](https://github.com/moby/moby/pull/32352)
- Remove Ubuntu 12.04 (Precise Pangolin) as supported platform. Ubuntu 12.04 is EOL, and no longer receives updates [#32520](https://github.com/moby/moby/pull/32520)

### Known issues

If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

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
