---
description: How to install Docker Compose
keywords: compose, orchestration, install, installation, docker, documentation
title: Install Docker Compose
toc_max: 2
---

This page contains information on how to install Docker Compose. You can run Compose on macOS, Windows, and 64-bit Linux.


## Prerequisites

* Docker Compose requires Docker Engine.
* Docker Compose plugin requires Docker CLI.

## Install Compose

Check what installation scenario fits your needs. Are you looking to:

* __Get latest Compose and pre-requirements of the bat__
Install Docker Desktop for your platform. This is the fastest route and you get Docker Engine, Docker CLI and Compose plugin. 
Docker Desktop is available for Mac, Windows and Linux.

* __Install Compose plug in:__
  + If you got __Docker Desktop__ you're all set, no need to install anything separatedly.
  + __Windows Server__: If you want to run the Docker daemon and client directly on Microsoft Windows Server, follow the [Windows Server install instructions](#install-compose-on-windows-server)
  + __Linux systems__: 
     + To install Docker Compose with the Docker Engine we recommend [setting up Docker's repositories](#install-using-the-repository) or,
     + Using the [convenience scripts](../../engine/install/#server){: target="_blank" rel="noopener" class="_"} offered per Linux distro from the Engine install section. 
    + Other scenarios, check the [Linux install](#installing-compose-on-linux-systems)

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#dockerdesktop">Docker Desktop Install</a></li>
<li><a data-toggle="tab" data-target="#linux">Linux Install</a></li>
<li><a data-toggle="tab" data-target="#windows-server">Windows Server Install</a></li>
<li><a data-toggle="tab" data-target="#uninstall">Uninstall</a></li>
</ul>

<div class="tab-content">
<div id="dockerdesktop" class="tab-pane fade in active" markdown="1">

### Installing Docker Desktop

With Docker Desktop you get Docker Engine, Docker CLI with Compose plugin as well as other components and tools. 
Check a list of what's shipped with Docker and a list of key features in the [Docker Desktop Overview](../desktop/index.md).

Docker Desktop is available for Mac, Windows and Linux.
For download information, system requirements, and installation instructions, see:

* [Docker Desktop for Linux](../../desktop/linux/install.md)
* [Docker Desktop for Mac](../../desktop/mac/install.md)
* [Docker Desktop for Windows](../../desktop/windows/install.md)

For information about Docker Desktop licensing, see [Docker Desktop License Agreement](../../subscription/index.md#docker-desktop-license-agreement).
</div>

<div id="linux" class="tab-pane fade" markdown="1">
{: id="tabtest"}

### Installing Compose on Linux systems

* Installing Docker Desktop for Linux is the easiest and recommended route. Check the [supported platforms](../desktop/linux/install.md/#supported-platforms) page to verify the supported Linux distributions and architectures.

You can also:

* __Use the automated convenience scripts__ (for testing and development environments). 
These scripts will get you Docker Engine and Docker CLI with the Compose plugin. 
Choose your Linux distribution from the [Docker Engine install](..engine/install/) and follow the instructions available there.

* __Set up Dockers repositories__ and install from them. This is option is second-best regarding ease of installation and upgrading. See the [Install using the repository](#install-using-the-repository) section in this page.

* __Use a package manager__. See [convenience scripts](#install-using-the-convenience-script) for instructions.

* By downloading and __installing the compose plugin binary manually__. Note that this option requires you to manage upgrades manually as well. See the [Install the Plugin manually](#install-the-binary-manually) section in this page.


#### Install using the repository


>These instructions assume you already have Docker Engine and Docker CLI installed and now want to install the Compose plugin. 
For other scenarios check [this summary](#install-compose).

>To run Compose as a non-root user, see [Manage Docker as a non-root user](../engine/install/linux-postinstall.md).


If you have already set up the Docker repository jump to step 2.

1. Go to the `Set up the repository` section of your corresponding [Linux distribution](../engine/install/index.md#server) found in the Docker Engine installation pages.

2. Update the `apt` package index, and install the _latest version_ of Docker Compose:

    ```console
    $ sudo apt-get update
    $ sudo apt-get install docker-compose-plugin
    ```
    
    Alternatively, to install a specific version of Docker Engine:
      
    a. List the versions available in your repo:


      ```console
      $ apt-cache madison docker-compose-plugin

        docker-compose-plugin | 2.3.3~ubuntu-focal | https://download.docker.com/linux/ubuntu focal/stable arm64 Packages
      ```

    b. From the list obtained use the version string you can in the second column to specify the version you wish to install. 
      
    c. Install the selected version:


      ```console
      $ sudo apt-get install docker-compose-plugin=<VERSION_STRING>
      ```
    where <VERSION_STRING> is, for example,`2.3.3~ubuntu-focal`.

3.  Verify that Docker Compose is installed correctly by checking the version.

    ```console
    $ docker compose version
    Docker Compose version v2.3.3
    ```

#### Install the Plugin manually

>
>These instructions assume you already have Docker Engine and Docker CLI installed and now want to install the Compose plugin. 
For other scenarios check [this summary](#install-compose).

>To run Compose as a non-root user, see [Manage Docker as a non-root user](../engine/install/linux-postinstall.md).


1.  To download and install the Compose plugin, run:

    ```console
    $ DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
    $ mkdir -p $DOCKER_CONFIG/cli-plugins
    $ curl -SL https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/cli-plugins/docker-compose
    ```

    This command downloads the current stable release of Docker Compose (from the Compose releases repository) and installs Compose for _the active user_ under `$HOME` directory. 
    
    > To install:
    >* Docker Compose for _all users_ on your system, replace `~/.docker/ cli-plugins` with `/usr/local/lib/docker/cli-plugins`.
    >* A different version of Compose, substitute `{{site.compose_version}}` with the version of Compose you want to use.

2. Apply executable permissions to the binary:

     ```console
    $ chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose
    ```
    or, if you chose to install Compose for all users:

    ```console
    $ sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
    ```

3. Test the installation.

    ```console
    $ docker compose version
    Docker Compose version {{site.compose_version}}
    ```

> __Compose standalone__: If you need to use Compose without installing the Docker CLI, the instructions for installing the standalone binary are similar. 

1. To download and install Compose Standalone, run:
  ```console
  $ curl -SL https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
  ```
2. Apply executable permissions to the standalone binary in target path for the installation.
3. Test and execute compose commands use, `docker-compose`.
  > **Note**:
  > If the command `docker-compose` fails after installation, check your path.
  > You can also create a symbolic link to `/usr/bin` or any other directory in your path. For example:
  > ```console
  > $ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
  > ```
</div>

<div id="windows-server" class="tab-pane fade" markdown="1">

### Install Compose on Windows Server

Follow these instructions if you are running the Docker daemon and client directly
on Microsoft Windows Server and want to install Docker Compose.


1.  Start an "elevated" PowerShell (run it as administrator).
    Search for PowerShell, right-click, and choose
    **Run as administrator**. When asked if you want to allow this app
    to make changes to your device, click **Yes**.
    
2.  In PowerShell, since GitHub now requires TLS1.2, run the following:
    
    ```powershell  
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    ```

    Then run the following command to download the current stable release of
    Compose ({{site.compose_version}}):

    ```powershell
    Invoke-WebRequest "https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-Windows-x86_64.exe" -UseBasicParsing -OutFile $Env:ProgramFiles\Docker\docker-compose.exe
    ```

    > **Note**
    >
    > On Windows Server 2019, you can add the Compose executable to `$Env:ProgramFiles\Docker`. Because this directory is  registered in the system `PATH`, you can run the `docker-compose --version` command on the subsequent step with no additional configuration.

    > To install a different version of Compose, substitute `{{site.compose_version}}`
    > with the version of Compose you want to use.

3.  Test the installation.

    ```console
    $ docker compose version
    Docker Compose version {{site.compose_version}}
    ```
</div>

<div id="uninstall" class="tab-pane fade" markdown="1">

### Uninstalling Docker Desktop

* See Uninstall Docker Desktop for:
  + [Mac](../../desktop/mac/install/#uninstall-docker-desktop)
  + [Windows](../../desktop/windows/install/#uninstall-docker-desktop)
  + [Linux](../../desktop/linux/install.md/#uninstall-docker-desktop)



* __Other scenarios__:
To uninstall Docker Compose if you installed using `curl`:

```console
$ rm $DOCKER_CONFIG/cli-plugins/docker-compose
```
    
or if you choose to install Compose for all users  

```console
$ sudo rm /usr/local/lib/docker/cli-plugins/docker-compose
```

> Got a "Permission denied" error?
>
> If you get a "Permission denied" error using either of the above
> methods, you probably do not have the proper permissions to remove
> `docker-compose`. To force the removal, prepend `sudo` to either of the above
> commands and run again.
</div>
</div>

----




## Where to go next

- [User guide](index.md)
- [Getting Started](gettingstarted.md)
- [Command line reference](reference/index.md)
- [Compose file reference](compose-file/index.md)
- [Sample apps with Compose](samples-for-compose.md)
