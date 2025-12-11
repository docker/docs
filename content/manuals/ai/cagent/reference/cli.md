---
title: CLI reference
linkTitle: CLI
description: Complete reference for cagent command-line interface
keywords: [ai, agent, cagent, cli, command line]
weight: 30
---

Command-line interface for running, managing, and deploying AI agents.

For agent configuration file syntax, see the [Configuration file
reference](./config.md). For toolset capabilities, see the [Toolsets
reference](./toolsets.md).

## Synopsis

```console
$ cagent [command] [flags]
```

## Global flags

Work with all commands:

| Flag            | Type    | Default | Description          |
| --------------- | ------- | ------- | -------------------- |
| `-d`, `--debug` | boolean | false   | Enable debug logging |
| `-o`, `--otel`  | boolean | false   | Enable OpenTelemetry |
| `--log-file`    | string  | -       | Debug log file path  |

Debug logs write to `~/.cagent/cagent.debug.log` by default. Override with
`--log-file`.

## Runtime flags

Work with most commands. Supported commands link to this section.

| Flag                | Type    | Default | Description                          |
| ------------------- | ------- | ------- | ------------------------------------ |
| `--models-gateway`  | string  | -       | Models gateway address               |
| `--env-from-file`   | array   | -       | Load environment variables from file |
| `--code-mode-tools` | boolean | false   | Enable JavaScript tool orchestration |
| `--working-dir`     | string  | -       | Working directory for the session    |

Set `--models-gateway` via `CAGENT_MODELS_GATEWAY` environment variable.

## Commands

### a2a

Expose agent via the Agent2Agent (A2A) protocol. Allows other A2A-compatible
systems to discover and interact with your agent. Auto-selects an available
port if not specified.

```console
$ cagent a2a agent-file|registry-ref
```

> [!NOTE]
> A2A support is currently experimental and needs further work. Tool calls are
> handled internally and not exposed as separate ADK events. Some ADK features
> are not yet integrated.

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)

Flags:

| Flag            | Type    | Default | Description       |
| --------------- | ------- | ------- | ----------------- |
| `-a`, `--agent` | string  | root    | Agent name        |
| `--port`        | integer | 0       | Port (0 = random) |

Supports [runtime flags](#runtime-flags).

Examples:

```console
$ cagent a2a ./agent.yaml --port 8080
$ cagent a2a agentcatalog/pirate --port 9000
```

### acp

Start agent as ACP (Agent Client Protocol) server on stdio for editor integration.
See [ACP integration](../integrations/acp.md) for setup guides.

```console
$ cagent acp agent-file|registry-ref
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)

Supports [runtime flags](#runtime-flags).

### alias add

Create alias for agent.

```console
$ cagent alias add name target
```

Arguments:

- `name` - Alias name (required)
- `target` - Path to YAML or registry reference (required)

Examples:

```console
$ cagent alias add dev ./dev-agent.yaml
$ cagent alias add prod docker.io/user/prod-agent:latest
$ cagent alias add default ./agent.yaml
```

Setting alias name to "default" lets you run `cagent run` without arguments.

### alias list

List all aliases.

```console
$ cagent alias list
$ cagent alias ls
```

### alias remove

Remove alias.

```console
$ cagent alias remove name
$ cagent alias rm name
```

Arguments:

- `name` - Alias name (required)

### api

HTTP API server.

```console
$ cagent api agent-file|agents-dir
```

Arguments:

- `agent-file|agents-dir` - Path to YAML or directory with agents (required)

Flags:

| Flag                 | Type    | Default    | Description                       |
| -------------------- | ------- | ---------- | --------------------------------- |
| `-l`, `--listen`     | string  | :8080      | Listen address                    |
| `-s`, `--session-db` | string  | session.db | Session database path             |
| `--pull-interval`    | integer | 0          | Auto-pull OCI ref every N minutes |

Supports [runtime flags](#runtime-flags).

Examples:

```console
$ cagent api ./agent.yaml
$ cagent api ./agents/ --listen :9000
$ cagent api docker.io/user/agent --pull-interval 10
```

The `--pull-interval` flag works only with OCI references. Automatically pulls and reloads at the specified interval.

### build

Build Docker image for agent.

```console
$ cagent build agent-file|registry-ref [image-name]
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)
- `image-name` - Docker image name (optional)

Flags:

| Flag         | Type    | Default | Description                |
| ------------ | ------- | ------- | -------------------------- |
| `--dry-run`  | boolean | false   | Print Dockerfile only      |
| `--push`     | boolean | false   | Push image after build     |
| `--no-cache` | boolean | false   | Build without cache        |
| `--pull`     | boolean | false   | Pull all referenced images |

Example:

```console
$ cagent build ./agent.yaml myagent:latest
$ cagent build ./agent.yaml --dry-run
```

### catalog list

List catalog agents.

```console
$ cagent catalog list [org]
```

Arguments:

- `org` - Organization name (optional, default: `agentcatalog`)

Queries Docker Hub for agent repositories.

### debug config

Show resolved agent configuration.

```console
$ cagent debug config agent-file|registry-ref
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)

Supports [runtime flags](#runtime-flags).

Shows canonical configuration in YAML after all processing and defaults.

### debug toolsets

List agent tools.

```console
$ cagent debug toolsets agent-file|registry-ref
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)

Supports [runtime flags](#runtime-flags).

Lists all tools for each agent in the configuration.

### eval

Run evaluation tests.

```console
$ cagent eval agent-file|registry-ref [eval-dir]
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)
- `eval-dir` - Evaluation files directory (optional, default: `./evals`)

Supports [runtime flags](#runtime-flags).

### exec

Single message execution without TUI.

```console
$ cagent exec agent-file|registry-ref [message|-]
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)
- `message` - Prompt, or `-` for stdin (optional)

Same flags as [run](#run).

Supports [runtime flags](#runtime-flags).

Examples:

```console
$ cagent exec ./agent.yaml
$ cagent exec ./agent.yaml "Check for security issues"
$ echo "Instructions" | cagent exec ./agent.yaml -
```

### feedback

Submit feedback.

```console
$ cagent feedback
```

Shows link to submit feedback.

### mcp

MCP (Model Context Protocol) server on stdio. Exposes agents as tools to MCP
clients. See [MCP integration](../integrations/mcp.md) for setup guides.

```console
$ cagent mcp agent-file|registry-ref
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (required)

Supports [runtime flags](#runtime-flags).

Examples:

```console
$ cagent mcp ./agent.yaml
$ cagent mcp docker.io/user/agent:latest
```

### new

Create agent configuration interactively.

```console
$ cagent new [message...]
```

Flags:

| Flag               | Type    | Default | Description                     |
| ------------------ | ------- | ------- | ------------------------------- |
| `--model`          | string  | -       | Model as `provider/model`       |
| `--max-iterations` | integer | 0       | Maximum agentic loop iterations |

Supports [runtime flags](#runtime-flags).

Opens interactive TUI to configure and generate agent YAML.

### pull

Pull agent from OCI registry.

```console
$ cagent pull registry-ref
```

Arguments:

- `registry-ref` - OCI registry reference (required)

Flags:

| Flag      | Type    | Default | Description                 |
| --------- | ------- | ------- | --------------------------- |
| `--force` | boolean | false   | Pull even if already exists |

Example:

```console
$ cagent pull docker.io/user/agent:latest
```

Saves to local YAML file.

### push

Push agent to OCI registry.

```console
$ cagent push agent-file registry-ref
```

Arguments:

- `agent-file` - Path to local YAML (required)
- `registry-ref` - OCI reference like `docker.io/user/agent:latest` (required)

Example:

```console
$ cagent push ./agent.yaml docker.io/myuser/myagent:latest
```

### run

Interactive terminal UI for agent sessions.

```console
$ cagent run [agent-file|registry-ref] [message|-]
```

Arguments:

- `agent-file|registry-ref` - Path to YAML or OCI registry reference (optional)
- `message` - Initial prompt, or `-` for stdin (optional)

Flags:

| Flag            | Type    | Default | Description                  |
| --------------- | ------- | ------- | ---------------------------- |
| `-a`, `--agent` | string  | root    | Agent name                   |
| `--yolo`        | boolean | false   | Auto-approve all tool calls  |
| `--attach`      | string  | -       | Attach image file            |
| `--model`       | array   | -       | Override model (repeatable)  |
| `--dry-run`     | boolean | false   | Initialize without executing |
| `--remote`      | string  | -       | Remote runtime address       |

Supports [runtime flags](#runtime-flags).

Examples:

```console
$ cagent run ./agent.yaml
$ cagent run ./agent.yaml "Analyze this codebase"
$ cagent run ./agent.yaml --agent researcher
$ echo "Instructions" | cagent run ./agent.yaml -
$ cagent run
```

Running without arguments uses the default agent or a "default" alias if configured.

Shows interactive TUI in a terminal. Falls back to exec mode otherwise.

#### Interactive commands

TUI slash commands:

| Command    | Description                      |
| ---------- | -------------------------------- |
| `/exit`    | Exit                             |
| `/reset`   | Clear history                    |
| `/eval`    | Save conversation for evaluation |
| `/compact` | Compact conversation             |
| `/yolo`    | Toggle auto-approval             |

### version

Print version information.

```console
$ cagent version
```

Shows cagent version and commit hash.

## Environment variables

| Variable                       | Description                     |
| ------------------------------ | ------------------------------- |
| `CAGENT_MODELS_GATEWAY`        | Models gateway address          |
| `TELEMETRY_ENABLED`            | Telemetry control (set `false`) |
| `CAGENT_HIDE_TELEMETRY_BANNER` | Hide telemetry banner (set `1`) |
| `OTEL_EXPORTER_OTLP_ENDPOINT`  | OpenTelemetry endpoint          |

## Model overrides

Override models specified in your configuration file using the `--model` flag.

Format: `[agent=]provider/model`

Without an agent name, the model applies to all agents. With an agent name, it applies only to that specific agent.

Apply to all agents:

```console
$ cagent run ./agent.yaml --model gpt-5
$ cagent run ./agent.yaml --model anthropic/claude-sonnet-4-5
```

Apply to specific agents only:

```console
$ cagent run ./agent.yaml --model researcher=gpt-5
$ cagent run ./agent.yaml --model "agent1=gpt-5,agent2=claude-sonnet-4-5"
```

Providers: `openai`, `anthropic`, `google`, `dmr`

Omit provider for automatic selection based on model name.
