---
title: Prerequisites for Using Docker Compose with Laravel
description: Ensure you have the required tools and knowledge before setting up Laravel with Docker Compose.
weight: 10
---

<!-- vale Docker.HeadingLength = NO -->

---

## Prerequisites

Before you begin setting up Laravel with Docker Compose, make sure you meet the following prerequisites:

### 1. Docker and Docker Compose Installed

You need Docker and Docker Compose installed on your system. Docker allows you to containerize applications, and Docker Compose helps you manage multi-container applications.

- **Docker**: Make sure Docker is installed and running on your machine. Refer to the [Docker Installation Guide](https://docs.docker.com/get-docker/) to install Docker.
- **Docker Compose**: Docker Compose is often included with Docker Desktop, but you can also follow the [Docker Compose Installation Guide](https://docs.docker.com/compose/install/) if needed.

### 2. Basic Understanding of Docker and Containers

A fundamental understanding of Docker and how containers work will be helpful. If you're new to Docker, consider reviewing the [Docker Overview](https://docs.docker.com/get-started/overview/) to familiarize yourself with containerization concepts.

### 3. Basic Knowledge of Laravel

This guide assumes you have a basic understanding of Laravel and PHP. Familiarity with Laravel’s command-line tools, such as `Artisan`, and its project structure is important for following the instructions.

- **Laravel CLI**: You should be comfortable using Laravel’s command-line tool (`artisan`).
- **Laravel Project Structure**: Familiarize yourself with Laravel’s folder structure (`app`, `config`, `routes`, `tests`, etc.).

### 4. Source Code for Laravel Application

You need the source code for a Laravel application to follow along with this guide. You can either use an existing Laravel project or create a new one. If you don't have a Laravel project, you can create one using Composer:

```sh
$ composer create-project --prefer-dist laravel/laravel example-app
```

### 5. Basic Command-Line Knowledge

You'll need basic command-line skills to navigate your system and run Docker commands. Commands provided in this guide are for Unix-based systems, but they should work similarly in PowerShell or Command Prompt on Windows.

### 6. Text Editor or IDE

To modify configuration files and Dockerfiles, you need a text editor or an IDE. Popular choices include [Visual Studio Code](https://code.visualstudio.com/), [PHPStorm](https://www.jetbrains.com/phpstorm/), or any editor that supports YAML and PHP.

### 7. Internet Access

Internet access is required to pull Docker images, install dependencies, and fetch Laravel packages during the setup process.

<div id="compose-lp-survey-anchor"></div>