---
title: GPU access
weight: 40
description: How to access NVIDIA GPUs from a container
keywords: docker, GPU, NVIDIA, cuda, nvidia-smi, device, container toolkit
---

## Access an NVIDIA GPU

### Prerequisites

Visit the official [NVIDIA drivers page](https://www.nvidia.com/Download/index.aspx)
to download and install the proper drivers. Reboot your system once you have
done so.

Verify that your GPU is running and accessible.

### Install NVIDIA Container Toolkit

Follow the official NVIDIA Container Toolkit [installation instructions](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html).

### Expose GPUs for use

Include the `--gpus` flag when you start a container to access GPU resources.

To expose all available GPUs:

```console
$ docker run -it --rm --gpus all ubuntu nvidia-smi
```

The output looks like the following:

```text
+---------------------------------------------------------------------------------------+
| NVIDIA-SMI 535.288.01             Driver Version: 535.288.01   CUDA Version: 12.2     |
|-----------------------------------------+----------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
|                                         |                      |               MIG M. |
|=========================================+======================+======================|
|   0  NVIDIA L4                      Off | 00000000:31:00.0 Off |                    0 |
| N/A   40C    P0              27W /  72W |      0MiB / 23034MiB |      4%      Default |
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

The leftmost column in the GPU table shows the index of each GPU (`0` for the
NVIDIA L4 in the previous example). Use these index numbers to target specific GPUs
with the `device` option.

To expose a single GPU by index:

```console
$ docker run -it --rm --gpus device=0 ubuntu nvidia-smi
```

To expose a GPU by its UUID, first list UUIDs with `nvidia-smi -L`:

```console
$ nvidia-smi -L
GPU 0: NVIDIA L4 (UUID: GPU-3a23c669-1f69-c64e-cf85-44e9b07e7a2a)
```

Then pass the UUID to `--gpus`:

```console
$ docker run -it --rm --gpus device=GPU-3a23c669-1f69-c64e-cf85-44e9b07e7a2a ubuntu nvidia-smi
```

On systems with multiple GPUs, you can expose several by index. The `device`
value must be quoted because it contains a comma:

```console
$ docker run -it --rm --gpus '"device=0,2"' ubuntu nvidia-smi
```

This exposes the GPUs at index `0` and `2` — the first and third GPUs listed in
`nvidia-smi` output.

> [!NOTE]
>
> NVIDIA GPUs can only be accessed by systems running a single engine.

### Set NVIDIA capabilities

You can set capabilities manually. For example, on Ubuntu you can run the
following:

```console
$ docker run --gpus 'all,capabilities=utility' --rm ubuntu nvidia-smi
```

This enables the `utility` driver capability which adds the `nvidia-smi` tool to
the container.

Capabilities and other configurations can be set in images via environment
variables. For valid variables, see the
[nvidia-container-toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/docker-specialized.html)
documentation. These variables can be set in a Dockerfile.

You can also use CUDA images, which set these variables automatically. See the
official [CUDA images](https://catalog.ngc.nvidia.com/orgs/nvidia/containers/cuda)
NGC catalog page.
