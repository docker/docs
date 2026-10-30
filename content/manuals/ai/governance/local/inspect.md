---
title: Inspect policy status and rules
linkTitle: Inspect policies
weight: 40
description: List the policies and rules that apply to your sandboxes, see where they come from, and check organization governance status with sbx policy ls.
keywords: docker sandboxes, sbx policy ls, sbx policy inspect, policy rules, inactive rules, organization governance, policy debugging
aliases:
  - /ai/sandboxes/governance/monitoring/
  - /ai/sandboxes/governance/monitor-and-enforce/monitoring/
---

Use `sbx policy ls` to see which policies and rules apply to your sandboxes and
where they come from, and `sbx policy inspect` to see one policy or rule in
full. A listing includes rules from your local policy, from kits, and from your
organization, so you can confirm the rules you added and work out why a request
is allowed or blocked. To watch network activity instead, see
[Test and monitor network policy](test-monitor.md).

## List the policies that apply

Run `sbx policy ls` to see every active policy:

```console
$ sbx policy ls
POLICY                                 SOURCE   APPLIES TO          SUMMARY
local-policy                           local    all                 network: 42 allow, 1 deny; filesystem read: 1 allow; filesystem write: 1 allow
1b2633ea-e604-48bb-a5e6-3ac86ba383fe   kit      sandbox:my-sandbox  network: 3 allow
```

The columns are:

- `POLICY`: the policy name.
- `SOURCE`: where the policy came from. `local` means your local policy, which
  holds your preset and the rules you added with `sbx policy`. `kit` means a
  [kit](../../sandboxes/customize/use-kits.md). `org` means your organization.
- `APPLIES TO`: which sandboxes the policy applies to. `all` means the policy
  is global, and `sandbox:<name>` scopes it to one sandbox.
- `SUMMARY`: a count of rule entries by type and decision, for example
  `network: 5 allow, 1 deny`. A rule that names several destinations
  contributes one entry per destination.

## List the policies for one sandbox

Pass a sandbox name to see only the global policies and the policies scoped to
that sandbox:

```console
$ sbx policy ls my-sandbox
```

The output summarizes the rules by type, for example
`allow 46 hosts from local-policy`. To see each rule, add `--wide`.

## Show rule details

Pass `--wide` to list each rule on its own row, with its rule ID, type,
decision, protocols, and resources:

```console
$ sbx policy ls --wide
```

Use the rule ID from this listing to remove a rule with
`sbx policy rm network --id`. For more details, see
[Remove a rule](network-rules.md#remove-a-rule).

To see one policy or rule in full, pass its name or rule ID to
`sbx policy inspect`:

```console
$ sbx policy inspect local-policy
```

## Filter the listing

Filters work with both the summary and the `--wide` listing:

- `--type`: show only `network`, `http`, or `filesystem` policies.
- `--source`: show only policies from `local`, `kit`, or `org`.
- `--decision`: show only `allow` or `deny` rules.
- `--protocol`: show only network rules for `tcp` or `udp`.
- `--created-via`: show only rules created a certain way. Pass `default` for
  preset rules, `added` for rules you added yourself, `provisioned` for rules a
  kit or application added, or `approval` for rules recorded when you approved
  a destination.

For example, to list the destinations you've approved:

```console
$ sbx policy ls --wide --created-via approval
```

For more details, see [Respond to approval requests](approvals.md).

To list the rules a kit added to a sandbox:

```console
$ sbx policy ls my-sandbox --source kit --wide
```

## List HTTP rules

Rules that match an HTTP method and path have the type `http`. To list only
those rules, pass `--type http`:

```console
$ sbx policy ls --type http
```

In a `--wide` listing, HTTP rules fill the `METHOD` and `PATH` columns. Rules
that match a whole destination show `-` in both columns.

In the `SUMMARY` column, HTTP rules count as network rules, and each part is
labeled by the network layer it matches on. `L4` counts entries that match a
whole destination, and `L7` counts entries that also match an HTTP method and
path:

```console
$ sbx policy ls
POLICY         SOURCE   APPLIES TO   SUMMARY
local-policy   local    all          network: 2 allow (L4), 1 deny (L7)
```

When the same decision has entries at both layers, each layer gets its own
count, L4 first. Two host allows and one HTTP allow read
`network: 2 allow (L4), 1 allow (L7)`.

The labels appear only when the listing includes at least one HTTP rule.
Because filters and hidden inactive rules change what the listing contains, a
filtered listing with no HTTP rules shows an unlabeled count, such as
`network: 42 allow`.

For more details, see
[HTTP method and path](../reference/network-rules.md#http-method-and-path).

## List filesystem mount rules

Filesystem mount rules control which host paths a sandbox can mount as a workspace.
To list only those rules, pass `--type filesystem`:

```console
$ sbx policy ls --type filesystem
POLICY         SOURCE   APPLIES TO   SUMMARY
local-policy   local    all          filesystem read: 1 allow; filesystem write: 1 allow
```

A writable workspace mount must be allowed by both a `filesystem:read` and a
`filesystem:write` rule. A read-only mount needs only `filesystem:read`. Your
local policy allows read and write access to all paths by default, shown as the
two allow entries above. For more details, see
[Filesystem mount rule syntax](../reference/filesystem-rules.md).

## Check organization governance status

When an organization policy is enforced for you, `sbx policy ls` starts with a
status line:

```console
$ sbx policy ls
Governance: Managed by my-org · Sync: OK, last synced 08:21:01 · Hidden: 9 inactive rules. Show with: sbx policy ls --include-inactive

POLICY               SOURCE   APPLIES TO   SUMMARY
default filesystem   org      all          filesystem read: 2 allow; filesystem write: 7 allow, 2 deny
default network      org      all          network: 38 allow, 4 deny
```

A `--wide` listing shows the same fields as a `POLICY RULES` block above the
table.

- `Governance` shows which organization manages the policy.
- `Sync` shows when the daemon last pulled the organization policy. If it shows
  an error or an old timestamp, run `sbx policy refresh` to fetch and apply the
  latest organization policy right away, without restarting the daemon or
  stopping sandboxes. The refreshed policy applies to new connections. To close
  existing connections that the updated policy denies, restart the sandbox.
- `Hidden` shows how many inactive rules the listing leaves out. For more
  details, see [Show inactive rules](#show-inactive-rules).

If the status shows `Governance: Unresolved`, see
[`sbx policy ls` shows Governance: Unresolved](../troubleshooting.md#sbx-policy-ls-shows-governance-unresolved).

## Show inactive rules

When an organization policy is enforced for you, local and kit allow rules
aren't evaluated, so `sbx policy ls` hides them. To list them, for example to
confirm which of your allow rules the organization policy overrides, pass
`--include-inactive` with `--wide`:

```console
$ sbx policy ls --wide --include-inactive
```

The `STATUS` column shows `active` or `inactive` for each rule. Inactive rules
have no effect while the organization policy is enforced. Local and kit deny
rules stay active and aren't hidden, because a deny still applies on top of the
organization policy. For more details, see
[How organization policy relates to local and kit rules](../organization/policies/_index.md#how-organization-policy-relates-to-local-and-kit-rules).
