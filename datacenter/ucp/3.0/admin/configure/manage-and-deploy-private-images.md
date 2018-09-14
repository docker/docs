---
title: Manage and deploy private images
description: Learn how to push an image to Docker Trusted Registry and deploy it to a Kubernetes cluster managed by Docker Enterprise Edition.
keywords: Docker EE, DTR, UCP, image, Kubernetes, orchestration, cluster
---

Docker Enterprise Edition (EE) has its own image registry (DTR) so that
you can store and manage the images that you deploy to your cluster.
In this topic, you push an image to DTR and later deploy it to your cluster,
using the Kubernetes orchestrator.

## Open the DTR web UI

1.  In the Docker EE web UI, click **Admin Settings**.
2.  In the left pane, click **Docker Trusted Registry**.
3.  In the **Installed DTRs** section, note the URL of your cluster's DTR
    instance.

    ![](../../images/manage-and-deploy-private-images-1.png){: .with-border}

4.  In a new browser tab, enter the URL to open the DTR web UI.

## Create an image repository

1.  In the DTR web UI, click **Repositories**.
2.  Click **New Repository**, and in the **Repository Name** field, enter
    "wordpress".
3.  Click **Save** to create the repository.

    ![](../../images/manage-and-deploy-private-images-2.png){: .with-border}

## Push an image to DTR

Instead of building an image from scratch, we'll pull the official WordPress
image from Docker Hub, tag it, and push it to DTR. Once that WordPress version
is in DTR, only authorized users can change it.

To push images to DTR, you need CLI access to a licensed installation of
Docker EE.

- [License your installation](license-your-installation.md).
- [Set up your Docker CLI](../../user-acccess/cli.md).

When you're set up for CLI-based access to a licensed Docker EE instance,
you can push images to DTR.

1.  Pull the public WordPress image from Docker Hub:

    ```bash
    docker pull wordpress
    ```

2.  Tag the image, using the IP address or DNS name of your DTR instance:

    ```bash
    docker tag wordpress:latest <dtr-url>:<port>/admin/wordpress:latest
    ```
3.  Log in to a Docker EE manager node.
4.  Push the tagged image to DTR:

    ```bash
    docker image push <dtr-url>:<port>/admin/wordpress:latest
    ```

## Confirm the image push

In the DTR web UI, confirm that the `wordpress:latest` image is store in your
DTR instance.

1.  In the DTR web UI, click **Repositories**.
2.  Click **wordpress** to open the repo.
3.  Click **Images** to view the stored images.
4.  Confirm that the `latest` tag is present.

    ![](../../images/manage-and-deploy-private-images-3.png){: .with-border}

You're ready to deploy the `wordpress:latest` image into production.

## Deploy the private image to UCP

With the WordPress image stored in DTR, Docker EE can deploy the image to a
Kubernetes cluster with a simple Deployment object:

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: wordpress-deployment
spec:
  selector:
    matchLabels:
      app: wordpress
  replicas: 2
  template:
    metadata:
      labels:
        app: wordpress
    spec:
      containers:
      - name: wordpress
        image: <dtr-url>:<port>/admin/wordpress:latest
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: wordpress-service
  labels:
    app: wordpress
spec:
  type: NodePort
  ports:
    - port: 80
      nodePort: 30081
  selector:
    app: wordpress
```

The Deployment object's YAML specifies your DTR image in the pod template spec:
`image: <dtr-url>:<port>/admin/wordpress:latest`. Also, the YAML file defines
a `NodePort` service that exposes the WordPress application, so it's accessible
from outside the cluster.

1.  Open the Docker EE web UI, and in the left pane, click **Kubernetes**.
2.  Click **Create** to open the **Create Kubernetes Object** page.
3.  In the **Namespace** dropdown, select **default**.
4.  In the **Object YAML** editor, paste the Deployment object's YAML.
5.  Click **Create**. When the Kubernetes objects are created,
    the **Load Balancers** page opens.
6.  Click **wordpress-service**, and in the details pane, find the **Ports**
    section.
7.  Click the URL to open the default WordPress home page.  

    ![](../../images/manage-and-deploy-private-images-4.png){: .with-border}

