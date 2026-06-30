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
- Memory instructions to give the agent

You declare these in a single `spec.yaml` file, point the CLI at the
directory (or a ZIP, OCI artifact, or Git URL), and the sandbox applies
and enforces them at runtime. Credentials stay on the host and go through
a proxy instead of entering the VM, and outbound traffic is restricted to
the domains permitted by the kit's network rules.

A kit is either a mixin or a sandbox:

- Mixin kits (`kind: mixin`) extend an existing agent with extra
  capabilities. Stack several on the same sandbox.
- Sandbox kits (`kind: sandbox`) define a full agent from scratch: its image,
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
idempotent — see the [`startup`](kit-reference.md#startup) spec reference:

```yaml
commands:
  startup:
    - command: ["my-daemon"]
      background: true
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

See [`initFiles`](kit-reference.md#initfiles) in the spec reference for all fields.

Sandboxes seed settings files for some built-in agents during setup.
For example, the sandbox writes `/home/agent/.claude/settings.json`
for the `claude` agent. This happens after the kit's static files and
`initFiles`, so kit-injected files at those paths get overwritten.
Workspace files (such as `<workspace>/.claude/settings.local.json`)
aren't affected, and you can ship them under `files/workspace/` as
usual. To override a path the sandbox writes to, use a
[`commands.startup`](kit-reference.md#startup) script instead. See
[Override agent settings](kit-examples.md#override-agent-settings) for
an example.

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

> [!IMPORTANT]
> The sandbox manages proxy settings for you. It sets `HTTP_PROXY`,
> `HTTPS_PROXY`, `NO_PROXY`, and their lowercase equivalents automatically so
> that traffic flows through its built-in forward proxy, which enforces
> network policy and injects credentials. Leave these variables to the
> sandbox — setting them in a kit points traffic away from the forward proxy,
> so it can no longer apply network policy or inject credentials, and those
> requests typically fail to connect. To send sandbox traffic through an
> upstream corporate proxy, configure it on the host. See
> [Upstream proxy](../architecture.md#upstream-proxy).

### Control network access

Network rules define which domains the sandbox can reach or block. Kit
network rules apply only to sandboxes that use the kit:

```yaml
caps:
  network:
    allow:
      - api.example.com
      - "*.cdn.example.com"
    deny:
      - telemetry.example.com
```

Use `allow` for hosts the agent needs, such as package
registries, install endpoints, or external APIs. Use `deny` for
hosts the agent should not reach, such as telemetry endpoints. If a domain
matches both an allow rule and a deny rule, the deny rule wins.

> [!IMPORTANT]
> Kit network rules don't apply when organization governance is active. In
> that case, only organization rules are evaluated, so kit-defined allow and
> deny rules are ignored — including any domains a kit allows for the agent
> to reach. For details, see
> [Policy precedence](../governance/concepts.md#precedence).

For authenticated services, see
[Authenticate to external services](#authenticate-to-external-services).

### Authenticate to external services

A kit can attach credentials to outbound requests through the
host-side proxy. The agent inside the VM works with a sentinel value;
the proxy reads the real credential on the host and overwrites the
auth header before the request leaves the sandbox.

A kit declares the service, the in-container environment variable, and how
to inject the credential. It does not declare where the value comes from —
that's the user's
[credential binding](../security/credentials.md):

```yaml
credentials:
  - service: my-service
    apiKey:
      name: MY_SERVICE_API_KEY # in-VM env var, set to a sentinel
      inject:
        - domain: api.example.com # inject on requests to this domain
          header: Authorization # overwrite this header
          format: "Bearer %s"

caps:
  network:
    allow:
      - api.example.com # the domain must also be reachable
```

The agent boots with `MY_SERVICE_API_KEY=proxy-managed`, sends a
request with that sentinel in `Authorization`, and the proxy overwrites
the header with the real credential before forwarding. The real
secret never enters the VM.

See [Credentials](../security/credentials.md) for how to provide the
credential value on your host, other approaches for cases the example
above doesn't fit, and what the proxy does at request time. To scope where
a kit-declared credential is sourced or which domains it's injected into,
see [Credential bindings](../security/credentials.md).

### Inject agent memory

A kit can append content to the agent's memory file, such as `CLAUDE.md`
or `AGENTS.md`. The agent reads this file at startup. Use it to give
the agent project conventions, usage tips for a tool the kit installs,
or other guidance that should be in scope when the sandbox runs.

```yaml
agentContext: |
  Ruff is installed. Run `ruff check` before committing.
  Shared config lives at `/workspace/ruff.toml`.
```

Both mixin and sandbox kits can declare `agentContext:`. The content is written
only when the active sandbox kit sets [`sandbox.aiFilename`](kit-reference.md#sandbox-block),
which determines the memory file's name.

When more than one loaded kit declares an `agentContext:` block, each kit's
content is written to its own `<kit-name>.md` file under a sibling
`kits-agent-context/` directory. The main memory file gets a `## Kits`
section that points to each kit file:

```text
/Users/you/
├── myproject/              # workspace
├── AGENTS.md               # main memory file with a "## Kits" index
└── kits-agent-context/
    ├── ruff-lint.md
    ├── vale.md
    └── git-ssh-sign.md
```

See [`agentContext`](kit-reference.md#agent-context) in the spec reference for the full field schema.

### Define an agent

Sandbox kits declare a `sandbox:` block with the image the agent runs in and
the command the user attaches to when they launch the sandbox:

```yaml
sandbox:
  image: "my-registry/my-agent:latest"
  entrypoint:
    run: [my-agent, "--yolo"]
```

See [Sandbox kits](#sandbox-kits) for use cases and an example.

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
schemaVersion: "2"
kind: mixin
name: ruff-lint
displayName: Ruff Linter
description: Python linting with shared team config

caps:
  network:
    allow:
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
> The templates for the built-in agents (`claude`, `codex`, and so on)
> already include `uv`, so this mixin can use it without installing it
> separately.

To start a new sandbox with this mixin:

```console
$ sbx run claude --kit /path/to/ruff-lint/
```

To apply the mixin to a sandbox that's already running, use
[`sbx kit add`](#local) instead. The `--kit` flag only takes effect when a
sandbox is created.

## Sandbox kits

A sandbox kit defines a full agent from scratch — image, entrypoint, and
everything the agent needs. Common use cases:

- Package a custom agent you've built so others can run it
- Ship a team-internal agent with defaults baked in
- Run a fork of an existing agent with your own config
- Prototype a new agent integration

Sandbox kits declare everything a mixin kit can, plus an
[`sandbox:` block](kit-reference.md#sandbox-block) that tells the sandbox how to launch the
agent. For a step-by-step walkthrough, see
[Build your own agent kit](build-an-agent.md).

### Example: the built-in `claude` agent

The `claude` agent you get from `sbx run claude` is defined as a kit. Here
is an abbreviated version of its spec, showing how the sandbox block combines
with network, credentials, environment, and commands:

```yaml {title="claude/spec.yaml"}
schemaVersion: "2"
kind: sandbox
name: claude
sandbox:
  image: "docker/sandbox-templates:claude-code-docker"
  aiFilename: CLAUDE.md
  entrypoint:
    run: [claude, "--dangerously-skip-permissions"]

caps:
  network:
    allow:
      - "claude.com:443"

credentials:
  - service: anthropic
    apiKey:
      name: ANTHROPIC_API_KEY
      inject:
        - domain: api.anthropic.com
          header: x-api-key
          format: "%s"
        - domain: console.anthropic.com
          header: x-api-key
          format: "%s"

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
> For Docker Hub, `sbx` reuses your `sbx login` session to pull private
> kits. For other registries, store pull credentials with
> [`sbx secret set --registry`](../security/credentials.md#registry-credentials)
> before running the sandbox:
>
> ```console
> $ gh auth token | sbx secret set --registry ghcr.io --password-stdin
> ```
>
> Without stored credentials, pulls from non-Docker Hub registries are
> anonymous and private kits fail to pull.

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

`sbx kit pull` prefers credentials stored with
[`sbx secret set --registry`](../security/credentials.md#registry-credentials),
falling back to the Docker credential store. `sbx kit push` only uses the
Docker credential store, so pushing to a private registry requires a prior
`docker login`.

## Spec reference

For a field-by-field reference of every `spec.yaml` block — top-level
fields, credentials, network, environment, commands, static files,
agent context, and the sandbox block — see [Kit spec reference](kit-reference.md).

## Debugging

When a kit doesn't behave as expected, start with the network policy log
and direct inspection inside the sandbox:

- `sbx policy log` shows every outbound request the sandbox proxy saw,
  the rule it matched, extra context when available, and its `PROXY`
  value, such as `forward`, `forward-bypass`, `transparent`, or
  `browser-open`. Use it to diagnose install-time download failures,
  blocked domains, and unexpected TLS interception. If downloads fail or
  arrive corrupted after you add a credential's `apiKey.inject`, check
  whether an injection domain is too broad. Inject only on the hosts that
  need credentials.
- `sbx exec <sandbox> -- <cmd>` runs an arbitrary command inside an
  existing sandbox. Useful for inspecting post-install state without
  recreating: `which mytool`, `ls /home/agent/.local/bin/`,
  `cat /home/agent/.config/...`, and so on.

Install and startup command output is only emitted during `sbx run` or
`sbx create`; `sbx` doesn't retain it for later inspection. To repeat
setup with fresh output, remove and recreate the sandbox:
`sbx rm <sandbox> && sbx run ...`.
