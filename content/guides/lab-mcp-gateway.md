---
title: "Lab: Docker MCP Gateway"
linkTitle: "Lab: Docker MCP Gateway"
description: |
  Run containerized MCP servers safely and securely with the Docker MCP Gateway
  in this hands-on interactive lab.
summary: |
  Hands-on lab: Configure, secure, and connect MCP servers to your agentic
  applications using the Docker MCP Gateway.
keywords: AI, Docker, MCP, MCP Gateway, MCP servers, lab, labspace
aliases:
  - /labs/docker-for-ai/mcp-gateway/
params:
  tags: [ai, labs]
  time: 30 minutes
  resource_links:
    - title: Docker MCP Gateway docs
      url: /ai/mcp-gateway/
    - title: MCP Gateway GitHub
      url: https://github.com/docker/mcp-gateway
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-mcp-gateway
---

This lab provides a comprehensive, hands-on overview of the Docker MCP Gateway,
which enables you to run containerized MCP servers safely and securely. Learn
how to configure, secure, and connect MCP servers to your agentic applications.

## What you'll learn

- Learn about the Docker MCP Gateway and its architecture
- Run the MCP Gateway with a simple MCP server
- Securely inject secrets into MCP servers
- Filter tools to reduce noise and save tokens
- Connect the MCP Gateway to your application using popular agentic frameworks
- Configure and use custom MCP servers

## Modules

| # | Module | Description |
|---|--------|-------------|
| 1 | Introduction | Overview of the MCP Gateway and why it matters |
| 2 | Adding a Simple MCP Server | Get started with a basic MCP server configuration |
| 3 | Adding a Complex MCP Server | Configure MCP servers with secrets and advanced options |
| 4 | Filtering Available Tools | Reduce noise and save tokens by filtering tool availability |
| 5 | Connecting MCP Gateway to Your App | Integrate the MCP Gateway with agentic frameworks |
| 6 | Using a Custom MCP Server | Build and run your own custom MCP server |
| 7 | Conclusion | Summary and next steps |

## Prerequisites

- Docker Desktop (latest version) with Docker Model Runner enabled
- Basic familiarity with Docker and Docker Compose

## Launch the lab

Start the labspace:

```console
$ docker compose -f oci://dockersamples/labspace-mcp-gateway up -d
```

Then open your browser to [http://localhost:3030](http://localhost:3030).
