---
title: Kit spec reference
linkTitle: Spec reference
description: Field-by-field reference for a kit's spec.yaml — credentials, network rules, environment, commands, files, agent context, and the sandbox block.
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
recommended; `"1"` is still accepted. Both are parsed and validated the same way,
and a v1 spec is automatically normalized into the v2 model — so existing v1 kits
keep working unchanged.

You don't have to migrate a kit all at once. Field validity isn't tied to the
version — you can adopt v2 fields incrementally, or even mix v1 and v2 spellings
in the same file, and `sbx kit validate` reports a deprecation warning naming the
v2 replacement for each legacy field.

What changed in v2:

| v1                                         | v2                                        |
| ------------------------------------------ | ----------------------------------------- |
| `credentials.sources.<id>`                 | `credentials:` list entry with `service`  |
| `network.allowedDomains` / `deniedDomains` | `caps.network.allow` / `deny`             |
| `network.serviceDomains` / `serviceAuth`   | `credentials[].apiKey.inject`             |
| standalone `oauth:` block                  | `credentials[].oauth`                     |
| `environment.proxyManaged`                 | automatic — `credentials[].apiKey.name`   |
| `memory`                                   | `agentContext`                            |
| `kind: agent` / `agent:` block             | `kind: sandbox` / `sandbox:` block        |
| `tmpfs:`                                    | `volumes:` entries with `type: tmpfs`     |
| `settings:` / `kitDir` / `persistence`     | removed (no replacement)                  |

Credential **discovery** also moved out of the kit in v2: a kit declares which
credentials it needs and how to inject them, but where each value comes from is
controlled by the user through
[credential bindings](../security/credentials.md#credential-bindings).

> [!NOTE]
> `mixins` and `sandbox.build` are accepted by the parser but not yet applied by
> the runtime — `sbx kit validate` reports them as accepted but not yet
> implemented. Don't rely on them yet.

## Changelog

Renamed fields are still accepted for backward compatibility, but
`sbx kit validate` reports a deprecation warning for each, and a future
release may stop accepting them. Update kits to the current names. For the full
v1-to-v2 field mapping, see [Schema versions](#schema-versions).

### v0.32.0

The following `spec.yaml` fields were renamed:

| Previous       | Current          |
| -------------- | ---------------- |
| `memory`       | `agentContext`   |
| `kind: agent`  | `kind: sandbox`  |
| `agent:` block | `sandbox:` block |

The per-kit directory was also renamed from `kits-memory/` to
`kits-agent-context/`. An existing `kits-memory/` directory is migrated
automatically the next time the sandbox starts.

## Top-level fields

```yaml
schemaVersion: "2"
kind: <mixin | sandbox>
name: <name>
displayName: <name>
description: <text>
```

| Field           | Required | Description                                                                |
| --------------- | -------- | -------------------------------------------------------------------------- |
| `schemaVersion` | Yes      | Spec schema version. Use `"2"`; `"1"` is still accepted. See [Schema versions](#schema-versions). |
| `kind`          | Yes      | `mixin` for kits that extend an agent; `sandbox` for kits that define one. |
| `name`          | Yes      | Unique identifier. Lowercase, alphanumeric, hyphens.                       |
| `displayName`   | No       | Human-readable name.                                                       |
| `description`   | No       | Short description.                                                         |

The sections below apply to both kinds. Sandbox kits also declare a
[`sandbox:` block](#sandbox-block).

## Credentials

A kit declares the credentials it needs and how the proxy injects them into
outbound requests. It does not declare where the value comes from — discovery is
controlled by the user through
[credential bindings](../security/credentials.md#credential-bindings), so a kit
can't read arbitrary host environment variables or files.

```yaml
credentials:
  - service: <service-id>
    description: <text>          # optional
    required: <true | false>     # optional, default false
    apiKey:
      name: <ENV_VAR>
      inject:
        - domain: <domain>
          header: <header>
          format: <format>
          username: <user>       # optional, for HTTP basic auth
    oauth:
      tokenEndpoint:
        host: <host>
        path: <path>
      sentinels:
        accessToken: <sentinel>
        refreshToken: <sentinel>
      credentialFile:
        path: <path>
        template: <json-template>
```

`credentials` is a list; each entry names a `service` and configures one or more
auth mechanisms.

| Field         | Description                                                                                                                          |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `service`     | Credential identifier. Known providers (`anthropic`, `github`, `openai`, `google`, ...) auto-expand to their `apiKey` injection config. |
| `description` | Optional. Shown to the user when approving a [binding](../security/credentials.md#credential-bindings).                              |
| `required`    | If `true`, sandbox creation fails when the credential is unavailable. Default `false`.                                              |
| `apiKey`      | API-key injection (see [apiKey](#apikey)).                                                                                           |
| `oauth`       | OAuth interception (see [oauth](#oauth)).                                                                                            |

For a known provider, `- service: anthropic` is enough — `apiKey.name` and
`inject` are filled in from the provider registry. Custom services must declare
`apiKey.name` and `apiKey.inject` (or `oauth`) themselves.

### apiKey

| Field               | Description                                                                                                                          |
| ------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `name`              | Environment variable set inside the container. The agent sees a sentinel (`proxy-managed`); the proxy injects the real value. Auto-derived for known providers. |
| `inject[].domain`   | Domain to inject the credential into. Must also be allowed in [`caps.network`](#network).                                           |
| `inject[].header`   | HTTP header the proxy sets (for example, `x-api-key`, `Authorization`).                                                             |
| `inject[].format`   | Header value format, with one `%s` placeholder (for example, `"%s"` or `"Bearer %s"`).                                             |
| `inject[].username` | Optional. Use HTTP basic auth with this username instead of a bearer header (for example, `x-access-token` for git over HTTPS).     |

### oauth

For agents that authenticate with OAuth (for example, Claude Code), the proxy
intercepts token responses and replaces real tokens with sentinels, then swaps
the real token back in on outbound requests. The token never enters the sandbox.

| Field                                     | Description                                                                                          |
| ----------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `tokenEndpoint.host` / `path`             | The OAuth token endpoint the proxy intercepts.                                                       |
| `sentinels.accessToken` / `refreshToken`  | Sentinel values written into the container in place of the real tokens.                              |
| `credentialFile.path`                     | Where to write the credential file inside the container (`~` expands).                               |
| `credentialFile.template`                 | JSON template for that file. `{{.AccessToken}}`, `{{.RefreshToken}}`, `{{.ExpiresAt}}`, and `{{.ScopesJSON}}` are substituted at runtime. |

### file.parser

A credential sourced from a file — through a
[credential binding](../security/credentials.md#credential-bindings) `file`
source, or a legacy `credentials.sources` entry — can pull its value from a JSON
field. Omit the parser for plain-text files; set `json:<dot.path>` to extract a
field from a JSON file.

| Value             | Behavior                                                                             |
| ----------------- | ------------------------------------------------------------------------------------ |
| omitted or empty  | Reads the entire file as the credential. Leading and trailing whitespace is trimmed. |
| `json:<dot.path>` | Parses the file as JSON and returns the value at the dot-separated path.             |
| any other value   | Rejected — `unsupported parser: <value>`.                                            |

For `json:` paths, segments are separated by `.` (for example, `json:credentials.github.token`).
Only object keys can be navigated — arrays are not supported and there is no `[0]`-style indexing.
Keys that contain a literal `.` cannot be referenced. The resolved value must be a string, number,
or boolean; numbers and booleans are converted to strings. Objects, arrays, and null are rejected.

Common errors when using `json:` parsers:

| Error message                                 | Cause                                                               |
| --------------------------------------------- | ------------------------------------------------------------------- |
| `field 'X' not found in JSON`                 | The path doesn't exist in the file.                                 |
| `cannot navigate to field 'X': not an object` | A path segment hit a string, array, or scalar instead of an object. |
| `field 'X' is not a string value`             | The resolved value is an object, array, or null.                    |
| `failed to parse JSON: ...`                   | The file is not valid JSON.                                         |

## Network

Network egress is declared under `caps.network`. Credentials no longer carry
their own domain mapping — the proxy injects a credential only into the domains
its [`apiKey.inject`](#apikey) (or provider default) lists, and every domain the
sandbox reaches must be allowed here.

```yaml
caps:
  network:
    allow: [<domain>, ...]
    deny: [<domain>, ...]
```

| Field                | Description                                                                                                     |
| -------------------- | --------------------------------------------------------------------------------------------------------------- |
| `caps.network.allow` | Domains the sandbox can reach. Wildcards supported.                                                             |
| `caps.network.deny`  | Domains the sandbox is blocked from reaching. Deny takes precedence over allow, including across composed kits. |

In v1 this was the `network:` block (`allowedDomains` / `deniedDomains`, plus
`serviceDomains` / `serviceAuth`). Those still parse with a deprecation warning:
domain lists fold into `caps.network`, and `serviceDomains` / `serviceAuth` fold
into [`credentials[].apiKey.inject`](#apikey).

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

In v1, `proxyManaged` listed the variables the proxy populated at request time.
That's now automatic: declaring a credential with `apiKey.name: <VAR>` sets
`<VAR>` to a sentinel in the container and injects the real value at the proxy.
`proxyManaged` still parses with a deprecation warning.

## Commands

```yaml
commands:
  install:
    - command: <shell-string>
      user: <uid>
      description: <text>
  startup:
    - command: [<argv>, ...]
      user: <uid>
      background: <true | false>
      description: <text>
  initFiles:
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
refreshing config — and use `commands.initFiles` for any value that
needs to land on disk before the agent runs.

Startup commands must be idempotent. They run on every sandbox start
and replay on container restarts, so a command that fails or
misbehaves on a second invocation breaks the restart path. Guard
work with existence checks, use upserts instead of inserts, and
prefer commands that converge to the same end state regardless of
how many times they run.

### initFiles

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

## Agent context

```yaml
agentContext: |
  <markdown>
```

Top-level field. Available in both mixin and sandbox kits. Markdown
appended to the agent's memory file at sandbox creation. The agent reads
this content at startup. Write it as instructions or notes the agent
should follow when working in the sandbox. Applied only when the active
sandbox kit sets [`sandbox.aiFilename`](#sandbox-block).

The file is written to the parent of the workspace path inside the
sandbox, not to the workspace itself. For a workspace mounted at
`/Users/you/myproject`, the memory file lands at
`/Users/you/AGENTS.md` (or whatever `aiFilename` is set to). It exists
only inside the sandbox. Nothing is written to the host.

When several loaded kits declare `agentContext:` blocks, the content is split
across files instead of being concatenated into the main one:

- Each kit's agent context is written to `<kit-name>.md` in a sibling
  `kits-agent-context/` directory next to the main memory file.
- The main memory file gets a `## Kits` section listing every kit with
  a pointer to its file. The section is delimited by
  `<!-- sbx:kits-section start -->` and `<!-- sbx:kits-section end -->`
  markers so it can be regenerated when kits are added or removed.

## Sandbox block

Required for `kind: sandbox`.

```yaml
sandbox:
  image: <image-ref>
  aiFilename: <filename>
  entrypoint:
    run: [<argv>, ...]
    args: [<arg>, ...]
```

| Field                     | Required | Description                                                                                                 |
| ------------------------- | -------- | ----------------------------------------------------------------------------------------------------------- |
| `sandbox.image`           | Yes      | Docker image reference. See [Base image requirements](#base-image-requirements).                            |
| `sandbox.aiFilename`      | No       | Memory filename (for example, `AGENTS.md`). Appends top-level [`agentContext`](#agent-context) at creation. |
| `sandbox.entrypoint.run`  | No       | Command and args as a string array. Replaces the image's entrypoint.                                        |
| `sandbox.entrypoint.args` | No       | Args appended to the image's existing entrypoint.                                                           |

### Base image requirements

The agent's container image must provide:

- A non-root `agent` user at UID 1000 with passwordless sudo.
- A `/home/agent/` home directory owned by `agent`.
- HTTP proxy environment variables (`HTTP_PROXY`, `HTTPS_PROXY`,
  `NO_PROXY`) preserved across sudo.
- The agent binary (baked in, or installed via
  [`commands.install`](#commands)).

Build on top of `docker/sandbox-templates:shell-docker` to get these for
free.
