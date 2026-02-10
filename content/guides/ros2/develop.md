---
title: Build and develop a ROS 2 workspace
linkTitle: Set Up ROS 2 workspace
weight: 15
keywords: ros2, robotics, docker, dockerfile, devcontainer, vscode, workspace
description: Learn how to develop ROS 2 applications using a Docker based workspace and development containers.

---

## Overview

In this section, you will set up a ROS 2 workspace using Docker and development containers, review the workspace layout, open the workspace in Visual Studio Code, and edit and build ROS 2 projects inside the container.

---

## Get the sample ROS 2 workspace

A consistent workspace simplifies managing ROS 2 projects and build artifacts across different distributions.

1. Open a terminal and clone the sample workspace repository:

    ```console
    $ git clone https://github.com/shakirth-anisha/docker-ros2-workspace.git
    $ cd docker-ros2-workspace

    ```

    Moving forward, Linux users can use the `ws_linux` folder, and macOS users can use `ws_mac`.

2. Verify the workspace structure:

    ```text
    ws/
    ├── cache/
    |   ├── [ROS2_DISTRO]/
    |   |   ├── build/
    |   |   ├── install/
    |   |   └── log/
    |   └── ...
    |
    ├── src/
        ├── .devcontainer/
        │   ├── devcontainer.json
        │   └── Dockerfile
        ├── package1/
        └── package2/

    ```

3. Explore the workspace layout
- `cache` : Stores build, install, and log artifacts for each ROS 2 distribution. Keeping these directories outside the source tree lets you change distributions without rebuilding everything from scratch.
- `src` : Contains all ROS 2 packages and the development container configuration. This directory is mounted into the container as the active workspace.
- `.devcontainer/` : Holds the configuration used by Visual Studio Code to open the workspace in a container.
  - `Dockerfile` : Builds the ROS 2 development image. It uses an official ROS 2 base image, creates a non-root development user, and installs required system and ROS 2 dependencies.
  - `devcontainer.json` : Defines how Visual Studio Code builds, runs, and connects to the container. It configures the workspace location, user context, mounted directories, environment variables, networking options, installed extensions, and post-create setup commands.

## Open and build dev container

1. Install the `devcontainer` CLI:

    - Linux / Windows
    ```console
    $ sudo apt update && sudo apt upgrade
    $ sudo apt install -y nodejs npm
    $ sudo npm install -g @devcontainers/cli
    ```
    - macOS
    ```console
    $ brew install devcontainer
    ```

2. Execute the following commands to build and start the container:

    ```console
    $ devcontainer up --workspace-folder ws/src

    ```

    This command builds the Docker image defined in your `.devcontainer` folder and starts the container in the background.

    > [!NOTE]
    >
    > Building the image may take several minutes during the first run 
    > as the CLI pulls the base ROS 2 image and installs required dependencies. 
    > Subsequent starts will be significantly faster.

3. Once the container is running, execute commands inside it using `exec`:

    ```console
    devcontainer exec --workspace-folder ws/src /bin/bash
    ```

4. Inside the container terminal, verify the environment:

```console
$ ros2 --version
$ which colcon
```

All commands should execute successfully inside the container.

## Switch ROS 2 distributions

1. Create a new cache directory for the target distribution:

    ```console
    $ mkdir -p ~/ws/cache/[ROS_DISTRIBUTION]/{build,install,log}
    ```

2. Update the cache paths in `src/.devcontainer/Dockerfile` and `src/.devcontainer/devcontainer.json`, changing from `humble` to `ROS_DISTRIBUTION`.

## Summary

In this section, you learned how to create a structured workspace with cache directories, write a Dockerfile with development tools and configure a dev container with `devcontainer.json`. Your ROS 2 development environment is now ready with a consistent, reproducible setup across any machine.

## Next steps

In the next section, you'll run a complete end-to-end example with Turtlesim.
