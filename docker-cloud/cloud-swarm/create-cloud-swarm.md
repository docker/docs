---
description: Create new swarms in Docker Cloud
keywords: swarm mode, swarms, create swarm, Cloud
title: Create a new swarm in Docker Cloud
---

<b>Note</b>: All Swarm management features in Docker Cloud are free while in Beta.
{: style="text-align:center" }

--------

You can now create _new_ Docker Swarms from within Docker Cloud as well as
register existing swarms.

When you create a swarm, Docker Cloud connects to the provider on your behalf,
and uses the provider's APIs and a provider-specific template to launch Docker
instances. The instances are then joined to a swarm and the swarm is configured
using your input. When you access the swarm from Docker Cloud, the system
forwards your commands directly to the Docker instances running in the swarm.

> **Note**: The Create Swarm functionality is only available for AWS at this time. Additional provider support is coming soon.

## Create a swarm

> **Note**: To create a swarm, you need to give Docker Cloud permission to deploy swarm nodes on your behalf in your cloud services provider account. See the [AWS with swarm instructions](link-aws-swarm.md) to learn more.

1. If necessary, log in to Docker Cloud and switch to Swarm mode.
2. Click Swarms in the top navigation.
3. Click **Create**.
4. Enter a name for the new swarm.
5. Select a connected cloud services provider.

    Additional options appear depending on which provider you select.

6. Enter any additional provider-specific information, such as region.

    > **Note**: The SSH keys that appear in this wizard filtered by the region you select. Ensure that you have appropriate SSH keys available on the region you select.

7. Choose how many swarm managers and swarm worker nodes to deploy.
8. Select the instance sizes for the managers, and for the workers.
9. Select the SSH key to use to connect to the nodes.

    The list contains any SSH keys that you have access to on your linked cloud services provider. Select the one for which you have the private key locally.

10. Click **Create**.

Docker Cloud connects to your AWS account, deploys Docker for AWS instances, forms a Swarm, and joins the instances to it.

> **Note**: At this time, you cannot dissolve swarms or delete swarm nodes from within Docker Cloud - you can only unregister the swarm from the Docker Cloud UI. To delete the swarm and its members, log in to your AWS account and delete them manually.
