---
title: Sandbox environment files
linkTitle: Environment files
weight: 20
description: Use a declarative sbxenv.yaml file to describe and share your sandbox configuration.
keywords:
  - docker sandboxes
  - sbx env
  - sbxenv
  - environment file
  - sandbox configuration
  - declarative
aliases:
  - /ai/sandboxes/sandbox-environments/
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

A sandbox environment file captures the setup for a project in a
`sbxenv.yaml` file. Share the file with project contributors so they use the
same agent, tools, resources, and credentials without reproducing CLI flags and
setup steps.

> [!NOTE]
> `sbx env` is experimental. The command interface and file format may change.

## Start an environment

Keep the environment file outside the directories you mount into the sandbox.
That includes the primary workspace and every `additionalWorkspaces` mount.
For example, place it beside your project:

```text
web-app-env/
├── sbxenv.yaml
└── web-app/
```

Create `web-app-env/sbxenv.yaml`. This example gives the agent a shared
environment variable and the Playwright browser-testing tools. It also
publishes the application's development port:

```yaml
schemaVersion: "1"
name: web-app
agent: claude
workspace: ./web-app

kits:
  - docker.io/sbx/playwright-kit:latest

env:
  NODE_ENV: test

ports:
  - sandbox: 3000
    host: 3000
```

From `web-app-env`, run the environment:

```console
$ sbx env run
```

`sbx` shows an environment plan and asks you to approve it. If you approve the
plan, the `web-app` directory becomes the workspace, while `sbxenv.yaml` remains
outside the sandbox. If the environment doesn't exist, `sbx` creates a sandbox
named `web-app`, installs Playwright and Chromium, and publishes sandbox port
`3000` on the host. It then attaches to the agent. Later runs attach to the
existing sandbox.

This placement keeps the environment file outside the agent's writable
workspace. If you later add `additionalWorkspaces`, keep `sbxenv.yaml`
outside those directories too. See the [`workspace` guidance](#workspace)
for details.

## Commands

| Command                                                                       | Description                                                                          |
| ----------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| `sbx env plan [PATH...]`                                                       | Shows what applying the environment would change without changing or approving it    |
| [`sbx env run`](/reference/cli/sbx/env/run/) `[PATH...]`                      | Applies the approved plan, creates the environment if needed, and attaches            |
| [`sbx env create`](/reference/cli/sbx/env/create/) `[PATH...]`                | Applies the approved plan and creates the environment without attaching              |
| [`sbx env exec`](/reference/cli/sbx/env/exec/) `[PATH...] -- COMMAND [ARG...]` | Runs a command in an existing environment without running lifecycle commands         |
| [`sbx env rm`](/reference/cli/sbx/env/rm/) `[PATH...]`                        | Shows a destroy plan, then removes the sandbox and resources named in the plan        |

`PATH` can be a directory or a direct path to an environment file. A directory
always resolves to `sbxenv.yaml`. To use another filename, pass the file path
explicitly. With no path, `sbx` reads `sbxenv.yaml` from the working directory.

Pass the same set of paths to each `sbx env` command so they resolve the same
sandbox.

### Set user defaults

Create `~/.sbxenv.yaml` to define defaults shared across projects. When you run
an `sbx env` command without a path, `sbx` merges this file beneath the
project's `sbxenv.yaml`. Passing any path skips the user-level file.

The user-level file must not set `name` or `workspace`, because those fields
identify a project. Lists such as `ports` and `mcp.servers` concatenate across
the user and project files.

### Parameterize an environment

Declare inputs in a top-level `args` block when values need to vary between
uses of the same environment file. Each argument must have exactly one of
`default` or `required: true`:

```yaml
schemaVersion: "1"
name: web-app
agent: claude

args:
  channel:
    default: stable
    description: Release channel
    enum:
      - stable
      - beta
  endpoint:
    required: true
    description: API endpoint
  cpus:
    default: "4"
    pattern: "[1-9][0-9]*"

env:
  RELEASE_CHANNEL: ${{ env.args.channel }}
  API_ENDPOINT: ${{ env.args.endpoint }}

sandboxOptions:
  cpus: ${{ env.args.cpus }}
```

Reference a declared argument as `${{ env.args.NAME }}` anywhere a YAML value
can appear. References can't be used in field names or within the `args` block.
An unquoted reference is interpreted as a YAML value after substitution, so
the `cpus` value in this example becomes an integer. Quote a reference to
preserve it as a string.

All `sbx env` commands accept repeatable `--env-arg NAME=VALUE` flags. Values
provided with a flag replace defaults from the environment file:

```console
$ sbx env run --env-arg endpoint=https://api.example.com --env-arg channel=beta
```

Use `--env-args-file` to load values from a file. Each non-empty, non-comment
line must have the form `NAME=VALUE`:

```text
# production.args
channel=beta
endpoint=https://api.example.com
```

```console
$ sbx env run --env-args-file production.args
```

You can pass multiple argument files. Later files take precedence over earlier
files, and `--env-arg` flags take precedence over every argument file. The
command rejects undeclared arguments, missing required values, and values that
don't satisfy an argument's `enum` or `pattern`. Values can contain `=`, and
values in an argument file are read literally rather than expanded by a shell.

Argument references are the only variable expressions expanded in an
environment file. Shell-style expressions such as `${VAR}` aren't expanded
from the host environment. Other dollar signs remain literal, so a value such
as `$PATH:/opt/bin` is passed unchanged. Use `$${{ env.args.NAME }}` to produce
the literal text `${{ env.args.NAME }}`. Substituted values aren't expanded a
second time.

## Common workflows

The following examples combine environment file fields into configurations you
can adapt for a project.

### Combine team defaults and personal settings

Keep the shared configuration in a version-controlled environment directory
outside the mounted workspace. Put machine-specific settings in a file excluded
from version control. For example, commit `base.sbxenv.yaml` beside the
`web-app` directory:

```yaml
schemaVersion: "1"
name: web-app
agent: claude
workspace: ./web-app

env:
  NODE_ENV: development

sandboxOptions:
  cpus: 4
  memory: 8g
```

Add `local.sbxenv.yaml` to `.gitignore`, then use it for personal settings:

```yaml
env:
  LOG_LEVEL: debug

sandboxOptions:
  memory: 12g
```

Pass both files in merge order:

```console
$ sbx env run base.sbxenv.yaml local.sbxenv.yaml
```

Nested mappings merge by key, lists concatenate, and values from later files
replace earlier scalar values. In this example, the sandbox has four CPUs,
12 GB of memory, and both environment variables. The first file controls the
base directory for relative workspace paths and the default sandbox name.

### Work across multiple repositories

Mount related repositories alongside the primary project when the agent needs
to coordinate changes or consult shared code and documentation:

```yaml
# sbxenv.yaml in the directory above the repositories
schemaVersion: "1"
name: web-platform
agent: codex

workspace: ./web-app

additionalWorkspaces:
  - path: ./shared-components
  - path: ./architecture-docs
    readOnly: true
```

The agent starts in `web-app`, can modify `shared-components`, and can read
`architecture-docs` without changing it. The environment file stays outside all
three workspaces. Relative paths resolve from the directory of the first
environment file. Additional workspaces are mounted directly even when the
primary workspace uses clone mode.

### Reuse an environment in automation

Use the same committed environment for interactive development and automated
tasks. Developers attach to the agent with `run`:

```console
$ sbx env run
```

Automation can create the sandbox without attaching, run commands in it, and
remove it afterward:

```console
$ sbx env create --auto-approve
$ sbx env exec -- npm test
$ sbx env rm --force
```

`--auto-approve` approves the plan for that invocation without recording
consent for later invocations. Use the flag for each unattended `create` or
`run`. The `--force` flag approves the destroy plan and removes the sandbox even
when it is in use.

Commands and vault references under `secrets` resolve on the host, so the
automation runner must provide the referenced tools and authentication. The
secret values remain outside the environment file.

## Review an environment plan

`sbx env create`, `sbx env run`, and `sbx env rm` show the changes an
environment makes outside its sandbox and ask for approval before applying
them. The plan includes host commands, credentials, bindings, MCP
registrations, directories, kits, published ports, sandbox options, and
environment variables.

Run `sbx env plan` to inspect the apply plan without changing, approving, or
recording anything:

```console
$ sbx env plan
```

The plan compares the environment file with what the environment last applied
and with resources on the host. It omits resources that are unchanged and
already approved. Literal secret values appear as SHA-256 digests. Secret
references, host commands, environment variables, ports, paths, and binding
domains remain visible so you can review them.

Interactive approval is recorded for the environment under the `sbx` state
directory. Plans without host commands apply silently on later invocations until
the environment changes or a resource is missing. An approval provided with
`--auto-approve` applies only to that invocation.

## Update an environment

For an existing sandbox, `sbx env run` applies changes to `env` to the new
session and reconciles declared MCP servers. The plan marks changes to
workspaces, kits, ports, secrets, bindings, and `sandboxOptions` as waiting for
the next sandbox creation. Remove the environment with `sbx env rm`, then create
it again to apply those changes.

## Remove an environment

Secrets and registry credentials are sandbox-scoped. Credential bindings and
MCP server registrations are host-global and can be shared by multiple
sandboxes.

`sbx env rm` builds a destroy plan from the resources on the host. The plan
includes all credentials stored at the sandbox's scope, including credentials
that the environment file no longer declares. After approval, `sbx` removes
only the resources named in the plan.

Global credential bindings remain unless you pass `--prune-bindings`. MCP
registrations remain available to other sandboxes.

### Clean up after a failed create

Secret provisioning, binding updates, and MCP server registration occur before
the sandbox is created. If sandbox creation fails, scoped secrets remain, and
bindings and MCP registrations may also remain. Run `sbx env rm` with the same
paths to remove the scoped secrets. Pass `--prune-bindings` if you also want to
remove the declared global bindings. MCP registrations are host-global and
remain after cleanup.

## File reference

The loader rejects unknown fields and unsupported schema versions.

### Top-level fields

| Field                  | Type             | Required | Default                        | Description                                                                     |
| ---------------------- | ---------------- | -------- | ------------------------------ | ------------------------------------------------------------------------------- |
| `schemaVersion`        | string           | Yes      | None                           | Schema version. The supported value is `"1"`                                   |
| `name`                 | string           | No       | `<agent>-<workspace-basename>` | Sandbox name                                                                    |
| `agent`                | string           | Yes      | None                           | Built-in agent or the name of an agent kit                                      |
| `args`                 | map              | No       | None                           | Environment arguments. See [`args`](#args)                                      |
| `kits`                 | list             | No       | None                           | Kits to install at creation. See [`kits`](#kits)                                 |
| `workspace`            | string or object | No       | No host mount                  | Primary workspace. See [`workspace`](#workspace)                                |
| `additionalWorkspaces` | list             | No       | None                           | Extra directories to mount. See [`additionalWorkspaces`](#additionalworkspaces) |
| `env`                  | map of strings   | No       | None                           | Environment variables for the sandbox                                           |
| `sandboxOptions`       | object           | No       | None                           | Creation options. See [`sandboxOptions`](#sandboxoptions)                        |
| `secrets`              | map              | No       | None                           | Service credentials. See [`secrets`](#secrets)                                  |
| `bindings`             | map              | No       | None                           | Credential injection approvals. See [`bindings`](#bindings)                     |
| `registries`           | map              | No       | None                           | Registry pull credentials. See [`registries`](#registries)                      |
| `mcp`                  | object           | No       | None                           | MCP servers. See [`mcp`](#mcp)                                                  |
| `ports`                | list             | No       | None                           | Port mappings. See [`ports`](#ports)                                            |
| `lifecycle`            | object           | No       | None                           | Host commands. See [`lifecycle`](#lifecycle)                                    |

### `args`

`args` maps argument names to their declarations. Names must start with a
letter or underscore and can contain letters, numbers, underscores, and
hyphens. Each declaration must set exactly one of `default` or `required:
true`.

| Field         | Type            | Default | Description                                                       |
| ------------- | --------------- | ------- | ----------------------------------------------------------------- |
| `default`     | string          | None    | Value used when the command doesn't supply the argument           |
| `required`    | boolean         | `false` | Require the command to supply the argument                         |
| `description` | string          | None    | Explanation shown in command output                               |
| `enum`        | list of strings | None    | Values accepted for the argument                                  |
| `pattern`     | string          | None    | Go (`RE2`) expression matched against the complete argument value |

`enum` and `pattern` can't be used together. A default value must satisfy the
declared `enum` or `pattern`.

### `kits`

`kits` accepts local directories, ZIP archives, OCI registry references, and
Git URLs prefixed with `git+https://` or `git+ssh://`. Kits can install tools,
configure the sandbox, and give the agent project-specific instructions. See
[Kits](../customize/kits.md) for details.

Explicit relative paths resolve from the directory of the environment file
that declares them. These include `.`, `..`, paths that start with `./` or
`../`, and relative paths that end in `.zip`. Bare references such as
`organization/kit` remain registry references.

Use an object entry to pass arguments to a kit. Set `source` to the kit
reference and map the kit's argument names to values under `args`:

```yaml
kits:
  - source: ./kits/tool
    args:
      version: ${{ env.args.channel }}
```

Remote kit sources must match the
[kit source allowlist](../customize/kits.md#restrict-kit-sources). Docker Hub is
allowed by default. To use Git kits from `docker/sbx-kits-contrib`, add its
source:

```console
$ sbx settings set kit.allowedSources '["docker.io/","github.com/docker/"]'
```

The setting replaces the complete allowlist, so include any existing sources
you want to keep. For reproducible setup, pin Git kits with the `ref` URL
parameter and OCI kits with an immutable tag or digest.

### `workspace`

When specified as a string, `workspace` is the path. Use the object form for
clone mode. Omit `workspace` to create a sandbox without a host bind mount. Set
`workspace: .` to mount the directory that contains the first environment
file. If `workspace` is present, its path can't be empty or contain only
whitespace.

When an environment file is inside a writable workspace, `sbx` binds the file
read-only at its path in the sandbox. Other files in the workspace remain
writable.

> [!WARNING]
> With [direct mount](../security/isolation.md#direct-mount-default), an agent
> can rename a subdirectory that contains an environment file and reach the
> underlying writable file. Store the environment file outside direct-mounted
> workspaces or directly in a workspace's root directory. Set
> `sandboxOptions.writableEnvFiles` only when the agent must edit the file.

| Field   | Type    | Required | Default | Description                                                     |
| ------- | ------- | -------- | ------- | --------------------------------------------------------------- |
| `path`  | string  | Yes      | None    | Workspace directory. Relative paths resolve from the first file |
| `clone` | boolean | No       | `false` | Use a private clone, equivalent to `sbx create --clone`          |

You can override `workspace.clone` for one `create` or `run` invocation with
`--clone` or `--clone=false`.

### `additionalWorkspaces`

Each additional workspace is mounted after the primary workspace. Relative
paths resolve from the directory of the first environment file.

| Field      | Type    | Required | Default | Description                   |
| ---------- | ------- | -------- | ------- | ----------------------------- |
| `path`     | string  | Yes      | None    | Directory to mount            |
| `readOnly` | boolean | No       | `false` | Mount the directory read-only |

### `sandboxOptions`

| Field              | Type            | Default  | Description                                                                       |
| ------------------ | --------------- | -------- | --------------------------------------------------------------------------------- |
| `template`         | string          | None     | Custom sandbox template image                                                     |
| `memory`           | string          | None     | Memory limit, such as `8g` or `512m`                                               |
| `cpus`             | integer         | `0`      | Number of CPUs. `0` allocates all host CPUs                                        |
| `pullPolicy`       | string          | `always` | Image pull policy: `always`, `missing`, or `never`                                 |
| `profile`          | string          | None     | Governance profile name                                                           |
| `shareSkills`      | boolean         | `true`   | Mount the shared agent skills store                                                |
| `display`          | boolean         | `false`  | Provision a display socket for graphical applications                             |
| `gpu`              | boolean         | `false`  | Pass the host GPU through to the sandbox                                           |
| `usb`              | list of strings | None     | USB device selectors to pass through to the sandbox                               |
| `writableEnvFiles` | boolean         | `false`  | Let the sandbox modify loaded environment files that would otherwise be read-only |

Imported [agent skills](../workflows/agent-skills.md) are shared with the
sandbox by default. Set `shareSkills: false` to opt out.

### `lifecycle`

The `lifecycle` block declares commands that run on the host with your user
privileges. Use lifecycle commands for work that must happen outside the
sandbox, such as creating a workspace, seeding fixtures, or archiving state.

```yaml
lifecycle:
  initialize:
    - name: Prepare workspace
      command: test -d web-app || git clone https://github.com/example/web-app
      timeout: 5m
  postCreate:
    - command: ./scripts/seed-fixtures.sh
      workdir: web-app
  preRemove:
    - command: ./scripts/archive-state.sh
```

Lifecycle phases run at the following points:

| Phase        | Timing                                                                                                           |
| ------------ | ---------------------------------------------------------------------------------------------------------------- |
| `initialize` | Before `sbx` resolves, provisions, creates, or attaches. Runs for every `create` and `run`                        |
| `postCreate` | After a new sandbox is created. For `run`, before attachment. Does not run when attaching to an existing sandbox   |
| `preRemove`  | After you approve removal and before `sbx` deletes resources. A failure produces a warning and removal continues  |

`sbx env exec` doesn't run lifecycle commands. Commands within a phase run in
order and stop at the first failure. Write `initialize` commands so repeated
runs produce the same result.

After `preRemove` finishes, `sbx` checks the destroy plan again. Removal stops
if the command added or changed a resource that the approved plan didn't name.

Each command requires `command` and accepts the following fields:

| Field     | Type   | Default           | Description                                                                    |
| --------- | ------ | ----------------- | ------------------------------------------------------------------------------ |
| `name`    | string | Command text      | Label shown in progress and plan output                                        |
| `command` | string | None              | Command passed to the user's shell                                             |
| `workdir` | string | Project directory | Host working directory. Relative paths resolve from the first file's directory |
| `timeout` | string | None              | Maximum runtime, such as `90s` or `5m`                                         |

Commands inherit the environment of the `sbx` process and receive the following
variables:

- `SBX_LIFECYCLE_PHASE`
- `SBX_ENV_FILE` and `SBX_ENV_FILES`
- `SBX_ENV_DIR`
- `SBX_SANDBOX_NAME`
- `SBX_AGENT`
- `SBX_WORKSPACE`

The environment file's `env` values and resolved secrets aren't passed to host
commands.

Plans containing lifecycle commands or credential `command` sources require
approval for every invocation by default, even when the command text hasn't
changed. Approve one invocation with `--auto-approve`, skip lifecycle commands
with `--skip-host-commands`, or remember approval until the commands change:

```console
$ sbx settings set env.rememberHostCommands true
```

### `secrets`

`secrets` maps service names to secret sources. Each entry must set exactly one
of `value`, `ref`, or `command`. The secret is stored at the sandbox scope when
the environment is created.

| Field      | Type    | Default | Description                                                                  |
| ---------- | ------- | ------- | ---------------------------------------------------------------------------- |
| `value`    | string  | None    | Literal secret value                                                         |
| `ref`      | string  | None    | Vault URI, such as `op://Vault/Item/field`                                    |
| `command`  | string  | None    | Host shell command whose standard output becomes the secret                   |
| `refresh`  | string  | None    | Resolution policy for `ref` or `command`, such as `on-demand` or `55m`        |
| `backend`  | string  | Automatic | Resolver for `ref`: `sdk` or `cli`                                          |
| `noVerify` | boolean | `false` | Skip verifying that a `ref` or `command` resolves during provisioning          |

> [!WARNING]
> A literal `value` is visible to anyone with read access to the file. Use a
> vault URI with `ref` or obtain the value at runtime with `command`.

```yaml
secrets:
  anthropic:
    ref: op://Private/Anthropic/api-key
    refresh: 55m
  github:
    command: gh auth token
```

### `bindings`

`bindings` approves credential injection domains for each service. The
environment merges these approvals into the user's global
`credentials.yaml`. Each service can contain an `apiKey` block, an `oauth`
block, or both. Each block contains a `domains` list:

```yaml
bindings:
  github:
    apiKey:
      domains:
        - api.github.com
```

`sbx env rm` preserves global bindings by default. Pass `--prune-bindings` to
remove every service binding declared by the environment file.

> [!WARNING]
> `--prune-bindings` deletes the complete global binding entry for every
> service declared in the environment file. This can affect other sandboxes
> that share those service bindings.

### `registries`

`registries` maps registry hostnames to pull credentials. Each entry requires
`secret` and accepts an optional `username`. Both fields accept a secret source
with exactly one of `value`, `ref`, or `command`.

When `username` is omitted, `sbx` stores a token-only credential. Registries
such as GHCR and GitLab accept token-only credentials.

```yaml
registries:
  ghcr.io:
    secret:
      command: gh auth token
```

### `mcp`

The `mcp.servers` list registers servers with the built-in
[MCP gateway](../mcp-gateway.md) and adds them to the sandbox. MCP registrations
are host-global and remain after `sbx env rm`.

| Field     | Type            | Required | Default | Description                                                     |
| --------- | --------------- | -------- | ------- | --------------------------------------------------------------- |
| `name`    | string          | Yes      | None    | Server name                                                     |
| `url`     | string          | No       | None    | Remote server URL, registry reference, or OCI reference          |
| `command` | string          | No       | None    | Command for a local stdio server                                 |
| `args`    | list of strings | No       | None    | Arguments passed to `command`                                   |

Each server must set exactly one of `url` or `command`.

### `ports`

`ports` publishes sandbox ports when the environment is created. Ports exposed
by a kit but omitted from this list receive an ephemeral host port.

| Field      | Type    | Required | Default                              | Description                                     |
| ---------- | ------- | -------- | ------------------------------------ | ----------------------------------------------- |
| `sandbox`  | integer | Yes      | None                                 | Sandbox port from 1 through 65535               |
| `host`     | integer | No       | Ephemeral                            | Host port from 1 through 65535                  |
| `protocol` | string  | No       | `tcp4`, or `tcp6` for IPv6 `hostIP` | `tcp`, `tcp4`, `tcp6`, `udp`, `udp4`, or `udp6` |
| `hostIP`   | string  | No       | Loopback                             | Host interface to bind                          |

Set `protocol: tcp` to bind both IPv4 and IPv6. Leave `hostIP` unset for a
dual-stack binding because an explicit address binds only its own IP family.

If a port can't be published, sandbox creation fails and removes the new
sandbox.
