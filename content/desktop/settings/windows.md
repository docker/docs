---
description: Change your Docker Desktop settings on Windows
keywords: settings, preferences, proxy, file sharing, resources, kubernetes, Docker
  Desktop, Windows
title: Change Docker Desktop settings on Windows
---

This page provides information on how to configure and manage your Docker Desktop settings.

To navigate to **Settings** either:

- Select the Docker menu {{< inline-image src="../images/whale-x.svg" alt="whale menu" >}} and then **Settings**
- Select the **Settings** icon from the Docker Dashboard.

You can also locate the `settings.json` file at `C:\Users\[USERNAME]\AppData\Roaming\Docker\settings.json`.

## General

On the **General** tab, you can configure when to start Docker and specify other settings:

- **Start Docker Desktop when you log in**. Select to automatically start Docker
  Desktop when you sign in to your machine.

- **Choose Theme for Docker Desktop**. Choose whether you want to apply a **Light** or **Dark** theme to Docker Desktop. Alternatively you can set Docker Desktop to **Use System Settings**.

- **Use integrated container terminal**. Select to execute commands in a running container straight from the Docker Dashboard. For more information, see [Explore containers](../use-desktop/container.md).

- **Expose daemon on tcp://localhost:2375 without TLS**. Check this option to
  enable legacy clients to connect to the Docker daemon. You must use this option
  with caution as exposing the daemon without TLS can result in remote code
  execution attacks.

- **Use the WSL 2 based engine**. WSL 2 provides better performance than the
  Hyper-V backend. For more information, see [Docker Desktop WSL 2 backend](../wsl/index.md).

- **Add `*.docker.internal` to the host's `/etc/hosts` file. (Password required)**. Lets you resolve `*.docker.internal` DNS names from both the host and your containers.

- **Send usage statistics**. Select so Docker Desktop sends diagnostics,
  crash reports, and usage data. This information helps Docker improve and
  troubleshoot the application. Clear the check box to opt out. Docker may
  periodically prompt you for more information.

- **Show weekly tips**. Select to display useful advice and suggestions about
  using Docker.

- **Open Docker Desktop dashboard at startup**. Select to automatically open the
  dashboard when starting Docker Desktop.

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
Different settings are available for configuration depending on whether you are
using Linux containers in WSL 2 mode, Linux containers in Hyper-V mode, or Windows
containers.

### Advanced

> **Note**
>
> The **Resource allocation** options in the **Advanced** tab are only available in Hyper-V mode, because Windows manages
> the resources in WSL 2 mode and Windows container mode. In WSL 2
> mode, you can configure limits on the memory, CPU, and swap size allocated
> to the [WSL 2 utility VM](https://docs.microsoft.com/en-us/windows/wsl/wsl-config#configure-global-options-with-wslconfig).

On the **Advanced** tab, you can limit resources available to the Docker Linux VM.

Advanced settings are:

- **CPU limit**. Specify the maximum number of CPUs to be used by Docker Desktop.
  By default, Docker Desktop is set to use all the processors available on the host machine.

- **Memory limit**. By default, Docker Desktop is set to use up to `2` GB of your host's
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

> **Note**
>
> The **File sharing** tab is only available in Hyper-V mode because the files
> are automatically shared in WSL 2 mode and Windows container mode.

Use File sharing to allow local directories on your machine to be shared with
Linux containers. This is especially useful for editing source code in an IDE on
the host while running and testing the code in a container.

Note that configuring file sharing is not necessary for Windows containers,
only [Linux containers](../faqs/windowsfaqs.md#how-do-i-switch-between-windows-and-linux-containers).
If a directory is not shared with a Linux container you may get `file not found`
or `cannot start service` errors at runtime. See [Volume mounting requires shared folders for Linux containers](../troubleshoot/topics.md).

File share settings are:

- **Add a Directory**. select `+` and navigate to the directory you want to add.

- **Remove a Directory**. select `-` next to the directory you want to remove

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
> * Docker Desktop sets permissions to read/write/execute for users, groups and
>   others [0777 or a+rwx](https://chmodcommand.com/chmod-0777/).
>   This is not configurable. See [Permissions errors on data directories for shared volumes](../troubleshoot/topics.md).
> * Windows presents a case-insensitive view of the filesystem to applications while Linux is case-sensitive.
>   On Linux, it is possible to create two separate files: `test` and `Test`,
>   while on Windows these filenames would actually refer to the same underlying
>   file. This can lead to problems where an app works correctly on a developer's
>   machine (where the file contents are shared) but fails when run in Linux in
>   production (where the file contents are distinct). To avoid this, Docker Desktop
>   insists that all shared files are accessed as their original case. Therefore,
>   if a file is created called `test`, it must be opened as `test`. Attempts to
>   open `Test` will fail with the error "No such file or directory". Similarly,
>   once a file called `test` is created, attempts to create a second file called
>   `Test` will fail.
{ .tip }

#### Shared folders on demand

You can share a folder "on demand" the first time a particular folder is used by a container.

If you run a Docker command from a shell with a volume mount (as shown in the
example below) or kick off a Compose file that includes volume mounts, you get a
popup asking if you want to share the specified folder.

You can select to **Share it**, in which case it is added to your Docker Desktop Shared Folders list and available to
containers. Alternatively, you can opt not to share it by selecting **Cancel**.

![Shared folder on demand](../images/shared-folder-on-demand.png)

### Proxies

HTTP/HTTPS proxies can be used when:

- Signing in to Docker
- Pulling or pushing images
- Fetching artifacts during image builds
- Containers interact with the external network
- Scanning images

If the host uses a HTTP/HTTPS proxy configuration (static or via Proxy Auto-Configuration), Docker Desktop reads
this configuration
and automatically uses these settings for signing into Docker, for pulling and pushing images, and for
container Internet access. If the proxy requires authorization then Docker Desktop dynamically asks
the developer for a username and password. All passwords are stored securely in the OS credential store.
Note that only the `Basic` proxy authentication method is supported so we recommend using an `https://`
URL for your HTTP/HTTPS proxies to protect passwords while in transit on the network. Docker Desktop
supports TLS 1.3 when communicating with proxies.

To set a different proxy for Docker Desktop, turn on **Manual proxy configuration** and enter a single
upstream proxy URL of the form `http://proxy:port` or `https://proxy:port`.

To prevent developers from accidentally changing the proxy settings, see
[Settings Management](../hardened-desktop/settings-management/index.md#what-features-can-i-configure-with-settings-management).

The HTTPS proxy settings used for scanning images are set using the `HTTPS_PROXY` environment variable.

If you are running Windows containers in Docker, you can allow the Windows Docker daemon to use Docker Desktop's internal proxy, with the **Use proxy for Windows Docker daemon** setting.
This is useful when a corporate proxy that requires authentication is manually configured or set at the system level. If you are an admin for your organization and have a Docker Business subscription, you can control this setting with [Settings management](../hardened-desktop/settings-management/configure.md) using the `windowsDockerdPort` parameter.

### Network

> **Note**
>
> The **Network** tab isn't available in the Windows container mode because
> Windows manages networking.

{{< include "desktop-network-setting.md" >}}

### WSL Integration

In WSL 2 mode, you can configure which WSL 2 distributions will have the Docker
WSL integration.

By default, the integration is enabled on your default WSL distribution.
To change your default WSL distro, run `wsl --set-default <distro name>`. (For example,
to set Ubuntu as your default WSL distro, run `wsl --set-default ubuntu`).

You can also select any additional distributions you would like to enable the WSL 2 integration on.

For more details on configuring Docker Desktop to use WSL 2, see
[Docker Desktop WSL 2 backend](../wsl/index.md).

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

To see the full list of possible configuration options, see the
[dockerd command reference](/engine/reference/commandline/dockerd/).

Select **Apply & Restart** to save your settings and restart Docker Desktop.

## Builders

{{< include "desktop-builders-setting.md" >}}

## Kubernetes

> **Note**
>
> The **Kubernetes** tab is not available in Windows container mode.

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
When there's a new update, you can choose to download the update right away, or
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

## Features in development

Use the **Extensions** tab to:

- **Enable Docker Extensions**
- **Allow only extensions distributed through the Docker Marketplace**
- **Show Docker Extensions system containers**

For more information about Docker extensions, see [Extensions](../extensions/index.md).

## Feature control

On the **Feature control** tab you can control your settings for **Beta features** and **Experimental features**.

You can also sign up to the [Developer Preview Program](https://www.docker.com/community/get-involved/developer-preview/) from the **Features in development** tab.

### Beta features

{{< include "beta.md" >}}

#### Use containerd for pulling and storing images

Turns on the containerd image store. This brings new features like faster container startup performance by lazy-pulling images, and the ability to run Wasm applications with Docker. For more information, see [containerd image store](../containerd/index.md).

### Experimental features

{{< include "desktop-experimental.md" >}}

## Notifications

{{< include "desktop-notifications-settings.md" >}}