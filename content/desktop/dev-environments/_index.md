---
description: Dev Environments
keywords: Dev Environments, share, local, Compose
title: Overview of Dev Environments
---

{{< include "dev-envs-changing.md" >}}

> **Beta**
>
> The Dev Environments feature is currently in [Beta](../../release-lifecycle.md#beta). We recommend that you do not use this in production environments.
{ .experimental }

Dev Environments let you create a configurable developer environment with all the code and tools you need to quickly get up and running. 

It uses tools built into code editors that allows Docker to access code mounted into a container rather than on your local host. This isolates the tools, files and running services on your machine allowing multiple versions of them to exist side by side.

You can use Dev Environments through the intuitive GUI in Docker Dashboard or straight from your terminal with the new [`docker dev` CLI plugin](dev-cli.md).

![Dev environments tab in Docker Desktop](../images/dev-env.PNG)

## How does it work?

>**Changes to Dev Environments with Docker Desktop 4.13**
>
>Docker has simplified how you configure your dev environment project. All you need to get started is a `compose-dev.yaml` file. If you have an existing project with a `.docker/` folder this is automatically migrated the next time you launch.

Dev Environments is powered by [Docker Compose](/compose/). This allows Dev Environments to take advantage of all the benefits and features of Compose whilst adding an intuitive GUI where you can launch environments with the click of a button.

Every dev environment you want to run needs a `compose-dev.yaml` file which configures your application's services and lives in your project directory. You don't need to be an expert in Docker Compose or write a `compose-dev.yaml` file from scratch as Dev Environments creates a starter `compose-dev.yaml` files based on the main language in your project. 

You can also use the many [sample dev environments](https://github.com/docker/awesome-compose) as a starting point for how to integrate different services. Alternatively, see [Set up a dev environment](set-up.md) for more information. 

## What's next?

Learn how to:
- [Launch a dev environment](create-dev-env.md)
- [Set up a dev environment](set-up.md)
- [Distribute your dev environment](share.md)
