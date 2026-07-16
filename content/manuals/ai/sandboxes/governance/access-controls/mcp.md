---
title: MCP access policies
linkTitle: MCP access
weight: 50
description: Use Cedar-based MCP policies to control sandbox MCP server registration, tool calls, prompts, resources, and approval gates.
keywords: docker sandboxes, MCP policy, MCP access, Cedar policy, requireApproval, AI Governance
---

MCP policies control Model Context Protocol activity made available to a
sandbox through Docker's MCP gateway. Use them to govern server registration,
tool calls, gateway meta-tools, resources, prompts, and approval requirements.
To register MCP servers and connect them to sandboxes, see
[MCP gateway](../../mcp-gateway.md).

Unlike [network access policies](network.md) and
[filesystem access policies](filesystem.md), MCP policies are organization
policies written in Cedar. Docker defines the `MCP` namespace, including the
actions, resource types, attributes, and approval behavior that policies can
match. For Docker's MCP policy actions, resources, attributes, and context
fields, see the [MCP policy reference](../reference/mcp-policy.md). For Cedar
syntax and language semantics, see the
[Cedar documentation](https://docs.cedarpolicy.com/).

## How MCP policy works

Cedar evaluates each request against a principal, action, resource, and context.
For Docker MCP policies:

- Policy scope supplies the principal. Use organization or team scope instead
  of matching users, teams, tenants, or roles in the policy.
- Actions use the `MCP::Action` namespace, such as
  `MCP::Action::"invokeTool"`.
- Resources use MCP entity types, such as `MCP::Server`, `MCP::Tool`,
  `MCP::Resource`, `MCP::Prompt`, and `MCP::Primordial`.
- A matching `forbid` overrides any `permit`, including a permit with
  `@requireApproval`.
- When MCP policy enforcement is active, evaluation is fail closed: server
  registration and governed MCP requests are denied unless a matching `permit`
  allows them.

If MCP policy enforcement isn't active for a user, the MCP gateway doesn't
evaluate Cedar policy and MCP activity is allowed by the gateway. MCP doesn't
have a local preset equivalent to network policy.

MCP policies are enforced on the MCP gateway path, not by the sandbox network
proxy. During `sbx mcp add`, Docker Sandboxes evaluates the resolved server
definition against `MCP::Action::"register"` before storing the registration.
When a sandbox uses MCP, the gateway evaluates governed MCP requests before
tool calls, resource reads, prompt retrieval, and gateway meta-tool execution.

Tool and resource listings can include entries that policy later denies at use
time. Use the [MCP policy reference](../reference/mcp-policy.md) for the exact
actions and limitations.

## Start with a baseline

To allow all MCP activity while you build a narrower policy, use an actionless
`permit`:

```plaintext
permit (principal, action, resource);
```

This is useful as a temporary baseline. For production policies, prefer
patterns that limit access to approved servers, read-only tools, or
approval-gated tool calls.

## Allow read-only tools

Servers can declare tool annotations such as `readOnly` and `destructive`. To
allow tools that a server marks read-only:

```plaintext
permit (principal, action == MCP::Action::"invokeTool", resource)
when { resource.readOnly == true };
```

`readOnly` defaults to `false` for tools that don't declare it, so this pattern
fails closed for unannotated tools.

## Require approval for writes

Use `@requireApproval` on a `permit` statement to require user approval before
a matching request runs:

```plaintext
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

```plaintext
permit (principal, action == MCP::Action::"invokeTool", resource)
when {
  resource.readOnly == true &&
  (resource in MCP::Server::"github" || resource in MCP::Server::"notion")
};
```

Server names must match the registered MCP server names.

## Block destructive tools

Use `forbid` for controls that must always win over permits:

```plaintext
permit (principal, action == MCP::Action::"invokeTool", resource);

forbid (principal, action == MCP::Action::"invokeTool", resource)
when { resource.destructive == true };
```

`destructive` defaults to `true` for tools that don't declare it, so this
pattern also blocks unannotated tools.

## Control server registration

Use the `register` action to control which MCP servers can be registered:

```plaintext
permit (principal, action == MCP::Action::"register", resource)
when { resource in MCP::Server::"github" };

permit (principal, action == MCP::Action::"register", resource)
when { resource in MCP::Server::"notion" };
```

Server names are chosen at registration time. To block a remote server
regardless of the name a developer chooses, match the server's
`resource.identityURL`:

```plaintext
permit (principal, action == MCP::Action::"register", resource);

forbid (principal, action == MCP::Action::"register", resource)
when { resource.identityURL == "https://mcp.example.com/mcp" };
```

This pattern prevents future `sbx mcp add` registrations for that identity URL.
It doesn't remove registrations that already exist or stop an already-loaded
server by itself. To govern existing registrations, add use-time rules for the
registered server name. Replace `example` with the name used in the existing
registration:

```plaintext
forbid (principal, action == MCP::Action::"invokeTool", resource)
when { resource in MCP::Server::"example" };

forbid (principal, action == MCP::Action::"readResource", resource)
when { resource in MCP::Server::"example" };

forbid (principal, action == MCP::Action::"getPrompt", resource)
when { resource in MCP::Server::"example" };
```

For remote server registration, match attributes the gateway provides, such as
server type, identity URL, OAuth requirement, or network requirement. Local
server command and argument attributes don't apply to remote servers.

For local stdio servers, `resource.type == "local-stdio"` matches host-side
local servers. Use `resource.command` and `resource.args` only when the
resolved registration includes command details.

## Next steps

- Use [Policy concepts](../concepts.md#mcp-policies) for the MCP policy model.
- Use the [MCP policy reference](../reference/mcp-policy.md) for exact action,
  resource, attribute, context, and approval behavior.
- Use [Organization policies](organization.md) to manage policy scope.
- Use [Audit logs](../monitor-and-enforce/audit.md) to collect policy decision
  records.
