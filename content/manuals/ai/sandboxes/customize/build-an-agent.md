---
title: Build your own agent kit
linkTitle: Build an agent
description: Build a schema v3 Claude Code workload kit with a pinned agent binary, runtime configuration, proxy-managed credentials, and agent instructions.
keywords: sandboxes, sbx, kits, agent, tutorial, claude, workload, build
weight: 30
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> Kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change. Share
> feedback in the [docker/sbx-releases](https://github.com/docker/sbx-releases)
> repository.

Build a schema v3 workload kit that runs Claude Code with a pinned binary,
configurable model, and an Anthropic API key held on the host. The same steps
apply to other agents: build the software into an image, declare its runtime
requirements, and provide instructions about the environment. Start from Docker's
shell template, or [use your own base image](#use-your-own-base-image), such as
Red Hat Universal Base Image 9.

This example uses API-key authentication. If you're starting with kits, read
[Kits v3](kits.md) for the file layout and the roles of a workload and a mixin.
The walkthrough builds up one kit, adding each part of its descriptor as it is
needed. For field definitions, see the [Kit spec reference](kit-reference.md).

This creates an independent v3 workload. It doesn't inherit the built-in
`claude` kit, which uses v2. Any mixins you add must also use v3. For existing
customizations in the deprecated format, see [Kits v2](kits-v2/_index.md).

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

## Build the agent into the image

First, define what goes into the image. The Dockerfile installs Claude Code
and launches it with a settings file. The descriptor you write next will
supply the agent version and create that settings file in the sandbox.

Save the following as `claude-team/claude-team.dockerfile`:

```dockerfile {title="claude-team/claude-team.dockerfile"}
FROM docker/sandbox-templates:shell

USER agent
ARG CLAUDE_VERSION
ENV PATH="/home/agent/.local/bin:${PATH}" \
    IS_SANDBOX=1 \
    CLAUDE_ENV_FILE=/etc/sandbox-persistent.sh

RUN curl -fsSL https://claude.ai/install.sh -o /tmp/install-claude.sh \
    && bash /tmp/install-claude.sh "${CLAUDE_VERSION}" \
    && rm /tmp/install-claude.sh

WORKDIR /home/agent/workspace
ENTRYPOINT ["claude", "--settings", "/home/agent/.config/claude-team/settings.json"]
CMD []
```

The `shell` template supplies the sandbox environment, including Bash, Git,
curl, certificates, and the `agent` user at UID 1000. Installing as `agent`
puts Claude Code under `/home/agent/`, where the launch user can access it.
The `CLAUDE_VERSION` build argument receives its value from the descriptor in
the next step.

The Dockerfile owns the image's `ENTRYPOINT`, `CMD`, environment, user, and
working directory. `CMD []` clears any inherited arguments. Claude Code's
`--settings` option reads an additional settings file that the kit writes
during sandbox creation.

Installing Claude Code belongs in the build recipe because the binary is the
same in every sandbox using this kit. BuildKit can cache that work. Reserve
lifecycle install hooks for configuration that depends on an individual
sandbox, such as registering a runtime endpoint.

### Use your own base image

Your workload can start from Red Hat Universal Base Image (UBI), or an image
maintained by your organization. A Docker sandbox template is optional. The
workload's Dockerfile prepares the operating system and installs the agent,
while its v3 descriptor declares network access, credentials, and other
runtime configuration.

To use UBI 9 for this walkthrough, replace `claude-team.dockerfile` with the
following. Keep the directory layout and follow the remaining descriptor
steps on this page.

```dockerfile {title="claude-team/claude-team.dockerfile"}
FROM registry.access.redhat.com/ubi9/ubi:9

USER root
RUN dnf install -y bash ca-certificates curl-minimal git shadow-utils tar gzip \
    && dnf clean all

RUN groupadd --gid 1000 agent \
    && useradd --uid 1000 --gid 1000 --create-home --shell /bin/bash agent \
    && mkdir -p /home/agent/workspace /home/agent/.local/bin \
        /home/agent/.local/share /home/agent/.local/state \
        /home/agent/.config/claude-team /home/agent/.docker/sandbox/locks \
    && chown -R agent:agent /home/agent

# Prepare the certificate bundle used by the sandbox proxy and clients.
RUN update-ca-trust \
    && mkdir -p /usr/local/share/ca-certificates /etc/ssl/certs \
    && cp /etc/pki/tls/certs/ca-bundle.crt /etc/ssl/certs/ca-certificates.crt

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt \
    CURL_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt \
    REQUESTS_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt \
    NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

# Give the agent a writable file for persistent shell exports.
RUN touch /etc/sandbox-persistent.sh \
    && chown agent:agent /etc/sandbox-persistent.sh \
    && chmod 0644 /etc/sandbox-persistent.sh \
    && printf '%s\n' '. /etc/sandbox-persistent.sh' \
        > /etc/profile.d/sandbox-persistent.sh \
    && printf '%s\n' '. /etc/sandbox-persistent.sh' >> /home/agent/.bashrc

ENV HOME=/home/agent \
    PATH="/home/agent/.local/bin:${PATH}" \
    BASH_ENV=/etc/sandbox-persistent.sh \
    CLAUDE_ENV_FILE=/etc/sandbox-persistent.sh \
    IS_SANDBOX=1

USER agent
ARG CLAUDE_VERSION
RUN curl -fsSL https://claude.ai/install.sh -o /tmp/install-claude.sh \
    && bash /tmp/install-claude.sh "${CLAUDE_VERSION}" \
    && rm /tmp/install-claude.sh

WORKDIR /home/agent/workspace
ENTRYPOINT ["claude", "--settings", "/home/agent/.config/claude-team/settings.json"]
CMD []
```

The preparation covers three parts of the environment that the shell template
would otherwise supply:

- A non-root `agent` account with UID 1000, home directory `/home/agent`, and
  writable directories for agent configuration and state
- Bash, Git, curl, CA certificates, and the archive tools used during agent
  installation
- A persistent environment file, sourced by Bash and used by Claude Code,
  matching the instructions added later in this walkthrough

The certificate step copies UBI's public CA roots to the bundle path used by
`sbx`. At sandbox startup, Docker Sandboxes adds its proxy CA to this bundle.
The environment variables direct clients to the combined bundle so HTTPS
requests can trust the proxy. If you add corporate CA certificates during the
build, add them to `/etc/pki/ca-trust/source/anchors/` before running
`update-ca-trust` and copying the bundle.

For an organization-maintained image, adapt the package installation and
account creation to its existing contents. Preserve the
[base image requirements](kit-reference.md#base-image-requirements), and
install the libraries and tools your agent and projects need. This example
runs the agent without `sudo`. Install additional system packages in the
Dockerfile, or use a runtime install hook when configuration depends on the
individual sandbox.

The base image supplies filesystem content and image settings. It doesn't
supply a kit's runtime declarations: continue below to configure API access,
credentials, settings, and agent instructions. Any mixins added to this
workload must use v3 and provide tools compatible with its operating system
and architecture.

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
[Agent instructions](kits.md#give-the-agent-instructions) for how workload
and mixin instructions contribute to the profile.

## Store the key and run

Store your Anthropic API key on the host:

```console
$ sbx secret set anthropic
```

From the directory containing `claude-team`, launch the kit against your
project:

```console
$ sbx run --name claude-team ./claude-team <PROJECT_PATH>
```

`sbx` builds the local directory, loads the kit, and launches Claude Code.
Approve the kit's credential request to connect the stored key to this kit,
then follow Claude Code's first-run prompts. Without an approved credential
binding, storing a key alone doesn't authenticate the agent. See
[Credential bindings](../configuration/credentials.md#credential-bindings).

Choose a different model when creating a sandbox:

```console
$ sbx run --name claude-team-opus ./claude-team \
    --kit-arg claude-team.model=opus <PROJECT_PATH>
```

The argument prefix is the local kit directory's name. The model value is
validated against the descriptor's `enum` and written to the settings file
before Claude Code starts.

## Iterate and publish

Edit the descriptor, Dockerfile, or context file and create another sandbox
with a different name to test the changes:

```console
$ sbx run --name claude-team-test-2 ./claude-team <PROJECT_PATH>
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
$ sbx run --name claude-team-shared docker.io/<NAMESPACE>/claude-team:1.0.0 <PROJECT_PATH>
```

For build layouts, multi-platform images, and distribution details, see
[Kits v3](kits.md). To add tools or shared configuration to this workload,
see [Kit examples](kit-examples.md).
