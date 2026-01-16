---
title: Docker Sandboxes
description: Run AI agents in isolated environments
weight: 20
params:
  sidebar:
    group: AI
    badge:
      color: violet
      text: Experimental
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

Docker Sandboxes simplifies running AI agents securely on your local machine.
Designed for developers building with coding agents like Claude Code, Sandboxes
isolate your agents from your local machine while preserving a familiar
development experience. With Docker Sandboxes, agents can execute commands,
install packages, and modify files inside a containerized workspace that
mirrors your local directory. This gives you full agent autonomy without
compromising safety.

## How it works

When you run `docker sandbox run <agent>`:

1. Docker creates a container from a template image and mounts your current
   working directory into the container at the same path.

2. Docker discovers your Git `user.name` and `user.email` configuration and
   injects it into the container so commits made by the agent are attributed
   to you.

3. On first run, you're prompted to authenticate. Credentials are stored in a
   Docker volume and reused for future sandboxed agents.

4. The agent starts inside the container with bypass permissions enabled.

### Workspace mounting

Your workspace directory is mounted into the container at the same absolute path
(on macOS and Linux). For example, `/Users/alice/projects/myapp` on your host
is also `/Users/alice/projects/myapp` in the container. This means:

- File paths in error messages match your host
- Scripts with hard-coded paths work as expected
- Changes to workspace files are immediately visible on both host and container

### One sandbox per workspace

Docker enforces one sandbox per workspace. When you run `docker sandbox run
<agent>` in the same directory, Docker reuses the existing container. This
means state (installed packages, temporary files) persists across agent sessions
in that workspace.

> [!NOTE]
> To change a sandbox's configuration (environment variables, mounted volumes,
> etc.), you need to remove and recreate it. See
> [Managing sandboxes](advanced-config.md#managing-sandboxes) for details.

## Release status

Docker Sandboxes is an experimental feature. Features and setup are subject to
change.

Report issues on the [Docker Desktop issue tracker](https://github.com/docker/desktop-feedback).

## Get started

Head to the [Get started guide](get-started.md) to run your first sandboxed agent.
