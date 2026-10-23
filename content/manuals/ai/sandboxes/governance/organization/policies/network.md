---
title: Organization network access policies
linkTitle: Network
weight: 10
description: Control outbound network access from Docker Sandboxes with organization network policies.
keywords: docker sandboxes, network access, network rules, HTTP rules, approval, governance, organization policy
aliases:
  - /ai/sandboxes/governance/access-controls/network/
---

{{< summary-bar feature_name="AI Governance" >}}

Network access policies are organization policies that control outbound
connections from sandboxes. When an organization policy is enforced for a
developer, their sandboxes can reach only the destinations your network rules
allow. Each policy contains one or more rules that allow the domains, IP
ranges, and ports a workflow needs, or block destinations that should stay
unavailable. Rules can also match the HTTP method and path of a request, so a
policy can allow part of an API without allowing all of it.

For destinations you want a person to confirm, a policy can require approval
instead of granting access outright. This page explains what a rule can match,
how approval works, and when your changes reach developers.

To create a network access policy, see
[Create and manage organization policies](create-manage.md). For hostname
patterns, CIDR ranges, and how HTTP rules combine with network rules, see
[Network and HTTP rule syntax](../../reference/network-rules.md). Developers can
also add network rules on their own machines with
[local policy](../../local/network-rules.md). For how local and organization
rules combine, see
[How organization policy relates to local and kit rules](_index.md#how-organization-policy-relates-to-local-and-kit-rules).

## What a network rule can match

Each rule has an optional **Rule name**, a **Type** that decides what the rule
matches, and a **Decision** of **Allow** or **Deny**. The type is either
**All traffic**, which covers whole destinations, or **HTTP**, which covers
specific requests to one destination.

### All traffic rules

An **All traffic** rule matches every request to the destinations you list, on
any port, method, and path.

- Under **Protocols**, select **TCP**, **UDP**, or **Both**.
- Under **Destinations**, add the hosts, IP addresses, or CIDR ranges the rule
  covers. A destination matches any port unless you add one, such as
  `example.com:8080`.

### HTTP rules

An **HTTP** rule matches only HTTP requests with the methods and paths you
specify.

- In **Destination**, enter the host or IP address the rule covers. It matches
  any port unless you add one. Enter the destination with no scheme and no
  path, so `api.github.com` rather than `https://api.github.com/repos`.
- Under **HTTP methods**, select the methods the rule applies to. Use
  **Select all** to select every method, or **Read-only** to select `GET`,
  `HEAD`, and `OPTIONS`. Select at least one method. A rule with no methods
  selected fails validation. The composer doesn't list `CONNECT` or `TRACE`,
  which differs from the CLI, where `--method ANY` matches every HTTP method.
- Under **Path patterns**, add one or more paths the rule covers, such as
  `/repos/*` and `/v1/**`. Leave it empty to match any path.

An HTTP rule's paths all belong to its one destination, so to cover paths on a
second host, add a second rule. For the pattern syntax and how HTTP rules
combine with **All traffic** rules, see
[HTTP rules](../../reference/network-rules.md#http-method-and-path).

## Approval-required access

An organization network policy can require approval instead of granting access
outright. Destinations the policy allows aren't reachable until the developer
confirms them, which keeps an allowlist broad enough to be usable while still
putting a person in front of each destination an agent reaches for. To turn
it on when you create or edit a policy, see
[Create and manage organization policies](create-manage.md).

Under organization governance, approval is a property of the policy rather
than of individual rules, so turning it on applies it to every allow rule in
that policy. A destination stays directly reachable only when no policy that
allows it requires approval. If a policy that requires approval also matches,
the request needs approval regardless of what the other policies allow.

Only a destination the developer has already approved satisfies the requirement.
Preset rules, kit-defined rules, and rules the developer added with `sbx policy
allow network` don't answer it. An approval also can't reach a destination the
organization doesn't allow at all, and it can't override a deny rule. To
withdraw a destination, add a deny rule, which takes precedence over any
approval already recorded.

Developers respond to approval requests from the command line. For what they
see and how they respond, see
[Respond to approval requests](../../local/approvals.md).

Network approval is separate from the MCP `@requireApproval` annotation. An MCP
approval confirms a single call within the session and creates no rule, while an
approved network destination stays allowed until you remove the rule. For more
details, see [Organization MCP access policies](mcp.md).

## When changes take effect

Network policy is evaluated on every outbound request. Once a policy change has
[synced to the developer's machine](_index.md#when-policy-changes-reach-developers),
it applies immediately to subsequent requests. HTTP rules are evaluated per
request in the same way.

An approval requirement applies from the point the policy change syncs.
Destinations a developer already approved stay reachable, because the approval
is recorded on the developer's machine. To withdraw one, add a deny rule. A
deny takes precedence over a recorded approval.
