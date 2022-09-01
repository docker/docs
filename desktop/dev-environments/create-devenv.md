---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local
title: Create a dev environment
---

This page contains information on how to create a dev environment. It walks you through the recommended steps and links to useful resources that are great starting points. 

## Step one: Preparation and gathering requirements

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

## Step two: Dive into Compose

Once you've gathered your dev environment requirements you're ready to start building your compose.yaml file which powers your dev environment. 

You can either start your compose.yaml file from scratch or use [sample dev environment](https://github.com/docker/awesome-compose) as a starting point then add or amend accordingly. 

## Step three: Test your dev environment

Once you've built your dev environment using Compose, add it to the Dev Environments tab in Docker Desktop.

To add your environment from a Git repository or subfolder of a Git repository:
1. Copy the link of your Git repository or subfolder.
2. In the **Add the URL** field in the Dev Environments tab, paste the link. 
3. Select **Add project**.

## Step four: Distribute to your team

Dev Environments can be easily shared with your team members and provides a snapshot of the current state of your project. They can access the code, any dependencies, and the current Git branch you are working on. They can also review your changes and provide feedback before you create a pull request.

{% include upgrade-cta.html
  body="Docker Pro, Team, and Business users can now share Dev Environments with their team members."
  header-text="This feature requires a paid Docker subscription"
  target-url="https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade"
%}

To share your environment, hover over your Dev Environment, select the **Share** icon, and specify the Docker Hub namespace where you’d like to push your Dev Environment to.

This creates an image of your Dev Environment, uploads it to the Docker Hub namespace you have specified, and provides a tiny URL to share with your team members.

## Step five: Iterate and improve

You can use the dev environment that you have created to test out upgrades in and then roll it out to everyone when you're ready. 

## What's next?

- Learn how to create a project with multiple dev environments

<div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">
    <div class="panel panel-default">
      <div class="panel-heading" role="tab" id="headingSeven">
        <h5 class="panel-title">
          <a role="button" data-toggle="collapse" data-parent="#accordion" href="#collapseSeven" aria-expanded="true" aria-controls="collapseSeven">
            Still to do
            <i class="fa fa-chevron-down"></i>
          </a>
        </h5>
      </div>
      <div id="collapseSeven" class="panel-collapse collapse" role="tabpanel" aria-labelledby="headingSeven">
        <div class="panel-body">
            <p>
            Questions: is sharing still going to be a paid feature?
            To-dos: iron out the appropriate compose links to direct users to.  flesh out each step 
            
            </p>
        </div>
      </div>
    </div>
  </div>
