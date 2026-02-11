---
title: Run a complete example with Turtlesim
linkTitle: Turtlesim example
weight: 20
keywords: ros2, turtlesim, example, nodes, topics, teleop
description: Run a complete end-to-end ROS 2 example with Turtlesim.
---

## Overview

Turtlesim is a simple simulation tool that demonstrates fundamental ROS 2 concepts such as nodes, topics, and services. In this section, you'll run a complete example with Turtlesim, control the turtle, monitor topics, and visualize the system with rqt.

---

## Configure display forwarding 

### Linux

Allow Docker access to your X server:

```console
$ xhost +local:docker
```

### macOS

On macOS, use XQuartz to provide X11 support. Install XQuartz using Homebrew:

1. Install XQuartz using Homebrew:

    ```console
    $ brew install --cask xquartz
    ```

2. Open XQuartz from Applications, then navigate to `Preferences > Security` and enable `Allow connections from network clients`. Restart your computer to ensure the changes take effect.

3. After rebooting, open a terminal and allow local connections:

    ```console
    $ xhost +localhost
    ```

> [!Note]
>
> Some ROS 2 visualization tools, such as RViz may be unavailable when using macOS.

## Start the dev container

Start the container using the same dev container setup from the workspace section.

For Linux:
```console
$ devcontainer up --workspace-folder ws_linux/src
$ devcontainer exec --workspace-folder ws_linux/src /bin/bash
```

For macOS:
```console
$ devcontainer up --workspace-folder ws_mac/src
$ devcontainer exec --workspace-folder ws_mac/src /bin/bash
```

## Install and Run Turtlesim

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

## Control the turtle

1. Open a new terminal and connect to the same container, then start the keyboard teleop node:

    ```console
    $ ros2 run turtlesim turtle_teleop_key
    ```

    This node allows you to control the turtle using your keyboard. Use the arrow keys to move the turtle forward, backward, left, and right. Press `Ctrl+C` to stop the teleop node.

2. Move the turtle around the window. You should see it draw a path as it moves.

## Monitor topics

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

## Visualize the system with rqt

1. Open another terminal and connect to the same container, then update the package manager:

    ```console
    $ sudo apt update
    ```

2. Install rqt:

    ```console
    $ sudo apt install 'ros-humble-rqt*'
    ```

3. Start rqt:

    ```console
    $ ros2 run rqt_gui rqt_gui
    ```

An rqt window should appear. rqt provides several useful plugins for visualizing and monitoring ROS 2 systems.

### Node Graph

You can explore the node graph by navigating to **Plugins > Introspection > Node Graph**. A new tab opens showing nodes and topics with connections illustrated as lines. This visualization demonstrates how the teleop node sends velocity commands to the Turtlesim node, and how the Turtlesim node publishes position data back through topics.

### Topic Monitor

You can monitor active topics by navigating to **Plugins > Topics > Topic Monitor**. A new tab opens displaying all active topics and their current values. Select the eye icon next to `/turtle1/pose` to monitor it. As you move the turtle, watch the pose values update in real time, showing the position of the turtle and orientation changing based on your commands.

### Message Publisher

You can also publish messages manually using **Plugins > Topics > Message Publisher**. The Message Publisher plugin allows you to manually publish messages to any topic. You can use it to send velocity commands directly without writing a script.

### Plots

To plot topic data over time navigate to **Plugins > Visualization > Plot**. For example, in the Plot window, type `/turtle1/pose/x` in the Topic field and press Enter. Move the turtle and watch the X position displayed as a graph over time. 

## Call ROS 2 services

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

    The turtle should instantly move to the specified position (3.0, 4.0).

## Create a simple publisher

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

## Summary

In this section, you configured display forwarding, used the Turtlesim nodes, inspected nodes and topics, and visualized the system using rqt. Finally, you interacted with ROS 2 services and created a simple publisher to move the turtle programmatically.

These fundamental concepts apply directly to real-world robotics applications with actual sensors and actuators.

## Related resources

- [ROS 2 Turtlesim tutorials](https://docs.ros.org/en/humble/Tutorials/Beginner-CLI-Tools/Understanding-ROS2-Topics/Understanding-ROS2-Topics.html)
- [ROS 2 Concepts](https://docs.ros.org/en/humble/Concepts.html)
- [Geometry Messages](https://github.com/ros2/geometry2/tree/humble/geometry_msgs)
