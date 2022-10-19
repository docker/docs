---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local
title: Overview
---
> **Beta**
>
> The Dev Environments feature is currently in [Beta](../../release-lifecycle.md#beta). We recommend that you do not use this in production environments.

Dev Environments lets you create a configurable developer environment with all the code and tools you need to quickly get up and running. 

It uses tools built into code editors that allows Docker to access code mounted into a container rather than on your local host. This isolates the tools, files and running services on your machine allowing multiple versions of them to exist side by side.

>**Changes to Dev Environments with Docker Desktop 4.13**
>
>Docker has simplified how you configure your dev environment project. All you need to get started is a `compose-dev.yaml` file. If you have an existing project with a `.docker/` folder this is automatically migrated the next time you launch.
{: .important}

![Dev environment intro](../images/dev-env.PNG){:width="700px"}

## Prerequisites

Dev Environments is available as part of Docker Desktop 3.5.0 release. Download and install **Docker Desktop 3.5.0** or higher:

- [Docker Desktop](../release-notes.md)

To get started with Dev Environments, you must also install the following tools and extension on your machine:

- [Git](https://git-scm.com){:target="_blank" rel="noopener" class="_"}
- [Visual Studio Code](https://code.visualstudio.com/){:target="_blank" rel="noopener" class="_"}
- [Visual Studio Code Remote Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers){:target="_blank" rel="noopener" class="_"}

### Add Git to your PATH on Windows

If you have already installed Git, and it's not detected properly, run the following command to check whether you can use Git with the CLI or PowerShell:

`$ git --version`

If it doesn't detect Git as a valid command, you must reinstall Git and ensure you choose the option  **Git from the command line...** or the **Use Git and optional Unix tools...**  on the **Adjusting your PATH environment**  step.

![Windows add Git to path](../images/dev-env-gitbash.png){:width="300px"}

> **Note**
>
> After Git is installed, restart Docker Desktop. Select **Quit Docker Desktop**, and then start it again.

## Known issues

The following section lists known issues and workarounds:

1. When sharing a dev environment between Mac and Windows, the VS Code terminal may not function correctly in some cases. To work around this issue, use the Exec in CLI option in the Docker Dashboard.

## What's next?

Learn how to:
- [Create a simple dev environment](create-dev-env.md)
- [Create an advanced dev environment](create-compose-dev-env.md)
- [Distribute your dev environment](share.md)
