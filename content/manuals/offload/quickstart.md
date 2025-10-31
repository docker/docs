---
title: Docker Offload quickstart
linktitle: Quickstart
weight: 10
description: Learn how to use Docker Offload to build and run your container images faster, both locally and in CI.
keywords: cloud, quickstart, Docker Desktop, gpu, offload
---

{{< summary-bar feature_name="Docker Offload" >}}

[Docker Offload](./about.md) lets you build and run containers in the cloud while using your local Docker Desktop tools
and workflow. This means faster builds, access to powerful cloud resources like GPUs, and a seamless development
experience.

This quickstart covers the steps developers need to get started with Docker Offload.

> [!NOTE]
>
> If you're an organization owner, to get started you must [sign up](https://www.docker.com/products/docker-offload/)
> and subscribe your organization to use Docker Offload. After subscribing, see the following:
>
> - [Manage Docker products](../admin/organization/manage-products.md) to learn how to manage access for the developers
> in your organization.
> - [Usage and billing](./usage.md) to learn how set up billing and monitor usage.

## Prerequisites

- You must have [Docker Desktop](/desktop/) installed. Docker Offload works with Docker Desktop version 4.50 or later.
- You must have access to Docker Offload. Your organization owner must [sign
  up](https://www.docker.com/products/docker-offload/) your organization.
- You must have committed usage available or on-demand usage enabled for your organization. This is set up by your
  organization owner. For more details, see [Docker Offload usage and billing](/offload/usage/).

## Step 1: Verify access to Docker Offload

To access Docker Offload, you must be part of an organization that has subscribed to Docker Offload. As a developer, you
can verify this by checking if the Docker Offload toggle appears in the Docker Desktop Dashboard header.

1. Start Docker Desktop and sign in.
2. In the Docker Desktop Dashboard header, look for the Docker Offload toggle.

![Offload toggle](./images/offload-toggle.png)

If you see the Docker Offload toggle, you have access to Docker Offload and can proceed to the next step. If you don't
see the Docker Offload toggle, check if Docker Offload is disabled in your [Docker Desktop
settings](./configuration.md), and then contact your administrator to verify that your organization has subscribed to
Docker Offload and that they have enabled access for your organization.

## Step 2: Start Docker Offload

You can start Docker Offload from the CLI or in the header of the Docker Desktop Dashboard. The following steps describe
how to start Docker Offload using the CLI.

1. Start Docker Desktop and sign in.
2. Open a terminal and run the following command to start Docker Offload:

   ```console
   $ docker offload start
   ```

   > [!TIP]
   >
   > To learn more about the Docker Offload CLI commands, see the [Docker Offload CLI
   > reference](/reference/cli/docker/offload/).

3. If you are a member of multiple organizations that have access to Docker Offload, you have the option to select a
   profile. The selected organization is responsible for any usage.


4. If GPU-accelerated instances are enabled for your organization, you have the option to select **GPU Supported**, or
   **CPU-only support**.

When Docker Offload is started, you'll see a cloud icon
({{< inline-image src="./images/cloud-mode.png" alt="Offload mode icon" >}})
in the Docker Desktop Dashboard header, and the Docker Desktop Dashboard appears purple. You can run
`docker offload status` in a terminal to check the status of Docker Offload.

## Step 3: Run a container with Docker Offload

After starting Docker Offload, Docker Desktop connects to a secure cloud environment
that mirrors your local experience. When you run builds or containers, they
execute remotely, but behave just like local ones.

To verify that Docker Offload is working, run a container:

```console
$ docker run --rm hello-world
```

If you enabled GPU support, you can also run a GPU-enabled container:

```console
$ docker run --rm --gpus all hello-world
```

If Docker Offload is working, you'll see `Hello from Docker!` in the terminal output.

## Step 4: Monitor your Offload usage

When Docker Offload is started and you have started session (for example, you've ran a container), then you can see
current session duration estimate in the Docker Desktop Dashboard footer next to the hourglass icon
({{< inline-image src="./images/hourglass-icon.png" alt="Offload session duration" >}}).

Also, when Docker Offload is started, you can view detailed session information by selecting **Offload** > **Insights**
in the left navigation of the Docker Desktop Dashboard.

## Step 5: Stop Docker Offload

Docker Offload automatically [idles](./configuration.md#understand-active-and-idle-states) after a period of inactivity. You can stop it at
any time. To stop Docker Offload:

```console
$ docker offload stop
```

When you stop Docker Offload, the cloud environment is terminated and all running containers and images are removed.
When Docker Offload has been idle for around 5 minutes, the environment is also terminated and all running containers
and images are removed.

To start Docker Offload again, run the `docker offload start` command.

## What's next

Configure your idle timeout in Docker Desktop. For more information, see [Configure Docker Offload](./configuration.md).