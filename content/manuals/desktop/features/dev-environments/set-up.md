---
description: Set up a dev Environments
keywords: Dev Environments, share, set up, Compose, Docker Desktop
title: Set up a dev environment
weight: 20
aliases:
- /desktop/dev-environments/set-up/
---

{{% include "dev-envs-changing.md" %}}

>**Changes to Dev Environments with Docker Desktop 4.13**
>
>Docker has simplified how you configure your dev environment project. All you need to get started is a `compose-dev.yaml` file. If you have an existing project with a `.docker/` folder this is automatically migrated the next time you launch.
>
> If you are using `.docker/docker-compose.yaml`, we move it to `../compose-dev.yaml`.
>If you are using `.docker/config.json`, we create a `../compose-dev.yaml` file with a single service named "app”. It is configured to use the image or Dockerfile referenced in the JSON as a starting point.

To set up a dev environment, there are additional configuration steps to tell Docker Desktop how to build, start, and use the right image for your services.

Dev Environments use an `compose-dev.yaml` file located at the root of your project. This file allows you to define the image required for a dedicated service, the ports you'd like to expose, along with additional configuration options.

The following is an example `compose-dev.yaml` file.

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
# syntax=docker/dockerfile:1
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

The `development` target uses a `golang:1.16-alpine` image with all dependencies you need for development. You can start your project directly from VS Code and interact with the others applications or services such as the database or the frontend.

In the example, the Docker Compose files are the same. However, they could be different and the services defined in the main Compose file may use other targets to build or directly reference other images.

## What's next?

Learn how to [distribute your dev environment](share.md)
