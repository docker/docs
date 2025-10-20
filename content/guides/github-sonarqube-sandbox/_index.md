---
title: Automate code quality workflows with GitHub and SonarQube in E2B sandboxes
linkTitle: GitHub and SonarQube quality checks
summary: Build AI-powered code quality workflows using E2B sandboxes with Docker's MCP catalog to automate GitHub and SonarQube integration.
description: Learn how to create E2B sandboxes with MCP servers, analyze code quality with SonarQube, and generate quality-gated pull requests using GitHub—all through natural language interactions with Claude.
tags: [devops]
params:
  featured: true
  time: 40 minutes
  image:
  resource_links:
    - title: E2B Documentation
      url: https://e2b.dev/docs
    - title: Docker MCP Catalog
      url: https://hub.docker.com/mcp
    - title: Sandboxes
      url: https://docs.docker.com/ai/mcp-catalog-and-toolkit/sandboxes/
---

This guide demonstrates how to build an AI-powered code quality workflow using
[E2B sandboxes](https://e2b.dev/docs) with Docker’s MCP catalog. You’ll create
a system that automatically analyzes code quality issues in GitHub repositories
using SonarQube, then generate pull requests with fixes.

## What you'll build

You'll build a Node.js script that spins up an E2B sandbox, connects GitHub
and SonarQube MCP servers, and uses Claude to analyze code quality and propose
improvements.

## What you'll learn

In this guide, you'll learn:

- How to create E2B sandboxes with multiple MCP servers
- How to configure GitHub and SonarQube MCP servers for AI workflows
- How to use Claude CLI inside sandboxes to interact with external tools
- How to build automated code review workflows that create quality-gated
pull requests

## Why use E2B sandboxes?

Running this workflow in E2B sandboes provides several advantages over
local execution:

- Security: AI-generated code runs in isolated containers, protecting your
local environment and credentials
- Zero setup: No need to install SonarQube, GitHub CLI, or manage dependencies
locally
- Scalability: Resource-intensive operations like code scanning run in the
cloud without consuming local resources
