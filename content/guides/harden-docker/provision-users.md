---
title: Provision users with least privilege
description:
weight: 30
---

Granting the right level of access to each user in your Docker organization
helps prevent accidental misconfiguration, enforces separation of duties, and
reduces risk.

This module shows you how to assign roles, restrict self-invites, and use group
mapping for automated, least-privilege access.

## Prerequisites

Before you begin, ensure you have:

- A Docker Business subscription
- Organization owner access to your Docker organization
- A verified domain and SSO configured
- [Group mapping](https://docs.docker.com/security/for-admins/provisioning/group-mapping/) (optional but recommended)

## Step one: Assign roles based on job function

Docker provides three organization-level roles:

- Member: Default role for most users.
- Editor: Manage repositories, settings, and access for assigned teams.
- Organization owner: Full administrator access.

To assign roles:

1. Sign in to the [Admin Console](https://app.docker.com/admin) and choose your
organization from the **Choose profile** page.
2. Navigate to **Members** and select the **Actions** menu next to a user.
3. Select **Edit** and assign a role.

For more details, see [Roles and permissions]().

## Step two: Use teams to organize access

Organize users into teams to control access to Docker Hub repositories,
products, and settings:

1. In the [Admin Console](https://app.docker.com/admin), navigate to **Teams**.
2. Create a new team or edit an existing one.
3. Add members manually or automatically via group mapping.

## Step three: Automate role assignment with group mapping

If you set up SSO and group mapping, Docker can assign users to teams automatically based on their group in your IdP.

Use the `org:team` naming convention in Docker and your IdP for group names.
