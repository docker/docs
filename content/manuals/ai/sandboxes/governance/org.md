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
sandboxes can access. Organization policy moves that control to the admin level:
rules defined in the [Docker Admin Console](https://app.docker.com/admin) apply
uniformly to every sandbox in the organization. When organization governance is
active, it replaces local `sbx policy` rules entirely — local rules are no
longer evaluated and can't be used to supplement or override the organization
policy.

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

When organization governance is active, local network rules are not evaluated.
The organization policy is the only policy in effect. Local rules still appear
in `sbx policy ls` but with an `inactive` status. See [Monitoring](monitoring.md)
for how to read the rule view.

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

Within the active policy, deny rules beat allow rules. If a domain matches both,
it's blocked regardless of specificity. Outbound traffic is blocked unless a
rule allows it.

When organization governance is active, local rules are not evaluated. Only
organization rules set in the Admin Console determine what is allowed or
denied, and they can't be supplemented or overridden from a developer's machine.
The same model applies to filesystem policies: organization rules replace local
behavior entirely.

To unblock a domain when organization governance is active, update the rule in
the Admin Console or via the [API](/reference/api/ai-governance/). Without
organization governance, remove the local rule with `sbx policy rm`.

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
