---
title: Develop and Deploy Laravel applications with Docker Compose
linkTitle: Laravel applications with Docker Compose
summary: Learn how to efficiently set up Laravel development and production environments using Docker Compose.
description: A guide on using Docker Compose to manage Laravel applications for development and production, covering container configurations and service management.
tags: [frameworks]
languages: [php]
aliases:
  - /frameworks/laravel/
params:
  time: 30 minutes
  resource_links:
    - title: Laravel
      url: https://laravel.com/
    - title: Docker Compose
      url: /compose/
    - title: Use Compose in production
      url: /compose/how-tos/production/
    - title: Repository with examples
      url: https://github.com/dockersamples/laravel-docker-examples
---

Laravel is a popular PHP framework that allows developers to build web applications quickly and effectively. Docker Compose simplifies the management of development and production environments by defining essential services, like PHP, a web server, and a database, in a single YAML file. This guide provides a streamlined approach to setting up a robust Laravel environment using Docker Compose, focusing on simplicity and efficiency.

> **Acknowledgment**
>
> Docker would like to thank [Sergei Shitikov](https://github.com/rw4lll) for
> his contribution to this guide.

The demonstrated examples can be found in [this GitHub repository](https://github.com/dockersamples/laravel-docker-examples). Docker Compose offers a straightforward approach to connecting multiple containers for Laravel, though similar setups can also be achieved using tools like Docker Swarm, Kubernetes, or individual Docker containers.

This guide is intended for educational purposes, helping developers adapt and optimize configurations for their specific use cases. Additionally, there are existing tools that support Laravel in containers:

- [Laravel Sail](https://laravel.com/docs/11.x/sail): An official package for easily starting Laravel in Docker.
- [Laradock](https://github.com/laradock/laradock): A community project that helps run Laravel applications in Docker.

## What you’ll learn

- How to use Docker Compose to set up a Laravel development and production environment.
- Defining services that make Laravel development easier, including PHP-FPM, Nginx, and database containers.
- Best practices for managing Laravel environments using containerization.

## Who’s this for?

- Developers who work with Laravel and want to streamline environment management.
- DevOps engineers seeking efficient ways to manage and deploy Laravel applications.
