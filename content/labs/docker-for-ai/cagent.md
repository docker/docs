---
title: Getting Started with cagent
description: |
  Build intelligent multi-agent teams with cagent and Docker.
weight: 40
---

This lab walks you through building intelligent agents using cagent — from basic
agent concepts to building sophisticated multi-agent teams that handle complex
real-world tasks. Learn how to create, share, and orchestrate AI agents with
Docker.

## What you'll learn

- Create simple agents with cagent
- Use built-in generic agentic tools for common tasks
- Integrate MCP servers from the MCP Toolkit
- Share agents using the Docker Registry
- Build multi-agent systems for complex workflows
- Use Docker Model Runner with cagent (preview)

## Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | Introduction | Overview of cagent and intelligent agent concepts |
| 2 | Getting Started | Create your first agent with cagent |
| 3 | Using Built-in Tools | Leverage the generic agentic tools in cagent |
| 4 | Using MCP | Integrate MCP servers from the MCP Toolkit |
| 5 | Sharing Agents | Package and share agents via Docker Registry |
| 6 | Introduction to Sub-agents | Build multi-agent systems with sub-agent orchestration |
| 7 | Conclusion | Summary and next steps |

## Launch the lab

1. Clone the repository:

```console
$ git clone https://github.com/ajeetraina/labspace-cagent
$ cd labspace-cagent
```

2. Start the labspace:

```console
$ docker compose -f oci://dockersamples/labspace-cagent up -d
```

Then open your browser to [http://localhost:3030](http://localhost:3030).

## Resources

- [Labspace repository](https://github.com/ajeetraina/labspace-cagent)
- [cagent documentation](https://github.com/docker/cagent)
- [Docker MCP Toolkit](https://hub.docker.com/extensions/docker/labs-ai-tools-for-devs)

<div id="docker-ai-labs-survey-anchor"></div>
