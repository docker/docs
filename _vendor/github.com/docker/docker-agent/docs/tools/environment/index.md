---
title: "Environment Tool"
description: "Report the OS and resolved shell so the model can pick the right syntax."
keywords: docker agent, ai agents, tools, toolsets, environment tool
linkTitle: "Environment"
weight: 145
canonical: https://docs.docker.com/ai/docker-agent/tools/environment/
---

_Report the OS and resolved shell so the model can pick the right syntax._

## Overview

The environment tool exposes a single read-only call — `get_environment_info` — that returns the operating system and the resolved shell that will run shell-tool commands, e.g. `{"os":"Windows","shell":"powershell"}`.

It takes no arguments, has no side effects, and reads only `runtime.GOOS` and the shell binary that the shell tool has already resolved. Its `ReadOnlyHint` annotation classifies it as safe. Balanced, restricted, and autonomous modes auto-approve it unless a permission rule or preempting hook intervenes. The legacy default auto-approves it only after ordinary pre-tool hooks run. Strict mode does not auto-approve it merely because it is read-only; see [Permissions](../../configuration/permissions/index.md).

## When to use

Pair with the shell toolset when the model may run on hosts with unfamiliar shells (Windows PowerShell, cmd.exe, fish, nushell). The shell tool description already names the resolved interpreter, but that hint sits in the tool schema; giving the model a callable tool provides a first-class fallback for edge cases such as multi-agent handoffs or long sessions where the schema hint has slid out of attention.

## Configuration

```yaml
toolsets:
  - type: environment
  - type: shell
```

No configuration options.

## Output shape

```json
{
  "os": "Windows",
  "shell": "powershell"
}
```

- `os` — friendly label for `runtime.GOOS` (`Windows`, `macOS`, `Linux`, or the raw `GOOS` for anything exotic).
- `shell` — lowercase basename of the resolved shell (`powershell`, `pwsh`, `cmd`, `bash`, `zsh`, `fish`, …).

The output shape is fixed. No working directory, no full paths, no username — anything user-controlled would land in the conversation transcript.
