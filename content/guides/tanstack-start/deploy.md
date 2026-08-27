---
title: Test your TanStack Start deployment
linkTitle: Test your deployment
weight: 70
keywords: deploy, kubernetes, tanstack start
description: Learn how to deploy locally to test and debug your Kubernetes deployment
---

## Prerequisites

Before you begin, make sure you've completed the following:

- Complete all the previous sections of this guide, starting with
  [Containerize TanStack Start application](containerize.md).
- [Enable Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

> [!NOTE]
> New to Kubernetes? Visit the
> [Kubernetes basics tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/)
> to learn how clusters, pods, deployments, and services work.

---

## Overview

This section guides you through deploying your containerized TanStack Start
application locally using
[Docker Desktop's built-in Kubernetes](/desktop/kubernetes/). A local cluster
lets you test and debug workloads before promoting them to staging or
production.

---

## Create a Kubernetes YAML file

Follow these steps to define your deployment configuration:

1. In the root of your project, create a file named
   `tanstack-start-kubernetes.yaml`.

2. Open the file in your IDE or preferred text editor.

3. Add the following configuration. Replace `{DOCKER_USERNAME}` and
   `{DOCKERHUB_PROJECT_NAME}` with your Docker Hub username and repository
   name from
   [Automate your builds with GitHub Actions](configure-github-actions.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tanstack-start
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tanstack-start
  template:
    metadata:
      labels:
        app: tanstack-start
    spec:
      containers:
        - name: tanstack-start-container
          image: {DOCKER_USERNAME}/{DOCKERHUB_PROJECT_NAME}:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 3000
          env:
            - name: NODE_ENV
              value: "production"
            - name: HOST
              value: "0.0.0.0"
            - name: PORT
              value: "3000"
---
apiVersion: v1
kind: Service
metadata:
  name: tanstack-start-service
  namespace: default
spec:
  type: NodePort
  selector:
    app: tanstack-start
  ports:
    - port: 3000
      targetPort: 3000
      nodePort: 30001
```

This manifest defines two Kubernetes resources, separated by `---`:

- **Deployment** — Runs a single replica of your TanStack Start application.
  The pod uses the Docker image built and pushed by your GitHub Actions
  workflow. The container listens on port `3000`.

- **Service (NodePort)** — Exposes the pod on port `30001` on your host,
  forwarding traffic to port `3000` in the container. Open
  [http://localhost:30001](http://localhost:30001) in your browser.

> [!NOTE]
> To learn more about Kubernetes objects, see the
> [Kubernetes documentation](https://kubernetes.io/docs/home/).

---

## Deploy and check your application

### Step 1: Apply the Kubernetes configuration

From the directory that contains `tanstack-start-kubernetes.yaml`, run:

```console
$ kubectl apply -f tanstack-start-kubernetes.yaml
```

Expected output:

```shell
deployment.apps/tanstack-start created
service/tanstack-start-service created
```

### Step 2: Check the deployment status

```console
$ kubectl get deployments
```

Example output:

```shell
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
tanstack-start     1/1     1            1           14s
```

### Step 3: Verify the service exposure

```console
$ kubectl get services
```

Example output:

```shell
NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
tanstack-start-service   NodePort    10.100.244.65    <none>        3000:30001/TCP   1m
```

### Step 4: Access your app in the browser

Open [http://localhost:30001](http://localhost:30001). You should see your
TanStack Start application served from the local Kubernetes cluster.

### Step 5: Clean up Kubernetes resources

When you're done testing:

```console
$ kubectl delete -f tanstack-start-kubernetes.yaml
```

Expected output:

```shell
deployment.apps "tanstack-start" deleted
service "tanstack-start-service" deleted
```

---

## Summary

In this section, you deployed your TanStack Start application to a local
Kubernetes cluster using Docker Desktop.

What you accomplished:

- Created a Kubernetes Deployment and NodePort Service for your app
- Used `kubectl apply` to deploy the application locally
- Verified the app at `http://localhost:30001`
- Cleaned up Kubernetes resources after testing

---

## Related resources

- [Kubernetes documentation](https://kubernetes.io/docs/home/) – Core concepts,
  workloads, and services
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md) –
  Use Docker Desktop's built-in Kubernetes support
- [`kubectl` CLI reference](https://kubernetes.io/docs/reference/kubectl/) –
  Manage clusters from the command line
- [Kubernetes Deployment resource](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) –
  Manage and scale applications
- [Kubernetes Service resource](https://kubernetes.io/docs/concepts/services-networking/service/) –
  Expose applications to internal and external traffic
