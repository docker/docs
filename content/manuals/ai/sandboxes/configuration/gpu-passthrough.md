---
title: Enable NVIDIA GPU passthrough
linkTitle: GPU passthrough
description: Run GPU-accelerated workloads in a sandbox by building and installing the NVIDIA driver bundle for VFIO passthrough.
keywords: docker sandboxes, sbx, gpu, gpu passthrough, nvidia, vfio, iommu, driver bundle, cuda
weight: 30
aliases:
  - /ai/sandboxes/gpu-passthrough/
---

> [!IMPORTANT]
> GPU passthrough is experimental. The `--gpu` flag, the driver bundle, and the
> setup steps on this page are subject to change.

GPU passthrough in local Docker Sandboxes uses
[VFIO](https://www.kernel.org/doc/html/latest/driver-api/vfio.html) to assign a
physical NVIDIA GPU directly to the sandbox. The host unbinds the GPU from its
driver, and the sandbox drives the hardware itself.

## Requirements

- A Linux host with an NVIDIA GPU. Only NVIDIA GPUs and a single GPU per
  sandbox are supported.
- An `x86_64` or `arm64` CPU. On `arm64`, the kernel must use 4k or 64k pages.
  Kernels with 16k pages aren't supported.
- A GPU that nothing else is using: a headless host, or an additional GPU.
- IOMMU turned on in the BIOS or UEFI.
- The `iommufd` and `vfio_pci` kernel modules loaded:

  ```console
  sudo modprobe -a iommufd vfio_pci
  ```

  To load them at every boot, see
  [Troubleshooting](#sandbox-creation-fails-with-host-not-ready-for-vfio).

- The NVIDIA driver bundle, an EROFS image at
  `/usr/libexec/nerdbox-nvidia-bundle.erofs`. On `arm64` hosts with 64k pages,
  the file name is `nerdbox-nvidia-bundle_64k.erofs`. Docker Sandboxes doesn't
  include the bundle. See [Build the bundle](#build-the-bundle).

## Turn on the feature

The `--gpu` flag is hidden until you turn on experimental features and the GPU
feature flag:

```console
sbx settings set platform.allowExperimentalFeatures true
sbx settings set feature.sandbox-gpu true
```

## Build the bundle

The bundle contains the `nvidia` and `nvidia-uvm` kernel modules built for the
Docker Sandboxes guest kernel, plus the NVIDIA userspace libraries, firmware,
and `nvidia-cdi-hook`.

Each Docker Sandboxes release publishes an archive named
`nerdbox-nvidia-modules-<arch>.zip`. The archive holds the kernel modules, the
driver version they were built against (`VERSION`), and a build script. Use
`<arch>` as follows:

| Host               | `<arch>`    |
| ------------------ | ----------- |
| `x86_64`           | `x86_64`    |
| `arm64`, 4k pages  | `arm64`     |
| `arm64`, 64k pages | `arm64_64k` |

The script runs a container for the host architecture. The container downloads
the matching NVIDIA driver and the pinned `nvidia-cdi-hook` release, verifies
the checksum, builds the bundle, and installs it to `/usr/libexec`. You need
Docker and network access to `download.nvidia.com` and `github.com`.

Download and unpack the archive, then run the script as root:

```console
ARCH=$(uname -m | sed 's/aarch64/arm64/')
[ "$ARCH" = arm64 ] && [ "$(getconf PAGESIZE)" = 65536 ] && ARCH=arm64_64k
curl -fSLO https://github.com/docker/sbx-releases/releases/latest/download/nerdbox-nvidia-modules-$ARCH.zip
unzip nerdbox-nvidia-modules-$ARCH.zip -d nvidia-modules
cd nvidia-modules
sudo ./prepare-nvidia-bundle.sh
```

If the script can't write to the output path, it leaves the bundle in the
current directory and prints the `install` command to finish the job.

The script reads three environment variables:

- `OUTPUT`: where to write the bundle. The default is
  `/usr/libexec/<bundle file>`.
- `ACCEPT_NVIDIA_LICENSE`: set to `1` to accept the NVIDIA license without a
  prompt, for example in CI.
- `TARGET_PAGE_SIZE`: `4k` or `64k`. On `arm64`, the bundle is named for the
  page size in the archive's `VERSION` file. Set this to override it, for
  example if you created the sandbox with `sbx create --gpu --page-size 4k` on
  a 64k host.

Each Docker Sandboxes release uses a specific guest kernel, and the kernel
modules in the bundle must match it. After you upgrade Docker Sandboxes,
download the new archive and rebuild the bundle.

### Install the bundle on another host

If you built the bundle on another host with the same architecture and page
size, copy the bundle file to `/usr/libexec` on the GPU host.

On hosts where `/usr` is read-only, such as Fedora Silverblue or bootc-based
systems, put the bundle in a writable directory such as `/var/lib/sbx`. Add that
directory to the `SAILOR_PATH` environment variable of the daemon, then restart
the daemon.

## Run a sandbox with a GPU

Pass the `--gpu` flag when you create the sandbox:

```console
sbx create --gpu claude .
```

`sbx run --gpu claude` takes the same flag. The flag takes effect only when the
sandbox is created. Passing it when you re-attach to an existing sandbox has no
effect.

## Use the GPU in containers

Containers that you start inside the sandbox can use the GPU through the NVIDIA
CDI hook. Pass `--gpus all` or `--device nvidia.com/gpu=all` to `docker run`:

```console
docker run --rm --gpus all ubuntu:24.04 nvidia-smi -L
```

## Troubleshooting

### Sandbox creation fails with host not ready for VFIO

Creating a sandbox with `--gpu` fails with an error like this:

```text
host not ready for VFIO: /dev/iommu (iommufd not loaded), vfio-pci driver not registered; fix: sudo modprobe -a iommufd vfio_pci
```

The host is missing the `iommufd` or `vfio_pci` kernel module. Load both, then
create the sandbox again:

```console
sudo modprobe -a iommufd vfio_pci
```

Modules loaded with `modprobe` don't persist across reboots. On systemd-based
distributions, list them in a file to load them at boot:

```console
printf 'iommufd\nvfio_pci\n' | sudo tee /etc/modules-load.d/sbx-gpu.conf
```

If the error persists, turn on IOMMU in the BIOS or UEFI.

### Sandbox fails with GPU in use

Creating or starting a sandbox with `--gpu` fails with an error like this:

```text
GPU 0000:01:00.0 is in use by kwin_wayland (pid 52911); stop these processes (for a desktop session, log out or stop the display manager) and try again
```

Docker Sandboxes detaches the GPU from the host driver. This fails while another
process, such as a desktop session or nvtop, holds the GPU open. Stop the
processes that the message lists, then start the sandbox again. For a desktop
session, log out or stop the display manager. For example, on a host that uses
SDDM:

```console
sudo systemctl stop sddm
```

The unit name depends on your distribution and desktop environment.

### `docker run --gpus all` reports no GPU vendor

The bundle was built without `nvidia-cdi-hook`. Rebuild it with the current
script from the latest release.

### The script reports a file not found next to this script

The script reads `nvidia.ko`, `nvidia-uvm.ko`, and `VERSION` from its own
directory. On `x86_64`, run it from the unzipped archive. On `arm64`, put the
matching `arm64` guest modules and their `VERSION` file beside the script.

### The script reports a file not found in the driver download

The extracted driver didn't contain an expected library or firmware file.
Confirm that the download completed. If your driver version names a library
differently, adjust `DRIVER_LIB_FAMILIES` in the script.

### GPU workloads fail after an upgrade

The guest kernel or the pinned driver version likely changed. Download the new
release's archive, run the script again, and reinstall the bundle.
