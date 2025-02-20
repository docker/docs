---
title: Docker Projects
params:
  sidebar:
    group: Products
    badge:
      color: blue
      text: Beta
weight: 50
description: Learn how to use Docker Projects which provides a unified, project-based workflow to run your containerized projects. 
keywords: Docker, projects, docker deskotp, containerization
grid:
- title: Open a new project
  description: Learn how to open a new local or remote projects.
  icon: checklist
  link: /projects/open/
- title: Edit your project
  description: Edit your project's run commands and setup. N
  icon: design_services
  link: /projects/edit/
- title: Manage your projects
  description: Run or remove your projects.
  icon: tune
  link: /projects/manage/
---

Docker Projects provides a simplified, project-based workflow for running and managing containerized applications. It organizes your code, configurations, and logs across local and cloud environments into a single view, making it easy to collaborate and share across teams.

A project organizes your code and Docker artifacts into a single object. These artifacts include logs as well as customizable run commands. These artifacts can persist remotely in the cloud, which lets you access your projects from any device that has Docker Desktop.

### Key features and benefits

 - One-click project setup: Open a local folder or clone a Git repository and run your project instantly.
 - Minimal Docker expertise required: Ideal for both beginners and experienced engineers.
 - Custom `run` commands for your projects: Define and store preconfigured `run` commands that are equivalent to running `docker compose up`.
 - Local & remote projects:  Work on projects locally or sync artifacts to the cloud for cross-device access and easy collaboration.

## How it works 

Docker Projects requires a Compose file (docker-compose.yml) to define your application's services, networks, and configurations. When you open a project, Docker Projects automatically detects the Compose file, allowing you to configure and run services with pre-set commands. 

By integrating with Docker Compose, Docker Projects ensures a consistent, easy-to-manage workflow for both individual developers and teams. Whether you're starting a new project, configuring it, or collaborating with a team, Docker Projects keeps the process simple.

1. Create or open a project. You can:

 - Open a local project: Select a folder on your machine that contains your project code.
 - Clone a Git repository:Provide a repository URL and clone the project into a local directory.

Once a project is opened, Docker Desktop detects the Compose file and prepares the project for execution.

2. Configure and run your project with pre-configured commands. These commands:

 - Work like `docker compose up`, launching services based on the Compose file.
 - Can be customized with additional flags, multiple Compose files, and environment variables.
 - Allow pre-run tasks, such as executing scripts before starting the services.

All of which means you can fine-tune your configurations without manually running complex CLI commands.

3. Collaborate and share with teams. For projects linked to a Git repository, Docker Projects stores artifacts in the cloud, enabling easy collaboration:

 - Work across devices: Open a project from any machine and instantly access stored configurations.
 - Share configurations: Team members can access predefined run commands, reducing setup time.

Collaboration is seamless—new developers can join a team, open a project, and start working without complex setup steps.

4. Manage and iterate. Once a project is up and running, Docker Projects makes it easy to monitor, update, and troubleshoot:

 - View logs to debug issues and track service activity.
 - Edit configurations and run commands as requirements change.

{{< grid >}}
