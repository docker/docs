---
title: Configuring Azure Disk Storage for Kubernetes
description: Learn how to add persistent storage to your Docker Enterprise clusters running on Azure with Azure Disk. 
keywords: Universal Control Plane, UCP, Docker EE, Kubernetes, storage, volume
redirect_from:
---

Platform operators can provide persistent storage for workloads running on
Docker Enterprise when running on Microsoft Azure by using Azure Disk. Platform
operators can either pre-provision Azure Disks to be consumed by Kubernetes
Pods, or can use the Azure Kubernetes integration to dynamically provision Azure
Disks on demand.


## Prerequisites

This guide assumes you have already provisioned a UCP environment on to
Microsoft Azure. The Cluster must be provisioned after meeting all of the
prerequisites within [Install UCP on
Azure](/ee/ucp/admin/install/install-on-azure.md).

Additionally, the following guide is using the Kubernetes Command Line tool `$
kubectl` to provision Kubernetes objects within a UCP cluster. Therefore this
tool needs to be downloaded, along with a UCP client bundle. For more
information on configuring CLI access to UCP see [CLI Based
Access](/ee/ucp/user-access/cli/).

## Manually Provisioning Azure Disks

An operator can use existing Azure Disks or manually provision new ones to
provide persistent storage for Kubernetes Pods. Azure Disks can be manually
provisioned in the Azure Portal, using ARM Templates or using the Azure CLI. In
the below example we have used the Azure CLI to manually provision an Azure
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

Using the above Azure CLI command should return the Azure ID of the Azure Disk
Object. If you are provisioning Azure resources using an alternative method,
make sure you retrieve the Azure ID of the Azure Disk as you will need it later
on. 

```
/subscriptions/<subscriptionID>/resourceGroups/<resourcegroup>/providers/Microsoft.Compute/disks/<diskname>
```

You can now create Kubernetes Objects referring to this Azure Disk. The below
example is using a Kubernetes Pod, however the same Azure Disk syntax can be
used for DaemonSets, Deployments and StatefulSets. In the below example the
Azure Disk Name and ID to reflect the manually created Azure Disk. 

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

## Dynamically Provisioning Azure Disks

### Defining the Azure Disk Storage Class

Kubernetes can dynamically provision Azure Disks using the Azure Kubernetes
integration which was configured for you when UCP was installed. For Kubernetes
to know which APIs to use when provisioning storage, Platform operators need to
create Kubernetes Storage Classes specific to each storage backend. For more
information on Kubernetes Storage Classes, see [Storage
Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/).

In Azure there are 2 different types of Azure Disk that can be consumed by
Kubernetes: Azure Disk Standard Volumes and Azure Disk Premium Volumes. For more
information on their differences see [Azure
Disks](https://docs.microsoft.com/en-us/azure/virtual-machines/windows/disks-types).

Depending on your use case, you may choose to deploy 1 or both of the Azure Disk storage Classes (Standard and Advanced).

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

To see which Storage Classes have been provisioned:

```bash
$ kubectl get storageclasses
NAME       PROVISIONER                AGE
premium    kubernetes.io/azure-disk   1m
standard   kubernetes.io/azure-disk   1m
```

### Creating an Azure Disk with a Persistent Volume Claim

After an Operator has created a Storage Class, they can then use Kubernetes
Objects to dynamically provision Azure Disks. This is done using Kubernetes
Persistent Volumes Claims For more information on Kubernetes Persistent Volume
Claims, see
[PVCs](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#introduction).

The below example will use the standard storage class, and create a 5 GiB Azure Disk. These values can be altered to fit your use case. 

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

At this point you should see a new Persistent Volume Claim and Persistent Volume
inside of Kubernetes, you should also see a new Azure Disk created in the Azure
Portal.

```bash
$ kubectl get persistentvolumeclaim
NAME              STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
azure-disk-pvc    Bound     pvc-587deeb6-6ad6-11e9-9509-0242ac11000b   5Gi        RWO            standard       1m

$ kubectl get persistentvolume
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                     STORAGECLASS   REASON    AGE
pvc-587deeb6-6ad6-11e9-9509-0242ac11000b   5Gi        RWO            Delete           Bound     default/azure-disk-pvc    standard                 3m
```

### Attach the new Azure Disk to a Kubernetes Pod

Now that a Kubernetes Persistent Volume has been created, we can mount this into
a Kubernetes Pod. The disk can be consumed by any Kubernetes object type such
as a Deployment, DaemonSet or StatefulSet. However in the following example we
are just mounting the persistent volume into a standalone pod.

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

### Azure Virtual Machine Data Disk Capacity

In Azure there are limits to the number of data disks that can be attached to
each Virtual Machine. This data is shown in the [Azure Virtual Machine
Sizes](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/sizes-general)
. Kubernetes is aware of these restrictions, and will prevent pods from
deploying on Nodes that have reached their maximum Azure Disk Capacity. 

This can be seen if a pod is stuck in the `ContainerCreating` stage.

```bash
$ kubectl get pods
NAME                  READY     STATUS              RESTARTS   AGE
mypod-azure-disk      0/1       ContainerCreating   0          4m
```

Describing the pod will display troubleshooting logs, showing the node has
reached its capacity.

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
