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

## Audit logs event definitions

Refer to the following section for a list of events and their descriptions:

### Organization events

| Event                                                          | Description                                   |
|:------------------------------------------------------------------|:------------------------------------------------|
| Team Created | Activities related to the creation of a team |
| Team Updated | Activities related to the modification of a team |
| Team Deleted | Activities related to the deletion of a team |
| Team Member Added | Details of the member added to your team |
| Team Member Removed | Details of the member removed from your team |
| Team Member Invited | Details of the member invited to your team |
| Organization Member Added | Details of the member added to your organization |
| Organization Member Removed | Details about the member removed from your organization |
| Organization Created | Activities related to the creation of a new organization |
| Organization Settings Updated | Details related to the organization setting that was updated |
| Registry Access Management enabled | Activities related to enabling Registry Access Management |
| Registry Access Management disabled | Activities related to disabling Registry Access Management |
| Registry Access Management registry added | Activities related to the addition of a registry |
| Registry Access Management registry removed | Activities related to the removal of a registry |
| Registry Access Management registry updated | Details related to the registry that was updated |
| Single Sign-On domain added | Details of the single sign-on domain added to your organization |
| Single Sign-On domain removed | Details of the single sign-on domain removed from your organization |
| Single Sign-On domain verified | Details of the single sign-on domain verified for your organization |

### Repository events

| Event                                                          | Description                                   |
|:------------------------------------------------------------------|:------------------------------------------------|
| Repository Created | Activities related to the creation of a new repository |
| Repository Deleted | Activities related to the deletion of a repository |
| Privacy Changed | Details related to the privacy policies that were updated |
| Tag Pushed | Activities related to the tags pushed |
| Tag Deleted | Activities related to the tags deleted |

### Billing events

| Event                                                          | Description                                   |
|:------------------------------------------------------------------|:------------------------------------------------|
| Plan Upgraded | Occurs when your organization’s billing plan is upgraded to a higher tier plan.|
| Plan Downgraded | Occurs when your organization’s billing plan is downgraded to a lower tier plan. |
| Seat Added | Occurs when a seat is added to your organization’s billing plan. |
| Seat Removed | Occurs when a seat is removed from your organization’s billing plan. |
| Billing Cycle Changed | Occurs when there is a change in the recurring interval that your organization is charged.|
| Plan Downgrade Canceled | Occurs when a scheduled plan downgrade for your organization is canceled.|
| Seat Removal Canceled | Occurs when a scheduled seat removal for an organization’s billing plan is canceled. |