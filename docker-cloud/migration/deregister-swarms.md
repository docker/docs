---
description: Migrating from Docker Cloud
keywords: cloud, swarm, migration
title: Deregister swarms from Docker Cloud
---

## Deregistering Swarms

This page explains how to deregister a Swarm cluster from Docker Cloud so that it can be managed independently. We show how to do this for Swarms on AWS and Azure.

You do not need to migrate or configure your applications as part of this procedure. The only thing that changes is that your Swarm cluster no longer integrates with Docker Cloud, Docker for Mac, or Docker for Windows, and other Docker services.

## Prerequisites

To complete this procedure you need:

- An AWS or Azure account that lets you inspect resources such as instances

## High-level steps

High-level steps to deregister your Swarm cluster from Docker Cloud are:

- Verify you can SSH to your Swarm
- Deregister your Swarm from Docker Cloud
- Clean up old Docker Cloud resources.

### Verify you can SSH to your Swarm

Your Docker Cloud Swarm runs on either AWS or Azure, and it is vital that you can SSH to it before you deregister it from Docker Cloud.

To SSH to your Swarm nodes, you must know the public IP addresses or public DNS names of your nodes. The simplest way to find this information is via the native AWS or Azure tools.

#### AWS example

1.  Log on to the AWS console and open the EC2 Dashboard for the region host your Swarm nodes.

2.  Locate your instances and note their hostnames and IPs.

    By default, AWS labels your Swarm nodes as _swarm-name_-worker or _swarm-name_-manager. For example, a Swarm called "prod-equus" in Docker Cloud, has manager and worker nodes in AWS labelled, "prod-equus-manager" and "prod-equus-worker" respectively.

    You will also have a load balancer (type=classic) that includes the name of the Swarm. This is if accepts Docker commands on port 2376 and balances them to all nodes in the Swarm.

3. Note the public IP addresses and DNS names of each node in your Swarm.

	 You can also see the name of the public key used to access each node.

4.  Open an SSH session to a node in the cluster (you may want to test this on all nodes).

    The example below opens an SSH session to a Swarm node with a public DNS name of “ec2-34-244-56-42.eu-west-1.compute.amazonaws.com” using a private key called “awskey.pem”. It also logs on using the “docker” username.

    ```
    $ ssh -i ./awskey.pem docker@ec2-34-244-56-42.eu-west-1.compute.amazonaws.com
    ```

#### Azure example

1.  Log on to the Azure portal and click **Resource groups**.

2.  Click on the RG that contains your Swarm.

    If you don’t know which resource group has your Swarm, you’ll need to perform the following against each RG until you find it.

3.  From within the RG, click **Deployments**.

    The “DEPLOYMENT NAME” on the following screen will match the name of your Swarm.

4.  Click into the deployment with the name of your Swarm and verify the values displayed.

    For example, the “DOCKERCLOUDCLUSTERNAME” value under “Inputs” should exactly match the name of your Swarm as shown in Docker Cloud.

5. Copy the value from “SSH TARGETS” under “Outputs” and paste it into a new browser tab.

This takes you to the inbound NAT Rules for the external load balancer that provides SSH access to your Swarm. It displays a list of all of the Swarm managers including public IP address (DESTINATION) and port (SERVICE) that you can use to gain SSH access.

6.  Test an SSH connection.

    You use a combination of the public IP and port to make an SSH connection. The following command shows how to create an SSH session to a manager at 51.140.229.154 on port 50000 using the “azkey.pem” private key in the current directory. It also logs on using the “docker” username.

    ```
    ssh -i ./azkey.pem -p 50000 docker@51.140.229.154
    ```

If you are not certain which private key you use to connect to your Swarm, you can see the public key under “SSHPUBLICKEY” in the Outputs section of the Deployment. You can compare this value to the contents of public keys you have on file.

Once you are certain you can gain SSH access to the nodes in your Swarm you can proceed to the next section

### Deregister your Swarm from Docker Cloud

Only proceed with this section if you know the details of your Swarm nodes (cloud provider, public DNS, public IP etc.) and you have verified that you can SSH to each node with your private key.

1.  Open the Docker Cloud web UI and click `Swarms`.

2.  Click the three dots to the right of the Swarm you want to deregister and select `Unregister`.

3.  You are prompted to confirm the de-registration process.

The Swarm is now deregistered from the Docker Cloud web UI and will no longer be visible in other products such as Docker for Mac and Docker for Windows.

However, there are still a few leftovers that need cleaning up.

### Clean-up old Docker Cloud resources

The final step to deregister your Swarm with Docker Cloud is to clean-up old Docker cloud resources.

Docker Cloud deployed a service on your Swarm called **dockercloud-server-proxy**. This was used to proxy and load balance incoming Docker commands on port 2376 across all manager nodes. It has a network called “dockercloud-server-proxy-network” and secret called “dockercloud-server-proxy-secret”. All of these should be removed.

The following commands show you how to remove the service, network and secret. Be sure to perform these commands against the correct Swarm.

1. If you aren’t already, open an SSH session to a Swarm manager.

2. Remove the service.

```
$ docker service rm dockercloud-server-proxy
```

3. Remove the network.

```
$ docker network rm dockercloud-server-proxy-network
```

4. Remove the secret.

```
$ docker secret rm dockercloud-server-proxy-secret
```

Your Docker Swarm cluster is now deregistered from Docker Cloud and you can manage it independently.


THE END

The following can be tightened up and used if we feel there is a need to show customers how to deploy new public keys to Swarm nodes in the event that customers do not have the keys to access their nodes. However, i don’t think it should be needed as surely all customers will know how to log on to their Swarms - otherwise how would they deploy apps to them?


The easiest way to connect to your Swarm and run a `docker node ls` is via Docker for Mac or Docker for Windows. Click the Docker whale, choose `Swarms` and select the Swarm you want to connect to. This will open a terminal connected to the Swarm.




If you don’t have Docker for Mac or Docker for Windows, you can accomplish the same by clicking the desired Swarm in the Docker Cloud web UI, and running the resulting command in a terminal window.



Selecting a node in the AWS or Azure console's will allow you to get hostname, IP and SSH key information.

### Copy new keys to cluster nodes

You only need to complete this section if you do not have the private key required to SSH to your Swarm nodes.

This section will show you how to run a containerized application that will add a new public key to each node in your Swarm. They key will be added as an extra key to the node's `authorized_keys` file.

The high-level procedure will be as follows:

1.  Connect a Docker client to your Swarm
2.  Create a service stack to copy new keys
3.  Deploy the stack
4.  Test connectivity

You eed a valid key pair to be able to complete this procedure.

1. Connect a Docker client to your Swarm.

  If your are using Docker for Mac or Docker for Windows, click the Swarm from the list available and you will get a terminal configured to send `docker` commands to your Swarm.

  If you're not using Docker for Mac or Docker for Windows, you can click the Swarm name in the Docker Cloud web UI and use the `docker run` command to configure an existing Docker client to talk to the Swarm.

2. Verify that the your client is talking to the correct Swarm.

  One of the easiest ways to do this is to run a `docker node ls`, `docker service ls`, or `docker stack ls` command. If the output returned matches what you expect, then you are connected to the right Swarm.

3. Create the following stackfile in your system's PATH. You will use the file in the next step.

  ```
  version: "3.5"
  services:
    keys:
      image: dockercloud/authorizedkeys:latest
      deploy:
        mode: global
        restart_policy:
          condition: none
      environment:
        - AUTHORIZED_KEYS=ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCFUzmdk+fJLUm7mqzh3yvBhjl//7RNoNmTAL4QyRyZZcctRNEWvNxkbe1jfbsnd9hks/4kfmdTxP3czacdpHBeb9H40AYVo7SoC5qQ4ZraFc58ni1LkaqYpeYME1pRxjuEBBw2gQsI7qrpt5HHr1QWGsVKWiPBdXFz1ciG5unupVSP7bgFo5Qkwd8ymBT8RFn/bzRUxHW8U79m62cQMUsC5FrGWzKlVItNnj/zG7nvVGIdkBtlq9KPuKaijY0PTPXpA/9OY+MWOz4V9GI/UCHiHHHllosD/TcWIIioUq8hp6sGpg5d9ioOIGThejoceRTgZCd/NXLNwu1SFkrVMz imported-openssh-key
      volumes:
        - type: bind
          source: /home/docker
          target: /user
```

  A few notes about the file.

  - It creates a global service on the cluster. Global services run on every node, ensuring the new key gets copied to the `authorized_keys` file on every node.
  - This file assumes your Swarm nodes are running a version of *Moby Linux* that uses the **docker** user account. If you are using a different Linux distro, such as Ubuntu, you may need to change the `services.keys.volumes.source` value. For example, if your Swarm nodes are running on AWS using an Ubuntu image and the **ubuntu** user account, you would change the value to `services.keys.volumes.source=/home/ubuntu`.
  - You will need to substitute the value of **your** public key in the `services.keys.environment` field. Just replace the text after "AUTHORIZED_KEYS".

4. Deploy the stack.

  The following command assumes you created the stackfile above as a file called "keys.yml" and it's in your system's PATH.

  ```
  $ docker stack deploy keys.yml keys
  Creating network keys_default
  Creating service keys_keys
  ```

5. Verify that the stack deployed and added the keys.

  The `docker service logs` command for the stack/service should show two lines of output per service replica. The first line will show "Found authorized keys" and the second will show "Adding public key to ..."

```
$ docker service logs keys_keys
keys_keys.0.txaev739chr9@ip-172-31-6-250.eu-west-1.compute.internal    | => Found authorized keys
keys_keys.0.txaev739chr9@ip-172-31-6-250.eu-west-1.compute.internal    | => Adding public key to .ssh/authorized_keys: ssh-rsa AAAAB...FkrVMz imported-openssh-key
keys_keys.0.hpyx5dd87wm0@ip-172-31-32-56.eu-west-1.compute.internal    | => Found authorized keys
keys_keys.0.hpyx5dd87wm0@ip-172-31-32-56.eu-west-1.compute.internal    | => Adding public key to .ssh/authorized_keys: ssh-rsa AAAAB...FkrVMz imported-openssh-key
keys_keys.0.yzr2d2s2k73w@ip-172-31-24-225.eu-west-1.compute.internal    | => Found authorized keys
keys_keys.0.yzr2d2s2k73w@ip-172-31-24-225.eu-west-1.compute.internal    | => Adding public key to .ssh/authorized_keys: ssh-rsa AAAAB...FkrVMz imported-openssh-key
keys_keys.0.2m8hgpfd9hos@ip-172-31-44-247.eu-west-1.compute.internal    | => Found authorized keys
keys_keys.0.2m8hgpfd9hos@ip-172-31-44-247.eu-west-1.compute.internal    | => Adding public key to .ssh/authorized_keys: ssh-rsa AAAAB...FkrVMz imported-openssh-key
```

6.  Delete the "keys" stack.

    Once the stack has added the new keys it can be deleted. Remember to substitute the name of your stack in the example below.

    ```
    $ docker stack rm keys
    ```

7.  Test logging on with the new private key.

    Use your favourite SSH tool to connect to each node with the private key associated with the new public key you just added.

If you can login with the new private key, you can proceed to the next step.
