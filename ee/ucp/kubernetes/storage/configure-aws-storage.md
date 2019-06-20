---
title: Configure AWS EBS Storage for Kubernetes
description: Learn how configure AWS EBS storage for Kubernetes clusters.
keywords: UCP, Docker Enterprise, Kubernetes, storage, AWS, ELB
redirect_from:
- /ee/ucp/kubernetes/configure-aws-storage/
---

[AWS Elastic Block Store](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AmazonEBS.html) (EBS) can be deployed with Kubernetes in Docker Enterprise 2.1 to use AWS volumes as peristent storage for applications. Before using EBS volumes, configure UCP and the AWS infrastructure for storage orchestration to function.

## Configure AWS Infrastructure for Kubernetes

Kubernetes [Cloud Providers](https://kubernetes.io/docs/concepts/cluster-administration/cloud-providers/) provide a method of provisioning cloud resources through Kubernetes via the `--cloud-provider` option. In AWS, this flag allows the [provisioning of EBS volumes](#) and cloud load balancers.

Configuring a cluster for AWS requires several specific configuration parameters in the infrastructure before installing UCP.

### AWS IAM Permissions

Instances must have the following [AWS Identity and Access Management](https://docs.aws.amazon.com/IAM/latest/UserGuide/introduction.html) permissions configured to provision EBS volumes through Kubernetes PVCs.


| Master     | Worker |
|------------|--------|
| ec2:DescribeInstances         | ec2:DescribeInstances      |
| ec2:AttachVolume           	| ec2:AttachVolume           |
| ec2:DetachVolume     			| ec2:DetachVolume           |
| ec2:DescribeVolumes  			| ec2:DescribeVolumes        |
| ec2:CreateVolume				| ec2:DescribeSecurityGroups |
| ec2:DeleteVolume				|                            |
| ec2:CreateTags                |                            |
| ec2:DescribeSecurityGroups   |                            |


### Infrastructure Configuration

- Apply the roles and policies to Kubernetes masters and workers as indicated in the above chart.
- Set the hostname of the EC2 instances to the private DNS hostname of the instance. See [DNS Hostnames](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-dns.html#vpc-dns-hostnames) and [To change the system hostname without a public DNS name](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/set-hostname.html#set-hostname-system) for more details.
- Label the EC2 instances with the key `KubernetesCluster` and assign the same value across all nodes, for example, `UCPKubenertesCluster`.

### Cluster Configuration

- In addition to your existing [install flags](https://docs.docker.com/reference/ucp/3.1/cli/install/) the cloud provider flag `--cloud-provider=aws` is required at install time.
- The cloud provider can also be enabled post-install through the UCP config. The `ucp-agent` needs to be updated to propogate the new config, as described in [UCP configuration file](https://docs.docker.com/ee/ucp/admin/configure/ucp-configuration-file/#inspect-and-modify-existing-configuration).

```
[cluster_config]

...

  cloud_provider = "aws"
```

## Deploy AWS EBS Volumes

After configuring UCP for the AWS cloud provider, you can create persistent volumes that deploy EBS volumes attached to hosts and mounted inside pods. The EBS volumes are provisioned dynamically such they are created, attached, destroyed along with the lifecycle of the persistent volumes. This does not require users to directly access to the AWS as you request these resources directly through Kubernetes primitives.

We recommend you use the `StorageClass` and `PersistentVolumeClaim` resources as these abstraction layers provide more portability as well as control over the storage layer across environments.

To learn more about storage concepts in Kubernetes, see [Storage - Kubernetes](https://kubernetes.io/docs/concepts/storage/).

### Creating a Storage Class

A `StorageClass` lets administrators describe “classes” of storage available in which classes map to quality-of-service levels, or backup policies, or any policies required by cluster administrators. The following `StorageClass` maps a "standard" class of storage to the `gp2` type of storage in AWS EBS.

```
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: standard
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
reclaimPolicy: Retain
mountOptions:
  - debug
```

For descriptions of AWS EBS parameters, see [Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/#aws).

### Creating a Persistent Volume Claim

A `PersistentVolumeClaim` (PVC) is a claim for storage resources that are bound to a `PersistentVolume` (PV) when storage resources are granted. The following PVC makes a request for `1Gi` of storage from the `standard` storage class.

```
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: task-pv-claim
spec:
  storageClassName: standard
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
 ```

### Deploying a Persistent Volume

The following Pod spec references the PVC `task-pv-claim` from above which references the `standard` storage class in this cluster.

```
kind: Pod
apiVersion: v1
metadata:
  name: task-pv-pod
spec:
  volumes:
    - name: task-pv-storage
      persistentVolumeClaim:
       claimName: task-pv-claim
  containers:
    - name: task-pv-container
      image: nginx
      ports:
        - containerPort: 80
          name: "http-server"
      volumeMounts:
        - mountPath: "/usr/share/nginx/html"
          name: task-pv-storage
 ```

### Inspecting and Using PVs

 Once the pod is deployed, run the following `kubectl` command to verify the PV was created and bound to the PVC.

```
kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                   STORAGECLASS   REASON    AGE
pvc-751c006e-a00b-11e8-8007-0242ac110012   1Gi        RWO            Retain           Bound     default/task-pv-claim   standard                 3h
```

The AWS console shows a volume has been provisioned having a matching name with type `gp2` and a `1GiB` size.

![](../images/aws-ebs.png)

## Where to go next

- [Deploy an Ingress Controller on Kubernetes](/ee/ucp/kubernetes/layer-7-routing/)
- [Discover Network Encryption on Kubernetes](/ee/ucp/kubernetes/kubernetes-network-encryption/)