---
title: Containerize a Golang application
linkTitle: Containerize your app
weight: 10
keywords: go, golang, containerize, initialize
description: Learn how to containerize a Golang application.
aliases:
  - /guides/go-monitoring/containerize/
---

To containerize a Golang application, we first need to create a Dockerfile. The Dockerfile contains instructions to build and run the application in a container. Also, when creating a Dockerfile, we can follow different sets of best practices to optimize the image size and make it more secure.

## Creating a Dockerfile

```dockerfile
# Use the official Go image with version 1.23 on Alpine Linux as the build stage
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container for the build process
WORKDIR /build

# Copy Go module files to the container to set up dependencies
COPY go.mod go.sum ./

# Download the module dependencies
RUN go mod download

# Copy the rest of the application source code into the container
COPY . .

# Build the Go application and output the binary to /app
RUN go build -o /app .

# Use the official Alpine Linux image as the final lightweight runtime image
FROM alpine:3.17 AS final

# Expose port 8080 for the application
EXPOSE 8080

# Create a non-root user and group for better security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Switch to the non-root user to run the application with limited permissions
USER appuser

# Copy the built application binary from the builder stage to the runtime image
COPY --from=builder /app /bin/app

# Set the command to run the application binary
CMD ["bin/app"]
```

## Dockerfile explanation

The Dockerfile consists of two stages:

1. **Build stage**: This stage uses the official Golang image with version 1.23 on Alpine Linux. It sets up the working directory, copies the Go module files, downloads the module dependencies, copies the application source code, builds the application, and outputs the binary to `/app`.

2. **Final stage**: This stage uses the official Alpine Linux image as the final lightweight runtime image. It exposes port 8080 for the application, creates a non-root user and group for better security, switches to the non-root user, copies the built application binary from the builder stage to the runtime image, and sets the command to run the application binary.

By using multi-stage builds, we can keep the final image size small and secure. The build stage contains all the necessary tools and dependencies to build the application, while the final stage only includes the application binary and runtime dependencies.

Apart from the multi-stage build, the Dockerfile also follows best practices such as using the official images, setting the working directory, creating a non-root user, and copying only the necessary files to the final image.

## Build the Docker image and run the application

One we have the Dockerfile, we can build the Docker image and run the application in a container.

To build the Docker image, run the following command in the terminal:

```console
$ docker build -t go-monitoring .
```

After building the image, you can run the application in a container using the following command:

```console
$ docker run -p 8080:8080 go-monitoring
```

The application will start running inside the container, and you can access it at [http://localhost:8080](http://localhost:8080). To verify the application is working you can 

## Summary

In this section, you learned how to containerize a Golang application using a Dockerfile. You created a multi-stage Dockerfile to build and run the application in a container. By following best practices and using multi-stage builds, you can create optimized and secure Docker images for your applications.

Related information:

 - [Dockerfile reference](/reference/dockerfile.md)
 - [.dockerignore file](/reference/dockerfile.md#dockerignore-file)
 - [Docker Compose overview](/manuals/compose/_index.md)
 - [Compose file reference](/reference/compose-file/_index.md)

## Next steps

In the next section, you'll learn how you can develop your application using
containers.
