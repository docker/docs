---
title: Access controls
weight: 20
description: Configure local and organization controls for sandbox network, filesystem, and MCP access.
keywords: docker sandboxes, access controls, governance, network access, filesystem access, MCP access
---

Access controls are expressed as policies. Local and organization pages
describe where policies apply. Network and filesystem pages describe the rules
inside those policies. MCP policies use Cedar statements instead of the network
and filesystem rule format.

## Policy scope

- [Local policy](local.md): configure network rules on a developer machine with
  the `sbx policy` CLI.
- [Organization policies](organization.md): manage centralized policies for an
  organization or team in the Docker Admin Console.

## Access surfaces

- [Network access policies](network.md): control outbound network access from
  sandboxes.
- [Filesystem access policies](filesystem.md): control which host paths
  sandboxes can mount as workspaces.
- [MCP access policies](mcp.md): control MCP server registration, tool calls,
  resources, prompts, and approval gates with Cedar policy.
