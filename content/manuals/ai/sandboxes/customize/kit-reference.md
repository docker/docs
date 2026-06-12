---
title: Kit spec reference
linkTitle: Spec reference
description: Field-by-field reference for a kit's spec.yaml — credentials, network rules, environment, commands, files, memory, and the agent block.
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

## Top-level fields

```yaml
schemaVersion: "1"
kind: <mixin | agent>
name: <name>
displayName: <name>
description: <text>
```

| Field           | Required | Description                                                              |
| --------------- | -------- | ------------------------------------------------------------------------ |
| `schemaVersion` | Yes      | Spec schema version. Set to `"1"`.                                       |
| `kind`          | Yes      | `mixin` for kits that extend an agent; `agent` for kits that define one. |
| `name`          | Yes      | Unique identifier. Lowercase, alphanumeric, hyphens.                     |
| `displayName`   | No       | Human-readable name.                                                     |
| `description`   | No       | Short description.                                                       |

The sections below apply to both kinds. Agent kits also declare an
[`agent:` block](#agent-block).

## Credentials

```yaml
credentials:
  sources:
    <service-id>:
      env: [<env-var>, ...]
      file:
        path: <path>
        parser: <parser>
      priority: <env-first | file-first>
```

| Field                      | Description                                                   |
| -------------------------- | ------------------------------------------------------------- |
| `sources`                  | Map of service identifier to credential source.               |
| `sources.<id>.env`         | Environment variables to read on the host, in priority order. |
| `sources.<id>.file.path`   | Path on host. `~` expands to home directory.                  |
| `sources.<id>.file.parser` | How to extract the credential value from the file.            |
| `sources.<id>.priority`    | `env-first` (default) or `file-first`.                        |

Service identifiers link credentials to [network rules](#network).

### file.parser

`file.parser` tells the proxy how to extract a credential from the file at `file.path`.
Omit it for plain-text files; set it to `json:<dot.path>` to extract a field from a JSON file.

| Value             | Behavior                                                                             |
| ----------------- | ------------------------------------------------------------------------------------ |
| omitted or empty  | Reads the entire file as the credential. Leading and trailing whitespace is trimmed. |
| `json:<dot.path>` | Parses the file as JSON and returns the value at the dot-separated path.             |
| any other value   | Rejected — `unsupported parser: <value>`.                                            |

For `json:` paths, segments are separated by `.` (for example, `json:credentials.github.token`).
Only object keys can be navigated — arrays are not supported and there is no `[0]`-style indexing.
Keys that contain a literal `.` cannot be referenced. The resolved value must be a string, number,
or boolean; numbers and booleans are converted to strings. Objects, arrays, and null are rejected.

When a source has both `env` and `file` defined, `priority` controls which is tried first. The
preferred source is used when it exists — the environment variable is set, or the file is
present on disk. If it doesn't, the other source is used instead. The choice is made once at
discovery time, so parser errors (missing JSON field, wrong value type, invalid JSON) surface
as errors rather than triggering a fallback.

Plain-text token file:

```yaml
credentials:
  sources:
    openai:
      file:
        path: "~/.openai/token"
```

Nested JSON field, with an environment variable as fallback:

```yaml
credentials:
  sources:
    github:
      env:
        - GH_TOKEN
      file:
        path: "~/.config/myapp/creds.json"
        parser: "json:credentials.github.token"
      priority: file-first
```

Given `~/.config/myapp/creds.json`:

```json
{
  "credentials": {
    "github": { "token": "ghp_xyz", "expires": "2026-12-31" }
  }
}
```

The proxy resolves the credential to `ghp_xyz`, falling back to `GH_TOKEN` if the file is
missing. If the file exists but the JSON path doesn't resolve, the request fails with the
parser error below instead of falling back.

Common errors when using `json:` parsers:

| Error message                                 | Cause                                                               |
| --------------------------------------------- | ------------------------------------------------------------------- |
| `field 'X' not found in JSON`                 | The path doesn't exist in the file.                                 |
| `cannot navigate to field 'X': not an object` | A path segment hit a string, array, or scalar instead of an object. |
| `field 'X' is not a string value`             | The resolved value is an object, array, or null.                    |
| `failed to parse JSON: ...`                   | The file is not valid JSON.                                         |

## Network

```yaml
network:
  allowedDomains: [<domain>, ...]
  deniedDomains: [<domain>, ...]
  serviceDomains:
    <domain>: <service-id>
  serviceAuth:
    <service-id>:
      headerName: <header>
      valueFormat: <format>
```

| Field                     | Description                                                                                                                          |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `allowedDomains`          | Domains the sandbox can reach. Wildcards supported.                                                                                  |
| `deniedDomains`           | Domains the sandbox is blocked from reaching. Deny rules take precedence over allow rules, including those from other composed kits. |
| `serviceDomains`          | Map of domain to service identifier from `credentials.sources`.                                                                      |
| `serviceAuth.headerName`  | HTTP header the proxy sets (for example, `Authorization`).                                                                           |
| `serviceAuth.valueFormat` | Format string for the header value (for example, `"Bearer %s"`).                                                                     |

## Environment

```yaml
environment:
  variables:
    <NAME>: <value>
  proxyManaged: [<NAME>, ...]
```

| Field          | Description                                                                                                         |
| -------------- | ------------------------------------------------------------------------------------------------------------------- |
| `variables`    | Key-value pairs set directly in the container.                                                                      |
| `proxyManaged` | Environment variable names populated by the proxy at request time. Pair with [`credentials.sources`](#credentials). |

Variable names must be valid shell identifiers (`[A-Za-z_][A-Za-z0-9_]*`).

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

## Memory

```yaml
memory: |
  <markdown>
```

Top-level field. Available in both mixin and agent kits. Markdown
appended to the agent's memory file at sandbox creation. The agent reads
this content at startup. Write it as instructions or notes the agent
should follow when working in the sandbox. Applied only when the active
agent kit sets [`agent.aiFilename`](#agent-block).

The file is written to the parent of the workspace path inside the
sandbox, not to the workspace itself. For a workspace mounted at
`/Users/you/myproject`, the memory file lands at
`/Users/you/AGENTS.md` (or whatever `aiFilename` is set to). It exists
only inside the sandbox. Nothing is written to the host.

When several loaded kits declare `memory:` blocks, the content is split
across files instead of being concatenated into the main one:

- Each kit's memory is written to `<kit-name>.md` in a sibling
  `kits-memory/` directory next to the main memory file.
- The main memory file gets a `## Kits` section listing every kit with
  a pointer to its file. The section is delimited by
  `<!-- sbx:kits-section start -->` and `<!-- sbx:kits-section end -->`
  markers so it can be regenerated when kits are added or removed.

## Agent block

Required for `kind: agent`.

```yaml
agent:
  image: <image-ref>
  aiFilename: <filename>
  entrypoint:
    run: [<argv>, ...]
    args: [<arg>, ...]
```

| Field                   | Required | Description                                                                                    |
| ----------------------- | -------- | ---------------------------------------------------------------------------------------------- |
| `agent.image`           | Yes      | Docker image reference. See [Base image requirements](#base-image-requirements).               |
| `agent.aiFilename`      | No       | Memory filename (for example, `AGENTS.md`). Appends top-level [`memory`](#memory) at creation. |
| `agent.entrypoint.run`  | No       | Command and args as a string array. Replaces the image's entrypoint.                           |
| `agent.entrypoint.args` | No       | Args appended to the image's existing entrypoint.                                              |

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
