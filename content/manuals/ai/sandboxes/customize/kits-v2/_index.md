---
title: Kits v2
linkTitle: Kits v2
description: Reference for v2 kits, including usage, schema fields, maintenance examples, signing, and migration to v3.
keywords: sandboxes, sbx, kits, v2, migration, spec.yaml
weight: 60
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

V2 kits remain supported. This page covers v2 usage, configuration, and the
specification. For new kit development, use [v3 kits](/manuals/ai/sandboxes/customize/_index.md), which are
experimental.

Built-in shortcuts such as `claude` and `codex` select v2 kits and still work
with v2 mixins. V3 workloads and mixins can't be combined with v1 or v2 kits.
V1 also remains supported.

## Move an environment to v3

Select a v3 workload, convert or replace its mixins, and create a separate
sandbox with a different `--name`. Use the explicit workload reference in
place of the built-in shortcut. Every selected kit must use v3.

Changing `schemaVersion` alone doesn't convert a kit. See the
[v2-to-v3 field mapping](/manuals/ai/sandboxes/customize/author/kit-reference.md#move-from-v2-to-v3), the
[v3 runtime support table](/manuals/ai/sandboxes/customize/author/kit-reference.md#runtime-capabilities), and the
[agent authoring tutorial](/manuals/ai/sandboxes/customize/author/build-an-agent.md). Running an existing sandbox
keeps its recorded configuration; it doesn't migrate the kit set.

## Use existing kits

A v2 kit contains `spec.yaml` and an optional `files/` tree. Pass a sandbox kit
in place of the agent name and add mixins with `--kit`:

```console
$ sbx run ./my-agent --name my-project --kit ./team-config <PROJECT_PATH>
$ sbx run claude --name claude-project --kit ./team-config <PROJECT_PATH>
```

Use `sbx create` with the same arguments to create without launching the
agent. References can be local directories, ZIP files, OCI artifacts, or Git
URLs. Start relative paths with `./` or `../` and include `docker.io/` for
Docker Hub references. In Git URLs, `ref` selects a revision and `dir` the kit
directory. Quote URLs containing `&`:

```console
$ sbx run "git+https://github.com/<ORG>/<REPOSITORY>.git#ref=<COMMIT>&dir=my-agent" <PROJECT_PATH>
```

`git+ssh://` URLs work with your local SSH agent and Git credentials. See
[Restrict kit sources](/manuals/ai/sandboxes/customize/use-kits.md#restrict-kit-sources) for source policies;
`kit.allowLocalKits` also governs v2 ZIP files. For private registries, see
[Registry credentials](../../configuration/credentials.md#registry-credentials).

Kit selection with `--kit` applies at creation. Recreate the sandbox to change
its kit set, except for the limited updates supported by
[`sbx kit add`](#execution-order). That command restarts the sandbox while
preserving packages, images, volumes, and agent history. Kits can't be
removed from a running sandbox.

## Schema versions

Schema v2 is supported starting with Docker Sandboxes version 0.36. Use
`schemaVersion: "2"` for the syntax on this page. Version `"1"` also remains
accepted. V3 is a separate format for environments built entirely
with v3 workloads and mixins. V3 kits can't compose with v1 or v2 kits.
See [Kits v3](/manuals/ai/sandboxes/customize/_index.md) for that workflow.

The loader forks on `schemaVersion`. A v2 spec uses the v2 grammar only. Legacy
v1 fields in a `schemaVersion: "2"` spec are rejected during decode instead of
being folded into the v2 model. Keep each `spec.yaml` on one grammar.

What changed in v2:

| v1                                          | v2                                       |
| ------------------------------------------- | ---------------------------------------- |
| `credentials.sources.<id>`                  | `credentials:` list entry with `service` |
| `network.allowedDomains` / `deniedDomains`  | `permissions.network.allow` / `deny`     |
| `network.serviceDomains` / `serviceAuth`    | `credentials[].apiKey.inject`            |
| `network.publishedPorts` / `publishedPorts` | top-level `ports`                        |
| standalone `oauth:` block                   | `credentials[].oauth`                    |
| `oauth.skipIfEnv`                           | Accepted but ignored                     |
| `environment.proxyManaged`                  | `credentials[].apiKey.proxyManaged`      |
| `memory` / `agentContext`                   | `agentInstructions.content`              |
| `kind: agent` / `agent:` block              | `kind: sandbox` / `sandbox:` block       |
| `sandbox.aiFilename`                        | `agentInstructions.filename`             |
| `sandbox.entrypoint.run`                    | `sandbox.entrypoint`                     |
| `sandbox.entrypoint.args`                   | `sandbox.command.default`                |
| `sandbox.entrypoint.ttyArgs`                | `sandbox.command.interactive`            |
| `tmpfs:`                                    | `volumes:` entries with `type: tmpfs`    |
| `volumes:` (mapping form)                   | `volumes:` sequence (`- path: <path>`)   |
| `commands:` / `commands.initFiles`          | `setup:` / `setup.files`                 |
| `settings:` / `kitDir` / `persistence`      | Removed                                  |

Credential discovery also moved out of the kit in v2: a kit declares which
credentials it needs and how to inject them, but where each value comes from is
controlled by the user through
[credential bindings](../../configuration/credentials.md#credential-bindings).

> [!NOTE]
> `mixins` and `sandbox.build` are accepted by the parser, but runtime support
> is pending. A kit that sets `sandbox.build` must also set `sandbox.image`.

## Top-level fields

For the normative grammar, see the
[v2 specification](https://github.com/docker/sbx-kits-contrib/blob/main/spec/SPEC-v2.md).

| Field           | Required | Description                                                                                     |
| --------------- | -------- | ----------------------------------------------------------------------------------------------- |
| `schemaVersion` | Yes      | Spec schema version. Use `"2"` for this grammar.                                                |
| `kind`          | Yes      | `mixin` for kits that extend an agent; `sandbox` for kits that define one.                      |
| `name`          | Yes      | Unique identifier. Lowercase alphanumeric with hyphens, 1 to 64 characters.                     |
| `version`       | No       | Kit version.                                                                                    |
| `displayName`   | No       | Human-readable name.                                                                            |
| `description`   | No       | Short description.                                                                              |
| `sourceURL`     | No       | Source repository or documentation URL.                                                         |
| `licenses`      | No       | SPDX license identifiers.                                                                       |
| `locked`        | No       | Dotted paths child kits may not override.                                                       |
| `security`      | No       | Container security settings. `security.privileged: true` runs the container in privileged mode. |
| `args`          | No       | Arguments supplied when the kit is loaded. Schema v2 only.                                      |

A kit also declares behavior blocks such as `agentInstructions`,
`permissions`, `ports`, `credentials`, `environment`, `setup`, and `volumes`.

## Arguments

A schema v2 kit can declare arguments and reference them anywhere in
`spec.yaml` or under `files/` as `${{ kit.args.<name> }}`. Substitution happens
before the spec is decoded.

```yaml
args:
  version:
    default: latest
    description: Tool version to install
    pattern: '^(latest|[0-9]+\.[0-9]+\.[0-9]+)$'
  channel:
    default: stable
    enum: [stable, beta, nightly]
  target:
    required: true
    description: Build target

environment:
  variables:
    TOOL_VERSION: "${{ kit.args.version }}"
```

Don't use kit arguments for API tokens, passwords, or other secrets. Use
[Credentials](../../configuration/credentials.md) to provide sensitive values to
a sandbox.

| Field         | Description                                                                                                  |
| ------------- | ------------------------------------------------------------------------------------------------------------ |
| Argument name | Starts with a letter or underscore and contains only letters, digits, underscores, and hyphens.             |
| `default`     | String to use when the caller supplies no value. Mutually exclusive with `required: true`.                   |
| `required`    | Set to `true` when the caller must supply a value. Mutually exclusive with `default`.                        |
| `description` | Optional help text shown when a required value is missing.                                                   |
| `enum`        | Optional list of accepted values. Mutually exclusive with `pattern`.                                         |
| `pattern`     | Optional Go RE2 regular expression matched against the complete value. Mutually exclusive with `enum`.       |

Each argument must declare either `default`, including an empty-string
default, or `required: true`. A declared default must satisfy its own `enum` or
`pattern`. Every `${{ kit.args.<name> }}` reference must have a matching
declaration.

Argument values are strings, but substitution happens before YAML decoding.
Quote a placeholder in a string-valued field so a value such as `1.20` isn't
decoded as a number.

### Pass arguments to kits

Use `--kit-arg name=value` for every kit declaring that argument, or prefix
with the kit's `name` to target one kit. Scoped values override shared values:

```console
$ sbx run ./my-agent --kit ./my-mixin --kit-arg channel=stable \
    --kit-arg my-mixin.channel=beta <PROJECT_PATH>
```

`--kit-args-file <FILE>` reads `name=value` entries, ignoring blank lines and
`#` comments. Later files override earlier files; `--kit-arg` overrides files.
For repeated CLI keys, the last value wins. Missing required values, unknown
arguments, undeclared placeholders, and invalid values fail before creation.
Pass the same flags to `sbx kit validate` or `sbx kit inspect` when needed.
Argument values can remain in shell history and are stored unencrypted in
argument files.

## Kit kinds

### `kind: mixin`

A mixin layers capabilities onto an existing sandbox. It must not declare a
`sandbox:` block, `extends:`, or `mixins:`. A mixin can declare `requires:` to
pin the base agent it is designed for:

```yaml
schemaVersion: "2"
kind: mixin
name: github-tools
requires:
  agent: claude
```

`requires.agent` takes one base-agent name. It is validated as a kit name and
enforced during composition.

### `kind: sandbox`

A sandbox kit defines a full agent. A root sandbox must declare a `sandbox:`
block. A sandbox that uses `extends:` can inherit the parent image and omit its
own `sandbox:` block:

```yaml
schemaVersion: "2"
kind: sandbox
name: claude-safe
extends: claude
```

`extends:` is sandbox-only. The parent must resolve to a sandbox kit. `mixins:`
is also sandbox-only and accepted by the parser, but runtime composition support
is pending.

## Sandbox block

```yaml
sandbox:
  image: <image-ref>
  build:
    context: .
    dockerfile: Dockerfile
    args:
      AGENT_VERSION: "1.0.0"
    target: runtime
    platforms:
      - linux/amd64
  entrypoint: [my-agent, "--flag"]
  command:
    default: ["--task-mode"]
    interactive: []
  resources:
    cpu: 2
    memory: 4g
    gpu: "1"
```

| Field                | Required | Description                                                                                                     |
| -------------------- | -------- | --------------------------------------------------------------------------------------------------------------- |
| `sandbox.image`      | When `extends:` is omitted | Docker image reference.                                                                                         |
| `sandbox.build`      | No       | Build configuration. Runtime support is pending, so a kit with `build:` must also set `image:`.                 |
| `sandbox.entrypoint` | No       | Fixed process prefix as a string array. The first element is the agent binary.                                  |
| `sandbox.command`    | No       | Mode-specific argument tail. Use a list shorthand for `default`, or a mapping with `default` and `interactive`. |
| `sandbox.resources`  | No       | Optional CPU, memory, and GPU constraints. Memory uses byte-size strings such as `4096m` or `4g`.               |

The effective command is `entrypoint` plus `command.default` for non-interactive
launches, and `entrypoint` plus `command.interactive` for TTY sessions. If
`interactive` is omitted, it falls back to `default`.

For a kit that uses `extends:`, `sandbox.command` replaces the full inherited
argument tail, including flags after the binary in the parent's
`sandbox.entrypoint`. It doesn't append to that tail. Define every argument the
child needs. For example, a child of `claude` that adds `--settings` must also
include `--dangerously-skip-permissions` to preserve that behavior.

The agent's container image must provide:

- A non-root `agent` user at UID 1000 with passwordless sudo.
- A `/home/agent/` home directory owned by `agent`.
- HTTP proxy environment variables (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`) preserved across sudo.
- The agent binary, either baked in or installed with [`setup.install`](#setup).

Build on top of `docker/sandbox-templates:shell-docker` to get these base
requirements.

## Agent instructions

Declare these fields under `agentInstructions`:

| Field      | Description                                                                                         |
| ---------- | --------------------------------------------------------------------------------------------------- |
| `filename` | AI profile filename. Meaningful for `kind: sandbox`; ignored with a warning for `kind: mixin`.      |
| `content`  | Markdown instructions. For a sandbox, inlined into the profile. For a mixin, written to kit memory. |

For mixins, the engine writes `content` to
`<dir-of-AI-file>/kits-memory/<kit-name>.md` and adds a `## Kits` pointer
section to the base AI file. This keeps each mixin's instructions in a separate
file.

The generated profile lives in the parent directory of the mounted workspace
inside the sandbox. It sits outside the mount and doesn't replace an
instruction file in the project. The sandbox kit's inline instructions go
directly into that profile.

## Credentials

A kit declares the credentials it needs and how the proxy injects them into
outbound requests. It does not declare a host discovery source. The user
provides the value through the secret store or the first-run prompt, and a
[credential binding](../../configuration/credentials.md) authorizes its use. A kit
can't read arbitrary host environment variables or files.

`credentials` is a list; each entry names a `service` and configures one or more
auth mechanisms.

| Field         | Description                                                                                                                                 |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `service`     | Credential identifier, matched against the value stored with `sbx secret set`. Lowercase kebab-case.                                        |
| `description` | Optional. Shown to the user when approving a [binding](../../configuration/credentials.md#credential-bindings).                                     |
| `required`    | Marks the credential as essential to the agent. If it has no binding, `sbx` warns and starts with the credential withheld. Default `false`. |
| `provider`    | Reserved for a provider registry. Accepted with a warning and no runtime effect.                                                            |
| `apiKey`      | API-key injection (see [apiKey](#apikey)).                                                                                                  |
| `oauth`       | OAuth interception (see [oauth](#oauth)).                                                                                                   |

Each service must declare `apiKey`, `oauth`, or both. When both resolve at
runtime, the API key takes precedence and OAuth acts as the fallback.

### `apiKey`

| Field               | Description                                                                                                                                       |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`              | Environment variable name for the credential (for example, `ANTHROPIC_API_KEY`).                                                                  |
| `proxyManaged`      | If `true`, `sbx` sets `name` inside the container to the `proxy-managed` sentinel. Default `false`.                                               |
| `inject[].domain`   | Domain to inject the credential into. Must also be allowed in [`permissions.network`](#network).                                                  |
| `inject[].header`   | HTTP header the proxy sets (for example, `x-api-key`, `Authorization`).                                                                           |
| `inject[].format`   | Header value format, with one `%s` placeholder (for example, `"%s"` or `"Bearer %s"`). Mutually exclusive with `scheme`.                          |
| `inject[].scheme`   | Shorthand for common auth schemes. `bearer` expands to `Authorization: Bearer %s`; `basic` requires `username`. Mutually exclusive with `format`. |
| `inject[].username` | Username for HTTP Basic auth, for example `x-access-token` for Git over HTTPS.                                                                    |

### `oauth`

For agents that authenticate with OAuth (for example, Claude Code), the proxy
intercepts token responses and replaces real tokens with sentinels, then swaps
the real token back in on outbound requests. By default, the token never enters
the sandbox. Setting `passthrough: true` opts out of sentinel masking and sends
the real token response into the sandbox.

| Field                                    | Description                                                                                                                                                                                       |
| ---------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `tokenEndpoint.host` / `path`            | The OAuth token endpoint the proxy intercepts.                                                                                                                                                    |
| `sentinels.accessToken` / `refreshToken` | Sentinel values written into the container in place of the real tokens.                                                                                                                           |
| `credentialFile.path`                    | Where to write the credential file inside the container (`~` expands).                                                                                                                            |
| `credentialFile.structure`               | Declarative JSON shape. Supports `{{.AccessToken}}`, `{{.RefreshToken}}`, `{{.ExpiresAt}}`, and `{{.Scopes}}`.                                                                                   |
| `credentialFile.template`                | Go template. Supports `{{.AccessToken}}`, `{{.RefreshToken}}`, `{{.ExpiresAt}}`, `{{.Scopes}}`, and `{{.ScopesJSON}}`.                                                                          |
| `resourceHosts`                          | API hosts where the proxy attaches the token on outbound requests, distinct from the token endpoint host.                                                                                         |
| `skipIfEnv`                              | Accepted for compatibility, but ignored for schema v2. A v2 binding is authoritative instead of host environment variables.                                                                       |
| `responseFields`                         | Overrides the default field names the proxy reads from the token response.                                                                                                                        |
| `passthrough`                            | If `true`, the proxy passes the token response through unchanged instead of replacing the tokens with sentinels.                                                                                  |

`credentialFile.structure` provides a declarative alternative to
`credentialFile.template`. The engine renders it as well-formed JSON. If both
fields are set, `structure` takes precedence.

## Network

Network egress is declared under `permissions.network`. Credentials no longer carry
their own domain mapping — the proxy injects a credential only into the domains
its [`apiKey.inject`](#apikey) lists, and every domain the
sandbox reaches must be allowed here.

| Field                       | Description                                                                                                     |
| --------------------------- | --------------------------------------------------------------------------------------------------------------- |
| `permissions.network.allow` | Domains the sandbox can reach.                                                                                  |
| `permissions.network.deny`  | Domains the sandbox is blocked from reaching. Deny takes precedence over allow, including across composed kits. |

Allow and deny patterns:

| Pattern               | Example                  | Status                      |
| --------------------- | ------------------------ | --------------------------- |
| Exact host            | `api.example.com`        | Enforced                    |
| Exact host and port   | `api.example.com:8080`   | Enforced                    |
| Single-label wildcard | `*.example.com`          | Enforced                    |
| Multi-label wildcard  | `**.example.com`         | Parsed; enforcement pending |
| Port range            | `api.example.com:80-443` | Parsed; enforcement pending |
| Port wildcard         | `api.example.com:*`      | Parsed; enforcement pending |
| CIDR                  | `10.0.0.0/8`             | Parsed; enforcement pending |

In v1 this was the `network:` block (`allowedDomains` / `deniedDomains`, plus
`serviceDomains` / `serviceAuth`). In v2, those fields are decode errors.

## Ports

Declare `ports` as a list of entries to expose sandbox services to the host:

| Field       | Description                                                         |
| ----------- | ------------------------------------------------------------------- |
| `container` | Container port, 1 to 65535.                                         |
| `protocol`  | `tcp` or `udp`. Empty publishes one family; see below.              |
| `name`      | Optional label surfaced by tools that list published port bindings. |

Host ports are allocated ephemerally. Leave `protocol` empty unless the service
listens on IPv6: an empty value publishes IPv4 only (`127.0.0.1`), which is what
a service bound to `0.0.0.0` needs, while `tcp` publishes both `127.0.0.1` and
`::1` — and a client arriving over `::1` is accepted and then reset if nothing
in the sandbox is listening there. Users can pin host ports with
`sbx ports --publish <host>:<container>`.

## Environment

| Field       | Description                                    |
| ----------- | ---------------------------------------------- |
| `environment.variables` | Key-value pairs set directly in the container. |

Variable names must be valid shell identifiers (`[A-Za-z_][A-Za-z0-9_]*`).

Do not set `DASH_`, `SBX_`, or `DOCKER_` variables, and avoid overriding
`HOME`, `USER`, `SHELL`, `PATH`, `LD_PRELOAD`, and `LD_LIBRARY_PATH`. The
runtime reserves these names and may override them.

## Setup

`setup.install`, `setup.startup`, and `setup.files` are lists of commands or
files with the fields described here.

### Execution order

When a sandbox is created, kit content is applied in this order:

1. Network permissions and environment variables.
2. Static files under `files/home/`.
3. `setup.install` commands, in declaration order.
4. `setup.files` entries.
5. `setup.startup` commands are registered for each sandbox start.
6. Static files under `files/workspace/`, after the workspace is ready. With
   `--clone`, this means after the repository has been cloned.

For stacked kits, entries in each stage are applied in `--kit` order. An install
command can consume a bundled file from `files/home/`, but not one from
`files/workspace/` or `setup.files`, because those files land later.

`sbx kit add` recreates the sandbox rather than modifying it in place. It
supports mixin kits limited to
`environment.variables`, `setup.install`, and `permissions.network.allow`,
which follow the same order as sandbox creation. It rejects a kit that declares
static files, `setup.startup`, or `setup.files`. To use those fields, recreate
the sandbox with the kit.

### install

Runs synchronously when a kit is applied, either during sandbox creation or
through `sbx kit add`. Shell strings are passed to `sh -c`.

Kit install commands start in the template image's configured `WORKDIR`.
Docker-provided templates use `/home/agent/workspace`, which isn't necessarily
the primary workspace in a direct-mounted or clone-mode sandbox. Don't rely on
the current directory to locate workspace files. Use absolute paths for bundled
assets from `files/home/`.

| Field         | Default | Description                   |
| ------------- | ------- | ----------------------------- |
| `command`     | —       | Shell command string.         |
| `user`        | `"0"`   | User to run as. `"0"` = root. |
| `description` | —       | Human-readable description.   |

### startup

Runs at every sandbox start. String array, not interpreted by a shell.

| Field         | Default  | Description                         |
| ------------- | -------- | ----------------------------------- |
| `command`     | —        | Command and args as a string array. |
| `user`        | `"1000"` | User to run as. `"1000"` = agent.   |
| `background`  | `false`  | Block later startup commands until this command finishes. Set to `true` to let later commands run without waiting. |
| `description` | —        | Human-readable description.         |

Startup commands are non-interactive. They run before the agent
attaches, with no terminal connected, so they can't prompt the user
(for example, an interactive `aws login` will hang or fail). They also
don't gate the agent's entrypoint: the agent launches once startup
commands have been dispatched, regardless of `background`. A value of
`false` waits within the startup dispatcher before it runs the next command;
it doesn't delay the agent entrypoint. Use startup commands
for work that can run alongside the agent. Use `setup.files` for any value that
needs to land on disk before the agent runs.

Startup commands must be idempotent. They run on every sandbox start
and replay on container restarts, so a command that fails or
misbehaves on a second invocation breaks the restart path. Guard
work with existence checks, use upserts instead of inserts, and
prefer commands that converge to the same end state regardless of
how many times they run.

### files

Files written at sandbox start, with runtime substitution.

| Field           | Default  | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `path`          | —        | Absolute container path.                                  |
| `content`       | —        | File content. `${WORKDIR}` expands to the workspace path. |
| `mode`          | `"0644"` | File permissions in octal.                                |
| `onlyIfMissing` | `false`  | Skip if the file already exists.                          |

The runtime writes these files as the agent user with UID 1000. The target
path must be writable by that user. To write to a root-owned path such as
`/etc`, use an `install` command, which runs as root by default. Set ownership
in the install command if the agent needs to modify the file later.

### Shell initialization and service logs

With Docker templates, append shell initialization to
`/etc/sandbox-persistent.sh` in an install command. Keep existing content and
omit completion scripts: interactive and non-interactive Bash commands source
this file. For a background service, redirect startup output to a file and
read it with `sbx exec`. Use `background: true` instead of a trailing `&`.

## Static files

```text
my-kit/files/
├── home/       → /home/agent/
└── workspace/  → primary workspace path
```

| Kit path           | Container destination                   |
| ------------------ | --------------------------------------- |
| `files/home/`      | `/home/agent/` (config files, dotfiles) |
| `files/workspace/` | The primary workspace path              |

Parent directories are created automatically. Existing files are
overwritten. Absolute paths and path-traversal sequences (`../../`) are
rejected.

Static files can supply linter settings, helper scripts, or agent skills.
For example, a Claude Code project skill belongs at
`files/workspace/.claude/skills/<NAME>/SKILL.md`.

## Volumes

Declare `volumes` as a list of mounts with these fields:

| Field  | Description                                                         |
| ------ | ------------------------------------------------------------------- |
| `path` | Required absolute container path.                                   |
| `type` | Empty for a block-backed volume, or `tmpfs` for RAM-backed storage. |
| `size` | Optional byte-size string.                                          |
| `mode` | Optional octal permissions.                                         |

Volumes are applied only when a sandbox is created. `sbx kit add` cannot attach
volumes to a running container.

## Fork an existing agent

Sandbox kits (`kind: sandbox`) define a full agent from scratch. The most
common variant is a fork of a built-in agent. Use `extends:` to inherit the
parent's complete configuration and declare only the fields you want to change.
This example replaces the built-in `claude` entrypoint so Claude Code uses
manual permission mode instead of bypassing approval prompts:

```yaml {title="claude-safe/spec.yaml"}
schemaVersion: "2"
kind: sandbox
name: claude-safe
displayName: Claude Code (with approval prompts)
description: Claude Code in manual permission mode

extends: claude

sandbox:
  entrypoint: [claude, "--permission-mode", "manual"]
```

The child inherits the built-in image, credentials, network permissions,
persistent volumes, settings, MCP integration, agent instructions, setup
entries, and environment variables. Its `sandbox.entrypoint` replaces the
inherited entrypoint.

Launch by passing the sandbox kit in place of a built-in agent name:

```console
$ sbx run ./claude-safe
```

## Install an internal CA certificate

Put each PEM-encoded root certificate under `files/home/` with a `.crt`
extension. For `files/home/internal-ca.crt`, use:

```yaml {title="internal-ca/spec.yaml"}
schemaVersion: "2"
kind: mixin
name: internal-ca
setup:
  install:
    - command: "install -m 0644 /home/agent/internal-ca.crt /usr/local/share/ca-certificates/internal-ca.crt && update-ca-certificates"
      user: "0"
```

This updates the system trust store. For several CAs, install every
certificate before running `update-ca-certificates`.

## Sandbox-managed agent configuration

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

Use separate settings files when supported: Claude Code accepts `--settings`,
and OpenCode reads `OPENCODE_CONFIG`. Don't use `setup.startup` for settings
the agent must read during initialization; startup commands don't gate the
entrypoint.

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
