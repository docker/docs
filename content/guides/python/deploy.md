---
title: Test your Python deployment
linkTitle: Test your deployment
weight: 50
keywords: deploy, kubernetes, python
description: Learn how to develop locally using Kubernetes
aliases:
  - /language/python/deploy/
  - /guides/language/python/deploy/
---

## Prerequisites

- Complete all the previous sections of this guide, starting with [Use containers for Python development](develop.md).
- [Turn on Kubernetes](/manuals/desktop/features/kubernetes.md#install-and-turn-on-kubernetes) in Docker Desktop.

## Overview

In this section, you'll learn how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine. This allows you to test and debug your workloads on Kubernetes locally before deploying.

## Create a Kubernetes YAML file

In your `python-docker-dev-example` directory, create a file named `docker-postgres-kubernetes.yaml`. Open the file in an IDE or text editor and add
the following contents.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              value: example
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: POSTGRES_PASSWORD
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: default
spec:
  ports:
    - port: 5432
  selector:
    app: postgres
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: Secret
metadata:
  name: postgres-secret
  namespace: default
type: Opaque
data:
  POSTGRES_PASSWORD: cG9zdGdyZXNfcGFzc3dvcmQ= # Base64 encoded password (e.g., 'postgres_password')
```

In your `python-docker-dev-example` directory, create a file named
`docker-python-kubernetes.yaml`. Replace `DOCKER_USERNAME/REPO_NAME` with your
Docker username and the repository name that you created in [Configure CI/CD for
your Python application](./configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-python-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: fastapi
  template:
    metadata:
      labels:
        service: fastapi
    spec:
      containers:
        - name: fastapi-service
          image: DOCKER_USERNAME/REPO_NAME
          imagePullPolicy: Always
          env:
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-secret
                  key: POSTGRES_PASSWORD
            - name: POSTGRES_USER
              value: postgres
            - name: POSTGRES_DB
              value: example
            - name: POSTGRES_SERVER
              value: postgres
            - name: POSTGRES_PORT
              value: "5432"
          ports:
            - containerPort: 8001
---
apiVersion: v1
kind: Service
metadata:
  name: service-entrypoint
  namespace: default
spec:
  type: NodePort
  selector:
    service: fastapi
  ports:
    - port: 8001
      targetPort: 8001
      nodePort: 30001
```

In these Kubernetes YAML file, there are various objects, separated by the `---`:

- A Deployment, describing a scalable group of identical pods. In this case,
  you'll get just one replica, or copy of your pod. That pod, which is
  described under `template`, has just one container in it. The
  container is created from the image built by GitHub Actions in [Configure CI/CD for
  your Python application](configure-ci-cd.md).
- A Service, which will define how the ports are mapped in the containers.
- A PersistentVolumeClaim, to define a storage that will be persistent through restarts for the database.
- A Secret, Keeping the database password as a example using secret kubernetes resource.
- A NodePort service, which will route traffic from port 30001 on your host to
  port 8001 inside the pods it routes to, allowing you to reach your app
  from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

> [!NOTE]
>
> - The `NodePort` service is good for development/testing purposes. For production you should implement an [ingress-controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/).

## Deploy and check your application

1. In a terminal, navigate to `python-docker-dev-example` and deploy your database to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-postgres-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```console
   deployment.apps/postgres created
   service/postgres created
   persistentvolumeclaim/postgres-pvc created
   secret/postgres-secret created
   ```

   Now, deploy your python application.

   ```console
   kubectl apply -f docker-python-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```console
   deployment.apps/docker-python-demo created
   service/service-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```console
   NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
   docker-python-demo   1/1     1            1           48s
   postgres             1/1     1            1           2m39s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```console
   NAME                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
   kubernetes           ClusterIP   10.43.0.1      <none>        443/TCP          13h
   postgres             ClusterIP   10.43.209.25   <none>        5432/TCP         3m10s
   service-entrypoint   NodePort    10.43.67.120   <none>        8001:30001/TCP   79s
   ```

   In addition to the default `kubernetes` service, you can see your `service-entrypoint` service, accepting traffic on port 30001/TCP and the internal `ClusterIP` `postgres` with the port `5432` open to accept connections from you python app.

3. In a terminal, curl the service. Note that a database was not deployed in
   this example.

   ```console
   $ curl http://localhost:30001/
   Hello, Docker!!!
   ```

4. Run the following commands to tear down your application.

   ```console
   $ kubectl delete -f docker-python-kubernetes.yaml
   $ kubectl delete -f docker-postgres-kubernetes.yaml
   ```

## Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:

- [Kubernetes documentation](https://kubernetes.io/docs/home/)
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/features/kubernetes.md)
- [Swarm mode overview](/manuals/engine/swarm/_index.md)
