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
name must match the registry host as it appears in image references:

| Image reference                         | Directory                    |
| --------------------------------------- | ---------------------------- |
| `docker.io/myorg/myimage:latest`        | `docker.io/`                 |
| `registry.example.com/myimage:latest`   | `registry.example.com/`      |
| `registry.example.com:5000/myimage:tag` | `registry.example.com:5000/` |

Each directory contains a `hosts.toml` file:

```text
/etc/docker/certs.d/
├── docker.io/
│   └── hosts.toml
├── registry.example.com/
│   └── hosts.toml
└── registry.example.com:5000/
    └── hosts.toml
```

Changes to `hosts.toml` files take effect immediately, without restarting
Docker.

## hosts.toml format

Each `hosts.toml` file configures the behavior for one registry. The `server`
field sets the upstream registry URL. The `[host]` section configures specific
endpoints, including what operations they're allowed to perform using the
`capabilities` field.

Valid capabilities are:

| Capability | Description          |
| ---------- | -------------------- |
| `pull`     | Allow pulling images |
| `resolve`  | Allow tag resolution |
| `push`     | Allow pushing images |

## Examples

### Disable push to a registry

To prevent Docker from pushing images to a specific registry, omit `push` from
the capabilities:

```toml {title="/etc/docker/certs.d/docker.io/hosts.toml"}
server = "https://registry-1.docker.io"

[host."https://registry-1.docker.io"]
  capabilities = ["pull", "resolve"]
```

With this configuration, `docker pull` from Docker Hub works normally, but
`docker push` to Docker Hub returns an error.

### Redirect pulls to a mirror

To route pull traffic through a registry mirror:

```toml {title="/etc/docker/certs.d/docker.io/hosts.toml"}
server = "https://registry-1.docker.io"

[host."https://mirror.example.com"]
  capabilities = ["pull", "resolve"]

[host."https://registry-1.docker.io"]
  capabilities = ["pull", "resolve", "push"]
```

Docker tries the mirror first for pulls, and falls back to Docker Hub if the
mirror doesn't have the image. Pushes always go to Docker Hub directly.

### Internal registry only

To restrict Docker to only push and pull from an internal registry, and block
access to all public registries:

```toml {title="/etc/docker/certs.d/docker.io/hosts.toml"}
server = "https://registry-1.docker.io"

[host."https://registry-1.docker.io"]
  capabilities = []
```

With no capabilities, all operations to that registry fail.

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
