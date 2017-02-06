---
description: Frequently asked questions
keywords: azure faqs
title: Docker for Azure Frequently asked questions (FAQ)
---

## FAQ

## Can I use my own VHD?
No, at this time we only support the default Docker for Azure VHD.

## Can I specify the type of Storage Account I use for my VM instances?

Not at this time, but it is on our roadmap for future releases.

## Which Azure regions will Docker for Azure work with?

Docker for Azure should work with all supported Azure Marketplace regions.

## Where are my container logs?

All container logs are aggregated within the `xxxxlog` storage account.

## Where do I report problems or bugs?

Send an email to <docker-for-iaas@docker.com> or post to the [Docker for Azure](https://github.com/docker/for-azure) GitHub repositories.

In Azure, if your resource group is misbehaving, please run the following diagnostic tool from one of the managers - this will collect your docker logs and send them to Docker:

```bash
$ docker-diagnose
OK hostname=manager1
OK hostname=worker1
OK hostname=worker2
Done requesting diagnostics.
Your diagnostics session ID is 1234567890-xxxxxxxxxxxxxx
Please provide this session ID to the maintainer debugging your issue.
```

> **Note**: Your output will be slightly different from the above, depending on your swarm configuration.

## Analytics

Docker for Azure sends anonymized minimal analytics to Docker (heartbeat). These analytics are used to monitor adoption and are critical to improve Docker for Azure.

## How do I run administrative commands?

By default when you SSH into a manager, you will be logged in as the regular username: `docker` - It is possible however to run commands with elevated privileges by using `sudo`.
For example to ping one of the nodes, after finding its IP via the Azure/Azure portal (e.g. 10.0.0.4), you could run:

```bash
$ sudo ping 10.0.0.4
```

> **Note**: Access to Docker for Azure and Azure happens through a shell container that itself runs on Docker.


## What are the different Azure Regions?
All regions can be found here: [Microsoft Azure Regions](https://azure.microsoft.com/en-us/regions/)
An exerpt of the above regions to use when you create your service principal are:

```none
australiaeast
australiasoutheast
brazilsouth
canadacentral
canadaeast
centralindia
centralus
eastasia
eastus
eastus2
japaneast
japanwest
northcentralus
northeurope
southcentralus
southeastasia
southindia
uksouth
ukwest
westcentralus
westeurope
westindia
westus
westus2
```
