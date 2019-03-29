---
<<<<<<< HEAD
title: Docker Desktop Enterprise overview
=======
title: Docker Desktop Enterprise
>>>>>>> 1013: Move desktop ent content to docs-private
description: Learn about Docker Desktop Enterprise
keywords: Docker EE, Windows, Mac, Docker Desktop, Enterprise
---

<<<<<<< HEAD
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
=======
# Overview

Docker Desktop Enterprise provides local development, testing, and building of Docker applications on Mac or Windows. With work performed locally, developers can leverage a rapid feedback loop before pushing code or docker images to shared servers / continuous integration infrastructure.

Docker Desktop Enterprise takes Docker Desktop Community, formerly known as Docker for Windows and Docker for Mac, a step further with simplified enterprise application development and maintenance. With Docker Desktop Enterprise, IT organizations can ensure developers are working with the same version of Docker Desktop Enterprise and can easily distribute Docker Desktop Enterprise to large teams using a number of third-party endpoint management applications. With the Docker Desktop Enterprise graphical user interface (GUI), developers are no longer required to work with lower-level Docker commands and can auto-generate Docker artifacts.

Installed with a single click or command line command, Docker Desktop Enterprise is integrated with the host OS framework, networking, and filesystem. Docker Desktop Enterprise is also designed to integrate with existing development environments (IDEs) such as Visual Studio and IntelliJ. With support for defined application templates, Docker Desktop Enterprise allows organizations to specify the look and feel of their applications.

Feature comparison of Docker Desktop Community versus Docker Desktop Enterprise:

  | Feature                 | Community version | Docker Desktop Enterprise |
  | :-----------------------|:-----------------:|:-------------------------:|
  | Docker Engine           | X                 |  X                        |
  | Docker Compose          | X                 |  X                        |
  | CLI                     | X                 |  X                        |
  | Windows and Mac support | X                 |  X                        |
  | Version selection       |                   |  X                        |
  | Application Designer    |                   |  X                        |
  | Device management       |                   |  X                        |
  | Administrative control  |                   |  X                        |
>>>>>>> 1013: Move desktop ent content to docs-private

## Docker Desktop Enterprise features

The following section lists features that are exclusive to Docker Desktop Enterprise:

### Version Selection

Configurable version packs ensure the local instance of Docker Desktop Enterprise is a precise copy of the production environment where applications are deployed.

System administrators can install version packs using a built-in command line tool. Once installed, developers can switch between versions of Docker and Kubernetes with a single click and ensure Docker and Kubernetes versions match UCP cluster versions.

### Application Designer

<<<<<<< HEAD
 Application Designer provides a library of application and service templates to help developers quickly create new Docker applications.

### Application templates
=======
 Application Designer provides a library of application and service templates to help Docker developers quickly create new Docker applications.

### Application Templates
>>>>>>> 1013: Move desktop ent content to docs-private

Application templates allow you to choose a technology stack and focus on business logic and code, and require only minimal Docker syntax knowledge. Template support includes .NET, Spring, and more.

### Device management

<<<<<<< HEAD
The Docker Desktop Enterprise installer is available as standard MSI (Windows) and PKG (Mac) downloads, which allows administrators to script an installation across many developer workstations.

### Administrative control

IT organizations can specify and lock configuration parameters for the creation of  standardized development environments, including disabling drive sharing.
=======
The Docker Desktop Enterprise installer is available as standard MSI (Win) and PKG (Mac) downloads, which allows administrators to script an installation across many developer workstations.

### Administrative control

IT organizations can specify and lock configuration parameters for creation of  standardized development environment, including disabling drive sharing and limiting version pack installations.
>>>>>>> 1013: Move desktop ent content to docs-private

Developers can then run commands using the command line without worrying about configuration settings.