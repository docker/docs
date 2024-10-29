---
title: Laravel Development Setup with Docker Compose
description: Set up a Laravel development environment using Docker Compose.
weight: 40
---

## Development Environment Setup

In this section, you'll learn how to set up a development environment for a Laravel application using Docker Compose. This setup will provide all necessary services such as PHP-FPM, Nginx, and a database, enabling you to develop locally in an isolated environment.

### Project Structure

To start, create a project structure that will include both the Laravel application and Docker-related files:

```
laravel-dockerize
├── development
│   └── nginx-fpm
│       ├── php-fpm
│       │   └── Dockerfile
│       ├── workspace
│       │   └── Dockerfile
│       └── compose.yaml
├── example-app
│   ├── app, config, routes, tests, etc.
```

This structure contains:
- **development**: Docker configurations for development.
- **example-app**: The Laravel application code.

### Writing the Dockerfile for PHP-FPM

The PHP-FPM Dockerfile defines the environment in which PHP will run. Here is an example:

```dockerfile
# development/nginx-fpm/php-fpm/Dockerfile
FROM php:8.3-fpm

# Install necessary PHP extensions
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    unzip \
    libpq-dev \
    libssl-dev \
    && docker-php-ext-install pdo_mysql pdo_pgsql zip bcmath \
    && pecl install redis xdebug \
    && docker-php-ext-enable redis xdebug

# Set working directory
WORKDIR /var/www
```

This Dockerfile installs the necessary PHP extensions required by Laravel, including database drivers and the Xdebug extension for debugging.

### Writing the Dockerfile for Workspace

The workspace container is used to run Artisan commands, Composer, and NPM. Here's the Dockerfile for the workspace:

```dockerfile
# development/nginx-fpm/workspace/Dockerfile
FROM php:8.3-cli

# Set environment variables for user and group ID
ARG UID=1000
ARG GID=1000

# Install dependencies and PHP extensions
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    unzip \
    libpq-dev \
    libssl-dev \
    && docker-php-ext-install pdo_mysql pdo_pgsql zip bcmath \
    && pecl install redis xdebug \
    && docker-php-ext-enable redis xdebug

# Set the working directory
WORKDIR /var/www

# Install Composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
```

### Docker Compose Configuration for Development

Here's the `compose.yaml` file to set up the development environment:

```yaml
version: '3.8'
services:
  web:
    image: nginx:latest
    volumes:
      - ./example-app:/var/www
      - ./example-app/nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "80:80"
    depends_on:
      - php-fpm
    networks:
      - laravel

  php-fpm:
    build:
      context: ./development/nginx-fpm/php-fpm
    volumes:
      - ./example-app:/var/www
    environment:
      - XDEBUG_ENABLED=true
    networks:
      - laravel

  workspace:
    build:
      context: ./development/nginx-fpm/workspace
    volumes:
      - ./example-app:/var/www
    tty: true
    stdin_open: true
    networks:
      - laravel

  postgres:
    image: postgres:16
    environment:
      - POSTGRES_DB=example
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=secret
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - laravel

networks:
  laravel:

volumes:
  postgres-data:
```

### Running Your Development Environment

To start your Laravel development environment, run the following command in your terminal:

```sh
$ docker compose -f development/nginx-fpm/compose.yaml up --build
```

This command will build and start all the required services, including PHP, Nginx, and the PostgreSQL database. You can now access your Laravel application at `http://localhost`.

### Summary

By setting up a Docker Compose environment for Laravel development, you ensure that your development setup is consistent and easily reproducible. This makes it easier for you and your team to collaborate on the same project, without worrying about differences in local environments.

<div id="compose-lp-survey-anchor"></div>
