---
description: Learn how to use the Docker Desktop Dashboard within Docker Desktop, including Quick search, the Docker menu, and more
keywords: Docker Desktop Dashboard, manage, containers, gui, dashboard, images, user manual,
  whale menu
title: Explore Docker Desktop
weight: 30
aliases:
- /desktop/dashboard/
---

When you open Docker Desktop, the Docker Desktop Dashboard displays.

![Docker Desktop Dashboard on Containers view](../images/dashboard.png)

It provides a centralized interface to manage your [containers](container.md), [images](images.md), [volumes](volumes.md), and [builds](builds.md).

In addition, the Docker Desktop Dashboard lets you:

- Use [Ask Gordon](/manuals/desktop/features/gordon/_index.md), a personal AI assistant embedded in Docker Desktop and the Docker CLI. It's designed to streamline your workflow and help you make the most of the Docker ecosystem.
- Navigate to the **Settings** menu to configure your Docker Desktop settings. Select the **Settings** icon in the Dashboard header.
- Access the **Troubleshoot** menu to debug and perform restart operations. Select the **Troubleshoot** icon in the Dashboard header.
- Be notified of new releases, installation progress updates, and more in the **Notifications center**. Select the bell icon in the bottom-right corner of the Docker Desktop Dashboard to access the notification center.
- Access the **Learning center** from the Dashboard header. It helps you get started with quick in-app walkthroughs and provides other resources for learning about Docker. 

  For a more detailed guide about getting started, see [Get started](/get-started/introduction/_index.md).
- Access [Docker Hub](/manuals/docker-hub/_index.md) to search, browse, pull, run, or view details
  of images.
- Get to the [Docker Scout](../../scout/_index.md) dashboard.
- Navigate to [Docker Extensions](/manuals/extensions/_index.md).

## Docker terminal

From the Docker Dashboard footer, you can use the integrated terminal directly within Docker Desktop. 

The integrated terminal:

- Persists your session if you navigate to another
  part of the Docker Desktop Dashboard and then return.
- Supports copy, paste, search, and clearing your session.

#### Open the integrated terminal

To open the integrated terminal, either:

- Hover over your running container and under the **Actions** column, select the **Show container actions**
  menu. From the drop-down menu, select **Open in terminal**.
- Or, select the **Terminal** icon located in the bottom-right corner, next to the version number.

To use your external terminal, navigate to the **General** tab in **Settings**
and select the **System default** option under **Choose your terminal**.

## Quick search

Use Quick Search, which is located in the Docker Dashboard header, to search for:

- Any container or Compose application on your local system. You can see an overview of associated environment variables or perform quick actions, such as start, stop, or delete.

- Public Docker Hub images, local images, and images from remote repositories (private repositories from organizations you're a part of in Hub). Depending on the type of image you select, you can either pull the image by tag, view documentation, go to Docker Hub for more details, or run a new container using the image.

- Extensions. From here, you can learn more about the extension and install it with a single click. Or, if you already have an extension installed, you can open it straight from the search results.

- Any volume. From here you can view the associated container.

- Docs. Find help from Docker's official documentation straight from Docker Desktop. 

## The Docker menu

Docker Desktop also includes a tray icon, referred to as the Docker menu {{< inline-image src="../../assets/images/whale-x.svg" alt="whale menu" >}} for quick access.

Select the {{< inline-image src="../../assets/images/whale-x.svg" alt="whale menu" >}} icon in your taskbar to open options such as:

- **Dashboard**. This takes you to the Docker Desktop Dashboard.
- **Sign in/Sign up**
- **Settings**
- **Check for updates**
- **Troubleshoot**
- **Give feedback**
- **Switch to Windows containers** (if you're on Windows)
- **About Docker Desktop**. Contains information on the versions you are running, and links to the Subscription Service Agreement for example.
- **Docker Hub**
- **Documentation**
- **Extensions**
- **Kubernetes**
- **Restart**
- **Quit Docker Desktop**
