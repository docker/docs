---
description: How to uninstall Docker Desktop 
keywords: windows, install, download, run, docker, local
title: Uninstall Docker Desktop
---

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab3">Windows</a></li>
  <li><a data-toggle="tab" data-target="#tab4">Mac</a></li>
  <li><a data-toggle="tab" data-target="#tab5">Linux</a></li>
  <li><a data-toggle="tab" data-target="#tab6">Ubuntu</a></li>
  <li><a data-toggle="tab" data-target="#tab7">Debian</a></li>
  <li><a data-toggle="tab" data-target="#tab8">Fedora</a></li>
  <li><a data-toggle="tab" data-target="#tab9">Arch</a></li>
</ul>
<div class="tab-content">
<div id="tab3" class="tab-pane fade in active" markdown="1">
<br>
To uninstall Docker Desktop from your Windows machine:

1. From the Windows **Start** menu, select **Settings** > **Apps** > **Apps & features**.
2. Select **Docker Desktop** from the **Apps & features** list and then select **Uninstall**.
3. Click **Uninstall** to confirm your selection.

> **Important**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.
<hr>
</div>
<div id="tab4" class="tab-pane fade" markdown="1">
<br>
To uninstall Docker Desktop from your Mac:

1. From the Docker menu, select **Troubleshoot** and then select **Uninstall**.
2. Click **Uninstall** to confirm your selection.

> Uninstall Docker Desktop from the command line
>
> To uninstall Docker Desktop from a terminal, run: `<DockerforMacPath>
> --uninstall`. If your instance is installed in the default location, this
> command provides a clean uninstall:
>
> ```console
> $ /Applications/Docker.app/Contents/MacOS/Docker --uninstall
> Docker is running, exiting...
> Docker uninstalled successfully. You can move the Docker application to the trash.
> ```
>
> You might want to use the command-line uninstall if, for example, you find that
> the app is non-functional, and you cannot uninstall it from the menu.

> **Note**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.

<hr>
</div>
<div id="tab5" class="tab-pane fade in active" markdown="1">
<br>
Docker Desktop can be removed from a Linux host using the package manager.

Once Docker Desktop has been removed, users must remove the `credsStore` and `currentContext` properties from the `~/.docker/config.json`.

> **Note**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.
<hr>
</div>
<div id="tab6" class="tab-pane fade" markdown="1">
<br>

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
<hr>
</div>
<div id="tab7" class="tab-pane fade" markdown="1">
<br>

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

<hr>
</div>
<div id="tab8" class="tab-pane fade" markdown="1">
<br>
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

<hr>
</div>
<div id="tab9" class="tab-pane fade" markdown="1">
<br>

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

<hr>
</div>
</div>