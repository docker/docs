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

Use `sbx policy ls` to see all active rules and their current status:

```console
$ sbx policy ls
NAME                  TYPE      ORIGIN               DECISION   STATUS   RESOURCES
balanced-dev          network   local                allow      active   api.anthropic.com
ads-block             network   local                deny       active   ads.example.com
kit:my-sandbox        network   sandbox:my-sandbox   allow      active   api.example.com
kit:my-sandbox:deny   network   sandbox:my-sandbox   deny       active   telemetry.example.com
```

The columns are:

- `NAME`: the rule name.
- `TYPE`: the rule domain, such as `network`.
- `ORIGIN`: where the rule was configured. `local` means the rule is global
  and applies to all sandboxes. `sandbox:<name>` means the rule is scoped to
  the named sandbox. `remote` means the rule was set by your organization.
- `DECISION`: whether the rule allows or denies the resource.
- `STATUS`: whether the rule is in effect. A rule may be `inactive` if it's
  overridden or suppressed (for example, when organization governance is
  active, local rules are not evaluated and show as `inactive`).
- `RESOURCES`: the hosts or patterns the rule applies to.

When organization governance is active, the output starts with a governance
header showing which organization manages the policy and when it last synced:

```console
$ sbx policy ls
Governance: managed by my-org
[OK] last synced 13:54:21
NAME                  TYPE      ORIGIN               DECISION   STATUS     RESOURCES
balanced-dev          network   local                allow      inactive   api.anthropic.com
allow AI services     network   remote               allow      active     api.anthropic.com
                                                                           api.openai.com
allow Docker services network   remote               allow      active     *.docker.com
                                                                           *.docker.io
```

The governance header shows which organization is managing the policy and
confirms the daemon has successfully pulled the latest rules. If the sync
status shows an error or a stale timestamp, the daemon may not have the most
recent org policy. Run `sbx policy reset` to force a fresh pull.

Use `--type network` to show only network rules. Without a sandbox argument,
`sbx policy ls` shows every rule across all sandboxes. Pass a sandbox name to
filter to global rules and rules scoped to that sandbox:

```console
$ sbx policy ls my-sandbox
```

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
