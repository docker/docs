---
description: Frequently asked questions
keywords: aws faqs
title: Docker for AWS frequently asked questions (FAQ)
toc_max: 2
---

## Stable and edge channels

Two different download channels are available for Docker for AWS:

* The **stable channel** provides a general availability release-ready deployment for a fully baked and
tested, more reliable cluster. The stable version of Docker for AWS comes with
the latest released version of Docker Engine. The release schedule is synched
with Docker Engine releases and hotfixes. On the stable channel, you can select
whether to send usage statistics and other data.

* The **edge channel** provides a deployment with new features we are
working on, but is not necessarily fully tested. It comes with the
experimental version of Docker Engine. Bugs, crashes, and issues are
more likely to occur with the edge cluster, but you get a chance to preview
new functionality, experiment, and provide feedback as the deployment
evolve. Releases are typically more frequent than for stable, often one
or more per month. Usage statistics and crash reports are sent by default.
You do not have the option to disable this on the edge  schannel.

## Can I use my own AMI?

No, at this time we only support the default Docker for AWS AMI.

## How can I use Docker for AWS with an AWS account in an EC2-Classic region?

If you have an AWS account that was created before **December 4th, 2013** you
have what is known as an **EC2-Classic** account on regions where you have
previously deployed resources. **EC2-Classic** accounts don't have default VPC's
or the associated subnets, etc. This causes a problem when using our
CloudFormation template because we are using the
[Fn:GetAZs](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getavailabilityzones.html)
function they provide to determine which availability zones you have access to.
When used in a region where you have **EC2-Classic**, this function returns
all availability zones for a region, even ones you don't have access to. When
you have an **EC2-VPC** account, it returns only the availability zones you
have access to.

This causes an error like the following:

> "Value (us-east-1a) for parameter availabilityZone is invalid.
Subnets can currently only be created in the following availability
zones: us-east-1d, us-east-1c, us-east-1b, us-east-1e."

If you have an **EC2-Classic** account, and you don't have access to the `a` and
`b` availability zones for that region.

There isn't anything we can do right now to fix this issue. We have contacted
Amazon to provide a solution.

### How to tell if you are in the EC2-Classic region.

[This AWS documentation
page](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-supported-platforms.html)
describes how you can tell if you have EC2-Classic, EC2-VPC or both.

### Possible fixes to the EC2-Classic region issue:
There are a few workarounds that you can try to get Docker for AWS up and running for you.

1. Create your own VPC, then [install Docker for AWS with a pre-existing VPC](/docker-for-aws/index.md#install-with-an-existing-vpc).
2. Use a region that doesn't have **EC2-Classic**. The most common region with this issue is `us-east-1`. So try another region, `us-west-1`, `us-west-2`, or the new `us-east-2`. These regions should be set up with **EC2-VPC** and the issue shouldn't occur.
3. Create an new AWS account, all new accounts are setup using **EC2-VPC** and do not have this problem.
4. Contact AWS support to convert your **EC2-Classic** account to a **EC2-VPC** account. For more information checkout the following answer for **"Q. I really want a default VPC for my existing EC2 account. Is that possible?"** on https://aws.amazon.com/vpc/faqs/#Default_VPCs

### Helpful links:
- http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/default-vpc.html
- http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-supported-platforms.html
- http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-vpc.html
- https://aws.amazon.com/vpc/faqs/#Default_VPCs
- https://aws.amazon.com/blogs/aws/amazon-ec2-update-virtual-private-clouds-for-everyone/


## Can I use my existing VPC?

Yes, see [install Docker for AWS with a pre-existing VPC](/docker-for-aws/index.md#install-with-an-existing-vpc) for more info.

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
If you are using the `10.0.0.0/16` CIDR in your VPC. When you create a docker network, you need to pick a subnet (using `docker network create —subnet` option) that doesn't conflict with the `10.0.0.0` network.

## Which AWS regions does this work with?

Docker for AWS should work with all regions except for AWS US Gov Cloud (us-gov-west-1) and AWS China, which are a little different than the other regions.

## How many Availability Zones does Docker for AWS use?

Docker for AWS determines the correct amount of Availability Zone's to use based on the region. In regions that support it, we use 3 Availability Zones, and 2 for the rest of the regions. We recommend running production workloads only in regions that have at least 3 Availability Zones.

## What do I do if I get `KeyPair error` on AWS?
As part of the prerequisites, you need to have an SSH key uploaded to the AWS region you are trying to deploy to.
For more information about adding an SSH key pair to your account, refer to the [Amazon EC2 Key Pairs docs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html).

## Where are my container logs?

All container logs are aggregated within [AWS CloudWatch](https://aws.amazon.com/cloudwatch/).

## Best practice to deploy a large cluster

When deploying a cluster of more than 20 workers, it can take a very long time for AWS to deploy all of the instances (1+hrs).
It is best to deploy a cluster of 20 workers then scale it up in the Auto-Scaling Group (ASG) once it's been deployed.

Benchmark of 3 Managers (m4.large) + 200 workers (t2.medium):

* Deploying (~3.1hrs)
	* Deployment: 3 Managers + 200 workers = ~190mins
* Scaling (~35mins)
	* Deployment: 3 Managers + 20 workers = ~20mins
	* Scaling: 20 workers -> 200 workers via ASG = ~15mins


> **Note**: During a Stack upgrade, you need to match the Auto-Scaling Group worker count, otherwise AWS scales it back down (aka type 200 workers in the input box)


## Where do I report problems or bugs?

Search for existing issues, or create a new one, within the [Docker for AWS](https://github.com/docker/for-aws) GitHub repositories.

In AWS, if your stack is misbehaving, run the following diagnostic tool from one of the managers - this collects your docker logs and send them to Docker:

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

Docker for AWS sends anonymized minimal metrics to Docker (heartbeat). These metrics are used to monitor adoption and are critical to improve Docker for AWS.

## How do I run administrative commands?

By default when you SSH into a manager, you are logged in as the regular username: `docker` - It is possible however to run commands with elevated privileges by using `sudo`.
For example to ping one of the nodes, after finding its IP via the Azure/AWS portal, such as 10.0.0.4, you could run:

```bash
$ sudo ping 10.0.0.4
```

> **Note**: Access to Docker for AWS and Azure happens through a shell container that itself runs on Docker.


## What are the Editions containers running after deployment?

In order for our editions to deploy properly and for load balancer integrations to happen, we run a few containers. They are as follow:

| Container name | Description |
|---|---|
| `init`  | Sets up the swarm and makes sure that the stack came up properly. (checks manager+worker count).|
| `shell` | This is our shell/ssh container. When you SSH into an instance, you're actually in this container.|
| `meta`  | Assist in creating the swarm cluster, giving privileged instances the ability to join the swarm.|
| `l4controller` | Listens for ports exposed at the docker CLI level and opens them in the load balancer. |

## How do I uninstall Docker for AWS?

You can remove the Docker for AWS setup and stacks through the [AWS
Console](https://console.aws.amazon.com/console/home){: target="_blank"
class="_"} on the CloudFormation page. See [Uninstalling or removing a
stack](/docker-for-aws/index.md#uninstalling-or-removing-a-stack).
