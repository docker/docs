---
description: Frequently asked questions
keywords: aws faqs
title: Docker for AWS Frequently asked questions (FAQ)
toc_max: 2
---

## Stable and beta channels

Two different download channels are available for Docker for AWS:

* The **stable channel** provides a general availability release-ready deployment
  for a fully baked and tested, more reliable cluster. The stable version of Docker
  for AWS comes with the latest released version of Docker Engine. The release
  schedule is synched with Docker Engine releases and hotfixes. On the stable
  channel, you can select whether to send usage statistics and other data.

* The **beta channel** provides a deployment with new features we are working on,
  but is not necessarily fully tested. It comes with the experimental version of
  Docker Engine. Bugs, crashes and issues are more likely to occur with the beta
  cluster, but you get a chance to preview new functionality, experiment, and provide
  feedback as the deployment evolve. Releases are typically more frequent than for
  stable, often one or more per month. Usage statistics and crash reports are sent
  by default. You do not have the option to disable this on the beta channel.

## Can I use my own AMI?

No, at this time we only support the default Docker for AWS AMI.

## How can I use Docker for AWS with an AWS account in an EC2-Classic region?

If you have an AWS account that was created before **December 4th, 2013** you have what is known as an **EC2-Classic** account on regions where you have previously deployed resources. **EC2-Classic** accounts don't have default VPC's or the associated subnets, etc. This causes a problem when using our CloudFormation template  because we are using the [Fn:GetAZs](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getavailabilityzones.html) function they provide to determine which availability zones you have access too. When used in a region where you have **EC2-Classic**, this function will return all availability zones for a region, even ones you don't have access too. When you have an **EC2-VPC** account, it will return only the availability zones you have access to.

This will cause an error like the following:

> "Value (us-east-1a) for parameter availabilityZone is invalid. Subnets can currently only be created in the following availability zones: us-east-1d, us-east-1c, us-east-1b, us-east-1e."

If you have an **EC2-Classic** account, and you don't have access to the `a` and `b` availability zones for that region.

There isn't anything we can do right now to fix this issue, we have contacted Amazon, and we are hoping they will be able to provide us with a way to determine if an account is either **EC2-Classic** or **EC2-VPC**, so we can act accordingly.

### How to tell if you are in the EC2-Classic region.

This AWS documentation page will describe how you can tell if you have EC2-Classic, EC2-VPC or both.  http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-supported-platforms.html

### Possible fixes to the EC2-Classic region issue:
There are a few work arounds that you can try to get Docker for AWS up and running for you.

1. Create your own VPC, then [install Docker for AWS with a pre-existing VPC](index.md#install-with-an-existing-vpc).
2. Use a region that doesn't have **EC2-Classic**. The most common region with this issue is `us-east-1`. So try another region, `us-west-1`, `us-west-2`, or the new `us-east-2`. These regions will more then likely be setup with **EC2-VPC** and you will not longer have this issue.
3. Create an new AWS account, all new accounts will be setup using **EC2-VPC** and will not have this problem.
4. Contact AWS support to convert your **EC2-Classic** account to a **EC2-VPC** account. For more information checkout the following answer for **"Q. I really want a default VPC for my existing EC2 account. Is that possible?"** on https://aws.amazon.com/vpc/faqs/#Default_VPCs

### Helpful links:
- http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/default-vpc.html
- http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-supported-platforms.html
- http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-vpc.html
- https://aws.amazon.com/vpc/faqs/#Default_VPCs
- https://aws.amazon.com/blogs/aws/amazon-ec2-update-virtual-private-clouds-for-everyone/


## Can I use my existing VPC?

Yes, see [install Docker for AWS with a pre-existing VPC](index.md#install-with-an-existing-vpc) for more info.

## Recommended VPC and subnet setup

#### VPC

* **CIDR:** 172.31.0.0/16
* **DNS hostnames:** yes
* **DNS resolution:** yes
* **DHCP option set:** DHCP Options (Below)

#### Internet gateway
* **VPC:** VPC (above)

#### DHCP option set

* **domain-name:** ec2.internal
* **domain-name-servers:** AmazonProvidedDNS

#### Subnet1
* **CIDR:** 172.31.16.0/20
* **Auto-assign public IP:** yes
* **Availability-Zone:** A

#### Subnet2
* **CIDR:** 172.31.32.0/20
* **Auto-assign public IP:** yes
* **Availability-Zone:** B

#### Subnet3
* **CIDR:** 172.31.0.0/20
* **Auto-assign public IP:** yes
* **Availability-Zone:** C

#### Route table
* **Destination CIDR block:** 0.0.0.0/0
* **Subnets:** Subnet1, Subnet2, Subnet3

##### Subnet note:
If you are using the `10.0.0.0/16` CIDR in your VPC. When you create a docker network, you will need to make sure you pick a subnet (using `docker network create —subnet` option) that doesn't conflict with the `10.0.0.0` network.

## Which AWS regions will this work with?

Docker for AWS should work with all regions except for AWS China, which is a little different than the other regions.

## How many Availability Zones does Docker for AWS use?

Docker for AWS determines the correct amount of Availability Zone's to use based on the region. In regions that support it, we will use 3 Availability Zones, and 2 for the rest of the regions. We recommend running production workloads only in regions that have at least 3 Availability Zones.

## What do I do if I get `KeyPair error` on AWS?
As part of the prerequisites, you need to have an SSH key uploaded to the AWS region you are trying to deploy to.
For more information about adding an SSH key pair to your account, please refer to the [Amazon EC2 Key Pairs docs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html)

## Where are my container logs?

All container logs are aggregated within [AWS CloudWatch](https://aws.amazon.com/cloudwatch/).

## Where do I report problems or bugs?

Send an email to <docker-for-iaas@docker.com> or post to the [Docker for AWS](https://github.com/docker/for-aws) GitHub repositories.

In AWS, if your stack is misbehaving, please run the following diagnostic tool from one of the managers - this will collect your docker logs and send them to Docker:

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

## Metrics

Docker for AWS sends anonymized minimal metrics to Docker (heartbeat). These metrics are used to monitor adoption and are critical to improve Docker for AWS.

## How do I run administrative commands?

By default when you SSH into a manager, you will be logged in as the regular username: `docker` - It is possible however to run commands with elevated privileges by using `sudo`.
For example to ping one of the nodes, after finding its IP via the Azure/AWS portal (e.g. 10.0.0.4), you could run:

```bash
$ sudo ping 10.0.0.4
```

> **Note**: Access to Docker for AWS and Azure happens through a shell container that itself runs on Docker.
