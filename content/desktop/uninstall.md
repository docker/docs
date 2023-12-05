---
description: How to uninstall Docker Desktop
keywords: Windows, unintall, Mac, Linux, Docker Desktop
title: Uninstall Docker Desktop
---

{{< tabs >}}
{{< tab name="Windows" >}}

To uninstall Docker Desktop from your Windows machine:

1. From the Windows **Start** menu, select **Settings** > **Apps** > **Apps & features**.
2. Select **Docker Desktop** from the **Apps & features** list and then select **Uninstall**.
3. Select **Uninstall** to confirm your selection.

You can also uninstall Docker Desktop from the CLI:

1. Locate the installer:
   ```console
   $ C:\Program Files\Docker\Docker\Docker Desktop Installer.exe
   ```
2. Uninstall Docker Desktop. 
 - In PowerShell, run:
    ```console
    $ Start-Process 'Docker Desktop Installer.exe' -Wait uninstall
    ```
 - In the Command Prompt, run:
    ```console
    $ start /w "" "Docker Desktop Installer.exe" uninstall
    ```

After uninstalling Docker Desktop, there may be some residual files left behind which you can remove manually. These are:

```console
C:\ProgramData\Docker
C:\ProgramData\DockerDesktop
C:\Program Files\Docker
C:\Users\<your user name>\AppData\Local\Docker
C:\Users\<your user name>\AppData\Roaming\Docker
C:\Users\<your user name>\AppData\Roaming\Docker Desktop
C:\Users\<your user name>\.docker
```
 
{{< /tab >}}
{{< tab name="Mac" >}}

To uninstall Docker Desktop from your Mac:

1. From the Docker menu, select the **Troubleshoot** icon in the top-right corner of Docker Dashboard and then select **Uninstall**.
2. Select **Uninstall** to confirm your selection.

You can also uninstall Docker Desktop from the CLI. Run:

```console
$ /Applications/Docker.app/Contents/MacOS/uninstall
```

After uninstalling Docker Desktop, there may be some residual files left behind which you can remove:

```console
$ rm -rf ~/Library/Group\ Containers/group.com.docker
$ rm -rf ~/Library/Containers/com.docker.docker
$ rm -rf ~/.docker
```

{{< /tab >}}
{{< tab name="Linux" >}}

Docker Desktop is removed from a Linux host using the package manager.

Once Docker Desktop is removed, users must delete the `credsStore` and `currentContext` properties from the `~/.docker/config.json`.

{{< /tab >}}
{{< tab name="Ubuntu" >}}

To remove Docker Desktop for Ubuntu, run:

```console
$ sudo apt remove docker-desktop
```

For a complete cleanup, remove configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purge
the remaining systemd service files.

```console
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
$ sudo apt purge docker-desktop
```

Remove the `credsStore` and `currentContext` properties from `$HOME/.docker/config.json`. Additionally, you must delete any edited configuration files manually. 

{{< /tab >}}
{{< tab name="Debian" >}}

To remove Docker Desktop for Debian, run:

```console
$ sudo apt remove docker-desktop
```

For a complete cleanup, remove configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purge
the remaining systemd service files.

```console
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
$ sudo apt purge docker-desktop
```

Remove the `credsStore` and `currentContext` properties from `$HOME/.docker/config.json`. Additionally, you must delete any edited configuration files manually.  preserve important data before uninstalling.

{{< /tab >}}
{{< tab name="Fedora" >}}

To remove Docker Desktop for Fedora, run:

```console
$ sudo dnf remove docker-desktop
```

For a complete cleanup, remove configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purge
the remaining systemd service files.

```console
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
```

Remove the `credsStore` and `currentContext` properties from `$HOME/.docker/config.json`. Additionally, you must delete any edited configuration files manually. 

{{< /tab >}}
{{< tab name="Arch" >}}

To remove Docker Desktop for Arch, run:

```console
$ sudo pacman -R docker-desktop
```

For a complete cleanup, remove configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purge
the remaining systemd service files.

```console
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
$ sudo pacman -Rns docker-desktop
```

Remove the `credsStore` and `currentContext` properties from `$HOME/.docker/config.json`. Additionally, you must delete any edited configuration files manually. 

{{< /tab >}}
{{< /tabs >}}

> **Important**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker-related data local to the machine, and removes the files generated
> by the application. To learn how to preserve important data before uninstalling, refer to the [back up and restore data](backup-and-restore.md) section .
{ .important }
