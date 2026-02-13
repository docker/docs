---
title: Architecture
description: Technical architecture of Docker Sandboxes including microVM isolation, private Docker daemon, and workspace syncing.
weight: 60
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This page explains how Docker Sandboxes works and the design decisions behind
it.

## Why microVMs?

AI coding agents need to build images, run containers, and use Docker Compose.
Giving an agent access to your host Docker daemon means it can see your
containers, pull images, and run workloads directly on your system. That's too
much access for autonomous code execution.

Running the agent in a container doesn't solve this. Containers share the host
kernel (or in the case of Docker Desktop, share the same virtual machine) and
can't safely isolate something that needs its own Docker daemon.
Docker-in-Docker approaches either compromise isolation (privileged mode with
host socket mounting) or create nested daemon complexity.

MicroVMs provide the isolation boundary needed. Each sandbox gets its own VM
with a private Docker daemon. The agent can build images, start containers, and
run tests without any access to your host Docker environment. When you remove
the sandbox, everything inside - images, containers, packages - is gone.

## Isolation model

### Private Docker daemon per sandbox

Each sandbox runs a complete Docker daemon inside its VM. This daemon is
isolated from your host and from other sandboxes.

```plaintext
Host system (your Docker Desktop)
  ├── Your containers and images
  │
  ├── Sandbox VM 1
  │   ├── Docker daemon (isolated)
  │   ├── Agent container
  │   └── Other containers (created by agent)
  │
  └── Sandbox VM 2
      ├── Docker daemon (isolated)
      └── Agent container
```

When an agent runs `docker build` or `docker compose up`, those commands
execute inside the sandbox using the private daemon. The agent sees only
containers it creates. It cannot access your host containers, images, or
volumes.

This architecture solves a fundamental constraint: autonomous agents need full
Docker capabilities but cannot safely share your Docker daemon.

### Hypervisor-level isolation

Sandboxes use your system's native virtualization:

- macOS: virtualization.framework
- Windows: Hyper-V {{< badge color=violet text=Experimental >}}

This provides hypervisor-level isolation between the sandbox and your host.
Unlike containers (which share the host kernel), VMs have separate kernels and
cannot access host resources outside their defined boundaries.

### What this means for security

The VM boundary provides:

- Process isolation - Agent processes run in a separate kernel
- Filesystem isolation - Only your workspace is accessible
- Network isolation - Sandboxes cannot reach each other
- Docker isolation - No access to host daemon, containers, or images

Network filtering adds an additional control layer for HTTP/HTTPS traffic. See
[Network policies](network-policies.md) for details on that mechanism.

## Workspace syncing

### Bidirectional file sync

Your workspace syncs to the sandbox at the same absolute path:

- Host: `/Users/alice/projects/myapp`
- Sandbox: `/Users/alice/projects/myapp`

Changes sync both ways. Edit a file on your host, and the agent sees it. The
agent modifies a file, and you see the change on your host.

This is file synchronization, not volume mounting. Files are copied between
host and VM. This approach works across different filesystems and maintains
consistent paths regardless of platform differences.

### Path preservation

Preserving absolute paths means:

- File paths in error messages match between host and sandbox
- Hard-coded paths in configuration files work correctly
- Build outputs reference paths you can find on your host

The agent sees the same directory structure you see, reducing confusion when
debugging issues or reviewing changes.

## Storage and persistence

### What persists

When you create a sandbox, these persist until you remove it:

- Docker images and containers - Built or pulled by the agent
- Installed packages - System packages added with apt, yum, etc.
- Agent state - Credentials, configuration, history
- Workspace changes - Files created or modified sync back to host

### What's ephemeral

Sandboxes are lightweight but not stateless. They persist between runs but are
isolated from each other. Each sandbox maintains its own:

- Docker daemon state
- Image cache
- Package installations

When you remove a sandbox with `docker sandbox rm`, the entire VM and its
contents are deleted. Images built inside the sandbox, packages installed, and
any state not synced to your workspace are gone.

### Disk usage

Each sandbox consumes disk space for:

- VM disk image (grows as you build images and install packages)
- Docker images pulled or built inside the sandbox
- Container layers and volumes

Multiple sandboxes do not share images or layers. Each has its own isolated
Docker daemon and storage.

## Networking

### Internet access

Sandboxes have outbound internet access through your host's network connection.
Agents can install packages, pull images, and access APIs.

An HTTP/HTTPS filtering proxy runs on your host and is available at
`host.docker.internal:3128`. Agents automatically use this proxy for outbound
web requests. You can configure network policies to control which destinations
are allowed. See [Network policies](network-policies.md).

### Credential injection

The HTTP/HTTPS proxy automatically injects credentials into API requests for
supported providers (OpenAI, Anthropic, Google, GitHub, etc.). When you set
environment variables like `OPENAI_API_KEY` or `ANTHROPIC_API_KEY` on your
host, the proxy intercepts outbound requests to those services and adds the
appropriate authentication headers.

This approach keeps credentials on your host system - they're never stored
inside the sandbox VM. The agent makes API requests without credentials, and
the proxy injects them transparently. When the sandbox is removed, no
credentials remain inside.

For multi-provider agents (OpenCode, cagent), the proxy automatically selects
the correct credentials based on the API endpoint being called. See individual
[agent configuration](agents/) for credential setup instructions.

When building custom templates or installing agents manually in the shell
sandbox, some agents may require environment variables like `OPENAI_API_KEY`
to be set before they start. Set these to placeholder values (e.g.,
`proxy-managed`) if needed - the proxy will inject actual credentials
regardless of the environment variable value.

### Sandbox isolation

Sandboxes cannot communicate with each other. Each VM has its own private
network namespace. An agent in one sandbox cannot reach services or containers
in another sandbox.

Sandboxes also cannot access your host's `localhost` services. The VM boundary
prevents direct access to services running on your host machine.

## Lifecycle

### Creating and running

`docker sandbox run` initializes a VM with a workspace for a specified agent,
and starts the agent inside an existing sandbox. You can stop and restart the
agent without recreating the VM, preserving installed packages and Docker
images.

`docker sandbox create` initializes the VM with a workspace but doesn't start
the agent automatically. This separates environment setup from agent execution.

### State management

Sandboxes persist until explicitly removed. Stopping an agent doesn't delete
the VM. This means:

- Installed packages remain available
- Built images stay cached
- Environment setup persists between runs

Use `docker sandbox rm` to delete a sandbox and reclaim disk space.

## Comparison to alternatives

Understanding when to use sandboxes versus other approaches:

| Approach                    | Isolation         | Agent Docker access      | Host impact                         | Use case                                      |
| --------------------------- | ----------------- | ------------------------ | ----------------------------------- | --------------------------------------------- |
| Sandboxes (microVMs)        | Hypervisor-level  | Private daemon           | None - fully isolated               | Autonomous agents building/running containers |
| Container with socket mount | Kernel namespaces | Host daemon (shared)     | Agent sees all host containers      | Trusted tools that need Docker CLI            |
| Docker-in-Docker            | Nested containers | Private daemon (complex) | Moderate - privileged mode required | CI/CD environments                            |
| Host execution              | None              | Host daemon              | Full - direct system access         | Manual development by trusted humans          |

Sandboxes trade higher resource overhead (VM + daemon) for complete isolation.
Use containers when you need lightweight packaging without Docker access. Use
sandboxes when you need to give something autonomous full Docker capabilities
without trusting it with your host environment.

## Next steps

- [Network policies](network-policies.md)
- [Custom templates](templates.md)
- [Using sandboxes effectively](workflows.md)
