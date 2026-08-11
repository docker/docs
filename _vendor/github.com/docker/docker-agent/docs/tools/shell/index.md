---
title: "Shell Tool"
description: "Execute arbitrary shell commands in the user's environment."
keywords: docker agent, ai agents, tools, toolsets, shell tool
linkTitle: "Shell"
weight: 20
canonical: https://docs.docker.com/ai/docker-agent/tools/shell/
---

_Execute arbitrary shell commands in the user's environment._

## Overview

The shell tool allows agents to execute arbitrary shell commands synchronously. This is one of the most powerful tools — it lets agents run builds, install dependencies, query APIs, and interact with the system. Each call runs in a fresh, isolated shell session — no state persists between calls.

Commands have a default 30-second timeout and require user confirmation unless `--yolo` is used. For servers, watchers, and other long-running commands, add the [`background_jobs`](../background-jobs/index.md) toolset alongside `shell`.

## Configuration

```yaml
toolsets:
  - type: shell
```

### Options

| Property       | Type    | Description                                                                                          |
| -------------- | ------- | --------------------------------------------------------------------------------------------------- |
| `env`          | object  | Environment variables to set for all shell commands                                                 |
| `safer`        | boolean | Detect destructive shell commands and force confirmation regardless of `--yolo` or permission rules (see [Safer mode](#safer-mode)). Default `false`. |
| `sudo_askpass` | boolean | Opt in to prompting for a `sudo` password (see [Sudo support](#sudo-support)). Default `false`.     |

### Custom Environment Variables

```yaml
toolsets:
  - type: shell
    env:
      MY_VAR: "value"
      PATH: "${env.PATH}:/custom/bin"
```

### Safer mode

Set `safer: true` to enable destructive-command detection for the shell toolset:

```yaml
toolsets:
  - type: shell
    safer: true
```

This auto-registers the [`safer_shell`](../../configuration/hooks/index.md#built-in-hooks) builtin
under `pre_tool_use` with `preempt_yolo: true` so the entry fires
before `Decide()` / `--yolo`. Three behaviors:

- **Destructive matches** (`rm -rf <path>`, `docker volume rm`, `mkfs`, `dd if=… of=/dev/<disk>`, …) get a forced user confirmation carrying a `blast_radius` classification (`low` / `medium` / `high` / `unknown`) and a `category` tag. The TUI confirmation dialog renders the blast radius with a color badge.
- **Known-safe reads** (`ls`, `cat`, `git status`, `git diff`, `docker ps`, `docker logs`, `kubectl get`, …) flow through silently — they're treated as no-opinion and follow the regular approval pipeline (`--yolo`, permission rules, read-only hint).
- **Everything else** asks with `blast_radius=unknown`. Safer mode is conservative by default: unrecognised commands surface to the user before `--yolo` or permission allow-rules can auto-approve them.

The verdict cannot be bypassed by `--yolo` or by a `permission_request` hook that returns `allow` — the `preempt_yolo` lane runs before both. Compound shell (`a && b`, `a; b`, `a | b`) is never matched against the safe allowlist; any destructive segment falls through to ask. The full taxonomy lives in [`pkg/hooks/builtins/safety_patterns.json`](https://github.com/docker/docker-agent/blob/main/pkg/hooks/builtins/safety_patterns.json).

See [`examples/shell_safer.yaml`](https://github.com/docker/docker-agent/blob/main/examples/shell_safer.yaml) for a full example. Under the hood, `safer: true` is a sugar that appends one entry under `hooks.pre_tool_use` with `preempt_yolo: true`; writing the entry by hand achieves the same thing.

### Sudo support

By default a shell command has no controlling terminal, so a `sudo` command that needs a password hangs until it times out (the agent usually gives up and falls back to printing manual instructions).

Set `sudo_askpass: true` to enable a sudo privilege escalation flow:

```yaml
toolsets:
  - type: shell
    sudo_askpass: true
```

When enabled, `sudo` commands prompt you for your password through the host UI (the input is masked). The password is handed to `sudo` over a private, per-session socket via the standard `SUDO_ASKPASS` mechanism — it is never written to the command line, the logs, or stored by the agent.

The bridge environment variables (`SUDO_ASKPASS`, `CAGENT_ASKPASS_SOCKET`, `CAGENT_ASKPASS_TOKEN`) are added only to commands that invoke `sudo`, but within such a command they are visible to every child process, not just `sudo`. They carry a socket path and a session token, not the password; the socket lives in a `0700` directory, so only your own user can reach it.

Notes and limitations:

- Unix only. The flag has no effect on Windows.
- Interactive UI only. In headless / non-interactive runs the prompt is declined automatically and `sudo` fails as before.
- Only a bare `sudo ...` invocation in a POSIX shell (`sh`, `bash`, `zsh`, ...) is handled. `sudo` called by absolute path (`/usr/bin/sudo`), via `env sudo`, from inside a nested script, or under a non-POSIX shell (e.g. `fish`) is not intercepted and behaves as before.
- Caching is `sudo`'s own. Because each shell tool call runs in a fresh shell with no controlling terminal, `sudo`'s credential cache does not persist across separate tool calls: you are prompted once per shell command that uses `sudo`. Within a single command, multiple `sudo` calls (e.g. `sudo a && sudo b`) usually share one prompt, subject to `sudo`'s own timestamp configuration.
- The prompt must be answered within the command's timeout; raise the `timeout` parameter for `sudo` commands that may wait on input.
- Prompts are serialized: if a single command runs two `sudo` calls in parallel (e.g. `sudo a & sudo b`), the second waits for the first prompt to be answered rather than opening two dialogs at once.

## Available Tools

The shell toolset exposes one tool:

| Tool Name | Description                                                                  |
| --------- | ---------------------------------------------------------------------------- |
| `shell`   | Run a command synchronously and return its combined output when it finishes. |

### `shell` parameters

| Parameter | Type    | Required | Description                                                               |
| --------- | ------- | -------- | ------------------------------------------------------------------------- |
| `cmd`     | string  | ✓        | The shell command to execute.                                             |
| `cwd`     | string  | ✗        | Working directory to run the command in (default: `.`).                   |
| `timeout` | integer | ✗        | Per-call execution timeout in seconds (default: `30`).                    |

> [!WARNING]
> **Safety**
>
> The shell tool gives agents full access to the system shell. Always set `max_iterations` on agents that use the shell tool to prevent infinite loops. A value of 20–50 is typical for development agents. Use [Sandbox Mode](../../configuration/sandbox/index.md) for additional isolation.

> [!NOTE]
> **Tool Confirmation**
>
> By default, docker-agent asks for user confirmation before executing shell commands. Use `--yolo` to auto-approve all tool calls.
