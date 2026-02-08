---
title: Toolsets reference
linkTitle: Toolsets
description: Complete reference for cagent toolsets and their capabilities
keywords: [ai, agent, cagent, tools, toolsets]
weight: 20
---

This reference documents the toolsets available in cagent and what each one
does. Tools give agents the ability to take action—interacting with files,
executing commands, accessing external resources, and managing state.

For configuration file syntax and how to set up toolsets in your agent YAML,
see the [Configuration file reference](./config.md).

## How agents use tools

When you configure toolsets for an agent, those tools become available in the
agent's context. The agent can invoke tools by name with appropriate parameters
based on the task at hand.

Tool invocation flow:

1. Agent analyzes the task and determines which tool to use
2. Agent constructs tool parameters based on requirements
3. cagent executes the tool and returns results
4. Agent processes results and decides next steps

Agents can call multiple tools in sequence or make decisions based on tool
results. Tool selection is automatic based on the agent's understanding of the
task and available capabilities.

## Tool types

cagent supports three types of toolsets:

Built-in toolsets
: Core functionality built directly into cagent (`filesystem`, `shell`,
`memory`, etc.). These provide essential capabilities for file operations,
command execution, and state management.
MCP toolsets
: Tools provided by Model Context Protocol servers, either local processes
(stdio) or remote servers (HTTP/SSE). MCP enables access to a wide ecosystem
of standardized tools.
Custom toolsets
: Shell scripts wrapped as tools with typed parameters (`script`). This
lets you define domain-specific tools for your use case.

## Configuration

Toolsets are configured in your agent's YAML file under the `toolsets` array:

```yaml
agents:
  my_agent:
    model: anthropic/claude-sonnet-4-5
    description: A helpful coding assistant
    toolsets:
      # Built-in toolset
      - type: filesystem

      # Built-in toolset with configuration
      - type: memory
        path: ./memories.db

      # Local MCP server (stdio)
      - type: mcp
        command: npx
        args: ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/dir"]

      # Remote MCP server (SSE)
      - type: mcp
        remote:
          url: https://mcp.example.com/sse
          transport_type: sse
          headers:
            Authorization: Bearer ${API_TOKEN}

      # Custom shell tools
      - type: script
        tools:
          build:
            cmd: npm run build
            description: Build the project
```

### Common configuration options

All toolset types support these optional properties:

| Property      | Type             | Description                                                                                                                                                                                                                         |
| ------------- | ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `instruction` | string           | Additional instructions for using the toolset                                                                                                                                                                                       |
| `tools`       | array            | Specific tool names to enable (defaults to all)                                                                                                                                                                                     |
| `env`         | object           | Environment variables for the toolset                                                                                                                                                                                               |
| `toon`        | string           | Comma-delimited regex patterns matching tool names whose JSON outputs should be compressed. Reduces token usage by simplifying/compressing JSON responses from matched tools using automatic encoding. Example: `"search.*,list.*"` |
| `defer`       | boolean or array | Control which tools load into initial context. Set to `true` to defer all tools, or array of tool names to defer specific tools. Deferred tools don't consume context until explicitly loaded via `search_tool`/`add_tool`.         |

### Tool selection

By default, agents have access to all tools from their configured toolsets. You
can restrict this using the `tools` option:

```yaml
toolsets:
  - type: filesystem
    tools: [read_file, write_file, list_directory]
```

This is useful for:

- Limiting agent capabilities for security
- Reducing context size for smaller models
- Creating specialized agents with focused tool access

### Deferred loading

Deferred loading keeps tools out of the initial context window, loading them
only when explicitly requested. This is useful for large toolsets where most
tools won't be used, significantly reducing context consumption.

Defer all tools in a toolset:

```yaml
toolsets:
  - type: mcp
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path"]
    defer: true # All tools load on-demand
```

Or defer specific tools while loading others immediately:

```yaml
toolsets:
  - type: mcp
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path"]
    defer: [search_files, list_directory] # Only these are deferred
```

Agents can discover deferred tools via `search_tool` and load them into context
via `add_tool` when needed. Best for toolsets with dozens of tools where only a
few are typically used.

### Output compression

The `toon` property compresses JSON outputs from matched tools to reduce token
usage. When a tool's output is JSON, it's automatically compressed using
efficient encoding before being returned to the agent:

```yaml
toolsets:
  - type: mcp
    command: npx
    args: ["-y", "@modelcontextprotocol/server-github"]
    toon: "search.*,list.*" # Compress outputs from search/list tools
```

Useful for tools that return large JSON responses (API results, file listings,
search results). The compression is transparent to the agent but can
significantly reduce context consumption for verbose tool outputs.

### Per-agent tool configuration

Different agents can have different toolsets:

```yaml
agents:
  coordinator:
    model: anthropic/claude-sonnet-4-5
    sub_agents: [code_writer, code_reviewer]
    toolsets:
      - type: filesystem
        tools: [read_file]

  code_writer:
    model: openai/gpt-5-mini
    toolsets:
      - type: filesystem
      - type: shell

  code_reviewer:
    model: anthropic/claude-sonnet-4-5
    toolsets:
      - type: filesystem
        tools: [read_file, read_multiple_files]
```

This allows specialized agents with focused capabilities, security boundaries,
and optimized performance.

## Built-in tools reference

### Filesystem

The `filesystem` toolset gives your agent the ability to work with
files and directories. Your agent can read files to understand
context, write new files, make targeted edits to existing files,
search for content, and explore directory structures. Essential for
code analysis, documentation updates, configuration management, and
any agent that needs to understand or modify project files.

Access is restricted to the current working directory by default. Agents can
request access to additional directories at runtime, which requires your
approval.

#### Configuration

```yaml
toolsets:
  - type: filesystem

  # Optional: restrict to specific tools
  - type: filesystem
    tools: [read_file, write_file, edit_file]
```

### Shell

The `shell` toolset lets your agent execute commands in your system's shell
environment. Use this for agents that need to run builds, execute tests, manage
processes, interact with CLI tools, or perform system operations. The agent can
run commands in the foreground or background.

Commands execute in the current working directory and inherit environment
variables from the cagent process. This toolset is powerful but should be used
with appropriate security considerations.

#### Configuration

```yaml
toolsets:
  - type: shell
```

### Think

The `think` toolset provides your agent with a reasoning scratchpad. The agent
can record thoughts and reasoning steps without taking actions or modifying
data. Particularly useful for complex tasks where the agent needs to plan
multiple steps, verify requirements, or maintain context across a long
conversation.

Agents use this to break down problems, list applicable rules, verify they have
all needed information, and document their reasoning process before acting.

#### Configuration

```yaml
toolsets:
  - type: think
```

### Todo

The `todo` toolset gives your agent task-tracking capabilities for managing
multi-step operations. Your agent can break down complex work into discrete
tasks, track progress through each step, and ensure nothing is missed before
completing a request. Especially valuable for agents handling complex
workflows with multiple dependencies.

The `shared` option allows todos to persist across different agents in a
multi-agent system, enabling coordination.

#### Configuration

```yaml
toolsets:
  - type: todo

  # Optional: share todos across agents
  - type: todo
    shared: true
```

### Memory

The `memory` toolset allows your agent to store and retrieve information across
conversations and sessions. Your agent can remember user preferences, project
context, previous decisions, and other information that should persist. Useful
for agents that interact with users over time or need to maintain state about
a project or environment.

Memories are stored in a local database file and persist across cagent
sessions.

#### Configuration

```yaml
toolsets:
  - type: memory

  # Optional: specify database location
  - type: memory
    path: ./agent-memories.db
```

### Fetch

The `fetch` toolset enables your agent to retrieve content from HTTP/HTTPS URLs.
Your agent can fetch documentation, API responses, web pages, or any content
accessible via HTTP GET requests. Useful for agents that need to access
external resources, check API documentation, or retrieve web content.

The agent can specify custom HTTP headers when needed for authentication or
other purposes.

#### Configuration

```yaml
toolsets:
  - type: fetch
```

### User Prompt

The `user_prompt` toolset lets your agent ask you questions during task
execution. When the agent needs clarification, decisions, or additional
information it can't determine on its own, it displays a dialog and waits for
your response.

You'll see a prompt with the agent's question. Depending on what the agent
needs, you might provide free-form text, select from options, or fill out a
form with multiple fields. You can accept and provide the information, decline
to answer, or cancel the operation entirely.

#### Configuration

```yaml
toolsets:
  - type: user_prompt
```

No additional configuration is required. The tool becomes available to the
agent once configured. When the agent calls this tool, the user sees a dialog
with the prompt. The user can:

- **Accept**: Provide the requested information
- **Decline**: Refuse to provide the information
- **Cancel**: Cancel the operation

### API

The `api` toolset lets you define custom tools that call HTTP APIs. Similar to
`script` but for web services, this allows you to expose REST APIs,
webhooks, or any HTTP endpoint as a tool your agent can use. The agent sees
these as typed tools with automatic parameter validation.

Use this to integrate with external services, call internal APIs, trigger
webhooks, or interact with any HTTP-based system.

#### Configuration

Each API tool is defined with an `api_config` containing the endpoint, HTTP method, and optional typed parameters:

```yaml
toolsets:
  - type: api
    api_config:
      name: search_docs
      endpoint: https://api.example.com/search
      method: GET
      instruction: Search the documentation database
      headers:
        Authorization: Bearer ${API_TOKEN}
      args:
        query:
          type: string
          description: Search query
        limit:
          type: number
          description: Maximum results
      required: [query]

  - type: api
    api_config:
      name: create_ticket
      endpoint: https://api.example.com/tickets
      method: POST
      instruction: Create a support ticket
      args:
        title:
          type: string
          description: Ticket title
        description:
          type: string
          description: Ticket description
      required: [title, description]
```

For GET requests, parameters are interpolated into the endpoint URL. For POST
requests, parameters are sent as JSON in the request body.

Supported argument types: `string`, `number`, `boolean`, `array`, `object`.

### Script

The `script` toolset lets you define custom tools by wrapping shell
commands with typed parameters. This allows you to expose domain-specific
operations to your agent as first-class tools. The agent sees these custom
tools just like built-in tools, with parameter validation and type checking
handled automatically.

Use this to create tools for deployment scripts, build commands, test runners,
or any operation specific to your project or workflow.

#### Configuration

Each custom tool is defined with a command, description, and optional typed
parameters:

```yaml
toolsets:
  - type: script
    tools:
      deploy:
        cmd: ./deploy.sh
        description: Deploy the application to an environment
        args:
          environment:
            type: string
            description: Target environment (dev, staging, prod)
          version:
            type: string
            description: Version to deploy
        required: [environment]

      run_tests:
        cmd: npm test
        description: Run the test suite
        args:
          filter:
            type: string
            description: Test name filter pattern
```

Supported argument types: `string`, `number`, `boolean`, `array`, `object`.

#### Tools

The tools you define become available to your agent. In the previous example,
the agent would have access to `deploy` and `run_tests` tools.

## Automatic tools

Some tools are automatically added to agents based on their configuration. You
don't configure these explicitly—they appear when needed.

### transfer_task

Automatically available when your agent has `sub_agents` configured. Allows
the agent to delegate tasks to sub-agents and receive results back.

### handoff

Automatically available when your agent has `handoffs` configured. Allows the
agent to transfer the entire conversation to a different agent.

## What's next

- Read the [Configuration file reference](./config.md) for YAML file structure
- Review the [CLI reference](./cli.md) for running agents
- Explore [MCP servers](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md) for extended capabilities
- Browse [example configurations](https://github.com/docker/cagent/tree/main/examples)
