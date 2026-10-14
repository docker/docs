---
title: Build a tool mixin
description: Package Claude Code as a reusable mixin with API access and credentials, test it in a shell workload, and publish it.
keywords: sandboxes, sbx, kits, mixins, tools, credentials, build
weight: 30
aliases:
  - /ai/sandboxes/customize/kit-examples/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

A tool mixin adds an executable and the settings it needs to an existing
workload. This tutorial packages Claude Code as a mixin, adds it to a shell
workload, and tests a request to the Anthropic API. The sandbox starts a shell,
and you choose when to run Claude Code.

The example adapts the upstream
[Claude Code mixin](https://github.com/docker/sandbox-kit-spec/tree/main/examples/claude-mixin)
to use API-key authentication. The upstream example also includes OAuth,
session storage, and MCP configuration.

You need `sbx`, Docker with Buildx, and an Anthropic API key. To publish the
mixin, you also need a registry namespace you can push to. Use a v3 workload
with this mixin. Built-in shortcuts such as `shell` and `claude` use v2. See
[Version compatibility](/manuals/ai/sandboxes/customize/_index.md#version-compatibility).

## Package the executable

Create a `claude-mixin` directory with a Dockerfile and YAML descriptor:

```text
claude-mixin/
├── claude-mixin.dockerfile
└── claude-mixin.yaml
```

The Dockerfile downloads the native Claude Code executable in a build stage,
then copies it into an empty image:

```dockerfile {title="claude-mixin/claude-mixin.dockerfile"}
FROM debian:trixie-slim AS build
ARG CLAUDE_VERSION
ARG TARGETARCH
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates
RUN case "$TARGETARCH" in \
      amd64) platform=linux-x64 ;; \
      arm64) platform=linux-arm64 ;; \
      *) echo "unsupported TARGETARCH: $TARGETARCH" >&2; exit 1 ;; \
    esac \
    && mkdir -p /out/usr/local/bin \
    && curl -fsSL "https://downloads.claude.ai/claude-code-releases/${CLAUDE_VERSION}/${platform}/claude" \
        -o /out/usr/local/bin/claude \
    && chmod 0755 /out/usr/local/bin/claude

FROM scratch
COPY --from=build /out/ /
```

BuildKit supplies `TARGETARCH` to select the download for the image's Linux
architecture. You'll set `CLAUDE_VERSION` in the descriptor. The final stage
contains only the executable, so the mixin doesn't add Debian or the download
tools to the workload.

Docker Sandboxes adds the mixin's files to the workload's environment. The
workload keeps its own startup command, user, and working directory. This
example uses Docker's shell workload. When choosing another workload, check
that it supplies the Linux libraries and shell your tool needs.

A mixin includes only the files its Dockerfile adds or changes. If you use
a larger base image instead of `scratch`, the mixin won't include that
image's unchanged files.

## Declare runtime access

Create the descriptor with the version to install and the API access Claude
Code needs:

```yaml {title="claude-mixin/claude-mixin.yaml"}
# syntax=docker/sandbox-kit:3
schemaVersion: "3"
kind: mixin
displayName: Claude Code mixin

args:
  version:
    default: "2.1.278"
    pattern: '^[0-9]+\.[0-9]+\.[0-9]+$'
    buildArg: CLAUDE_VERSION

provides: ["claude@${{ kit.args.version }}"]

capabilities:
  - type: com.docker.sandbox/network-policy@1
    config:
      runtime:
        allow:
          - api.anthropic.com:443
  - type: com.docker.sandbox/credential@1
    config:
      service: anthropic
      phase: runtime
      apiKey:
        name: ANTHROPIC_API_KEY
        proxyManaged: true
        inject:
          - domain: api.anthropic.com
            header: x-api-key
            format: "%s"
```

The `version` argument supplies `CLAUDE_VERSION` to the Dockerfile. The
`provides` entry identifies the installed Claude Code version so other kits
can check their requirements.

The network rule permits HTTPS requests to the Anthropic API. The credential
entry requests the key stored under `anthropic` on your host. Inside the
sandbox, `ANTHROPIC_API_KEY` contains a placeholder. The host proxy inserts
the real key into the `x-api-key` header when Claude Code calls the API.

Network access and credential access are separate: declaring a credential
doesn't allow connections to its service. This example pairs `phase: runtime`
with `runtime.allow`. The download in the Dockerfile uses the builder's
network and doesn't need a rule in the descriptor.

Keep the tool's access requirements in the mixin so they follow it when you
use it with another workload. Access shared by a project's tools, such as a
package registry, can go in a
[kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md).

## Try the mixin

Store your Anthropic API key on the host:

```console
$ sbx secret set anthropic
```

From the directory containing `claude-mixin`, add the mixin to Docker's
published shell workload:

```console
$ sbx run docker.io/docker/sbx-kit-shell:1.0.0 --name claude-mixin-test \
    --kit ./claude-mixin
```

Approve the mixin's credential request when prompted. `sbx` builds the mixin,
adds its executable to the workload, and opens a shell. The mixin doesn't
replace the workload's launch command.

Inside the sandbox, check that Claude Code runs:

```console
$ claude --version
```

Then send a request to check network access and authentication. This command
uses the Anthropic API and incurs usage charges:

```console
$ claude -p "Reply with the word hello."
```

A response confirms that Claude Code can reach the API and authenticate with
your stored key. If authentication fails, check that you stored the key and
approved the kit's request to use it. See
[Credential bindings](/manuals/ai/sandboxes/configuration/credentials.md#credential-bindings).
Network requests must also meet the sandbox's
[network policy](/manuals/ai/sandboxes/governance/concepts.md#precedence).

Run `claude` without arguments to start an interactive session. Follow its
first-run prompts. Exit Claude Code and the shell to return to your host.

After editing the kit's files, create a sandbox with a different name to test
your changes. Reopening an existing sandbox uses the kits it was created with.

## Publish the mixin

Build and publish the mixin so others can use it without building from source:

```console
$ docker login
$ docker buildx build ./claude-mixin -f ./claude-mixin/claude-mixin.yaml \
    -t docker.io/<NAMESPACE>/claude-mixin:1.0.0 --push
```

Replace `<NAMESPACE>` with a Docker Hub namespace you can push to. Use the
published reference with `--kit` in place of `./claude-mixin`.

Choose a workload that doesn't already provide Claude Code. Two kits that
declare the same feature in `provides` conflict. To share the shell workload
and this mixin as one kit, see
[Compose a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md).

## More source examples

The [sandbox-kit-spec examples](https://github.com/docker/sandbox-kit-spec/tree/main/examples)
contain complete kits you can study and adapt:

- [Claude Code mixin](https://github.com/docker/sandbox-kit-spec/tree/main/examples/claude-mixin)
  includes additional authentication options, session storage, and agent
  instructions.
- [GitHub CLI](https://github.com/docker/sandbox-kit-spec/tree/main/examples/gh)
  packages `gh` and its dependencies with network rules, credentials, and
  agent instructions.
- [Message of the day](https://github.com/docker/sandbox-kit-spec/tree/main/examples/motd)
  uses an inline Dockerfile and a build argument in a single-file mixin.

For examples of building a tool from source, running setup scripts, and other
reusable approaches, see
[Kit authoring patterns](/manuals/ai/sandboxes/customize/author/patterns.md).
