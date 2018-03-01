---
description: Migrating from Docker Cloud
keywords: cloud, swarm, migration
title: Deregister swarms from Docker Cloud
---

## De-registering Swarms from Docker Cloud

A Docker Cloud Swarm runs on either AWS or Azure. Docker Cloud just adds integrations with other Docker tools and services that make it easier to connect and manage.

This page explains how you can de-register a Swarm cluster from Docker Cloud so that it can be managed independently. We demonstrate the process using a Swarm cluster managed by Docker Cloud that is running on AWS. The overall process is the same if your Swarm is running on Azure.

There is no application migration or reconfiguration required as part of this procedure. The only thing that will change will be that your Swarm cluster will no longer be integrated with Docker Cloud, Docker for Mac/Windows, and other Docker services.

## Pre-requisites

To complete this procedure you will need:

- An account to log on to AWS or Azure and view instances
- An SSH key pair

## High-level overview

De-registering your Swarm cluster from Docker Cloud involves the following high-level steps:

1. Identify your Swarm nodes
2. Identify key pair used to log-on to nodes
3. Copy new keys to cluster nodes
4. De-register from the Docker Cloud UI
5. Delete the dockercloud-server-proxy service

### Identify your Swarm nodes

You do not need to complete this step if your already know the DNS names and IP addresses of the nodes in your Swarm.

Your Docker Cloud Swarm runs on either AWS or Azure, and it's vital that you know the DNS names or IP addresses of your Swarm nodes before you de-register it. Failure to know this information might result in you not being able to locate and connect to your Swarm after the de-registering process.

The simplest way to find this information is via the native AWS or Azure tools. For example, log on to the AWS console and open the EC2 Dashboard for the region that your Swarm nodes are in. Locate your instances and note their DNS names and IPs.

If you cannot locate your Swarm nodes using native cloud tools, there are a few ways to connect to your Swarm cluster and get hints about cloud instance details. These include:

- Connecting to your Swarm via Docker for Mac and Docker for Windows
- Connecting to your Swarm using the `docker run` command from Docker Cloud
- Inspecting the Swarm endpoint from Docker Cloud

If you are using DfM or DfW, you can click `Swarms` and select the Swarm you plan to de-register. This will open a new terminal window configured to connect to your Swarm. The example below shows how to connect to a Swarm called **prod-equus** using Docker for Windows.



Running a `docker node ls` lists the nodes in the Swarm, including their names. The names in this example tell us that the nodes are in the **eu-west-1** region. You should now be able to locate them in the AWS EC2 dashboard for that region.

```
$ docker node ls
ID                            HOSTNAME                                      STATUS              AVAILABILITY        MANAGER STATUS
926r8aql8inreo3s5029gvalw *   ip-172-31-6-250.eu-west-1.compute.internal    Ready               Active              Leader
zl0g4g54akozhqqsmol6i4058     ip-172-31-24-225.eu-west-1.compute.internal   Ready               Active              Reachable
b1u5549n1j9lhwdqmcaabjya7     ip-172-31-32-56.eu-west-1.compute.internal    Ready               Active              Reachable
sgp9v372hu63w11hfjxfwcsyd     ip-172-31-44-247.eu-west-1.compute.internal   Ready               Active
```

By default, AWS will label your Swarm nodes as <swarm-name>-worker or <swarm-name>-Manager. For example, a Swarm called **prod-equus** in Docker Cloud will have *manager* nodes labelled as **prod-equus-Manager**, and *worker* nodes labelled as **prod-equus-worker** in AWS.

Selecting a node in the AWS or Azure console's will allow you to get hostnameDNS, IP and SSH key information.

If you are not using DfM or DfW, you can click the Swarm in Docker Cloud UI and use the resulting `docker run` command to configure a Docker client and then run the same `docker node ls` to get DNS names and region info.

You should now know the important details of Cloud-managed Swarm --- cloud provider, instance details etc.

### Identify key pair used to connect to nodes

You only need to perform this step if you do not have the private key required to connect to your Swarm nodes.

It's possible to install additional keys on your Swarm nodes before de-registering them from Docker Cloud.

We'll walk through the process of identifying and testing the private key for a Cloud-managed Swarm node running on AWS.

1. Open the AWS console and drill down to the EC2 dashboard for the Region that your Swarm is hosted in.

2. Select `Instances` from the left-hand pane.

  Your Swarm instances should be listed in the main area of the screen.

3. Select one of your Swarm instances.

  Take note of the `Key pair name` value on the `Description` tab in the bottom half of the screen. The example below shows a key pair called **eu-west-1-key**.



If you are in possession of the private key from the key pair, you can open an SSH session and test that you can connect to the instance. You should perform this test for each Swarm *manager* and *worker* node.

If you can successfully SSH to the nodes in your Swarm, you can skip the next section and go straight to de-registering your cluster from Docker Cloud.

### Copy new keys to cluster nodes

You only need to complete this section if you do not have the private key required to SSH to your Swarm nodes.

This section will show you how to run a containerized application that will add a new public key to each node in your Swarm. They key will be added as an extra key to the node's `authorized_keys` file.

The high-level procedure will be as follows:

1. Connect a Docker client to your Swarm
2. Create a service stack to copy new keys
3. Deploy the stack
4. Test connectivity

You will need a valid key pair to be able to complete this procedure.

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

6. Delete the "keys" stack.

  Once the stack has added the new keys it can be deleted. Remember to substitute the name of your stack in the example below.

  ```
  $ docker stack rm keys
  ```

7. Test logging on with the new private key.

  Use your favourite SSH tool to connect to each node with the private key associated with the new public key you just added.

If you can login with the new private key, you can proceed to the next step.

### De-register from the Docker Cloud UI

Only proceed with this section if you know the details of your Swarm nodes (cloud provider, public DNS, public IP etc.) and you have verified that you can SSH to each node with your private key.

1. Open the Docker Cloud web UI and click `Swarms`.

2. Click the three dots to the right of the Swarm you want to de-register and select `Unregister`.

3. You will be prompted to confirm the de-registration process.

The Swarm is now de-registered from the Docker Cloud web UI and will no longer be visible in other products such as Docker for Mac and Docker for Windows.

### Delete the proxy service

The final step to de-register your Swarm with Docker Cloud is to delete the **dockercloud-server-proxy** service. This is a global service constrained to all manager nodes in the Swarm.

Execute the following command from a manager node in your Swarm.

```
$ docker service rm dockercloud-server-proxy
```

Confirm that the service is deleted.

```
$ docker service ls
```

The **dockercloud-server-proxy** service should no longer appear in the list.

## Summary

Congratulations, your Docker Swarm cluster is now de-registered from Docker Cloud and you can manage it independently.
