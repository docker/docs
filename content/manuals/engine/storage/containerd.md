---
title: containerd image store with Docker Engine
linkTitle: containerd image store
weight: 50
keywords: containerd, snapshotters, image store, docker engine
description: Learn how to enable the containerd image store on Docker Engine
aliases:
  - /storage/containerd/
---

Starting with Docker Engine v29, Docker uses [`containerd`](https://containerd.io/)
for managing container storage and images. `containerd` is the industry-standard container runtime.

## Main benefits

containerd offers the following benefits:

- Shared maintenance: containerd is an open-source project maintained by a large community.
- Customizability: use [snapshotters](snapshotters.md)
  with unique characteristics, such as:
  - [stargz](https://github.com/containerd/stargz-snapshotter) for lazy-pulling images on container
  startup.
  - [nydus](https://github.com/containerd/nydus-snapshotter) or [dragonfly](https://github.com/dragonflyoss/nydus) for peer-to-peer image distribution.
- Portability: containerd is lighter and more modular.
- Multi-platform support: Build multi-platform images and images with attestations.
- WebAssembly: Ability to run Wasm containers.

For more information about the containerd image store and its benefits, refer to
[containerd image store on Docker Desktop](/manuals/desktop/features/containerd.md).

## Migrate to containerd image store on Docker Engine

{{< summary-bar feature_name="Containerd migration" >}}

You can enable the auto-migration feature for containerd snapshotters once you
update to Docker Engine v29. The migration mechanism handles overlay and vfs
images. To enable it:

1. Add the following configuration to your `/etc/docker/daemon.json`
   configuration file:

   ```json
   {
     "features": {
         "containerd-migration": true
     }
   }
   ```

1. Save the file.

1. Restart the daemon for the changes to take effect.

   Switching to containerd snapshotters causes you to temporarily lose images and
   containers created using the classic storage drivers.
   If you use other graph storage, repull or rebuild your images.

1. To display which driver you are using, run:

   ```console
   $ docker info -f '{{ .DriverStatus }}'
   ```

> [!TIP]
> Those resources still exist on your filesystem, and you can retrieve them by
> [turning off the containerd snapshotters feature](./drivers/_index.md#enable-legacy-storage-drivers).

## Related pages

- [Legacy storage drivers](drivers/_index.md)
