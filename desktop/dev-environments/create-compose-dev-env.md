---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local, compose
title: Create an advanced dev environment
---

Create an advanced dev environment such as a microservice with a server, proxy and DB. 

As with a simple dev environment, you can create a more advanced dev environment from a:
- Git repository
- Branch or tag of a Git repository
- Subfolder of a Git repository
- Local folder

> **Note**
>
> When cloning a Git repository using SSH, ensure you've added your SSH key to the ssh-agent. To do this, open a terminal and run `ssh-add <path to your private ssh key>`.

## Create a Compose Dev Environment

The example below, taken from the `compose-dev-env` project from the [Docker Samples](https://github.com/dockersamples/compose-dev-env){:target="_blank" rel="noopener" class="_"} GitHub repository, demonstrates how to create a Compose Dev Environment from a Git repository. 

>Note
>
>Currently, Dev Environments is not able to detect the main language of the subdirectory. You need to define your own base image or compose services in a `compose-dev.yaml` located in your subdirectory.
>
>For more information on how to configure, see the [React application with a Spring backend and a MySQL database sample](https://github.com/docker/awesome-compose/tree/master/react-java-mysql) or the [Go server with an Nginx proxy and a Postgres database sample](https://github.com/docker/awesome-compose/tree/master/nginx-golang-postgres). 

1. From **Dev Environments**, select **Create**. The **Create a Dev Environment** dialog displays. 
2. Click **Get Started** and then copy `https://github.com/dockersamples/compose-dev-env.git` and add it to the **Enter the Git Repository** field with **Existing Git repo** as the source.
3. Click **Continue**. This initializes the project, clones the Git code, and builds the Compose application. This:

    - Builds local images for services that are defined in the Compose file
    - Pulls images required for other services
    - Creates volumes and networks
    - Starts the Compose stack

Once your application is up and running, you can check by opening [http://localhost:8080](http://localhost:8080) in your browser.

The time taken to start the Compose application depends on how your application is configured, whether the images have been built, and the number of services you have defined, for example.

Note that VS Code doesn't open directly, unlike a simple Dev Environment, as there are multiple services configured. You can hover over a service and then click on the **Open in VS Code** button to open a specific service in VS Code. This stops the existing container and creates a new container which allows you to develop and update your service in VS Code.

You can now update your service and test it against your Compose application.

## Set up your own dev environment

>**Changes to Dev Environments with Docker Desktop 4.13**
>
>Docker has simplified how you configure your dev environment project. All you need to get started is a `compose-dev.yaml` file. If you have an existing project with a `.docker/` folder this is automatically migrated the next time you launch.
>
> If you are using `.docker/docker-compose.yaml`, we move it to `../compose-dev.yaml`.
>If you are using `.docker/config.json`, we create a `../compose-dev.yaml` file with a single service named "app”. It is configured to use the image or Dockerfile referenced in the JSON as a starting point.
{: .important}

To set up a dev environment, there are additional configuration steps to tell Docker Desktop how to build, start, and use the right image for your services.

Dev Environments use an `compose-dev.yaml` file located in the at the root of your project. This file allows you to define the image required for a dedicated service, the ports you'd like to expose, along with additional configuration options dedicated to Dev Environments coming in the future.

Take a detailed look at the `compose-dev.yaml` file used in the [compose-dev-env](https://github.com/dockersamples/compose-dev-env/blob/main/.docker/docker-compose.yaml){:target="_blank" rel="noopener" class="_"} sample project.

```yaml
version: "3.7"
services:
  backend:
    build:
      context: backend
      target: development
    secrets:
      - db-password
    depends_on:
      - db
  db:
    image: mariadb
    restart: always
    healthcheck:
      test: [ "CMD", "mysqladmin", "ping", "-h", "127.0.0.1", "--silent" ]
      interval: 3s
      retries: 5
      start_period: 30s
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/mysql
    environment:
      - MYSQL_DATABASE=example
      - MYSQL_ROOT_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 3306
  proxy:
    build: proxy
    ports:
      - 8080:80
    depends_on:
      - backend
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

In the yaml file, the build context `backend` specifies that that the container should be built using the `development` stage (`target` attribute) of the Dockerfile located in the `backend` directory (`context` attribute)

The `development` stage of the Dockerfile is defined as follows:

```dockerfile
FROM golang:1.16-alpine AS build
WORKDIR /go/src/github.com/org/repo
COPY . .

RUN go build -o server .

FROM build AS development
RUN apk update \
    && apk add git
CMD ["go", "run", "main.go"]

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/org/repo/server /server
CMD ["/server"]
```

The `development`target uses a `golang:1.16-alpine` image with all dependencies you need for development. You can start your project directly from VS Code and interact with the others applications or services such as the database or the frontend.

In the example, the Docker Compose files are the same. However, they could be different and the services defined in the main Compose file may use other targets to build or directly reference other images.

## What's next?

Learn how to [distribute your dev environment](share.md)
