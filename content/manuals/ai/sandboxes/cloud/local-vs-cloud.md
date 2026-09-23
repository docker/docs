---
title: Compare local and cloud sandboxes
linkTitle: Local and cloud
description: Compare local and cloud Docker Sandboxes, including workspaces, host integrations, ports, secrets, storage, and lifecycle behavior.
keywords: docker sandboxes, local sandbox, cloud sandbox, sbx cloud, comparison
weight: 10
---

Local and cloud sandboxes provide isolated environments for AI agents, but
they run against different resources and stores. This comparison helps you
choose an environment and identify workflows that need cloud-specific setup.

| Capability | Local sandbox | Cloud sandbox |
| --- | --- | --- |
| Compute | Uses resources from the host | Uses Docker-managed cloud resources |
| CPU architecture | Uses the platform supported by the local runtime | Supports `linux/amd64` and `linux/arm64`, subject to account availability; moves must retain the source platform |
| Workspace | Mounts host paths or uses a private Git clone backed by the host repository | Has no access to host paths; transfer files or clone a repository inside the sandbox |
| Hardware | Can use supported host integrations, such as GPU, USB, display, and nested virtualization | Has no access to host hardware |
| Ports | Binds sandbox ports to host addresses and ports | Exposes a sandbox port through a public HTTPS URL |
| Secrets | Reads from the local `sbx` secret store | Reads from the separate `sbx --cloud` secret store |
| Network policy | Uses local and organization policy sources supported by the local runtime | Uses separate, network-only account and sandbox policy; local policies aren't copied |
| MCP servers | Uses servers registered in the local MCP store | Uses servers connected through Docker Agentic Platform |
| Storage | Persists in the local sandbox and its attached host resources | Persists in the cloud sandbox and optional cloud volumes |
| Lifetime | Persists across stops until you remove it | Expires according to its time-to-live and timeout action |
| Billing | No metered sandbox compute charge | [Pay-as-you-go compute](/manuals/agentic-platform/signup.md#billing) |

## Host-dependent features

A cloud sandbox has no path back to the machine where you run `sbx`. The
following local features don't apply in cloud mode:

- Workspace paths, bind mounts, `--clone`, and Git worktrees created with
  `--branch`
- GPU, USB, display, and nested virtualization options
- Model selection through `--model` and `--provider`
- Host port bindings
- Host-backed agent skills and the full local MCP management workflow

The CLI rejects local-only flags used with `--cloud` instead of ignoring them.

You can use [environment files](../configuration/environment-files.md#use-a-cloud-environment)
with `sbx --cloud env`. Cloud environments support a subset of the local
configuration fields, including agents, environment variables, resource limits,
and credentials. Remove host workspace and port mappings before using a local
environment file in the cloud.

## Separate resources

Adding `--cloud` changes the backend for the command. A sandbox shown by
`sbx ls` doesn't appear in `sbx --cloud ls`, and resources created for one
backend don't automatically become available to the other.

This separation applies to sandboxes, templates, secrets, volumes, and network
policy. Use [`sbx move`](move.md) when you need to copy a sandbox filesystem
between backends. Moving doesn't unify the resource stores or transfer
host-mounted files or managed secrets. Credentials saved in copied files can
still travel with the sandbox.

## Network policy differences

Configure and verify [cloud network policy](network-policy.md) separately.
Moving a sandbox doesn't carry its local policy configuration into the cloud.
The cloud CLI supports account and sandbox network rules, but rejects local
governance profiles and the `--protocol` option.

HTTP method and path restrictions, also called L7 filtering, aren't supported
in cloud sandboxes. If these rules apply to a local sandbox, `sbx move` warns
and asks for confirmation before moving it to the cloud.

## What changes when you move

[`sbx move`](move.md) copies a filesystem snapshot and creates a separate
destination sandbox. It doesn't transfer running processes or memory, and it
leaves the source sandbox in place.

| State | Local to cloud | Cloud to local | What to do |
| --- | --- | --- | --- |
| Network policy | Local rules aren't copied; the CLI warns that cloud policy applies | Cloud rules aren't copied; the CLI notice explains that the host's default egress policy applies | Configure destination rules and verify allowed and blocked connections before continuing work |
| Managed secrets | Uses applicable credentials in the cloud secret store | Uses applicable credentials in the local secret store | Configure credentials in the destination store; the CLI warns that secrets aren't copied and you may need to sign in again |
| Workspace and attached storage | Host mounts and clone-mode volumes aren't copied; a workspace triggers a warning and confirmation | Cloud volumes aren't copied and no host folder is mounted | Transfer needed files separately with `sbx cp` |
| Published ports | TCP ports get cloud URLs; refused ports are skipped with a warning | Ports use local loopback bindings; cloud URLs aren't retained | Inspect destination ports and update clients |
| CPU architecture | Cloud access must support the source platform | The local runtime must support the source platform | Match platforms; the CLI checks compatibility before transfer |

The managed-secret warning doesn't itself block a move. Configure
[cloud credentials](credentials.md) or
[local credentials](../configuration/credentials.md) before using an agent
that needs them. Credentials written to files inside the sandbox can be
included in the snapshot; remove those files before moving if you don't want
the credentials copied.

Cloud-to-local moves display a notice and ask for confirmation. Local-to-cloud
moves ask for confirmation when workspace files or L7 filtering would be left
behind. If a confirmation is required and standard input isn't a terminal,
the move fails unless you pass `--force`. This flag skips confirmation but
retains warnings and doesn't change what transfers.

Verify the destination's files, credentials, network access, ports, and
[expiration settings](usage.md#configure-expiration) before removing the source.
