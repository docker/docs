---
title: Build your own agent kit
linkTitle: Build an agent
description: Build a schema v3 agent workload from your own Linux base image, prepare its sandbox environment, and configure credentials and agent instructions.
keywords: sandboxes, sbx, kits, agent, tutorial, claude, workload, build
weight: 20
aliases:
  - /ai/sandboxes/customize/build-an-agent/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> V3 kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change. Share
> feedback in the [docker/sbx-releases](https://github.com/docker/sbx-releases)
> repository.

Build a sandbox environment from a Linux base image you choose, then package
it as a schema v3 workload kit. You'll prepare the operating system, install
Claude Code, and configure its model, API access, and agent instructions.
The same approach applies to other agents and organization-maintained images.

For a shorter example that extends an existing agent environment, see
[Build a workload](/manuals/ai/sandboxes/customize/author/_index.md#build-a-workload). This tutorial builds up the
whole environment and its descriptor step by step. If you're starting with
kits, read [Author kits](/manuals/ai/sandboxes/customize/author/_index.md) for the file layout and descriptor concepts. For field definitions, see the [Kit spec reference](/manuals/ai/sandboxes/customize/author/kit-reference.md).

This creates an independent v3 workload. It doesn't inherit the built-in
`claude` kit, which uses v2. Any mixins you add must also use v3. For
v2 customizations, see [Kits v2](/manuals/ai/sandboxes/customize/kits-v2/_index.md).

## Prepare the kit directory

You need `sbx` with schema v3 support, Docker with Buildx, and an Anthropic
API key. Create a directory beside the project you want the agent to work on:

```console
$ mkdir claude-team
```

The completed directory contains three files:

```text
claude-team/
├── claude-team.yaml
├── claude-team.dockerfile
└── context.md
```

The YAML descriptor declares the kit's requirements. Its companion Dockerfile
builds the agent and defines the launch command. The Markdown file contains
instructions the agent can read. Matching the YAML and Dockerfile stems lets
the kit frontend find the recipe.

## Use your own base image

A Docker sandbox template is optional. Your Dockerfile can prepare a Linux
base image with the tools, account, and certificate store that the sandbox
needs. This tutorial uses Red Hat Universal Base Image (UBI) 9 as a concrete
example. For another base, adapt the package installation and account creation
to its existing contents while preserving the
[base image requirements](/manuals/ai/sandboxes/customize/author/kit-reference.md#base-image-requirements).

### Install the system packages

Create `claude-team/claude-team.dockerfile` with the base image and packages:

```dockerfile {title="claude-team/claude-team.dockerfile"}
FROM registry.access.redhat.com/ubi9/ubi:9

USER root
RUN dnf install -y bash ca-certificates curl-minimal git shadow-utils tar gzip \
    && dnf clean all
```

Bash, Git, curl, and CA certificates provide the shell, source control, and
HTTPS tools used in the sandbox. The `shadow-utils` package provides account
creation commands. The archive tools support agent installation. Add any
compilers, libraries, or other tools your projects need here.

### Create the agent account

The sandbox needs a non-root `agent` user with UID 1000 and home directory
`/home/agent`. Append the following to the Dockerfile to create that account
and writable directories for its workspace, configuration, and state:

```dockerfile {title="Append to claude-team/claude-team.dockerfile"}
RUN groupadd --gid 1000 agent \
    && useradd --uid 1000 --gid 1000 --create-home --shell /bin/bash agent \
    && mkdir -p /home/agent/workspace /home/agent/.local/bin \
        /home/agent/.local/share /home/agent/.local/state \
        /home/agent/.config/claude-team /home/agent/.docker/sandbox/locks \
    && chown -R agent:agent /home/agent
```

Creating these directories as part of the build gives them the intended
ownership before runtime mounts are applied. This example runs the agent
without `sudo`. Install system packages in the Dockerfile. Reserve runtime
install hooks for configuration that depends on an individual sandbox.

### Prepare certificate trust

Docker Sandboxes adds its proxy CA at runtime so clients can make HTTPS
requests through the sandbox proxy. UBI stores its public CA roots at a
different path from the bundle used by `sbx`. Append the following to copy
those roots to that path and direct clients to the combined bundle:

```dockerfile {title="Append to claude-team/claude-team.dockerfile"}
RUN update-ca-trust \
    && mkdir -p /usr/local/share/ca-certificates /etc/ssl/certs \
    && cp /etc/pki/tls/certs/ca-bundle.crt /etc/ssl/certs/ca-certificates.crt

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt \
    CURL_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt \
    REQUESTS_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt \
    NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt
```

If you add corporate CA certificates during the build, add them to
`/etc/pki/ca-trust/source/anchors/` before running `update-ca-trust` and
copying the bundle. For another distribution, use its certificate management
command and source bundle path.

### Prepare the shell environment

Give the agent a writable file for environment exports that need to persist
across shell commands. Append the following to create the file, source it
from login and interactive shells, and set `BASH_ENV` for non-interactive Bash:

```dockerfile {title="Append to claude-team/claude-team.dockerfile"}
RUN touch /etc/sandbox-persistent.sh \
    && chown agent:agent /etc/sandbox-persistent.sh \
    && chmod 0644 /etc/sandbox-persistent.sh \
    && printf '%s\n' '. /etc/sandbox-persistent.sh' \
        > /etc/profile.d/sandbox-persistent.sh \
    && printf '%s\n' '. /etc/sandbox-persistent.sh' >> /home/agent/.bashrc

ENV HOME=/home/agent \
    PATH="/home/agent/.local/bin:${PATH}" \
    BASH_ENV=/etc/sandbox-persistent.sh
```

The agent instructions later in this tutorial explain how to use this file.
This completes the operating system preparation. The next step installs the
agent into that environment.

## Build the agent into the image

Append the following to the Dockerfile to install Claude Code as `agent` and
set its launch command:

```dockerfile {title="Append to claude-team/claude-team.dockerfile"}
USER agent
ENV CLAUDE_ENV_FILE=/etc/sandbox-persistent.sh \
    IS_SANDBOX=1
ARG CLAUDE_VERSION
RUN curl -fsSL https://claude.ai/install.sh -o /tmp/install-claude.sh \
    && bash /tmp/install-claude.sh "${CLAUDE_VERSION}" \
    && rm /tmp/install-claude.sh

WORKDIR /home/agent/workspace
ENTRYPOINT ["claude", "--settings", "/home/agent/.config/claude-team/settings.json"]
CMD []
```

Installing as `agent` puts Claude Code under `/home/agent/`, where the launch
user can access it. The `CLAUDE_VERSION` build argument receives its value
from the descriptor in the next step. `CLAUDE_ENV_FILE` points Claude Code
to the persistent environment file you prepared above.

The Dockerfile owns the image's `ENTRYPOINT`, `CMD`, environment, user, and
working directory. `CMD []` clears inherited arguments. This tutorial uses
Claude Code's `--settings` option to load a model setting chosen when creating
the sandbox. You'll declare that settings file in
[Write the model settings](#write-the-model-settings).

Installing Claude Code belongs in the build recipe because the binary is the
same in every sandbox using this kit. BuildKit can cache that work. The image
provides the software and launch command; the descriptor that follows declares
its runtime needs. Any mixins added to this workload must use v3 and provide
tools compatible with its operating system and architecture.

## Describe the workload and its inputs

Create `claude-team/claude-team.yaml`. Start by identifying the workload and
declaring its two inputs: the agent version to build and the model to use
when the sandbox runs.

```yaml {title="claude-team/claude-team.yaml"}
# syntax=docker/runtime-kit:3
schemaVersion: "3"
kind: workload
displayName: Team Claude Code
description: Claude Code with team defaults and API-key authentication

args:
  version:
    default: "2.1.259"
    pattern: '^[0-9]+\.[0-9]+\.[0-9]+$'
    buildArg: CLAUDE_VERSION
  model:
    default: sonnet
    enum: [sonnet, opus, haiku]

provides: ["claude@${{ kit.args.version }}"]
```

The `version` input maps to `CLAUDE_VERSION` in the Dockerfile. The build
validates the value and substitutes it into `provides`, which tells other
kits which agent version this workload supplies.

The `model` input will be used in the settings file. It is chosen when creating
a sandbox and doesn't change the installed agent binary. For local sources,
changing this argument can still invoke Buildx, which can reuse cached image
layers.

## Allow access to the API

Claude Code needs to reach the Anthropic API. Add a `capabilities` list at
the top level of `claude-team.yaml`, after `provides`. Its first entry declares
that network access:

```yaml {title="Add to claude-team/claude-team.yaml"}
capabilities:
  - type: com.docker.runtime/network-policy@1
    config:
      runtime:
        allow:
          - api.anthropic.com:443
```

This rule applies when the sandbox runs. Downloading Claude Code in the
Dockerfile happens during the build and uses the builder's network.

## Declare the credential

The network rule permits a connection. A credential capability tells Docker
Sandboxes how to authenticate requests on that connection using a key stored
on the host.

Append this entry to the same `capabilities` list:

```yaml {title="Append under capabilities"}
  - type: com.docker.runtime/credential@1
    description: Anthropic API access
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

The sandbox receives a placeholder in `ANTHROPIC_API_KEY`. When Claude Code
makes a request to `api.anthropic.com`, the host proxy inserts the real API
key into the `x-api-key` header. The key stays on the host. You will supply
its value and approve its use when launching the kit.

## Write the model settings

The Dockerfile's launch command reads
`/home/agent/.config/claude-team/settings.json`. Use the lifecycle capability
to create that file with the model chosen for the sandbox.

Append this entry to `capabilities`:

```yaml {title="Append under capabilities"}
  - type: com.docker.runtime/lifecycle@1
    config:
      files:
        - path: /home/agent/.config/claude-team/settings.json
          content: |
            {"model": "${{ kit.args.model }}"}
          mode: "0644"
```

Docker Sandboxes substitutes the `model` argument and writes the file before
Claude Code starts. This work belongs to sandbox creation because the value
can differ between sandboxes using the same image.

## Add agent instructions

The kit can also give Claude Code instructions about the environment. Save the
following Markdown alongside the descriptor and Dockerfile:

```markdown {title="claude-team/context.md"}
## Team workflow

Read the project's README before changing code. Run the project's checks
before reporting a task complete, and report any checks you couldn't run.

Claude Code is installed in this sandbox. Its additional settings are at
`/home/agent/.config/claude-team/settings.json`.

Use `/etc/sandbox-persistent.sh` for environment exports needed by later
Bash commands. Keep shell completion scripts out of that file because
non-interactive commands also source it.
```

Append an agent-context entry to `capabilities` to include these instructions:

```yaml {title="Append under capabilities"}
  - type: com.docker.runtime/agent-context@1
    config:
      filename: CLAUDE.md
      contentFile: ./context.md
```

The build includes `context.md` in the image. At runtime, `sbx` generates a
`CLAUDE.md` profile in the parent directory of the mounted workspace inside
the sandbox. The profile points to the packaged `context.md` file for the
agent to read. It doesn't replace a `CLAUDE.md` in your project. See
[Agent instructions](/manuals/ai/sandboxes/customize/author/_index.md#give-the-agent-instructions) for how workload
and mixin instructions contribute to the profile.

## Store the key and run

Store your Anthropic API key on the host:

```console
$ sbx secret set anthropic
```

From the directory containing `claude-team`, launch the kit:

```console
$ sbx run --name claude-team ./claude-team
```

`sbx` builds the kit, mounts your current directory as the workspace, and
launches Claude Code. To work on another project, append its path to the command.
Approve the kit's credential request to connect the stored key to this kit,
then follow Claude Code's first-run prompts. Without an approved credential
binding, storing a key alone doesn't authenticate the agent. See
[Credential bindings](/manuals/ai/sandboxes/configuration/credentials.md#credential-bindings).

Choose a different model when creating a sandbox:

```console
$ sbx run --name claude-team-opus ./claude-team \
    --kit-arg claude-team.model=opus
```

The argument prefix is the local kit directory's name. The model value is
validated against the descriptor's `enum` and written to the settings file
before Claude Code starts.

## Iterate and publish

Edit the descriptor, Dockerfile, or context file and create another sandbox
with a different name to test the changes:

```console
$ sbx run --name claude-team-test-2 ./claude-team
```

Running an existing sandbox keeps its recorded configuration. During sandbox
creation, `sbx` reuses a local build when both the source content and supplied
kit arguments are unchanged. Changing either can invoke Buildx, including
changes to create-time arguments such as `model`. BuildKit can still reuse
unchanged image layers. To change the installed agent version for local runs,
update `args.version.default` in the descriptor.

When the kit is ready to share, sign in to Docker Hub, then build and push it
with Docker Buildx. Replace `<NAMESPACE>` with a Docker Hub namespace you can
push to:

```console
$ docker login
$ docker buildx build ./claude-team \
    --file ./claude-team/claude-team.yaml \
    --tag docker.io/<NAMESPACE>/claude-team:1.0.0 \
    --push
```

Buildx reads the descriptor as the build file. Its syntax directive selects
the kit frontend, which builds the companion Dockerfile and publishes the
declarations with the image. To override the binary version for a published
build, add `--build-arg version=<CLAUDE_VERSION>`. Use the kit argument name
`version` in this flag, rather than the Dockerfile's `CLAUDE_VERSION` name.

Run the published kit by its image reference:

```console
$ sbx run --name claude-team-shared docker.io/<NAMESPACE>/claude-team:1.0.0
```

For multi-platform images and distribution details, see
[Build and distribute kits](/manuals/ai/sandboxes/customize/author/distribute.md).
To add tools or shared configuration to this workload,
see [Kit examples](/manuals/ai/sandboxes/customize/author/kit-examples.md).
