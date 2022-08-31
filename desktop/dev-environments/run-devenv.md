---
description: Dev Environments
keywords: Dev Environments, local, run
title: Launch a dev environment
---

Remove the challenge of figuring out what you need in order to get set up properly with a new project or to review a team member's PR.

This page contains information on how to add and launch a single dev environment.

You can launch a dev environment from a:
- Git repository
- Branch or tag of a Git repository
- Subfolder of a Git repository
- Local folder

This did not conflict with any of the local files or local tooling set up on your host.  

## Step one: Preparation 

Dev Environments is available as part of Docker Desktop 3.5.0 release. Download and install **Docker Desktop 3.5.0** or higher:

- [Docker Desktop](../release-notes.md)

To get started with Dev Environments, you must also install the following tools and extension on your machine:

- [Git](https://git-scm.com){:target="_blank" rel="noopener" class="_"}
- [Visual Studio Code](https://code.visualstudio.com/){:target="_blank" rel="noopener" class="_"}
- [Visual Studio Code Remote Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers){:target="_blank" rel="noopener" class="_"}

#### Add Git to your PATH on Windows

If you have already installed Git, and it's not detected properly, run the following command to check whether you can use Git with the CLI or PowerShell:

`$ git --version`

If it doesn't detect Git as a valid command, you must reinstall Git and ensure you choose the option  **Git from the command line...** or the **Use Git and optional Unix tools...**  on the **Adjusting your PATH environment**  step.

> **Note**
>
> After Git is installed, restart Docker Desktop. Select **Quit Docker Desktop**, and then start it again.

## Step two: Add your environment to the Dev Environment tab in Docker Desktop

To add your environment from a Git repository or subfolder of a Git repository:
1. Copy the link of your Git repository or subfolder.
2. In the **Add the URL** field in the Dev Environments tab, paste the link. 
3. Select **Add project**.

>Note
>
>Currently, Dev Environments is not able to detect the main language of the subdirectory. You need to define your own base image or compose file in a .docker folder located in your subdirectory. For more information on how to configure, see the [React application with a Spring backend and a MySQL database sample](https://github.com/docker/awesome-compose/tree/master/react-java-mysql) or the [Go server with an Nginx proxy and a Postgres database sample](https://github.com/docker/awesome-compose/tree/master/nginx-golang-postgres). 

To add your environment from a specific branch or tag:
1. Copy the link of your Git repository and add `@mybranch` or `@tag` as a suffix to your Git URL.
2. In the **Add the URL** field in the Dev Environments tab, paste the link. 
3. Select **Add project**.

Your environment is added and Docker detects the main language of your repository, clones the Git code inside a volume, determines the best image for your Dev Environment, and opens VS Code inside the Dev Environment container. (IS THIS STILL CORRECT?)

To add your environment from a local folder, select the **Browse** button to the right of the **Add the URL** field. Once your environment has been added, Docker detects the main language of your local folder, creates a Dev Environment using your local folder, and bind-mounts your local code in the Dev Environment. It then opens VS Code inside the Dev Environment container.

> **Note**
>
> When using a local folder for a Dev Environment, file changes are synchronized between your Dev Environment container and your local files. This can affect the performance inside the container, depending on the number of files in your local folder and the operations performed in the container.

## Step three: Launch your environment



## Optional step four: Share your Dev Environment

Dev Environments can be easily shared with your team members and provides a snapshot of the current state of your project. They can access the code, any dependencies, and the current Git branch you are working on. They can also review your changes and provide feedback before you create a pull request.

{% include upgrade-cta.html
  body="Docker Pro, Team, and Business users can now share Dev Environments with their team members."
  header-text="This feature requires a paid Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade"
%}

To share your environment, hover over your Dev Environment, select the **Share** icon, and specify the Docker Hub namespace where you’d like to push your Dev Environment to.

This creates an image of your Dev Environment, uploads it to the Docker Hub namespace you have specified, and provides a tiny URL to share with your team members.

## What's next?
- Learn how to add and launch a project with multiple dev environments




