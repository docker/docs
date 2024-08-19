---
title: Test your .NET deployment
keywords: deploy, .net, local, development
description: Learn how to deploy your application
---

## Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize
  a .NET application](containerize.md).
- [Turn on Kubernetes](/desktop/kubernetes/#install-and-turn-on-kubernetes) in Docker
  Desktop.

## Overview

In this section, you'll learn how to use Docker Desktop to deploy your
application to a fully-featured Kubernetes environment on your development
machine. This allows you to test and debug your workloads on Kubernetes locally
before deploying.

## Create a Kubernetes YAML file

In your `docker-dotnet-sample` directory, create a file named
`docker-dotnet-kubernetes.yaml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your .NET application](configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    service: server
  name: server
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: server
  strategy: {}
  template:
    metadata:
      labels:
        service: server
    spec:
      initContainers:
        - name: wait-for-db
          image: busybox:1.28
          command: ['sh', '-c', 'until nc -zv db 5432; do echo "waiting for db"; sleep 2; done;']
      containers:
        - image: DOCKER_USERNAME/REPO_NAME
          name: server
          imagePullPolicy: Always
          ports:
            - containerPort: 80
              hostPort: 8080
              protocol: TCP
          resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    service: db
  name: db
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: db
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        service: db
    spec:
      containers:
        - env:
            - name: POSTGRES_DB
              value: example
            - name: POSTGRES_PASSWORD
              value: example
          image: postgres
          name: db
          ports:
            - containerPort: 5432
              protocol: TCP
          resources: {}
      restartPolicy: Always
status: {}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    service: server
  name: server
  namespace: default
spec:
  type: NodePort
  ports:
    - name: "8080"
      port: 8080
      targetPort: 80
      nodePort: 30001
  selector:
    service: server
status:
  loadBalancer: {}
---
apiVersion: v1
kind: Service
metadata:
  labels:
    service: db
  name: db
  namespace: default
spec:
  ports:
    - name: "5432"
      port: 5432
      targetPort: 5432
  selector:
    service: db
status:
  loadBalancer: {}
```

In this Kubernetes YAML file, there are four objects, separated by the `---`. In addition to a Service and Deployment for the database, the other two objects are:

 - A Deployment, describing a scalable group of identical pods. In this case,
   you'll get just one replica, or copy of your pod. That pod, which is
   described under `template`, has just one container in it. The container is
    created from the image built by GitHub Actions in [Configure CI/CD for your
    .NET application](configure-ci-cd.md).
 - A NodePort service, which will route traffic from port 30001 on your host to
   port 8080 inside the pods it routes to, allowing you to reach your app
   from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

## Deploy and check your application

1. In a terminal, navigate to the `docker-dotnet-sample` directory
   and deploy your application to Kubernetes.

   ```console
   $ kubectl apply -f docker-dotnet-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```shell
   deployment.apps/db created
   service/db created
   deployment.apps/server created
   service/server created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME     READY   UP-TO-DATE   AVAILABLE   AGE
   db       1/1     1            1           76s
   server   1/1     1            1           76s
   ```

   This indicates all of the pods are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
   db           ClusterIP   10.96.156.90    <none>        5432/TCP         2m8s
   kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP          164m
   server       NodePort    10.102.94.225   <none>        8080:30001/TCP   2m8s
   ```

   In addition to the default `kubernetes` service, you can see your `server` service and `db` service. The `server` service is accepting traffic on port 30001/TCP.

3. Open a browser and visit your app at `localhost:30001`. You should see your
   application.

4. Run the following command to tear down your application.

   ```console
   $ kubectl delete -f docker-dotnet-kubernetes.yaml
   ```

## Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:
   - [Kubernetes documentation](https://kubernetes.io/docs/home/)
   - [Deploy on Kubernetes with Docker Desktop](../../desktop/kubernetes.md)
   - [Swarm mode overview](../../engine/swarm/_index.md)
