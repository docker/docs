---
description: How to install Docker Compose
keywords: compose, orchestration, install, installation, docker, documentation
title: Install Docker Compose
toc_max: 2
redirect_from:
- /compose/compose-desktop/
---

On this page you can find information on how to get and install Compose.

## Install Compose

**To get Compose:**

* If you have Docker Desktop, you've got a full Docker installation, including Compose.

* If you don't have Docker yet installed, or just need to install Compose, check the options below:
{% assign yes = '![yes](/images/green-check.svg){: .inline style="height: 14px; margin: 0 auto; align=right"}' %}

| Platform       | Docker Desktop                 | Compose CLI Plugin      | Compose Standalone  |
|:---------------|:------------------------------:|:-----------------------:|:-------------------:|
|Linux (64b)     |{{ yes }} [Install](../../desktop/install/linux-install.md)|{{ yes }} [Install](./compose-linux.md)|{{ yes }} [Install](./compose-linux.md#install-compose-standalone)|
|Mac             |{{ yes }} [Install](../../desktop/install/mac-install.md)    | - | - |
|Windows         |{{ yes }} [Install](../../desktop/install/windows-install.md)| - | - |
|Windows Server  | - | - |{{ yes }} [Install](#install-compose-on-windows-server)|

**Table description:**
* **Linux(64bit):** *option1)* Install [Docker Desktop for Linux](../../desktop/install/linux-install.md) or, *option2)* [Install Docker Engine and the CLI](../../engine/install/#server) and then [install Compose](compose-linux.md).
* **macOS:** [Install Docker Desktop](../../desktop/install/mac-install/).
* **Windows:** [Install Docker Desktop](../../desktop/install/windows-install/).
* **Windows Server:** [Install Docker Engine and the CLI](../../engine/install/binaries/#install-server-and-client-binaries-on-windows) and then [install Compose](#install-compose-on-windows-server).

<!--
 * The best way to get all this is Docker Desktop. See [Docker Desktop's Overview](../../desktop/) for more information, or click on [Download and Install](../../desktop/#download-and-install){:target="_blank" rel="noopener" class="_"} to go to the install instructions.
* To perform a complete installation of Docker from the command line, go to the [Docker Engine install](../../../engine/install/){:target="_blank" rel="noopener" class="_"} page and follow the provided instructions.
-->

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

4.  Test the installation.

    ```console
    $ docker compose version
    Docker Compose version {{site.compose_version}}
    ```

## Where to go next

- [Getting Started](../gettingstarted.md)
- [Command line reference](../../reference/index.md)
- [Compose file reference](../compose-file/index.md)
- [Sample apps with Compose](../samples-for-compose.md)

