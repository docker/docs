---
title: Develop and run Laravel applications with Docker Compose
linkTitle: Laravel applications with Docker Compose
summary: Learn how to efficiently set up Laravel development and production environments using Docker Compose.
description: A guide on using Docker Compose to manage Laravel applications for development and production, covering container configurations and service management.
languages: [php]
levels: [intermediate]
products: [compose]
aliases:
  - /frameworks/laravel/
params:
  time: 30 minutes
  resource_links:
    - title: Laravel
      url: https://laravel.com/
      external: true
    - title: Docker Compose
      url: /compose/
    - title: Use Compose in production
      url: /compose/how-tos/production/
    - title: Repository with examples
      url: https://github.com/rw4lll/laravel-docker-examples
      external: true
---

Laravel is a popular PHP framework that allows developers to build web applications quickly and effectively. Docker Compose helps manage development and production environments by defining all necessary services, like PHP, web server and database, in a single YAML file. This guide provides an overview of setting up a robust environment for Laravel using Docker Compose, with a focus on simplicity and efficiency.

The demonstrated examples can be found in the [GitHub repository](https://github.com/rw4lll/laravel-docker-examples). Docker Compose is used as a straightforward way to connect multiple containers to meet Laravel's requirements, but similar techniques can also be applied using Docker Swarm, Kubernetes, or even individual Docker containers.

This guide is intended for educational purposes, helping developers adapt and optimize configurations for their specific use cases. Additionally, there are existing tools that support Laravel in containers:

1. **Laravel Sail**: An official package for easily starting Laravel in Docker ([Laravel Sail documentation](https://laravel.com/docs/11.x/sail)).
2. **Laradock**: A community project that helps run Laravel applications in Docker ([Laradock GitHub repository](https://github.com/laradock/laradock)).

## What you’ll learn

- How to use Docker Compose to set up a Laravel development and production environment.
- Defining services that make Laravel development easier, including PHP-FPM, Nginx, and database containers.
- Best practices for managing Laravel environments using containerization.

## Who’s this for?

- Developers who work with Laravel and want to streamline environment management.
- DevOps engineers seeking efficient ways to manage and deploy Laravel applications.

## Tools integration

Docker Compose integrates seamlessly with the Docker CLI, CI/CD tools, and container orchestration platforms like Docker Swarm and Kubernetes. This makes defining, managing, and orchestrating Laravel applications consistent across all environments.
