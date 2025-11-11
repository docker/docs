---
title: Advanced configurations
linkTitle: Advanced
description: Docker access, volume mounting, environment variables, custom templates, and sandbox management.
weight: 40
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

This guide covers advanced configurations for sandboxed agents running locally.

## Managing sandboxes

### Recreating sandboxes

Since Docker enforces one sandbox per workspace, the same sandbox is reused
each time you run `docker sandbox run <agent>` in a given directory. To create
a fresh sandbox, you need to remove the existing one first:

```console
$ docker sandbox ls  # Find the sandbox ID
$ docker sandbox rm <sandbox-id>
$ docker sandbox run <agent>  # Creates a new sandbox
```

### When to recreate sandboxes

Sandboxes remember their initial configuration and don't pick up changes from subsequent `docker sandbox run` commands. You must recreate the sandbox to modify:

- Environment variables (the `-e` flag)
- Volume mounts (the `-v` flag)
- Docker socket access (the `--mount-docker-socket` flag)
- Credentials mode (the `--credentials` flag)

### Listing and inspecting sandboxes

View all your sandboxes:

```console
$ docker sandbox ls
```

Get detailed information about a specific sandbox:

```console
$ docker sandbox inspect <sandbox-id>
```

This shows the sandbox's configuration, including environment variables, volumes, and creation time.

### Removing sandboxes

Remove a specific sandbox:

```console
$ docker sandbox rm <sandbox-id>
```

Remove all sandboxes at once:

```console
$ docker sandbox rm $(docker sandbox ls -q)
```

This is useful for cleanup when you're done with a project or want to start fresh.

## Giving agents access to Docker

Mount the Docker socket to give agents access to Docker commands inside the
container. The agent can build images, run containers, and work with Docker
Compose setups.

> [!CAUTION]
> Mounting the Docker socket grants the agent full access to your Docker daemon,
> which has root-level privileges on your system. The agent can start or stop
> any container, access volumes, and potentially escape the sandbox. Only use
> this option when you fully trust the code the agent is working with.

### Enable Docker socket access

Use the `--mount-docker-socket` flag:

```console
$ docker sandbox run --mount-docker-socket claude
```

This mounts your host's Docker socket (`/var/run/docker.sock`) into the
container, giving the agent access to Docker commands.

> [!IMPORTANT]
> The agent can see and interact with all containers on your host, not just
> those created within the sandbox.

### Example: Testing a containerized application

If your project has a Dockerfile, the agent can build and test it:

```console
$ cd ~/my-docker-app
$ docker sandbox run --mount-docker-socket claude
```

Example conversation:

```plaintext
You: "Build the Docker image and run the tests"

Claude: *runs*
  docker build -t myapp:test .
  docker run myapp:test npm test
```

### What agents can do with Docker socket access

With Docker access enabled, agents can:

- Start multi-container applications with Docker Compose
- Build images for multiple architectures
- Manage existing containers on your host
- Validate Dockerfiles and test build processes

## Environment variables

Pass environment variables to configure the sandbox environment with the `-e`
flag:

```console
$ docker sandbox run \
  -e NODE_ENV=development \
  -e DATABASE_URL=postgresql://localhost/myapp_dev \
  -e DEBUG=true \
  claude
```

These variables are available to all processes in the container, including the
agent and any commands it runs. Use multiple `-e` flags for multiple variables.

### Example: Development environment setup

Set up a complete development environment:

```console
$ docker sandbox run \
  -e NODE_ENV=development \
  -e DATABASE_URL=postgresql://localhost/myapp_dev \
  -e REDIS_URL=redis://localhost:6379 \
  -e LOG_LEVEL=debug \
  claude
```

Example conversation:

```plaintext
You: "Run the database migrations and start the development server"

Claude: *uses DATABASE_URL and other environment variables*
  npm run migrate
  npm run dev
```

### Common use cases

API keys for testing:

```console
$ docker sandbox run \
  -e STRIPE_TEST_KEY=sk_test_xxx \
  -e SENDGRID_API_KEY=SG.xxx \
  claude
```

> [!CAUTION]
> Only use test/development API keys in sandboxes, never production keys.

Loading from .env files:

Sandboxes don't automatically load `.env` files from your workspace, but you can ask Claude to use them:

```plaintext
You: "Load environment variables from .env.development and start the server"
```

Claude can use `dotenv` tools or source the file directly.

## Volume mounting

Mount additional directories or files to share data beyond your main workspace.
Use the `-v` flag with the syntax `host-path:container-path`:

```console
$ docker sandbox run -v ~/datasets:/data claude
```

This makes `~/datasets` available at `/data` inside the container. The agent
can read and write files in this location.

Read-only mounts:

Add `:ro` to prevent modifications:

```console
$ docker sandbox run -v ~/configs/app.yml:/config/app.yml:ro claude
```

Multiple mounts:

Use multiple `-v` flags to mount several locations:

```console
$ docker sandbox run \
  -v ~/datasets:/data:ro \
  -v ~/models:/models \
  -v ~/.cache/pip:/root/.cache/pip \
  claude
```

### Example: Machine learning workflow

Set up an ML environment with shared datasets, model storage, and persistent
caches:

```console
$ docker sandbox run \
  -v ~/datasets:/data:ro \
  -v ~/models:/models \
  -v ~/.cache/pip:/root/.cache/pip \
  claude
```

This provides read-only access to datasets (preventing accidental modifications),
read-write access to save trained models, and a persistent pip cache for faster
package installs across sessions.

Example conversation:

```plaintext
You: "Train a model on the MNIST dataset and save it to /models"

Claude: *runs*
  python train.py --data /data/mnist --output /models/mnist_model.h5
```

### Common use cases

Shared configuration files:

```console
$ docker sandbox run -v ~/.aws:/root/.aws:ro claude
```

Build caches:

```console
$ docker sandbox run \
  -v ~/.cache/go-build:/root/.cache/go-build \
  -v ~/go/pkg/mod:/go/pkg/mod \
  claude
```

Custom tools:

```console
$ docker sandbox run -v ~/bin:/shared-bin:ro claude
```

## Custom templates

Create custom sandbox templates to reuse configured environments. Instead of
installing tools every time you start an agent, build a Docker image with
everything pre-installed:

```dockerfile
# syntax=docker/dockerfile:1
FROM docker/sandbox-templates:claude-code
RUN <<EOF
curl -LsSf https://astral.sh/uv/install.sh | sh
. ~/.local/bin/env
uv tool install ruff@latest
EOF
ENV PATH="$PATH:~/.local/bin"
```

Build the image, and use the [`docker sandbox run --template`](/reference/cli/docker/sandbox/run#template)
flag to start a new sandbox based on the image.

```console
$ docker build -t my-dev-env .
$ docker sandbox run --template my-dev-env claude
```

### Using standard images

You can use standard Docker images as sandbox templates, but they don't include
agent binaries, shell configuration, or runtime dependencies that Docker's
sandbox templates provide. Using a standard Python image directly fails:

```console
$ docker sandbox run --template python:3-slim claude
The claude binary was not found in the sandbox; please check this is the correct sandbox for this agent.
```

To use a standard image, create a Dockerfile that installs the agent binary,
dependencies, and shell configuration on top of your base image. This approach
makes sense when you need a specific base image (for example, an exact OS
version or a specialized image with particular build tools).
