---
title: "Agent Distribution"
description: "Package, share, and run agents via OCI-compatible registries — just like container images."
keywords: docker agent, ai agents, concepts, agent distribution
weight: 50
canonical: https://docs.docker.com/ai/docker-agent/concepts/distribution/
aliases:
  - /ai/docker-agent/sharing-agents/
---

_Package, share, and run agents via OCI-compatible registries — just like container images._

## Overview

docker-agent agents can be pushed to any OCI-compatible registry (Docker Hub, GitHub Container Registry, etc.) and pulled/run anywhere. This makes sharing agents as easy as sharing Docker images.

> [!TIP]
> For CLI commands related to distribution, see [CLI Reference](../../features/cli/index.md) (`docker agent share push`, `docker agent share pull`, `docker agent alias`).

## Pushing Agents

```bash
# Push to Docker Hub
$ docker agent share push ./agent.yaml docker.io/username/my-agent:latest

# Push to GitHub Container Registry
$ docker agent share push ./agent.yaml ghcr.io/username/my-agent:v1.0
```

## Pulling Agents

```bash
# Pull an agent
$ docker agent share pull docker.io/username/my-agent:latest

# Pull from the agent catalog
$ docker agent share pull agentcatalog/pirate
```

## Running from a Registry

Run agents directly from a registry without pulling first:

```bash
# Run directly from Docker Hub
$ docker agent run docker.io/username/my-agent:latest

# Run from the agent catalog
$ docker agent run agentcatalog/pirate

# Run with a specific agent from a multi-agent config
$ docker agent run docker.io/username/dev-team:latest -a developer
```

## Agent Catalog

The `agentcatalog` namespace on Docker Hub hosts pre-built agents you can try:

```bash
# Try the pirate-themed assistant
$ docker agent run agentcatalog/pirate

# Try the coding agent
$ docker agent run agentcatalog/coder
```

## Using as Sub-Agents

Registry agents can be used directly as sub-agents in a multi-agent configuration — no need to define them locally:

```yaml
agents:
  root:
    model: openai/gpt-5
    description: Coordinator
    instruction: Delegate tasks to the right sub-agent.
    sub_agents:
      - agentcatalog/pirate         # auto-named "pirate"
      - my_reviewer:myorg/reviewer  # explicitly named "my_reviewer"
```

External sub-agents are automatically named after their last path segment. Use the `name:reference` syntax to give them a custom name.

Tag references are checked against the registry on every `docker agent run`, which adds a network round-trip per sub-agent at startup. Pin them to a digest (`agentcatalog/pirate@sha256:…`) to serve them from cache instead.

See [Pin external sub-agents to a digest](../multi-agent/index.md#pin-external-sub-agents-to-a-digest) and [External Sub-Agents](../multi-agent/index.md#external-sub-agents-from-registries) for details.

## Using with Aliases

Combine OCI references with aliases for convenient access:

```bash
# Create an alias for a registry agent
$ docker agent alias add coder agentcatalog/coder --yolo

# Now just run
$ docker agent run coder
```

## Using with API Server

The API server supports OCI references with auto-refresh:

```bash
# Start API from registry, auto-pull every 10 minutes
$ docker agent serve api docker.io/username/agent:latest --pull-interval 10
```

## Private Repositories

docker-agent supports pulling from private GitHub repositories and registries that require authentication. Use standard Docker login or GitHub authentication:

```bash
# Login to a registry
$ docker login docker.io

# Now push/pull works with private repos
$ docker agent share push ./agent.yaml docker.io/myorg/private-agent:latest
$ docker agent run docker.io/myorg/private-agent:latest
```

> [!NOTE]
> **Docker Desktop credentials**
>
> When pulling or running an agent from a `docker.com` or `*.docker.com` HTTPS URL (e.g. `desktop.docker.com`), docker-agent automatically forwards your Docker Desktop JWT for authentication — no explicit login required when Docker Desktop is running and signed in.
>
> Note: `docker.io` (the standard Docker Hub registry domain) is a separate domain and is **not** covered by automatic JWT forwarding. Agents pulled from `docker.io` or `registry-1.docker.io` still require `docker login docker.io` for private repositories.

> [!NOTE]
> **Troubleshooting**
>
> Having issues with push/pull? See [Troubleshooting](../../community/troubleshooting/index.md) for common registry issues.

## Local Development

For local development and testing, you can run an agent directly from a local HTTP server without a registry:

```bash
# Serve an agent config locally
$ python3 -m http.server 8080

# Run it directly via HTTP
$ docker agent run http://localhost:8080/agent.yaml
$ docker agent run http://127.0.0.1:8080/agent.yaml
```

This is useful for iterating on agent configs served from a local dev server before pushing to a registry. Both `localhost` and `127.0.0.1` addresses are supported with plain `http://` URLs.
