---
title: Deploy an appliction package
description: Learn how to deploy an appliction package in UCP
keywords: ucp, swarm, kubernetes, application
---

> Beta disclaimer
>
> This is beta content. It is not yet complete and should be considered a work in progress. This content is subject to change without notice.

In Docker 2.1 EE, an application package has one of these formats:

- **Three-file format**: Defined by a metadata.yml, a docker-compose.yml, and a settings.yml files inside a `my-app.dockerapp` folder. This is also called the folder format.
- **Single-file format**: Defined by a data from the previously three files concatenated in the order givem and separated by `---\n` in a single file named named 'my-app.dockerapp'.
- **Image forma   t**: Defined by a Docker image in the engine store or exported as a tarball.

The docker-app binary lets a user render an application package to a Compose file using the settings values in the settings file or those specified by the user. This Compose file can then be deployed to a cluster running in Swarm mode or Kubernetes using `docker stack deploy` or to a single engine or Swarm classic cluster using `docker-compose up`.

The docker-app binary also lets a user deploy an application package directly, essentially executing a `docker-app render` command followed by a `docker stack deploy` command.

Once a stack has been deployed, you must use the `docker stack` or `docker-compose` commands to manipulate and manage the stack.

## Creating a stack in the UCP web interface

To create a stack in the UCP web interface, follow these steps:

1. Go to the UCP web interface.
2. In the lefthand menu, first select **Shared Resources**, then **Stacks**.

    ![Create stacks in UCP](/ee/ucp/images/ucp-create-stack.png)

3. Select **Create Stack** to display the stack creation dialog.

    ![Configure stacks in UCP](/ee/ucp/images/ucp-config-stack.png)

4. Enter a name for the stack in the **Name** field.
5. Select either **Swarm Services** or **Kubernetes Workloads** for the orchestrator mode.
6. If you selected Swarm in the previous step, select either **Compose File** or **App Package** for the **Application File Mode**.
7. Select **Next**.

To specify a Compose file:

1. After selecting Swarm and Compose file, enter or upload your `docker-compose.yml` in **2. Add Application File**.

    ![Provide docker-compose.yml in UCP](/ee/ucp/images/ucp-stack-compose.png)

2. Select **Create**.


### Deploying on Swarm

[placeholder]

### Deploying on Kubernetes

[placeholder]
