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

> **Important**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.
{ .important }

{{< /tab >}}
{{< tab name="Mac" >}}

To uninstall Docker Desktop from your Mac:

1. From the Docker menu, select the **Troubleshoot** icon in the top-right corner of Docker Dashboard and then select **Uninstall**.
2. Select **Uninstall** to confirm your selection.

> Uninstall Docker Desktop from the command line
>
> To uninstall Docker Desktop from a terminal, run: `<path to Docker app>/Contents/MacOS/uninstall`.
> If your instance is installed in the default location, this
> command provides a clean uninstall:
>
> ```console
> $ /Applications/Docker.app/Contents/MacOS/uninstall
> Uninstalling Docker Desktop...
> Docker uninstalled successfully. You can move the Docker application to the trash.
> ```
>
> You might want to use the command line to uninstall if, for example, you find that
> the app is non-functional, and you cannot uninstall it from the menu.

{{< /tab >}}
{{< tab name="Linux" >}}

Docker Desktop can be removed from a Linux host using the package manager.

Once Docker Desktop has been removed, users must remove the `credsStore` and `currentContext` properties from the `~/.docker/config.json`.

{{< /tab >}}
{{< tab name="Ubuntu" >}}

### Remove Docker Desktop

1. Run the following command to remove Docker Desktop:
```
sudo apt remove docker-desktop
```
2. Remove configuration and data files:
```
rm -r $HOME/.docker/desktop
```
3. Remove the symlink:
```
sudo rm /usr/local/bin/com.docker.cli
```
4. Purge remaining systemd service files:
```
sudo apt purge docker-desktop
```
5. Remove the credsStore and currentContext properties from `$HOME/.docker/config.json`:
```
sudo nano $HOME/.docker/config.json
```
Remove the following lines:
```
{
  "credsStore": "desktop",
  "currentContext": "desktop"
}
```
6. Save and exit the file.
7. Delete any edited configuration files manually.

That's it! You have successfully removed Docker Desktop for Ubuntu.

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

> **Important**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.
{ .important }

{{< /tab >}}
{{< /tabs >}}
