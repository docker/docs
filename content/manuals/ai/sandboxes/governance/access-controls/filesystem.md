---
title: Filesystem access policies
linkTitle: Filesystem access
weight: 40
description: Control which host paths Docker Sandboxes can mount as workspaces with organization filesystem policies.
keywords: docker sandboxes, filesystem access, filesystem rules, workspace mount, organization policy, governance
---

Filesystem access policies control which host paths a sandbox can mount as a
workspace. Each policy contains one or more rules that restrict sandbox
workspaces to approved directories.

Filesystem access is managed with [organization policies](organization.md). When
organization governance is active, filesystem rules replace local behavior for
workspace mounts.

## Rule syntax

Filesystem rules use the actions `read` and `write`. Resources are host path
patterns.

A writable workspace mount must be allowed by both a `read` rule and a `write`
rule. A read-only workspace needs only `read`.

Examples:

- `~/**`
- `/data/project/**`
- `C:\data\project\**`
- `\\wsl.localhost\<distro>\data\project\**`

Use `**` to match a directory tree recursively. A single `*` matches only one
path segment. For exact path matching behavior across macOS, Linux, Windows,
and WSL, see [Filesystem rules](../concepts.md#filesystem-rules).

## Organization filesystem rules

Organization filesystem rules belong to policies that can apply to the whole
organization or to selected teams. For setup steps and team scoping, see
[Organization policies](organization.md).

Filesystem policy is checked when a workspace is mounted, which happens when a
sandbox is created. To apply a filesystem policy change to a running workflow,
remove the sandbox and create a new one.

## Troubleshooting

### Sandbox cannot mount workspace

If a sandbox fails to mount with a `mount policy denied` error, verify that the
filesystem allow rule uses `**` rather than `*`. A single `*` doesn't match
across directory separators.
