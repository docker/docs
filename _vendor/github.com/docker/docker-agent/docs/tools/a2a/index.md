---
title: "A2A Tool"
description: "Connect to remote agents via the Agent-to-Agent protocol."
keywords: docker agent, ai agents, tools, toolsets, a2a tool
linkTitle: "A2A"
weight: 60
canonical: https://docs.docker.com/ai/docker-agent/tools/a2a/
---

_Connect to remote agents via the Agent-to-Agent protocol._

## Overview

The A2A tool connects to a remote agent exposed over the A2A (Agent-to-Agent) protocol. Unlike [`handoff`](../handoff/index.md), which only targets local agents declared in the same config, `a2a` reaches out to an agent running on the network.

## Configuration

```yaml
toolsets:
  - type: a2a
    url: "http://localhost:8080/a2a"
    # Optional: custom tool name (defaults to a sanitized form of the URL / agent card name)
    name: research_agent
    # Optional: custom HTTP headers (typically for auth)
    headers:
      Authorization: "Bearer ${env.A2A_TOKEN}"
      X-Tenant: "acme"
```

## Properties

| Property   | Type             | Required | Description                                                                                              |
| ---------- | ---------------- | -------- | -------------------------------------------------------------------------------------------------------- |
| `url`      | string           | ✓        | A2A server endpoint URL (must include scheme).                                                           |
| `name`     | string           | ✗        | Tool name registered for the remote agent. Defaults to a name derived from the server's agent card.     |
| `headers`  | map\[string\]string | ✗     | Extra HTTP headers sent with every request (useful for `Authorization`, tenant selection, tracing, \u2026). |

> [!TIP]
> **See also**
>
> For full details on the A2A protocol and serving agents as A2A endpoints, see [A2A Protocol](../../features/a2a/index.md).
