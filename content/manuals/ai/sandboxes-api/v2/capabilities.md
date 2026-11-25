---
title: Local and cloud capabilities in API v2
description: Understand the shared API contract and differences between execution environments.
keywords: docker sandboxes, API, SDK, local sandboxes, cloud sandboxes
weight: 30
linkTitle: Local and cloud capabilities
---

Local and cloud sandboxes share an API contract, but they do not expose every
operation or accept every configuration option. Check the capabilities of the
endpoint your application uses before choosing a workflow.

The API provides capability discovery for served operations, optional inputs,
and limits. An operation appearing in the API specification does not establish
that a particular endpoint serves it. Account permissions and resource
availability can also affect a cloud request.

## Shared workflows

The local and cloud implementations support sandbox creation, listing,
inspection, starting, stopping, deletion, and synchronous command execution.
These operations form the basis of the [getting started guides](get-started/_index.md).

## Environment-specific workflows

Local workspaces can use directories on the machine running the daemon. Cloud
sandboxes use remote storage and require a different approach to getting project
files into the environment.

Interactive processes, file transfers, managed secrets and volumes, and managed
MCP gateways are not part of the shared local API subset. A feature available
through the `sbx` CLI is not necessarily exposed through the API.

<!-- Before publication: verify a capability matrix against the release builds
of both hosts. Cover process interaction, files, storage, credentials, MCP,
policies, timeouts, ports, snapshots, and SDK transport support. Document the
actual feature-discovery request only after its release contract is stable. -->
