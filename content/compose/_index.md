---
description: Learn how to use Docker Compose to define and run multi-container applications
  with this detailed introduction to the tool.
keywords: docker compose, docker-compose, docker compose command, docker compose files,
  docker compose documentation, using docker compose, compose container, docker compose
  service
title: Docker Compose overview
grid:
- title: Install Compose
  description: Follow the instructions on how to install Docker Compose.
  icon: download
  link: /compose/install
- title: Try Compose
  description: Learn the key concepts of Docker Compose whilst building a simple Python
    web application.
  icon: explore
  link: /compose/gettingstarted
- title: View the release notes
  description: Find out about the latest enhancements and bug fixes.
  icon: note_add
  link: /compose/release-notes
- title: Understand key features of Compose
  description: Understand its key features and explore common use cases.
  icon: category
  link: /compose/features-uses/
- title: Explore the Compose file reference
  description: Find information on defining services, networks, and volumes for a
    Docker application.
  icon: feature_search
  link: /compose/compose-file
- title: Browse common FAQs
  description: Explore general FAQs and find out how to give feedback.
  icon: help
  link: /compose/faq
aliases:
- /compose/cli-command/
- /compose/networking/swarm/
- /compose/overview/
- /compose/swarm/
- /compose/completion/
---

{{< include "compose-eol.md" >}}

Compose is a tool for defining and running multi-container Docker applications.
With Compose, you use a YAML file to configure your application's services.
Then, with a single command, you create and start all the services
from your configuration.

Compose works in all environments; production, staging, development, testing, as
well as CI workflows. It also has commands for managing the whole lifecycle of your application:

 * Start, stop, and rebuild services
 * View the status of running services
 * Stream the log output of running services
 * Run a one-off command on a service

The key features of Compose that make it effective are:

* [Have multiple isolated environments on a single host](features-uses.md#have-multiple-isolated-environments-on-a-single-host)
* [Preserve volume data when containers are created](features-uses.md#preserves-volume-data-when-containers-are-created)
* [Only recreate containers that have changed](features-uses.md#only-recreate-containers-that-have-changed)
* [Support variables and moving a composition between environments](features-uses.md#supports-variables-and-moving-a-composition-between-environments)

{{< grid >}}