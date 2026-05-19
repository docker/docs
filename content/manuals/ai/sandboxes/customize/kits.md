---
title: Kits
description: Extend a sandbox with tools, credentials, network rules, and configuration using declarative YAML artifacts.
keywords: sandboxes, sbx, kits, mixins, customization, extensions, agents
weight: 20
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

> [!NOTE]
> Kits are experimental. The kit file format, CLI commands, and experience
> for creating, loading, and managing kits are subject to change as the
> feature evolves. Share feedback and bug reports in the
> [docker/sbx-releases](https://github.com/docker/sbx-releases) repository.

A kit packages a set of capabilities a sandbox can use, such as:

- Tools to install
- Environment variables to set
- Credentials to inject
- Network rules to allow or deny domains
- Files to drop in
- Startup commands to run

You declare these in a single `spec.yaml` file, point the CLI at the
directory (or a ZIP, OCI artifact, or Git URL), and the sandbox applies
and enforces them at runtime. Credentials stay on the host and go through
a proxy instead of entering the VM, and outbound traffic is restricted to
the domains permitted by the kit's network rules.

A kit is either a mixin or an agent:

- Mixin kits (`kind: mixin`) extend an existing agent with extra
  capabilities. Stack several on the same sandbox.
- Agent kits (`kind: agent`) define a full agent from scratch: its image,
  entrypoint, network policies, and everything else the agent needs.

## What kits can do

### Run commands

A kit can run commands inside the sandbox automatically. **Install
commands** run once at creation; **startup commands** run each time
the sandbox starts.

Install commands are the place to put anything an agent needs into the
image, via `apt`, `pip`, `npm`, `curl | bash`, or whatever fits:

```yaml
commands:
  install:
    - command: "apt-get update && apt-get install -y jq"
```

Startup commands cover things like launching background services,
warming caches, or refreshing config on each start. They must be
idempotent — see the [`startup`](#startup) spec reference:

```yaml
commands:
  startup:
    - command: ["sh", "-c", "my-daemon &"]
```

### Inject files

Kits can inject files into the sandbox in two ways: **static files** bundled
with the kit, and **`initFiles`** written at startup with runtime values
substituted in.

Static files work well for content that doesn't vary between sandboxes, such
as tool configurations, shared linter rules, helper scripts the agent can
invoke, or reference material like a style guide or API cheatsheet.

```text
my-kit/
├── spec.yaml
└── files/
    ├── home/
    │   └── .config/my-tool/settings.json
    └── workspace/
        └── .editorconfig
```

`initFiles` cover content that depends on runtime values, such as an
absolute workspace path that a tool needs to bake into its config file
at startup:

```yaml
commands:
  initFiles:
    - path: /home/agent/.my-tool/config.json
      content: '{"workspace": "${WORKDIR}"}'
      onlyIfMissing: true
```

See [`initFiles`](#initfiles) in the spec reference for all fields.

### Set environment variables

Environment variables set by the kit are available to the agent at
runtime:

```yaml
environment:
  variables:
    MY_TOOL_WORKSPACE: /home/agent/my-tool
```

For credentials, see
[Authenticate to external services](#authenticate-to-external-services).
Don't put secret values directly in `environment.variables` — they'd
be visible inside the sandbox VM.

### Control network access

Network rules define which domains the sandbox can reach or block. Kit
network rules apply only to sandboxes that use the kit:

```yaml
network:
  allowedDomains:
    - api.example.com
    - "*.cdn.example.com"
  deniedDomains:
    - telemetry.example.com
```

Use `allowedDomains` for hosts the agent needs, such as package
registries, install endpoints, or external APIs. Use `deniedDomains` for
hosts the agent should not reach, such as telemetry endpoints. If a domain
matches both an allow rule and a deny rule, the deny rule wins.

For authenticated services, see
[Authenticate to external services](#authenticate-to-external-services).

### Authenticate to external services

A kit can attach credentials to outbound requests through the
host-side proxy. The agent inside the VM works with a sentinel value;
the proxy reads the real credential on the host and overwrites the
auth header before the request leaves the sandbox.

The standard pattern uses four blocks tied to a service identifier
you choose (here, `my-service`):

```yaml
network:
  allowedDomains:
    - api.example.com
  serviceDomains:
    api.example.com: my-service # Tag traffic to this domain
  serviceAuth:
    my-service:
      headerName: Authorization # Overwrite this header
      valueFormat: "Bearer %s"

credentials:
  sources:
    my-service:
      env:
        - MY_SERVICE_API_KEY # Host-side credential lookup

environment:
  proxyManaged:
    - MY_SERVICE_API_KEY # Set the in-VM env var to "proxy-managed"
```

The agent boots with `MY_SERVICE_API_KEY=proxy-managed`, sends a
request with that value in `Authorization`, and the proxy overwrites
the header with the real credential before forwarding. The real
secret never enters the VM.

See [Credentials](../security/credentials.md) for how to provide the
credential value on your host, other approaches for cases the example
above doesn't fit, and what the proxy does at request time.

### Define an agent

Agent kits declare an `agent:` block with the image the agent runs in and
the command the user attaches to when they launch the sandbox:

```yaml
agent:
  image: "my-registry/my-agent:latest"
  entrypoint:
    run: [my-agent, "--yolo"]
```

See [Agent kits](#agent-kits) for use cases and an example.

## Mixin kits

A mixin kit extends an existing agent with extra capabilities. Common use
cases:

- Pre-install tools: linters, libraries, or other custom programs
- Grant the agent access to a new authenticated service (a database, a
  vendor API)
- Inject shared team config (linter rules, editor settings, dotfiles)

### Example: Python linting kit

This kit installs [Ruff](https://docs.astral.sh/ruff/) and injects a shared
configuration file, so every sandbox starts with the same linting setup.

```text
ruff-lint/
├── spec.yaml
└── files/
    └── workspace/
        └── ruff.toml
```

```yaml {title="ruff-lint/spec.yaml"}
schemaVersion: "1"
kind: mixin
name: ruff-lint
displayName: Ruff Linter
description: Python linting with shared team config

network:
  allowedDomains:
    - pypi.org
    - files.pythonhosted.org

commands:
  install:
    - command: "uv tool install ruff@latest"
      user: "1000"
      description: Install Ruff
```

```toml {title="ruff-lint/files/workspace/ruff.toml"}
line-length = 100

[lint]
select = ["E", "F", "I"]
```

> [!TIP]
> The templates for the built-in agents (`claude`, `codex`, etc) already
> includes `uv`, so this mixin can use it without installing it separately.

To start a new sandbox with this mixin:

```console
$ sbx run claude --kit /path/to/ruff-lint/
```

To apply the mixin to a sandbox that's already running, use
[`sbx kit add`](#local) instead. The `--kit` flag only takes effect when a
sandbox is created.

## Agent kits

An agent kit defines a full agent from scratch — image, entrypoint, and
everything the agent needs. Common use cases:

- Package a custom agent you've built so others can run it
- Ship a team-internal agent with defaults baked in
- Run a fork of an existing agent with your own config
- Prototype a new agent integration

Agent kits declare everything a mixin kit can, plus an
[`agent:` block](#agent-block) that tells the sandbox how to launch the
agent. For a step-by-step walkthrough, see
[Build your own agent kit](build-an-agent.md).

### Example: the built-in `claude` agent

The `claude` agent you get from `sbx run claude` is defined as a kit. Here
is an abbreviated version of its spec, showing how the agent block combines
with network, credentials, environment, and commands:

```yaml {title="claude/spec.yaml"}
schemaVersion: "1"
kind: agent
name: claude
agent:
  image: "docker/sandbox-templates:claude-code-docker"
  aiFilename: CLAUDE.md
  persistence: persistent
  entrypoint:
    run: [claude, "--dangerously-skip-permissions"]

network:
  serviceDomains:
    api.anthropic.com: anthropic
    console.anthropic.com: anthropic
  serviceAuth:
    anthropic:
      headerName: x-api-key
      valueFormat: "%s"
  allowedDomains:
    - "claude.com:443"

credentials:
  sources:
    anthropic:
      env:
        - ANTHROPIC_API_KEY

environment:
  variables:
    IS_SANDBOX: "1"

commands:
  install:
    - command: "curl -fsSL https://claude.ai/install.sh | bash"
      user: "1000"
      description: Install Claude Code
```

## Using kits

Kits can be loaded from a local path (a directory or ZIP file), a Git
repository, or an OCI registry. Pass `--kit` more than once to stack
several kits on the same sandbox.

> [!IMPORTANT]
> `--kit` only takes effect when a sandbox is created. Passing it
> against an existing sandbox name fails with
> `--kit can only be used when creating a new sandbox`. To extend a
> running sandbox with a kit, use [`sbx kit add`](#local) instead.

### Local

Point `--kit` at a directory or ZIP file on disk:

```console
$ sbx run claude --kit ./my-kit/
$ sbx run claude --kit ./my-kit-1.0.zip
```

While iterating on a kit, apply changes to a running sandbox with
`sbx kit add` instead of recreating it:

```console
$ sbx kit add my-sandbox ./my-kit/
```

`kit add` re-runs install commands and re-copies files. Kits can't be
removed from a running sandbox — remove and recreate it to start clean.

### Git repository

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

```console
$ sbx run claude --kit ghcr.io/myorg/my-kit:1.0
```

For Docker Hub, include the full `docker.io` prefix. See
[Packaging and distribution](#packaging-and-distribution) for publishing.

> [!IMPORTANT]
> Private kits are only supported on Docker Hub. `sbx` reuses your
> `sbx login` session to pull private artifacts from Docker Hub. Other
> registries are pulled anonymously, so private kits hosted on
> registries other than Docker Hub fail to pull.

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

## Spec reference

A kit directory has a required `spec.yaml` and an optional `files/` tree:

```text
my-kit/
├── spec.yaml       # required
└── files/          # optional — static files to inject
    ├── home/
    └── workspace/
```

### Top-level fields

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

### Credentials

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
| `sources.<id>.file.path`   | Path on host. `~` expands to home.                            |
| `sources.<id>.file.parser` | How to extract the value (for example, `"json:apiKey"`).      |
| `sources.<id>.priority`    | `env-first` (default) or `file-first`.                        |

Service identifiers link credentials to [network rules](#network).

### Network

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

| Field                     | Description                                                      |
| ------------------------- | ---------------------------------------------------------------- |
| `allowedDomains`          | Domains the sandbox can reach. Wildcards supported.              |
| `deniedDomains`           | Domains the sandbox can't reach. Deny rules take precedence.     |
| `serviceDomains`          | Map of domain to service identifier from `credentials.sources`.  |
| `serviceAuth.headerName`  | HTTP header the proxy sets (for example, `Authorization`).       |
| `serviceAuth.valueFormat` | Format string for the header value (for example, `"Bearer %s"`). |

### Environment

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

### Commands

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

#### `install`

Runs once during sandbox creation. Shell strings passed to `sh -c`.

| Field         | Default | Description                   |
| ------------- | ------- | ----------------------------- |
| `command`     | —       | Shell command string.         |
| `user`        | `"0"`   | User to run as. `"0"` = root. |
| `description` | —       | Human-readable description.   |

#### `startup`

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

#### `initFiles`

Files written at sandbox start, with runtime substitution.

| Field           | Default  | Description                                               |
| --------------- | -------- | --------------------------------------------------------- |
| `path`          | —        | Absolute container path.                                  |
| `content`       | —        | File content. `${WORKDIR}` expands to the workspace path. |
| `mode`          | `"0644"` | File permissions in octal.                                |
| `onlyIfMissing` | `false`  | Skip if the file already exists.                          |

### Static files

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

### Agent block

Required for `kind: agent`.

```yaml
agent:
  image: <image-ref>
  aiFilename: <filename>
  persistence: <persistent | ephemeral>
  entrypoint:
    run: [<argv>, ...]
    args: [<arg>, ...]
```

| Field                   | Required | Description                                                                                    |
| ----------------------- | -------- | ---------------------------------------------------------------------------------------------- |
| `agent.image`           | Yes      | Docker image reference. See [Base image requirements](#base-image-requirements).               |
| `agent.aiFilename`      | No       | Memory filename (for example, `AGENTS.md`). Appends top-level [`memory`](#memory) at creation. |
| `agent.persistence`     | No       | `persistent` (named volume across restarts) or `ephemeral` (default).                          |
| `agent.entrypoint.run`  | No       | Command and args as a string array. Replaces the image's entrypoint.                           |
| `agent.entrypoint.args` | No       | Args appended to the image's existing entrypoint.                                              |

#### Base image requirements

The agent's container image must provide:

- A non-root `agent` user at UID 1000 with passwordless sudo.
- A `/home/agent/` home directory owned by `agent`.
- HTTP proxy environment variables (`HTTP_PROXY`, `HTTPS_PROXY`,
  `NO_PROXY`) preserved across sudo.
- The agent binary (baked in, or installed via
  [`commands.install`](#commands)).

Build on top of `docker/sandbox-templates:shell-docker` to get these for
free.

#### Memory

```yaml
memory: |
  <markdown>
```

Top-level field. Markdown appended to the agent's memory file at sandbox
creation. The agent reads this content at startup, so write it as
instructions or notes the agent should follow when working in the
sandbox. Applied only when `agent.aiFilename` is set.

The file is written to the parent of the workspace path inside the
sandbox, not to the workspace itself. For a workspace mounted at
`/Users/you/myproject`, the memory file lands at
`/Users/you/AGENTS.md` (or whatever `aiFilename` is set to). It exists
only inside the sandbox — nothing is written to the host.

## Debugging

When a kit doesn't behave as expected, start with the network policy log
and direct inspection inside the sandbox:

- `sbx policy log` shows every outbound request the sandbox proxy saw,
  the rule it matched, extra context when available, and its `PROXY`
  value, such as `forward`, `forward-bypass`, `transparent`, or
  `browser-open`. Use it to diagnose install-time download failures,
  blocked domains, and unexpected TLS interception. If downloads fail or
  arrive corrupted after you add `serviceDomains`, check whether the
  service mapping is too broad. Map only the hosts that need credential
  injection.
- `sbx exec <sandbox> -- <cmd>` runs an arbitrary command inside an
  existing sandbox. Useful for inspecting post-install state without
  recreating: `which mytool`, `ls /home/agent/.local/bin/`,
  `cat /home/agent/.config/...`, and so on.

Install and startup command output is only emitted during `sbx run` or
`sbx create`; `sbx` doesn't retain it for later inspection. To repeat
setup with fresh output, remove and recreate the sandbox:
`sbx rm <sandbox> && sbx run ...`.
