---
title: Run Mistral Vibe in a Docker Sandbox
description: Package Mistral's Vibe coding agent as a Docker Sandbox kit so it runs in an isolated microVM and reaches the Mistral API through the sandbox proxy, keeping your API key off the VM.
summary: |
  Build a pinned, reusable image for Mistral's Vibe CLI and wire it into a
  Docker Sandbox agent kit. Vibe runs in an isolated microVM and reaches the
  Mistral API through the sandbox proxy, so your API key never enters the VM.
keywords: ai, mistral, vibe, docker sandboxes, sbx, coding agent, microvm, security, kits, uv
tags: [ai]
params:
  time: 25 minutes
---

Mistral Vibe is Mistral's open source coding agent. This guide shows how to
package it as a Docker Sandbox agent so it runs in an isolated microVM instead
of directly on your host. The agent reaches the Mistral API through the
sandbox proxy, so your API key stays on the host and never enters the VM.

Rather than install the agent fresh on every run, you'll bake a pinned image
once and reuse it. Pinning the agent version and building a dedicated image
gives you reproducible sandboxes and faster startups, and it's the approach
Docker recommends for [building your own agent](../manuals/ai/sandboxes/customize/build-an-agent.md).

In this guide, you'll learn how to:

- Store your Mistral API key on the host as a sandbox secret
- Build a pinned, multi-architecture image that ships Vibe on the `shell` template
- Work around a proxy parsing issue with a small startup wrapper
- Write an agent kit that wires Vibe to the Mistral API through the proxy
- Validate, launch, and iterate on the sandbox

## How isolation works

Every outbound request from a sandbox passes through a proxy that runs on
your host. The proxy enforces network policy and injects credentials, so the
agent inside the VM never handles the real key.

Mistral is a
[built-in service](../manuals/ai/sandboxes/security/credentials.md#built-in-services):
`sbx` already maps the `mistral` service name to the `MISTRAL_API_KEY`
environment variable and the `api.mistral.ai` domain. Inside the VM, Vibe
sees only a sentinel value for `MISTRAL_API_KEY`. The proxy swaps in the real
key, and only for requests to `api.mistral.ai`. If the agent reads the
variable for any other purpose, it gets the sentinel.

Built-in doesn't mean automatic. The kit you write in
[Step 5](#step-5-write-the-agent-kit) still declares where the key comes from
and how the proxy attaches it to requests.

## Prerequisites

Before you start, make sure you have:

- [Docker Desktop](../get-started/get-docker.md) or Docker Engine installed
- [Docker Sandboxes (`sbx`) installed and signed in](../manuals/ai/sandboxes/get-started.md#install-and-sign-in)
- A [Mistral API key](https://console.mistral.ai/)
- A Docker Hub namespace, or another registry, to publish the image to

## Step 1: Store the Mistral key on the host

Provide the key once on the host. Because Mistral is a built-in service,
`sbx` resolves it under the `mistral` name without any extra wiring:

```console
$ sbx secret set -g mistral
```

The `-g` flag stores the secret globally, so any sandbox that declares the
`mistral` service can use it. For how the proxy resolves and injects
credentials, see
[Credentials](../manuals/ai/sandboxes/security/credentials.md).

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

# The startup wrapper from Step 3.
USER root
COPY --chmod=0755 start.sh /usr/local/bin/vibe-sandbox

# Install Vibe as the non-root agent user. The socks extra is installed
# explicitly so the agent works through the sandbox proxy.
USER agent
RUN uv tool install "mistral-vibe==${VIBE_VERSION}" --with "httpx[socks]" \
    && vibe --version

CMD ["vibe-sandbox"]
```

Two choices are worth calling out:

- Pinning `mistral-vibe` to an explicit version keeps the image
  reproducible. To move to a newer release, change `VIBE_VERSION` and rebuild.
- Installing the `httpx[socks]` extra explicitly avoids a startup failure
  when the agent runs behind the sandbox proxy. Some Vibe releases don't pull
  it in on their own.

## Step 3: Add a startup wrapper for the proxy

Vibe uses the `httpx` HTTP library. The sandbox injects a `NO_PROXY`
variable that includes a bracketed IPv6 entry without a port, which older
`httpx` versions can't parse:

```text
NO_PROXY=localhost,127.0.0.1,::1,[::1],gateway.docker.internal
```

At startup, Vibe fails with:

```text
Background initialization failed: Invalid port: ':1]'
```

You can't rewrite an environment variable from `spec.yaml`, so add a wrapper
that strips the `[::1]` token before launching Vibe. Create `start.sh` next
to the `Dockerfile`:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Drop only the "[::1]" token, keep every other NO_PROXY entry.
_np="${NO_PROXY:-${no_proxy:-}}"
_np="$(printf '%s' "$_np" | sed 's/\[::1\]//g; s/,,*/,/g; s/^,//; s/,$//')"
export NO_PROXY="$_np" no_proxy="$_np"

exec vibe "$@"
```

The `exec` call replaces the wrapper with `vibe` under the same process ID,
and `"$@"` forwards any arguments `sbx` appends, such as in `--task` mode.
When `sbx` normalizes `NO_PROXY` on the runtime side, you can drop the
wrapper and launch `vibe` directly.

## Step 4: Build and publish the image

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

## Step 5: Write the agent kit

The kit ties the image, launch command, network policy, and credentials
together. Create a directory for the kit with a `spec.yaml` inside. Replace
`YOUR_NAMESPACE` with the namespace you published to:

```yaml {title="mistral-vibe/spec.yaml"}
schemaVersion: "1"
kind: agent
name: mistral-vibe
displayName: Mistral Vibe

agent:
  image: docker.io/YOUR_NAMESPACE/sbx-mistral-vibe:0.1.0
  aiFilename: AGENTS.md
  persistence: persistent
  entrypoint:
    run: [vibe-sandbox, "--agent", "auto-approve"]

network:
  serviceDomains:
    api.mistral.ai: mistral
  serviceAuth:
    mistral:
      headerName: Authorization
      valueFormat: "Bearer %s"

credentials:
  sources:
    mistral:
      env:
        - MISTRAL_API_KEY

environment:
  proxyManaged:
    - MISTRAL_API_KEY

memory: |
  You are running inside an isolated Docker Sandbox microVM.
  Network access is restricted to the Mistral API. Prefer tools and
  packages already available in the workspace.
```

Each field does the following:

| Field                      | Purpose                                                                                                   |
| -------------------------- | --------------------------------------------------------------------------------------------------------- |
| `kind: agent`              | Declares a standalone agent: a complete image plus its launch configuration.                              |
| `name`                     | The kit's identifier, reused in the `sbx run` command.                                                    |
| `agent.image`              | The pinned image you published in Step 4.                                                                 |
| `agent.aiFilename`         | The instructions file Vibe reads in the project.                                                          |
| `agent.persistence`        | `persistent` keeps agent state across restarts in a named volume.                                         |
| `agent.entrypoint.run`     | The command run at start. `--agent auto-approve` runs Vibe with automatic tool approvals.                 |
| `network.serviceDomains`   | Maps `api.mistral.ai` to the `mistral` service. Keep this narrow to avoid intercepting unrelated traffic. |
| `network.serviceAuth`      | The header format the proxy injects: `Authorization: Bearer <key>`.                                       |
| `credentials.sources`      | Where the key comes from. Without it, requests to `api.mistral.ai` return 401.                            |
| `environment.proxyManaged` | Places a sentinel `MISTRAL_API_KEY` in the VM so Vibe sees a non-empty key and skips its setup screen.    |
| `memory`                   | Markdown appended to `AGENTS.md` at creation to prime the agent about its environment.                    |

For the full kit format, see
[Kits](../manuals/ai/sandboxes/customize/kits.md).

> [!WARNING]
> `--agent auto-approve` runs Vibe in a mode that approves every tool
> execution without prompting. The sandbox isolates the agent from your host,
> but review the agent's actions before you run it against sensitive
> workspaces.

## Step 6: Validate and run

Validate the kit before you launch it:

```console
$ sbx kit validate ./mistral-vibe
```

Then, from your project directory, launch the agent with the kit:

```console
$ sbx run --kit ./mistral-vibe mistral-vibe .
```

- `--kit ./mistral-vibe` points to the folder that contains `spec.yaml`.
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
it. Use it to spot a blocked domain or a `serviceDomains` mapping that's too
broad. After you change `spec.yaml`, recreate the sandbox for a clean start:

```console
$ sbx rm mistral-vibe && sbx run --kit ./mistral-vibe mistral-vibe .
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
- [Credentials and built-in services](../manuals/ai/sandboxes/security/credentials.md#built-in-services)
- [Mistral Vibe](https://github.com/mistralai/mistral-vibe)
