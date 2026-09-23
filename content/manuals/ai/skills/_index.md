---
title: Docker Skills
description: Use Docker's official skills to give compatible AI coding agents guidance for Docker tasks.
keywords: [docker skills, agent skills, ai coding agents, dockerfile, compose]
weight: 70
params:
  sidebar:
    group: AI and agents
---

Docker Skills are Docker's official, open-source guidance for compatible AI
coding agents. They help agents with Docker tasks, such as improving a Dockerfile
or setting up a Compose application. The collection evolves: browse the
[repository catalog](https://github.com/docker/skills#readme) for the current
list of skills. Docker product coverage is non-exhaustive and evolves with the
collection.

## Get started

1. [Install Docker Skills](install.md) with a method supported by your agent,
   and start a new agent session.
2. Ask for an outcome in plain language, for example, “Make this Dockerfile
   smaller and run as a non-root user.”
3. Review the proposed changes and commands, run your project's tests, and
   confirm destructive operations before approving them. Skills guide the
   agent; they don't replace validation.

## How agents use skills

A compatible agent matches your request to descriptions of skills it can
access, then loads the relevant instructions. You don't need to name a skill in
your prompt. A skill is a directory containing a `SKILL.md` file that describes
when and how to use it, with optional references or assets. An agent can use
more than one skill when a task spans Docker products.
