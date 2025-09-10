---
description: How to uninstall Docker Desktop
keywords: Windows, uninstall, Mac, Linux, Docker Desktop
title: Uninstall Docker Desktop
linkTitle: Uninstall
weight: 210
---

> [!WARNING]
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker-related data local to the machine, and removes the files generated
> by the application. To learn how to preserve important data before uninstalling, refer to the [back up and restore data](/manuals/desktop/settings-and-maintenance/backup-and-restore.md) section.

{{< tabs >}}
{{< tab name="Windows" >}}

#### From the GUI

1. From the Windows **Start** menu, select **Settings** > **Apps** > **Apps & features**.
2. Select **Docker Desktop** from the **Apps & features** list and then select **Uninstall**.
3. Select **Uninstall** to confirm your selection.

#### From the CLI

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

After uninstalling Docker Desktop, some residual files may remain which you can remove manually. These are:

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

#### From the GUI

1. Open Docker Desktop. 
2. In the top-right corner of the Docker Desktop Dashboard, select the **Troubleshoot** icon.
3. Select **Uninstall**.
4. When prompted, confirm by selecting **Uninstall** again.

You can then move the Docker application to the trash. 

#### From the CLI

Run:

```console
$ /Applications/Docker.app/Contents/MacOS/uninstall
```

You can then move the Docker application to the trash. 

> [!NOTE]
> You may encounter the following error when uninstalling Docker Desktop using the uninstall command.
>
> ```console
> $ /Applications/Docker.app/Contents/MacOS/uninstall
> Password:
> Uninstalling Docker Desktop...
> Error: unlinkat /Users/<USER_HOME>/Library/Containers/com.docker.docker/.com.apple.containermanagerd.metadata.plist: > operation not permitted
> ```
>
> The operation not permitted error is reported either on the file `.com.apple.containermanagerd.metadata.plist` or on the parent directory `/Users/<USER_HOME>/Library/Containers/com.docker.docker/`. This error can be ignored as you have successfully uninstalled Docker Desktop.
> You can remove the directory `/Users/<USER_HOME>/Library/Containers/com.docker.docker/` later by allowing **Full Disk Access** to the terminal application you are using (**System Settings** > **Privacy & Security** > **Full Disk Access**).

After uninstalling Docker Desktop, some residual files may remain which you can remove:

```console
$ rm -rf ~/Library/Group\ Containers/group.com.docker
$ rm -rf ~/.docker
```

With Docker Desktop version 4.36 and earlier, the following files may also be left on the file system. You can remove these with administrative privileges:

```console
/Library/PrivilegedHelperTools/com.docker.vmnetd
/Library/PrivilegedHelperTools/com.docker.socket
```

{{< /tab >}}
{{< tab name="Ubuntu" >}}

To uninstall Docker Desktop for Ubuntu:

1. Remove the Docker Desktop application. Run:

   ```console
   $ sudo apt remove docker-desktop
   ```

   This removes the Docker Desktop package itself but doesn’t delete all of its files or settings.

2. Manually remove leftover file.

   ```console
   $ rm -r $HOME/.docker/desktop
   $ sudo rm /usr/local/bin/com.docker.cli
   $ sudo apt purge docker-desktop
   ```

   This removes configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purges the remaining systemd service files.

3. Clean up Docker config settings. In `$HOME/.docker/config.json`, remove the `credsStore` and `currentContext` properties.

   These entries tell Docker where to store credentials and which context is active. If they remain after uninstalling Docker Desktop, they may conflict with a future Docker setup.

{{< /tab >}}
{{< tab name="Debian" >}}

To uninstall Docker Desktop for Debian, run:

1. Remove the Docker Desktop application:

   ```console
   $ sudo apt remove docker-desktop
   ```

   This removes the Docker Desktop package itself but doesn’t delete all of its files or settings.

2. Manually remove leftover file.

   ```console
   $ rm -r $HOME/.docker/desktop
   $ sudo rm /usr/local/bin/com.docker.cli
   $ sudo apt purge docker-desktop
   ```

   This removes configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purges the remaining systemd service files.

3. Clean up Docker config settings. In `$HOME/.docker/config.json`, remove the `credsStore` and `currentContext` properties.

   These entries tell Docker where to store credentials and which context is active. If they remain after uninstalling Docker Desktop, they may conflict with a future Docker setup.

{{< /tab >}}
{{< tab name="Fedora" >}}

To uninstall Docker Desktop for Fedora:

1. Remove the Docker Desktop application. Run:

   ```console
   $ sudo dnf remove docker-desktop
   ```

   This removes the Docker Desktop package itself but doesn’t delete all of its files or settings.

2. Manually remove leftover file.

   ```console
   $ rm -r $HOME/.docker/desktop
   $ sudo rm /usr/local/bin/com.docker.cli
   $ sudo dnf remove docker-desktop
   ```

   This removes configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purges the remaining systemd service files.

3. Clean up Docker config settings. In `$HOME/.docker/config.json`, remove the `credsStore` and `currentContext` properties.

   These entries tell Docker where to store credentials and which context is active. If they remain after uninstalling Docker Desktop, they may conflict with a future Docker setup.

{{< /tab >}}
{{< tab name="Arch" >}}

To uninstall Docker Desktop for Arch:

1. Remove the Docker Desktop application. Run:

   ```console
   $ sudo pacman -Rns docker-desktop
   ```

   This removes the Docker Desktop package along with its configuration files and dependencies not required by other packages.

2. Manually remove leftover files.

   ```console
   $ rm -r $HOME/.docker/desktop
   ```

   This removes configuration and data files at `$HOME/.docker/desktop`.

3. Clean up Docker config settings. In `$HOME/.docker/config.json`, remove the `credsStore` and `currentContext` properties.

   These entries tell Docker where to store credentials and which context is active. If they remain after uninstalling Docker Desktop, they may conflict with a future Docker setup.

{{< /tab >}}
{{< /tabs >}}


