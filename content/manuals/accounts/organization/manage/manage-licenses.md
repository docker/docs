---
title: Manage license assignment
linkTitle: License assignment
description: Manage product licenses for your organization, including invite-time assignment, revocation, and automatic assignment.
keywords: licenses, organization, members, invite, invitation, Docker Core, Docker Offload, AI Governance, license assignment, docker home
weight: 30
aliases:
  - /admin/organization/manage/manage-licenses/
---

Licenses let you choose which organization members can access supported Docker
products. Organization owners can manage active licenses for their members, or
configure licenses to assign automatically when members access supported Docker
products.

> [!TIP]
> To learn more about product licenses, Docker Core seats, and other Docker
> add-ons, see [Docker plans](/manuals/subscription-billing/plans/_index.md),
> or
> <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_admin_licenses" class="link" rel="noopener">contact sales</a>
> to purchase licenses.

## Licenses and invites

When you invite someone, you can select a product license to assign when they
accept. Docker doesn't reserve or deduct the license at invite time; assignment
happens on acceptance:

- If a license is available when they accept, Docker assigns it to them and the
  number of available licenses decreases by one.
- If no licenses remain when they accept, they still join your organization,
  but without a license.

> [!NOTE]
> Docker doesn't notify you or the invitee when a selected license is
> unavailable at acceptance.

Licenses aren't reserved for pending invitations, so they must still be
available when each invitee accepts. Monitor availability on the Members page
while invitations are pending.

### Select licenses when inviting

Selecting licenses when you invite is an alternative to assigning one manually
after they join, or, where available, relying on automatic license assignment
the first time they use a supported product. To select licenses when you invite
a member:

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Members** from the left navigation, then select **Invite**.
1. Select **Emails or usernames**.
1. Enter the email addresses or Docker IDs of the people you want to invite,
   then select the product licenses you want to assign.
1. Select **Send invitation**.

The selected licenses are assigned when the invitees accept their invitations,
provided those licenses are available at that time.

## Assign and revoke licenses

You can assign available licenses to members after they join your organization.
You can also revoke assigned licenses when a member no longer needs access.

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Members** from the left navigation.
1. Find the member whose licenses you want to manage.
1. Select the member's **Actions** menu.
1. Select **Manage licenses**.
1. Select or clear the licenses you want to assign or revoke.
1. Select **Save**.

When you revoke a license, the member remains in the organization unless you
remove them separately.

## Automatic license assignment

Automatic license assignment gives members a product license when they use a
supported product for the first time. Automatic license assignment is available
for AI Governance licenses. Only organizations that purchase AI Governance can
set up auto-assignment for Docker Core as well.

- When you purchase AI Governance, signing in to
  [Docker Sandboxes](/manuals/ai/sandboxes/_index.md) with the `sbx login`
  command automatically provisions AI Governance licenses on a first-come,
  first-served basis.
- Similarly, signing in to Docker Desktop automatically provisions Docker Core
  for AI Governance license-holding organizations that have available Docker
  Core seats.
- Licenses are assigned until exhausted.
  - Once the available licenses are exhausted, automatic license assignment
    stops until you purchase more licenses or revoke assigned licenses.
  - Members can still use Docker Sandboxes or Docker Desktop, but organization
    policies for those products won't affect their usage.

AI Governance licenses can be assigned independently of your Docker Core
subscription. However, SSO is a Docker Business feature. To configure
automatic license assignment, set up
[SSO](/manuals/security/authentication/single-sign-on/connect.md), then
configure [provisioning](/manuals/security/provisioning/_index.md) with System
for Cross-domain Identity Management (SCIM) or Just-in-Time (JIT).
