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
PROVENANCE   APPLIES_TO     POLICY/RULE                  TYPE               DECISION   RESOURCES
local        all            default-ai-services          network            allow      api.anthropic.com:443
                                                                                       api.openai.com:443
local        all            default-fs-read-allow-all    filesystem:read    allow      **
local        all            default-fs-write-allow-all   filesystem:write   allow      **
kit          sandbox:docs   kit:docs                     network            allow      api.github.com
                                                                                       registry.npmjs.org
```

The columns are:

- `PROVENANCE`: where the rule came from. `local` is a rule from your local
  policy — a preset default or one you added with `sbx policy`. `kit` is a rule
  added by a [kit](../customize/kits.md#control-network-access). `remote` is a
  rule set by your organization.
- `APPLIES_TO`: which sandboxes the rule applies to. `all` means the rule is
  global. `sandbox:<name>` means it's scoped to the named sandbox.
- `POLICY/RULE`: the rule's identity. Organization rules show as
  `<policy> / <rule>`. Local and kit rules show the rule name.
- `TYPE`: the rule domain. Network rules show as `network`. Filesystem rules
  show as `filesystem:read` or `filesystem:write`, depending on the access the
  rule controls.
- `DECISION`: whether the rule allows or denies the resource.
- `RESOURCES`: the hosts or patterns the rule applies to.

A `STATUS` column also appears when you pass `--include-inactive`; see
[Showing inactive rules](#showing-inactive-rules).

When organization governance is active, the output starts with a `Policy rules`
header showing which organization manages the policy, the sync state, and how
many inactive rules are hidden:

```console
$ sbx policy ls
Policy rules
------------
Governance  Managed by my-org
Sync        OK, last synced 08:21:01
Hidden      9 inactive rules. Show with: sbx policy ls --include-inactive

PROVENANCE   APPLIES_TO   POLICY/RULE                                      TYPE               DECISION   RESOURCES
remote       all          default filesystem / allow home subdirectories   filesystem:write   allow      ~/**
remote       all          default filesystem / deny home directory         filesystem:write   deny       ~/
remote       all          default network / allow AI services              network            allow      api.anthropic.com
                                                                                                         api.openai.com
remote       all          default network / allow Docker services          network            allow      *.docker.com
                                                                                                         *.docker.io
```

The `Governance` line shows which organization manages the policy, and `Sync`
confirms the daemon has pulled the latest rules. If the sync state shows an
error or a stale timestamp, the daemon may not have the most recent org policy.
Run `sbx policy reset` to force a fresh pull. The `Hidden` line reports how many
inactive rules are suppressed and how to reveal them.

### Showing inactive rules

When organization governance is active, local and kit-defined rules are not
evaluated, so `sbx policy ls` hides them by default. To list them too — for
example, to confirm which local rules the organization policy overrides — pass
`--include-inactive`. This adds a `STATUS` column:

```console
$ sbx policy ls --include-inactive
Policy rules
------------
Governance  Managed by my-org
Sync        OK, last synced 08:41:06

PROVENANCE   APPLIES_TO   POLICY/RULE                                      TYPE               DECISION   STATUS                        RESOURCES
local                     default-fs-read-allow-all                        filesystem:read    allow      inactive — corporate policy   **
                                                                                                         takes precedence and does
                                                                                                         not delegate this rule type
                                                                                                         to local policy.
local                     default-fs-write-allow-all                       filesystem:write   allow      inactive — corporate policy   **
                                                                                                         takes precedence and does
                                                                                                         not delegate this rule type
                                                                                                         to local policy.
remote       all          default filesystem / allow home subdirectories   filesystem:write   allow      active                        ~/**
remote       all          default filesystem / deny home directory         filesystem:write   deny       active                        ~/
```

Inactive rules show `inactive` in the `STATUS` column, along with the reason.
They have no effect while organization governance is active.

Use `--type network` or `--type filesystem` to show only rules of that type.
Without a sandbox argument, `sbx policy ls` shows every rule across all
sandboxes. Pass a sandbox name to filter to global rules and rules scoped to
that sandbox:

```console
$ sbx policy ls my-sandbox
```

### Filesystem rules

`sbx policy ls` lists filesystem rules alongside network rules. Filesystem
rules control which host paths a sandbox can mount as a workspace. Pass
`--type filesystem` to show only them:

```console
$ sbx policy ls --type filesystem
PROVENANCE   APPLIES_TO   POLICY/RULE                  TYPE               DECISION   RESOURCES
local        all          default-fs-read-allow-all    filesystem:read    allow      **
local        all          default-fs-write-allow-all   filesystem:write   allow      **
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
