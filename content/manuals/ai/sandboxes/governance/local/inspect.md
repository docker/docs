---
title: Inspect policy status and rules
linkTitle: Inspect policies
weight: 40
description: Inspect active policy rules and monitor sandbox network traffic with sbx policy ls and sbx policy log.
keywords: docker sandboxes, policy monitoring, sbx policy ls, sbx policy log, network traffic, policy debugging
aliases:
  - /ai/sandboxes/governance/monitoring/
  - /ai/sandboxes/governance/monitor-and-enforce/monitoring/
---

`sbx policy ls` and `sbx policy log` give you a combined view of all active
policy rules and sandbox network activity, regardless of whether those rules
come from local configuration or organization governance. They're useful both
for verifying rules you've written and for debugging why a request is being
blocked or allowed.

## Listing rules

Use `sbx policy ls` to see all active policies and their current status:

```console
$ sbx policy ls
POLICY                                 SOURCE   APPLIES TO          SUMMARY
local-policy                           local    all                 network: 42 allow, 1 deny; filesystem read: 1 allow; filesystem write: 1 allow
1b2633ea-e604-48bb-a5e6-3ac86ba383fe   kit      sandbox:my-sandbox  network: 3 allow
```

The columns are:

- `POLICY`: the policy name.
- `SOURCE`: where the policy came from. `local` means your local configuration
  — a preset or rules you added with `sbx policy`. `kit` means a
  [kit](/manuals/ai/sandboxes/governance/organization/policies/_index.md#precedence). `org` means your
  organization.
- `APPLIES TO`: which sandboxes the policy applies to. `all` means the policy
  is global. `sandbox:<name>` scopes it to a single sandbox; a profile name
  scopes it to sandboxes using that profile.
- `SUMMARY`: a count of rule entries by type and decision, for example
  `network: 5 allow, 1 deny`. A rule that names several destinations
  contributes one entry per destination. When the listing includes rules that
  match an HTTP method and path, the network count labels each part `(L4)` or
  `(L7)`. See [HTTP rules](#http-rules).

To see full rule-level detail including rule IDs and resources, pass `--wide`.
To inspect a single policy or rule, use `sbx policy inspect`:

```console
$ sbx policy inspect Balanced
```

Use `--source` to filter by origin (`local`, `org`, or `kit`) and `--decision`
to filter by outcome (`allow` or `deny`).

Use `--protocol tcp` or `--protocol udp` to filter network rules.

Use `--created-via` to filter by how a rule was created. Pass `default` for
preset rules, `added` for rules you added yourself, `provisioned` for rules a
kit or application added, or `approval` for rules recorded when you approved a
destination. A wide listing shows the same information per rule:

```console
$ sbx policy ls --wide --created-via approval
```

See [Approval-required access](approvals.md).

A `STATUS` column also appears when you pass `--include-inactive`; see
[Showing inactive rules](#showing-inactive-rules).

When organization governance is active, the output starts with a summary block
showing which organization manages the policy, the sync state, and how many
inactive rules are hidden:

```console
$ sbx policy ls
POLICY RULES
  Governance: Managed by my-org
  Sync:       OK, last synced 08:21:01
  Hidden:     9 inactive rules. Show with: sbx policy ls --include-inactive

POLICY               SOURCE   APPLIES TO   SUMMARY
default filesystem   org      all          filesystem read: 2 allow; filesystem write: 7 allow, 2 deny
default network      org      all          network: 38 allow, 4 deny
```

`Governance` shows which organization manages the policy, and `Sync` confirms
the daemon has pulled the latest rules. `Hidden` reports how many inactive rules
are suppressed and how to reveal them.

If the sync state shows an error or a stale timestamp, run `sbx policy refresh`
to fetch and apply the latest policy without restarting the daemon or stopping
sandboxes. The policy applies to new connections. To close existing connections
that the updated policy denies, restart the sandbox.

If Docker can't determine which organization governs your account, policy
output shows `Governance: Unresolved`, and the dashboard shows the same
unresolved state. For example, this happens when your account belongs to
multiple organizations with governance enabled. Policy enforcement fails closed
until the conflict is resolved, so local allow rules can't grant access. Contact
an administrator for the affected organizations to resolve the conflicting
governance configuration.

### Showing inactive rules

When organization governance is active, local and kit-defined allow rules are
not evaluated, so `sbx policy ls` hides them by default. To list them too — for
example, to confirm which allow rules the organization policy overrides — pass
`--include-inactive`. This adds a `STATUS` column:

```console
$ sbx policy ls --include-inactive
POLICY RULES
  Governance: Managed by my-org
  Sync:       OK, last synced 08:41:06

POLICY                       SOURCE   APPLIES TO   SUMMARY                                                    STATUS
default filesystem           org      all          filesystem read: 2 allow; filesystem write: 7 allow, 2 deny   active
default network              org      all          network: 38 allow, 4 deny                                   active
default-fs-read-allow-all    local    all          filesystem read: 1 allow                                    inactive
default-fs-write-allow-all   local    all          filesystem write: 1 allow                                   inactive
```

Inactive policies show `inactive` in the `STATUS` column. They have no effect
while organization governance is active. Local and kit-defined deny rules stay
active and aren't hidden, because a deny still applies on top of the
organization policy. See [Precedence](../organization/policies/_index.md#precedence).

Use `--type network`, `--type filesystem`, or `--type http` to show only
policies of that type. Without a sandbox argument, `sbx policy ls` shows every
policy across all sandboxes. Pass a sandbox name to filter to global policies
and those scoped to that sandbox:

```console
$ sbx policy ls my-sandbox
```

### Filesystem rules

`sbx policy ls` lists filesystem policies alongside network policies. Filesystem
rules control which host paths a sandbox can mount as a workspace. Pass
`--type filesystem` to show only them:

```console
$ sbx policy ls --type filesystem
POLICY         SOURCE   APPLIES TO   SUMMARY
local-policy   local    all          filesystem read: 1 allow; filesystem write: 1 allow
```

A writable workspace mount must be allowed by both a `filesystem:read` and a
`filesystem:write` rule; a read-only mount needs only `filesystem:read`. The
default local policy allows read and write access to all paths, shown as the
two `default-fs-*` rules above. For the rule syntax and path patterns, see
[Policy concepts](../reference/filesystem-rules.md).

### HTTP rules

Rules that match an HTTP method and path are listed as type `http`. Pass
`--wide` to see the `METHOD` and `PATH` columns alongside network rules:

```console
$ sbx policy ls --wide
TYPE      METHOD   PATH
network   -        -
http      GET      /repos/org/project/**
http      POST     /admin/**
```

Rules that match a whole destination show `-` in both columns. To list only
HTTP rules, pass `--type http`.

HTTP rules are counted as network rules in the `SUMMARY` column, with each part
labeled by the network layer it matches on. `L4` counts entries that match a
whole destination, and `L7` counts those that also match an HTTP method and
path:

```console
$ sbx policy ls
POLICY         SOURCE   APPLIES TO   SUMMARY
local-policy   local    all          network: 2 allow (L4), 1 deny (L7)
```

When the same decision has entries at both layers, each layer gets its own
count, L4 first. Two host allows and one HTTP allow read
`network: 2 allow (L4), 1 allow (L7)`.

The labels appear when the current listing includes at least one HTTP rule.
Because filters and hidden inactive rules change what the listing contains, a
filtered listing with no HTTP rules shows an unlabeled count, such as
`network: 42 allow`.

For the rule syntax, see
[HTTP method and path](../reference/network-rules.md#http-method-and-path).

## Inspecting rules

To inspect which policies are active and where they come from, use
`sbx policy ls`. Use `--source` to filter by origin (`local`, `org`, `kit`),
`--decision` to filter by outcome (`allow`, `deny`), and `--wide` for
rule-level detail including rule IDs. To inspect a single policy or rule in
full, use `sbx policy inspect`.

## Troubleshooting

### Local allow rules have no effect

If rules you add with `sbx policy allow` don't change sandbox behavior, your
organization likely has governance enabled. Run `sbx policy ls` to check: if
the summary block contains `Governance: Managed by <org>`,
org governance is active. When it's active, local allow rules are inactive.
You can't use them to loosen restrictions the org policy imposes.

Inactive allow rules are hidden from `sbx policy ls` by default; run
`sbx policy ls --include-inactive` to see them with an `inactive` status in
the `STATUS` column.

When organization governance is active, only organization allow rules can grant
access. Ask your admin to update the organization policy if you need access to
an additional resource. Local deny rules remain active, so you can use
`sbx policy deny` to restrict access further.

### A domain is still blocked after adding an allow rule

If a domain remains blocked after you add a local allow rule, your organization
likely enforces governance, which makes local allow rules inactive. Run `sbx
policy ls` to check whether org governance is active; if the output starts with
a `Governance:` status line showing `Managed by <org>`, it is. Add
`--include-inactive` to confirm your rule shows an `inactive` status. If so,
the block can only be lifted by updating the org policy in Docker Home or via
the [API](/reference/api/ai-governance/).
