---
description: Configuring the Docker daemon
keywords: docker, daemon, configuration
title: Docker daemon configuration overview
linkTitle: Daemon
weight: 60
aliases:
  - /articles/chef/
  - /articles/configuring/
  - /articles/dsc/
  - /articles/puppet/
  - /config/thirdparty/
  - /config/thirdparty/ansible/
  - /config/thirdparty/chef/
  - /config/thirdparty/dsc/
  - /config/thirdparty/puppet/
  - /engine/admin/
  - /engine/admin/ansible/
  - /engine/admin/chef/
  - /engine/admin/configuring/
  - /engine/admin/dsc/
  - /engine/admin/puppet/
  - /engine/articles/chef/
  - /engine/articles/configuring/
  - /engine/articles/dsc/
  - /engine/articles/puppet/
  - /engine/userguide/
  - /config/daemon/
---

This page shows you how to customize the Docker daemon, `dockerd`.

> [!NOTE]
>
> This page is for users who've installed Docker Engine manually. If you're
> using Docker Desktop, refer to the [settings page](/manuals/desktop/settings-and-maintenance/settings.md#docker-engine).

## Configure the Docker daemon

There are two ways to configure the Docker daemon:

- Use a JSON configuration file. This is the preferred option, since it keeps
  all configurations in a single place.
- Use flags when starting `dockerd`.

You can use both of these options together as long as you don't specify the same
option both as a flag and in the JSON file. If that happens, the Docker daemon
won't start and prints an error message.

### Configuration file

The following table shows the location where the Docker daemon expects to find
the configuration file by default, depending on your system and how you're
running the daemon.

| OS and configuration | File location                              |
| -------------------- | ------------------------------------------ |
| Linux, regular setup | `/etc/docker/daemon.json`                  |
| Linux, rootless mode | `~/.config/docker/daemon.json`             |
| Windows              | `C:\ProgramData\docker\config\daemon.json` |

For rootless mode, the daemon respects the `XDG_CONFIG_HOME` variable. If set,
the expected file location is `$XDG_CONFIG_HOME/docker/daemon.json`.

You can also explicitly specify the location of the configuration file on
startup, using the `dockerd --config-file` flag.

Learn about the available configuration options in the
[dockerd reference docs](/reference/cli/dockerd.md#daemon-configuration-file)

### Configuration using flags

You can also start the Docker daemon manually and configure it using flags.
This can be useful for troubleshooting problems.

Here's an example of how to manually start the Docker daemon, using the same
configurations as shown in the previous JSON configuration:

```console
$ dockerd --debug \
  --tls=true \
  --tlscert=/var/docker/server.pem \
  --tlskey=/var/docker/serverkey.pem \
  --host tcp://192.168.59.3:2376
```

#### Listen on a Unix socket on Windows

Starting with Docker Engine 29.5, you can configure `dockerd` on Windows to
listen on a Unix socket. This setting is optional. The default Windows endpoint
remains the `npipe` named pipe.

Run `dockerd` from an elevated PowerShell session:

```powershell
PS C:\> dockerd -H unix://C:/Users/<USER>/docker.sock --group docker-users
```

Replace `<USER>` with the name of your Windows user profile directory.
Without `--group`, only members of the Administrators group and the `SYSTEM`
account can access the socket. Use a comma-separated list to grant access to
multiple users or groups.

Use a path in the user profile when granting access with `--group`. A socket
created directly under `C:\` can inherit a high integrity level that blocks
access for the configured users and groups.

Configure the Docker client to use the same socket:

```powershell
PS C:\> docker -H unix://C:/Users/<USER>/docker.sock info
```

Learn about the available configuration options in the
[dockerd reference docs](/reference/cli/dockerd.md), or by
running:

```console
$ dockerd --help
```

## Daemon data directory

The Docker daemon persists all data in a single directory. This tracks
everything related to Docker, including containers, images, volumes, service
definition, and secrets.

By default the daemon stores data in:

- `/var/lib/docker` on Linux
- `C:\ProgramData\docker` on Windows

When using the [containerd image store](/manuals/engine/storage/containerd.md)
(the default for Docker Engine 29.0 and later on fresh installations), image
contents and container snapshots are stored in `/var/lib/containerd`. Other
daemon data (volumes, configs) remains in `/var/lib/docker`.

When using [classic storage drivers](/manuals/engine/storage/drivers/_index.md)
like `overlay2` (the default for upgraded installations), all data is stored in
`/var/lib/docker`.

### Configure the data directory location

You can configure the Docker daemon to use a different storage directory using
the `data-root` configuration option.

```json
{
  "data-root": "/mnt/docker-data"
}
```

The `data-root` option does not affect image and container data stored in
`/var/lib/containerd` when using the containerd image store. To change the
storage location of containerd snapshotters, use the system containerd
configuration file:

```toml {title="/etc/containerd/config.toml"}
version = 2
root = "/mnt/containerd-data"
```

Make sure you use a dedicated directory for each daemon. If two daemons share
the same directory, for example an NFS share, you will experience errors that
are difficult to troubleshoot.

## Next steps

Many specific configuration options are discussed throughout the Docker
documentation. Some places to go next include:

- [Automatically start containers](/manuals/engine/containers/start-containers-automatically.md)
- [Run containerd in the Docker daemon](/manuals/engine/daemon/embedded-containerd.md)
- [Limit a container's resources](/manuals/engine/containers/resource_constraints.md)
- [Configure storage drivers](/manuals/engine/storage/drivers/select-storage-driver.md)
- [Container security](/manuals/engine/security/_index.md)
- [Configure the Docker daemon to use a proxy](./proxy.md)
