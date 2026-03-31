---
title: Docker Desktop sandboxes
linkTitle: Docker Desktop
description: Run sandboxed AI coding agents using Docker Desktop and the docker sandbox CLI.
weight: 80
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

Docker Desktop includes a built-in sandbox feature for running AI coding
agents in isolated microVMs using the `docker sandbox` command. This is a
convenience integration. For full functionality, use the standalone `sbx`
CLI instead.

> [!NOTE]
> The standalone `sbx` CLI provides more features, more flexibility, and doesn't
> require Docker Desktop. If you're setting up sandboxed agents for the first
> time, consider using the standalone CLI instead.

## Prerequisites

- Docker Desktop 4.58 or later
- macOS or Windows
- API keys for your chosen agent

## Quick start

1. Set your API key in your shell configuration file:

   ```plaintext {title="~/.bashrc or ~/.zshrc"}
   export ANTHROPIC_API_KEY=sk-ant-api03-xxxxx
   ```

   Source your shell configuration and restart Docker Desktop so the daemon
   picks up the variable.

2. Create and run a sandbox:

   ```console
   $ cd ~/my-project
   $ docker sandbox run claude
   ```

   The first run takes longer while Docker initializes the microVM.

Replace `claude` with a different [agent](#supported-agents) if needed.

## Supported agents

| Agent                             | Command        | Notes                                |
| --------------------------------- | -------------- | ------------------------------------ |
| Claude Code                       | `claude`       | Most tested implementation           |
| Codex                             | `codex`        |                                      |
| Copilot                           | `copilot`      |                                      |
| Gemini                            | `gemini`       |                                      |
| [Docker Agent](/ai/docker-agent/) | `docker-agent` | Also available as a standalone tool  |
| Kiro                              | `kiro`         |                                      |
| OpenCode                          | `opencode`     |                                      |
| Custom shell                      | `shell`        | Minimal environment for manual setup |

All agents are experimental. The agent type is specified when creating a
sandbox and can't be changed later.

## Authentication

Each agent requires its own API key or credentials. Docker Sandboxes uses a
daemon that doesn't inherit environment variables from your shell session, so
you must set keys in your shell configuration file (not just export them in
your terminal).

Common environment variables by agent:

| Agent        | Environment variable(s)                           |
| ------------ | ------------------------------------------------- |
| Claude Code  | `ANTHROPIC_API_KEY`                               |
| Codex        | `OPENAI_API_KEY`                                  |
| Copilot      | `GH_TOKEN` or `GITHUB_TOKEN`                      |
| Gemini       | `GEMINI_API_KEY` or `GOOGLE_API_KEY`              |
| Docker Agent | `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, and others |
| OpenCode     | `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, and others |
| Kiro         | Device flow (interactive browser login)           |
| Shell        | Any provider keys needed                          |

After setting variables, source your shell configuration and restart Docker
Desktop. The sandbox proxy injects credentials into API requests so keys stay
on your host and are never stored inside the sandbox.

## Commands

```console
$ docker sandbox run AGENT [PATH]                  # Create and run
$ docker sandbox ls                                # List sandboxes
$ docker sandbox exec -it <name> bash              # Shell into a sandbox
$ docker sandbox rm <name>                         # Remove a sandbox
$ docker sandbox reset                             # Remove all sandboxes
$ docker sandbox network proxy <name> --policy … # Set network policy
$ docker sandbox network log                       # View network log
```

Sandboxes don't appear in `docker ps` because they're microVMs, not
containers. For the full command reference, see the
[CLI reference](/reference/cli/docker/sandbox/).

Pass agent-specific CLI options after the sandbox name with a `--` separator:

```console
$ docker sandbox run <name> -- --continue
```

## Architecture

Each sandbox is a lightweight microVM with its own kernel, using your system's
native virtualization (macOS virtualization.framework, Windows Hyper-V). The
default agent templates include a private Docker daemon, so `docker build` and
`docker compose up` run inside the sandbox without affecting your host.

```plaintext
Host system
  ├── Your containers and images
  ├── Sandbox VM 1
  │   ├── Docker daemon (isolated)
  │   ├── Agent container
  │   └── Containers created by agent
  └── Sandbox VM 2
      ├── Docker daemon (isolated)
      └── Agent container
```

Your workspace syncs bidirectionally between host and sandbox at the same
absolute path. Outbound internet goes through an HTTP/HTTPS filtering proxy on
the host. See [Network policies](#network-policies) for configuration.

## Network policies

The filtering proxy controls what a sandbox can access. By default, all
traffic is allowed except private networks and localhost.

Allow mode (block specific destinations):

```console
$ docker sandbox network proxy my-sandbox \
  --policy allow \
  --block-cidr 10.0.0.0/8
```

Deny mode (allow specific destinations):

```console
$ docker sandbox network proxy my-sandbox \
  --policy deny \
  --allow-host api.anthropic.com \
  --allow-host "*.npmjs.org"
```

View what an agent is accessing:

```console
$ docker sandbox network log
```

## Custom templates

Build custom templates to pre-install tools:

```dockerfile
FROM docker/sandbox-templates:claude-code

USER root
RUN apt-get update && apt-get install -y build-essential \
    && rm -rf /var/lib/apt/lists/*
USER agent
```

```console
$ docker build -t my-template:v1 .
$ docker sandbox run -t my-template:v1 claude ~/project
```

## Base environment

All agent templates share a common environment:

- Ubuntu 25.10
- Docker CLI (with Buildx and Compose), Git, GitHub CLI, Node.js, Go, Python 3, uv, make, jq, ripgrep
- Non-root `agent` user with sudo access
- Package managers: apt, pip, npm

## Troubleshooting

<!-- vale off -->

### 'sandbox' is not a docker command

<!-- vale on -->

The CLI plugin isn't installed or isn't in the correct location. Verify the
plugin exists at `~/.docker/cli-plugins/docker-sandbox` and restart Docker
Desktop.

### Beta features need to be enabled

If your Docker Desktop is managed by an administrator with
[Settings Management](/enterprise/security/hardened-desktop/settings-management/),
ask them to
[allow beta features](/enterprise/security/hardened-desktop/settings-management/configure-json-file/#beta-features).

### Authentication failure

Verify your API key is valid and set in your shell configuration file (not
just exported in the current session). Source the file and restart Docker
Desktop.

### Permission denied on workspace files

Check **Docker Desktop** > **Settings** > **Resources** > **File Sharing** and
ensure your workspace path is listed. Verify file permissions with `ls -la`.

### Sandbox crashes on Windows

If launching multiple sandboxes causes crashes, end all `docker.openvmm.exe`
processes in Task Manager and restart Docker Desktop. Launch sandboxes one at a
time.

### Persistent issues

Reset all sandbox state:

```console
$ docker sandbox reset
```

This stops all VMs and deletes all sandbox data. Create fresh sandboxes
afterward.
