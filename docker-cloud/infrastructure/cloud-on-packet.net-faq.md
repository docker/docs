---
description: Docker Cloud and Packet.net
keywords: Packet.net, Cloud, drives
redirect_from:
- /docker-cloud/faq/cloud-on-packet.net-faq/
title: Use Docker Cloud and Packet.net
---

This page answers frequently asked questions about using Docker Cloud with Packet.net.

## What does Docker Cloud create in my Packet.net account?

Docker Cloud creates a project named "Docker Cloud" which contains all the devices that Docker Cloud deploys, no matter what type of device you chose.

Device storage is organized as follows:

- Type 1 devices have a RAID 1 of two SSD drives mounted in `/`.
- Type 3 devices also have a RAID 1 of two SSD drives mounted in `/`, and also offer two NVMe drives without being mounted. Docker Cloud mounts a RAID 1 in `/var/lib/docker`.

An SSH keypair named `dockercloud-<uuid>` is created if no key is found in your account.

## How long does it take to deploy a Packet.net device?

Docker Cloud deploys Ubuntu 14.04 LTS images on both types. Type 1 takes between
5 and 10 minutes to initialize, while type 3 can take up to 15 minutes. The Packet.net engineering team is working to reduce these deployment times.

## What happens if I restart a node in the Packet.net portal?

After the node boots up, the Docker Cloud Agent contacts Docker Cloud using the
API and registers itself with its new IP. Cloud then automatically updates the
DNS of the node and the containers on it to use the new IP. The node changes
state from `Unreachable` to `Deployed`.

## Can I terminate a node from the Packet.net portal?

If you create a node using Docker Cloud but terminate it from the Packet.net
portal, all data in the node is destroyed. Docker Cloud detects the termination
and marks the node as `Terminated`.

If you turn off the device, Docker Cloud marks it as `Unreachable` because the
node has not been terminated, but Cloud cannot contact it.

If you created the host yourself, added it to Docker Cloud as a "Bring Your Own
Node" and then terminated it, the node is marked as `Unreachable` until you
manually remove it.

## How can I log in to a Packet.net node managed by Docker Cloud?

Packet.net copies SSH keys into the created device. This means you can upload your own SSH public key to Packet.net's portal and then SSH into the node using the `root` user. You can also log in to the node from Packet's console, or use a container to copy your SSH keys into the node, as explained in [Sshing into a node](../infrastructure/ssh-into-a-node.md).

## Packet has returned an error, what can I do?

Here is a list of known errors thrown by Packet.net API:

- **You have reached the maximum number of projects you can create (number)**. Please contact `help@packet.net` -> Packet.net limits the number of projects that an account can create. Delete projects in the account or contact [Packet.net](https://www.packet.net/) support to increase the limit.
- **There is an error with your Packet.net account**. Please contact `help@packet.net` -> There is something else wrong with your Packet.net account. Contact [Packet.net](https://www.packet.net/) for more details.
