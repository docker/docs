---
description: Understand how to use the Desktop Settings Compliance reporting dashboard
keywords: Settings Management, docker desktop, hardened desktop, compliance reporting
title: Compliance reporting
linkTitle: Compliance Reporting
weight: 30
---

Compliance reporting is a feature of Desktop Settings Management that tracks and reports user compliance with the settings that are assigned to them. This lets administrators to track the application of settings and monitor what actions they need to take to make users compliant.

This guide provides steps for accessing the reporting dashboard, viewing compliance status, and resolving non-compliant users.

## Access compliance reporting dashboard

1. Sign in to the [Admin Console](https://app.docker.com/admin).
2. Select your **organization** from the **Choose profile** page.
3. Under Security and access, select **Desktop Settings reporting**.

This opens the Desktop Settings Management Compliance dashboard. From here you can:

- Use the **Search bar** to search by username or email address
- Filter by policies
- Hide or un-hide compliant users
- View a user’s compliance status and what policy is assigned to the user
- Download a CSV file of user compliance information

## View compliance status

1. Sign in to the [Admin Console](https://app.docker.com/admin).
2. Select your **organization** from the **Choose profile** page.
3. Under Security and access, select **Desktop Settings reporting**.
4. By default, non-compliant users are displayed.
5. Optional. Select the **Hide compliant users** checkbox to show compliant and non-compliant users.
6. Use the search bar to search by username or email address.
7. Hover over a user’s compliance status indicator to quickly view their status.
8. Select a **username** to view more details about their compliance status, and for steps to resolve non-compliant users.

## Resolve compliance status

To resolve compliance status, you must open a user's compliance status details. This window provides the following information:

- **Compliance status**: Indicates whether the user is compliant with the settings applied to them
- **Domain status**: Indicates whether the user’s email address is associated with a verified domain
- **Settings status**: Indicates whether the user has settings applied to them
- **Resolution steps**: If a user is non-compliant, this provides information on how to resolve the user’s compliance status

### Compliant

When a user is compliant, a Compliant icon appears next to their name on the Desktop Settings Management Compliance ****dashboard. Select a compliant user to open their compliance status details, they will have the following status details:

- **Compliance status**: Compliant
- **Domain status**: Verified domain
- **Settings status**: Compliant
- **User is compliant** indicator

No resolution steps are needed for compliant users.

### Non-compliant

When a user is non-compliant, a Non-compliant or Unknown icon appears next to their name on the Desktop Settings Management Compliance dashboard. Non-compliant and Unknown statuses need resolved to become compliant.

