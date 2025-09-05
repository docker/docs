---
title: Docker Offload quickstart
linktitle: Quickstart
weight: 10
description: Learn how to use Docker Offload to run your container images faster.
keywords: cloud, quickstart, cloud mode, Docker Desktop, GPU support
---

{{< summary-bar feature_name="Docker Offload" >}}

This quickstart helps you get started with Docker Offload. Docker Offload lets
you run container images faster by offloading resource-intensive tasks
to the cloud. It provides a cloud-based environment that mirrors your local
Docker Desktop experience.

To use Docker Offload, you must have:

   - [A Docker account](/accounts/create-account/)
   - [Docker Offload enabled](./configuration.md#offload-access). For users that
     are part of an organization with a Team or Business subscription, Docker
     Offload is disabled by default. An organization owner must enable Docker
     Offload for the organization.

## Step 1: Start Docker Offload

To start Docker Offload:

1. Start Docker Desktop and sign in.
2. Open a terminal and run the following command to start Docker Offload:

   ```console
   $ docker offload start
   ```

3. When prompted, select your account to use for Docker Offload. This account
   will consume credits for your Docker Offload usage.

4. When prompted, select whether to enable GPU support. If you choose to enable
   GPU support, Docker Offload will run in an instance with an NVIDIA L4 GPU,
   which is useful for machine learning or compute-intensive workloads.

   > [!NOTE]
   >
   > Enabling GPU support consumes more budget. For more details, see [Docker
   > Offload usage](/offload/usage/).

When Docker Offload is started, you'll see a cloud icon ({{< inline-image
src="./images/cloud-mode.png" alt="Offload mode icon" >}})
in the Docker Desktop Dashboard header, and the Docker Desktop Dashboard appears purple.
You can run `docker offload status` in a terminal to check the status of
Docker Offload.

## Step 2: Run a container with Docker Offload

After starting Docker Offload, Docker Desktop connects to a secure cloud
environment that mirrors your local experience. When you run containers, they
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

## Step 3: Stop Docker Offload

When you're done using Docker Offload, you can stop it. When stopped, you run
containers locally.

```console
$ docker offload stop
```

To start Docker Offload again, run the `docker offload start` command.

## What's next

- [Configure Docker Offload](configuration.md).
- Try [Docker Model Runner](../ai/model-runner/_index.md) or
  [Compose](../ai/compose/models-and-compose.md) to run AI models using Docker Offload.