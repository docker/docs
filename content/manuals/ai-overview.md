---
title: Docker AI overview
linkTitle: Overview
description: Docker's AI tools help you build, run, and manage AI-powered applications and workflows.
weight: 1
keywords: docker, ai, gordon, docker agent, sandboxes, governance, agentic platform, model runner, mcp
params:
  sidebar:
    group: AI and agents
grid:
  - title: "Docker Sandboxes"
    description: "Run coding agents in isolated environments with the `sbx` CLI."
    link: /ai/sandboxes/
    icon: terminal
  - title: "Sandboxes API and SDK"
    description: "Create and manage cloud sandboxes from your applications."
    link: /ai/sandboxes-api/
    icon: code
  - title: "AI governance"
    description: "Control sandbox access and audit actions with `sbx policy` and Docker Home."
    link: /ai/governance/
    icon: settings
  - title: "Docker Agentic Platform"
    description: "Run agents in cloud sandboxes through an experimental web Console."
    link: /agentic-platform/
    icon: cloud
  - title: "Gordon"
    description: "Get AI help with Docker tasks in Docker Desktop or with `docker ai`."
    link: /ai/gordon/
    icon: terminal
  - title: "Model Runner"
    description: "Run AI models locally with an OpenAI-compatible API and `docker model`."
    link: /ai/model-runner/
    icon: box
  - title: "MCP Catalog and Toolkit"
    description: "Connect AI tools to external services through MCP and `docker mcp`."
    link: /ai/mcp-catalog-and-toolkit/
    icon: box
  - title: "Docker Agent"
    description: "Build and orchestrate custom agent teams with `docker agent`."
    link: /ai/docker-agent/
    icon: code
  - title: "Docker Skills"
    description: "Give compatible coding agents guidance for Docker tasks."
    link: /ai/skills/
    icon: book-open
---

Docker provides tools for working with AI across your development workflow.
Each tool serves a different purpose.

## Which tool do I need?

{{< grid >}}

## How these tools relate

**Gordon** is Docker's built-in AI assistant. It helps with Docker-specific
tasks like debugging containers, writing Dockerfiles, and managing images. You
interact with it through Docker Desktop or the `docker ai` command.

**Docker Agent** is an open-source framework for defining teams of AI agents
in YAML. You configure agents with specific roles, models, and tools, then
run them from your terminal. Docker Agent is a general-purpose agent runtime,
not specific to Docker tasks.

**Docker Skills** are Docker's official, open-source instructions for compatible
coding agents working on Docker tasks. Install them through your agent's
supported method; they can also be used by Docker Agent and agents running in
Docker Sandboxes. Browse the [Docker Skills
catalog](https://github.com/docker/skills#readme) for current guidance.

Docker Sandboxes provides isolated environments for running coding agents
[locally](./ai/sandboxes/get-started.md) or
[in the cloud](./ai/sandboxes/cloud/_index.md). Sandboxes is the isolation layer;
the agents themselves are separate tools. Agent configuration and supported
features differ between local and cloud sandboxes.

Docker Agentic Platform is an experimental service for running agents in
Docker-managed cloud sandboxes. Its web Console provides kit selection,
credentials, network policies, MCP tools, and sandbox lifecycle controls.
[Activate a subscription](./agentic-platform/signup.md) to use cloud compute,
billed on a pay-as-you-go basis.

**Model Runner** lets you run LLMs locally. Other tools like Docker Agent can
use Model Runner as a model provider.

**MCP Catalog and Toolkit** manages connections between AI tools and external
services using the Model Context Protocol. Gordon, Docker Agent, and
third-party tools can all use MCP servers configured through the Toolkit.
