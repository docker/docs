---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Overview
redirect_from:
- /desktop/dashboard/
---

When you open Docker Desktop, the Docker Dashboard displays.

![dashboard](../images/dashboard.PNG){:width="750px"}

The **Containers** view provides a runtime view of all your containers and applications. It allows you to interact with containers and applications, and manage the lifecycle of your applications directly from your machine. This view also provides an intuitive interface to perform common actions to inspect, interact with, and manage your Docker objects including containers and Docker Compose-based applications. For more information, see [Explore running containers and applications](container.md).

The **Images** view displays a list of your Docker images and allows you to run an image as a container, pull the latest version of an image from Docker Hub, and inspect images. It also displays a summary of the vulnerability scanning report using Snyk. In addition, the **Images** view contains clean-up options to remove unwanted images from the disk to reclaim space. If you are logged in, you can also see the images you and your organization have shared on Docker Hub. For more information, see [Explore your images](images.md).

The **Volumes** view displays a list of volumes and allows you to easily create and delete volumes and see which ones are being used. For more information, see [Explore volumes](volumes.md).

In addition, the Docker Dashboard allows you to:

- Easily navigate to the **Preferences** (**Settings** in Windows) menu to configure Docker Desktop preferences. Select the **Preferences** or **Settings** icon in the Dashboard header.
- Access the **Troubleshoot** menu to debug and perform restart operations. Select the **Troubleshoot** icon in the Dashboard header.

## The Whale menu

Docker Desktop also provides an easy-access tray icon that appears in the taskbar and is reffered to as the Whale menu ![whale menu](../../assets/images/whale-x.svg){: .inline}.

To display the Whale menu, right-click on the ![whale menu](../../assets/images/whale-x.svg){: .inline} icon. It displays the following options:

- **Dashboard**. This takes you to the Docker Dashboard.
- **Sign in/Create Dock Id**
- **Settings**. Takes you to the Settings, or Preferences is you use a macOS.
- **Switch to Windows containers**
- **About Docker Desktop**. Contains information on the versions you are running, and links to the Subscription Service Agreement for example. 
- **Documentation**
- **Quick Start Guide**. Launches the Quick Start Guide.
- **Docker Hub**
- **Extensions**
- **Kubernetes**
- **Pause**
- **Restart**
- **Quit Docker Desktop**