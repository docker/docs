---
description: Audit your domains for uncaptured users.
keywords: domain audit, security
title: Domain audit
---

> **Note**
>
> Domain audit is currently in [Early Access](../release-lifecycle.md/#early-access-ea).
> The feature is enabled for specific user groups as part of an incremental roll-out strategy.

When your organization has configured SSO, and you have verified your domains, you can audit your domains. Auditing your domains will identify uncaptured users that have authenticated with an email associated with one of your verified domains.

Uncaptured users can pose a security threat to your environment since your organization's security settings aren't applied to the user's sessions who aren't part of your organization. In addition, you won't have visibility into the activity of uncaptured users.

You can add uncaptured users to your organization to gain visibility into their activity and apply your organization's security settings. Additionally, you can enforce sign-in to ensure that only members of your organization can sign in to Docker Desktop in your environment. For more details about enforcing sign-in, see [Configure registry.json to enforce sign-in](../docker-hub/configure-sign-in.md).

## Audit your domains for uncaptured users

To audit your domains:

1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator of your organization.

2. Select **Organizations** and then select your organization.

3. Select **Settings** and then select **Security**.

4. In **Domain Audit**, select **Export Users** to export a CSV file of uncaptured users with the following columns:
  - Name: The name of the user.
  - Username: The Docker ID of the user.
  - Email: The email address of the user.
  - Date Joined: The date the user created their Docker account.

You can invite all the uncaptured users to your organization using the exported CSV file. For more details, see [Invite members via CSV file](../docker-hub/members.md/#invite-members-via-csv-file).
