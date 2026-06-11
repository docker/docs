---
description: Containerize and develop Bun applications using Docker.
keywords: getting started, bun
title: Bun language-specific guide
summary: |
  Learn how to containerize JavaScript applications with the Bun runtime.
linkTitle: Bun
aliases:
  - /guides/bun/configure-ci-cd/
  - /guides/bun/containerize/
  - /guides/bun/deploy/
  - /guides/bun/develop/
params:
  tags: [cicd]
  time: 10 minutes
---


The Bun getting started guide teaches you how to create a containerized Bun application using Docker. 

> **Acknowledgment**
>
> Docker would like to thank [Pradumna Saraf](https://twitter.com/pradumna_saraf) for his contribution to this guide.

## What will you learn?

* Containerize and run a Bun application using Docker
* Set up a local environment to develop a Bun application using containers
* Configure a CI/CD pipeline for a containerized Bun application using GitHub Actions
* Deploy your containerized application locally to Kubernetes to test and debug your deployment

## Prerequisites

- Basic understanding of JavaScript is assumed.
- You must have familiarity with Docker concepts like containers, images, and Dockerfiles. If you are new to Docker, you can start with the [Docker basics](/get-started/docker-concepts/the-basics/what-is-a-container.md) guide.

After completing the Bun getting started modules, you should be able to containerize your own Bun application based on the examples and instructions provided in this guide.

Start by containerizing an existing Bun application.

## Containerize a Bun application

### Prerequisites

* You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

### Overview

For a long time, Node.js has been the de-facto runtime for server-side
JavaScript applications. Recent years have seen a rise in new alternative
runtimes in the ecosystem, including [Bun website](https://bun.sh/). Like
Node.js, Bun is a JavaScript runtime. Bun is a comparatively lightweight
runtime that is designed to be fast and efficient.

Why develop Bun applications with Docker? Having multiple runtimes to choose
from is great. But as the number of runtimes increases, it becomes challenging
to manage the different runtimes and their dependencies consistently across
environments. This is where Docker comes in. Creating and destroying containers
on demand is a great way to manage the different runtimes and their
dependencies. Also, as it's fairly a new runtime, getting a consistent
development environment for Bun can be challenging. Docker can help you set up
a consistent development environment for Bun.

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following
command to clone the repository:

```console
$ git clone https://github.com/dockersamples/bun-docker.git && cd bun-docker
```

You should now have the following contents in your `bun-docker` directory.

```text
├── bun-docker/
│ ├── compose.yml
│ ├── Dockerfile
│ ├── LICENSE
│ ├── server.js
│ └── README.md
```

### Create a Dockerfile

Before creating a Dockerfile, you need to choose a base image. You can either use the [Bun Docker Official Image](https://hub.docker.com/r/oven/bun) or a Docker Hardened Image (DHI) from the [Hardened Image catalog](https://hub.docker.com/hardened-images/catalog).

Choosing DHI offers the advantage of a production-ready image that is lightweight and secure. For more information, see [Docker Hardened Images](https://docs.docker.com/dhi/).

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}

Docker Hardened Images (DHIs) are available for Bun in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog/dhi/bun). You can pull DHIs directly from the `dhi.io` registry.

1. Sign in to the DHI registry:
   ```console
   $ docker login dhi.io
   ```

2. Pull the Bun DHI as `dhi.io/bun:1`. The tag (`1`) in this example refers to the version to the latest 1.x version of Bun.

   ```console
   $ docker pull dhi.io/bun:1
   ```

For other available versions, refer to the [catalog](https://hub.docker.com/hardened-images/catalog/dhi/bun).

```dockerfile
# Use the DHI Bun image as the base image
FROM dhi.io/bun:1

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Expose the port on which the API will listen
EXPOSE 3000

# Run the server when the container launches
CMD ["bun", "server.js"]
```

{{< /tab >}}
{{< tab name="Using the official image" >}}

Using the Docker Official Image is straightforward. In the following Dockerfile, you'll notice that the `FROM` instruction uses `oven/bun` as the base image.

You can find the image on [Docker Hub](https://hub.docker.com/r/oven/bun). This is the Docker Official Image for Bun created by Oven, the company behind Bun, and it's available on Docker Hub.

```dockerfile
# Use the official Bun image
FROM oven/bun:latest

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Expose the port on which the API will listen
EXPOSE 3000

# Run the server when the container launches
CMD ["bun", "server.js"]
```

{{< /tab >}}
{{< /tabs >}}

In addition to specifying the base image, the Dockerfile also:

- Sets the working directory in the container to `/app`.
- Copies the content of the current directory to the `/app` directory in the container.
- Exposes port 3000, where the API is listening for requests.
- And finally, starts the server when the container launches with the command `bun server.js`.

### Run the application

Inside the `bun-docker` directory, run the following command in a terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000). You will see a message `{"Status" : "OK"}` in the browser.

In the terminal, press `ctrl`+`c` to stop the application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `bun-docker` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000).


In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

### Summary

In this section, you learned how you can containerize and run your Bun
application using Docker.

Related information:

 - [Dockerfile reference](/reference/dockerfile.md)
 - [.dockerignore file](/reference/dockerfile.md#dockerignore-file)
 - [Docker Compose overview](/manuals/compose/_index.md)
 - [Compose file reference](/reference/compose-file/_index.md)
 - [Docker Hardened Images](/dhi/)

### Next steps

In the next section, you'll learn how you can develop your application using
containers.

## Use containers for Bun development

### Prerequisites

Complete [Containerize a Bun application](containerize.md).

### Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Configuring Compose to automatically update your running Compose services as you edit and save your code

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/dockersamples/bun-docker.git && cd bun-docker
```

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](/manuals/compose/how-tos/file-watch.md).

Open your `compose.yml` file in an IDE or text editor and then add the Compose Watch instructions. The following example shows how to add Compose Watch to your `compose.yml` file.

```yaml {hl_lines="9-12",linenos=true}
services:
  server:
    image: bun-server
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    develop:
      watch:
        - action: rebuild
          path: .
```

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

Now, if you modify your `server.js` you will see the changes in real time without re-building the image.

To test it out, open the `server.js` file in your favorite text editor and change the message from `{"Status" : "OK"}` to `{"Status" : "Updated"}`. Save the file and refresh your browser at `http://localhost:3000`. You should see the updated message.

Press `ctrl+c` in the terminal to stop your application.

### Summary

In this section, you also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:
 - [Compose file reference](/reference/compose-file/)
 - [Compose file watch](/manuals/compose/how-tos/file-watch.md)
 - [Multi-stage builds](/manuals/build/building/multi-stage.md)

### Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

## Configure CI/CD for your Bun application

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a Bun application](containerize.md). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

### Overview

In this section, you'll learn how to set up and use GitHub Actions to build and test your Docker image as well as push it to Docker Hub. You will complete the following steps:

1. Create a new repository on GitHub.
2. Define the GitHub Actions workflow.
3. Run the workflow.

### Step one: Create the repository

Create a GitHub repository, configure the Docker Hub credentials, and push your source code.

1. [Create a new repository](https://github.com/new) on GitHub.

2. Open the repository **Settings**, and go to **Secrets and variables** >
   **Actions**.

3. Create a new **Repository variable** named `DOCKER_USERNAME` and your Docker ID as a value.

4. Create a new [Personal Access Token (PAT)](/manuals/security/access-tokens.md#create-an-access-token)for Docker Hub. You can name this token `docker-tutorial`. Make sure access permissions include Read and Write.

5. Add the PAT as a **Repository secret** in your GitHub repository, with the name
   `DOCKERHUB_TOKEN`.

6. In your local repository on your machine, run the following command to change
   the origin to the repository you just created. Make sure you change
   `your-username` to your GitHub username and `your-repository` to the name of
   the repository you created.

   ```console
   $ git remote set-url origin https://github.com/your-username/your-repository.git
   ```

7. Run the following commands to stage, commit, and push your local repository to GitHub.

   ```console
   $ git add -A
   $ git commit -m "my commit"
   $ git push -u origin main
   ```

### Step two: Set up the workflow

Set up your GitHub Actions workflow for building, testing, and pushing the image
to Docker Hub.

1. Go to your repository on GitHub and then select the **Actions** tab.

2. Select **set up a workflow yourself**.

   This takes you to a page for creating a new GitHub actions workflow file in
   your repository, under `.github/workflows/main.yml` by default.

3. In the editor window, copy and paste the following YAML configuration and commit the changes.

   ```yaml
   name: ci

   on:
     push:
       branches:
         - main

   jobs:
     build:
       runs-on: ubuntu-latest
       steps:
         - name: Login to Docker Hub
           uses: docker/login-action@{{% param "login_action_version" %}}
           with:
             username: ${{ vars.DOCKER_USERNAME }}
             password: ${{ secrets.DOCKERHUB_TOKEN }}

         - name: Set up Docker Buildx
           uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

         - name: Build and push
           uses: docker/build-push-action@{{% param "build_push_action_version" %}}
           with:
             platforms: linux/amd64,linux/arm64
             push: true
             tags: ${{ vars.DOCKER_USERNAME }}/${{ github.event.repository.name }}:latest
   ```

   For more information about the YAML syntax for `docker/build-push-action`,
   refer to the [GitHub Action README](https://github.com/docker/build-push-action/blob/master/README.md).

### Step three: Run the workflow

Save the workflow file and run the job.

1. Select **Commit changes...** and push the changes to the `main` branch.

   After pushing the commit, the workflow starts automatically.

2. Go to the **Actions** tab. It displays the workflow.

   Selecting the workflow shows you the breakdown of all the steps.

3. When the workflow is complete, go to your
   [repositories on Docker Hub](https://hub.docker.com/repositories).

   If you see the new repository in that list, it means the GitHub Actions
   successfully pushed the image to Docker Hub.

### Summary

In this section, you learned how to set up a GitHub Actions workflow for your Bun application.

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

### Next steps

Next, learn how you can locally test and debug your workloads on Kubernetes before deploying.

## Test your Bun deployment

### Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a Bun application](containerize.md).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

### Overview

In this section, you'll learn how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine. This allows you to test and debug your workloads on Kubernetes locally before deploying.

### Create a Kubernetes YAML file

In your `bun-docker` directory, create a file named
`docker-kubernetes.yml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your Bun application](configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-bun-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: bun-api
  template:
    metadata:
      labels:
        app: bun-api
    spec:
      containers:
       - name: bun-api
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
    app: bun-api
  ports:
  - port: 3000
    targetPort: 3000
    nodePort: 30001
```

In this Kubernetes YAML file, there are two objects, separated by the `---`:

 - A Deployment, describing a scalable group of identical pods. In this case,
   you'll get just one replica, or copy of your pod. That pod, which is
   described under `template`, has just one container in it. The
    container is created from the image built by GitHub Actions in [Configure CI/CD for
    your Bun application](configure-ci-cd.md).
 - A NodePort service, which will route traffic from port 30001 on your host to
   port 3000 inside the pods it routes to, allowing you to reach your app
   from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

### Deploy and check your application

1. In a terminal, navigate to `bun-docker` and deploy your application to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-kubernetes.yml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```text
   deployment.apps/docker-bun-demo created
   service/service-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                 READY   UP-TO-DATE   AVAILABLE    AGE
   docker-bun-demo       1/1     1            1           10s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
   kubernetes           ClusterIP   10.96.0.1        <none>        443/TCP          88m
   service-entrypoint   NodePort    10.105.145.223   <none>        3000:30001/TCP   83s
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

### Summary

In this section, you learned how to use Docker Desktop to deploy your Bun application to a fully-featured Kubernetes environment on your development machine. 

Related information:
   - [Kubernetes documentation](https://kubernetes.io/docs/home/)
   - [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
   - [Swarm mode overview](/manuals/engine/swarm/_index.md)
