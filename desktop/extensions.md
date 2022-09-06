---
description: Extensions
keywords: Docker Extensions, Docker Desktop, Linux, Mac, Windows
title: Docker Extensions
toc_min: 1
toc_max: 2
---

> **Beta**
>
> The Docker Extensions feature is currently in [Beta](../release-lifecycle.md#beta). We recommend that you do not use Docker Extensions in production environments.

Docker Extensions lets you use third-party tools within Docker Desktop to extend its functionality. There is no limit to the number of extensions you can install.

![extenstions](images/extensions-marketplace.PNG){:width="750px"}

You can explore the list of available extensions in [Docker Hub](https://hub.docker.com/search?q=&type=extension) or in the Extensions Marketplace within Docker Desktop.

Docker Community members and partners can use our [SDK](extensions-sdk/index.md) to create new extensions.

To find out more about Docker Extensions, we recommend the video walkthrough from DockerCon 2022:

<iframe width="560" height="315" src="https://www.youtube.com/embed/3rAGXS8pszQ" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

## Prerequisites

Docker Extensions is available as part of Docker Desktop 4.8.0 or a later release. [Download and install Docker Desktop 4.8.0 or later](release-notes.md).

## Add an extension

>**Note**
>
> For some extensions, a separate account needs to be created before use.

To add Docker Extensions:

1. Open Docker Desktop.
2. From the Dashboard, select **Add Extensions** in the menu bar. 
The Extensions Marketplace opens. 
3. Browse the available extensions.
    You can sort the list of extensions by **Recently added** or alphabetically. 
4. Click **Install**.

From here, you can click **Open** to access the extension or install additional extensions. The extension also appears in the menu bar.

## See containers created by extensions

By default, containers created by extensions are hidden from the list of containers in Docker Dashboard and the Docker CLI. To make them visible 
update your settings:

1. Navigate to  **Settings**, or **Preferences** if you're a Mac user.
2. Select the **Extensions** tab.
3. Next to **Show Docker Extensions system containers**, select or clear the checkbox to set your desired state.
4. In the bottom-right corner, click **Apply & Restart**.

## Enable or disable extensions available in the Marketplace

Docker Extensions are switched on by default. To change your settings:

1. Navigate to  **Settings**, or **Preferences** if you're a Mac user.
2. Select the **Extensions** tab.
3. Next to **Enable Docker Extensions**, select or clear the checkbox to set your desired state.
4. In the bottom-right corner, click **Apply & Restart**.

## Enable or disable extensions not available in the Marketplace

You can install Docker Extensions through the Marketplace or through the Extensions SDK tools. You can choose to only allow published extensions (that have been published in the Extensions Marketplace).

1. Navigate to **Settings**, or **Preferences** if you're a Mac user.
2. Select the **Extensions** tab.
3. Next to **Allow only extensions distributed through the Docker Marketplace**, select or clear the checkbox to set your desired state.
4. In the bottom-right corner, click **Apply & Restart**.

## Update an extension
You can update Docker Extensions outside of Docker Desktop releases. To update an extension to the latest version:

1. Navigate to Docker Dashboard, and from the menu bar select the ellipsis to the right of **Extensions**.
2. Click **Manage Extensions**.
If an extension has a new version available, it displays an **Update** button.
3. Click **Update**.

## Submit feedback
Feedback can be given to an extension author through a dedicated Slack channel or Github. To submit feedback about a particular extension:

1. Navigate to Docker Dashboard and from the menu bar select the ellipsis to the right of **Extensions**.
2. Click **Manage Extensions**.
3. Select the extension you want to provide feedback on. 
4. Scroll down to the bottom of the extension's description and, depending on the 
extension, select:
    - Support
    - Slack
    - Issues. You'll be sent to a page outside of Docker Desktop to submit your feedback.

If an extension does not provide a way for you to give feedback, contact us and we'll pass on the feedback for you. To provide feedback, select the **Give feedback** to the right of **Extensions Marketplace**

## Uninstall an extension
 You can uninstall an extension at any time. 
 
 > **Note**  
 >
 > Any data used by the extension that is stored in a volume must be manually deleted. 

1. From the menu bar, select the ellipsis to the right of **Extensions**.
2. Click **Manage Extensions**. This displays a list of extensions you've installed.
3. Click **Uninstall**.

## Install an extension not available in the Marketplace

> **Warning**
>
> Docker Extensions that are not in the Marketplace haven't gone through Docker's review process.
> Extensions can install binaries, invoke commands and access files on your machine. Installing them is at your own risk.
{: .warning}

The Extensions Marketplace is the trusted and official place to install extensions from within Docker Desktop. These extensions have gone through a review process by Docker. However, other extensions can also be installed in Docker Desktop if you trust the extension author.

Given the nature of a Docker Extension (i.e. a Docker image) you can find other places where users have their extension's source code published. For example on GitHub, GitLab or even hosted in image registries like DockerHub or GHCR.
You can install an extension that has been developed by the community or internally at your company from a teammate. You are not limited to installing extensions just from the Marketplace.

> **Note**
>
> Ensure the option **Allow only extensions distributed through the Docker Marketplace** is disabled. Otherwise, this prevents any extension not listed in the Marketplace, via the Extension SDK tools from, being installed.
> You can change this option in the Settings > Extensions.

To install an extension which is not present in the Marketplace, you can use the Extensions CLI that is bundled with Docker Desktop.

In a terminal, type `docker extension install IMAGE[:TAG]` to install an extension by its image reference and optionally a tag. Use the `-f` or `--force` flag to avoid interactive confirmation.

Go to the Docker Dashboard to see the new extension installed.

## List installed extensions

Regardless whether the extension was installed from the Marketplace or manually by using the Extensions CLI, you can use the `docker extension ls` command to display the list of extensions installed.
As part of the output you'll see the extension ID, the provider, version, the title and whether it runs a backend container or has deployed binaries to the host, for example:

```
$ docker extension ls
ID                  PROVIDER            VERSION             UI                    VM                  HOST
john/my-extension   John                latest              1 tab(My-Extension)   Running(1)          -
```

Go to the Docker Dashboard, click on **Add Extensions** and on the **Installed** tab to see the new extension installed.
Notice that an `UNPUBLISHED` label displays which indicates that the extension has not been installed from the Marketplace.

## Update an extension which is not in the Marketplace

To update an extension which is not present in the Marketplace, in a terminal type `docker extension update IMAGE[:TAG]` where the `TAG` should be different from the extension that is already installed.

For instance, if you installed an extension with `docker extension install john/my-extension:0.0.1`, you can update it by running `docker extension update john/my-extension:0.0.2`.
Go to the Docker Dashboard to see the new extension updated.

> **Note**
>
> Extensions that have not been installed through the Marketplace will not receive update notifications from Docker Desktop.

## Uninstall an extension which is not in the Marketplace

To uninstall an extension which is not present in the Marketplace, you can either navigate to the **Installed** tab in the Marketplace and select the **Uninstall** button, or from a terminal type `docker extension uninstall IMAGE[:TAG]`.
