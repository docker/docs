---
title: Network access and policy evaluation
linkTitle: Network access and policies
description: Distinguish local and cloud network controls from organization governance.
keywords: docker sandboxes, sbx, cloud sandboxes
weight: 50
---

Network policies determine which destinations a sandbox can reach. Configuring
a secret or connecting an MCP server does not replace those controls.

## Personal controls

Local CLI policies apply on your machine, with global and sandbox-specific
scopes. Cloud policies use Docker account and sandbox scopes. Local rules are
not copied to the cloud when you create or move a sandbox there.

An explicit deny takes precedence over an allow. A default-deny preset is
different: applicable sandbox or kit allow rules can still grant access.
Inspect the rules that apply to your sandbox when a connection behaves
differently from what you expect.

See [local CLI policies](/manuals/ai/sandboxes/cli/network-local.md),
[cloud CLI policies](/manuals/ai/sandboxes/cli/network-cloud.md), or
[Console policies](/manuals/ai/sandboxes/console/policies.md) for setup and inspection.

## Organization governance

Organization governance applies to local sandboxes in the documented product
scope. When it is active, organization allow rules grant access; local allow
rules cannot expand that access. Local deny rules can restrict it further.

Organization administrators can also govern filesystem access and MCP
operations. These controls are distinct from personal cloud account policies.
See [AI Governance](/manuals/ai/governance/_index.md) for administration.

Delivering audit records to a cloud service does not extend organization
policy enforcement to cloud sandboxes.

## Routing and proxies

Policy decides whether traffic is permitted. An upstream proxy determines how
permitted traffic leaves your local machine. Upstream proxy configuration is a
local feature; it is not a cloud network-policy setting.

See [Configure an upstream proxy](/manuals/ai/sandboxes/cli/local/upstream-proxy.md) for host
configuration, and [Local network isolation](/manuals/ai/sandboxes/concepts/isolation/isolation.md#network-isolation)
for enforcement details.
