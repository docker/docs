---
title: Network access policies
linkTitle: Network access
weight: 30
description: Control outbound network access from Docker Sandboxes with local and organization policy rules.
keywords: docker sandboxes, network access, network rules, governance, local policy, organization policy
---

Network access policies control outbound connections from sandboxes. Each
policy contains one or more rules that allow the domains, IP ranges, and ports a
workflow needs, or block destinations that should stay unavailable.

You can configure network access in two places:

- [Local policy](local.md), which applies to sandboxes on one developer machine
  when organization governance is not active.
- [Organization policies](organization.md), which apply centrally across an
  organization or to selected teams.

When organization governance is active, organization network rules replace
local rules. Local rules are inactive until organization governance no longer
applies.

## Rule syntax

Network rules use the actions `connect:tcp` and `connect:udp`. Resources are
hostnames, CIDR ranges, ports, or hostnames with ports.

Examples:

- `api.example.com`
- `*.example.com`
- `**.example.com`
- `example.com:443`
- `10.0.0.0/8`

For exact wildcard behavior and CIDR support, see
[Network rules](../concepts.md#network-rules).

## Local network rules

Use `sbx policy allow network` and `sbx policy deny network` to manage local
network rules:

```console
$ sbx policy allow network api.example.com
$ sbx policy deny network ads.example.com
```

For presets, sandbox-scoped rules, testing, and troubleshooting, see
[Local policy](local.md).

## Organization network rules

Admins manage organization network rules in the Docker Admin Console. A policy
can apply to the whole organization or to selected teams. For setup steps and
team scoping, see [Organization policies](organization.md).

Use [Monitoring policies](../monitor-and-enforce/monitoring.md) to inspect
which network rules are active on a developer machine.
