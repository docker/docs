---
title: Organization policy limits
linkTitle: Policy limits
weight: 35
description: Limits on the number and size of Docker Sandboxes organization policies and their rules.
keywords: docker sandboxes, organization policy, policy limits, rules per policy, AI Governance
---

Organization policies have the following limits, which help ensure fair usage
and resource availability across organizations:

| Limit                     | Value  | Applies to                                                                    |
|---------------------------|--------|-------------------------------------------------------------------------------|
| Policies per organization | 100    | All organization policies combined, whatever their type                       |
| Rules per policy          | 250    | Each network or filesystem policy. MCP policies don't have a rule limit.      |
| Policy size               | 400 KB | Each policy, whatever its type                                                |

The policies per organization limit is one shared total, not a separate
allowance for each policy type. For example, an organization with 60 network
policies, 30 filesystem policies, and 10 MCP policies has reached the limit.

The policy size limit applies to each policy as Docker stores it. This
includes the policy's content, such as its name, description, scope, and
rules or Cedar statements, plus a small amount of metadata that Docker stores
with every policy. Typical policies use only a small fraction of this limit.
Domain and file path values have no separate length limit beyond valid format.
Because MCP policies have no rule limit, the policy size limit is what caps how
many Cedar statements an MCP policy can hold.

If these limits don't fit your organization's needs,
[contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
to discuss options.
