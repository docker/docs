---
title: Builtin tools
description: cagent's builtin tools
keywords: [ai, agent, cagent]
weight: 10
---

Tools are what make agents useful. They give your agent the ability to interact
with external systems, execute commands, access files, fetch web content, and
much more. Without tools, an agent can only generate text-with tools, it can
take action.

## What are Toolsets?

`cagent` organizes tools into **toolsets** - logical groups of related tools
that work together to accomplish specific tasks. For example:

- The **filesystem** toolset provides tools for reading, writing, searching, and
  managing files
- The **shell** toolset enables command execution in your terminal
- The **memory** toolset allows agents to remember information across
  conversations
- An **MCP server** is a toolset that can provide any number of custom tools

Toolsets are configured in your agent YAML file and determine what capabilities
your agent has access to.

## Types of Toolsets

`cagent` supports several types of toolsets:

### Builtin Toolsets

Built directly into cagent, these toolsets provide core functionality:

| Toolset          | Description                                                   |
| ---------------- | ------------------------------------------------------------- |
| `filesystem`     | Read, write, search, and manage files and directories         |
| `shell`          | Execute shell commands in your environment                    |
| `memory`         | Store and retrieve persistent information about users         |
| `fetch`          | Retrieve content from HTTP/HTTPS URLs                         |
| `think`          | Reasoning scratchpad for complex planning                     |
| `todo`           | Task tracking for multi-step operations                       |
| `script_shell`   | Define custom parameterized shell commands as tools           |

[Learn more about builtin toolsets →](/docs/tools/builtin/filesystem)

### MCP Servers

The [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) is an open
standard for connecting AI assistants to external systems. MCP servers act as
toolsets in cagent, providing standardized access to a wide ecosystem of tools.

**Local MCP Servers (stdio):** Run as local processes that communicate via
standard input/output. Great for accessing local resources and services.

**Remote MCP Servers (SSE/HTTP):** Connect to remote servers over HTTP, enabling
access to web APIs, databases, and cloud services.

[Learn more about MCP servers →](/docs/mcp)

### Custom Shell Scripts

Using the `script_shell` toolset, you can define your own custom tools that
execute shell commands with typed parameters:

```yaml
toolsets:
  - type: script_shell
    tools:
      deploy:
        cmd: "./deploy.sh"
        description: "Deploy the application"
        args:
          environment:
            type: string
            description: "Target environment"
        required:
          - environment
```

[Learn more about custom shell tools →](/docs/tools/builtin/script_shell)

## Configuring Toolsets

Toolsets are configured in your agent's YAML file under the `toolsets` array:

```yaml
agents:
  my_agent:
    model: gpt-4o
    description: "A helpful coding assistant"
    toolsets:
      # Builtin toolset - simple type reference
      - type: filesystem
      
      # Builtin toolset with configuration
      - type: memory
        path: "./memories.db"
      
      # Local MCP server (stdio)
      - type: mcp
        command: npx
        args:
          - "-y"
          - "@modelcontextprotocol/server-filesystem"
          - "/path/to/directory"
      
      # Remote MCP server (SSE)
      - type: mcp
        remote:
          url: "https://api.example.com/mcp"
          transport_type: sse
          headers:
            Authorization: "Bearer ${API_TOKEN}"
      
      # Custom shell tools
      - type: script_shell
        tools:
          build:
            cmd: "npm run build"
            description: "Build the project"
```

## Toolset Configuration Options

Each toolset type may have specific configuration options. Common options are:

- `instruction`: Additional instructions for using the toolset (optional)
- `tools`: Array of specific tool names to enable (optional, defaults to all)
- `env`: Environment variables for the toolset (optional)

## Tool Selection

By default, agents have access to all tools provided by their configured
toolsets. You can restrict this using the `tools` option:

```yaml
toolsets:
  - type: filesystem
    tools:
      - read_file
      - write_file
      - list_directory
    # Agent only gets these three filesystem tools
```

This is useful for:
- Limiting agent capabilities for security
- Reducing context size for smaller models
- Creating specialized agents with focused tool access

## Best Practices

### Performance

- **Choose appropriate toolsets**: Don't load toolsets the agent won't use
- **Limit tool selection**: Use the `tools` array to restrict available tools
- **Consider model capabilities**: Smaller models may struggle with too many
  tools

## Multi-Agent Systems

Different agents in a multi-agent system can have different toolsets:

```yaml
agents:
  coordinator:
    model: gpt-4o
    sub_agents:
      - code_writer
      - code_reviewer
    toolsets:
      - type: transfer_task
  
  code_writer:
    model: gpt-4o
    toolsets:
      - type: filesystem
      - type: shell
  
  code_reviewer:
    model: gpt-4o
    toolsets:
      - type: filesystem
        tools:
          - read_file
          - read_multiple_files
```

This allows you to:
- Create specialized agents with focused capabilities
- Implement security boundaries between agents
- Optimize performance by limiting each agent's toolset

