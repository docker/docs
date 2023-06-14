---
description: Audit logs
keywords: Team, organization, activity, log, audit, activities
title: Audit logs
---

> **Note**
>
> Audit logs requires a [Docker Team, or Business subscription](../subscription/index.md).

Audit logs display a chronological list of activities that occur at organization and repository levels. It provides a report to owners of Docker Team on all their team member activities. 

With audit logs, team owners can view and track:
 - What changes were made
 - The date when a change was made
 - Who initiated the change
 
 For example, Audit logs display activities such as the date when a repository was created or deleted, the team member who created the repository, the name of the repository, and when there was a change to the privacy settings.

Team owners can also see the audit logs for their repository if the repository is part of the organization subscribed to a Docker Team plan.

Audit logs began tracking activities from the date the feature went live, that is from 25 January 2021. Activities that took place before this date are not captured.

## View the audit logs

To view the audit logs:

1. Sign in to Docker Hub.
2. Select your organization from the list and then select the **Activity** tab.

> **Note**
>
> Docker retains the activity data for a period of three months.

## Customize the audit logs

By default, all activities that occur at organization and repository levels are displayed on the **Activity** tab. Use the calendar option to select a date range and customize your results. After you have selected a date range, the **Activity** tab displays the audit logs of all the activities that occurred during that period.

![Activities list](images/activity-list.png){:width="600px"}

> **Note**
>
> Activities created by the Docker Support team as part of resolving customer issues appear in the audit logs as **dockersupport**.

Select the **All Activities** dropdown to view activities that are specific to an organization, repository, or billing. If you select the **Activities** tab from the **Repository** view, you can only filter repository-level activities.

After choosing **Organization**, **Repository**, or **Billing**, you can further refine the results using the **All Actions** dropdown. 

{% include admin-org-audit-log-events.md %}