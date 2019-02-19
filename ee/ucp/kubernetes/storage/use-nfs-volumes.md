---
title: Configuring NFS Storage for Kubernetes
description: Learn how to add support for NFS persistent storage by adding a default storage class.
keywords: Universal Control Plane, UCP, Docker EE, Kubernetes, storage, volume
redirect_from:
- /ee/ucp/admin/configure/use-nfs-volumes/
---

Users can provide persistent storage for workloads running on Docker Enterprise
by using NFS storage. These NFS shares, when mounted into the running container,
providing state to the application, managing data externally to the container's
lifecycle. 

> Note: Provisioning an NFS server and/or exporting an NFS share is out of scope
> of this guide. Additionally, using external [Kubernetes
> plugins](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
> to dynamically provision NFS shares, is also out of scope for this guide. 

To mount existing NFS shares within Kubernetes Pods, we have 2 options:
 - We can define NFS shares within our Pod definitions. NFS shares are defined
   manually by each tenant when creating a workload.
 - We can define NFS shares as a Cluster object through Persistent Volumes, with
   its lifecycle handled separately to the workload. This is common if an
   operator wanted to define a range of NFS shares, for tenants to request and
   consume.

## Defining NFS Shares in the Pod Spec

When defining workloads in Kubernetes manifest files, an end user can directly
reference the NFS shares to mount inside of each Pod. The NFS share is defined
within the Pod specification, this could be a standalone pod, or could be
wrapped in a higher-level object like a Deployment, Daemonset or StatefulSet. 

In the following example, we have a running UCP cluster, and have downloaded a
[client bundle](../../user-access/cli/#download-client-certificates), with
permission to schedule pods in a namespace. 

An example pod specification with an NFS volume defined:

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

To deploy the pod, and ensure that is has started up correctly we will use [kubectl](../../user-access/kubectl/) command line. 

```bash
$ kubectl create -f nfsinapod.yaml

$ kubectl get pods
NAME                     READY     STATUS      RESTARTS   AGE
nfs-in-a-pod             1/1       Running     0          6m
```

We can check everything has been mounted correctly by getting a shell prompt
within the container, and searching for our mount. 

```bash
$ kubectl exec -it pod-using-nfs sh
/ #
/ # mount | grep nfs.example.com
nfs.example.com://share1 on /var/nfs type nfs4 (rw,relatime,vers=4.0,rsize=262144,wsize=262144,namlen=255,hard,proto=tcp,timeo=600,retrans=2,sec=sys,clientaddr=172.31.42.23,local_lock=none,addr=nfs.example.com)
/ #
```

As we have defined the NFS share as part of the Pod Spec, UCP or Kubernetes
doesn't know anything about this NFS share. This means that when the pod gets
deleted, the NFS share will be unattached from the Cluster. The data will of
course still remain in the NFS share.

## Exposing NFS shares as a Cluster Object

For this method we will use the Kubernetes Objects [Persistent
Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes)
and [Persistent Volume
Claims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims)
to manage the lifecycle and access to NFS Shares. 

Here an operator could define multiple shares for a tenant to use within the
cluster. The [Persistent
Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes)
is a cluster wide object so could be pre-provisioned by an operator. A
[Persistent Volume
Claims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims)
is a claim by a tenant, for use of a PV within their namespace. 

> Note: NFS Share Lifecycle in this sense, is referring to granting and removing
> end user's ability to consume NFS storage, rather than managing the lifecycle
> of the NFS Server.

### Persistent Volume

As an operator define the persistent volume at the cluster level: 

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

To create a Physical Volume on the cluster, an operator would need a [Cluster
Role
Binding](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding)
grant, to create persistent volume objects at the Cluster level. Once again a we
will use the [kubectl](../../user-access/kubectl/) command line to create the
volume.

```
$ kubectl create -f pvwithnfs.yaml

$ kubectl get pv
NAME           CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                       STORAGECLASS   REASON    AGE

my-nfs-share   5Gi        RWO            Recycle          Available                               slow                     7s
```

#### Access Modes

The access mode for a NFS persistent volume can either be:

- ReadWriteOnce – the volume can be mounted as read-write by a single node
- ReadOnlyMany – the volume can be mounted read-only by many nodes
- ReadWriteMany – the volume can be mounted as read-write by many nodes 

The access mode in the Persistent Volume definition is used to match a
Persistent Volume to a Claim. When a Persistent Volume is defined and created
inside of Kubernetes, a Volume is not mounted. For more information on [access
modes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes)
see the kubernetes documentation. 

#### Reclaim

The [reclaim
policy](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#reclaiming)
is used to define what the cluster should do after a persistent volume has been
released from a claim. A persistent volume reclaim policy could be: Reclaim,
Recycle and Delete. Please see the [Kubernetes
documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#reclaiming)
for a deeper understanding.

### Persistent Volume Claim

A tenant can now "claim" that persistent volume for use within their workloads
by using a Kubernetes persistent volume claim. A persistent volume claim will
live within a namespace, and it will try and match available persistent volumes
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

A tenant, with a
[RoleBinding](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding)
to create persistent volume claims, will now be able to deploy this persistent
volume claim. Assuming there is a persistent volume that meets the tenants
criteria, Kubernetes will now bind the persistent volume to the Claim. Once
again, this is not mounting the share.

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

### Defining a Workload

Finally, a tenant can deploy a workload to consume this persistent volume claim.
The persistent volume claim is defined within the Pod specification, this could
be a standalone pod, or could be wrapped in a higher-level object like a
Deployment, Daemonset or StatefulSet. 

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

The pod can be deployed by a tenant, using the
[kubectl](../../user-access/kubectl/) command line tool. Additionally, we can
check that the pod is running successfully, and the NFS share has been mounted
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