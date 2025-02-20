---
title: Manage your projects
description: Learn how to delete or completely remove a project.
keywords: Docker, projects, docker deskotp, containerization, open, remote, local
weight: 40
---

## Run a project or service

1. Open a new or existing project.

2. Select a run command from the dropdown menu.

3. Select the **Run** button for the project or the **Play** button next to the service you'd like to run.

## Stop or restart a project or service

1. Open an existing project that is running.

2. Select the **Stop** or **Restart** button for the project or the appropriate button next to the service.

## Remove a project from Docker Desktop

If a project is associated with a Git repository, you can remove it from Docker Desktop.  When a project is deleted, you can no longer run the project from the **Projects** view, but its run configuration still exists remotely in the cloud. 

This means that you can later [open the project](/manuals/projects/open.md#open-an-existing-remote-project) and associate it with the remote run configuration without having to specify the run command again.
None of your local code is deleted when removing a project from Docker Desktop.

To remove a project from Docker Desktop:
Sign in to Docker Desktop, and go to Projects.
Select the options menu () next to the project, and then select Remove from Docker Desktop.

## Delete a project

Deleting a project removes it from Docker Desktop and deletes all configuration locally and remotely from the cloud. When a project is deleted, you can no longer run the project from the **Projects** view. None of your local code is deleted when you delete a project from Docker Desktop.
