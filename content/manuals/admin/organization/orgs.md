---
title: Create your organization
weight: 10
description: Learn how to create an organization.
keywords: Docker, docker, registry, teams, organizations, plans, Dockerfile, Docker
  Hub, docs, documentation
aliases:
- /docker-hub/orgs/
---

{{< summary-bar feature_name="Admin orgs" >}}

This section describes how to create an organization. Before you begin:

- You need a [Docker ID](/accounts/create-account/)
- Review the [Docker subscriptions and features](../../subscription/details.md) to determine what plan to choose for your organization

## Create an organization

{{< summary-bar feature_name="Admin console early access" >}}

There are multiple ways to create an organization. You can either:
- Create a new organization using the **Create Organization** option in Docker Hub
- Convert an existing user account to an organization

The following section contains instructions on how to create a new organization. For prerequisites and
detailed instructions on converting an existing user account to an organization, see
[Convert an account into an organization](/manuals/admin/organization/convert-account.md).

{{< tabs >}}
{{< tab name="Docker Hub" >}}

1. Sign in to [Docker Hub](https://hub.docker.com/) using your Docker ID, your email address, or your social provider.
2. Select **Organizations** and then **Create Organization** to create a new organization.
3. Choose a plan for your organization, a billing cycle, and specify how many seats you need. See [Docker Pricing](https://www.docker.com/pricing/) for details on the features offered in the Team and Business plan.
4. Select **Continue to profile**.
5. Enter an **Organization namespace**. This is the official, unique name for
your organization in Docker Hub. It's not possible to change the name of the
organization after you've created it.

   > [!NOTE]
   >
   > You can't use the same name for the organization and your Docker ID. If you want to use your Docker ID as the organization name, then you must first [convert your account into an organization](/manuals/admin/organization/convert-account.md).

6. Enter your **Company name**. This is the full name of your company. Docker
displays the company name on your organization page and in the details of any
public images you publish. You can update the company name anytime by navigating
to your organization's **Settings** page.
7. Select **Continue to billing** to continue.
8. Enter your organization's billing information and select **Continue to payment** to continue to the billing portal.
9. Provide your card details and select **Purchase**.

You've now created an organization.

{{< /tab >}}
{{< tab name="Admin Console" >}}

To create an organization:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Go to Admin Console**.
3. Select the **Organization** drop-down in the left-hand navigation and then **Create Organization**.
4. Choose a plan for your organization, a billing cycle, and specify how many seats you need. See [Docker Pricing](https://www.docker.com/pricing/) for details on the features offered in the Team and Business plan.
5. Select **Continue to profile**.
6. Enter an **Organization namespace**. This is the official, unique name for
your organization in Docker Hub. It's not possible to change the name of the
organization after you've created it.

   > [!NOTE]
   >
   > You can't use the same name for the organization and your Docker ID. If you want to use your Docker ID as the organization name, then you must first [convert your account into an organization](/manuals/admin/organization/convert-account.md).

7. Enter your **Company name**. This is the full name of your company. Docker
displays the company name on your organization page and in the details of any
public images you publish. You can update the company name anytime by navigating
to your organization's **Settings** page.
8. Select **Continue to billing** to continue.
9. Enter your organization's billing information and select **Continue to payment** to continue to the billing portal.
10. Provide your card details and select **Purchase**.

You've now created an organization.

{{< /tab >}}
{{< /tabs >}}

## View an organization

{{< summary-bar feature_name="Admin console early access" >}}

{{< tabs >}}
{{< tab name="Docker Hub" >}}

To view an organization:

1. Sign in to [Docker Hub](https://hub.docker.com) with a user account that is a member of any team in the
   organization.

      > [!NOTE]
      >
      > You can't *directly* sign in to an organization. This is especially
      > important to note if you create an organization by
      [converting a user account](/manuals/admin/organization/convert-account.md), as conversion means you lose the ability to log into that
      > "account", since it no longer exists. To view the organization you
      > need to sign in with the new owner account assigned during the
      > conversion or another account that was added as a member. If you
      > don't see the organization after logging in,
      > then you are neither a member or an owner of it. An organization
      > administrator needs to add you as a member of the organization.

2. Select **Organizations** in the top navigation bar, then choose your
   organization from the list.

The organization landing page displays various options that let you to
configure your organization.

- **Members**: Displays a list of team members. You
  can invite new members using the **Invite members** button. See [Manage members](./members.md) for details.

- **Teams**: Displays a list of existing teams and the number of
  members in each team. See [Create a team](./manage-a-team.md) for details.

- **Repositories**: Displays a list of repositories associated with the
  organization. See [Repositories](../../docker-hub/repos/_index.md) for detailed information about
  working with repositories.

- **Activity** Displays the audit logs, a chronological list of activities that
  occur at organization and repository levels. It provides the org owners a
  report of all their team member activities. See [Audit logs](./activity-logs.md) for
  details.

- **Settings**: Displays information about your
  organization, and you to view and change your repository privacy
  settings, configure org permissions such as
  [Image Access Management](/manuals/security/for-admins/hardened-desktop/image-access-management.md), configure notification settings, and [deactivate](../deactivate-account.md#deactivate-an-organization) You can also update your organization name and company name that appear on your organization landing page. You must be an owner to access the
   organization's **Settings** page.

- **Billing**: Displays information about your existing
[Docker subscription (plan)](../../subscription/_index.md), including the number of seats and next payment due date. For how to access the billing history and payment methods for your organization, see [View billing history](../../billing/history.md).

{{< /tab >}}
{{< tab name="Admin Console" >}}

To view an organization in the Admin Console:

1. Sign in to [Docker Home](https://app.docker.com).
2. Under Settings and administration, select **Go to Admin Console**.
3. Select your organization from the **Organization** drop-down in the left-hand navigation.

The Admin Console displays various options that let you to
configure your organization.

- **Members**: Displays a list of team members. You
  can invite new members using the **Invite members** button. See [Manage members](./members.md) for details.

- **Teams**: Displays a list of existing teams and the number of
  members in each team. See [Create a team](./manage-a-team.md) for details.

- **Activity** Displays the audit logs, a chronological list of activities that
  occur at organization and repository levels. It provides the org owners a
  report of all their team member activities. See [Audit logs](./activity-logs.md) for
  details.

- **Security and access**: Manage security settings. For more information, see [Security](/manuals/security/_index.md).

- **Organization settings**: Update general settings, manage your company settings, or [deactivate your organization](/manuals/admin/deactivate-account.md).

{{< /tab >}}
{{< /tabs >}}

## Merge organizations

> [!WARNING]
>
> If you are merging organizations, it is recommended to do so at the *end* of
> your billing cycle. When you merge an organization and downgrade another, you
> will lose seats on your downgraded organization. Docker does not offer
> refunds for downgrades.

If you have multiple organizations that you want to merge into one, complete the following:

1. Based on the number of seats from the secondary organization, [purchase additional seats](../../subscription/manage-seats.md) for the primary organization account that you want to keep.
2. Manually add users to the primary organization and remove existing users from the secondary organization.
3. Manually move over your data, including all repositories.
4. Once you're done moving all of your users and data, [downgrade](../../subscription/change.md) the secondary account to a free subscription. Note that Docker does not offer refunds for downgrading organizations mid-billing cycle.

> [!TIP]
>
> If your organization has a Docker Business subscription with a purchase order, contact Support or your Account Manager at Docker.

## More resources

- [Video: Docker Hub Organizations](https://www.youtube.com/watch?v=WKlT1O-4Du8)
