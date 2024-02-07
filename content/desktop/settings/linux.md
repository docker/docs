---
description: Change your Docker Desktop settings on Linux
keywords: settings, preferences, proxy, file sharing, resources, kubernetes, Docker
  Desktop, Linux
title: Change Docker Desktop settings on Linux
---

This page provides information on how to configure and manage your Docker Desktop settings.

To navigate to **Settings** either:

- Select the Docker menu {{< inline-image src="images/desktop/whale-x.svg" alt="whale menu" >}} and then **Settings**
- Select the **Settings** icon from the Docker Dashboard.

You can also locate the `settings.json` file at `~/.docker/desktop/settings.json`.

## General

On the **General** tab, you can configure when to start Docker and specify other settings:

- **Start Docker Desktop when you sign in**. Select to automatically start Docker
  Desktop when you sign in to your machine.

- **Open Docker Dashboard when Docker Desktop starts**. Select to automatically open the
  dashboard when starting Docker Desktop.

- **Choose theme for Docker Desktop**. Choose whether you want to apply a **Light** or **Dark** theme to Docker Desktop. Alternatively you can set Docker Desktop to **Use System Settings**.

- **Use containerd for pulling and storing images**.
  Turns on the containerd image store.
  This brings new features like faster container startup performance by lazy-pulling images,
  and the ability to run Wasm applications with Docker.
  For more information, see [containerd image store](../containerd.md).

- **Choose container terminal**. Determines which terminal is launched when opening the terminal from a container.
If you choose the integrated terminal, you can run commands in a running container straight from the Docker Dashboard. For more information, see [Explore containers](../use-desktop/container.md).

- **Send usage statistics**. Select so Docker Desktop sends diagnostics,
  crash reports, and usage data. This information helps Docker improve and
  troubleshoot the application. Clear the check box to opt out. Docker may
  periodically prompt you for more information.

- **Show weekly tips**. Select to display useful advice and suggestions about
  using Docker.

- **Use Enhanced Container Isolation**. Select to enhance security by preventing containers from breaching the Linux VM. For more information, see [Enhanced Container Isolation](../hardened-desktop/enhanced-container-isolation/index.md)
    >**Note**
    >
    > This setting is only available if you are signed in to Docker Desktop and have a Docker Business subscription.

- **Show CLI hints**. Displays CLI hints and tips when running Docker commands in the CLI. This is turned on by default. To turn CLI hints on or off from the CLI, set `DOCKER_CLI_HINTS` to `true` or `false` respectively.

- **SBOM Indexing**. When this option is enabled, inspecting an image in Docker Desktop shows a **Start analysis** button that, when selected, analyzes the image with Docker Scout.

- **Enable background SBOM indexing**. When this option is enabled, Docker Scout automatically analyzes images that you build or pull.

## Resources

The **Resources** tab allows you to configure CPU, memory, disk, proxies,
network, and other resources.

### Advanced

On the **Advanced** tab, you can limit resources available to the Docker Linux VM.

Advanced settings are:

- **CPU limit**. Specify the maximum number of CPUs to be used by Docker Desktop.
  By default, Docker Desktop is set to use all the processors available on the host machine.

- **Memory limit**. By default, Docker Desktop is set to use up to 25% of your host's
  memory. To increase the RAM, set this to a higher number; to decrease it,
  lower the number.

- **Swap**. Configure swap file size as needed. The default is 1 GB.

- **Virtual disk limit**. Specify the maximum size of the disk image.

- **Disk image location**. Specify the location of the Linux volume where containers and images are stored.

  You can also move the disk image to a different location. If you attempt to
  move a disk image to a location that already has one, you are asked if you
  want to use the existing image or replace it.

>**Tip**
>
> If you feel Docker Desktop starting to get slow or you're running
> multi-container workloads, increase the memory and disk image space allocation
{ .tip }

- **Resource Saver**. Enable or disable [Resource Saver mode](../use-desktop/resource-saver.md),
  which significantly reduces CPU and memory utilization on the host by
  automatically turning off the Linux VM when Docker Desktop is idle (i.e., no
  containers are running).

  You can also configure the Resource Saver timeout which indicates how long
  should Docker Desktop be idle before Resource Saver mode kicks in. Default is
  5 minutes.

  >**Note**
  >
  > Exit from Resource Saver mode occurs automatically when containers run. Exit
  > may take a few seconds (~3 to 10 secs) as Docker Desktop restarts the Linux VM.

### File sharing

Use File sharing to allow local directories on your machine to be shared with
Linux containers. This is especially useful for editing source code in an IDE on
the host while running and testing the code in a container.

By default the `/home/<user>` directory is shared.
If your project is outside this directory then it must be added to the list,
otherwise you may get `Mounts denied` or `cannot start service` errors at runtime.

File share settings are:

- **Add a Directory**. Select `+` and navigate to the directory you want to add.

- **Remove a Directory**. Select `-` next to the directory you want to remove

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
{ .tip }

### Proxies

To configure HTTP proxies, switch on the **Manual proxy configuration** setting.
This setting is used for logging into Docker, for pulling and pushing images, and for
container Internet access. If the proxy requires authorization then Docker Desktop dynamically asks
the developer for a username and password. All passwords are stored securely in the OS credential store.
Note that only the `Basic` proxy authentication method is supported so we recommend using an `https://`
URL for your HTTP/HTTPS proxies to protect passwords while in transit on the network. Docker Desktop
supports TLS 1.3 when communicating with proxies.

To prevent developers from accidentally changing the proxy settings, see
[Settings Management](../hardened-desktop/settings-management/index.md#what-features-can-i-configure-with-settings-management).

### Network

{{< include "desktop-network-setting.md" >}}

## Docker Engine

The **Docker Engine** tab allows you to configure the Docker daemon used to run containers with Docker Desktop.

You configure the daemon using a JSON configuration file. Here's what the file might look like:

```json
{
  "builder": {
    "gc": {
      "defaultKeepStorage": "20GB",
      "enabled": true
    }
  },
  "experimental": false
}
```

You can find this file at `$HOME/.docker/daemon.json`. To change the configuration, either
edit the JSON configuration directly from the dashboard in Docker Desktop, or open and
edit the file using your favorite text editor.

> **Note**
>
> Only the **Docker Desktop** daemon uses the configuration file under `$HOME/.docker`.
> If you manually install Docker Engine alongside Docker Desktop, the manually
> installed instance uses a `daemon.json` configuration file in a  different location.
> Refer to [Configure the daemon](../../config/daemon/index.md) for more information
> about how to configure the Docker daemon on a manually installed Docker Engine.

To see the full list of possible configuration options, see the
[dockerd command reference](/engine/reference/commandline/dockerd/).

Select **Apply & Restart** to save your settings and restart Docker Desktop.

## Builders

{{< include "desktop-builders-setting.md" >}}

## Kubernetes

Docker Desktop includes a standalone Kubernetes server, so that you can test
deploying your Docker workloads on Kubernetes. To turn on Kubernetes support and
install a standalone instance of Kubernetes running as a Docker container,
select **Enable Kubernetes**.

Select **Show system containers (advanced)** to view internal containers when
using Docker commands.

Select **Reset Kubernetes cluster** to delete all stacks and Kubernetes resources.

For more information about using the Kubernetes integration with Docker Desktop,
see [Deploy on Kubernetes](../kubernetes.md).

## Software Updates

The **Software Updates** tab notifies you of any updates available to Docker Desktop.
When there's a new update,
select the **Release Notes** option to learn what's included in the updated version.

Turn off the check for updates by clearing the **Automatically check for updates**
check box. This disables notifications in the Docker menu and the notification
badge that appears on the Docker Dashboard. To check for updates manually, select
the **Check for updates** option in the Docker menu.

To allow Docker Desktop to automatically download new updates in the background,
select **Always download updates**. This downloads newer versions of Docker Desktop
when an update becomes available. After downloading the update, select
**Apply and Restart** to install the update. You can do this either through the
Docker menu or in the **Updates** section in the Docker Dashboard.

## Extensions

Use the **Extensions** tab to:

- **Enable Docker Extensions**
- **Allow only extensions distributed through the Docker Marketplace**
- **Show Docker Extensions system containers**

For more information about Docker extensions, see [Extensions](../extensions/index.md).

## Features in development

On the **Feature control** tab you can control your settings for **Beta features** and **Experimental features**.

You can also sign up to the [Developer Preview program](https://www.docker.com/community/get-involved/developer-preview/) from the **Features in development** tab.

### Beta features

{{< include "beta.md" >}}

### Experimental features

{{< include "desktop-experimental.md" >}}

## Notifications

{{< include "desktop-notifications-settings.md" >}}
