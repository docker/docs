---
title: Control MCP tool access
linkTitle: MCP tool access
weight: 50
description: Use Cedar-based MCP policies to control sandbox MCP server registration, tool calls, prompts, resources, and approval gates.
keywords: docker sandboxes, MCP policy, MCP tool access, Cedar policy, requireApproval, AI Governance
---

MCP access rules control Model Context Protocol activity made available to a
sandbox through Docker's MCP gateway. Use them to govern server registration,
tool calls, gateway meta-tools, resources, prompts, and approval requirements.

Unlike [network access rules](network.md) and
[filesystem access rules](filesystem.md), MCP access rules are organization
policies written in Cedar. Docker defines the `MCP` namespace, including the
actions, resource types, attributes, and approval behavior that policies can
match. For Docker's MCP policy actions, resources, attributes, context fields,
and approval behavior, see the
[MCP policy reference](../reference/mcp-policy.md). For Cedar syntax and
language semantics, see the [Cedar documentation](https://docs.cedarpolicy.com/).

## How MCP policy works

Cedar evaluates each request against a principal, action, resource, and context.
For Docker MCP policies:

- Policy scope supplies the principal. Use organization or team scope instead
  of matching users, teams, tenants, or roles in the policy.
- Actions use the `MCP::Action` namespace, such as
  `MCP::Action::"invokeTool"`.
- Resources use MCP entity types, such as `MCP::Server`, `MCP::Tool`,
  `MCP::Resource`, `MCP::Prompt`, and `MCP::Primordial`.
- Governed MCP activity is default deny. A request is blocked unless a
  matching `permit` allows it.
- A matching `forbid` overrides any `permit`, including a permit with
  `@requireApproval`.

## Start with a baseline

To allow all MCP activity while you build a narrower policy, use an actionless
`permit`:

```cedar
permit (principal, action, resource);
```

This is useful as a temporary baseline. For production policies, prefer
patterns that limit access to approved servers, read-only tools, or
approval-gated tool calls.

## Allow read-only tools

Servers can declare tool annotations such as `readOnly` and `destructive`. To
allow tools that a server marks read-only:

```cedar
permit (principal, action == MCP::Action::"invokeTool", resource)
when { resource.readOnly == true };
```

`readOnly` defaults to `false` for tools that don't declare it, so this pattern
fails closed for unannotated tools.

## Require approval for writes

Use `@requireApproval` on a `permit` statement to require user approval before
a matching request runs:

```cedar
permit (principal, action == MCP::Action::"invokeTool", resource)
when { resource.readOnly == true };

@requireApproval("write tool call")
permit (principal, action == MCP::Action::"invokeTool", resource)
when { resource.readOnly == false };
```

If the sandbox can't ask a user for approval, the request is denied. A matching
`forbid` still denies the request without showing an approval prompt.

## Limit tools to approved servers

Use `MCP::Server` to match the registered server that exposes a tool:

```cedar
permit (principal, action == MCP::Action::"invokeTool", resource)
when {
  resource.readOnly == true &&
  (resource in MCP::Server::"github" || resource in MCP::Server::"notion")
};
```

Server names must match the registered MCP server names.

## Block destructive tools

Use `forbid` for controls that must always win over permits:

```cedar
permit (principal, action == MCP::Action::"invokeTool", resource);

forbid (principal, action == MCP::Action::"invokeTool", resource)
when { resource.destructive == true };
```

`destructive` defaults to `true` for tools that don't declare it, so this
pattern also blocks unannotated tools.

## Control server registration

Use the `register` action to control which MCP servers can be registered:

```cedar
permit (principal, action == MCP::Action::"register", resource)
when { resource in MCP::Server::"github" };

permit (principal, action == MCP::Action::"register", resource)
when { resource in MCP::Server::"notion" };
```

For remote server registration, match attributes the gateway provides, such as
server type, identity URL, OAuth requirement, or network requirement. Local
server command and argument attributes don't apply to remote servers.

## Next steps

- Use [Policy concepts](../concepts.md#mcp-policies) for the MCP policy model.
- Use the [MCP policy reference](../reference/mcp-policy.md) for exact action,
  resource, attribute, context, and approval behavior.
- Use [Organization policies](organization.md) to manage policy scope.
- Use [Audit logs](../monitor-and-enforce/audit.md) to collect policy decision
  records.
