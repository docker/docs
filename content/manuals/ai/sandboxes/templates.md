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

## Building templates with Dockerfiles

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

Official templates include the agent binary, Ubuntu base, and development tools
like Git, Docker CLI, Node.js, Python, and Go. They run as the non-root
`agent` user with sudo access.

### The USER pattern

Switch to `root` for system-level installations, then back to `agent` at the
end. The base template defaults to `USER agent`, so you need to explicitly
switch to root for package installations. Always switch back to `agent` before
the end of your Dockerfile so the agent runs with the correct permissions.

### Using templates built from Dockerfiles

Build your template:

```console
$ docker build -t my-template:v1 .
```

Use it directly from your local Docker daemon:

```console
$ docker sandbox run --load-local-template -t my-template:v1 claude ~/project
```

The `--load-local-template` flag tells the sandbox to use an image from your
local Docker daemon. Without it, the sandbox looks for the image in a registry.

To share the template with others, push it to a registry:

```console
$ docker tag my-template:v1 myorg/my-template:v1
$ docker push myorg/my-template:v1
$ docker sandbox run -t myorg/my-template:v1 claude ~/project
```

Once pushed to a registry, you don't need `--load-local-template`.

## Creating templates from existing sandboxes

Rather than writing a Dockerfile, you can start with a sandbox, configure it,
then save it as a template. This is convenient when you already have a working
environment set up by the agent.

Start a sandbox and have the agent install what you need:

```console
$ docker sandbox run claude ~/project
```

Inside the sandbox, ask the agent to install tools and configure the
environment. Once everything works, exit and save the sandbox as a template:

```console
$ docker sandbox save claude-sandbox-2026-02-02-123456 my-template:v1
✓ Saved sandbox as my-template:v1
```

This saves the image to your local Docker daemon. Use `--load-local-template`
to create new sandboxes from it:

```console
$ docker sandbox run --load-local-template -t my-template:v1 claude ~/other-project
```

To save as a tar file instead (for example, to transfer to another machine):

```console
$ docker sandbox save -o template.tar claude-sandbox-2026-02-02-123456 my-template:v1
```

Use a Dockerfile when you want a clear record of how the environment is built.
Use `docker sandbox save` when you already have a working sandbox and want to
reuse it.

## Example: Node.js template

This template adds Node.js 20 and common development tools:

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
    typescript@5.7.2 \
    eslint@9.17.0 \
    prettier@3.4.2

USER agent
```

Pin specific versions for reproducible builds across your team.

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

To share templates, push them to a registry with version tags:

```console
$ docker build -t myorg/sandbox-templates:python-v1.0 .
$ docker push myorg/sandbox-templates:python-v1.0
```

Or tag and push a saved sandbox:

```console
$ docker tag my-template:v1 myorg/my-template:v1.0
$ docker push myorg/my-template:v1.0
```

Team members use the template by referencing the registry image:

```console
$ docker sandbox run -t myorg/sandbox-templates:python-v1.0 claude ~/project
```

Use version tags like `:v1.0` instead of `:latest` for consistency across your
team.

## Next steps

- [Using sandboxes effectively](workflows.md)
- [Architecture](architecture.md)
- [Network policies](network-policies.md)
