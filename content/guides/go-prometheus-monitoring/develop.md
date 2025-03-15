---
title: Developing your application
linkTitle: Develop your app
weight: 40
keywords: go, golang, containerize, initialize
description: Learn how to develop the Golang application with Docker.
aliases:
  - /guides/go-prometheus-monitoring/develop/
---

In the last section, we saw how using Docker Compose, we can connect our services together. In this section, we will learn how to develop the Golang application with Docker.

Now, if we make any changes to our golang application, it needs to reflect in the container, right? To do that, one approach is use --build flag in Docker Compose. It will rebuild all the services which have `build` instruction in the `compose.yml` file. This is how you can use it:

```
docker compose up --build
```

But, this is not the best approach. It will rebuild all the services which have `build` instruction. This is not efficient. Every time you make a change in the code, we need to rebuild manually. This is not is not very good flow for development. Another approach is to use Docker Compose Watch. In the `compose.yml` file, under the service `api`, we have added the `develop` section. So, it's more like a hot reloading. Whenever we make changes to code (defined in `path`), it will rebuild the image (or restart depending on the action). This is how you can use it:

```yaml {hl_lines="17-20",linenos=true}
services:
  api:
    container_name: go-api
    build:
      context: .
      dockerfile: Dockerfile
    image: go-api:latest
    ports:
      - 8000:8000
    networks:
      - go-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 5
    develop:
      watch:
        - path: .
          action: rebuild
```

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

Now, if you modify your `main.go` you will see the changes in real time without re-building the image.

