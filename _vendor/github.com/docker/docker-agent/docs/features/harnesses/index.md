---
title: "Coding Harnesses"
description: "Delegate coding tasks to external AI coding CLIs (Claude Code, Codex, opencode) as sub-agents."
keywords: docker agent, ai agents, features, coding harnesses
weight: 50
canonical: https://docs.docker.com/ai/docker-agent/features/harnesses/
---

_Delegate coding tasks to external AI coding CLIs (Claude Code, Codex, opencode) as sub-agents._

## Overview

A **harness** agent delegates its work to an external coding CLI — `claude` (Claude Code), `codex` (OpenAI Codex), `opencode`, or `pi` — instead of calling a model API directly. The external CLI drives the coding loop while Docker Agent provides orchestration, hooks, permissions, and distribution.

This pattern gives you the best of both worlds:

- **External CLI strengths** — deep IDE integration, specialized coding workflows, CLI-native tool access
- **Docker Agent strengths** — multi-agent coordination, hook-based auditing and policy enforcement, permission controls, OCI distribution, and the full agent config schema

> [!NOTE]
> **When to use harnesses**
>
> Use a harness when you want a Claude Code / Codex / opencode session to act as a sub-agent inside a larger Docker Agent workflow — for example, an orchestrator that plans work and delegates coding tasks to specialized harness agents.

## Prerequisites

The external CLI must be installed and available on `PATH` before starting Docker Agent:

| Harness type | Required binary | Install |
| --- | --- | --- |
| `claude-code` | `claude` | [docs.anthropic.com/en/docs/claude-code](https://docs.anthropic.com/en/docs/claude-code) |
| `codex` | `codex` | [github.com/openai/codex](https://github.com/openai/codex) |
| `opencode` | `opencode` | [opencode.ai](https://opencode.ai) |
| `pi` | `pi` | Refer to the `pi` CLI documentation for installation instructions. |

Docker Agent will report an error at session start if the required binary is not found. The `claude-code` harness additionally needs the CLI to be logged in — see [Authentication](#authentication-claude-subscription-no-api-key) below.

## Configuration

Add a `harness:` block to any agent definition to make it harness-backed:

```yaml
agents:
  coder:
    description: A Claude Code harness agent
    harness:
      type: claude-code
```

Harness agents do **not** need a `model:` field — the external CLI manages its own model selection.

### Field Reference

| Field | Applies to | Type | Description |
| --- | --- | --- | --- |
| `type` | all | string | **Required.** One of `claude-code`, `codex`, `opencode`, `pi` |
| `model` | all | string | Optional model override forwarded to the CLI. When omitted, the CLI uses its own default model. |
| `effort` | `claude-code` | string | Reasoning effort: `low` \| `medium` \| `high` \| `xhigh` \| `max` — forwarded as `--effort`. When omitted, Claude Code uses its own default. |
| `agent` | `opencode` | string | opencode agent profile name |
| `thinking` | `opencode` | boolean | Enable extended thinking — forwarded as `--thinking` |

### Claude Code

```yaml
agents:
  coder:
    description: Claude Code coding agent
    harness:
      type: claude-code
      model: claude-sonnet-4-5 # optional: alias (sonnet, opus, haiku) or full model ID
      effort: high # low | medium | high | xhigh | max
```

- `model` accepts whatever the `claude` CLI accepts for `--model`: an alias
  like `sonnet`, `opus`, or `haiku`, or a full model ID like
  `claude-sonnet-4-5`. Omit it to use Claude Code's own default model. This is
  **not** a Docker Agent `provider/model` reference — the harness model never
  goes through Docker Agent's model providers or routing.
- `effort` is forwarded as `--effort` and must be one of `low`, `medium`,
  `high`, `xhigh`, or `max`. Omit it to use Claude Code's own default.

#### Authentication (Claude subscription, no API key)

The `claude-code` harness runs the official CLI with the CLI's **own login**.
A Claude (claude.ai) subscription sign-in is all it needs — no
`ANTHROPIC_API_KEY` and no Docker Agent model credential are required for the
harness agent itself. Log in once:

```bash
$ claude auth login --claudeai   # interactive, opens a browser
$ claude auth status --text      # verify the login
```

The login is stored per OS user by the CLI (in its own configuration and, on
macOS, the keychain) and is found via `HOME` and the process environment. Run
the login **as the same OS user and environment that run `docker agent`** — a
login made under another user, container, or `sudo` context is invisible to
the harness. Docker Agent never reads, copies, or stores the CLI's tokens; it
only launches `claude`, which authenticates itself.

If the CLI is not logged in when a harness agent runs, the `claude` subprocess
fails at session start and the agent reports a harness error — Docker Agent
never opens a browser or starts a login on its own. Diagnose and fix with:

```bash
$ docker agent doctor ./agent.yaml   # checks install + login for claude-code harness files
$ docker agent setup                 # pick "Claude Code harness" to be walked through it
```

`docker agent doctor <file>` probes the CLI only when the file declares a
`claude-code` harness, and reports installation, version, and safe login
metadata (auth method, API provider, subscription type — never your email,
organization, or tokens). `docker agent setup` offers to run the official
`claude auth login --claudeai` for you (only after you confirm) and writes a
ready-to-run `claude-code-agent.yaml`.

> [!WARNING]
> **The harness bypasses Claude Code's permission prompts**
>
> Docker Agent runs the CLI non-interactively with its own tools and passes
> `--dangerously-skip-permissions`: Claude Code edits files and runs commands
> without asking. Only point a harness agent at a repository you trust, and
> prefer isolation — `docker agent run --worktree` runs it on an isolated git
> worktree, keeping its changes off your checkout (a worktree with work, or
> from a non-interactive run, is kept for inspection per the normal cleanup
> rules). `docker agent run --sandbox` does not automatically carry the
> `claude` CLI or its login into the sandbox, so it cannot isolate the
> harness unless the sandbox image is separately provisioned and
> authenticated.

### Codex

```yaml
agents:
  coder:
    description: Codex coding agent
    harness:
      type: codex
      model: o4-mini     # optional model override
```

### opencode

```yaml
agents:
  coder:
    description: opencode coding agent
    harness:
      type: opencode
      agent: my-profile  # optional agent profile
      thinking: true     # enable extended thinking
```

## What Does NOT Work

Harness agents bypass the Docker Agent model pipeline entirely. As a result:

- **Docker Agent toolsets are inactive.** The external CLI provides its own tools — filesystem, shell, etc. Any `toolsets:` defined on a harness agent are silently ignored.
- **`model:` routing is unavailable.** The harness CLI manages model selection; Docker Agent's `models:` configuration and routing rules do not apply to harness agents.
- **Token usage tracking depends on the external CLI.** Docker Agent records usage when the CLI reports it (Claude Code and Codex both report usage data). If the CLI does not emit usage data, the session will show zero token usage.

> [!WARNING]
> **No Docker Agent toolsets inside a harness**
>
> Do not configure `toolsets:` on a harness agent — they are silently ignored. If you need Docker Agent toolsets alongside external coding capabilities, use a standard sub-agent with `transfer_task` rather than a harness.

## Hook Behavior

Hooks work normally on harness agents, including `before_llm_call` and `after_llm_call`. `before_llm_call` runs before the prompt is forwarded to the external CLI and can block or rewrite the run; `after_llm_call` fires after the CLI returns its final response.

The `model_id` field in hook payloads is set to the harness label (e.g. `claude-code`) rather than a canonical `provider/model` string. This applies to `before_llm_call`, `after_llm_call`, and any other event that carries `model_id`.

See [Hooks](../../configuration/hooks/index.md) for the full hook reference.

## Recipe: Root Harness Agent

The simplest setup: a single root agent that hands everything to Claude Code.
No `models:` section, no API key — the CLI's subscription login does the work.
This is exactly the file `docker agent setup` generates for the Claude Code
harness path (as `claude-code-agent.yaml`).

```yaml
# claude-code-agent.yaml
agents:
  root:
    description: Claude Code running on your Claude subscription
    harness:
      type: claude-code
      effort: medium # low | medium | high | xhigh | max; omit for the Claude Code default
      # model: claude-sonnet-4-5   # optional; omit for the Claude Code default
```

```bash
$ docker agent run claude-code-agent.yaml
$ docker agent doctor claude-code-agent.yaml   # verify the CLI is installed and logged in
```

## Recipe: Orchestrator + Harness Sub-Agents (Sequential)

An orchestrator plans the work and delegates to specialized harness agents one at a time. Each coding agent runs in its own sub-session and reports results back.

```yaml
# examples/coding_harnesses.yaml

models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5

agents:
  root:
    model: claude
    description: Orchestrator that plans and delegates coding tasks
    instruction: |
      You are a project orchestrator. Break down coding requests into
      focused tasks and delegate each task to the most appropriate
      coding agent. Collect results and synthesize a final summary.
    sub_agents:
      - claude-coder
      - codex-coder

  claude-coder:
    description: Claude Code specialist for complex refactors
    harness:
      type: claude-code
      model: claude-sonnet-4-5
      effort: xhigh

  codex-coder:
    description: Codex specialist for code generation
    harness:
      type: codex
```

Only the orchestrator's `model: claude` needs an Anthropic API key — the
`claude-coder` harness agent authenticates through the CLI's own subscription
login.

The root agent calls `transfer_task` to send work to a harness sub-agent, waits for the result, and continues. See the [full example on GitHub](https://github.com/docker/docker-agent/blob/main/examples/coding_harnesses.yaml).

## Recipe: Parallel Harness Dispatch

Combine the `background_agents` toolset with harness sub-agents to dispatch multiple coding tasks simultaneously:

```yaml
# examples/coding_harness_background_agents.yaml

models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-5

agents:
  root:
    model: claude
    description: Orchestrator that fans out coding tasks in parallel
    instruction: |
      Use background agents to run multiple coding tasks at once.
      Dispatch all tasks, then collect results when each finishes.
    sub_agents:
      - claude-coder
      - codex-coder
    toolsets:
      - type: background_agents

  claude-coder:
    description: Frontend specialist (Claude Code)
    harness:
      type: claude-code
      effort: medium

  codex-coder:
    description: Backend specialist (Codex)
    harness:
      type: codex
```

The orchestrator calls `run_background_agent` for each task, monitors progress with `list_background_agents`, and collects results with `view_background_agent`. See the [full example on GitHub](https://github.com/docker/docker-agent/blob/main/examples/coding_harness_background_agents.yaml).

For the general background agents reference, see [Background Agents](../../tools/background-agents/index.md).

## See Also

- [Multi-Agent Systems](../../concepts/multi-agent/index.md) — orchestration patterns
- [Background Agents](../../tools/background-agents/index.md) — parallel task dispatch
- [Hooks](../../configuration/hooks/index.md) — auditing and policy enforcement
- [Agent Configuration](../../configuration/agents/index.md) — full agent schema reference
