---
title: Manage access and security
linkTitle: Access and security
description: Manage MCP Gateway roles and resource permissions, block servers, and understand authentication and security boundaries.
keywords: mcp gateway, permissions, roles, secrets, security, governance
weight: 40
---

Use roles and resource permissions in the gateway console to delegate server
publishing and gateway management to teams. These controls govern management
of console resources. They are separate from signing in to upstream services
and from policies that govern requests at runtime.

The procedures below apply to deployments with permissions enabled. You need
superadmin membership to grant creation permissions. Ask your deployment
administrator to set up superadmin membership.

## Set up team access

A custom role binds to a group from your identity provider (IdP). Manage group
membership in the IdP. Owning or editing a role doesn't make you a member.

1. Open **Roles** and select **Create role**. You need permission to create
   roles or superadmin membership.
2. Enter a name, an optional description, and the **IdP group** exactly as it
   appears in the group claim. Select **Save role**.
3. As a superadmin, open the role's **Creation Permissions** tab.
4. Select which resources members can create: **Servers**, **Gateways**,
   **Roles**, or **Secrets**. Select **Save**.
5. Share existing resources with the role through their **Permissions**
   controls. Creation permissions don't grant access to existing resources.

For example, grant server creation to a developer role so its members can
publish servers. Share a gateway with a team role using **Read** and **Update**
so members can find it and maintain its server selection.

The built-in **Everyone** role includes all signed-in people and permits
gateway creation by default. To restrict creation to selected teams, clear
**Gateways** for **Everyone** and enable it for those teams' roles. A creation
grant from any role is sufficient, so clearing a box on one role doesn't
cancel a grant from another.

Treat permission to edit a role carefully: changing its IdP group changes who
gets access to resources shared with or owned by that role.

## Resource permissions and ownership

| Permission | What it grants |
| --- | --- |
| **Read** | View a resource. For secrets, view metadata but not the value. |
| **Update** | Change configuration. For secrets, replace the value. |
| **Delete** | Remove the resource. |
| **Read Secret Value** | Reveal a secret's value and use it in a server. |

Permissions are independent. Grant **Read** alongside **Update** or **Delete**
if the recipient needs to find the resource in the console.

Creators own their resources unless they choose an owning role. Members of
an owning role get owner rights. Only owners and superadmins can change
sharing or transfer ownership. **Update** doesn't grant those actions.

Access from ownership, direct grants, roles, and superadmin membership adds
up. Removing one grant doesn't remove access provided through another source.
Superadmins can manage resources throughout their tenant, including changing
ownership when an owner leaves.

For secrets, granting **Read** and **Read Secret Value** lets developers select
the secret for a server, but also lets them reveal its value. **Update** can
permit replacing a secret's value without reading it. Don't treat permission
to use a secret as a way to hide that secret from its user.

Resources managed in Git are read-only in the console. Ask your deployment
administrator to change them at their source.

## Block a server

The kill switch blocks calls through every gateway that includes the server,
and through its own endpoint. It doesn't stop the server itself.

1. Open the server's **Overview** tab. You need **Update** on the server.
2. In **Kill switch**, select **Activate kill switch**.
3. Enter a reason and confirm. Allow a few seconds for gateways to apply it.

To restore access, select **Deactivate kill switch** and confirm.
For a server managed in Git, superadmins can use **Activate emergency block**
instead. That block remains until a superadmin clears it, even if Git changes.

## Governance integration

Where configured, the gateway can enforce your Docker organization's MCP
policies and send audit records to a security information and event management
system (SIEM). Resource sharing in the console
doesn't configure those policies.

A superadmin configures the integration in **Platform settings**. It requires
governance enabled for the Docker organization, AI Governance seats, an
organization access token with the required scopes, and an `email` claim in
client access tokens. Audit delivery also requires a configured SIEM destination
and audit support in the deployment.

Use the console's **Platform settings** guide for setup requirements and steps.
Don't assume that enabling governance confirms audit delivery.

## Security boundaries

The gateway authenticates clients, exposes the selected servers' tools, and
holds upstream credentials for gateway-managed authentication. Signing in to
the console, permission to manage a resource, and authorization to an upstream
service are separate concerns. Sharing a gateway doesn't share personal
service accounts.

The gateway doesn't detect or block prompt injection in tool calls or results.
A tool result can contain instructions intended to manipulate your agent.
Use your client's tool approval controls and sandboxing where available, and
include only the servers needed for the task. Stop and inspect unexpected
actions after a tool call.
