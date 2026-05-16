---
title: Trust model for Compose files
weight: 70
description: Learn how Docker Compose treats Compose files as trusted input and what this means when using files you did not author.
keywords: compose, security, trust model, oci, remote, registry, include, extends, supply chain, trust, best practices
---

Docker Compose treats every Compose file as trusted input. When a Compose file
requests elevated privileges, host filesystem access, or any other
configuration, Compose applies it as written. This is the same behavior as
passing flags directly to `docker run`.

This means that any Compose file you run, whether it lives on your local
filesystem, in a Git repository, or in an OCI registry, has full control over
how containers interact with your host. The security boundary is not where the file comes from but whether you trust the author.

Evaluating trust means asking: do you know who authored this file, can you verify it hasn't changed since you last reviewed it, and do you understand every privilege it requests?

## The dependency chain

A Compose application can be assembled from multiple sources. The
[`include`](/reference/compose-file/include.md) directive imports entire Compose
files, while [`extends`](/reference/compose-file/services.md#extends) inherits
configuration from a specific service in another file. Both support remote
references and can be chained:

```text
Your command
  └─ compose.yaml                                    (local or remote)
       ├─ services, volumes, networks                (direct config)
       ├─ include:
       │    └─ oci://registry.example.com/base:v2   (remote dependency)
       │         └─ services, volumes, networks      (indirect config)
       └─ services:
            └─ app:
                 └─ extends:
                      └─ file: oci://registry.example.com/templates:v1
                           └─ service: webapp        (inherited config)
```

Each level has the same capabilities. The top-level file you inspect may appear
safe while a nested `include` or `extends` introduces services with elevated
privileges, host bind mounts, or untrusted images. These dependencies can also
change independently. Risky settings can be introduced by a nested dependency that you never
see unless you inspect the fully resolved output.

> [!IMPORTANT]
>
> Compose warns you when a configuration references remote sources. Do not
> accept this without understanding every reference in the chain.

## Best practices

### Inspect the full configuration

To see exactly what Compose applies, including all resolved `includes`,
`extends`, merged overrides, and interpolated variables, use:

```console
$ docker compose config
```

For remote references:

```console
$ docker compose -f oci://registry.example.com/myapp:latest config
```

Review this output before running `up` or `create`, especially when the
configuration comes from a source you have not audited.

#### Fields to look out for

A Compose configuration has broad control over how containers interact with the
host. The following is a non-exhaustive list of fields that carry security
implications when set by an untrusted author:

| Field | Effect |
|-------|--------|
| `privileged` | Grants the container full access to the host |
| `cap_add` | Adds Linux capabilities such as `SYS_ADMIN` or `NET_RAW` |
| `security_opt` | Configures security profiles including seccomp and AppArmor |
| `volumes` / bind mounts | Mounts host directories into the container |
| `network_mode: host` | Shares the host network stack |
| `pid: host` | Shares the host PID namespace |
| `devices` | Exposes host devices to the container |
| `image` | Pulls and runs an arbitrary container image |

When in doubt, look up the effect of any unfamiliar field before running the configuration.

### CI/CD environments

Automated pipelines are particularly sensitive because they often run with
access to credentials, cloud provider tokens, or Docker sockets.

- Avoid referencing public or unverified Compose configurations in automated
  pipelines.
- Gate updates behind your normal code review process.
- Use read-only Docker socket mounts where possible to limit your risk.

### Pin remote references to digests

Tags are mutable meaning anyone with push access to a registry can overwrite a tag silently, so a reference you reviewed last week may point to different content today.

Digests are immutable. Instead of referencing by tag, pin to the digest. 

```yaml
include:
  - oci://registry.example.com/base@sha256:a1b2c3d4...
```

Treat any update to a pinned digest as a code change. Make sure you review the new content before updating the reference.

### Other

- Use a private registry: Host OCI artifacts on a registry your
  organization controls. Restrict who can push to it.
- Audit transitive dependencies: Check every remote `include` and `extends`
  reference in the chain, not just the top-level file.
- Review all Compose confirmation prompts: When loading remote Compose files,
  Compose displays confirmation prompts for interpolation variables, environment
  values, and remote includes. Review these before accepting.

## Further reading

- [OCI artifact applications](/manuals/compose/how-tos/oci-artifact.md)
- [Use Compose in production](/manuals/compose/how-tos/production.md)
- [`include` reference](/reference/compose-file/include.md)
- [`extends` reference](/reference/compose-file/services.md#extends)
- [Manage secrets in Compose](/manuals/compose/how-tos/use-secrets.md)