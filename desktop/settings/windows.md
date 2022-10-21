---
description: Docker Desktop settings
keywords: settings, preferences, proxy, file sharing, resources, kubernetes, Docker Desktop, Windows
title: Change Docker Desktop settings on Windows
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

- **Use integrated container terminal**. Select to execute commands in a running container straight from the Docker Dashboard. For more information, see [Explore containers](../use-desktop/container.md).

- **Expose daemon on tcp://localhost:2375 without TLS**. Check this option to
  enable legacy clients to connect to the Docker daemon. You must use this option
  with caution as exposing the daemon without TLS can result in remote code
  execution attacks.

- **Use the WSL 2 based engine**. WSL 2 provides better performance than the
  legacy Hyper-V backend. For more information, see [Docker Desktop WSL 2 backend](../windows/wsl.md).

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
Different settings are available for configuration depending on whether you are
using Linux containers in WSL 2 mode, Linux containers in Hyper-V mode, or Windows
containers.

### Advanced

> **Note**
>
> The Advanced tab is only available in Hyper-V mode, because Windows manages
> the resources in WSL 2 mode and Windows container mode. In WSL 2
> mode, you can configure limits on the memory, CPU, and swap size allocated
> to the [WSL 2 utility VM](https://docs.microsoft.com/en-us/windows/wsl/wsl-config#configure-global-options-with-wslconfig){:target="_blank"
class="_"}.

On the **Advanced** tab, you can limit resources available to Docker.

Advanced settings are:

- **CPUs**. By default, Docker Desktop is set to use half the number of processors
  available on the host machine. To increase processing power, set this to a
  higher number; to decrease, lower the number.

- **Memory**. By default, Docker Desktop is set to use `2` GB  of your host's
  memory. To increase the RAM, set this to a higher number; to decrease it,
  lower the number.

- **Swap**. Configure swap file size as needed. The default is 1 GB.

- **Disk image size**. Specify the size of the disk image.

- **Disk image location**. Specify the location of the Linux volume where containers and images are stored.

You can also move the disk image to a different location. If you attempt to move a disk image to a location that already has one, you are asked if you want to use the existing image or replace it.

### File sharing

> **Note**
>
> The File sharing tab is only available in Hyper-V mode because the files
> are automatically shared in WSL 2 mode and Windows container mode.

Use File sharing to allow local directories on your machine to be shared with
Linux containers. This is especially useful for editing source code in an IDE on
the host while running and testing the code in a container.

Note that configuring file sharing is not necessary for Windows containers,
only [Linux containers](../faqs/windowsfaqs.md#how-do-i-switch-between-windows-and-linux-containers).
If a directory is not shared with a Linux container you may get `file not found`
or `cannot start service` errors at runtime. See [Volume mounting requires shared folders for Linux containers](../troubleshoot/topics.md).

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

#### Shared folders on demand

You can share a folder "on demand" the first time a particular folder is used by a container.

If you run a Docker command from a shell with a volume mount (as shown in the
example below) or kick off a Compose file that includes volume mounts, you get a
popup asking if you want to share the specified folder.

You can select to **Share it**, in which case it is added to your Docker Desktop Shared Folders list and available to
containers. Alternatively, you can opt not to share it by selecting **Cancel**.

![Shared folder on demand](../images/shared-folder-on-demand.png){:width="600px"}

### Proxies

HTTP/HTTPS proxies can be used when:

- Logging in to Docker
- Pulling or pushing images
- Fetching artifacts during image builds
- Containers interact with the external network
- Scanning images

Each use case above is configured slightly differently.

If the host uses a static HTTP/HTTPS proxy configuration, Docker Desktop reads this configuration
and automatically uses these settings for logging into Docker and for pulling and pushing images.

If the host uses a more sophisticated HTTP/HTTPS configuration, enable **Manual proxy configuration** and enter a single upstream proxy URL
of the form `http://username:password@proxy:port`.

The HTTP/HTTPS proxy settings used for fetching artifacts during builds and for running containers
are set via the `.docker/config.json` file, see [Configure the Docker client](../../network/proxy.md#configure-the-docker-client).

The HTTPS proxy settings used for scanning images are set using the `HTTPS_PROXY` environment variable.

### Network

> **Note**
>
> The Network tab is not available in the Windows container mode because
> Windows manages networking.

You can configure Docker Desktop networking to work on a virtual private network (VPN). Specify a network address translation (NAT) prefix and subnet mask to enable Internet connectivity.

### WSL Integration

In WSL 2 mode, you can configure which WSL 2 distributions will have the Docker
WSL integration.

By default, the integration is enabled on your default WSL distribution.
To change your default WSL distro, run `wsl --set-default <distro name>`. (For example,
to set Ubuntu as your default WSL distro, run `wsl --set-default ubuntu`).

You can also select any additional distributions you would like to enable the WSL 2 integration on.

For more details on configuring Docker Desktop to use WSL 2, see
[Docker Desktop WSL 2 backend](../windows/wsl.md).

## Docker Engine

The **Docker Engine** tab allows you to configure the Docker daemon to determine how your containers run.

Type a JSON configuration file in the box to configure the daemon settings. For a full list of options, see the Docker Engine
[dockerd commandline reference](/engine/reference/commandline/dockerd/){:target="_blank" rel="noopener" class="_"}.

Click **Apply & Restart** to save your settings and restart Docker Desktop.

## Beta features

{% include beta.md %}

From the **Beta features** tab, you can sign up to the [Developer Preview program](https://www.docker.com/community/get-involved/developer-preview/).

On the **Beta features** tab, you also have the option to allow version 4.13 feature flags, which are product features Docker is currently experimenting with. This is switched on by default. 

## Kubernetes

> **Note**
>
> The **Kubernetes** tab is not available in Windows container mode.

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
When there's a new update, you can choose to download the update right away, or
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
