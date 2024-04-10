---
description: Learn how to manage organization in a company.
keywords: company, multiple organizations, manage organizations
title: Manage organizations
---

You can manage the organizations in a company in Docker Hub and the Docker Admin Console. In Docker Hub, you can manage seats and members at the organization-level. When you use the Admin Console, you can do more to manage organizations at the company-level.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

## View all organizations

1. In Docker Hub, select **Organizations**.
2. Select your company.
3. From the company page, you can view all organizations in the **Overview** tab.

## Add organizations to a company

>**Important**
>
> You must be a company owner to add an organization to a company. You must also be an organization owner of the organization you want to add.
{ .important }

There is no limit to the number of organizations you can have under a company layer. All organizations must have a Business subscription.

>**Important**
>
> Once you add an organization to a company, you can't remove it from the company.
{ .important }

1. In Docker Hub, select **Organizations**.
2. Select your company.
3. From the company page, select **Add organization**.
4. Choose the organization you want to add from the drop-down menu.
5. Select **Add organization** to confirm.

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

## View all organizations

1. Sign in to the [Admin Console](https://admin.docker.com).
2. In the left navigation, select your company in the drop-down menu.
3. Under **Organizations**, select **Overview**.

## Add seats to an organization

When you have a [self-serve](../../subscription/core-subscription/details.md#self-serve) subscription that has no pending subscription changes, you can add seats using the following steps.

1. Sign in to the [Admin Console](https://admin.docker.com).
2. In the left navigation, select your company in the drop-down menu.
3. Under **Organizations**, select **Overview**.
4. Select the action icon in the organization's card, and then select **Get more seats**.

## Add organizations to a company

>**Important**
>
> You must be a company owner to add an organization to a company. You must also be an organization owner of the organization you want to add.
{ .important }

There is no limit to the number of organizations you can have under a company layer. All organizations must have a Business subscription.

>**Important**
>
> Once you add an organization to a company, you can't remove it from the company.
{ .important }

1. Sign in to the [Admin Console](https://admin.docker.com).
2. In the left navigation, select your company in the drop-down menu.
3. Select **Add organization**.
4. Choose the organization you want to add from the drop-down menu.
5. Select **Add organization** to confirm.

## Manage an organization

1. Sign in to the [Admin Console](https://admin.docker.com).
2. In the left navigation, select your company in the drop-down menu.
3. Select the organization that you want to manage.

For more details about managing an organization, see [Organization administration](../organization/index.md).
{{< /tab >}}

{{< /tabs >}}

## More resources

- [Video demo: Managing a company and nested organizations](https://youtu.be/XZ5_i6qiKho?feature=shared&t=229)
- [Video demo: Adding nested organizations to a company](https://youtu.be/XZ5_i6qiKho?feature=shared&t=454)
