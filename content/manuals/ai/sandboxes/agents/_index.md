---
title: Supported agents
linkTitle: Agents
description: AI coding agents supported by Docker Sandboxes with experimental status and configuration details.
weight: 50
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

Docker Sandboxes supports multiple AI coding agents. All agents run isolated
inside microVMs with private Docker daemons.

## Supported agents

| Agent       | Command    | Status       | Notes                                     |
| ----------- | ---------- | ------------ | ----------------------------------------- |
| Claude Code | `claude`   | Experimental | Most tested implementation                |
| Codex       | `codex`    | Experimental | In development                            |
| Copilot     | `copilot`  | Experimental | In development                            |
| Gemini      | `gemini`   | Experimental | In development                            |
| cagent      | `cagent`   | Experimental | In development                            |
| Kiro        | `kiro`     | Experimental | In development                            |
| OpenCode    | `opencode` | Experimental | In development                            |
| Custom shell | `shell`   | Experimental | Minimal environment for manual setup      |

## Experimental status

All agents are experimental features. This means:

- Breaking changes may occur between Docker Desktop versions
- Features may be incomplete or change significantly
- Stability and performance are not production-ready
- Limited support and documentation

Use sandboxes for development and testing, not production workloads.

## Using different agents

The agent type is specified when creating a sandbox:

```console
$ docker sandbox create AGENT [PATH] [PATH...]
```

Each agent runs in its own isolated sandbox. The agent type is bound to the
sandbox when created and cannot be changed later.

## Template environment

All agent templates share a common base environment:

- Ubuntu 25.10 base
- Development tools: Docker CLI (with Buildx and Compose), Git, GitHub CLI, Node.js, Go, Python 3, uv, make, jq, ripgrep
- Non-root `agent` user with sudo access
- Private Docker daemon for running additional containers
- Package managers: apt, pip, npm

Individual agents add their specific CLI tools on top of this base. See
[Custom templates](../templates.md) to build your own agent images.

## Agent-specific configuration

Each agent has its own credential requirements and authentication flow.
Credentials are scoped per agent and must be provided specifically for that
agent (no fallback authentication methods are used).

See the agent-specific documentation:

- [Claude Code](./claude-code.md)
- [cagent](./cagent.md)
- [Codex](./codex.md)
- [Copilot](./copilot.md)
- [Gemini](./gemini.md)
- [Kiro](./kiro.md)
- [OpenCode](./opencode.md)
- [Custom shell](./shell.md)

## Requirements

- Docker Desktop 4.58 or later
- Platform support:
  - macOS with virtualization.framework
  - Windows with Hyper-V {{< badge color=violet text=Experimental >}}
- API keys or credentials for your chosen agent
