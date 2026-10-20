---
title: Network access policies
linkTitle: Network
weight: 20
description: Control outbound network access from Docker Sandboxes with organization network policies.
keywords: docker sandboxes, network access, network rules, HTTP rules, approval, governance, organization policy
aliases:
  - /ai/sandboxes/governance/access-controls/network/
---

Network access policies are organization policies that control outbound
connections from sandboxes. Each policy contains one or more rules that allow
the domains, IP ranges, and ports a workflow needs, or block destinations that
should stay unavailable. Rules can also match the HTTP method and path of a
request, so a policy can allow part of an API without allowing all of it.

To create a network access policy, see
[Create and manage policies](create-manage.md). For hostname patterns, CIDR
ranges, and how HTTP rules combine with network rules, see
[Network and HTTP rule syntax](../../reference/network-rules.md). Developers can
also add network rules on their own machines with
[local policy](../../local/network-rules.md). For how local and organization
rules combine, see [Precedence](_index.md#precedence).

## Network rule fields

Each rule has an optional **Rule name**, a **Type** that decides what the rule
matches, and a **Decision** of **Allow** or **Deny**.

- **HTTP** matches only HTTP requests with the methods and paths you specify.
  - In **Destination**, enter the host or IP address the rule covers. It
    matches any port unless you add one. Enter the destination with no scheme
    and no path, so `api.github.com` rather than
    `https://api.github.com/repos`.
  - Under **HTTP methods**, select the methods the rule applies to. Use
    **Select all** to select every method, or **Read-only** to select `GET`,
    `HEAD`, and `OPTIONS`. Select at least one method. A rule with no methods
    selected fails validation. The composer doesn't list `CONNECT` or `TRACE`,
    which differs from the CLI, where `--method ANY` matches every HTTP method.
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

Network approval is separate from the MCP `@requireApproval` annotation. An MCP
approval confirms a single call within the session and creates no rule, while
an approved network destination stays allowed until you remove the rule. See
[MCP access policies](mcp.md).

### Require approval for a network policy

Turning on **Require approval before access** means the destinations a network
policy allows aren't reachable until the developer confirms each one.

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

## When changes take effect

- Network policy is evaluated on every outbound request. Once a policy
  change has synced to the developer's machine (up to 5 minutes), it applies
  immediately to subsequent requests. HTTP rules are evaluated per request in
  the same way.

- An approval requirement applies from the point the policy change syncs.
  Destinations a developer already approved stay reachable, because the
  approval is recorded on the developer's machine. To withdraw one, add a deny
  rule. A deny takes precedence over a recorded approval.
