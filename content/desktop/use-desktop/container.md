---
description: Understand what you can do with the Containers view on Docker Dashboard
keywords: Docker Dashboard, manage, containers, gui, dashboard, images, user manual
title: Explore the Containers view in Docker Desktop
---

The **Containers** view lists all your running containers and applications. You must have running or stopped containers and applications to see them listed.

## Container actions

Use the **Search** field to search for any specific container.

From the **Containers** view you can perform the following actions:
- Pause/Resume
- Stop/Start/Restart
- View image packages and CVEs
- Delete
- Open the application in VS code
- Open the port exposed by the container in a browser
- Copy docker run. This lets you share container run details or modify certain parameters.

## Resource usage

From the **Containers** view you can monitor your containers' CPU and memory usage over time. This can help you understand if something is wrong with your containers or if you need to allocate additional resources. 

When you [inspect a container](#inspect-a-container), the **Stats** tab displays further information about a container's resource utilization. You can see how much CPU, memory, network and disk space your container is using over time.

## Inspect a container

You can obtain detailed information about the container when you select it.

From here, you can use the quick action buttons to perform various actions such as pause, resume, start or stop, or explore the **Logs**, **Inspect**, **Bind mounts**, **Terminal**, **Files**, and **Stats** tabs. 

### Logs

- Select **Logs** to see logs from the container. You can also:
    - Use `Cmd + f`/`Ctrl + f` to open the search bar and find specific entries. Search matches are highlighted in yellow.
    - Press `Enter` or `Shift + Enter` to jump to the next or previous search match respectively. 
    - Use the **Copy** icon in the top right-hand corner to copy all the logs to your clipboard.
    - Automatically copy any logs content by highlighting a few lines or a section of the logs.
    - Use the **Clear terminal** icon in the top right-hand corner to clear the logs terminal. 
    - Select and view external links that may be in your logs. 

### Inspect

Select **Inspect** to view low-level information about the container. It displays the local path, version number of the image, SHA-256, port mapping, and other details.

### Integrated terminal

From the **Exec** tab, lets you use the integrated terminal, on a running container, directly within Docker Desktop. You are able to quickly run commands within your container so you can understand its current state or debug when something goes wrong.

Using the integrated terminal is the same as running `docker exec -it <container-id> /bin/sh`, or `docker exec -it <container-id> cmd.exe` if you are using Windows containers, in your external terminal. It also:

- Automatically detects the default user for a running container from the image's Dockerfile. If no user is specified it defaults to `root`.
- Persists your session if you navigate to another part of the Docker Dashboard and then return.
- Supports copy, paste, search, and clearing your session.

To open the integrated terminal, either:
- Hover over your running container and select the **Show container actions** menu. From the drop-down menu, select **Open in terminal**.
- Select the container and then select the **Exec** tab.

To use your external terminal, navigate to the **General** tab in **Settings** and select the **System default** option under **Choose your terminal**.

### Files

- Select **Files** to explore the filesystem of running or stopped containers. You can also:
    - See which files have been recently added, modified, or deleted
    - Edit a file straight from the built-in editor
    - Drag and drop files and folders between the host and the container
    - Delete unnecessary files when you right-click on a file
    - Download file and folders from the container straight to the host

## Additional resources
- [What is a container](../../guides/walkthroughs/what-is-a-container.md)
- [Run a container](../../guides/walkthroughs/run-a-container.md)
- [Run multi-container applications](../../guides/walkthroughs/multi-container-apps.md)