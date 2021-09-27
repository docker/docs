---
description: Docker Desktop for Windows and Docker Toolbox
keywords: windows, alpha, beta, toolbox, docker-machine, tutorial
sitemap: false
title: Docker Toolbox
redirect_from:
- /docker-for-mac/docker-toolbox/
- /docker-for-windows/docker-toolbox/
- /mackit/docker-toolbox/
- /toolbox/faqs/
- /toolbox/faqs/troubleshoot/
- /toolbox/overview/
- /toolbox/toolbox_install_mac/
- /toolbox/toolbox_install_windows/
toc_min: 1
toc_max: 2
---

>**Deprecated**
>
> Docker Toolbox has been deprecated and is no longer in active development. Please
> use Docker Desktop instead. See [Docker Desktop for Mac](../desktop/mac/index.md)
> and [Docker Desktop for Windows](../desktop/windows/index.md).
{: .warning }

This page explains how to migrate your Docker Toolbox installation to Docker Desktop.
It also contains instructions on how to uninstall Docker Toolbox from Mac and Windows machines.

## Migrate from Docker Toolbox to Docker Desktop

Uninstalling Docker Toolbox will remove your local image cache, volumes, containers,
and other data stored in Docker. Refer to the [back up and restore data](../desktop/backup-and-restore.md)
documentation before uninstalling Docker Toolbox, to learn how to back up your
data.

## Uninstall Docker Toolbox

Removing Toolbox involves removing all the Docker components it includes.

A full uninstall also includes removing the local and remote machines
you created with Docker Machine. If you have remote machines on a cloud provider and
you plan to manage them using the provider, you wouldn't want to remove
them. So, the step to remove machines is described here as optional.

### Uninstall Docker Toolbox on Mac

To uninstall Docker Toolbox on Mac:

1.  List your machines.

    ```console
    $ docker-machine ls
    NAME                ACTIVE   DRIVER       STATE     URL                        SWARM
    dev                 *        virtualbox   Running   tcp://192.168.99.100:2376
    my-docker-machine            virtualbox   Stopped
    default                      virtualbox   Stopped
    ```

2.  Optionally, remove each machine. For example:

    ```console
    $ docker-machine rm my-docker-machine
    Successfully removed my-docker-machine
    ```

3.  In your "Applications" folder, remove the "Docker" directory,
    which contains "Docker Quickstart Terminal" and "Kitematic".

4.  Run the following in a command shell to fully remove Kitematic:

    ```console
    $ rm -fr ~/Library/Application\ Support/Kitematic
    ```

5.  Remove the `docker`, `docker-compose`, and `docker-machine` commands from
    the `/usr/local/bin` folder.  Docker Desktop for Mac and Brew may also have
    installed them; in case of doubt leave them, or reinstall them via Brew, or
    rerun Docker Desktop for Mac (no need to reinstall it).

    ```console
    $ rm -f /usr/local/bin/docker
    $ rm -f /usr/local/bin/docker-compose
    $ rm -f /usr/local/bin/docker-machine
    ```

6.  Optionally, remove the `~/.docker/machine` directory.

    This directory stores some configuration and/or state, such as information
    about created machines and certificates.

7.  Uninstall Oracle VirtualBox, which is installed as a part of the
    Toolbox install.

### Uninstall Docker Toolbox on Windows

To uninstall Toolbox on Windows:

1.  List your machines.

    ```console
    $ docker-machine ls
    NAME                ACTIVE   DRIVER       STATE     URL                        SWARM
    dev                 *        virtualbox   Running   tcp://192.168.99.100:2376
    my-docker-machine            virtualbox   Stopped
    default                      virtualbox   Stopped
    ```

2.  Remove each machine. For example:

    ```console
    $ docker-machine rm my-docker-machine
    Successfully removed my-docker-machine
    ```

3. Uninstall Docker Toolbox using Window's standard process for uninstalling programs through the control panel (programs and features).

    >**Note**: This process does not remove the `docker-install.exe` file. You must delete that file yourself.

4. Optionally, remove the `C:\Users\<your-user>\.docker` directory.

    If you want to remove Docker entirely, you
    can verify that the uninstall removed
    the `.docker` directory under your user path.
    If it is still there, remove it manually.
    This directory stores some Docker
    program configuration and state, such as
    information about created machines and
    certificates. You usually don't need to remove this directory.

5. Uninstall Oracle VirtualBox, which is
  installed as a part of the Toolbox install.
