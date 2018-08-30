---
description: How to link Docker Cloud to a hosted cloud services provider or your own hosts
keywords: node, create, understand
redirect_from:
- /docker-cloud/getting-started/use-hosted/
title: Link to your infrastructure
---

To deploy Docker Cloud nodes, you first need to grant Docker Cloud access to your infrastructure.

This could mean granting access to a cloud services provider such as AWS or Azure, or installing the Docker Cloud Agent on your own hosts. Once this is done, you can provision nodes directly from within Docker Cloud using the Web UI, CLI, or API.

## Link to a cloud service provider
To link your cloud provider accounts, first go to your [Docker Cloud dashboard](https://cloud.docker.com/).

Then, use one of the detailed tutorials below to link your account. You should open the detailed linking tutorial in a new tab or window so you can continue the tutorial when you're finished.

  - [Amazon Web Services](../infrastructure/link-aws.md) (uses an Access Key ID + Secret Access Key)
  - [DigitalOcean](../infrastructure/link-do.md) (uses OAuth)
  - [Microsoft Azure](../infrastructure/link-azure.md) (uses OAuth)
  - [IBM SoftLayer](../infrastructure/link-softlayer.md) (uses an API key)
  - [Packet.net](../infrastructure/link-packet.md) (uses an API key)

  You can always come back and link more cloud service providers later.

## Link to your own hosts (Bring Your Own Node - BYON)

If you are not using a cloud services provider but using your own hosts, install the Docker Cloud Agent on those hosts so that Docker Cloud can communicate with them. Follow the directions at [Bring Your Own Node instructions](../infrastructure/byoh.md). Open these instructions in a new window or tab so you can return to this tutorial once you're done linking your hosts.

## Ready to go?
Once you've linked to your cloud services provider or to your own hosts, [continue the tutorial and deploy your first node](your_first_node.md).
