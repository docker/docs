---
title: Network policies
linkTitle: Policies
description: Control outbound network access from Docker Agentic Platform sandboxes.
keywords: docker agentic platform, network policies, network access, allow rules, deny rules
weight: 50
aliases:
  - /agentic-platform/concepts/policies/
  - /agentic-platform/guides/create-and-apply-policies/
---

Network policies control the external destinations that sandboxes can reach.
They are separate from MCP server connections and secret bindings.

Docker Agentic Platform uses two types of network policy:

- Kit policies are read-only policies for each sandbox type. The corresponding
  kit policy is applied automatically when you create a sandbox. These rules
  are the sandbox type's kit defaults. Review them under **Kit policies** on the
  **Policies** page.
- User policies are policies that you can select when you create a sandbox.
  Docker provides the read-only **Open** and **Balanced** presets, and you can
  create custom policies. **Open** allows all outbound destinations.
  **Balanced** allows a curated set of destinations.

You can select no user policies, one policy, or multiple policies. The selected
user policies are combined with the sandbox type's kit policy. Docker evaluates
all applicable rules, and a deny rule takes precedence over an allow rule.

If you select no user policies, only the kit policy applies. Network access is
default-deny, so the sandbox can reach only destinations that the kit policy
explicitly allows. If its kit policy has no network rules, all outbound
destinations are blocked.

## Policy rules

Policy rules allow or deny network destinations. Use the following hostname
patterns:

| Pattern               | Example           | Matches                                            |
| --------------------- | ----------------- | -------------------------------------------------- |
| Exact hostname        | `example.com`     | `example.com` only, not subdomains                 |
| Single-level wildcard | `*.example.com`   | One subdomain level, such as `api.example.com`     |
| Multi-level wildcard  | `**.example.com`  | Any depth: `api.example.com`, `v2.api.example.com` |
| Hostname with port    | `example.com:443` | `example.com` on port 443 only                     |

`example.com` does not match subdomains, and `*.example.com` does not match the
root domain. Add each pattern required by your destinations.

You can also match IPv4 and IPv6 CIDR ranges, such as `10.0.0.0/8`,
`192.168.1.0/24`, and `2001:db8::/32`.

When defining access, include every service the sandbox needs during startup
and operation. Depending on the workload, these services can include source
control hosts, package registries, and model providers.

## Create a policy

1. Open **Policies** and select **New policy**.
2. Enter a name that identifies the intended workload or access level.
3. Add allow and deny rules for the required destinations.
4. Review the rules and save the policy.

To apply the policy, select it under **Egress policy for this sandbox** when you
create a sandbox.

You can edit, copy, or delete a custom policy. Docker-managed policies cannot
be edited or deleted.

## Understand blocked access

An agent might report an HTTP 403 response, a connection failure, or another
service error when a required destination is not allowed. The error output from
the agent or tool is the primary source for identifying the destination.

For a sandbox that uses limited access, account for every host and port the
workload needs, including source control, package registries, model providers,
and supporting APIs. Add the narrowest allow rule that covers the required
destination. Use **Open** when broad outbound access is appropriate for the
workload.

Network policy controls outbound destinations. It does not grant MCP tools or
supply credentials. Configure those separately under **MCP** and **Secrets**.
