---
description: Understand how to use the Desktop settings reporting dashboard
keywords: Settings Management, docker desktop, hardened desktop, reporting, compliance
title: Desktop settings reporting
linkTitle: Desktop settings reporting
weight: 30
params:
  sidebar:
    badge:
      color: violet
      text: EA
---

{{< summary-bar feature_name="Compliance reporting" >}}

Desktop settings reporting is a feature of Desktop Settings Management that
tracks and reports user compliance with the settings policies that are assigned
to them. This lets administrators track the application of settings and
monitor what actions they need to take to make users compliant.

This guide provides steps for accessing Desktop settings reporting, viewing
compliance status, and resolving non-compliant users.

## Access Desktop settings reporting

> [!IMPORTANT]
>
> Desktop settings reporting is in Early Access and is being rolled out
> gradually. You may not see this setting in the Admin Console yet.

1. Sign in to the [Admin Console](https://app.docker.com/admin).
2. Select your organization or company from the **Choose profile** page.
3. Under Docker Desktop, select **Reporting**.

This opens the Desktop settings reporting page. From here you can:

- Use the **Search** field to search by username or email address
- Filter by policies
- Hide or un-hide compliant users
- View a user’s compliance status and what policy is assigned to the user
- Download a CSV file of user compliance information

## View compliance status

> [!WARNING]
>
> Users on Docker Desktop versions older than 4.40 may appear non-compliant
> because older versions can't report compliance. To ensure accurate
> compliance status, users must update to Docker Desktop version 4.40 and later.

1. Sign in to the [Admin Console](https://app.docker.com/admin).
2. Select your organization or company from the **Choose profile** page.
3. Under **Docker Desktop**, select **Reporting**. By default, non-compliant users
are displayed.
4. Optional. Select the **Hide compliant users** checkbox to show both compliant
and non-compliant users.
5. Use the **Search** field to search by username or email address.
6. Hover over a user’s compliance status indicator to quickly view their status.
7. Select a username to view more details about their compliance status, and for
steps to resolve non-compliant users.

## Understand compliance status

Docker evaluates compliance status based on:

- Compliance status: Whether a user has fetched and applied the latest settings. This is the primary label shown on the reporting page.
- Domain status: Whether the user's email matches a verified domain.
- Settings status: Whether a settings policy is applied to the user.

The combination of these statuses determines what actions you need to take.

### Compliance status reference

This reference explains how each status is determined in the reporting dashboard
based on user domain and settings data. The Admin Console displays the
highest-priority applicable status according to the following rules.

#### Compliance status

| Compliance status | Description |
|-------------------|---------------|
| Uncontrolled domain | The user's email domain is not verified. |
| No policy assigned | The user does not have any policy assigned to them. |
| Non-compliant | The user fetched the correct policy, but hasn't applied it. |
| Outdated | The user fetched a previous version of the policy. |
| Unknown | The user hasn't fetched any policy yet, or their compliance can't be determined. |
| Compliant | The user fetched and applied the latest assigned policy. |

#### Domain status

| Domain status | Description |
|---------------|---------------|
| Verified | The user’s email domain is verified. |
| Guest user | The user's email domain is not verified. |
| Domainless | Your organization has no verified domains, and the user's domain is unknown. |
| Unknown user | Your organization has verified domains, but the user's domain is unknown. |

#### Settings status
| Settings status | Description |
|-----------------|---------------|
| Global policy | The user is assigned your organzation's default policy. |
| User policy | The user is assigned a specific custom policy. |
| No policy assigned | The user is not assigned to any policy. |


## Resolve compliance status

To resolve compliance status, you must view a user's compliance status details
by selecting their username from the Desktop settings reporting page.
These details include the following information:

- **Compliance status**: Indicates whether the user is compliant with the
settings applied to them
- **Domain status**: Indicates whether the user’s email address is associated
with a verified domain
- **Settings status**: Indicates whether the user has settings applied to them
- **Resolution steps**: If a user is non-compliant, this provides information
on how to resolve the user’s compliance status

### Compliant

When a user is compliant, a **Compliant** icon appears next to their name on the
Desktop settings reporting dashboard. Select a compliant user to open their
compliance status details. Compliant users have the following status details:

- **Compliance status**: Compliant
- **Domain status**: Verified
- **Settings status**: Global policy or user policy
- **User is compliant** indicator

No resolution steps are needed for compliant users.

### Non-compliant

When a user is non-compliant, a **Non-compliant** or **Unknown** icon appears
next to their name on the Desktop settings reporting dashboard. Non-compliant
users must have their compliance status resolved:

1. Select a username from the Desktop settings reporting dashboard.
2. On the compliance status details page, follow the resolution steps provided
to resolve the compliance status.
3. Refresh the page to ensure the resolution steps resolved the compliance
status.
