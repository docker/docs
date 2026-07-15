---
title: MCP policy reference
linkTitle: MCP policy
weight: 20
description: Reference for Docker MCP policy actions, resources, attributes, context fields, and approval behavior.
keywords: docker sandboxes, MCP policy, Cedar policy, MCP actions, MCP resources, requireApproval, AI Governance
---

MCP policies are organization policies written in Cedar using Docker's `MCP`
namespace. This reference defines the Docker-specific policy surface for Model
Context Protocol (MCP) activity made available to sandboxes through Docker's
[MCP gateway](../../mcp-gateway.md).

Use this reference with [MCP access policies](../access-controls/mcp.md) for
common policy patterns. For the Cedar language, see the
[Cedar documentation](https://docs.cedarpolicy.com/).

## Evaluation model

Cedar evaluates MCP requests against a principal, action, resource, and context.
For Docker MCP policies, policy scope supplies the principal. Write policies
against the action, resource, and context. Clauses that reference principal
attributes, such as `principal in ...`, `principal.role`, or
`principal.tenant`, don't match.

Governed MCP activity is default deny. A request is blocked unless a matching
`permit` allows it. A matching `forbid` overrides any `permit`, including a
permit annotated with `@requireApproval`.

Default deny applies only to MCP activity that reaches Cedar evaluation. If MCP
policy enforcement isn't active for a user, the MCP gateway doesn't evaluate
Cedar policy and MCP activity is allowed by the gateway. If enforcement is
active but no effective MCP policy permits a request, the request is denied.

An actionless `permit` matches every MCP action that reaches Cedar evaluation:

```cedar
permit (principal, action, resource);
```

## Actions

| Action              | Governs                   | Notes                                                                                                |
| ------------------- | ------------------------- | ---------------------------------------------------------------------------------------------------- |
| `register`          | MCP server registration   | Remote server registration needs an explicit `permit`. Use server attributes to scope registration.  |
| `invokeTool`        | MCP tool calls            | Most tool access policies target this action.                                                        |
| `invokePrimordial`  | Gateway meta-tool calls   | Applies to gateway-injected meta-tools such as `mcp-add` and `code-mode`.                            |
| `readResource`      | MCP resource reads        | Rules match `MCP::Resource` and `resource.uri`.                                                      |
| `getPrompt`         | MCP prompt retrieval      | Rules match `MCP::Prompt` and `resource.name`.                                                       |
| `listTools`         | MCP tool listing          | Defined in the schema but not Cedar-gated. Tool listings can include tools denied at invocation.     |
| `listResources`     | MCP resource listing      | Defined in the schema but not Cedar-gated. Resource listings can include resources denied by policy. |
| `subscribeResource` | MCP resource subscription | Defined in the schema but not Cedar-gated.                                                           |

## Resources

Match resources with the MCP entity type and attributes for the request.

| Entity            | Match with             | Notes                                                                                                                            |
| ----------------- | ---------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `MCP::Server`     | Registered server name | The server URL or source is the `resource.identityURL` attribute, not the entity ID.                                             |
| `MCP::Tool`       | Bare tool name         | Use `resource.name`. Display prefixes aren't included. A bare-name match applies to every server exposing a tool with that name. |
| `MCP::Resource`   | Resource URI           | Use `resource.uri`.                                                                                                              |
| `MCP::Prompt`     | Prompt name            | Use `resource.name`.                                                                                                             |
| `MCP::Primordial` | Gateway meta-tool name | Match a specific primordial with an entity reference.                                                                            |

Examples:

```cedar
resource in MCP::Server::"notion"
resource.name == "move_file"
resource.uri like "*/docs/*"
resource in MCP::Primordial::"code-mode"
```

## Resource attributes

Tool annotation attributes come from MCP tool annotations or catalog metadata
and are advisory.

| Attribute                  | Applies to        | Notes                                                                                                                                                 |
| -------------------------- | ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `resource.name`            | Tools and prompts | For tools, this is the bare tool name, not a display-prefixed name.                                                                                   |
| `resource.uri`             | Resources         | Use with string operators such as `like`.                                                                                                             |
| `resource.readOnly`        | Tools             | Defaults to `false` when a tool doesn't declare it.                                                                                                   |
| `resource.destructive`     | Tools             | Defaults to `true` when a tool doesn't declare it.                                                                                                    |
| `resource.idempotent`      | Tools             | Defaults to `false` when a tool doesn't declare it.                                                                                                   |
| `resource.openWorld`       | Tools             | Defaults to `true` when a tool doesn't declare it.                                                                                                    |
| `resource.category`        | Tools             | Server or catalog category copied onto the tool resource. Tools don't self-declare categories.                                                        |
| `resource.type`            | Servers           | Use for server registration rules.                                                                                                                    |
| `resource.identityURL`     | Servers           | The server URL or source. Use for server registration rules.                                                                                          |
| `resource.requiresOAuth`   | Servers           | Use for server registration rules.                                                                                                                    |
| `resource.requiresNetwork` | Servers           | Use for server registration rules.                                                                                                                    |
| `resource.command`         | Servers           | Local stdio server command, such as `npx` or `docker`, when available. Empty for remote servers and registrations that don't include command details. |
| `resource.args`            | Servers           | Local stdio server arguments when available. This is a set, so `.contains()` can match values. Empty when no command details are available.           |

Use `like` for string attributes. In Cedar, `like` uses `*` as its wildcard,
matches the full string, treats `?` as a literal character, and treats `\*` as
a literal asterisk.

Use `.contains()` only on set attributes, such as `resource.args`. On string
attributes, use `like`.

## Context fields

| Field                  | Notes                                                                                             |
| ---------------------- | ------------------------------------------------------------------------------------------------- |
| `context.request_time` | Bound at each enforcement point. Use it for time-window rules.                                    |
| `context.oauth_scopes` | OAuth scopes for the caller. Present as a set, even when empty.                                   |
| `context.args`         | Tool-call arguments for `invokeTool`. Present when arguments are available as a supported object. |

Guard tool-call argument rules with `context has args` and a field check:

```cedar
permit (principal, action == MCP::Action::"invokeTool", resource)
when {
  resource.name == "approve_expense" &&
  context has args &&
  context.args has amount &&
  context.args.amount <= 500
};
```

A `permit` gated on missing arguments doesn't match, so the request falls to
default deny. A `forbid` gated on missing arguments doesn't match, so it
doesn't block the request.

Only object-shaped tool arguments are represented in `context.args`.
Unsupported, malformed, or too deeply nested arguments are omitted.

## Approval annotation

Use `@requireApproval("reason")` on a `permit` statement to require user
approval before a matching request runs:

```cedar
@requireApproval("write tool call")
permit (principal, action == MCP::Action::"invokeTool", resource)
when { resource.readOnly == false };
```

When a request matches the annotated `permit` and no `forbid` overrides it, the
policy engine returns an approval-required outcome. The annotation string is
shown as the approval reason.

Approval requires a client session that can present an MCP elicitation request
to the user. If the request can't be presented for approval, the request is
denied. Approval is an in-session confirmation, not an out-of-band approval
workflow.

Only the exact annotation name `@requireApproval` applies approval behavior.
Other annotation names, such as `@requireConsent` or `@requireConfirmation`,
don't require approval.

## Limitations

- Tool and resource listing actions aren't Cedar-gated. Listings can include
  entries that a policy denies when the sandbox tries to use them.
- Server command and argument rules using `resource.command` or `resource.args`
  apply only when the resolved server registration includes local stdio command
  details. Remote servers and metadata-resolved local servers can have empty
  values for those attributes. Use `resource.type` to match all local stdio
  servers.
- Principal-based rules don't take effect. Use organization and team policy
  scope to target users.
- Server groups aren't supported in MCP policy. Reference servers individually.
- Tool categories aren't self-declared by MCP tools. When available,
  `resource.category` is server or catalog metadata copied onto the tool
  resource.
