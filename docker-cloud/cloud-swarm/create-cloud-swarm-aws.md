---
previewflag: cloud-swarm
description: Create new swarms on AWS with Docker Cloud
keywords: swarm mode, swarms, create swarm, Cloud, AWS
title: Create a new swarm on Amazon Web Services in Docker Cloud
---

{% include content/cloud-swarm-overview.md %}

## Link your service provider to Docker Cloud

To create a swarm, you need to give Docker Cloud permission to deploy swarm
nodes on your behalf in your cloud services provider account.

If you haven't yet linked Docker Cloud to AWS, follow the steps in [Link Amazon Web Services to Docker Cloud](link-aws-swarm.md). Once it's
linked, it will show up on the **Swarms -> Create** page as a connected service
provider.

![](images/aws-creds-cloud.png)

## Create a swarm

1. If necessary, log in to Docker Cloud and switch to Swarm Mode

2. Click **Swarms** in the top navigation, then click **Create**.

    >**Tip:** Alternatively, you can select **+ -> Swarm** from the top navigation to get to the same page.

3. Enter a name for the new swarm.

4. Select a connected cloud services provider.

    <font style="color:red;">TBD: MAKE THIS SPECIFIC TO AWS AND ADD IMAGES THROUGHOUT</font>

    <font style="color:red;">TBD from here down, Add info re: VPC for Region Advanced Settings, and other generic or AWS specific configurations. Generic configurations seem to be Swarm size and properties, along with manager and worker properties. (See Azure steps, which briefly cover these generic configs.)</font>

5. Enter any additional provider-specific information, such as region.

    > **Note:** For Amazon Web Services, the SSH keys that appear in this wizard are filtered by the region you select. Make sure that you have appropriate SSH keys available on the region you select.

6. Choose how many swarm managers and swarm worker nodes to deploy.

7. Select the instance sizes for the managers, and for the workers.

8. Select the SSH key to use to connect to the nodes.

    The list contains any SSH keys that you have access to on your linked cloud services provider. Select the one for which you have the private key locally.

9. Click **Create**.

    Docker for AWS bootstraps all of the recommended infrastructure to
    start using Docker on AWS automatically. You don't need to worry
    about rolling your own instances, security groups, or load balancers
    when using Docker for AWS. (To learn more, see
    [Why Docker for AWS](/docker-for-aws/why.md).)

> **Note**: At this time, you cannot add or remove nodes from a swarm from within Docker Cloud. To add new nodes or remove nodes from an existing swarm,
log in to your AWS account, and add or delete nodes manually. (You can
unregister or dissolve swarms directly from Docker Cloud.)

## Where to go next

Learn how to [connect to a swarm through Docker Cloud](connect-to-swarm.md).

Learn how to [register existing swarms](register-swarms.md).

You can get an overivew of topics on [swarms in Docker Cloud](index.md).

To find out more about Docker swarm in general, see the Docker engine
[Swarm Mode overview](/engine/swarm/).
