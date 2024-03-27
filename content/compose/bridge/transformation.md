---
title: Compose Bridge default transformation
description: Learn about the Compose Bridge default transformation
keywords: compose, bridge, kubernetes
---

Compose Bridge default transformation produces Kubernetes manifests so you can deploy
your compose application to Kubernetes enabled inside Docker Desktop.

Based on an arbitrary compose.yaml project, Compose Bridge will produce:

- A [Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) so all your resources are isolated and don't colide with another deployment
- A [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/) with an entry for each and every config resource in your compose model
- [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)s for application services 
- [Service](https://kubernetes.io/docs/concepts/services-networking/service/)s for ports exposed by your services, used for service-to-service communication
- [Service](https://kubernetes.io/docs/concepts/services-networking/service/)s for ports published by your services, with type `LoadBalancer` so that Docker Desktop will also expose same port on host
- [Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/) to replicate the networking topology expressed in compose
- [PersistentVolumeClaim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)s for your volumes, using `hostpath` storage class so that Docker Desktop manage volume creation
- [Secret](https://kubernetes.io/docs/concepts/configuration/secret/)s with your secret encoded - this is designed for local use in a testing environment

And a Kustomize overlay dedicated to Docker Desktop with:
 - Loadbalancer for services which need to expose ports on host
 - A PersistentVolumeClaim to use the Docker Desktop storage provisioner `desktop-storage-provisioner`
 - A Kustomize file to link the all those resources together
