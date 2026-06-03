---
title: Governance
weight: 55
description: Control what sandboxes can access, from local developer rules to org-wide enforcement.
keywords: docker sandboxes, governance, policy, network access, filesystem access, organization policy
---

Sandbox governance covers the policy system that controls what sandboxes can
access over the network and on the filesystem. It operates at two layers, and
only one applies at a time:

**Local policy** is configured per machine using the `sbx policy` CLI. It
lets individual developers customize which domains their sandboxes can reach.
See [Local policy](local.md).

**Organization policy** is configured centrally in the Docker Admin Console or
via the [Governance API](/reference/api/ai-governance/). Rules defined at the org level apply
uniformly across every sandbox in the organization. When organization
governance is active, it replaces local policy entirely: local `sbx policy`
rules are no longer evaluated. See [Organization policy](org.md).

> [!NOTE]
> Organization governance is available on a separate paid subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Learn more

- [Policy concepts](concepts.md): resource model, rule syntax, evaluation,
  and precedence
- [Local policy](local.md): configure network and filesystem rules on your
  machine with the `sbx policy` CLI
- [Organization policy](org.md): centrally manage sandbox policies across
  your organization from the Admin Console
- [Monitoring](monitoring.md): inspect active rules and monitor sandbox
  network traffic with `sbx policy ls` and `sbx policy log`
- [API reference](/reference/api/ai-governance/): manage org policies
  programmatically via the Governance API
