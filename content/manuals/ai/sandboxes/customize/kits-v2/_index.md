---
title: Kits v2
linkTitle: Kits v2
description: Maintain schema v2 sandbox and mixin kits with the spec.yaml format, runtime setup, built-in agent inheritance, and v2 packaging commands.
keywords: sandboxes, sbx, kits, mixins, customization, extensions, agents
weight: 60
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

This page covers kits with `schemaVersion: "2"`. Use [Kits v3](../kits.md)
for authoring kits with image builds and runtime capabilities. Docker Sandboxes
also supports v1 and v2 kits. A composition cannot combine v3 kits with v1 or v2
kits.

> [!NOTE]
> Kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change as the
> feature evolves. Share feedback and bug reports in the
> [docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

For kit concepts and use cases, see [What kits can do](../kits.md#what-kits-can-do).
The following sections describe the v2 fields and commands for maintaining
existing kits.

V2 kits use a `spec.yaml` file. A `kind: sandbox` kit selects the agent's image
and launch configuration. A `kind: mixin` kit adds configuration to an agent.
Load a kit from a directory, ZIP file, OCI artifact, or Git URL.

## V2 configuration

### Run commands

V2 declares runtime commands under `setup`. `setup.install` runs once at
sandbox creation and can install software:

```yaml
setup:
  install:
    - command: "apt-get update && apt-get install -y jq"
```

`setup.startup` runs on each start alongside the agent. Commands must tolerate
repeated execution. See [`startup`](kit-reference.md#startup) for its fields:

```yaml
setup:
  startup:
    - command: ["my-daemon"]
      background: true
```

### Inject files

V2 places bundled files from `files/home/` in the agent's home directory and
from `files/workspace/` in the workspace:

```text
my-kit/
├── spec.yaml
└── files/
    ├── home/
    │   └── .config/my-tool/settings.json
    └── workspace/
        └── .editorconfig
```

`setup.files` cover content that depends on runtime values, such as an
absolute workspace path that a tool needs to bake into its config file
at startup:

```yaml
setup:
  files:
    - path: /home/agent/.my-tool/config.json
      content: '{"workspace": "${WORKDIR}"}'
      onlyIfMissing: true
```

See [`setup.files`](kit-reference.md#files) in the spec reference for all
fields.

#### Sandbox-managed agent configuration

Built-in agent kits reserve the following paths for sandbox setup. Treat these
paths as sandbox-managed, even if a file is only needed for a particular
feature. Don't target them with static files, `setup.files`, or install
commands. Later setup can replace your content or depend on settings that your
file removes. In this table, `~` is `/home/agent`.

| Built-in agent kit | Managed configuration paths |
| ------------------ | --------------------------- |
| `claude` | `~/.claude.json`, `~/.claude/settings.json`, `~/.claude/.config.json` |
| `codex` | `~/.codex/config.toml` |
| `copilot` | `~/.copilot/config.json` |
| `cursor` | `~/.cursor/cli-config.json` |
| `devin` | `~/.config/devin/config.json`, `~/.config/devin/mcp_config.json` |
| `gemini` | `~/.gemini/settings.json` |
| `kiro` | `~/.kiro/settings/mcp.json` |
| `opencode` | `~/.config/opencode/opencode.json` |

Use a separate settings layer when the agent supports one. For example, Claude
Code can load an additional settings file with `--settings`, and OpenCode can
load one from the path in `OPENCODE_CONFIG`. See
[Customize agent settings](kit-examples.md#customize-agent-settings) for
examples. Don't use `setup.startup` for settings the agent must read during
initialization because startup commands don't gate the agent entrypoint.

### Set environment variables

Declare environment variables under `environment.variables`:

```yaml
environment:
  variables:
    MY_TOOL_WORKSPACE: /home/agent/my-tool
```

For credentials, see
[Authenticate to external services](#authenticate-to-external-services).
Don't put secret values directly in `environment.variables` — they'd
be visible inside the sandbox VM.

Leave proxy variables to the sandbox. See
[Set environment variables](../kits.md#set-environment-variables) for shared
environment and proxy guidance.

### Control network access

V2 declares network rules under `permissions.network`:

```yaml
permissions:
  network:
    allow:
      - api.example.com
      - "*.cdn.example.com"
    deny:
      - telemetry.example.com
```

Kit deny rules take precedence over kit allow rules. Organization governance
can also override kit allow rules. See [Control network access](../kits.md#control-network-access)
for policy behavior shared by kit versions.

### Authenticate to external services

V2 declares credentials in a top-level `credentials` list. The injection
domain must also be allowed by `permissions.network`:

```yaml
credentials:
  - service: my-service
    apiKey:
      name: MY_SERVICE_API_KEY # in-VM env var, set to a sentinel
      proxyManaged: true
      inject:
        - domain: api.example.com # inject on requests to this domain
          header: Authorization # overwrite this header
          format: "Bearer %s"

permissions:
  network:
    allow:
      - api.example.com # the domain must also be reachable
```

See [Credential configuration](../../configuration/credentials.md) for storing
secrets and approving credential bindings. The proxy behavior is shared with
v3; the descriptor fields differ.

### Inject agent memory

V2 uses `agentInstructions` for guidance such as tool usage and project
conventions:

```yaml
agentInstructions:
  content: |
    Ruff is installed. Run `ruff check` before committing.
    Shared config lives at `/workspace/ruff.toml`.
```

Both mixin and sandbox kits can declare `agentInstructions.content`. The active
sandbox kit sets `agentInstructions.filename`, which determines the memory
file's name. The sandbox kit's content is written inline in that file. Each
mixin's content is written to its own `<kit-name>.md` file under a sibling
`kits-memory/` directory, and the main memory file gets a `## Kits` section that
points to each mixin file:

```text
/Users/you/
├── myproject/              # workspace
├── AGENTS.md               # main memory file with a "## Kits" index
└── kits-memory/
    ├── ruff-lint.md
    ├── vale.md
    └── git-ssh-sign.md
```

See [`agentInstructions`](kit-reference.md#agent-instructions) in the spec
reference for the full field schema.

### Define an agent

Sandbox kits declare a `sandbox:` block with the image the agent runs in and
the command the user attaches to when they launch the sandbox:

```yaml
sandbox:
  image: "my-registry/my-agent:latest"
  entrypoint: [my-agent, "--yolo"]
```

See [Sandbox kits](#sandbox-kits) for inheritance from a built-in agent.

## Mixin kits

A v2 `kind: mixin` kit extends a built-in agent or a v2 sandbox kit. Pass it
with `--kit` when creating the sandbox.

See [Drop a shared config file](kit-examples.md#drop-a-shared-config-file) and
[Install a tool at sandbox creation](kit-examples.md#install-a-tool-at-sandbox-creation)
for complete mixin examples.

## Sandbox kits

Sandbox kits declare everything a mixin kit can, plus an
[`sandbox:` block](kit-reference.md#sandbox-block) that tells the sandbox how to launch the
agent. For a walkthrough of the v3 workload format, see
[Build your own agent kit](../build-an-agent.md).

### Extend a built-in agent

Use `extends:` to create a variant of a built-in agent without reproducing its
configuration. The child kit inherits the parent's image, credentials, network
permissions, persistent volumes, settings, MCP integration, and agent
instructions. It also inherits the parent's environment variables and all
`setup.install`, `setup.startup`, and `setup.files` entries. Parent setup entries
run before child entries. If both kits set the same environment variable, the
child's value wins. Use `extends:` for a single parent agent; use a mixin to add
an independent capability that can work with one or more agents. See
[Fork an existing agent](kit-examples.md#fork-an-existing-agent) for an example
that changes Claude Code's permission mode.

## Using kits

Kits can be loaded from a local path (a directory or ZIP file), a Git
repository, or an OCI registry. To launch a sandbox kit, pass its reference in
place of a built-in agent name to `sbx run` or `sbx create`. Use `--kit` for
mixins, and repeat the flag to apply multiple mixins to the same sandbox.

Starting with Docker Sandboxes version 0.42.0, pass the sandbox kit reference
as the first argument:

```console
$ sbx run <sandbox-kit-ref> [PATH...]
$ sbx create <sandbox-kit-ref> [PATH...]
```

The previous form, `sbx run <sandbox-kit-name> --kit <sandbox-kit-ref>`, is
deprecated.

> [!IMPORTANT]
> A mixin passed with `--kit` only takes effect when a sandbox is created.
> Passing it against an
> existing sandbox name fails with
> `--kit can only be used when creating a new sandbox`. To add a supported
> mixin kit to a running sandbox, use [`sbx kit add`](#local) instead.
> `sbx kit add` restarts the sandbox to apply the updated kit set.
> VM state — installed packages, Docker images, volumes, and agent history
> — is preserved across the restart. It supports mixin kits limited to
> `environment.variables`, `setup.install`, and `permissions.network.allow`.
> To use other fields, recreate the sandbox with the mixin.

### Pass arguments to kits

A schema v2 kit can declare inputs in a top-level `args:` block and reference
them in `spec.yaml` or static files with `${{ kit.args.<name> }}`. Supply a
value with `--kit-arg name=value`:

```console
$ sbx run ./my-agent/ --kit-arg channel=beta
```

Kit argument values are plain text. Values supplied with `--kit-arg` can remain
in your shell history, and argument files store their values unencrypted. Don't
use kit arguments for secrets. Use [Credentials](../../configuration/credentials.md)
instead.

An argument without a kit name prefix applies to every kit that declares it.
To target one kit, prefix the argument with the value of that kit's `name`
field and a period:

```console
$ sbx run ./my-agent/ \
    --kit ./my-mixin/ \
    --kit-arg version=1.2.3 \
    --kit-arg my-mixin.version=2.0.0
```

The kit-specific value takes precedence over the shared value for `my-mixin`.

Use `--kit-args-file` for a reusable set of `name=value` entries. Blank lines
and lines that start with `#` are ignored:

```text {title="kit.args"}
version=1.2.3
my-mixin.channel=beta
```

```console
$ sbx create ./my-agent/ . \
    --kit ./my-mixin/ \
    --kit-args-file ./kit.args \
    --kit-arg my-mixin.channel=stable
```

When you pass multiple argument files, a value in a later file overrides the
same key in an earlier file. Values passed with `--kit-arg` override every
file. For repeated `--kit-arg` entries with the same key, the last value wins.

Argument validation happens before the sandbox is created. `sbx` rejects a
missing required value, a value outside its declared `enum` or `pattern`, a
placeholder without a declaration, and a supplied argument that no resolved
kit declares. Pass the same argument flags to `sbx kit validate` or
`sbx kit inspect` when the kit requires values. See
[Kit arguments](kit-reference.md#arguments) for the declaration fields.

### Local

Launch a local sandbox kit by passing its directory or ZIP file in place of the
agent name. Relative paths must start with `./` or `../` so `sbx` can
distinguish them from agent and sandbox names:

```console
$ sbx run ./my-agent/
$ sbx create ../my-agent-1.0.zip .
```

Pass a local mixin with `--kit`:

```console
$ sbx run claude --kit ./my-mixin/
$ sbx run claude --kit ../my-mixin-1.0.zip
```

While iterating on a supported mixin kit, apply changes to a running sandbox
with `sbx kit add`:

```console
$ sbx kit add my-sandbox ./my-kit/
```

`sbx kit add` restarts the sandbox to apply the updated kit set.
VM state — installed packages, Docker images, volumes, and agent history — is
preserved across the restart. Kits can't be removed from a running sandbox —
remove and recreate it to start clean.

### Git repository

Launch a sandbox kit from a Git repository:

```console
$ sbx run "git+https://github.com/docker/sbx-kits-contrib.git#ref=v0.1.0&dir=amp"
```

Pass a Git-hosted mixin with `--kit`:

```console
$ sbx run claude --kit "git+https://github.com/docker/sbx-kits-contrib.git#ref=v0.1.0&dir=code-server"
```

- `#ref=<branch|tag|commit>` pins to a specific revision. Defaults to the
  repository's default branch.
- `#dir=<path>` loads a kit from a subdirectory.
- `git+ssh://` URLs also work, using your local SSH agent, Git credential
  helpers, and `.netrc`.
- Quote the URL in shells where `&` starts a background job.

### OCI registry

Launch a sandbox kit from an OCI registry:

```console
$ sbx run docker.io/sbx/droid-kit:latest
```

Pass an OCI-hosted mixin with `--kit`:

```console
$ sbx run claude --kit ghcr.io/myorg/my-kit:1.0
```

For Docker Hub, include the full `docker.io` prefix. See
[Packaging and distribution](#packaging-and-distribution) for publishing.

> [!IMPORTANT]
> For Docker Hub, `sbx` reuses your `sbx login` session to pull private
> kits. For other registries, store pull credentials with
> [`sbx secret set --registry`](../../configuration/credentials.md#registry-credentials)
> before running the sandbox. These credentials take priority over credentials
> in the Docker credential store:
>
> ```console
> $ gh auth token | sbx secret set --registry ghcr.io --password-stdin
> ```
>
> Without credentials from either store, pulls from non-Docker Hub registries
> are anonymous and private kits fail to pull.

### Restrict kit sources

`sbx` restricts which sources a kit can install from. A kit's install
commands run with root privileges inside the sandbox, so limiting where kits
come from reduces supply-chain risk. By default, only kits hosted on Docker
Hub (`docker.io/`) are allowed. Loading a kit from any other source fails:

```console
$ sbx run claude --kit "git+https://github.com/docker/sbx-kits-contrib.git#dir=vale"
ERROR: resolve kits: kit "git+https://github.com/docker/sbx-kits-contrib.git#dir=vale" cannot be installed — its source is not in your allowlist.
```

To allow another publisher, add its host or host/path prefix to the
`kit.allowedSources` setting. The setting replaces the whole list, so include
the entries you want to keep:

```console
$ sbx settings set kit.allowedSources '["docker.io/","github.com/docker/"]'
```

Entries match as prefixes on a path-segment boundary, so `github.com/docker/`
allows `github.com/docker/sbx-kits-contrib` but not `github.com/docker-evil/kit`.
To remove the restriction and allow any remote source, set the list to
`["*"]`. This isn't recommended.

Installing from a local directory or ZIP file is governed separately by the
`kit.allowLocalKits` setting, which defaults to `true`. Set it to `false` to
require a remote source:

```console
$ sbx settings set kit.allowLocalKits false
```

For non-interactive use, both settings have environment-variable equivalents:
`DOCKER_SANDBOXES_KIT_ALLOWED_SOURCES` and `DOCKER_SANDBOXES_KIT_ALLOW_LOCAL`.

## Sign and verify kits

Use cosign-compatible Sigstore signatures to verify who approved a kit and
that its signed content hasn't changed. Signing is keyless by default. Verify a
keyless signature with the certificate identity and OpenID Connect (OIDC)
issuer:

```console
$ sbx kit sign ./my-kit/
$ sbx kit verify \
    --certificate-identity user@example.com \
    --certificate-oidc-issuer https://accounts.google.com \
    ./my-kit/
```

For key-based signing, use an ECDSA P-256 key pair:

```console
$ sbx kit sign --key cosign.key ./my-kit/
$ sbx kit verify --key cosign.pub ./my-kit/
```

For a local directory, `sbx kit sign` writes a `kit.sig.bundle` file next to
`spec.yaml`. Commit this file so consumers can verify a kit loaded from the Git
repository. For an OCI kit, the signature is stored as an OCI referrer. You can
sign an OCI kit after pushing it, or push and sign it in one step:

```console
$ sbx kit push ./my-kit/ ghcr.io/myorg/my-kit:1.0 --sign
```

ZIP kits can't carry verifiable signatures.

### Require signed kits

Set a trusted signer policy for the identities or keys you trust before
requiring signatures. Otherwise, `sbx` uses the default policy, which trusts
Docker employee identities attested by Google's OpenID Connect issuer. A
keyless policy must specify both the certificate identity and its OpenID
Connect issuer:

```console
$ sbx settings set kit.trustedSigners \
    '[{"identity":"release-bot@example.com","issuer":"https://accounts.google.com"}]'
$ sbx settings set kit.requireSignature true
```

To trust a key-based signature, set the policy to the public key path:

```console
$ sbx settings set kit.trustedSigners '[{"key":"/path/to/cosign.pub"}]'
$ sbx settings set kit.requireSignature true
```

When `kit.requireSignature` is `true`, `sbx` rejects unsigned kits, signatures
that don't match `kit.trustedSigners`, and ZIP kits. This policy applies when a
kit is loaded from a local directory, Git repository, or OCI registry.

The signature covers `spec.yaml` and the kit's `files/` content, but not mutable
dependencies such as image tags or content downloaded by install and startup
commands. Pin those dependencies by digest or checksum when they must remain
immutable.

## Packaging and distribution

The `sbx kit` subcommands validate, inspect, and publish kits:

- `sbx kit validate <path>` — check that a kit directory or ZIP is
  well-formed.
- `sbx kit inspect <path>` — display kit details. Add `--json` for
  machine-readable output.
- `sbx kit pack <path> -o <file.zip>` — package a directory as a ZIP file
  for sharing.
- `sbx kit push <path> <ref>` — publish to an OCI registry (for example,
  `ghcr.io/myorg/my-kit:1.0`).
- `sbx kit pull <ref>` — download a kit from a registry as a ZIP file to
  the working directory.

For Docker Hub, include the full `docker.io` prefix — `sbx` doesn't add it
automatically.

For Docker Hub, `sbx kit pull` and `sbx kit push` use the session from
`sbx login`. For other registries, they prefer credentials stored with
[`sbx secret set --registry`](../../configuration/credentials.md#registry-credentials).
Both commands fall back to the Docker credential store, so credentials from
`docker login` also work.

## Spec reference

For a field-by-field reference of every `spec.yaml` block — top-level
fields, arguments, credentials, network, environment, setup, static files,
agent instructions, and the sandbox block — see [Kit spec reference](kit-reference.md).

## Debugging

See [Debug kits](../kits.md#debug-kits) for network policy logs, inspecting
files and tools inside the sandbox, and capturing background-service output.
Use `sbx kit validate` and `sbx kit inspect` to check a v2 descriptor before
launching it.
