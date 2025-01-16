---
title: Use the default Compose Bridge transformation
linkTitle: Usage
weight: 10
description: Learn about and use the Compose Bridge default transformation
keywords: compose, bridge, kubernetes
---

{{< summary-bar feature_name="Compose bridge" >}}

Compose Bridge supplies an out-of-the box transformation for your Compose configuration file. Based on an arbitrary `compose.yaml` file, Compose Bridge produces:

- A [Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) so all your resources are isolated and don't conflict with resources from other deployments.
- A [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) with an entry for each and every [config](/reference/compose-file/configs.md) resource in your Compose application.
- [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) for application services. This ensures that the specified number of instances of your application are maintained in the Kubernetes cluster.
- [Services](https://kubernetes.io/docs/concepts/services-networking/service/) for ports exposed by your services, used for service-to-service communication.
- [Services](https://kubernetes.io/docs/concepts/services-networking/service/) for ports published by your services, with type `LoadBalancer` so that Docker Desktop will also expose the same port on the host.
- [Network policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/) to replicate the networking topology defined in your `compose.yaml` file. 
- [PersistentVolumeClaims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) for your volumes, using `hostpath` storage class so that Docker Desktop manages volume creation.
- [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/) with your secret encoded. This is designed for local use in a testing environment.

It also supplies a Kustomize overlay dedicated to Docker Desktop with:
 - `Loadbalancer` for services which need to expose ports on host.
 - A `PersistentVolumeClaim` to use the Docker Desktop storage provisioner `desktop-storage-provisioner` to handle volume provisioning more effectively.
 - A Kustomize file to link all the resources together.

## Use the default Compose Bridge transformation

To use the default transformation run the following command:

```console
$ compose-bridge convert
```

Compose looks for a `compose.yaml` file inside the current directory and then converts it.

The following output is displayed 
```console
$ compose-bridge convert -f compose.yaml 
Kubernetes resource api-deployment.yaml created
Kubernetes resource db-deployment.yaml created
Kubernetes resource web-deployment.yaml created
Kubernetes resource api-expose.yaml created
Kubernetes resource db-expose.yaml created
Kubernetes resource web-expose.yaml created
Kubernetes resource 0-avatars-namespace.yaml created
Kubernetes resource default-network-policy.yaml created
Kubernetes resource private-network-policy.yaml created
Kubernetes resource public-network-policy.yaml created
Kubernetes resource db-db_data-persistentVolumeClaim.yaml created
Kubernetes resource api-service.yaml created
Kubernetes resource web-service.yaml created
Kubernetes resource kustomization.yaml created
Kubernetes resource db-db_data-persistentVolumeClaim.yaml created
Kubernetes resource api-service.yaml created
Kubernetes resource web-service.yaml created
Kubernetes resource kustomization.yaml created
```

These files are then stored within your project in the `/out` folder. 

The Kubernetes manifests can then be used to run the application on Kubernetes using
the standard deployment command `kubectl apply -k out/overlays/desktop/`.

> [!NOTE]
>
> Make sure you have enabled Kubernetes in Docker Desktop before you deploy your Compose Bridge transformations.

If you want to convert a `compose.yaml` file that is located in another directory, you can run:

```console
$ compose-bridge convert -f <path-to-file>/compose.yaml 
```

To see all available flags, run:

```console
$ compose-bridge convert --help
```

> [!TIP]
>
> You can now convert and deploy your Compose project to a Kubernetes cluster from the Compose file viewer.
> 
> Make sure you are signed in to your Docker account, navigate to your container in the **Containers** view, and in the top-right corner select **View configurations** and then **Convert and Deploy to Kubernetes**. 

## What's next?

- [Explore how you can customize Compose Bridge](customize.md)
- [Explore the advanced integration](advanced-integration.md)
