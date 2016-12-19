---
description: Frequently asked questions
keywords: aws faqs
title: Docker for AWS Frequently asked questions (FAQ)
---

## Can I use my own AMI?

No, at this time we only support the default Docker for AWS AMI.

## How to use Docker for AWS with an AWS account in an EC2-Classic region.

If you have an AWS account that was created before **December 4th, 2013** you have what is known as an **EC2-Classic** account on regions where you have previously deployed resources. **EC2-Classic** accounts don't have default VPC's or the associated subnets, etc. This causes a problem when using our CloudFormation template  because we are using the [Fn:GetAZs](http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getavailabilityzones.html) function they provide to determine which availability zones you have access too. When used in a region where you have **EC2-Classic**, this function will return all availability zones for a region, even ones you don't have access too. When you have an **EC2-VPC** account, it will return only the availability zones you have access to.

This will cause an error like the following:

> "Value (us-east-1a) for parameter availabilityZone is invalid. Subnets can currently only be created in the following availability zones: us-east-1d, us-east-1c, us-east-1b, us-east-1e."

If you have an **EC2-Classic** account, and you don't have access to the `a` and `b` availability zones for that region.

There isn't anything we can do right now to fix this issue, we have contacted Amazon, and we are hoping they will be able to provide us with a way to determine if an account is either **EC2-Classic** or **EC2-VPC**, so we can act accordingly.

#### How to tell if you have this issue.

This AWS documentation page will describe how you can tell if you have EC2-Classic, EC2-VPC or both.  http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-supported-platforms.html

#### How to fix:
There are a few work arounds that you can try to get Docker for AWS up and running for you.

1. Use a region that doesn't have **EC2-Classic**. The most common region with this issue is `us-east-1`. So try another region, `us-west-1`, `us-west-2`, or the new `us-east-2`. These regions will more then likely be setup with **EC2-VPC** and you will not longer have this issue.
2. Create an new AWS account, all new accounts will be setup using **EC2-VPC** and will not have this problem.
3. You can try and contact AWS support to convert your **EC2-Classic** account to a **EC2-VPC** account. For more information checkout the following answer for **"Q. I really want a default VPC for my existing EC2 account. Is that possible?"** on https://aws.amazon.com/vpc/faqs/#Default_VPCs

#### Helpful links:
- http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/default-vpc.html
- http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-supported-platforms.html
- http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-vpc.html
- https://aws.amazon.com/vpc/faqs/#Default_VPCs
- https://aws.amazon.com/blogs/aws/amazon-ec2-update-virtual-private-clouds-for-everyone/


## Can I use my existing VPC?

Not at this time, but it is on our roadmap for future releases.

## Which AWS regions will this work with.

Docker for AWS should work with all regions except for AWS China, which is a little different than the other regions.

## How many Availability Zones does Docker for AWS use?

All of Amazons regions have at least 2 AZ's, and some have more. To make sure Docker for AWS works in all regions, only 2 AZ's are used even if more are available.

## What do I do if I get "KeyPair error" on AWS?
As part of the prerequisites, you need to have an SSH key uploaded to the AWS region you are trying to deploy to.
For more information about adding an SSH key pair to your account, please refer to the [Amazon EC2 Key Pairs docs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html)

## Where are my container logs?

All container logs are aggregated within [AWS CloudWatch](https://aws.amazon.com/cloudwatch/).

## I have a problem/bug where do I report it?

Send an email to <docker-for-iaas@docker.com> or post to the [Docker for AWS](https://github.com/docker/for-aws) GitHub repositories.

In AWS, if your stack is misbehaving, please run the following diagnostic tool from one of the managers - this will collect your docker logs and send them to Docker:

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

Docker for AWS sends anonymized minimal analytics to Docker (heartbeat). These analytics are used to monitor adoption and are critical to improve Docker for AWS.

## How to run administrative commands?

By default when you SSH into a manager, you will be logged in as the regular username: `docker` - It is possible however to run commands with elevated privileges by using `sudo`.
For example to ping one of the nodes, after finding its IP via the Azure/AWS portal (e.g. 10.0.0.4), you could run:
```
$ sudo ping 10.0.0.4
```

Note that access to Docker for AWS and Azure happens through a shell container that itself runs on Docker.
