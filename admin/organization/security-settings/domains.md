---
description: Domain management in Docker Admin
keywords: domains, SCIM, SSO, Docker Admin, domain audit
title: Domain management
---

{% include admin-early-access.md %}

Use domain management to manage your domains for Single Sign-On and SCIM, as well as audit your domains for uncaptured users.

## Add and verify a domain

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Domain management**.
4. Select **Add a domain** and continue with the on-screen instructions to add the TXT Record Value to your domain name system (DNS).

    >**Note**
    >
    > Format your domains without protocol or www information, for example, `yourcompany.example`. This should include all email domains and subdomains users will use to access Docker, for example `yourcompany.example` and `us.yourcompany.example`. Public domains such as `gmail.com`, `outlook.com`, etc. aren’t permitted. Also, the email domain should be set as the primary email.

5. Once you have waited 72 hours for the TXT Record verification, you can then select **Verify** next to the domain you've added, and follow the on-screen instructions.


## Domain audit

> **Note**
>
> Domain audit is currently in [early access](../../../release-lifecycle.md#early-access-ea).

Domain audit identifies uncaptured users. Uncaptured users are Docker users who have authenticated to Docker using an email address associated with one of your verified domains, but they're not a member of your organization in Docker. You can audit domains on organizations that are part of the Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](../../../subscription/upgrade.md).

Uncaptured users who access Docker Desktop in your environment may pose a security risk because your organization's security settings, like Image Access Management and Registry Access Management, aren't applied to a user's session. In addition, you won't have visibility into the activity of uncaptured users. You can add uncaptured users to your organization to gain visibility into their activity and apply your organization's security settings.

Domain audit can't identify the following Docker users in your environment:
   * Users who access Docker Desktop without authenticating
   * Users who authenticate using an account that doesn't have an email address associated with one of your verified domains

Although domain audit can't identify all Docker users in your environment, you can enforce sign-in to prevent unidentifiable users from accessing Docker Desktop in your environment. For more details about enforcing sign-in, see [Configure registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md).

### Audit your domains for uncaptured users

Before you audit your domains, the following prerequisites are required:
   * Your organization must be part of a Docker Business subscription. To upgrade your existing account to a Docker Business subscription, see [Upgrade your subscription](../../../subscription/upgrade.md).
   * You must add and verify your domains. To add and verify a domain, see [Add and verify a domain](#add-and-verify-a-domain).

To audit your domains:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Domain management**.
4. In **Domain Audit**, select **Export Users** to export a CSV file of uncaptured users with the following columns:
  - Name: The name of the user.
  - Username: The Docker ID of the user.
  - Email: The email address of the user.

You can invite all the uncaptured users to your organization using the exported CSV file. For more details, see [Invite members via CSV file](../../../docker-hub/members.md#invite-members-via-csv-file). Optionally, enforce single sign-on or enable SCIM to add users to your organization automatically. For more details, see [Single Sign-on](sso.md) or [SCIM](scim.md).

> **Note**
>
> Domain audit may identify accounts of users who are no longer a part of your organization. If you don't want to add a user to your organization and you don't want the user to appear in future domain audits, you must deactivate the account or update the associated email address.
>
> Only someone with access to the Docker account can deactivate the account or update the associated email address. For more details, see [Deactivating an account](../../../docker-hub/deactivate-account.md).
>
> If you don't have access to the account, you can contact [Docker support](../../../support/index.md) to discover if more options are available.