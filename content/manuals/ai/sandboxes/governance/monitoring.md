---
title: Monitoring policies
weight: 25
description: Inspect active policy rules and monitor sandbox network traffic with sbx policy ls and sbx policy log.
keywords: docker sandboxes, policy monitoring, sbx policy ls, sbx policy log, network traffic, policy debugging
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
  [kit](../customize/kits.md#control-network-access). `org` means your
  organization.
- `APPLIES TO`: which sandboxes the policy applies to. `all` means the policy
  is global. A sandbox name or profile name scopes it to specific sandboxes.
- `SUMMARY`: a count of rules by type and decision — for example,
  `network: 5 allow, 1 deny`.

To see full rule-level detail including rule IDs and resources, pass `--wide`.
To inspect a single policy or rule, use `sbx policy inspect`:

```console
$ sbx policy inspect Balanced
```

Use `--source` to filter by origin (`local`, `org`, or `kit`) and `--decision`
to filter by outcome (`allow` or `deny`).

A `STATUS` column also appears when you pass `--include-inactive`; see
[Showing inactive rules](#showing-inactive-rules).

When organization governance is active, the output starts with a summary line
showing which organization manages the policy, the sync state, and how many
inactive rules are hidden:

```console
$ sbx policy ls
Governance: Managed by my-org | Sync: OK, last synced 08:21:01 | Hidden: 9 inactive rules. Show with: sbx policy ls --include-inactive

POLICY               SOURCE   APPLIES TO   SUMMARY
default filesystem   org      all          filesystem read: 2 allow; filesystem write: 7 allow, 2 deny
default network      org      all          network: 38 allow, 4 deny
```

`Governance` shows which organization manages the policy, and `Sync` confirms
the daemon has pulled the latest rules. If the sync state shows an error or a
stale timestamp, the daemon may not have the most recent org policy. Run
`sbx policy reset` to force a fresh pull. `Hidden` reports how many inactive
rules are suppressed and how to reveal them.

### Showing inactive rules

When organization governance is active, local and kit-defined rules are not
evaluated, so `sbx policy ls` hides them by default. To list them too — for
example, to confirm which local rules the organization policy overrides — pass
`--include-inactive`. This adds a `STATUS` column:

```console
$ sbx policy ls --include-inactive
Governance: Managed by my-org | Sync: OK, last synced 08:41:06

POLICY                       SOURCE   APPLIES TO   SUMMARY                                                    STATUS
default filesystem           org      all          filesystem read: 2 allow; filesystem write: 7 allow, 2 deny   active
default network              org      all          network: 38 allow, 4 deny                                   active
default-fs-read-allow-all    local    all          filesystem read: 1 allow                                    inactive
default-fs-write-allow-all   local    all          filesystem write: 1 allow                                   inactive
```

Inactive policies show `inactive` in the `STATUS` column. They have no effect
while organization governance is active.

Use `--type network` or `--type filesystem` to show only policies of that type.
Without a sandbox argument, `sbx policy ls` shows every policy across all
sandboxes. Pass a sandbox name to filter to global policies and those scoped to
that sandbox:

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
[Policy concepts](concepts.md#filesystem-rules).

## Monitoring traffic

Use `sbx policy log` to see which hosts your sandboxes have contacted and
which rules matched:

```console
$ sbx policy log
Blocked requests:
SANDBOX      TYPE     HOST                   PROXY        RULE            REASON         LAST SEEN        COUNT
my-sandbox   network  blocked.example.com    transparent  domain-blocked  default-deny   10:15:25 29-Jan  1

Allowed requests:
SANDBOX      TYPE     HOST                   PROXY          RULE             REASON   LAST SEEN        COUNT
my-sandbox   network  api.anthropic.com      forward        domain-allowed            10:15:23 29-Jan  42
my-sandbox   network  registry.npmjs.org     forward-bypass domain-allowed            10:15:20 29-Jan  18
my-sandbox   network  app.example.com        browser-open                             10:15:10 29-Jan  1
```

The `PROXY` column shows how the request left the sandbox:

| Value | Description |
| ----- | ----------- |
| `forward` | Routed through the forward proxy. Supports [credential injection](../security/credentials.md). |
| `forward-bypass` | Routed through the forward proxy without credential injection. |
| `transparent` | Intercepted by the transparent proxy. Policy is enforced but credential injection is not available. |
| `network` | Non-HTTP traffic (raw TCP, UDP, ICMP). TCP can be allowed with a policy rule; UDP and ICMP are always blocked. |
| `browser-open` | A sandbox process requested opening a URL in the host browser. Policy is enforced before opening the URL. |

The `RULE` column identifies the policy rule that matched the request. The
`REASON` column includes extra context when the daemon records one.

Filter by sandbox name by passing it as an argument:

```console
$ sbx policy log my-sandbox
```

Use `--limit N` to show only the last `N` entries, `--json` for
machine-readable output, or `--type network` to filter by policy type.
`sbx policy log` records network traffic only; filesystem mount decisions
aren't available in the log yet.
