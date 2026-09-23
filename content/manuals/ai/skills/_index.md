---
title: Docker Skills
description: Use Docker's official skills to give compatible AI coding agents guidance for Docker tasks.
keywords: [docker skills, agent skills, ai coding agents, dockerfile, compose]
weight: 70
params:
  sidebar:
    group: AI and agents
---

Docker Skills are Docker's official, open-source collection of guidance for
compatible AI coding agents. A skill is a directory with a `SKILL.md` file
describing when and how to use it, and optional references or assets. Explore
the source in the [Docker Skills repository](https://github.com/docker/skills).

> [!NOTE]
> Docker Skills are the content an agent uses. Docker Sandboxes
> [shared agent skills](/ai/sandboxes/workflows/agent-skills/) are a way to
> install and share skills in sandboxes. Docker Agent's
> [skills feature](/ai/docker-agent/features/skills/) lets an agent consume
> installed skills. Neither is a separate Docker Skills catalog.

## How agents use skills

A compatible agent matches your request to descriptions of the skills it can
access, then loads the relevant instructions. You don't need an entry-point
skill or to name a skill in your prompt. For work across products, an agent can
combine guidance from more than one skill when a task spans Docker products.

## What Docker publishes

Docker Skills can help with tasks involving, for example, Docker Sandboxes,
Docker Engine, and Docker Compose. These are illustrations, not a complete list
of supported products or skills: coverage changes as the collection evolves.
Browse the [repository catalog](https://github.com/docker/skills#readme) for
the current skills and their descriptions.

## Get started

1. [Install Docker Skills](install.md) through your agent's supported
   installation method. Select only the skills you need, and start a new agent
   session after installation.
2. Ask for the outcome in plain language, for example:
   - “Dockerize this application for local development.”
   - “Make this Dockerfile smaller and run as a non-root user.”
   - “Add a database health check to this Compose application.”
3. Review proposed changes and commands, run your project's tests, and confirm
   destructive operations before approving them. Skills guide the agent; they
   don't replace validation.
4. Keep skills current through the installer that owns them. Update pinned or
   manually copied skills deliberately from a reviewed release.

## Learn more

- [Install Docker Skills](install.md) — choose an installation method and verify
  that your agent can find the skills.
- [Docker Skills repository catalog](https://github.com/docker/skills#readme)
  — browse the authoritative list and descriptions.
- [Docker Skills releases](https://github.com/docker/skills/releases) — review
  source revisions before pinning an installation.
