---
title: Manage unassociated machines
description: Learn how to manage unassociated machines using the Docker Admin Console
keywords: unassociated machines, insights, manage users, enforce sign-in
weight: 56
---

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

### In the Insights dashboard

The Insights dashboard shows high-level metrics for both associated and
unassociated machines:

1. Navigate to your organization in Docker Hub.
2. Go to the Insights tab.
3. View the summary of:
	- Total active users (associated + unassociated)
	- Associated organization members
	- Unassociated machines detected

> [!NOTE]
>
> Full Insights dashboard features become available when you enable
sign-in enforcement for your organization.

For more information, see [Insights](/manuals/admin/organization/insights.md).

### In the machines management view

To see detailed information about unassociated machines:

1. Navigate to your organization in Docker Hub.
2. Go to Settings > General.
3. Select Unassociated machines.

The machine list displays:

- Machine ID (Docker-generated identifier)
- Registry addresses accessed (when available)
- Last activity date
- Docker Desktop version
- User email (after sign-in enforcement)

You can:

- Export the list as CSV
- Filter and sort machines by activity
- Take actions on individual or multiple machines

## Manage unassociated machines

### Enable sign-in enforcement

You can require users to sign in to Docker Desktop.

For all unassociated machines:

1. In your organization settings, go to Unassociated machines.
2. Select Enforce sign-in for all machines.
3. Confirm the action.

For specific machines:

1. In the unassociated machines list, select individual machines.
2. Choose Require sign-in from the actions menu.

### Manually add users

To manually add users:

1. Go to Settings > General > Unassociated machines.
2. Review users who have signed in (identified by email addresses).
3. Select users to add to your organization.
4. Choose Add to organization.

## User experience

Sign in enforcement only take effect after a Docker Desktop restart. The
following sections outline the user experience after sign in is enforced
and Docker Desktop restarted.

### First time sign in on enforced machine

When a user opens Docker Desktop on an enforced machine:

1. They see a customizable prompt explaining that their organization requires
sign-in.
2. The prompt includes information that their email will be shared with
organization administrators.
3. Users can continue using Docker Desktop immediately after signing in.
4. Users are not blocked based on license availability.

### After sign in

Once users sign in to enforced machines:

- With verified domains and auto-provisioning enabled: Users are automatically
added to your organization.
	- For more information on verifying a domain and enabling auto-provisioning,
    see [Domain management](/manuals/security/for-admins/domain-management.md).
- Without auto-provisioning: User emails appear in your the machines management
view for manual review and addition.

## Troubleshooting

For common issues and solutions, see [Troubleshoot unassociated machines](/manuals/security/troubleshoot/troubleshoot-unassociated-machines.md).
