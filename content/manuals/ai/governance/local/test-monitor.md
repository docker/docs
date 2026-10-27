---
title: Test and monitor network policy
linkTitle: Test and monitor
weight: 50
description: Check whether policy allows a destination, and view the network requests your sandboxes make, with sbx policy check and sbx policy log.
keywords: docker sandboxes, test network policy, sbx policy check, sbx policy log, network traffic, local policy
---

Use `sbx policy check` to find out whether policy allows a destination before a
sandbox tries to reach it, and `sbx policy log` to see the requests your
sandboxes made and which rules matched. Both work the same way whether the
rules come from your local policy, a kit, or your organization.

## Check whether a destination is allowed

Run `sbx policy check network` with the destination you want to test:

```console
$ sbx policy check network api.anthropic.com
$ sbx policy check network blocked.example.com
```

The output reports `Allowed` or `Denied`, along with governance and context
details. A denied check returns exit status 1, so you can use the command in
scripts.

The destination can be a hostname, a `host:port` pair, an IP address, or a URL.
Hostnames and IP addresses without a port are checked on port 443. A URL
supplies only its host and port, so the check ignores its path.

To check a destination for one sandbox, including the rules scoped to it, pass
`--sandbox`:

```console
$ sbx policy check network --sandbox my-sandbox api.example.com
```

A check never creates an approval request. A destination that a sandbox would
ask you to approve shows as `Denied`, with a `Reason` line of
`no matching allow rule (default deny)`, or `approval required by policy` under
organization governance.

A check reports the decision for the host, not for an HTTP method or path. To
test an HTTP rule, send a request from inside the sandbox. For more details,
see
[Allow or block HTTP methods and paths](network-rules.md#allow-or-block-http-methods-and-paths).

## View network activity

Run `sbx policy log` to see which hosts your sandboxes contacted and which
rules matched:

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

The `RULE` column identifies the policy rule that matched the request. The
`REASON` column includes extra context when the daemon records one. The `PROXY`
column shows how the request left the sandbox:

| Value            | Description                                                                                               |
| ---------------- | --------------------------------------------------------------------------------------------------------- |
| `forward`        | Routed through the forward proxy. Supports credential injection.                                          |
| `forward-bypass` | Routed through the forward proxy without credential injection.                                            |
| `transparent`    | Intercepted by the transparent proxy. Policy is enforced but credential injection is not available.       |
| `network`        | Non-HTTP traffic. TCP and experimental UDP egress follow network policy. ICMP is blocked.                 |
| `browser-open`   | A sandbox process requested opening a URL in the host browser. Policy is enforced before opening the URL. |

For how the forward and transparent proxies handle sandbox traffic, see
[Networking](../../sandboxes/architecture.md#networking).

The log records network requests only. Filesystem mount decisions aren't in
the log.

## Filter the log

To see only one sandbox's requests, pass its name:

```console
$ sbx policy log my-sandbox
```

To show only the most recent entries, pass `--limit`:

```console
$ sbx policy log --limit 20
```

## Get the log as JSON or YAML

Pass `--json` or `--format yaml` for output that scripts can read:

```console
$ sbx policy log --json
```

The document has the form `{updated_at, allowed: [...], blocked: [...]}`. Each
entry includes `sandbox_id`, `host`, `proxy_type`, `rule`, `reason`, `count`,
`first_seen`, and `last_seen`. The `--limit` flag applies to structured output
too.
