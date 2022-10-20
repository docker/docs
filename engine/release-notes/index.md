---
title: Docker Engine release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Engine
keywords: docker, docker engine, ce, whats new, release notes
toc_min: 1
toc_max: 2
skip_read_time: true
redirect_from:
  - /release-notes/docker-ce/
  - /release-notes/docker-engine/
---

This document describes the latest changes, additions, known issues, and fixes
for Docker Engine.

# Version 20.10

## 20.10.20
2022-10-18

This release of Docker Engine contains partial mitigations for a Git vulnerability
([CVE-2022-39253](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-39253){:target="_blank" rel="noopener"}),
and has updated handling of `image:tag@digest` image references.

The Git vulnerability allows a maliciously crafted Git repository, when used as a
build context, to copy arbitrary filesystem paths into resulting containers/images;
this can occur in both the daemon, and in API clients, depending on the versions and
tools in use.

The mitigations available in this release and in other consumers of the daemon API
are partial and only protect users who build a Git URL context (e.g. `git+protocol://`).
As the vulnerability could still be exploited by manually run Git commands that interact
with and check out submodules, users should immediately upgrade to a patched version of
Git to protect against this vulernability. Further details are available from the GitHub
blog (["Git security vulnerabilities announced"](https://github.blog/2022-10-18-git-security-vulnerabilities-announced/){:target="_blank" rel="noopener"}).


### Client

- Added a mitigation for [CVE-2022-39253](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-39253){:target="_blank" rel="noopener"},
  when using the classic Builder with a Git URL as the build context.

### Daemon

- Updated handling of `image:tag@digest` references. When pulling an image using
  the `image:tag@digest` ("pull by digest"), image resolution happens through
  the content-addressable digest and the `image` and `tag` are not used. While
  this is expected, this could lead to confusing behavior, and could potentially
  be exploited through social engineering to run an image that is already present
  in the local image store. Docker now checks if the digest matches the repository
  name used to pull the image, and otherwise will produce an error.


### Builder

- Updated handling of `image:tag@digest` references. Refer to the "Daemon" section
  above for details.
- Added a mitigation to the classic Builder and updated BuildKit to [v0.8.3-31-gc0149372](https://github.com/moby/buildkit/commit/c014937225cba29cfb1d5161fd134316c0e9bdaa){:target="_blank" rel="noopener"},
  for [CVE-2022-39253](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-39253){:target="_blank" rel="noopener"}.

### Packaging

- Update Docker Compose to [v2.12.0](https://github.com/docker/compose/releases/tag/v2.12.0){:target="_blank" rel="noopener"}.

## 20.10.19
2022-10-14

This release of Docker Engine comes with some bug-fixes, and an updated version
of Docker Compose.

### Builder

- Fix an issue that could result in a panic during `docker builder prune` or
  `docker system prune` [moby/moby#44122](https://github.com/moby/moby/pull/44122){:target="_blank" rel="noopener"}.

### Daemon

- Fix a bug where using `docker volume prune` would remove volumes that were
  still in use if the daemon was running with "live restore" and was restarted
  [moby/moby#44238](https://github.com/moby/moby/pull/44238){:target="_blank" rel="noopener"}.

### Packaging

- Update Docker Compose to [v2.11.2](https://github.com/docker/compose/releases/tag/v2.11.2){:target="_blank" rel="noopener"}.
- Update Go runtime to [1.18.7](https://go.dev/doc/devel/release#go1.18.minor){:target="_blank" rel="noopener"},
  which contains fixes for [CVE-2022-2879](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-2879){:target="_blank" rel="noopener"},
  [CVE-2022-2880](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-2880){:target="_blank" rel="noopener"},
  and [CVE-2022-41715](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-41715){:target="_blank" rel="noopener"}.

## 20.10.18
2022-09-09

This release of Docker Engine comes with a fix for a low-severity security issue,
some minor bug fixes, and updated versions of Docker Compose, Docker Buildx,
`containerd`, and `runc`.

### Client

- Add Bash completion for Docker Compose [docker/cli#3752](https://github.com/docker/cli/pull/3752){:target="_blank" rel="noopener"}.

### Builder

- Fix an issue where file-capabilities were not preserved during build
  [moby/moby#43876](https://github.com/moby/moby/pull/43876){:target="_blank" rel="noopener"}.
- Fix an issue that could result in a panic caused by a concurrent map read and
  map write [moby/moby#44067](https://github.com/moby/moby/pull/44067){:target="_blank" rel="noopener"}.

### Daemon

- Fix a security vulnerability relating to supplementary group permissions, which
  could allow a container process to bypass primary group restrictions within the
  container [CVE-2022-36109](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-36109){:target="_blank" rel="noopener"},
  [GHSA-rc4r-wh2q-q6c4](https://github.com/moby/moby/security/advisories/GHSA-rc4r-wh2q-q6c4){:target="_blank" rel="noopener"}.
- seccomp: add support for Landlock syscalls in default policy [moby/moby#43991](https://github.com/moby/moby/pull/43991){:target="_blank" rel="noopener"}.
- seccomp: update default policy to support new syscalls introduced in kernel 5.12 - 5.16 [moby/moby#43991](https://github.com/moby/moby/pull/43991){:target="_blank" rel="noopener"}.
- Fix an issue where cache lookup for image manifests would fail, resulting
  in a redundant round-trip to the image registry [moby/moby#44109](https://github.com/moby/moby/pull/44109){:target="_blank" rel="noopener"}.
- Fix an issue where `exec` processes and healthchecks were not terminated
  when they timed out [moby/moby#44018](https://github.com/moby/moby/pull/44018){:target="_blank" rel="noopener"}.

### Packaging

- Update Docker Buildx to [v0.9.1](https://github.com/docker/buildx/releases/tag/v0.9.1){:target="_blank" rel="noopener"}.
- Update Docker Compose to [v2.10.2](https://github.com/docker/compose/releases/tag/v2.10.2){:target="_blank" rel="noopener"}.
- Update containerd (`containerd.io` package) to [v1.6.8](https://github.com/containerd/containerd/releases/tag/v1.6.8){:target="_blank" rel="noopener"}.
- Update runc version to [v1.1.4](https://github.com/opencontainers/runc/releases/tag/v1.1.4){:target="_blank" rel="noopener"}.
- Update Go runtime to [1.18.6](https://go.dev/doc/devel/release#go1.18.minor){:target="_blank" rel="noopener"},
  which contains fixes for [CVE-2022-27664](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-27664){:target="_blank" rel="noopener"} and
  [CVE-2022-32190](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-32190){:target="_blank" rel="noopener"}.

## 20.10.17
2022-06-06

This release of Docker Engine comes with updated versions of Docker Compose and the
`containerd`, and `runc` components, as well as some minor bug fixes.

### Client

- Remove asterisk from docker commands in zsh completion script [docker/cli#3648](https://github.com/docker/cli/pull/3648){:target="_blank" rel="noopener"}.

### Networking

- Fix Windows port conflict with published ports in host mode for overlay [moby/moby#43644](https://github.com/moby/moby/pull/43644){:target="_blank" rel="noopener"}.
- Ensure performance tuning is always applied to libnetwork sandboxes [moby/moby#43683](https://github.com/moby/moby/pull/43683){:target="_blank" rel="noopener"}.

### Packaging

- Update Docker Compose to [v2.6.0](https://github.com/docker/compose/releases/tag/v2.6.0){:target="_blank" rel="noopener"}.
- Update containerd (`containerd.io` package) to [v1.6.6](https://github.com/containerd/containerd/releases/tag/v1.6.6),
  which contains a fix for [CVE-2022-31030](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-31030){:target="_blank" rel="noopener"}
- Update runc version to [v1.1.2](https://github.com/opencontainers/runc/releases/tag/v1.1.2), which contains a fix for
  [CVE-2022-29162](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-29162){:target="_blank" rel="noopener"}.
- Update Go runtime to [1.17.11](https://go.dev/doc/devel/release#go1.17.minor){:target="_blank" rel="noopener"},
  which contains fixes for [CVE-2022-30634](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-30634){:target="_blank" rel="noopener"},
  [CVE-2022-30629](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-30629){:target="_blank" rel="noopener"},
  [CVE-2022-30580](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-30580){:target="_blank" rel="noopener"} and
  [CVE-2022-29804](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-29804){:target="_blank" rel="noopener"}

## 20.10.16
2022-05-12

This release of Docker Engine fixes a regression in the Docker CLI builds for
macOS, fixes an issue with `docker stats` when using containerd 1.5 and up,
and updates the Go runtime to include a fix for [CVE-2022-29526](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-29526){:target="_blank" rel="noopener"}.

### Client

- Fixed a regression in binaries for macOS introduced in [20.10.15](#201015), which
  resulted in a panic [docker/cli#43426](https://github.com/docker/cli/pull/3592){:target="_blank" rel="noopener"}.
- Update golang.org/x/sys dependency which contains a fix for
  [CVE-2022-29526](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-29526){:target="_blank" rel="noopener"}.

### Daemon

- Fixed an issue where `docker stats` was showing empty stats when running with
  containerd 1.5.0 or up [moby/moby#43567](https://github.com/moby/moby/pull/43567){:target="_blank" rel="noopener"}.
- Updated the `golang.org/x/sys` build-time dependency which contains a fix for [CVE-2022-29526](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-29526){:target="_blank" rel="noopener"}.

### Packaging

- Updated Go runtime to [1.17.10](https://go.dev/doc/devel/release#go1.17.minor){:target="_blank" rel="noopener"},
  which contains a fix for [CVE-2022-29526](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-29526){:target="_blank" rel="noopener"}.
- Used "weak" dependencies for the `docker scan` CLI plugin, to prevent a
  "conflicting requests" error when users performed an off-line installation from
  downloaded RPM packages [docker/docker-ce-packaging#659](https://github.com/docker/docker-ce-packaging/pull/659){:target="_blank" rel="noopener"}.

## 20.10.15
2022-05-05

This release of Docker Engine comes with updated versions of the `compose`,
`buildx`, `containerd`, and `runc` components, as well as some minor bug fixes.

> **Known issues**
> 
> We've identified an issue with the [macOS CLI binaries](https://download.docker.com/mac/static/stable/){:target="_blank" rel="noopener" class="_"}
> in the 20.10.15 release. This issue has been resolved in the [20.10.16](#201016) release.
{:.important}

### Daemon

- Use a RWMutex for stateCounter to prevent potential locking congestion [moby/moby#43426](https://github.com/moby/moby/pull/43426).
- Prevent an issue where the daemon was unable to find an available IP-range in
  some conditions [moby/moby#43360](https://github.com/moby/moby/pull/43360) 

### Packaging

- Update Docker Compose to [v2.5.0](https://github.com/docker/compose/releases/tag/v2.5.0).
- Update Docker Buildx to [v0.8.2](https://github.com/docker/buildx/releases/tag/v0.8.2).
- Update Go runtime to [1.17.9](https://go.dev/doc/devel/release#go1.17.minor).
- Update containerd (`containerd.io` package) to [v1.6.4](https://github.com/containerd/containerd/releases/tag/v1.6.4).
- Update runc version to [v1.1.1](https://github.com/opencontainers/runc/releases/tag/v1.1.1).
- Add packages for CentOS 9 stream and Fedora 36.

## 20.10.14
2022-03-23

This release of Docker Engine updates the default inheritable capabilities for
containers to address [CVE-2022-24769](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-24769),
a new version of the `containerd.io` runtime is also included to address the same
issue.

### Daemon

- Update the default inheritable capabilities.

### Builder

- Update the default inheritable capabilities for containers used during build.

### Packaging

- Update containerd (`containerd.io` package) to [v1.5.11](https://github.com/containerd/containerd/releases/tag/v1.5.11).
- Update `docker buildx` to [v0.8.1](https://github.com/docker/buildx/releases/tag/v0.8.1).

## 20.10.13
2022-03-10

This release of Docker Engine contains some bug-fixes and packaging changes,
updates to the `docker scan` and `docker buildx` commands, an updated version of
the Go runtime, and new versions of the `containerd.io` runtime.
Together with this release, we now also provide `.deb` and `.rpm` packages of
Docker Compose V2, which can be installed using the (optional) `docker-compose-plugin`
package.

### Builder

- Updated the bundled version of buildx to [v0.8.0](https://github.com/docker/buildx/releases/tag/v0.8.0).

### Daemon

- Fix a race condition when updating the container's state [moby/moby#43166](https://github.com/moby/moby/pull/43166).
- Update the etcd dependency to prevent the daemon from incorrectly holding file locks [moby/moby#43259](https://github.com/moby/moby/pull/43259)
- Fix detection of user-namespaces when configuring the default `net.ipv4.ping_group_range` sysctl [moby/moby#43084](https://github.com/moby/moby/pull/43084).

### Distribution

- Retry downloading image-manifests if a connection failure happens during image
  pull [moby/moby#43333](https://github.com/moby/moby/pull/43333).

### Documentation

- Various fixes in command-line reference and API documentation.

### Logging

- Prevent an OOM when using the "local" logging driver with containers that produce
  a large amount of log messages [moby/moby#43165](https://github.com/moby/moby/pull/43165).
- Updates the fluentd log driver to prevent a potential daemon crash, and prevent
  containers from hanging when using the `fluentd-async-connect=true` and the
  remote server is unreachable [moby/moby#43147](https://github.com/moby/moby/pull/43147).

### Packaging

- Provide `.deb` and `.rpm` packages for Docker Compose V2. [Docker Compose v2.3.3](https://github.com/docker/compose/releases/tag/v2.3.3)
  can now be installed on Linux using the `docker-compose-plugin` packages, which
  provides the `docker compose` subcommand on the Docker CLI. The Docker Compose
  plugin can also be installed and run standalone to be used as a drop-in replacement
  for `docker-compose` (Docker Compose V1) [docker/docker-ce-packaging#638](https://github.com/docker/docker-ce-packaging/pull/638).
  The `compose-cli-plugin` package can also be used on older version of the Docker
  CLI with support for CLI plugins (Docker CLI 18.09 and up).
- Provide packages for the upcoming Ubuntu 22.04 "Jammy Jellyfish" LTS release [docker/docker-ce-packaging#645](https://github.com/docker/docker-ce-packaging/pull/645), [docker/containerd-packaging#271](https://github.com/docker/containerd-packaging/pull/271).
- Update `docker buildx` to [v0.8.0](https://github.com/docker/buildx/releases/tag/v0.8.0).
- Update `docker scan` (`docker-scan-plugin`) to [v0.17.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.17.0).
- Update containerd (`containerd.io` package) to [v1.5.10](https://github.com/containerd/containerd/releases/tag/v1.5.10).
- Update the bundled runc version to [v1.0.3](https://github.com/opencontainers/runc/releases/tag/v1.0.3).
- Update Golang runtime to Go 1.16.15.


## 20.10.12
2021-12-13

This release of Docker Engine contains changes in packaging only, and provides
updates to the `docker scan` and `docker buildx` commands. Versions of `docker scan`
before v0.11.0 are not able to detect the [Log4j 2 CVE-2021-44228](https://nvd.nist.gov/vuln/detail/CVE-2021-44228).
We are shipping an updated version of `docker scan` in this release to help you
scan your images for this vulnerability.

> **Note**
>
> The `docker scan` command on Linux is currently only supported on x86 platforms.
> We do not yet provide a package for other hardware architectures on Linux.

The `docker scan` feature is provided as a separate package and, depending on your
upgrade or installation method, 'docker scan' may not be updated automatically to
the latest version. Use the instructions below to update `docker scan` to the latest
version. You can also use these instructions to install, or upgrade the `docker scan`
package without upgrading the Docker Engine:

On `.deb` based distros, such as Ubuntu and Debian:

```console
$ apt-get update && apt-get install docker-scan-plugin
```

On rpm-based distros, such as CentOS or Fedora:

```console
$ yum install docker-scan-plugin
```

After upgrading, verify you have the latest version of `docker scan` installed:

```console
$ docker scan --accept-license --version
Version:    v0.12.0
Git commit: 1074dd0
Provider:   Snyk (1.790.0 (standalone))
```

[Read our blog post on CVE-2021-44228](https://www.docker.com/blog/apache-log4j-2-cve-2021-44228/)
to learn how to use the `docker scan` command to check if images are vulnerable.

### Packaging

- Update `docker scan` to [v0.12.0](https://github.com/docker/scan-cli-plugin/releases/tag/v0.12.0).
- Update `docker buildx` to [v0.7.1](https://github.com/docker/buildx/releases/tag/v0.7.1).
- Update Golang runtime to Go 1.16.12.


## 20.10.11
2021-11-17

> **IMPORTANT**
>
> Due to [net/http changes](https://github.com/golang/go/issues/40909) in [Go 1.16](https://golang.org/doc/go1.16#net/http),
> HTTP proxies configured through the `$HTTP_PROXY` environment variable are no
> longer used for TLS (`https://`) connections. Make sure you also set an `$HTTPS_PROXY`
> environment variable for handling requests to `https://` URLs.
>
> Refer to the [HTTP/HTTPS proxy section](../../config/daemon/systemd.md#httphttps-proxy)
> to learn how to configure the Docker Daemon to use a proxy server.
{: .important }


### Distribution

- Handle ambiguous OCI manifest parsing to mitigate [CVE-2021-41190](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-41190) / [GHSA-mc8v-mgrf-8f4m](https://github.com/opencontainers/distribution-spec/security/advisories/GHSA-mc8v-mgrf-8f4m).
  See [GHSA-xmmx-7jpf-fx42](https://github.com/moby/moby/security/advisories/GHSA-xmmx-7jpf-fx42) for details.

### Windows

- Fix panic.log file having read-only attribute set [moby/moby#42987](https://github.com/moby/moby/pull/42987).

### Packaging

- Update containerd to [v1.4.12](https://github.com/containerd/containerd/releases/tag/v1.4.12) to mitigate [CVE-2021-41190](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-41190).
- Update Golang runtime to Go 1.16.10.


## 20.10.10
2021-10-25

> **IMPORTANT**
>
> Due to [net/http changes](https://github.com/golang/go/issues/40909) in [Go 1.16](https://golang.org/doc/go1.16#net/http),
> HTTP proxies configured through the `$HTTP_PROXY` environment variable are no
> longer used for TLS (`https://`) connections. Make sure you also set an `$HTTPS_PROXY`
> environment variable for handling requests to `https://` URLs.
>
> Refer to the [HTTP/HTTPS proxy section](../../config/daemon/systemd.md#httphttps-proxy)
> to learn how to configure the Docker Daemon to use a proxy server.
{: .important }


### Builder

- Fix platform-matching logic to fix `docker build` using not finding images in
  the local image cache on Arm machines when using BuildKit [moby/moby#42954](https://github.com/moby/moby/pull/42954)

### Runtime

- Add support for `clone3` syscall in the default seccomp policy to support running
  containers based on recent versions of Fedora and Ubuntu. [moby/moby/#42836](https://github.com/moby/moby/pull/42836).
- Windows: update hcsshim library to fix a bug in sparse file handling in container
  layers, which was exposed by recent changes in Windows [moby/moby#42944](https://github.com/moby/moby/pull/42944).
- Fix some situations where `docker stop` could hang forever [moby/moby#42956](https://github.com/moby/moby/pull/42956).

### Swarm

- Fix an issue where updating a service did not roll back on failure [moby/moby#42875](https://github.com/moby/moby/pull/42875).

### Packaging

- Add packages for Ubuntu 21.10 "Impish Indri" and Fedora 35.
- Update `docker scan` to v0.9.0
- Update Golang runtime to Go 1.16.9.

## 20.10.9
2021-10-04

This release is a security release with security fixes in the CLI, runtime, as
well as updated versions of the containerd.io package.

> **IMPORTANT**
>
> Due to [net/http changes](https://github.com/golang/go/issues/40909) in [Go 1.16](https://golang.org/doc/go1.16#net/http),
> HTTP proxies configured through the `$HTTP_PROXY` environment variable are no
> longer used for TLS (`https://`) connections. Make sure you also set an `$HTTPS_PROXY`
> environment variable for handling requests to `https://` URLs.
>
> Refer to the [HTTP/HTTPS proxy section](../../config/daemon/systemd.md#httphttps-proxy)
> to learn how to configure the Docker Daemon to use a proxy server.
{: .important }

### Client

- [CVE-2021-41092](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-41092)
  Ensure default auth config has address field set, to prevent credentials being
  sent to the default registry.

### Runtime

- [CVE-2021-41089](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-41089)
  Create parent directories inside a chroot during `docker cp` to prevent a specially
  crafted container from changing permissions of existing files in the host’s filesystem.
- [CVE-2021-41091](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-41091)
  Lock down file permissions to prevent unprivileged users from discovering and
  executing programs in `/var/lib/docker`.

### Packaging

> **Known issue**
>
> The `ctr` binary shipping with the static packages of this release is not
> statically linked, and will not run in Docker images using alpine as a base
> image. Users can install the `libc6-compat` package, or download a previous
> version of the `ctr` binary as a workaround. Refer to the containerd ticket
> related to this issue for more details: [containerd/containerd#5824](https://github.com/containerd/containerd/issues/5824).

- Update Golang runtime to Go 1.16.8, which contains fixes for [CVE-2021-36221](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-36221)
  and [CVE-2021-39293](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-39293)
- Update static binaries and containerd.io rpm and deb packages to containerd
  v1.4.11 and runc v1.0.2 to address [CVE-2021-41103](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-41103).
- Update the bundled buildx version to v0.6.3 for rpm and deb packages.

## 20.10.8
2021-08-03

> **IMPORTANT**
>
> Due to [net/http changes](https://github.com/golang/go/issues/40909) in [Go 1.16](https://golang.org/doc/go1.16#net/http),
> HTTP proxies configured through the `$HTTP_PROXY` environment variable are no
> longer used for TLS (`https://`) connections. Make sure you also set an `$HTTPS_PROXY`
> environment variable for handling requests to `https://` URLs.
>
> Refer to the [HTTP/HTTPS proxy section](../../config/daemon/systemd.md#httphttps-proxy)
> to learn how to configure the Docker Daemon to use a proxy server.
{: .important }

### Deprecation

- Deprecate support for encrypted TLS private keys. Legacy PEM encryption as
  specified in RFC 1423 is insecure by design. Because it does not authenticate
  the ciphertext, it is vulnerable to padding oracle attacks that can let an
  attacker recover the plaintext. Support for encrypted TLS private keys is now
  marked as deprecated, and will be removed in an upcoming release. [docker/cli#3219](https://github.com/docker/cli/pull/3219)
- Deprecate Kubernetes stack support. Following the deprecation of [Compose on Kubernetes](https://github.com/docker/compose-on-kubernetes),
  support for Kubernetes in the `stack` and `context` commands in the Docker CLI
  is now marked as deprecated, and will be removed in an upcoming release [docker/cli#3174](https://github.com/docker/cli/pull/3174).

### Client

- Fix `Invalid standard handle identifier` errors on Windows [docker/cli#3132](https://github.com/docker/cli/pull/3132).

### Rootless

- Avoid `can't open lock file /run/xtables.lock: Permission denied` error on
  SELinux hosts [moby/moby#42462](https://github.com/moby/moby/pull/42462).
- Disable overlay2 when running with SELinux to prevent permission denied errors [moby/moby#42462](https://github.com/moby/moby/pull/42462).
- Fix `x509: certificate signed by unknown authority` error on openSUSE Tumbleweed [moby/moby#42462](https://github.com/moby/moby/pull/42462).

### Runtime

- Print a warning when using the `--platform` option to pull a single-arch image
  that does not match the specified architecture [moby/moby#42633](https://github.com/moby/moby/pull/42633).
- Fix incorrect `Your kernel does not support swap memory limit` warning when
  running with cgroups v2 [moby/moby#42479](https://github.com/moby/moby/pull/42479).
- Windows: Fix a situation where containers were not stopped if `HcsShutdownComputeSystem`
  returned an `ERROR_PROC_NOT_FOUND` error [moby/moby#42613](https://github.com/moby/moby/pull/42613) 

### Swarm

- Fix a possibility where overlapping IP addresses could exist as a result of the
  node failing to clean up its old loadbalancer IPs [moby/moby#42538](https://github.com/moby/moby/pull/42538)
- Fix a deadlock in log broker ("dispatcher is stopped") [moby/moby#42537](https://github.com/moby/moby/pull/42537)

### Packaging

> **Known issue**
>
> The `ctr` binary shipping with the static packages of this release is not
> statically linked, and will not run in Docker images using alpine as a base
> image. Users can install the `libc6-compat` package, or download a previous
> version of the `ctr` binary as a workaround. Refer to the containerd ticket
> related to this issue for more details: [containerd/containerd#5824](https://github.com/containerd/containerd/issues/5824).

- Remove packaging for Ubuntu 16.04 "Xenial" and Fedora 32, as they reached EOL [docker/docker-ce-packaging#560](https://github.com/docker/docker-ce-packaging/pull/560)
- Update Golang runtime to Go 1.16.6
- Update the bundled buildx version to v0.6.1 for rpm and deb packages [docker/docker-ce-packaging#562](https://github.com/docker/docker-ce-packaging/pull/562)
- Update static binaries and containerd.io rpm and deb packages to containerd v1.4.9 and runc v1.0.1: [docker/containerd-packaging#241](https://github.com/docker/containerd-packaging/pull/241), [docker/containerd-packaging#245](https://github.com/docker/containerd-packaging/pull/245), [docker/containerd-packaging#247](https://github.com/docker/containerd-packaging/pull/247).

## 20.10.7
2021-06-02

### Client

* Suppress warnings for deprecated cgroups [docker/cli#3099](https://github.com/docker/cli/pull/3099).
* Prevent sending `SIGURG` signals to container on Linux and macOS. The Go runtime
  (starting with Go 1.14) uses `SIGURG` signals internally as an interrupt to
  support preemptable syscalls. In situations where the Docker CLI was attached
  to a container, these interrupts were forwarded to the container. This fix
  changes the Docker CLI to ignore `SIGURG` signals [docker/cli#3107](https://github.com/docker/cli/pull/3107),
  [moby/moby#42421](https://github.com/moby/moby/pull/42421).

### Builder

* Update BuildKit to version v0.8.3-3-g244e8cde [moby/moby#42448](https://github.com/moby/moby/pull/42448):
    * Transform relative mountpoints for exec mounts in the executor to work around
      a breaking change in runc v1.0.0-rc94 and up. [moby/buildkit#2137](https://github.com/moby/buildkit/pull/2137).
    * Add retry on image push 5xx errors. [moby/buildkit#2043](https://github.com/moby/buildkit/pull/2043).
    * Fix build-cache not being invalidated when renaming a file that is copied using
      a `COPY` command with a wildcard. Note that this change invalidates
      existing build caches for copy commands that use a wildcard. [moby/buildkit#2018](https://github.com/moby/buildkit/pull/2018).
    * Fix build-cache not being invalidated when using mounts [moby/buildkit#2076](https://github.com/moby/buildkit/pull/2076).
* Fix build failures when `FROM` image is not cached when using legacy schema 1 images [moby/moby#42382](https://github.com/moby/moby/pull/42382).

### Logging

* Update the hcsshim SDK to make daemon logs on Windows less verbose [moby/moby#42292](https://github.com/moby/moby/pull/42292).

### Rootless

* Fix capabilities not being honored when an image was built on a daemon with
  user-namespaces enabled [moby/moby#42352](https://github.com/moby/moby/pull/42352).

### Networking

* Update libnetwork to fix publishing ports on environments with kernel boot
  parameter `ipv6.disable=1`, and to fix a deadlock causing internal DNS lookups
  to fail [moby/moby#42413](https://github.com/moby/moby/pull/42413).

### Contrib

* Update rootlesskit to v0.14.2 to fix a timeout when starting the userland proxy
  with the `slirp4netns` port driver [moby/moby#42294](https://github.com/moby/moby/pull/42294).
* Fix "Device or resource busy" errors when running docker-in-docker on a rootless
  daemon [moby/moby#42342](https://github.com/moby/moby/pull/42342).

### Packaging

* Update containerd to v1.4.6, runc v1.0.0-rc95 to address [CVE-2021-30465](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-30465)
  [moby/moby#42398](https://github.com/moby/moby/pull/42398), [moby/moby#42395](https://github.com/moby/moby/pull/42395),
  [ocker/containerd-packaging#234](https://github.com/docker/containerd-packaging/pull/234)
* Update containerd to v1.4.5, runc v1.0.0-rc94 [moby/moby#42372](https://github.com/moby/moby/pull/42372),
  [moby/moby#42388](https://github.com/moby/moby/pull/42388), [docker/containerd-packaging#232](https://github.com/docker/containerd-packaging/pull/232).
* Update Docker Scan plugin packages (`docker-scan-plugin`) to v0.8 [docker/docker-ce-packaging#545](https://github.com/docker/docker-ce-packaging/pull/545).


## 20.10.6
2021-04-12

### Client

* Apple Silicon (darwin/arm64) support for Docker CLI [docker/cli#3042](https://github.com/docker/cli/pull/3042)
* config: print deprecation warning when falling back to pre-v1.7.0 config file `~/.dockercfg`. Support for this file will be removed in a future release [docker/cli#3000](https://github.com/docker/cli/pull/3000)

### Builder

* Fix classic builder silently ignoring unsupported Dockerfile options and prompt to enable BuildKit instead [moby/moby#42197](https://github.com/moby/moby/pull/42197)

### Logging

* json-file: fix sporadic unexpected EOF errors [moby/moby#42174](https://github.com/moby/moby/pull/42174)

### Networking

* Fix a regression in docker 20.10, causing  IPv6 addresses no longer to be bound by default when mapping ports [moby/moby#42205](https://github.com/moby/moby/pull/42205)
* Fix implicit IPv6 port-mappings not included in API response. Before docker 20.10, published ports were accessible through both IPv4 and IPv6 by default, but the API only included information about the IPv4 (0.0.0.0) mapping [moby/moby#42205](https://github.com/moby/moby/pull/42205)
* Fix a regression in docker 20.10, causing the docker-proxy  to not be terminated in all cases [moby/moby#42205](https://github.com/moby/moby/pull/42205)
* Fix iptables forwarding rules not being cleaned up upon container removal [moby/moby#42205](https://github.com/moby/moby/pull/42205)

### Packaging

* Update containerd to [v1.4.4](https://github.com/containerd/containerd/releases/tag/v1.4.4) for static binaries. The containerd.io package on apt/yum repos already had this update out of band. Includes a fix for [CVE-2021-21334](https://github.com/containerd/containerd/security/advisories/GHSA-6g2q-w5j3-fwh4). [moby/moby#42124](https://github.com/moby/moby/pull/42124)
* Packages for Debian/Raspbian 11 Bullseye, Ubuntu 21.04 Hirsute Hippo and Fedora 34 [docker/docker-ce-packaging#521](https://github.com/docker/docker-ce-packaging/pull/521) [docker/docker-ce-packaging#522](https://github.com/docker/docker-ce-packaging/pull/522) [docker/docker-ce-packaging#533](https://github.com/docker/docker-ce-packaging/pull/533)
* Provide the [Docker Scan CLI](https://github.com/docker/scan-cli-plugin) plugin on Linux amd64 via a `docker-scan-plugin` package as a recommended dependency for the `docker-ce-cli` package [docker/docker-ce-packaging#537](https://github.com/docker/docker-ce-packaging/pull/537)
* Include VPNKit binary for arm64 [moby/moby#42141](https://github.com/moby/moby/pull/42141)

### Plugins

* Fix docker plugin create making plugins that were incompatible with older versions of Docker [moby/moby#42256](https://github.com/moby/moby/pull/42256)

### Rootless

* Update RootlessKit to [v0.14.1](https://github.com/rootless-containers/rootlesskit/releases/tag/v0.14.1) (see also [v0.14.0](https://github.com/rootless-containers/rootlesskit/releases/tag/v0.14.0) [v0.13.2](https://github.com/rootless-containers/rootlesskit/releases/tag/v0.13.2)) [moby/moby#42186](https://github.com/moby/moby/pull/42186) [moby/moby#42232](https://github.com/moby/moby/pull/42232)
* dockerd-rootless-setuptool.sh: create CLI context "rootless" [moby/moby#42109](https://github.com/moby/moby/pull/42109)
* dockerd-rootless.sh: prohibit running as root [moby/moby#42072](https://github.com/moby/moby/pull/42072)
* Fix "operation not permitted" when bind mounting existing mounts [moby/moby#42233](https://github.com/moby/moby/pull/42233)
* overlay2: fix "createDirWithOverlayOpaque(...) ... input/output error" [moby/moby#42235](https://github.com/moby/moby/pull/42235)
* overlay2: support "userxattr" option (kernel 5.11) [moby/moby#42168](https://github.com/moby/moby/pull/42168)
* btrfs: allow unprivileged user to delete subvolumes (kernel >= 4.18) [moby/moby#42253](https://github.com/moby/moby/pull/42253)
* cgroup2: Move cgroup v2 out of experimental [moby/moby#42263](https://github.com/moby/moby/pull/42263)


## 20.10.5
2021-03-02

### Client

* Revert [docker/cli#2960](https://github.com/docker/cli/pull/2960) to fix hanging in `docker start --attach` and remove spurious `Unsupported signal: <nil>. Discarding` messages. [docker/cli#2987](https://github.com/docker/cli/pull/2987).

## 20.10.4
2021-02-26

### Builder

* Fix incorrect cache match for inline cache import with empty layers [moby/moby#42061](https://github.com/moby/moby/pull/42061)
* Update BuildKit to v0.8.2 [moby/moby#42061](https://github.com/moby/moby/pull/42061)
  * resolver: avoid error caching on token fetch
  * fileop: fix checksum to contain indexes of inputs preventing certain cache misses
  * Fix reference count issues on typed errors with mount references (fixing `invalid mutable ref` errors)
  * git: set token only for main remote access allowing cloning submodules with different credentials
* Ensure blobs get deleted in /var/lib/docker/buildkit/content/blobs/sha256 after pull. To clean up old state run `builder prune` [moby/moby#42065](https://github.com/moby/moby/pull/42065)
* Fix parallel pull synchronization regression [moby/moby#42049](https://github.com/moby/moby/pull/42049)
* Ensure libnetwork state files do not leak [moby/moby#41972](https://github.com/moby/moby/pull/41972)

### Client

* Fix a panic on `docker login` if no config file is present [docker/cli#2959](https://github.com/docker/cli/pull/2959)
* Fix `WARNING: Error loading config file: .dockercfg: $HOME is not defined` [docker/cli#2958](https://github.com/docker/cli/pull/2958)

### Runtime

* docker info: silence unhandleable warnings [moby/moby#41958](https://github.com/moby/moby/pull/41958)
* Avoid creating parent directories for XGlobalHeader [moby/moby#42017](https://github.com/moby/moby/pull/42017)
* Use 0755 permissions when creating missing directories [moby/moby#42017](https://github.com/moby/moby/pull/42017)
* Fallback to manifest list when no platform matches in image config [moby/moby#42045](https://github.com/moby/moby/pull/42045) [moby/moby#41873](https://github.com/moby/moby/pull/41873)
* Fix a daemon panic on setups with a custom default runtime configured [moby/moby#41974](https://github.com/moby/moby/pull/41974)
* Fix a panic when daemon configuration is empty [moby/moby#41976](https://github.com/moby/moby/pull/41976)
* Fix daemon panic when starting container with invalid device cgroup rule [moby/moby#42001](https://github.com/moby/moby/pull/42001)
* Fix userns-remap option when username & UID match [moby/moby#42013](https://github.com/moby/moby/pull/42013)
* static: update runc binary to v1.0.0-rc93 [moby/moby#42014](https://github.com/moby/moby/pull/42014)

### Logger

* Honor `labels-regex` config even if `labels` is not set [moby/moby#42046](https://github.com/moby/moby/pull/42046)
* Handle long log messages correctly preventing awslogs in non-blocking mode to split events bigger than 16kB [mobymoby#41975](https://github.com/moby/moby/pull/41975)

### Rootless

* Prevent the service hanging when stopping by setting systemd KillMode to mixed [moby/moby#41956](https://github.com/moby/moby/pull/41956)
* dockerd-rootless.sh: add typo guard [moby/moby#42070](https://github.com/moby/moby/pull/42070)
* Update rootlesskit to v0.13.1 to fix handling of IPv6 addresses [moby/moby#42025](https://github.com/moby/moby/pull/42025)
* allow mknodding FIFO inside userns [moby/moby#41957](https://github.com/moby/moby/pull/41957)

### Security

* profiles: seccomp: update to Linux 5.11 syscall list [moby/moby#41971](https://github.com/moby/moby/pull/41971)

### Swarm

* Fix issue with heartbeat not persisting upon restart [moby/moby#42060](https://github.com/moby/moby/pull/42060)
* Fix potential stalled tasks [moby/moby#42060](https://github.com/moby/moby/pull/42060)
* Fix `--update-order` and `--rollback-order` flags when only `--update-order` or `--rollback-order` is provided [docker/cli#2963](https://github.com/docker/cli/pull/2963)
* Fix `docker service rollback` returning a non-zero exit code in some situations [docker/cli#2964](https://github.com/docker/cli/pull/2964)
* Fix inconsistent progress-bar direction on `docker service rollback` [docker/cli#2964](https://github.com/docker/cli/pull/2964)


## 20.10.3
2021-02-01

### Security

* [CVE-2021-21285](https://github.com/moby/moby/security/advisories/GHSA-6fj5-m822-rqx8) Prevent an invalid image from crashing docker daemon
* [CVE-2021-21284](https://github.com/moby/moby/security/advisories/GHSA-7452-xqpj-6rpc) Lock down file permissions to prevent remapped root from accessing docker state
* Ensure AppArmor and SELinux profiles are applied when building with BuildKit

### Client

* Check contexts before importing them to reduce risk of extracted files escaping context store
* Windows: prevent executing certain binaries from current directory [docker/cli#2950](https://github.com/docker/cli/pull/2950)

## 20.10.2
2021-01-04

### Runtime

- Fix a daemon start up hang when restoring containers with restart policies but that keep failing to start [moby/moby#41729](https://github.com/moby/moby/pull/41729)
- overlay2: fix an off-by-one error preventing to build or run containers when data-root is 24-bytes long [moby/moby#41830](https://github.com/moby/moby/pull/41830)
- systemd: send `sd_notify STOPPING=1` when shutting down [moby/moby#41832](https://github.com/moby/moby/pull/41832)

### Networking

- Fix IPv6 port forwarding [moby/moby#41805](https://github.com/moby/moby/pull/41805) [moby/libnetwork#2604](https://github.com/moby/libnetwork/pull/2604)

### Swarm

- Fix filtering for `replicated-job` and `global-job` service modes [moby/moby#41806](https://github.com/moby/moby/pull/41806)

### Packaging

- buildx updated to [v0.5.1](https://github.com/docker/buildx/releases/tag/v0.5.1) [docker/docker-ce-packaging#516](https://github.com/docker/docker-ce-packaging/pull/516)

## 20.10.1
2020-12-14

### Builder

- buildkit: updated to [v0.8.1](https://github.com/moby/buildkit/releases/tag/v0.8.1) with various bugfixes [moby/moby#41793](https://github.com/moby/moby/pull/41793)

### Packaging

- Revert a change in the systemd unit that could prevent docker from starting due to a startup order conflict [docker/docker-ce-packaging#514](https://github.com/docker/docker-ce-packaging/pull/514)
- buildx updated to [v0.5.0](https://github.com/docker/buildx/releases/tag/v0.5.0) [docker/docker-ce-packaging#515](https://github.com/docker/docker-ce-packaging/pull/515)

## 20.10.0
2020-12-08

### Deprecation / Removal

For an overview of all deprecated features, refer to the [Deprecated Engine Features](/engine/deprecated/) page.

- Warnings and deprecation notice when `docker pull`-ing from non-compliant registries not supporting pull-by-digest [docker/cli#2872](https://github.com/docker/cli/pull/2872)
- Sterner warnings and deprecation notice for unauthenticated tcp access [moby/moby#41285](https://github.com/moby/moby/pull/41285)
- Deprecate KernelMemory (`docker run --kernel-memory`) [moby/moby#41254](https://github.com/moby/moby/pull/41254) [docker/cli#2652](https://github.com/docker/cli/pull/2652)
- Deprecate `aufs` storage driver [docker/cli#1484](https://github.com/docker/cli/pull/1484)
- Deprecate host-discovery and overlay networks with external k/v stores [moby/moby#40614](https://github.com/moby/moby/pull/40614) [moby/moby#40510](https://github.com/moby/moby/pull/40510)
- Deprecate Dockerfile legacy 'ENV name value' syntax, use `ENV name=value` instead [docker/cli#2743](https://github.com/docker/cli/pull/2743)
- Remove deprecated "filter" parameter for API v1.41 and up [moby/moby#40491](https://github.com/moby/moby/pull/40491)
- Disable distribution manifest v2 schema 1 on push [moby/moby#41295](https://github.com/moby/moby/pull/41295)
- Remove hack MalformedHostHeaderOverride breaking old docker clients (<= 1.12) in which case, set `DOCKER_API_VERSION` [moby/moby#39076](https://github.com/moby/moby/pull/39076)
- Remove "docker engine" subcommands [docker/cli#2207](https://github.com/docker/cli/pull/2207)
- Remove experimental "deploy" from "dab" files [docker/cli#2216](https://github.com/docker/cli/pull/2216)
- Remove deprecated `docker search --automated` and `--stars` flags [docker/cli#2338](https://github.com/docker/cli/pull/2338)
- No longer allow reserved namespaces in engine labels [docker/cli#2326](https://github.com/docker/cli/pull/2326)

### API

- Update API version to v1.41
- Do not require "experimental" for metrics API [moby/moby#40427](https://github.com/moby/moby/pull/40427)
- `GET /events` now returns `prune` events after pruning resources have completed [moby/moby#41259](https://github.com/moby/moby/pull/41259)
  - Prune events are returned for `container`, `network`, `volume`, `image`, and `builder`, and have a `reclaimed` attribute, indicating the amount of space reclaimed (in bytes)
- Add `one-shot` stats option to not prime the stats [moby/moby#40478](https://github.com/moby/moby/pull/40478)
- Adding OS version info to the system info's API (`/info`) [moby/moby#38349](https://github.com/moby/moby/pull/38349)
- Add DefaultAddressPools to docker info [moby/moby#40714](https://github.com/moby/moby/pull/40714)
- Add API support for PidsLimit on services [moby/moby#39882](https://github.com/moby/moby/pull/39882)

### Builder

- buildkit,dockerfile: Support for `RUN --mount` options without needing to specify experimental dockerfile `#syntax` directive. [moby/buildkit#1717](https://github.com/moby/buildkit/pull/1717)
- dockerfile: `ARG` command now supports defining multiple build args on the same line similarly to `ENV` [moby/buildkit#1692](https://github.com/moby/buildkit/pull/1692)
- dockerfile: `--chown` flag in `ADD` now allows parameter expansion [moby/buildkit#1473](https://github.com/moby/buildkit/pull/1473)
- buildkit: Fetching authorization tokens has been moved to client-side (if the client supports it). Passwords do not leak into the build daemon anymore and users can see from build output when credentials or tokens are accessed. [moby/buildkit#1660](https://github.com/moby/buildkit/pull/1660)
- buildkit: Connection errors while communicating with the registry for push and pull now trigger a retry [moby/buildkit#1791](https://github.com/moby/buildkit/pull/1791)
- buildkit: Git source now supports token authentication via build secrets [moby/moby#41234](https://github.com/moby/moby/pull/41234) [docker/cli#2656](https://github.com/docker/cli/pull/2656) [moby/buildkit#1533](https://github.com/moby/buildkit/pull/1533)
- buildkit: Building from git source now supports forwarding SSH socket for authentication [moby/buildkit#1782](https://github.com/moby/buildkit/pull/1782)
- buildkit: Avoid builds that generate excessive logs to cause a crash or slow down the build. Clipping is performed if needed. [moby/buildkit#1754](https://github.com/moby/buildkit/pull/1754)
- buildkit: Change default Seccomp profile to the one provided by Docker [moby/buildkit#1807](https://github.com/moby/buildkit/pull/1807)
- buildkit: Support for exposing SSH agent socket on Windows has been improved [moby/buildkit#1695](https://github.com/moby/buildkit/pull/1695)
- buildkit: Disable truncating by default when using --progress=plain [moby/buildkit#1435](https://github.com/moby/buildkit/pull/1435)
- buildkit: Allow better handling client sessions dropping while it is being shared by multiple builds [moby/buildkit#1551](https://github.com/moby/buildkit/pull/1551)
- buildkit: secrets: allow providing secrets with env [moby/moby#41234](https://github.com/moby/moby/pull/41234) [docker/cli#2656](https://github.com/docker/cli/pull/2656) [moby/buildkit#1534](https://github.com/moby/buildkit/pull/1534)
  - Support `--secret id=foo,env=MY_ENV` as an alternative for storing a secret value to a file.
  - `--secret id=GIT_AUTH_TOKEN` will load env if it exists and the file does not.
- buildkit: Support for mirrors fallbacks, insecure TLS and custom TLS config [moby/moby#40814](https://github.com/moby/moby/pull/40814)
- buildkit: remotecache: Only visit each item once when walking results [moby/moby#41234](https://github.com/moby/moby/pull/41234) [moby/buildkit#1577](https://github.com/moby/buildkit/pull/1577)
  - Improves performance and CPU use on bigger graphs
- buildkit: Check remote when local image platform doesn't match [moby/moby#40629](https://github.com/moby/moby/pull/40629)
- buildkit: image export: Use correct media type when creating new layer blobs [moby/moby#41234](https://github.com/moby/moby/pull/41234) [moby/buildkit#1541](https://github.com/moby/buildkit/pull/1541)
- buildkit: progressui: fix logs time formatting [moby/moby#41234](https://github.com/moby/moby/pull/41234) [docker/cli#2656](https://github.com/docker/cli/pull/2656) [moby/buildkit#1549](https://github.com/moby/buildkit/pull/1549)
- buildkit: mitigate containerd issue on parallel push [moby/moby#41234](https://github.com/moby/moby/pull/41234) [moby/buildkit#1548](https://github.com/moby/buildkit/pull/1548)
- buildkit: inline cache: fix handling of duplicate blobs [moby/moby#41234](https://github.com/moby/moby/pull/41234) [moby/buildkit#1568](https://github.com/moby/buildkit/pull/1568)
  - Fixes https://github.com/moby/buildkit/issues/1388 cache-from working unreliably
  - Fixes https://github.com/moby/moby/issues/41219 Image built from cached layers is missing data
- Allow ssh:// for remote context URLs [moby/moby#40179](https://github.com/moby/moby/pull/40179)
- builder: remove legacy build's session handling (was experimental) [moby/moby#39983](https://github.com/moby/moby/pull/39983)

### Client

- Add swarm jobs support to CLI [docker/cli#2262](https://github.com/docker/cli/pull/2262)
- Add `-a/--all-tags` to docker push [docker/cli#2220](https://github.com/docker/cli/pull/2220)
- Add support for Kubernetes username/password auth [docker/cli#2308](https://github.com/docker/cli/pull/2308)
- Add `--pull=missing|always|never` to `run` and `create` commands [docker/cli#1498](https://github.com/docker/cli/pull/1498)
- Add `--env-file` flag to `docker exec` for parsing environment variables from a file [docker/cli#2602](https://github.com/docker/cli/pull/2602)
- Add shorthand `-n` for `--tail` option [docker/cli#2646](https://github.com/docker/cli/pull/2646)
- Add log-driver and options to service inspect "pretty" format [docker/cli#1950](https://github.com/docker/cli/pull/1950)
- docker run: specify cgroup namespace mode with `--cgroupns` [docker/cli#2024](https://github.com/docker/cli/pull/2024)
- `docker manifest rm` command to remove manifest list draft from local storage [docker/cli#2449](https://github.com/docker/cli/pull/2449)
- Add "context" to "docker version" and "docker info" [docker/cli#2500](https://github.com/docker/cli/pull/2500)
- Propagate platform flag to container create API [docker/cli#2551](https://github.com/docker/cli/pull/2551)
- The `docker ps --format` flag now has a `.State` placeholder to print the container's state without additional details about uptime and health check [docker/cli#2000](https://github.com/docker/cli/pull/2000)
- Add support for docker-compose schema v3.9 [docker/cli#2073](https://github.com/docker/cli/pull/2073)
- Add support for docker push `--quiet` [docker/cli#2197](https://github.com/docker/cli/pull/2197)
- Hide flags that are not supported by BuildKit, if BuildKit is enabled [docker/cli#2123](https://github.com/docker/cli/pull/2123)
- Update flag description for `docker rm -v` to clarify the option only removes anonymous (unnamed) volumes [docker/cli#2289](https://github.com/docker/cli/pull/2289)
- Improve tasks printing for docker services [docker/cli#2341](https://github.com/docker/cli/pull/2341)
- docker info: list CLI plugins alphabetically [docker/cli#2236](https://github.com/docker/cli/pull/2236)
- Fix order of processing of `--label-add/--label-rm`, `--container-label-add/--container-label-rm`, and `--env-add/--env-rm` flags on `docker service update` to allow replacing existing values [docker/cli#2668](https://github.com/docker/cli/pull/2668)
- Fix `docker rm --force` returning a non-zero exit code if one or more containers did not exist [docker/cli#2678](https://github.com/docker/cli/pull/2678)
- Improve memory stats display by using `total_inactive_file` instead of `cache` [docker/cli#2415](https://github.com/docker/cli/pull/2415)
- Mitigate against YAML files that has excessive aliasing [docker/cli#2117](https://github.com/docker/cli/pull/2117)
- Allow using advanced syntax when setting a config or secret with only the source field [docker/cli#2243](https://github.com/docker/cli/pull/2243)
- Fix reading config files containing `username` and `password` auth even if `auth` is empty [docker/cli#2122](https://github.com/docker/cli/pull/2122)
- docker cp: prevent NPE when failing to stat destination [docker/cli#2221](https://github.com/docker/cli/pull/2221)
- config: preserve ownership and permissions on configfile [docker/cli#2228](https://github.com/docker/cli/pull/2228)

### Logging

- Support reading `docker logs` with all logging drivers (best effort) [moby/moby#40543](https://github.com/moby/moby/pull/40543)
- Add `splunk-index-acknowledgment` log option to work with Splunk HECs with index acknowledgment enabled [moby/moby#39987](https://github.com/moby/moby/pull/39987)
- Add partial metadata to journald logs [moby/moby#41407](https://github.com/moby/moby/pull/41407)
- Reduce allocations for logfile reader [moby/moby#40796](https://github.com/moby/moby/pull/40796)
- Fluentd: add fluentd-async, fluentd-request-ack, and deprecate fluentd-async-connect [moby/moby#39086](https://github.com/moby/moby/pull/39086)

### Runtime

- Support cgroup2 [moby/moby#40174](https://github.com/moby/moby/pull/40174) [moby/moby#40657](https://github.com/moby/moby/pull/40657) [moby/moby#40662](https://github.com/moby/moby/pull/40662)
- cgroup2: use "systemd" cgroup driver by default when available [moby/moby#40846](https://github.com/moby/moby/pull/40846)
- new storage driver: fuse-overlayfs [moby/moby#40483](https://github.com/moby/moby/pull/40483)
- Update containerd binary to v1.4.3 [moby/moby#41732](https://github.com/moby/moby/pull/41732)
- `docker push` now defaults to `latest` tag instead of all tags [moby/moby#40302](https://github.com/moby/moby/pull/40302)
- Added ability to change the number of reconnect attempts during connection loss while pulling an image by adding max-download-attempts to the config file [moby/moby#39949](https://github.com/moby/moby/pull/39949)
- Add support for containerd v2 shim by using the now default `io.containerd.runc.v2` runtime [moby/moby#41182](https://github.com/moby/moby/pull/41182)
- cgroup v1: change the default runtime to io.containerd.runc.v2. Requires containerd v1.3.0 or later. v1.3.5 or later is recommended [moby/moby#41210](https://github.com/moby/moby/pull/41210)
- Start containers in their own cgroup namespaces [moby/moby#38377](https://github.com/moby/moby/pull/38377)
- Enable DNS Lookups for CIFS Volumes [moby/moby#39250](https://github.com/moby/moby/pull/39250)
- Use MemAvailable instead of MemFree to estimate actual available memory [moby/moby#39481](https://github.com/moby/moby/pull/39481)
- The `--device` flag in `docker run` will now be honored when the container is started in privileged mode [moby/moby#40291](https://github.com/moby/moby/pull/40291)
- Enforce reserved internal labels [moby/moby#40394](https://github.com/moby/moby/pull/40394)
- Raise minimum memory limit to 6M, to account for higher memory use by runtimes during container startup [moby/moby#41168](https://github.com/moby/moby/pull/41168)
- vendor runc v1.0.0-rc92 [moby/moby#41344](https://github.com/moby/moby/pull/41344) [moby/moby#41317](https://github.com/moby/moby/pull/41317)
- info: add warnings about missing blkio cgroup support [moby/moby#41083](https://github.com/moby/moby/pull/41083)
- Accept platform spec on container create [moby/moby#40725](https://github.com/moby/moby/pull/40725)
- Fix handling of looking up user- and group-names with spaces [moby/moby#41377](https://github.com/moby/moby/pull/41377)

### Networking

- Support host.docker.internal in dockerd on Linux [moby/moby#40007](https://github.com/moby/moby/pull/40007)
- Include IPv6 address of linked containers in /etc/hosts [moby/moby#39837](https://github.com/moby/moby/pull/39837)
- `--ip6tables` enables IPv6 iptables rules (only if experimental) [moby/moby#41622](https://github.com/moby/moby/pull/41622)
- Add alias for hostname if hostname != container name [moby/moby#39204](https://github.com/moby/moby/pull/39204)
- Better selection of DNS server (with systemd) [moby/moby#41022](https://github.com/moby/moby/pull/41022)
- Add docker interfaces to firewalld docker zone [moby/moby#41189](https://github.com/moby/moby/pull/41189) [moby/libnetwork#2548](https://github.com/moby/libnetwork/pull/2548)
  - Fixes DNS issue on CentOS8 [docker/for-linux#957](https://github.com/docker/for-linux/issues/957)
  - Fixes Port Forwarding on RHEL 8 with Firewalld running with FirewallBackend=nftables [moby/libnetwork#2496](https://github.com/moby/libnetwork/issues/2496)
- Fix an issue reporting 'failed to get network during CreateEndpoint' [moby/moby#41189](https://github.com/moby/moby/pull/41189) [moby/libnetwork#2554](https://github.com/moby/libnetwork/pull/2554)
- Log error instead of disabling IPv6 router advertisement failed [moby/moby#41189](https://github.com/moby/moby/pull/41189) [moby/libnetwork#2563](https://github.com/moby/libnetwork/pull/2563)
- No longer ignore `--default-address-pool` option in certain cases [moby/moby#40711](https://github.com/moby/moby/pull/40711)
- Produce an error with invalid address pool [moby/moby#40808](https://github.com/moby/moby/pull/40808) [moby/libnetwork#2538](https://github.com/moby/libnetwork/pull/2538)
- Fix `DOCKER-USER` chain not created when IPTableEnable=false [moby/moby#40808](https://github.com/moby/moby/pull/40808) [moby/libnetwork#2471](https://github.com/moby/libnetwork/pull/2471)
- Fix panic on startup in systemd environments [moby/moby#40808](https://github.com/moby/moby/pull/40808) [moby/libnetwork#2544](https://github.com/moby/libnetwork/pull/2544)
- Fix issue preventing containers to communicate over macvlan internal network [moby/moby#40596](https://github.com/moby/moby/pull/40596) [moby/libnetwork#2407](https://github.com/moby/libnetwork/pull/2407)
- Fix InhibitIPv4 nil panic [moby/moby#40596](https://github.com/moby/moby/pull/40596)
- Fix VFP leak in Windows overlay network deletion [moby/moby#40596](https://github.com/moby/moby/pull/40596) [moby/libnetwork#2524](https://github.com/moby/libnetwork/pull/2524)

### Packaging

- docker.service: Add multi-user.target to After= in unit file [moby/moby#41297](https://github.com/moby/moby/pull/41297)
- docker.service: Allow socket activation [moby/moby#37470](https://github.com/moby/moby/pull/37470)
- seccomp: Remove dependency in dockerd on libseccomp [moby/moby#41395](https://github.com/moby/moby/pull/41395)

### Rootless

- rootless: graduate from experimental [moby/moby#40759](https://github.com/moby/moby/pull/40759)
- Add dockerd-rootless-setuptool.sh [moby/moby#40950](https://github.com/moby/moby/pull/40950)
- Support `--exec-opt native.cgroupdriver=systemd` [moby/moby#40486](https://github.com/moby/moby/pull/40486)

### Security

- Fix CVE-2019-14271 loading of nsswitch based config inside chroot under Glibc [moby/moby#39612](https://github.com/moby/moby/pull/39612)
- seccomp: Whitelist `clock_adjtime`. `CAP_SYS_TIME` is still required for time adjustment [moby/moby#40929](https://github.com/moby/moby/pull/40929)
- seccomp: Add openat2 and faccessat2 to default seccomp profile [moby/moby#41353](https://github.com/moby/moby/pull/41353)
- seccomp: allow 'rseq' syscall in default seccomp profile [moby/moby#41158](https://github.com/moby/moby/pull/41158)
- seccomp: allow syscall membarrier [moby/moby#40731](https://github.com/moby/moby/pull/40731)
- seccomp: whitelist io-uring related system calls [moby/moby#39415](https://github.com/moby/moby/pull/39415)
- Add default sysctls to allow ping sockets and privileged ports with no capabilities [moby/moby#41030](https://github.com/moby/moby/pull/41030)
- Fix seccomp profile for clone syscall [moby/moby#39308](https://github.com/moby/moby/pull/39308)

### Swarm

- Add support for swarm jobs [moby/moby#40307](https://github.com/moby/moby/pull/40307)
- Add capabilities support to stack/service commands [docker/cli#2687](https://github.com/docker/cli/pull/2687) [docker/cli#2709](https://github.com/docker/cli/pull/2709) [moby/moby#39173](https://github.com/moby/moby/pull/39173) [moby/moby#41249](https://github.com/moby/moby/pull/41249)
- Add support for sending down service Running and Desired task counts [moby/moby#39231](https://github.com/moby/moby/pull/39231)
- service: support `--mount type=bind,bind-nonrecursive` [moby/moby#38788](https://github.com/moby/moby/pull/38788)
- Support ulimits on Swarm services. [moby/moby#41284](https://github.com/moby/moby/pull/41284) [docker/cli#2712](https://github.com/docker/cli/pull/2712)
- Fixed an issue where service logs could leak goroutines on the worker [moby/moby#40426](https://github.com/moby/moby/pull/40426)
