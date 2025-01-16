---
title: containerd image store with Docker Engine
linkTitle: containerd image store
weight: 50
keywords: containerd, snapshotters, image store, docker engine
description: Learn how to enable the containerd image store on Docker Engine
aliases:
  - /storage/containerd/
---

{{< summary-bar feature_name="containerd" >}}

containerd, the industry-standard container runtime, uses snapshotters instead
of the classic storage drivers for storing image and container data.
While the `overlay2` driver still remains the default driver for Docker Engine,
you can opt in to using containerd snapshotters as an experimental feature.

To learn more about the containerd image store and its benefits, refer to
[containerd image store on Docker Desktop](/manuals/desktop/features/containerd.md).

## Enable containerd image store on Docker Engine

Switching to containerd snapshotters causes you to temporarily lose images and
containers created using the classic storage drivers.
Those resources still exist on your filesystem, and you can retrieve them by
turning off the containerd snapshotters feature.

The following steps explain how to enable the containerd snapshotters feature.

1. Add the following configuration to your `/etc/docker/daemon.json`
   configuration file:

   ```json
   {
     "features": {
       "containerd-snapshotter": true
     }
   }
   ```

2. Save the file.
3. Restart the daemon for the changes to take effect.

   ```console
   $ sudo systemctl restart docker
   ```

After restarting the daemon, running `docker info` shows that you're using
containerd snapshotter storage drivers.

```console
$ docker info -f '{{ .DriverStatus }}'
[[driver-type io.containerd.snapshotter.v1]]
```

Docker Engine uses the `overlayfs` containerd snapshotter by default.
