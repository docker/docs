---
title: Configuring vSphere Volumes for Kubernetes
description: Learn how to add persistent storage to your Docker Enterprise clusters using vSphere Volumes.
keywords: Universal Control Plane, UCP, Docker Enterprise, Kubernetes, storage, volume
---

>{% include enterprise_label_shortform.md %}

The [vSphere Storage for Kubernetes driver](https://vmware.github.io/vsphere-storage-for-kubernetes/documentation/index.html) enables customers to address persistent storage requirements for Kubernetes pods  in vSphere environments. The driver allows you to create a persistent  volume on a Virtual Machine File System (VMFS), and use it to manage persistent storage requirements independent of pod and VM lifecycle.

> Note
>
> Of the three main storage backends offered by vSphere on Kubernetes (VMFS, vSAN, and NFS), Docker supports VMFS.

You can use the vSphere Cloud Provider to manage storage with Kubernetes in UCP 3.1 and later. This includes support for:

* Volumes
* Persistent volumes
* StorageClasses and provisioning volumes

## Prerequisites
* Ensure that `vsphere.conf` is populated according to the [vSphere Cloud Provider Configuration Deployment Guide](https://vmware.github.io/vsphere-storage-for-kubernetes/documentation/existing.html#create-the-vsphere-cloud-config-file-vsphereconf).
* The `disk.EnableUUID` value on the worker VMs must be set to `True`.

## Configure for Kubernetes

Kubernetes cloud providers provide a method of provisioning cloud resources through Kubernetes via the `--cloud-provider` option. This is to ensure that the kubelet is aware that it must be initialized by the ucp-kube-controller-manager before any work is scheduled.

```bash
docker container run --rm -it --name ucp -e REGISTRY_USERNAME=$REGISTRY_USERNAME -e REGISTRY_PASSWORD=$REGISTRY_PASSWORD \
  -v /var/run/docker.sock:/var/run/docker.sock \
  "dockereng/ucp:3.1.0-tp2" \
  install \
  --host-address <HOST_ADDR> \
  --admin-username admin \
  --admin-password XXXXXXXX \
  --cloud-provider=vsphere \
  --image-version latest:
```

## Create a StorageClass

1. Create a StorageClass with a user specified disk format.
```bash
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast
provisioner: kubernetes.io/vsphere-volume
parameters:
  diskformat: zeroedthick
```
For example, `diskformat` can be `thin`, `zeroedthick`, or `eagerzeroedthick`. The default format is `thin`.
2. Create a StorageClass with a disk format on a user-specified datastore.
```bash
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast
provisioner: kubernetes.io/vsphere-volume
parameters:
    diskformat: zeroedthick
    datastore: VSANDatastore
```
You can also specify the `datastore` in the StorageClass. The volume will be created on the datastore specified in the StorageClass, which in this case is `VSANDatastore`. This field is optional. If the datastore is not specified, then the volume will be created on the datastore specified in the vSphere configuration file used to initialize the vSphere Cloud Provider.

For more information on Kubernetes StorageClasses, see [Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/).

## Deploy vSphere Volumes

After you create a StorageClass, you can create [PersistentVolumes (PV)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#introduction) that deploy volumes attached to hosts and mounted inside pods.  A PersistentVolumeClaim (PVC) is a claim for storage resources that are bound to a PV when storage resources are granted. 

We recommend that you use the StorageClass and PVC resources as these abstraction layers provide more portability as well as control over the storage layer across environments.

To deploy vSphere volumes:

1. [Create a PVC from the plugin](https://vmware.github.io/vsphere-storage-for-kubernetes/documentation/policy-based-mgmt.html). When you define a PVC to use the StorageClass, a PV is created and bound.
2. [Create a reference to the PVC from the Pod](https://vmware.github.io/vsphere-storage-for-kubernetes/documentation/policy-based-mgmt.html).
3. Start a Pod using the PVC that you defined.

## Where to go next
* [Configuring iSCSI for Kubernetes](https://docs.docker.com/ee/ucp/kubernetes/storage/use-iscsi/)
* [Using CSI Drivers](https://docs.docker.com/ee/ucp/kubernetes/storage/use-csi/)
