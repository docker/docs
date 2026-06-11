---
title: R language-specific guide
linkTitle: R
description: Containerize R apps using Docker
keywords: Docker, getting started, R, language
summary: |
  This guide details how to containerize R applications using Docker.
aliases:
  - /languages/r/
  - /guides/languages/r/
  - /language/R/build-images/
  - /language/R/run-containers/
  - /language/r/containerize/
  - /language/r/develop/
  - /language/r/configure-ci-cd/
  - /language/r/deploy/
  - /guides/r/configure-ci-cd/
  - /guides/r/containerize/
  - /guides/r/deploy/
  - /guides/r/develop/
params:
  tags: [cicd]
  time: 10 minutes
---


The R language-specific guide teaches you how to containerize a R application using Docker. In this guide, you’ll learn how to:

- Containerize and run a R application
- Set up a local environment to develop a R application using containers
- Configure a CI/CD pipeline for a containerized R application using GitHub Actions
- Deploy your containerized R application locally to Kubernetes to test and debug your deployment

Start by containerizing an existing R application.

## Containerize a R application

### Prerequisites

- You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

### Overview

This section walks you through containerizing and running a R application.

### Get the sample application

The sample application uses the popular [Shiny](https://shiny.posit.co/) framework.

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/mfranzon/r-docker-dev.git && cd r-docker-dev
```

You should now have the following contents in your `r-docker-dev`
directory.

```text
├── r-docker-dev/
│ ├── src/
│ │ └── app.R
│ ├── src_db/
│ │ └── app_db.R
│ ├── compose.yaml
│ ├── Dockerfile
│ └── README.md
```

To learn more about the files in the repository, see the following:

- [Dockerfile](/reference/dockerfile.md)
- [.dockerignore](/reference/dockerfile.md#dockerignore-file)
- [compose.yaml](/reference/compose-file/_index.md)

### Run the application

Inside the `r-docker-dev` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3838](http://localhost:3838). You should see a simple Shiny application.

In the terminal, press `ctrl`+`c` to stop the application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `r-docker-dev` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:3838](http://localhost:3838).

You should see a simple Shiny application.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/).

### Summary

In this section, you learned how you can containerize and run your R
application using Docker.

Related information:

- [Docker Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll learn how you can develop your application using
containers.

## Use containers for R development

### Prerequisites

Complete [Containerize a R application](containerize.md).

### Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Adding a local database and persisting data
- Configuring Compose to automatically update your running Compose services as you edit and save your code

### Get the sample application

You'll need to clone a new repository to get a sample application that includes logic to connect to the database.

Change to a directory where you want to clone the repository and run the following command.

```console
$ git clone https://github.com/mfranzon/r-docker-dev.git
```

### Configure the application to use the database

To try the connection between the Shiny application and the local database you have to modify the `Dockerfile` changing the `COPY` instruction:

```diff
-COPY src/ .
+COPY src_db/ .
```

### Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor.

In the `compose.yaml` file, you need to un-comment the properties for configuring the database. You must also mount the database password file and set an environment variable on the `shiny-app` service pointing to the location of the file in the container.

The following is the updated `compose.yaml` file.

```yaml
services:
  shiny-app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 3838:3838
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
  db:
    image: postgres:18
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

> [!NOTE]
>
> To learn more about the instructions in the Compose file, see [Compose file
> reference](/reference/compose-file/).

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```text
mysecretpassword
```

Save and close the `password.txt` file.

You should now have the following contents in your `r-docker-dev`
directory.

```text
├── r-docker-dev/
│ ├── db/
│ │ └── password.txt
│ ├── src/
│ │ └── app.R
│ ├── src_db/
│ │ └── app_db.R
│ ├── requirements.txt
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── Dockerfile
│ └── README.md
```

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

Now test your DB connection opening a browser at:

```console
http://localhost:3838
```

You should see a pop-up message:

```text
DB CONNECTED
```

Press `ctrl+c` in the terminal to stop your application.

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](/manuals/compose/how-tos/file-watch.md).

Lines 15 to 18 in the `compose.yaml` file contain properties that trigger Docker
to rebuild the image when a file in the current working directory is changed:

```yaml {hl_lines="15-18",linenos=true}
services:
  shiny-app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 3838:3838
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    develop:
      watch:
        - action: rebuild
          path: .
  db:
    image: postgres:18
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

Now, if you modify your `app.R` you will see the changes in real time without re-building the image!

Press `ctrl+c` in the terminal to stop your application.

### Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose file watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)

### Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

## Configure CI/CD for your R application

### Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a R application](containerize.md). You must have a [GitHub](https://github.com/signup) account and a verified [Docker](https://hub.docker.com/signup) account to complete this section.

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

3. In the editor window, copy and paste the following YAML configuration.

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

In this section, you learned how to set up a GitHub Actions workflow for your R application.

Related information:

- [Introduction to GitHub Actions](/guides/gha.md)
- [Docker Build GitHub Actions](/manuals/build/ci/github-actions/_index.md)
- [Workflow syntax for GitHub Actions](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

### Next steps

Next, learn how you can locally test and debug your workloads on Kubernetes before deploying.

## Test your R deployment

### Prerequisites

- Complete all the previous sections of this guide, starting with [Containerize a R application](containerize.md).
- [Turn on Kubernetes](/manuals/desktop/use-desktop/kubernetes.md#enable-kubernetes) in Docker Desktop.

### Overview

In this section, you'll learn how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine. This allows you to test and debug your workloads on Kubernetes locally before deploying.

### Create a Kubernetes YAML file

In your `r-docker-dev` directory, create a file named
`docker-r-kubernetes.yaml`. Open the file in an IDE or text editor and add
the following contents. Replace `DOCKER_USERNAME/REPO_NAME` with your Docker
username and the name of the repository that you created in [Configure CI/CD for
your R application](configure-ci-cd.md).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-r-demo
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      service: shiny
  template:
    metadata:
      labels:
        service: shiny
    spec:
      containers:
        - name: shiny-service
          image: DOCKER_USERNAME/REPO_NAME
          imagePullPolicy: Always
          env:
            - name: POSTGRES_PASSWORD
              value: mysecretpassword
---
apiVersion: v1
kind: Service
metadata:
  name: service-entrypoint
  namespace: default
spec:
  type: NodePort
  selector:
    service: shiny
  ports:
    - port: 3838
      targetPort: 3838
      nodePort: 30001
```

In this Kubernetes YAML file, there are two objects, separated by the `---`:

- A Deployment, describing a scalable group of identical pods. In this case,
  you'll get just one replica, or copy of your pod. That pod, which is
  described under `template`, has just one container in it. The
  container is created from the image built by GitHub Actions in [Configure CI/CD for
  your R application](configure-ci-cd.md).
- A NodePort service, which will route traffic from port 30001 on your host to
  port 3838 inside the pods it routes to, allowing you to reach your app
  from the network.

To learn more about Kubernetes objects, see the [Kubernetes documentation](https://kubernetes.io/docs/home/).

### Deploy and check your application

1. In a terminal, navigate to `r-docker-dev` and deploy your application to
   Kubernetes.

   ```console
   $ kubectl apply -f docker-r-kubernetes.yaml
   ```

   You should see output that looks like the following, indicating your Kubernetes objects were created successfully.

   ```text
   deployment.apps/docker-r-demo created
   service/service-entrypoint created
   ```

2. Make sure everything worked by listing your deployments.

   ```console
   $ kubectl get deployments
   ```

   Your deployment should be listed as follows:

   ```shell
   NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
   docker-r-demo   1/1     1            1           15s
   ```

   This indicates all one of the pods you asked for in your YAML are up and running. Do the same check for your services.

   ```console
   $ kubectl get services
   ```

   You should get output like the following.

   ```shell
   NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
   kubernetes           ClusterIP   10.96.0.1       <none>        443/TCP          23h
   service-entrypoint   NodePort    10.99.128.230   <none>        3838:30001/TCP   75s
   ```

   In addition to the default `kubernetes` service, you can see your `service-entrypoint` service, accepting traffic on port 30001/TCP.

3. In a browser, visit the following address. Note that a database was not deployed in
   this example.

   ```console
   http://localhost:30001/
   ```

4. Run the following command to tear down your application.

   ```console
   $ kubectl delete -f docker-r-kubernetes.yaml
   ```

### Summary

In this section, you learned how to use Docker Desktop to deploy your application to a fully-featured Kubernetes environment on your development machine.

Related information:

- [Kubernetes documentation](https://kubernetes.io/docs/home/)
- [Deploy on Kubernetes with Docker Desktop](/manuals/desktop/use-desktop/kubernetes.md)
- [Swarm mode overview](/manuals/engine/swarm/_index.md)
