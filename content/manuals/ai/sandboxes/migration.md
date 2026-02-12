---
title: Migrating from legacy sandboxes
description: Migrate from container-based sandboxes to microVM-based sandboxes
weight: 100
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

The most recent versions of Docker Desktop create microVM-based sandboxes,
replacing the container-based implementation released in earlier versions. This
guide helps you migrate from legacy sandboxes to the new architecture.

## What changed

Docker Sandboxes now run in lightweight microVMs instead of containers. Each
sandbox has a private Docker daemon, better isolation, and network filtering
policies.

> [!NOTE]
> If you need to use legacy container-based sandboxes, install
> [Docker Desktop 4.57](/desktop/release-notes/#4570).

After upgrading to Docker Desktop 4.58 or later:

- Old sandboxes don't appear in `docker sandbox ls`
- They still exist as regular Docker containers and volumes
- You can see them with `docker ps -a` and `docker volume ls`
- Old sandboxes won't work with the new CLI commands
- Running `docker sandbox run` creates a new microVM-based sandbox

## Migration options

Choose the approach that fits your situation:

### Option 1: Start fresh (recommended)

This is the simplest approach for experimental features. You'll recreate your
sandbox with the new architecture.

1. Note any important configuration or installed packages in your old sandbox.

2. Remove the old sandbox containers:

   ```console
   $ docker rm -f $(docker ps -q -a --filter="label=docker/sandbox=true")
   ```

3. Remove the credential volume:

   ```console
   $ docker volume rm docker-claude-sandbox-data
   ```

4. Create a new microVM sandbox:

   ```console
   $ docker sandbox create claude ~/project
   $ docker sandbox run <sandbox-name>
   ```

5. Reinstall dependencies. Ask the agent to install needed tools:

   ```plaintext
   You: "Install all the tools needed to build and test this project"
   Claude: [Installs tools]
   ```

What you lose:

- API keys (re-authenticate on first run, or set `ANTHROPIC_API_KEY`)
- Installed packages (reinstall via the agent)
- Custom configuration (reconfigure as needed)

What you gain:

- Better isolation (microVM versus container)
- Private Docker daemon for test containers
- Network filtering policies
- Improved security

### Option 2: Migrate configuration

If you have extensive customization, preserve your setup by creating a custom
template.

1. Inspect your old sandbox to see what's installed:

   ```console
   $ docker exec <old-sandbox-container> dpkg -l
   ```

2. Create a custom template with your tools:

   ```dockerfile
   FROM docker/sandbox-templates:claude-code

   USER root

   # Install your tools
   RUN apt-get update && apt-get install -y \
       build-essential \
       nodejs \
       npm

   # Install language-specific packages
   RUN npm install -g typescript eslint

   # Add any custom configuration
   ENV EDITOR=vim

   USER agent
   ```

3. Build your template:

   ```console
   $ docker build -t my-sandbox-template:v1 .
   ```

4. Create a new sandbox with your template:

   ```console
   $ docker sandbox create --template my-sandbox-template:v1 \
       --pull-template=never \
       claude ~/project
   ```

   > [!NOTE]
   > The `--pull-template` flag was introduced in Docker Desktop 4.61 (Sandbox
   > version 0.12). For Docker Desktop 4.58–4.60, substitute
   > `--pull-template=never` with `--load-local-template`.

5. Run the sandbox:

   ```console
   $ docker sandbox run <sandbox-name>
   ```

If you want to share this template with your team, push it to a registry. See
[Custom templates](templates.md) for details.

## Cleanup old resources

After migrating, clean up legacy containers and volumes:

Remove specific sandbox:

```console
$ docker rm -f <old-sandbox-container>
$ docker volume rm docker-claude-sandbox-data
```

Remove all stopped containers and unused volumes:

```console
$ docker container prune
$ docker volume prune
```

> [!WARNING]
> `docker volume prune` removes ALL unused volumes, not just sandbox volumes.
> Make sure you don't have other important unused volumes before running this
> command.

## Understanding the differences

### Architecture

Old (container-based):

- Sandboxes were Docker containers
- Appeared in `docker ps`
- Mounted host Docker socket for container access
- Stored credentials in Docker volume

New (microVM-based):

- Sandboxes are lightweight microVMs
- Use `docker sandbox ls` to see them
- Private Docker daemon inside VM
- Credentials via `ANTHROPIC_API_KEY` environment variable or interactive auth

### CLI changes

Old command structure:

```console
$ docker sandbox run ~/project
```

New command structure:

```console
$ docker sandbox run claude ~/project
```

The agent name (`claude`, `codex`, `gemini`, `cagent`, `kiro`) is now a
required parameter when creating sandboxes, and you run the sandbox by name.
