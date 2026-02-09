---
title: Supported agents
description: AI coding agents supported by Docker Sandboxes with experimental status and configuration details.
weight: 50
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

Docker Sandboxes supports multiple AI coding agents. All agents run isolated
inside microVMs with private Docker daemons.

## Supported agents

| Agent       | Command    | Status       | Notes                      |
| ----------- | ---------- | ------------ | -------------------------- |
| Claude Code | `claude`   | Experimental | Most tested implementation |
| Codex       | `codex`    | Experimental | In development             |
| Copilot     | `copilot`  | Experimental | In development             |
| Gemini      | `gemini`   | Experimental | In development             |
| cagent      | `cagent`   | Experimental | In development             |
| Kiro        | `kiro`     | Experimental | In development             |

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
$ docker sandbox create AGENT [PATH]
```

Each agent runs in its own isolated sandbox. The agent type is bound to the
sandbox when created and cannot be changed later.

## Agent-specific configuration

Different agents may require different authentication methods or configuration.
See the agent-specific documentation:

- [Claude Code configuration](claude-code.md)

## Requirements

- Docker Desktop 4.58 or later
- Platform support:
  - macOS with virtualization.framework
  - Windows with Hyper-V {{< badge color=violet text=Experimental >}}
- API keys or credentials for your chosen agent

## Next steps

- [Claude Code configuration](claude-code.md)
- [Custom templates](templates.md)
- [Using sandboxes effectively](workflows.md)
