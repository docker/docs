---
title: Container Device Interface (CDI)
weight: 60
description: Using CDI to access GPUs and other devices in your builds
keywords: build, buildkit, buildx, guide, tutorial, cdi, device, gpu, nvidia, cuda, amd, rocm
---

<!-- vale Docker.We = NO -->

The [Container Device Interface (CDI)](https://github.com/cncf-tags/container-device-interface/blob/main/SPEC.md)
is a specification designed to standardize how devices (like GPUs, FPGAs, and
other hardware accelerators) are exposed to and used by containers. The aim is
to provide a more consistent and secure mechanism for using hardware devices in
containerized environments, addressing the challenges associated with
device-specific setups and configurations.

In addition to enabling the container to interact with the device node, CDI also
lets you specify additional configuration for the device, such as environment
variables, host mounts (such as shared objects), and executable hooks.

## Getting started

To get started with CDI, you need to have a compatible environment set up. This
includes having Docker v27+ installed with [CDI configured](/reference/cli/dockerd.md#configure-cdi-devices)
and Buildx v0.22+.

You also need to create the [device specifications using JSON or YAML files](https://github.com/cncf-tags/container-device-interface/blob/main/SPEC.md#cdi-json-specification)
in one of the following locations:

* `/etc/cdi`
* `/var/run/cdi`
* `/etc/buildkit/cdi`

> [!NOTE]
> Location can be changed by setting the `specDirs` option in the `cdi` section
> of the [`buildkitd.toml` configuration file](../buildkit/configure.md) if you
> are using BuildKit directly. If you're building using the Docker Daemon with
> the `docker` driver, see [Configure CDI devices](/reference/cli/dockerd.md#configure-cdi-devices)
> documentation.

> [!NOTE]
> If you are creating a container builder on WSL, you need to ensure that
> [Docker Desktop](../../desktop/_index.md) is installed and [WSL 2 GPU Paravirtualization](../../desktop/features/gpu.md#prerequisites)
> is enabled. Buildx v0.27+ is also required to mount the WSL libraries in the
> container.

## Building with a simple CDI specification

Let's start with a simple CDI specification that injects an environment variable
into the build environment and write it to `/etc/cdi/foo.yaml`:

```yaml {title="/etc/cdi/foo.yaml"}
cdiVersion: "0.6.0"
kind: "vendor1.com/device"
devices:
- name: foo
  containerEdits:
    env:
    - FOO=injected
```

Inspect the `default` builder to verify that `vendor1.com/device` is detected
as a device:

```console
$ docker buildx inspect
Name:   default
Driver: docker

Nodes:
Name:             default
Endpoint:         default
Status:           running
BuildKit version: v0.23.2
Platforms:        linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/amd64/v4, linux/386
Labels:
 org.mobyproject.buildkit.worker.moby.host-gateway-ip: 172.17.0.1
Devices:
 Name:                  vendor1.com/device=foo
 Automatically allowed: false
GC Policy rule#0:
 All:            false
 Filters:        type==source.local,type==exec.cachemount,type==source.git.checkout
 Keep Duration:  48h0m0s
 Max Used Space: 658.9MiB
GC Policy rule#1:
 All:            false
 Keep Duration:  1440h0m0s
 Reserved Space: 4.657GiB
 Max Used Space: 953.7MiB
 Min Free Space: 2.794GiB
GC Policy rule#2:
 All:            false
 Reserved Space: 4.657GiB
 Max Used Space: 953.7MiB
 Min Free Space: 2.794GiB
GC Policy rule#3:
 All:            true
 Reserved Space: 4.657GiB
 Max Used Space: 953.7MiB
 Min Free Space: 2.794GiB
```

Now let's create a Dockerfile to use this device:

```dockerfile
# syntax=docker/dockerfile:1-labs
FROM busybox
RUN --device=vendor1.com/device \
  env | grep ^FOO=
```

Here we use the [`RUN --device` command](/reference/dockerfile.md#run---device)
and set `vendor1.com/device` which requests the first device available in the
specification. In this case it uses `foo`, which is the first device in
`/etc/cdi/foo.yaml`.

> [!NOTE]
> [`RUN --device` command](/reference/dockerfile.md#run---device) is only
> featured in [`labs` channel](../buildkit/frontend.md#labs-channel) since
> [Dockerfile frontend v1.14.0-labs](../buildkit/dockerfile-release-notes.md#1140-labs)
> and not yet available in stable syntax.

Now let's build this Dockerfile:

```console
$ docker buildx build .
[+] Building 0.4s (5/5) FINISHED                                                                                                        docker:default
 => [internal] load build definition from Dockerfile                                                                                    0.0s 
 => => transferring dockerfile: 155B                                                                                                    0.0s
 => resolve image config for docker-image://docker/dockerfile:1-labs                                                                    0.1s 
 => CACHED docker-image://docker/dockerfile:1-labs@sha256:9187104f31e3a002a8a6a3209ea1f937fb7486c093cbbde1e14b0fa0d7e4f1b5              0.0s
 => [internal] load metadata for docker.io/library/busybox:latest                                                                       0.1s 
 => [internal] load .dockerignore                                                                                                       0.0s
 => => transferring context: 2B                                                                                                         0.0s 
ERROR: failed to build: failed to solve: failed to load LLB: device vendor1.com/device=foo is requested by the build but not allowed
```

It fails because the device `vendor1.com/device=foo` is not automatically
allowed by the build as shown in the `buildx inspect` output above:

```text
Devices:
 Name:                  vendor1.com/device=foo
 Automatically allowed: false
```

To allow the device, you can use the [`--allow` flag](/reference/cli/docker/buildx/build.md#allow)
with the `docker buildx build` command:

```console
$ docker buildx build --allow device .
```

Or you can set the `org.mobyproject.buildkit.device.autoallow` annotation in
the CDI specification to automatically allow the device for all builds:

```yaml {title="/etc/cdi/foo.yaml"}
cdiVersion: "0.6.0"
kind: "vendor1.com/device"
devices:
- name: foo
  containerEdits:
    env:
    - FOO=injected
annotations:
  org.mobyproject.buildkit.device.autoallow: true
```

Now running the build again with the `--allow device` flag:

```console
$ docker buildx build --progress=plain --allow device .
#0 building with "default" instance using docker driver

#1 [internal] load build definition from Dockerfile
#1 transferring dockerfile: 159B done
#1 DONE 0.0s

#2 resolve image config for docker-image://docker/dockerfile:1-labs
#2 DONE 0.1s

#3 docker-image://docker/dockerfile:1-labs@sha256:9187104f31e3a002a8a6a3209ea1f937fb7486c093cbbde1e14b0fa0d7e4f1b5
#3 CACHED

#4 [internal] load metadata for docker.io/library/busybox:latest
#4 DONE 0.1s

#5 [internal] load .dockerignore
#5 transferring context: 2B done
#5 DONE 0.0s

#6 [1/2] FROM docker.io/library/busybox:latest@sha256:f85340bf132ae937d2c2a763b8335c9bab35d6e8293f70f606b9c6178d84f42b
#6 CACHED

#7 [2/2] RUN --device=vendor1.com/device   env | grep ^FOO=
#7 0.155 FOO=injected
#7 DONE 0.2s
```

The build is successful and the output shows that the `FOO` environment variable
was injected into the build environment as specified in the CDI specification.

## Set up a container builder with GPU support

In this section, we will show you how to set up a [container builder](../builders/drivers/docker-container.md)
using NVIDIA GPUs. Since Buildx v0.22, when creating a new container builder, a
GPU request is automatically added to the container builder if the host has GPU
drivers installed in the kernel. This is similar to using [`--gpus=all` with the `docker run`](/reference/cli/docker/container/run.md#gpus)
command.

> [!NOTE]
> We made a specially crafted BuildKit image because the current BuildKit
> release image is based on Alpine that doesn’t support NVIDIA drivers. The
> following image is based on Ubuntu and installs the NVIDIA client libraries
> and generates the CDI specification for your GPU in the container builder if
> a device is requested during a build. This image is temporarily hosted on
> Docker Hub under `crazymax/buildkit:v0.23.2-ubuntu-nvidia`.

Now let's create a container builder named `gpubuilder` using Buildx:

```console
$ docker buildx create --name gpubuilder --driver-opt "image=crazymax/buildkit:v0.23.2-ubuntu-nvidia" --bootstrap
#1 [internal] booting buildkit
#1 pulling image crazymax/buildkit:v0.23.2-ubuntu-nvidia
#1 pulling image crazymax/buildkit:v0.23.2-ubuntu-nvidia 1.0s done
#1 creating container buildx_buildkit_gpubuilder0
#1 creating container buildx_buildkit_gpubuilder0 8.8s done
#1 DONE 9.8s
gpubuilder
```

Let's inspect this builder:

```console
$ docker buildx inspect gpubuilder
Name:          gpubuilder
Driver:        docker-container
Last Activity: 2025-07-10 08:18:09 +0000 UTC

Nodes:
Name:                  gpubuilder0
Endpoint:              unix:///var/run/docker.sock
Driver Options:        image="crazymax/buildkit:v0.23.2-ubuntu-nvidia"
Status:                running
BuildKit daemon flags: --allow-insecure-entitlement=network.host
BuildKit version:      v0.23.2
Platforms:             linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/arm64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/arm/v7, linux/arm/v6
Labels:
 org.mobyproject.buildkit.worker.executor:         oci
 org.mobyproject.buildkit.worker.hostname:         d6aa9cbe8462
 org.mobyproject.buildkit.worker.network:          host
 org.mobyproject.buildkit.worker.oci.process-mode: sandbox
 org.mobyproject.buildkit.worker.selinux.enabled:  false
 org.mobyproject.buildkit.worker.snapshotter:      overlayfs
Devices:
 Name:      nvidia.com/gpu
 On-Demand: true
GC Policy rule#0:
 All:            false
 Filters:        type==source.local,type==exec.cachemount,type==source.git.checkout
 Keep Duration:  48h0m0s
 Max Used Space: 488.3MiB
GC Policy rule#1:
 All:            false
 Keep Duration:  1440h0m0s
 Reserved Space: 9.313GiB
 Max Used Space: 93.13GiB
 Min Free Space: 188.1GiB
GC Policy rule#2:
 All:            false
 Reserved Space: 9.313GiB
 Max Used Space: 93.13GiB
 Min Free Space: 188.1GiB
GC Policy rule#3:
 All:            true
 Reserved Space: 9.313GiB
 Max Used Space: 93.13GiB
 Min Free Space: 188.1GiB
```

We can see `nvidia.com/gpu` vendor is detected as a device in the builder which
means that drivers were detected.

Optionally you can check if NVIDIA GPU devices are available in the container
using `nvidia-smi`:

```console
$ docker exec -it buildx_buildkit_gpubuilder0 nvidia-smi -L
GPU 0: Tesla T4 (UUID: GPU-6cf00fa7-59ac-16f2-3e83-d24ccdc56f84)
```

## Building with GPU support

Let's create a simple Dockerfile that will use the GPU device:

```dockerfile
# syntax=docker/dockerfile:1-labs
FROM ubuntu
RUN --device=nvidia.com/gpu nvidia-smi -L
```

Now run the build using the `gpubuilder` builder we created earlier:

```console
$ docker buildx --builder gpubuilder build --progress=plain .
#0 building with "gpubuilder" instance using docker-container driver
...

#7 preparing device nvidia.com/gpu
#7 0.000 > apt-get update
...
#7 4.872 > apt-get install -y gpg
...
#7 10.16 Downloading NVIDIA GPG key
#7 10.21 > apt-get update
...
#7 12.15 > apt-get install -y --no-install-recommends nvidia-container-toolkit-base
...
#7 17.80 time="2025-04-15T08:58:16Z" level=info msg="Generated CDI spec with version 0.8.0"
#7 DONE 17.8s

#8 [2/2] RUN --device=nvidia.com/gpu nvidia-smi -L
#8 0.527 GPU 0: Tesla T4 (UUID: GPU-6cf00fa7-59ac-16f2-3e83-d24ccdc56f84)
#8 DONE 1.6s
```

As you might have noticed, the step `#7` is preparing the `nvidia.com/gpu`
device by installing client libraries and the toolkit to generate the CDI
specifications for the GPU.

The `nvidia-smi -L` command is then executed in the container using the GPU
device. The output shows the GPU UUID.

You can check the generated CDI specification within the container builder with
the following command:

```console
$ docker exec -it buildx_buildkit_gpubuilder0 cat /etc/cdi/nvidia.yaml
```

For the EC2 instance [`g4dn.xlarge`](https://aws.amazon.com/ec2/instance-types/g4/)
used here, it looks like this:

```yaml {collapse=true}
cdiVersion: 0.6.0
containerEdits:
  deviceNodes:
  - path: /dev/nvidia-modeset
  - path: /dev/nvidia-uvm
  - path: /dev/nvidia-uvm-tools
  - path: /dev/nvidiactl
  env:
  - NVIDIA_VISIBLE_DEVICES=void
  hooks:
  - args:
    - nvidia-cdi-hook
    - create-symlinks
    - --link
    - ../libnvidia-allocator.so.1::/usr/lib/x86_64-linux-gnu/gbm/nvidia-drm_gbm.so
    hookName: createContainer
    path: /usr/bin/nvidia-cdi-hook
  - args:
    - nvidia-cdi-hook
    - create-symlinks
    - --link
    - libcuda.so.1::/usr/lib/x86_64-linux-gnu/libcuda.so
    hookName: createContainer
    path: /usr/bin/nvidia-cdi-hook
  - args:
    - nvidia-cdi-hook
    - enable-cuda-compat
    - --host-driver-version=570.133.20
    hookName: createContainer
    path: /usr/bin/nvidia-cdi-hook
  - args:
    - nvidia-cdi-hook
    - update-ldcache
    - --folder
    - /usr/lib/x86_64-linux-gnu
    hookName: createContainer
    path: /usr/bin/nvidia-cdi-hook
  mounts:
  - containerPath: /run/nvidia-persistenced/socket
    hostPath: /run/nvidia-persistenced/socket
    options:
    - ro
    - nosuid
    - nodev
    - bind
    - noexec
  - containerPath: /usr/bin/nvidia-cuda-mps-control
    hostPath: /usr/bin/nvidia-cuda-mps-control
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/bin/nvidia-cuda-mps-server
    hostPath: /usr/bin/nvidia-cuda-mps-server
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/bin/nvidia-debugdump
    hostPath: /usr/bin/nvidia-debugdump
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/bin/nvidia-persistenced
    hostPath: /usr/bin/nvidia-persistenced
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/bin/nvidia-smi
    hostPath: /usr/bin/nvidia-smi
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libcuda.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libcuda.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libcudadebugger.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libcudadebugger.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-allocator.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-allocator.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-cfg.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-cfg.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-gpucomp.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-gpucomp.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-ml.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-nscq.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-nscq.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-nvvm.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-nvvm.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-opencl.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-opencl.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-pkcs11-openssl3.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-pkcs11-openssl3.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-pkcs11.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-pkcs11.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/x86_64-linux-gnu/libnvidia-ptxjitcompiler.so.570.133.20
    hostPath: /usr/lib/x86_64-linux-gnu/libnvidia-ptxjitcompiler.so.570.133.20
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /lib/firmware/nvidia/570.133.20/gsp_ga10x.bin
    hostPath: /lib/firmware/nvidia/570.133.20/gsp_ga10x.bin
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /lib/firmware/nvidia/570.133.20/gsp_tu10x.bin
    hostPath: /lib/firmware/nvidia/570.133.20/gsp_tu10x.bin
    options:
    - ro
    - nosuid
    - nodev
    - bind
devices:
- containerEdits:
    deviceNodes:
    - path: /dev/nvidia0
  name: "0"
- containerEdits:
    deviceNodes:
    - path: /dev/nvidia0
  name: GPU-6cf00fa7-59ac-16f2-3e83-d24ccdc56f84
- containerEdits:
    deviceNodes:
    - path: /dev/nvidia0
  name: all
kind: nvidia.com/gpu
```

Congrats on your first build using a GPU device with BuildKit and CDI.
