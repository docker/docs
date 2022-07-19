---
description: Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Explore Containers
---

The **Containers** view lists all your running containers and applications. You must have running or stopped containers and applications to see them listed.

## Container actions

Use the **Search** field to search for any specific container.

From the **Containers** view you can perform the following actions on one or more containers at once:
- Pause
- Resume
- Stop
- Start
- Delete

When you hover over individual containers, you can also:
- Click **Open in Visual Studio Code** to open the application in VS Code.
- Open the port exposed by the container in a browser.

### Inspect a container

You can obtain detailed information about the container when you select a container.

The **container view** displays **Logs**, **Inspect**, and **Stats** tabs and provides quick action buttons to perform various actions.

- Select **Logs** to see logs from the container. You can also:
    - Use `Cmd + f`/`Ctrl + f` to open the search bar and find specific entries. Search matches are highlighted in yellow.
    - Press `Enter` or `Shit + Enter` to jump to the next or previous search match respectively. 
    - Use the **Copy** icon in the top right-hand corner to copy all the logs to your clipboard.
    - Automatically copy any logs content by highlighting a few lines or a section of the logs.
    - Use the **Clear terminal** icon in the top right-hand corner to clear the logs terminal. 
    - Select and view external links that may be in your logs. 


- Select **Inspect** to view low-level information about the container. You can see the local path, version number of the image, SHA-256, port mapping, and other details.

- Select **Stats** to view information about the container resource utilization. You can see the amount of CPU, disk I/O, memory, and network I/O used by the container.
