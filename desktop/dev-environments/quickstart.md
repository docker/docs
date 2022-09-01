---
description: quickstart guide for dev environments
keywords: quickstart, dev environments
title: Quickstart guide
---

The quickstart guide is a useful way of learning the basics of Dev Environments before getting into it more deeply.

The quickstart guide walks you through a scenario where you've been onboarded onto a new project that is building a mobile app and you have multiple environments in which to test your code changes. As a new joiner youve been tasked with TASK A. 

## Step one: Preparation

Dev Environments is available as part of Docker Desktop 3.5.0 release. Download and install **Docker Desktop 3.5.0** or higher:

- [Docker Desktop](../release-notes.md)

To get started with Dev Environments, you must also install the following tools and extension on your machine:

- [Git](https://git-scm.com){:target="_blank" rel="noopener" class="_"}
- [Visual Studio Code](https://code.visualstudio.com/){:target="_blank" rel="noopener" class="_"}
- [Visual Studio Code Remote Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers){:target="_blank" rel="noopener" class="_"}

### Add Git to your PATH on Windows

If you have already installed Git, and it's not detected properly, run the following command to check whether you can use Git with the CLI or PowerShell:

`$ git --version`

If it doesn't detect Git as a valid command, you must reinstall Git and ensure you choose the option  **Git from the command line...** or the **Use Git and optional Unix tools...**  on the **Adjusting your PATH environment**  step.

![Windows add Git to path](../images/dev-env-gitbash.png){:width="300px"}

> **Note**
>
> After Git is installed, restart Docker Desktop. Select **Quit Docker Desktop**, and then start it again.

## Step two: Explore the project information 

Showcase the README, contributing.md etc. 

## Step three: Launch the sample project

From the Dev Environments tab in Docker Desktop, hover over the sample project and select the **Run** icon. 

## Step four: Make changes to the code in the sandbox environment

## Step five: Switch to the product environment and add your changes

## Step six: Share your code changes 






## What's next
- Increase your familiarity with Dev Environments by running other sample projects. Select the **See more samples** link at the bottom-right corner. 
- Learn how to add a dev environment
- Learn how to add a project with multiple dev environments


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
            Questions: is the Git PATH item still relevant? I'm not sure it is but need to check. So we can slim down the preparation step. 
            To-dos: need to flesh out the quickstart in more detail, am awaiting user research results
            Thoughts: the quickstart would also familiarise the user with the language of Dev Envs 
            </p>
        </div>
      </div>
    </div>
  </div>
