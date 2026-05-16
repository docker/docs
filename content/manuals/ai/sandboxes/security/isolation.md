---
title: Isolation layers
weight: 10
description: How Docker Sandboxes isolate AI agents using hypervisor, network, Docker Engine, and credential boundaries.
keywords: docker sandboxes, isolation, hypervisor, network, credentials
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

AI coding agents need to execute code, install packages, and run tools on
your behalf. Docker Sandboxes run each agent in its own microVM with four
isolation layers: hypervisor, network, Docker Engine, and credential proxy.

## Hypervisor isolation

Every sandbox runs inside a lightweight microVM with its own Linux kernel.
Unlike containers, which share the host kernel, a sandbox VM cannot access host
processes, files, or resources outside its defined boundaries.

- **Process isolation:** separate kernel per sandbox; processes inside the VM
  are invisible to your host and to other sandboxes
- **Filesystem isolation:** only your workspace directory is shared with the
  host. The rest of the VM filesystem persists across restarts but is removed
  when you delete the sandbox. Symlinks pointing outside the workspace scope
  are not followed.
- **Full cleanup:** when you remove a sandbox with `sbx rm`, the VM and
  everything inside it is deleted

The agent runs as a non-root user with sudo privileges inside the VM. The
hypervisor boundary is the isolation control, not in-VM privilege separation.

## Network isolation

Each sandbox has its own isolated network. Sandboxes cannot communicate with
each other and cannot reach your host's localhost. There is no shared network
between sandboxes or between a sandbox and your host.

All HTTP and HTTPS traffic leaving a sandbox passes through a proxy on your
host that enforces the [network policy](policy.md). The sandbox routes
traffic through either a forward proxy or a transparent proxy depending on the
client's configuration. Both enforce the network policy; only the forward proxy
[injects credentials](credentials.md) for AI services.

Raw TCP connections, UDP, and ICMP are blocked at the network layer. DNS
resolution is handled by the proxy; the sandbox cannot make raw DNS queries.
Traffic to private IP ranges, loopback, and link-local addresses is also
blocked. Only domains explicitly listed in the policy are reachable.

For the default set of allowed domains, see
[Default security posture](defaults.md).

## Docker Engine isolation

Agents often need to build images, run containers, and use Docker Compose.
Mounting your host Docker socket into a container would give the agent full
access to your environment.

Docker Sandboxes avoid this by running a separate [Docker
Engine](https://docs.docker.com/engine/) inside the sandbox environment, isolated from
your host. When the agent runs `docker build` or `docker compose up`, those
commands execute against that engine. The agent has no path to your host Docker
daemon.

```plaintext
Host system
  ├── Host Docker daemon
  │   └── Your containers and images
  │
  └── Sandbox Docker engine (isolated from host)
      ├── [VM] Agent container — sandbox 1
      │    └── [VM] Containers created by agent
      └── [VM] Agent container — sandbox 2
           └── [VM] Containers created by agent
```

## Credential isolation

Most agents need API keys for their model provider. Rather than passing keys
into the sandbox, the host-side proxy intercepts outbound API requests and
injects authentication headers before forwarding each request.

Credential values are never stored inside the VM. They are not available as
environment variables or files inside the sandbox unless you explicitly set
them. This means a compromised sandbox cannot read API keys from the local
environment.

For how to store and manage credentials, see [Credentials](credentials.md).
