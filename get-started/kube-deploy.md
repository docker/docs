---
title: "Deploy to Kubernetes"
keywords: kubernetes, pods, deployments, kubernetes services
description: Learn how to describe and deploy a simple application on Kubernetes.
---

## Prerequisites

- Download and install Docker Desktop as described in [Get Docker](../get-docker.md).
- Work through containerizing an application in [Part 2](02_our_app.md).
- Make sure that Kubernetes is enabled on your Docker Desktop:
  - **Mac**: Click the Docker icon in your menu bar, navigate to **Preferences** and make sure there's a green light beside 'Kubernetes'.
  - **Windows**: Click the Docker icon in the system tray and navigate to **Settings** and make sure there's a green light beside 'Kubernetes'.

  If Kubernetes isn't running, follow the instructions in [Orchestration](orchestration.md) of this tutorial to finish setting it up.

## Introduction

Now that we've demonstrated that the individual components of our application run as stand-alone containers, it's time to arrange for them to be managed by an orchestrator like Kubernetes. Kubernetes provides many tools for scaling, networking, securing and maintaining your containerized applications, above and beyond the abilities of containers themselves.

In order to validate that our containerized application works well on Kubernetes, we'll use Docker Desktop's built in Kubernetes environment right on our development machine to deploy our application, before handing it off to run on a full Kubernetes cluster in production. The Kubernetes environment created by Docker Desktop is _fully featured_, meaning it has all the Kubernetes features your app will enjoy on a real cluster, accessible from the convenience of your development machine.

## Describing apps using Kubernetes YAML

All containers in Kubernetes are scheduled as _pods_, which are groups of co-located containers that share some resources. Furthermore, in a realistic application we almost never create individual pods; instead, most of our workloads are scheduled as _deployments_, which are scalable groups of pods maintained automatically by Kubernetes. Lastly, all Kubernetes objects can and should be described in manifests called _Kubernetes YAML_ files. These YAML files describe all the components and configurations of your Kubernetes app, and can be used to easily create and destroy your app in any Kubernetes environment.

1.  You already wrote a very basic Kubernetes YAML file in the Orchestration overview part of this tutorial. Now, let's write a slightly more sophisticated YAML file to run and manage our Todo app, the container `getting-started` image created in [Part 2](02_our_app.md) of the Quickstart tutorial. Place the following in a file called `bb.yaml`:

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: bb-demo
      namespace: default
    spec:
      replicas: 1
      selector:
        matchLabels:
          bb: web
      template:
        metadata:
          labels:
            bb: web
        spec:
          containers:
          - name: bb-site
            image: getting-started
            imagePullPolicy: Never
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: bb-entrypoint
      namespace: default
    spec:
      type: NodePort
      selector:
        bb: web
      ports:
      - port: 3000
        targetPort: 3000
        nodePort: 30001
    ```

    In this Kubernetes YAML file, we have two objects, separated by the `---`:
    - A `Deployment`, describing a scalable group of identical pods. In this case, you'll get just one `replica`, or copy of your pod, and that pod (which is described under the `template:` key) has just one container in it, based off of your `getting-started` image from the previous step in this tutorial.
    - A `NodePort` service, which will route traffic from port 30001 on your host to port 3000 inside the pods it routes to, allowing you to reach your Todo app from the network.

    Also, notice that while Kubernetes YAML can appear long and complicated at first, it almost always follows the same pattern:
    - The `apiVersion`, which indicates the Kubernetes API that parses this object
    - The `kind` indicating what sort of object this is
    - Some `metadata` applying things like names to your objects
    - The `spec` specifying all the parameters and configurations of your object.

## Deploy and check your application

1.  In a terminal, navigate to where you created `bb.yaml` and deploy your application to Kubernetes:

    ```console
    $ kubectl apply -f bb.yaml
    ```

    you should see output that looks like the following, indicating your Kubernetes objects were created successfully:

    ```shell
    deployment.apps/bb-demo created
    service/bb-entrypoint created
    ```

2.  Make sure everything worked by listing your deployments:

    ```console
    $ kubectl get deployments
    ```

    if all is well, your deployment should be listed as follows:

    ```shell
    NAME      READY   UP-TO-DATE   AVAILABLE   AGE
    bb-demo   1/1     1            1           40s
    ```

    This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services:

    ```console
    $ kubectl get services

    NAME            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
    bb-entrypoint   NodePort    10.106.145.116   <none>        3000:30001/TCP   53s
    kubernetes      ClusterIP   10.96.0.1        <none>        443/TCP          138d
    ```

    In addition to the default `kubernetes` service, we see our `bb-entrypoint` service, accepting traffic on port 30001/TCP.

3.  Open a browser and visit your Todo app at `localhost:30001`; you should see your Todo application, the same as when we ran it as a stand-alone container in [Part 2](02_our_app.md) of the Quickstart tutorial.

4.  Once satisfied, tear down your application:

    ```console
    $ kubectl delete -f bb.yaml
    ```

## Conclusion

At this point, we have successfully used Docker Desktop to deploy our application to a fully-featured Kubernetes environment on our development machine. We haven't done much with Kubernetes yet, but the door is now open; you can begin adding other components to your app and taking advantage of all the features and power of Kubernetes, right on your own machine.

In addition to deploying to Kubernetes, we have also described our application as a Kubernetes YAML file. This simple text file contains everything we need to create our application in a running state. We can check it into version control and share it with our colleagues, allowing us to distribute our applications to other clusters (like the testing and production clusters that probably come after our development environments) easily.

## Kubernetes references

Further documentation for all new Kubernetes objects used in this article are available here:

 - [Kubernetes Pods](https://kubernetes.io/docs/concepts/workloads/pods/pod/)
 - [Kubernetes Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
 - [Kubernetes Services](https://kubernetes.io/docs/concepts/services-networking/service/)
