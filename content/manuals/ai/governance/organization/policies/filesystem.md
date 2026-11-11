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

Filesystem access policies control which host paths a sandbox can mount as a
workspace. When an organization policy is enforced for a developer, a sandbox
can mount a path only if one of your filesystem mount rules allows it, so a
developer can't create a sandbox for a directory you haven't approved. The
rules work at the level of the mount. Once a path is mounted, they don't
control which files inside it an agent can read or change.

Without an organization policy, a sandbox can mount any host path the developer
chooses. Developers can't add filesystem mount rules with local policy, but it
includes two built-in rules that allow read and write access to every path,
which developers see as `default-fs-*` rules in `sbx policy ls --wide`. Under
an enforced organization policy, those built-in rules stop applying, and only
your filesystem mount rules decide which paths a sandbox can mount. For more
details, see
[How organization policy relates to local and kit rules](_index.md#how-organization-policy-relates-to-local-and-kit-rules).

This page explains what a filesystem mount rule can match and when your changes
reach developers. To create a filesystem access policy, see
[Create and manage organization policies](create-manage.md).

## What a filesystem mount rule can match

Each filesystem mount rule allows or denies `read`, `write`, or both on the
host paths it lists. A read-only workspace needs only `read`, and a writable
workspace needs both. Paths can use wildcards, and `~` expands to each user's
home directory, so a single `~/**` rule matches every developer's home tree.

POSIX path matching is case-sensitive, including on macOS. Windows drive paths
and UNC paths, including WSL UNC paths, match without regard to case. For path
formats and matching details, see
[Filesystem mount rule syntax](../../reference/filesystem-rules.md).

## When changes take effect

Filesystem policy is checked when a sandbox is created and each time it starts.
Once a policy change has
[synced to the developer's machine](_index.md#when-policy-changes-reach-developers),
it applies the next time a sandbox starts. A sandbox that's already running
keeps its mount until it stops.
