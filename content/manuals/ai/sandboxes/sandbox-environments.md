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
> `sbx env` is experimental. The command interface and file format may change
> in future releases.

A sandbox environment file (`.sbxenv.yaml`) describes a sandbox in code:
the agent, kits, workspace, secrets, ports, and resource limits. Commit one
alongside your project and every team member runs an identical sandbox with a
single command, without sharing flag combinations or setup instructions.

The environment file doesn't need to live in the same directory as the
workspace. You can place it anywhere and point `workspace.path` at the target
directory:

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

| Command                    | Description                                                                       |
| -------------------------- | --------------------------------------------------------------------------------- |
| `sbx env run [PATH...]`    | Provisions the environment if it doesn't exist, then opens an interactive session |
| `sbx env create [PATH...]` | Provisions without attaching; use in scripts and CI                               |
| `sbx env rm [PATH...]`     | Removes the sandbox and all resources provisioned by the environment              |

`PATH` can be a directory (reads the `.sbxenv.yaml` inside it) or a direct
path to a file. With no argument, `sbx env` reads `.sbxenv.yaml` in the
current directory.

Passing multiple paths merges the files in order — see
[Multiple files](#multiple-files).

`sbx env rm` removes all secrets and registry credentials that were
provisioned when the environment was created, not just the sandbox container.

## Multiple files

`sbx env run` and `sbx env create` accept multiple paths and deep-merge them
in order. Later files override scalar values; lists concatenate. This follows
the same semantics as Docker Compose's multiple `-f` files:

```console
$ sbx env run base.sbxenv.yaml local.sbxenv.yaml
```

A common pattern is to commit a `base.sbxenv.yaml` with shared team
configuration and add `local.sbxenv.yaml` to `.gitignore` for personal
overrides: a different workspace path, additional secrets, or adjusted
resource limits.

## Variable interpolation

Host environment variables are expanded before the file is parsed:

| Syntax             | Behavior                                        |
| ------------------ | ----------------------------------------------- |
| `$VAR` or `${VAR}` | Expands to the value of `VAR`; fails if unset   |
| `${VAR:-default}`  | Uses `default` if `VAR` is unset or empty       |
| `${VAR:?message}`  | Fails with `message` if `VAR` is unset or empty |
| `$$`               | Literal `$`                                     |

```yaml
workspace:
  path: $HOME/src/myproject

secrets:
  my-token:
    value: ${MY_TOKEN:?MY_TOKEN must be set}
```

## File reference

### Top-level fields

| Field                  | Type             | Required | Default                        | Description                                                                     |
| ---------------------- | ---------------- | -------- | ------------------------------ | ------------------------------------------------------------------------------- |
| `schemaVersion`        | string           | Yes      | —                              | Schema version. Currently `"1"`                                                 |
| `name`                 | string           | No       | `<agent>-<workspace-basename>` | Sandbox name                                                                    |
| `agent`                | string           | Yes      | —                              | Agent type, e.g. `claude`                                                       |
| `workspace`            | string or object | No       | `.`                            | Workspace path or configuration; see [`workspace`](#workspace)                  |
| `additionalWorkspaces` | list             | No       | —                              | Extra directories to mount; see [`additionalWorkspaces`](#additionalworkspaces) |
| `kits`                 | list of strings  | No       | —                              | Kit references to install at sandbox creation                                   |
| `env`                  | map              | No       | —                              | Environment variables to inject into the sandbox                                |
| `sandboxOptions`       | object           | No       | —                              | Resource and image-pull settings; see [`sandboxOptions`](#sandboxoptions)       |
| `secrets`              | map              | No       | —                              | Service credentials; see [`secrets`](#secrets)                                  |
| `registries`           | map              | No       | —                              | Registry pull credentials; see [`registries`](#registries)                      |
| `ports`                | list             | No       | —                              | Port mappings; see [`ports`](#ports)                                            |

### `workspace`

When specified as a string, `workspace` is treated as the path. Use the
object form to enable clone mode or when the file doesn't live next to the
workspace:

| Field   | Type    | Default | Description                                                             |
| ------- | ------- | ------- | ----------------------------------------------------------------------- |
| `path`  | string  | `.`     | Path to the workspace directory                                         |
| `clone` | boolean | `false` | Mount the workspace as a private clone, equivalent to `sbx run --clone` |

### `additionalWorkspaces`

A list of extra directories to mount alongside the primary workspace:

| Field      | Type    | Required | Description                   |
| ---------- | ------- | -------- | ----------------------------- |
| `path`     | string  | Yes      | Path to the directory         |
| `readOnly` | boolean | No       | Mount the directory read-only |

### `sandboxOptions`

| Field        | Type   | Default  | Description                                                     |
| ------------ | ------ | -------- | --------------------------------------------------------------- |
| `memory`     | string | —        | Memory limit, e.g. `8g`, `512m`                                 |
| `cpus`       | number | —        | CPU limit                                                       |
| `pullPolicy` | string | `always` | When to pull the sandbox image: `always`, `missing`, or `never` |
| `template`   | string | —        | Custom sandbox template image                                   |
| `profile`    | string | —        | Governance profile name                                         |

### `secrets`

A map of secret names to secret sources. Each secret is provisioned when the
environment is created, scoped to the sandbox. `sbx env rm` removes all
secrets in this map.

| Field     | Description                                                                                      |
| --------- | ------------------------------------------------------------------------------------------------ |
| `ref`     | A vault URI, e.g. `op://Vault/Item/field` (1Password). Resolved from the vault at creation time. |
| `command` | A shell command whose stdout becomes the secret value.                                           |
| `value`   | A plaintext secret value.                                                                        |
| `refresh` | Re-fetch interval for `ref`-based secrets, e.g. `55m`.                                           |

> [!WARNING]
> Avoid setting real credentials as a plaintext `value`. The plaintext is visible to
> anyone with read access to the file. Use `ref` (vault URI) or `command`
> to source the value at runtime, or use variable interpolation to read it
> from the environment: `value: ${MY_TOKEN}`.

```yaml
secrets:
  anthropic:
    ref: op://Private/Anthropic/api-key
    refresh: 55m
  github:
    command: gh auth token
```

### `registries`

A map of registry hostnames to pull credentials. Each entry requires
`secret` and accepts an optional `username`. Both fields use the same secret
source forms as [`secrets`](#secrets) (`ref`, `command`, or `value`).
Omitting `username` stores the credential as token-only, which registries
like GHCR and GitLab accept.

```yaml
registries:
  ghcr.io:
    secret:
      command: gh auth token
```

### `ports`

A list of port mappings between the sandbox and the host:

| Field      | Type    | Required | Default | Description                                                        |
| ---------- | ------- | -------- | ------- | ------------------------------------------------------------------ |
| `sandbox`  | integer | Yes      | —       | Port number inside the sandbox                                     |
| `host`     | integer | No       | —       | Port number on the host. Omit to expose without a fixed host port. |
| `protocol` | string  | No       | `tcp`   | Protocol: `tcp` or `udp`                                           |
