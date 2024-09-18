---
description: Learn how to add and remove company owners.
keywords: company, owners
title: Manage company owners
aliases:
- /docker-hub/company-owner/
---

A company can have multiple owners. Company owners have company-wide
observability and can manage company-wide settings that apply to all associated
organizations. In addition, company owners have the same access as organization
owners for all associated organizations. Unlike organization owners, company
owners don't need to be member of an organization. When company owners aren't a
member in an organization, they don't occupy a seat.

As a company owner, you can configure [Single Sign-on (SSO)](../../security/for-admins/single-sign-on/_index.md) and [System for Cross-domain Identity Management (SCIM)](../../security/for-admins/provisioning/scim.md) for all organizations under the company.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

## Add a company owner

1. Sign in to [Docker Hub](https://hub.docker.com/), navigate to the **Organizations** page, and select your company.
2. In the **Overview** tab, select **Add owner** and enter their Docker ID or email address.
3. After you find the user, select **Add company owner**.

Once you've added the new owner, you can view the user in the **Company owners** table.

## Remove a company owner

1. Sign in to [Docker Hub](https://hub.docker.com/), navigate to the **Organizations** page, and select your company.
2. In the **Overview** tab, find the company owner you want to remove.
3. Select the **Action** icon and remove the owner.

    ![remove-owner](../images/remove-owner.png)

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

## Add a company owner

1. Sign in to the [Admin Console](https://admin.docker.com).
2. In the left navigation, select your company in the drop-down menu.
3. Select **Company owners**.
4. Select **Add owner**.
5. Specify the user's Docker ID to search for the user.
6. After you find the user, select **Add company owner**.

## Remove a company owner

1. Sign in to the [Admin Console](https://admin.docker.com).
2. In the left navigation, select your company in the drop-down menu.
3. Select **Company owners**.
4. Select the **Action** icon in the row of the company owner that your want to remove.
5. Select **Remove as company owner**.

{{< /tab >}}
{{< /tabs >}}
