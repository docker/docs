---
title: Organization filesystem access policies
linkTitle: Filesystem
weight: 20
description: Control which host paths Docker Sandboxes can mount as workspaces with organization filesystem policies.
keywords: docker sandboxes, filesystem access, filesystem mount rules, workspace mount, organization policy, governance
aliases:
  - /ai/sandboxes/governance/access-controls/filesystem/
---

{{< summary-bar feature_name="AI Governance" >}}

Filesystem access policies are organization policies that control which host
paths a sandbox can mount as a workspace. When an organization policy is
enforced for a developer, a sandbox can mount a host path only if a filesystem
rule allows it, so a developer can't create a sandbox for a workspace your
policies don't allow. Each policy contains one or more filesystem mount rules
that restrict sandbox workspaces to approved directories. These rules control
which host paths a sandbox can mount. They don't control access to individual
files inside an allowed mount.

This page explains what a filesystem rule can match, how organization rules
change what developers can mount, and when your changes reach developers.

To create a filesystem access policy, see
[Create and manage organization policies](create-manage.md). For path patterns
and the `read` and `write` actions, see
[Filesystem mount rule syntax](../../reference/filesystem-rules.md).

## What a filesystem rule can match

Each filesystem mount rule allows or denies one or both of two actions on the
host paths it lists:

- `read` lets a sandbox mount the path as a read-only workspace.
- `write`, together with `read`, lets a sandbox mount the path as a writable
  workspace.

A writable workspace needs both a `read` and a `write` rule that allow it. A
read-only workspace needs only `read`.

Paths can use wildcards, and `~` expands to each user's home directory on
every platform, so a single `~/**` rule matches every developer's home tree.
For the path formats on each operating system and how wildcards match, see
[Filesystem mount rule syntax](../../reference/filesystem-rules.md).

## Built-in filesystem rules

Without an organization policy, a sandbox can mount any host path the developer
chooses. Developers can't add filesystem rules with local policy, which manages
network rules only. Local policy does include two built-in rules that allow
read and write access to every path, which developers see as `default-fs-*`
rules in `sbx policy ls --wide`.

When an organization policy is enforced for a developer, organization filesystem
mount rules decide which paths a sandbox can mount instead, and those built-in
rules stop applying. For more details, see
[How organization policy relates to local and kit rules](_index.md#how-organization-policy-relates-to-local-and-kit-rules).

## When changes take effect

Filesystem policy is checked when a sandbox is created and each time it starts.
Once a policy change has
[synced to the developer's machine](_index.md#when-policy-changes-reach-developers),
it applies the next time a sandbox starts. A sandbox that's already running
keeps its mount until it stops.
