---
title: Open a new project
description: Learn how to open a new local or remote project, or an existing project in Docker Projects. 
keywords: Docker, projects, docker deskotp, containerization, open, remote, local
weight: 20
---

> [!NOTE]
> 
> In order to use Docker Projects, make sure you have enabled the **Access experimental features** setting on the **Features in development** tab. 

## New projects

In order to run a new project, it must be stored locally. In the **Projects** view in Docker Desktop, local projects display the local path under the project.

### Open a new local project

A project consists of your code and at least one Compose file. Ensure that you have a Compose file before trying to open a new project.

To open a new project:

1. Sign in to Docker Desktop, and go to **Projects**.

2. Select **Open a local folder**. This lets you select a local folder that contains your project’s code and a Compose file.

   > [!NOTE]
   >
   > A local folder can also be the folder of a Git repository that you have already cloned. 

3. Configure your project by giving it a name and setting the owner, then select **Next**.

   > [!NOTE]
   >
   > If you are part of a Docker organization you have the option to [share your project](share.md) with the organization. 

4. Specify how to run your project by selecting **New run command**:

   > [!TIP]
   >
   > While configuring your run command, you can view the equivalent `docker compose up` command in the **Run command** section on the configuration page. You can also use this command to run your project from the command line. You can refer to the [`docker compose up` reference documentation](/reference/cli/docker/compose/up.md) to learn more about the options you configure. 

   - **Name**: Specify a name to identify the run command.
   - **Compose files**: Select one or more Compose files from your project. 
   - **Flags**: Optionally, select one or more flags for your run command.

   > [!TIP]
   > 
   > While the `--env-file` flag isn't currently supported, you can specify environment variables in your Compose file, or use the **Tasks** option to run a script that sets your environment variables. 

   - **Services that will run**: After selecting one or more Compose files, the services defined in the files will appear here. If there is more than one service, you can optionally choose to not run a service by deselecting the checkbox.
   - **Tasks (Advanced options)**: Optionally specify a command to run before running the project. For example, if you want to run a bash script from the project directory named `set-vars.sh`, you can specify bash `set-vars.sh`. Or, on Windows to run a script with `cmd.exe` named `set-vars.bat`, specify `set-vars.bat`. Note that a task can access environment variables from your terminal profile, but it can't access local shell functions nor aliases.

5. Select **Save changes**.

Your project is now ready to run. 

## Open a new remote project

The following steps prompt you to clone the Git repository for your project. 

If you have already cloned the repository outside of Docker Projects, then you can open the project as a new project and Docker Projects will automatically detect and link the repository.

To clone and open a remote project:

1. Sign in to Docker Desktop, and go to **Projects**.

2. Select **Clone a git repository**. This lets you specify a Git repository and a local folder to clone that repository to. The repository must contain at least your project’s code and a Compose file.

3. Enter the remote source and choose the local destination to clone to. 

4. Select **Clone project**.

5. Configure your project by giving it a name and setting the owner, then select **Next**.

   > [!NOTE]
   >
   > If you are part of a Docker organization you have the option to [share your project](share.md) with the organization. 

6. Specify how to run your project by selecting **New run command**:

   > [!TIP]
   >
   > While configuring your run command, you can view the equivalent `docker compose up` command in the **Run command** section on the configuration page. You can also use this command to run your project from the command line. You can refer to the [`docker compose up` reference documentation](/reference/cli/docker/compose/up.md) to learn more about the options you configure. 

   - **Name**: Specify a name to identify the run command.
   - **Compose files**: Select one or more Compose files from your project. 
   - **Flags**: Optionally, select one or more flags for your run command.

   > [!TIP]
   > 
   > While the `--env-file` flag isn't currently supported, you can specify environment variables in your Compose file, or use the **Tasks** option to run a script that sets your environment variables. 

   - **Services that will run**: After selecting one or more Compose files, the services defined in the files will appear here. If there is more than one service, you can optionally choose to not run a service by deselecing the checkbox.
   - **Tasks (Advanced options)**: Optionally specify a command to run before running the project. For example, if you want to run a bash script from the project directory named `set-vars.sh`, you can specify bash `set-vars.sh`. Or, on Windows to run a script with `cmd.exe` named `set-vars.bat`, specify `set-vars.bat`. Note that a task can access environment variables from your terminal profile, but it can't access local shell functions nor aliases.

7. Select **Save changes**.

## Existing projects

### Open an existing local project

1. Sign in to Docker Desktop, and go to **Projects**.

2. Open your project by selecting your project under **Recents**, or by selecting the specific owner that your project is associated with and then select your project. 

### Open an existing remote project

In the **Projects** view in Docker Desktop, existing remote projects display **No local copy** under the project. 

You’ll see remote projects when you are new to the team and are accessing a shared project, remove a project from Docker Desktop, or access Docker Desktop from a new device after creating a project associated with a Git repository.

To open an existing remote project, you can choose between:

   - Cloning the project into a local destination. 
   - Linking to an existing folder where the project has already been cloned

## What's next?

 - [View your project](/manuals/projects/view.md)
 - [Add or edit your run commands](/manuals/projects/edit.md)
 - [Manage your projects](/manuals/projects/manage.md)
