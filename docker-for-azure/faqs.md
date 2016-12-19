---
description: Frequently asked questions
keywords: azure faqs
title: Docker for Azure Frequently asked questions (FAQ)
---

## Can I use my own VHD?
No, at this time we only support the default Docker for Azure VHD.

## Can I specify the type of Storage Account I use for my VM instances?

Not at this time, but it is on our roadmap for future releases.

## Which Azure regions will Docker for Azure work with.

Docker for Azure should work with all supported Azure Marketplace regions.

## Where are my container logs?

All container logs are aggregated within the `xxxxlog` storage account.

## I have a problem/bug where do I report it?

Send an email to <docker-for-iaas@docker.com> or post to the [Docker for Azure](https://github.com/docker/for-azure) GitHub repositories.

In Azure, if your resource group is misbehaving, please run the following diagnostic tool from one of the managers - this will collect your docker logs and send them to Docker:

```
$ docker-diagnose
OK hostname=manager1
OK hostname=worker1
OK hostname=worker2
Done requesting diagnostics.
Your diagnostics session ID is 1234567890-xxxxxxxxxxxxxx
Please provide this session ID to the maintainer debugging your issue.
```

_Please note that your output will be slightly different from the above, depending on your swarm configuration_

## Analytics

Docker for Azure sends anonymized minimal analytics to Docker (heartbeat). These analytics are used to monitor adoption and are critical to improve Docker for Azure.

## How to run administrative commands?

By default when you SSH into a manager, you will be logged in as the regular username: `docker` - It is possible however to run commands with elevated privileges by using `sudo`.
For example to ping one of the nodes, after finding its IP via the Azure/AWS portal (e.g. 10.0.0.4), you could run:
```
$ sudo ping 10.0.0.4
```

Note that access to Docker for AWS and Azure happens through a shell container that itself runs on Docker.
