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

- Compliance status: Whether a user has fetched and applied the latest settings.
- Domain status: Whether the user's email matches a verified domain.
- Settings status: Whether a settings policy is applied to the user.

The combination of these statuses determines what actions an administrator needs to take.

### Compliance status reference

Use the following table to understand how a user’s compliance status is
determined based on their domain status and settings status. Each row represents
a combination of statuses that may appear in the reporting dashboard.

> [!TIP]
>
> If a combination is marked "Not possible" in the reference table,
it means that Docker does not currently report that combination due to how
compliance is evaluated. These are included here to help you interpret
unexpected or unclear data.

| Domain status   | Settings status     | Compliance: Unknown                                                                 | Compliance: Outdated                                                                                      | Compliance: Non-compliant                                                                                      | Compliance: Compliant                                            |
|-----------------|---------------------|--------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------|
| Unknown         | Unknown             | The user does not exist in the system, is inactive, or has never been a member of the organization. | The user was removed from the organization after fetching settings that are now outdated.                 | The user was removed from the organization after fetching valid settings.                                       | Not possible                                                     |
| Uncontrolled    | Unknown             | The user’s email domain is not verified and the user has never fetched settings. May be a guest user or inactive. | The user was removed or changed their email after fetching outdated settings.                             | The user was removed or changed their email after fetching valid settings.                                      | Not possible                                                     |
| Controlled      | Uncontrolled        | The user does not have a settings package assigned or applicable to them.            | The user is assigned to a settings package but has not fetched the latest version.                         | Not possible                                                                                                     | Not possible                                                     |
| Controlled      | Controlled          | The user has not logged into Docker Desktop to fetch settings.                       | The user fetched an outdated version of the settings.                                                      | The user fetched the latest settings but has not applied them yet.                                              | The user has fetched and applied the latest settings.            |
| Uncontrolled  | Controlled       | The user is associated with a verified settings policy, but their email domain is not verified. May be a guest user. | The user fetched outdated settings but their domain is still not verified. | The user fetched valid settings but their domain is not verified. | Not possible |

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
- **Domain status**: Verified domain
- **Settings status**: Compliant
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
