---
title: Test your Deno deployment
linkTitle: Test your deployment
weight: 50
keywords: deploy, kubernetes, deno
description: Learn how to develop locally using Kubernetes
aliases:
- /language/deno/deploy/
---

## Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a Deno application](containerize.md).
- [Turn on Kubernetes](/manuals//desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

## Overview

In this section, you'll learn how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine. This allows you to test and debug your workloads on Kubernetes locally before deploying.

## Create a Kubernetes YAML file

In your `deno-docker` directory, create a file named
`docker-kubernetes.yml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your Deno application](configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-deno-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: deno-api
  template:
    metadata:
      labels:
        app: deno-api
    spec:
      containers:
       - name: deno-api
         image: DOCKER_USERNAME/REPO_NAME
         imagePullPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: service-entrypoint
  namespace: default
spec:
  type: NodePort
  selector:
    app: deno-api
  ports:
  - port: 8000
    targetPort: 8000
    nodePort: 30001
```

In this Kubernetes YAML file, there are two objects, separated by the `---`:

 - A Deployment, describing a scalable group of identical pods. In this case,
   you'll get just one replica, or copy of your pod. That pod, which is
   described under `template`, has just one container in it. The
    container is created from the image built by GitHub Actions in [Configure CI/CD for
    your Deno application](configure-ci-cd.md).
 - A NodePort service, which will route traffic from port 30001 on your host to
   port 8000 inside the pods it routes to, allowing you to reach your app
   from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

## Deploy and check your application

1. In a terminal, navigate to `deno-docker` and deploy your application to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-kubernetes.yml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```text
   deployment.apps/docker-deno-demo created
   service/service-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                 READY   UP-TO-DATE   AVAILABLE    AGE
   docker-deno-demo       1/1     1            1           10s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
   kubernetes           ClusterIP   10.96.0.1        <none>        443/TCP          88m
   service-entrypoint   NodePort    10.105.145.223   <none>        8000:30001/TCP   83s
   ```

   In addition to the default `kubernetes` service, you can see your `service-entrypoint` service, accepting traffic on port 30001/TCP.

3. In a browser, visit the following address. You should see the message `{"Status" : "OK"}`.

   ```console
   http://localhost:30001/
   ```

4. Run the following command to tear down your application.

   ```console
   $ kubectl delete -f docker-kubernetes.yml
   ```

## Summary

In this section, you learned how to use Docker Desktop to deploy your Deno application to a fully-featured Kubernetes environment on your development machine. 

Related information:
   - [Kubernetes documentation](https://kubernetes.io/docs/home/)
   - [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
   - [Swarm mode overview](/manuals/engine/swarm/_index.md)
