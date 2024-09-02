---
description: Learn how to use Docker Compose to define and run multi-container applications
  with this detailed introduction to the tool.
keywords: docker compose, docker-compose, docker compose command, docker compose files,
  docker compose documentation, using docker compose, compose container, docker compose
  service
title: Docker Compose overview
grid:
- title: Why use Compose?
  description: Understand Docker Compose's key benefits
  icon: feature_search
  link: /compose/intro/features-uses/
- title: How Compose works 
  description: Understand how Compose works
  icon: category
  link: /compose/compose-application-model/
- title: Install Compose
  description: Follow the instructions on how to install Docker Compose.
  icon: download
  link: /compose/install
- title: Quickstart
  description: Learn the key concepts of Docker Compose whilst building a simple Python
    web application.
  icon: explore
  link: /compose/gettingstarted
- title: View the release notes
  description: Find out about the latest enhancements and bug fixes.
  icon: note_add
  link: /compose/release-notes
- title: Explore the Compose file reference
  description: Find information on defining services, networks, and volumes for a
    Docker application.
  icon: polyline
  link: /reference/compose-file
- title: Browse common FAQs
  description: Explore general FAQs and find out how to give feedback.
  icon: help
  link: /compose/faq
- title: Migrate to Compose V2
  description: Learn how to migrate from Compose V1 to V2
  icon: folder_delete
  link: /compose/migrate/
aliases:
- /compose/cli-command/
- /compose/networking/swarm/
- /compose/overview/
- /compose/swarm/
- /compose/completion/
---

Docker Compose is a tool for defining and running multi-container applications. 
It is the key to unlocking a streamlined and efficient development and deployment experience. 

Compose simplifies the control of your entire application stack, making it easy to manage services, networks, and volumes in a single, comprehensible YAML configuration file. Then, with a single command, you create and start all the services
from your configuration file.

Compose works in all environments; production, staging, development, testing, as
well as CI workflows. It also has commands for managing the whole lifecycle of your application:

 * Start, stop, and rebuild services
 * View the status of running services
 * Stream the log output of running services
 * Run a one-off command on a service

{{< grid >}}
