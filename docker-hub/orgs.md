---
description: Docker Hub Teams & Organizations
keywords: Docker, docker, registry, teams, organizations, plans, Dockerfile, Docker Hub, docs, documentation
title: Create your organization
redirect_from:
- /docker-cloud/orgs/
---

An organization in Docker Hub is a collection of teams and repositories
that can be managed together. A team is a group of Docker members that belong to an organization. 
An organization can have multiple teams. 

Docker users become members of an organization
when they are assigned to at least one team in the organization. When you first
create an organization, you have one team, the "owners" team, that has a single member. An organization owner is someone that is part of the
owners team. They can create new teams and add
members to an existing team using their Docker ID or email address and by
selecting a team the user should be part of. An organization owner can also add
additional owners to help them manage users, teams, and repositories in the
organization.

## Create an organization

There are multiple ways to create an organization. You can create a brand new
organization using the **Create Organization** option in Docker Hub, or you can
convert an existing user account to an organization. The following section
contains instructions on how to create new organization. For prerequisites and
detailed instructions on converting an existing user account to an org, see
[Convert an account into an organization](convert-account.md).

To create an organization:

1. Sign into [Docker Hub](https://hub.docker.com/){: target="_blank"
rel="noopener" class="_"} using your [Docker ID](../docker-id/index.md) or your email address.
2. Select **Organizations** and then **Create Organization** to create a new organization.
3. Choose a plan for your organization. See [Docker Pricing](https://www.docker.com/pricing/){: target="_blank" rel="noopener"
class="_" id="dkr_docs_subscription_btl"} for details on the features offered
in the Team and Business plan.
4. Enter a name for your organization. This is the official, unique name for
your organization in Docker Hub. It is not possible to change the name
of the organization after you've created it.

      > **Note**
      >
      > The organization name cannot be the same as your Docker ID.

5. Enter the name of your company. This is the full name of your company.
This info is displayed on your organization page, and in the details of any
public images you publish. You can update the company name anytime by navigating
to your organization's **Settings** page. 
6. Select **Continue to Org size** and then specify the number of users (seats) you'd
require.
7. Select **Continue to payment** and follow the onscreen instructions. 

You've now created an organization with one team, the owners team, with you as the single member.

## View an organization

To view an organization:

1. Sign in to Docker Hub with a user account that is a member of any team in the
   organization. 

      > **Note**
      >
      > You can't _directly_ log in to an organization. This is especially
      > important to note if you create an organization by
      [converting a user account](convert-account.md), as conversion means you lose the ability to log into that
      > "account", since it no longer exists. To view the organization you 
      > need to log in with the new owner account assigned during the
      > conversion or another account that was added as a member. If you 
      > don't see the organization after logging in,
      > then you are neither a member or an owner of it. An organization
      > administrator needs to add you as a member of the organization.

2. Select **Organizations** in the top navigation bar, then choose your
   organization from the list.

The organization landing page displays various options that allow you to
configure your organization.

- **Members**: Displays a list of team members. You
  can invite new members using the **Invite members** button. See [Manage members](../docker-hub/members.md) for details.

- **Teams**: Displays a list of existing teams and the number of
  members in each team. See [Create a team](manage-a-team.md) for details.

- **Repositories**: Displays a list of repositories associated with the
  organization. See [Repositories](../docker-hub/repos/index.md) for detailed information about
  working with repositories.

- **Activity** Displays the audit logs, a chronological list of activities that
  occur at organization and repository levels. It provides the org owners a
  report of all their team member activities. See [Audit logs](audit-log.md) for
  details.

- **Settings**: Displays information about your
  organization, and allows you to view and change your repository privacy
  settings, configure org permissions such as
  [Image Access Management](image-access-management.md), configure notification settings, and [deactivate](deactivate-account.md#deactivate-an-organization) your
  organization. You can also update your organization name and company name that appear on your organization landing page. You must be part of the owners team to access the
   organization's **Settings** page.

- **Billing**: Displays information about your existing
[Docker subscription (plan)](../subscription/index.md) and your billing history.
You can also access your invoices from this tab.

## Videos

You can also check out the following videos for information about creating Teams
and Organizations in Docker Hub.

- [Overview of organizations](https://www.youtube-nocookie.com/embed/G7lvSnAqed8){: target="_blank" rel="noopener" class="_"}
- [Create an organization](https://www.youtube-nocookie.com/embed/b0TKcIqa9Po){: target="_blank" rel="noopener" class="_"}
- [Working with Teams](https://www.youtube-nocookie.com/embed/MROKmtmWCVI){: target="_blank" rel="noopener" class="_"}
- [Create Teams](https://www.youtube-nocookie.com/embed/78wbbBoasIc){: target="_blank" rel="noopener" class="_"}

