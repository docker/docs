---
title: Connect MCP tools with API v2
description: Understand the managed cloud MCP gateway workflow for sandbox agents.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 50
linkTitle: Connect MCP tools
---

A managed MCP gateway connects an agent in a cloud sandbox to tools provided by
MCP servers. The cloud API can request a gateway when creating a sandbox,
inspect that gateway, and attach supported servers to it.

Attaching a server and authorizing access to that service are separate steps.
Do not assume that adding a server signs your account in to it.

This recipe covers the managed cloud gateway. The local API subset does not
expose gateway management. For the local CLI workflow, see
[MCP gateway](/manuals/ai/sandboxes/mcp-gateway.md).

<!-- Add request examples after the API and MCP documentation are stable.
Verify supported server sources, credential prerequisites, authorization paths,
and static gateway behavior. Cross-link the dedicated MCP documentation when
its destination is established. Do not infer local API support from sbx mcp. -->
