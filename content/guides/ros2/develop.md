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
    ws_linux/
    ├── compose.yml
    ├── Dockerfile
    └── src/
        ├── package1/
        └── package2/

    ws_mac/
    ├── compose.yml
    ├── Dockerfile
    └── src/
        ├── package1/
        └── package2/

    ```

3. Explore the workspace layout

- `compose.yml` : Defines how Docker Compose builds and runs the ROS 2 container, including mounts, environment variables, and networking settings.
- `Dockerfile` : Builds the ROS 2 development image. It uses an official ROS 2 base image, creates a non-root development user, and installs required system and ROS 2 dependencies.
- `src` : Contains all ROS 2 packages. This directory is mounted into the container as the active workspace.

## Open and build the container

1. Execute the following commands to build and start the container:

    For Linux:

    ```console
    $ cd ws_linux
    $ docker compose up -d
    $ docker compose exec ros2 /bin/bash
    ```

    For macOS:

    ```console
    $ cd ws_mac
    $ docker compose up -d
    $ docker compose exec ros2 /bin/bash
    ```

    This command builds the Docker image defined in your `Dockerfile` and starts the container in the background.

    > [!NOTE]
    >
    > Building the image may take several minutes during the first run 
    > as the CLI pulls the base ROS 2 image and installs required dependencies. 
    > Subsequent starts will be significantly faster.

2. Once the container is running, execute commands inside it using `exec`:

    ```console
    $ docker compose exec ros2 /bin/bash
    ```

3. Inside the container terminal, verify the environment:

```console
$ echo $ROS_VERSION
$ which colcon
```

All commands should execute successfully inside the container.

## Switch ROS 2 distributions

Update the base image in your `Dockerfile`, changing from `humble` to another distribution like `rolling`, `jazzy`, or `iron`.

## Summary

In this section, you learned how to create a structured workspace, write a Dockerfile with development tools, and configure a Docker Compose setup. Your ROS 2 development environment is now ready with a consistent, reproducible setup across any machine.

## Next steps

In the next section, you'll run a complete end-to-end example with Turtlesim.
