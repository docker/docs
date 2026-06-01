---
title: Docker Compose
weight: 30
description: Learn how to use Docker Compose to define and run multi-container applications
  with this detailed introduction to the tool.
keywords: docker compose, docker-compose, compose.yaml, docker compose command, multi-container applications, container orchestration, docker cli
params:
  sidebar:
    group: Application development
grid:
- title: Why use Compose?
  description: Understand Docker Compose's key benefits
  icon: magnifying-glass
  link: /compose/intro/features-uses/
- title: How Compose works 
  description: Understand how Compose works
  icon: squares-2x2
  link: /compose/intro/compose-application-model/
- title: Install Compose
  description: Follow the instructions on how to install Docker Compose.
  icon: arrow-down-tray
  link: /compose/install
- title: Quickstart
  description: Learn the key concepts of Docker Compose whilst building a simple Python
    web application.
  icon: magnifying-glass-plus
  link: /compose/gettingstarted
- title: View the release notes
  description: Find out about the latest enhancements and bug fixes.
  icon: document-plus
  link: "https://github.com/docker/compose/releases"
- title: Explore the Compose file reference
  description: Find information on defining services, networks, and volumes for a
    Docker application.
  icon: arrows-right-left
  link: /reference/compose-file
- title: Use Compose Bridge
  description: Transform your Compose configuration file into configuration files for different platforms, such as Kubernetes.
  icon: arrow-down
  link: /compose/bridge
- title: Browse common FAQs
  description: Explore general FAQs and find out how to give feedback.
  icon: question-mark-circle
  link: /compose/faq
aliases:
- /compose/overview/
- /compose/swarm/
- /compose/releases/migrate/
---

Docker Compose is a tool for defining and running multi-container applications. 
It is the key to unlocking a streamlined and efficient development and deployment experience. 

Compose simplifies the control of your entire application stack, making it easy to manage services, networks, and volumes in a single YAML configuration file. Then, with a single command, you create and start all the services
from your configuration file.

Compose works in all environments - production, staging, development, testing, as
well as CI workflows. It also has commands for managing the whole lifecycle of your application:

 - Start, stop, and rebuild services
 - View the status of running services
 - Stream the log output of running services
 - Run a one-off command on a service

{{< grid >}}
