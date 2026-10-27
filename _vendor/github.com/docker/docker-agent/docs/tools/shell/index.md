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

### Shell interpreter detection

The shell tool automatically detects and names the resolved shell interpreter (e.g., `bash`, `zsh`, `powershell`, `pwsh`, `cmd`) in its description to the model, along with the operating system (Linux, macOS, Windows). This helps models use the correct shell syntax for the host environment.

For example:

- On Linux with bash: "Executes the given shell command with bash on Linux."
- On Windows with PowerShell 5.1: "Executes the given shell command with powershell on Windows. Windows PowerShell dialect. Chain commands with ";" (not "&&"; "&&" is a parse error on 5.1). POSIX utilities are not available — use Select-String (not grep), Select-Object -First N (not head), -Last N (not tail), Measure-Object -Line (not wc -l), Get-Content (not cat), Expand-Archive (not gunzip/tar). "ls -la" fails — use "ls" or "Get-ChildItem"."
- On Windows with cmd.exe: "Executes the given shell command with cmd on Windows. cmd.exe syntax, not POSIX. No POSIX utilities (grep, head, tail, wc, cat), no POSIX flags on built-ins, and variable expansion uses %VAR% (not $VAR)."

Each description also points the model at `get_environment_info` for edge cases the static hint does not cover. This reduces wasted turns where models assume POSIX syntax on Windows or vice versa.

### ACP client execution

In ACP mode, clients advertising terminal support execute the configured `shell` tool through client terminals rather than a host subprocess. The adapter uses `/bin/sh` or the absolute Windows system `cmd.exe`, requires matching client/agent platform and paths, forwards only configured environment overrides, retains a bounded output tail, and does not support `sudo_askpass`. Without terminal capability, native behavior is unchanged; client failures never trigger native fallback. See [ACP client terminals](../../features/acp/index.md#client-terminals) for cleanup, cancellation, and unknown-create limitations.

### Dialect-error correction

If a command still fails with a known shell-dialect error (for example PowerShell's `'&&' is not a valid statement separator`, a `cmd.exe`/`pwsh` "not recognized" error, or a POSIX-style `2>/dev/null` redirection on Windows), the tool prepends a `[shell-hint]` line to the output pointing the model at the corrective syntax before it retries. Hints are gated on the resolved shell so they only fire for genuine dialect mismatches (e.g. a PowerShell `&&` hint never fires for `pwsh`, which supports `&&` natively).

### Command preview in the TUI

In the TUI, multi-line shell commands are shown with syntax highlighting. When a command wraps a script for a common interpreter — `python`/`python3`, `node`/`nodejs`, `ruby`, `perl`, or `sh`/`bash`/`zsh` — through a heredoc (for example `python3 - <<'PY'`) or an inline flag (`-c` for Python and shell, `-e`/`--eval` for Node, `-e` for Ruby and Perl), the embedded script is highlighted in its own language and the surrounding command as bash. Highlighting is display-only and does not change the command that runs.

## Configuration

```yaml
toolsets:
  - type: shell
```

### Options

| Property       | Type    | Description                                                                                          |
| -------------- | ------- | --------------------------------------------------------------------------------------------------- |
| `env`          | object  | Environment variables to set for all shell commands                                                 |
| `sudo_askpass` | boolean | Opt in to prompting for a `sudo` password (see [Sudo support](#sudo-support)). Default `false`.     |

### Custom Environment Variables

```yaml
toolsets:
  - type: shell
    env:
      MY_VAR: "value"
      PATH: "${env.PATH}:/custom/bin"
```

### Command classification

Every shell command is classified against an embedded taxonomy before the approval decision — no opt-in required. The same classification applies to the `cmd` of [`run_background_job`](../background-jobs/index.md#command-classification):

- **Destructive matches** (`rm -rf <path>`, `docker volume rm`, `mkfs`, `dd if=… of=/dev/<disk>`, …) are labelled `destructive` with a `blast_radius` (`low` / `medium` / `high`) and a `category` tag. The TUI confirmation dialog renders the blast radius with a color badge.
- **Known-safe reads** (`ls`, `cat`, `git status`, `git diff`, `gh pr view`, `go env`, `docker ps`, `docker logs`, `kubectl get`, …) are labelled `safe`.
- **Everything else** is labelled `unknown`.

The session's [safety mode](../../configuration/permissions/index.md#safety-modes) decides what each label means: `strict` asks about everything, `balanced` auto-runs safe commands and asks about destructive/unknown ones, `restricted` auto-runs safe commands and denies destructive/unknown ones without asking (fail-closed for unattended runs), `autonomous` runs everything. Custom permission rules always win over the mode.

Read patterns also cover numeric `head`/`tail` counts and numeric sed print addresses such as `sed -n '10,20p' file`. Arbitrary sed programs are not safe-listed. GitHub read commands exclude browser/watch flags; Go environment queries exclude configuration writes (`-w`/`-u`). Flags that write output or run external diff/text-conversion helpers are excluded from Git log-family reads. Expansions that could hide a denied flag remain unknown.

In-place formatters (`gofmt -w`, `goimports -w`), `tee` writes, process termination (`pkill`, `killall`), and additional Git force/delete variants receive destructive labels. Ordinary builds, tests, scripts, and task runners remain unknown unless they contain a recognized destructive operation.

Compound shell (`a && b`, `a; b`, `a | b`) and command substitutions (including zsh `=(...)` and fish `(...)`) are never matched against the safe allowlist; any destructive segment falls through to ask. The full taxonomy lives in [`pkg/safety/safety_patterns.json`](https://github.com/docker/docker-agent/blob/main/pkg/safety/safety_patterns.json).

See [`examples/safety_modes.yaml`](https://github.com/docker/docker-agent/blob/main/examples/safety_modes.yaml) for a full example. The legacy `safer: true` toolset flag was removed in config version 15: a config declaring version 15 or later (including a version-less config, which resolves to the latest schema) now fails to load with `unknown field "safer"` if the flag is present — delete it, it has had no effect since v1.117.0. Configs pinned to `version: "14"` or lower still accept the flag and silently ignore it.

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
> By default, Docker Agent asks for user confirmation before executing shell commands. Use `--yolo` to auto-approve all tool calls.
