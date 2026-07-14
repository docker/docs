---
title: Access controls
weight: 20
description: Configure local and organization controls for sandbox network, filesystem, and MCP access.
keywords: docker sandboxes, access controls, governance, network access, filesystem access, MCP access
---

Access controls define what sandboxes can reach or use. Choose a policy scope,
then choose the access surface you want to govern.

## Policy scope

- [Local access rules](local.md): configure network rules on a developer
  machine with the `sbx policy` CLI.
- [Organization policies](organization.md): manage centralized policies for an
  organization or team in the Docker Admin Console.

## Access surfaces

- [Network access rules](network.md): control outbound network access from
  sandboxes.
- [Filesystem access rules](filesystem.md): control which host paths sandboxes
  can mount as workspaces.
- [MCP tool access](mcp.md): control MCP server registration, tool calls,
  resources, prompts, and approval gates with Cedar policy.
