---
title: Sandbox environment files
linkTitle: Environment files
weight: 35
description: Use a declarative .sbxenv.yaml file to describe and share your sandbox configuration.
keywords:
  - docker sandboxes
  - sbx env
  - sbxenv
  - environment file
  - sandbox configuration
  - declarative
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
---

A sandbox environment file captures the setup for a project in a
`.sbxenv.yaml` file. Commit the file with your project so contributors use the
same agent, tools, resources, and credentials without reproducing CLI flags and
setup steps.

> [!NOTE]
> `sbx env` requires `sbx` 0.39.0 or later. The feature is experimental, so the
> command interface and file format may change in future releases.

## Start an environment

Create a `.sbxenv.yaml` file in your project directory. This example gives the
agent a shared environment variable and the Playwright browser-testing tools.
It also publishes the application's development port:

```yaml
schemaVersion: "1"
name: web-app
agent: claude

kits:
  - docker.io/sbx/playwright-kit:latest

env:
  NODE_ENV: test

ports:
  - sandbox: 3000
    host: 3000
```

From the same directory, run the environment:

```console
$ sbx env run
```

The project directory becomes the workspace. If the environment doesn't exist,
`sbx` creates a sandbox named `web-app`, installs Playwright and Chromium, and
publishes sandbox port `3000` on the host. It then attaches to the agent. Later
runs attach to the existing sandbox.

## Commands

| Command                                                                                 | Description                                                                                          |
| --------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| [`sbx env run`](/reference/cli/sbx/env/run/) `[PATH...]`                                | Creates the environment if needed, then attaches. Re-runs apply only [`env` and MCP changes](#update-an-environment) |
| [`sbx env create`](/reference/cli/sbx/env/create/) `[PATH...]`                          | Creates the environment without attaching                                                            |
| [`sbx env exec`](/reference/cli/sbx/env/exec/) `[PATH...] -- COMMAND [ARG...]`           | Runs a command in an existing environment                                                            |
| [`sbx env rm`](/reference/cli/sbx/env/rm/) `[PATH...]`                                  | Removes the sandbox and its scoped credentials                                                       |

`PATH` can be a directory or a direct path to an environment file. When you
pass a directory, `sbx` reads `.sbxenv.yaml` and falls back to `.sbxenv.yml`.
With no path, `sbx` searches the working directory.

Pass the same set of paths to each lifecycle command so they resolve the same
sandbox.

## Common workflows

The following examples combine environment file fields into configurations you
can adapt for a project.

### Combine team defaults and personal settings

Keep the shared configuration in a committed file and put machine-specific
settings in a file excluded from version control. For example, commit
`base.sbxenv.yaml`:

```yaml
schemaVersion: "1"
name: web-app
agent: claude

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
# .sbxenv.yaml in the web-app repository
schemaVersion: "1"
name: web-platform
agent: codex

workspace: .

additionalWorkspaces:
  - path: ../shared-components
  - path: ../architecture-docs
    readOnly: true
```

The agent starts in `web-app`, can modify `shared-components`, and can read
`architecture-docs` without changing it. Relative paths resolve from the
directory of the first environment file. Additional workspaces are mounted
directly even when the primary workspace uses clone mode.

### Reuse an environment in automation

Use the same committed environment for interactive development and automated
tasks. Developers attach to the agent with `run`:

```console
$ sbx env run
```

Automation can create the sandbox without attaching, run commands in it, and
remove it afterward:

```console
$ sbx env create
$ sbx env exec -- npm test
$ sbx env rm --force
```

Commands and vault references under `secrets` resolve on the host, so the
automation runner must provide the referenced tools and authentication. The
secret values remain outside the environment file.

## Update an environment

`sbx env run` starts and attaches to an existing sandbox without provisioning
its secrets and bindings again. For an existing sandbox, the command applies
changes to `env` to the new session and reconciles declared MCP servers. Changes
to workspaces, kits, ports, secrets, bindings, and `sandboxOptions` require you
to remove the environment with `sbx env rm` and create it again.

## Remove an environment

Secrets and registry credentials are sandbox-scoped. Credential bindings and
MCP server registrations are host-global and can be shared by multiple
sandboxes.

`sbx env rm` removes the sandbox and its scoped credentials. Global credential
bindings remain unless you pass `--prune-bindings`. MCP registrations remain
available to other sandboxes.

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
| `kits`                 | list of strings  | No       | None                           | Kits to install at creation. See [`kits`](#kits)                                 |
| `workspace`            | string or object | No       | First file's directory         | Primary workspace. See [`workspace`](#workspace)                                |
| `additionalWorkspaces` | list             | No       | None                           | Extra directories to mount. See [`additionalWorkspaces`](#additionalworkspaces) |
| `env`                  | map of strings   | No       | None                           | Environment variables for the sandbox                                           |
| `sandboxOptions`       | object           | No       | None                           | Creation options. See [`sandboxOptions`](#sandboxoptions)                        |
| `secrets`              | map              | No       | None                           | Service credentials. See [`secrets`](#secrets)                                  |
| `bindings`             | map              | No       | None                           | Credential injection approvals. See [`bindings`](#bindings)                     |
| `registries`           | map              | No       | None                           | Registry pull credentials. See [`registries`](#registries)                      |
| `mcp`                  | object           | No       | None                           | MCP servers. See [`mcp`](#mcp)                                                  |
| `ports`                | list             | No       | None                           | Port mappings. See [`ports`](#ports)                                            |

### `kits`

`kits` accepts local directories, ZIP archives, OCI registry references, and
Git URLs prefixed with `git+https://` or `git+ssh://`. Kits can install tools,
configure the sandbox, and give the agent project-specific instructions. See
[Kits](customize/kits.md) for details.

Remote kit sources must match the
[kit source allowlist](customize/kits.md#restrict-kit-sources). Docker Hub is
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
clone mode:

| Field   | Type    | Default                | Description                                                             |
| ------- | ------- | ---------------------- | ----------------------------------------------------------------------- |
| `path`  | string  | First file's directory | Workspace directory. Relative paths resolve from the first file         |
| `clone` | boolean | `false`                | Use a private clone, equivalent to `sbx create --clone`                  |

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

| Field        | Type    | Default  | Description                                                     |
| ------------ | ------- | -------- | --------------------------------------------------------------- |
| `template`   | string  | None     | Custom sandbox template image                                   |
| `memory`     | string  | None     | Memory limit, such as `8g` or `512m`                             |
| `cpus`       | integer | `0`      | CPU limit. `0` selects the automatic value                       |
| `pullPolicy` | string  | `always` | Image pull policy: `always`, `missing`, or `never`               |
| `profile`    | string  | None     | Governance profile name                                         |

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
| `noVerify` | boolean | `false` | Skip resolving a `ref` or `command` once when provisioning the secret         |

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
[MCP gateway](mcp-gateway.md) and adds them to the sandbox. MCP registrations
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

| Field      | Type    | Required | Default          | Description                                                           |
| ---------- | ------- | -------- | ---------------- | --------------------------------------------------------------------- |
| `sandbox`  | integer | Yes      | None             | Sandbox port from 1 through 65535                                     |
| `host`     | integer | No       | Ephemeral        | Host port from 1 through 65535                                        |
| `protocol` | string  | No       | `tcp`            | `tcp`, `tcp4`, `tcp6`, `udp`, `udp4`, or `udp6`                       |
| `hostIP`   | string  | No       | Loopback         | Host interface. The default uses available IPv4 and IPv6 loopback     |

If a port can't be published, sandbox creation fails and removes the new
sandbox.
