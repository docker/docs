---
title: View your project
description: View information about your project or the services within your project.
keywords: containers, docker projects, local, remote, docker desktop
weight: 30
---

## View a project’s README

If a project has a README file, you can view it via Docker Projects. Note that the README tab is only visible if the project has a README file.

To view a project’s README, open a new or existing project and then select the README tab.

## View logs for a project

1. Open a new or existing project.

2. Select the **Logs** tab to see all project logs.

3. Optionally, use the menu in the top right corner of the logs to copy the logs to your clipboard or clear the logs.

## View service-level information 

With Docker Projects, you can view the following information about your containers within your project:

 - Logs
 - Image
 - Files
 - Network
 - Environment variables

From the **Exec** tab, you can use the integrated terminal, on a running container, directly within Docker Desktop. You are able to quickly run commands within your container so you can understand its current state or debug when something goes wrong.

### Logs

Select **Logs** to see logs from the containers in your project. You can also:

- Use `Cmd + f`/`Ctrl + f` to open the search bar and find specific entries.
  Search matches are highlighted in yellow.
- Press `Enter` or `Shift + Enter` to jump to the next or previous search match
  respectively. 
- Use the **Copy** icon in the top right-hand corner to copy all the logs to
  your clipboard.
- Automatically copy any logs content by highlighting a few lines or a section
  of the logs.
- Use the **Clear terminal** icon in the top right-hand corner to clear the
  logs terminal. 
- Select and view external links that may be in your logs. 

### Image

The **Image** tab in Docker Projects provides details about the Docker image associated with a service. It helps you verify which image is being used, when it was last built, and where the corresponding Dockerfile is located.

It also provides quick access to inspect the image or open the Dockerfile for modifications.

### Files

Select **Files** to explore the filesystem of running or stopped containers in your project. You
can also:

   - See which files have been recently added, modified, or deleted
   - Edit a file straight from the built-in editor
   - Drag and drop files and folders between the host and the container
   - Delete unnecessary files when you right-click on a file
   - Download files and folders from the container straight to the host

### Network

The **Network** tab in Docker Projects provides an overview of how the containerized services communicate with each other and the host system. It displays the assigned network name, connected services, and mapped container ports.

If a service is mapped to a host port, you can select the link to open it in a browser

### Environment variables

The **Env** tab in Docker Projects displays the environment variables available to a service. These variables help configure the runtime environment without modifying the container image.

## What's next?

 - [Add or edit your run commands](/manuals/projects/edit.md)
 - [Manage your projects](/manuals/projects/manage.md)
