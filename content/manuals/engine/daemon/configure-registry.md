---
title: Registry host configuration
description: Configure per-registry behavior for Docker Engine using hosts.toml files
keywords: containerd, registry, hosts, push, pull, mirror, configuration, daemon
weight: 25
---

When using the [containerd image store](/manuals/engine/storage/containerd.md),
you can configure per-registry behavior using `hosts.toml` files. This lets
you restrict push or pull access, redirect traffic to a mirror, or customize
TLS settings on a per-registry basis.

## Configuration directory

Docker Engine reads registry host configuration from the following directory:

| Setup         | Directory                   |
| ------------- | --------------------------- |
| Regular       | `/etc/docker/certs.d/`      |
| Rootless mode | `~/.config/docker/certs.d/` |

Create a subdirectory for each registry you want to configure. The directory
name must match the registry host as it appears in image references. On
Windows, replace the colon with underscores for registries that use a
non-standard port, since colons are not valid in Windows directory names:

| Image reference                         | Directory (Linux)                    | Directory (Windows)           |
| --------------------------------------- | ------------------------------------ | ----------------------------- |
| `docker.io/myorg/myimage:latest`        | `docker.io/`                         | `docker.io\`                  |
| `registry.example.com/myimage:latest`   | `registry.example.com/`              | `registry.example.com\`       |
| `registry.example.com:5000/myimage:tag` | `registry.example.com:5000/`         | `registry.example.com_5000_\` |

Each directory contains a `hosts.toml` file:

```text
/etc/docker/certs.d/
├── _default/
│   └── hosts.toml
├── docker.io/
│   └── hosts.toml
├── registry.example.com/
│   └── hosts.toml
└── registry.example.com:5000/
    └── hosts.toml
```

The `_default/` directory is optional. If present, its `hosts.toml` settings
apply to any registry that doesn't have its own directory. Use it to set a
global mirror or change the default behavior for all registries.

Changes to `hosts.toml` files take effect immediately, without restarting
Docker.

## hosts.toml format

Each `hosts.toml` file configures one registry. The file has two levels of
configuration:

- **Top-level fields** — apply to the registry's default endpoint (the server
  itself). Common fields are `server`, `capabilities`, `ca`, and `skip_verify`.
- **`[host]` sections** — configure additional endpoints such as mirrors.
  Hosts are tried in the order listed before falling back to the `server`.

The `server` field sets the upstream registry URL. If omitted, Docker uses the
registry name from the image reference.

Valid values for `capabilities` are:

| Capability | Description                       |
| ---------- | --------------------------------- |
| `pull`     | Allow pulling images              |
| `resolve`  | Allow resolving a tag to a digest |
| `push`     | Allow pushing images              |

If `capabilities` is not set, all three are enabled by default.

## Examples

### Disable push to a registry

To prevent Docker from pushing images to a specific registry, set `capabilities`
at the top level and omit `push`:

```toml {title="/etc/docker/certs.d/docker.io/hosts.toml"}
server = "https://registry-1.docker.io"
capabilities = ["pull", "resolve"]
```

With this configuration, `docker pull` from Docker Hub works normally, but
`docker push` to Docker Hub returns an error.

### Redirect pulls to a mirror

To route pull traffic through a registry mirror, add the mirror as a `[host]`
entry. Docker tries the mirror first for pulls and falls back to the upstream
registry if the mirror doesn't have the image.

Pushes always bypass mirrors and go directly to the upstream registry because
the mirror entry only has `pull` and `resolve` capabilities:

```toml {title="/etc/docker/certs.d/docker.io/hosts.toml"}
server = "https://registry-1.docker.io"

[host."https://mirror.example.com"]
  capabilities = ["pull", "resolve"]
```

### Use an internal registry

To route all traffic for a registry namespace through an internal host, set
`server` to the internal registry URL. No `[host]` entries are needed:

```toml {title="/etc/docker/certs.d/docker.io/hosts.toml"}
server = "https://registry.internal.example.com"
```

With this configuration, Docker sends all push and pull operations for
`docker.io` images to `registry.internal.example.com` instead.

> [!NOTE]
> This configuration controls behavior at the daemon level, not as a security
> boundary. Builds, containers, and other mechanisms can still interact with
> registries. For strict registry access control, consider
> [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md)
> in Docker Business.

## Relation to daemon.json registry settings

Docker daemon also supports registry configuration through `daemon.json` options
like `insecure-registries` and `registry-mirrors`. These settings interact with
`hosts.toml` as follows:

- If a `hosts.toml` file configures **two or more** endpoints for a registry
  (such as a mirror and an upstream fallback), the daemon.json settings for that
  registry are **ignored**. The `hosts.toml` configuration takes full control.
- If `hosts.toml` is absent or configures only a single endpoint, the
  daemon.json settings are applied on top.

If you're using `hosts.toml` to configure mirrors for a registry, include all
TLS and authentication settings in the `hosts.toml` file rather than relying on
`insecure-registries` in `daemon.json`.

## Reference

For the full `hosts.toml` specification, see the
[containerd registry hosts documentation](https://github.com/containerd/containerd/blob/main/docs/hosts.md).
