---
description: Docker Cloud on AWS
keywords: Cloud, AWS, faq
redirect_from:
- /docker-cloud/faq/cloud-on-aws-faq/
title: Use Docker Cloud on AWS
---

This section answers frequently asked questions about using Docker Cloud with
Amazon Web Services (AWS).

## I can't get my account to link to Docker Cloud. How do I troubleshoot it?

To validate your AWS Security Credentials, Docker Cloud tries to dry-run an
instance on every region. Credentials are marked as valid if the operation
succeeds at least in one of the regions. If you get the following message
`Invalid AWS credentials or insufficient EC2 permissions` follow these steps to
troubleshoot it:

1. [Download AWS CLI](https://aws.amazon.com/cli/){: target="_blank" class="_"}.
2. [Configure the CLI](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html){: target="_blank" class="_"} with your security credentials.
2.  Run the following command:

    ```
    aws ec2 run-instances --dry-run --image-id ami-4d883350 --instance-type m3.medium
    ```

This tries to dry-run an Ubuntu 14.04 LTS 64-bit in `sa-east-1` (Sao Paulo,
South America). You can look for the AMI in the region you want to deploy to
[here](http://cloud-images.ubuntu.com/locator/ec2/){: target="_blank"
class="_"}. It should show you the error message. If your configuration is
correct, you see the following message:

```
A client error (DryRunOperation) occurred when calling the RunInstances operation: Request would have succeeded, but DryRun flag is set.
```

## "AWS returned an error: unauthorized operation" using instance profiles to deploy node clusters

This error occurs when you are using an instance profile that has more
permissions than the IAM user you are using with Docker Cloud. You can fix this
by adding the `"Action":"iam:PassRole"` permission to the IAM policy for the
`dockercloud` service user. You can read more about this
[here](http://blogs.aws.amazon.com/security/post/Tx3M0IFB5XBOCQX/Granting-Permission-to-Launch-EC2-Instances-with-IAM-Roles-PassRole-Permission){:
target="_blank" class="_"}.

## What objects does Docker Cloud create in my EC2 account?

If you decide to let Docker Cloud create elements for you, it creates:

- A VPC with the tag name `dc-vpc` and CIDR range `10.78.0.0/16`.
- A set of subnets if there are no subnets already created in the VPC. Docker Cloud creates a subnet in every Availability Zone (AZ) possible, and leaves enough CIDR space for the user to create customized subnets. Every subnet created is tagged with `dc-subnet`.
- An internet gateway named `dc-gateway` attached to the VPC.
- A route table named `dc-route-table` in the VPC, associating the subnet with the gateway.

## How can I customize VPC/IAM elements in Docker Cloud through the AWS dashboard?

Users with AWS EC2-VPC accounts can customize any of the elements explained
above through the AWS API or the dashboard.

In the launch node cluster view, you can choose:

- VPC dropdown:
    1. `Auto` - Delegates creation of the VPC to Docker Cloud.
    2. `vpc-XXXX (dc-vpc)` - Docker Cloud's default VPC. This only appears if you have already deployed nodes to that region. You can choose subnets and security groups with the VPC. See "Which objects does Docker Cloud create in my EC2 account" for detailed info.
    3. `vpc-XXXX` - You can select one of the VPCs already created by you. If you tag name them, it is displayed too.
- Subnets dropdown:
    1. `Auto` - Delegates the management of the subnets to Docker Cloud. Creates them if they do not exist or uses the ones tagged with `dc-subnet`.
    2. Multiple selection of existing subnets. See `How does Docker Cloud balance my nodes among different availability zones?` section for detailed info.
- Security groups dropdown:
    1. `Auto`
    2. Multiple selection of existing security groups.
- IAM roles dropdown:
    1. `None` - Docker Cloud does not apply any instance profiles to the node.
    2. `my_instance_role_name` - You can select one of the IAM roles already created by you.

## How do I customize VPC/IAM elements in Docker Cloud using the API?

Add the following section to your body parameters:

```json
"provider_options" = {
    "vpc": {                                                 # optional
        "id": "vpc-xxxxxxxx",                                # required
        "subnets": ["subnet-xxxxxxxx", "subnet-yyyyyyyy"],   # optional
        "security_groups": ["sg-xxxxxxxx"]                   # optional
    },
    "iam": {                                                 # optional
        "instance_profile_name": "my_instance_profile_name"  # required
    }
}
```

## How does Docker Cloud balance my nodes among different availability zones? (high availability schema)

By default, Docker Cloud tries to deploy your node cluster using a high
availability strategy. To do this, it places every instance one by one in the
less populated availability zone for that node cluster. We can see this behavior
with some examples:

### We allow Docker Cloud to manage VPCs and subnets

Docker Cloud can take over VPC and subnet management for you when you deploy a
node cluster.

For example, assume this is the first time you're deploying a node cluster. You
delegate deployment management to Docker Cloud in the Sao Paulo (South America,
`sa-east-1`) region. You don't send any `provider_options` using the API, and
you leave the VPC, subnet, security groups and IAM role values set to their
defaults on the dashboard. In this situation:

1. Docker Cloud looks for a VPC called `dc-vpc`. The VPC does not exist on the first try, so Docker Cloud creates it and a `dc-gateway`, which attaches to the VPC.
2. Docker Cloud retrieves all subnets in the VPC. No subnets exist on the first try.
3. Docker Cloud creates the subnet.
4. For every availability zone (AZ), Docker Cloud splits the VPC CIDR IP space in (# of AZs + 1) blocks and tries to create (# of AZs) subnets. Remember, we left space for custom subnets.
5. For every subnet, Docker Cloud tries to dry-run an instance of the selected type and creates it if the operation succeeds, creating and associating a `dc-route-table` to the subnet.
6. Once all subnets have been created, Docker Cloud deploys every node of the cluster using a round-robin pattern.

> **Note** If the `dry-run` fails on any of the availability zones, you may see fewer subnets than were originally specified by the number of zones.

### Scaling a node cluster

Following the example in the previous section, you have a node cluster deployed and want to scale it up. Docker Cloud:

1. Looks for `dc-vpc`. Found!
2. Looks for `dc-subnet`s. Found!
3. Counts the nodes in every subnet.
4. Chooses the less populated subnet and deploys the next node there.
4. Repeats until all nodes are deployed.

### We choose where to deploy

What if you have another VPC for some other purpose, (the components already exist) and you want to deploy a node cluster in that VPC.

Docker Cloud:

1. Looks for the selected VPC. Found!
2. Looks for selected subnets. If you do not select any subnets, Docker Cloud tries to create them using the rules previously described.
3. If you selected more than one subnet, Docker Cloud distributes the nodes in the cluster among those subnets. If not, all nodes are placed in the same subnet.

## What happens if I restart a node in the AWS console?

After the node boots up, the Docker Cloud Agent tries to contact the Cloud API
and register itself with its new IP. Once it registers, Docker Cloud
automatically updates the DNS of the node and the containers on it to use the
new IP. The node's state changes from `Unreachable` to `Deployed`.

## Can I use an elastic IP for my nodes?

Yes. However, you must restart the Docker Cloud Agent (or the host) for the
changes to take effect in Docker Cloud.

## What happens when I terminate a node from the AWS console?

If you created the node using Docker Cloud, but you terminate it in the AWS
console, all data in that node is destroyed as the volume attached to it is set
to destroy on node termination. As long as the Docker Cloud IAM user still has
access, Cloud detects the termination and marks the node as `Terminated`.

If you created the host yourself, added it to Docker Cloud as a "Bring Your Own
Node" and then terminated it, the node stays in the `Unreachable` state until
you manually remove it.

## How do I SSH into a node?

Use the instructions [here](ssh-into-a-node.md) to access your nodes over SSH.
If you chose a custom security group, remember to open port 22.

## How do I back up my Docker container volumes to AWS S3?

Use the [dockercloud/dockup](https://hub.docker.com/r/dockercloud/dockup/){:
target="_blank" class="_"} utility image to back up your volumes. You only need
to run it taking the volumes of the container you want to back up with
`volumes-from` and pass it the environment configuration of the container. You
can find more information in its Github repository.
