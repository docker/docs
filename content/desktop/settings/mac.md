---
description: Change your Docker Desktop settings on Mac
keywords: settings, preferences, proxy, file sharing, resources, kubernetes, Docker
  Desktop, Mac
title: Change Docker Desktop settings on Mac
aliases:
- /docker-for-mac/mutagen-caching/
- /docker-for-mac/mutagen/
- /docker-for-mac/osxfs-caching/
- /docker-for-mac/osxfs/
---

This page provides information on how to configure and manage your Docker Desktop for Mac settings.

To navigate to **Settings** either:

- Select the Docker menu {{< inline-image src="../images/whale-x.svg" alt="whale menu" >}} and then **Settings**
- Select the **Settings** icon from the Docker Dashboard.

You can also locate the `settings.json` file at `~/Library/Group Containers/group.com.docker/settings.json`.

## General

On the **General** tab, you can configure when to start Docker and specify other settings:

- **Start Docker Desktop when you log in**. Select to automatically start Docker
  Desktop when you sign in to your machine.

- **Choose Theme for Docker Desktop**. Choose whether you want to apply a **Light** or **Dark** theme to Docker Desktop. Alternatively you can set Docker Desktop to **Use System Settings**.

- **Choose container terminal**. Determines which terminal is launched when opening the terminal from a container.
If you choose the integrated terminal, you can run commands in a running container straight from the Docker Dashboard. For more information, see [Explore containers](../use-desktop/container.md).

- **Include VM in Time Machine backups**. Select to back up the Docker Desktop
  virtual machine. This option is turned off by default.

- **Use Virtualization framework**. Select to allow Docker Desktop to use the `virtualization.framework` instead of the `hypervisor.framework`.
    > **Tip**
    >
    > Turn this setting on to make Docker Desktop run faster.
    { .tip }

- **Choose file sharing implementation for your containers**. Choose whether you want to share files using **VirtioFS**, **gRPC FUSE**, or **osxfs**. VirtioFS is only available for macOS versions 12.5 and above, and is turned on by default.
    >**Tip**
    >
    > Use VirtioFS for speedy file sharing. VirtioFS has reduced the time taken to complete filesystem operations by [up to 98%](https://github.com/docker/roadmap/issues/7#issuecomment-1044452206)
    { .tip }


- **Use Rosetta for x86/AMD64 emulation on Apple Silicon**. Turns on Rosetta to accelerate x86/AMD64 binary emulation on Apple Silicon. This option is only available if you have turned on **Virtualization framework** in the **General** settings tab. You must also be on macOS Ventura or later. 

- **Send usage statistics**. Select so Docker Desktop sends diagnostics,
  crash reports, and usage data. This information helps Docker improve and
  troubleshoot the application. Clear the checkbox to opt out. Docker may
  periodically prompt you for more information.

- **Show weekly tips**. Select to display useful advice and suggestions about
  using Docker.

- **Open Docker Desktop dashboard at startup**. Select to automatically open the
  dashboard when starting Docker Desktop.

- **Use Enhanced Container Isolation**. Select to enhance security by preventing containers from breaching the Linux VM. For more information, see [Enhanced Container Isolation](../hardened-desktop/enhanced-container-isolation/index.md).
    >**Note**
    >
    > This setting is only available if you are signed in to Docker Desktop and have a Docker Business subscription.

- **Show CLI hints**. Displays CLI hints and tips when running Docker commands in the CLI. This is turned on by default. To turn CLI hints on or off from the CLI, set `DOCKER_CLI_HINTS` to `true` or `false` respectively.

## Resources

The **Resources** tab allows you to configure CPU, memory, disk, proxies,
network, and other resources.

### Advanced

On the **Advanced** tab, you can limit resources available to the Docker Linux VM.

Advanced settings are:

- **CPU limit**. Specify the maximum number of CPUs to be used by Docker Desktop.
  By default, Docker Desktop is set to use all the processors available on the host machine.

- **Memory limit**. By default, Docker Desktop is set to use up to 50% of your host's
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

By default the `/Users`, `/Volume`, `/private`, `/tmp` and `/var/folders` directory are shared.
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
{ .tip }

### Proxies

HTTP/HTTPS proxies can be used when:

- Signing in to Docker
- Pulling or pushing images
- Fetching artifacts during image builds
- Containers interact with the external network
- Scanning images

If the host uses a HTTP/HTTPS proxy configuration (static or via Proxy Auto-Configuration), Docker Desktop reads
this configuration
and automatically uses these settings for signing in to Docker, for pulling and pushing images, and for
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

#### Use containerd for pulling and storing images

Turns on the containerd image store. This brings new features like faster container startup performance by lazy-pulling images, and the ability to run Wasm applications with Docker. For more information, see [containerd image store](../containerd/index.md).

### Experimental features

{{< include "desktop-experimental.md" >}}

## Advanced

On the **Advanced** tab, you can reconfigure your initial installation settings:

- **Choose how to configure the installation of Docker's CLI tools**.
  - **System**: Docker CLI tools are installed in the system directory under `/usr/local/bin`
  - **User**: Docker CLI tools are installed in the user directory under `$HOME/.docker/bin`. You must then add `$HOME/.docker/bin` to your PATH. To add `$HOME/.docker/bin` to your path:
      1. Open your shell configuration file. This is `~/.bashrc` if you're using a bash shell, or `~/.zshrc` if you're using a zsh shell.
      2. Copy and paste the following:
            ```console
            $ export PATH=$PATH:~/.docker/bin
            ```
     3. Save and the close the file. Restart your shell to apply the changes to the PATH variable.

- **Enable default Docker socket (Requires password)**. Creates `/var/run/docker.sock` which some third party clients may use to communicate with Docker Desktop. For more information, see [permission requirements for macOS](../mac/permission-requirements.md#installing-symlinks).

- **Enable privileged port mapping (Requires password)**. Starts the privileged helper process which binds the ports that are between 1 and 1024. For more information, see [permission requirements for macOS](../mac/permission-requirements.md#binding-privileged-ports).

  For more information on each configuration and use case, see [Permission requirements](../mac/permission-requirements.md).

- **Automatically check configuration**. Regularly checks your configuration to ensure no unexpected changes have been made by another application.
