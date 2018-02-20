---
description: Migrating from Docker Cloud
keywords: cloud, kubernetes, migration
title: Migrate Docker Cloud stack to Kubernetes on GKE
---

## Intro

Docker Cloud applications are defined in declarative configuration files and run on cloud instances in the public cloud. Kubernetes applications are defined in declarative manifest files and can also run on cloud instances in the public cloud.

This document will guide you through the process of converting your Docker Cloud application into a Kubernetes application. It will then show you how to deploy and test it on a Google Container Engine (GKE) cluster.

GKE is a hosted Kubernetes service on the Google public cloud. It exposes all of the standard Kubernetes API endpoints, meaning you can use standard Kubernetes tools and manifests. This means that the manifest files you create can be used to run your application on any Kubernetes infrastructure --- different public clouds and on-premises.

The high-level process will involve:

- Building a Kubernetes cluster on GKE (this is equivalent of your *node cluster* in Docker Cloud)
- Converting your Docker Cloud stack into a Kubernetes application defined in a Kubernetes manifest file
- Testing the Kubernetes application on GKE
- Deploying the Kubernetes application on GKE

The actual process of performing the migration (switching customers from using your Docker Cloud application to using your Kubernetes application on GKE) will vary considerably between applications, and is therefore outside the scope of this document. This document will walk you through the process of converting your Docker Cloud app into a Kubernetes app, and testing it works on GKE. It will not outline migration steps.

## Things that will be different after migration

As part of a migration like this, some of the ways you manage your Docker application will change. However, the way that users and external systems interact with it will not change.

Here are some examples of things that will change:

- You will no longer be able to deploy and manage your application using the intuitive Docker Cloud web UI.
- You will no longer have authorization integrated with the Docker platform
- You will lose the Docker Cloud **autoredeploy** feature.

> Autoredeploy is a Docker Cloud feature that automatically updates your running application every time an updated image is pushed or built. This advanced feature is not included with Docker CE, but you may be able to regain it by using a 3rd party product.

Many things about your application will stay the same. For example, your Docker images and overall application functionality will not change --- if your application uses a Docker image called `myorg/webfe:v3` and is accessible on port `80`, none of this will change. Also, the way your customers and other applications interact with your application will not change.

## Pre-reqs

To complete the migration, you will need:

- An active Google Cloud subscription with billing enabled.



## Build target environment

Google Container Engine (GKE) is a managed Kubernetes service. As such, it takes care of all of the Kubernetes control plane management --- delivering the control plane APIs, managing control plane HA, managing control plane upgrades etc. You only need to manage worker nodes --- how many, size & spec, where to deploy them etc.

We'll complete the following three high-level steps as part of building your GKE cluster:

1. Create a new project
2. Build a GKE cluster
3. Connect to the GKE cluster

### Create a new project

Everything in the Google Cloud Platform has to sit inside of a *project*. Let's create one.

1. Log in to the [Google Cloud Platform Console](https://console.cloud.google.com).

2. Create new project.

    Either of the following options will take you to the `New Project` screen.

    - Select the `Create an empty project` option from the GCP Console home screen (only works if you aren't already in a project).
    - Click the `Select a project` dropdown from the top of the screen and then click the `+` button to create a new project.

3. Give the project a name and click `Create`.

    The examples in this document will assume a project called `proj-k8s-vote`.

    It may take a minute for the project to create.

### Build a GKE cluster

In this section, we'll create a new 3-node GKE cluster in the `europe-west` Region. Your cluster will be different, and should probably be based on the configuration of your Docker Cloud *node cluster*.

You can use one of the following methods to see the configuration of your Docker Cloud *node cluster*:

- **Docker Cloud classic mode users:** Log-in to Docker Cloud and select `Node Clusters` under `INFRASTRUCTURE`. From here you can select any of your node clusters to find information such as; cloud provider, number of nodes, spec of nodes, region and availability zone.
- **Docker Cloud Swarm mode users:** Log-in to Docker Cloud and select `Swarms` from the top menu. Select one of your Swarms.

Before continuing, you will need to know all of the following:

- The GCP *Region* and *Zone* that you want to deploy your GKE cluster to
- Number of nodes
- Size and spec of nodes

1. If you aren't already, log-in to the [GCP Console](https://console.cloud.google.com).

2. Select the desired project from the `Select a project` dropdown menu at the top of the Console screen.

3. Click `Kubernetes Engine` from the left-hand menu.

    It may take a minute for the Kubernetes Engine to start.

4. Click `Create Cluster`.

5. Configure the required cluster options:

    - **Name:** An arbitrary name you give the cluster.
    - **Description:** An arbitrary description for the cluster.
    - **Location:** Whether you want the Kubernetes control plane nodes (masters) in a single availability zone or spread across availability zones within a GCP Region.
    - **Zone/Region:** The zone or region to deploy the cluster to.
    - **Cluster version:** The version of Kubernetes to use. You should probably use a 1.8.x version.
    - **Machine type:** The type of GCE VM to use as worker nodes in the cluster. This should probably match your Docker Cloud node cluster.
    - **Node image:** The OS to run on each Kubernetes worker node. Use Ubuntu if you require NFS, glusterfs, Sysdig, or Debian packages.
    - **Size:** How many nodes to have in the cluster.

    There are various other options that you can configure, and you should carefully consider them. However, the majority of deployments should be OK with default values for these.

6. Click `Create`.

    It'll take a minute or two for the cluster to create.

Once the cluster is created, you can click it's name to see more details.

### Connect to the GKE cluster

You can use the web-based *Google cloud shell*, or a properly configured local terminal to connect to your cluster. We'll show you how to configure your laptop, or other local terminal, to connect to your cluster.

We'll complete the following high-level steps:

- Install and configure the `gcloud` CLI tool
- Install `kubectl` --- the Kubernetes command line tool
- Configure `kubectl` to connect to your cluster.

The `gcloud` tool is the command-line tool for interacting with the Google Cloud Platform. It is installed as part of the Google Cloud SDK.

1. Click this [link](https://cloud.google.com/sdk/) and select the option for your Operating System.

2. Click the link to download the SDK installer.

3. Install the SDK.

    Follow the prompts.

4. Configure `gcloud`.

    Open a terminal window and enter the following command:

    ```
    > gcloud init --console-only
    ```

    Follow all prompts, including the one to open a web browser and approve the requested authorizations.

    As part of the procedure you will need to copy and paste a code into the terminal window to authorize `gcloud`.

5. Install `kubectl`.

    `kubectl` is the Kubernetes command line tool and instructions on how to download and install it can be found [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

6. Configure `kubectl` to talk to your GKE cluster.

    Go back to the GKE section of the GCP Console and click the `Connect` button at the end of the line representing your cluster.

    Copy the long command and paste it into your local terminal window. Your command might be slightly different.

    ```
    > gcloud container clusters get-credentials clus-k8s-vote --zone europe-west2-c --project proj-k8s-vote

    Fetching cluster endpoint and auth data.
    kubeconfig entry generated for clus-k8s-vote.
    ```

7. Run a command to test the `kubectl` configuration.

    ```
    > kubectl get nodes
    NAME                                           STATUS    ROLES     AGE       VERSION
    gke-clus-k8s-vote-default-pool-81bd226c-2jtp   Ready     <none>    1h        v1.9.2-gke.1
    gke-clus-k8s-vote-default-pool-81bd226c-mn4k   Ready     <none>    1h        v1.9.2-gke.1
    gke-clus-k8s-vote-default-pool-81bd226c-qjm2   Ready     <none>    1h        v1.9.2-gke.1
    ```

    The `kubectl get nodes` command returns information about the nodes in the cluster it is configured to manage. If the values returned match your GKE cluster (number of nodes, age, and version) then you have successfully configured `kubectl` to manage your AKS cluster.

You now have a GKE cluster and have configured `kubectl` to manage it. Let's look at how to convert your Docker Cloud app into a Kubernetes app.


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

Let's step throught he process of converting each one.

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

By naming the *Service* "db", we've ensured that a cluster-wide DNS mapping has been created for `db` pointing to this *Service*. We've then defined a label selector (`Service.spec.selector`) with a value of "app=db". This means that the service will provide stable networking a load-balancing for all *Pods* (containers) on the cluster with matching labels. Notice that the *Pods* defined in the *Deployment* section of the file are all laballed as "app-db". It is this mapping between the *Service's* label selector and the *Pod* labels that tell the *Service* object which Pods to provide networking for.


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
