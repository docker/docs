---
title: Manage license assignment
linkTitle: License assignment
description: Manage product licenses for your organization, including assignment, revocation, and automatic assignment.
keywords: licenses, organization, members, Docker Offload, AI governance, license assignment, admin console
weight: 30
---

Licenses let you selectively choose which of your organization members have access to supported Docker products. Organization owners can oversee who on their team has active licenses, or configure licenses to assign automatically when members access supported Docker products. 

## Prerequisites

- You can purchase licenses for AI Governance or Docker Offload by [contacting sales](https://www.docker.com/pricing/contact-sales/). 
- To learn more about licenses, see [scale your subscription](/manuals/subscription/scale.md).

## Manage licenses

The **Members** page lets you track the number of available licenses for your organization and who currently holds a license. You can also assign or revoke licenses from this page. 

To manage licenses for your organization: 

1. Sign in to [Docker Home](https://app.docker.com), then choose your organization.
1. Select **Members** from the left navigation.
1. Select the action menu at the end of the row to assign or revoke an active license.
1. Optional. To bulk assign or revoke licenses, choose the members you want to bulk manage, then select the **Bulk** menu. 
1. Optional. To manage automatic license assignment, turn on or turn off on a per-product basis with the **Automatically assign licenses** toggle. 

You must assign licenses manually, or configure automatic license assignment to consume a license. Inviting a new member to your organization may consume a seat for your Docker Team or Docker Business subscription, but won't auto-assign product licenses by default. Conversely, purchasing a set of licenses doesn't automatically assign licenses to existing members.

## Automated license assignment

Automated license assignment gives members a product license when they use a supported product for the first time.

- For AI Governance and Docker Core licenses, invoking `sbx` or signing into Docker Desktop (respectively) triggers an event that provisions licenses on a first come, first serve basis.
- Licenses are assigned until exhausted. Once the available licenses are exhausted, automatic license assignment will stop until more licenses are purchased or assigned licenses are revoked.
- Members may still use Docker Sandbox or Docker Desktop, but won't be held to the policies of those product licenses. 

Automated license assignment requires [setting up SSO](/manuals/enterprise/security/single-sign-on/connect.md), then [provisioning with SCIM or JIT](/manuals/enterprise/security/provisioning/_index.md). 

> [!NOTE]
> AI Governance licenses include SSO and 
> provisioning features regardless of 
> your Docker Core subscription.  

## What’s next

To learn more about Docker products:

- Learn more about [AI Governance](/manuals/ai/sandboxes/security/governance.md).
- Learn more about [Docker Offload](/manuals/offload/about.md).
- Enable [SSO](/manuals/enterprise/security/single-sign-on/_index.md) and set up [provisioning](/manuals/enterprise/security/provisioning/_index.md).
