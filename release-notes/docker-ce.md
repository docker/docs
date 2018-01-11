---
title: Docker CE release notes
description: Release notes for Docker CE
keywords: release notes, community
toc_max: 2
toc_min: 1
---

{% raw %}
For Docker Enterprise Edition, see [Docker EE](/ee/index.md).

For Docker releases prior to 17.03.0, see
[Docker Engine release notes](/release-notes/docker-engine.md).

[Learn about Docker releases](/engine/installation.md).

Release notes for stable versions are listed first. You can
[go straight to the Edge release notes](#edge-releases) or
[learn more about Stable and Edge releases](/engine/installation/).

# Stable releases

## 18.03.0-ce (2018-03-21)

### Builder

* Switch to -buildmode=pie [moby/moby#34369](https://github.com/moby/moby/pull/34369)
* Allow Dockerfile to be outside of build-context [docker/cli#886](https://github.com/docker/cli/pull/886)
* Builder: fix wrong cache hits building from tars [moby/moby#36329](https://github.com/moby/moby/pull/36329)
- Fixes files leaking to other images in a multi-stage build [moby/moby#36338](https://github.com/moby/moby/pull/36338)

### Client

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

### Logging

* AWS logs - don't add new lines to maximum sized events [moby/moby#36078](https://github.com/moby/moby/pull/36078)
* Move log validator logic after plugins are loaded [moby/moby#36306](https://github.com/moby/moby/pull/36306)
* Support a proxy in Splunk log driver [moby/moby#36220](https://github.com/moby/moby/pull/36220)
- Fix log tail with empty logs [moby/moby#36305](https://github.com/moby/moby/pull/36305)

### Networking

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

### Runtime

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

### Swarm Mode

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

## 17.12.1-ce (2018-02-27)

### Client
- Fix `node-generic-resource` typo [moby/moby#35970](https://github.com/moby/moby/pull/35970) and [moby/moby#36125](https://github.com/moby/moby/pull/36125)
* Return errors from daemon on stack deploy configs create/update [docker/cli#757](https://github.com/docker/cli/pull/757)

### Logging
- awslogs: fix batch size calculation for large logs [moby/moby#35726](https://github.com/moby/moby/pull/35726)
* Support a proxy in splunk log driver [moby/moby#36220](https://github.com/moby/moby/pull/36220)

### Networking
- Fix ingress network when upgrading from 17.09 to 17.12 [moby/moby#36003](https://github.com/moby/moby/pull/36003)
* Add verbose info to partial overlay ID [moby/moby#35989](https://github.com/moby/moby/pull/35989)
- Fix IPv6 networking being deconfigured if live-restore is being enabled [docker/libnetwork#2043](https://github.com/docker/libnetwork/pull/2043)
- Fix watchMiss thread context [docker/libnetwork#2051](https://github.com/docker/libnetwork/pull/2051)

### Packaging
- Set TasksMax in docker.service [docker/docker-ce-packaging#78](https://github.com/docker/docker-ce-packaging/pull/78)

### Runtime
* Bump Golang to 1.9.4
* Bump containerd to 1.0.1
- Fix dockerd not being able to reconnect to containerd when it is restarted [moby/moby#36173](https://github.com/moby/moby/pull/36173)
- Fix containerd events from being processed twice [moby/moby#35891](https://github.com/moby/moby/issues/35891)
- Fix vfs graph driver failure to initialize because of failure to setup fs quota [moby/moby#35827](https://github.com/moby/moby/pull/35827)
- Fix regression of health check not using container's working directory [moby/moby#35845](https://github.com/moby/moby/pull/35845)
- Honor `DOCKER_RAMDISK` with containerd 1.0 [moby/moby#35957](https://github.com/moby/moby/pull/35957)
- Update runc to fix hang during start and exec [moby/moby#36097](https://github.com/moby/moby/pull/36097)
- Windows: Vendor of Microsoft/hcsshim @v.0.6.8 partial fix for import layer failing [moby/moby#35924](https://github.com/moby/moby/pull/35924)
* Do not make graphdriver homes private mounts [moby/moby#36047](https://github.com/moby/moby/pull/36047)
* Use rslave propogation for mounts from daemon root [moby/moby#36055](https://github.com/moby/moby/pull/36055)
* Set daemon root to use shared mount propagation [moby/moby#36096](https://github.com/moby/moby/pull/36096)
* Validate that mounted paths exist when container is started, not just during creation  [moby/moby#35833](https://github.com/moby/moby/pull/35833)
* Add `REMOVE` and `ORPHANED` to TaskState [moby/moby#36146](https://github.com/moby/moby/pull/36146)
- Fix issue where network inspect does not show Created time for networks in swarm scope [moby/moby#36095](https://github.com/moby/moby/pull/36095)
* Nullify container read write layer upon release [moby/moby#36130](https://github.com/moby/moby/pull/36160) and [moby/moby#36343](https://github.com/moby/moby/pull/36242)

### Swarm
* Remove watchMiss from swarm mode [docker/libnetwork#2047](https://github.com/docker/libnetwork/pull/2047)

## 17.12.0-ce (2017-12-27)

### Known Issues
* AWS logs batch size calculation [moby/moby#35726](https://github.com/moby/moby/pull/35726)
* Health check no longer uses the container's working directory [moby/moby#35843](https://github.com/moby/moby/issues/35843)
* Errors not returned from client in stack deploy configs [moby/moby#757](https://github.com/docker/cli/pull/757)
* Daemon aborts when project quota fails [moby/moby#35827](https://github.com/moby/moby/pull/35827)
* Docker cannot use memory limit when using systemd options [moby/moby#35123](https://github.com/moby/moby/issues/35123)

### Builder

- Fix build cache hash for broken symlink [moby/moby#34271](https://github.com/moby/moby/pull/34271)
- Fix long stream sync [moby/moby#35404](https://github.com/moby/moby/pull/35404)
- Fix dockerfile parser failing silently on long tokens [moby/moby#35429](https://github.com/moby/moby/pull/35429)

### Client

* Remove secret/config duplication in cli/compose [docker/cli#671](https://github.com/docker/cli/pull/671)
* Add `--local` flag to `docker trust sign` [docker/cli#575](https://github.com/docker/cli/pull/575)
* Add `docker trust inspect` [docker/cli#694](https://github.com/docker/cli/pull/694)
+ Add `name` field to secrets and configs to allow interpolation in Compose files [docker/cli#668](https://github.com/docker/cli/pull/668)
+ Add `--isolation` for setting swarm service isolation mode [docker/cli#426](https://github.com/docker/cli/pull/426)
* Remove deprecated "daemon" subcommand [docker/cli#689](https://github.com/docker/cli/pull/689)
- Fix behaviour of `rmi -f` with unexpected errors [docker/cli#654](https://github.com/docker/cli/pull/654)
* Integrated Generic resource in service create [docker/cli#429](https://github.com/docker/cli/pull/429)
- Fix external networks in stacks [docker/cli#743](https://github.com/docker/cli/pull/743)
* Remove support for referencing images by image shortid [docker/cli#753](https://github.com/docker/cli/pull/753) and [moby/moby#35790](https://github.com/moby/moby/pull/35790)
* Use commit-sha instead of tag for containerd [moby/moby#35770](https://github.com/moby/moby/pull/35770)

### Documentation

* Update API version history for 1.35 [moby/moby#35724](https://github.com/moby/moby/pull/35724)

### Logging

* Logentries driver line-only=true []byte output fix [moby/moby#35612](https://github.com/moby/moby/pull/35612)
* Logentries line-only logopt fix to maintain backwards compatibility [moby/moby#35628](https://github.com/moby/moby/pull/35628)
+ Add `--until` flag for docker logs [moby/moby#32914](https://github.com/moby/moby/pull/32914)
+ Add gelf log driver plugin to Windows build [moby/moby#35073](https://github.com/moby/moby/pull/35073)
* Set timeout on splunk batch send [moby/moby#35496](https://github.com/moby/moby/pull/35496)
* Update Graylog2/go-gelf [moby/moby#35765](https://github.com/moby/moby/pull/35765)

### Networking

* Move load balancer sandbox creation/deletion into libnetwork [moby/moby#35422](https://github.com/moby/moby/pull/35422)
* Only chown network files within container metadata [moby/moby#34224](https://github.com/moby/moby/pull/34224)
* Restore error type in FindNetwork [moby/moby#35634](https://github.com/moby/moby/pull/35634)
- Fix consumes MIME type for NetworkConnect [moby/moby#35542](https://github.com/moby/moby/pull/35542)
+ Added support for persisting Windows network driver specific options [moby/moby#35563](https://github.com/moby/moby/pull/35563)
- Fix timeout on netlink sockets and watchmiss leak [moby/moby#35677](https://github.com/moby/moby/pull/35677)
+ New daemon config for networking diagnosis [moby/moby#35677](https://github.com/moby/moby/pull/35677)
- Clean up node management logic [docker/libnetwork#2036](https://github.com/docker/libnetwork/pull/2036)
- Allocate VIPs when endpoints are restored [docker/swarmkit#2474](https://github.com/docker/swarmkit/pull/2474)

### Runtime

* Update to containerd v1.0.0 [moby/moby#35707](https://github.com/moby/moby/pull/35707)
* Have VFS graphdriver use accelerated in-kernel copy [moby/moby#35537](https://github.com/moby/moby/pull/35537)
* Introduce `workingdir` option for docker exec [moby/moby#35661](https://github.com/moby/moby/pull/35661)
* Bump Go to 1.9.2 [moby/moby#33892](https://github.com/moby/moby/pull/33892) [docker/cli#716](https://github.com/docker/cli/pull/716)
* `/dev` should not be readonly with `--readonly` flag [moby/moby#35344](https://github.com/moby/moby/pull/35344)
+ Add custom build-time Graphdrivers priority list [moby/moby#35522](https://github.com/moby/moby/pull/35522)
* LCOW: CLI changes to add platform flag - pull, run, create and build [docker/cli#474](https://github.com/docker/cli/pull/474)
* Fix width/height on Windoes for `docker exec` [moby/moby#35631](https://github.com/moby/moby/pull/35631)
* Detect overlay2 support on pre-4.0 kernels [moby/moby#35527](https://github.com/moby/moby/pull/35527)
* Devicemapper: remove container rootfs mountPath after umount [moby/moby#34573](https://github.com/moby/moby/pull/34573)
* Disallow overlay/overlay2 on top of NFS [moby/moby#35483](https://github.com/moby/moby/pull/35483)
- Fix potential panic during plugin set. [moby/moby#35632](https://github.com/moby/moby/pull/35632)
- Fix some issues with locking on the container [moby/moby#35501](https://github.com/moby/moby/pull/35501)
- Fixup some issues with plugin refcounting [moby/moby#35265](https://github.com/moby/moby/pull/35265)
+ Add missing lock in ProcessEvent [moby/moby#35516](https://github.com/moby/moby/pull/35516)
+ Add vfs quota support [moby/moby#35231](https://github.com/moby/moby/pull/35231)
* Skip empty directories on prior graphdriver detection [moby/moby#35528](https://github.com/moby/moby/pull/35528)
* Skip xfs quota tests when running in user namespace [moby/moby#35526](https://github.com/moby/moby/pull/35526)
+ Added SubSecondPrecision to config option. [moby/moby#35529](https://github.com/moby/moby/pull/35529)
* Update fsnotify to fix deadlock in removing watch [moby/moby#35453](https://github.com/moby/moby/pull/35453)
- Fix "duplicate mount point" when `--tmpfs /dev/shm` is used [moby/moby#35467](https://github.com/moby/moby/pull/35467)
- Fix honoring tmpfs-size for user `/dev/shm` mount [moby/moby#35316](https://github.com/moby/moby/pull/35316)
- Fix EBUSY errors under overlayfs and v4.13+ kernels [moby/moby#34948](https://github.com/moby/moby/pull/34948)
* Container: protect health monitor channel [moby/moby#35482](https://github.com/moby/moby/pull/35482)
* Container: protect the health status with mutex [moby/moby#35517](https://github.com/moby/moby/pull/35517)
* Container: update real-time resources [moby/moby#33731](https://github.com/moby/moby/pull/33731)
* Create labels when volume exists only remotely [moby/moby#34896](https://github.com/moby/moby/pull/34896)
- Fix leaking container/exec state [moby/moby#35484](https://github.com/moby/moby/pull/35484)
* Disallow using legacy (v1) registries [moby/moby#35751](https://github.com/moby/moby/pull/35751) and [docker/cli#747](https://github.com/docker/cli/pull/747)
- Windows: Fix case insensitive filename matching against builder cache [moby/moby#35793](https://github.com/moby/moby/pull/35793)
- Fix race conditions around process handling and error checks [moby/moby#35809](https://github.com/moby/moby/pull/35809)
* Ensure containers are stopped on daemon startup [moby/moby#35805](https://github.com/moby/moby/pull/35805)
* Follow containerd namespace conventions [moby/moby#35812](https://github.com/moby/moby/pull/35812)

### Swarm Mode

+ Added support for swarm service isolation mode [moby/moby#34424](https://github.com/moby/moby/pull/34424)
- Fix task clean up for tasks that are complete [docker/swarmkit#2477](https://github.com/docker/swarmkit/pull/2477)

### Packaging

+ Add Packaging for Fedora 27 [docker/docker-ce-packaging#59](https://github.com/docker/docker-ce-packaging/pull/59)
* Change default versioning scheme to 0.0.0-dev unless specified for packaging [docker/docker-ce-packaging#67](https://github.com/docker/docker-ce-packaging/pull/67)
* Pass Version to engine static builds [docker/docker-ce-packaging#70](https://github.com/docker/docker-ce-packaging/pull/70)
+ Added support for aarch64 on Debian (stretch/jessie) and Ubuntu Zesty or newer [docker/docker-ce-packaging#35](https://github.com/docker/docker-ce-packaging/pull/35)

## 17.09.1-ce (2017-12-07)

### Builder

- Fix config leakage on shared parent stage [moby/moby#33753](https://github.com/moby/moby/issues/33753)
- Warn on empty continuation lines only, not on comment-only lines [moby/moby#35004](https://github.com/moby/moby/pull/35004)

### Client

- Set API version on Client even when Ping fails [docker/cli#546](https://github.com/docker/cli/pull/546)

### Networking

- Overlay fix for transient IP reuse [docker/libnetwork#2016](https://github.com/docker/libnetwork/pull/2016)
- Fix reapTime logic in NetworkDB and handle DNS cleanup for attachable container [docker/libnetwork#2017](https://github.com/docker/libnetwork/pull/2017)
- Disable hostname lookup on chain exists check [docker/libnetwork#2019](https://github.com/docker/libnetwork/pull/2019)
- Fix lint issues [docker/libnetwork#2020](https://github.com/docker/libnetwork/pull/2020)
- Restore error type in FindNetwork [moby/moby#35634](https://github.com/moby/moby/pull/35634)

### Runtime

- Protect `health monitor` Go channel [moby/moby#35482](https://github.com/moby/moby/pull/35482)
- Fix leaking container/exec state [moby/moby#35484](https://github.com/moby/moby/pull/35484)
- Add /proc/scsi to masked paths (patch to work around [CVE-2017-16539](http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2017-16539)) [moby/moby/#35399](https://github.com/moby/moby/pull/35399)
- Vendor tar-split: fix to prevent memory exhaustion issue that could crash Docker daemon [moby/moby/#35424](https://github.com/moby/moby/pull/35424) Fixes [CVE-2017-14992](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2017-14992)
- Fix P/Z HubPullSuite tests  [moby/moby#34837](https://github.com/moby/moby/pull/34837)
+ Windows: Add support for version filtering on pull [moby/moby#35090](https://github.com/moby/moby/pull/35090)
- Windows: Stop filtering Windows manifest lists by version [moby/moby#35117](https://github.com/moby/moby/pull/35117)
- Use rslave instead of rprivate in chroot archive [moby/moby/#35217](https://github.com/moby/moby/pull/35217)
- Remove container rootfs mountPath after unmount [moby/moby#34573](https://github.com/moby/moby/pull/34573)
- Fix honoring tmpfs size of user /dev/shm mount [moby/moby#35316](https://github.com/moby/moby/pull/35316)
- Don't abort when setting may_detach_mounts (log the error instead)  [moby/moby#35172](https://github.com/moby/moby/pull/35172)
- Fix version comparison when negotiating the API version [moby/moby#35008](https://github.com/moby/moby/pull/35008)

### Swarm mode

* Increase gRPC request timeout when sending snapshots [docker/swarmkit#2404](https://github.com/docker/swarmkit/pull/2404)
- Fix node filtering when there is no log driver [docker/swarmkit#2442](https://github.com/docker/swarmkit/pull/2442)
- Add an error on attempt to change cluster name [docker/swarmkit/#2454](https://github.com/docker/swarmkit/pull/2454)
- Delete node attachments when node is removed [docker/swarmkit/#2456](https://github.com/docker/swarmkit/pull/2456)
- Provide custom gRPC dialer to override default proxy dialer [docker/swarmkit/#2457](https://github.com/docker/swarmkit/pull/2457)
- Avoids recursive readlock on swarm info [moby/moby#35388](https://github.com/moby/moby/pull/35388)

## 17.09.0-ce (2017-09-26)

### Builder

+ Add `--chown` flag to `ADD/COPY` commands in Dockerfile [moby/moby#34263](https://github.com/moby/moby/pull/34263)
* Fix cloning unneeded files while building from git repositories [moby/moby#33704](https://github.com/moby/moby/pull/33704)

### Client

* Allow extension fields in the v3.4 version of the compose format [docker/cli#452](https://github.com/docker/cli/pull/452)
* Make compose file allow to specify names for non-external volume [docker/cli#306](https://github.com/docker/cli/pull/306)
* Support `--compose-file -` as stdin [docker/cli#347](https://github.com/docker/cli/pull/347)
* Support `start_period` for healthcheck in Docker Compose [docker/cli#475](https://github.com/docker/cli/pull/475)
+ Add support for `stop-signal` in docker stack commands [docker/cli#388](https://github.com/docker/cli/pull/388)
+ Add support for update order in compose deployments [docker/cli#360](https://github.com/docker/cli/pull/360)
+ Add ulimits to unsupported compose fields [docker/cli#482](https://github.com/docker/cli/pull/482)
+ Add `--format` to `docker-search` [docker/cli#440](https://github.com/docker/cli/pull/440)
* Show images digests when `{{.Digest}}` is in format [docker/cli#439](https://github.com/docker/cli/pull/439)
* Print output of `docker stack rm` on `stdout` instead of `stderr` [docker/cli#491](https://github.com/docker/cli/pull/491)
- Fix `docker history --format {{json .}}'` printing human-readable timestamps instead of ISO8601 when `--human=true` [docker/cli#438](https://github.com/docker/cli/pull/438)
- Fix idempotence of `docker stack deploy` when secrets or configs are used [docker/cli#509](https://github.com/docker/cli/pull/509)
- Fix presentation of random host ports [docker/cli#404](https://github.com/docker/cli/pull/404)
- Fix redundant service restarts when service created with multiple secrets [moby/moby#34746](https://github.com/moby/moby/issues/34746)

### Logging

- Fix Splunk logger not transmitting log data when tag is empty and raw-mode is used [moby/moby#34520](https://github.com/moby/moby/pull/34520)

### Networking

+ Add the control plane MTU option in the daemon config [moby/moby#34103](https://github.com/moby/moby/pull/34103)
+ Add service virtual IP to sandbox's loopback address [docker/libnetwork#1877](https://github.com/docker/libnetwork/pull/1877)

### Runtime

* Graphdriver: promote overlay2 over aufs [moby/moby#34430](https://github.com/moby/moby/pull/34430)
* LCOW: Additional flags for VHD boot [moby/moby#34451](https://github.com/moby/moby/pull/34451)
* LCOW: Don't block export [moby/moby#34448](https://github.com/moby/moby/pull/34448)
* LCOW: Dynamic sandbox management [moby/moby#34170](https://github.com/moby/moby/pull/34170)
* LCOW: Force Hyper-V Isolation [moby/moby#34468](https://github.com/moby/moby/pull/34468)
* LCOW: Move toolsScratchPath to /tmp [moby/moby#34396](https://github.com/moby/moby/pull/34396)
* LCOW: Remove hard-coding [moby/moby#34398](https://github.com/moby/moby/pull/34398)
* LCOW: WORKDIR correct handling [moby/moby#34405](https://github.com/moby/moby/pull/34405)
* Windows: named pipe mounts [moby/moby#33852](https://github.com/moby/moby/pull/33852)
- Fix "permission denied" errors when accessing volume with SELinux enforcing mode [moby/moby#34684](https://github.com/moby/moby/pull/34684)
- Fix layers size reported as `0` in `docker system df` [moby/moby#34826](https://github.com/moby/moby/pull/34826)
- Fix some "device or resource busy" errors when removing containers on RHEL 7.4 based kernels [moby/moby#34886](https://github.com/moby/moby/pull/34886)

### Swarm mode

* Include whether the managers in the swarm are autolocked as part of `docker info` [docker/cli#471](https://github.com/docker/cli/pull/471)
+ Add 'docker service rollback' subcommand [docker/cli#205](https://github.com/docker/cli/pull/205)
- Fix managers failing to join if the gRPC snapshot is larger than 4MB [docker/swarmkit#2375](https://github.com/docker/swarmkit/pull/2375)
- Fix "permission denied" errors for configuration file in SELinux-enabled containers [moby/moby#34732](https://github.com/moby/moby/pull/34732)
- Fix services failing to deploy on ARM nodes [moby/moby#34021](https://github.com/moby/moby/pull/34021)

### Packaging

+ Build scripts for ppc64el on Ubuntu [docker/docker-ce-packaging#43](https://github.com/docker/docker-ce-packaging/pull/43)

### Deprecation

+ Remove deprecated `--enable-api-cors` daemon flag [moby/moby#34821](https://github.com/moby/moby/pull/34821)

## 17.06.2-ce (2017-09-05)

### Client

- Enable TCP keepalive in the client to prevent loss of connection [docker/cli#415](https://github.com/docker/cli/pull/415)

### Runtime

- Devmapper: ensure UdevWait is called after calls to setCookie [moby/moby#33732](https://github.com/moby/moby/pull/33732)
- Aufs: ensure diff layers are correctly removed to prevent leftover files from using up storage [moby/moby#34587](https://github.com/moby/moby/pull/34587)

### Swarm mode

- Ignore PullOptions for running tasks [docker/swarmkit#2351](https://github.com/docker/swarmkit/pull/2351)

## 17.06.1-ce (2017-08-15)

### Builder

* Fix a regression, where `ADD` from remote URL's extracted archives [#89](https://github.com/docker/docker-ce/pull/89)
* Fix handling of remote "git@" notation [#100](https://github.com/docker/docker-ce/pull/100)
* Fix copy `--from` conflict with force pull [#86](https://github.com/docker/docker-ce/pull/86)

### Client

* Make pruning volumes optional when running `docker system prune`, and add a `--volumes` flag [#109](https://github.com/docker/docker-ce/pull/109)
* Show progress of replicated tasks before they are assigned [#97](https://github.com/docker/docker-ce/pull/97)
* Fix `docker wait` hanging if the container does not exist [#106](https://github.com/docker/docker-ce/pull/106)
* If `docker swarm ca` is called without the `--rotate` flag, warn if other flags are passed [#110](https://github.com/docker/docker-ce/pull/110)
* Fix API version negotiation not working if the daemon returns an error [#115](https://github.com/docker/docker-ce/pull/115)
* Print an error if "until" filter is combined with "--volumes" on system prune [#154](https://github.com/docker/docker-ce/pull/154)

### Logging

* Fix stderr logging for `journald` and `syslog` [#95](https://github.com/docker/docker-ce/pull/95)
* Fix log readers can block writes indefinitely [#98](https://github.com/docker/docker-ce/pull/98)
* Fix `awslogs` driver repeating last event [#151](https://github.com/docker/docker-ce/pull/151)

### Networking

* Fix issue with driver options not received by network drivers [#127](https://github.com/docker/docker-ce/pull/127)

### Plugins

* Make plugin removes more resilient to failure [#91](https://github.com/docker/docker-ce/pull/91)

### Runtime

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

* Redact secret data on secret creation [#99](https://github.com/docker/docker-ce/pull/99)

### Swarm mode

* Do not add duplicate platform information to service spec [#107](https://github.com/docker/docker-ce/pull/107)
* Cluster update and memory issue fixes [#114](https://github.com/docker/docker-ce/pull/114)
* Changing get network request to return predefined network in swarm [#150](https://github.com/docker/docker-ce/pull/150)

## 17.06.0-ce (2017-06-28)

> **Note**: Docker 17.06.0 has an issue in the image builder causing a change in the behavior
> of the `ADD` instruction of Dockerfile when referencing a remote `.tar.gz` file. The issue will be
> fixed in Docker 17.06.1.

> **Note**: Starting with Docker CE 17.06, Ubuntu packages are also available
> for IBM Z using the s390x architecture.

> **Note**: Docker 17.06 by default disables communication with legacy (v1)
> registries. If you require interaction with registries that have not yet
> migrated to the v2 protocol, set the `--disable-legacy-registry=false` daemon
> option. Interaction with v1 registries will be removed in Docker 17.12.

### Builder

+ Add `--iidfile` option to docker build. It allows specifying a location where to save the resulting image ID
+ Allow specifying any remote ref in git checkout URLs [#32502](https://github.com/moby/moby/pull/32502)

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

### Distribution

* Select digest over tag when both are provided during a pull [#33214](https://github.com/moby/moby/pull/33214)

### Logging

+ Add monitored resource type metadata for GCP logging driver [#32930](https://github.com/moby/moby/pull/32930)
+ Add multiline processing to the AWS CloudWatch logs driver [#30891](https://github.com/moby/moby/pull/30891)

### Networking

+ Add Support swarm-mode services with node-local networks such as macvlan, ipvlan, bridge, host [#32981](https://github.com/moby/moby/pull/32981)
+ Pass driver-options to network drivers on service creation [#32981](https://github.com/moby/moby/pull/33130)
+ Isolate Swarm Control-plane traffic from Application data traffic using --data-path-addr [#32717](https://github.com/moby/moby/pull/32717)
* Several improvments to Service Discovery [#docker/libnetwork/1796](https://github.com/docker/libnetwork/pull/1796)

### Packaging

+ Rely on `container-selinux` on Centos/Fedora/RHEL when available [#32437](https://github.com/moby/moby/pull/32437)

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
- Prevent a containerd crash when journald is restarted [#containerd/930](https://github.com/containerd/containerd/pull/930)
- Fix healthcheck failures due to invalid environment variables [#33249](https://github.com/moby/moby/pull/33249)
- Prevent a directory to be created in lieu of the daemon socket when a container mounting it is to be restarted during a shutdown [#30348](https://github.com/moby/moby/pull/33330)
- Prevent a container to be restarted upon stop if its stop signal is set to `SIGKILL` [#33335](https://github.com/moby/moby/pull/33335)
- Ensure log drivers get passed the same filename to both StartLogging and StopLogging endpoints [#33583](https://github.com/moby/moby/pull/33583)
- Remove daemon data structure dump on `SIGUSR1` to avoid a panic [#33598](https://github.com/moby/moby/pull/33598)

### Security

+ Allow personality with UNAME26 bit set in default seccomp profile [#32965](https://github.com/moby/moby/pull/32965)

### Swarm Mode

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

### Deprecation

* Disable legacy registry (v1) by default [#33629](https://github.com/moby/moby/pull/33629)

## 17.03.2-ce (2017-05-29)

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

## 17.03.1-ce (2017-03-27)

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

* Update swarmkit to 17756457ad6dc4d8a639a1f0b7a85d1b65a617bb [#31807](https://github.com/docker/docker/pull/31807)
 * Scheduler now correctly considers tasks which have been assigned to a node but aren't yet running [docker/swarmkit#1980](https://github.com/docker/swarmkit/pull/1980)
 * Allow removal of a network when only dead tasks reference it [docker/swarmkit#2018](https://github.com/docker/swarmkit/pull/2018)
 * Retry failed network allocations less aggressively [docker/swarmkit#2021](https://github.com/docker/swarmkit/pull/2021)
 * Avoid network allocation for tasks that are no longer running [docker/swarmkit#2017](https://github.com/docker/swarmkit/pull/2017)
 * Bookkeeping fixes inside network allocator allocator [docker/swarmkit#2019](https://github.com/docker/swarmkit/pull/2019) [docker/swarmkit#2020](https://github.com/docker/swarmkit/pull/2020)

### Windows

* Cleanup HCS on restore [#31503](https://github.com/docker/docker/pull/31503)

## 17.03.0-ce (2017-03-01)

**IMPORTANT**: Starting with this release, Docker is on a monthly release cycle and uses a
new YY.MM versioning scheme to reflect this. Two channels are available: monthly and quarterly.
Any given monthly release will only receive security and bugfixes until the next monthly
release is available. Quarterly releases receive security and bugfixes for 4 months after
initial release. This release includes bugfixes for 1.13.1 but
there are no major feature additions and the API version stays the same.
Upgrading from Docker 1.13.1 to 17.03.0 is expected to be simple and low-risk.

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

# Edge releases

## 18.04.0-ce (2018-04-10)

### Builder

- Fix typos in builder and client. [moby/moby#36424](https://github.com/moby/moby/pull/36424)

### Client

* Print Stack API and Kubernetes versions in version command. [docker/cli#898](https://github.com/docker/cli/pull/898)
- Fix Kubernetes duplication in version command. [docker/cli#953](https://github.com/docker/cli/pull/953)
* Use HasAvailableFlags instead of HasFlags for Options in help. [docker/cli#959](https://github.com/docker/cli/pull/959)
+ Add support for mandatory variables to stack deploy. [docker/cli#893](https://github.com/docker/cli/pull/893)
- Fix docker stack services command Port output. [docker/cli#943](https://github.com/docker/cli/pull/943)
* Deprecate unencrypted storage. [docker/cli#561](https://github.com/docker/cli/pull/561)
* Don't set a default filename for ConfigFile. [docker/cli#917](https://github.com/docker/cli/pull/917)
- Fix compose network name. [docker/cli#941](https://github.com/docker/cli/pull/941)

### Logging

* Silent login: use credentials from cred store to login. [docker/cli#139](https://github.com/docker/cli/pull/139)
+ Add support for compressibility of log file. [moby/moby#29932](https://github.com/moby/moby/pull/29932)
- Fix empty LogPath with non-blocking logging mode. [moby/moby#36272](https://github.com/moby/moby/pull/36272)

### Networking

- Prevent explicit removal of ingress network. [moby/moby#36538](https://github.com/moby/moby/pull/36538)

### Runtime

* Devmapper cleanup improvements. [moby/moby#36307](https://github.com/moby/moby/pull/36307)
* Devmapper.Mounted: remove. [moby/moby#36437](https://github.com/moby/moby/pull/36437)
* Devmapper/Remove(): use Rmdir, ignore errors. [moby/moby#36438](https://github.com/moby/moby/pull/36438)
* LCOW - Change platform parser directive to FROM statement flag. [moby/moby#35089](https://github.com/moby/moby/pull/35089)
* Split daemon service code to windows file. [moby/moby#36653](https://github.com/moby/moby/pull/36653)
* Windows: Block pulling uplevel images. [moby/moby#36327](https://github.com/moby/moby/pull/36327)
* Windows: Hyper-V containers are broken after 36586 was merged. [moby/moby#36610](https://github.com/moby/moby/pull/36610)
* Windows: Move kernel_windows to use golang registry functions. [moby/moby#36617](https://github.com/moby/moby/pull/36617)
* Windows: Pass back system errors on container exit. [moby/moby#35967](https://github.com/moby/moby/pull/35967)
* Windows: Remove servicing mode. [moby/moby#36267](https://github.com/moby/moby/pull/36267)
* Windows: Report Version and UBR. [moby/moby#36451](https://github.com/moby/moby/pull/36451)
* Bump Runc to 1.0.0-rc5. [moby/moby#36449](https://github.com/moby/moby/pull/36449)
* Mount failure indicates the path that failed. [moby/moby#36407](https://github.com/moby/moby/pull/36407)
* Change return for errdefs.getImplementer(). [moby/moby#36489](https://github.com/moby/moby/pull/36489)
* Client: fix hijackedconn reading from buffer. [moby/moby#36663](https://github.com/moby/moby/pull/36663)
* Content encoding negotiation added to archive request. [moby/moby#36164](https://github.com/moby/moby/pull/36164)
* Daemon/stats: more resilient cpu sampling. [moby/moby#36519](https://github.com/moby/moby/pull/36519)
* Daemon/stats: remove obnoxious types file. [moby/moby#36494](https://github.com/moby/moby/pull/36494)
* Daemon: use context error rather than inventing new one. [moby/moby#36670](https://github.com/moby/moby/pull/36670)
* Enable CRIU on non-amd64 architectures (v2). [moby/moby#36676](https://github.com/moby/moby/pull/36676)
- Fixes intermittent client hang after closing stdin to attached container [moby/moby#36517](https://github.com/moby/moby/pull/36517)
- Fix daemon panic on container export after restart [moby/moby#36586](https://github.com/moby/moby/pull/36586)
- Follow-up fixes on multi-stage moby's Dockerfile. [moby/moby#36425](https://github.com/moby/moby/pull/36425)
* Freeze busybox and latest glibc in Docker image. [moby/moby#36375](https://github.com/moby/moby/pull/36375)
* If container will run as non root user, drop permitted, effective caps early. [moby/moby#36587](https://github.com/moby/moby/pull/36587)
* Layer: remove metadata store interface. [moby/moby#36504](https://github.com/moby/moby/pull/36504)
* Minor optimizations to dockerd. [moby/moby#36577](https://github.com/moby/moby/pull/36577)
* Whitelist statx syscall. [moby/moby#36417](https://github.com/moby/moby/pull/36417)
+ Add missing error return for plugin creation. [moby/moby#36646](https://github.com/moby/moby/pull/36646)
- Fix AppArmor not being applied to Exec processes. [moby/moby#36466](https://github.com/moby/moby/pull/36466)
* Daemon/logger/ring.go: log error not instance. [moby/moby#36475](https://github.com/moby/moby/pull/36475)
- Fix stats collector spinning CPU if no stats are collected. [moby/moby#36609](https://github.com/moby/moby/pull/36609)
- Fix(distribution): digest cache should not be moved if it was an auth. [moby/moby#36509](https://github.com/moby/moby/pull/36509)
- Make sure plugin container is removed on failure. [moby/moby#36715](https://github.com/moby/moby/pull/36715)
* Bump to containerd 1.0.3. [moby/moby#36749](https://github.com/moby/moby/pull/36749)
* Don't sort plugin mount slice. [moby/moby#36711](https://github.com/moby/moby/pull/36711)

### Swarm Mode

* Fixes for synchronizing the dispatcher shutdown with in-progress rpcs. [moby/moby#36371](https://github.com/moby/moby/pull/36371)
* Increase raft ElectionTick to 10xHeartbeatTick. [moby/moby#36672](https://github.com/moby/moby/pull/36672)
* Make Swarm manager Raft quorum parameters configurable in daemon config. [moby/moby#36726](https://github.com/moby/moby/pull/36726)
* Ingress network should not be attachable. [docker/swarmkit#2523](https://github.com/docker/swarmkit/pull/2523)
* [manager/state] Add fernet as an option for raft encryption. [docker/swarmkit#2535](https://github.com/docker/swarmkit/pull/2535)
* Log GRPC server errors.  [docker/swarmkit#2541](https://github.com/docker/swarmkit/pull/2541)
* Log leadership changes at the manager level. [docker/swarmkit#2542](https://github.com/docker/swarmkit/pull/2542)
* Remove the containerd executor. [docker/swarmkit#2568](https://github.com/docker/swarmkit/pull/2568)
* Agent: backoff session when no remotes are available. [docker/swarmkit#2570](https://github.com/docker/swarmkit/pull/2570)
* [ca/manager] Remove root CA key encryption support entirely. [docker/swarmkit#2573](https://github.com/docker/swarmkit/pull/2573)
- Fix agent logging race. [docker/swarmkit#2578](https://github.com/docker/swarmkit/pull/2578)
* Adding logic to restore networks in order. [docker/swarmkit#2571](https://github.com/docker/swarmkit/pull/2571)

## 18.02.0-ce (2018-02-07)

### Builder

- Gitutils: fix checking out submodules [moby/moby#35737](https://github.com/moby/moby/pull/35737)

### Client

* Attach: Ensure attach exit code matches container's [docker/cli#696](https://github.com/docker/cli/pull/696)
+ Added support for tmpfs-mode in compose file [docker/cli#808](https://github.com/docker/cli/pull/808)
+ Adds a new compose file version 3.6 [docker/cli#808](https://github.com/docker/cli/pull/808)
- Fix issue of filter in `docker ps` where `health=starting` returns nothing [moby/moby#35940](https://github.com/moby/moby/pull/35940)
+ Improve presentation of published port ranges [docker/cli#581](https://github.com/docker/cli/pull/581)
* Bump Go to 1.9.3 [docker/cli#827](https://github.com/docker/cli/pull/827)
- Fix broken Kubernetes stack flags [docker/cli#831](https://github.com/docker/cli/pull/831)
* Annotate "stack" commands to be "swarm" and "kubernetes" [docker/cli#804](https://github.com/docker/cli/pull/804)

### Experimental

+ Add manifest command [docker/cli#138](https://github.com/docker/cli/pull/138)
* LCOW remotefs - return error in Read() implementation [moby/moby#36051](https://github.com/moby/moby/pull/36051)
+ LCOW: Coalesce daemon stores, allow dual LCOW and WCOW mode [moby/moby#34859](https://github.com/moby/moby/pull/34859)
- LCOW: Fix OpenFile parameters [moby/moby#36043](https://github.com/moby/moby/pull/36043)
* LCOW: Raise minimum requirement to Windows RS3 RTM build (16299) [moby/moby#36065](https://github.com/moby/moby/pull/36065)

### Logging

* Improve daemon config reload; log active configuration [moby/moby#36019](https://github.com/moby/moby/pull/36019)
- Fixed error detection using IsErrNotFound and IsErrNotImplemented for the ContainerLogs method [moby/moby#36000](https://github.com/moby/moby/pull/36000)
+ Add journald tag as SYSLOG_IDENTIFIER [moby/moby#35570](https://github.com/moby/moby/pull/35570)
* Splunk: limit the reader size on error responses [moby/moby#35509](https://github.com/moby/moby/pull/35509)

### Networking

* Disable service on release network results in zero-downtime deployments with rolling upgrades [moby/moby#35960](https://github.com/moby/moby/pull/35960)
- Fix services failing to start if multiple networks with the same name exist in different spaces [moby/moby#30897](https://github.com/moby/moby/pull/30897)
- Fix duplicate networks being added with `docker service update --network-add` [docker/cli#780](https://github.com/docker/cli/pull/780)
- Fixing ingress network when upgrading from 17.09 to 17.12. [moby/moby#36003](https://github.com/moby/moby/pull/36003)
- Fix ndots configuration [docker/libnetwork#1995](https://github.com/docker/libnetwork/pull/1995)
- Fix IPV6 networking being deconfigured if live-restore is enabled [docker/libnetwork#2043](https://github.com/docker/libnetwork/pull/2043)
+ Add support for MX type DNS queries in the embedded DNS server [docker/libnetwork#2041](https://github.com/docker/libnetwork/pull/2041)

### Packaging

+ Added packaging for Fedora 26, Fedora 27, and Centos 7 on aarch64 [docker/docker-ce-packaging#71](https://github.com/docker/docker-ce-packaging/pull/71)
- Removed support for Ubuntu Zesty [docker/docker-ce-packaging#73](https://github.com/docker/docker-ce-packaging/pull/73)
- Removed support for Fedora 25 [docker/docker-ce-packaging#72](https://github.com/docker/docker-ce-packaging/pull/72)

### Runtime

- Fixes unexpected Docker Daemon shutdown based on pipe error [moby/moby#35968](https://github.com/moby/moby/pull/35968)
- Fix some occurrences of hcsshim::ImportLayer failed in Win32: The system cannot find the path specified [moby/moby#35924](https://github.com/moby/moby/pull/35924)
* Windows: increase the maximum layer size during build to 127GB [moby/moby#35925](https://github.com/moby/moby/pull/35925)
- Fix Devicemapper: Error running DeleteDevice dm_task_run failed [moby/moby#35919](https://github.com/moby/moby/pull/35919)
+ Introduce « exec_die » event [moby/moby#35744](https://github.com/moby/moby/pull/35744)
* Update API to version 1.36 [moby/moby#35744](https://github.com/moby/moby/pull/35744)
- Fix `docker update` not updating cpu quota, and cpu-period of a running container [moby/moby#36030](https://github.com/moby/moby/pull/36030)
* Make container shm parent unbindable [moby/moby#35830](https://github.com/moby/moby/pull/35830)
+ Make image (layer) downloads faster by using pigz [moby/moby#35697](https://github.com/moby/moby/pull/35697)
+ Protect the daemon from volume plugins that are slow or deadlocked [moby/moby#35441](https://github.com/moby/moby/pull/35441)
- Fix `DOCKER_RAMDISK` environment variable not being honoured [moby/moby#35957](https://github.com/moby/moby/pull/35957)
* Bump containerd to 1.0.1 (9b55aab90508bd389d7654c4baf173a981477d55) [moby/moby#35986](https://github.com/moby/moby/pull/35986)
* Update runc to fix hang during start and exec [moby/moby#36097](https://github.com/moby/moby/pull/36097)
- Fix "--node-generic-resource" singular/plural [moby/moby#36125](https://github.com/moby/moby/pull/36125)

## 18.01.0-ce (2018-01-10)

### Builder

* Fix files not being deleted if user-namespaces are enabled [moby/moby#35822](https://github.com/moby/moby/pull/35822)
- Add support for expanding environment-variables in `docker commit --change ...` [moby/moby#35582](https://github.com/moby/moby/pull/35582)

### Client

* Return errors from client in stack deploy configs [docker/cli#757](https://github.com/docker/cli/pull/757)
- Fix description of filter flag in prune commands [docker/cli#774](https://github.com/docker/cli/pull/774)
+ Add "pid" to unsupported options list [docker/cli#768](https://github.com/docker/cli/pull/768)
+ Add support for experimental Cli configuration [docker/cli#758](https://github.com/docker/cli/pull/758)
+ Add support for generic resources to bash completion [docker/cli#749](https://github.com/docker/cli/pull/749)
- Fix error in zsh completion script for docker exec [docker/cli#751](https://github.com/docker/cli/pull/751)
+ Add a debug message when client closes websocket attach connection [moby/moby#35720](https://github.com/moby/moby/pull/35720)
- Fix bash completion for `"docker swarm"` [docker/cli#772](https://github.com/docker/cli/pull/772)


### Documentation
* Correct references to `--publish` long syntax in docs [docker/cli#746](https://github.com/docker/cli/pull/746)
* Corrected descriptions for MAC_ADMIN and MAC_OVERRIDE [docker/cli#761](https://github.com/docker/cli/pull/761)
* Updated developer doc to explain external CLI [moby/moby#35681](https://github.com/moby/moby/pull/35681)
- Fix `"on-failure"` restart policy being documented as "failure" [docker/cli#754](https://github.com/docker/cli/pull/754)
- Fix anchors to "Storage driver options" [docker/cli#748](https://github.com/docker/cli/pull/748)

### Experimental

+ Add kubernetes support to `docker stack` command [docker/cli#721](https://github.com/docker/cli/pull/721)
* Don't append the container id to custom directory checkpoints. [moby/moby#35694](https://github.com/moby/moby/pull/35694)

### Logging

* Fix daemon crash when using the GELF log driver over TCP when the GELF server goes down [moby/moby#35765](https://github.com/moby/moby/pull/35765)
- Fix awslogs batch size calculation for large logs [moby/moby#35726](https://github.com/moby/moby/pull/35726)

### Networking

- Windows: Fix to allow docker service to start on Windows VM [docker/libnetwork#1916](https://github.com/docker/libnetwork/pull/1916)
- Fix for docker intercepting DNS requests on ICS network [docker/libnetwork#2014](https://github.com/docker/libnetwork/pull/2014)
+ Windows: Added a new network creation driver option [docker/libnetwork#2021](https://github.com/docker/libnetwork/pull/2021)


### Runtime

* Validate Mount-specs on container start to prevent missing host-path [moby/moby#35833](https://github.com/moby/moby/pull/35833)
- Fix overlay2 storage driver inside a user namespace [moby/moby#35794](https://github.com/moby/moby/pull/35794)
* Zfs: fix busy error on container stop [moby/moby#35674](https://github.com/moby/moby/pull/35674)
- Fix health checks not using the container's working directory [moby/moby#35845](https://github.com/moby/moby/pull/35845)
- Fix VFS graph driver failure to initialize because of failure to setup fs quota [moby/moby#35827](https://github.com/moby/moby/pull/35827)
- Fix containerd events being processed twice [moby/moby#35896](https://github.com/moby/moby/pull/35896)

### Swarm mode

- Fix published ports not being updated if a service has the same number of host-mode published ports with Published Port 0 [docker/swarmkit#2376](https://github.com/docker/swarmkit/pull/2376)
* Make the task termination order deterministic [docker/swarmkit#2265](https://github.com/docker/swarmkit/pull/2265)

## 17.11.0-ce (2017-11-20)

> **Important**: Docker CE 17.11 is the first Docker release based on
[containerd 1.0 beta](https://github.com/containerd/containerd/releases/tag/v1.0.0-beta.2).
Docker CE 17.11 and later don't recognize containers started with
previous Docker versions. If using
[Live Restore](https://docs.docker.com/engine/admin/live-restore/#enable-the-live-restore-option),
you must stop all containers before upgrading to Docker CE 17.11.
If you don't, any containers started by Docker versions that predate
17.11 aren't recognized by Docker after the upgrade and keep
running, un-managed, on the system.
{:.important}

### Builder

* Test & Fix build with rm/force-rm matrix [moby/moby#35139](https://github.com/moby/moby/pull/35139)
- Fix build with `--stream` with a large context [moby/moby#35404](https://github.com/moby/moby/pull/35404)

### Client

* Hide help flag from help output [docker/cli#645](https://github.com/docker/cli/pull/645)
* Support parsing of named pipes for compose volumes [docker/cli#560](https://github.com/docker/cli/pull/560)
* [Compose] Cast values to expected type after interpolating values [docker/cli#601](https://github.com/docker/cli/pull/601)
+ Add output for "secrets" and "configs" on `docker stack deploy` [docker/cli#593](https://github.com/docker/cli/pull/593)
- Fix flag description for `--host-add` [docker/cli#648](https://github.com/docker/cli/pull/648)
* Do not truncate ID on docker service ps --quiet [docker/cli#579](https://github.com/docker/cli/pull/579)

### Deprecation

* Update bash completion and deprecation for synchronous service updates [docker/cli#610](https://github.com/docker/cli/pull/610)

### Logging

* copy to log driver's bufsize, fixes #34887 [moby/moby#34888](https://github.com/moby/moby/pull/34888)
+ Add TCP support for GELF log driver [moby/moby#34758](https://github.com/moby/moby/pull/34758)
+ Add credentials endpoint option for awslogs driver [moby/moby#35055](https://github.com/moby/moby/pull/35055)

### Networking

- Fix network name masking network ID on delete [moby/moby#34509](https://github.com/moby/moby/pull/34509)
- Fix returned error code for network creation from 500 to 409 [moby/moby#35030](https://github.com/moby/moby/pull/35030)
- Fix tasks fail with error "Unable to complete atomic operation, key modified" [docker/libnetwork#2004](https://github.com/docker/libnetwork/pull/2004)

### Runtime

* Switch to Containerd 1.0 client [moby/moby#34895](https://github.com/moby/moby/pull/34895)
* Increase container default shutdown timeout on Windows [moby/moby#35184](https://github.com/moby/moby/pull/35184)
* LCOW: API: Add `platform` to /images/create and /build [moby/moby#34642](https://github.com/moby/moby/pull/34642)
* Stop filtering Windows manifest lists by version [moby/moby#35117](https://github.com/moby/moby/pull/35117)
* Use windows console mode constants from Azure/go-ansiterm [moby/moby#35056](https://github.com/moby/moby/pull/35056)
* Windows Daemon should respect DOCKER_TMPDIR [moby/moby#35077](https://github.com/moby/moby/pull/35077)
* Windows: Fix startup logging [moby/moby#35253](https://github.com/moby/moby/pull/35253)
+ Add support for Windows version filtering on pull [moby/moby#35090](https://github.com/moby/moby/pull/35090)
- Fixes LCOW after containerd 1.0 introduced regressions [moby/moby#35320](https://github.com/moby/moby/pull/35320)
* ContainerWait on remove: don't stuck on rm fail [moby/moby#34999](https://github.com/moby/moby/pull/34999)
* oci: obey CL_UNPRIVILEGED for user namespaced daemon [moby/moby#35205](https://github.com/moby/moby/pull/35205)
* Don't abort when setting may_detach_mounts [moby/moby#35172](https://github.com/moby/moby/pull/35172)
- Fix panic on get container pid when live restore containers [moby/moby#35157](https://github.com/moby/moby/pull/35157)
- Mask `/proc/scsi` path for containers to prevent removal of devices (CVE-2017-16539) [moby/moby#35399](https://github.com/moby/moby/pull/35399)
* Update to github.com/vbatts/tar-split@v0.10.2 (CVE-2017-14992) [moby/moby#35424](https://github.com/moby/moby/pull/35424)

### Swarm Mode

* Modifying integration test due to new ipam options in swarmkit [moby/moby#35103](https://github.com/moby/moby/pull/35103)
- Fix deadlock on getting swarm info [moby/moby#35388](https://github.com/moby/moby/pull/35388)
+ Expand the scope of the `Err` field in `TaskStatus` to also cover non-terminal errors that block the task from progressing [docker/swarmkit#2287](https://github.com/docker/swarmkit/pull/2287)

### Packaging

+ Build packages for Debian 10 (Buster) [docker/docker-ce-packaging#50](https://github.com/docker/docker-ce-packaging/pull/50)
+ Build packages for Ubuntu 17.10 (Artful) [docker/docker-ce-packaging#55](https://github.com/docker/docker-ce-packaging/pull/55)

## 17.10.0-ce (2017-10-17)

> **Important**: Starting with this release, `docker service create`, `docker service update`,
`docker service scale` and `docker service rollback` use non-detached mode as default,
use `--detach` to keep the old behaviour.
{: .important }

### Builder

* Reset uid/gid to 0 in uploaded build context to share build cache with other clients [docker/cli#513](https://github.com/docker/cli/pull/513)
+ Add support for `ADD` urls without any sub path [moby/moby#34217](https://github.com/moby/moby/pull/34217)

### Client

* Move output of `docker stack rm` to stdout [docker/cli#491](https://github.com/docker/cli/pull/491)
* Use natural sort for secrets and configs in cli [docker/cli#307](https://github.com/docker/cli/pull/307)
* Use non-detached mode as default for `docker service` commands [docker/cli#525](https://github.com/docker/cli/pull/525)
* Set APIVersion on the client, even when Ping fails [docker/cli#546](https://github.com/docker/cli/pull/546)
- Fix loader error with different build syntax in `docker stack deploy` [docker/cli#544](https://github.com/docker/cli/pull/544)
* Change the default output format for `docker container stats` to show `CONTAINER ID` and `NAME` [docker/cli#565](https://github.com/docker/cli/pull/565)
+ Add `--no-trunc` flag to `docker container stats` [docker/cli#565](https://github.com/docker/cli/pull/565)
+ Add experimental `docker trust`: `view`, `revoke`, `sign` subcommands [docker/cli#472](https://github.com/docker/cli/pull/472)
- Various doc and shell completion fixes [docker/cli#610](https://github.com/docker/cli/pull/610) [docker/cli#611](https://github.com/docker/cli/pull/611) [docker/cli#618](https://github.com/docker/cli/pull/618) [docker/cli#580](https://github.com/docker/cli/pull/580) [docker/cli#598](https://github.com/docker/cli/pull/598) [docker/cli#603](https://github.com/docker/cli/pull/603)

### Networking

* Enabling ILB/ELB on windows using per-node, per-network LB endpoint [moby/moby#34674](https://github.com/moby/moby/pull/34674)
* Overlay fix for transient IP reuse [docker/libnetwork#1935](https://github.com/docker/libnetwork/pull/1935)
* Serializing bitseq alloc [docker/libnetwork#1788](https://github.com/docker/libnetwork/pull/1788)
- Disable hostname lookup on chain exists check [docker/libnetwork#1974](https://github.com/docker/libnetwork/pull/1974)

### Runtime

* LCOW: Add UVM debuggability by grabbing logs before tear-down [moby/moby#34846](https://github.com/moby/moby/pull/34846)
* LCOW: Prepare work for bind mounts [moby/moby#34258](https://github.com/moby/moby/pull/34258)
* LCOW: Support for docker cp, ADD/COPY on build [moby/moby#34252](https://github.com/moby/moby/pull/34252)
* LCOW: VHDX boot to readonly [moby/moby#34754](https://github.com/moby/moby/pull/34754)
* Volume: evaluate symlinks before relabeling mount source [moby/moby#34792](https://github.com/moby/moby/pull/34792)
- Fixing ‘docker cp’ to allow new target file name in a host symlinked directory [moby/moby#31993](https://github.com/moby/moby/pull/31993)
+ Add support for Windows version filtering on pull [moby/moby#35090](https://github.com/moby/moby/pull/35090)

### Swarm mode

* Produce an error if `docker swarm init --force-new-cluster` is executed on worker nodes [moby/moby#34881](https://github.com/moby/moby/pull/34881)
+ Add support for `.Node.Hostname` templating in swarm services [moby/moby#34686](https://github.com/moby/moby/pull/34686)
* Increase gRPC request timeout to 20 seconds for sending snapshots [docker/swarmkit#2391](https://github.com/docker/swarmkit/pull/2391)
- Do not filter nodes if logdriver is set to `none` [docker/swarmkit#2396](https://github.com/docker/swarmkit/pull/2396)
+ Adding ipam options to ipam driver requests [docker/swarmkit#2324](https://github.com/docker/swarmkit/pull/2324)

## 17.07.0-ce (2017-08-29)

### API & Client

* Add support for proxy configuration in config.json [docker/cli#93](https://github.com/docker/cli/pull/93)
* Enable pprof/debug endpoints by default [moby/moby#32453](https://github.com/moby/moby/pull/32453)
* Passwords can now be passed using `STDIN` using the new  `--password-stdin` flag on `docker login` [docker/cli#271](https://github.com/docker/cli/pull/271)
+ Add `--detach` to docker scale [docker/cli#243](https://github.com/docker/cli/pull/243)
* Prevent `docker logs --no-stream` from hanging due to non-existing containers [moby/moby#34004](https://github.com/moby/moby/pull/34004)
- Fix `docker stack ps` printing error to `stdout` instead of `stderr` [docker/cli#298](https://github.com/docker/cli/pull/298)
* Fix progress bar being stuck on `docker service create` if an error occurs during deploy [docker/cli#259](https://github.com/docker/cli/pull/259)
* Improve presentation of progress bars in interactive mode [docker/cli#260](https://github.com/docker/cli/pull/260) [docker/cli#237](https://github.com/docker/cli/pull/237)
* Print a warning if `docker login --password` is used, and recommend `--password-stdin` [docker/cli#270](https://github.com/docker/cli/pull/270)
* Make API version negotiation more robust [moby/moby#33827](https://github.com/moby/moby/pull/33827)
* Hide `--detach` when connected to daemons older than Docker 17.05 [docker/cli#219](https://github.com/docker/cli/pull/219)
+ Add `scope` filter in `GET /networks/(id or name)` [moby/moby#33630](https://github.com/moby/moby/pull/33630)

### Builder

* Implement long running interactive session and sending build context incrementally [moby/moby#32677](https://github.com/moby/moby/pull/32677) [docker/cli#231](https://github.com/docker/cli/pull/231) [moby/moby#33859](https://github.com/moby/moby/pull/33859)
* Warn on empty continuation lines [moby/moby#33719](https://github.com/moby/moby/pull/33719)
- Fix `.dockerignore` entries with a leading `/` not matching anything [moby/moby#32088](https://github.com/moby/moby/pull/32088)

### Logging

- Fix wrong filemode for rotate log files [moby/moby#33926](https://github.com/moby/moby/pull/33926)
- Fix stderr logging for journald and syslog [moby/moby#33832](https://github.com/moby/moby/pull/33832)

### Runtime

* Allow stopping of paused container [moby/moby#34027](https://github.com/moby/moby/pull/34027)
+ Add quota support for the overlay2 storage driver [moby/moby#32977](https://github.com/moby/moby/pull/32977)
* Remove container locks on `docker ps` [moby/moby#31273](https://github.com/moby/moby/pull/31273)
* Store container names in memdb [moby/moby#33886](https://github.com/moby/moby/pull/33886)
* Fix race condition between `docker exec` and `docker pause` [moby/moby#32881](https://github.com/moby/moby/pull/32881)
* Devicemapper: Rework logging and add `--storage-opt dm.libdm_log_level` [moby/moby#33845](https://github.com/moby/moby/pull/33845)
* Devicemapper: Prevent "device in use" errors if deferred removal is enabled, but not deferred deletion [moby/moby#33877](https://github.com/moby/moby/pull/33877)
* Devicemapper: Use KeepAlive to prevent tasks being garbage-collected while still in use [moby/moby#33376](https://github.com/moby/moby/pull/33376)
* Report intermediate prune results if prune is cancelled [moby/moby#33979](https://github.com/moby/moby/pull/33979)
- Fix run `docker rename <container-id> new_name` concurrently resulting in the having multiple names [moby/moby#33940](https://github.com/moby/moby/pull/33940)
* Fix file-descriptor leak and error handling [moby/moby#33713](https://github.com/moby/moby/pull/33713)
- Fix SIGSEGV when running containers [docker/cli#303](https://github.com/docker/cli/pull/303)
* Prevent a goroutine leak when healthcheck gets stopped [moby/moby#33781](https://github.com/moby/moby/pull/33781)
* Image: Improve store locking [moby/moby#33755](https://github.com/moby/moby/pull/33755)
* Fix Btrfs quota groups not being removed when container is destroyed [moby/moby#29427](https://github.com/moby/moby/pull/29427)
* Libcontainerd: fix defunct containerd processes not being properly reaped [moby/moby#33419](https://github.com/moby/moby/pull/33419)
* Preparations for Linux Containers on Windows
  * LCOW: Dedicated scratch space for service VM utilities [moby/moby#33809](https://github.com/moby/moby/pull/33809)
  * LCOW: Support most operations excluding remote filesystem [moby/moby#33241](https://github.com/moby/moby/pull/33241) [moby/moby#33826](https://github.com/moby/moby/pull/33826)
  * LCOW: Change directory from lcow to "Linux Containers" [moby/moby#33835](https://github.com/moby/moby/pull/33835)
  * LCOW: pass command arguments without extra quoting [moby/moby#33815](https://github.com/moby/moby/pull/33815)
  * LCOW: Updates necessary due to platform schema change [moby/moby#33785](https://github.com/moby/moby/pull/33785)

### Swarm mode

* Initial support for plugable secret backends [moby/moby#34157](https://github.com/moby/moby/pull/34157) [moby/moby#34123](https://github.com/moby/moby/pull/34123)
* Sort swarm stacks and nodes using natural sorting [docker/cli#315](https://github.com/docker/cli/pull/315)
* Make engine support cluster config event [moby/moby#34032](https://github.com/moby/moby/pull/34032)
* Only pass a join address when in the process of joining a cluster [moby/moby#33361](https://github.com/moby/moby/pull/33361)
* Fix error during service creation if a network with the same name exists both as "local" and "swarm" scoped network [docker/cli#184](https://github.com/docker/cli/pull/184)
* (experimental) Add support for plugins on swarm [moby/moby#33575](https://github.com/moby/moby/pull/33575)

## 17.05.0-ce (2017-05-04)

### Builder

+ Add multi-stage build support [#31257](https://github.com/docker/docker/pull/31257) [#32063](https://github.com/docker/docker/pull/32063)
+ Allow using build-time args (`ARG`) in `FROM` [#31352](https://github.com/docker/docker/pull/31352)
+ Add an option for specifying build target [#32496](https://github.com/docker/docker/pull/32496)
* Accept `-f -` to read Dockerfile from `stdin`, but use local context for building [#31236](https://github.com/docker/docker/pull/31236)
* The values of default build time arguments (e.g `HTTP_PROXY`) are no longer displayed in docker image history unless a corresponding `ARG` instruction is written in the Dockerfile. [#31584](https://github.com/docker/docker/pull/31584)
- Fix setting command if a custom shell is used in a parent image [#32236](https://github.com/docker/docker/pull/32236)
- Fix `docker build --label` when the label includes single quotes and a space [#31750](https://github.com/docker/docker/pull/31750)

### Client

* Add `--mount` flag to `docker run` and `docker create` [#32251](https://github.com/docker/docker/pull/32251)
* Add `--type=secret` to `docker inspect` [#32124](https://github.com/docker/docker/pull/32124)
* Add `--format` option to `docker secret ls` [#31552](https://github.com/docker/docker/pull/31552)
* Add `--filter` option to `docker secret ls` [#30810](https://github.com/docker/docker/pull/30810)
* Add `--filter scope=<swarm|local>` to `docker network ls` [#31529](https://github.com/docker/docker/pull/31529)
* Add `--cpus` support to `docker update` [#31148](https://github.com/docker/docker/pull/31148)
* Add label filter to `docker system prune` and other `prune` commands [#30740](https://github.com/docker/docker/pull/30740)
* `docker stack rm` now accepts multiple stacks as input [#32110](https://github.com/docker/docker/pull/32110)
* Improve `docker version --format` option when the client has downgraded the API version [#31022](https://github.com/docker/docker/pull/31022)
* Prompt when using an encrypted client certificate to connect to a docker daemon [#31364](https://github.com/docker/docker/pull/31364)
* Display created tags on successful `docker build` [#32077](https://github.com/docker/docker/pull/32077)
* Cleanup compose convert error messages [#32087](https://github.com/moby/moby/pull/32087)

### Contrib

+ Add support for building docker debs for Ubuntu 17.04 Zesty on amd64 [#32435](https://github.com/docker/docker/pull/32435)

### Daemon

- Fix `--api-cors-header` being ignored if `--api-enable-cors` is not set [#32174](https://github.com/docker/docker/pull/32174)
- Cleanup docker tmp dir on start [#31741](https://github.com/docker/docker/pull/31741)
- Deprecate `--graph` flag in favor or `--data-root` [#28696](https://github.com/docker/docker/pull/28696)

### Logging

+ Add support for logging driver plugins [#28403](https://github.com/docker/docker/pull/28403)
* Add support for showing logs of individual tasks to `docker service logs`, and add `/task/{id}/logs` REST endpoint [#32015](https://github.com/docker/docker/pull/32015)
* Add `--log-opt env-regex` option to match environment variables using a regular expression [#27565](https://github.com/docker/docker/pull/27565)

### Networking

+ Allow user to replace, and customize the ingress network [#31714](https://github.com/docker/docker/pull/31714)
- Fix UDP traffic in containers not working after the container is restarted [#32505](https://github.com/docker/docker/pull/32505)
- Fix files being written to `/var/lib/docker` if a different data-root is set [#32505](https://github.com/docker/docker/pull/32505)

### Runtime

- Ensure health probe is stopped when a container exits [#32274](https://github.com/docker/docker/pull/32274)

### Swarm Mode

+ Add update/rollback order for services (`--update-order` / `--rollback-order`) [#30261](https://github.com/docker/docker/pull/30261)
+ Add support for synchronous `service create` and `service update` [#31144](https://github.com/docker/docker/pull/31144)
+ Add support for "grace periods" on healthchecks through the `HEALTHCHECK --start-period` and `--health-start-period` flag to
  `docker service create`, `docker service update`, `docker create`, and `docker run` to support containers with an initial startup
  time [#28938](https://github.com/docker/docker/pull/28938)
* `docker service create` now omits fields that are not specified by the user, when possible. This will allow defaults to be applied inside the manager [#32284](https://github.com/docker/docker/pull/32284)
* `docker service inspect` now shows default values for fields that are not specified by the user [#32284](https://github.com/docker/docker/pull/32284)
* Move `docker service logs` out of experimental [#32462](https://github.com/docker/docker/pull/32462)
* Add support for Credential Spec and SELinux to services to the API [#32339](https://github.com/docker/docker/pull/32339)
* Add `--entrypoint` flag to `docker service create` and `docker service update` [#29228](https://github.com/docker/docker/pull/29228)
* Add `--network-add` and `--network-rm` to `docker service update` [#32062](https://github.com/docker/docker/pull/32062)
* Add `--credential-spec` flag to `docker service create` and `docker service update` [#32339](https://github.com/docker/docker/pull/32339)
* Add `--filter mode=<global|replicated>` to `docker service ls` [#31538](https://github.com/docker/docker/pull/31538)
* Resolve network IDs on the client side, instead of in the daemon when creating services [#32062](https://github.com/docker/docker/pull/32062)
* Add `--format` option to `docker node ls` [#30424](https://github.com/docker/docker/pull/30424)
* Add `--prune` option to `docker stack deploy` to remove services that are no longer defined in the docker-compose file [#31302](https://github.com/docker/docker/pull/31302)
* Add `PORTS` column for `docker service ls` when using `ingress` mode [#30813](https://github.com/docker/docker/pull/30813)
- Fix unnescessary re-deploying of tasks when environment-variables are used [#32364](https://github.com/docker/docker/pull/32364)
- Fix `docker stack deploy` not supporting `endpoint_mode` when deploying from a docker compose file [#32333](https://github.com/docker/docker/pull/32333)
- Proceed with startup if cluster component cannot be created to allow recovering from a broken swarm setup [#31631](https://github.com/docker/docker/pull/31631)

### Security

* Allow setting SELinux type or MCS labels when using `--ipc=container:` or `--ipc=host` [#30652](https://github.com/docker/docker/pull/30652)


### Deprecation

- Deprecate `--api-enable-cors` daemon flag. This flag was marked deprecated in Docker 1.6.0 but not listed in deprecated features [#32352](https://github.com/docker/docker/pull/32352)
- Remove Ubuntu 12.04 (Precise Pangolin) as supported platform. Ubuntu 12.04 is EOL, and no longer receives updates [#32520](https://github.com/docker/docker/pull/32520)

## 17.04.0-ce (2017-04-05)

### Builder

* Disable container logging for build containers [#29552](https://github.com/docker/docker/pull/29552)
* Fix use of `**/` in `.dockerignore` [#29043](https://github.com/docker/docker/pull/29043)

### Client

+ Sort `docker stack ls` by name [#31085](https://github.com/docker/docker/pull/31085)
+ Flags for specifying bind mount consistency [#31047](https://github.com/docker/docker/pull/31047)
* Output of docker CLI --help is now wrapped to the terminal width [#28751](https://github.com/docker/docker/pull/28751)
* Suppress image digest in docker ps [#30848](https://github.com/docker/docker/pull/30848)
* Hide command options that are related to Windows [#30788](https://github.com/docker/docker/pull/30788)
* Fix `docker plugin install` prompt to accept "enter" for the "N" default [#30769](https://github.com/docker/docker/pull/30769)
+ Add `truncate` function for Go templates [#30484](https://github.com/docker/docker/pull/30484)
* Support expanded syntax of ports in `stack deploy` [#30476](https://github.com/docker/docker/pull/30476)
* Support expanded syntax of mounts in `stack deploy` [#30597](https://github.com/docker/docker/pull/30597) [#31795](https://github.com/docker/docker/pull/31795)
+ Add `--add-host` for docker build [#30383](https://github.com/docker/docker/pull/30383)
+ Add `.CreatedAt` placeholder for `docker network ls --format` [#29900](https://github.com/docker/docker/pull/29900)
* Update order of `--secret-rm` and `--secret-add` [#29802](https://github.com/docker/docker/pull/29802)
+ Add `--filter enabled=true` for `docker plugin ls` [#28627](https://github.com/docker/docker/pull/28627)
+ Add `--format` to `docker service ls` [#28199](https://github.com/docker/docker/pull/28199)
+ Add `publish` and `expose` filter for `docker ps --filter` [#27557](https://github.com/docker/docker/pull/27557)
* Support multiple service IDs on `docker service ps` [#25234](https://github.com/docker/docker/pull/25234)
+ Allow swarm join with `--availability=drain` [#24993](https://github.com/docker/docker/pull/24993)
* Docker inspect now shows "docker-default" when AppArmor is enabled and no other profile was defined [#27083](https://github.com/docker/docker/pull/27083)

### Logging

+ Implement optional ring buffer for container logs [#28762](https://github.com/docker/docker/pull/28762)
+ Add `--log-opt awslogs-create-group=<true|false>` for awslogs (CloudWatch) to support creation of log groups as needed [#29504](https://github.com/docker/docker/pull/29504)
- Fix segfault when using the gcplogs logging driver with a "static" binary [#29478](https://github.com/docker/docker/pull/29478)


### Networking

* Check parameter `--ip`, `--ip6` and `--link-local-ip` in `docker network connect` [#30807](https://github.com/docker/docker/pull/30807)
+ Added support for `dns-search` [#30117](https://github.com/docker/docker/pull/30117)
+ Added --verbose option for docker network inspect to show task details from all swarm nodes [#31710](https://github.com/docker/docker/pull/31710)
* Clear stale datapath encryption states when joining the cluster [docker/libnetwork#1354](https://github.com/docker/libnetwork/pull/1354)
+ Ensure iptables initialization only happens once [docker/libnetwork#1676](https://github.com/docker/libnetwork/pull/1676)
* Fix bad order of iptables filter rules [docker/libnetwork#961](https://github.com/docker/libnetwork/pull/961)
+ Add anonymous container alias to service record on attachable network [docker/libnetwork#1651](https://github.com/docker/libnetwork/pull/1651)
+ Support for `com.docker.network.container_interface_prefix` driver label [docker/libnetwork#1667](https://github.com/docker/libnetwork/pull/1667)
+ Improve network list performance by omitting network details that are not used [#30673](https://github.com/docker/docker/pull/30673)

### Runtime

* Handle paused container when restoring without live-restore set [#31704](https://github.com/docker/docker/pull/31704)
- Do not allow sub second in healthcheck options in Dockerfile [#31177](https://github.com/docker/docker/pull/31177)
* Support name and id prefix in `secret update` [#30856](https://github.com/docker/docker/pull/30856)
* Use binary frame for websocket attach endpoint [#30460](https://github.com/docker/docker/pull/30460)
* Fix linux mount calls not applying propagation type changes [#30416](https://github.com/docker/docker/pull/30416)
* Fix ExecIds leak on failed `exec -i` [#30340](https://github.com/docker/docker/pull/30340)
* Prune named but untagged images if `danglingOnly=true` [#30330](https://github.com/docker/docker/pull/30330)
+ Add daemon flag to set `no_new_priv` as default for unprivileged containers [#29984](https://github.com/docker/docker/pull/29984)
+ Add daemon option `--default-shm-size` [#29692](https://github.com/docker/docker/pull/29692)
+ Support registry mirror config reload [#29650](https://github.com/docker/docker/pull/29650)
- Ignore the daemon log config when building images [#29552](https://github.com/docker/docker/pull/29552)
* Move secret name or ID prefix resolving from client to daemon [#29218](https://github.com/docker/docker/pull/29218)
+ Allow adding rules to `cgroup devices.allow` on container create/run [#22563](https://github.com/docker/docker/pull/22563)
- Fix `cpu.cfs_quota_us` being reset when running `systemd daemon-reload` [#31736](https://github.com/docker/docker/pull/31736)

### Swarm Mode

+ Topology-aware scheduling [#30725](https://github.com/docker/docker/pull/30725)
+ Automatic service rollback on failure [#31108](https://github.com/docker/docker/pull/31108)
+ Worker and manager on the same node are now connected through a UNIX socket [docker/swarmkit#1828](https://github.com/docker/swarmkit/pull/1828), [docker/swarmkit#1850](https://github.com/docker/swarmkit/pull/1850), [docker/swarmkit#1851](https://github.com/docker/swarmkit/pull/1851)
* Improve raft transport package [docker/swarmkit#1748](https://github.com/docker/swarmkit/pull/1748)
* No automatic manager shutdown on demotion/removal [docker/swarmkit#1829](https://github.com/docker/swarmkit/pull/1829)
* Use TransferLeadership to make leader demotion safer [docker/swarmkit#1939](https://github.com/docker/swarmkit/pull/1939)
* Decrease default monitoring period [docker/swarmkit#1967](https://github.com/docker/swarmkit/pull/1967)
+ Add Service logs formatting [#31672](https://github.com/docker/docker/pull/31672)
* Fix service logs API to be able to specify stream [#31313](https://github.com/docker/docker/pull/31313)
+ Add `--stop-signal` for `service create` and `service update` [#30754](https://github.com/docker/docker/pull/30754)
+ Add `--read-only` for `service create` and `service update` [#30162](https://github.com/docker/docker/pull/30162)
+ Renew the context after communicating with the registry [#31586](https://github.com/docker/docker/pull/31586)
+ (experimental) Add `--tail` and `--since` options to `docker service logs` [#31500](https://github.com/docker/docker/pull/31500)
+ (experimental) Add `--no-task-ids` and `--no-trunc` options to `docker service logs` [#31672](https://github.com/docker/docker/pull/31672)

### Windows

* Block pulling Windows images on non-Windows daemons [#29001](https://github.com/docker/docker/pull/29001)

{% endraw %}
