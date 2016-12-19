---
description: Frequently asked questions
keywords: azure faqs
title: Docker for Azure Frequently asked questions (FAQ)
---

## How long will it take before I get accepted into the Docker for Azure private beta?

Docker for Azure is built on top of Docker 1.13, but as with all Beta, things are still changing, which means things can break between release candidates.

We are currently rolling it out slowly to make sure everything is working as it should. This is to ensure that if there are any issues we limit the number of people that are affected.

## Why do you need my Azure Subscription ID?

We are using a private custom VHD, and in order to give you access to this VHD, we need your Azure Subscription ID.

## How do I find my Azure Subscription ID?

You can find this information your Azure Portal Subscription. For more info, look at the directions on [this page](../index.md).

## I use more than one Azure Subscription ID, how do I get access to all of them.

Use the beta sign up form, and put the subscription ID that you need to use most there. Then email us <docker-for-iaas@docker.com> with your information and your other Azure Subscription ID, and we will do our best to add those accounts as well. But due to the large amount of requests, it might take a while before those subscriptions to get added, so be sure to include the important one in the sign up form, so at least you will have that one.

## Can I use my own VHD?
No, at this time we only support the default Docker for Azure VHD.

## Can I specify the type of Storage Account I use for my VM instances?

Not at this time, but it is on our roadmap for future releases.

## Which Azure regions will Docker for Azure work with.

Docker for Azure should work with all supported Azure Marketplace regions.

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

The beta versions of Docker for AWS and Azure send anonymized analytics to Docker. These analytics are used to monitor beta adoption and are critical to improve Docker for AWS and Azure.

## How to run administrative commands?

By default when you SSH into a manager, you will be logged in as the regular username: `docker` - It is possible however to run commands with elevated privileges by using `sudo`.
For example to ping one of the nodes, after finding its IP via the Azure/AWS portal (e.g. 10.0.0.4), you could run:
```
$ sudo ping 10.0.0.4
``` 

Note that access to Docker for AWS and Azure happens through a shell container that itself runs on Docker.