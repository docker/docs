---
previewflag: cloud-swarm
description: How to set up SSH keys
keywords: Cloud, SSH keys, Azure, link
title: Set up SSH keys
---

You can link your Cloud and Service providers so that Docker Cloud can provision and manage swarms on your behalf. For this, you need an SSH key to authenticate Docker to your provider.

## About SSH

{% include content/ssh/ssh-overview.md %}

## Check for existing SSH keys

You may not need to generate a new SSH key if you have an existing key that you
want to reuse.

{% include content/ssh/ssh-find-keys.md %}

If you find an existing key you want to use, skip to the topic that describes
how to [copy your public key for use with Docker
Cloud](#copy-your-public-key-for-use-with-docker-cloud).

Otherwise, [create a new SSH
key](#create-a-new-ssh-key-for-use-by-docker-cloud).

## Create a new SSH key

{% include content/ssh/ssh-gen-keys.md %}

## Add your key to the ssh-agent

{% include content/ssh/ssh-add-keys-to-agent.md %}

## Copy your public key for use with Docker Cloud

You need your SSH public key to provide to Docker Cloud. When you are ready
to add it, you can copy the public key as follows.

{% include content/ssh/ssh-copy-key.md %}

## Related topics

* [Swarms in Docker Cloud](index.md)

* [Link to Docker Cloud to Amazon Web Services](link-aws-swarm.md)

* [Link Docker Cloud to Microsoft Azure Cloud Services](link-azure-swarm.md)

* [Create a new swarm on Microsoft Azure in Docker Cloud](create-cloud-swarm-azure.md)

* [Create a new swarm on AWS in Docker Cloud](create-cloud-swarm-azure.md)
