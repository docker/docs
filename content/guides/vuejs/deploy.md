---
title: Test your Vue.js deployment
linkTitle: Test your deployment
weight: 60
keywords: deploy, kubernetes, vue, vue.js
description: Learn how to deploy locally to test and debug your Kubernetes deployment

---

## Prerequisites

Before you begin, make sure you’ve completed the following:
- Complete all the previous sections of this guide, starting with [Containerize Vue.js application](containerize.md).
- [Enable Kubernetes](/manuals/desktop/features/kubernetes.md#install-and-turn-on-kubernetes) in Docker Desktop.

> **New to Kubernetes?**  
> Visit the [Kubernetes basics tutorial](https://kubernetes.io/docs/tutorials/kubernetes-basics/) to get familiar with how clusters, pods, deployments, and services work.

---

## Overview

This section guides you through deploying your containerized Vue.js application locally using [Docker Desktop’s built-in Kubernetes](/desktop/kubernetes/). Running your app in a local Kubernetes cluster closely simulates a real production environment, enabling you to test, validate, and debug your workloads with confidence before promoting them to staging or production.

---

## Create a Kubernetes YAML file

Follow these steps to define your deployment configuration:

1. In the root of your project, create a new file named: vuejs-sample-kubernetes.yaml

2. Open the file in your IDE or preferred text editor.

3. Add the following configuration, and be sure to replace `{DOCKER_USERNAME}` and `{DOCKERHUB_PROJECT_NAME}` with your actual Docker Hub username and repository name from the previous [Automate your builds with GitHub Actions](configure-github-actions.md).


```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vuejs-sample
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vuejs-sample
  template:
    metadata:
      labels:
        app: vuejs-sample
    spec:
      containers:
        - name: vuejs-container
          image: {DOCKER_USERNAME}/{DOCKERHUB_PROJECT_NAME}:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
          resources:
            limits:
              cpu: "500m"
              memory: "256Mi"
            requests:
              cpu: "250m"
              memory: "128Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: vuejs-sample-service
  namespace: default
spec:
  type: NodePort
  selector:
    app: vuejs-sample
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 30001
```

This manifest defines two key Kubernetes resources, separated by `---`:

- Deployment
  Deploys a single replica of your Vue.js application inside a pod. The pod uses the Docker image built and pushed by your GitHub Actions CI/CD workflow  
  (refer to [Automate your builds with GitHub Actions](configure-github-actions.md)).  
  The container listens on port `8080`, which is typically used by [Nginx](https://nginx.org/en/docs/) to serve your production Vue.js app.

- Service (NodePort) 
  Exposes the deployed pod to your local machine.  
  It forwards traffic from port `30001` on your host to port `8080` inside the container.  
  This lets you access the application in your browser at [http://localhost:30001](http://localhost:30001).

> [!NOTE]
> To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

---

## Deploy and check your application

Follow these steps to deploy your containerized Vue.js app into a local Kubernetes cluster and verify that it’s running correctly.

### Step 1. Apply the Kubernetes configuration

In your terminal, navigate to the directory where your `vuejs-sample-kubernetes.yaml` file is located, then deploy the resources using:

```console
  $ kubectl apply -f vuejs-sample-kubernetes.yaml
```

If everything is configured properly, you’ll see confirmation that both the Deployment and the Service were created:

```shell
  deployment.apps/vuejs-sample created
  service/vuejs-sample-service created
```
   
This confirms that both the Deployment and the Service were successfully created and are now running inside your local cluster.

### Step 2. Check the Deployment status

Run the following command to check the status of your deployment:
   
```console
  $ kubectl get deployments
```

You should see output similar to the following:

```shell
  NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
  vuejs-sample         1/1     1            1           1m14s
```

This confirms that your pod is up and running with one replica available.

### Step 3. Verify the Service exposure

Check if the NodePort service is exposing your app to your local machine:

```console
$ kubectl get services
```

You should see something like:

```shell
NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
vuejs-sample-service     NodePort    10.98.233.59    <none>        8080:30001/TCP   1m
```

This output confirms that your app is available via NodePort on port 30001.

### Step 4. Access your app in the browser

Open your browser and navigate to [http://localhost:30001](http://localhost:30001).

You should see your production-ready Vue.js Sample application running — served by your local Kubernetes cluster.

### Step 5. Clean up Kubernetes resources

Once you're done testing, you can delete the deployment and service using:

```console
  $ kubectl delete -f vuejs-sample-kubernetes.yaml
```

Expected output:

```shell
  deployment.apps "vuejs-sample" deleted
  service "vuejs-sample-service" deleted
```

This ensures your cluster stays clean and ready for the next deployment.
   
---

## Summary

In this section, you learned how to deploy your Vue.js application to a local Kubernetes cluster using Docker Desktop. This setup allows you to test and debug your containerized app in a production-like environment before deploying it to the cloud.

What you accomplished:

- Created a Kubernetes Deployment and NodePort Service for your Vue.js app  
- Used `kubectl apply` to deploy the application locally  
- Verified the app was running and accessible at `http://localhost:30001`  
- Cleaned up your Kubernetes resources after testing

---

## Related resources

Explore official references and best practices to sharpen your Kubernetes deployment workflow:

- [Kubernetes documentation](https://kubernetes.io/docs/home/) – Learn about core concepts, workloads, services, and more.  
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/features/kubernetes.md) – Use Docker Desktop’s built-in Kubernetes support for local testing and development.
- [`kubectl` CLI reference](https://kubernetes.io/docs/reference/kubectl/) – Manage Kubernetes clusters from the command line.  
- [Kubernetes Deployment resource](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) – Understand how to manage and scale applications using Deployments.  
- [Kubernetes Service resource](https://kubernetes.io/docs/concepts/services-networking/service/) – Learn how to expose your application to internal and external traffic.