---
title: Containerize a Golang application
linkTitle: Containerize your app
weight: 20
keywords: go, golang, containerize, initialize
description: Learn how to containerize a Golang application.
---

Containerization helps you bundle the application and its dependencies into a single package called a container. This package can run on any platform without worrying about the environment. In this section, you will learn how to containerize a Golang application using Docker.

To containerize a Golang application, you first need to create a Dockerfile. The Dockerfile contains instructions to build and run the application in a container. Also, when creating a Dockerfile, you can follow different sets of best practices to optimize the image size and make it more secure.

## Creating a Dockerfile

Create a new file named `Dockerfile` in the root directory of your Golang application. The Dockerfile contains instructions to build and run the application in a container.

The following is a Dockerfile for a Golang application. You will also find this file in the `go-prometheus-monitoring` directory.

```dockerfile
# Use the official Golang image as the base
FROM golang:1.24-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /build

# Copy go.mod and go.sum files for dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application source
COPY . .

# Build the Go binary
RUN go build -o /app .

# Final lightweight stage
FROM alpine:3.17 AS final

# Copy the compiled binary from the builder stage
COPY --from=builder /app /bin/app

# Expose the application's port
EXPOSE 8000

# Run the application
CMD ["bin/app"]
```

## Understanding the Dockerfile

The Dockerfile consists of two stages:

1. **Build stage**: This stage uses the official Golang image as the base and sets the necessary environment variables. It also sets the working directory inside the container, copies the `go.mod` and `go.sum` files for dependency installation, downloads the dependencies, copies the entire application source, and builds the Go binary.

    We use the `golang:1.24-alpine` image as the base image for the build stage. The `CGO_ENABLED=0` environment variable disables CGO, which is useful for building static binaries. We also set the `GOOS` and `GOARCH` environment variables to `linux` and `amd64`, respectively, to build the binary for the Linux platform.

2. **Final stage**: This stage uses the official Alpine image as the base and copies the compiled binary from the build stage. It also exposes the application's port and runs the application.

    We use the `alpine:3.17` image as the base image for the final stage. We copy the compiled binary from the build stage to the final image. We expose the application's port using the `EXPOSE` instruction and run the application using the `CMD` instruction.

    Apart from the multi-stage build, the Dockerfile also follows best practices such as using the official images, setting the working directory, and copying only the necessary files to the final image. We can further optimize the Dockerfile by other best practices.

## Build the Docker image and run the application

One we have the Dockerfile, we can build the Docker image and run the application in a container.

To build the Docker image, run the following command in the terminal:

```console
$ docker build -t go-api:latest .
```

After building the image, you can run the application in a container using the following command:

```console
$ docker run -p 8000:8000 go-api:latest
```

The application will start running inside the container, and you can access it at `http://localhost:8000`. You can also check the running containers using the `docker ps` command.

```console
$ docker ps
```

## Summary

In this section, you learned how to containerize a Golang application using a Dockerfile. You created a multi-stage Dockerfile to build and run the application in a container. You also learned about best practices to optimize the Docker image size and make it more secure.

Related information:

 - [Dockerfile reference](/reference/dockerfile.md)
 - [.dockerignore file](/reference/dockerfile.md#dockerignore-file)

## Next steps

In the next section, we will learn how to use Docker Compose to connect and run multiple services together to monitor a Golang application with Prometheus and Grafana.
