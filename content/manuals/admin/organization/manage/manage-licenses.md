---
title: Manage license assignment
linkTitle: License assignment
description: Manage product licenses for your organization, including assignment, revocation, and automatic assignment.
keywords: licenses, organization, members, Docker Core, Docker Offload, AI governance, license assignment, admin console
weight: 30
---

Licenses let you selectively choose which of your organization members have access to supported Docker products. Organization owners can oversee who on their team has active licenses, or configure licenses to assign automatically when members access supported Docker products. Like Docker Core seats, licenses can be configured on a per member basis. 

> [!TIP]
> To learn more about product licenses, Docker Core seats, and other Docker add-ons see [scale your subscription](/manuals/subscription/scale.md), 
> or <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_admin_licenses" class="link" rel="noopener">contact sales</a> to purchase licenses.

## Manage licenses

The **Members** page lets you track the number of available licenses for your organization and who currently holds a license. You can also assign or revoke licenses from this page. 

To manage licenses for your organization: 

1. Sign in to [Docker Home](https://app.docker.com), then choose your organization.
1. Select **Members** from the left navigation.
1. Select the action menu at the end of the row to assign or revoke an active license.
1. Optional. To bulk assign or revoke licenses, choose the members you want to bulk manage, then select the **Bulk actions** menu. 
1. Optional. To manage automatic license assignment, turn off or turn on with the **Automatically assign licenses** toggle. 

You must assign licenses manually, or configure automatic license assignment to consume a license. Inviting a new member to your organization consumes a seat or license if you select a product in **Licenses (optional)** during the [invite flow](/manuals/admin/organization/manage/members.md), but won't auto-assign product licenses by default. Conversely, purchasing a set of licenses won't trigger automatic assignment to existing members.

## Automatic license assignment

Automatic license assignment gives members a product license when they use a supported product for the first time. Automatic license assignment is available for AI Governance licenses. Only organizations that purchase AI Governance can set up auto-assignment for Docker Core as well. 

- When you purchase AI Governance, signing into [Docker Sandboxes](https://docs.docker.com/ai/sandboxes/) with `login` command in `sbx` CLI (`sbx login`) automatically provisions AI Governance licenses on a first-come, first served basis. 
- Similarly, logins to Docker Desktop will automatically provision Docker Core for AI Governance license-holding organizations that have available Docker Core seats.
- Licenses are assigned until exhausted. 
   - Once the available licenses are exhausted, automatic license assignment will stop until you purchase more licenses or revoke assigned licenses.
   - Members can still use Docker Sandbox or Docker Desktop, but organization policies for those products won't affect their usage. 

AI Governance licenses include single sign-on (SSO) and provisioning features regardless of your Docker Core subscription. Automatic license assignment requires [setting up SSO](/manuals/enterprise/security/single-sign-on/connect.md), then [provisioning with System for Cross-domain Identity Management (SCIM) or Just-in-Time (JIT)](/manuals/enterprise/security/provisioning/_index.md).

## What's next

See these docs to explore Docker Core add-ons, or products that need licenses:

- [Scale your subscription](/manuals/subscription/scale.md) to learn about different add-ons
- [Manage seats](/manuals/admin/organization/manage/manage-seats.md) to add more seats to your Docker Core subscription
- [AI Governance](/manuals/ai/sandboxes/governance/org.md) to set up organization policies for your organization members
- [Docker Offload](/manuals/offload/about.md) to let your developers offload building and running containers to the cloud
