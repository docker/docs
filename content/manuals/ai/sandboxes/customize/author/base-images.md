---
title: Base images for sandbox workloads
linkTitle: Base images
weight: 10
description: Choose a Docker-provided agent image or your own Linux base image for a sandbox workload kit, and use image overrides with built-in agents.
keywords: sandboxes, sbx, kits, base images, templates, dockerfile, custom agents
aliases:
  - /ai/sandboxes/customize/templates/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

A workload kit's Dockerfile defines the sandbox's filesystem. Start from a
Docker-provided agent image to reuse its tools and setup, or choose another
Linux image to control the operating system and prepare the environment yourself.
The kit descriptor declares the workload's runtime requirements alongside that
build recipe.

For a v3 workload, use the image in its Dockerfile. For a built-in agent such
as `claude`, use an [image override](#image-overrides-for-built-in-agents).
A kit set gets its base environment from its workload component.

To save and reuse an environment you've configured interactively, see
[Save a sandbox as a template](/manuals/ai/sandboxes/usage.md#saving-a-sandbox-as-a-template).

## Base images

All sandbox templates are published as
`docker/sandbox-templates:<variant>`. They are based on Ubuntu and run as a
non-root `agent` user with sudo access. Most variants include Git, Docker
CLI, and common development tools like Node.js, Python, Go, and Java.

| Variant               | Agent                                                                |
| --------------------- | -------------------------------------------------------------------- |
| `claude-code`         | [Claude Code](https://claude.ai/download)                            |
| `claude-code-minimal` | Claude Code with a minimal toolset (no Node.js, Python, Go, or Java) |
| `codex`               | [OpenAI Codex](https://github.com/openai/codex)                      |
| `copilot`             | [GitHub Copilot](https://github.com/github/copilot-cli)              |
| `cursor-agent`        | [Cursor](https://cursor.com/cli)                                     |
| `devin`               | [Devin CLI](https://docs.devin.ai/work-with-devin/devin-cli)         |
| `docker-agent`        | [Docker Agent](https://github.com/docker/docker-agent)               |
| `droid`               | [Droid](https://www.factory.ai)                                      |
| `gemini`              | [Gemini CLI](https://github.com/google-gemini/gemini-cli)            |
| `kiro`                | [Kiro](https://kiro.dev)                                             |
| `opencode`            | [OpenCode](https://opencode.ai)                                      |
| `shell`               | No agent pre-installed. Use for manual agent setup.                  |

## Use an image in a v3 workload

Reference the image in your workload's Dockerfile with `FROM`, then add the
tools and configuration you need. System package installations run as `root`;
switch back to `agent` before installing tools into the agent's home directory.
Otherwise, user-level installers put files under `/root/`, where the agent
can't use them.

The [workload example](/manuals/ai/sandboxes/customize/author/_index.md#build-a-workload) pairs the OpenCode image
with a kit descriptor and an explicit launch command. `sbx` builds the
image as part of the kit, so you don't need to build and distribute a separate
template first.

The base image supplies filesystem content and image settings. It doesn't
replace the kit's declarations for network access, credentials, storage, or
lifecycle hooks. Using a Docker-provided template image as a build base also
doesn't select a built-in v2 kit: the workload descriptor defines the kit
version. See [Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility).

## Use your own Linux image

You don't have to derive your workload from a Docker-provided image. Choose a
Linux base image that suits your operating system and package requirements,
and prepare it for the sandbox runtime.

See [Base image requirements](/manuals/ai/sandboxes/customize/author/kit-reference.md#base-image-requirements) for the
required tools, user account, and launch configuration. The
[Build an agent workload](/manuals/ai/sandboxes/customize/author/build-an-agent.md) tutorial walks through preparing a base
image, installing an agent, and declaring its runtime requirements.

## Image overrides for built-in agents

Built-in agent shortcuts such as `claude` and `codex` use v2 kits. Their
`--template` option selects a replacement image for that agent. Installing a
different agent binary in the image doesn't change which agent the built-in
kit launches. For v3 customization, define the image build and launch command
in a workload kit instead.

Each variant also has a `-docker` version (for example, `claude-code-docker`)
that includes a full Docker Engine running inside the sandbox — no local Docker
daemon required. When you pick a built-in agent without specifying a custom
template, `sbx run` and `sbx create` use the `-docker` template variants by
default.

The agent containers created from the `-docker` templates run in privileged
mode inside the microVM (not on your host), with a dedicated block volume at
`/var/lib/docker`, and `dockerd` starts automatically inside the sandbox. The
block volume defaults to 10 GB and uses a sparse file, so it only consumes
disk space as Docker writes to it.

To change the volume size for a sandbox, set
`DOCKER_SANDBOXES_DOCKER_SIZE` when you create it:

```console
$ DOCKER_SANDBOXES_DOCKER_SIZE=20g sbx run claude
```

The volume size must be at least 512 MiB. The environment variable doesn't
resize existing volumes.

Use the non-Docker variant if you don't need to build or run containers
inside the sandbox and want a lighter, non-privileged environment. Specify
it explicitly with `--template`:

```console
$ sbx run claude --template docker.io/docker/sandbox-templates:claude-code
```

### Build a custom template

Building a custom template requires
[Docker Desktop](/manuals/desktop/_index.md).

Write a Dockerfile that extends one of the base images. Pick the variant
that matches the agent you plan to run. For example, extend `claude-code`
to customize a Claude Code environment, or `codex` to customize an OpenAI
Codex environment.

The following example creates a Claude Code template with Rust and
protocol buffer tools pre-installed:

```dockerfile
FROM docker/sandbox-templates:claude-code
USER root
RUN apt-get update && apt-get install -y protobuf-compiler
USER agent
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
```

Build the image and push it to an OCI registry, such as Docker Hub:

```console
$ docker build -t my-org/my-template:v1 --push .
```

For registry credentials and loading a locally built image, see
[Load a template](/manuals/ai/sandboxes/usage.md#load-a-template).

Unless you use the permissive `allow-all` network policy, you may also need
to allow-list any domains that your custom tools depend on:

```console
$ sbx policy allow network "*.example.com:443,example.com:443"
```

Then run a sandbox with your template. The agent you specify must match
the base image variant your template extends:

```console
$ sbx run --template docker.io/my-org/my-template:v1 claude
```

Because this template extends the `claude-code` base image, you run it
with `claude`. If you extend `codex`, use `codex`; if you extend `shell`,
use `shell` (which drops you into a Bash shell with no agent).

> [!NOTE]
> Unlike Docker commands, `sbx` does not automatically resolve the Docker
> Hub domain (`docker.io`) in image references.
