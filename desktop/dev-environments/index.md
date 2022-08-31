---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local
title: Overview
---

> **Beta**
>
> The Dev Environments feature is currently in [Beta](../../release-lifecycle.md#beta). We recommend that you do not use this in production environments.

With Dev Environments, you can:

- Customize dev environments with everything you need for your current project.
- Reduce the time investment needed to onboard onto new projects.
- Hit the ground running with ready-to-code environments.
- Avoid dependency conflicts when switching between projects or reviewing a team member's code.

### Who is Dev Environments for?

You'll find Dev Environments useful if you are a developer that:
- Needs to get up and running with a project or multiple projects quickly and easily.
- Has limited time to set up all the individual project dependencies.
- Wants a self-service experience that helps you get set up and to your code faster.

You'll also find Dev Environments useful if you are a developer that sets up your application's services for developer teams and wants to standardize the onboarding experience of new team members.

Dev Environments provides an easy way for developer teams of all skill levels to specify simple or complex requirements for their projects, and then quickly hit the ground running.

### How does it work?

Dev Environments is powered by [Docker Compose](../../compose/index.md). This allows Dev Envs to take advantage of all the benefits and features of Compose whilst adding an intuitive GUI where you can launch projects with the click of a button and have a centralized place for you to manage more than one project at a time.

Every dev environment you want to run needs a compose.YAML which configures your application's services and lives in your project directory. But you don't need to be an expert in Docker Compose or write a compose.YAML file from scratch - there are many [sample dev environments](https://github.com/docker/awesome-compose) that provide a useful starting point for how to integrate different services.

<div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
             <div class="component-icon">
                 <img src="/images/quickstart.svg" alt="Quickstart guide" width="45" height="45">
             </div>
                 <h2 id="docker-for-mac"><a href="/desktop/dev-environments/quickstart/">Quickstart guide </a></h2>
                <p>Step-by-step instructions on how to easily get started with Dev Environments </p>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <img src="/images/icon-machine@2X.png" alt="Run a dev environment" width="45" height="45">
            </div>
                <h2 id="docker-for-mac"><a href="/desktop/dev-environments/run-devenv/">Run a dev environment</a></h2>
                <p>Get up and running with a project or multiple projects quickly and easily.</p>
         </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <img src="/images/compose_48.svg" alt="build a dev environment" width="45" height="45">
            </div>
                <h2 id="docker-for-linux"><a href="/desktop/dev-environments/create-devenv/">Build a dev environment</a></h2>
                <p>Learn how to build a dev environment and onboard team members to your projects.</p>
        </div>
    </div>
</div>
</div>









