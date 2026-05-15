---
title: "Lab: Building Agentic Apps with Docker"
linkTitle: "Lab: Building Agentic Apps"
description: |
  Build agentic applications with Docker Model Runner, MCP Gateway, and Compose
  in this hands-on interactive lab.
summary: |
  Hands-on lab: Build agentic apps with Docker Model Runner, MCP Gateway, and
  Compose. Learn about models, tools, and agentic frameworks.
keywords: AI, Docker, Model Runner, MCP Gateway, agentic apps, lab, labspace
aliases:
  - /labs/docker-for-ai/agentic-apps/
params:
  tags: [ai, labs]
  time: 20 minutes
  resource_links:
    - title: Docker Model Runner docs
      url: /ai/model-runner/
    - title: Docker MCP Gateway docs
      url: /ai/mcp-gateway/
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-agentic-apps-with-docker
---

Get up and running with building agentic applications using Compose, Docker
Model Runner, and the Docker MCP Gateway. This hands-on lab takes you from
understanding AI models to building complete agentic applications.

## Launch the lab

{{< labspace-launch image="dockersamples/labspace-agentic-apps-with-docker" model-download="true" >}}

## What you'll learn

This lab covers three core areas of agentic application development:

**Models**: What models are, how to interact with them, configuring Docker
Model Runner in Compose, and writing code that connects to the Model Runner

**Tools**: Understanding tools and how they work, how MCP (Model Context
Protocol) fits in, configuring the Docker MCP Gateway, and connecting to the
MCP Gateway in code

**Code**: What agentic frameworks are, defining models and tools in a Compose
file, and configuring your app to use those models and tools

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
