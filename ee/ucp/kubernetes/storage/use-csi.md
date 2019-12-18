---
title: Using CSI drivers
description: Learn how to deploy a CSI driver with UCP.
keywords: Universal Control Plane, UCP, Docker EE, Kubernetes, storage, volume
---

>{% include enterprise_label_shortform.md %}

The Container Storage Interface (CSI) is a specification for container orchestrators to manage block and file-based 
volumes for storing data. Storage vendors can each create a single CSI driver that works with multiple 
container orchestrators. The Kubernetes community maintains sidecar containers that can be used by a containerized 
CSI driver to interface with Kubernetes controllers in charge of managing persistent volumes, attaching volumes to 
nodes (if applicable), mounting volumes to pods, taking snapshots, and more. These sidecar containers include 
a driver registrar, external attacher, external provisioner, and external snapshotter.

Docker Enterprise 3.0 supports version 1.0+ of the CSI specification. Therefore, UCP 3.2 (as part of Docker Enterprise 3.0) can manage storage backends that ship with an associated CSI driver, as illustrated in the following diagram.

![Kubernetes and CSI components](/ee/ucp/images/csi-plugins.png){: .with-border}

> Note
>
> Docker Enterprise does not provide CSI drivers. CSI drivers are provided by enterprise storage vendors. Kubernetes does not enforce a specific procedure for how Storage Providers (SP) should bundle and distribute CSI drivers.

Review the [Kubernetes CSI Developer Documentation](https://kubernetes-csi.github.io/docs/) for CSI architecture, 
security, and deployment details.

## Prerequisites

1. Select a CSI driver to use with Kubernetes. See [certified CSI drivers](#certified-csi-drivers) for the list.
2. Install Docker Enterprise 3.0, which includes UCP 3.2. 
3. Optionally, set the `--storage-expt-enabled` flag in the UCP install configuration if you want to enable
experimental storage features in Kubernetes 1.14. Features that are enabled include `VolumeSnapshotDataSource`,
`ExpandCSIVolumes`, `CSIMigration`, `CSIMigrationAWS`, `CSIMigrationGCE`, `CSIMigrationOpenStack`, 
and `VolumeSubpathEnvExpansion`. For details on these features, refer to [Feature Gates](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/#feature-gates).
4. Install the CSI plugin from your storage provider. For notes regarding installation, refer to your 
storage provider’s user manual.
5. Apply RBAC for sidecars and the CSI driver. For details on how to apply RBAC for your specific storage provider, 
refer to the storage vendor documentation for specific permissions and roles required for deploying CSI plugins 
on the cluster.
6. Perform static or dynamic provisioning of PVs using the CSI plugin as the provisioner. For details on how 
to provision volumes for your specific storage provider, refer to the storage provider’s user manual.

### Certified CSI drivers
The following table lists the UCP certified CSI drivers.

| Partner name | Kubernetes on Docker Enterprise 3.0 |
|--------------|-------------------------------------|
| NetApp       | Certified (Trident - CSI)           |
| EMC/Dell     | Certified (VxFlexOS CSI driver)     |
| VMware       | Certified (CSI)                     |
| Portworx     | Certified (CSI)                     |
| Nexenta      | Certified (CSI)                     |
| Blockbridge  | Certified (CSI)                     |
| Storidge     | Certified (CSI)                     |


## CSI driver deployment
Refer to documentation from your storage vendor around how to deploy the desired CSI driver. 
For easy deployment, storage vendors can package the CSI driver in containers. In the context of 
Kubernetes clusters, containerized CSI drivers are typically deployed as `StatefulSets` for 
managing the cluster-wide logic and `DaemonSets` for managing node-specific logic.

You can deploy multiple CSI drivers for different storage backends in the same cluster.

> Note
> - To avoid credential leak to user processes, Kubernetes recommends running CSI Controllers on master nodes and the CSI node plugin on worker nodes. 
> - UCP allows running privileged pods. This is needed to run CSI drivers.
> - The Docker daemon on the hosts must be configured with Shared Mount propagation for CSI. This is to allow the sharing of volumes mounted by one container into other containers in the same pod or even to other pods on the same node. By default, Docker daemon in UCP enables "Bidirectional Mount Propagation".

For additional information, refer to the [Kubernetes CSI documentation](https://kubernetes-csi.github.io/docs/deploying.html).

### Role-based access control (RBAC)
Pods containing CSI plugins need the appropriate permissions to access and manipulate Kubernetes objects. The desired cluster roles and bindings for service accounts associated with CSI driver pods can be configured through YAML files distributed by the storage vendor. UCP administrators must apply those YAML files to properly configure RBAC for the service accounts associated with CSI pods.

## Usage 

### Dynamic provisioning

Dynamic provisioning of persistent storage depends on the capabilities of the CSI driver and underlying storage backend. The provider of the CSI driver should document the parameters available for configuration. 
Refer to [CSI HostPath Driver](https://github.com/kubernetes-csi/csi-driver-host-path) for a generic CSI plugin example.

### Manage CSI deployment
The UCP user interface (UI) provides information about your CSI deployments, as shown in the following screen capture. In this example, a CSI Host Path Plugin was deployed as a `Pod`:

![UCP UI with CSI host plugin](/ee/ucp/images/csi-host-path-plugin.png)

In the UCP UI, you can navigate to **Kubernetes** > **Storage** for information about persistent storage objects such as `StorageClass`, `PersistentVolumeClaim`, and `PersistentVolume`. The following example provides information for objects specifically created using a CSI HostPath plugin:

![UCP UI with persistent storage object information](/ee/ucp/images/persistent-storage-object.png)

The **Volumes** section on the Pod details page shows that the Pod using this CSI HostPath plugin has a volume mounted into the Pod:

![UCP UI with CSI volume mount information](/ee/ucp/images/csi-volume-mounted.png)
