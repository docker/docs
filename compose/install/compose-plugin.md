---
description: How to install Docker Compose CLI plugin
keywords: compose, orchestration, install, installation, docker, documentation
toc_max: 3

title: Install Docker Compose CLI plugin
---

On this page you can find instructions on how to install the Compose plugin for Docker CLI on Linux and Windows Server operating systems.
>Note that installing Docker Compose as a plugin requires Docker CLI.

## Installing Compose on Linux systems

In this section, you can find various methods for installing Compose on Linux.

### Installation methods

* [Installing Docker Desktop for Linux](../../desktop/install/linux-install.md/){:target="_blank" rel="noopener" class="_"} is the easiest and recommended installation route. 
Check the Desktop for Linux [supported platforms](../../desktop/install/linux-install.md/#supported-platforms){:target="_blank" rel="noopener" class="_"} page to verify the supported Linux distributions and architectures.


The following other methods are possible:

* __Using the automated convenience scripts__ (for testing and development environments). 
These scripts install Docker Engine and Docker CLI with the Compose plugin. 
For this route, go to the [Docker Engine install](../../../engine/install/){:target="_blank" rel="noopener" class="_"} page and follow the provided instructions. _After installing Desktop for Linux, this is the recommended route._

* __Setting up Docker's repository__ and using it to install Docker CLI Compose plugin. See the [Install using the repository](#install-using-the-repository) section on this page. _This is the second best route._

* __Installing the Docker CLI Compose plugin manually__. See the [Install the plugin manually](#install-the-plugin-manually) section on this page. _Note that this option requires you to manage upgrades manually as well._ 


### Install using the repository

> **Note**
>
>These instructions assume you already have Docker Engine and Docker CLI installed and now want to install the Compose plugin. 
For other Linux installation methods see [this summary](#installation-methods).

>To run Compose as a non-root user, see [Manage Docker as a non-root user](../../engine/install/linux-postinstall.md){:target="_blank" rel="noopener" class="_"}.


If you have already set up the Docker repository jump to step 2.

1. Set up the repository. Go to the "Set up the repository" section of the chosen [Linux distribution](../../engine/install/index.md#server){:target="_blank" rel="noopener" class="_"}. found on the Docker Engine installation pages to check the instructions.

2. Update the `apt` package index, and install the _latest version_ of Docker Compose:

    > Or, if using a different distro, use the equivalent package manager instructions. 


    ```console
    $ sudo apt-get update
    $ sudo apt-get install docker-compose-plugin
    ```
    
    Alternatively, to install a specific version of Compose CLI plugin:
      
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
    where `<VERSION_STRING>` is, for example,`2.3.3~ubuntu-focal`.

3.  Verify that Docker Compose is installed correctly by checking the version.

    ```console
    $ docker compose version
    Docker Compose version v2.3.3
    ```

### Install the plugin manually

> **Note**
>
> These instructions assume you already have Docker Engine and Docker CLI installed and now want to install the Compose plugin. 
>
> Note as well this option requires you to manage upgrades manually. Whenever possible we recommend any of the other installation methods listed. For other Linux installation methods see [this summary](#installation-methods).

>To run Compose as a non-root user, see [Manage Docker as a non-root user](../../engine/install/linux-postinstall.md).


1.  To download and install the Compose CLI plugin, run:

    ```console
    $ DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
    $ mkdir -p $DOCKER_CONFIG/cli-plugins
    $ curl -SL https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/cli-plugins/docker-compose
    ```

    This command downloads the latest release of Docker Compose (from the Compose releases repository) and installs Compose for the active user under `$HOME` directory.
    
    > To install:
    >* Docker Compose for _all users_ on your system, replace `~/.docker/cli-plugins` with `/usr/local/lib/docker/cli-plugins`.
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

> **Note**
>
>__Compose standalone__: If you need to use Compose without installing the Docker CLI, the instructions for the standalone scenario are similar. 
> Note the target folder for the binary's installation is different as well as the compose syntax used with the plugin (_space compose_) or the standalone version (_dash compose_).

1. To download and install Compose standalone, run:
  ```console
  $ curl -SL https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
  ```
2. Apply executable permissions to the standalone binary in the target path for the installation.
3. Test and execute compose commands using `docker-compose`.

> **Note**
>
> If the command `docker-compose` fails after installation, check your path.
> You can also create a symbolic link to `/usr/bin` or any other directory in your path. 
> For example:
> ```console
> $ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
> ```


## Install Compose on Windows Server

Follow these instructions if you are running the Docker daemon and client directly
on Microsoft Windows Server and want to install Docker Compose.


1.  Run a PowerShell as an administrator. 
When asked if you want to allow this app to make changes to your device, click **Yes** in order to proceed with the installation.
    
2.  GitHub now requires TLS1.2. In PowerShell, run the following:
    
    ```powershell  
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    ```
3. Run the following command to download the latest release of Compose ({{site.compose_version}}):

    ```powershell
    Invoke-WebRequest "https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-Windows-x86_64.exe" -UseBasicParsing -OutFile $Env:ProgramFiles\Docker\docker-compose.exe
    ```

    > **Note**
    >
    > On Windows Server 2019 you can add the Compose executable to `$Env:ProgramFiles\Docker`.
     Because this directory is registered in the system `PATH`, you can run the `docker-compose --version` 
     command on the subsequent step with no additional configuration.

    > To install a different version of Compose, substitute `{{site.compose_version}}`
    > with the version of Compose you want to use.

3.  Test the installation.

    ```console
    $ docker compose version
    Docker Compose version {{site.compose_version}}
    ```

