---
title: "Lab: Getting Started with Docker Agent"
linkTitle: "Lab: Getting Started with Docker Agent"
description: |
  Build intelligent multi-agent teams with Docker Agent and Docker in this hands-on
  interactive lab.
summary: |
  Hands-on lab: Create, share, and orchestrate intelligent AI agents using
  Docker Agent, MCP Toolkit, and Docker.
keywords: AI, Docker, docker agent, agents, multi-agent, MCP Toolkit, lab, labspace
aliases:
  - /labs/docker-for-ai/cagent/
  - /guides/lab-cagent/
params:
  tags: [ai, labs]
  time: 20 minutes
  resource_links:
    - title: Docker Agent documentation
      url: https://github.com/docker/docker-agent
    - title: Docker MCP Toolkit
      url: https://docs.docker.com/ai/mcp-catalog-and-toolkit/toolkit/
    - title: Labspace repository
      url: https://github.com/ajeetraina/labspace-cagent
---

This lab walks you through building intelligent agents with Docker Agent. You'll learn beginner
agent concepts, then build sophisticated multi-agent teams that handle complex
real-world tasks. Learn how to create, share, and orchestrate AI agents with
Docker.

## What you'll learn

- Create simple agents with Docker Agent
- Use built-in generic agentic tools for common tasks
- Integrate MCP servers from the MCP Toolkit
- Share agents using the Docker Registry
- Build multi-agent systems for complex workflows
- Use Docker Model Runner with Docker Agent (preview)

## Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | Introduction | Overview of Docker Agent and intelligent agent concepts |
| 2 | Getting Started | Create your first agent with Docker Agent |
| 3 | Using Built-in Tools | Leverage the generic agentic tools in Docker Agent |
| 4 | Using MCP | Integrate MCP servers from the MCP Toolkit |
| 5 | Sharing Agents | Package and share agents via Docker Registry |
| 6 | Introduction to Sub-agents | Build multi-agent systems with sub-agent orchestration |
| 7 | Conclusion | Summary and next steps |

## Prerequisites

- Latest version of Docker Desktop
- Basic familiarity with Docker

## Launch the lab

Start the labspace:

```console
$ docker compose -f oci://dockersamples/labspace-cagent up -d
```

Then open your browser to [http://localhost:3030](http://localhost:3030).
