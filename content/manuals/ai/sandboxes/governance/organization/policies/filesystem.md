---
title: Filesystem access policies
linkTitle: Filesystem
weight: 30
description: Control which host paths Docker Sandboxes can mount as workspaces with organization filesystem policies.
keywords: docker sandboxes, filesystem access, filesystem mount rules, workspace mount, organization policy, governance
aliases:
  - /ai/sandboxes/governance/access-controls/filesystem/
---

Filesystem access policies are organization policies that control which host
paths a sandbox can mount as a workspace. Each policy contains one or more
filesystem mount rules that restrict sandbox workspaces to approved
directories. These rules control which host paths a sandbox can mount. They
don't control access to individual files inside an allowed mount.

To create a filesystem access policy, see
[Create and manage policies](create-manage.md). For path patterns and the
`read` and `write` actions, see
[Filesystem mount rule syntax](../../reference/filesystem-rules.md).

When an organization policy is enforced for you, organization filesystem mount
rules decide which paths a sandbox can mount, and the filesystem allow rules
from your local preset stop applying. See [Precedence](_index.md#precedence).

## When changes take effect

Filesystem policy is checked only when a workspace is mounted, which happens
when a sandbox is created. Changing the filesystem policy has no effect on a
running sandbox, which keeps access to the previously allowed path. To apply a
filesystem policy change, remove the sandbox and create a new one.

## Troubleshooting

### Sandbox cannot mount workspace

If a sandbox fails to mount with a `mount policy denied` error, verify that the
filesystem allow rule uses `**` rather than `*`. A single `*` doesn't match
across directory separators.
