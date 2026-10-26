---
title: Troubleshoot governance
linkTitle: Troubleshooting
weight: 40
description: Resolve blocked requests, inactive rules, approval prompts, policy sync, and missing audit records for Docker Sandboxes governance.
keywords: docker sandboxes, governance, troubleshooting, network policy, sbx policy, approval, organization policy, audit logs
---

These problems come from the governance rules that control what sandboxes can
reach and mount, and from the audit logs that record those decisions. For other
Docker Sandboxes problems, see
[Troubleshooting](../sandboxes/troubleshooting.md).

## Agent can't install packages or reach an API

Sandboxes use [network rules](local/network-rules.md) to control
outbound traffic. If the agent fails to install packages or call an external
API, the destination probably isn't allowed. Check which requests are being
blocked:

```console
$ sbx policy log
```

Then allow the domains your workflow needs:

```console
$ sbx policy allow network "*.npmjs.org,*.pypi.org,files.pythonhosted.org"
```

To allow all outbound traffic instead:

```console
$ sbx policy allow network "**"
```

If `sbx policy allow` doesn't unblock the request, see
[A local allow rule has no effect](#a-local-allow-rule-has-no-effect).

## Every request from a sandbox is blocked

When policy enforcement can't get a valid organization policy, it fails closed
and blocks every network request and workspace mount, whatever your rules
allow. The message a blocked request returns tells you why:

- `Denied because governance could not be resolved`: run `sbx policy ls`. If
  the status line shows `Governance: Unresolved`, see
  [`sbx policy ls` shows Governance: Unresolved](#sbx-policy-ls-shows-governance-unresolved).
- `Denied because the policy snapshot is stale`: the daemon lost its connection
  to the governance service, or your Docker session expired. Check your network
  connection, then run `sbx login` to sign in again if the problem continues.

## A local allow rule has no effect

If a destination is still blocked after you add a local allow rule, an
organization policy is probably enforced for you. Run `sbx policy ls` and check
for `Governance: Managed by <org>` in the status line. When an organization
policy is enforced, local allow rules are inactive and can't loosen what the
organization policy allows.

To confirm, run `sbx policy ls --wide --include-inactive` and check that your
rule shows `inactive` in the `STATUS` column.

Ask your administrator to update the organization policy if you need access to
another destination. For more details, see
[How organization policy relates to local and kit rules](organization/policies/_index.md#how-organization-policy-relates-to-local-and-kit-rules).

## `sbx policy ls` shows Governance: Unresolved

Docker can't determine which organization governs your account. For example,
this happens when your account belongs to more than one organization with
governance enabled. A `Details` line below the status explains the cause.

Policy enforcement fails closed until the conflict is resolved, so local allow
rules can't grant access. Contact an administrator for the affected
organizations to resolve the conflicting governance configuration.

## A request is blocked with "Approval required"

The destination needs your confirmation. Either no allow or deny rule matches
it and your machine isn't under organization governance, or an organization
policy allows it but requires approval first.

Run `sbx policy approval ls` to see the pending request and respond to it with
`sbx policy approval respond`. Approving applies to later requests, not the one
that was blocked, so run the operation again afterward. For more details, see
[Respond to approval requests](local/approvals.md).

## An HTTP method or path is blocked on an allowed host

A host that a network rule allows can still have individual methods or paths
denied by an HTTP rule. Run `sbx policy ls --type http` to see which HTTP rules
apply. `sbx policy check network` reports the decision for the host only, so it
shows a host as allowed even when the specific request is denied. For more
details, see
[Allow or block HTTP methods and paths](local/network-rules.md#allow-or-block-http-methods-and-paths).

## A sandbox can't mount its workspace

If a sandbox fails to mount with a `mount policy denied` error, an organization
filesystem policy doesn't allow the path. Check that the filesystem allow rule
uses `**` rather than `*`. A single `*` doesn't match across directory
separators. For more details, see
[Organization filesystem access policies](organization/policies/filesystem.md).

## Organization policy changes aren't taking effect

After an administrator updates organization policies, changes take up to 5
minutes to reach developer machines. To apply them right away on a developer
machine, run `sbx policy refresh`. It fetches and applies the latest
organization policy without restarting the daemon or stopping running
sandboxes.

The refreshed policy applies to new connections only. A connection that a
sandbox opened while the previous policy allowed it stays open even if the new
policy denies it. Restart the sandbox to close those connections.

When a synced change takes effect also depends on the policy type. For more
details, see
[Organization network access policies](organization/policies/network.md#when-changes-take-effect),
[Organization filesystem access policies](organization/policies/filesystem.md#when-changes-take-effect),
and
[Organization MCP access policies](organization/policies/mcp.md#when-changes-take-effect).

## Audit records are missing

If an organization's audit logs don't include records from a developer, check
the following:

- The developer is governed. Records come only from signed-in users who have an
  AI Governance license and are governed by an enforced organization policy. On
  the developer's machine, run `sbx policy ls` and check for
  `Governance: Managed by <org>` in the status line.
- The developer runs Docker Sandboxes 0.39.0 or later. Run `sbx version` to
  check, and [install](../sandboxes/install.md) a newer version if needed.
- Cloud delivery is on. Records reach audit events in Docker Home and SIEM
  streaming only through cloud delivery. Organizations that used local audit
  logging before hosted audit logs were available start with cloud delivery
  off. For more details, see
  [Configure audit delivery](organization/audit/configure.md).
- Local files are finalized. With local delivery, the daemon writes records to
  a `.tmp` file and renames it to `.jsonl` after 5 minutes, 1000 events, or 50
  MiB, or when the daemon shuts down cleanly. For more details, see
  [Where records are stored](organization/audit/_index.md#where-records-are-stored).
