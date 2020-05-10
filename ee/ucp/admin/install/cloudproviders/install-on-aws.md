---
title: Install UCP on AWS
description: Learn how to install Docker Universal Control Plane in an Amazon Web Services (AWS) environment.
keywords: Universal Control Plane, UCP, install, Docker EE, AWS, Kubernetes
---

>{% include enterprise_label_shortform.md %}

Docker Universal Control Plane (UCP) can be installed on top of AWS without any
customisation following the UCP [install documentation](../index.md). Therefore
this document is **optional**, however if you are deploying Kubernetes
workloads with UCP and want to leverage the [AWS kubernetes cloud
provider](https://github.com/kubernetes/cloud-provider-aws), which provides
dynamic volume and loadbalancer provisioning then you should follow this guide.
This guide is not required if you are only deploying swarm workloads.

The requirements for installing UCP on AWS are included in the following sections:

## Instances
### Hostnames 
The instance's host name must be named `ip-<private ip>.<region>.compute.internal`. For example:
`ip-172-31-15-241.us-east-2.compute.internal`

### Instance tags
The instance must be tagged with ` kubernetes.io/cluster/<UniqueID for Cluster>` and given a 
value of `owned` or `shared`. If the resources created by the cluster is considered owned and 
managed by the cluster, the value should be owned.  If the resources can be shared between multiple 
clusters, it should be tagged as shared.

`kubernetes.io/cluster/1729543642a6` `owned`

### Instance profile for managers
Manager nodes must have an instance profile with appropriate policies attached to enable 
introspection and provisioning of resources. The following example is very permissive:

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [ "ec2:*" ],
      "Resource": [ "*" ]
    },
    {
      "Effect": "Allow",
      "Action": [ "elasticloadbalancing:*" ],
      "Resource": [ "*" ]
    },
    {
      "Effect": "Allow",
      "Action": [ "route53:*" ],
      "Resource": [ "*" ]
    },
    {
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": [ "arn:aws:s3:::kubernetes-*" ]
    }
  ]
}
```

### Instance profile for workers
Worker nodes must have an instance profile with appropriate policies attached to enable access to 
dynamically provisioned resources. The following example is very permissive:

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": [ "arn:aws:s3:::kubernetes-*" ]
    },
    {
      "Effect": "Allow",
      "Action": "ec2:Describe*",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:AttachVolume",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:DetachVolume",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [ "route53:*" ],
      "Resource": [ "*" ]
    }
} 
```

## VPC
### VPC tags
The VPC must be tagged with ` kubernetes.io/cluster/<UniqueID for Cluster>` and given a value 
of `owned` or `shared`.  If the resources created by the cluster is considered owned and managed by the cluster, 
the value should be owned.  If the resources can be shared between multiple clusters, it should be tagged shared.

`kubernetes.io/cluster/1729543642a6` `owned`

### Subnet tags
Subnets  must be tagged with ` kubernetes.io/cluster/<UniqueID for Cluster>` and given a value of `owned` or `shared`. If the resources created by the cluster is considered owned and managed by the cluster, the value should be owned.  If the resources may be shared between multiple clusters, it should be tagged shared.  For example:

`kubernetes.io/cluster/1729543642a6` `owned`

## UCP 
### UCP install 

Once all pre-requisities have been met, run the following command to install
UCP on a manager node. The `--host-address` flag maps to the private IP address
of the master node.

```bash
$ docker container run --rm -it \
  --name ucp \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <ucp-ip> \
  --cloud-provider aws \
  --interactive
```

The full Docker UCP install CLI reference can be found
[here](/reference/ucp/3.2/cli/install/)
