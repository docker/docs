---
title: View organization Insights for Docker Desktop
linkTitle: Insights
weight: 50
description: >-
  View Docker Desktop usage for your organization in Insights, including
  active users, builds, containers, images, and extensions.
keywords: organization, Insights, Desktop insights, Docker Desktop analytics,
  Send usage statistics, analyticsEnabled, active users, builds, containers,
  images, extensions, export Docker Desktop users, Docker Business
aliases:
  - /admin/insights/
---

{{< summary-bar feature_name="Insights" >}}

Insights shows organization administrators how people in their organization
use Docker Desktop. Use it to track adoption, spot version and platform
patterns, and see where builds, containers, images, and extensions are in
use.

> [!TIP]
>
> After
> [upgrading to Docker Business](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdminInsights),
> you must reach out to your Docker account team to enable Insights for
> your organization.

## Prerequisites

To use Insights, you must meet the following requirements:

- A Docker Business subscription
- Administrators must
  [enforce sign-in](/manuals/desktop/enterprise/enforce-sign-in/_index.md)
  for users
- Docker Desktop users must have
  [Send usage
  statistics](/manuals/desktop/enterprise/hardened-desktop/settings-management/settings-reference.md#send-usage-statistics)
  enabled (`analyticsEnabled`). You can set this for every user through
  [Settings
  Management](/manuals/desktop/enterprise/hardened-desktop/settings-management/_index.md).

## View Insights {#view-insights-for-organization-users}

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Docker Desktop**, then **Desktop insights**.
1. Select a time period: **1 Week**, **1 Month**, or **3 Months**. Data
   may be delayed 24–48 hours. **Last updated** shows the sync date.

## Export Docker Desktop user data

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Docker Desktop**, then **Desktop insights**.
1. Select **Export Docker Desktop users**. The CSV file downloads
   automatically. It includes the standard
   [member CSV fields](/manuals/accounts/organization/manage/members.md#export-a-member-list-csv),
   plus these Insights columns:

   | Column                 | Description                                                                                                                                                                                 |
   | :--------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
   | Last Logged-in Date    | Last browser sign-in to Docker Hub or Docker Home                                                                                                                                           |
   | Docker Desktop Version | Last Docker Desktop version the user ran                                                                                                                                                    |
   | Last Seen Date         | Last time the user ran Docker Desktop                                                                                                                                                       |
   | Opted Out Analytics    | Whether the user turned off [Send usage statistics](/manuals/desktop/enterprise/hardened-desktop/settings-management/settings-reference.md#send-usage-statistics) |

## Reference

The Insights dashboard is organized into sections that capture each area of
Docker Desktop activity. Refer to these tables when a value looks
unexpected, when you need to explain a metric to someone else, or when
you're deciding which section answers a question about adoption, build
health, or image usage.

{{< accordion title="Docker Desktop users" >}}

The Docker Desktop users section shows how many people in your organization
run Docker Desktop and how that number changes over time. Use it to measure
adoption and to find users missing from the active counts. Active counts
include only users who send analytics.

| Data | Description |
| :--------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Active members | Users who used Docker Desktop and signed in with a licensed account in your organization, or with an email address from a domain associated with your organization. |
| Total org members | Users who have used Docker Desktop, regardless of their Insights activity. |
| Members opted out of analytics | Organization members who opted out of sending analytics. |
| Registered users | Users linked to unassociated machines whose email addresses Docker captured, and those emails are outside your organization's domain. |
| Active users (graph) | Total active users over time. |

{{< /accordion >}}

{{< accordion title="Builds" >}}

The Builds section shows how often users build images, how long builds take,
and how many succeed. Use it to spot slow or failing builds across the
organization and to gauge how much of your Docker Desktop usage is build
activity.

| Data | Description |
| :--------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Average build per user | Average builds per active user. A build includes any time a user runs one of: `docker buildx bake`, `docker buildx build`, `docker buildx f`, `docker builder build`, `docker builder f`, `docker compose build`, or `docker image build`. |
| Average build time | Average build time per build. |
| Build success rate | Percentage of builds that exited normally. A failed build is any build that exits abnormally. |
| Total builds (graph) | Successful builds versus failed builds over time. |

{{< /accordion >}}

{{< accordion title="Containers" >}}

The Containers section shows how many containers active users run and how
that activity trends over time. Use it to see how heavily your organization
relies on Docker Desktop for running workloads, not only for building them.

| Data | Description |
| :------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Total containers run | Containers run by active users in the Docker Desktop GUI, with `docker run`, or with `docker compose`. |
| Average number of containers run | Average containers run per active user. |
| Containers run by active users (graph) | Containers run over time by active users. |

{{< /accordion >}}

{{< accordion title="Docker Desktop usage" >}}

The Docker Desktop usage section shows which versions, operating systems,
and interfaces your users run. Use it to plan upgrades, find users on
outdated versions, and understand whether your organization works mainly
through the CLI or the GUI.

| Data | Description |
| :-------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Most used version | Most used Docker Desktop version in your organization. |
| Most used OS | Most used operating system. |
| Versions by active users (graph) | Active users per Docker Desktop version. See the [Docker Desktop release notes](/manuals/desktop/release-notes.md). |
| Interface by active users (graph) | Active users by interface. A CLI user has run a `docker` command. A GUI user has used the Docker Desktop graphical interface. |

{{< /accordion >}}

{{< accordion title="Docker Hub images" >}}

The Docker Hub images section shows which images your users pull and push,
and how often. Use it to see which images matter most to your organization
and to confirm that users pull from the sources you expect. Image data
covers Docker Hub only. Third-party registries and mirrors are excluded.

| Data | Description |
| :----------------- | :------------------------------------------------------------------------------------------------------ |
| Total pulled images | Images pulled by active users. |
| Total pushed images | Images pushed. |
| Top pulled images | Top images pulled, and how many times each was pulled. |

{{< /accordion >}}

{{< accordion title="Extensions" >}}

The Extensions section shows how widely Docker Desktop extensions are used
in your organization and which ones users install most. Use it to find
popular tools worth supporting and to spot extensions you did not
expect.

| Data | Description |
| :---------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------------- |
| Percentage of org with extensions installed | Percentage of users with at least one extension installed. |
| Top extensions installed in the organization | Top extensions installed, and how many users installed each one. |

{{< /accordion >}}

## Next steps

- [Activity logs](/manuals/accounts/organization/activity-logs.md) for a
  record of organization and repository changes.
- [Insights FAQs](/manuals/faqs/accounts.md#why-is-insights-data-missing-for-some-users)
  to troubleshoot missing data.
