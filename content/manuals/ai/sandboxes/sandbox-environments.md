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

> [!NOTE]
> `sbx env` requires `sbx` 0.39.0 or later and supports local sandboxes only.
> The feature is experimental, so the command interface and file format may
> change in future releases.

A sandbox environment file describes the agent, kits, workspaces, environment
variables, credentials, MCP servers, ports, and resource limits for a sandbox.
Commit the file with your project so team members can run the same sandbox
configuration without sharing flag combinations or setup instructions.

The environment file doesn't need to be in the workspace. You can store it in
another directory and set `workspace.path` to the workspace:

```yaml
# .sbxenv.yaml
schemaVersion: "1"
name: docs-env
agent: claude

workspace:
  path: $HOME/src/github.com/docker/docs
  clone: true

kits:
  - "git+https://github.com/docker/sbx-kits-contrib.git#dir=vale"
  - "git+https://github.com/docker/sbx-kits-contrib.git#dir=git-ssh-sign"
  - "git+https://github.com/docker/sbx-kits-contrib.git#dir=github-ssh"

secrets:
  github:
    command: gh auth token

ports:
  - sandbox: 1313
    host: 1313
```

## Commands

| Command                                                                                 | Description                                      |
| --------------------------------------------------------------------------------------- | ------------------------------------------------ |
| [`sbx env run`](/reference/cli/sbx/env/run/) `[PATH...]`                                | Creates the environment if needed, then attaches |
| [`sbx env create`](/reference/cli/sbx/env/create/) `[PATH...]`                          | Creates the environment without attaching        |
| [`sbx env exec`](/reference/cli/sbx/env/exec/) `[PATH...] -- COMMAND [ARG...]`           | Runs a command in an existing environment        |
| [`sbx env rm`](/reference/cli/sbx/env/rm/) `[PATH...]`                                  | Removes the sandbox and its scoped credentials   |

`PATH` can be a directory or a direct path to an environment file. When you
pass a directory, `sbx` reads `.sbxenv.yaml` and falls back to `.sbxenv.yml`.
With no path, `sbx` searches the working directory.

Pass the same set of paths to each lifecycle command so they resolve the same
sandbox. See [Multiple files](#multiple-files) for details.

`sbx env run` starts and attaches to an existing sandbox without provisioning
its secrets and bindings again. For an existing sandbox, the command applies
changes to `env` to the new session and reconciles declared MCP servers. Changes
to workspaces, kits, ports, secrets, bindings, and `sandboxOptions` require you
to remove the environment with `sbx env rm` and create it again.

For `sbx env exec`, arguments before `--` are environment-file paths and
arguments after `--` form the command. Without `--`, all arguments form the
command and `sbx` reads the environment file from the working directory:

```console
$ sbx env exec .sbxenv.yaml -- go test ./...
$ sbx env exec go test ./...
```

`sbx env rm` removes the sandbox and its scoped service and registry
credentials. Global credential bindings remain unless you pass
`--prune-bindings`. Host-global MCP registrations remain available to other
sandboxes.

Secret provisioning, binding updates, and MCP server registration occur before
the sandbox is created. If sandbox creation fails, scoped secrets remain, and
bindings and MCP registrations may also remain. Run `sbx env rm` with the same
paths to remove the scoped secrets. Pass `--prune-bindings` if you also want to
remove the declared global bindings. MCP registrations are host-global and
remain after cleanup.

## Multiple files

Pass multiple paths to merge environment files in order. Nested mappings merge
by key, lists concatenate, and values from later files replace earlier scalar
values. Relative workspace paths and the default sandbox name use the directory
of the first file.

```console
$ sbx env run base.sbxenv.yaml local.sbxenv.yaml
```

Commit a `base.sbxenv.yaml` with shared configuration and add
`local.sbxenv.yaml` to `.gitignore` for personal overrides, such as another
workspace path, additional secrets, or different resource limits.

## Variable interpolation

Host environment variables are expanded in each file before the files are
merged and parsed:

| Syntax                   | Behavior                                                    |
| ------------------------ | ----------------------------------------------------------- |
| `$VAR` or `${VAR}`       | Uses the value of `VAR`, or an empty string when unset      |
| `${VAR:-default}`        | Uses `default` when `VAR` is unset or empty                 |
| `${VAR-default}`         | Uses `default` when `VAR` is unset                          |
| `${VAR:+replacement}`    | Uses `replacement` when `VAR` is set and non-empty          |
| `${VAR+replacement}`     | Uses `replacement` when `VAR` is set                        |
| `${VAR:?message}`        | Fails with `message` when `VAR` is unset or empty           |
| `${VAR?message}`         | Fails with `message` when `VAR` is unset                    |
| `$$`                     | Inserts a literal `$`                                       |

Default, replacement, and error values can contain nested variable
expressions.

```yaml
workspace:
  path: ${WORKSPACE:-$HOME/src/myproject}

secrets:
  my-token:
    value: ${MY_TOKEN:?MY_TOKEN must be set}
```

## File reference

The loader rejects unknown fields and unsupported schema versions.

### Top-level fields

| Field                  | Type             | Required | Default                        | Description                                                                     |
| ---------------------- | ---------------- | -------- | ------------------------------ | ------------------------------------------------------------------------------- |
| `schemaVersion`        | string           | Yes      | None                           | Schema version. The supported value is `"1"`                                   |
| `name`                 | string           | No       | `<agent>-<workspace-basename>` | Sandbox name                                                                    |
| `agent`                | string           | Yes      | None                           | Built-in agent or the name of an agent kit                                      |
| `kits`                 | list of strings  | No       | None                           | Directory, ZIP, or OCI kit references to install at creation                    |
| `workspace`            | string or object | No       | First file's directory         | Primary workspace. See [`workspace`](#workspace)                                |
| `additionalWorkspaces` | list             | No       | None                           | Extra directories to mount. See [`additionalWorkspaces`](#additionalworkspaces) |
| `env`                  | map of strings   | No       | None                           | Environment variables for the sandbox                                           |
| `sandboxOptions`       | object           | No       | None                           | Creation options. See [`sandboxOptions`](#sandboxoptions)                        |
| `secrets`              | map              | No       | None                           | Service credentials. See [`secrets`](#secrets)                                  |
| `bindings`             | map              | No       | None                           | Credential injection approvals. See [`bindings`](#bindings)                     |
| `registries`           | map              | No       | None                           | Registry pull credentials. See [`registries`](#registries)                      |
| `mcp`                  | object           | No       | None                           | MCP servers. See [`mcp`](#mcp)                                                  |
| `ports`                | list             | No       | None                           | Port mappings. See [`ports`](#ports)                                            |

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
> vault URI with `ref`, obtain the value at runtime with `command`, or use
> variable interpolation such as `value: ${MY_TOKEN}`.

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
