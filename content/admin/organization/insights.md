---
description: Gain insights about your organization's users and their Docker usage.
keywords: organization, insights
title: Insights
sitemap: false
---

> **Early Access**
>
> Insights is an [early access](/release-lifecycle#early-access-ea) feature and
> is only available to those in the early access feedback program.
{ .restricted }

> **Note**
>
> Insights requires a [Docker Business subscription](/subscription/core-subscription/details/#docker-business).

Insights helps administrators visualize and understand how Docker is used within
their organizations. With insights, administrators can ensure their teams are
fully equipped to utilize Docker to its fullest potential, leading to improved
productivity and efficiency across the organization.

Key benefits include:

 * Uniform working environment. Establish and maintain standardized
   configurations across teams.
 * Best practices. Promote and enforce usage guidelines to ensure optimal
   performance.
 * Increased visibility. Monitor and drive adoption of organizational
   configurations and policies.
 * Optimized license use. Ensure that developers have access to advanced
   features provided by a Docker subscription.

## View insights for organization users

1. Go to the [Admin Console](https://app.docker.com/admin/) and sign in to an
   account that is an organization owner.
2. In the Admin Console, select your organization from the drop-down in the left
   navigation.
3. Select **Insights**.
4. On the **Insights** page, select the period of time for the data.

> **Note**
>
> Insights data is not real-time. The data is at least one day behind.

You can view insights in the following charts:

 * [Docker Desktop users](#docker-desktop-users)
 * [Builds](#builds)
 * [Containers](#containers)
 * [Docker Desktop usage](#docker-desktop-usage)
 * [Images](#images)
 * [Extensions](#extensions)

### Docker Desktop users

Track active Docker Desktop users in your domain, differentiated by license
status. This chart helps you understand the engagement levels within your
organization, providing insights into how many users are actively using Docker
Desktop. Note that users who opt out of analytics aren't included in the active
counts.

The chart contains the following insights.

| Insight                      | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
|:-----------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Total active users           | The number of users that have opened Docker Desktop and are linked to a Docker account with an email address from a domain associated with your organization. <br><br>Users who don’t sign in to an account associated with your organization are not represented in the insights. To ensure users sign in with an account associated with your organization, you can [enforce sign-in](/security/for-admins/enforce-sign-in/).                                                                                                                                                                            |
| Active with license          | The number of users that have opened Docker Desktop, are linked to a Docker account with an email address from a domain associated with your organization, and have a license assigned to their account.                                                                                                                                                                                                                                                                                                                                                                                                   |
| Active without license       | The number of users that have opened Docker Desktop, are linked to a Docker account with an email address from a domain associated with your organization, and don’t have a license assigned to their account. <br><br>Users without a license don’t receive the benefits of your subscription. You can use [domain audit](/security/for-admins/domain-audit/) to identify users without a license. You can also use [Just-in-Time provisioning](/security/for-admins/provisioning/just-in-time/) or [SCIM](/security/for-admins/provisioning/scim/) to help automatically provision users with a license. |
| Users opted out of analytics | The number of users that have opened Docker Desktop, are linked to a Docker account with an email address from a domain associated with your organization, and have opted out of sending analytics. <br><br>When users opt out of sending analytics, you will see inaccurate insights. To ensure that your insights are accurate, you can use [Settings Management](/desktop/hardened-desktop/settings-management/) to set `analyticsEnabled` for all your users.                                                                                                                                          |
| Active users (graph)         | The number of users over time that have opened Docker Desktop and are linked to a Docker account with an email address from a domain associated with your organization.                                                                                                                                                                                                                                                                                                                                                                                                                                    |


### Builds

Monitor development efficiency and the time your team invests in builds with
this chart. It provides a clear view of the build activity, helping you identify
patterns, optimize build times, and enhance overall development productivity.

The chart contains the following insights.

| Insight                | Description                                                                                                                                                                                        |
|:-----------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Average build per user | The average number of builds per active user. A build includes any time a user runs a Docker command that builds an image.                                                                         |
| Average build time     | The average build time per build.                                                                                                                                                                  |
| Build success rate     | The percentage of builds that were successful out of the total number of builds. A successful build includes any build that exits normally.                                                        |
| Total builds (graph)   | The total number of builds separated into successful builds and failed builds. A successful build includes any build that exits normally. A failed build includes any build that exits abnormally. |

### Containers

View the total and average number of containers run by users with this chart. It
lets you gauge container usage across your organization, helping you understand
usage trends and manage resources effectively.

The chart contains the following insights.

| Insight                                | Description                                                                                                                  |
|:---------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------|
| Total containers run                   | The total number of containers run by active users. Containers run include those run using `docker run` or `docker compose`. |
| Average number of containers run       | The average number of containers run per active user.                                                                        |
| Containers run by active users (graph) | The number of containers run over time by active users.                                                                      |

### Docker Desktop usage

Explore Docker Desktop usage patterns with this chart to optimize your team's
workflows and ensure compatibility. It provides valuable insights into how
Docker Desktop is being utilized, enabling you to streamline processes and
improve efficiency.

The chart contains the following insights.

| Insight                           | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
|:----------------------------------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Most used version                 | The most used version of Docker Desktop by users.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| Most used OS                      | The most used operating system by users.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| Versions by active users (graph)  | The number of active users using each version of Docker Desktop. <br><br>To learn more about each version and release dates, see the [Docker Desktop release notes](/desktop/release-notes/). |
| Interface by active users (graph) | The number of active users grouped into the type of interface they used to interact with Docker Desktop. <br><br>A CLI user is any active user who has run a `docker` command. A GUI user is any active user who has interacted with the Docker Desktop graphical user interface.                                                                                                                                                                                                                                                                                                                                                                                   |

### Images

Analyze image distribution activity with this chart and view the most utilized
Docker Hub images within your domain. This information helps you manage image
usage, ensuring that the most critical resources are readily available and
efficiently used.

The chart contains the following insights.

| Insight              | Description                                                                                                     |
|:---------------------|:----------------------------------------------------------------------------------------------------------------|
| Total pulled images  | The total number of images pulled by users from Docker Hub.                                                     |
| Total pushed images  | The total number of images pushed by users to Docker Hub.                                                       |
| Top 10 pulled images | A list of the top 10 images pulled by users from Docker Hub and the number of times each image has been pulled. |

### Extensions

Monitor extension installation activity with this chart. It provides visibility
into the Docker Desktop extensions your team are using, letting you track
adoption and identify popular tools that enhance productivity.

The chart contains the following insights.

| Insight                                        | Description                                                                                                                          |
|:-----------------------------------------------|:-------------------------------------------------------------------------------------------------------------------------------------|
| Percentage of org with extensions installed    | The percentage of users in your organization with at least one Docker Desktop extension installed.                                   |
| Top 5 extensions installed in the organization | A list of the top 5 Docker Desktop extensions installed by users in your organization and the number of installs for each extension. |


## Troubleshoot insights

If you’re experiencing issues viewing accurate insights, consider the following
solutions to resolve common problems.

* Update users to the latest version of Docker Desktop.

    Insights are not shown for users using versions 4.16 or lower of Docker
    Desktop.  In addition, older versions may not provide all insights. Ensure
    all users have installed the latest version of Docker Desktop.

* Enable **Send usage statistics** in Docker Desktop for all your users.

   If users have opted out of sending usage statistics for Docker Desktop, then
   their usage data will not be a part of your insights. To manage the setting
   at scale for all your users, you can use [Settings
   Management](/desktop/hardened-desktop/settings-management/) and enable the
   `analyticsEnabled` setting.

* Ensure that users are using Docker Desktop and aren't using the standalone version
   of Docker Engine.

   Only Docker Desktop can provide usage statistics for insights. If a user
   installs and uses Docker Engine outside of Docker Desktop, some insights
   won’t appear for that user.

* Ensure that users are signing in to an account associated with your
  organization.

   Users who don’t sign in to an account associated with your organization are
   not represented in the insights. To ensure users sign in with an account
   associated with your organization, you can [enforce
   sign-in](/security/for-admins/enforce-sign-in/).