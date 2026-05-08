---
title: Isolation layers
weight: 10
description: How Docker Sandboxes isolate AI agents using hypervisor, network, Docker Engine, source-repository, and credential boundaries.
keywords: docker sandboxes, isolation, hypervisor, network, credentials, git
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

AI coding agents need to execute code, install packages, and run tools on
your behalf. Docker Sandboxes run each agent in its own microVM with five
isolation layers: hypervisor, network, Docker Engine, source-repository
(in branch mode), and credential proxy.

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

## Source-repository isolation

When you start a sandbox with `--branch` (see the
[branch-mode workflow](../usage.md#branch-mode)), the agent never works
directly against your host repository. Even with full root inside the VM,
it cannot corrupt your local Git state.

The boundary works like this:

- Your repository's Git root is bind-mounted into the sandbox at
  `/run/sandbox/source` as a read-only mount. The agent — and anything it
  spawns — cannot write to your `.git` directory, your working tree, or
  any tracked file via that mount.
- The agent's working copy is a private `git clone --reference` populated
  on the sandbox's overlay filesystem. The clone has its own index, its
  own refs, and its own working tree. Object storage is shared via
  `.git/objects/info/alternates`, so the clone is space-efficient and
  full history is walkable, but writes to the clone never reach your
  host's object database.
- Your host pulls the agent's commits over a `git-daemon` exposed by the
  sandbox on `127.0.0.1:<ephemeral-port>`. The CLI registers it as a
  `sandbox-<sandbox-name>` remote on your host repository. Fetching from that
  remote uses the same trust model as fetching from any third-party
  remote: nothing is integrated until you explicitly merge or check out
  the fetched refs.

```plaintext
Host repository                            Sandbox VM
  .git/                                      /run/sandbox/source/  (RO bind mount)
    objects/  ◄─────── alternates ─────────  clone/.git/objects/
    refs/                                    clone/.git/refs/      (private)
    HEAD                                     clone/.git/HEAD       (private)
    working tree                             clone/working tree    (overlay FS)
    remote sandbox-<name>  ──── git:// ────► git-daemon :9418
                                             (published 127.0.0.1:<ephemeral>)
```

The practical guarantees:

- Index and ref corruption can't happen — concurrent `git` commands on the
  host and inside the sandbox don't race on a shared `.git/index` or shared
  refs because there is no shared writable state.
- The agent can't write back to your working tree. A compromised or buggy
  agent can't drop a `.git/hooks/pre-commit`, modify `.github/workflows/`,
  or edit any other tracked file in a way that affects your host until you
  fetch and merge from the `sandbox-<name>` remote.
- Credentials, signing keys, and global settings declared in your
  repository's `.git/config` stay on the host. The agent's clone has its
  own independent configuration.
- Cleanup is automatic: `sbx rm` deletes the clone, the published port,
  and the `sandbox-<name>` remote on your host. Nothing leaks outside the
  sandbox lifecycle.

In direct mode (no `--branch`), the agent edits your working tree directly
and this isolation does not apply. Use branch mode whenever you want a
strong boundary between the agent's Git activity and your host
repository.

## Credential isolation

Most agents need API keys for their model provider. Rather than passing keys
into the sandbox, the host-side proxy intercepts outbound API requests and
injects authentication headers before forwarding each request.

Credential values are never stored inside the VM. They are not available as
environment variables or files inside the sandbox unless you explicitly set
them. This means a compromised sandbox cannot read API keys from the local
environment.

For how to store and manage credentials, see [Credentials](credentials.md).
