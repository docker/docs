---
title: Configuring Azure File Storage for Kubernetes
description: Learn how to add persistent storage to your Docker Enterprise clusters running on Azure with Azure File.
keywords: Universal Control Plane, UCP, Docker EE, Kubernetes, storage, volume
redirect_from:
---

Operators can provide persistent storage for workloads running on Docker
Enterprise when running on Microsoft Azure by using Azure File. Operators can
either pre-provision Azure File Shares to be consumed by Kubernetes Pods or can
you use the Azure Kubernetes integration to dynamically provision Azure File
Shares on demand.

## Pre-Requisites

This guide assumes you have already provisioned a UCP environment on to
Microsoft Azure. The Cluster must be provisioned after meeting all of the
pre-requisites within [Install UCP on
Azure](/ee/ucp/admin/install/install-on-azure.md).

Additionally the following guide is using the Kubernetes Command Line tool `$
kubectl` to provision Kubernetes objects within a UCP cluster. Therefore this
tool needs to be downloaded, along with a UCP client bundle. For more
information on configuring CLI access to UCP see [CLI Based
Access](/ee/ucp/user-access/cli/).

## Manually Provisioning Azure File

An operator can use existing Azure File Shares or manually provision new ones to
provide persistent storage for Kubernetes Pods. Azure File Shares can be
manually provisioned in the Azure Portal, using ARM Templates or using the Azure
Cli. In the below example we have used the Azure Cli to manually provision an
Azure File Shares. 

### Creating an Azure Storage Account

The first step when manually creating an Azure File Share, is to create an Azure
Storage Account for the file shares to live in. If you have already provisioned
a Storage Account, you can move on to [Creating an Azure File
Share](#creating-an-azure-file-share).

> Note the Azure Kubernetes Driver does not support Azure Storage Accounts
> created using Azure Premium Storage. 

```
$ REGION=ukwest
$ SA=mystorageaccount
$ RG=myresourcegroup

$ az storage account create \
 --name $SA \
 --resource-group $RG \
 --location $REGION \
 --sku Standard_LRS
```

### Creating an Azure File Share

Next we will provision an Azure File Share, the size of this share can be
adjusted to fit the end user's requirements. If you have already created an
Azure File Share, you can move on to [Configuring a Kubernetes
Secret](#configuring-a-kubernetes-secret)

```
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

After a File Share has been created, an Operator needs to load the Azure Storage
Account Access key as a Kubernetes Secret into UCP, this will provide access to
the file share when Kubernetes attempts to mount the share into a pod. This key
can be found in the Azure Portal, or retrieved as shown below by the Azure Cli. 

```
$ SA=mystorageaccount
$ RG=myresourcegroup
$ FS=myfileshare

# The Azure Storage Account Access Key can also be found in the Azure Portal
$ STORAGE_KEY=$(az storage account keys list --resource-group $RG --account-name $SA --query "[0].value" -o tsv)

$ kubectl create secret generic azure-secret \
  --from-literal=azurestorageaccountname=$SA \
  --from-literal=azurestorageaccountkey=$STORAGE_KEY
```

### Mount the Azure File Share into a Kubernetes Pod

The final step is to mount the Azure File Share using the Kubernetes Secret into
a Kubernetes Pod. Below we have created a standalone Kubernetes pod, however you
could use alternative Kubernetes Objects such as Deployments, Daemonsets or
Statefulsets with the existing Azure File Share.

```
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

## Dynamically Provisioning Azure File Shares

### Defining the Azure Disk Storage Class

Kubernetes can dynamically provision Azure File Shares using the Azure
Kubernetes integration which was configured for you when UCP was installed. For
Kubernetes to know which APIs to use when provisioning storage, Operators need
to create Kubernetes Storage Classes specific to each storage backend. For more
information on Kubernetes Storage Class, see [Storage
Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/).

> Today only the Standard Storage Class is supported when using the Azure
> Kubernetes Plugin, File shares using the Premium Storage Class will fail to
> mount. 

```
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

To see which Storage Classes have been provisioned.

```
$ kubectl get storageclasses
NAME       PROVISIONER                AGE
azurefile  kubernetes.io/azure-file   1m
```

### Creating an Azure File Share using a Persistent Volume Claim

After an Operator has created a Storage Class, they can then use Kubernetes
Objects to dynamically provision Azure File Shares. This is done using
Kubernetes Persistent Volumes Claims
[PVCs](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#introduction).
Kubernetes will try and use an existing Azure Storage Account if one exists in
side of the Azure Resource Group, if an Azure Storage Account does not exist
Kubernetes will create one. 

The below example will use the standard storage class, and create a 5 GiB Azure
File Share. These values can be altered to fit your use case. 

```
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

At this point you should see a new Persistent Volume Claim and Persistent Volume
have been created. 

```
$ kubectl get pvc
NAME             STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
azure-file-pvc   Bound     pvc-f7ccebf0-70e0-11e9-8d0a-0242ac110007   5Gi        RWX            standard       22s

$ kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                    STORAGECLASS   REASON    AGE
pvc-f7ccebf0-70e0-11e9-8d0a-0242ac110007   5Gi        RWX            Delete           Bound     default/azure-file-pvc   standard                 2m
```

### Attach the new Azure File Share to a Kubernetes Pod

Now that a Kubernetes Persistent Volume has been created, we can mount this into
a Kubernetes Pod. The file share can be consumed by any Kubernetes object type
such as a Deployment, Daemonset or Stateful set. However in the following
example we are just mounting the persistent volume into a standalone pod.

```
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
