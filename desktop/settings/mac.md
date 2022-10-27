---
description: Docker Desktop settings
keywords: settings, preferences, proxy, file sharing, resources, kubernetes, Docker Desktop, Mac
title: Change Docker Desktop preferences on Mac
redirect_from:
- /docker-for-mac/mutagen-caching/
- /docker-for-mac/mutagen/
- /docker-for-mac/osxfs-caching/
- /docker-for-mac/osxfs/
---

This page provides information on how to configure and manage your Docker Desktop settings.

To navigate to **Preferences** either:

- Select the Docker menu ![whale menu](../images/whale-x.svg){: .inline} and then **Preferences**
- Select the **Preferences** icon from the Docker Dashboard.

## General

On the **General** tab, you can configure when to start Docker and specify other settings:

- **Start Docker Desktop when you log in**. Select to automatically start Docker
  Desktop when you log into your machine.

- **Choose Theme for Docker Desktop**. Choose whether you want to apply a **Light** or **Dark** theme to Docker Desktop. Alternatively you can set Docker Desktop to **Use System Settings**.

- **Use integrated container terminal**. Select to execute commands in a running container straight from the Docker Dashboard. For more information, see [Explore containers](../use-desktop/container.md).

- **Include VM in Time Machine backups**. Select to back up the Docker Desktop
  virtual machine. This option is disabled by default.

- **Use gRPC FUSE for file sharing**. Clear this check box to use the legacy
  osxfs file sharing instead.

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

- **Memory**. By default, Docker Desktop is set to use `2` GB  of your host's
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

By default the `/Users`, `/Volume`, `/private`, `/tmp` and `/var/folders` directory are shared.
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
> * If you share the whole of your home directory into a container, MacOS may
>   prompt you to give Docker access to personal areas of your home directory such as
>   your Reminders or Downloads.
> * By default, Mac file systems are case-insensitive while Linux is case-sensitive.
>   On Linux, it is possible to create two separate files: `test` and `Test`,
>   while on Mac these filenames would actually refer to the same underlying
>   file. This can lead to problems where an app works correctly on a developer's
>   machine (where the file contents are shared) but fails when run in Linux in
>   production (where the file contents are distinct). To avoid this, Docker Desktop
>   insists that all shared files are accessed as their original case. Therefore,
>   if a file is created called `test`, it must be opened as `test`. Attempts to
>   open `Test` will fail with the error "No such file or directory". Similarly,
>   once a file called `test` is created, attempts to create a second file called
>   `Test` will fail.
>
> For more information, see [Volume mounting requires file sharing for any project directories outside of `/Users`](../troubleshoot/topics.md)

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

HTTP/HTTPS traffic from image builds and running containers is forwarded transparently to the same
upstream proxy used for logging in and image pulls.
If you want to override this behaviour and use different HTTP/HTTPS proxies for image builds and
running containers, see [Configure the Docker client](../../network/proxy.md#configure-the-docker-client).

The HTTPS proxy settings used for scanning images are set using the `HTTPS_PROXY` environment variable.

### Network

You can configure Docker Desktop networking to work on a virtual private network (VPN). Specify a network address translation (NAT) prefix and subnet mask to enable Internet connectivity.

## Docker Engine

The **Docker Engine** tab allows you to configure the Docker daemon to determine how your containers run.

Type a JSON configuration file in the box to configure the daemon settings. For a full list of options, see the Docker Engine
[dockerd commandline reference](/engine/reference/commandline/dockerd/){:target="_blank" rel="noopener" class="_"}.

Click **Apply & Restart** to save your settings and restart Docker Desktop.

## Beta Features

{% include beta.md %}

On the **Beta features** tab, you also have the option to allow version 4.13 feature flags, which are product features Docker is currently experimenting with. This is switched on by default. 

### Enable the new Apple Virtualization framework

Select **Use the new Virtualization framework** to allow Docker Desktop to use the new `virtualization.framework` instead of the ‘hypervisor.framework’. Ensure to reset your Kubernetes cluster when you enable the new Virtualization framework for the first time.

### Enable VirtioFS

 Docker Desktop for Mac lets developers use a new experimental file-sharing implementation called [virtiofS](https://virtio-fs.gitlab.io/){: target='_blank' rel='noopener' class='_'}; the current default is gRPC-FUSE. virtiofs has been found to significantly improve file sharing performance on macOS. For more details, see our blog post [Speed boost achievement unlocked on Docker Desktop 4.6 for Mac](https://www.docker.com/blog/speed-boost-achievement-unlocked-on-docker-desktop-4-6-for-mac/){:target="_blank" rel="noopener" class="_"}.

To enable virtioFS:

1. Verify that you are on the following macOS version:
   - macOS 12.2 or later (for Apple Silicon)
   - macOS 12.3 or later (for Intel)

2. Select **Enable VirtioFS accelerated directory sharing** to enable virtioFS.

3. Click **Apply & Restart**.

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
