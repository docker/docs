---
title: cagent
description: cagent lets you build, orchestrate, and share AI agents that work together as a team.
weight: 60
params:
  sidebar:
    group: Open source
keywords: [ai, agent, cagent]
---

`cagent` is lets you build, orchestrate, and share AI agents. Use it to build AI
agents that can work as a team. Build a root agent that can delegate
tasks to sub-agents. Each agent can use the model of your choice, with the
parameters of your choice.

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
   > You may need to give the release executable permissions.
   > On MacOS and Linux, do so by running:

     ```console
     chmod +x /path/to/downloads/cagent-linux-<arm/amd>64
     ```

1. Optional: Rename the binary as needed and update your PATH to include
   cagent's executable.

1. Set the following environment variables:

   ```bash
   # If using the Docker AI Gateway, set this env var or use
   # the `--models-gateway url_to_docker_ai_gateway` CLI flag

   export CAGENT_MODELS_GATEWAY=url_to_docker_ai_gateway

   # Alternatively, you to need set keys for remote inference services
   # These are not needed if you are using Docker AI Gateway

   export OPENAI_API_KEY=your_api_key_here    # For OpenAI models
   export ANTHROPIC_API_KEY=your_api_key_here # For Anthopic models
   export GOOGLE_API_KEY=your_api_key_here    # For Gemini models
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
   cagent run my-agent.yaml
   ```

## Create an agentic team

```yaml {title="agentic-team.yaml"}
agents:
  root:
    model: claude
    description: "Main coordinator agent that delegates tasks and manages workflow"
    instruction: |
      You are the root coordinator agent. Your job is to:
      1. Understand user requests and break them down into manageable tasks
      2. Delegate appropriate tasks to your helper agent
      3. Coordinate responses and ensure tasks are completed properly
      4. Provide final responses to the user
      When you receive a request, analyze what needs to be done and decide whether to:
      - Handle it yourself if it's simple
      - Delegate to the helper agent if it requires specific assistance
      - Break complex requests into multiple sub-tasks
    sub_agents: ["helper"]

  helper:
    model: claude
    description: "Assistant agent that helps with various tasks as directed by the root agent"
    instruction: |
      You are a helpful assistant agent. Your role is to:
      1. Complete specific tasks assigned by the root agent
      2. Provide detailed and accurate responses
      3. Ask for clarification if tasks are unclear
      4. Report back to the root agent with your results

      Focus on being thorough and helpful in whatever task you're given.

models:
  claude:
    provider: anthropic
    model: claude-sonnet-4-0
    max_tokens: 64000
```

[See the reference documentation](https://github.com/docker/cagent/blob/main/docs/user-guide.md#configuration-reference).

## Learn more about cagent

For more information, see the following documentation in the
[GitHub repository](https://github.com/docker/cagent):

- [User guide](https://github.com/docker/cagent/blob/main/docs/user-guide.md)
- [Configuration reference](https://github.com/docker/cagent/blob/main/docs/user-guide.md#configuration-reference)
