---
title: What next after the Docker workshop
weight: 100
linkTitle: "Part 9: What next"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers,
  docker desktop, AI, model runner, MCP, agents, hardened images, security
description: Explore what to do next after completing the Docker workshop, including securing your images, AI development, and language-specific guides.
aliases:
 - /get-started/11_what_next/
 - /guides/workshop/10_what_next/
summary: |
  Now that you've completed the Docker workshop, you're ready to explore
  securing your images with Docker Hardened Images, building AI-powered
  applications, and diving into language-specific guides.
notoc: true

secure-images:
- title: What are Docker Hardened Images?
  description: Understand secure, minimal, production-ready base images with near-zero CVEs.
  link: /dhi/explore/what/
- title: Get started with DHI
  description: Pull and run your first Docker Hardened Image in minutes.
  link: /dhi/get-started/
- title: Use hardened images
  description: Learn how to use DHI in your Dockerfiles and CI/CD pipelines.
  link: /dhi/how-to/use/
- title: Explore the DHI catalog
  description: Browse available hardened images, variants, and security attestations.
  link: /dhi/how-to/explore/

ai-development:
- title: Docker Model Runner
  description: Run and manage AI models locally using familiar Docker commands with OpenAI-compatible APIs.
  link: /ai/model-runner/
- title: MCP Toolkit
  description: Set up, manage, and run containerized MCP servers to power your AI agents.
  link: /ai/mcp-catalog-and-toolkit/toolkit/
- title: Build AI agents with cagent
  description: Create teams of specialized AI agents that collaborate to solve complex problems.
  link: /ai/cagent/
- title: Use AI models in Compose
  description: Define AI model dependencies in your Docker Compose applications.
  link: /compose/how-tos/model-runner/

language-guides:
- title: Node.js
  description: Learn how to containerize and develop Node.js applications.
  link: /guides/language/nodejs/
- title: Python
  description: Build and run Python applications in containers.
  link: /guides/language/python/
- title: Java
  description: Containerize Java applications with best practices.
  link: /guides/language/java/
- title: Go
  description: Develop and deploy Go applications using Docker.
  link: /guides/language/golang/
---

Congratulations on completing the Docker workshop. You've learned how to containerize applications, work with multi-container setups, use Docker Compose, and apply image-building best practices.

Here's what to explore next.

## Secure your images

Take your image-building skills to the next level with Docker Hardened Images—secure, minimal, and production-ready base images that are now free for everyone.

{{< grid items="secure-images" >}}

## Build with AI

Docker makes it easy to run AI models locally and build agentic AI applications. Explore Docker's AI tools and start building AI-powered apps.

{{< grid items="ai-development" >}}

## Language-specific guides

Apply what you've learned to your preferred programming language with hands-on tutorials.

{{< grid items="language-guides" >}}
