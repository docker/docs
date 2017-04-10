---
previewflag: cloud-swarm
description: Link your Microsoft Azure account
keywords: Azure, Cloud, link
title: Link Microsoft Azure Cloud Services to Docker Cloud
---

You can link your Microsoft Azure account so that Docker Cloud can provision and
manage swarms on your behalf.

For this, you will need your Azure subscription ID to authenticate Docker to your service provider.

> **Note**: When you are ready to create and deploy swarms, you must also have an [SSH key](ssh-key-setup.md).

## Find your Azure subscription ID

You will also need your Azure Cloud Services subscription ID to provide to
Docker Cloud. There are a few ways to navigate to it on Azure.

You can click a resource from the Dashboard and find the subscription ID under
"Essentials" on the resulting display. Alternatively, from the left menu, go to
**Billing -> Subscriptions -> Subscription ID** or simply click
**Subscriptions**, then click a subscription in the list to drill down.

![](images/azure-subscription-id.png)

When you are ready to add your subscription ID to Docker Cloud,
copy it from your Azure Dashboard.

## Add you Azure account credentials to Docker Cloud

Go to Docker Cloud to connect the account.

1. In Docker Cloud, click the account menu at upper right and select **Cloud settings**.
2. In the **Service Providers** section, click the plug icon next to Microsoft Azure.

    ![](images/azure-id-wizard.png)

3. Provide your subscription ID.

    You will be redirected to [Azure Cloud Services](portal.azure.com).

4. Log in to your Azure account.

5. Click **Accept** to grant Docker Cloud access to your Microsoft Azure account.

    ![](images/azure-permissions.png)

6. Your Microsoft Azure login credentials will automatically populate to
Docker Cloud under **Service Providers -> Microsoft Azure**.

7. Click **Save**.

You're now ready to deploy a swarm!

## Where to go next

You'll need an SSH key to provide to Docker Cloud
during the swarm create process. See [Setting up SSH keys](ssh-key-setup.md).

**Ready to create swarms on Azure?** See [Create a new swarm in Docker
Cloud](create-cloud-swarm.md).

You can get an overivew of topics on [swarms in Docker Cloud](index.md).
