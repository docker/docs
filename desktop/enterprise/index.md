---
title: Docker Desktop Enterprise overview
description: Learn about Docker Desktop Enterprise
keywords: Docker EE, Windows, Mac, Docker Desktop, Enterprise
redirect_from:
- /ee/desktop/
---

Welcome to Docker Desktop Enterprise. This page contains information about the Docker Desktop Enterprise (DDE) release. For information about Docker Desktop Community, see:

- [Docker Desktop for Mac (Community)](/docker-for-mac/){: target="_blank" class="_"}

- [Docker Desktop for Windows (Community)](/docker-for-windows/){: target="_blank" class="_"}

Docker Desktop Enterprise provides local development, testing, and building of Docker applications on Mac and Windows. With work performed locally, developers can leverage a rapid feedback loop before pushing code or Docker images to shared servers / continuous integration infrastructure.

Docker Desktop Enterprise takes Docker Desktop Community, formerly known as Docker for Windows and Docker for Mac, a step further with simplified enterprise application development and maintenance. With DDE, IT organizations can ensure developers are working with the same version of Docker Desktop and can easily distribute Docker Desktop to large teams using third-party endpoint management applications. With the Docker Desktop graphical user interface (GUI), developers do not have to work with lower-level Docker commands and can auto-generate Docker artifacts.

Installed with a single click or command line command, Docker Desktop Enterprise is integrated with the host OS framework, networking, and filesystem. DDE is also designed to integrate with existing development environments (IDEs) such as Visual Studio and IntelliJ. With support for defined application templates, Docker Desktop Enterprise allows organizations to specify the look and feel of their applications.

Feature comparison of Docker Desktop Community versus Docker Desktop Enterprise:

  | Feature                     | Docker Desktop (Community) | Docker Desktop Enterprise |
  | :-------------------------  |:--------------------------:|:-------------------------:|
  | Docker Engine               | X                          |  X                        |
  | Certified Kubernetes        | X                          |  X                        |
  | Docker Compose              | X                          |  X                        |
  | CLI                         | X                          |  X                        |
  | Windows and Mac support     | X                          |  X                        |
  | Version Selection           |                            |  X                        |
  | Application Designer        |                            |  X                        |
  | Custom application templates|                            |  X                        |
  | Docker Assemble             |                            |  X                        |
  | Device management           |                            |  X                        |
  | Administrative control      |                            |  X                        |

## Docker Desktop Enterprise features

The following section lists features that are exclusive to Docker Desktop Enterprise:

### Version Selection

Configurable version packs ensure the local instance of Docker Desktop Enterprise is a precise copy of the production environment where applications are deployed.

System administrators can install version packs using a built-in command line tool. Once installed, developers can switch between versions of Docker and Kubernetes with a single click and ensure Docker and Kubernetes versions match UCP cluster versions.

### Application Designer

 Application Designer provides a library of application and service templates to help developers quickly create new Docker applications.

### Application templates

Application templates allow you to choose a technology stack and focus on business logic and code, and require only minimal Docker syntax knowledge. Template support includes .NET, Spring, and more.

### Device management

The Docker Desktop Enterprise installer is available as standard MSI (Windows) and PKG (Mac) downloads, which allows administrators to script an installation across many developer workstations.

### Administrative control

IT organizations can specify and lock configuration parameters for the creation of  standardized development environments, including disabling drive sharing.

Developers can then run commands using the command line without worrying about configuration settings.
