---
description: How to deregister swarms on Docker Cloud
keywords: cloud, swarm, migration
title: Deregister Swarms on Docker Cloud
---

## Introduction

This page explains how to deregister a Swarm cluster from Docker Cloud so that it can be managed independently. We explain how to deregister on both Amazon Web Services (AWS) and Microsoft Azure (because Docker Cloud swarms run on either AWS or Azure behind the scenes).

You do not need to migrate or reconfigure your applications as part of this procedure. The only thing that changes is that your Swarm cluster no longer integrates with Docker services (such as Docker Cloud, Docker for Mac, or Docker for Windows).

### Prerequisites

To complete this procedure you need:

- An AWS or Azure account that lets you inspect resources such as instances.

### High-level steps

- Verify that you can SSH to your Swarm nodes (on AWS and Azure).
- Deregister your Swarm from Docker Cloud.
- Clean up old Docker Cloud resources.

## SSH to your Swarm

It is vital that you can SSH to your Docker Cloud Swarm before you deregister it from Docker Cloud.

Your Docker Cloud Swarm runs on either AWS or Azure, so to SSH to your Swarm nodes, you must know the public IP addresses or public DNS names of your nodes. The simplest way to find this information is with the native AWS or Azure tools.

### How to SSH to AWS nodes

1.  Log on to the AWS console and open the **EC2 Dashboard** for the **region** that hosts your Swarm nodes.

2.  Locate your instances and note their DNS names and IPs.

    By default, AWS labels your Swarm nodes as _swarm-name_-worker or _swarm-name_-manager. For example, a Swarm called "prod-equus" in Docker Cloud, has manager and worker nodes in AWS labelled, "prod-equus-manager" and "prod-equus-worker" respectively.

    You will also have a load balancer (type=classic) that includes the name of the Swarm. It accepts Docker commands on port 2376 and balances them to the manager nodes in the Swarm (as the server proxy is only deployed on the managers).

3.  Open an SSH session to each node in the cluster.

    This example opens an SSH session to a Swarm node with:

    - Private key = “awskey.pem”
    - Username = “docker”
    - Public DNS name = “ec2-34-244-56-42.eu-west-1.compute.amazonaws.com”

    ```
    $ ssh -i ./awskey.pem docker@ec2-34-244-56-42.eu-west-1.compute.amazonaws.com
    ```

Once you are certain that you are able to SSH to _all nodes_ in your Swarm, you can [deregister from Docker Cloud](#deregister-swarm-from-docker-cloud).

> If you do not have the keys required to SSH on to your nodes, you can deploy new public keys to your nodes using [this procedure](https://github.com/docker/dockercloud-authorizedkeys/blob/master/README.md){: target="_blank" class="_"}. You should perform this operation before deregistering your Swarm from Docker Cloud.

### How to SSH to Azure nodes

In Azure, you can only SSH to manager nodes because worker nodes do not get public IPs and public DNS names. If you need to log on to worker nodes, you can use your manager nodes as jump hosts.

1.  Log on to the Azure portal and click **Resource groups**.

2.  Click on the resource group that contains your Swarm. The `DEPLOYMENT NAME` should match the name of your Swarm.

3.  Click into the deployment with the name of your Swarm and verify the values. For example, the `DOCKERCLOUDCLUSTERNAME` value under **Inputs** should exactly match the name of your Swarm as shown in Docker Cloud.

4.  Copy the value from `SSH TARGETS` under **Outputs** and paste it into a new browser tab.

    This takes you to the inbound NAT Rules for the external load balancer that provides SSH access to your Swarm. It displays a list of all of the **Swarm managers** (not workers) including public IP address (`DESTINATION`) and port (`SERVICE`) that you can use to gain SSH access.

5.  Open an SSH session to each manager in the cluster. Use public IP and port to connect.

    This example creates an SSH session with user `docker` to a swarm manager at `51.140.229.154` on port `50000` with the `azkey.pem` private key in the current directory.

    ```
    ssh -i ./azkey.pem -p 50000 docker@51.140.229.154
    ```

    > If you do not know which private key to use, you can see the public key under `SSHPUBLICKEY` in the **Outputs** section of the Deployment. You can compare this value to the contents of public keys you have on file.

6.  Log on to your worker nodes by using your manager nodes as jump hosts. With
    [SSH agent forwarding enabled](https://docs.docker.com/docker-for-azure/deploy/#connecting-to-your-linux-worker-nodes-using-ssh), SSH from the manager nodes to the workers nodes over the private network.

Once you are certain that you are able to SSH to the manager nodes in your Swarm you can [deregister from Docker Cloud](#deregister-swarm-from-docker-cloud).

> If you do not have the keys required to SSH on to your nodes, you can deploy new public keys to your nodes using [this procedure](https://github.com/docker/dockercloud-authorizedkeys/blob/master/README.md){: target="_blank" class="_"}. You should perform this operation before deregistering your Swarm from Docker Cloud.

## Deregister swarm from Docker Cloud

> Proceed with caution
>
> Only deregister if you know the details of your Swarm nodes (cloud provider, public DNS names, public IP address, etc.) and you have verified that you can SSH to each node with your private key.
{: .warning}

1.  Open the Docker Cloud web UI and click **Swarms**.

2.  Click the three dots to the right of the Swarm you want to deregister and select **Unregister**.

3.  Confirm the deregistration process.

The Swarm is now deregistered from the Docker Cloud web UI and no longer is visible in other products such as Docker for Mac and Docker for Windows.

## Clean up Docker Cloud resources

The final step  is to clean up old Docker cloud resources such as the service, network and secret.

Docker Cloud deployed a service on your Swarm called `dockercloud-server-proxy` to proxy and load balance incoming Docker commands on port 2376 across all manager nodes. It has a network called `dockercloud-server-proxy-network` and a secret called `dockercloud-server-proxy-secret`.

All of these should be removed:

1.  Open an SSH session to a Swarm manager _for the correct swarm!_

2.  Remove the service:

    ```
    $ docker service rm dockercloud-server-proxy
    ```

3.  Remove the network:

    ```
    $ docker network rm dockercloud-server-proxy-network
    ```

4.  Remove the secret:

    ```
    $ docker secret rm dockercloud-server-proxy-secret
    ```

Your Docker Swarm cluster is now deregistered from Docker Cloud and you can manage it independently.
