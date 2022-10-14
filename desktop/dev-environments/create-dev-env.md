---
description: Dev Environments
keywords: Dev Environments, share, collaborate, local
title: Create a simple dev environment
---

You can create a dev environment from a:
- Git repository
- Branch or tag of a Git repository
- Subfolder of a Git repository
- Local folder

This did not conflict with any of the local files or local tooling set up on your host. 

## Create a dev environment from a Git repository

The simplest way to get started with Dev Environments is to create a new environment by cloning the Git repository of the project you are working on. 

For example, create a new Dev Environment using the simple `single-dev-env` project from the [Docker Samples](https://github.com/dockersamples/single-dev-env){:target="_blank" rel="noopener" class="_"} GitHub repository.

> **Note**
>
> When cloning a Git repository using SSH, ensure you've added your SSH key to the ssh-agent. To do this, open a terminal and run `ssh-add <path to your private ssh key>`.

> **Important**
>
> If you have enabled the WSL 2 integration in Docker Desktop for Windows, make sure you have an SSH agent running in your WSL 2 distribution.
{: .important}

<div class="panel panel-default">
    <div class="panel-heading collapsed" data-toggle="collapse" data-target="#collapse-wsl2-ssh" style="cursor: pointer">
    How to start an SSH agent in WSL2
    <i class="chevron fa fa-fw"></i></div>
    <div class="collapse block" id="collapse-wsl2-ssh">
    If your WSL 2 distribution doesn't have an `ssh-agent` running, you can append this script at the end of your profile file (that is: ~/.profile, ~/.zshrc, ...).
<pre><code>
SSH_ENV="$HOME/.ssh/agent-environment"
function start_agent {
    echo "Initialising new SSH agent..."
    /usr/bin/ssh-agent | sed 's/^echo/#echo/' > "${SSH_ENV}"
    echo succeeded
    chmod 600 "${SSH_ENV}"
    . "${SSH_ENV}" > /dev/null
}
# Source SSH settings, if applicable
if [ -f "${SSH_ENV}" ]; then
    . "${SSH_ENV}" > /dev/null
    ps -ef | grep ${SSH_AGENT_PID} | grep ssh-agent$ > /dev/null || {
        start_agent;
    }
else
    start_agent;
fi
</code></pre>
    </div>
</div>

To create a dev environment:

1. From **Under Dev Environments** in Docker Dashboard, click **Create**. The **Create a Dev Environment** dialog displays.
2. Select **Get Started** and then copy `https://github.com/dockersamples/single-dev-env.git` and add it to the **Enter the Git Repository** field with **Existing Git repo** as the source.
3. Select **Continue**.

    This detects the main language of your repository, clones the Git code inside a volume, determines the best image for your Dev Environment, and opens VS Code inside the Dev Environment container.

4. Hover over the container and select **Open in VS Code** to start working. You can also open a terminal in VS Code, and use Git to push or pull code to your repository, or switch between branches and work as you would normally.

5. To launch the application, run the command `make run` in your terminal. This opens an http server on port 8080. Open [http://localhost:8080](http://localhost:8080) in your browser to see the running application.


## Create a dev environment from a specific branch or tag

You can create a dev environment from a specific branch (for example, a branch corresponding to a Pull Request) or a tag by adding `@mybranch` or `@tag` as a suffix to your Git URL:

 `https://github.com/dockersamples/single-dev-env@mybranch`

 or

 `git@github.com:dockersamples/single-dev-env.git@mybranch`

Docker then clones the repository with your specified branch or tag.

## Create a dev environment from a subdirectory of a Git repository

>Note
>
>Currently, Dev Environments is not able to detect the main language of the subdirectory. You need to define your own base image or services in a `compose-dev.yaml`file located in your subdirectory. For more information on how to configure, see the [React application with a Spring backend and a MySQL database sample](https://github.com/docker/awesome-compose/tree/master/react-java-mysql) or the [Go server with an Nginx proxy and a Postgres database sample](https://github.com/docker/awesome-compose/tree/master/nginx-golang-postgres). 

1. From **Dev Environments** in Docker Dashboard, click **Create**. The **Create a Dev Environment** dialog displays.
2. Select **Get Started** and then copy your Git subfolder link into the **Enter the Git Repository** field with **Existing Git repo** as the source.
3. Select **Continue**.

    This clones the Git code inside a volume, determines the best image for your Dev Environment, and opens VS Code inside the Dev Environment container.

4. Hover over the container and select **Open in VS Code** to start working. You can also open a terminal in VS Code, and use Git to push or pull code to your repository, or switch between branches and work as you would normally.

5. To launch the application, run the command `make run` in your terminal. This opens an http server on port 8080. Open [http://localhost:8080](http://localhost:8080) in your browser to see the running application.

## Create a dev environment from a local folder

1. From **Dev Environments** in Docker Dashboard, click **Create**. The **Create a Dev Environment** dialog displays.
2. Select **Get Started** and then choose **Local Folder** as the source.
3. Select **Select directory** to open the root of the code that you would like to work on.
4. Select **Continue**.

    This detects the main language of your local folder, creates a Dev Environment using your local folder, and bind-mounts your local code in the Dev Environment. It then opens VS Code inside the Dev Environment container.

> **Note**
>
> When using a local folder for a dev environment, file changes are synchronized between your environment container and your local files. This can affect the performance inside the container, depending on the number of files in your local folder and the operations performed in the container.

## What's next?

Learn how to [distribute your dev environment](share.md)
