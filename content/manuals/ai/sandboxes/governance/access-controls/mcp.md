---
title: MCP access policies
linkTitle: MCP access
weight: 50
description: Use Cedar-based MCP policies to control sandbox MCP server registration, tool calls, prompts, resources, and approval gates.
keywords: docker sandboxes, MCP policy, MCP access, Cedar policy, requireApproval, AI Governance
---

MCP access policies let organization administrators control which Model
Context Protocol (MCP) servers developers can register and what agents can do
through Docker's MCP gateway. Use these policies to approve trusted servers,
withdraw access to a server, require approval for tool calls, and restrict
host-run servers. To register MCP servers and connect them to sandboxes, see
[MCP gateway](../../mcp-gateway.md).

Unlike [network access policies](network.md) and
[filesystem access policies](filesystem.md), MCP policies are organization
policies written in Cedar. Docker defines the `MCP` namespace, including the
actions, resource types, attributes, and approval behavior that policies can
match. This page focuses on representative access patterns. For Docker's exact
policy surface, see the [MCP policy reference](../reference/mcp-policy.md). For
Cedar syntax and language semantics, see the
[Cedar documentation](https://docs.cedarpolicy.com/).

## Govern the server lifecycle

MCP policy applies at two points in a server's lifecycle. A rule for one point
doesn't automatically govern the other.

| Admin decision                             | Evaluation point                            | Match with                                                                 |
| ------------------------------------------ | ------------------------------------------- | -------------------------------------------------------------------------- |
| Whether a server can be registered         | When a developer runs `sbx mcp add`         | The registered name and resolved server attributes, such as `identityURL`  |
| What agents can do through the MCP gateway | When the gateway handles a governed request | The registered server name, tool annotations, resource URI, or prompt name |

Registration rules affect future registrations. They don't remove a saved
registration or prevent an existing registration from being loaded with
`sbx mcp load`. Use-time rules govern tool calls, resource reads, and prompt
retrieval from servers that are already registered or loaded.

Server names are chosen during registration. Registration rules can match the
chosen name and resolved server identity together. At use time, tools,
resources, and prompts are associated with the registered name, so rules for an
existing server must match every name under which it was registered.

Built-in gateway tools, such as `code-mode` and OAuth authorization helpers,
are also governed at use time. They are `MCP::Primordial` resources rather than
tools associated with a registered server. For details, see
[Built-in gateway tools](../../mcp-gateway.md#built-in-gateway-tools).

Use-time policy doesn't hide or remove existing registrations. Tool and
resource listings can also include entries that policy denies when an agent
tries to use them.

## Choose an access posture

When MCP policy enforcement is active for a user, registration and governed MCP
requests are denied unless a matching `permit` allows them. A matching `forbid`
overrides any `permit`, including a permit with `@requireApproval`.

Use permits for an allowlist policy. For a blocklist policy that grants MCP
activity except for explicit restrictions, start with an actionless permit:

```plaintext
permit (principal, action, resource);
```

This statement permits every MCP action that reaches Cedar evaluation. Add
`forbid` statements for the restrictions the policy must enforce.

Policy scope supplies the principal. Use organization or team scope instead of
matching users, teams, tenants, or roles in Cedar. If MCP policy enforcement
isn't active for a user, the gateway doesn't evaluate Cedar policy and permits
MCP activity. MCP doesn't have a local preset equivalent to network policy.

## Approve a server

For an allowlist, approve both the server registration and its use-time
capabilities. The following policy approves a remote server only when it is
registered as `example` with the expected identity URL. It permits read-only
tool calls, resource reads, and prompt retrieval from that registered server:

```plaintext
permit (principal, action == MCP::Action::"register", resource)
when {
  resource in MCP::Server::"example" &&
  resource.identityURL == "https://mcp.example.com/mcp"
};

permit (principal, action == MCP::Action::"invokeTool", resource)
when {
  resource in MCP::Server::"example" &&
  resource.readOnly == true
};

permit (principal, action == MCP::Action::"readResource", resource)
when { resource in MCP::Server::"example" };

permit (principal, action == MCP::Action::"getPrompt", resource)
when { resource in MCP::Server::"example" };
```

Matching both the name and identity URL establishes a canonical registration.
It prevents a developer from registering another endpoint under the approved
name or registering the approved endpoint under another name. Remove the
resource or prompt permit if users don't need that capability.

## Require approval for non-read-only tools

Add an approval-gated permit to the previous policy when users should be able to
run other tools after confirmation:

```plaintext
@requireApproval("non-read-only tool call")
permit (principal, action == MCP::Action::"invokeTool", resource)
when {
  resource in MCP::Server::"example" &&
  resource.readOnly == false
};
```

Tool annotations are supplied by the server and are advisory. `readOnly`
defaults to `false` for tools that don't declare it, so this pattern requires
approval for unannotated tools. If the client session can't ask the user for
approval, the request is denied. A matching `forbid` also denies the request
without showing an approval prompt.

Use approval for in-session gateway requests. `sbx mcp add` can't present an
approval request, so a registration permit with `@requireApproval` results in a
denial.

## Withdraw server access

To withdraw access from a server that broader rules permit, address both
enforcement points. First, prevent future registrations of the server by
matching its identity URL:

```plaintext
forbid (principal, action == MCP::Action::"register", resource)
when { resource.identityURL == "https://mcp.example.com/mcp" };
```

Then deny use-time requests for each registered name that refers to the server:

```plaintext
forbid (principal, action == MCP::Action::"invokeTool", resource)
when { resource in MCP::Server::"example" };

forbid (principal, action == MCP::Action::"readResource", resource)
when { resource in MCP::Server::"example" };

forbid (principal, action == MCP::Action::"getPrompt", resource)
when { resource in MCP::Server::"example" };
```

The registration remains saved and can still be listed or loaded. These rules
prevent another registration for the identity URL and deny governed use under
the registered name. If the server was registered under other names, add
use-time rules for those names as well.

An OAuth authorization helper is a built-in gateway tool, not a child of the
registered server. To prevent agents from starting authorization for the
server, govern the helper separately:

```plaintext
forbid (principal, action == MCP::Action::"invokePrimordial", resource)
when { resource in MCP::Primordial::"example-authorize" };
```

## Restrict host-run servers

Local stdio servers run on the host, outside the sandbox VM. This includes
explicit host commands and OCI-packaged stdio servers started with host Docker.
For details about this boundary, see
[Docker Engine isolation](../../security/isolation.md#docker-engine-isolation).

In a blocklist policy that otherwise permits registration, deny both host-run
server types:

```plaintext
forbid (principal, action == MCP::Action::"register", resource)
when {
  resource.type == "local-stdio" ||
  resource.type == "container-stdio"
};
```

`local-stdio` covers explicit commands, including commands that start a Docker
container. `container-stdio` covers OCI-packaged stdio servers resolved from
registry or manifest metadata.

## Related information

- Use [Policy concepts](../concepts.md#mcp-policies) for the MCP policy model.
- Use the [MCP policy reference](../reference/mcp-policy.md) for exact action,
  resource, attribute, context, and approval behavior.
- Use [Organization policies](organization.md) to manage policy scope.
- Use [Audit logs](../monitor-and-enforce/audit.md) to collect policy decision
  records.
