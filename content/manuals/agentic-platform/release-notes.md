---
title: Docker Agentic Platform release notes
linkTitle: Release notes
description: Review feature updates, behavior changes, and fixes in Docker Agentic Platform.
keywords: docker agentic platform, release notes, updates, fixes
weight: 90
---

## August 26, 2026

Docker Agentic Platform is available for running agent and tool workloads in
sandboxes hosted on Docker-managed cloud infrastructure.

The initial release includes:

- Predefined sandbox types for Claude Code, Codex, OpenCode, Copilot, and
  Gemini CLI
- A live terminal for interacting with the workload running in a sandbox
- Pause, resume, and delete controls for sandboxes
- Automatic sandbox stop or deletion after a timer from 1 to 24 hours
- Predefined MCP servers and custom servers added by URL, with connection and
  authorization flows
- Provider and service credential management under **Secrets**, with values
  kept outside sandboxes
- Read-only kit policies that apply automatically, Docker-managed Open and
  Balanced user policies, selectable custom policies, and deny-over-allow
  precedence
- Metered sandbox compute with selectable instance sizes
