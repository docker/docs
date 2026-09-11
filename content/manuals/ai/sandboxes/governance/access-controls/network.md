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

When organization governance is active, only organization allow rules grant
network access. Local allow rules are inactive until organization governance no
longer applies, while local deny rules still apply on top of the organization
policy. See [Precedence](../concepts.md#precedence).

## Rule syntax

Network rules use the action `connect:tcp`. Resources are hostnames, CIDR
ranges, ports, or hostnames with ports. The governance policy schema also
accepts `connect:udp`, but Docker Sandboxes always blocks direct external UDP
and ICMP. `connect:udp` rules have no effect.

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

Organization network rules belong to policies that can apply to the whole
organization or to selected teams. For setup steps and team scoping, see
[Organization policies](organization.md).

Use [Monitoring policies](../monitor-and-enforce/monitoring.md) to inspect
which network rules are active on a developer machine.

## Approval-required access

An organization network policy can require approval instead of granting access
outright. Destinations the policy allows aren't reachable until the developer
confirms them, which keeps an allowlist broad enough to be usable while still
putting a person in front of each destination an agent reaches for.

Approval is a property of the policy rather than of individual rules, so
turning it on applies it to every allow rule in that policy. A destination stays
directly reachable only when every policy that allows it is one that doesn't
require approval. If a policy that requires approval also matches, the request
needs approval regardless of what the other policies allow.

Only a destination the developer has already approved satisfies the
requirement. Preset rules, [kit-defined rules](../../customize/kits.md#control-network-access),
and rules the developer added with `sbx policy allow network` don't answer it.
An approval also can't reach a destination the organization doesn't allow at
all, and it can't override a deny rule. To withdraw a destination, add a deny
rule, which takes precedence over any approval already recorded.

Requiring approval on a policy is available only with organization
governance. To turn it on, see
[Organization policies](organization.md#require-approval-for-a-network-policy).

### Respond to an approval request

When an organization policy requires approval, a sandbox can't reach a
destination until you confirm it. The request is blocked and the sandbox
receives a message naming the destination:

```plaintext
Approval required for api.example.com.

Review and respond with:
  sbx policy approval ls
```

If your organization
[configures a support message](organization.md#configure-a-support-message), it
appears after the approval instructions.

The request that triggers the prompt doesn't wait for an answer. It's denied,
and approving the destination affects later requests. Agents that retry a
failed request pick up the new access on their next attempt. For others, run
the operation again.

List the destinations waiting for a response:

```console
$ sbx policy approval ls
APPROVAL network:c2FuZGJveA  (sandbox: my-sandbox)
  api.example.com:443
  Protocol: TCP
  Resource type: domain
  approval required by policy "default network"

  OPTION    LABEL
  allow     Allow
  dismiss   Dismiss
```

Each entry names the destination, the sandbox that asked for it, and the policy
that requires the approval. To inspect a single entry, pass its ID to
`sbx policy approval inspect`.

Respond by selecting one of the options the entry offers:

```console
$ sbx policy approval respond network:c2FuZGJveA --option allow
Recorded: Allow
```

Choosing `allow` grants access to that destination. Choosing `dismiss` leaves
it blocked, and the destination is requested again the next time the sandbox
tries to reach it. Repeated attempts collapse into a single entry, so a sandbox
retrying in a loop doesn't produce a queue of duplicates.

#### What approving grants

Approving records a rule that allows the destination the sandbox actually
asked for, scoped to the sandbox that asked. Three things follow from that:

- The rule covers one destination, not the pattern the policy rule used. A
  policy that allows `*.example.com` with approval asks about
  `api.example.com` and `cdn.example.com` separately.
- A destination includes its port, so `api.example.com:443` and
  `api.example.com:8443` are approved separately.
- Another sandbox reaching the same destination asks again.

Approved destinations stay allowed until the rule is removed. List them with
`sbx policy ls --wide --created-via approval`, and remove one the same way as
any other local rule, with [`sbx policy rm network`](local.md#managing-rules).

Approvals live in the local policy store, so [`sbx policy reset`](local.md#resetting)
removes all of them along with your other local rules. Each destination is
requested again the next time a sandbox reaches it.

> [!NOTE]
> To manage Model Context Protocol (MCP) server registration and requests
> through Docker's MCP gateway, use [MCP access policies](mcp.md). These
> policies apply only to the gateway. Direct MCP connections from a sandbox
> don't use the gateway, but you can control access to remote MCP servers with
> network policy.
