---
title: Configure a registry mirror
linkTitle: Registry mirror
description: Route Docker Hub template, kit, and in-sandbox Docker image pulls through an organization's registry mirror.
keywords: docker sandboxes, sbx, registry mirror, docker hub, templates, kits, image pulls
weight: 50
---

A registry mirror routes Docker Hub pulls for sandbox templates and OCI kits
through your organization's registry infrastructure. If the mirror meets
Docker Engine's requirements, Docker inside the sandbox uses it too.

## Configure the mirror

Set `platform.images.registryMirror` to the mirror host. Include a port when
needed, but omit the URL scheme:

```console
$ sbx settings set platform.images.registryMirror registry.example.com
```

You can include a path prefix for registries that store mirrored Docker Hub
content below a repository path:

```console
$ sbx settings set platform.images.registryMirror registry.example.com/docker-remote
```

Docker Sandboxes redirects image references that resolve to Docker Hub and
preserves their repository path, tag, and digest. References that explicitly
name another registry remain unchanged.

If the mirror requires authentication, configure
[registry credentials](credentials.md#registry-credentials) for the mirror
host.

## Mirror Docker pulls inside the sandbox

Docker Sandboxes configures Docker Engine inside a sandbox to use the same
mirror when the setting contains a reachable host without a path prefix.

| Mirror setting                              | Template and OCI kit pulls | Docker pulls inside the sandbox |
| ------------------------------------------- | -------------------------- | ------------------------------- |
| `registry.example.com`                      | Mirrored                   | Mirrored                        |
| `registry.example.com:5000`                 | Mirrored                   | Mirrored                        |
| `registry.example.com/docker-remote`        | Mirrored                   | Not mirrored                    |
| `localhost:5000`, `127.0.0.1`, or `0.0.0.0` | Mirrored when host-reachable | Not mirrored                    |

Loopback and wildcard addresses refer to the sandbox itself from inside its
network namespace, so Docker Sandboxes doesn't add them to the sandbox's Docker
Engine configuration. A path prefix is also excluded because Docker Engine
interprets mirror URL paths differently from image repository prefixes.

Docker Engine connects to the mirror over HTTPS, so the sandbox must trust the
certificate that the mirror presents. For a mirror that uses an internal
certificate authority, add the CA to the sandbox's system trust store. See
[Install an internal CA certificate](../customize/kit-examples.md#install-an-internal-ca-certificate).

Template and kit pulls use the changed setting immediately. Existing sandboxes
retain the Docker Engine mirror configuration with which they were created.
Recreate a sandbox to apply a changed mirror to Docker inside it.

## Disable the mirror

Unset the setting to disable mirroring:

```console
$ sbx settings unset platform.images.registryMirror
```

An empty setting value also disables mirroring. Recreate existing sandboxes to
remove a mirror from their Docker Engine configuration.
