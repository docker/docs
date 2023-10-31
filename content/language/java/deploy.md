---
title: Test your Java deployment
keywords: deploy, kubernetes, java
description: Learn how to develop locally using Kubernetes
---

## Prerequisites

- Complete all the previous sections of this guide, starting with [Build your Java image](build-images.md).
- [Turn on Kubernetes](/desktop/kubernetes/#turn-on-kubernetes) in Docker Desktop.

## Overview

In this section, you'll learn how to use Docker Desktop to deploy your
application to a fully-featured Kubernetes environment on your development
machine. This lets you test and debug your workloads on Kubernetes locally
before deploying.

## Create a Kubernetes YAML file

In your `spring-petclinic` directory, create a file named
`docker-java-kubernetes.yaml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your Java application](configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-java-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: server
  template:
    metadata:
      labels:
        service: server
    spec:
      containers:
       - name: server-service
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
    service: server
  ports:
  - port: 8080
    targetPort: 8080
    nodePort: 30001
```

In this Kubernetes YAML file, there are two objects, separated by the `---`:

 - A Deployment, describing a scalable group of identical pods. In this case,
   you'll get just one replica, or copy of your pod. That pod, which is
   described under `template`, has just one container in it. The
    container is created from the image built by GitHub Actions in [Configure CI/CD for
    your Java application](configure-ci-cd.md).
 - A NodePort service, which will route traffic from port 30001 on your host to
   port 8080 inside the pods it routes to, allowing you to reach your app
   from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

## Deploy and check your application

1. In a terminal, navigate to `spring-petclinic` and deploy your application to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-java-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```shell
   deployment.apps/docker-java-demo created
   service/service-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
   docker-java-demo     1/1     1            1           15s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
   kubernetes           ClusterIP   10.96.0.1       <none>        443/TCP          23h
   service-entrypoint   NodePort    10.99.128.230   <none>        8080:30001/TCP   75s
   ```

   In addition to the default `kubernetes` service, you can see your `service-entrypoint` service, accepting traffic on port 30001/TCP.

3. In a terminal, curl the service. Note that a database wasn't deployed in
   this example.

   ```console
   $ curl --request GET \
     --url http://localhost:30001/actuator/health \
     --header 'content-type: application/json'
   ```

   You should get output like the following.
   ```console
   {"status":"UP","groups":["liveness","readiness"]}
   ```

4. Run the following command to tear down your application.

   ```console
   $ kubectl delete -f docker-java-kubernetes.yaml
   ```

## Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:
   - [Kubernetes documentation](https://kubernetes.io/docs/home/)
   - [Deploy on Kubernetes with Docker Desktop](../../desktop/kubernetes.md)
   - [Swarm mode overview](../../engine/swarm/_index.md)