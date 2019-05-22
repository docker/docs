---
title: Configuring Azure Files Storage for Kubernetes
description: Learn how to add persistent storage to your Docker Enterprise clusters running on Azure with Azure Files.
keywords: Universal Control Plane, UCP, Docker EE, Kubernetes, storage, volume
redirect_from:
---

Platform operators can provide persistent storage for workloads running on
Docker Enterprise and Microsoft Azure by using Azure Files. You can either 
pre-provision Azure Files Shares to be consumed by
Kubernetes Pods or can you use the Azure Kubernetes integration to dynamically
provision Azure Files Shares on demand.

## Prerequisites

This guide assumes you have already provisioned a UCP environment on 
Microsoft Azure. The cluster must be provisioned after meeting all 
prerequisites listed in [Install UCP on
Azure](/ee/ucp/admin/install/install-on-azure.md).

Additionally, this guide uses the Kubernetes Command Line tool `$
kubectl` to provision Kubernetes objects within a UCP cluster. Therefore, you must download
this tool along with a UCP client bundle. For more
information on configuring CLI access to UCP, see [CLI Based
Access](/ee/ucp/user-access/cli/).

## Manually Provisioning Azure Files

You can use existing Azure Files Shares or manually provision new ones to
provide persistent storage for Kubernetes Pods. Azure Files Shares can be
manually provisioned in the Azure Portal using ARM Templates or using the Azure
CLI. The following example uses the Azure CLI to manually provision 
Azure Files Shares. 

### Creating an Azure Storage Account

When manually creating an Azure Files Share, first create an Azure
Storage Account for the file shares. If you have already provisioned
a Storage Account, you can skip to [Creating an Azure Files
Share](#creating-an-azure-file-share).

> **Note**: the Azure Kubernetes Driver does not support Azure Storage Accounts
> created using Azure Premium Storage. 

```bash
$ REGION=ukwest
$ SA=mystorageaccount
$ RG=myresourcegroup

$ az storage account create \
 --name $SA \
 --resource-group $RG \
 --location $REGION \
 --sku Standard_LRS
```

### Creating an Azure Files Share

Next, provision an Azure Files Share. The size of this share can be
adjusted to fit the end user's requirements. If you have already created an
Azure Files Share, you can skip to [Configuring a Kubernetes
Secret](#configuring-a-kubernetes-secret).

```bash
$ SA=mystorageaccount
$ RG=myresourcegroup
$ FS=myfileshare
$ SIZE=5

# This Azure Collection String can also be found in the Azure Portal
$ export AZURE_STORAGE_CONNECTION_STRING=`az storage account show-connection-string --name $SA --resource-group $RG -o tsv`

$ az storage share create \
  --name $FS \
  --quota $SIZE \
  --connection-string $AZURE_STORAGE_CONNECTION_STRING
```

### Configuring a Kubernetes Secret

After a File Share has been created, you must load the Azure Storage
Account Access key as a Kubernetes Secret into UCP. This provides access to
the file share when Kubernetes attempts to mount the share into a pod. This key
can be found in the Azure Portal or retrieved as shown in the following example by the Azure CLI: 

```bash
$ SA=mystorageaccount
$ RG=myresourcegroup
$ FS=myfileshare

# The Azure Storage Account Access Key can also be found in the Azure Portal
$ STORAGE_KEY=$(az storage account keys list --resource-group $RG --account-name $SA --query "[0].value" -o tsv)

$ kubectl create secret generic azure-secret \
  --from-literal=azurestorageaccountname=$SA \
  --from-literal=azurestorageaccountkey=$STORAGE_KEY
```

### Mount the Azure Files Share into a Kubernetes Pod

The final step is to mount the Azure Files Share, using the Kubernetes Secret, into
a Kubernetes Pod. The following code creates a standalone Kubernetes pod, but you
can also use alternative Kubernetes Objects such as Deployments, DaemonSets, or
StatefulSets, with the existing Azure Files Share.

```bash
$ FS=myfileshare

$ cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: mypod-azurefile
spec:
  containers:
  - image: nginx
    name: mypod
    volumeMounts:
      - name: mystorage
        mountPath: /data
  volumes:
  - name: mystorage
    azureFile:
      secretName: azure-secret
      shareName: $FS
      readOnly: false
EOF
```

## Dynamically Provisioning Azure Files Shares

### Defining the Azure Disk Storage Class

Kubernetes can dynamically provision Azure Files Shares using the Azure
Kubernetes integration, which was configured when UCP was installed. For
Kubernetes to know which APIs to use when provisioning storage, you must
create Kubernetes Storage Classes specific to each storage
backend. For more information on Kubernetes Storage Classes, see [Storage
Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/).

> Today, only the Standard Storage Class is supported when using the Azure
> Kubernetes Plugin. File shares using the Premium Storage Class will fail to
> mount. 

```bash
$ cat <<EOF | kubectl create -f -
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: standard
provisioner: kubernetes.io/azure-file
mountOptions:
  - dir_mode=0777
  - file_mode=0777
  - uid=1000
  - gid=1000
parameters:
  skuName: Standard_LRS
EOF
```

To see which Storage Classes have been provisioned:

```bash
$ kubectl get storageclasses
NAME       PROVISIONER                AGE
azurefile  kubernetes.io/azure-file   1m
```

### Creating an Azure Files Share using a Persistent Volume Claim

After you create a Storage Class, you can use Kubernetes
Objects to dynamically provision Azure Files Shares. This is done using
Kubernetes Persistent Volumes Claims
[PVCs](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#introduction).
Kubernetes uses an existing Azure Storage Account if one exists inside of the 
Azure Resource Group. If an Azure Storage Account does not exist,
Kubernetes creates one. 

The following example uses the standard storage class and creates a 5 GB Azure
File Share. Alter these values to fit your use case. 

```bash
$ cat <<EOF | kubectl create -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: azure-file-pvc
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: standard
  resources:
    requests:
      storage: 5Gi
EOF
```

At this point, you should see a newly created Persistent Volume Claim and Persistent Volume: 

```bash
$ kubectl get pvc
NAME             STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
azure-file-pvc   Bound     pvc-f7ccebf0-70e0-11e9-8d0a-0242ac110007   5Gi        RWX            standard       22s

$ kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                    STORAGECLASS   REASON    AGE
pvc-f7ccebf0-70e0-11e9-8d0a-0242ac110007   5Gi        RWX            Delete           Bound     default/azure-file-pvc   standard                 2m
```

### Attach the new Azure Files Share to a Kubernetes Pod

Now that a Kubernetes Persistent Volume has been created, mount this into
a Kubernetes Pod. The file share can be consumed by any Kubernetes object type
such as a Deployment, DaemonSet, or StatefulSet. However, the following
example just mounts the persistent volume into a standalone pod.

```bash
$ cat <<EOF | kubectl create -f -
kind: Pod
apiVersion: v1
metadata:
  name: mypod
spec:
  containers:
    - name: task-pv-container
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
       claimName: azure-file-pvc
EOF
```

## Where to go next

- [Deploy an Ingress Controller on
  Kubernetes](/ee/ucp/kubernetes/layer-7-routing/)
- [Discover Network Encryption on
  Kubernetes](/ee/ucp/kubernetes/kubernetes-network-encryption/)
