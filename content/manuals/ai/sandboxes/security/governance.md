---
title: Organization governance
linkTitle: Org governance
weight: 35
description: Centrally manage sandbox network and filesystem policies for your organization.
keywords: docker sandboxes, governance, organization policy, AI governance, admin console, network access, filesystem access
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

This page covers how to configure organization policies in the Docker Admin
Console under AI governance settings. For local sandbox policies that
individual users configure on their own machine, see [Policies](policy.md).

Sandbox network and filesystem policies defined in the
[Docker Admin Console](https://app.docker.com/admin) apply uniformly to every
sandbox in the organization. Rules are enforced across all developers'
machines, take precedence over local `sbx policy` rules, and can't be
overridden by individual users. Admins can optionally
[delegate](#delegate-rules-to-local-policy) specific rule types back to local
control so developers can add additional allow rules.

> [!NOTE]
> Sandbox organization governance is available on a separate paid
> subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Network policies

### Configuring org-level network rules

Define network allow and deny rules in the Admin Console under
**AI governance > Network access**. Each rule takes a network target (domain,
wildcard, or CIDR range) and an action (allow or deny). You can add multiple
entries at once, one per line.

Rules support exact domains (`example.com`), wildcard subdomains
(`*.example.com`), and optional port suffixes (`example.com:443`).

`example.com` doesn't match subdomains, and `*.example.com` doesn't match
the root domain. Specify both to cover both.

### Delegate rules to local policy

When organization governance is active, local rules are ignored by default —
only the organization policy is in effect. Admins can delegate a rule type
back to local policy by turning on the **User defined** setting for that
rule type in AI governance settings. Turning the setting on delegates the
rule type: local `sbx policy` rules of that type are evaluated alongside
organization rules, letting users add hosts to the allowlist from their own
machine.

If a rule type isn't delegated, local rules of that type still appear in
`sbx policy ls` but with an `inactive` status and a note that the
organization hasn't delegated the rule type to local policy:

```console
$ sbx policy ls
NAME                  TYPE      ORIGIN               DECISION   STATUS                                                  RESOURCES
balanced-dev          network   local                allow      inactive — corporate policy takes precedence and does   api.anthropic.com
                                                                not delegate this rule type to local policy.
allow AI services     network   remote               allow      active                                                  api.anthropic.com
                                                                                                                        api.openai.com
allow Docker services network   remote               allow      active                                                  *.docker.com
                                                                                                                        *.docker.io
```

Organization rules show up with `remote` in the `ORIGIN` column.

Delegated local rules can expand access for domains the organization hasn't
explicitly denied, but can't override organization-level deny rules. This
applies to exact matches and wildcard matches alike; if the organization denies
`*.example.com`, a local allow for `api.example.com` has no effect because the
org-level wildcard deny covers it.

For example, given an organization policy that allows `api.anthropic.com`
and denies `*.corp.internal`:

- `sbx policy allow network -g api.example.com` — works, because the
  organization hasn't denied `api.example.com`
- `sbx policy allow network -g build.corp.internal` — no effect, because the
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

> [!CAUTION]
> Use `**` (double wildcard) rather than `*` (single wildcard) when writing
> path patterns to match path segments recursively. A single `*` only matches
> within a single path segment. For example, `~/**` matches all paths under
> the user's home directory, whereas `~/*` matches only paths directly
> under `~`.

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
the rule in the Admin Console.

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
