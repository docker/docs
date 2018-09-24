---
title: Deploy an appliction package
description: Learn how to deploy an appliction package in UCP
keywords: ucp, swarm, kubernetes, application
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

An application package has one of these formats:

- **Three-file format**: Defined by a metadata.yml, a docker-compose.yml, and a settings.yml files inside a `my-app.dockerapp` folder. This is also called the folder format.
- **Single-file format**: Defined by a data from the previously three files concatenated in the order givem and separated by `---\n` in a single file named named 'my-app.dockerapp'.
- **Image format**: Defined by a Docker image in the engine store or exported as a tarball.

The docker-app binary lets a user render an application package to a Compose file using the settings values in the settings file or those specified by the user. This Compose file can then be deployed to a cluster running in Swarm mode or Kubernetes using `docker stack deploy` or to a single engine or Swarm classic cluster using `docker-compose up`.
