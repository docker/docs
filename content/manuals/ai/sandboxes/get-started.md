---
title: Get started with Docker Sandboxes
linkTitle: Get started
weight: 10
description: Install the sbx CLI and run an AI coding agent in an isolated sandbox.
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Docker Sandboxes run AI coding agents in isolated microVM sandboxes. Each
sandbox gets its own Docker daemon, filesystem, and network — the agent can
build containers, install packages, and modify files without touching your host
system.

## Prerequisites

- macOS (Apple silicon) or Windows (x86_64, Windows 11 required)
- If you're on Windows, enable Windows Hypervisor Platform. Open an elevated
  PowerShell prompt (Run as Administrator) and run:
  ```powershell
  Enable-WindowsOptionalFeature -Online -FeatureName HypervisorPlatform -All
  ```
- An API key or authentication method for the agent you want to use. Most
  agents require an API key for their model provider (Anthropic, OpenAI,
  Google, and others). See the [agent pages](agents/) for provider-specific
  instructions, and [Credentials](security/credentials.md) for how to store
  and manage keys.

Docker Desktop is not required to use `sbx`.

## Install and sign in

{{< tabs >}}
{{< tab name="macOS" >}}

```console
$ brew install docker/tap/sbx
$ sbx login
```

{{< /tab >}}
{{< tab name="Windows" >}}

```powershell
> winget install -h Docker.sbx
> sbx login
```

{{< /tab >}}
{{< /tabs >}}

`sbx login` opens a browser for Docker OAuth. On first login (and after `sbx
policy reset`), the CLI prompts you to choose a default network policy for your
sandboxes:

```plaintext
Choose a default network policy:

     1. Open         — All network traffic allowed, no restrictions.
     2. Balanced     — Default deny, with common dev sites allowed.
     3. Locked Down  — All network traffic blocked unless you allow it.

Use ↑/↓ to navigate, Enter to select, or press 1–3.
```

See [Policies](security/policy.md) for a full description of each option.

> [!NOTE]
> See the [FAQ](faq.md) for details on why sign-in is required and what
> happens with your data.

## Run your first sandbox

Pick a project directory and launch an agent with [`sbx run`](/reference/cli/sbx/run/):

```console
$ cd ~/my-project
$ sbx run claude
```

Replace `claude` with the agent you want to use — see [Agents](agents/) for the
full list.

The first run takes a little longer while the agent image is pulled.
Subsequent runs reuse the cached image and start in seconds.

You can check what's running at any time:

```console
$ sbx ls
NAME                 STATUS   UPTIME
claude-my-project    running  12s
```

The agent can modify files in your project directory, so review changes before
merging. See [Workspace trust](security/workspace.md) for details.

> [!CAUTION]
> Your network policy controls what the sandbox can reach. With **Locked
> Down**, even your model provider API is blocked. With **Balanced**, a broad
> set of common development services is allowed by default — add other hosts
> with `sbx policy allow`. See [Policies](security/policy.md) for details.

## Next steps

- [Usage guide](usage.md) — common patterns and workflows
- [Agents](agents/) — supported agents and configuration
- [Custom environments](agents/custom-environments.md) — build your own sandbox
  images
- [Policies](security/policy.md) — control outbound access
