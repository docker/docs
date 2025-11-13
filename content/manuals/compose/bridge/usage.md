---
title: Use the default Compose Bridge transformation
linkTitle: Usage
weight: 10
description: Learn how to use the default Compose Bridge transformation to convert Compose files into Kubernetes manifests
keywords: docker compose bridge, compose kubernetes transform, kubernetes from compose, compose bridge convert, compose.yaml to kubernetes
---

{{< summary-bar feature_name="Compose bridge" >}}

Compose Bridge includes a built-in transformation that automatically converts your Compose configuration into a set of Kubernetes manifests.

Based on your `compose.yaml` file, it produces:

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
 - A `Kustomization.yaml` file to link all the resources together.

If your Compose file defines a `models` section for a service, Compose Bridge automatically configures your deployment so your service can locate and use its models via Docker Model Runner.

For each declared model, the transformation injects two environment variables:

- `<MODELNAME>_URL`: The endpoint for Docker Model Runner serving that model  
- `<MODELNAME>_MODEL`: The model’s name or identifier

You can optionally customize these variable names using `endpoint_var` and `model_var`.

The default transformation generates two different overlays - one for Docker Desktop whilst using a local instance of Docker Model Runner, and a `model-runner` overlay with all the relevant Kubernetes resources to deploy Docker Model Runner in a pod. 

| Environment    | Endpoint                                        |
| -------------- | ----------------------------------------------- |
| Docker Desktop | `http://host.docker.internal:12434/engines/v1/` |
| Kubernetes     | `http://model-runner/engines/v1/`               |


For more details, see [Use Model Runner](use-model-runner.md).

## Use the default Compose Bridge transformation

To convert your Compose file using the default transformation:

```console
$ docker compose bridge convert
```

Compose looks for a `compose.yaml` file inside the current directory and generates Kubernetes manifests.

Example output:

```console
$ docker compose -f compose.yaml bridge convert
Kubernetes resource backend-deployment.yaml created
Kubernetes resource frontend-deployment.yaml created
Kubernetes resource backend-expose.yaml created
Kubernetes resource frontend-expose.yaml created
Kubernetes resource 0-my-project-namespace.yaml created
Kubernetes resource default-network-policy.yaml created
Kubernetes resource backend-service.yaml created
Kubernetes resource frontend-service.yaml created
Kubernetes resource kustomization.yaml created
Kubernetes resource backend-deployment.yaml created
Kubernetes resource frontend-deployment.yaml created
Kubernetes resource backend-service.yaml created
Kubernetes resource frontend-service.yaml created
Kubernetes resource kustomization.yaml created
Kubernetes resource model-runner-configmap.yaml created
Kubernetes resource model-runner-deployment.yaml created
Kubernetes resource model-runner-service.yaml created
Kubernetes resource model-runner-volume-claim.yaml created
Kubernetes resource kustomization.yaml created
```

All generated files are stored in the `/out` directory in your project.

## Deploy the generated manifests

> [!IMPORTANT]
>
> Before you deploy your Compose Bridge transformations, make sure you have [enabled Kubernetes](/manuals/desktop/settings-and-maintenance/settings.md#kubernetes) in Docker Desktop.

Once the manifests are generated, deploy them to your local Kubernetes cluster:

```console
$ kubectl apply -k out/overlays/desktop/
```

> [!TIP]
>
> You can convert and deploy your Compose project to a Kubernetes cluster from the Compose file viewer.
> 
> Make sure you are signed in to your Docker account, navigate to your container in the **Containers** view, and in the top-right corner select **View configurations** and then **Convert and Deploy to Kubernetes**. 

## Additional commands

Convert a `compose.yaml` file located in another directory:

```console
$ docker compose -f <path-to-file>/compose.yaml bridge convert
```

To see all available flags, run:

```console
$ docker compose bridge convert --help
```

## What's next?

- [Explore how you can customize Compose Bridge](customize.md)
- [Use Model Runner](use-model-runner.md).