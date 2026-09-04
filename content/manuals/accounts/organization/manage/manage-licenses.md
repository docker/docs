---
title: Manage licenses
linkTitle: Licenses
description: View your organization's license inventory and assign licenses to teams or individual members, including invite-time, automatic, and queued assignment.
keywords: licenses, organization, teams, members, invite, Docker Core, Docker Offload, AI Governance, license assignment, team assignment, license queue, docker home
weight: 30
aliases:
  - /admin/organization/manage/manage-licenses/
---

Licenses control which organization members can access supported Docker
products. View your license inventory on the Licenses page, assign a license to
a whole team from the Teams view, or assign it to one person from the Members
view. Automatic license assignment can also grant a license the first time a
member uses a supported product.

> [!TIP]
> To learn more about product licenses, Docker Core seats, and other Docker
> add-ons, see [Docker plans](/manuals/subscription-billing/plans/_index.md),
> or
> <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_admin_licenses" class="link" rel="noopener">contact sales</a>
> to purchase licenses.

## How licenses get assigned

A member can receive a license from three sources:

- A team assignment, which grants the license to every member of that team,
  including members who join the team later
- A direct assignment to an individual member, including a license you select
  when you invite them
- Automatic license assignment, which grants a license the first time the
  member uses a supported product

A member can access the product if any source grants the license. Each member
holds at most one license per product, so assigning a license through a second
source doesn't consume another license. Members who already hold a license keep
the license they have.

## View licenses

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Licenses** from the left navigation.

Products your organization hasn't purchased appear as cards with a description,
a link to learn more, and **Add licenses**. Products you own appear as license
cards that show your plan and the following counts:

- **Available** reports how many licenses remain, together with how many of
  your total licenses are assigned and a link to view every member that holds
  the license
- **Team assignment** reports how many licenses come from teams, with a link to
  view the teams that hold the license
- **Direct assignment** reports how many licenses are assigned to members
  individually

The member and team lists show who holds a license. They don't identify which
source granted it.

Each license card also has the
[automatic license assignment](#automatic-license-assignment) toggle and
**Manage subscription**.

When no licenses remain, **Available** reports 0 and a banner reports how many
members in teams that assign the license are waiting for one. To buy more
licenses, select **Add licenses** or **Manage subscription**. These actions
change your subscription. To assign licenses you already own, use the Teams or
Members view.

## Assign licenses to a team

When you assign a license to a team, every member of that team receives it, and
members who join the team later receive it too.

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Teams**, then select the team name.
1. On the **Licenses** card, select the edit icon to open **Add licenses**.
1. Select one or more licenses. Each license shows how many are available, and
   the modal reports how many members receive each license you add.
1. Follow the on-screen instructions to save.

Members who already hold one of the selected licenses keep the license they
have. If the team has more members than you have available licenses, you can
still save. See [When licenses run out](#when-licenses-run-out).

## Assign licenses to a member

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Members** from the left navigation.
1. Select the action menu at the end of the member's row to assign or revoke an
   active license.
1. Optional. To assign or revoke licenses for several members, select the
   members you want to manage, then select the **Bulk actions** menu.

To assign a license when you invite someone, see
[Licenses and invites](#licenses-and-invites).

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
available when each invitee accepts. Monitor availability on the Licenses page
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
   then assign their
   [role](/manuals/security/roles-and-permissions/_index.md).
1. Under **Licenses (optional)**, select one or more licenses that are
   available to your organization.
1. Select **Invite** to send the invite.

For more about sending, resending, and removing invitations, including CSV
file limits, see
[Manage organization members](/manuals/accounts/organization/manage/members.md).

### Accept invites

A user can accept from the link in their invitation email or from their
**Notifications Center**. If the selected license is available, Docker assigns
it automatically upon acceptance.

## Automatic license assignment

Automatic license assignment gives members a product license when they use a
supported product for the first time. To turn it on, use the
**Automatic license assignment** toggle on the product's license card on the
Licenses page.

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

## When licenses run out

Docker assigns licenses in the order they're requested. If a request needs more
licenses than you have available, Docker assigns the available licenses and
queues the remaining members.

- When you assign a license to a team with more members than you have available
  licenses, you can still save. Docker warns you how many members are queued.
- The Licenses page reports how many members in teams that assign the license
  are waiting for one.
- When you add licenses, Docker assigns them to the queued members and resumes
  automatic license assignment.
- Licenses you revoke, and licenses that return to your pool when you remove a
  team assignment, become available again and automatic license assignment can
  claim them.

## Remove or lose a license

When you remove a license from a team, the team stops granting it. Members lose
access immediately unless another team still assigns the license or they hold a
direct assignment. The removed licenses return to your available pool.

1. Sign in to [Docker Home](https://app.docker.com), then choose your
   organization.
1. Select **Teams**, then select the team name.
1. On the **Licenses** card, select the edit icon.
1. Remove the license you no longer want the team to assign.
1. Review the confirmation message, which reports how many members lose access
   and how many keep the license from another source, then select
   **Remove license**.

The same rule applies when a member leaves a team and when you
[delete a team](/manuals/accounts/organization/manage/manage-a-team.md#delete-a-team):
members lose the licenses that team granted unless another team assigns the
same license or they hold a direct assignment.

To revoke a license from one member, use the action menu on the Members page.
See [Assign licenses to a member](#assign-licenses-to-a-member).

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
