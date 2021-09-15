---
description: How to install Docker Compose
keywords: compose, orchestration, install, installation, docker, documentation
title: Install Docker Compose
toc_max: 2
---

You can run Compose on macOS, Windows, and 64-bit Linux.

## Prerequisites

Docker Compose relies on Docker Engine for any meaningful work, so make sure you
have Docker Engine installed either locally or remote, depending on your setup.

- On desktop systems like Docker Desktop for Mac and Windows, Docker Compose is
included as part of those desktop installs.

- On Linux systems, first install the
[Docker Engine](../engine/install/index.md#server){: target="_blank" rel="noopener" class="_"}
for your OS as described on the Get Docker page, then come back here for
instructions on installing Compose on
Linux systems.

- To run Compose as a non-root user, see [Manage Docker as a non-root user](../engine/install/linux-postinstall.md).

## Install Compose

Follow the instructions below to install Compose on Mac, Windows, Windows Server
2016, or Linux systems.

> Install a different version
> 
> The instructions below outline installation of the current stable release
> (**v{{site.compose_version}}**) of Compose. To install a different version of
> Compose, replace the given release number with the one that you want. Compose
> releases are also listed and available for direct download on the
> [Compose repository release page on GitHub](https://github.com/docker/compose/releases){:target="_blank" rel="noopener" class="_"}.
> To install a **pre-release** of Compose, refer to the [install pre-release builds](#install-pre-release-builds)
> section.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#macOS">Mac</a></li>
<li><a data-toggle="tab" data-target="#windows">Windows</a></li>
<li><a data-toggle="tab" data-target="#windows-server">Windows Server</a></li>
<li><a data-toggle="tab" data-target="#linux">Linux</a></li>
<li><a data-toggle="tab" data-target="#pre-release">Pre-release builds</a></li>
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
    Compose (v{{site.compose_version}}):

    ```powershell
    Invoke-WebRequest "https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-Windows-x86_64.exe" -UseBasicParsing -OutFile $Env:ProgramFiles\Docker\docker-compose.exe
    ```

**Note**: On Windows Server 2019, you can add the Compose executable to `$Env:ProgramFiles\Docker`. Because this directory is  registered in the system `PATH`, you can run the `docker compose --version` command on the subsequent step with no additional configuration.

    > To install a different version of Compose, substitute `{{site.compose_version}}`
    > with the version of Compose you want to use.

3.  Test the installation.

    ```powershell
    docker compose --version

    docker compose version {{site.compose_version}}, build 01110ad01
    ```

</div>
<div id="linux" class="tab-pane fade" markdown="1">

### Install Compose on Linux systems

On supported Linux distributions, you can install Docker Compose using
the `docker-compose-plugin` system package. See [Engine installation instruction](/engine/install/)

Alternatively, you can download the Docker Compose binary from the
[Compose repository release page on GitHub](https://github.com/docker/compose/releases){:target="_blank" rel="noopener" class="_"}.
Follow the instructions from the link, which involve running the `curl` command
in your terminal to download the binaries. These step-by-step instructions are
also included below.

1.  Run this command to download the current stable release of Docker Compose:

    ```console
    $ mkdir -p ~/.docker/cli-plugins
    $ sudo curl -L "https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-$(uname -s)-$(uname -m)" -o ~/.docker/cli-plugins/docker-compose
    ```

    > This command installs Compose V2 for the active user under $HOME directory. To install for all 
    > users on your system, replace `~/.docker/cli-plugins` by `/usr/local/lib/docker/cli-plugins`. 

2.  Apply executable permissions to the binary:

    ```console
    $ sudo chmod +x ~/.docker/cli-plugins/docker-compose
    ```
    
3.  Test the installation.

    ```console
    $ docker compose --version
    Docker Compose version {{site.compose_version}}
    ```

## Uninstallation

To uninstall Docker Compose if you installed using `curl`:

```console
$ sudo rm ~/.docker/cli-plugins/docker-compose
```

</div>

<div id="pre-release" class="tab-pane fade" markdown="1">

### Install pre-release builds

If you're interested in trying out a pre-release build, you can download release
candidates from the [Compose repository release page on GitHub](https://github.com/docker/compose/releases){: target="_blank" rel="noopener" class="_"}.
Follow the instructions from the link, which involves running the `curl` command
in your terminal to download the binaries.

> Pre-release builds allow you to try out new features before they are released,
> but may be less stable.
{: .important}

</div>
</div>

----


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
