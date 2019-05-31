---
title: Configuring Azure Disk Storage for Kubernetes
description: Learn how to add persistent storage to your Docker Enterprise clusters running on Azure with Azure Disk. 
keywords: Universal Control Plane, UCP, Docker EE, Kubernetes, storage, volume
redirect_from:
---

Platform operators can provide persistent storage for workloads running on
Docker Enterprise and Microsoft Azure by using Azure Disk. Platform
operators can either pre-provision Azure Disks to be consumed by Kubernetes
Pods, or can use the Azure Kubernetes integration to dynamically provision Azure
Disks on demand.


## Prerequisites

This guide assumes you have already provisioned a UCP environment on 
Microsoft Azure. The Cluster must be provisioned after meeting all of the
prerequisites listed in [Install UCP on
Azure](/ee/ucp/admin/install/install-on-azure.md).

Additionally, this guide uses the Kubernetes Command Line tool `$
kubectl` to provision Kubernetes objects within a UCP cluster. Therefore, this
tool must be downloaded, along with a UCP client bundle. For more
information on configuring CLI access for UCP, see [CLI Based
Access](/ee/ucp/user-access/cli/).

## Manually provision Azure Disks

An operator can use existing Azure Disks or manually provision new ones to
provide persistent storage for Kubernetes Pods. Azure Disks can be manually
provisioned in the Azure Portal, using ARM Templates or the Azure CLI. The 
following example uses the Azure CLI to manually provision an Azure
Disk. 

```bash
$ RG=myresourcegroup

$ az disk create \
  --resource-group $RG \
  --name k8s_volume_1  \
  --size-gb 20 \
  --query id \
  --output tsv
```

Using the Azure CLI command in the previous example should return the Azure ID of the Azure Disk
Object. If you are provisioning Azure resources using an alternative method,
make sure you retrieve the Azure ID of the Azure Disk, because it is needed for another step. 

```
/subscriptions/<subscriptionID>/resourceGroups/<resourcegroup>/providers/Microsoft.Compute/disks/<diskname>
```

You can now create Kubernetes Objects that refer to this Azure Disk. The following
example uses a Kubernetes Pod. However, the same Azure Disk syntax can be
used for DaemonSets, Deployments, and StatefulSets. In the following example, the
Azure Disk Name and ID reflect the manually created Azure Disk. 

```bash
$ cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: mypod-azuredisk
spec:
  containers:
  - image: nginx
    name: mypod
    volumeMounts:
      - name: mystorage
        mountPath: /data
  volumes:
      - name: mystorage
        azureDisk:
          kind: Managed
          diskName: k8s_volume_1
          diskURI: /subscriptions/<subscriptionID>/resourceGroups/<resourcegroup>/providers/Microsoft.Compute/disks/<diskname>
EOF
```

## Dynamically provision Azure Disks

### Define the Azure Disk Storage Class

Kubernetes can dynamically provision Azure Disks using the Azure Kubernetes
integration, which was configured when UCP was installed. For Kubernetes
to determine which APIs to use when provisioning storage, you must 
create Kubernetes Storage Classes specific to each storage backend. For more
information on Kubernetes Storage Classes, see [Storage
Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/).

In Azure there are 2 different Azure Disk types that can be consumed by
Kubernetes: Azure Disk Standard Volumes and Azure Disk Premium Volumes. For more
information on their differences, see [Azure
Disks](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/disks-types).

Depending on your use case, you can deploy one or both of the Azure Disk storage Classes (Standard and Advanced).

To create a Standard Storage Class:

```bash
$ cat <<EOF | kubectl create -f -
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: standard
provisioner: kubernetes.io/azure-disk
parameters:
  storageaccounttype: Standard_LRS
  kind: Managed
EOF
```

To Create a Premium Storage Class:

```bash
$ cat <<EOF | kubectl create -f -
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: premium
provisioner: kubernetes.io/azure-disk
parameters:
  storageaccounttype: Premium_LRS
  kind: Managed
EOF
```

To determine which Storage Classes have been provisioned:

```bash
$ kubectl get storageclasses
NAME       PROVISIONER                AGE
premium    kubernetes.io/azure-disk   1m
standard   kubernetes.io/azure-disk   1m
```

### Create an Azure Disk with a Persistent Volume Claim

After you create a Storage Class, you can use Kubernetes
Objects to dynamically provision Azure Disks. This is done using Kubernetes
Persistent Volumes Claims. For more information on Kubernetes Persistent Volume
Claims, see
[PVCs](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#introduction).

The following example uses the standard storage class and creates a 5 GiB Azure Disk. Alter these values to fit your use case. 

```bash
$ cat <<EOF | kubectl create -f -
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: azure-disk-pvc
spec:
  storageClassName: standard
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
EOF
```

At this point, you should see a new Persistent Volume Claim and Persistent Volume
inside of Kubernetes. You should also see a new Azure Disk created in the Azure
Portal.

```bash
$ kubectl get persistentvolumeclaim
NAME              STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
azure-disk-pvc    Bound     pvc-587deeb6-6ad6-11e9-9509-0242ac11000b   5Gi        RWO            standard       1m

$ kubectl get persistentvolume
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                     STORAGECLASS   REASON    AGE
pvc-587deeb6-6ad6-11e9-9509-0242ac11000b   5Gi        RWO            Delete           Bound     default/azure-disk-pvc    standard                 3m
```

### Attach the new Azure Disk to a Kubernetes pod

Now that a Kubernetes Persistent Volume has been created, you can mount this into
a Kubernetes Pod. The disk can be consumed by any Kubernetes object type, including
a Deployment, DaemonSet, or StatefulSet. However, the following example just mounts
the persistent volume into a standalone pod.

```bash
$ cat <<EOF | kubectl create -f -
kind: Pod
apiVersion: v1
metadata:
  name: mypod-dynamic-azuredisk
spec:
  containers:
    - name: mypod
      image: nginx
      ports:
        - containerPort: 80
          name: "http-server"
      volumeMounts:
        - mountPath: "/usr/share/nginx/html"
          name: storage
  volumes:
    - name: storage 
      persistentVolumeClaim:
        claimName: azure-disk-pvc
EOF
```

### Azure Virtual Machine data disk capacity

In Azure, there are limits to the number of data disks that can be attached to
each Virtual Machine. This data is shown in  [Azure Virtual Machine
Sizes](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/sizes-general). 
Kubernetes is aware of these restrictions, and prevents pods from
deploying on Nodes that have reached their maximum Azure Disk Capacity. 

This can be seen if a pod is stuck in the `ContainerCreating` stage:

```bash
$ kubectl get pods
NAME                  READY     STATUS              RESTARTS   AGE
mypod-azure-disk      0/1       ContainerCreating   0          4m
```

Describing the pod displays troubleshooting logs, showing the node has
reached its capacity:

```bash
$ kubectl describe pods mypod-azure-disk
<...>
  Warning  FailedAttachVolume  7s (x11 over 6m)  attachdetach-controller  AttachVolume.Attach failed for volume "pvc-6b09dae3-6ad6-11e9-9509-0242ac11000b" : Attach volume "kubernetes-dynamic-pvc-6b09dae3-6ad6-11e9-9509-0242ac11000b" to instance "/subscriptions/<sub-id>/resourceGroups/<rg>/providers/Microsoft.Compute/virtualMachines/worker-03" failed with compute.VirtualMachinesClient#CreateOrUpdate: Failure sending request: StatusCode=409 -- Original Error: failed request: autorest/azure: Service returned an error. Status=<nil> Code="OperationNotAllowed" Message="The maximum number of data disks allowed to be attached to a VM of this size is 4." Target="dataDisks"
```

## Where to go next

- [Deploy an Ingress Controller on
  Kubernetes](/ee/ucp/kubernetes/layer-7-routing/)
- [Discover Network Encryption on
  Kubernetes](/ee/ucp/kubernetes/kubernetes-network-encryption/)
