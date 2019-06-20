---
title: Deploy a workload to a Kubernetes cluster
description: Use Docker Enterprise Edition to deploy Kubernetes workloads from yaml files.
keywords: UCP, Docker EE, orchestration, Kubernetes, cluster
redirect_from:
  - /ee/ucp/user/services/deploy-kubernetes-workload/
---

The Docker EE web UI enables deploying your Kubernetes YAML files. In most
cases, no modifications are necessary to deploy on a cluster that's managed by
Docker EE.

## Deploy an NGINX server

In this example, a simple Kubernetes Deployment object for an NGINX server is
defined in YAML:

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

The YAML specifies an earlier version of NGINX, which will be updated in a
later section.

1. Open the Docker EE web UI, and in the left pane, click **Kubernetes**.
2. Click **Create** to open the **Create Kubernetes Object** page.
3. In the **Namespace** dropdown, select **default**.
4. In the **Object YAML** editor, paste the previous YAML.
5. Click **Create**.

![](../images/deploy-kubernetes-workload-1.png){: .with-border}

## Inspect the deployment

The Docker EE web UI shows the status of your deployment when you click the
links in the **Kubernetes** section of the left pane.

1.  In the left pane. click **Controllers** to see the resource controllers
    that Docker EE created for the NGINX server.
2.  Click the **nginx-deployment** controller, and in the details pane, scroll
    to the **Template** section. This shows the values that Docker EE used to
    create the deployment.
3.  In the left pane, click **Pods** to see the pods that are provisioned for
    the NGINX server. Click one of the pods, and in the details pane, scroll to
    the **Status** section to see that pod's phase, IP address, and other
    properties.

![](../images/deploy-kubernetes-workload-2.png){: .with-border}

## Expose the server

The NGINX server is up and running, but it's not accessible from outside of the
cluster. Add a `NodePort` service to expose the server on a specified port:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  type: NodePort
  ports:
    - port: 80
      nodePort: 32768
  selector:
    app: nginx
```

The service connects the cluster's internal port 80 to the external port
32768.

1.  Repeat the previous steps and copy-paste the YAML that defines the `nginx`
    service into the **Object YAML** editor on the
    **Create Kubernetes Object** page. When you click **Create**, the
    **Load Balancers** page opens.
2.  Click the **nginx** service, and in the details pane, find the **Ports**
    section.

    ![](../images/deploy-kubernetes-workload-3.png){: .with-border}

3.  Click the link that's labeled **URL** to view the default NGINX page.

The YAML definition connects the service to the NGINX server by using the
app label `nginx` and a corresponding label selector.
[Learn about using a service to expose your app](https://v1-8.docs.kubernetes.io/docs/tutorials/kubernetes-basics/expose-intro/).

## Update the deployment

Update an existing deployment by applying an updated YAML file. In this
example, the server is scaled up to four replicas and updated to a later
version of NGINX.

```yaml
...
spec:
  progressDeadlineSeconds: 600
  replicas: 4
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: nginx
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx:1.8
...
```

1.  In the left pane, click **Controllers** and select **nginx-deployment**.
2.  In the details pane, click **Configure**, and in the **Edit Deployment**
    page, find the **replicas: 2** entry.
3.  Change the number of replicas to 4, so the line reads **replicas: 4**.
4.  Find the **image: nginx:1.7.9** entry and change it to **image: nginx:1.8**.

    ![](../images/deploy-kubernetes-workload-4.png){: .with-border}

5.  Click **Save** to update the deployment with the new YAML.
6.  In the left pane, click **Pods** to view the newly created replicas.

    ![](../images/deploy-kubernetes-workload-5.png){: .with-border}

## Use the CLI to deploy Kubernetes objects

With Docker EE, you deploy your Kubernetes objects on the command line by using
`kubectl`. [Install and set up kubectl](https://v1-8.docs.kubernetes.io/docs/tasks/tools/install-kubectl/).

Use a client bundle to configure your client tools, like Docker CLI and `kubectl`
to communicate with UCP instead of the local deployments you might have running.
[Get your client bundle by using the Docker EE web UI or the command line](../user-access/cli.md).

When you have the client bundle set up, you can deploy a Kubernetes object
from YAML.

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  type: NodePort
  ports:
    - port: 80
      nodePort: 32768
  selector:
    app: nginx
```

Save the previous YAML to a file named "deployment.yaml", and use the following
command to deploy the NGINX server:

```bash
kubectl apply -f deployment.yaml
```

## Inspect the deployment

Use the `describe deployment` option to inspect the deployment:

```bash
kubectl describe deployment nginx-deployment
```

Also, you can use the Docker EE web UI to see the deployment's pods and
controllers.

## Update the deployment

Update an existing deployment by applying an updated YAML file.

Edit deployment.yaml and change the following lines:

- Increase the number of replicas to 4, so the line reads **replicas: 4**.
- Update the NGINX version by specifying **image: nginx:1.8**.

Save the edited YAML to a file named "update.yaml", and use the following
command to deploy the NGINX server:

```bash
kubectl apply -f update.yaml
```

Check that the deployment was scaled out by listing the deployments in the
cluster:

```bash
 kubectl get deployments
```

You should see four pods in the deployment:

```bash
NAME                   DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment       4         4         4            4           2d
```

Check that the pods are running the updated image:

```bash
kubectl describe deployment nginx-deployment | grep -i image
```

You should see the currently running image:

```bash
    Image:        nginx:1.8
```

