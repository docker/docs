---
title: "Permissions"
description: "Control which tools can execute automatically, require confirmation, or are blocked entirely."
keywords: docker agent, ai agents, configuration, yaml, permissions
weight: 70
canonical: https://docs.docker.com/ai/docker-agent/configuration/permissions/
---

_Control which tools can execute automatically, require confirmation, or are blocked entirely._

## Overview

Permissions provide fine-grained control over tool execution. You can configure which tools are auto-approved (run without asking), which require user confirmation, and which are completely blocked.

> [!NOTE]
> **Evaluation Order**
>
> Permissions are evaluated in this order: **Deny → Allow → Ask**. Deny patterns take priority, then allow patterns, and anything else falls through to the session's [safety mode](#safety-modes).

## Safety Modes

Every session runs in a **safety mode** that decides what happens when no permission rule matched a tool call. The runtime labels each call `safe` (safe-listed shell command such as `ls` or `git status`, or a read-only-annotated tool), `destructive` (destructive shell command such as `rm -rf`, or a destructive-annotated tool), or `unknown` — and the mode gates on that label:

| Mode | safe | destructive | unknown |
| ---- | ---- | ----------- | ------- |
| `strict` | ask | ask | ask |
| `balanced` | **allow** | ask | ask |
| `autonomous` | **allow** | **allow** | **allow** |

- **`strict`** prompts for every tool call, read-only ones included. Only an `allow:` rule silences a prompt.
- **`balanced`** runs safe calls silently and asks about everything else.
- **`autonomous`** is the legacy `--yolo` behavior: everything runs. Only `deny:` rules, session-scoped `ask:` rules, and `preempt_yolo` hooks still gate.

Pick a mode with the `--safety` flag (`docker-agent run --safety balanced ...`), the `safety_policy` field on session create (`POST /api/sessions`) or mid-session (`PATCH /api/sessions/:id/safety-policy`), or escalate directly from a confirmation prompt (`B` switches to balanced, `A` to autonomous). Sessions that never choose a mode keep the historical default: read-only tools auto-approve, everything else asks.

### Declarative Safety Defaults

Safety modes can also be declared as **defaults** in YAML, at four scopes:

| Scope | Location | Owner |
| ----- | -------- | ----- |
| Alias | `aliases.<name>.safety` in `~/.config/cagent/config.yaml` (or `docker agent alias add ... --safety <mode>`) | User |
| Global settings | `settings.safety` in `~/.config/cagent/config.yaml` | User |
| Per-agent | `agents.<name>.safety` in the agent YAML | Agent author |
| Config-wide | `runtime.safety` in the agent YAML | Agent author |

```yaml
# Agent YAML (author-declared defaults)
runtime:
  safety: balanced # config-wide default for new sessions

agents:
  root:
    safety: strict # overrides runtime.safety for this agent
```

All four fields accept only the three canonical modes — `strict`, `balanced`, `autonomous` (yes, an author may declare `autonomous`) — and any other value fails loading with an error naming the field. The legacy spellings remain as aliases for `autonomous`: `settings.YOLO`, the alias `yolo` option, and the `--yolo` flag. When both are set at the same scope, `safety` wins over the legacy `YOLO`/`yolo`.

For a **new** root session the first source in this order wins:

1. explicit `--safety` flag
2. explicit `--yolo` flag
3. alias `safety`/`yolo` option
4. `settings.safety`/`settings.YOLO` (user config)
5. selected agent's `agents.<name>.safety`
6. `runtime.safety`
7. the historical default (read-only tools auto-approve, everything else asks)

**Resuming a session never re-applies defaults**: the stored mode is kept unless you pass an explicit `--safety` or `--yolo` flag for that run. Agent switches, handoffs, and delegated sub-agent sessions inherit the active session's mode rather than resetting it.

Sessions created through the API (`POST /api/sessions`) without a `safety_policy` receive the author-declared defaults (5–6) when their first run starts — the earliest point the agent configuration is loaded. If the server restarts before that first run, the session keeps the historical unset default (7).

> [!WARNING]
> **Trust: author defaults never outrank you.** `runtime.safety` and `agents.<name>.safety` are written by the agent's author — which may be a config you pulled from a URL or an OCI registry. They only fill the gap when you expressed no preference: any user-owned source (CLI flag, alias option, user settings) always takes precedence, and a resumed session keeps its stored mode. Still, an author default of `autonomous` means a fresh session runs every tool call unprompted — review third-party configs before running them, or pin your own floor with `settings.safety` / `--safety`.

**Custom rules always win over the mode**, with one asymmetry: `ask:` rules written in an agent's YAML (or global config) are agent-author advisories and yield to a user-chosen `balanced`/`autonomous` mode, while `ask:` rules granted at the session level (interactive "always ask" decisions, the session permissions API) always prompt.

## Permission Levels

Permissions can be defined at two levels:

| Level | Location | Scope |
| ----- | -------- | ----- |
| **Agent-level** | Agent YAML config (`permissions:` section) | Applies to that specific agent config |
| **Global (user-level)** | `~/.config/cagent/config.yaml` under `settings.permissions` | Applies to every agent you run |

Hooks follow the same user-config pattern: agent-level hooks live under `agents.<name>.hooks`, and global hooks live under `settings.hooks`. See [Hooks](../hooks/index.md#global-user-level-hooks).

Both levels use the same `allow`/`ask`/`deny` pattern syntax. When both are present, they are **merged** at startup -- patterns from both sources are combined into a single checker. See [Merging Behavior](#merging-behavior) for details.

## Agent-Level Configuration

```yaml
agents:
  root:
    model: openai/gpt-4o
    description: Agent with permission controls
    instruction: You are a helpful assistant.

permissions:
  # Auto-approve these tools (no confirmation needed)
  allow:
    - "read_file"
    - "read_*" # Glob patterns
    - "shell:cmd=ls*" # With argument matching

  # Always ask before running these tools, even if an allow pattern would match
  ask:
    - "shell:cmd=git push*"
    - "write_file:path=/home/user/important/*"

  # Block these tools entirely
  deny:
    - "shell:cmd=sudo*"
    - "shell:cmd=rm*-rf*"
    - "dangerous_tool"
```

The three lists are evaluated in order `deny` → `allow` → `ask`, so an `ask:` entry lets you add a confirmation layer on top of an otherwise-allowed tool.

## Global Permissions

Global permissions let you enforce rules across **all** agents, regardless of which agent config you run. Define them in your user config file:

```yaml
# ~/.config/cagent/config.yaml
settings:
  permissions:
    deny:
      - "shell:cmd=sudo*"
      - "shell:cmd=rm*-rf*"
    allow:
      - "read_*"
      - "shell:cmd=ls*"
      - "shell:cmd=cat*"
```

This is useful for setting personal safety guardrails that apply everywhere -- for example, always blocking `sudo` or always auto-approving read-only tools -- without relying on each agent config to include those rules.

### Merging Behavior

When both global and agent-level permissions are present, they are merged into a single set of patterns before evaluation. The merge works as follows:

- **Deny patterns from either source block the tool.** A global deny cannot be overridden by an agent-level allow, and vice versa.
- **Allow patterns from either source auto-approve the tool** (as long as no deny pattern matches).
- **Ask patterns from either source force confirmation** (as long as no deny or allow pattern matches).

The evaluation order remains the same after merging: **Deny > Allow > Ask > default Ask**.

> [!TIP]
> **Example: Global deny + agent allow**
>
> If your global config denies `shell:cmd=sudo*` and an agent config allows `shell:cmd=sudo apt update`, the deny wins. Deny patterns always take priority regardless of source.

## Pattern Syntax

Permissions support glob-style patterns with optional argument matching:

### Simple Patterns

| Pattern        | Matches                        |
| -------------- | ------------------------------ |
| `shell`        | Exact match for `shell` tool   |
| `read_*`       | Any tool starting with `read_` |
| `github_*`     | Any GitHub MCP tool            |
| `*`            | All tools                      |

### Argument Matching

You can match tools based on their argument values using `tool:arg=pattern` syntax:

```yaml
permissions:
  allow:
    # Allow shell only when cmd starts with "ls" or "cat"
    - "shell:cmd=ls*"
    - "shell:cmd=cat*"

    # Allow edit_file only in specific directory
    - "edit_file:path=/home/user/safe/*"

  deny:
    # Block shell with sudo
    - "shell:cmd=sudo*"

    # Block writes to system directories
    - "write_file:path=/etc/*"
    - "write_file:path=/usr/*"
```

> [!NOTE]
> **Colons inside argument values are preserved.** Only the `:key=` token boundaries between
> argument conditions split a pattern — colons that appear inside a value are treated as
> ordinary characters and do not start a new condition. Check a tool’s actual argument names
> (and whether they accept a string or a list) before writing an argument-matching pattern.

### Multiple Argument Conditions

Chain multiple argument conditions with colons. All conditions must match:

```yaml
permissions:
  allow:
    # Allow shell with ls in current directory
    - "shell:cmd=ls*:cwd=."

  deny:
    # Block shell with rm -rf anywhere
    - "shell:cmd=rm*:cmd=*-rf*"
```

## Glob Pattern Rules

Patterns follow filepath.Match semantics with some extensions:

- `*` — matches any sequence of characters (including spaces)
- `?` — matches any single character
- `[abc]` — matches any character in the set
- `[a-z]` — matches any character in the range

Matching is **case-insensitive**.

> [!TIP]
> **Trailing Wildcards**
>
> Trailing wildcards like `sudo*` match any characters including spaces, so `sudo*` matches `sudo rm -rf /`.

## Decision Types

| Decision  | Behavior                                            |
| --------- | --------------------------------------------------- |
| **Allow** | Tool executes immediately without user confirmation |
| **Ask**   | User must confirm before tool executes (default)    |
| **Deny**  | Tool is blocked and returns an error to the agent   |

## Examples

### Read-Only Agent

Allow all read operations, block all writes:

```yaml
permissions:
  allow:
    - "read_file"
    - "read_multiple_files"
    - "list_directory"
    - "directory_tree"
    - "search_files_content"
  deny:
    - "write_file"
    - "edit_file"
    - "shell"
```

### Safe Shell Agent

Allow specific safe commands, block dangerous ones:

```yaml
permissions:
  allow:
    - "shell:cmd=ls*"
    - "shell:cmd=cat*"
    - "shell:cmd=grep*"
    - "shell:cmd=find*"
    - "shell:cmd=head*"
    - "shell:cmd=tail*"
    - "shell:cmd=wc*"
  deny:
    - "shell:cmd=sudo*"
    - "shell:cmd=rm*"
    - "shell:cmd=mv*"
    - "shell:cmd=chmod*"
    - "shell:cmd=chown*"
```

### MCP Tool Permissions

Control MCP tools by their qualified names:

```yaml
permissions:
  allow:
    # Allow all GitHub read operations
    - "github_get_*"
    - "github_list_*"
    - "github_search_*"
  deny:
    # Block destructive GitHub operations
    - "github_delete_*"
    - "github_close_*"
```

## Combining with Hooks

Permissions work alongside [hooks](../hooks/index.md). The evaluation order is:

1. Run **`preempt_yolo` pre_tool_use hooks** — security-critical checks that no mode or allow rule can bypass
2. Check **deny** patterns — if matched, tool is blocked
3. Check **allow** patterns — if matched, tool is auto-approved
4. Apply the **[safety mode](#safety-modes)** to the call's safety label — may auto-approve
5. Run **pre_tool_use hooks** — hooks can allow, deny, or ask
6. If no decision, **ask user** for confirmation

Default-lane hooks only see calls the rules and the mode routed to "ask"; they cannot override deny decisions.

> [!WARNING]
> **Security Note**
>
> Permissions are enforced client-side. They help prevent accidental operations but should not be relied upon as a security boundary for untrusted agents. For stronger isolation, use [sandbox mode](../sandbox/index.md).
