---
description: Migrating from Docker Cloud
keywords: cloud, swarm, migration
title: Deregister swarms on Docker Cloud
---

## Deregistering Swarms

This page explains how to deregister a Swarm cluster from Docker Cloud so that it can be managed independently. We explain how to deregister on both Amazon Web Services (AWS) and Microsoft Azure (because Docker Cloud swarms run on either AWS or Azure behind the scenes).

You do not need to migrate or configure your applications as part of this procedure. The only thing that changes is that your Swarm cluster no longer integrates with Docker Cloud, Docker for Mac, or Docker for Windows, and other Docker services.

### Prerequisites

To complete this procedure you need:

- An AWS or Azure account that lets you inspect resources such as instances.

### High-level steps

- Verify that you can SSH to your Swarm (on AWS and Azure)
- Deregister your Swarm from Docker Cloud
- Clean up old Docker Cloud resources.

## Verify you can SSH to your Swarm

> Your Docker Cloud Swarm runs on either AWS or Azure, and it is vital that you can SSH to it before you deregister it from Docker Cloud. To SSH to your Swarm nodes, you must know the public IP addresses or public DNS names of your nodes. The simplest way to find this information is via the native AWS or Azure tools.

### How to SSH to AWS nodes

1.  Log on to the AWS console and open the EC2 Dashboard for the region that hosts your Swarm nodes.

2.  Locate your instances and note their hostnames, DNS names, and IPs.

    By default, AWS labels your Swarm nodes as _swarm-name_-worker or _swarm-name_-manager. For example, a Swarm called "prod-equus" in Docker Cloud, has manager and worker nodes in AWS labelled, "prod-equus-manager" and "prod-equus-worker" respectively.

    You will also have a load balancer (type=classic) that includes the name of the Swarm. It accepts Docker commands on port 2376 and balances them to all nodes in the Swarm.

3.  Open an SSH session to each node in the cluster.

    The example below opens an SSH session to a Swarm node with a public DNS name of “ec2-34-244-56-42.eu-west-1.compute.amazonaws.com” using a private key called “awskey.pem”. It also logs on using the “docker” username.

    ```
    $ ssh -i ./awskey.pem docker@ec2-34-244-56-42.eu-west-1.compute.amazonaws.com
    ```

### How to SSH to Azure nodes

1.  Log on to the Azure portal and click **Resource groups**.

2.  Click on the resource group that contains your Swarm. The “DEPLOYMENT NAME” should match the name of your Swarm.

3.  Click into the deployment with the name of your Swarm and verify the values.  For example, the “DOCKERCLOUDCLUSTERNAME” value under “Inputs” should exactly match the name of your Swarm as shown in Docker Cloud.

4.  Copy the value from “SSH TARGETS” under “Outputs” and paste it into a new browser tab.

   This takes you to the inbound NAT Rules for the external load balancer that provides SSH access to your Swarm. It displays a list of all of the Swarm managers including public IP address (DESTINATION) and port (SERVICE) that you can use to gain SSH access.

5.  Open an SSH session to each node in the cluster.

    Use a combination of the public IP and port to make an SSH connection. The following command shows how to create an SSH session to a manager at 51.140.229.154 on port 50000 using the “azkey.pem” private key in the current directory. It also logs on using the “docker” username.

    ```
    ssh -i ./azkey.pem -p 50000 docker@51.140.229.154
    ```

If you are not certain which private key you use to connect to your Swarm, you can see the public key under “SSHPUBLICKEY” in the Outputs section of the Deployment. You can compare this value to the contents of public keys you have on file.

Once you are certain you can gain SSH access to the nodes in your Swarm you can proceed to the next section

## Deregister swarm from Docker Cloud

> Proceed with caution
>
> Only proceed with this section if you know the details of your Swarm nodes (cloud provider, public DNS, public IP etc.) and you have verified that you can SSH to each node with your private key.
{: .warning}

1.  Open the Docker Cloud web UI and click **Swarms**.

2.  Click the three dots to the right of the Swarm you want to deregister and select **Unregister**.

3.  Confirm the deregistration process.

The Swarm is now deregistered from the Docker Cloud web UI and no longer is visible in other products such as Docker for Mac and Docker for Windows.

## Clean up Docker Cloud resources

The final step  is to clean up old Docker cloud resources.

Docker Cloud deployed a service on your Swarm called **dockercloud-server-proxy**. This was used to proxy and load balance incoming Docker commands on port 2376 across all manager nodes. It has a network called `dockercloud-server-proxy-network` and secret called `dockercloud-server-proxy-secret`. All of these should be removed.

The following commands show you how to remove the service, network and secret. Be sure to perform these commands against the correct Swarm.

1.  Open an SSH session to a Swarm manager.

2.  Remove the service.

    ```
    $ docker service rm dockercloud-server-proxy
    ```

3.  Remove the network.

    ```
    $ docker network rm dockercloud-server-proxy-network
    ```

4.  Remove the secret.

    ```
    $ docker secret rm dockercloud-server-proxy-secret
    ```

Your Docker Swarm cluster is now deregistered from Docker Cloud and you can manage it independently.
