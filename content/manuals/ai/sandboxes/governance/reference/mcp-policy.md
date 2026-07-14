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
MCP gateway.

Use this reference with [Control MCP tool access](../access-controls/mcp.md)
for common policy patterns. For the Cedar language, see the
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

An actionless `permit` covers every MCP action:

```cedar
permit (principal, action, resource);
```

## Actions

| Action              | Governs                   | Notes                                                                                                                  |
| ------------------- | ------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `register`          | MCP server registration   | Remote server registration needs an explicit `permit`. Use `forbid` rules to block local server registration patterns. |
| `invokeTool`        | MCP tool calls            | Most tool access policies target this action.                                                                          |
| `invokePrimordial`  | Gateway meta-tool calls   | Applies to gateway-injected meta-tools such as `mcp-add` and `code-mode`.                                              |
| `readResource`      | MCP resource reads        | Cedar policy is enforced inside the sandbox. Remote server reads aren't evaluated by Cedar policy.                     |
| `getPrompt`         | MCP prompt retrieval      | Remote prompt access is governed by gateway config, not Cedar policy.                                                  |
| `listTools`         | MCP tool listing          | Not Cedar-gated. Tool listings can include tools denied at invocation.                                                 |
| `listResources`     | MCP resource listing      | Not Cedar-gated.                                                                                                       |
| `subscribeResource` | MCP resource subscription | Not Cedar-gated.                                                                                                       |

## Resources

Match resources with the MCP entity type and attributes for the request.

| Entity            | Match with             | Notes                                                                                                                            |
| ----------------- | ---------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `MCP::Server`     | Registered server name | The server URL is the `resource.identityURL` attribute, not the entity ID.                                                       |
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

Tool annotation attributes are server-declared and advisory.

| Attribute                  | Applies to                | Notes                                                                           |
| -------------------------- | ------------------------- | ------------------------------------------------------------------------------- |
| `resource.name`            | Tools and prompts         | For tools, this is the bare tool name, not a display-prefixed name.             |
| `resource.uri`             | Resources                 | Use with string operators such as `like`.                                       |
| `resource.readOnly`        | Tools                     | Defaults to `false` when a tool doesn't declare it.                             |
| `resource.destructive`     | Tools                     | Defaults to `true` when a tool doesn't declare it.                              |
| `resource.idempotent`      | Tools                     | Server-declared tool annotation.                                                |
| `resource.openWorld`       | Tools                     | Server-declared tool annotation.                                                |
| `resource.type`            | Server registration       | Use for remote server registration rules.                                       |
| `resource.identityURL`     | Server registration       | The server URL. Use for remote server registration rules.                       |
| `resource.requiresOAuth`   | Server registration       | Use for remote server registration rules.                                       |
| `resource.requiresNetwork` | Server registration       | Use for remote server registration rules.                                       |
| `resource.command`         | Local server registration | Local stdio server command. Remote servers don't provide this value.            |
| `resource.args`            | Local server registration | Local stdio server arguments. This is a set, so `.contains()` can match values. |

Use `like` for string attributes. In Cedar, `like` uses `*` as its wildcard,
matches the full string, treats `?` as a literal character, and treats `\*` as
a literal asterisk.

Use `.contains()` only on set attributes, such as `resource.args`. On string
attributes, use `like`.

## Context fields

| Field                  | Notes                                                                                          |
| ---------------------- | ---------------------------------------------------------------------------------------------- |
| `context.request_time` | Bound at each enforcement point. Use it for time-window rules.                                 |
| `context.args`         | Tool-call arguments. Bound inside the sandbox. Gateway enforcement doesn't receive this field. |

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
- `readResource` is evaluated by Cedar policy inside the sandbox, not for
  remote server reads.
- `getPrompt` for remote servers is governed by gateway config, not Cedar
  policy.
- Tool-call argument rules using `context.args` apply inside the sandbox only.
- Server command and argument rules using `resource.command` or `resource.args`
  apply to local stdio servers only.
- Principal-based rules don't take effect. Use organization and team policy
  scope to target users.
- Server groups and tool categories aren't supported in MCP policy. Reference
  servers and tools individually.
