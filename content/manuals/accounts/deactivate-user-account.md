---
title: Deactivate a Docker account
linkTitle: Deactivate an account
weight: 30
description: Learn how to deactivate a Docker user account.
keywords: Docker Hub, delete, deactivate, account, account management, delete Docker account, close Docker account, disable Docker account
---

Learn how to deactivate an individual Docker account, including prerequisites required
for deactivation.

For information on deactivating an organization,
see [Deactivating an organization](../admin/organization/deactivate-account.md).

> [!WARNING]
>
> All Docker products and services that use your Docker account are
inaccessible after deactivating your account.

## Prerequisites

Before deactivating your Docker account, ensure you meet the following requirements:

- If you are an organization or company owner, you must leave your organization
or company before deactivating your Docker account:
    1. Sign in to [Docker Home](https://app.docker.com/admin) and choose
    your organization.
    1. Select **Members** and find your username.
    1. Select the **Actions** menu and then select **Leave organization**.
- If you are the sole owner of an organization, you must assign the owner role
to another member of the organization and then remove yourself from the
organization, or deactivate the organization. Similarly, if you are the sole
owner of a company, either add someone else as a company owner and then remove
yourself, or deactivate the company.
- If you have an active Docker subscription, [downgrade it to a Docker Personal subscription](../subscription/change.md).
- Download any images and tags you want to keep. Use `docker pull -a <image>:<tag>`.
- Unlink your [GitHub and account](../docker-hub/repos/manage/builds/link-source.md#unlink-a-github-user-account).

## Deactivate

Once you have completed all the previous steps, you can deactivate your account.

> [!WARNING]
>
> Deactivating your account is permanent and can't be undone. Make sure
to back up any important data.

1. Sign in to [Docker Home](https://app.docker.com/login).
1. Select your avatar to open the drop-down menu.
1. Select **Account settings**.
1. Select **Deactivate**.
1. Select **Deactivate account**, then select again to confirm.

## Delete personal data

Deactivating your account does not delete your personal data. To request
personal data deletion, fill out Docker's
[Privacy request form](https://preferences.docker.com/).
