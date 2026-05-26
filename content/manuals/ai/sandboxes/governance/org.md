---
title: Organization policy
linkTitle: Org policy
weight: 20
description: Centrally manage sandbox network and filesystem policies for your organization.
keywords: docker sandboxes, governance, organization policy, AI governance, admin console, network access, filesystem access
aliases:
  - /ai/sandboxes/security/governance/
---

[Local policies](local.md) give individual developers control over what their
sandboxes can access. Organization policy extends this to the admin level:
rules defined in the [Docker Admin Console](https://app.docker.com/admin) apply
uniformly to every sandbox in the organization, take precedence over local
`sbx policy` rules, and can't be overridden by individual users.

Admins can manage organization policies through the Admin Console UI or
programmatically using the [Governance API](/reference/api/ai-governance/).

> [!NOTE]
> Sandbox organization governance is available on a separate paid
> subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Network policies

### Configuring org-level network rules

Define network allow and deny rules in the Admin Console under
**AI governance > Network access**. Each rule takes a network target and an
action (allow or deny). You can add multiple entries at once, one per line.

For the full syntax reference (exact hostnames, wildcard subdomains, port
suffixes, and CIDR ranges), see [Policy concepts](concepts.md#network-rules).

### Delegate rules to local policy

When organization governance is active, local rules are ignored by default.
Only the organization policy is in effect. Admins can delegate a rule type
back to local policy by turning on the **User defined** setting for that
rule type in AI governance settings. Turning the setting on delegates the
rule type: local `sbx policy` rules of that type are evaluated alongside
organization rules, letting users add hosts to the allowlist from their own
machine.

When a rule type isn't delegated, local rules of that type still appear in
`sbx policy ls` but with an `inactive` status. See [Monitoring](monitoring.md)
for how to read the combined rule view.

Delegated local rules can expand access for domains the organization hasn't
explicitly denied, but can't override organization-level deny rules. This
applies to exact matches and wildcard matches alike; if the organization denies
`*.example.com`, a local allow for `api.example.com` has no effect because the
org-level wildcard deny covers it.

For example, given an organization policy that allows `api.anthropic.com`
and denies `*.corp.internal`:

- `sbx policy allow network -g api.example.com`: works, because the
  organization hasn't denied `api.example.com`
- `sbx policy allow network -g build.corp.internal`: no effect, because the
  organization denies `*.corp.internal`

#### Blocked values in delegated rules

To prevent overly broad rules from undermining the organization's policy,
certain catch-all values are blocked in delegated local rules:

- Domain patterns: `*`, `**`, `*.com`, `**.com`, `*.*`, `**.**`
- CIDR ranges: `0.0.0.0/0`, `::/0`

Scoped wildcards like `*.example.com` are still allowed. If a user attempts
to use a blocked value, `sbx policy` returns an error immediately.

## Filesystem policies

Filesystem policies control which host paths a sandbox can mount as
workspaces. By default, sandboxes can mount any directory the user has
access to.

Admins can restrict which paths are mountable by defining filesystem allow
and deny rules in the Admin Console under **AI governance > Filesystem
access**. Each rule takes a path pattern and an action (allow or deny).

For path pattern syntax including the difference between `*` and `**`, see
[Policy concepts](concepts.md#filesystem-rules).

## Precedence

Within any layer, deny rules beat allow rules. If a domain matches both, it's
blocked regardless of specificity. Outbound traffic is blocked unless a rule
allows it.

When organization governance is active, local rules are not evaluated. Only
organization rules set in the Admin Console determine what is allowed or
denied. Organization-level denials can't be overridden locally.

If the admin [delegates](#delegate-rules-to-local-policy) a rule type to
local policy by turning on the **User defined** setting, local rules of
that type are also evaluated alongside organization rules. Delegated local
rules can expand access for domains the organization hasn't explicitly
denied, but can't override organization-level denials.

The same model applies to filesystem policies: organization-level rules take
precedence over local behavior.

To unblock a domain, identify where the deny rule comes from. For local
rules, remove it with `sbx policy rm`. For organization-level rules, update
the rule in the Admin Console or via the [API](/reference/api/ai-governance/).

## Troubleshooting

### Policy changes not taking effect

After updating organization policies in the Admin Console, changes take up
to 5 minutes to propagate to developer machines. To apply changes
immediately, users can run `sbx policy reset`, which stops the daemon and
forces it to pull the latest organization policies on the next `sbx`
command.

> [!WARNING]
> `sbx policy reset` deletes all locally configured policy rules. The command
> prompts for confirmation before proceeding.

### Sandbox cannot mount workspace

If a sandbox fails to mount with a `mount policy denied` error, verify that
the filesystem allow rule in the Admin Console uses `**` rather than `*`. A
single `*` doesn't match across directory separators.
