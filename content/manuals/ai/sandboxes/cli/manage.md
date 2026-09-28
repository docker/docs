---
title: Create and manage sandboxes
linkTitle: Manage sandboxes
description: Create, inspect, stop, resume, and remove local or cloud sandboxes with sbx.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 30
---

Use `sbx` for local sandboxes and add `--cloud` for cloud sandboxes. Choose the
environment before creating a sandbox: credentials and network policies must
be configured in that environment.

## Before you start

[Install the CLI](/manuals/ai/sandboxes/cli/install.md) and sign in with `sbx login`. Configure
[local credentials](/manuals/ai/sandboxes/cli/credentials.md) or [cloud credentials](/manuals/ai/sandboxes/cli/credentials-cloud.md)
for the agent. Cloud use also requires [cloud access](/manuals/ai/sandboxes/cloud-access/_index.md).

The examples use Claude Code and the name `my-project`.

## Create a sandbox

{{< tabs group="sandbox-environment" >}}
{{< tab name="Local" >}}

Create a sandbox with your project directory mounted:

```console
$ sbx create --name my-project claude ~/my-project
$ sbx run --name my-project
```

Edits to the mounted directory affect host files. For a private clone or a
sandbox without a mount, see [Workspace modes](/manuals/ai/sandboxes/cli/usage.md#choose-a-workspace).

{{< /tab >}}
{{< tab name="Cloud" >}}

Create the sandbox, then connect to its agent:

```console
$ sbx --cloud create --name my-project claude
$ sbx --cloud attach my-project
```

Cloud creation does not accept a local workspace path. Clone a repository or
[transfer files](/manuals/ai/sandboxes/cli/cloud-usage.md#transfer-files) into the sandbox. Check its
[expiration action](/manuals/ai/sandboxes/cli/cloud-usage.md#configure-expiration) before relying on it
to retain your work.

{{< /tab >}}
{{< /tabs >}}

## List sandboxes

Each command lists sandboxes in its selected environment:

{{< tabs group="sandbox-environment" >}}
{{< tab name="Local" >}}

```console
$ sbx ls
```

{{< /tab >}}
{{< tab name="Cloud" >}}

```console
$ sbx --cloud ls
```

{{< /tab >}}
{{< /tabs >}}

## Stop and resume

Stop a sandbox to pause execution while retaining its filesystem. Memory
preservation and connection continuity are separate concerns; see
[Lifecycle and persistence](/manuals/ai/sandboxes/concepts/lifecycle.md).

{{< tabs group="sandbox-environment" >}}
{{< tab name="Local" >}}

```console
$ sbx stop my-project
$ sbx run --name my-project
```

{{< /tab >}}
{{< tab name="Cloud" >}}

```console
$ sbx --cloud stop my-project
$ sbx --cloud ls
$ sbx --cloud attach my-project
```

Wait until the listing confirms the sandbox has stopped before resuming it.
Volume-backed sandboxes cannot be stopped. Check the expiration after resuming
with `sbx --cloud ttl my-project`.

{{< /tab >}}
{{< /tabs >}}

## Remove a sandbox

Copy or commit work you need to keep before removal. Removing a sandbox deletes
its internal filesystem. Files in a mounted local workspace remain on the host.

{{< tabs group="sandbox-environment" >}}
{{< tab name="Local" >}}

```console
$ sbx rm my-project
```

{{< /tab >}}
{{< tab name="Cloud" >}}

```console
$ sbx --cloud rm my-project
```

{{< /tab >}}
{{< /tabs >}}

For commands, ports, templates, and other operations, see
[Local operations](/manuals/ai/sandboxes/cli/usage.md) or [Cloud operations](/manuals/ai/sandboxes/cli/cloud-usage.md).
