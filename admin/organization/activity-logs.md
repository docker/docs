---
description: Learn about activity logs.
keywords: team, organization, activity, log, audit, activities
title: Activity logs
---

{% include admin-early-access.md %}

Activity logs are a chronological list of activities that occur at organization and repository levels. The feature provides information to organization owners on all their team member activities.

With activity logs, owners can view and track:
 - What changes were made
 - The date when a change was made
 - Who initiated the change
 
 For example, activity logs display activities such as the date when a repository was created or deleted, the team member who created the repository, the name of the repository, and when there was a change to the privacy settings.

Owners can also see the activity logs for their repository if the repository is part of the organization.

## View the activity logs

To view the activity logs:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Activity Logs**.

> **Note**
>
> Docker retains the activity data for a period of three months.

## Customize the activity logs

By default, all activities that occur at organization and repository levels are displayed. Use the calendar option to select a date range and customize your results. After you have selected a date range, Docker Admin displays the activity logs of all the activities that occurred during that period.

> **Note**
>
> Activities created by the Docker Support team as part of resolving customer issues appear in the activity logs as **dockersupport**.

Select the **All Activities** dropdown to view activities that are specific to an organization, repository, or billing.

After choosing **Organization**, **Repository**, or **Billing**, you can further refine the results using the **All Actions** dropdown.

{% include admin-org-audit-log-events.md %}