---
title: Manage unassociated machines
description: Learn how to manage unassociated machines using the Docker Admin Console
keywords: unassociated machines, insights, manage users, enforce sign-in
weight: 56
---

{{< summary-bar feature_name="Unassociated machines" >}}

Docker administrators can identify, view, and manage Docker Desktop machines
that should be associated with their organization but aren't currently linked
to user accounts. This self-service capability helps you understand Docker
Desktop usage across your organization and streamline user onboarding without
IT involvement.

## Prerequisites

- Docker Business subscription
- Organization owner access to your Docker organization

## About unassociated machines

Docker Desktop machines in your organization may be:

- Associated: The user has signed in to Docker Desktop and is a member of
your organization
- Unassociated: Docker has identified machines likely belonging to your
organization based on usage patterns, but the users haven't signed in or
joined your organization

## How Docker identifies unassociated machines

Docker uses telemetry data to identify which machines belong to your
organization:

- Private registry usage: Machines accessing your organization's private
container registries
- Domain matching: Users signed in with email domains associated with your
organization
- Registry patterns: Analysis of container registry access patterns that
indicate organizational usage

## View unassociated machines

To see detailed information about unassociated machines:

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select
your organization.
1. In **User management**, select **Unassociated**.

The machine list displays:

- Machine ID (Docker-generated identifier)
- Registry addresses accessed (when available)
- User email
- Docker Desktop version
- Last activity date
- Sign-in enforced status

You can:

- Export the list as CSV
- Take actions on individual or multiple machines

## Add unassociated machines to your organization

You can add unassociated machines by:
- [Auto-provisiong](/manuals/security/for-admins/domain-management.md#auto-provisioning)
- [SSO user provisioning](/manuals/security/for-admins/provisioning/_index.md)
- [Manually adding them](#add-unassociated-machines-to-your-organization)

> [!NOTE]
>
> If you add users and do not have enough seats in your organization, a
pop-up will appear prompting you to **Get more seats**.

### Add individual users

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select
your organization.
1. In **User management**, select **Unassociated**.
1. Locate the machine you want to add to your organization.
1. Select the **Actions** menu and choose **Add to organization**.
1. In the pop-up modal, select **Add user**.

### Bulk add users

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select
your organization.
1. In **User management**, select **Unassociated**.
1. Use the **checkboxes** to select the machines you want to add to your
organizations.
1. Select the **Add to organization** button.
1. In the pop-up modal, select **Add users** to confirm.

## Enable sign-in enforcement

### Enable for all unassociated machines

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select
your organization.
1. In **User management**, select **Unassociated**.
1. Turn on the **Require sign-in** toggle.
1. In the pop-up modal, select **Require sign-in** to confirm.

The **Sign-in required** status will update for all unassociated machines to
**Yes**.

> [!NOTE]
>
> Sign-in enforcement requires Docker Desktop version 4.37 or later. If you
enable enforcement for a user with an older version, their status shows
as **Pending** until they update Docker Desktop.

### Enable for individual unassociated machines

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select
your organization.
1. In **User management**, select **Unassociated**.
1. Locate the machine you want to enable sign-in enforcement for.
1. Select the **Actions** menu and choose **Turn on sign-in enforcement**.
1. In the pop-up modal, select **Require sign-in** to confirm.

The **Sign-in required** status will update for the individual machine to
**Yes**.

> [!NOTE]
>
> Sign-in enforcement works with Docker Desktop versions 4.37 and later. If you
enable sign-in enforcement for a user using an older version of Docker Desktop,
their **Sign-in required** status will display as **Pending**.

## Disable sign-in enforcement

### Disable for all unassociated machines

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select
your organization.
1. In **User management**, select **Unassociated**.
1. Turn off the **Require sign-in** toggle.
1. In the pop-up modal, select **Turn off sign-in** to confirm.

The **Sign-in required** status will update for all unassociated machines to
**No**.

### Disable for specific unassociated machines

1. Sign in to the [Admin Console](https://app.docker.com/admin) and select
your organization.
1. In **User management**, select **Unassociated**.
1. Locate the machine you want to disable sign-in enforcement for.
1. Select the **Actions** menu and choose **Turn off sign-in enforcement**.
1. In the pop-up modal, select **Turn off sign-in** to confirm.

The **Sign-in required** status will update for the individual machine to
**No**.

## Developer experience

Sign in enforcement only takes effect after a Docker Desktop restart. The
following sections outline the developer experience after sign in is enforced
and Docker Desktop is restarted.

### First time sign in on enforced machine

When a user opens Docker Desktop on an enforced machine, they see a sign-in
prompt explaining that their organization requires authentication. After
signing in, users can continue using Docker Desktop immediately.

> [!NOTE]
>
> Sign-in enforcement only takes effect after Docker Desktop is restarted.

### After sign in

Once users sign in to enforced machines:

- With verified domains and auto-provisioning enabled: Users are automatically
added to your organization. For more information on verifying a domain and
enabling auto-provisioning, see [Domain management](/manuals/security/for-admins/domain-management.md).
- Without auto-provisioning: User emails appear in your the machines management
view for manual review and addition. To add a user to your organization,
see [Add unassociated machines to your organization](#add-unassociated-machines-to-your-organization).

## Troubleshooting

For common issues and solutions, see [Troubleshoot unassociated machines](/manuals/security/troubleshoot/troubleshoot-unassociated-machines.md).
