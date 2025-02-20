---
title: How it works
description: Understand how Docker Projects works 
keywords: docker projects
weight: 10
---

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