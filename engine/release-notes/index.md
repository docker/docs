---
title: Docker Engine release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Engine - Community
keywords: docker, docker engine, ce, whats new, release notes
toc_min: 1
toc_max: 2
skip_read_time: true
redirect_from:
 - /release-notes/docker-ce/
---

This document describes the latest changes, additions, known issues, and fixes
for Docker Engine - Community.

> **Note:**
> The client and container runtime are now in separate packages from the daemon
> in Docker Engine 18.09. Users should install and update all three packages at
> the same time to get the latest patch releases. For example, on Ubuntu:
> `sudo apt install docker-ce docker-ce-cli containerd.io`. See the install
> instructions for the corresponding linux distro for details.

# Version 19.03

## 19.03.8
2020-03-10

### Runtime

- Improve mitigation for [CVE-2019-14271](https://nvd.nist.gov/vuln/detail/CVE-2019-14271) for some nscd configuration.

## 19.03.7
2020-03-03

### Builder

- builder-next: Fix deadlock issues in corner cases. [moby/moby#40557](https://github.com/moby/moby/pull/40557)

### Runtime

* overlay: remove modprobe execs. [moby/moby#40462](https://github.com/moby/moby/pull/40462)
* selinux: display better error messages when setting file labels. [moby/moby#40547](https://github.com/moby/moby/pull/40547)
* Speed up initial stats collection. [moby/moby#40549](https://github.com/moby/moby/pull/40549)
- rootless: use certs.d from XDG_CONFIG_HOME. [moby/moby#40461](https://github.com/moby/moby/pull/40461)
- Bump Golang 1.12.17. [moby/moby#40533](https://github.com/moby/moby/pull/40533) 
- Bump google.golang.org/grpc to v1.23.1. [moby/moby#40566](https://github.com/moby/moby/pull/40566)
- Update containerd binary to v1.2.13. [moby/moby#40540](https://github.com/moby/moby/pull/40540)
- Prevent showing stopped containers as running in an edge case. [moby/moby#40555](https://github.com/moby/moby/pull/40555)
- Prevent potential lock. [moby/moby#40604](https://github.com/moby/moby/pull/40604)

### Client

- Bump Golang 1.12.17. [docker/cli#2342](https://github.com/docker/cli/pull/2342)
- Bump google.golang.org/grpc to v1.23.1. [docker/cli#1884](https://github.com/docker/cli/pull/1884) [docker/cli#2373](https://github.com/docker/cli/pull/2373)

## 19.03.6
2020-02-12

### Builder

- builder-next: Allow modern sign hashes for ssh forwarding. [docker/engine#453](https://github.com/docker/engine/pull/453)
- builder-next: Clear onbuild rules after triggering. [docker/engine#453](https://github.com/docker/engine/pull/453)
- builder-next: Fix issue with directory permissions when usernamespaces is enabled. [moby/moby#40440](https://github.com/moby/moby/pull/40440)
- Bump hcsshim to fix docker build failing on Windows 1903. [docker/engine#429](https://github.com/docker/engine/pull/429)

### Networking

- Shorten controller ID in exec-root to not hit UNIX_PATH_MAX. [docker/engine#424](https://github.com/docker/engine/pull/424)
- Fix panic in drivers/overlay/encryption.go. [docker/engine#424](https://github.com/docker/engine/pull/424)
- Fix hwaddr set race between us and udev. [docker/engine#439](https://github.com/docker/engine/pull/439)

### Runtime

* Bump Golang 1.12.16. [moby/moby#40433](https://github.com/moby/moby/pull/40433)
* Update containerd binary to v1.2.12. [moby/moby#40433](https://github.com/moby/moby/pull/40453)
* Update to runc v1.0.0-rc10. [moby/moby#40433](https://github.com/moby/moby/pull/40453)
- Fix possible runtime panic in Lgetxattr. [docker/engine#454](https://github.com/docker/engine/pull/454)
- rootless: fix proxying UDP packets. [docker/engine#434](https://github.com/docker/engine/pull/434)

## 19.03.5
2019-11-14

### Builder

* builder-next: Added `entitlements` in builder config. [docker/engine#412](https://github.com/docker/engine/pull/412)
* Fix builder-next: permission errors on using build secrets or ssh forwarding with userns-remap. [docker/engine#420](https://github.com/docker/engine/pull/420)
* Fix builder-next: copying a symlink inside an already copied directory. [docker/engine#420](https://github.com/docker/engine/pull/420)

### Packaging

* Support RHEL 8 packages

### Runtime

* Bump Golang to 1.12.12. [docker/engine#418](https://github.com/docker/engine/pull/418)
* Update to RootlessKit to v0.7.0 to harden slirp4netns with mount namespace and seccomp. [docker/engine#397](https://github.com/docker/engine/pull/397)
* Fix to propagate GetContainer error from event processor. [docker/engine#407](https://github.com/docker/engine/pull/407)
* Fix push of OCI image. [docker/engine#405](https://github.com/docker/engine/pull/405)

## 19.03.4
2019-10-17

### Networking

* Rollback libnetwork changes to fix `DOCKER-USER` iptables chain issue. [docker/engine#404](https://github.com/docker/engine/pull/404)

### Known Issues

#### Existing

* In some circumstances with large clusters, Docker information might, as part of the Swarm section,
  include the error `code = ResourceExhausted desc = grpc: received message larger than
  max (5351376 vs. 4194304)`. This does not indicate any failure or misconfiguration by the user,
  and requires no response.
* Orchestrator port conflict can occur when redeploying all services as new. Due to many Swarm manager
  requests in a short amount of time, some services are not able to receive traffic and are causing a `404`
  error after being deployed.
     - **Workaround:** restart all tasks via `docker service update --force`.
* [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack with directory traversal. Workaround until proper fix is available in upcoming patch release: `docker pause` container before doing file operations. [moby/moby#39252](https://github.com/moby/moby/pull/39252)
* `docker cp` regression due to CVE mitigation. An error is produced when the source of `docker cp` is set to `/`.
* Install Docker Engine - Enterprise fails to install on RHEL on Azure. This affects any RHEL version that uses an Extended Update Support (EUS) image. At the time of this writing, known versions affected are RHEL 7.4, 7.5, and 7.6.

     - **Workaround options:**
         - Use an older image and don't get updates. Examples of EUS images are here: https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#rhel-images-with-eus.
         - Import your own RHEL images into Azure and do not rely on the Extended Update Support (EUS) RHEL images.
         - Use a RHEL image that does not contain a minor version in the SKU. These are not attached to EUS repositories. Some examples of those are the first three images (SKUs: 7-RAW, 7-LVM, 7-RAW-CI) listed here : https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#list-of-rhel-images-available.


## 19.03.3
2019-10-08

### Security

* Patched `runc` in containerd. [CVE-2017-18367](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2017-18367)

### Builder

* Fix builder-next: resolve digest for third party registries. [docker/engine#339](https://github.com/docker/engine/pull/339)

* Fix builder-next: user namespace builds when daemon started with socket activation. [docker/engine#373](https://github.com/docker/engine/pull/373)

* Fix builder-next; session: release forwarded ssh socket connection per connection. [docker/engine#373](https://github.com/docker/engine/pull/373)

* Fix build-next: llbsolver: error on multiple cache importers. [docker/engine#373](https://github.com/docker/engine/pull/373)

### Client

* Added support for Docker Template 0.1.6.

* Mitigate against YAML files that have excessive aliasing. [docker/cli#2119](https://github.com/docker/cli/pull/2119)

### Runtime

* Bump Golang to 1.12.10. [docker/engine#387](https://github.com/docker/engine/pull/387)

* Bump containerd to 1.2.10. [docker/engine#385](https://github.com/docker/engine/pull/385)

* Distribution: modify warning logic when pulling v2 schema1 manifests. [docker/engine#368](https://github.com/docker/engine/pull/368)

* Fix `POST /images/create` returning a 500 status code when providing an incorrect platform option. [docker/engine#365](https://github.com/docker/engine/pull/365)

* Fix `POST /build` returning a 500 status code when providing an incorrect platform option. [docker/engine#365](https://github.com/docker/engine/pull/365)

* Fix panic on 32-bit ARMv7 caused by misaligned struct member. [docker/engine#363](https://github.com/docker/engine/pull/363)

* Fix to return "invalid parameter" when linking to non-existing container. [docker/engine#352](https://github.com/docker/engine/pull/352)

* Fix overlay2: busy error on mount when using kernel >= 5.2. [docker/engine#332](https://github.com/docker/engine/pull/332)

* Fix `docker rmi` stuck in certain misconfigured systems, e.g. dead NFS share. [docker/engine#335](https://github.com/docker/engine/pull/335)

* Fix handling of blocked I/O of exec'd processes. [docker/engine#296](https://github.com/docker/engine/pull/296)

* Fix jsonfile logger: follow logs stuck when `max-size` is set and `max-file=1`. [docker/engine#378](https://github.com/docker/engine/pull/378)

### Known Issues

#### New

* `DOCKER-USER` iptables chain is missing: [docker/for-linux#810](https://github.com/docker/for-linux/issues/810).
  Users cannot perform additional container network traffic filtering on top of
  this iptables chain. You are not affected by this issue if you are not
  customizing iptable chains on top of `DOCKER-USER`.
     - **Workaround:** Insert the iptables chain after the docker daemon starts.
       For example:
       ```
       iptables -N DOCKER-USER
       iptables -I FORWARD -j DOCKER-USER
       iptables -A DOCKER-USER -j RETURN
       ```

#### Existing

* In some circumstances with large clusters, docker information might, as part of the Swarm section,
  include the error `code = ResourceExhausted desc = grpc: received message larger than
  max (5351376 vs. 4194304)`. This does not indicate any failure or misconfiguration by the user,
  and requires no response.
* Orchestrator port conflict can occur when redeploying all services as new. Due to many swarm manager
  requests in a short amount of time, some services are not able to receive traffic and are causing a `404`
  error after being deployed.
     - **Workaround:** restart all tasks via `docker service update --force`.
* [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack with directory traversal. Workaround until proper fix is available in upcoming patch release: `docker pause` container before doing file operations. [moby/moby#39252](https://github.com/moby/moby/pull/39252)
* `docker cp` regression due to CVE mitigation. An error is produced when the source of `docker cp` is set to `/`.
* Install Docker Engine - Enterprise fails to install on RHEL on Azure. This affects any RHEL version that uses an Extended Update Support (EUS) image. At the time of this writing, known versions affected are RHEL 7.4, 7.5, and 7.6.

     - **Workaround options:**
         - Use an older image and don't get updates. Examples of EUS images are here: https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#rhel-images-with-eus.
         - Import your own RHEL images into Azure and do not rely on the Extended Update Support (EUS) RHEL images.
         - Use a RHEL image that does not contain a minor version in the SKU. These are not attached to EUS repositories. Some examples of those are the first three images (SKUs: 7-RAW, 7-LVM, 7-RAW-CI) listed here : https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#list-of-rhel-images-available.

## 19.03.2
2019-09-03

### Builder

* Fix `COPY --from` to non-existing directory on Windows. [moby/moby#39695](https://github.com/moby/moby/pull/39695)

* Fix builder-next: metadata commands not having created time in history. [moby/moby#39456](https://github.com/moby/moby/issues/39456)

* Fix builder-next: close progress on layer export error. [moby/moby#39782](https://github.com/moby/moby/pull/39782)

* Update buildkit to 588c73e1e4. [moby/moby#39781](https://github.com/moby/moby/pull/39781)

### Client

* Fix Windows absolute path detection on non-Windows [docker/cli#1990](https://github.com/docker/cli/pull/1990)

* Fix to zsh completion script for `docker login --username`.

* Fix context: produce consistent output on `context create`. [docker/cli#1985](https://github.com/docker/cli/pull/1874)

* Fix support for HTTP proxy env variable. [docker/cli#2059](https://github.com/docker/cli/pull/2059)

### Logging

* Fix for reading journald logs. [moby/moby#37819](https://github.com/moby/moby/pull/37819) [moby/moby#38859](http://github.com/moby/moby/pull/38859)

### Networking

* Prevent panic on network attached to a container with disabled networking. [moby/moby#39589](https://github.com/moby/moby/pull/39589)

### Runtime

* Bump Golang to 1.12.8.

* Fix a potential engine panic when using XFS disk quota for containers. [moby/moby#39644](https://github.com/moby/moby/pull/39644)

### Swarm

* Fix an issue where nodes with several tasks could not be removed. [docker/swarmkit#2867](https://github.com/docker/swarmkit/pull/2867)

### Known issues

* In some circumstances with large clusters, docker information might, as part of the Swarm section,
  include the error `code = ResourceExhausted desc = grpc: received message larger than
  max (5351376 vs. 4194304)`. This does not indicate any failure or misconfiguration by the user,
  and requires no response.
* Orchestrator port conflict can occur when redeploying all services as new. Due to many swarm manager
  requests in a short amount of time, some services are not able to receive traffic and are causing a `404`
  error after being deployed.
     - Workaround: restart all tasks via `docker service update --force`.

* Traffic cannot egress the HOST because of missing Iptables rules in the FORWARD chain
  The missing rules are :
     ```
     /sbin/iptables --wait -C FORWARD -o docker_gwbridge -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
     /sbin/iptables --wait -C FORWARD -o docker0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
     ```
     - Workaround: Add these rules back using a script and cron definitions. The script
     must contain '-C' commands to check for the presence of a rule and '-A' commands to add
     rules back. Run the script on a cron in regular intervals, for example, every <x> minutes.
     - Affected versions: 17.06.2-ee-16, 18.09.1, 19.03.0
* [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack with directory traversal. Workaround until proper fix is available in upcoming patch release: `docker pause` container before doing file operations. [moby/moby#39252](https://github.com/moby/moby/pull/39252)
* `docker cp` regression due to CVE mitigation. An error is produced when the source of `docker cp` is set to `/`.
* Install Docker Engine - Enterprise fails to install on RHEL on Azure. This affects any RHEL version that uses an Extended Update Support (EUS) image. At the time of this writing, known versions affected are RHEL 7.4, 7.5, and 7.6.

     - Workaround options:
         - Use an older image and don't get updates. Examples of EUS images are here: https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#rhel-images-with-eus.
         - Import your own RHEL images into Azure and do not rely on the Extended Update Support (EUS) RHEL images.
         - Use a RHEL image that does not contain a minor version in the SKU. These are not attached to EUS repositories. Some examples of those are the first three images (SKUs: 7-RAW, 7-LVM, 7-RAW-CI) listed here : https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#list-of-rhel-images-available.

## 19.03.1
2019-07-25

### Security

 * Fixed loading of nsswitch based config inside chroot under Glibc. [CVE-2019-14271](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-14271)

### Known issues

 * In some circumstances, in large clusters, docker information might, as part of the Swarm section,
 include the error `code = ResourceExhausted desc = grpc: received message larger than
 max (5351376 vs. 4194304)`. This does not indicate any failure or misconfiguration by the user,
 and requires no response.
 * Orchestrator port conflict can occur when redeploying all services as new. Due to many swarm manager
 requests in a short amount of time, some services are not able to receive traffic and are causing a `404`
 error after being deployed.
    - Workaround: restart all tasks via `docker service update --force`.

 * Traffic cannot egress the HOST because of missing Iptables rules in the FORWARD chain
 The missing rules are :
     ```
     /sbin/iptables --wait -C FORWARD -o docker_gwbridge -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
     /sbin/iptables --wait -C FORWARD -o docker0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
     ```
     - Workaround: Add these rules back using a script and cron definitions. The script
     must contain '-C' commands to check for the presence of a rule and '-A' commands to add
     rules back. Run the script on a cron in regular intervals, for example, every <x> minutes.
     - Affected versions: 17.06.2-ee-16, 18.09.1, 19.03.0
 * [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack with directory traversal. Workaround until proper fix is available in upcoming patch release: `docker pause` container before doing file operations. [moby/moby#39252](https://github.com/moby/moby/pull/39252)
 * `docker cp` regression due to CVE mitigation. An error is produced when the source of `docker cp` is set to `/`.
 * Install Docker Engine - Enterprise fails to install on RHEL on Azure. This affects any RHEL version that uses an Extended Update Support (EUS) image. At the time of this writing, known versions affected are RHEL 7.4, 7.5, and 7.6.

     - Workaround options:
         - Use an older image and don't get updates. Examples of EUS images are here: https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#rhel-images-with-eus.
         - Import your own RHEL images into Azure and do not rely on the Extended Update Support (EUS) RHEL images.
         - Use a RHEL image that does not contain a minor version in the SKU. These are not attached to EUS repositories. Some examples of those are the first three images (SKUs: 7-RAW, 7-LVM, 7-RAW-CI) listed here : https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#list-of-rhel-images-available.

## 19.03.0
2019-07-22

### Builder

* Fixed `COPY --from` to preserve ownership. [moby/moby#38599](https://github.com/moby/moby/pull/38599)
* builder-next:

    - Added inline cache support `--cache-from`. [docker/engine#215](https://github.com/docker/engine/pull/215)
    - Outputs configuration allowed. [moby/moby#38898](https://github.com/moby/moby/pull/38898)
    - Fixed gcr workaround token cache. [docker/engine#212](https://github.com/docker/engine/pull/212)
    - `stopprogress` called on download error. [docker/engine#215](https://github.com/docker/engine/pull/215)
    - Buildkit now uses systemd's `resolv.conf`. [docker/engine#260](https://github.com/docker/engine/pull/260).
    - Setting buildkit outputs now allowed. [docker/cli#1766](https://github.com/docker/cli/pull/1766)
    - Look for Dockerfile specific dockerignore file (for example, Dockerfile.dockerignore) for
        ignored paths. [docker/engine#215](https://github.com/docker/engine/pull/215)
    - Automatically detect if process execution is possible for x86, arm, and arm64 binaries.
        [docker/engine#215](https://github.com/docker/engine/pull/215)
    - Updated buildkit to 1f89ec1. [docker/engine#260](https://github.com/docker/engine/pull/260)
    - Use Dockerfile frontend version `docker/dockerfile:1.1` by default.
        [docker/engine#215](https://github.com/docker/engine/pull/215)
    - No longer rely on an external image for COPY/ADD operations.
        [docker/engine#215](https://github.com/docker/engine/pull/215)

### Client

* Added `--pids-limit` flag to `docker update`. [docker/cli#1765](https://github.com/docker/cli/pull/1765)
* Added systctl support for services. [docker/cli#1754](https://github.com/docker/cli/pull/1754)
* Added support for `template_driver` in compose files. [docker/cli#1746](https://github.com/docker/cli/pull/1746)
* Added `--device` support for Windows. [docker/cli#1606](https://github.com/docker/cli/pull/1606)
* Added support for Data Path Port configuration. [docker/cli#1509](https://github.com/docker/cli/pull/1509)
* Added fast context switch: commands. [docker/cli#1501](https://github.com/docker/cli/pull/1501)
* Support added for `--mount type=bind,bind-nonrecursive,...` [docker/cli#1430](https://github.com/docker/cli/pull/1430)
* Added maximum replicas per node. [docker/cli#1612](https://github.com/docker/cli/pull/1612)
* Added option to pull images quietly. [docker/cli#882](https://github.com/docker/cli/pull/882)
* Added a separate `--domainname` flag. [docker/cli#1130](https://github.com/docker/cli/pull/1130)
* Added support for secret drivers in `docker stack deploy`. [docker/cli#1783](https://github.com/docker/cli/pull/1783)
* Added ability to use swarm `Configs` as `CredentialSpecs` on services.
[docker/cli#1781](https://github.com/docker/cli/pull/1781)
* Added `--security-opt systempaths=unconfined` support. [docker/cli#1808](https://github.com/docker/cli/pull/1808)
* Added basic framework for writing and running CLI plugins. [docker/cli#1564](https://github.com/docker/cli/pull/1564)
  [docker/cli#1898](https://github.com/docker/cli/pull/1898)
* Bumped Docker App to v0.8.0. [docker/docker-ce-packaging#341](https://github.com/docker/docker-ce-packaging/pull/341)
* Added support for Docker buildx. [docker/docker-ce-packaging#336](https://github.com/docker/docker-ce-packaging/pull/336)
* Added support for Docker Assemble v0.36.0.
* Added support for Docker Cluster v1.0.0-rc2.
* Added support for Docker Template v0.1.4.
* Added support for Docker Registry v0.1.0-rc1.
* Bumped google.golang.org/grpc to v1.20.1. [docker/cli#1884](https://github.com/docker/cli/pull/1884)
* CLI changed to pass driver specific options to `docker run`. [docker/cli#1767](https://github.com/docker/cli/pull/1767)
* Bumped Golang 1.12.5. [docker/cli#1875](https://github.com/docker/cli/pull/1875)
* `docker system info` output now segregates information relevant to the client and daemon.
[docker/cli#1638](https://github.com/docker/cli/pull/1638)
* (Experimental) When targeting Kubernetes, added support for `x-pull-secret: some-pull-secret` in
compose-files service configs. [docker/cli#1617](https://github.com/docker/cli/pull/1617)
* (Experimental) When targeting Kubernetes, added support for `x-pull-policy: <Never|Always|IfNotPresent>`
in compose-files service configs. [docker/cli#1617](https://github.com/docker/cli/pull/1617)
* cp, save, export: Now preventing overwriting irregular files. [docker/cli#1515](https://github.com/docker/cli/pull/1515)
* npipe volume type on stack file now allowed. [docker/cli#1195](https://github.com/docker/cli/pull/1195)
* Fixed tty initial size error. [docker/cli#1529](https://github.com/docker/cli/pull/1529)
* Fixed problem with labels copying value from environment variables.
[docker/cli#1671](https://github.com/docker/cli/pull/1671)

### API

* Updated API version to v1.40. [moby/moby#38089](https://github.com/moby/moby/pull/38089)
* Added warnings to `/info` endpoint, and moved detection to the daemon.
[moby/moby#37502](https://github.com/moby/moby/pull/37502)
* Added HEAD support for `/_ping` endpoint. [moby/moby#38570](https://github.com/moby/moby/pull/38570)
* Added `Cache-Control` headers to disable caching `/_ping` endpoint.
[moby/moby#38569](https://github.com/moby/moby/pull/38569)
* Added `containerd`, `runc`, and `docker-init` versions to `/version`.
[moby/moby#37974](https://github.com/moby/moby/pull/37974)
* Added undocumented `/grpc` endpoint and registered BuildKit's controller.
[moby/moby#38990](https://github.com/moby/moby/pull/38990)

### Experimental
* Enabled checkpoint/restore of containers with TTY. [moby/moby#38405](https://github.com/moby/moby/pull/38405)
* LCOW: Added support for memory and CPU limits. [moby/moby#37296](https://github.com/moby/moby/pull/37296)
* Windows: Added ContainerD runtime. [moby/moby#38541](https://github.com/moby/moby/pull/38541)
* Windows: LCOW now requires Windows RS5+. [moby/moby#39108](https://github.com/moby/moby/pull/39108)

### Security

* mount: added BindOptions.NonRecursive (API v1.40). [moby/moby#38003](https://github.com/moby/moby/pull/38003)
* seccomp: whitelisted `io_pgetevents()`. [moby/moby#38895](https://github.com/moby/moby/pull/38895)
* seccomp: `ptrace(2)` for 4.8+ kernels now allowed. [moby/moby#38137](https://github.com/moby/moby/pull/38137)

### Runtime

* Running `dockerd` as a non-root user (Rootless mode) is now allowed.
[moby/moby#380050](https://github.com/moby/moby/pull/38050)
* Rootless: optional support provided for `lxc-user-nic` SUID binary.
[docker/engine#208](https://github.com/docker/engine/pull/208)
* Added DeviceRequests to HostConfig to support NVIDIA GPUs. [moby/moby#38828](https://github.com/moby/moby/pull/38828)
* Added `--device` support for Windows. [moby/moby#37638](https://github.com/moby/moby/pull/37638)
* Added `memory.kernelTCP` support for linux. [moby/moby#37043](https://github.com/moby/moby/pull/37043)
* Windows credential specs can now be passed directly to the engine.
[moby/moby#38777](https://github.com/moby/moby/pull/38777)
* Added pids-limit support in docker update. [moby/moby#32519](https://github.com/moby/moby/pull/32519)
* Added support for exact list of capabilities. [moby/moby#38380](https://github.com/moby/moby/pull/38380)
* daemon: Now use 'private' ipc mode by default. [moby/moby#35621](https://github.com/moby/moby/pull/35621)
* daemon: switched to semaphore-gated WaitGroup for startup tasks. [moby/moby#38301](https://github.com/moby/moby/pull/38301)
* Now use `idtools.LookupGroup` instead of parsing `/etc/group` file for docker.sock ownership to
fix: `api.go doesn't respect nsswitch.conf`. [moby/moby#38126](https://github.com/moby/moby/pull/38126)
* cli: fixed images filter when using multi reference filter. [moby/moby#38171](https://github.com/moby/moby/pull/38171)
* Bumped Golang to 1.12.5. [docker/engine#209](https://github.com/docker/engine/pull/209)
* Bumped `containerd` to 1.2.6. [moby/moby#39016](https://github.com/moby/moby/pull/39016)
* Bumped `runc` to 1.0.0-rc8, opencontainers/selinux v1.2.2. [docker/engine#210](https://github.com/docker/engine/pull/210)
* Bumped `google.golang.org/grpc` to v1.20.1. [docker/engine#215](https://github.com/docker/engine/pull/215)
* Performance optimized in aufs and layer store for massively parallel container creation/removal.
[moby/moby#39135](https://github.com/moby/moby/pull/39135) [moby/moby#39209](https://github.com/moby/moby/pull/39209)
* Root is now passed to chroot for chroot Tar/Untar (CVE-2018-15664)
[moby/moby#39292](https://github.com/moby/moby/pull/39292)
* Fixed `docker --init` with /dev bind mount. [moby/moby#37665](https://github.com/moby/moby/pull/37665)
* The right device number is now fetched when greater than 255 and using the `--device-read-bps` option.
[moby/moby#39212](https://github.com/moby/moby/pull/39212)
* Fixed `Path does not exist` error when path definitely exists. [moby/moby#39251](https://github.com/moby/moby/pull/39251)

### Networking

* Moved IPVLAN driver out of experimental.
[moby/moby#38983](https://github.com/moby/moby/pull/38983)
* Added support for 'dangling' filter. [moby/moby#31551](https://github.com/moby/moby/pull/31551)
[docker/libnetwork#2230](https://github.com/docker/libnetwork/pull/2230)
* Load balancer sandbox is now deleted when a service is updated with `--network-rm`.
[docker/engine#213](https://github.com/docker/engine/pull/213)
* Windows: Now forcing a nil IP specified in `PortBindings` to IPv4zero (0.0.0.0).
[docker/libnetwork#2376](https://github.com/docker/libnetwork/pull/2376)

### Swarm

* Added support for maximum replicas per node. [moby/moby#37940](https://github.com/moby/moby/pull/37940)
* Added support for GMSA CredentialSpecs from Swarmkit configs. [moby/moby#38632](https://github.com/moby/moby/pull/38632)
* Added support for sysctl options in services. [moby/moby#37701](https://github.com/moby/moby/pull/37701)
* Added support for filtering on node labels. [moby/moby#37650](https://github.com/moby/moby/pull/37650)
* Windows: Support added for named pipe mounts in docker service create + stack yml.
[moby/moby#37400](https://github.com/moby/moby/pull/37400)
* VXLAN UDP Port configuration now supported. [moby/moby#38102](https://github.com/moby/moby/pull/38102)
* Now using Service Placement Constraints in Enforcer. [docker/swarmkit#2857](https://github.com/docker/swarmkit/pull/2857)
* Increased max recv gRPC message size for nodes and secrets.
[docker/engine#256](https://github.com/docker/engine/pull/256)

### Logging

* Enabled gcplogs driver on Windows. [moby/moby#37717](https://github.com/moby/moby/pull/37717)
* Added zero padding for RFC5424 syslog format. [moby/moby#38335](https://github.com/moby/moby/pull/38335)
* Added `IMAGE_NAME` attribute to `journald` log events. [moby/moby#38032](https://github.com/moby/moby/pull/38032)

### Deprecation

* Deprecate image manifest v2 schema1 in favor of v2 schema2. Future version of Docker will remove
support for v2 schema1 althogether. [moby/moby#39365](https://github.com/moby/moby/pull/39365)
* Removed v1.10 migrator. [moby/moby#38265](https://github.com/moby/moby/pull/38265)
* Now skipping deprecated storage-drivers in auto-selection. [moby/moby#38019](https://github.com/moby/moby/pull/38019)
* Deprecated `aufs` storage driver and added warning. [moby/moby#38090](https://github.com/moby/moby/pull/38090)
* Removed support for 17.09.
* SLES12 is deprecated from Docker Enterprise 3.0, and EOL of SLES12 as an operating system will occur
in Docker Enterprise 3.1. Upgrade to SLES15 for continued support on Docker Enterprise.
* Windows 2016 is formally deprecated from Docker Enterprise 3.0. Only non-overlay networks are supported
on Windows 2016 in Docker Enterprise 3.0. EOL of Windows Server 2016 support will occur in Docker
Enterprise 3.1. Upgrade to Windows Server 2019 for continued support on Docker Enterprise.

For more information on deprecated flags and APIs, refer to
https://docs.docker.com/engine/deprecated/ for target removal dates.

### Known issues

* In some circumstances with large clusters, docker information might, as part of the Swarm section,
include the error `code = ResourceExhausted desc = grpc: received message larger than
max (5351376 vs. 4194304)`. This does not indicate any failure or misconfiguration by the user,
and requires no response.
* Orchestrator port conflict can occur when redeploying all services as new. Due to many swarm manager
requests in a short amount of time, some services are not able to receive traffic and are causing a `404`
error after being deployed.
   - Workaround: restart all tasks via `docker service update --force`.

* Traffic cannot egress the HOST because of missing Iptables rules in the FORWARD chain
The missing rules are :
    ```
    /sbin/iptables --wait -C FORWARD -o docker_gwbridge -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
    /sbin/iptables --wait -C FORWARD -o docker0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
    ```
    - Workaround: Add these rules back using a script and cron definitions. The script
    must contain '-C' commands to check for the presence of a rule and '-A' commands to add
    rules back. Run the script on a cron in regular intervals, for example, every <x> minutes.
    - Affected versions: 17.06.2-ee-16, 18.09.1, 19.03.0
* [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack with directory traversal. Workaround until proper fix is available in upcoming patch release: `docker pause` container before doing file operations. [moby/moby#39252](https://github.com/moby/moby/pull/39252)
* `docker cp` regression due to CVE mitigation. An error is produced when the source of `docker cp` is set to `/`.
* Install Docker Engine - Enterprise fails to install on RHEL on Azure. This affects any RHEL version that uses an Extended Update Support (EUS) image. At the time of this writing, known versions affected are RHEL 7.4, 7.5, and 7.6.

    - Workaround options:
        - Use an older image and don't get updates. Examples of EUS images are here: https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#rhel-images-with-eus.
        - Import your own RHEL images into Azure and do not rely on the Extended Update Support (EUS) RHEL images.
        - Use a RHEL image that does not contain a minor version in the SKU. These are not attached to EUS repositories. Some examples of those are the first three images (SKUs: 7-RAW, 7-LVM, 7-RAW-CI) listed here : https://docs.microsoft.com/en-us/azure/virtual-machines/linux/rhel-images#list-of-rhel-images-available.

# Older Docker Engine EE Release notes

## 18.03.1-ee-12
2019-11-14

### Client

* Fix potential out of memory in CLI when running `docker image prune`. [docker/cli#1423](https://github.com/docker/cli/pull/1423)

### Logging

* Fix jsonfile logger: follow logs stuck when `max-size` is set and `max-file=1`. [moby/moby#39969](https://github.com/moby/moby/pull/39969)

### Runtime

* Update to Go 1.12.12.
* Seccomp: add sigprocmask (used by x86 glibc) to default seccomp profile. [moby/moby#39824](https://github.com/moby/moby/pull/39824)

## 18.03.1-ee-11

2019-09-03

### Runtime

* Fix [CVE-2019-14271](https://nvd.nist.gov/vuln/detail/CVE-2019-14271) loading of nsswitch based config inside chroot under Glibc.

* Fix a potential engine panic when using XFS disk quota for containers. [moby/moby#39644](https://github.com/mony/moby/pull/39644)

* Fix overlay2 storage driver getting "device or resource busy" on mount. [moby/moby#37993](https://github.com/moby/moby/pull/37993)

* Update to Go 1.11.13.

### Logging

* Fix for reading journald logs. [moby/moby#37819](https://github.com/moby/moby/pull/37819) [moby/moby#38859](https://github.com/moby/moby/pull/38859)

### Networking

* Fix cluster connectivity issue caused by high qLen in networkdb. [docker/libnetwork#2216](https://github.com/docker/libnetwork/pull/2216)

* Fix possible nil pointer exception. [docker/libnetwork#2325](https://github.com/docker/libnetwork/pull/2325)

* Fix service port for an application becomes unavailable randomly. [docker/libnetwork#2069](https://github.com/docker/libnetwork/pull/2069)

### Swarm

* Fix swarm overlay networking not working after `--force-new-cluster`. [docker/libnetwork#2307](https://github.com/docker/libnetwork/pull/2307)

## 18.03.1-ee-10

2019-07-17

### Runtime

* Masked the secrets updated to the log files when running Docker Engine in debug mode. [CVE-2019-13509](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-13509): If a Docker engine is running in debug mode, and `docker stack deploy` is used to redeploy a stack which includes non-external secrets, the logs will contain the secret.

## 18.03.1-ee-9

2019-06-25

### Client

* Fixed annnotation on `docker config create --template-driver`. [docker/cli#1769](https://github.com/docker/cli/pull/1769)
* Fixed annnotation on `docker secret create --template-driver`. [docker/cli#1785](https://github.com/docker/cli/pull/1785)

### Runtime

* Performance optimized in aufs and layer store for massively parallel container creation/removal.
[moby/moby#39107](https://github.com/moby/moby/pull/39107)
* Windows: fixed support for `docker service create --limit-cpu`.
[moby/moby#39190](https://github.com/moby/moby/pull/39190)
* Now using original process spec for execs. [moby/moby#38871](https://github.com/moby/moby/pull/38871)
* Fixed [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack
with directory traversal. [moby/moby#39357](https://github.com/moby/moby/pull/39357)

## 18.03.1-ee-8

2019-03-28

### Builder

* Added validation for `git ref` to avoid misinterpreation as a flag. [moby/moby#38944](https://github.com/moby/moby/pull/38944)

### Runtime

* Fixed `docker cp` error for filenames greater than 100 characters. [moby/moby#38634]
* Fixed `layer/layer_store` to ensure `NewInputTarStream` resources are released. [moby/moby#38413]

### Swarm Mode

* Fixed issue for swarm nodes not being able to join as masters if http proxy is set. [moby/moby#36951]

## 18.03.1-ee-7

2019-02-28

### Runtime

* Updated to Go version 1.10.8.
* Updated to containerd version 1.1.6.
- When copying existing folder, xattr set errors when the target filesystem doesn't support xattr are now ignored. [moby/moby#38316](https://github.com/moby/moby/pull/38316)
- Fixed FIFO, sockets, and device files in userns, and fixed device mode not being detected. [moby/moby#38758](https://github.com/moby/moby/pull/38758)
- Deleted stale containerd object on start failure. [moby/moby#38364](https://github.com/moby/moby/pull/38364)

### Bug fixes

* Fixed an issue to address the IPAM state from networkdb if manager is not attached to the overlay network. (docker/escalation#1049)

## 18.03.1-ee-6
2019-02-11

### Security fixes for Docker Engine - Enterprise
* Update `runc` to address a critical vulnerability that allows specially-crafted containers to gain administrative privileges on the host. [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)
* Ubuntu 14.04 customers using a 3.13 kernel will need to upgrade to a supported Ubuntu 4.x kernel

## 18.03.1-ee-5
2019-01-09

### Security fixes
* Upgraded Go language to 1.10.6 to resolve CVE-2018-16873, CVE-2018-16874, and CVE-2018-16875.
* Added `/proc/asound` to masked paths
* Fixed authz plugin for 0-length content and path validation.

### Fixes for Docker Engine - Enterprise
* Disable kmem accounting in runc on RHEL/CentOS (docker/escalation#614, docker/escalation#692)
* Fix resource leak on `docker logs --follow` [moby/moby#37576](https://github.com/moby/moby/pull/37576)
* Mask proxy credentials from URL when displayed in system info (docker/escalation#879)

## 18.03.1-ee-4
2018-10-25

  > **Note**: If you're deploying UCP or DTR, use Docker EE Engine 18.09 or higher. 18.03 is an engine only release.

### Client
* Fixed help message flags on docker stack commands and child commands. [docker/cli#1251](https://github.com/docker/cli/pull/1251)
* Fixed typo breaking zsh docker update autocomplete. [docker/cli#1232](https://github.com/docker/cli/pull/1232)

### Networking
* Added optimizations to reduce the messages in the NetworkDB queue. [docker/libnetwork#2225](https://github.com/docker/libnetwork/pull/2225)
* Fixed a very rare condition where managers are not correctly triggering the reconnection logic. [docker/libnetwork#2226](https://github.com/docker/libnetwork/pull/2226)
* Changed loglevel from error to warning for missing disable_ipv6 file. [docker/libnetwork#2224](https://github.com/docker/libnetwork/pull/2224)

### Runtime
* Fixed denial of service with large numbers in cpuset-cpus and cpuset-mems. [moby/moby#37967](https://github.com/moby/moby/pull/37967)
* Added stability improvements for devicemapper shutdown. [moby/moby#36307](https://github.com/moby/moby/pull/36307) [moby/moby#36438](https://github.com/moby/moby/pull/36438)

### Swarm Mode
* Fixed the logic used for skipping over running tasks. [docker/swarmkit#2724](https://github.com/docker/swarmkit/pull/2724)
* Addressed unassigned task leak when a service is removed. [docker/swarmkit#2709](https://github.com/docker/swarmkit/pull/2709)

## 18.03.1-ee-3
2018-08-30

### Builder
* Fix: no error if build args are missing during docker build. [docker/engine#25](https://github.com/docker/engine/pull/25)
* Ensure RUN instruction to run without healthcheck. [moby/moby#37413](https://github.com/moby/moby/pull/37413)

### Client
* Fix manifest list to always use correct size. [docker/cli#1156](https://github.com/docker/cli/pull/1156)
* Various shell completion script updates. [docker/cli#1159](https://github.com/docker/cli/pull/1159) [docker/cli#1227](https://github.com/docker/cli/pull/1227)
* Improve version output alignment. [docker/cli#1204](https://github.com/docker/cli/pull/1204)

### Runtime
* Disable CRI plugin listening on port 10010 by default. [docker/engine#29](https://github.com/docker/engine/pull/29)
* Update containerd to v1.1.2. [docker/engine#33](https://github.com/docker/engine/pull/33)
* Windows: Pass back system errors on container exit. [moby/moby#35967](https://github.com/moby/moby/pull/35967)
* Windows: Fix named pipe support for hyper-v isolated containers. [docker/engine#2](https://github.com/docker/engine/pull/2) [docker/cli#1165](https://github.com/docker/cli/pull/1165)
* Register OCI media types. [docker/engine#4](https://github.com/docker/engine/pull/4)

### Swarm Mode
* Clean up tasks in dirty list for which the service has been deleted. [docker/swarmkit#2694](https://github.com/docker/swarmkit/pull/2694)
* Propagate the provided external CA certificate to the external CA object in swarm. [docker/cli#1178](https://github.com/docker/cli/pull/1178)

## 18.03.1-ee-2
2018-07-10

> ### Important notes about this release
>
> If you're deploying UCP or DTR, use Docker Engine EE `17.06` or `18.09`. See [Docker Compatibility Matrix](https://success.docker.com/article/compatibility-matrix) for more information.
{: .important}

#### Runtime

+ Add /proc/acpi to masked paths [(CVE-2018-10892)](https://cve.mitre.org/cgi-bin/cvename.cgi?name=2018-10892). [moby/moby#37404](https://github.com/moby/moby/pull/37404)

## 18.03.1-ee-1
2018-06-27

> ### Important notes about this release
>
> If you're deploying UCP or DTR, use Docker Engine EE `17.06` or `18.09`. See [Docker Compatibility Matrix](https://success.docker.com/article/compatibility-matrix) for more information.
{: .important}

### Client

+ Update to docker-ce 18.03.1 client.
+ Add `docker trust` command for image signing and enabling the secure supply chain from development to deployment.
+ Add docker compose on Kubernetes.

### Runtime

+ Update to docker-ce 18.03.1 engine.
+ Add support for FIPS 140-2 on x86_64.
+ Add support for Microsoft Windows Server 1709 and 1803 with support for [swarm ingress routing mesh](https://docs.docker.com/engine/swarm/ingress/), [VIP service discovery](https://docs.docker.com/v17.09/engine/swarm/networking/#configure-service-discovery), and [named pipe mounting](https://blog.docker.com/2017/09/docker-windows-server-1709/).
+ Add support for Ubuntu 18.04.
+ Windows opt-out telemetry stream.
+ Support for `--chown` with `COPY` and `ADD` in `Dockerfile`.
+ Added functionality for the `docker logs` command to include the output of multiple logging drivers.

## 17.06.2-ee-25
2019-11-19

### Builder

* Fix for ENV in multi-stage builds not being isolated. [moby/moby#35456](https://github.com/moby/moby/pull/35456)

### Client

* Fix potential out of memory in CLI when running `docker image prune`. [docker/cli#1423](https://github.com/docker/cli/pull/1423)
* Fix compose file schema to prevent invalid properties in `deploy.resources`. [docker/cli#455](https://github.com/docker/cli/pull/455)

### Logging

* Fix jsonfile logger: follow logs stuck when `max-size` is set and `max-file=1`. [moby/moby#39969](https://github.com/moby/moby/pull/39969)

### Runtime

* Update to Go 1.12.12.
* Seccomp: add sigprocmask (used by x86 glibc) to default seccomp profile. [moby/moby#39824](https://github.com/moby/moby/pull/39824)
* Fix "device or resource busy" error on container removal with devicemapper. [moby/moby#34573](https://github.com/moby/moby/pull/34573)
* Fix `daemon.json` configuration `default-ulimits` not working. [moby/moby#32547](https://github.com/moby/moby/pull/32547)
* Fix denial of service with large numbers in `--cpuset-cpus` and `--cpuset-mems`. [moby/moby#37967](https://github.com/moby/moby/pull/37967)
* Fix for `docker start` creates host-directory for bind mount, but shouldn't. [moby/moby#35833](https://github.com/moby/moby/pull/35833)
* Fix OCI image media types. [moby/moby#37359](https://github.com/moby/moby/pull/37359)

### Windows

* Windows: bump RW layer size to 127GB. [moby/moby#35925](https://github.com/moby/moby/pull/35925)

## 17.06.2-ee-24
2019-09-03

### Runtime

* Fix [CVE-2019-14271](https://nvd.nist.gov/vuln/detail/CVE-2019-14271) loading of nsswitch based config inside chroot under Glibc.
* Fix Fix a potential engine panic when using XFS disk quota for containers. [moby/moby#39644](https://github.com/moby/moby/pull/39644)
* Update to Go 1.11.13.

### Logging

* Fix for reading journald logs. [moby/moby#37819](https://github.com/moby/moby/pull/37819) [moby/moby#38859](https://github.com/moby/moby/pull/38859)

### Networking

* Fix cluster connectivity issue caused by high qLen in networkdb. [docker/libnetwork#2216](https://github.com/docker/libnetwork/pull/2216)
* Fix service port for an application becomes unavailable randomly. [docker/libnetwork#2069](docker/libnetwork#2069)

## 17.06.2-ee-23
2019-07-17

### Runtime

* Masked the secrets updated to the log files when running Docker Engine in debug mode. [CVE-2019-13509](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-13509): If a Docker engine is running in debug mode, and `docker stack deploy` is used to redeploy a stack which includes non-external secrets, the logs will contain the secret.

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-22
2019-06-27

### Networking

* Fixed a bug where if a service has the same number of host-mode published ports with PublishedPort 0, changes to the spec is not reflected in the service object. [docker/swarmkit#2376](https://github.com/docker/swarmkit/pull/2376)

### Runtime

* Added performance optimizations in aufs and layer store that helps in the creation and removal of massively parallel containers. [moby/moby#39107](https://github.com/moby/moby/pull/39107)
* Fixed [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack with directory traversal. [moby/moby#39357](https://github.com/moby/moby/pull/39357)
* Windows: fixed support for docker service `create --limit-cpu`. [moby/moby#39190](https://github.com/moby/moby/pull/39190)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-21
2019-04-11

### Builder

* Added validation for git ref so it can't be misinterpreted as a flag. [moby/moby#38944](https://github.com/moby/moby/pull/38944)

### Runtime

* Fixed `docker cp` error with filenames greater than 100 characters. [moby/moby#38634](https://github.com/moby/moby/pull/38634)
* Removed temporary hot-fix and applied latest upstream patches for CVE-2019-5736. [docker/runc#9](https://github.com/docker/runc/pull/9)
* Fixed rootfs: umount all procfs and sysfs with `--no-pivot`. [docker/runc#10](https://github.com/docker/runc/pull/10)

## 17.06.2-ee-20
2019-02-28

### Bug fixes
* Fixed an issue to address the IPAM state from networkdb if manager is not attached to the overlay network. (docker/escalation#1049)

### Runtime

* Updated to Go version 1.10.8.
+ Added cgroup namespace support. [docker/runc#7](https://github.com/docker/runc/pull/7)

### Windows

* Fixed `failed to register layer` bug on `docker pull` of windows images.

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-19

2019-02-11

### Security fixes for Docker Engine - Enterprise
* Update `runc` to address a critical vulnerability that allows specially-crafted containers to gain administrative privileges on the host. [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)
* Ubuntu 14.04 customers using a 3.13 kernel will need to upgrade to a supported Ubuntu 4.x kernel

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-18
2019-01-09

### Security fixes
* Upgraded Go language to 1.10.6 to resolve CVE-2018-16873, CVE-2018-16874, and CVE-2018-16875.
* Added `/proc/asound` to masked paths
* Fixed authz plugin for 0-length content and path validation.

### Fixes for Docker Engine Engine - Enterprise
* Disable kmem accounting in runc on RHEL/CentOS (docker/escalation#614, docker/escalation#692)
* Fix resource leak on `docker logs --follow` [moby/moby#37576](https://github.com/moby/moby/pull/37576)
* Mask proxy credentials from URL when displayed in system info (docker/escalation#879)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-17
2018-10-25

### Networking

* Changed loglevel from error to warning for missing disable_ipv6 file. [docker/libnetwork#2223](https://github.com/docker/libnetwork/pull/2223)
* Fixed subnet allocation to avoid reallocating recently freed subnets. [docker/libnetwork#2255](https://github.com/docker/libnetwork/pull/2255)
* Fixed libnetwork issue which caused errors to be returned when iptables or firewalld issues transient warnings. [docker/libnetwork#2218](https://github.com/docker/libnetwork/pull/2218)

### Plugins

* Fixed too many "Plugin not found" error messages. [moby/moby#36119](https://github.com/moby/moby/pull/36119)

### Swarm mode

* Added failed allocations retry immediately upon a deallocation to overcome IP exhaustion. [docker/swarmkit#2711](https://github.com/docker/swarmkit/pull/2711)
* Fixed leaking task resources. [docker/swarmkit#2755](https://github.com/docker/swarmkit/pull/2755)
* Fixed deadlock in dispatcher that could cause node to crash. [docker/swarmkit#2753](https://github.com/docker/swarmkit/pull/2753)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-16
2018-07-26

### Client

- Fix service rollback options being cross-wired. [docker/cli#1052](https://github.com/docker/cli/pull/1052)

### Networking

* Protect against possible race on ingress programming. [docker/libnetwork#2195](https://github.com/docker/libnetwork/pull/2195)
* Add a recovery mechanism for a split gossip cluster. [docker/libnetwork#2169](https://github.com/docker/libnetwork/pull/2169)

### Packaging

* Update packaging description and license to Docker EUSA.

### Runtime

* Update overlay2 to use naive diff for changes. [moby/moby#37313](https://github.com/moby/moby/pull/37313)

### Swarm mode

- Fix task reaper batching. [docker/swarmkit#2678](https://github.com/docker/swarmkit/pull/2678)
* RoleManager will remove deleted nodes from the cluster membership. [docker/swarmkit#2607](https://github.com/docker/swarmkit/pull/2607)
- Fix unassigned task leak when service is removed. [docker/swarmkit#2708](https://github.com/docker/swarmkit/pull/2708)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-15
2018-07-10

### Runtime

- Add /proc/acpi to masked paths [(CVE-2018-10892)](https://cve.mitre.org/cgi-bin/cvename.cgi?name=2018-10892). [moby/moby#37404](https://github.com/moby/moby/pull/37404)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-14
2018-06-21

### Client

* Set a 30s timeout for HTTP client communication with plugin backend. [docker/cli#883](https://github.com/docker/cli/pull/883)
- Fix `docker version` output misaligned. [docker/cli#965](https://github.com/docker/cli/pull/965)

### Runtime

- Fix Windows docker daemon crash when docker stats is used. [moby/moby#35968](https://github.com/moby/moby/pull/35968)
* Add `/proc/keys` to masked paths. [moby/moby#36368](https://github.com/moby/moby/pull/36368)
* Added support for persisting Windows network driver options. [moby/moby#35563](https://github.com/moby/moby/pull/35563)
- Fix to ensure graphdriver dir is a shared mount. [moby/moby#36047](https://github.com/moby/moby/pull/36047)

### Swarm mode

- Fix `docker stack deploy --prune` with empty name removes all swarm services. [moby/moby#36776](https://github.com/moby/moby/issues/36776)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-13
2018-06-04

### Networking

- Fix attachable containers that may leave DNS state when exiting. [docker/libnetwork#2175](https://github.com/docker/libnetwork/pull/2175)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-12
2018-05-29

### Networking

- Fix to allow service update with no connection loss. [docker/libnetwork#2157](https://github.com/docker/libnetwork/pull/2157)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-11
2018-05-17

### Client

- Fix presentation of published "random" host ports. [docker/cli#404](https://github.com/docker/cli/pull/404)

### Networking

* Fix concurrent CreateNetwork in bridge driver. [docker/libnetwork#2127](https://github.com/docker/libnetwork/pull/2127)

### Runtime

* Use rslave propagation for mounts from daemon root. [moby/moby#36055](https://github.com/moby/moby/pull/36055)
* Use rslave instead of rprivate in choortarchive. [moby/moby#35217](https://github.com/moby/moby/pull/35217)
* Set daemon root to use shared propagation. [moby/moby#36096](https://github.com/moby/moby/pull/36096)
* Windows: Increase container default shutdown timeout. [moby/moby#35184](https://github.com/moby/moby/pull/35184)
* Avoid using all system memory with authz plugins. [moby/moby#36595](https://github.com/moby/moby/pull/36595)
* Daemon/stats: more resilient cpu sampling. [moby/moby#36519](https://github.com/moby/moby/pull/36519)

### Known issues

* When all Swarm managers are stopped at the same time, the swarm might end up in a
split-brain scenario. [Learn more](https://success.docker.com/article/KB000759).
* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-10
2018-04-27

### Runtime

* Fix version output to not have `-dev`.

### Known issues

* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-9
2018-04-26

### Runtime

* Make Swarm manager Raft quorum parameters configurable in daemon config. [moby/moby#36726](https://github.com/moby/moby/pull/36726)
* Windows: Ignore missing tombstone files when closing an image.
* Windows: Fix directory deletes when a container sharing a base image is running.

### Swarm mode

- Increase raft ElectionTick to 10xHeartbeatTick. [docker/swarmkit#2564](https://github.com/docker/swarmkit/pull/2564)
- Adding logic to restore networks in order. [docker/swarmkit#2584](https://github.com/docker/swarmkit/pull/2584)

### Known issues

* Under certain conditions, swarm leader re-election may timeout
  prematurely. During this period, docker commands may fail. Also during
  this time, creation of globally-scoped networks may be unstable. As a
  workaround, wait for leader election to complete before issuing commands
  to the cluster.
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-8
2018-04-17

### Runtime

* Update `hcsshim` to v0.6.10 to address [CVE-2018-8115](https://portal.msrc.microsoft.com/en-us/security-guidance/advisory/CVE-2018-8115)

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
* It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
* Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-7
2018-03-19

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

### Known issues

 * It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
 * Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-6
2017-11-27

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

### Known issues

 * It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
 * Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-5
2017-11-02

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

#### Swarm mode

* Increase gRPC request timeout to 20 seconds for sending snapshots to prevent `context deadline exceeded` errors [docker/swarmkit#2391](https://github.com/docker/swarmkit/pull/2391)
* When a node is removed, delete all of its attachment tasks so networks used by those tasks can be removed [docker/swarmkit#2414](https://github.com/docker/swarmkit/pull/2414)

#### Known issues

 * It's recommended that users create overlay networks with `/24` blocks (the default) of 256 IP addresses when networks are used by services created using VIP-based endpoint-mode (the default). This is because of limitations with Docker Swarm [moby/moby#30820](moby/moby/issues/30820). Users should _not_ work around this by increasing the IP block size. To work around this limitation, either use `dnsrr` endpoint-mode or use multiple smaller overlay networks.
 * Docker may experience IP exhaustion if many tasks are assigned to a single overlay network, for example if many services are attached to that network or because services on the network are scaled to many replicas. The problem may also manifest when tasks are rescheduled because of node failures. In case of node failure, Docker currently waits 24h to release overlay IP addresses. The problem can be diagnosed by looking for `failed to allocate network IP for task` messages in the Docker logs.
* SELinux enablement is not supported for containers on IBM Z on RHEL because of missing Red Hat package.
* If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-4
2017-10-12

### Client

* Fix idempotence of `docker stack deploy` when secrets or configs are used [docker/cli#509](https://github.com/docker/cli/pull/509)

### Logging

* Avoid using a map for log attributes to prevent panic [moby/moby#34174](https://github.com/moby/moby/pull/34174)

### Networking

* Fix for garbage collection logic in NetworkDB. Entries were not properly garbage collected and deleted within the expected time [docker/libnetwork#1944](https://github.com/docker/libnetwork/pull/1944) [docker/libnetwork#1960](https://github.com/docker/libnetwork/pull/1960)
* Allow configuration of max packet size in network DB to use the full available MTU. This requires a configuration in the docker daemon and need a dockerd restart [docker/libnetwork#1839](https://github.com/docker/libnetwork/pull/1839)
* Overlay fix for transient IP reuse [docker/libnetwork#1935](https://github.com/docker/libnetwork/pull/1935) [docker/libnetwork#1968](https://github.com/docker/libnetwork/pull/1968)
* Serialize IP allocation [docker/libnetwork#1788](https://github.com/docker/libnetwork/pull/1788)

### Known issues

If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.2-ee-3
2017-09-22

### Swarm mode

- Increase max message size to allow larger snapshots [docker/swarmkit#131](https://github.com/docker/swarmkit/pull/131)

### Known issues

If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.1-ee-2
2017-08-24

### Client

- Enable TCP Keep-Alive in Docker client [#415](https://github.com/docker/cli/pull/415)

### Networking

- Lock goroutine to OS thread while changing NS [#1911](https://github.com/docker/libnetwork/pull/1911)

### Runtime

- devmapper: ensure that UdevWait is called after calls to setCookie [#33732](https://github.com/moby/moby/pull/33732)
- aufs: ensure diff layers are correctly removed to prevent leftover files from using up storage [#34587](https://github.com/moby/moby/pull/34587)

### Swarm mode

- Ignore PullOptions for running tasks [#2351](https://github.com/docker/swarmkit/pull/2351)

### Known issues

If a container is spawned on node A, using the same IP of a container destroyed
on nodeB within 5 min from the time that it exit, the container on node A is
not reachable until one of these 2 conditions happens:

1. Container on A sends a packet out,
2. The timer that cleans the arp entry in the overlay namespace is triggered (around 5 minutes).

As a workaround, send at least a packet out from each container like
(ping, GARP, etc).

## 17.06.1-ee-1
2017-08-16

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

## Docker EE 17.03.2-ee-8
2017-12-13

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

## Docker EE 17.03.2-ee-7
2017-10-04

* Fix logic in network resource reaping to prevent memory leak [docker/libnetwork#1944](https://github.com/docker/libnetwork/pull/1944) [docker/libnetwork#1960](https://github.com/docker/libnetwork/pull/1960)
* Increase max GRPC message size to 128MB for larger snapshots so newly added managers can successfully join [docker/swarmkit#2375](https://github.com/docker/swarmkit/pull/2375)

## Docker EE 17.03.2-ee-6
2017-08-24

* Fix daemon panic on docker image push [moby/moby#33105](https://github.com/moby/moby/pull/33105)
* Fix panic in concurrent network creation/deletion operations [docker/libnetwork#1861](https://github.com/docker/libnetwork/pull/1861)
* Improve network db stability under stressful situations [docker/libnetwork#1860](https://github.com/docker/libnetwork/pull/1860)
* Enable TCP Keep-Alive in Docker client [docker/cli#415](https://github.com/docker/cli/pull/415)
* Lock goroutine to OS thread while changing NS [docker/libnetwork#1911](https://github.com/docker/libnetwork/pull/1911)
* Ignore PullOptions for running tasks [docker/swarmkit#2351](https://github.com/docker/swarmkit/pull/2351)

## Docker EE 17.03.2-ee-5
20 Jul 2017

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

## Docker EE 17.03.2-ee-4
2017-06-01

Refer to the [detailed list](https://github.com/moby/moby/releases/tag/v17.03.2-ce) of all changes since the release of Docker EE 17.03.1-ee-3

**Note**: This release includes a fix for potential data loss under certain
circumstances with the local (built-in) volume driver.

## Docker EE 17.03.1-ee-3
2017-03-30

* Fix an issue with the SELinux policy for Oracle Linux [#31501](https://github.com/docker/docker/pull/31501)

## Docker EE 17.03.1-ee-2
2017-03-28

* Fix issue with swarm CA timeouts [#2063](https://github.com/docker/swarmkit/pull/2063) [#2064](https://github.com/docker/swarmkit/pull/2064/files)

Refer to the [detailed list](https://github.com/moby/moby/releases/tag/v17.03.1-ce) of all changes since the release of Docker EE 17.03.0-ee-1

## Docker EE 17.03.0-ee-1 (2 Mar 2017)

Initial Docker EE release, based on Docker CE 17.03.0

* Optimize size calculation for `docker system df` container size [#31159](https://github.com/docker/docker/pull/31159)
