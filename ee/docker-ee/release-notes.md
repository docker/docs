---
title: Docker Engine - Enterprise release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Engine - Enterprise
keywords: docker, docker engine, ee, ce, whats new, release notes
toc_min: 1
toc_max: 2
skip_read_time: true
redirect_from:
  - /cs-engine/1.12/release-notes/
  - /cs-engine/1.12/release-notes/release-notes/
  - /cs-engine/1.12/release-notes/prior-release-notes/
  - /cs-engine/1.13/release-notes/
  - /ee/engine/release-notes/
---

>{% include enterprise_label_shortform.md %}

This document describes the latest changes, additions, known issues, and fixes
for Docker Engine - Enterprise.

Docker Engine - Enterprise builds upon the corresponding Docker Engine -
Community that it references. Docker Engine - Enterprise includes enterprise
features as well as back-ported fixes (security-related and priority defects)
from the open source. It also incorporates defect fixes for environments in
which new features cannot be adopted as quickly for consistency and
compatibility reasons.

> **Note**
>
> The client and container runtime are now in separate packages from the daemon
> since Docker Engine 18.09. Users should install and update all three packages at
> the same time to get the latest patch releases. For example, on Ubuntu:
> `sudo apt-get install docker-ee docker-ee-cli containerd.io`. See the install
> instructions for the corresponding linux distro for details.

# Version 19.03

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

# Version 18.09

> **Note**
>
> New in 18.09 is an aligned release model for Docker Engine - Community and
> Docker Engine - Enterprise. The new versioning scheme is YY.MM.x where x is an
> incrementing patch version. The enterprise engine is a superset of the
> community engine. They will ship concurrently with the same x patch version
> based on the same code base.

## 18.09.11
2019-11-14

### Builder

* Fix builder-next: filter type in BuildKit GC config. [docker/engine#409](https://github.com/docker/engine/pull/409)

### Runtime

* Bump Golang to 1.12.12.

### Swarm

* Fix update out of sequence and increase max recv gRPC message size for nodes and secrets. [docker/swarmkit#2900](https://github.com/docker/swarmkit/pull/2900)
* Fix for specifying `--default-addr-pool` for `docker swarm init` not picked up by ingress network. [docker/swarmkit#2892](https://github.com/docker/swarmkit/pull/2892)

## 18.09.10
2019-10-08

### Client

* Fix client version not being pinned when set. [docker/engine#118](https://github.com/docker/engine/pull/188)
* Improve error message shown on Windows when daemon is not running or client does not have elevated permissions. [docker/engine#343](https://github.com/docker/engine/pull/343)
* Mitigate against YAML files that have excessive aliasing. [docker/cli#2119](https://github.com/docker/cli/pull/2119)

### Runtime

* Send exec exit event even if the exec fails to find the binary. [docker/engine#357](https://github.com/docker/engine/pull/357)
* Devicemapper: use correct API to get the free loop device index. [docker/engine#348](https://github.com/docker/engine/pull/348)
* Fix overlay2 busy error on mount using kernels >=5.2. [docker/engine#333](https://github.com/docker/engine/pull/333)
* Sleep before attemping to restart event processing. [docker/engine#362](https://github.com/docker/engine/pull/362)
* Seccomp: add sigprocmask (used by x86 glibc) to default profile. [docker/engine#341](https://github.com/docker/engine/pull/341)
* Fix panic on 32-bit ARMv7 caused by misaligned struct member. [docker/engine#364](https://github.com/docker/engine/pull/364)
* Fix `docker rmi` stuck in case of misconfigured system (such as dead NFS share). [docker/engine#336](https://github.com/docker/engine/pull/336)
* Fix jsonfile logger: follow logs stuck when `max-size` is set and `max-file=1`. [docker/engine#377](https://github.com/docker/engine/pull/377)

## 18.09.9
2019-09-03

### Client

* Fix Windows absolute path detection on non-Windows. [docker/cli#1990](https://github.com/docker/cli/pull/1990)
* Fix Docker refusing to load key from delegation.key on Windows. [docker/cli#1968](https://github.com/docker/cli/pull/1968)
* Completion scripts updates for bash and zsh.

### Logging

* Fix for reading journald logs. [moby/moby#37819](https://github.com/moby/moby/pull/37819) [moby/moby#38859](https://github.com/moby/moby/pull/38859)

### Networking

* Prevent panic on network attached to a container with disabled networking. [moby/moby#39589](https://github.com/moby/moby/pull/39589)
* Fix service port for an application becomes unavailable randomly. [docker/libnetwork#2069](https://github.com/docker/libnetwork/pull/2069)
* Fix cleaning up `--config-only` networks `--config-from` networkshave ungracefully exited. [docker/libnetwork#2373](https://github.com/docker/libnetwork/pull/2373)

### Runtime

* Update to Go 1.11.13.
* Fix a potential engine panic when using XFS disk quota for containers. [moby/moby#39644](https://github.com/moby/moby/pull/39644)

### Swarm

* Fix "grpc: received message larger than max" errors. [moby/moby#39306](https://github.com/moby/moby/pull/39306)
* Fix an issue where nodes several tasks could not be removed. [docker/swarmkit#2867](https://github.com/docker/swarmkit/pull/2867)

## 18.09.8
2019-07-17

### Runtime

* Masked the secrets updated to the log files when running Docker Engine in debug mode. [CVE-2019-13509](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-13509): If a Docker engine is running in debug mode, and `docker stack deploy` is used to redeploy a stack which includes non-external secrets, the logs will contain the secret.


### Client

* Fixed rollback config type interpolation for `parallelism` and `max_failure_ratio` fields.

### Known Issue

* There are [important changes](/ee/upgrade) to the upgrade process that, if not correctly followed, can have an impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or later.

## 18.09.7
2019-06-27

### Builder

* Fixed a panic error when building dockerfiles that contain only comments. [moby/moby#38487](https://github.com/moby/moby/pull/38487)
* Added a workaround for GCR authentication issue. [moby/moby#38246](https://github.com/moby/moby/pull/38246)
* Builder-next: Fixed a bug in the GCR token cache implementation workaround. [moby/moby#39183](https://github.com/moby/moby/pull/39183)

### Networking
*  Fixed an error where `--network-rm` would fail to remove a network. [moby/moby#39174](https://github.com/moby/moby/pull/39174)

### Runtime

* Added performance optimizations in aufs and layer store that helps in massively parallel container creation and removal. [moby/moby#39107](https://github.com/moby/moby/pull/39107), [moby/moby#39135](https://github.com/moby/moby/pull/39135)
* Updated containerd to version 1.2.6. [moby/moby#39016](https://github.com/moby/moby/pull/39016)
* Fixed [CVE-2018-15664](https://nvd.nist.gov/vuln/detail/CVE-2018-15664) symlink-exchange attack with directory traversal. [moby/moby#39357](https://github.com/moby/moby/pull/39357)
* Windows: fixed support for `docker service create --limit-cpu`. [moby/moby#39190](https://github.com/moby/moby/pull/39190)
* daemon: fixed a mirrors validation issue.  [moby/moby#38991](https://github.com/moby/moby/pull/38991)
* Docker no longer supports sorting UID and GID ranges in ID maps. [moby/moby#39288](https://github.com/moby/moby/pull/39288)

### Logging

* Added a fix that now allows large log lines for logger plugins. [moby/moby#39038](https://github.com/moby/moby/pull/39038)

### Known Issue
* There are [important changes](/ee/upgrade) to the upgrade process that, if not correctly followed, can have an impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or later.

## 18.09.6
2019-05-06

### Builder
* Fixed `COPY` and `ADD` with multiple `<src>` to not invalidate cache if `DOCKER_BUILDKIT=1`.[moby/moby#38964](https://github.com/moby/moby/issues/38964)

### Networking
* Cleaned up the cluster provider when the agent is closed. [docker/libnetwork#2354](https://github.com/docker/libnetwork/pull/2354)
* Windows: Now selects a random host port if the user does not specify a host port. [docker/libnetwork#2369](https://github.com/docker/libnetwork/pull/2369)

### Known Issues
* There are [important changes](/ee/upgrade) to the upgrade process that, if not correctly followed, can have an impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or later.

## 18.09.5

2019-04-11

### Builder

* Fixed `DOCKER_BUILDKIT=1 docker build --squash ..` [docker/engine#176](https://github.com/docker/engine/pull/176)

### Client

* Fixed tty initial size error. [docker/cli#1775](https://github.com/docker/cli/pull/1775)
* Fixed dial-stdio goroutine leakage. [docker/cli#1795](https://github.com/docker/cli/pull/1795)
* Fixed the stack informer's selector used to track deployment. [docker/cli#1794](https://github.com/docker/cli/pull/1794)

### Networking

* Fixed `network=host` using wrong `resolv.conf` with `systemd-resolved`. [docker/engine#180](https://github.com/docker/engine/pull/180)
* Fixed Windows ARP entries getting corrupted randomly under load. [docker/engine#192](https://github.com/docker/engine/pull/192)

### Runtime
* Now showing stopped containers with restart policy as `Restarting`. [docker/engine#181](https://github.com/docker/engine/pull/181)
* Now using original process spec for execs. [docker/engine#178](https://github.com/docker/engine/pull/178)

### Swarm Mode

* Fixed leaking task resources when nodes are deleted. [docker/engine#185](https://github.com/docker/engine/pull/185)

### Known Issues

* There are [important changes](/ee/upgrade) to the upgrade process that, if not correctly followed, can have an impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or later.

## 18.09.4

 2019-03-28

### Builder

* Fixed [CVE-2019-13139](https://nvd.nist.gov/vuln/detail/CVE-2019-13139) by adding validation for `git ref` to avoid misinterpretation as a flag. [moby/moby#38944](https://github.com/moby/moby/pull/38944)

### Runtime

* Fixed `docker cp` error for filenames greater than 100 characters. [moby/moby#38634](https://github.com/moby/moby/pull/38634)
* Fixed `layer/layer_store` to ensure `NewInputTarStream` resources are released. [moby/moby#38413](https://github.com/moby/moby/pull/38413)
* Increased GRPC limit for `GetConfigs`. [moby/moby#38800](https://github.com/moby/moby/pull/38800)
* Updated `containerd` 1.2.5. [docker/engine#173](https://github.com/docker/engine/pull/173)

### Swarm Mode
* Fixed nil pointer exception when joining node to swarm. [moby/moby#38618](https://github.com/moby/moby/issues/38618)
* Fixed issue for swarm nodes not being able to join as masters if http proxy is set. [moby/moby#36951]

### Known Issues
* There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or later.

## 18.09.3

2019-02-28

### Networking fixes
* Windows: now avoids regeneration of network IDs to prevent broken references to networks. [docker/engine#149](https://github.com/docker/engine/pull/149)
* Windows: Fixed an issue to address `- restart always` flag on standalone containers not working when specifying a network. (docker/escalation#1037)
* Fixed an issue to address the IPAM state from networkdb if the manager is not attached to the overlay network. (docker/escalation#1049)

### Runtime fixes and updates

* Updated to Go version 1.10.8.
* Modified names in the container name generator. [docker/engine#159](https://github.com/docker/engine/pull/159)
* When copying an existing folder, xattr set errors when the target filesystem doesn't support xattr are now ignored. [docker/engine#135](https://github.com/docker/engine/pull/135)
* Graphdriver: fixed "device" mode not being detected if "character-device" bit is set. [docker/engine#160](https://github.com/docker/engine/pull/160)
* Fixed nil pointer derefence on failure to connect to containerd. [docker/engine#162](https://github.com/docker/engine/pull/162)
* Deleted stale containerd object on start failure. [docker/engine#154](https://github.com/docker/engine/pull/154)

### Known Issues
* There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or greater.

## 18.09.2

2019-02-11

### Security fixes for Docker Engine - Enterprise
* Update `runc` to address a critical vulnerability that allows specially-crafted containers to gain administrative privileges on the host. [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)
* Ubuntu 14.04 customers using a 3.13 kernel will need to upgrade to a supported Ubuntu 4.x kernel

For additional information, [refer to the Docker blog post](https://blog.docker.com/2019/02/docker-security-update-cve-2018-5736-and-container-security-best-practices/).

### Known Issues
* There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or greater.

## 18.09.1

2019-01-09

#### Important notes about this release

In Docker versions prior to 18.09, containerd was managed by the Docker engine daemon. In Docker Engine 18.09, containerd is managed by systemd. Since containerd is managed by systemd, any custom configuration to the `docker.service` systemd configuration which changes mount settings (for example, `MountFlags=slave`) breaks interactions between the Docker Engine daemon and containerd, and you will not be able to start containers.

Run the following command to get the current value of the `MountFlags` property for the `docker.service`:

```bash
sudo systemctl show --property=MountFlags docker.service
MountFlags=
```
Update your configuration if this command prints a non-empty value for `MountFlags`, and restart the docker service.

### Security fixes
* Upgraded Go language to 1.10.6 to resolve [CVE-2018-16873](https://nvd.nist.gov/vuln/detail/CVE-2018-16873), [CVE-2018-16874](https://nvd.nist.gov/vuln/detail/CVE-2018-16874), and [CVE-2018-16875](https://nvd.nist.gov/vuln/detail/CVE-2018-16875).
* Fixed authz plugin for 0-length content and path validation.
* Added `/proc/asound` to masked paths [docker/engine#126](https://github.com/docker/engine/pull/126)

### Improvements
* Updated to BuildKit 0.3.3 [docker/engine#122](https://github.com/docker/engine/pull/122)
* Updated to containerd 1.2.2 [docker/engine#144](https://github.com/docker/engine/pull/144)
* Provided additional warnings for use of deprecated legacy overlay and devicemapper storage drivers [docker/engine#85](https://github.com/docker/engine/pull/85)
* prune: perform image pruning before build cache pruning [docker/cli#1532](https://github.com/docker/cli/pull/1532)
* Added bash completion for experimental CLI commands (manifest) [docker/cli#1542](https://github.com/docker/cli/pull/1542)
* Windows: allow process isolation on Windows 10 [docker/engine#81](https://github.com/docker/engine/pull/81)

### Fixes
* Disable kmem accounting in runc on RHEL/CentOS (docker/escalation#614, docker/escalation#692) [docker/engine#121](https://github.com/docker/engine/pull/121)
* Fixed inefficient networking configuration [docker/engine#123](https://github.com/docker/engine/pull/123)
* Fixed docker system prune doesn't accept until filter [docker/engine#122](https://github.com/docker/engine/pull/122)
* Avoid unset credentials in `containerd` [docker/engine#122](https://github.com/docker/engine/pull/122)
* Fixed iptables compatibility on Debian [docker/engine#107](https://github.com/docker/engine/pull/107)
* Fixed setting default schema to tcp for docker host [docker/cli#1454](https://github.com/docker/cli/pull/1454)
* Fixed bash completion for `service update --force`  [docker/cli#1526](https://github.com/docker/cli/pull/1526)
* Windows: DetachVhd attempt in cleanup [docker/engine#113](https://github.com/docker/engine/pull/113)
* API: properly handle invalid JSON to return a 400 status [docker/engine#110](https://github.com/docker/engine/pull/110)
* API: ignore default address-pools on API < 1.39 [docker/engine#118](https://github.com/docker/engine/pull/118)
* API: add missing default address pool fields to swagger [docker/engine#119](https://github.com/docker/engine/pull/119)
* awslogs: account for UTF-8 normalization in limits [docker/engine#112](https://github.com/docker/engine/pull/112)
* Prohibit reading more than 1MB in HTTP error responses [docker/engine#114](https://github.com/docker/engine/pull/114)
* apparmor: allow receiving of signals from `docker kill` [docker/engine#116](https://github.com/docker/engine/pull/116)
* overlay2: use index=off if possible (fix EBUSY on mount) [docker/engine#84](https://github.com/docker/engine/pull/84)

### Packaging
* Add docker.socket requirement for docker.service. [docker/docker-ce-packaging#276](https://github.com/docker/docker-ce-packaging/pull/276)
* Add socket activation for RHEL-based distributions. [docker/docker-ce-packaging#274](https://github.com/docker/docker-ce-packaging/pull/274)
* Add libseccomp requirement for RPM packages. [docker/docker-ce-packaging#266](https://github.com/docker/docker-ce-packaging/pull/266)

### Known Issues
* When upgrading from 18.09.0 to 18.09.1, `containerd` is not upgraded to the correct version on Ubuntu.  [Learn more](https://success.docker.com/article/error-upgrading-to-engine-18091-with-containerd).
* There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or greater.

## 18.09.0

2018-11-08

### Important notes about this release

In Docker versions prior to 18.09, containerd was managed by the Docker engine daemon. In Docker Engine 18.09, containerd is managed by systemd. Since containerd is managed by systemd, any custom configuration to the `docker.service` systemd
configuration which changes mount settings (for example, `MountFlags=slave`) breaks interactions between the Docker Engine daemon and containerd, and you will not be able to start containers.

Run the following command to get the current value of the `MountFlags` property for the `docker.service`:

```bash
sudo systemctl show --property=MountFlags docker.service
MountFlags=
```

Update your configuration if this command prints a non-empty value for `MountFlags`, and restart the docker service.


### New features for Docker Engine EE

* [FIPS Compliance added for Windows Server 2016 and later](/ee/docker-ee/windows/docker-ee.md)
* [Docker Content Trust Enforcement](/engine/security/trust/content_trust.md) for the Enterprise Engine. This allows the Docker Engine - Enterprise to run containers not signed by a specific organization.

### New features

* Updated API version to 1.39 [moby/moby#37640](https://github.com/moby/moby/pull/37640)
* Added support for remote connections using SSH [docker/cli#1014](https://github.com/docker/cli/pull/1014)
* Builder: added prune options to the API [moby/moby#37651](https://github.com/moby/moby/pull/37651)
* Added "Warnings" to `/info` endpoint, and move detection to the daemon [moby/moby#37502](https://github.com/moby/moby/pull/37502)
* Allows BuildKit builds to run without experimental mode enabled. Buildkit can now be configured with an option in daemon.json [moby/moby#37593](https://github.com/moby/moby/pull/37593) [moby/moby#37686](https://github.com/moby/moby/pull/37686) [moby/moby#37692](https://github.com/moby/moby/pull/37692) [docker/cli#1303](https://github.com/docker/cli/pull/1303)  [docker/cli#1275](https://github.com/docker/cli/pull/1275)
* Added support for build-time secrets using a `--secret` flag when using BuildKit [docker/cli#1288](https://github.com/docker/cli/pull/1288)
* Added SSH agent socket forwarder (`docker build --ssh $SSHMOUNTID=$SSH_AUTH_SOCK`) when using BuildKit [docker/cli#1438](https://github.com/docker/cli/pull/1438) / [docker/cli#1419](https://github.com/docker/cli/pull/1419)
* Added `--chown` flag support for `ADD` and `COPY` commands on Windows [moby/moby#35521](https://github.com/moby/moby/pull/35521)
* Added `builder prune` subcommand to prune BuildKit build cache [docker/cli#1295](https://github.com/docker/cli/pull/1295) [docker/cli#1334](https://github.com/docker/cli/pull/1334)
* BuildKit: Adds configurable garbage collection policy for the BuildKit build cache [docker/engine#59](https://github.com/docker/engine/pull/59) / [moby/moby#37846](https://github.com/moby/moby/pull/37846)
* BuildKit: Adds support for `docker build --pull ...` when using BuildKit [moby/moby#37613](https://github.com/moby/moby/pull/37613)
* BuildKit: Adds support or "registry-mirrors" and "insecure-registries" when using BuildKit [docker/engine#59](https://github.com/docker/engine/pull/59) / [moby/moby#37852](https://github.com/moby/moby/pull/37852)
* BuildKit: Enables net modes and bridge. [moby/moby#37620](https://github.com/moby/moby/pull/37620)
* Added `docker engine` subcommand to manage the lifecycle of a Docker Engine running as a privileged container on top of containerd, and to allow upgrades to Docker Engine Enterprise [docker/cli#1260](https://github.com/docker/cli/pull/1260)
* Exposed product license in `docker info` output [docker/cli#1313](https://github.com/docker/cli/pull/1313)
* Showed warnings produced by daemon in `docker info` output [docker/cli#1225](https://github.com/docker/cli/pull/1225)
* Added "local" log driver [moby/moby#37092](https://github.com/moby/moby/pull/37092)
* Amazon CloudWatch: adds `awslogs-endpoint` logging option [moby/moby#37374](https://github.com/moby/moby/pull/37374)
* Added support for global default address pools [moby/moby#37558](https://github.com/moby/moby/pull/37558) [docker/cli#1233](https://github.com/docker/cli/pull/1233)
* Configured containerd log-level to be the same as dockerd [moby/moby#37419](https://github.com/moby/moby/pull/37419)
* Added configuration option for cri-containerd [moby/moby#37519](https://github.com/moby/moby/pull/37519)
* Updates containerd client to v1.2.0-rc.1 [moby/moby#37664](https://github.com/moby/moby/pull/37664), [docker/engine#75](https://github.com/docker/engine/pull/75) / [moby/moby#37710](https://github.com/moby/moby/pull/37710)
* Added support for global default address pools [moby/moby#37558](https://github.com/moby/moby/pull/37558) [docker/cli#1233](https://github.com/docker/cli/pull/1233)
* Moved the `POST /session` endpoint out of experimental. [moby/moby#40028](https://github.com/moby/moby/pull/40028)


### Improvements

* Does not return "`<unknown>`" in /info response [moby/moby#37472](https://github.com/moby/moby/pull/37472)
* BuildKit: Changes `--console=[auto,false,true]` to `--progress=[auto,plain,tty]` [docker/cli#1276](https://github.com/docker/cli/pull/1276)
* BuildKit: Sets BuildKit's ExportedProduct variable to show useful errors in the future. [moby/moby#37439](https://github.com/moby/moby/pull/37439)
* Hides `--data-path-addr` flags when connected to a daemon that doesn't support this option [docker/docker/cli#1240](https://github.com/docker/cli/pull/1240)
* Only shows buildkit-specific flags if BuildKit is enabled [docker/cli#1438](https://github.com/docker/cli/pull/1438) / [docker/cli#1427](https://github.com/docker/cli/pull/1427)
* Improves version output alignment [docker/cli#1204](https://github.com/docker/cli/pull/1204)
* Sorts plugin names and networks in a natural order [docker/cli#1166](https://github.com/docker/cli/pull/1166), [docker/cli#1266](https://github.com/docker/cli/pull/1266)
* Updates bash and zsh [completion scripts](https://github.com/docker/cli/issues?q=label%3Aarea%2Fcompletion+milestone%3A18.09.0+is%3Aclosed)
* Passes log-level to containerd. [moby/moby#37419](https://github.com/moby/moby/pull/37419)
* Uses direct server return (DSR) in east-west overlay load balancing [docker/engine#93](https://github.com/docker/engine/pull/93) / [docker/libnetwork#2270](https://github.com/docker/libnetwork/pull/2270)
* Builder: temporarily disables bridge networking when using buildkit. [moby/moby#37691](https://github.com/moby/moby/pull/37691)
* Blocks task starting until node attachments are ready [moby/moby#37604](https://github.com/moby/moby/pull/37604)
* Propagates the provided external CA certificate to the external CA object in swarm. [docker/cli#1178](https://github.com/docker/cli/pull/1178)
* Removes Ubuntu 14.04 "Trusty Tahr" as a supported platform [docker-ce-packaging#255](https://github.com/docker/docker-ce-packaging/pull/255) / [docker-ce-packaging#254](https://github.com/docker/docker-ce-packaging/pull/254)
* Removes Debian 8 "Jessie" as a supported platform [docker-ce-packaging#255](https://github.com/docker/docker-ce-packaging/pull/255) / [docker-ce-packaging#254](https://github.com/docker/docker-ce-packaging/pull/254)
* Removes 'docker-' prefix for containerd and runc binaries [docker/engine#61](https://github.com/docker/engine/pull/61) / [moby/moby#37907](https://github.com/moby/moby/pull/37907), [docker-ce-packaging#241](https://github.com/docker/docker-ce-packaging/pull/241)
* Splits "engine", "cli", and "containerd" to separate packages, and run containerd as a separate systemd service [docker-ce-packaging#131](https://github.com/docker/docker-ce-packaging/pull/131), [docker-ce-packaging#158](https://github.com/docker/docker-ce-packaging/pull/158)
* Builds binaries with Go 1.10.4 [docker-ce-packaging#181](https://github.com/docker/docker-ce-packaging/pull/181)
* Removes `-ce` / `-ee` suffix from version string [docker-ce-packaging#206](https://github.com/docker/docker-ce-packaging/pull/206)

### Fixes

* BuildKit: Do not cancel buildkit status request. [moby/moby#37597](https://github.com/moby/moby/pull/37597)
* Fixes no error is shown if build args are missing during docker build [moby/moby#37396](https://github.com/moby/moby/pull/37396)
* Fixes error "unexpected EOF" when adding an 8GB file [moby/moby#37771](https://github.com/moby/moby/pull/37771)
* LCOW: Ensures platform is populated on `COPY`/`ADD`. [moby/moby#37563](https://github.com/moby/moby/pull/37563)
* Fixes mapping a range of host ports to a single container port [docker/cli#1102](https://github.com/docker/cli/pull/1102)
* Fixes `trust inspect` typo: "`AdminstrativeKeys`" [docker/cli#1300](https://github.com/docker/cli/pull/1300)
* Fixes environment file parsing for imports of absent variables and those with no name. [docker/cli#1019](https://github.com/docker/cli/pull/1019)
* Fixes a potential "out of memory exception" when running `docker image prune` with a large list of dangling images [docker/cli#1432](https://github.com/docker/cli/pull/1432) / [docker/cli#1423](https://github.com/docker/cli/pull/1423)
* Fixes pipe handling in ConEmu and ConsoleZ on Windows [moby/moby#37600](https://github.com/moby/moby/pull/37600)
* Fixes long startup on windows, with non-hns governed Hyper-V networks [docker/engine#67](https://github.com/docker/engine/pull/67) / [moby/moby#37774](https://github.com/moby/moby/pull/37774)
* Fixes daemon won't start when "runtimes" option is defined both in config file and cli [docker/engine#57](https://github.com/docker/engine/pull/57) / [moby/moby#37871](https://github.com/moby/moby/pull/37871)
* Loosens permissions on `/etc/docker` directory to prevent "permission denied" errors when using `docker manifest inspect` [docker/engine#56](https://github.com/docker/engine/pull/56) / [moby/moby#37847](https://github.com/moby/moby/pull/37847)
* Fixes denial of service with large numbers in `cpuset-cpus` and `cpuset-mems` [docker/engine#70](https://github.com/docker/engine/pull/70) / [moby/moby#37967](https://github.com/moby/moby/pull/37967)
* LCOW: Add `--platform` to `docker import` [docker/cli#1375](https://github.com/docker/cli/pull/1375) / [docker/cli#1371](https://github.com/docker/cli/pull/1371)
* LCOW: Add LinuxMetadata support by default on Windows [moby/moby#37514](https://github.com/moby/moby/pull/37514)
* LCOW: Mount to short container paths to avoid command-line length limit [moby/moby#37659](https://github.com/moby/moby/pull/37659)
* LCOW: Fix builder using wrong cache layer [moby/moby#37356](https://github.com/moby/moby/pull/37356)
* Fixes json-log file descriptors leaking when using `--follow` [docker/engine#48](https://github.com/docker/engine/pull/48) [moby/moby#37576](https://github.com/moby/moby/pull/37576) [moby/moby#37734](https://github.com/moby/moby/pull/37734)
* Fixes a possible deadlock on closing the watcher on kqueue [moby/moby#37392](https://github.com/moby/moby/pull/37392)
* Uses poller based watcher to work around the file caching issue in Windows [moby/moby#37412](https://github.com/moby/moby/pull/37412)
* Handles systemd-resolved case by providing appropriate resolv.conf to networking layer [moby/moby#37485](https://github.com/moby/moby/pull/37485)
* Removes support for TLS < 1.2 [moby/moby#37660](https://github.com/moby/moby/pull/37660)
* Seccomp: Whitelist syscalls linked to `CAP_SYS_NICE` in default seccomp profile [moby/moby#37242](https://github.com/moby/moby/pull/37242)
* Seccomp: move the syslog syscall to be gated by `CAP_SYS_ADMIN` or `CAP_SYSLOG` [docker/engine#64](https://github.com/docker/engine/pull/64) / [moby/moby#37929](https://github.com/moby/moby/pull/37929)
* SELinux: Fix relabeling of local volumes specified via Mounts API on selinux-enabled systems [moby/moby#37739](https://github.com/moby/moby/pull/37739)
* Adds warning if REST API is accessible through an insecure connection [moby/moby#37684](https://github.com/moby/moby/pull/37684)
* Masks proxy credentials from URL when displayed in system info [docker/engine#72](https://github.com/docker/engine/pull/72) / [moby/moby#37934](https://github.com/moby/moby/pull/37934)
* Fixes mount propagation for btrfs [docker/engine#86](https://github.com/docker/engine/pull/86) / [moby/moby#38026](https://github.com/moby/moby/pull/38026)
* Fixes nil pointer dereference in node allocation [docker/engine#94](https://github.com/docker/engine/pull/94) / [docker/swarmkit#2764](https://github.com/docker/swarmkit/pull/2764)

### Known Issues

* There are [important changes to the upgrade process](/ee/upgrade) that, if not correctly followed, can have impact on the availability of applications running on the Swarm during upgrades. These constraints impact any upgrades coming from any version before 18.09 to version 18.09 or greater.
* With https://github.com/boot2docker/boot2docker/releases/download/v18.09.0/boot2docker.iso, connection is being refused from a node on the virtual machine. Any publishing of swarm ports in virtualbox-created docker-machine VM's will not respond. This is occurring on macOS and Windows 10, using docker-machine version 0.15 and 0.16.

   The following `docker run` command works, allowing access from host browser:

   `docker run -d -p 4000:80 nginx`

   However, the following `docker service` command fails, resulting in curl/chrome unable to connect (connection refused):

   `docker service create -p 5000:80 nginx`

   This issue is not apparent when provisioning 18.09.0 cloud VM's using docker-machine.

   Workarounds:
   * Use cloud VM's that don't rely on boot2docker.
   * `docker run` is unaffected.
   * For Swarm, set VIRTUALBOX_BOOT2DOCKER_URL=https://github.com/boot2docker/boot2docker/releases/download/v18.06.1-ce/boot2docker.iso.

   This issue is resolved in 18.09.1.

### Deprecation Notices

- As of EE 2.1, Docker has deprecated support for Device Mapper as a storage driver. It will continue to be
supported at this time, but support will be removed in a future release. Docker will continue to support
Device Mapper for existing EE 2.0 and 2.1 customers. Please contact Sales for more information.

    Docker recommends that existing customers
    [migrate to using Overlay2 for the storage driver](https://success.docker.com/article/how-do-i-migrate-an-existing-ucp-cluster-to-the-overlay2-graph-driver). The [Overlay2 storage driver](https://docs.docker.com/storage/storagedriver/overlayfs-driver/) is now the default for Docker engine implementations.
- As of EE 2.1, Docker has deprecated support for IBM Z (s390x). Refer to the
[Docker Compatibility Matrix](https://success.docker.com/article/compatibility-matrix) for detailed
compatibility information.

For more information on the list of deprecated flags and APIs, have a look at the [deprecation information](https://docs.docker.com/engine/deprecated/) where you can find the target removal dates.

### End of Life Notification

In this release, Docker has also removed support for TLS < 1.2 [moby/moby#37660](https://github.com/moby/moby/pull/37660),
Ubuntu 14.04 "Trusty Tahr" [docker-ce-packaging#255](https://github.com/docker/docker-ce-packaging/pull/255) / [docker-ce-packaging#254](https://github.com/docker/docker-ce-packaging/pull/254), and Debian 8 "Jessie" [docker-ce-packaging#255](https://github.com/docker/docker-ce-packaging/pull/255) / [docker-ce-packaging#254](https://github.com/docker/docker-ce-packaging/pull/254).

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
- Fix error with merge compose file with networks [docker/cli#983](https://github.com/docker/cli/pull/983)
* Fix docker stack deploy re-deploying services after the service was updated with `--force` [docker/cli#963](https://github.com/docker/cli/pull/963)
* Fix docker version output alignment [docker/cli#965](https://github.com/docker/cli/pull/965)
* Simplify the marshaling of compose types.Config [docker/cli#895](https://github.com/docker/cli/pull/895)
+ Add support for multiple composefile when deploying [docker/cli#569](https://github.com/docker/cli/pull/569)
- Fix broken Kubernetes stack flags [docker/cli#831](https://github.com/docker/cli/pull/831)
- Fix stack marshaling for Kubernetes [docker/cli#890](https://github.com/docker/cli/pull/890)
- Fix and simplify bash completion for service env, mounts and labels [docker/cli#682](https://github.com/docker/cli/pull/682)
- Fix `before` and `since` filter for `docker ps` [moby/moby#35938](https://github.com/moby/moby/pull/35938)
- Fix `--label-file` weird behavior [docker/cli#838](https://github.com/docker/cli/pull/838)
- Fix compilation of defaultCredentialStore() on unsupported platforms [docker/cli#872](https://github.com/docker/cli/pull/872)
* Improve and fix bash completion for images [docker/cli#717](https://github.com/docker/cli/pull/717)
+ Added check for empty source in bind mount [docker/cli#824](https://github.com/docker/cli/pull/824)
- Fix TLS from environment variables in client [moby/moby#36270](https://github.com/moby/moby/pull/36270)
* docker build now runs faster when registry-specific credential helper(s) are configured [docker/cli#840](https://github.com/docker/cli/pull/840)
* Update event filter zsh completion with `disable`, `enable`, `install` and `remove` [docker/cli#372](https://github.com/docker/cli/pull/372)
* Produce errors when empty ids are passed into inspect calls [moby/moby#36144](https://github.com/moby/moby/pull/36144)
* Marshall version for the k8s controller [docker/cli#891](https://github.com/docker/cli/pull/891)
* Set a non-zero timeout for HTTP client communication with plugin backend [docker/cli#883](https://github.com/docker/cli/pull/883)
+ Add DOCKER_TLS environment variable for --tls option [docker/cli#863](https://github.com/docker/cli/pull/863)
+ Add --template-driver option for secrets/configs [docker/cli#896](https://github.com/docker/cli/pull/896)
+ Move `docker trust` commands out of experimental [docker/cli#934](https://github.com/docker/cli/pull/934) [docker/cli#935](https://github.com/docker/cli/pull/935) [docker/cli#944](https://github.com/docker/cli/pull/944)

### Builder

* Switch to -buildmode=pie [moby/moby#34369](https://github.com/moby/moby/pull/34369)
* Allow Dockerfile to be outside of build-context [docker/cli#886](https://github.com/docker/cli/pull/886)
* Builder: fix wrong cache hits building from tars [moby/moby#36329](https://github.com/moby/moby/pull/36329)
- Fixes files leaking to other images in a multi-stage build [moby/moby#36338](https://github.com/moby/moby/pull/36338)

### Runtime

+ Update to docker-ce 18.03.1 engine.
+ Add support for FIPS 140-2 on x86_64.
+ Add support for Microsoft Windows Server 1709 and 1803 with support for [swarm ingress routing mesh](https://docs.docker.com/engine/swarm/ingress/), [VIP service discovery](https://docs.docker.com/v17.09/engine/swarm/networking/#configure-service-discovery), and [named pipe mounting](https://blog.docker.com/2017/09/docker-windows-server-1709/).
+ Add support for Ubuntu 18.04.
+ Windows opt-out telemetry stream.
+ Support for `--chown` with `COPY` and `ADD` in `Dockerfile`.
+ Added functionality for the `docker logs` command to include the output of multiple logging drivers.
- Fix AppArmor profiles not being applied to `docker exec` processes [moby/moby#36466](https://github.com/moby/moby/pull/36466)
- Don't sort plugin mount slice [moby/moby#36711](https://github.com/moby/moby/pull/36711)
- Daemon/cluster: handle partial attachment entries during configure [moby/moby#36769](https://github.com/moby/moby/pull/36769)
* Bump Golang to 1.9.5 [moby/moby#36779](https://github.com/moby/moby/pull/36779) [docker/cli#986](https://github.com/docker/cli/pull/986)
- Daemon/stats: more resilient cpu sampling [moby/moby#36519](https://github.com/moby/moby/pull/36519)
* Containerd: update to 1.0.3 release [moby/moby#36749](https://github.com/moby/moby/pull/36749)
- Fix Windows layer leak when write fails [moby/moby#36728](https://github.com/moby/moby/pull/36728)
* Don't make container mount unbindable [moby/moby#36768](https://github.com/moby/moby/pull/36768)
- Fix Daemon panics on container export after a daemon restart [moby/moby/36586](https://github.com/moby/moby/pull/36586)
- Fix digest cache being removed on autherrors [moby/moby#36509](https://github.com/moby/moby/pull/36509)
- Make sure plugin container is removed on failure [moby/moby#36715](https://github.com/moby/moby/pull/36715)
- Copy: avoid using all system memory with authz plugins [moby/moby#36595](https://github.com/moby/moby/pull/36595)
- Relax some libcontainerd client locking [moby/moby#36848](https://github.com/moby/moby/pull/36848)
- Update `hcsshim` to v0.6.10 to address [CVE-2018-8115](https://portal.msrc.microsoft.com/en-us/security-guidance/advisory/CVE-2018-8115)
* Enable HotAdd for Windows [moby/moby#35414](https://github.com/moby/moby/pull/35414)
* LCOW: Graphdriver fix deadlock in hotRemoveVHDs [moby/moby#36114](https://github.com/moby/moby/pull/36114)
* LCOW: Regular mount if only one layer [moby/moby#36052](https://github.com/moby/moby/pull/36052)
* Remove interim env var LCOW_API_PLATFORM_IF_OMITTED [moby/moby#36269](https://github.com/moby/moby/pull/36269)
* Revendor Microsoft/opengcs @ v0.3.6 [moby/moby#36108](https://github.com/moby/moby/pull/36108)
- Fix issue of ExitCode and PID not show up in Task.Status.ContainerStatus [moby/moby#36150](https://github.com/moby/moby/pull/36150)
- Fix issue with plugin scanner going too deep [moby/moby#36119](https://github.com/moby/moby/pull/36119)
* Do not make graphdriver homes private mounts [moby/moby#36047](https://github.com/moby/moby/pull/36047)
* Do not recursive unmount on cleanup of zfs/btrfs [moby/moby#36237](https://github.com/moby/moby/pull/36237)
* Don't restore image if layer does not exist [moby/moby#36304](https://github.com/moby/moby/pull/36304)
* Adjust minimum API version for templated configs/secrets [moby/moby#36366](https://github.com/moby/moby/pull/36366)
* Bump containerd to 1.0.2 (cfd04396dc68220d1cecbe686a6cc3aa5ce3667c) [moby/moby#36308](https://github.com/moby/moby/pull/36308)
* Bump Golang to 1.9.4 [moby/moby#36243](https://github.com/moby/moby/pull/36243)
* Ensure daemon root is unmounted on shutdown [moby/moby#36107](https://github.com/moby/moby/pull/36107)
* Update runc to 6c55f98695e902427906eed2c799e566e3d3dfb5 [moby/moby#36222](https://github.com/moby/moby/pull/36222)
- Fix container cleanup on daemon restart [moby/moby#36249](https://github.com/moby/moby/pull/36249)
* Support SCTP port mapping (bump up API to v1.37) [moby/moby#33922](https://github.com/moby/moby/pull/33922)
* Support SCTP port mapping [docker/cli#278](https://github.com/docker/cli/pull/278)
- Fix Volumes property definition in ContainerConfig [moby/moby#35946](https://github.com/moby/moby/pull/35946)
* Bump moby and dependencies [docker/cli#829](https://github.com/docker/cli/pull/829)
* C.RWLayer: check for nil before use [moby/moby#36242](https://github.com/moby/moby/pull/36242)
+ Add `REMOVE` and `ORPHANED` to TaskState [moby/moby#36146](https://github.com/moby/moby/pull/36146)
- Fixed error detection using `IsErrNotFound` and `IsErrNotImplemented` for `ContainerStatPath`, `CopyFromContainer`, and `CopyToContainer` methods [moby/moby#35979](https://github.com/moby/moby/pull/35979)
+ Add an integration/internal/container helper package [moby/moby#36266](https://github.com/moby/moby/pull/36266)
+ Add canonical import path [moby/moby#36194](https://github.com/moby/moby/pull/36194)
+ Add/use container.Exec() to integration [moby/moby#36326](https://github.com/moby/moby/pull/36326)
- Fix "--node-generic-resource" singular/plural [moby/moby#36125](https://github.com/moby/moby/pull/36125)
* Daemon.cleanupContainer: nullify container RWLayer upon release [moby/moby#36160](https://github.com/moby/moby/pull/36160)
* Daemon: passdown the `--oom-kill-disable` option to containerd [moby/moby#36201](https://github.com/moby/moby/pull/36201)
* Display a warn message when there is binding ports and net mode is host [moby/moby#35510](https://github.com/moby/moby/pull/35510)
* Refresh containerd remotes on containerd restarted [moby/moby#36173](https://github.com/moby/moby/pull/36173)
* Set daemon root to use shared propagation [moby/moby#36096](https://github.com/moby/moby/pull/36096)
* Optimizations for recursive unmount [moby/moby#34379](https://github.com/moby/moby/pull/34379)
* Perform plugin mounts in the runtime [moby/moby#35829](https://github.com/moby/moby/pull/35829)
* Graphdriver: Fix RefCounter memory leak [moby/moby#36256](https://github.com/moby/moby/pull/36256)
* Use continuity fs package for volume copy [moby/moby#36290](https://github.com/moby/moby/pull/36290)
* Use proc/exe for reexec [moby/moby#36124](https://github.com/moby/moby/pull/36124)
+ Add API support for templated secrets and configs [moby/moby#33702](https://github.com/moby/moby/pull/33702) and [moby/moby#36366](https://github.com/moby/moby/pull/36366)
* Use rslave propagation for mounts from daemon root [moby/moby#36055](https://github.com/moby/moby/pull/36055)
+ Add /proc/keys to masked paths [moby/moby#36368](https://github.com/moby/moby/pull/36368)
* Bump Runc to 1.0.0-rc5 [moby/moby#36449](https://github.com/moby/moby/pull/36449)
- Fixes `runc exec` on big-endian architectures [moby/moby#36449](https://github.com/moby/moby/pull/36449)
* Use chroot when mount namespaces aren't provided [moby/moby#36449](https://github.com/moby/moby/pull/36449)
- Fix systemd slice expansion so that it could be consumed by cAdvisor [moby/moby#36449](https://github.com/moby/moby/pull/36449)
- Fix devices mounted with wrong uid/gid [moby/moby#36449](https://github.com/moby/moby/pull/36449)
- Fix read-only containers with IPC private mounts `/dev/shm` read-only [moby/moby#36526](https://github.com/moby/moby/pull/36526)


### Logging

* AWS logs - don't add new lines to maximum sized events [moby/moby#36078](https://github.com/moby/moby/pull/36078)
* Move log validator logic after plugins are loaded [moby/moby#36306](https://github.com/moby/moby/pull/36306)
* Support a proxy in Splunk log driver [moby/moby#36220](https://github.com/moby/moby/pull/36220)
- Fix log tail with empty logs [moby/moby#36305](https://github.com/moby/moby/pull/36305)


### Networking

* Gracefully remove LB endpoints from services [docker/libnetwork#2112](https://github.com/docker/libnetwork/pull/2112)
* Retry other external DNS servers on ServFail [docker/libnetwork#2121](https://github.com/docker/libnetwork/pull/2121)
* Improve scalabiltiy of bridge network isolation rules [docker/libnetwork#2117](https://github.com/docker/libnetwork/pull/2117)
* Allow for larger preset property values, do not override [docker/libnetwork#2124](https://github.com/docker/libnetwork/pull/2124)
* Prevent panics on concurrent reads/writes when calling `changeNodeState` [docker/libnetwork#2136](https://github.com/docker/libnetwork/pull/2136)
* Libnetwork revendoring [moby/moby#36137](https://github.com/moby/moby/pull/36137)
- Fix for deadlock on exit with Memberlist revendor [docker/libnetwork#2040](https://github.com/docker/libnetwork/pull/2040)
* Fix user specified ndots option [docker/libnetwork#2065](https://github.com/docker/libnetwork/pull/2065)
- Fix to use ContainerID for Windows instead of SandboxID [docker/libnetwork#2010](https://github.com/docker/libnetwork/pull/2010)
* Verify NetworkingConfig to make sure EndpointSettings is not nil [moby/moby#36077](https://github.com/moby/moby/pull/36077)
- Fix `DockerNetworkInternalMode` issue [moby/moby#36298](https://github.com/moby/moby/pull/36298)
- Fix race in attachable network attachment [moby/moby#36191](https://github.com/moby/moby/pull/36191)
- Fix timeout issue of `InspectNetwork` on AArch64 [moby/moby#36257](https://github.com/moby/moby/pull/36257)
* Verbose info is missing for partial overlay ID [moby/moby#35989](https://github.com/moby/moby/pull/35989)
* Update `FindNetwork` to address network name duplications [moby/moby#30897](https://github.com/moby/moby/pull/30897)
* Disallow attaching ingress network [docker/swarmkit#2523](https://github.com/docker/swarmkit/pull/2523)
- Prevent implicit removal of the ingress network [moby/moby#36538](https://github.com/moby/moby/pull/36538)
- Fix stale HNS endpoints on Windows [moby/moby#36603](https://github.com/moby/moby/pull/36603)
- IPAM fixes for duplicate IP addresses [docker/libnetwork#2104](https://github.com/docker/libnetwork/pull/2104) [docker/libnetwork#2105](https://github.com/docker/libnetwork/pull/2105)

### Swarm Mode

* Increase raft Election tick to 10 times Heartbeat tick [moby/moby#36672](https://github.com/moby/moby/pull/36672)
* Replace EC Private Key with PKCS#8 PEMs [docker/swarmkit#2246](https://github.com/docker/swarmkit/pull/2246)
* Fix IP overlap with empty EndpointSpec [docker/swarmkit #2505](https://github.com/docker/swarmkit/pull/2505)
* Add support for Support SCTP port mapping [docker/swarmkit#2298](https://github.com/docker/swarmkit/pull/2298)
* Do not reschedule tasks if only placement constraints change and are satisfied by the assigned node [docker/swarmkit#2496](https://github.com/docker/swarmkit/pull/2496)
* Ensure task reaper stopChan is closed no more than once [docker/swarmkit #2491](https://github.com/docker/swarmkit/pull/2491)
* Synchronization fixes [docker/swarmkit#2495](https://github.com/docker/swarmkit/pull/2495)
* Add log message to indicate message send retry if streaming unimplemented [docker/swarmkit#2483](https://github.com/docker/swarmkit/pull/2483)
* Debug logs for session, node events on dispatcher, heartbeats [docker/swarmkit#2486](https://github.com/docker/swarmkit/pull/2486)
+ Add swarm types to bash completion event type filter [docker/cli#888](https://github.com/docker/cli/pull/888)
- Fix issue where network inspect does not show Created time for networks in swarm scope [moby/moby#36095](https://github.com/moby/moby/pull/36095)

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

## 17.03.2-ee-8
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

## 17.03.2-ee-7
2017-10-04

* Fix logic in network resource reaping to prevent memory leak [docker/libnetwork#1944](https://github.com/docker/libnetwork/pull/1944) [docker/libnetwork#1960](https://github.com/docker/libnetwork/pull/1960)
* Increase max GRPC message size to 128MB for larger snapshots so newly added managers can successfully join [docker/swarmkit#2375](https://github.com/docker/swarmkit/pull/2375)

## 17.03.2-ee-6
2017-08-24

* Fix daemon panic on docker image push [moby/moby#33105](https://github.com/moby/moby/pull/33105)
* Fix panic in concurrent network creation/deletion operations [docker/libnetwork#1861](https://github.com/docker/libnetwork/pull/1861)
* Improve network db stability under stressful situations [docker/libnetwork#1860](https://github.com/docker/libnetwork/pull/1860)
* Enable TCP Keep-Alive in Docker client [docker/cli#415](https://github.com/docker/cli/pull/415)
* Lock goroutine to OS thread while changing NS [docker/libnetwork#1911](https://github.com/docker/libnetwork/pull/1911)
* Ignore PullOptions for running tasks [docker/swarmkit#2351](https://github.com/docker/swarmkit/pull/2351)

## 17.03.2-ee-5
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

## 17.03.2-ee-4
2017-06-01

> **Note**
>
> This release includes a fix for potential data loss under certain
> circumstances with the local (built-in) volume driver.

### Networking

- Fix a concurrency issue preventing network creation [#33273](https://github.com/moby/moby/pull/33273)

### Runtime

- Relabel secrets path to avoid a Permission Denied on selinux enabled systems [#33236](https://github.com/moby/moby/pull/33236) (ref [#32529](https://github.com/moby/moby/pull/32529)
- Fix cases where local volume were not properly relabeled if needed [#33236](https://github.com/moby/moby/pull/33236) (ref [#29428](https://github.com/moby/moby/pull/29428))
- Fix an issue while upgrading if a plugin rootfs was still mounted [#33236](https://github.com/moby/moby/pull/33236) (ref [#32525](https://github.com/moby/moby/pull/32525))
- Fix an issue where volume wouldn't default to the `rprivate` propagation mode [#33236](https://github.com/moby/moby/pull/33236) (ref [#32851](https://github.com/moby/moby/pull/32851))
- Fix a panic that could occur when a volume driver could not be retrieved [#33236](https://github.com/moby/moby/pull/33236) (ref [#32347](https://github.com/moby/moby/pull/32347))
+ Add a warning in `docker info` when the `overlay` or `overlay2` graphdriver is used on a filesystem without `d_type` support [#33236](https://github.com/moby/moby/pull/33236) (ref [#31290](https://github.com/moby/moby/pull/31290))
- Fix an issue with backporting mount spec to older volumes [#33207](https://github.com/moby/moby/pull/33207)
- Fix issue where a failed unmount can lead to data loss on local volume remove [#33120](https://github.com/moby/moby/pull/33120)

### Swarm Mode

- Fix a case where tasks could get killed unexpectedly [#33118](https://github.com/moby/moby/pull/33118)
- Fix an issue preventing to deploy services if the registry cannot be reached despite the needed images being locally present [#33117](https://github.com/moby/moby/pull/33117)


## 17.03.1-ee-3
2017-03-30

* Fix an issue with the SELinux policy for Oracle Linux [#31501](https://github.com/docker/docker/pull/31501)

## 17.03.1-ee-2
2017-03-28

### Remote API (v1.27) & Client

* Fix autoremove on older api [#31692](https://github.com/docker/docker/pull/31692)
* Fix default network customization for a stack [#31258](https://github.com/docker/docker/pull/31258/)
* Correct CPU usage calculation in presence of offline CPUs and newer Linux [#31802](https://github.com/docker/docker/pull/31802)
* Fix issue where service healthcheck is `{}` in remote API [#30197](https://github.com/docker/docker/pull/30197)

### Runtime

* Update runc to 54296cf40ad8143b62dbcaa1d90e520a2136ddfe [#31666](https://github.com/docker/docker/pull/31666)
 * Ignore cgroup2 mountpoints [opencontainers/runc#1266](https://github.com/opencontainers/runc/pull/1266)
* Update containerd to 4ab9917febca54791c5f071a9d1f404867857fcc [#31662](https://github.com/docker/docker/pull/31662) [#31852](https://github.com/docker/docker/pull/31852)
 * Register healtcheck service before calling restore() [docker/containerd#609](https://github.com/docker/containerd/pull/609)
* Fix `docker exec` not working after unattended upgrades that reload apparmor profiles [#31773](https://github.com/docker/docker/pull/31773)
* Fix unmounting layer without merge dir with Overlay2 [#31069](https://github.com/docker/docker/pull/31069)
* Do not ignore "volume in use" errors when force-delete [#31450](https://github.com/docker/docker/pull/31450)

### Swarm Mode

* Fix issue with swarm CA timeouts [#2063](https://github.com/docker/swarmkit/pull/2063) [#2064](https://github.com/docker/swarmkit/pull/2064/files)
* Update swarmkit to 17756457ad6dc4d8a639a1f0b7a85d1b65a617bb [#31807](https://github.com/docker/docker/pull/31807)
 * Scheduler now correctly considers tasks which have been assigned to a node but aren't yet running [docker/swarmkit#1980](https://github.com/docker/swarmkit/pull/1980)
 * Allow removal of a network when only dead tasks reference it [docker/swarmkit#2018](https://github.com/docker/swarmkit/pull/2018)
 * Retry failed network allocations less aggressively [docker/swarmkit#2021](https://github.com/docker/swarmkit/pull/2021)
 * Avoid network allocation for tasks that are no longer running [docker/swarmkit#2017](https://github.com/docker/swarmkit/pull/2017)
 * Bookkeeping fixes inside network allocator allocator [docker/swarmkit#2019](https://github.com/docker/swarmkit/pull/2019) [docker/swarmkit#2020](https://github.com/docker/swarmkit/pull/2020)

### Windows

* Cleanup HCS on restore [#31503](https://github.com/docker/docker/pull/31503)

## 17.03.0-ee-1
2017-03-02

Initial Docker EE release, based on Docker CE 17.03.0

> **IMPORTANT**
>
> Starting with this release, Docker is on a monthly release cycle and uses a
> new YY.MM versioning scheme to reflect this. Two channels are available: monthly and quarterly.
> Any given monthly release will only receive security and bugfixes until the next monthly
> release is available. Quarterly releases receive security and bugfixes for 4 months after
> initial release. This release includes bugfixes for 1.13.1 but
> there are no major feature additions and the API version stays the same.
> Upgrading from Docker 1.13.1 to 17.03.0 is expected to be simple and low-risk.

### Client

* Fix panic in `docker stats --format` [#30776](https://github.com/docker/docker/pull/30776)

### Contrib

* Update various `bash` and `zsh` completion scripts [#30823](https://github.com/docker/docker/pull/30823), [#30945](https://github.com/docker/docker/pull/30945) and more...
* Block obsolete socket families in default seccomp profile - mitigates unpatched kernels' CVE-2017-6074 [#29076](https://github.com/docker/docker/pull/29076)

### Networking

* Fix bug on overlay encryption keys rotation in cross-datacenter swarm [#30727](https://github.com/docker/docker/pull/30727)
* Fix side effect panic in overlay encryption and network control plane communication failure ("No installed keys could decrypt the message") on frequent swarm leader re-election [#25608](https://github.com/docker/docker/pull/25608)
* Several fixes around system responsiveness and datapath programming when using overlay network with external kv-store [docker/libnetwork#1639](https://github.com/docker/libnetwork/pull/1639), [docker/libnetwork#1632](https://github.com/docker/libnetwork/pull/1632) and more...
* Discard incoming plain vxlan packets for encrypted overlay network [#31170](https://github.com/docker/docker/pull/31170)
* Release the network attachment on allocation failure [#31073](https://github.com/docker/docker/pull/31073)
* Fix port allocation when multiple published ports map to the same target port [docker/swarmkit#1835](https://github.com/docker/swarmkit/pull/1835)

### Runtime

* Optimize size calculation for `docker system df` container size [#31159](https://github.com/docker/docker/pull/31159)
* Fix a deadlock in docker logs [#30223](https://github.com/docker/docker/pull/30223)
* Fix CPU spin waiting for log write events [#31070](https://github.com/docker/docker/pull/31070)
* Fix a possible crash when using journald [#31231](https://github.com/docker/docker/pull/31231) [#31263](https://github.com/docker/docker/pull/31263)
* Fix a panic on close of nil channel [#31274](https://github.com/docker/docker/pull/31274)
* Fix duplicate mount point for `--volumes-from` in `docker run` [#29563](https://github.com/docker/docker/pull/29563)
* Fix `--cache-from` does not cache last step [#31189](https://github.com/docker/docker/pull/31189)

### Swarm Mode

* Shutdown leaks an error when the container was never started [#31279](https://github.com/docker/docker/pull/31279)
* Fix possibility of tasks getting stuck in the "NEW" state during a leader failover [docker/swarmkit#1938](https://github.com/docker/swarmkit/pull/1938)
* Fix extraneous task creations for global services that led to confusing replica counts in `docker service ls` [docker/swarmkit#1957](https://github.com/docker/swarmkit/pull/1957)
* Fix problem that made rolling updates slow when `task-history-limit` was set to 1 [docker/swarmkit#1948](https://github.com/docker/swarmkit/pull/1948)
* Restart tasks elsewhere, if appropriate, when they are shut down as a result of nodes no longer satisfying constraints [docker/swarmkit#1958](https://github.com/docker/swarmkit/pull/1958)
* (experimental)
