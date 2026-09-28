---
title: Access controls
weight: 20
description: Configure local and organization controls for sandbox network, filesystem, and MCP access.
keywords: docker sandboxes, access controls, governance, network access, filesystem access, MCP access
aliases:
  - /ai/sandboxes/governance/access-controls/
---

Access controls are expressed as policies. Local and organization pages
describe where policies apply. Network and filesystem pages describe the rules
inside those policies. MCP policies use Cedar statements instead of the network
and filesystem rule format.

## Policy scope

- [Local policy](/manuals/ai/sandboxes/cli/network-local.md): configure network rules on a developer machine with
  the `sbx policy` CLI.
- [Organization policies](/manuals/ai/governance/access-controls/organization.md): manage centralized policies for an
  organization or team.

## Access surfaces

- [Network access policies](/manuals/ai/governance/access-controls/network.md): control outbound network access from
  sandboxes. A local policy rule can match a host, or an HTTP method and path.
- [Filesystem access policies](/manuals/ai/governance/access-controls/filesystem.md): control which host paths
  sandboxes can mount as workspaces.
- [MCP access policies](/manuals/ai/governance/access-controls/mcp.md): control MCP server registration, tool calls,
  resources, prompts, and approval gates with Cedar policy.
