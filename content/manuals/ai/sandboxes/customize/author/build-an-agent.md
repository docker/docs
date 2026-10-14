---
title: Build an agent workload
linkTitle: Build an agent workload
description: Build a schema v3 agent workload from your own Linux base image, prepare its sandbox environment, and configure credentials and agent instructions.
keywords: sandboxes, sbx, kits, agent, tutorial, claude, workload, build
weight: 20
aliases:
  - /ai/sandboxes/customize/build-an-agent/
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Build a kit that runs Claude Code on a Linux base image of your choice.
You'll install the agent, give it access to the Anthropic API, and add model
settings and instructions for your team. The result is a v3 workload kit you
can run locally or publish for others to use. You can follow the same steps
for other agents or your organization's own base images.

This tutorial prepares the whole environment, from system packages to the
command that starts the agent. If you already have a suitable agent image,
see [Package an existing agent image](/manuals/ai/sandboxes/customize/author/base-images.md#package-an-existing-agent-image).
To combine a published agent kit with tools, [author a kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md).

If you're starting with kits, read [Author kits](/manuals/ai/sandboxes/customize/author/_index.md)
for an introduction to the files you'll create. Field definitions are in the
[upstream v3 specification](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md).

## Prepare the kit directory

You need `sbx`, Docker with Buildx, and an Anthropic
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

The Dockerfile installs the agent and sets the command to run it. The YAML
file, called the descriptor, tells Docker Sandboxes what the agent needs to
run, such as network access and credentials. The Markdown file contains
instructions for the agent. Give the YAML file and Dockerfile the same name
before the extension so the build can find both files.

This kit uses v3 and is independent of the built-in `claude` kit, which uses
v2. Any mixins you add must also use v3. To customize the built-in kit, see
[Kits v2](/manuals/ai/sandboxes/customize/kits-v2.md).

## Use your own base image

This tutorial starts from Red Hat Universal Base Image (UBI) 9. The following
steps add the tools, user account, and certificates the sandbox needs. You can
use another Linux image, including one maintained by your organization.
Adjust the package and account commands for that image, following the same
[base image requirements](/manuals/ai/sandboxes/customize/author/base-images.md#base-image-requirements).

### Install the system packages

Create `claude-team/claude-team.dockerfile` with the base image and packages:

```dockerfile {title="claude-team/claude-team.dockerfile"}
FROM registry.access.redhat.com/ubi9/ubi:9.8

USER root
RUN dnf install -y bash ca-certificates curl-minimal git shadow-utils tar gzip \
    && dnf clean all
```

Bash runs shell commands, Git accesses source repositories, and curl uses
the CA certificates to make HTTPS requests. The remaining packages create
the agent's user account and unpack its installer. Add any compilers,
libraries, or other tools your projects need here.

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

Creating these directories during the build makes the agent their owner
before Docker Sandboxes mounts the workspace and storage. The agent runs
without `sudo` in this example, so install system packages in the Dockerfile.
Use install hooks for setup that depends on an individual sandbox.

### Prepare certificate trust

HTTPS requests from the sandbox go through a proxy. Docker Sandboxes adds
the proxy's certificate authority (CA) to the sandbox's trusted certificates
when it starts. UBI keeps its certificates at a different path from the one
`sbx` uses. Copy them to the expected path, then tell tools to use that file
so they trust both public certificates and the proxy:

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

An environment variable exported in one shell isn't automatically available
in another. Give the agent a file where it can save variables for later
shell sessions. The following lines create that file and arrange for login,
interactive, and non-interactive Bash shells to read it:

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

`BASH_ENV` tells non-interactive Bash to read the file. You'll include
instructions for the agent to use it later in the tutorial.

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

Installing as `agent` puts Claude Code under `/home/agent/`, where that user
can access it. `CLAUDE_ENV_FILE` points Claude Code to the file you prepared
for persistent environment variables. You'll set `CLAUDE_VERSION` in the
descriptor in the next step.

The Dockerfile also sets the user, working directory, environment variables,
and launch command. `CMD []` clears any arguments inherited from the base
image. The `--settings` option loads the model you choose for each sandbox.
You'll create its settings file in [Write the model settings](#write-the-model-settings).

Each sandbox created from this image gets the same Claude Code binary.
The next steps configure how Docker Sandboxes runs it. Any mixins you add
must use v3 and contain tools compatible with this workload's operating
system and architecture.

## Describe the workload and its inputs

Create `claude-team/claude-team.yaml`. Start by identifying the workload and
declaring two arguments:

- `version` selects the Claude Code version to install during the image build.
- `model` selects the model to use when creating a sandbox.

```yaml {title="claude-team/claude-team.yaml"}
# syntax=docker/sandbox-kit:3
schemaVersion: "3"
kind: workload
displayName: Team Claude Code
description: Claude Code with team defaults and API-key authentication

args:
  version:
    default: "2.1.278"
    pattern: '^[0-9]+\.[0-9]+\.[0-9]+$'
    buildArg: CLAUDE_VERSION
  model:
    default: sonnet
    enum: [sonnet, opus, haiku]

provides: ["claude@${{ kit.args.version }}"]
```

`version` sets the Dockerfile's `CLAUDE_VERSION` build argument. The build
checks the version against `pattern`, then includes it in `provides` so
other kits can check which Claude Code version is installed.

`model` lets you choose a model when creating each sandbox. The choice goes
into a settings file, so it doesn't change the installed agent.

## Allow access to the API

With the agent installed, the next step is to give it access to the Anthropic
API. Add a `capabilities` list at the top level of `claude-team.yaml`, after
`provides`:

```yaml {title="Add to claude-team/claude-team.yaml"}
capabilities:
  - type: com.docker.sandbox/sbx@1
  - type: com.docker.sandbox/network-policy@1
    config:
      runtime:
        allow:
          - api.anthropic.com:443
```

The network rule permits HTTPS requests to the Anthropic API while the
sandbox runs. Downloading Claude Code in the Dockerfile uses the builder's
network, so it doesn't need a rule here.

The `sbx@1` entry declares that the workload is prepared for `sbx` to manage
the agent's launch. See [Declare the agent launch contract](/manuals/ai/sandboxes/customize/author/base-images.md#declare-the-agent-launch-contract)
for the shells and user account that this requires.

## Declare the credential

Claude Code can reach the API, but it also needs to authenticate. Add a
credential capability to tell Docker Sandboxes which key to use and how to
include it in requests. You'll store the key on your host, outside the kit.

Append this entry to the same `capabilities` list:

```yaml {title="Append under capabilities"}
  - type: com.docker.sandbox/credential@1
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
key into the `x-api-key` header. The key stays on the host. You'll supply
its value and approve its use when launching the kit.

## Write the model settings

The Dockerfile's launch command reads
`/home/agent/.config/claude-team/settings.json`. Use the lifecycle capability
to create that file with the model chosen for the sandbox.

Append this entry to `capabilities`:

```yaml {title="Append under capabilities"}
  - type: com.docker.sandbox/lifecycle@1
    config:
      files:
        - path: /home/agent/.config/claude-team/settings.json
          content: |
            {"model": "${{ kit.args.model }}"}
          mode: "0644"
```

Docker Sandboxes fills in the chosen `model` and writes the file before
Claude Code starts. Creating the file at this point lets each sandbox use a
different model with the same image.

Docker Sandboxes writes these files as UID 1000 after running install hooks.
Use an absolute path the agent can write to, as in this example. Shell
variables such as `$HOME` aren't expanded in the path.

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
  - type: com.docker.sandbox/agent-context@1
    config:
      filename: CLAUDE.md
      contentFile: ./context.md
```

The build includes `context.md` in the image. When the sandbox runs, `sbx`
writes a `CLAUDE.md` file in the parent directory of the mounted workspace.
That file points Claude Code to your `context.md`. A `CLAUDE.md` in your
project stays in place.

Use `contentFile`, as in this example, to keep longer instructions in a
separate Markdown file. You can also write instructions directly in the
workload's descriptor using `content` instead of `contentFile`.
For how instructions from multiple kits work together, see
[Runtime access and instructions](/manuals/ai/sandboxes/customize/use-kits.md#runtime-access-and-instructions).

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
When prompted, approve the kit's request to use your stored key, then follow
Claude Code's first-run prompts. The agent needs both the stored key and your
approval to use it. See
[Credential bindings](/manuals/ai/sandboxes/configuration/credentials.md#credential-bindings).

Choose a different model when creating a sandbox:

```console
$ sbx run --name claude-team-opus ./claude-team \
    --kit-arg claude-team.model=opus
```

The `claude-team` prefix matches the local kit directory's name. Docker
Sandboxes checks that `opus` is one of the choices in `enum`, then writes it
to the settings file before Claude Code starts.

## Iterate and publish

Edit the descriptor, Dockerfile, or context file and create another sandbox
with a different name to test the changes:

```console
$ sbx run --name claude-team-test-2 ./claude-team
```

Restarting an existing sandbox won't pick up your edits. To try another
agent version locally, update `args.version.default` in the descriptor and
create another sandbox.

`sbx` reuses the local build when the source files and kit arguments are
unchanged. Changing either can trigger a build, even for arguments such as
`model` that only affect sandbox setup. BuildKit can still reuse unchanged
image layers.

### Publish the workload

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

The `--file` option tells Buildx to read the descriptor. Its `syntax` line
selects the kit frontend, which reads the companion Dockerfile and includes
the descriptor in the published image. To publish a different agent version,
add `--build-arg version=<CLAUDE_VERSION>`. This flag uses the kit argument
name, `version`, which the descriptor maps to the Dockerfile's `CLAUDE_VERSION`.

Run the published kit by its image reference:

```console
$ sbx run --name claude-team-shared docker.io/<NAMESPACE>/claude-team:1.0.0
```

For multi-platform images and distribution details, see
[Build and distribute kits](/manuals/ai/sandboxes/customize/author/distribute.md).
To learn how to package a tool separately from its workload, see
[Build a tool mixin](/manuals/ai/sandboxes/customize/author/tool-mixins.md).
That tutorial also packages Claude Code, so use its mixin with a shell workload.
Combining it with this workload would give two kits that provide `claude`,
which Docker Sandboxes rejects.
You can also include the published workload in a
[kit set](/manuals/ai/sandboxes/customize/author/kit-sets.md) and add settings
and instructions there. The workload still defines how to prepare the base
image, install the agent, and launch it.
