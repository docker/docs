---
title: Test and monitor network policy
linkTitle: Test and monitor
weight: 50
description: Test which destinations a sandbox can reach and monitor sandbox network traffic.
keywords: docker sandboxes, test network policy, sbx policy log, network traffic, local policy
---

## Testing policy

Before running a sandbox, you can check whether the current policy would allow
a network request with `sbx policy check network`:

```console
$ sbx policy check network api.anthropic.com
$ sbx policy check network blocked.example.com
```

The output reports `Allowed` or `Denied`, along with governance and context
details. A denied check returns a non-zero exit status.

The target can be a hostname, a `host:port` pair, an IP address, or a URL.
Bare hostnames and IP addresses are evaluated against port 443. This is useful
for verifying custom rules or checking what the Locked Down preset blocks
before you start an agent.

A check never creates an approval request. A destination that a sandbox would
ask you to approve shows as `Denied:`, with a `Reason:` line of
`no matching allow rule (default deny)`, or `approval required by policy` under
organization governance.

To check policy in the context of a specific sandbox:

```console
$ sbx policy check network --sandbox my-sandbox api.example.com
```

## Monitoring traffic

Use `sbx policy log` to see which hosts your sandboxes have contacted and
which rules matched:

```console
$ sbx policy log
BLOCKED REQUESTS
SANDBOX      TYPE     HOST                   PROXY        RULE            REASON         LAST SEEN        COUNT
my-sandbox   network  blocked.example.com    transparent  domain-blocked  default-deny   10:15:25 29-Jan  1

ALLOWED REQUESTS
SANDBOX      TYPE     HOST                   PROXY          RULE             REASON   LAST SEEN        COUNT
my-sandbox   network  api.anthropic.com      forward        domain-allowed            10:15:23 29-Jan  42
my-sandbox   network  registry.npmjs.org     forward-bypass domain-allowed            10:15:20 29-Jan  18
my-sandbox   network  app.example.com        browser-open                             10:15:10 29-Jan  1
```

The `PROXY` column shows how the request left the sandbox:

| Value            | Description                                                                                                    |
| ---------------- | -------------------------------------------------------------------------------------------------------------- |
| `forward`        | Routed through the forward proxy. Supports [credential injection](../../configuration/credentials.md).              |
| `forward-bypass` | Routed through the forward proxy without credential injection.                                                 |
| `transparent`    | Intercepted by the transparent proxy. Policy is enforced but credential injection is not available.            |
| `network`        | Non-HTTP traffic. TCP and experimental UDP egress follow network policy. ICMP is blocked. |
| `browser-open`   | A sandbox process requested opening a URL in the host browser. Policy is enforced before opening the URL.      |

The `RULE` column identifies the policy rule that matched the request. The
`REASON` column includes extra context when the daemon records one.

Filter by sandbox name by passing it as an argument:

```console
$ sbx policy log my-sandbox
```

Use `--limit N` to show only the last `N` entries, `--json` or `--format yaml` for
machine-readable output, or `--type network` to filter by policy type.
`sbx policy log` records network traffic only; filesystem mount decisions
aren't available in the log yet.

Structured output uses the same document for local and cloud sandboxes:
`{updated_at, allowed: [...], blocked: [...]}`. Each entry includes `sandbox_id`,
`host`, `proxy_type`, `rule`, `reason`, `count`, `first_seen`, and `last_seen`.
The `--limit` flag applies to structured output too.
