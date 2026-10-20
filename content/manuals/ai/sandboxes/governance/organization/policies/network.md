---
title: Network access policies
linkTitle: Network
weight: 20
description: Control outbound network access from Docker Sandboxes with local and organization policy rules.
keywords: docker sandboxes, network access, network rules, governance, local policy, organization policy
aliases:
  - /ai/sandboxes/governance/access-controls/network/
---

Network access policies control outbound connections from sandboxes. Each
policy contains one or more rules that allow the domains, IP ranges, and ports a
workflow needs, or block destinations that should stay unavailable. Rules can
also match the HTTP method and path of a request, so a policy can allow part of
an API without allowing all of it.

You can configure network access in two places:

- [Local policy](../../local/_index.md), which applies to sandboxes on one developer machine
  when organization governance is not active.
- [Organization policies](../_index.md), which apply centrally across an
  organization or to selected teams.

When organization governance is active, only organization allow rules grant
network access. Local allow rules are inactive until organization governance no
longer applies, while local deny rules still apply on top of the organization
policy. See [Precedence](_index.md#precedence).

## Rule syntax

Network rules use `connect:tcp` for TCP and `connect:udp` for UDP. Resources are
hostnames, CIDR ranges, ports, or hostnames with ports. UDP requires
[experimental outbound UDP](../../local/network-rules.md#allow-outbound-udp). ICMP is blocked.

Examples:

- `api.example.com`
- `*.example.com`
- `**.example.com`
- `example.com:443`
- `10.0.0.0/8`

For exact wildcard behavior and CIDR support, see
[Network rules](../../reference/network-rules.md).

## HTTP method and path rules

A network rule matches a destination, so it allows or blocks everything a
sandbox sends there. An HTTP rule narrows the match to specific HTTP methods
and URL paths on that destination, which lets a policy allow reads from an API
without allowing writes to it.

HTTP rules layer on top of network rules. A network allow is the baseline for
a destination and HTTP rules carve into it, while a network deny blocks the
destination outright and no HTTP allow can reopen it. For the pattern syntax
and the full matching table, see
[HTTP rules](../../reference/network-rules.md#http-method-and-path).

Configure them in either place:

- Organization policies, in the network rule composer in Docker Home. Set the
  rule **Type** to **HTTP**, then select the methods and path patterns. See
  [Add a network rule](#add-a-network-rule).
- Local policies, with `--method` and `--path` on `sbx policy`. See
  [HTTP method and path rules](../../local/network-rules.md#http-method-and-path-rules).

## Add a network rule

Each rule has an optional **Rule name**, a **Type** that decides what the rule
matches, and a **Decision** of **Allow** or **Deny**.

- **HTTP** matches only HTTP requests with the methods and paths you specify.
  - In **Destination**, enter the host or IP address the rule covers. It
    matches any port unless you add one. Enter the destination with no scheme
    and no path, so `api.github.com` rather than
    `https://api.github.com/repos`. A local HTTP rule accepts only a host.
  - Under **HTTP methods**, select the methods the rule applies to. Use
    **Select all** to select every method, or **Read-only** to select `GET`,
    `HEAD`, and `OPTIONS`. A rule saved with no methods selected matches every
    method the composer lists. The composer doesn't list `CONNECT` or
    `TRACE`, which differs from the CLI, where `--method ANY` matches every
    HTTP method.
  - Under **Path patterns**, add one or more paths the rule covers, such as
    `/repos/*` and `/v1/**`. Leave it empty to match any path.
- **All traffic** matches every request to the destinations you list, on any
  port, method, and path.
  - Under **Protocols**, select **TCP**, **UDP**, or **Both**.
  - Under **Destinations**, add the hosts, IP addresses, or CIDR ranges the
    rule covers. A destination matches any port unless you add one, such as
    `example.com:8080`.

An HTTP rule's paths all belong to its one destination, so to cover paths on a
second host, add a second rule. For the pattern syntax and how HTTP rules
combine with **All traffic** rules, see
[HTTP rules](../../reference/network-rules.md#http-method-and-path).

## Approval-required access

An organization network policy can require approval instead of granting access
outright. Destinations the policy allows aren't reachable until the developer
confirms them, which keeps an allowlist broad enough to be usable while still
putting a person in front of each destination an agent reaches for.

Under organization governance, approval is a property of the policy rather
than of individual rules, so turning it on applies it to every allow rule in
that policy. A destination stays directly reachable only when no policy that
allows it requires approval. If a policy that requires approval also matches,
the request needs approval regardless of what the other policies allow.

Only a destination the developer has already approved satisfies the
requirement. Preset rules,
[kit-defined rules](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/network-policy@1.md),
and rules the developer added with `sbx policy allow network` don't answer it.
An approval also can't reach a destination the organization doesn't allow at
all, and it can't override a deny rule. To withdraw a destination, add a deny
rule, which takes precedence over any approval already recorded.

To require approval on an organization policy, see
[Organization policies](#require-approval-for-a-network-policy).

Network approval is separate from the MCP `@requireApproval` annotation. An MCP
approval confirms a single call within the session and creates no rule, while
an approved network destination stays allowed until you remove the rule. See
[MCP access policies](mcp.md).

### Require approval for a network policy

Turning on **Require approval before access** means the destinations a network
policy allows aren't reachable until the developer confirms each one. For how
approval behaves and what satisfies it, see
[Approval-required access](#approval-required-access).

To set it on an existing policy:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Network access**.
1. In the policy list, open the policy's action menu (⋮) and select **Edit**.
1. Turn on **Require approval before access**.
1. Select **Save**.

The policy's detail page reports approval as **Required** or **Not required**.
Editing a policy replaces it in full, so turning the setting off removes the
requirement from every rule in that policy.
