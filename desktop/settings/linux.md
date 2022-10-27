---
description: Docker Desktop settings
keywords: settings, preferences, proxy, file sharing, resources, kubernetes, Docker Desktop, Linux
title: Change Docker Desktop settings on Linux
---

This page provides information on how to configure and manage your Docker Desktop settings.

To navigate to **Settings** either:

- Select the Docker menu ![whale menu](../images/whale-x.svg){: .inline} and then **Settings**
- Select the **Settings** icon from the Docker Dashboard.

## General

On the **General** tab, you can configure when to start Docker and specify other settings:

- **Start Docker Desktop when you log in**. Select to automatically start Docker
  Desktop when you log into your machine.

- **Choose Theme for Docker Desktop**. Choose whether you want to apply a **Light** or **Dark** theme to Docker Desktop. Alternatively you can set Docker Desktop to **Use System Settings**.

- **Use integrated container terminal**. Select to execute commands in a running container straight from the Docker Dashboard. For more information, see [Explore containers](../use-desktop/container.md)

- **Send usage statistics**. Select so Docker Desktop sends diagnostics,
  crash reports, and usage data. This information helps Docker improve and
  troubleshoot the application. Clear the check box to opt out. Docker may
  periodically prompt you for more information.

- **Show weekly tips**. Select to display useful advice and suggestions about
  using Docker.

- **Open Docker Desktop dashboard at startup**. Select to automatically open the
  dashboard when starting Docker Desktop.

- **Use Docker Compose V2**. Select to enable the `docker-compose` command to
  use Docker Compose V2. For more information, see [Docker Compose V2](../../compose/compose-v2/index.md).

## Resources

The **Resources** tab allows you to configure CPU, memory, disk, proxies,
network, and other resources.

### Advanced

On the **Advanced** tab, you can limit resources available to Docker.

Advanced settings are:

- **CPUs**. By default, Docker Desktop is set to use half the number of processors
  available on the host machine. To increase processing power, set this to a
  higher number; to decrease, lower the number.

- **Memory**. By default, Docker Desktop is set to use 25% of your host's
  memory. To increase the RAM, set this to a higher number; to decrease it,
  lower the number.

- **Swap**. Configure swap file size as needed. The default is 1 GB.

- **Disk image size**. Specify the size of the disk image.

- **Disk image location**. Specify the location of the Linux volume where containers and images are stored.

You can also move the disk image to a different location. If you attempt to move a disk image to a location that already has one, you are asked if you want to use the existing image or replace it.

### File sharing

Use File sharing to allow local directories on your machine to be shared with
Linux containers. This is especially useful for editing source code in an IDE on
the host while running and testing the code in a container.

By default the `/home/<user>` directory is shared.
If your project is outside this directory then it must be added to the list,
otherwise you may get `Mounts denied` or `cannot start service` errors at runtime.

File share settings are:

- **Add a Directory**. Click `+` and navigate to the directory you want to add.

- **Remove a Directory**. Click `-` next to the directory you want to remove

- **Apply & Restart** makes the directory available to containers using Docker's
  bind mount (`-v`) feature.

> Tips on shared folders, permissions, and volume mounts
>
> * Share only the directories that you need with the container. File sharing
>   introduces overhead as any changes to the files on the host need to be notified
>   to the Linux VM. Sharing too many files can lead to high CPU load and slow
>   filesystem performance.
> * Shared folders are designed to allow application code to be edited
>   on the host while being executed in containers. For non-code items
>   such as cache directories or databases, the performance will be much
>   better if they are stored in the Linux VM, using a [data volume](../../storage/volumes.md)
>   (named volume) or [data container](../../storage/volumes.md).

### Proxies

To configure HTTP proxies, switch on the **Manual proxy configuration** setting.

Your proxy settings, however, are not propagated into the containers you start.
If you wish to set the proxy settings for your containers, you need to define
environment variables for them, just like you would do on Linux, for example:

```console
$ docker run -e HTTP_PROXY=http://proxy.example.com:3128 alpine env

PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=b7edf988b2b5
TERM=xterm
HOME=/root
HTTP_PROXY=http://proxy.example.com:3128
```

For more information on configuring the Docker CLI to automatically set proxy variables for both `docker run` and `docker build`
see [Configure the Docker client](../../network/proxy.md#configure-the-docker-client).

### Network

Docker Desktop uses a private IPv4 network for internal services such as a DNS server and an HTTP proxy. In case the choice of subnet clashes with something in your environment, specify a custom subnet using the **Network** setting.

### Docker Engine

The **Docker Engine** tab allows you to configure the Docker daemon to determine how your containers run.

Type a JSON configuration file in the box to configure the daemon settings. For a full list of options, see the Docker Engine
[dockerd commandline reference](/engine/reference/commandline/dockerd/){:target="_blank" rel="noopener" class="_"}.

Click **Apply & Restart** to save your settings and restart Docker Desktop.

## Beta Features

{% include beta.md %}

From the **Beta features** tab, you can sign up to the [Developer Preview program](https://www.docker.com/community/get-involved/developer-preview/).

On the **Beta features** tab, you also have the option to allow version 4.13 feature flags, which are product features Docker is currently experimenting with. This is switched on by default. 

## Kubernetes

Docker Desktop includes a standalone Kubernetes server, so that you can test
deploying your Docker workloads on Kubernetes. To enable Kubernetes support and
install a standalone instance of Kubernetes running as a Docker container,
select **Enable Kubernetes**.

Select **Show system containers (advanced)** to view internal containers when
using Docker commands.

Select **Reset Kubernetes cluster** to delete all stacks and Kubernetes resources.

For more information about using the Kubernetes integration with Docker Desktop,
see [Deploy on Kubernetes](../kubernetes.md){:target="_blank" rel="noopener" class="_"}.

## Software Updates

The **Software Updates** tab notifies you of any updates available to Docker Desktop.
When there's a new update,
click the **Release Notes** option to learn what's included in the updated version.

Turn off the check for updates by clearing the **Automatically check for updates**
check box. This disables notifications in the Docker menu and the notification
badge that appears on the Docker Dashboard. To check for updates manually, select
the **Check for updates** option in the Docker menu.

To allow Docker Desktop to automatically download new updates in the background,
select **Always download updates**. This downloads newer versions of Docker Desktop
when an update becomes available. After downloading the update, click
**Apply and Restart** to install the update. You can do this either through the
Docker menu or in the **Updates** section in the Docker Dashboard.

## Extensions 

Use the **Extensions** tab to:

- **Enable Docker Extensions**
- **Allow only extensions distributed through the Docker Marketplace**
- **Show Docker Extensions system containers**

For more information about Docker extensions, see [Extensions](../extensions/index.md).
