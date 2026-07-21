---
title: Architecture
weight: 40
description: Technical architecture of Docker Sandboxes; workspace mounting, storage, networking, and sandbox lifecycle.
keywords: docker sandboxes, architecture, microVM, workspace mounting, sandbox lifecycle
---

This page explains how Docker Sandboxes work under the hood. For the security
properties of the architecture, see [Sandbox isolation](security/isolation.md).

## Workspace mounting

Your workspace is mounted directly into the sandbox through a filesystem
passthrough. The sandbox sees your actual host files, so changes in either
direction are instant with no sync process involved.

Your workspace is mounted at the same absolute path as on your host. Preserving
absolute paths means error messages, configuration files, and build outputs all
reference paths you can find on your host. The agent sees exactly the directory
structure you see, which reduces confusion when debugging or reviewing changes.

> [!WARNING]
> Avoid mounting network-attached or remote storage (network drives, SMB/NFS
> shares, or cloud-synced folders) as a workspace. The sandbox accesses
> workspaces through a filesystem passthrough, so every file read and write
> goes over the network. This adds latency and slows agent performance.

## Storage and persistence

When you create a sandbox, everything inside it persists until you remove it:
Docker images and containers built or pulled by the agent, installed packages,
agent state and history, and workspace changes.

Sandboxes are isolated from each other. Each one maintains its own Docker
daemon state, image cache, and package installations. Multiple sandboxes don't
share images or layers.

Each sandbox consumes disk space for its VM image, Docker images, container
layers, and volumes, and this grows as you build images and install packages.

Virtiofs caching is enabled by default on all operating systems. File reads
from the sandbox VM are cached on the host side, reducing round-trips through
the filesystem passthrough and improving performance for read-heavy workloads
such as `git status` or directory scans. To opt out, set
`DOCKER_SANDBOXES_ENABLE_VIRTIOFS_CACHE=0` when creating the sandbox:

```console
$ DOCKER_SANDBOXES_ENABLE_VIRTIOFS_CACHE=0 sbx run <template>
```

## Networking

All outbound traffic from the sandbox routes through an HTTP/HTTPS proxy on
your host. Agents are configured to use the proxy automatically. The proxy
enforces [network access policies](governance/) and handles
[credential injection](security/credentials.md). See
[Network isolation](security/isolation.md#network-isolation) for how this
works and [Default security posture](security/defaults.md) for what is
allowed out of the box.

### Upstream proxy

The host-side proxy makes its outbound connections using your host's network
configuration and routing. When a destination is reachable through a direct
route, traffic follows that route. When reaching a destination requires an
upstream proxy, the host-side proxy forwards the request to it. Chaining to an
upstream proxy means sandbox traffic respects the same egress controls as other
applications on your host.

By default, both sandbox traffic and the daemon's own traffic follow your OS
system proxy, so this usually works without any configuration. To set a proxy
explicitly — with a proxy URL, a PAC file, a SOCKS5 proxy, or separate settings
for sandbox and daemon traffic — see
[Configure an upstream proxy](upstream-proxy.md).

Only HTTP and HTTPS traffic can be forwarded to an upstream proxy. Other TCP
traffic can't be redirected to a proxy.

## Lifecycle

`sbx run` initializes a VM with a workspace for a specified agent and starts
the agent. You can stop and restart without recreating the VM, preserving
installed packages and Docker images.

Sandboxes persist until explicitly removed. Stopping an agent doesn't delete
the VM; environment setup carries over between runs. Use `sbx rm` to delete
the sandbox, its VM, and all of its contents. If the sandbox used
[`--clone`](usage.md#clone-mode), the `sandbox-<name>` Git remote is also
removed from your host repository.

## Comparison to alternatives

| Approach                                            | Isolation            | Docker access      | Use case           |
| --------------------------------------------------- | -------------------- | ------------------ | ------------------ |
| Sandboxes (microVMs)                                | Full (hypervisor)    | Isolated daemon    | Autonomous agents  |
| Container with socket mount                         | Partial (namespaces) | Shared host daemon | Trusted tools      |
| [Docker-in-Docker](https://hub.docker.com/_/docker) | Partial (privileged) | Nested daemon      | CI/CD pipelines    |
| Host execution                                      | None                 | Host daemon        | Manual development |

Sandboxes trade higher resource overhead (a VM plus its own daemon) for
complete isolation. Use containers when you need lightweight packaging without
Docker access. Use sandboxes when you need to give something autonomous full
Docker capabilities without trusting it with your host environment.
