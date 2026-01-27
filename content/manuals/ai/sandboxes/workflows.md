---
title: Using sandboxes effectively
linkTitle: Workflows
description: Best practices and common workflows for Docker Sandboxes including dependency management, testing, and multi-project setups.
weight: 35
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers practical patterns for working with sandboxed agents.

## Basic workflow

Create a sandbox for your project:

```console
$ cd ~/my-project
$ docker sandbox run claude .
```

The sandbox persists. Stop and restart it without losing installed packages or
configuration:

```console
$ docker sandbox run <sandbox-name>  # Reconnect later
```

## Installing dependencies

Ask the agent to install what's needed:

```plaintext
You: "Install pytest and black"
Claude: [Installs packages via pip]

You: "Install build-essential"
Claude: [Installs via apt]
```

The agent has sudo access. Installed packages persist for the sandbox lifetime.
This works for system packages, language packages, and development tools.

For teams or repeated setups, use [Custom templates](templates.md) to
pre-install tools.

## Docker inside sandboxes

Agents can build images, run containers, and use Docker Compose. Everything
runs inside the sandbox's private Docker daemon.

### Testing containerized apps

```plaintext
You: "Build the Docker image and run the tests"

Claude: *runs*
  docker build -t myapp:test .
  docker run myapp:test npm test
```

Containers started by the agent run inside the sandbox, not on your host. They
don't appear in your host's `docker ps`.

### Multi-container stacks

```plaintext
You: "Start the application with docker-compose and run integration tests"

Claude: *runs*
  docker-compose up -d
  docker-compose exec api pytest tests/integration
  docker-compose down
```

Remove the sandbox, and all images, containers, and volumes are deleted.

## What persists

While a sandbox exists:

- Installed packages (apt, pip, npm, etc.)
- Docker images and containers inside the sandbox
- Configuration changes
- Command history

When you remove a sandbox:

- Everything inside is deleted
- Your workspace files remain on your host (synced back)

To preserve a configured environment, create a [Custom template](templates.md).

## Security considerations

Agents can create and modify any files in your mounted workspace, including
scripts, configuration files, and hidden files.

After an agent works in a workspace, review changes before performing actions
on your host that might execute code:

- Committing changes (executes Git hooks)
- Opening the workspace in an IDE (may auto-run scripts or extensions)
- Running scripts or executables the agent created or modified

Review what changed:

```console
$ git status                        # See modified and new files
$ git diff                          # Review changes to tracked files
```

Check for untracked files and be aware that some changes, like Git hooks in
`.git/hooks/`, won't appear in standard diffs.

This is the same trust model used by editors like Visual Studio Code, which
warn when opening new workspaces for similar reasons.

## Named sandboxes

Use meaningful names for sandboxes you'll reuse:

```console
$ docker sandbox run --name myproject claude ~/project
```

Create multiple sandboxes for the same workspace:

```console
$ docker sandbox create --name dev claude ~/project
$ docker sandbox create --name staging claude ~/project
$ docker sandbox run dev
```

Each maintains separate packages, Docker images, and state, but share the
workspace files.

## Debugging

Access the sandbox directly with an interactive shell:

```console
$ docker sandbox exec -it <sandbox-name> bash
```

Inside the shell, you can inspect the environment, manually install packages,
or check Docker containers:

```console
agent@sandbox:~$ docker ps
agent@sandbox:~$ docker images
```

List all sandboxes:

```console
$ docker sandbox ls
```

## Managing multiple projects

Create sandboxes for different projects:

```console
$ docker sandbox create claude ~/project-a
$ docker sandbox create claude ~/project-b
$ docker sandbox create claude ~/work/client-project
```

Each sandbox is completely isolated. Switch between them by running the
appropriate sandbox name.

Remove unused sandboxes to reclaim disk space:

```console
$ docker sandbox rm <sandbox-name>
```

## Next steps

- [Custom templates](templates.md)
- [Architecture](architecture.md)
- [Network policies](network-policies.md)
