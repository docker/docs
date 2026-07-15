---
title: "Sandbox Mode"
description: "Run agents in an isolated Docker sandbox VM for enhanced security."
keywords: docker agent, ai agents, configuration, yaml, sandbox mode
weight: 80
canonical: https://docs.docker.com/ai/docker-agent/configuration/sandbox/
---

_Run agents in an isolated Docker sandbox VM for enhanced security._

## Overview

Sandbox mode runs the entire agent inside a disposable sandbox VM instead of directly on the host system. All shell, filesystem, and process activity happens inside that VM, so a misbehaving agent cannot touch files outside the mounted working directory or reach long-lived host state.

The backend is provided by the [`docker sandbox`](https://docs.docker.com/ai/sandboxes/) CLI plugin (ships with Docker Desktop) or the standalone [`sbx`](https://github.com/docker/sbx) CLI if it is on `PATH`.

> [!NOTE]
> **Requirements**
>
> Sandbox mode requires Docker Desktop with sandbox support (or a working `sbx` CLI). docker-agent shells out to these tools, it does not start raw `docker run` containers.

## Usage

Enable sandbox mode with the `--sandbox` flag on the `docker agent run` command:

```bash
docker agent run --sandbox agent.yaml
```

docker-agent launches a sandbox VM, copies itself into it, mounts the current working directory, and re-runs the agent from inside.

## Flags

| Flag          | Default                                      | Description                                                                                               |
| ------------- | -------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `--sandbox`   | `false`                                      | Enable sandbox mode.                                                                                      |
| `--template`  | `docker/sandbox-templates:docker-agent`      | OCI image used as the sandbox template. Passed to `docker sandbox create -t` / `sbx create -t`.           |
| `--sbx`       | `true`                                       | Prefer the `sbx` CLI backend when it is available. Set `--sbx=false` to always use `docker sandbox`.      |
| `--no-kit`    | `false`                                      | Disable the [auto-kit](#auto-kit) — do not stage skills or prompt files into the sandbox.                 |

```bash
# Use a custom template image
docker agent run --sandbox --template myorg/custom-agent-template:latest agent.yaml

# Force the docker sandbox backend even if sbx is on PATH
docker agent run --sandbox --sbx=false agent.yaml

# Run without staging skills / prompt files into the sandbox
docker agent run --sandbox --no-kit agent.yaml
```

### Always sandbox a given agent

Add `--sandbox` to an [alias](../../features/cli/index.md#docker-agent-alias) so the
sandbox path is taken automatically whenever that alias is invoked:

```bash
docker agent alias add safe-coder agentcatalog/coder --sandbox
docker agent run safe-coder
```

An explicit `--sandbox=false` on the command line still wins, so you can opt
out of the sandbox for a single run without touching the alias.

### Bake the default into the agent config

Agent authors can declare a sandbox default in the YAML itself. Any caller of
the agent then gets the sandbox path automatically, without having to know
(or remember) to pass `--sandbox`:

```yaml
# agent.yaml
runtime:
  sandbox: true

agents:
  root:
    model: openai/gpt-4o
    description: A helpful assistant
    instruction: You are a helpful assistant.
    toolsets:
      - type: shell
```

```bash
docker agent run agent.yaml   # runs in a sandbox automatically
```

The rule is the same as for aliases: an explicit `--sandbox=false` on the
CLI overrides the config default, so you can debug an agent on the host
without editing its YAML.

### Declare a network allowlist

The runner already opens the [tool install hosts](#network-allowlist) and
the [models gateway](#how-it-works) automatically, but agents that talk
to endpoints those resolvers can't infer (custom MCP servers, third-party
APIs, registries not covered by the aqua resolver) would still see a 403
from the sandbox proxy on first contact.

Declare those hosts in `runtime.network_allowlist` and they are unioned
with the inferred set, so the agent can reach them on its first request:

```yaml
# agent.yaml
runtime:
  sandbox: true
  network_allowlist:
    - api.example.com
    - registry.npmjs.org
```

Each entry is a hostname with an optional `:port` suffix. Commas and
whitespace are rejected to keep a single entry from smuggling several
rules into the policy engine. The runner prints the resulting allowlist
before launch so you can audit exactly which hosts the run opens up.

### Persist your own allowlist

For hosts you keep needing across agents (a corporate proxy, a
self-hosted registry, ...) `docker agent sandbox allow` writes the
entry into `~/.config/cagent/config.yaml` once and unions it with the
inferred and agent-declared sets on every subsequent `--sandbox` run:

```bash
# I just got a `Blocked by network policy` 403 on api.example.com.
docker agent sandbox allow api.example.com

# See what's currently persisted.
docker agent sandbox list

# Drop a host you no longer need.
docker agent sandbox deny api.example.com
```

When the kit's per-toolset host resolver fails (the `! using fallback
host set` line in the launch summary), the runner now prints a hint
pointing at this command so you can turn the missing host into a
one-line, persistent fix instead of relying on the wider conservative
fallback host set.

## Example

```yaml
# agent.yaml
agents:
  root:
    model: openai/gpt-4o
    description: Agent with sandboxed shell
    instruction: You are a helpful assistant.
    toolsets:
      - type: shell
```

```bash
docker agent run --sandbox agent.yaml
```

## How It Works

1. `--sandbox` tells docker-agent to prefer the `sbx` CLI (if available and `--sbx` is true), otherwise it falls back to `docker sandbox`.
2. A new sandbox VM is created from the image passed via `--template`.
3. The current working directory is mounted into the VM; the agent binary is copied in.
4. The [auto-kit](#auto-kit) is staged on the host and bind-mounted read-only into the VM, so the agent sees its skills and prompt files inside the sandbox.
5. The default-deny network proxy is opened for the configured [models gateway](../../features/cli/index.md#runtime-configuration-flags) and any package hosts the auto-installer needs for the agent's MCP/LSP toolsets.
6. All tools (shell, filesystem, background jobs, etc.) run inside the VM.
7. When the session ends, docker-agent exits but does not stop or remove the sandbox VM; both the VM and the kit are kept around so subsequent runs from the same workspace can reuse them. A fresh sandbox is created only when the mount set has changed.

## Auto-Kit

The sandbox VM has its own filesystem and `$HOME` — none of the host's `~/.agents/skills/`, `~/.claude/skills/`, project-level `.agents/skills/`, or prompt files like `AGENTS.md` and `CLAUDE.md` are visible inside it. To bridge that gap, docker-agent automatically builds a **kit**: a self-contained directory staged on the host before the sandbox starts and bind-mounted read-only into the VM at the same path.

The kit is built whenever `--sandbox` is used with an agent reference. It is opt-out via `--no-kit`.

### What gets staged

For the agent referenced on the command line, the kit collects:

- **Local [skills](../../features/skills/index.md)** — every `SKILL.md` discovered on the host (global `~/.codex/skills/`, `~/.claude/skills/`, `~/.agents/skills/`, plus project `.claude/skills/` and `.agents/skills/`) is copied under `<kit>/skills/<skill-name>/`. The in-sandbox skills loader reads from the kit instead of the (non-existent) host `$HOME`.
- **Prompt files** — every file referenced via the agent's `add_prompt_files` (`AGENTS.md`, `CLAUDE.md`, …) is collected. Files that already live under the working directory are left alone (the live workspace mount surfaces them); files outside it (e.g. an `AGENTS.md` in `$HOME`) are copied under `<kit>/prompt_files/`.
- **A manifest** — `<kit>/manifest.json` records what was staged. The on-disk copy is sanitised so it cannot be used to map the host filesystem from inside the sandbox.

Before launch, docker-agent prints a summary of what was staged so you can see exactly which skills and prompt files the agent will have access to inside the sandbox.

### Secret redaction

Every text file copied into the kit is run through [portcullis](https://github.com/docker/portcullis), which redacts secrets that match its detection patterns (API keys, tokens, …) in the staged copy. The kit's printed summary marks files as `(redacted)` whenever at least one secret was replaced. Detection is best-effort — portcullis recognises common secret formats but novel or obfuscated tokens may slip through, so the kit is not a substitute for keeping secrets out of skill sources in the first place.

### Network allowlist

The sandbox templates ship with a default-deny network proxy that allows the major model providers but blocks `*.docker.com` and every package-registry / source host the auto-installer reaches for. When the agent declares MCP or LSP toolsets that have a `command` and an installable `version`, the kit build resolves each toolset's package against the [aqua](https://aquaproj.github.io/) registry and computes the minimal set of hosts the in-sandbox auto-installer will need (Go module proxy + toolchain bootstrap for `go_install` packages, GitHub release hosts for `github_release` packages, …). Those hosts, `models.dev` (needed so the in-sandbox agent can resolve model metadata such as context limits, pricing, and capabilities — without it the first catalog lookup fails with a `403 Blocked by network policy` error), and the configured [`--models-gateway`](../../features/cli/index.md#runtime-configuration-flags) — are then allow-listed on the sandbox proxy. If a per-toolset registry lookup fails, a conservative fallback union is used so the run can still succeed; the affected toolsets are surfaced in the printed summary.

### Caching

Kits are stored under the docker-agent cache directory (`~/Library/Caches/cagent/sandbox-kits/<hash>` on macOS) keyed by a content hash of the agent reference. Reusing the same agent across runs reuses the same kit directory in place; disk usage is bounded by the number of distinct agents you have run. Kits are deliberately kept on disk between runs because the reused sandbox VM holds a hard reference to the kit's bind-mount path — deleting it would leave the sandbox un-startable.

### Disabling the kit

Pass `--no-kit` to skip the kit build entirely. The agent then runs without any host-side skills or external prompt files visible inside the sandbox, and the network allowlist falls back to the template defaults. Useful for debugging the sandbox itself, or for agents that don't depend on host skills.

```bash
docker agent run --sandbox --no-kit agent.yaml
```

> [!WARNING]
> **Limitations**
>
> - Sandboxes are reused across runs from the same workspace; if the required mount set changes (e.g. a new kit is staged), the previous sandbox is removed and a fresh one is created.
> - Only the working directory, the agent config directory, and (when staged) the kit directory are mounted; other host files are not visible to the agent.
> - Network egress is constrained by the sandbox backend's default-deny policy plus the per-run allowlist described above.
