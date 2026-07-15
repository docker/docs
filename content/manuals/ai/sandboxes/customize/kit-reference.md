---
title: Kit spec reference
linkTitle: Spec reference
description: Field-by-field reference for a kit's spec.yaml, including credentials, network rules, environment, setup, files, agent instructions, and the sandbox block.
keywords: sandboxes, sbx, kits, spec.yaml, reference, schema, fields
weight: 22
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> Kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change as the
> feature evolves. Share feedback and bug reports in the
> [docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

This page documents every field in a kit's `spec.yaml`. For an overview of
what kits are and how to use them, see [Kits](kits.md).

A kit directory has a required `spec.yaml` and an optional `files/` tree:

```text
my-kit/
├── spec.yaml       # required
└── files/          # optional — static files to inject
    ├── home/
    └── workspace/
```

## Schema versions

Two schema versions are supported. `schemaVersion: "2"` is current and
recommended; `"1"` is still accepted through the legacy path.

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
[credential bindings](../security/credentials.md#credential-bindings).

> [!NOTE]
> `mixins` and `sandbox.build` are accepted by the parser, but runtime support
> is pending. A kit that sets `sandbox.build` must also set `sandbox.image`.

## Top-level fields

```yaml
schemaVersion: "2"
kind: <mixin | sandbox>
name: <name>
version: <version>
displayName: <name>
description: <text>
sourceURL: <url>
licenses:
  - MIT
locked:
  - sandbox.image
security:
  privileged: false
```

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

A kit also declares behavior blocks such as `agentInstructions`,
`permissions`, `ports`, `credentials`, `environment`, `setup`, and `volumes`.

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
sandbox:
  entrypoint: [claude]
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
| `sandbox.image`      | Yes      | Docker image reference. A sandbox that sets `extends:` can inherit this from its parent.                        |
| `sandbox.build`      | No       | Build configuration. Runtime support is pending, so a kit with `build:` must also set `image:`.                 |
| `sandbox.entrypoint` | No       | Fixed process prefix as a string array. The first element is the agent binary.                                  |
| `sandbox.command`    | No       | Mode-specific argument tail. Use a list shorthand for `default`, or a mapping with `default` and `interactive`. |
| `sandbox.resources`  | No       | Optional CPU, memory, and GPU constraints. Memory uses byte-size strings such as `4096m` or `4g`.               |

The effective command is `entrypoint` plus `command.default` for non-interactive
launches, and `entrypoint` plus `command.interactive` for TTY sessions. If
`interactive` is omitted, it falls back to `default`.

The agent's container image must provide:

- A non-root `agent` user at UID 1000 with passwordless sudo.
- A `/home/agent/` home directory owned by `agent`.
- HTTP proxy environment variables (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`) preserved across sudo.
- The agent binary, either baked in or installed with [`setup.install`](#setup).

Build on top of `docker/sandbox-templates:shell-docker` to get these base
requirements.

## Agent instructions

```yaml
agentInstructions:
  filename: CLAUDE.md
  content: |
    Ruff is installed. Run `ruff check` before committing.
```

| Field      | Description                                                                                         |
| ---------- | --------------------------------------------------------------------------------------------------- |
| `filename` | AI profile filename. Meaningful for `kind: sandbox`; ignored with a warning for `kind: mixin`.      |
| `content`  | Markdown instructions. For a sandbox, inlined into the profile. For a mixin, written to kit memory. |

For mixins, the engine writes `content` to
`<dir-of-AI-file>/kits-memory/<kit-name>.md` and adds a `## Kits` pointer
section to the base AI file. This keeps each mixin's instructions in a separate
file.

## Credentials

A kit declares the credentials it needs and how the proxy injects them into
outbound requests. It does not declare where the value comes from — discovery is
controlled by the user through
[credential bindings](../security/credentials.md), so a kit
can't read arbitrary host environment variables or files.

```yaml
credentials:
  - service: <service-id>
    description: <text> # optional
    required: <true | false> # optional, default false
    provider: <provider> # optional, reserved
    apiKey:
      name: <ENV_VAR>
      proxyManaged: true
      inject:
        - domain: <domain>
          header: <header>
          format: <format>
        - domain: <domain>
          scheme: bearer
        - domain: <domain>
          scheme: basic
          username: <user> # optional, for HTTP basic auth
    oauth:
      tokenEndpoint:
        host: <host>
        path: <path>
      sentinels:
        accessToken: <sentinel>
        refreshToken: <sentinel>
      credentialFile:
        path: <path>
        structure:
          service:
            accessToken: "{{.AccessToken}}"
```

`credentials` is a list; each entry names a `service` and configures one or more
auth mechanisms.

| Field         | Description                                                                                                  |
| ------------- | ------------------------------------------------------------------------------------------------------------ |
| `service`     | Credential identifier, matched against the value stored with `sbx secret set`. Lowercase kebab-case.         |
| `description` | Optional. Shown to the user when approving a [binding](../security/credentials.md#credential-bindings).      |
| `required`    | If `true`, the resolver fails fast when the credential has no binding and no host fallback. Default `false`. |
| `provider`    | Reserved for a provider registry. Accepted with a warning and no runtime effect.                             |
| `apiKey`      | API-key injection (see [apiKey](#apikey)).                                                                   |
| `oauth`       | OAuth interception (see [oauth](#oauth)).                                                                    |

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
the real token back in on outbound requests. The token never enters the sandbox.

| Field                                    | Description                                                                                                                                              |
| ---------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `tokenEndpoint.host` / `path`            | The OAuth token endpoint the proxy intercepts.                                                                                                           |
| `sentinels.accessToken` / `refreshToken` | Sentinel values written into the container in place of the real tokens.                                                                                  |
| `credentialFile.path`                    | Where to write the credential file inside the container (`~` expands).                                                                                   |
| `credentialFile.structure`               | Declarative JSON shape for the credential file. `{{.AccessToken}}`, `{{.RefreshToken}}`, `{{.ExpiresAt}}`, and `{{.Scopes}}` are substituted at runtime. |
| `credentialFile.template`                | Deprecated Go template form. If both `structure` and `template` are set, `structure` wins.                                                               |
| `resourceHosts`                          | API hosts where the proxy attaches the token on outbound requests, distinct from the token endpoint host.                                                |
| `skipIfEnv`                              | Environment variable names that, if set on the host, make a stored API key take precedence over the OAuth flow.                                          |
| `responseFields`                         | Overrides the default field names the proxy reads from the token response.                                                                               |
| `passthrough`                            | If `true`, the proxy passes the token response through unchanged instead of replacing the tokens with sentinels.                                         |

## Network

Network egress is declared under `permissions.network`. Credentials no longer carry
their own domain mapping — the proxy injects a credential only into the domains
its [`apiKey.inject`](#apikey) lists, and every domain the
sandbox reaches must be allowed here.

```yaml
permissions:
  network:
    allow: [<domain>, ...]
    deny: [<domain>, ...]
```

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

Use `ports` for inbound service exposure from the sandbox to the host:

```yaml
ports:
  - container: 8080
    protocol: tcp
    name: web
```

| Field       | Description                                                         |
| ----------- | ------------------------------------------------------------------- |
| `container` | Container port, 1 to 65535.                                         |
| `protocol`  | `tcp` or `udp`. Empty means `tcp`.                                  |
| `name`      | Optional label surfaced by tools that list published port bindings. |

Host ports are allocated ephemerally on `127.0.0.1`. Users can pin host ports
with `sbx ports --publish <host>:<container>`.

## Environment

```yaml
environment:
  variables:
    <NAME>: <value>
```

| Field       | Description                                    |
| ----------- | ---------------------------------------------- |
| `variables` | Key-value pairs set directly in the container. |

Variable names must be valid shell identifiers (`[A-Za-z_][A-Za-z0-9_]*`).

Do not set `DASH_`, `SBX_`, or `DOCKER_` variables, and avoid overriding
`HOME`, `USER`, `SHELL`, `PATH`, `LD_PRELOAD`, and `LD_LIBRARY_PATH`. The
runtime reserves these names and may override them.

## Setup

```yaml
setup:
  install:
    - command: <shell-string>
      user: <uid>
      description: <text>
  startup:
    - command: [<argv>, ...]
      user: <uid>
      background: <true | false>
      description: <text>
  files:
    - path: <path>
      content: <text>
      mode: <octal>
      onlyIfMissing: <true | false>
      description: <text>
```

### install

Runs once during sandbox creation. Shell strings passed to `sh -c`.

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
| `background`  | `false`  | Run in background.                  |
| `description` | —        | Human-readable description.         |

Startup commands are non-interactive. They run before the agent
attaches, with no terminal connected, so they can't prompt the user
(for example, an interactive `aws login` will hang or fail). They also
don't gate the agent's entrypoint: the agent launches once startup
commands have been dispatched, regardless of `background`. Use them
for non-interactive prep — launching daemons, warming caches,
refreshing config — and use `setup.files` for any value that
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

## Volumes

```yaml
volumes:
  - path: /workspace
    size: 10g
    mode: "0755"
  - path: /tmp/scratch
    type: tmpfs
    size: 512m
    mode: "1777"
```

| Field  | Description                                                         |
| ------ | ------------------------------------------------------------------- |
| `path` | Required absolute container path.                                   |
| `type` | Empty for a block-backed volume, or `tmpfs` for RAM-backed storage. |
| `size` | Optional byte-size string.                                          |
| `mode` | Optional octal permissions.                                         |

Volumes are applied only when a sandbox is created. `sbx kit add` cannot attach
volumes to a running container.
