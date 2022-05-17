---
description: How to install Docker Compose
keywords: compose, orchestration, install, installation, docker, documentation
title: Install Docker Compose
toc_max: 2
---

> **Accelerating new features in Docker Desktop**
>
> Docker Desktop helps you build, share, and run containers easily on Mac and Windows as you do on Linux. Docker handles the complex setup and allows you to focus on writing the code. Thanks to the positive support we received on the [subscription updates](https://www.docker.com/blog/updating-product-subscriptions/){: target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}, we've started working on [Docker Desktop for Linux](https://www.docker.com/blog/accelerating-new-features-in-docker-desktop/){: target="_blank" rel="noopener" class="_" id="dkr_docs_cta"} which is the second-most popular feature request in our public roadmap. If you are interested in early access, sign up for our [Developer Preview program](https://www.docker.com/community/get-involved/developer-preview){: target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}.
{: .important}

This page contains information on how to install Docker Compose. You can run Compose on macOS, Windows, and 64-bit Linux.

## Prerequisites

Docker Compose relies on Docker Engine for any meaningful work, so make sure you
have Docker Engine installed either locally or remote, depending on your setup.

- On desktop systems like Docker Desktop for Mac and Windows, Docker Compose is
included as part of those desktop installs.

- On Linux systems, you can install Docker Compose with the Docker Engine using the 
[convenience script](../engine/install/index.md#server){: target="_blank" rel="noopener" class="_"}. Select the install Docker Engine page for your distribution and then look for instructions on installing using the convenience script.  
Otherwise, you should first install the [Docker Engine](../engine/install/index.md#server){: target="_blank" rel="noopener" class="_"}
for your OS and then refer to this page for
instructions on installing Compose on
Linux systems.

- To run Compose as a non-root user, see [Manage Docker as a non-root user](../engine/install/linux-postinstall.md).

## Install Compose

Follow the instructions below to install Compose on Mac, Windows, Windows Server, or Linux systems.

> Install a different version
>
> The instructions below outline installation of the current stable release
> (**{{site.compose_version}}**) of Compose. To install a different version of
> Compose, replace the given release number with the one that you want.
>
> Compose releases are also listed and available for direct download on the
> [Compose repository release page on GitHub](https://github.com/docker/compose/releases){:target="_blank" rel="noopener" class="_"}.
> 
> To install the Python version of Compose, follow instructions in the [Compose v1 GitHub branch](https://github.com/docker/compose/blob/master/INSTALL.md){: target="_blank" rel="noopener" class="_"}.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#macOS">Mac</a></li>
<li><a data-toggle="tab" data-target="#windows">Windows</a></li>
<li><a data-toggle="tab" data-target="#windows-server">Windows Server</a></li>
<li><a data-toggle="tab" data-target="#linux">Linux</a></li>
<li><a data-toggle="tab" data-target="#linux-standalone">Linux Standalone binary</a></li>
</ul>
<div class="tab-content">
<div id="macOS" class="tab-pane fade in active" markdown="1">

### Install Compose on macOS

**Docker Desktop for Mac** includes Compose along
with other Docker apps, so Mac users do not need to install Compose separately.
For installation instructions, see [Install Docker Desktop on Mac](../desktop/mac/install.md).

</div>
<div id="windows" class="tab-pane fade" markdown="1">

### Install Compose on Windows desktop systems

**Docker Desktop for Windows** includes Compose
along with other Docker apps, so most Windows users do not need to
install Compose separately. For install instructions, see [Install Docker Desktop on Windows](../desktop/windows/install.md).

If you are running the Docker daemon and client directly on Microsoft
Windows Server, follow the instructions in the Windows Server tab.

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
    $ docker-compose version
    Docker Compose version {{site.compose_version}}
    ```
</div>
<div id="linux" class="tab-pane fade" markdown="1">

### Install Compose on Linux systems

You can install Docker Compose in different ways, depending on your needs:

- In testing and development environments, some users choose to use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.
- Most users [set up Docker's repositories](#install-using-the-repository) and
  install from them, for ease of installation and upgrade tasks. This is the
  recommended approach.
- Some users [download and install the binary](#install-the-binary-manually),
  and manage upgrades manually.


#### Install using the convenience script

As Docker Compose is now part of the Docker CLI it can be installed via a convenience script with Docker Engine and the CLI.  
[Choose your Linux distribution](../engine/install/index.md#server) and follow the instructions.


#### Install using the repository

If you already follow the instructions to install Docker Engine, Docker Compose should already be installed.   
Otherwise, you can set up the Docker repository as mentioned in the Docker Engine installation, [choose your Linux distribution](../engine/install/index.md#server) and go to the `Set up the repository` section.

When finished

1. Update the `apt` package index, and install the _latest version_ of Docker Compose, or go to the next step to install a specific version:

    ```console
    $ sudo apt-get update
    $ sudo apt-get install docker-compose-plugin
    ```


2.  To install a _specific version_ of Docker Engine, list the available versions
    in the repo, then select and install:

    a. List the versions available in your repo:

    ```console
    $ apt-cache madison docker-compose-plugin

      docker-compose-plugin | 2.3.3~ubuntu-focal | https://download.docker.com/linux/ubuntu focal/stable arm64 Packages
    ```

    b. Install a specific version using the version string from the second column,
    for example, `2.3.3~ubuntu-focal`.

    ```console
    $ sudo apt-get install docker-compose-plugin=<VERSION_STRING>
    ```

3.  Verify that Docker Compose is installed correctly by checking the version.

    ```console
    $ docker-compose version
    Docker Compose version v2.3.3
    ```


#### Install the binary manually

On Linux, you can download the Docker Compose binary from the
[Compose repository release page on GitHub](https://github.com/docker/compose/releases){:target="_blank" rel="noopener" class="_"} and copying it into `$HOME/.docker/cli-plugins` as `docker-compose`.
Follow the instructions from the link, which involve running the `curl` command
in your terminal to download the binaries. These step-by-step instructions are
also included below.

1.  Run this command to download the current stable release of Docker Compose:

    ```console
    $ DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
    $ mkdir -p $DOCKER_CONFIG/cli-plugins
    $ curl -SL https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/cli-plugins/docker-compose
    ```

    This command installs Compose for the active user under `$HOME` directory. To install Docker Compose for all users on your system, replace `~/.docker/cli-plugins` with `/usr/local/lib/docker/cli-plugins`.


    > To install a different version of Compose, substitute `{{site.compose_version}}`
    > with the version of Compose you want to use.

2. Apply executable permissions to the binary:

     ```console
    $ chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose
    ```
    or if you choose to install Compose for all users
    ```console
    $ sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
    ```

3. Test the installation.

    ```console
    $ docker-compose version
    Docker Compose version {{site.compose_version}}
    ```
</div>
<div id="linux-standalone" class="tab-pane fade" markdown="1">


### Install Compose as standalone binary on Linux systems

You can use Compose as a standalone binary without installing the Docker CLI.

1. Run this command to download the current stable release of Docker Compose:

  ```console
  $ curl -SL https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
  ```

  > To install a different version of Compose, substitute `{{site.compose_version}}`
  > with the version of Compose you want to use.

2. Apply executable permissions to the binary:

  ```console
  $ sudo chmod +x /usr/local/bin/docker-compose
  ```
  
  > **Note**:
  >
  > If the command `docker-compose` fails after installation, check your path.
  > You can also create a symbolic link to `/usr/bin` or any other directory in your path.
  >
  > For example:
  > ```console
  > $ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
  > ```

3. Test the installation.

    ```console
    $ docker-compose --version
    Docker Compose version {{site.compose_version}}
    ```
</div>
</div>

----

## Uninstallation

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


## Where to go next

- [User guide](index.md)
- [Getting Started](gettingstarted.md)
- [Command line reference](reference/index.md)
- [Compose file reference](compose-file/index.md)
- [Sample apps with Compose](samples-for-compose.md)
