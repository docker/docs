---
title: cagent
description: cagent lets you build, orchestrate, and share AI agents that work together as a team.
weight: 60
params:
  sidebar:
    group: Open source
keywords: [ai, agent, cagent]
---

{{< summary-bar feature_name="cagent" >}}

[cagent](https://github.com/docker/cagent) lets you build, orchestrate, and share
AI agents. You can use it to define AI agents that work as a team.

cagent relies on the concept of a _root agent_ that acts as a team lead and
delegates tasks to the sub-agents you define.
Each agent:
- uses the model of your choice, with the parameters of your choice.
- has access to the [built-in tools](#built-in-tools) and MCP servers
  configured in the [Docker MCP gateway](/manuals/ai/mcp-gateway/_index.md).
- works in its own context. They do not share knowledge.

The root agent is your main contact point. Each agent has its own context,
they don't share knowledge.

## Key features

- ️Multi-tenant architecture with client isolation and session management.
- Rich tool ecosystem via Model Context Protocol (MCP) integration.
- Hierarchical agent system with intelligent task delegation.
- Multiple interfaces including CLI, TUI, API server, and MCP server.
- Agent distribution via Docker registry integration.
- Security-first design with proper client scoping and resource isolation.
- Event-driven streaming for real-time interactions.
- Multi-model support (OpenAI, Anthropic, Gemini, DMR, Docker AI Gateway).

## Get started with cagent

1. Download the [latest release](https://github.com/docker/cagent/releases)
   for your operating system.

   > [!NOTE]
   > You might need to give the binary executable permissions.
   > On macOS and Linux, run:

     ```console
     chmod +x /path/to/downloads/cagent-linux-<arm/amd>64
     ```

   > [!NOTE]
   > You can also build cagent from the source. See the [repository](https://github.com/docker/cagent?tab=readme-ov-file#build-from-source).

1. Optional: Rename the binary as needed and update your PATH to include
   cagent's executable.

1. Set the following environment variables:

   ```bash
   # If using the Docker AI Gateway, set this environment variable or use
   # the `--models-gateway <url_to_docker_ai_gateway>` CLI flag

   export CAGENT_MODELS_GATEWAY=<url_to_docker_ai_gateway>

   # Alternatively, set keys for remote inference services.
   # These are not needed if you are using Docker AI Gateway.

   export OPENAI_API_KEY=<your_api_key_here>    # For OpenAI models
   export ANTHROPIC_API_KEY=<your_api_key_here> # For Anthropic models
   export GOOGLE_API_KEY=<your_api_key_here>    # For Gemini models
   ```

1. Create an agent by saving this sample as `assistant.yaml`:

   ```yaml {title="assistant.yaml"}
   agents:
     root:
       model: openai/gpt-5-mini
       description: A helpful AI assistant
       instruction: |
         You are a knowledgeable assistant that helps users with various tasks.
         Be helpful, accurate, and concise in your responses.
   ```

1. Start your prompt with your agent:

   ```bash
   cagent run assistant.yaml
   ```

## Create an agentic team

You can use AI prompting to generate a team of agents with the `cagent new`
command:

```console
$ cagent new

For any feedback, visit: https://docker.qualtrics.com/jfe/form/SV_cNsCIg92nQemlfw

Welcome to cagent! (Ctrl+C to exit)

What should your agent/agent team do? (describe its purpose):

> I need a cross-functional feature team. The team owns a specific product
  feature end-to-end. Include the key responsibilities of each of the roles
  involved (engineers, designer, product manager, QA). Keep the description
  short, clear, and focused on how this team delivers value to users and the business.
```

Alternatively, you can write your configuration file manually. For example:

```yaml {title="agentic-team.yaml"}
agents:
  root:
    model: claude
    description: "Main coordinator agent that delegates tasks and manages workflow"
    instruction: |
      You are the root coordinator agent. Your job is to:
      1. Understand user requests and break them down into manageable tasks.
      2. Delegate appropriate tasks to your helper agent.
      3. Coordinate responses and ensure tasks are completed properly.
      4. Provide final responses to the user.
      When you receive a request, analyze what needs to be done and decide whether to:
      - Handle it yourself if it's simple.
      - Delegate to the helper agent if it requires specific assistance.
      - Break complex requests into multiple sub-tasks.
    sub_agents: ["helper"]

  helper:
    model: claude
    description: "Assistant agent that helps with various tasks as directed by the root agent"
    instruction: |
      You are a helpful assistant agent. Your role is to:
      1. Complete specific tasks assigned by the root agent.
      2. Provide detailed and accurate responses.
      3. Ask for clarification if tasks are unclear.
      4. Report back to the root agent with your results.

      Focus on being thorough and helpful in whatever task you're given.

models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-0
    max_tokens: 64000
```

[See the reference documentation](https://github.com/docker/cagent?tab=readme-ov-file#-configuration-reference).

## Built-in tools

cagent includes a set of built-in tools that enhance your agents' capabilities.
You don't need to configure any external MCP tools to use them.

```yaml
agents:
  root:
    # ... other config
    toolsets:
      - type: todo
      - type: transfer_task
```

### Think tool

The think tool allows agents to reason through problems step by step:

```yaml
agents:
  root:
    # ... other config
    toolsets:
      - type: think
```

### Todo tool

The todo tool helps agents manage task lists:

```yaml
agents:
  root:
    # ... other config
    toolsets:
      - type: todo
```

### Memory tool

The memory tool provides persistent storage:

```yaml
agents:
  root:
    # ... other config
    toolsets:
      - type: memory
        path: "./agent_memory.db"
```

### Task transfer tool

The task transfer tool is an internal tool that allows an agent to delegate a task
to sub-agents. To prevent an agent from delegating work, make sure it doesn't have
sub-agents defined in its configuration.

### Using tools via the Docker MCP Gateway

If you use the [Docker MCP gateway](/manuals/ai/mcp-gateway.md),
you can configure your agent to interact with the
gateway and use the MCP servers configured in it. See [docker mcp
gateway run](/reference/cli/docker/mcp/gateway/gateway_run.md).

For example, to enable an agent to use Duckduckgo via the MCP Gateway:

```yaml
toolsets:
  - type: mcp
    command: docker
    args: ["mcp", "gateway", "run", "--servers=duckduckgo"]
```

## CLI interactive commands

You can use the following CLI commands, during
CLI sessions with your agents:

| Command  | Description                              |
|----------|------------------------------------------|
| /exit    | Exit the program                         |
| /reset   | Clear conversation history               |
| /eval    | Save current conversation for evaluation |
| /compact | Compact the current session              |

## Share your agents

Agent configurations can be packaged and shared via Docker Hub.
Before you start, make sure you have a [Docker repository](/manuals/docker-hub/repos/create.md).

To push an agent:

```bash
cagent push ./<agent-file>.yaml <namespace>/<reponame>
```

To pull an agent to the current directory:

```bash
cagent pull <namespace>/<reponame>
```

The agent's configuration file is named `<namespace>_<reponame>.yaml`. Run
it with the `cagent run <filename>` command.

## Related pages

- For more information about cagent, see the
[GitHub repository](https://github.com/docker/cagent).
- [Docker MCP Gateway](/manuals/ai/mcp-gateway/_index.md)