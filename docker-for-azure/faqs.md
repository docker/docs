---
description: Frequently asked questions
keywords: azure faqs
title: Docker for Azure frequently asked questions (FAQ)
toc_max: 2
---

## Stable and edge channels

Two different download channels are available for Docker for Azure:

* The **stable channel** provides a general availability release-ready deployment
  for a fully baked and tested, more reliable cluster. The stable version of Docker
  for Azure comes with the latest released version of Docker Engine. The release
  schedule is synched with Docker Engine releases and hotfixes. On the stable
  channel, you can select whether to send usage statistics and other data.

* The **edge channel** provides a deployment with new features we are working on,
  but is not necessarily fully tested. It comes with the experimental version of
  Docker Engine. Bugs, crashes, and issues are more likely to occur with the edge
  cluster, but you get a chance to preview new functionality, experiment, and provide
  feedback as the deployment evolve. Releases are typically more frequent than for
  stable, often one or more per month. Usage statistics and crash reports are sent
  by default. You do not have the option to disable this on the edge channel.

## Can I use my own VHD?
No, at this time we only support the default Docker for Azure VHD.

## Can I specify the type of Storage Account I use for my VM instances?

Not at this time, but it is on our roadmap for future releases.

## Which Azure regions does Docker for Azure work with?

Docker for Azure should work with all supported Azure Marketplace regions.

## Where are my container logs?

All container logs are aggregated within the `xxxxlog` storage account.

## Where do I report problems or bugs?

Search for existing issues, or create a new one, within the [Docker for Azure](https://github.com/docker/for-azure) GitHub repositories.

In Azure, if your resource group is misbehaving, run the following diagnostic tool from one of the managers - this collects your docker logs and sends them to Docker:

```bash
$ docker-diagnose
OK hostname=manager1
OK hostname=worker1
OK hostname=worker2
Done requesting diagnostics.
Your diagnostics session ID is 1234567890-xxxxxxxxxxxxxx
Please provide this session ID to the maintainer debugging your issue.
```

> **Note**: Your output may be slightly different from the above, depending on your swarm configuration.

## Metrics

Docker for Azure sends anonymized minimal metrics to Docker (heartbeat). These metrics are used to monitor adoption and are critical to improve Docker for Azure.

## How do I run administrative commands?

By default when you SSH into a manager, you are logged in as the regular username: `docker` - It is possible however to run commands with elevated privileges by using `sudo`.
For example to ping one of the nodes, after finding its IP via the Azure/Azure portal, such as, 10.0.0.4, you could run:

```bash
$ sudo ping 10.0.0.4
```

> **Note**: Access to Docker for Azure and Azure happens through a shell container that itself runs on Docker.


## What are the Editions containers running after deployment?

In order for our editions to deploy properly and for load balancer integrations to happen, we run a few containers. They are as follow:

| Container name | Description |
|---|---|
| `init`  | Sets up the swarm and makes sure that the stack came up properly. (checks manager+worker count).|
| `agent` | This is our shell/ssh container. When you SSH into an instance, you're actually in this container.|
| `meta`  | Assist in creating the swarm cluster, giving privileged instances the ability to join the swarm.|
| `l4controller` | Listens for ports exposed at the docker CLI level and opens them in the load balancer.|
| `logger` | Our log aggregator. This allows us to send all docker logs to the storage account.|


## What are the different Azure Regions?
All regions can be found here: [Microsoft Azure Regions](https://azure.microsoft.com/en-us/regions/).
An excerpt of the above regions to use when you create your service principal are:

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
koreacentral
koreasouth
northcentralus
northeurope
southcentralus
southeastasia
southindia
uksouth
ukwest
usgovvirginia
usgoviowa
westcentralus
westeurope
westindia
westus
westus2
```
