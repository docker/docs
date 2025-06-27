---
title: Docker Cloud quickstart
linktitle: Quickstart
weight: 10
description: Learn how to use Docker Cloud to build and run your container images faster, both locally and in CI.
keywords: cloud, quickstart, cloud mode, Docker Desktop, GPU support, cloud builder, usage
---

{{< summary-bar feature_name="Docker Cloud" >}}

This quickstart helps you get started with Cloud mode for Docker Cloud.

## Cloud mode

Cloud mode uses cloud-based resources while maintaining your familiar
local Docker experience. When started, Docker Desktop offloads builds and
container runs to high-performance cloud infrastructure.

This mode is especially useful in virtual desktop environments (VDIs) where
nested virtualization isn't supported, or when your local machine doesn't have
enough CPU, memory, or GPU to handle container workloads. Docker Cloud runs
everything in the cloud while preserving a local-like experience through Docker
Desktop.

Despite running remotely, features like bind mounts and port forwarding continue
to work seamlessly, providing a local-like experience from within Docker Desktop
and the CLI.

### Step 1: Start Docker Cloud

1. Open the Docker Desktop Dashboard and sign in.
2. Toggle the **Start Cloud mode** switch ({{< inline-image
   src="./images/cloud-mode-stopped.png" alt="Cloud mode icon"
   >}}) in the header bar. 

   When Docker Cloud is started, you'll see a cloud icon ({{< inline-image
   src="./images/cloud-mode.png" alt="Cloud mode icon"
   >}}), the Docker Desktop Dashboard appears purple, and **Cloud running**
   appears in the left side of the footer.

   > [!NOTE]
   > If you don't see the Cloud mode switch, ensure you have the latest version of
   > Docker Desktop installed and that Docker Cloud is enabled in **Settings** >
   > **Beta features**.

3. In the Docker Desktop Dashboard, select the dropdown at the top of the left
   navigation to choose an account. Docker Cloud will consume Cloud credits from
   the selected account.

After enabling Docker Cloud, Docker Desktop connects to a secure cloud environment
that mirrors your local experience. When you run builds or containers, they
execute remotely, but behave just like local ones.

To verify that Docker Cloud is active, run a container:

```console
$ docker run hello-world
```

If Docker Cloud is working, you'll see `Hello from Docker!` in the terminal output.

### Step 2: Optional. Enable GPU support

GPU support is useful for machine learning or compute-intensive workloads. Once
enabled, your cloud containers run in an instance with an NVIDIA L4 GPU.

1. In the Docker Desktop Dashboard, navigate to **Settings**.
2. Select **Beta features**.
3. Ensure that the **Enable Docker Cloud** option is enabled.
4. Check the **Enable GPU support** option.
5. Select **Apply & restart** to enable GPU support.
6. Ensure that the Cloud mode switch is enabled in the Docker Desktop Dashboard
   header.

To verify that GPU support is active, run an NVIDIA CUDA container and the
nvidia-smi tool that lets you check GPU status. Use the following command in a
terminal:

```console
$ docker run --rm --gpus all nvidia/cuda:12.9.0-base-ubuntu24.04 nvidia-smi
```

If GPU support is working, you'll see output from nvidia-smi listing available
NVIDIA GPUs, including device names and driver versions.

For example, the output might look like this:

```text
+---------------------------------------------------------------------------------------+
| NVIDIA-SMI 535.247.01             Driver Version: 535.247.01   CUDA Version: 12.9     |
|-----------------------------------------+----------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
|                                         |                      |               MIG M. |
|=========================================+======================+======================|
|   0  NVIDIA L4                      Off | 00000000:31:00.0 Off |                    0 |
| N/A   47C    P0              30W /  72W |      0MiB / 23034MiB |      5%      Default |
|                                         |                      |                  N/A |
+-----------------------------------------+----------------------+----------------------+

+---------------------------------------------------------------------------------------+
| Processes:                                                                            |
|  GPU   GI   CI        PID   Type   Process name                            GPU Memory |
|        ID   ID                                                             Usage      |
|=======================================================================================|
|  No running processes found                                                           |
+---------------------------------------------------------------------------------------+
```

If GPU support is not enabled or available in your environment, you'll see an
error indicating that the GPU runtime could not be initialized.

## What's next

- [Configure Docker Cloud](configuration.md).
- [Use Docker Cloud to build images in CI environments](ci-build.md).