---
description: Learn about the windowsfilter storage driver
keywords: container, storage, driver, windows, windowsfilter
title: windowsfilter storage driver
---

The windowsfilter storage driver is the default storage driver for Docker
Engine on Windows. The windowsfilter driver uses Windows-native file system
layers to for storing Docker layers and volume data on disk. The windowsfilter
storage driver only works on file systems formatted with NTFS.

## Configure the windowsfilter storage driver

For most use case, no configuring the windowsfilter storage driver is not
necessary.

The default storage limit for Docker Engine on Windows is 127GB. To use a
different storage size, set the `size` option for the windowsfilter storage
driver. See [windowsfilter options](/reference/cli/dockerd.md#windowsfilter-options).

Data is stored on the Docker host in `image` and `windowsfilter` subdirectories
within `C:\ProgramData\docker` by default. You can change the storage location
by configuring the `data-root` option in the [Daemon configuration file](/reference/cli/dockerd.md#on-windows):

```json
{
  "data-root": "d:\\docker"
}
```

You must restart the daemon for the configuration change to take effect.

## Additional information

For more information about how container storage works on Windows, refer to
Microsoft's [Containers on Windows documentation](https://learn.microsoft.com/en-us/virtualization/windowscontainers/manage-containers/container-storage).
