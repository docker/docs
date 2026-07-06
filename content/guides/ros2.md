---
title: Introduction to ROS 2 Development with Docker
linkTitle: ROS 2
description: Learn how to containerize and develop ROS 2 applications using Docker.
keywords: ros2, robotics, devcontainers, python, cpp, Dockerfile, rviz
summary: |
  This guide details how to containerize ROS 2 applications using Docker.
aliases:
  - /guides/ros2/develop/
  - /guides/ros2/run-ros2/
  - /guides/ros2/turtlesim-example/
params:
  tags: [deployment]
  time: 30 minutes
  image: /images/guides/ros2.jpg
  featured: false
---


> **Acknowledgment**
>
> This guide is a community contribution. Docker would like to thank
> [Shakirth Anisha](https://www.linkedin.com/in/shakirth-anisha/) for her contribution
> to this guide.

[ROS 2](https://www.ros.org/) is a set of software libraries and tools for building robot applications. It uses Data Distribution Service (DDS) for real-time, secure communication between distributed nodes, making it ideal for robotics and autonomous systems.

---

## What will you learn?

In this guide, you'll learn how to:

- Use official ROS 2 base images from Docker Hub
- Run ROS 2 in an Ubuntu container
- Install ROS 2 packages and dependencies
- Set up a development container for local development
- Run a complete end-to-end example with Turtlesim

## Prerequisites

Before you begin, make sure you're familiar with the following:

- [Docker Desktop](https://docs.docker.com/desktop/): You must have Docker Desktop installed and running.
- [Docker concepts](/get-started/docker-concepts/the-basics/what-is-a-container.md): You must understand core Docker concepts, such as images and containers.
- [ROS 2 concepts](https://www.ros.org): Basic understanding of concepts like nodes, packages, topics, and services.

## What's next?

Start by setting up your ROS 2 development environment using Docker and dev containers.

## Run ROS 2 in a container

### Overview

In this section, you will run ROS 2 in an isolated Docker container using official ROS 2 images, verify that ROS 2 is working, and install additional ROS 2 packages for development and testing.

---

### Run ROS 2 in a container

The fastest way to get started with ROS 2 is to use the [official Docker image](https://hub.docker.com/_/ros/). To pull an image, start a container, and open an interactive bash shell:

1. Pull and run the official ROS 2 Docker image:

   ```console
   $ docker run -it ros:humble
   ```

   This guide uses the Humble distribution. You can replace `humble` with another supported distribution such as `rolling`, `jazzy`, or `iron`.

   > [!NOTE]
   >
   > This environment is temporary and does not maintain persistence.
   > Any files you create or packages you install will be deleted once the container is stopped or removed.

2. Verify ROS 2 is working:

   ```console
   $ echo $ROS_DISTRO
   ```

   You should see output similar to:

   ```text
   humble
   ```

### Install ROS 2 packages

The official ROS 2 images include core packages. To install additional packages, use the `apt` package manager:

1. Update the package manager:

   ```console
   $ sudo apt update
   ```

2. Install the desired package:

   ```console
   $ sudo apt install $PACKAGE_NAME
   ```

Replace `$PACKAGE_NAME` with any package you want to install.

Some commonly used packages include:

- `ros-humble-turtlesim` - Visualization and simulation tool
- `ros-humble-rviz2` - 3D visualization tool
- `ros-humble-rqt` - Qt-based ROS graphical tools
- `ros-humble-demo-nodes-cpp` - C++ demo nodes
- `ros-humble-demo-nodes-py` - Python demo nodes
- `ros-humble-colcon-common-extensions` - Build system extensions

### Summary

In this section, you pulled an official ROS 2 Docker image, launched an interactive session, and extended the container's capabilities by installing additional ROS 2 packages using apt.

### Next steps

In the next section, you will configure a persistent workspace to ensure your code and modifications are saved across sessions.

## Build and develop a ROS 2 workspace

### Overview

In this section, you will set up a ROS 2 workspace using Docker and development containers, review the workspace layout, open the workspace in Visual Studio Code, and edit and build ROS 2 projects inside the container.

---

### Get the sample ROS 2 workspace

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

### Open and build the container

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

### Switch ROS 2 distributions

Update the base image in your `Dockerfile`, changing from `humble` to another distribution like `rolling`, `jazzy`, or `iron`.

### Summary

In this section, you learned how to create a structured workspace, write a Dockerfile with development tools, and configure a Docker Compose setup. Your ROS 2 development environment is now ready with a consistent, reproducible setup across any machine.

### Next steps

In the next section, you'll run a complete end-to-end example with Turtlesim.

## Run a complete example with Turtlesim

### Overview

Turtlesim is a simple simulation tool that demonstrates fundamental ROS 2 concepts such as nodes, topics, and services. In this section, you'll run a complete example with Turtlesim, control the turtle, monitor topics, and visualize the system with rqt.

---

### Configure display forwarding

#### Linux

Allow Docker access to your X server:

```console
$ xhost +local:docker
```

#### macOS

On macOS, use XQuartz to provide X11 support. Install XQuartz using Homebrew:

1. Install XQuartz using Homebrew:

   ```console
   $ brew install --cask xquartz
   ```

2. Open XQuartz from Applications, then navigate to `Preferences > Security` and enable `Allow connections from network clients`. Restart your computer to ensure the changes take effect.

3. After rebooting, open a terminal and allow local connections:

   ```console
   $ defaults write org.xquartz.X11 nolisten_tcp -bool false
   $ xhost +localhost
   $ xhost + 127.0.0.1
   ```

### Start the container

Start the container using the same Docker Compose setup from the workspace section.

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

### Install and run Turtlesim

Inside the container, install the Turtlesim package:

1. Update the package manager:

   ```console
   $ sudo apt update
   ```

2. Install the Turtlesim package:

   ```console
   $ sudo apt install -y ros-humble-turtlesim
   ```

3. Run the Turtlesim node:

   ```console
   $ ros2 run turtlesim turtlesim_node
   ```

A window should appear on your desktop showing a turtle in a grid.

### Control the turtle

1. Open a new terminal and connect to the same container, then start the keyboard teleop node:

   ```console
   $ ros2 run turtlesim turtle_teleop_key
   ```

   This node allows you to control the turtle using your keyboard. Use the arrow keys to move the turtle forward, backward, left, and right. Press `Ctrl+C` to stop the teleop node.

2. Move the turtle around the window. You should see it draw a path as it moves.

### Monitor topics

1. Open another terminal and connect to the same container, then list all active topics:

   ```console
   $ ros2 topic list
   ```

   You should see output similar to the following:

   ```text
   /parameter_events
   /rosout
   /turtle1/cmd_vel
   /turtle1/color_sensor
   /turtle1/pose
   ```

2. Get information about a specific topic:

   ```console
   $ ros2 topic info /turtle1/pose
   ```

   You'll see the topic type and which nodes publish and subscribe to it.

### Visualize the system with rqt

1. Open another terminal and connect to the same container, then update the package manager:

   ```console
   $ sudo apt update
   ```

2. Install rqt:

   ```console
   $ sudo apt install -y 'ros-humble-rqt*'
   ```

3. Start rqt:

   ```console
   $ ros2 run rqt_gui rqt_gui
   ```

An rqt window should appear. rqt provides several useful plugins for visualizing and monitoring ROS 2 systems.

#### Node Graph

You can explore the node graph by navigating to **Plugins > Introspection > Node Graph**. A new tab opens showing nodes and topics with connections illustrated as lines. This visualization demonstrates how the teleop node sends velocity commands to the Turtlesim node, and how the Turtlesim node publishes position data back through topics.

#### Topic Monitor

You can monitor active topics by navigating to **Plugins > Topics > Topic Monitor**. A new tab opens displaying all active topics and their current values. Select the eye icon next to `/turtle1/pose` to monitor it. As you move the turtle, watch the pose values update in real time, showing the position of the turtle and orientation changing based on your commands.

#### Service Caller

You can call services from rqt using **Plugins > Services > Service Caller**. Select a service such as `/turtle1/teleport_absolute`, enter values for the request fields, and select **Call** to send the request.

#### Plots

To plot topic data over time navigate to **Plugins > Visualization > Plot**. For example, in the Plot window, type `/turtle1/pose/x` in the Topic field and press Enter. Move the turtle and watch the X position displayed as a graph over time.

### Call ROS 2 services

Turtlesim provides services for actions such as repositioning the turtle and clearing the path.

1. List available services:

   ```console
   $ ros2 service list
   ```

   You should see services such as `/turtle1/set_pen` (to change pen color and width), `/turtle1/teleport_absolute` (to move the turtle to a specific position), and `/turtle1/teleport_relative` (to move the turtle relative to its current position).

2. Teleport the turtle to a new position:

   ```console
   $ ros2 service call /turtle1/teleport_absolute turtlesim/srv/TeleportAbsolute "
   x: 1.0
   y: 3.0
   theta: 0.0
   "
   ```

   The turtle should instantly move to the specified position (1.0, 3.0).

### Create a simple publisher

1. Create a Python script that publishes velocity commands to control the turtle programmatically. In a new terminal, create a file called `move_turtle.py`:

   ```python
   import rclpy
   from geometry_msgs.msg import Twist
   import time

   def main():
       rclpy.init()
       node = rclpy.create_node('turtle_mover')
       publisher = node.create_publisher(Twist, 'turtle1/cmd_vel', 10)

       # Create a twist message
       msg = Twist()
       msg.linear.x = 2.0  # Move forward at 2 m/s
       msg.angular.z = 1.0  # Rotate at 1 rad/s

       # Publish the message
       for i in range(50):
           publisher.publish(msg)
           time.sleep(0.1)

       # Stop the turtle
       msg.linear.x = 0.0
       msg.angular.z = 0.0
       publisher.publish(msg)

       node.destroy_node()
       rclpy.shutdown()

   if __name__ == '__main__':
       main()
   ```

2. Run the script:

   ```console
   $ python3 move_turtle.py
   ```

   The turtle should move in a circular motion for 5 seconds and then stop.

### Summary

In this section, you configured display forwarding, used the Turtlesim nodes, inspected nodes and topics, and visualized the system using rqt. Finally, you interacted with ROS 2 services and created a simple publisher to move the turtle programmatically.

These fundamental concepts apply directly to real-world robotics applications with actual sensors and actuators.

### Related resources

- [ROS 2 Turtlesim tutorials](https://docs.ros.org/en/humble/Tutorials/Beginner-CLI-Tools/Understanding-ROS2-Topics/Understanding-ROS2-Topics.html)
- [ROS 2 Concepts](https://docs.ros.org/en/humble/Concepts.html)
- [Geometry Messages](https://github.com/ros2/geometry2/tree/humble/geometry_msgs)
