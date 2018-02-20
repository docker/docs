---
description: Migrating from Docker Cloud
keywords: cloud, swarm, migration
title: Migrate Docker Cloud stack to Kubernetes on AKS
---

## Introduction

Docker Cloud applications are defined in declarative configuration files and run on cloud instances in the public cloud. Kubernetes applications are defined in declarative manifest files and can also run on cloud instances in the public cloud.

This document will guide you through the process of converting your Docker Cloud application into a Kubernetes application, and then running it on a Microsoft's Azure Container Service (AKS) cluster.

AKS is a hosted Kubernetes service on the Microsoft Azure public cloud. It exposes all of the standard Kubernetes API endpoints, meaning you can use standard Kubernetes tools and manifests. This means that the manifest files you create can be used to run your application on any Kubernetes infrastructure --- different public clouds and on-premises.

The high-level process will involve:

- Building a Kubernetes cluster on AKS (this is equivalent of your *node cluster* in Docker Cloud)
- Converting your Docker Cloud stack into a Kubernetes application defined in a Kubernetes manifest file
- Testing the Kubernetes application on AKS
- Deploying the Kubernetes application on AKS

The actual process of performing the migration (switching customers from using your Docker Cloud application to using your AKS application) will vary considerably between applications, and is therefore outside the scope of this document. This document will walk you through the process of converting your Docker Cloud app into a Kubernetes app, and testing it works. It will not outline migration steps.

## Things that will be different after migration

> %% I need to dig around and find out what AKS does and does not offer before we can document what will be lost. It's currently in public beta.

Many things about your application will stay the same. For example, your Docker images and overall application functionality will not change --- if your application uses a Docker image called `myorg/webfe:v3` and is accessible on port `80`, none of this will change. Also, the way your customers and other applications interact with your application will not change.

## Pre-reqs

To complete the migration, you will need:

- An active Azure subscription with billing enabled.

## Sample Application

We'll use the following sample application to walk you through the process of converting your application from Docker Cloud to Kubernetes.

https://github.com/dockersamples/example-voting-app/blob/master/dockercloud.yml

The Docker Cloud stackfile (`dockercloud.yml`) defines the application. As we can see, it is comprised of 6 application services:

- **vote:** Web front-end displaying voting options
- **redis:** In-memory k/v store that collects votes
- **worker:** Stores votes in database
- **db:** Persistent store for votes
- **result:** Web server that pulls and displays results from database
- **lb:** Container-based load balancer

The application accepts votes via the `vote` service and stores them in a persistent backend database (`db`). The `lb`, `redis`, and `worker` services assist in the process, and you can view results of the vote via the `result` service. The high-level architecture of the application is shown below.

![]()  < put image here.

Your application will be different, but the overall process of converting it into a Kubernetes application will be the same.

## Build target environment

Azure Container Service (AKS) is a managed Kubernetes service. As such, Azure takes care of all of the Kubernetes control plane management --- delivering the control plane APIs, managing control plane HA, managing control plane upgrades etc. You only need to look after worker nodes --- how many, size & spec, where to deploy them etc.

We'll complete the following three high-level steps as part of building your Azure AKS cluster:

1. Create an Azure Application Registration
2. Build an AKS cluster
3. Connect to the AKS cluster

### Create and Azure Application Registration

1. Log in to the [Azure portal](https://portal.azure.com).

2. Click `Azure Active Directory` > `App registrations` > `New application registration`.

3. Give it a `Name`, select `Web app/API` as the `Application type`, and enter a `Sign-on URL`

    The sign-on URL needs to be a valid DNS name but does not need to be resolvable. An example might be `https://k8s-vote.com`


4. Click `Create`.

5. Copy the Application ID to a safe place you will need it in a later step.

6. Click `Settings` > `Keys` and set a description and duration.

7. Click `Save`.

8. Copy contents of the `Value` field as this is a password that you will need in a later step.

The application registration is now complete. Time to build your AKS cluster.

### Build an AKS cluster

In this section, we'll create an Azure AKS cluster. We'll build a 3-node cluster in the `West Europe` Region. Your cluster will be different, and should probably be based on the configuration of your Docker Cloud *node cluster*.

You can use one of the following methods to see the configuration of your Docker Cloud *node cluster*:

- **Docker Cloud classic mode users:** Log-in to Docker Cloud and select `Node Clusters` under `INFRASTRUCTURE`. From here you can select any of your node clusters to find information such as; cloud provider, number of nodes, spec of nodes, region and availability zone.
- **Docker Cloud Swarm mode users:** Log-in to Docker Cloud and select `Swarms` from the top menu. Select one of your Swarms.

Before continuing, you will need to know all of the following:

- Your Azure subscription
- Azure *Region* that you want to deploy your AKS cluster to
- An SSH public key to use when connecting to AKS nodes
- The *Service principal ID* and *Service Principal secret* from the previous section
- Number of nodes
- Size and spec of nodes


1. Select `+New` from the Azure portal main dashboard.

2. Select `Containers` and choose `Azure Container Service - AKS (preview)`. **Do not select the other ACS option.**

3. Fill out the form with the required details.

    - **Cluster name:** This is the name you want to use to identify this AKS cluster.
    - **Kubernetes version:** Choose one of the 1.8.x versions.
    - **Subscription:** Choose the subscription that you want to use to pay for the cluster.
    - **Resource group:** Create a new resource group or choose one from your existing list.
    - **Location:** Select the Azure region that you want to deploy the AKS cluster to. AKS may not be available in all Azure locations.

4. Click `OK`.

5. Configure the additional AKS cluster parameters.

    - **User name:** The default option should be fine.
    - **SSH public key:** This is the public key (certificate) of a key-pair that you own that can be used for SSH. You can generate this using a number of different tools such as PuTTY. It should be a minimum of 2048 bits of type ssh-rsa.
    - **Service principal client ID:** This is the application ID that you copied in an earlier step.
    - **Service principal client secret:** This is password value that you copied in a previous step.
    - **Node count:** This is the number of nodes that you want in the cluster. It should probably match the number of nodes in your existing Docker Cloud node cluster.
    - **Node virtual machine size:** This is the size and specification of each AKS node. It should probably match the configuration of your existing Docker Cloud node cluster.

6. Click `OK`.

7. Review the configuration on the Summary screen and click `OK` to deploy the cluster.

It will take a few minutes for the cluster to deploy.

### Connect to the AKS cluster

You can use the web-based Azure cloud shell to connect to your cluster. However, this section will show you how to configure your laptop, or other local terminal, to connect to your cluster.

We'll complete the following high-level steps:

- Install the Azure `az` CLI tool
- Install `kubectl` --- the Kubernetes command line tool
- Configure `kubectl` to connect to your cluster.

Let's install the Azure CLI tool.

1. Click this [link](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) and download the Azure CLI tool for your relevant Operating System.

2. Install the Azure CLI.

Use the Azure CLI tool to install the Kubernete command line tool (`kubectl`).

```
> az aks install-cli
Downloading client to C:\Program Files (x86)\kubectl.exe from...
```

Now that `kubectl` is installed you need to configure it to connect to your AKS cluster.

1. Start the Azure login process.

    ```
    > az login
    To sign in, use a web browser to open the page https://aka.ms/devicelogin and enter...
    ```

2. Open the devicelogin page in a browser tab and paste in the authentication code.

    When complete it will return some JSON.

3. Get the credentials and use them to configure `kubectl`

    ```
    > az aks get-credentials --resource-group=k8s-vote --name=k8s-vote
    Merged "k8s-vote" as current context in C:\Users\nigel\.kube\config
    ```

4. Test that `kubectl` can connect to your cluster.

    ```
    > kubectl get nodes
    NAME                       STATUS    ROLES     AGE       VERSION
    aks-agentpool-29046111-0   Ready     agent     3m        v1.8.1
    aks-agentpool-29046111-1   Ready     agent     2m        v1.8.1
    aks-agentpool-29046111-2   Ready     agent     2m        v1.8.1
    ```

    The `kubectl get nodes` command lists information about the nodes in the cluster it is configured to communicate with. If the values  returned by the command match your AKS cluster (number of nodes, age, and version) then you have successfully configured `kubectl` to manage your AKS cluster.

You now have an AKS cluster and have configured the `kubectl` Kubernetes client to manage the cluster. Let's look at how to convert your Docker Cloud app into a Kubernetes app.

## Converting the application

Your application is currently deployed as a Docker Cloud stack on the Docker Cloud platform. The steps in this section will guide you through the process of converting your application into a Kubernetes application (one that can be deployed on a Kubernetes cluster).

We'll use a sample application to demonstrate the process of converting a Docker Cloud stack into a Kubernetes application. Your application will be different, but the overall process will be the same.

Your Docker Cloud application probably comprises several services defined in a single Docker Cloud stackfile called `dockercloud.yml`. Kubernetes applications are defined in *manifest* files that are also YAML. However, Kubernetes manifest files tend to be longer.

We'll refer to your existing Docker cloud stackfile as the **source** file, and the new Kubernetes manifest file as the **target**.

Your **source** Docker Cloud application (stack) will be defined in a YAML file. This is the Docker Cloud stackfile (`dockercloud.yml`) from the example app. Yours will be different, and you can find it in the Docker Cloud web UI by selecting your stack and clicking `Edit`.

The example Docker Cloud application comprises six Docker services, each of which is defined in the Docker Cloud stackfile as a top-level key:

```
db:
redis:
result:
lb:
vote:
worker:
```

Let's step through the process of converting each one.

We'll deploy the application on Kubernetes as a set of Kubernetes *Deployments* and Kubernetes *Services*. We'll create one *Deployment* and *Service* per Docker Cloud stack service.

A Kubernetes *Deployment* defines the application service. This includes things such as; which Docker image to use, and which container ports to map. You can also define how rolling updates work, rollbacks, and may other advanced features. However, these are beyond the scope of this document.

A Kubernetes *Service* provides in an abstraction that provides stable networking for a set of *Pods* defined in a *Deployment*. It also allows you to create cloud-native load-balancers.

### Converting the `db` service

The sample Docker Cloud stackfile defines an image and a restart policy for the `db` service.

```
db:
  image: 'postgres:9.4'
  restart: always
```

This can be represented in a Kubernetes manifest as follows:

```
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: db
  labels:
    app: db
spec:
  template:
    metadata:
      labels:
        app: db
    spec:
      containers:
      - image: postgres:9.4
        name: db
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: db
spec:
  clusterIP: None
  ports:
  - port: 55555
    targetPort: 0
  selector:
    app: db
```

Two things to note:

1. This is a lot longer than the Docker Cloud definition
2. The manifest defines a *Deployment* and a *Service* separated by three dashes `---`.

Every Kubernetes object (*Deployments*, *Services* etc.) needs to specify an `apiVersion` and `kind`.

The `apiVersion` tells Kubernetes which schema version to use when defining and managing the object. The versions used in the examples will work on all versions of Kubernetes currently supported on AKS (1.7.7 and 1.8.1).

The `kind` field tells Kubernetes what type of object is being defined. This guide will only define *Deployments* and *Services*.

The `metadata` section gives the object a name and a set of labels. The example is naming and labelling as everything as `db` to keep it in line with the Docker Cloud `db` service.

The `spec` sections is where the configuration of the object is defined.

In the *Deployment* section, it defines the Kubernetes *Pods* to deploy. Think of Kubernetes *Pods* as containers. In this example, we're defining a container called `db` based on the `postgres:9.4` Docker image, and defining a restart policy. We're also labelling all *Pods* (containers) with the `app=db` label.

It's important that the *Pod* (container) labels (`Deployment.spec.template.metadata.labels`) match the *Deployment* labels (`Deployment.metadata.labels`) as this is how the deployment knows which *Pods* on the cluster to manage.

It's possible to define a lot more as part of a Kubernetes *Deployment*, but the options we've configured are enough to replicate what was defined for the `db` service on Docker Cloud.

We've also defined a Kubernetes *Service* for the `db` service.

> %% May be add something explaining that the Service's label selector doesn't have to match every Pod label. E.g. Pods can have more labels, all that's required for the Service to manage them is that the Pods have all of the labels that the Service is selecting on

By naming the *Service* "db", we've ensured that a cluster-wide DNS mapping has been created for `db` pointing to this *Service*. We've then defined a label selector (`Service.spec.selector`) with a value of "app=db". This means that the service will provide stable networking a load-balancing for all *Pods* (containers) on the cluster with matching labels. Notice that the *Pods* defined in the *Deployment* section of the file are all labelled as "app-db". It is this mapping between the *Service's* label selector and the *Pod* labels that tell the *Service* object which Pods to provide networking for.


### Converting the `redis` service

The Docker Cloud stackfile defines an image and a restart policy for the `redis` service.

```
redis:
  image: 'redis:latest'
  restart: always
```

This can be represented in a Kubernetes manifest as follows:

```
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: redis
  name: redis
spec:
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - image: redis:alpine
        name: redis
        ports:
        - containerPort: 6379
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: redis
  name: redis
spec:
  ports:
  - port: 6379
    targetPort: 6379
  selector:
    app: redis
```

This is very similar to the `db` service.

The *Deployment* section names and labels the the *Deployment* as "redis", deploys a *Pod* (container) called "redis" from the `redis:alpine` image, and sets the container port to be `6379`. It also makes Sure the *Pod* labels match the *Deployment* labels to tie the two together.

The *Service* section of the file creates a Kubernetes *Service* and cluster-wide DNS mapping for the name "redis" on port 6379. This maps traffic for `tcp://redis:6379` will be routed to this *Service* and load-balanced across all *Pods* on the cluster with the "app=redis" label.

Again, it's vital that the *Service's* label selector (`Service.spec.selector`) matches the labels assigned to the `redis` *Pods*.

### Converting the `lb` service

The Docker Cloud stackfile defines an `lb` service to load balance traffic to the vote service. This is not needed in Kubernetes on AKS, because Kubernetes lets you define a *Service* object with `type=loadbalancer` that creates a native Azure load-balancer to do this job. You'll see it in the next section.

### Converting the `vote` service

The Docker Cloud stackfile defines an image, a restart policy, and a specific number of containers for the `vote` service. It also defines the Docker Cloud `autoredeploy` feature.

The Docker Cloud `autoredeploy` feature is no supported in Azure AKS....  <<  %% Need to look into this.

```
vote:
  autoredeploy: true
  image: 'docker/example-voting-app-vote:latest'
  restart: always
  target_num_containers: 5
```

This can be represented in a Kubernetes manifest as follows:

```
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: vote
  name: vote
spec:
  replicas: 5
  template:
    metadata:
      labels:
        app: vote
    spec:
      containers:
      - image: docker/example-voting-app-vote:latest
        name: vote
        ports:
        - containerPort: 80
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: vote
  name: vote
spec:
  type: LoadBalancer
  ports:
  - port: 80
  selector:
    app: vote
```

The *Deployment* section of the file defines the required labels and *Pod* (container) specification. Importantly, this one sets the number of *Pod* replicas to 5 (`Deployment.spec.replicas`). This is to match the `target_num_containers` from the Docker Cloud stackfile.

The *Services* section of the file defines the required name and labels to provide stable networking for all *Pods* in the cluster with the "app=vote" label. However, this time it defines the *Service* as "type=loadbalancer". This will create a native Azure load-balancer that will create a stable, publicly routable, IP for the service. It also maps port 80. This means that any traffic hitting the load-balancer's publically routable IP on port 80 will be load-balanced across all 5 *Pod replicas* in the cluster.  This is why the `lb` service from the Docker Cloud app is not needed.

### Converting the `worker` service

The definition of the `worker` service in the Docker Cloud stackfile is similar to the `vote` service. It defines an image, a restart policy, and a specified number of containers (replicas). It also defines the Docker Cloud `autoredeploy` policy. As previously mentioned...... <<  %% Need info on AKS equivalent of autoredeploy.

```
worker:
  autoredeploy: true
  image: 'docker/example-voting-app-worker:latest'
  restart: always
  target_num_containers: 3
```

This can be represented in a Kubernetes manifest as follows:

```
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: worker
  name: worker
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: worker
    spec:
      containers:
      - image: docker/example-voting-app-worker:latest
        name: worker
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: worker
  name: worker
spec:
  clusterIP: None
  ports:
  - port: 55555
    targetPort: 0
  selector:
    app: worker
```

The *Deployment* section defines the required names labels and *Pod* spec.

The *Service* section defines a service that will provide a stable internal networking endpoint for all *Pods* deployed as part of the *Deployment* (Deployment.spec.template.metadata.labels.app=worker).

### Converting the `result` service

The Docker Cloud stackfile defines an image and a restart policy for the `result` service.

```
result:
  autoredeploy: true
  image: 'docker/example-voting-app-result:latest'
  ports:
    - '80:80'
  restart: always
```

This can be represented in a Kubernetes manifest as follows:

```
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: result
  name: result
spec:
  template:
    metadata:
      labels:
        app: result
    spec:
      containers:
      - args:
        - nodemon
        - --debug
        - server.js
        image: docker/example-voting-app-result:latest
        name: result
        ports:
        - containerPort: 80
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: result
  name: result
spec:
  type: LoadBalancer
  ports:
  - port: 80
  selector:
    app: result
```

The *Deployment* section defines the usual names, labels and *Pod* (container) spec.

The *Service* section defines another Azure-native load-balancer to load balance external traffic to the cluster on port 80.

### The Kubernetes manifest file

You can choose to include all *Deployments* and *Services* in a single YAML file, or have one YAML file per Docker Cloud service. The choice is yours, but it's easier to deploy and manage as a single file.

The example below shows them all defined in a single long YMAL file called "k8s-vote.yml". You should manage your Kubernetes manifest files the way you manage your application code --- checking them in and out of version control repositories etc.

```
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: db
  labels:
    app: db
spec:
  template:
    metadata:
      labels:
        app: db
    spec:
      containers:
      - image: postgres:9.4
        name: db
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: db
spec:
  clusterIP: None
  ports:
  - port: 55555
    targetPort: 0
  selector:
    app: db
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: redis
  name: redis
spec:
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - image: redis:alpine
        name: redis
        ports:
        - containerPort: 6379
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: redis
  name: redis
spec:
  ports:
  - port: 6379
    targetPort: 6379
  selector:
    app: redis
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: vote
  name: vote
spec:
  replicas: 5
  template:
    metadata:
      labels:
        app: vote
    spec:
      containers:
      - image: docker/example-voting-app-vote:latest
        name: vote
        ports:
        - containerPort: 80
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: vote
  name: vote
spec:
  type: LoadBalancer
  ports:
  - port: 80
  selector:
    app: vote
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: worker
  name: worker
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: worker
    spec:
      containers:
      - image: docker/example-voting-app-worker:latest
        name: worker
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: worker
  name: worker
spec:
  clusterIP: None
  ports:
  - port: 55555
    targetPort: 0
  selector:
    app: worker
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  labels:
    app: result
  name: result
spec:
  template:
    metadata:
      labels:
        app: result
    spec:
      containers:
      - image: docker/example-voting-app-result:latest
        name: result
        ports:
        - containerPort: 80
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: result
  name: result
spec:
  type: LoadBalancer
  ports:
  - port: 80
  selector:
    app: result
```


Save the Kubernetes manifest file. The examples that follow will assume it's called `k8s-vote.yml`.

## Test the app on AKS

You should thoroughly test your application on your AKS cluster before starting the migration. This includes deploying the application from the new Kubernetes manifest file, performing scaling operations, updates and rollbacks. You should also manage the manifest file in a version control system.

The following steps will show you how to deploy your app from the Kubernetes manifest files and test it.

Perform the following from an Azure cloud shell or local terminal with the azure client and `kubectl` Kubernetes client installed. It will need to be configured to talk to your AKS cluster.

1. Verify that you shell/terminal is configured to talk to your AKS cluster.

    ```
    > kubectl get nodes
    NAME                       STATUS    ROLES     AGE       VERSION
    aks-agentpool-29046111-0   Ready     agent     6h        v1.8.1
    aks-agentpool-29046111-1   Ready     agent     6h        v1.8.1
    aks-agentpool-29046111-2   Ready     agent     6h        v1.8.1
    ```

    If the output matches your cluster you're ready to proceed with the next steps.

2. Deploy your Kubernetes application to your cluster.

    This example assumes the application is defined in a Kubernetes manifest file called `ks8-vote.yml` in your system's PATH. You will need to substitute the name of your manifest file.

    ```
    > kubectl create -f k8s-vote.yml

    deployment "db" created
    service "db" created
    deployment "redis" created
    service "redis" created
    deployment "vote" created
    service "vote" created
    deployment "worker" created
    service "worker" created
    deployment "result" created
    service "result" created
    ```

3. Check the status of the app.

    Use the following two commands to check the status of the *Deployments* and *Services*.

    ```
    > kubectl get deploy
    NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    db        1         1         1            1           43s
    redis     1         1         1            1           43s
    result    1         1         1            1           43s
    vote      5         5         5            5           43s
    worker    3         3         3            3           43s

    > kubectl get svc
    NAME         TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
    db           ClusterIP      None           <none>        55555/TCP      48s
    kubernetes   ClusterIP      10.0.0.1       <none>        443/TCP        6h
    redis        ClusterIP      10.0.168.188   <none>        6379/TCP       48s
    result       LoadBalancer   10.0.76.157    <pending>     80:31033/TCP   47s
    vote         LoadBalancer   10.0.244.254   <pending>     80:31330/TCP   48s
    worker       ClusterIP      None           <none>        55555/TCP      48s
    ```

    Notice that the two `LoadBalancer` services are `pending`. This is because it takes a minute or two to provision an Azure load-balancer. You can run the `kubectl get svc --watch` command to see when they have finished provisioning. Once they're provisioned the output will look like this (the external IPs will be different in your environment).

    ```
    > kubectl get svc
    <Snip>
    result    LoadBalancer   10.0.76.157   52.174.195.232   80:31033/TCP   7m
    vote      LoadBalancer   10.0.244.254  52.174.196.199   80:31330/TCP   8m
    ```

4. Check that the application works.

    Copy the `EXTERNAL-IP` value for the `vote` service, paste it into a browser tab and cast a vote.

    Copy the `EXTERNAL-IP` value for the `result` service, paste it into a browser tab and check that the vote registered.

You should thoroughly test the application now that the stack is deployed and all services running.

You can extend your Kubernetes manifest file to include advanced features such as how to perform rolling updates and simple rollbacks. But you should not do this until you have confirmed your application is working with the simplified manifest file.

You should also test failure scenarios, increasing load, scaling operations, updates, rollbacks, and any other operations that are considered important for the lifecycle of the application. These tests will be specific to each of your apps, and are beyond the scope of this document. However, you should be sure to complete them before beginning the migration of your application.

If you had a CI/CD pipeline with automated tests and deployments for your Docker Cloud stack, you should build, test, and implement one for the app on AKS.

## Start migration

You should not terminate your Docker Cloud stacks or node clusters until a while after the migration has been signed off as successful. Keeping them on hand is a good precaution in case you experience issues with the migration and need to switch back for any reason.
