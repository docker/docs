---
title: Building Agentic Apps with Docker
description: |
  Build agentic applications with Docker Model Runner, MCP Gateway, and Compose.
weight: 10
---

Get up and running with building agentic applications using Compose, Docker
Model Runner, and the Docker MCP Gateway. This hands-on lab takes you from
understanding AI models to building complete agentic applications.

## What you'll learn

This lab covers three core areas of agentic application development:

**Models** — What models are, how to interact with them, configuring Docker
Model Runner in Compose, and writing code that connects to the Model Runner.

**Tools** — Understanding tools and how they work, how MCP (Model Context
Protocol) fits in, configuring the Docker MCP Gateway, and connecting to the
MCP Gateway in code.

**Code** — What agentic frameworks are, defining models and tools in a Compose
file, and configuring your app to use those models and tools.

## Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | Introduction | Overview of agentic applications and the Docker AI stack |
| 2 | Understanding Model Interactions | Learn how to interact with AI models |
| 3 | The Docker Model Runner | Configure and use Docker Model Runner with Compose |
| 4 | Understanding Tools and MCP | Deep dive into tools, tool calling, and MCP |
| 5 | The Docker MCP Gateway | Set up and configure the MCP Gateway |
| 6 | Putting It All Together | Build a complete agentic application |
| 7 | Conclusion | Summary and next steps |

## Launch the lab

1. Clone the repository:

```console
$ git clone https://github.com/dockersamples/labspace-agentic-apps-with-docker
$ cd labspace-agentic-apps-with-docker
```

2. Start the labspace:

```console
$ docker compose -f oci://dockersamples/labspace-agentic-apps-with-docker up -d
```

Then open your browser to [http://localhost:3030](http://localhost:3030).

> **Note:** It may take a little while to start due to the AI model download.

## Resources

- [Labspace repository](https://github.com/dockersamples/labspace-agentic-apps-with-docker)
- [Docker Model Runner docs](/ai/model-runner/)
- [Docker MCP Gateway docs](/ai/mcp-gateway/)

<div id="docker-ai-labs-survey-anchor"></div>
