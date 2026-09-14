---
title: Run Mistral Vibe in a Docker Sandbox
description: Package Mistral's Vibe coding agent as a Docker Sandbox kit so it runs in an isolated microVM and reaches the Mistral API through the sandbox proxy, keeping your API key off the VM.
summary: |
  Build a pinned, reusable image for Mistral's Vibe CLI and wire it into a
  Docker Sandbox agent kit. Vibe runs in an isolated microVM and reaches the
  Mistral API through the sandbox proxy, so your API key never enters the VM.
keywords: ai, mistral, vibe, docker sandboxes, sbx, coding agent, microvm, security, kits, uv
params:
  tags: [ai]
  time: 20 minutes
---

Mistral Vibe is Mistral's open source coding agent. This guide shows how to
package it as a Docker Sandbox agent so it runs in an isolated microVM instead
of directly on your host. The agent reaches the Mistral API through the
sandbox proxy, so your API key stays on the host and never enters the VM.

Rather than install the agent fresh on every run, you'll bake a pinned image
once and reuse it. Pinning the agent version and building a dedicated image
gives you reproducible sandboxes and faster startups, and it's the approach
Docker recommends for building your own agent.

In this guide, you'll learn how to:

- Store your Mistral API key on the host as a sandbox secret
- Build a pinned, multi-architecture image that ships Vibe on the `shell` template
- Write an agent kit that wires Vibe to the Mistral API through the proxy
- Validate, launch, and iterate on the sandbox

## How isolation works

Every outbound request from a sandbox passes through a proxy that runs on
your host. The proxy enforces network policy and injects credentials, so the
agent inside the VM never handles the real key.

Mistral is a
[built-in service](../manuals/ai/sandboxes/configuration/credentials.md#built-in-services):
`sbx` already maps the `mistral` service name to the `MISTRAL_API_KEY`
environment variable and the `api.mistral.ai` domain. Inside the VM, Vibe
sees only a sentinel value for `MISTRAL_API_KEY`. The proxy swaps in the real
key, and only for requests to `api.mistral.ai`. If the agent reads the
variable for any other purpose, it gets the sentinel.

Built-in doesn't mean automatic. The kit you write in
[Step 4](#step-4-write-the-agent-kit) still declares where the key comes from
and how the proxy attaches it to requests.

## Prerequisites

Before you start, make sure you have:

- [Docker Desktop](../get-started/get-docker.md) or Docker Engine installed
- [Docker Sandboxes (`sbx`) installed and signed in](../manuals/ai/sandboxes/install.md)
- A [Mistral API key](https://console.mistral.ai/)
- A Docker Hub namespace, or another registry, to publish the image to

## Step 1: Store the Mistral key on the host

Provide the key once on the host. Because Mistral is a built-in service,
`sbx` resolves it under the `mistral` name without any extra wiring:

```console
$ sbx secret set mistral
```

Service secrets are global by default, so any sandbox that declares the
`mistral` service can use it. Use `--sandbox` to scope a secret to a single
sandbox instead. For how the proxy resolves and injects credentials, see
[Credentials](../manuals/ai/sandboxes/configuration/credentials.md).

## Step 2: Write a pinned Vibe image

Vibe is a Python application. The official `shell` template already ships
`uv`, `git`, `ripgrep`, and Python, so build on top of it and install a
pinned Vibe version with `uv`.

Create a `Dockerfile`:

```dockerfile
# syntax=docker/dockerfile:1
ARG BASE_IMAGE=docker/sandbox-templates:shell
FROM ${BASE_IMAGE}

# Pin the agent version for reproducible sandboxes.
# Check https://pypi.org/project/mistral-vibe/ and bump as needed.
ARG VIBE_VERSION=2.24.5

# Install Vibe as the non-root agent user. The socks extra is installed
# explicitly so the agent works through the sandbox proxy.
USER agent
RUN uv tool install "mistral-vibe==${VIBE_VERSION}" --with "httpx[socks]" \
    && vibe --version

CMD ["vibe", "--agent", "auto-approve"]
```

Three choices are worth calling out:

- Pinning `mistral-vibe` to an explicit version keeps the image
  reproducible. To move to a newer release, change `VIBE_VERSION` and rebuild.
- Installing the `httpx[socks]` extra explicitly avoids a startup failure
  when the agent runs behind the sandbox proxy. Some Vibe releases don't pull
  it in on their own.
- The `CMD` launches Vibe with `--agent auto-approve`, so the agent approves
  tool executions automatically. Baking the flags into the image keeps the
  kit's launch behavior in one place, so the kit doesn't override the
  entrypoint.

## Step 3: Build and publish the image

Build for both common architectures and publish to your namespace. The
build attaches provenance and SBOM attestations, which record how the image
was built and what it contains. Replace `YOUR_NAMESPACE` with your Docker Hub
username:

```console
$ docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --provenance=true \
  --sbom=true \
  -t YOUR_NAMESPACE/sbx-mistral-vibe:0.1.0 \
  --push .
```

Tag the image with a real version rather than a moving tag, so the kit in the
next step always resolves to the same build.

## Step 4: Write the agent kit

The kit ties the image, network policy, and credentials together. The image
already sets the launch command in its `CMD`. Create a directory for the kit
with a `spec.yaml` inside. Replace
`YOUR_NAMESPACE` with the namespace you published to:

```yaml {title="mistral-vibe/spec.yaml"}
schemaVersion: "2"
kind: sandbox
name: mistral-vibe
displayName: Mistral Vibe

sandbox:
  image: docker.io/YOUR_NAMESPACE/sbx-mistral-vibe:0.1.0

agentInstructions:
  filename: AGENTS.md
  content: |
    You are running inside an isolated Docker Sandbox microVM.
    Network access is restricted to the Mistral API. Prefer tools and
    packages already available in the workspace.

permissions:
  network:
    allow:
      - "api.mistral.ai:443"

credentials:
  - service: mistral
    apiKey:
      name: MISTRAL_API_KEY
      inject:
        - domain: api.mistral.ai
          scheme: bearer
```

Each field does the following:

| Field                       | Purpose                                                                                                        |
| --------------------------- | -------------------------------------------------------------------------------------------------------------- |
| `kind: sandbox`             | Declares a sandbox agent: a complete image plus its launch configuration.                                      |
| `name`                      | The kit's identifier, reused in the `sbx run` command.                                                         |
| `sandbox.image`             | The pinned image you published in Step 3. Its `CMD` launches Vibe, so the kit doesn't set an entrypoint.       |
| `agentInstructions.filename`| The instructions file Vibe reads in the project.                                                               |
| `agentInstructions.content` | Markdown appended to `AGENTS.md` at creation to prime the agent about its environment.                         |
| `permissions.network.allow` | The hosts the sandbox may reach. Without it, requests are blocked by the default deny policy.                  |
| `credentials[].service`     | The built-in service that supplies the key. `mistral` maps to `MISTRAL_API_KEY` and `api.mistral.ai`.          |
| `credentials[].apiKey.name` | The environment variable the proxy manages. Vibe sees a sentinel value; the proxy swaps in the real key.       |
| `credentials[].apiKey.inject`| Where and how the proxy attaches the key. `scheme: bearer` sets `Authorization: Bearer <key>` for the domain. |

For the full kit format, see
[Kits](../manuals/ai/sandboxes/customize/kits.md).

> [!WARNING]
> `--agent auto-approve` runs Vibe in a mode that approves every tool
> execution without prompting. The sandbox isolates the agent from your host,
> but review the agent's actions before you run it against sensitive
> workspaces.

## Step 5: Validate and run

Validate the kit before you launch it:

```console
$ sbx kit validate ./mistral-vibe
```

Then, from your project directory, launch the agent with the kit:

```console
$ sbx run --kit ./mistral-vibe --name mistral-vibe mistral-vibe .
```

- `--kit ./mistral-vibe` points to the folder that contains `spec.yaml`.
- `--name mistral-vibe` names the sandbox. Without it, `sbx` derives a name
  from the agent and the working directory, and the commands below won't match.
- `mistral-vibe` is the agent name from `spec.yaml`.
- `.` is the project directory to mount in the sandbox.

Vibe starts in an isolated microVM, talks to the Mistral API through the
proxy, and the real key never touches the container.

## Iterate on the kit

If the agent can't reach a domain it needs, or a request behaves
unexpectedly, inspect what the proxy saw:

```console
$ sbx policy log
```

Each entry shows the request, the rule it matched, and how the proxy handled
it. Use it to spot a host missing from `permissions.network.allow`. After you
change `spec.yaml`, recreate the sandbox for a clean start:

```console
$ sbx rm mistral-vibe && sbx run --kit ./mistral-vibe --name mistral-vibe mistral-vibe .
```

## Clean up

Sandboxes persist after Vibe exits. To stop one without deleting it:

```console
$ sbx stop mistral-vibe
```

To remove the sandbox and everything inside it:

```console
$ sbx rm mistral-vibe
```

Files in your workspace are unaffected.

## Learn more

- [Get started with Docker Sandboxes](../manuals/ai/sandboxes/get-started.md)
- [Build your own agent kit](../manuals/ai/sandboxes/customize/build-an-agent.md)
- [Customize sandboxes with kits](../manuals/ai/sandboxes/customize/kits.md)
- [Credentials and built-in services](../manuals/ai/sandboxes/configuration/credentials.md#built-in-services)
- [Mistral Vibe](https://github.com/mistralai/mistral-vibe)
