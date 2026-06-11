---
title: C++ language-specific guide
linkTitle: C++
description: Containerize and develop C++ applications using Docker.
keywords: getting started, c++
summary: |
  This guide explains how to containerize C++ applications using Docker.
aliases:
  - /language/cpp/
  - /guides/language/cpp/
  - /language/cpp/containerize/
  - /language/cpp/develop/
  - /language/cpp/configure-ci-cd/
  - /language/cpp/deploy/
  - /guides/cpp/configure-ci-cd/
  - /guides/cpp/containerize/
  - /guides/cpp/deploy/
  - /guides/cpp/develop/
  - /guides/cpp/multistage/
  - /guides/cpp/security/
params:
  tags: [cicd]
  time: 20 minutes
---


The C++ getting started guide teaches you how to create a containerized C++ application using Docker. In this guide, you'll learn how to:

> **Acknowledgment**
>
> Docker would like to thank [Pradumna Saraf](https://twitter.com/pradumna_saraf) and [Mohammad-Ali A'râbi](https://twitter.com/MohammadAliEN) for their contribution to this guide.

- Containerize and run a C++ application using a multi-stage Docker build
- Build and run a C++ application using Docker Compose
- Set up a local environment to develop a C++ application using containers
- Configure a CI/CD pipeline for a containerized C++ application using GitHub Actions
- Deploy your containerized application locally to Kubernetes to test and debug your deployment
- Use BuildKit to generate SBOM attestations during the build process

After completing the C++ getting started modules, you should be able to containerize your own C++ application based on the examples and instructions provided in this guide.

Start by containerizing an existing C++ application.

## Create a multi-stage build for your C++ application

### Prerequisites

- You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

### Overview

This section walks you through creating a multi-stage Docker build for a C++ application.
A multi-stage build is a Docker feature that allows you to use different base images for different stages of the build process,
so you can optimize the size of your final image and separate build dependencies from runtime dependencies.

The standard practice for compiled languages like C++ is to have a build stage that compiles the code and a runtime stage that runs the compiled binary,
because the build dependencies are not needed at runtime.

### Get the sample application

Let's use a simple C++ application that prints `Hello, World!` to the terminal. To do so, clone the sample repository to use with this guide:

```bash
$ git clone https://github.com/dockersamples/c-plus-plus-docker.git
```

The example for this section is under the `hello` directory in the repository. Get inside it and take a look at the files:

```bash
$ cd c-plus-plus-docker/hello
$ ls
```

You should see the following files:

```text
Dockerfile  hello.cpp
```

### Check the Dockerfile

Open the `Dockerfile` in an IDE or text editor. The `Dockerfile` contains the instructions for building the Docker image.

```Dockerfile
# Stage 1: Build stage
FROM ubuntu:latest AS build

# Install build-essential for compiling C++ code
RUN apt-get update && apt-get install -y build-essential

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY hello.cpp .

# Compile the C++ code statically to ensure it doesn't depend on runtime libraries
RUN g++ -o hello hello.cpp -static

# Stage 2: Runtime stage
FROM scratch

# Copy the static binary from the build stage
COPY --from=build /app/hello /hello

# Command to run the binary
CMD ["/hello"]
```

The `Dockerfile` has two stages:

1. **Build stage**: This stage uses the `ubuntu:latest` image to compile the C++ code and create a static binary.
2. **Runtime stage**: This stage uses the `scratch` image, which is an empty image, to copy the static binary from the build stage and run it.

### Build the Docker image

To build the Docker image, run the following command in the `hello` directory:

```bash
$ docker build -t hello .
```

The `-t` flag tags the image with the name `hello`.

### Run the Docker container

To run the Docker container, use the following command:

```bash
$ docker run hello
```

You should see the output `Hello, World!` in the terminal.

### Summary

In this section, you learned how to create a multi-stage build for a C++ application. Multi-stage builds help you optimize the size of your final image and separate build dependencies from runtime dependencies.
In this example, the final image only contains the static binary and doesn't include any build dependencies.

As the image has an empty base, the usual OS tools are also absent. So, for example, you can't run a simple `ls` command in the container:

```bash
$ docker run hello ls
```

This makes the image very lightweight and secure.

## Containerize a C++ application

### Prerequisites

- You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

### Overview

This section walks you through containerizing and running a C++ application, using Docker Compose.

### Get the sample application

We're using the same sample repository that you used in the previous sections of this guide. If you haven't already cloned the repository, clone it now:

```console
$ git clone https://github.com/dockersamples/c-plus-plus-docker.git
```

You should now have the following contents in your `c-plus-plus-docker` (root)
directory.

```text
├── c-plus-plus-docker/
│ ├── compose.yml
│ ├── Dockerfile
│ ├── LICENSE
│ ├── ok_api.cpp
│ └── README.md

```

To learn more about the files in the repository, see the following:

- [Dockerfile](/reference/dockerfile.md)
- [.dockerignore](/reference/dockerfile.md#dockerignore-file)
- [compose.yml](/reference/compose-file/_index.md)

### Run the application

Inside the `c-plus-plus-docker` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You will see a message `{"Status" : "OK"}` in the browser.

In the terminal, press `ctrl`+`c` to stop the application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `c-plus-plus-docker` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080).

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/).

### Summary

In this section, you learned how you can containerize and run your C++
application using Docker.

Related information:

- [Docker Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll learn how you can develop your application using
containers.

## Use containers for C++ development

### Prerequisites

Complete [Containerize a C++ application](containerize.md).

### Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Configuring Compose to automatically update your running Compose services as you edit and save your code

### Get the sample application

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/dockersamples/c-plus-plus-docker.git && cd c-plus-plus-docker
```

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](/manuals/compose/how-tos/file-watch.md).

Open your `compose.yml` file in an IDE or text editor and then add the Compose Watch instructions. The following example shows how to add Compose Watch to your `compose.yml` file.

```yaml {hl_lines="11-14",linenos=true}
services:
  ok-api:
    image: ok-api
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    develop:
      watch:
        - action: rebuild
          path: .
```

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

Now, if you modify your `ok_api.cpp` you will see the changes in real time without re-building the image.

To test it out, open the `ok_api.cpp` file in your favorite text editor and change the message from `{"Status" : "OK"}` to `{"Status" : "Updated"}`. Save the file and refresh your browser at [http://localhost:8080](http://localhost:8080). You should see the updated message.

Press `ctrl+c` in the terminal to stop your application.

### Summary

In this section, you also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose file watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)

### Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

## Configure CI/CD for your C++ application

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a C++ application](containerize.md). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

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

4. Create a new [Personal Access Token (PAT)](/manuals/security/access-tokens.md#create-an-access-token) for Docker Hub. You can name this token `docker-tutorial`. Make sure access permissions include Read and Write.

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

In this section, you learned how to set up a GitHub Actions workflow for your C++ application.

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

### Next steps

Next, learn how you can locally test and debug your workloads on Kubernetes before deploying.

## Test your C++ deployment

### Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a C++ application](containerize.md).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

### Overview

In this section, you'll learn how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine. This allows you to test and debug your workloads on Kubernetes locally before deploying.

### Create a Kubernetes YAML file

In your `c-plus-plus-docker` directory, create a file named
`docker-kubernetes.yml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your C++ application](configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-c-plus-plus-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: ok-api
  template:
    metadata:
      labels:
        service: ok-api
    spec:
      containers:
        - name: ok-api-service
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
    service: ok-api
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
  your C++ application](configure-ci-cd.md).
- A NodePort service, which will route traffic from port 30001 on your host to
  port 8080 inside the pods it routes to, allowing you to reach your app
  from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

### Deploy and check your application

1. In a terminal, navigate to `c-plus-plus-docker` and deploy your application to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-kubernetes.yml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```text
   deployment.apps/docker-c-plus-plus-demo created
   service/service-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                     READY   UP-TO-DATE   AVAILABLE    AGE
   docker-c-plus-plus-demo   1/1     1            1           10s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
   kubernetes           ClusterIP   10.96.0.1        <none>        443/TCP          88m
   service-entrypoint   NodePort    10.105.145.223   <none>        8080:30001/TCP   83s
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

In this section, you learned how to use Docker Desktop to deploy your C++ application to a fully-featured Kubernetes environment on your development machine.

Related information:

- [Kubernetes documentation](https://kubernetes.io/docs/home/)
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
- [Swarm mode overview](/manuals/engine/swarm/_index.md)

## Supply-chain security for C++ Docker images

### Prerequisites

- You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.
- You have a Docker Desktop installed, with containerd enabled for pulling and storing images (it's a checkbox in **Settings** > **General**). Otherwise, if you use Docker Engine:
  - You have the [Docker Scout CLI plugin](https://docs.docker.com/scout/install/) installed. To install it on Docker Engine, use the following command:

    ```bash
    $ curl -sSfL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh | sh -s --
    ```

  - You have [containerd enabled](https://docs.docker.com/engine/storage/containerd/) for Docker Engine.

### Overview

This section walks you through extracting Software Bill of Materials (SBOMs) from a C++ Docker image using Docker Scout. SBOMs provide a detailed list of all the components in a software package, including their versions and licenses. You can use SBOMs to track the provenance of your software and ensure that it complies with your organization's security and licensing policies.

### Generate an SBOM

Here we will use the Docker image that we built in the [Create a multi-stage build for your C++ application](/guides/language/cpp/multistage/) guide. If you haven't already built the image, follow the steps in that guide to build the image.
The image is named `hello`. To generate an SBOM for the `hello` image, run the following command:

```bash
$ docker scout sbom --format list hello
```

The command will say "No packages discovered". This is because the final image is a scratch image and doesn't have any packages.

### Generate an SBOM attestation

The SBOM can be generated during the build process and "attached" to the image. This is called an SBOM attestation.
To generate an SBOM attestation for the `hello` image, first let's change the Dockerfile:

```Dockerfile
ARG BUILDKIT_SBOM_SCAN_STAGE=true

FROM ubuntu:latest AS build

RUN apt-get update && apt-get install -y build-essential

WORKDIR /app

COPY hello.cpp .

RUN g++ -o hello hello.cpp -static

# --------------------
FROM scratch

COPY --from=build /app/hello /hello

CMD ["/hello"]
```

The first line `ARG BUILDKIT_SBOM_SCAN_STAGE=true` enables SBOM scanning in the build stage.
Now, build the image with the following command:

```bash
$ docker buildx build --sbom=true -t hello:sbom .
```

This command will build the image and generate an SBOM attestation. You can verify that the SBOM is attached to the image by running the following command:

```bash
$ docker scout sbom --format list hello:sbom
```

Docker Scout reads the SBOM attestation when one is available, so this command reports packages from the build-stage metadata instead of indexing only the final scratch image filesystem.

### Summary

In this section, you learned how to generate SBOM attestation for a C++ Docker image during the build process.
Image scanners that inspect only the final filesystem may not identify packages in scratch images.
Use SBOM attestations to preserve package metadata from the build.
