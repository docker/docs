---
title: "MCP Mode"
description: "Expose your docker-agent agents as MCP tools for use in Claude Desktop, Claude Code, and other MCP-compatible applications."
keywords: docker agent, ai agents, features, mcp mode
weight: 40
canonical: https://docs.docker.com/ai/docker-agent/features/mcp-mode/
---

_Expose your docker-agent agents as MCP tools for use in Claude Desktop, Claude Code, and other MCP-compatible applications._

## Why MCP Mode?

The `docker agent serve mcp` command makes your agents available to any application that supports the [Model Context Protocol](https://modelcontextprotocol.io/). This means you can:

- Use custom agents directly within **Claude Desktop** or **Claude Code**
- Share specialized agents across different applications
- Build reusable agent teams consumable from any MCP client
- Integrate domain-specific agents into existing workflows

> [!NOTE]
> **What is MCP?**
>
> The [Model Context Protocol](https://modelcontextprotocol.io/) is an open standard for connecting AI tools. See also [Remote MCP Servers](../remote-mcp/index.md) for connecting to cloud services.

## Basic Usage

```bash
# Expose a local config (stdio transport, the default)
$ docker agent serve mcp ./agent.yaml

# Expose from a registry
$ docker agent serve mcp agentcatalog/pirate

# Set the working directory
$ docker agent serve mcp ./agent.yaml --working-dir /path/to/project
```

## Transports

By default, `serve mcp` uses the stdio transport — ideal for clients that spawn the server as a subprocess (Claude Desktop, Claude Code, Cursor, …).

To expose the MCP server over streaming HTTP instead, pass `--http`:

```bash
# Streaming HTTP transport on the default 127.0.0.1:8081
$ docker agent serve mcp ./agent.yaml --http

# Override the listen address / port
$ docker agent serve mcp ./agent.yaml --http --listen 0.0.0.0:9090
```

| Flag                   | Default            | Description                                                                                                  |
| ---------------------- | ------------------ | ------------------------------------------------------------------------------------------------------------ |
| `--http`               | `false`            | Use streaming HTTP transport instead of stdio.                                                               |
| `-l`, `--listen`       | `127.0.0.1:8081`   | Address to listen on when `--http` is enabled.                                                               |
| `-a`, `--agent`        | all agents         | Expose a single named agent instead of every agent in the config.                                            |
| `--tool-name`          | (none)             | Override the MCP tool identifier clients call (defaults to agent name); only valid when exposing one agent.  |
| `--mcp-keepalive`      | `0`                | Interval between MCP keep-alive pings (e.g. `30s`); `0` disables keep-alive.                                 |

Runtime configuration flags such as `--working-dir`, `--env-from-file`, `--models-gateway`, and hook flags are also available — see the [CLI reference](../cli/index.md).

## Using with Claude Desktop

Add a configuration to your Claude Desktop MCP settings file:

- **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "myagent": {
      "command": "/usr/local/bin/docker",
      "args": [
        "agent", 
        "serve",
        "mcp",
        "agentcatalog/coder",
        "--working-dir",
        "/home/user/projects"
      ],
      "env": {
        "ANTHROPIC_API_KEY": "your_key_here",
        "OPENAI_API_KEY": "your_key_here"
      }
    }
  }
}
```

Restart Claude Desktop after updating the configuration.

## Using with Claude Code

```bash
$ claude mcp add --transport stdio myagent \
  --env OPENAI_API_KEY=$OPENAI_API_KEY \
  --env ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY \
  -- docker agent serve mcp agentcatalog/pirate --working-dir $(pwd)
```

## Multi-Agent in MCP Mode

When you expose a multi-agent configuration via MCP, each agent becomes a separate tool in the MCP client:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Main coordinator
    sub_agents: [designer, engineer]
  designer:
    model: openai/gpt-5-mini
    description: UI/UX design specialist
  engineer:
    model: anthropic/claude-sonnet-4-5
    description: Software engineer
```

All three agents (`root`, `designer`, `engineer`) appear as separate tools in Claude Desktop or Claude Code.

## Troubleshooting

- **Agents not appearing:** Verify the `docker-agent` binary path and restart the MCP client
- **Permission errors:** Ensure `docker-agent` has execute permissions (`chmod +x`)
- **Missing API keys:** Pass all required keys in the `env` section
- **Working directory issues:** Verify the `--working-dir` path exists and is accessible
