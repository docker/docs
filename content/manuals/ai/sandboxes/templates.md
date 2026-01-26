---
title: Custom templates
description: Create custom sandbox templates to standardize development environments with pre-installed tools and configurations.
weight: 45
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

Custom templates let you create reusable sandbox environments with
pre-installed tools and configuration. Instead of asking the agent to install
packages each time, build a template with everything ready.

## When to use custom templates

Use custom templates when:

- Multiple team members need the same environment
- You're creating many sandboxes with identical tooling
- Setup involves complex steps that are tedious to repeat
- You need specific versions of tools or libraries

For one-off work or simple setups, use the default template and ask the agent
to install what's needed.

## Building a template

Start from Docker's official sandbox templates:

```dockerfile
FROM docker/sandbox-templates:claude-code

USER root

# Install system packages
RUN apt-get update && apt-get install -y \
    build-essential \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Install development tools
RUN pip3 install --no-cache-dir \
    pytest \
    black \
    pylint

USER agent
```

Official templates include the agent binary, Ubuntu base, development tools
(Git, Docker CLI, Node.js, Python, Go), and the non-root `agent` user with
sudo access.

### The USER pattern

Switch to `root` for system-level installations, then back to `agent` at the
end. The base template defaults to `USER agent`, so you need to explicitly
switch to root for package installations. Always switch back to `agent` before
the end of your Dockerfile so the agent runs with the correct permissions.

### Using templates

Build your template:

```console
$ docker build -t my-template:v1 .
```

Then choose how to use it:

Option 1: Load from local images (quick, for personal use)

```console
$ docker sandbox create --template my-template:v1 \
    --load-local-template \
    claude ~/project
$ docker sandbox run <sandbox-name>
```

The `--load-local-template` flag loads the image from your local Docker daemon
into the sandbox VM. This works for quick iteration and personal templates.

Option 2: Push to a registry (for sharing and persistence)

```console
$ docker tag my-template:v1 myorg/my-template:v1
$ docker push myorg/my-template:v1
$ docker sandbox create --template myorg/my-template:v1 claude ~/project
$ docker sandbox run <sandbox-name>
```

Pushing to a registry makes templates available to your team and persists them
beyond your local machine.

## Example: Node.js template

```dockerfile
FROM docker/sandbox-templates:claude-code

USER root

RUN apt-get update && apt-get install -y curl \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js 20 LTS
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Install common tools
RUN npm install -g \
    typescript@5.1.6 \
    eslint@8.46.0 \
    prettier@3.0.0

USER agent
```

Pin versions for reproducible builds.

## Using standard images

You can use standard Docker images (like `python:3.11` or `node:20`) as a
base, but they don't include agent binaries or sandbox configuration.

Using a standard image directly creates the sandbox but fails at runtime:

```console
$ docker sandbox create --template python:3-slim claude ~/project
✓ Created sandbox claude-sandbox-2026-01-16-170525 in VM claude-project

$ docker sandbox run claude-project
agent binary "claude" not found in sandbox: verify this is the correct sandbox type
```

To use a standard image, you'd need to install the agent binary, add sandbox
dependencies, configure the shell, and set up the `agent` user. Building from
`docker/sandbox-templates:claude-code` is simpler.

## Sharing with teams

Push templates to a registry and version them:

```console
$ docker build -t myorg/sandbox-templates:python-v1.0 .
$ docker push myorg/sandbox-templates:python-v1.0
```

Team members can then use the template:

```console
$ docker sandbox create --template myorg/sandbox-templates:python-v1.0 claude ~/project
```

Using version tags (`:v1.0`, `:v2.0`) instead of `:latest` ensures stability
across your team.

## Next steps

- [Using sandboxes effectively](workflows.md)
- [Architecture](architecture.md)
- [Network policies](network-policies.md)
