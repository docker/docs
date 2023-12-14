---
title: Test your Node.js deployment
keywords: deploy, kubernetes, node, node.js
description: Learn how to deploy locally to test and debug your Kubernetes deployment
---

## Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a Node.js application](containerize.md).
- [Turn on Kubernetes](/desktop/kubernetes/#turn-on-kubernetes) in Docker Desktop.

## Overview

In this section, you'll learn how to use Docker Desktop to deploy your
application to a fully-featured Kubernetes environment on your development
machine. This allows you to test and debug your workloads on Kubernetes locally
before deploying.

## Create a Kubernetes YAML file

In the cloned repository's directory, create a file named
`docker-node-kubernetes.yaml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your Node.js application](configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-nodejs-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      todo: web
  template:
    metadata:
      labels:
        todo: web
    spec:
      containers:
      - name: todo-site
        image: DOCKER_USERNAME/REPO_NAME
        imagePullPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: todo-entrypoint
  namespace: default
spec:
  type: NodePort
  selector:
    todo: web
  ports:
  - port: 3000
    targetPort: 3000
    nodePort: 30001
```

In this Kubernetes YAML file, there are two objects, separated by the `---`:

 - A Deployment, describing a scalable group of identical pods. In this case,
   you'll get just one replica, or copy of your pod. That pod, which is
   described under `template`, has just one container in it. The container is
   created from the image built by GitHub Actions in [Configure CI/CD for your
   Node.js application](configure-ci-cd.md).
 - A NodePort service, which will route traffic from port 30001 on your host to
   port 3000 inside the pods it routes to, allowing you to reach your app
   from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

## Deploy and check your application

1. In a terminal, navigate to where you created `docker-node-kubernetes.yaml`
   and deploy your application to Kubernetes.

   ```console
   $ kubectl apply -f docker-node-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```shell
   deployment.apps/docker-nodejs-demo created
   service/todo-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
   docker-nodejs-demo   1/1     1            1           6s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME              TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
   kubernetes        ClusterIP   10.96.0.1        <none>        443/TCP          7d22h
   todo-entrypoint   NodePort    10.111.101.229   <none>        3000:30001/TCP   33s
   ```

   In addition to the default `kubernetes` service, you can see your `todo-entrypoint` service, accepting traffic on port 30001/TCP.

3. Open a browser and visit your app at `localhost:30001`. You should see your
   application.

4. Run the following command to tear down your application.

   ```console
   $ kubectl delete -f docker-node-kubernetes.yaml
   ```
5. If step 3 does not load, use `kubectl port-forward` command to instruct Kubernetes to set up port forwarding to route traffice from port `30001` on your local machine to port `3000` on the Kubernetes cluster.
   ```console
   $ kubectl port-forward docker-nodejs-demo-95ddcb6f4-qmb2d 30001:3000
   ```

## Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:
   - [Kubernetes documentation](https://kubernetes.io/docs/home/)
   - [Deploy on Kubernetes with Docker Desktop](../../desktop/kubernetes.md)
   - [Swarm mode overview](../../engine/swarm/_index.md)
