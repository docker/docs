---
description: Extensions
keywords: Extensions, Docker Desktop, Linux, Mac, Windows
title: Extensions
toc_min: 1
toc_max: 2
---

Docker Extensions enable you to use third-party tools within Docker Desktop to extend its functionality. Docker Community members and partners can use our SDK  to create new extensions. There is no limit to the number of extensions you can install.

> **Preview**
>
>The Docker Extensions feature is currently offered as a Preview. We recommend that you do not use this in production environments.

SCREENSHOT

## Prerequisites

Docker Extensions are available as part of Docker Desktop <insert release number> or a later release. Download and install Docker Desktop <insert release number> or later:

* [Mac](mac/release-notes/index.md)
* [Windows](windows/release-notes/index.md)
* [Linux](linux/index.md)

## Add a Docker Extension

To add Docker Extensions:

1. Open Docker Desktop.
2. From the Dashboard, select **Add extensions** in the menu bar. 
The Extensions Marketplace opens. 
2. Browse the available extensions.
3. Click **Install**.
From here, you can click **Open** to access the extension or install additional extensions. The extension also appears in the menu bar.

## Enable or disable extensions not available in the Marketplace

 Docker Extensions are switched on by default. To change your settings:

1. From the  Docker menu select  **Preferences**.
2. Navigate to the **Extensions** tab.
3. Next to **Enable Docker Extensions**, select the checkbox to to set your desired state.
4. Click **Apply & Restart**.

## Enable or Disable extensions not available in the Marketplace

You can install Docker Extensions through the Marketplace or through the Extensions SDK tools. You can choose to only allow published extensions (that have been published in the Extensions Marketplace).

1. From the Docker menu select **Preferences**.
2. Navigate to the Extensions tab.
3. Next to **Allow only extensions distributed through the Docker Marketplace**, select the checkbox to set your desired state.
4. Click **Apply & Restart**.

## Update a Docker Extension
You can update Docker Extensions outside of Docker Desktop releases. To update an extension to the latest version:

1. From the menu bar, select the ellipsis to the right of **Extensions**.
2. Click **Manage Extensions**.
If an extension has a new version available, an **Update** button is visible.
3. Click **Update**.

## Uninstall a Docker Extension
 You can uninstall an extension at any time. 
 
 > **Note**  
 >
 >Any data used by the extension that is stored in a volume must be manually deleted. 

1. From the menu bar, select the ellipsis to the right of **Extensions**.
2. Click **Manage Extensions**.
3. Click **Uninstall**.
