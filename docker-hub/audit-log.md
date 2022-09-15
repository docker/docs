---
description: Audit logs
keywords: Team, organization, activity, log, audit, activities
title: Audit logs
---

{% include upgrade-cta.html
  body="Audit logs are available for users subscribed to a Docker Team or a Business subscription. Upgrade now to start tracking events across your organization."
  header-text="This feature requires a paid Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade_audit_log"
%}

Audit logs display a chronological list of activities that occur at organization and repository levels. It provides owners of Docker Team accounts a report of all their team member activities. This allows the team owners to view and track what changes were made, the date when a change was made, and who initiated the change. For example, the audit logs display activities such as the date when a repository was created or deleted, the team member who created the repository, the name of the repository, and when there was a change to the privacy settings.

Team owners can also see the audit logs for their repository if the repository is part of the organization subscribed to a Docker Team plan.

## View the audit logs

To view the audit logs:

1. Sign into an owner account for the organization in Docker Hub.
2. Select your organization from the list and then click on the **Activity** tab.

    ![Organization activity tab](images/org-activity-tab.png){:width="700px"}

The audit logs begin tracking activities from the date the feature is live, that is from **25 January 2021**. Activities that took place before this date are not captured.

> **Note**
>
> Docker will retain the activity data for a period of three months.

## Customize the audit logs

By default, all activities that occur at organization and repository levels are displayed on the **Activity** tab. Use the calendar option to select a date range and customize your results. After you have selected a date range, the **Activity** tab displays the audit logs of all the activities that occurred during that period.

![Activities list](images/activity-list.png){:width="600px"}

> **Note**
>
> Activities created by the Docker Support team as part of resolving customer issues appear in the audit logs as **dockersupport**.

Click the **All Activities** drop-down list to view activities that are specific to an organization or a repository. After choosing **Organization** or **Repository**, you can further refine the results using the **All Actions** drop-down list. If you select the **Activities** tab from the **Repository** view, you can only filter repository-level activities.

![Refine org activities](images/org-all-actions.png){:width="600px"}


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