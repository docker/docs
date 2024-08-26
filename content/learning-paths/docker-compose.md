---
title: Defining and Running Multi-Container Applications with Docker Compose
summary: Simplify the process of defining, configuring, and running multi-container Docker applications to enable efficient development, testing, and deployment.
description: Learn how to use Docker Compose to define and run multi-container Docker applications.
params:
  image: images/learning-paths/compose.png
  skill: Beginner
  time: 5 minutes
  prereq: None
---

{{< columns >}}

Developers face challenges with multi-container Docker applications, including
complex configuration, dependency management, and maintaining consistent
environments. Networking, resource allocation, data persistence, logging, and
monitoring add to the difficulty. Security concerns and troubleshooting issues
further complicate the process, requiring effective tools and practices for
efficient management.

Docker Compose solves the problem of managing multi-container Docker
applications by providing a simple way to define, configure, and run all the
containers needed for an application using a single YAML file. This approach
helps developers to easily set up, share, and maintain consistent development,
testing, and production environments, ensuring that complex applications can be
deployed with all their dependencies and services properly configured and
orchestrated.

<!-- break -->

## What you’ll learn

- What Docker Compose is and what it does
- How to define services
- Use cases for Docker Compose
- How things would be different without Docker Compose

## Who’s this for?

- Developers and DevOps engineers who need to define, manage, and orchestrate
  multi-container Docker applications efficiently across multiple environments.
- Development teams that want to increase productivity by streamlining
  development workflows and reducing setup time.

## Tools integration

Works well with Docker CLI, CI/CD tools, and container orchestration tools.

{{< /columns >}}

## Modules

{{< accordion large=true title=`Why Docker Compose?` icon=`play_circle` >}}

Docker Compose is an essential tool for defining and running multi-container
Docker applications. Docker Compose simplifies the Docker experience, making it
easier for developers to create, manage, and deploy applications by using YAML
files to configure application services.

Docker Compose provides several benefits:

- Lets you define multi-container applications in a single YAML file.
- Ensures consistent environments across development, testing, and production.
- Manages the startup and linking of multiple containers effortlessly.
- Streamlines development workflows and reduces setup time.
- Ensures that each service runs in its own container, avoiding conflicts.

**Duration**: 2.5 minutes

{{< youtube-embed 2EqarOM2V4U >}}

{{< /accordion >}}

{{< accordion large=true title=`Demo: Set up and use Docker Compose` icon=`play_circle` >}}

This Docker Compose demo shows how to orchestrate a multi-container application
environment, streamlining development and deployment processes.

- Compare Docker Compose to the `docker run` command
- Configure a multi-container web app using a compose.yml file
- Run a multi-container web app using one command

**Duration**: 2.5 minutes

{{< youtube-embed P5RBKmOLPH4 >}}

{{< /accordion >}}

{{< accordion large=true title=`Common challenges and questions` icon=`quiz` >}}

<!-- vale Docker.HeadingLength = NO -->

### Do I need to maintain a separate Compose file for my development, testing, and staging environments?

You don't necessarily need to maintain entirely separate Compose files for your
development, testing, and staging environments. You can define all your
services in a single Compose file (`compose.yml`). You can use profiles to
group service configurations specific to each environment (`dev`, `test`,
`staging`).

When you need to spin up an environment, you can activate the corresponding
profiles. For example, to set up the development environment:

```console
$ docker compose --profile dev up
```

This command starts only the services associated with the `dev` profile,
leaving the rest inactive.

For more information on using profiles, see [Using profiles with
Compose](/compose/profiles/).

### How can I enforce the database service to start up before the frontend service?

Docker Compose ensures services start in a specific order by using the
`depends_on` property. This tells Compose to start the database service before
even attempting to launch the frontend service. This is crucial since
applications often rely on databases being ready for connections.

However, `depends_on` only guarantees the order, not that the database is fully
initialized. For a more robust approach, especially if your application relies
on a prepared database (e.g., after migrations), consider [health
checks](/reference/compose-file/services.md#healthcheck). Here, you can
configure the frontend to wait until the database passes its health check
before starting. This ensures the database is not only up but also ready to
handle requests.

For more information on setting the startup order of your services, see
[Control startup and shutdown order in Compose](/compose/startup-order/).

### Can I use Compose to build a Docker image?

Yes, you can use Docker Compose to build Docker images. Docker Compose is a
tool for defining and running multi-container applications. Even if your
application isn't a multi-container application, Docker Compose can make it
easier to run by defining all the `docker run` options in a file.

To use Compose, you need a `compose.yml` file. In this file, you can specify
the build context and Dockerfile for each service. When you run the command
`docker compose up --build`, Docker Compose will build the images for each
service and then start the containers.

For more information on building Docker images using Compose, see the [Compose
Build Specification](/compose/compose-file/build/).

### What is the difference between Docker Compose and Dockerfile?

A Dockerfile provides instructions to build a container image while a Compose
file defines your running containers. Quite often, a Compose file references a
Dockerfile to build an image to use for a particular service.

### What is the difference between the `docker compose up` and `docker compose run` commands?

The `docker compose up` command creates and starts all your services. It's
perfect for launching your development environment or running the entire
application. The `docker compose run` command focuses on individual services.
It starts a specified service along with its dependencies, allowing you to run
tests or perform one-off tasks within that container.

<!-- vale Docker.HeadingLength = YES -->

{{< /accordion >}}

{{< accordion large=true title=`Resources` icon=`link` >}}

- [Overview of Docker Compose CLI](/compose/reference/)
- [Overview of Docker Compose](/compose/)
- [How Compose works](/compose/compose-application-model/)
- [Using profiles with Compose](/compose/profiles/)
- [Control startup and shutdown order with Compose](/compose/startup-order/)
- [Compose Build Specification](/compose/compose-file/build/)

{{< /accordion >}}
