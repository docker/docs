---
title: Manage licenses
linkTitle: Licenses
description: View your organization's license inventory and assign licenses to teams or individual members, including invite-time, automatic, and queued assignment.
keywords: licenses, organization, teams, members, invite, Docker Core, Docker Offload, AI Governance, license assignment, team assignment, license queue, docker home
weight: 40
aliases:
  - /admin/organization/manage/manage-licenses/
---

Licenses control which organization members can use supported Docker products.
As an organization owner, you manage license availability for your
organization.

> [!TIP]
> To learn more about product licenses, Docker Team and Business seats, and other Docker
> add-ons, see [Docker plans](/manuals/subscription-billing/plans/_index.md),
> or
> <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_admin_licenses" class="link" rel="noopener">contact sales</a>
> to purchase licenses.

## License assignment

You have a few options for assigning a license to a member. You can assign:

- Through a team, so every member of that team gets the license, including
  people who join the team later
- Through the **Members** page with the action menu, or through invitations
- By turning on automatic assignment so members receive a license the first time they
  use a supported product.

A member can use the product if they have a license from their team, from an
individual assignment, or from automatic assignment. Each member uses one
license per product. Assigning the same product again through another team or
as an individual assignment does not consume a second license.

## Assign licenses

Assign licenses from the **Teams** view, the **Members** view, through
invitations, or with automatic assignment.

### Members

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Members** from the left navigation.
1. Select the **action menu** at the end of the member's row to assign or
   revoke an active license.
1. Optional. To assign or revoke licenses for several members, use multi-select to choose the
   members you want to manage, then select the **Bulk actions** menu.

### Teams

Assigning a license to a team ensures every member of that team receives the license, including
members who join the team later.

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Teams** from the left navigation, then select the team name.
1. On the **Licenses** card, select the **edit** icon to open **Add licenses**.
1. Under **Licenses**, select one or more licenses.
    - Each license shows how many are available
    - The modal reports how many members receive each license
1. Select **Save**.

Docker grants the license only to team members who don't already have it.

- If a member already holds that license, they keep access and the extra
  assignment does not consume another license.
- If there aren't enough licenses for every remaining member, selecting
  **Save** assigns the licenses you have. Docker queues the rest until more
  are available.

### Invitations

Assignment happens on acceptance if a license is available:

- If a license is available when they accept, Docker assigns it to them and the
  number of available licenses decreases by one.
- If no licenses remain when they accept, they still join your organization,
  but without a license.
- Docker doesn't reserve or deduct the license at invite time.

You can monitor availability on the **Licenses**
page while invitations are pending.

To select licenses through invitations:

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Members** from the left navigation, then select **Invite**.
1. Select **Emails or usernames**.
1. Enter the email addresses or Docker IDs of the people you want to invite,
   then assign their
   [role](/manuals/security/roles-and-permissions/_index.md).
1. Under **Licenses (optional)**, select one or more licenses that are
   available to your organization.
1. Select **Invite** to send the invite.

A user can accept from the link in their invitation email or from their
**Notifications Center**. For more about sending, resending, and removing invitations, including CSV
file limits, see
[Manage organization members](/manuals/accounts/organization/manage/members.md).

### Automatic assignment

Automatic license assignment gives members a product license when they use a
supported product for the first time. To turn it on, use the
**Automatic license assignment** toggle on the product's license card on the
**Licenses** page.

- Members receive a Docker Core license the first time they sign in to Docker
  Desktop.
- Signing in to
  [Docker Sandboxes](/manuals/ai/sandboxes/_index.md) with the `sbx login`
  command provisions AI Governance licenses on a first-come, first-served
  basis.
- Licenses are assigned until exhausted.
  - Once the available licenses are exhausted, automatic license assignment
    stops until more licenses are available.
  - Members can still use Docker Sandboxes or Docker Desktop, but organization
    policies for those products won't affect their usage.

AI Governance licenses include single sign-on (SSO) and provisioning features
regardless of your Docker Core subscription. Automatic license assignment for
AI Governance requires
[setting up SSO](/manuals/security/authentication/single-sign-on/connect.md), then
[provisioning](/manuals/security/provisioning/_index.md) with System
for Cross-domain Identity Management (SCIM) or Just-in-Time (JIT).

## View licenses

The **Licenses** page shows how many licenses you have, how many are assigned,
and whether those assignments are to teams or to individual members. Use it to
check remaining capacity, open the members or teams that have a license, add
licenses to your subscription, and turn automatic assignment on or off.

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Licenses** from the left navigation.

Products you haven't purchased appear as cards with **Learn more** and **Add
licenses**. Products you own show how many licenses are:

- Remaining versus total assigned under **Available**. Select **View all** to
  open the members who have the license.
- Assigned through teams under **Team assignment**. Select
  **View teams with this license** to open those teams.
- Assigned to members individually under **Direct assignment**

## Remove licenses

When you remove a license from a team or a member, it becomes available to members queued from a team assignment. The same rule applies when a member leaves a team or you
[delete a team](/manuals/accounts/organization/manage/manage-a-team.md#delete-a-team).

To remove licenses:

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. To remove licenses from a team:
    - Select **Teams**, then select the team name.
    - On the **Licenses** card, select the **edit** icon.
    - Remove the license you no longer want the team to assign.
1. To revoke a license from one member:
    - Select **Members**.
    - Use the **action menu** at the end of the member's row.
    - Select **Remove**.
1. Review the confirmation message, then select **Remove license**.

## Next steps

Explore Docker Core add-ons and products that need licenses:

- [Docker plans](/manuals/subscription-billing/plans/_index.md) to learn about different
  add-ons
- [Manage seats](/manuals/accounts/organization/manage/manage-seats.md) to add more
  seats to your Docker Core subscription
- [Create and manage a team](/manuals/accounts/organization/manage/manage-a-team.md)
  to group members and assign licenses to a team
- [AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md) to learn
  about AI Governance license usage and billing
- [Docker Offload](/manuals/offload/about.md) to let your developers offload
  building and running containers to the cloud
