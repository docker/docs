---
title: Move a sandbox
description: Transfer a Docker Sandbox filesystem between local and cloud environments and understand which state, files, policies, and secrets remain behind.
keywords: docker sandboxes, sbx move, cloud sandbox, local sandbox, migrate sandbox
weight: 50
---

The `sbx move` command transfers a filesystem snapshot between a local sandbox
and a cloud sandbox. Use it to continue from the captured filesystem in the
other environment, not as a live migration of the running sandbox.

## Transfer behavior

A move performs the following operations:

1. Captures the source sandbox filesystem as a template image
2. Transfers the image across the local and cloud boundary
3. Creates a destination sandbox from the image

The destination gets a different sandbox ID. The source isn't deleted, so the
two sandboxes have independent state after the transfer. A local-to-cloud move
stops the local source while capturing it. A cloud-to-local move attempts to
stop the cloud source after creating the destination. This is a best-effort
operation, so the cloud source can remain running if it can't be stopped.
The move can return before the stop finishes. Check `sbx --cloud ls` for its
state. A cloud source with no agent is left running.
Remove the source separately after you verify the destination.

> [!IMPORTANT]
>
> Managed secrets aren't copied from the source secret store. Credentials
> written inside a sandbox by an agent's interactive sign-in are ordinary
> filesystem files and are included in the snapshot. Remove in-sandbox
> credentials before moving a sandbox.

Moving transfers filesystem data stored in the sandbox container. It doesn't
transfer:

- Running processes, memory, or open sockets
- Local workspace mounts, bind mounts, or clone-mode volumes
- Managed secrets from the source secret store
- Other resources attached outside the sandbox filesystem

Changes made to either sandbox after the move aren't synchronized.

The snapshot retains the source sandbox's platform. The destination must
support the same platform because `sbx move` doesn't convert between
`linux/amd64` and `linux/arm64`. For a local-to-cloud move, your cloud account
must support the local sandbox's platform. For a cloud-to-local move, create
the cloud sandbox with `--platform linux/amd64` or `--platform linux/arm64`
to match your local sandbox runtime. The CLI checks compatibility before
transferring the snapshot.

## Move from local to cloud

Move a local sandbox to the cloud:

```console
$ sbx move local-project --to cloud
```

The `move` command spans both backends, so don't add the global `--cloud` flag.
Use `--name` to set the destination name prefix. The cloud destination appends
a short unique suffix:

```console
$ sbx move local-project --to cloud --name cloud-project
```

Local workspace files are mounted outside the sandbox filesystem and don't
appear in the cloud destination. When the source has a workspace, the CLI
warns and asks for confirmation. The `--force` flag skips the prompt but
doesn't include those files.

The destination uses cloud network policy. Local network rules don't transfer.
If the local source has HTTP method or path restrictions, the CLI warns and
asks for confirmation because those restrictions won't apply in the cloud.
`--force` skips the prompt but retains the warning.

Published TCP sandbox ports are published on the cloud destination with cloud
URLs. Ports that the cloud refuses are skipped with a warning. Host port
numbers and non-TCP mappings don't transfer.

The destination can use applicable credentials that already exist in the cloud
secret store. Credentials in the local secret store don't transfer.

The CLI rounds the source's recorded CPU and memory limits up to a supported
[cloud size](usage.md#choose-resources-and-platform). Missing limits use cloud
defaults with a warning. Limits above the largest cloud size prevent the move.

### Set destination expiration

The cloud destination has an expiration time, even if the local source had
none. Without `--ttl`, it uses the server default, typically one hour. The
move requests that the sandbox stop on expiration when the account and sandbox
support it. Otherwise, it falls back to the server's default timeout action,
which deletes the sandbox.

Set the expiration and action explicitly when moving work you want to retain:

```console
$ sbx move local-project --to cloud --ttl 2h --on-timeout stop
```

The `stop` action preserves state. An explicit request fails if stopping is
unavailable. Use `--on-timeout delete` to delete on expiration.
These flags apply only to moves to the cloud. After the move, inspect or extend
the expiration with [`sbx --cloud ttl`](usage.md#configure-expiration).

## Move from cloud to local

Move a cloud sandbox to the local runtime by ID or name:

```console
$ sbx move cloud-project --to local --name local-copy
```

The local destination starts with the host's default network policy. Cloud
network rules aren't copied back to the local runtime.

Published TCP ports are saved on the local sandbox and bound to loopback while
it runs. Host port numbers can change after a restart. Run
`sbx ports local-copy` to see the bindings.

The move doesn't create a host workspace. Use `sbx cp` to copy files from the
local sandbox to your host. Attached cloud volumes and environment variables
don't transfer. The local destination uses local CPU and memory defaults.

This direction requires a machine that meets the local sandbox requirements
because the command imports the snapshot and creates a local sandbox. The
cloud source remains available unless you remove it with `sbx --cloud rm`.

Layer reuse during download can require up to 32 GiB of temporary host disk
space in addition to the local runtime's image storage.

## Verify the result

List both backends after the move:

```console
$ sbx ls
$ sbx --cloud ls
```

Inspect the destination filesystem and verify its credentials, ports, volumes,
and other environment-specific resources before removing the source. Configure
anything that didn't carry over.
