---
title: Create a company
description: Learn how to create a company.
keywords: company, hub, organization, company owner, Admin Console, company management
---

You can create a new company in Docker Hub or Docker Admin Console. Before you begin, make sure you're the owner of the organization you want to add to the new company. The organization also needs to have a Docker Business subscription.

{{< tabs >}}

{{< tab name="Docker Hub" >}}

## Create a company

To create a new company:

1. In Docker Hub, navigate to the organization you want to place under a company. The organization must have a Business subscription, and you must be an owner of the organization.
2. Select **Settings**.
3. Near the bottom of the **General** tab, select **Create a company**.
4. Enter a unique name for your company, then select **Continue**.

    > **Tip**
    >
    > The name for your company can't be the same as an existing user, organization, or company namespace.
    { .tip }

5. Review the company migration details and then select **Create company**.

## Add organizations to a company

>**Important**
>
> You must be a company owner to add an organization to a company. You must also be an organization owner of the organization you wish to add.
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

## Create a company

To create a new company:

1. In the Admin Console, navigate to the organization you want to place under a company. The organization must have a Business subscription, and you must be an owner of the organization.
2. Under **Organization Settings**, select **General**.
3. In the **Organization management** section, select **Create a company**.
4. Enter a unique name for your company, then select **Continue**.

    > **Tip**
    >
    > The name for your company can't be the same as an existing user, organization, or company namespace.
    { .tip }

5. Review the company migration details and then select **Create company**.

For more information on how you can add organizations to your company, see [Add organizations to a company](./organizations.md#add-organizations-to-a-company).

{{< /tab >}}

{{< /tabs >}}
