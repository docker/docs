---
title: "Sandbox Mode"
description: "Run agents in an isolated Docker sandbox VM."
keywords: docker agent, sbx, ai agents, configuration, yaml, sandbox mode
weight: 80
canonical: https://docs.docker.com/ai/docker-agent/configuration/sandbox/
---

_Run agents in an isolated [Docker sandbox VM](https://docs.docker.com/ai/sandboxes/)._

## Overview

Sandbox mode integrates Docker Agent with [Docker Sandboxes](https://docs.docker.com/ai/sandboxes/).
Docker Agent prefers the standalone `sbx` CLI when available and otherwise
uses the `docker sandbox` CLI plugin. `--sbx=false` forces the plugin backend.
The `--sandbox` flag asks the selected backend to create or reuse a VM and
launches Docker Agent inside it.

For local sandboxes, agent processes run inside the VM and the working directory is mounted read-write.
The Docker Agent configuration directory and staged kit are mounted read-only;
resolved built-in agent definitions are also staged read-only so guest aliases cannot change the selected agent;
for a local agent config outside the working directory, its parent directory is
also mounted read-only. An agent config inside the working directory remains on
the read-write mount. Files exposed through these mounts remain accessible to the agent. Docker Agent does not implement the sandbox or
start a raw `docker run` container; it orchestrates the selected sandbox CLI.

> [!NOTE]
> **Requirements**
>
> Install and configure [Docker Sandboxes](https://docs.docker.com/ai/sandboxes/) before
> using `--sandbox`. The standalone `sbx` CLI is optional if `docker sandbox` is available.
> Cloud mode requires a cloud-enabled build of the selected CLI and a Docker login; it does not require a local Docker Desktop VM. Kits v3 require a v3-capable sbx release.

## Usage

Enable sandbox mode with the `--sandbox` flag on the `docker agent run` command:

```bash
docker agent run --sandbox agent.yaml
```

Docker Agent asks the selected sandbox backend to launch or reuse a sandbox VM, mounts the current
working directory, and runs the agent inside it.

## Flags

| Flag          | Default                                      | Description                                                                                               |
| ------------- | -------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `--sandbox`   | `false`                                      | Enable sandbox mode.                                                                                      |
| `--template`  | `docker/docker-agent-sbx-templates:latest`   | Local OCI image, or an existing cloud template name with `--cloud`. The local default is not forwarded in cloud mode. See [Sandbox templates](#sandbox-templates). |
| `--sbx` | `true` | Prefer standalone `sbx` when available; set `--sbx=false` to force `docker sandbox`. |
| `--cloud` | `false` | Run in the cloud; implies `--sandbox`. No automatic host workspace/config upload. |
| `--sandbox-kit` | unset | Workload kit reference, passed as sbx's positional launch source. Mutually exclusive with `--template`. |
| `--kit` | unset | Additional mixin kit reference; repeatable. |
| `--kit-arg` | unset | Kit argument `KEY=VALUE`; repeatable. Values are passed to sbx unchanged. |
| `--sandbox-ttl` | `1h` | Cloud lifetime (greater than zero, at most 24h); stop on expiry. |
| `--no-kit`    | `false`                                      | Disable the [auto-kit](#auto-kit) — do not stage skills or prompt files into the sandbox.                 |

```bash
# Use a custom template image
docker agent run --sandbox --template myorg/custom-agent-template:latest agent.yaml

# Run without staging skills / prompt files into the sandbox
docker agent run --sandbox --no-kit agent.yaml
```

## Cloud sandboxes

```bash
# Uses sbx's built-in docker-agent workload and cloud credentials.
docker agent run --cloud --exec --model openai/gpt-5.6 default "Explain the sandbox environment"

# Run a published Docker Agent config; it must be self-contained remotely.
docker agent run --cloud --sandbox-ttl 2h myorg/coder:latest

# Use an existing cloud template containing docker-agent (not an OCI image ref).
docker agent run --cloud --template my-cloud-template --model openai/gpt-5.6 default
```

Cloud runs use a **fresh, remote workspace**. Docker Agent does not upload the
current directory, local YAML, user config, host API keys, or auto-kit. Use a
built-in agent, URL, or published OCI agent reference; local-file references
(including aliases to local files) fail before provisioning. Host aliases to
remote references are resolved without copying the alias config. Host-only path
flags such as `--working-dir`, `--attach`, `--prompt-file`, and `--env-from-file`
are rejected. Put remote setup, files, and paths in a workload kit instead.

Built-in cloud agents require an explicit `--model` (or a model set on the host
alias). sbx may expose proxy-managed sentinels for providers without a configured
cloud secret, so automatic `first_available` selection cannot reliably discover
which provider you can use. Published agents should pin their model or accept
an explicit override for the same reason.

Configure secrets with `sbx --cloud secret` and declare egress in your kit/cloud
policy. Host user settings, gateway defaults, inferred tool-host allowances, and
`runtime.network_allowlist` are **not** replayed into cloud policy. This avoids
weakening a kit's restrictions or silently copying local credentials. Any
provider, registry, or tool endpoint the agent needs must be reachable under the
cloud policy. CLI-specified model, safety, flavor, and prompt options still apply.

Docker Agent uses `sbx --cloud run --detached`, then `exec`: this supports sbx's
v3 source builds and multi-kit assembly without launching a second agent TUI.
Provisioning messages go to stderr, leaving stdout for the agent (including
`--json` output). Piped/headless execution does not allocate a TTY.

Each cloud create receives a random invocation-specific name and `--new`, so it
cannot accidentally resume another run. Interrupted provisioning preserves any
returned ID or reconciles that name against `sbx --cloud ls` for cleanup. If the
server outcome is still unknown, the name and recovery instructions are printed;
TTL remains the fallback. No unrelated sandbox is stopped.

On exit, including an agent failure or cancellation, Docker Agent requests a
**stop**, preserving remote files and releasing compute. A one-hour TTL with
`--on-timeout stop` is the fallback if the client disappears; it bounds a run
unless overridden by `--sandbox-ttl`. Sandbox IDs and cleanup hints are printed
to stderr. Stopping is asynchronous: check `sbx --cloud ls` before resuming.
Stopped sandboxes and stored data may still incur storage charges.

To recover files, explicitly restart the stopped sandbox, copy the desired
paths, then stop or remove it:

```bash
sbx --cloud exec sbx_ID true
sbx --cloud cp sbx_ID:/home/agent/workspace ./sandbox-output
sbx --cloud stop sbx_ID
# When the data is no longer needed:
sbx --cloud rm --force sbx_ID
```

There is no automatic copyback over your checkout. `sbx cp` is a snapshot, not a
sync engine; its cloud file API has limitations around symlinks and empty
directories. Use a kit to clone a repository remotely or transfer files explicitly.

## Sandbox kits v3

`--sandbox-kit` selects a workload; repeated `--kit` flags add mixins. References
may be explicit local paths, ZIPs, git references, or OCI artifacts. sbx owns
builds, caching, composition order, argument validation, credential capabilities,
and permission checks. All members of a composition must use the same kit grammar.
The built-in sbx docker-agent workload is currently v2: select a v3 workload when
adding v3 mixins.

```bash
docker agent run --sandbox --sandbox-kit ./my-workload \
  --kit myorg/tooling:3 --kit-arg tool_version=1.2.3 agent.yaml

docker agent run --cloud --sandbox-kit ./my-workload \
  --kit myorg/tooling:3 myorg/coder:latest
```

The workload must provide `docker-agent` on PATH. `--template` cannot be combined
with `--sandbox-kit`; cloud template launches also cannot take mixins. See the
[example v3 workload](https://github.com/docker/docker-agent/tree/main/docs/configuration/sandbox/kit),
which declares network policy, provider credentials, shared skills, agent context,
and headless/resume session commands.

With an explicit workload, Docker Agent disables its host auto-kit so it does
not hide sbx-provided skills/context. It also leaves network policy to the kit
instead of adding post-create allowances. Local Docker gateway authentication
uses a v3 credential/network mixin; legacy default launches retain their v2
login kit. Cloud gateway credentials are never forwarded from the host.

> [!NOTE]
> Current sbx cloud uploaded/assembled-kit paths reject required credential
> capabilities and can drop optional bindings. Shared host skills are also
> unavailable in cloud mode. Follow sbx's diagnostics and configure account
> cloud secrets/policies separately; Docker Agent does not bypass those checks.

### Always sandbox a given agent

Add `--sandbox` to an [alias](../../features/cli/index.md#docker-agent-alias) so the
sandbox path is taken automatically whenever that alias is invoked:

```bash
docker agent alias add safe-coder myorg/coder --sandbox
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
the [models gateway](#how-local-sandbox-mode-works) automatically, but agents that talk
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

Each entry is a hostname with an optional `:port` suffix. Entries containing commas or
whitespace are ignored with a warning, preventing a single entry from
smuggling several rules into the policy engine. The runner prints the resulting allowlist
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

## Sandbox templates

A sandbox template is the OCI image the sandbox VM boots from. It determines
the base OS and the tools available inside the VM, including whether the
`docker-agent` binary is already there. `--template` (or `-t` on the sandbox backend's
`create` command) selects it.

### The default template

`--template` defaults to `docker/docker-agent-sbx-templates:latest`, the
template this repository's own CI builds and publishes from the most recent
`v*` release — see [The repo-published templates](#the-repo-published-templates)
for what it contains and the other available tags.

### The repo-published templates

This repository's own CI builds and publishes the sandbox template — from
the `template` stage of the
[`Dockerfile`](https://github.com/docker/docker-agent/blob/main/Dockerfile) —
as `docker/docker-agent-sbx-templates`:

| Tag       | Built from                    | Use it when                                                        |
| --------- | ------------------------------ | ------------------------------------------------------------------- |
| `:latest` | The most recent `v*` release   | You want the newest Docker Agent template with release stability (the default). |
| `:edge`   | The current `main` branch      | You want today's `main` build and can tolerate lower stability.   |

(Each release also publishes a matching version-pinned tag, e.g.
`docker/docker-agent-sbx-templates:1.2.3`.)

`:latest` is the default. Reach for `:edge` only when you specifically
need an unreleased fix or feature and can tolerate the occasional breakage.
For reproducible runs, pin to a digest instead of a tag, e.g.
`docker/docker-agent-sbx-templates@sha256:...`.

Select one with `--template`:

```bash
$ docker agent run --sandbox --template docker/docker-agent-sbx-templates:edge agent.yaml
```

Or point the [`sbx`](https://docs.docker.com/ai/sandboxes/) CLI at it
directly, without going through Docker Agent:

```bash
$ sbx create -t docker/docker-agent-sbx-templates:latest
```

> [!TIP]
> The [`sbx` documentation](https://docs.docker.com/ai/sandboxes/) covers the
> sandbox CLI and runtime independently of Docker Agent.

### What they contain

The `template` stage in this repository's
[`Dockerfile`](https://github.com/docker/docker-agent/blob/main/Dockerfile)
layers onto `docker/sandbox-templates:shell-docker` and adds:

- The `docker-agent` binary.
- `vim` and `tmux`, for interactive debugging inside the VM.
- `kubectl`, for interacting with Kubernetes clusters from inside the VM.
- The **AWS CLI** (`aws`), for working with AWS from inside the VM.
- The **`docker-mcp`** Docker CLI plugin, installed at `~/.docker/cli-plugins/docker-mcp`.

The image carries the label `com.docker.sandboxes.flavor=docker-agent-docker`
so sandbox tooling can identify it.

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

## How local sandbox mode works

1. `--sandbox` tells Docker Agent to invoke standalone `sbx` when available, or `docker sandbox` otherwise.
2. A new sandbox VM is created from the image passed via `--template`.
3. The current working directory is mounted into the VM; the image or workload kit supplies the agent binary.
4. The [auto-kit](#auto-kit) is staged on the host and bind-mounted read-only into the VM, so the agent sees its skills and prompt files inside the sandbox.
5. The default-deny network proxy is opened for the configured [models gateway](../../features/cli/index.md#runtime-configuration-flags) and any package hosts the auto-installer needs for the agent's MCP/LSP toolsets.
6. All tools (shell, filesystem, background jobs, etc.) run inside the VM.
7. Docker Agent leaves local sandbox lifecycle management to sbx. Named sandboxes are reused only for the same launch configuration (workspace, exact mount modes, template, ordered kits and arguments, and local kit contents). A changed configuration gets a new name; previous sandboxes are never deleted automatically. Remove unused ones explicitly with `sbx rm --force NAME`.

### What the default template includes

The default template (`docker/docker-agent-sbx-templates:latest`) is built
from this repository's own `Dockerfile` — see
[What they contain](#what-they-contain) for the exact contents, including
the `docker-mcp` CLI plugin.

## Auto-Kit

The sandbox VM has its own filesystem and `$HOME` — none of the host's `~/.agents/skills/`, `~/.claude/skills/`, project-level `.agents/skills/`, or prompt files like `AGENTS.md` and `CLAUDE.md` are visible inside it. To bridge that gap, Docker Agent automatically builds a **kit**: a self-contained directory staged on the host before the sandbox starts and bind-mounted read-only into the VM at the same path.

This host auto-kit is distinct from an sbx workload/mixin kit. It is built for local default-template launches with an agent reference, but not for `--cloud` or `--sandbox-kit`. It is opt-out via `--no-kit`.

### What gets staged

For the agent referenced on the command line, the kit collects:

- **Local [skills](../../features/skills/index.md)** — every `SKILL.md` discovered on the host (global `~/.codex/skills/`, `~/.claude/skills/`, `~/.agents/skills/`, plus project `.claude/skills/`, `.github/skills/` and `.agents/skills/`) is copied under `<kit>/skills/<skill-name>/`. The in-sandbox skills loader reads from the kit instead of the (non-existent) host `$HOME`.
- **Prompt files** — every file referenced via the agent's `add_prompt_files` (`AGENTS.md`, `CLAUDE.md`, …) is collected. Files that already live under the working directory are left alone (the live workspace mount surfaces them); files outside it (e.g. an `AGENTS.md` in `$HOME`) are copied under `<kit>/prompt_files/`.
- **A manifest** — `<kit>/manifest.json` records what was staged. The on-disk copy is sanitised so it cannot be used to map the host filesystem from inside the sandbox.

Before launch, Docker Agent prints a summary of what was staged so you can see exactly which skills and prompt files the agent will have access to inside the sandbox.

When building a kit, `add_prompt_files` entries must be local relative paths
(e.g. `AGENTS.md` or `instructions/AGENTS.md`), not absolute paths or paths
escaping via `..`. Staged writes are confined to the kit directory, including
when a destination contains a symlink.

Skill traversal and file reads are anchored to the skill directory, so replacing
a file or parent directory with an escaping symlink cannot copy outside content
into the kit. File symlinks targeting the skill directory (relative or absolute)
are supported; outside, dangling, and directory symlinks are skipped.

### Secret redaction

Every text file copied into the kit is run through [portcullis](https://github.com/docker/portcullis), which redacts secrets that match its detection patterns (API keys, tokens, …) in the staged copy. The kit's printed summary marks files as `(redacted)` whenever at least one secret was replaced. Detection is best-effort — portcullis recognises common secret formats but novel or obfuscated tokens may slip through, so the kit is not a substitute for keeping secrets out of skill sources in the first place.

### Network allowlist

The sandbox templates ship with a default-deny network proxy that allows the major model providers but blocks `*.docker.com` and every package-registry / source host the auto-installer reaches for. When the agent declares MCP or LSP toolsets that have a `command` and an installable `version`, the kit build resolves each toolset's package against the [aqua](https://aquaproj.github.io/) registry and computes the minimal set of hosts the in-sandbox auto-installer will need (Go module proxy + toolchain bootstrap for `go_install` packages, GitHub release hosts for `github_release` packages, …). Those hosts, `models.dev` (needed so the in-sandbox agent can resolve model metadata such as context limits, pricing, and capabilities — without it the first catalog lookup fails with a `403 Blocked by network policy` error), and the configured [`--models-gateway`](../../features/cli/index.md#runtime-configuration-flags) — are then allow-listed on the sandbox proxy. If a per-toolset registry lookup fails, a conservative fallback union is used so the run can still succeed; the affected toolsets are surfaced in the printed summary.

### Caching

Kits are stored under the Docker Agent cache directory (`~/Library/Caches/cagent/sandbox-kits/<hash>` on macOS) keyed by a content hash of the agent reference. Reusing the same agent across runs reuses the same kit directory in place; disk usage is bounded by the number of distinct agents you have run. Kits are deliberately kept on disk between runs because the reused sandbox VM holds a hard reference to the kit's bind-mount path — deleting it would leave the sandbox un-startable.

### Disabling the kit

Pass `--no-kit` to skip the kit build entirely. The agent then runs without Docker Agent-staged host skills or external prompt files. Tool-install host inference is skipped; gateway, models.dev, and declared/user allowances still apply to default local launches. Useful for debugging the sandbox itself, or for agents that don't depend on host skills.

```bash
docker agent run --sandbox --no-kit agent.yaml
```

> [!WARNING]
> **Limitations**
>
> - Legacy/user-created sandboxes are neither adopted nor removed merely because they share a workspace. Changed launch configurations leave the previous sandbox intact. Pin remote kits/images by digest for reproducibility; mutable tags are not refreshed on reuse.
> - Only the working directory, the agent config directory, and (when staged) the kit directory are mounted; other host files are not visible to the agent.
> - Network egress is constrained by the sandbox backend's default-deny policy plus the per-run allowlist described above.

Explicit local `--data-dir` (or `DOCKER_AGENT_DATA_DIR`) and `--cache-dir`
overrides are forwarded only when
inside the writable workspace; paths outside it are rejected rather than silently
using a different session database or exposing another host directory.
