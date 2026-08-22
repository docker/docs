---
title: Architecture
weight: 100
description: Technical architecture of Docker Sandboxes; workspace mounting, storage, networking, and sandbox lifecycle.
keywords: docker sandboxes, architecture, microVM, workspace mounting, sandbox lifecycle
---

This page explains how Docker Sandboxes work under the hood. For the security
properties of the architecture, see [Sandbox isolation](security/isolation.md).

## Workspace storage

Starting with `sbx` version 0.40.0, workspace paths are optional for
`sbx create`. When you omit them, the sandbox has no host workspace bind mount.
The agent works in the template image's configured `WORKDIR` instead. The
Docker-provided images for built-in agents set `WORKDIR` to
`/home/agent/workspace`. A custom template can set another absolute path. If
the daemon can't resolve a usable absolute `WORKDIR` from the image config, it
falls back to `/home/agent/workspace`. Files created there stay inside the
sandbox and persist across stops and restarts.

When you pass a workspace path to `sbx create` or `sbx run`, the directory is
mounted into the sandbox through a filesystem passthrough. `sbx run` uses the
current directory when you don't pass a path. The sandbox sees your actual
host files, so changes in either direction are instant with no sync process
involved.

A directly mounted workspace appears at the same absolute path as on your
host. Preserving absolute paths means error messages, configuration files, and
build outputs all reference paths you can find on your host. The agent sees the
same directory structure, which reduces confusion when debugging or reviewing
changes.

Clone mode uses a third storage layout. The host repository is mounted
read-only at `/run/sandbox/source`, and the agent works in a private clone
inside the sandbox. See [Clone mode](usage.md#clone-mode).

> [!WARNING]
> Avoid mounting network-attached or remote storage (network drives, SMB/NFS
> shares, or cloud-synced folders) as a workspace. The sandbox accesses
> workspaces through a filesystem passthrough, so every file read and write
> goes over the network. This adds latency and slows agent performance.

## Storage and persistence

When you create a sandbox, everything inside it persists until you remove it:
Docker images and containers built or pulled by the agent, installed packages,
agent state and history, and files in mountless or cloned workspaces. Files in
a directly mounted workspace live on the host instead.

Each sandbox maintains its own Docker daemon state, image cache, and package
installations. Multiple sandboxes don't share images or layers. The
[shared agent skills store](workflows/agent-skills.md) is an exception:
supported agents mount the same host-side store read-write unless you opt out
when creating the sandbox.

Each sandbox consumes disk space for its VM image, Docker images, container
layers, and volumes, and this grows as you build images and install packages.

Virtiofs caching is enabled by default for directly mounted workspaces on all
operating systems. File reads from the sandbox VM are cached on the host side,
reducing round-trips through the filesystem passthrough and improving
performance for read-heavy workloads such as `git status` or directory scans.
To opt out, set
`DOCKER_SANDBOXES_ENABLE_VIRTIOFS_CACHE=0` when creating the sandbox:

```console
$ DOCKER_SANDBOXES_ENABLE_VIRTIOFS_CACHE=0 sbx run <agent>
```

## Networking

All outbound TCP traffic from the sandbox routes through a proxy on your host.
Agents use a forward proxy for HTTP and HTTPS; other TCP traffic is forwarded
transparently. Both paths enforce
[network access policies](governance/access-controls/network.md). The forward
proxy also handles [credential injection](configuration/credentials.md). See
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
[Configure an upstream proxy](configuration/upstream-proxy.md). Upstream proxy support is
experimental and subject to change.

Only HTTP and HTTPS traffic can be forwarded to an upstream proxy. Other TCP
traffic can't be redirected to a proxy.

## MCP gateway

Supported agents connect to a single MCP gateway endpoint for the sandbox. The
gateway runs on the host side of the sandbox boundary and brokers access to
registered MCP servers.

Registered MCP servers can be remote endpoints, or they can be local stdio
servers launched on the host. Local stdio servers don't run inside the sandbox
VM. If a local stdio server is packaged as an OCI image, or if you register an
explicit `docker` command, it uses Docker on the host.

When MCP policies apply, enforcement happens on the MCP gateway path, separate
from the HTTP/HTTPS network proxy. Server registration is checked before the
server is stored, and governed MCP requests are checked by the gateway before
tool calls, resource reads, prompt retrieval, or gateway meta-tool execution.

## Lifecycle

`sbx run` initializes a VM for a specified agent and starts the agent. You can
stop and restart without recreating the VM, preserving installed packages,
Docker images, and in-sandbox files.

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
