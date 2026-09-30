---
title: Network access policies
linkTitle: Network access
weight: 30
description: Control outbound network access from Docker Sandboxes with local and organization policy rules.
keywords: docker sandboxes, network access, network rules, governance, local policy, organization policy
---

The governance described here applies to local sandboxes. Cloud sandboxes
use separate network policy configuration. See
[Cloud network policy](../../cloud/network-policy.md) for cloud controls.

Network access policies control outbound connections from sandboxes. Each
policy contains one or more rules that allow the domains, IP ranges, and ports a
workflow needs, or block destinations that should stay unavailable. A local
policy rule can also match the HTTP method and path of a request, so it can
allow part of an API without allowing all of it.

You can configure network access in two places:

- [Local policy](local.md), which applies to sandboxes on one developer machine
  when organization governance is not active.
- [Organization policies](organization.md), which apply centrally across an
  organization or to selected teams.

When organization governance is active, only organization allow rules grant
network access. Local allow rules are inactive until organization governance no
longer applies, while local deny rules still apply on top of the organization
policy. See [Precedence](../concepts.md#precedence).

## Rule syntax

Network rules use `connect:tcp` for TCP and `connect:udp` for UDP. Resources are
hostnames, CIDR ranges, ports, or hostnames with ports. Docker Sandboxes
always blocks external ICMP.

Examples:

- `api.example.com`
- `*.example.com`
- `**.example.com`
- `example.com:443`
- `10.0.0.0/8`

For exact wildcard behavior and CIDR support, see
[Network rules](../concepts.md#network-rules).

## HTTP method and path rules

A network rule matches a destination, so it allows or blocks everything a
sandbox sends there. An HTTP rule narrows the match to specific HTTP methods
and URL paths on that destination, which lets a policy allow reads from an API
without allowing writes to it.

HTTP rules layer on top of network rules. A network allow is the baseline for
a destination and HTTP rules carve into it, while a network deny blocks the
destination outright and no HTTP allow can reopen it. For the pattern syntax
and the full matching table, see
[HTTP rules](../concepts.md#http-method-and-path).

Add them to a local policy with `--method` and `--path` on `sbx policy`. See
[HTTP method and path rules](local.md#http-method-and-path-rules).

## Local network rules

Use `sbx policy allow network` and `sbx policy deny network` to manage local
network rules:

```console
$ sbx policy allow network api.example.com
$ sbx policy deny network ads.example.com
```

Allow rules cover TCP only by default, while deny rules cover TCP and UDP by
default. See [Allow outbound UDP](local.md#allow-outbound-udp) for how to change either
with `--protocol`.

For presets, sandbox-scoped rules, testing, and troubleshooting, see
[Local policy](local.md).

## Organization network rules

Organization network rules belong to policies that can apply to the whole
organization or to selected teams. For setup steps and team scoping, see
[Organization policies](organization.md).

Organization rules apply only to the protocols selected for the rule. Docker
Home and the Governance API store this selection explicitly. The defaults for
local CLI rules don't apply.

Use [Monitoring policies](../monitor-and-enforce/monitoring.md) to inspect
which network rules are active on a developer machine.

> [!NOTE]
> To manage Model Context Protocol (MCP) server registration and requests
> through Docker's MCP gateway, use [MCP access policies](mcp.md). These
> policies apply only to the gateway. Direct MCP connections from a sandbox
> don't use the gateway, but you can control access to remote MCP servers with
> network policy.
