---
title: Configuring NFS Storage for Kubernetes
description: Learn how to add support for NFS persistent storage by adding a default storage class.
keywords: Universal Control Plane, UCP, Docker EE, Kubernetes, storage, volume
redirect_from:
- /ee/ucp/admin/configure/use-nfs-volumes/
---

Users can provide persistent storage for workloads running on Docker Enterprise
by using NFS storage. These NFS shares, when mounted into the running container,
provide state to the application, managing data external to the container's
lifecycle. 

> Note: Provisioning an NFS server and exporting an NFS share are out of scope
> for this guide. Additionally, using external [Kubernetes
> plugins](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
> to dynamically provision NFS shares, is also out of scope for this guide. 

To mount existing NFS shares within Kubernetes Pods, we have 2 options:
 - Define NFS shares within the Pod definitions. NFS shares are defined
   manually by each tenant when creating a workload.
 - Define NFS shares as a Cluster object through Persistent Volumes, with
   the Cluster object lifecycle handled separately from the workload. This is common for
   operators who want to define a range of NFS shares for tenants to request and
   consume.

## Defining NFS Shares in the Pod definition

When defining workloads in Kubernetes manifest files, an end user can directly
reference the NFS shares to mount inside of each Pod. The NFS share is defined
within the Pod specification, which could be a standalone pod, or could be
wrapped in a higher-level object like a Deployment, Daemonset, or StatefulSet. 

The following example includes a running UCP cluster and a downloaded 
[client bundle](../../user-access/cli/#download-client-certificates) with
permission to schedule pods in a namespace. 

Here is an example pod specification with an NFS volume defined:

```bash
$ cat nfs-in-a-pod.yaml
kind: Pod
apiVersion: v1
metadata:
  name: nfs-in-a-pod
spec:
  containers:
    - name: app
      image: alpine
      volumeMounts:
        - name: nfs-volume
          mountPath: /var/nfs # Please change the destination you like the share to be mounted too
      command: ["/bin/sh"]
      args: ["-c", "sleep 500000"]
  volumes:
    - name: nfs-volume
      nfs:
        server: nfs.example.com # Please change this to your NFS server
        path: /share1 # Please change this to the relevant share
```

To deploy the pod, and ensure that it started up correctly, use the [kubectl](../../user-access/kubectl/) command line tool. 

```bash
$ kubectl create -f nfsinapod.yaml

$ kubectl get pods
NAME                     READY     STATUS      RESTARTS   AGE
nfs-in-a-pod             1/1       Running     0          6m
```

Verify everything was mounted correctly by getting a shell prompt
within the container and searching for your mount. 

```bash
$ kubectl exec -it nfs-in-a-pod sh
/ #
/ # mount | grep nfs.example.com
nfs.example.com://share1 on /var/nfs type nfs4 (rw,relatime,vers=4.0,rsize=262144,wsize=262144,namlen=255,hard,proto=tcp,timeo=600,retrans=2,sec=sys,clientaddr=172.31.42.23,local_lock=none,addr=nfs.example.com)
/ #
```

Because you defined the NFS share as part of the Pod spec, neither UCP nor Kubernetes
knows anything about this NFS share. This means that when the pod gets
deleted, the NFS share is unattached from the Cluster. However, the data remains in the NFS share.

## Exposing NFS shares as a Cluster Object

For this method, use the Kubernetes Objects [Persistent
Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes)
and [Persistent Volume
Claims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims)
to manage the lifecycle and access to NFS Shares. 

Here you can define multiple shares for a tenant to use within the
cluster. The [Persistent
Volume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes)
is a cluster wide object, so it can be pre-provisioned. A
[Persistent Volume
Claim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims)
is a claim by a tenant for use of a Persistent Volume within their namespace. 

> Note: In this case, 'NFS share lifecycle' is referring to granting and removing the
> end user's ability to consume NFS storage, not managing the lifecycle
> of the NFS Server.

### Persistent Volume

Define the Persistent Volume at the cluster level: 

```bash
$ cat pvwithnfs.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: my-nfs-share
spec:
  capacity:
    storage: 5Gi # This size is used to match a volume to a tenents claim
  accessModes:
    - ReadWriteOnce # Access modes are defined below
  persistentVolumeReclaimPolicy: Recycle # Reclaim policies are defined below 
  nfs:
    server: nfs.example.com # Please change this to your NFS server
    path: /share1 # Please change this to the relevant share
```

To create Persistent Volume objects at the Cluster level, you need a [Cluster
Role
Binding](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding)
grant. Again use the [kubectl](../../user-access/kubectl/) command line tool to create the
volume:

```
$ kubectl create -f pvwithnfs.yaml

$ kubectl get pv
NAME           CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                       STORAGECLASS   REASON    AGE

my-nfs-share   5Gi        RWO            Recycle          Available                               slow                     7s
```

#### Access Modes

The access mode for a NFS Persistent Volume can be any of the following modes:

- ***ReadWriteOnce*** – the volume can be mounted as read-write by a single node.
- ***ReadOnlyMany*** – the volume can be mounted read-only by many nodes.
- ***ReadWriteMany*** – the volume can be mounted as read-write by many nodes. 

The access mode in the Persistent Volume definition is used to match a
Persistent Volume to a Claim. When a Persistent Volume is defined and created
inside of Kubernetes, a Volume is not mounted. See [access
modes in the Kubernetes documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes),
for more details.

#### Reclaim

The [reclaim
policy](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#reclaiming)
is used to define what the cluster should do after a Persistent Volume has been
released from a Claim. A Persistent Volume Reclaim policy could be: Reclaim,
Recycle and Delete. See [Reclaiming in the Kubernetes
documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#reclaiming)
for a deeper understanding.

### Persistent Volume Claim

A tenant can now "claim" a Persistent Volume for use within their workloads
by using a Kubernetes Persistent Volume Claim. A Persistent Volume Claim resides within a namespace, 
and it attempts to match available Persistent Volumes
to what a tenant has requested.

``` bash
$ cat myapp-claim.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myapp-nfs
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce # Access modes for volumes is defined under Persistent Volumes
  resources:
    requests:
      storage: 5Gi # volume size requested
```

A tenant with a
[RoleBinding](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding)
to create Persistent Volume Claims can deploy this Persistent
Volume Claim. If there is a Persistent Volume that meets the tenant's
criteria, Kubernetes binds the Persistent Volume to the Claim. Again, this does not mount the share.

```bash
$ kubectl create -f myapp-claim.yaml
persistentvolumeclaim "myapp-nfs" created

$ kubectl get pvc
NAME        STATUS    VOLUME         CAPACITY   ACCESS MODES   STORAGECLASS   AGE
myapp-nfs   Bound     my-nfs-share   5Gi        RWO            slow           2s

$ kubectl get pv
NAME           CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                       STORAGECLASS   REASON    AGE
my-nfs-share   5Gi        RWO            Recycle          Bound     default/myapp-nfs           slow                     4m
```

### Defining a workload

Finally, a tenant can deploy a workload to consume the Persistent Volume Claim.
The Persistent Volume Claim is defined within the Pod specification, which could
be a standalone pod or could be wrapped in a higher-level object like a
Deployment, Daemonset, or StatefulSet. 

```bash
$ cat myapp-pod.yaml
kind: Pod
apiVersion: v1
metadata:
  name: pod-using-nfs
spec:
  containers:
    - name: app
      image: alpine
      volumeMounts:
      - name: data
          mountPath: /var/nfs # Please change the destination you like the share to be mounted too
      command: ["/bin/sh"]
      args: ["-c", "sleep 500000"]  
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: myapp-nfs
```

The pod can be deployed by a tenant using the
[kubectl](../../user-access/kubectl/) command line tool. Additionally, you can
verify that the pod is running successfully and that the NFS share has been mounted
inside of the container.

```bash
$ kubectl create -f myapp-pod.yaml

$ kubectl get pod
NAME                     READY     STATUS      RESTARTS   AGE
pod-using-nfs            1/1       Running     0          1m

$ kubectl exec -it pod-using-nfs sh
/ # mount | grep nfs.example.com
nfs.example.com://share1 on /var/nfs type nfs4 (rw,relatime,vers=4.1,rsize=262144,wsize=262144,namlen=255,hard,proto=tcp,timeo=600,retrans=2,sec=sys,clientaddr=172.31.42.23,local_lock=none,addr=nfs.example.com)
/ #
```

## Where to go next

- [Deploy an Ingress Controller on Kubernetes](/ee/ucp/kubernetes/layer-7-routing/)
- [Discover Network Encryption on Kubernetes](/ee/ucp/kubernetes/kubernetes-network-encryption/)
