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
$ docker sandbox run AGENT
```

Replace `AGENT` with your preferred agent (`claude`, `codex`, `copilot`, etc.).
The workspace defaults to your current directory when omitted. You can also
specify an explicit path:

```console
$ docker sandbox run AGENT ~/my-project
```

The `docker sandbox run` command is idempotent. Running the same command
multiple times reuses the existing sandbox instead of creating a new one:

```console
$ docker sandbox run AGENT ~/my-project  # Creates sandbox
$ docker sandbox run AGENT ~/my-project  # Reuses same sandbox
```

This works with workspace path (absolute or relative) or omitted workspace. The
sandbox persists. Stop and restart it without losing installed packages or
configuration:

```console
$ docker sandbox run <sandbox-name>  # Reconnect by name
```

When using the `--name` flag, the behavior is also idempotent based on the
name:

```console
$ docker sandbox run --name dev AGENT  # Creates sandbox named "dev"
$ docker sandbox run --name dev AGENT  # Reuses sandbox "dev"
```

## Installing dependencies

Ask the agent to install what's needed:

```plaintext
You: "Install pytest and black"
Agent: [Installs packages via pip]

You: "Install build-essential"
Agent: [Installs via apt]
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

Agent: *runs*
  docker build -t myapp:test .
  docker run myapp:test npm test
```

Containers started by the agent run inside the sandbox, not on your host. They
don't appear in your host's `docker ps`.

### Multi-container stacks

```plaintext
You: "Start the application with docker-compose and run integration tests"

Agent: *runs*
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

Agents running in sandboxes automatically trust the workspace directory without
prompting. This enables agents to work freely within the isolated environment.

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

## Managing multiple projects

Create sandboxes for different projects:

```console
$ docker sandbox create claude ~/project-a
$ docker sandbox create codex ~/project-b
$ docker sandbox create copilot ~/work/client-project
```

Each sandbox is completely isolated. Switch between them by running the
appropriate sandbox name.

Remove unused sandboxes to reclaim disk space:

```console
$ docker sandbox rm <sandbox-name>
```

## Named sandboxes

Docker automatically generates sandbox names based on the agent and workspace
directory (for example, `claude-my-project`). You can also specify custom names
using the `--name` flag:

```console
$ docker sandbox run --name myproject AGENT ~/project
```

Create multiple sandboxes for the same workspace:

```console
$ docker sandbox create --name dev claude ~/project
$ docker sandbox create --name staging codex ~/project
$ docker sandbox run dev
```

Each maintains separate packages, Docker images, and state, but share the
workspace files.

## Multiple workspaces

Mount multiple directories into a single sandbox for working with related
projects or when the agent needs access to documentation and shared libraries.

```console
$ docker sandbox run AGENT ~/my-project ~/shared-docs
```

The primary workspace (first argument) is always mounted read-write. Additional
workspaces are mounted read-write by default.

### Read-only mounts

Mount additional workspaces as read-only by appending `:ro` or `:readonly`:

```console
$ docker sandbox run AGENT . /path/to/docs:ro /path/to/lib:readonly
```

The primary workspace remains fully writable while read-only workspaces are
protected from changes.

### Path resolution

Workspaces are mounted at their absolute paths inside the sandbox. Relative
paths are resolved to absolute paths before mounting.

Example:

```console
$ cd /Users/bob/projects
$ docker sandbox run AGENT ./app ~/docs:ro
```

Inside the sandbox:

- `/Users/bob/projects/app` - Primary workspace (read-write)
- `/Users/bob/docs` - Additional workspace (read-only)

Changes to `/Users/bob/projects/app` sync back to your host, while
`/Users/bob/docs` remains read-only.

### Sharing workspaces across sandboxes

A single path can be included in multiple sandboxes simultaneously:

```console
$ docker sandbox create --name sb1 claude ./project-a
$ docker sandbox create --name sb2 claude ./project-a ./project-b
$ docker sandbox create --name sb3 cagent ./project-a
$ docker sandbox ls
SANDBOX   AGENT    STATUS    WORKSPACE
sb1       claude   running   /Users/bob/src/project-a
sb2       claude   running   /Users/bob/src/project-a, /Users/bob/src/project-b
sb3       cagent   running   /Users/bob/src/project-a
```

Each sandbox runs in isolation with separate configurations while sharing the
same workspace files.

## Resetting state

If you encounter issues with sandbox state, use the reset command to clean up
all VMs and registries:

```console
$ docker sandbox reset
```

This command:

- Stops all running sandbox VMs
- Deletes all VM state and registries
- Continues running the sandbox daemon (does not shut it down)
- Warns about directories it cannot remove

After reset, you can create fresh sandboxes. Use this when troubleshooting
persistent issues or reclaiming disk space from all sandboxes at once.

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

## Next steps

- [Custom templates](templates.md)
- [Architecture](architecture.md)
- [Network policies](network-policies.md)
