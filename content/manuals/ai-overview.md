---
title: Docker AI overview
linkTitle: Overview
description: Docker's AI tools help you build, run, and manage AI-powered applications and workflows.
weight: 1
keywords: docker, ai, gordon, docker agent, sandboxes, model runner, mcp
params:
  sidebar:
    group: AI
---

Docker provides tools for working with AI across your development workflow.
Each tool serves a different purpose.

## Which tool do I need?

| I want to...                                                    | Use                                                      | CLI command      |
| --------------------------------------------------------------- | -------------------------------------------------------- | ---------------- |
| Get AI help with Docker tasks (containers, images, Dockerfiles) | [Gordon](./ai/gordon/)                                   | `docker ai`      |
| Run AI models locally with an OpenAI-compatible API             | [Model Runner](./ai/model-runner/)                       | `docker model`   |
| Connect AI tools to external services via MCP                   | [MCP Catalog and Toolkit](./ai/mcp-catalog-and-toolkit/) | `docker mcp`     |
| Build and orchestrate custom multi-agent teams                  | [Docker Agent](./ai/docker-agent/)                       | `docker agent`   |
| Run coding agents in isolated environments                      | [Docker Sandboxes](./ai/sandboxes/)                      | `docker sandbox` |

## How these tools relate

**Gordon** is Docker's built-in AI assistant. It helps with Docker-specific
tasks like debugging containers, writing Dockerfiles, and managing images. You
interact with it through Docker Desktop or the `docker ai` command.

**Docker Agent** is an open-source framework for defining teams of AI agents
in YAML. You configure agents with specific roles, models, and tools, then
run them from your terminal. Docker Agent is a general-purpose agent runtime,
not specific to Docker tasks.

**Docker Sandboxes** provides isolated microVM environments for running coding
agents. It supports multiple agents including Claude Code, Codex, Copilot,
Gemini, and Docker Agent. Sandboxes is the isolation layer — the agents
themselves are separate tools.

**Model Runner** lets you run LLMs locally. Other tools like Docker Agent can
use Model Runner as a model provider.

**MCP Catalog and Toolkit** manages connections between AI tools and external
services using the Model Context Protocol. Gordon, Docker Agent, and
third-party tools can all use MCP servers configured through the Toolkit.
