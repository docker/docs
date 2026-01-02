---
title: Sharing agents
description: Distribute agent configurations through OCI registries
keywords: [cagent, oci, registry, docker hub, sharing, distribution]
weight: 50
---

Push your agent to a registry and share it by name. Your teammates
reference `agentcatalog/security-expert` instead of copying YAML files
around or asking you where your agent configuration lives.

When you update the agent in the registry, everyone gets the new version
the next time they pull or restart their client.

## Prerequisites

To push agents to a registry, authenticate first:

```console
$ docker login
```

For other registries, use their authentication method.

## Publishing agents

Push your agent configuration to a registry:

```console
$ cagent push ./agent.yml myusername/agent-name
```

Push creates the repository if it doesn't exist yet. Use Docker Hub or
any OCI-compatible registry.

Tag specific versions:

```console
$ cagent push ./agent.yml myusername/agent-name:v1.0.0
$ cagent push ./agent.yml myusername/agent-name:latest
```

## Using published agents

Pull an agent to inspect it locally:

```console
$ cagent pull agentcatalog/pirate
```

This saves the configuration as a local YAML file.

Run agents directly from the registry:

```console
$ cagent run agentcatalog/pirate
```

Or reference it directly in integrations:

### Editor integration (ACP)

Use registry references in ACP configurations so your editor always uses
the latest version:

```json
{
  "agent_servers": {
    "myagent": {
      "command": "cagent",
      "args": ["acp", "agentcatalog/pirate"]
    }
  }
}
```

### MCP client integration

Agents can be exposed as tools in MCP clients:

```json
{
  "mcpServers": {
    "myagent": {
      "command": "/usr/local/bin/cagent",
      "args": ["mcp", "agentcatalog/pirate"]
    }
  }
}
```

## What's next

- Set up [ACP integration](./integrations/acp.md) with shared agents
- Configure [MCP integration](./integrations/mcp.md) with shared agents
- Browse the [agent catalog](https://hub.docker.com/u/agentcatalog) for examples
