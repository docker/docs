---
title: Deactivate an organization
description: Learn how to deactivate a Docker organization.
keywords: Docker Hub, delete, deactivate organization, account, organization management
weight: 42
aliases:
- /docker-hub/deactivate-account/
---

{{< summary-bar feature_name="General admin" >}}

You can deactivate an account at any time. This section describes the prerequisites and steps to deactivate an organization account. For information on deactivating a user account, see [Deactivate a user account](../../accounts/deactivate-user-account.md).

> [!WARNING]
>
> All Docker products and services that use your Docker account or organization account will be inaccessible after deactivating your account.

## Prerequisites

Before deactivating an organization, complete the following:

- Download any images and tags you want to keep:
  `docker pull -a <image>:<tag>`.
- If you have an active Docker subscription, [downgrade it to a free subscription](../../subscription/change.md).
- Remove all other members within the organization.
- Unlink your [Github and Bitbucket accounts](../../docker-hub/repos/manage/builds/link-source.md#unlink-a-github-user-account).
- For Business organizations, [remove your SSO connection](../../security/for-admins/single-sign-on/manage/#remove-an-organization).

## Deactivate

Once you have completed all the previous steps, you can deactivate your organization.

> [!WARNING]
>
> This cannot be undone. Be sure you've gathered all the data you need from your organization before deactivating it.

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Sign in to [Docker Home](https://app.docker.com) and select the organization
you want to deactivate.
1. Select **Admin Console**, then **Deactivate**. If this button is greyed out,
you must complete the [Prerequisites](#prerequisites).
1. Enter the organization name to confirm deactivation.
1. Select **Deactivate organization**.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Choose the organization you want to deactivate.
1. In **Settings**, select **Deactivate org**.
1. Select **Deactivate organization**.

{{< /tab >}}
{{< /tabs >}}
