---
description: Audit your domains for uncaptured users.
keywords: domain audit, security
title: Domain audit
---

> **Note**
>
> Domain audit is currently in [Early Access](../release-lifecycle.md/#early-access-ea).

Domain audit identifies uncaptured users. Uncaptured users are Docker users that have authenticated to Docker using an email address associated with one of your verified domains, but they're not a member of your organization in Docker. You can audit domains on organizations that are part of the Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](../subscription/upgrade.md).

Uncaptured users who access Docker Desktop in your environment may pose a security threat because your organization's security settings, like Image Access Management and Registry Access Management, aren't applied to a user's session. In addition, you won't have visibility into the activity of uncaptured users. You can add uncaptured users to your organization to gain visibility into their activity and apply your organization's security settings.

Domain audit can't identify the following Docker users in your environment:
   * Users that access Docker Desktop without authenticating
   * Users that authenticate using an account that doesn't have an email address associated with one of your verified domains

Although domain audit can't identify all Docker users in your environment, you can enforce sign-in to prevent unidentifiable users from accessing Docker Desktop in your environment. For more details about enforcing sign-in, see [Configure registry.json to enforce sign-in](../docker-hub/configure-sign-in.md).

## Audit your domains for uncaptured users

Before you audit your domains, the following prerequisites are required:
   * Your organization must be part of a Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](../subscription/upgrade.md).
   * Single sign-on must be configured for your organization. To configure single sign-on, see [Configure Single Sign-on](../single-sign-on/configure/index.md).
   * You must add and verify your domains. To add and verify a domain, see [Domain control](../single-sign-on/configure/index.md/#domain-control).

To audit your domains:

1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator of your organization.

2. Select **Organizations** and then select your organization.

3. Select **Settings** and then select **Security**.

4. In **Domain Audit**, select **Export Users** to export a CSV file of uncaptured users with the following columns:
  - Name: The name of the user.
  - Username: The Docker ID of the user.
  - Email: The email address of the user.

You can invite all the uncaptured users to your organization using the exported CSV file. For more details, see [Invite members via CSV file](../docker-hub/members.md/#invite-members-via-csv-file). Optionally, enforce single sign-on or enable SCIM to add users to your organization automatically. For more details, see [Single Sign-on](../single-sign-on/index.md) or [SCIM](../docker-hub/scim.md).

> **Note**
>
> Domain audit may identify accounts of users who are no longer a part of your organization. If you don't want to add a user to your organization and you don't want the user to appear in future domain audits, you must deactivate the account or update the associated email address.
>
> Only someone with access to the Docker account can deactivate the account or update the associated email address. For more details, see [Deactivating an account](../docker-hub/deactivate-account.md/).
>
> If you don't have access to the account, you can contact [Docker support](../support/index.md) to discover if more options are available.


