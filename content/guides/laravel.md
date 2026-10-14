---
title: Develop and Deploy Laravel applications with Docker Compose
linkTitle: Laravel
summary: Learn how to efficiently set up Laravel development and production environments using Docker Compose.
description: A guide on using Docker Compose to manage Laravel applications for development and production, covering container configurations and service management.
keywords: laravel, php, docker compose, web framework, development, production
aliases:
  - /frameworks/laravel/
  - /guides/frameworks/laravel/
  - /guides/frameworks/laravel/common-questions/
  - /guides/frameworks/laravel/development-setup/
  - /guides/frameworks/laravel/prerequisites/
  - /guides/frameworks/laravel/production-setup/
params:
  tags: [languages]
  time: 30 minutes
---


Laravel is a popular PHP framework that allows developers to build web applications quickly and effectively. Docker Compose simplifies the management of development and production environments by defining essential services, like PHP, a web server, and a database, in a single YAML file. This guide provides a streamlined approach to setting up a robust Laravel environment using Docker Compose, focusing on simplicity and efficiency.

> **Acknowledgment**
>
> Docker would like to thank [Sergei Shitikov](https://github.com/rw4lll) for
> his contribution to this guide.

The demonstrated examples can be found in [this GitHub repository](https://github.com/dockersamples/laravel-docker-examples). Docker Compose offers a straightforward approach to connecting multiple containers for Laravel, though similar setups can also be achieved using tools like Docker Swarm, Kubernetes, or individual Docker containers.

This guide is intended for educational purposes, helping developers adapt and optimize configurations for their specific use cases. Additionally, there are existing tools that support Laravel in containers:

- [Laravel Sail](https://laravel.com/docs/12.x/sail): An official package for easily starting Laravel in Docker.
- [Laradock](https://github.com/laradock/laradock): A community project that helps run Laravel applications in Docker.

## What you’ll learn

- How to use Docker Compose to set up a Laravel development and production environment.
- Defining services that make Laravel development easier, including PHP-FPM, Nginx, and database containers.
- Best practices for managing Laravel environments using containerization.

## Who’s this for?

- Developers who work with Laravel and want to streamline environment management.
- DevOps engineers seeking efficient ways to manage and deploy Laravel applications.

## Prerequisites for Setting Up Laravel with Docker Compose

Before you begin setting up Laravel with Docker Compose, make sure you meet the following prerequisites:

### Docker and Docker Compose

You need Docker and Docker Compose installed on your system. Docker allows you to containerize applications, and Docker Compose helps you manage multi-container applications.

- Docker: Make sure Docker is installed and running on your machine. Refer to the [Docker installation guide](/get-docker/) to install Docker.
- Docker Compose: Docker Compose is included with Docker Desktop, but you can also follow the [Docker Compose installation guide](/compose/install/) if needed.

### Basic understanding of Docker and containers

A fundamental understanding of Docker and how containers work will be helpful. If you're new to Docker, consider reviewing the [Docker Overview](/get-started/overview/) to familiarize yourself with containerization concepts.

### Basic knowledge of Laravel

This guide assumes you have a basic understanding of Laravel and PHP. Familiarity with Laravel’s command-line tools, such as [Artisan](https://laravel.com/docs/12.x/artisan), and its project structure is important for following the instructions.

- Laravel CLI: You should be comfortable using Laravel’s command-line tool (`artisan`).
- Laravel Project Structure: Familiarize yourself with Laravel’s folder structure (`app`, `config`, `routes`, `tests`, etc.).

## Laravel Production Setup with Docker Compose

This guide demonstrates how to set up a production-ready Laravel environment using Docker and Docker Compose. This configuration is designed for streamlined, scalable, and secure Laravel application deployments.

> [!NOTE]
> To experiment with a ready-to-run configuration, download the [Laravel Docker Examples](https://github.com/dockersamples/laravel-docker-examples) repository. It contains pre-configured setups for both development and production.

### Project structure

```plaintext
my-laravel-app/
├── app/
├── bootstrap/
├── config/
├── database/
├── public/
├── docker/
│   ├── common/
│   │   └── php-fpm/
│   │       ├── Dockerfile
│   │       └── conf.d/
│   │           └── 20-status-path.conf
│   ├── development/
│   ├── production/
│   │   ├── php-fpm/
│   │   │   └── entrypoint.sh
│   │   └── nginx
│   │       ├── Dockerfile
│   │       └── nginx.conf
├── compose.dev.yaml
├── compose.prod.yaml
├── .dockerignore
├── .env
├── vendor/
├── ...
```

This layout represents a typical Laravel project, with Docker configurations stored in a unified `docker` directory. You’ll find **two** Compose files — `compose.dev.yaml` (for development) and `compose.prod.yaml` (for production) — to keep your environments separate and manageable.

### Create a Dockerfile for PHP-FPM (production)

For production, the `php-fpm` Dockerfile creates an optimized image with only the PHP extensions and libraries your application needs. As demonstrated in the [GitHub example](https://github.com/dockersamples/laravel-docker-examples), a single Dockerfile with multi-stage builds maintains consistency and reduces duplication between development and production. The following snippet shows only the production-related stages:

```dockerfile
# Stage 1: Build environment and Composer dependencies
FROM php:8.5-fpm AS builder

# Install system dependencies and PHP extensions for Laravel with MySQL/PostgreSQL support.
# Dependencies in this stage are only required for building the final image.
# Node.js and asset building are handled in the Nginx stage, not here.
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    unzip \
    libpq-dev \
    libonig-dev \
    libssl-dev \
    libxml2-dev \
    libcurl4-openssl-dev \
    libicu-dev \
    libzip-dev \
    && docker-php-ext-install -j$(nproc) \
    pdo_mysql \
    pdo_pgsql \
    pgsql \
    intl \
    zip \
    bcmath \
    soap \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && apt-get autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Set the working directory inside the container
WORKDIR /var/www

# Copy the entire Laravel application code into the container
# -----------------------------------------------------------
# In Laravel, `composer install` may trigger scripts
# needing access to application code.
# For example, the `post-autoload-dump` event might execute
# Artisan commands like `php artisan package:discover`. If the
# application code (including the `artisan` file) is not
# present, these commands will fail, leading to build errors.
#
# By copying the entire application code before running
# `composer install`, we ensure that all necessary files are
# available, allowing these scripts to run successfully.
# In other cases, it would be possible to copy composer files
# first, to leverage Docker's layer caching mechanism.
# -----------------------------------------------------------
COPY . /var/www

# Install Composer and dependencies
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && composer install --no-dev --optimize-autoloader --no-interaction --no-progress --prefer-dist

# Stage 2: Production environment
FROM php:8.5-fpm AS production

# Install only runtime libraries needed in production
# libfcgi-bin and procps are required for the php-fpm-healthcheck script
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq-dev \
    libicu-dev \
    libzip-dev \
    libfcgi-bin \
    procps \
    && apt-get autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Download and install php-fpm health check script
RUN curl -o /usr/local/bin/php-fpm-healthcheck \
    https://raw.githubusercontent.com/renatomefi/php-fpm-healthcheck/master/php-fpm-healthcheck \
    && chmod +x /usr/local/bin/php-fpm-healthcheck

# Copy the initialization script
COPY ./docker/php-fpm/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Copy the initial storage structure
COPY ./storage /var/www/storage-init

# Copy PHP extensions and libraries from the builder stage
COPY --from=builder /usr/local/lib/php/extensions/ /usr/local/lib/php/extensions/
COPY --from=builder /usr/local/etc/php/conf.d/ /usr/local/etc/php/conf.d/
COPY --from=builder /usr/local/bin/docker-php-ext-* /usr/local/bin/

# Use the recommended production PHP configuration
# -----------------------------------------------------------
# PHP provides development and production configurations.
# Here, we replace the default php.ini with the production
# version to apply settings optimized for performance and
# security in a live environment.
# -----------------------------------------------------------
RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"

# Keep the image-provided FPM global config intact and add pool overrides separately
COPY ./docker/common/php-fpm/conf.d/*.conf /usr/local/etc/php-fpm.d/
# Update the variables_order to include E (for ENV)
#RUN sed -i 's/variables_order = "GPCS"/variables_order = "EGPCS"/' "$PHP_INI_DIR/php.ini"

# Copy the application code and dependencies from the build stage
COPY --from=builder /var/www /var/www

# Set working directory
WORKDIR /var/www

# Ensure correct permissions
RUN chown -R www-data:www-data /var/www

# Switch to the non-privileged user to run the application
USER www-data

# Change the default command to run the entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Expose port 9000 and start php-fpm server
EXPOSE 9000
CMD ["php-fpm"]
```

### Create a Dockerfile for PHP-CLI (production)

For production, you often need a separate container to run Artisan commands, migrations, and other CLI tasks. In most cases you can run these commands by reusing existing PHP-FPM container:

```console
$ docker compose -f compose.prod.yaml exec php-fpm php artisan route:list
```

If you need a separate CLI container with different extensions or strict separation of concerns, consider a php-cli Dockerfile:

```dockerfile
# Stage 1: Build environment and Composer dependencies
FROM php:8.5-cli AS builder

# Install system dependencies and PHP extensions required for Laravel + MySQL/PostgreSQL support
# Some dependencies are required for PHP extensions only in the build stage
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    unzip \
    libpq-dev \
    libonig-dev \
    libssl-dev \
    libxml2-dev \
    libcurl4-openssl-dev \
    libicu-dev \
    libzip-dev \
    && docker-php-ext-install -j$(nproc) \
    pdo_mysql \
    pdo_pgsql \
    pgsql \
    intl \
    zip \
    bcmath \
    soap \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && apt-get autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Set the working directory inside the container
WORKDIR /var/www

# Copy the entire Laravel application code into the container
COPY . /var/www

# Install Composer and dependencies
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && composer install --no-dev --optimize-autoloader --no-interaction --no-progress --prefer-dist

# Stage 2: Production environment
FROM php:8.5-cli

# Install client libraries required for php extensions in runtime
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq-dev \
    libicu-dev \
    libzip-dev \
    && apt-get autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Copy PHP extensions and libraries from the builder stage
COPY --from=builder /usr/local/lib/php/extensions/ /usr/local/lib/php/extensions/
COPY --from=builder /usr/local/etc/php/conf.d/ /usr/local/etc/php/conf.d/
COPY --from=builder /usr/local/bin/docker-php-ext-* /usr/local/bin/

# Use the default production configuration for PHP runtime arguments
RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"

# Copy the application code and dependencies from the build stage
COPY --from=builder /var/www /var/www

# Set working directory
WORKDIR /var/www

# Ensure correct permissions
RUN chown -R www-data:www-data /var/www

# Switch to the non-privileged user to run the application
USER www-data

# Default command: Provide a bash shell to allow running any command
CMD ["bash"]
```

This Dockerfile is similar to the PHP-FPM Dockerfile, but it uses the `php:8.5-cli` image as the base image and sets up the container for running CLI commands.

### Create a Dockerfile for Nginx (production)

Nginx serves as the web server for the Laravel application. You can include static assets directly to the container. Here's an example of possible Dockerfile for Nginx:

```dockerfile
# docker/nginx/Dockerfile
# Stage 1: Build assets
FROM debian AS builder

# Install Node.js and build tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    nodejs \
    npm \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Set working directory
WORKDIR /var/www

# Copy Laravel application code
COPY . /var/www

# Install Node.js dependencies and build assets
RUN npm install && npm run build

# Stage 2: Nginx production image
FROM nginx:alpine

# Copy custom Nginx configuration
# -----------------------------------------------------------
# Replace the default Nginx configuration with our custom one
# that is optimized for serving a Laravel application.
# -----------------------------------------------------------
COPY ./docker/nginx/nginx.conf /etc/nginx/nginx.conf

# Copy Laravel's public assets from the builder stage
# -----------------------------------------------------------
# We only need the 'public' directory from our Laravel app.
# -----------------------------------------------------------
COPY --from=builder /var/www/public /var/www/public

# Set the working directory to the public folder
WORKDIR /var/www/public

# Expose port 80 and start Nginx
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

This Dockerfile uses a multi-stage build to separate the asset building process from the final production image. The first stage installs Node.js and builds the assets, while the second stage sets up the Nginx production image with the optimized configuration and the built assets.

### Create a Docker Compose configuration for production

To bring all the services together, create a `compose.prod.yaml` file that defines the services, volumes, and networks for the production environment. Here's an example configuration:

```yaml
services:
  web:
    build:
      context: .
      dockerfile: ./docker/production/nginx/Dockerfile
    restart: unless-stopped # Automatically restart unless the service is explicitly stopped
    volumes:
      # Mount the 'laravel-storage' volume to '/var/www/storage' inside the container.
      # -----------------------------------------------------------
      # This volume stores persistent data like uploaded files and cache.
      # The ':ro' option mounts it as read-only in the 'web' service because Nginx only needs to read these files.
      # The 'php-fpm' service mounts the same volume without ':ro' to allow write operations.
      # -----------------------------------------------------------
      - laravel-storage-production:/var/www/storage:ro
    networks:
      - laravel-production
    ports:
      # Map port 80 inside the container to the port specified by 'NGINX_PORT' on the host machine.
      # -----------------------------------------------------------
      # This allows external access to the Nginx web server running inside the container.
      # For example, if 'NGINX_PORT' is set to '8080', accessing 'http://localhost:8080' will reach the application.
      # -----------------------------------------------------------
      - "${NGINX_PORT:-80}:80"
    depends_on:
      php-fpm:
        condition: service_healthy # Wait for php-fpm health check

  php-fpm:
    # For the php-fpm service, we will create a custom image to install the necessary PHP extensions and setup proper permissions.
    build:
      context: .
      dockerfile: ./docker/common/php-fpm/Dockerfile
      target: production # Use the 'production' stage in the Dockerfile
    restart: unless-stopped
    volumes:
      - laravel-storage-production:/var/www/storage # Mount the storage volume
    env_file:
      - .env
    networks:
      - laravel-production
    healthcheck:
      test: ["CMD-SHELL", "php-fpm-healthcheck || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 3
    # The 'depends_on' attribute with 'condition: service_healthy' ensures that
    # this service will not start until the 'postgres' service passes its health check.
    # This prevents the application from trying to connect to the database before it's ready.
    depends_on:
      postgres:
        condition: service_healthy

  # The 'php-cli' service provides a command-line interface for running Artisan commands and other CLI tasks.
  # -----------------------------------------------------------
  # This is useful for running migrations, seeders, or any custom scripts.
  # It shares the same codebase and environment as the 'php-fpm' service.
  # -----------------------------------------------------------
  php-cli:
    build:
      context: .
      dockerfile: ./docker/php-cli/Dockerfile
    tty: true # Enables an interactive terminal
    stdin_open: true # Keeps standard input open for 'docker exec'
    env_file:
      - .env
    networks:
      - laravel-production

  postgres:
    image: postgres:18
    restart: unless-stopped
    user: postgres
    ports:
      - "${POSTGRES_PORT}:5432"
    environment:
      - POSTGRES_DB=${POSTGRES_DATABASE}
      - POSTGRES_USER=${POSTGRES_USERNAME}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres-data-production:/var/lib/postgresql
    networks:
      - laravel-production
    # Health check for PostgreSQL
    # -----------------------------------------------------------
    # Health checks allow Docker to determine if a service is operational.
    # The 'pg_isready' command checks if PostgreSQL is ready to accept connections.
    # This prevents dependent services from starting before the database is ready.
    # -----------------------------------------------------------
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:alpine
    restart: unless-stopped # Automatically restart unless the service is explicitly stopped
    networks:
      - laravel-production
    # Health check for Redis
    # -----------------------------------------------------------
    # Checks if Redis is responding to the 'PING' command.
    # This ensures that the service is not only running but also operational.
    # -----------------------------------------------------------
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3

networks:
  # Attach the service to the 'laravel-production' network.
  # -----------------------------------------------------------
  # This custom network allows all services within it to communicate using their service names as hostnames.
  # For example, 'php-fpm' can connect to 'postgres' by using 'postgres' as the hostname.
  # -----------------------------------------------------------
  laravel-production:

volumes:
  postgres-data-production:
  laravel-storage-production:
```

> [!NOTE]
> Ensure you have an `.env` file at the root of your Laravel project with the necessary configurations to match the Docker Compose setup.

### Running your production environment

To start the production environment, run:

```console
$ docker compose -f compose.prod.yaml up --build -d
```

This command will build and start all the services in detached mode, providing a scalable and production-ready setup for your Laravel application.

### Summary

By setting up a Docker Compose environment for Laravel in production, you ensure that your application is optimized for performance, scalable, and secure. This setup makes deployments consistent and easier to manage, reducing the likelihood of errors due to differences between environments.

## Laravel Development Setup with Docker Compose

This guide demonstrates how to configure a **development** environment for a Laravel application using Docker and Docker Compose. It builds **on top of** the production image for PHP-FPM and then adds developer-focused features—like Xdebug—to streamline debugging. By basing the development container on a known production image, you keep both environments closely aligned.

This setup includes PHP-FPM, Nginx, and PostgreSQL services (although you can easily swap PostgreSQL for another database, like MySQL or MariaDB). Everything runs in containers, so you can develop in isolation without altering your host system.

> [!NOTE]
> To experiment with a ready-to-run configuration, download the [Laravel Docker Examples](https://github.com/dockersamples/laravel-docker-examples) repository. It contains pre-configured setups for both development and production.

### Project structure

```plaintext
my-laravel-app/
├── app/
├── bootstrap/
├── config/
├── database/
├── public/
├── docker/
│   ├── common/
│   │   └── php-fpm/
│   │       └── Dockerfile
│   ├── development/
│   │   ├── php-fpm/
│   │   │   └── entrypoint.sh
│   │   ├── workspace/
│   │   │   └── Dockerfile
│   │   └── nginx
│   │       ├── Dockerfile
│   │       └── nginx.conf
│   └── production/
├── compose.dev.yaml
├── compose.prod.yaml
├── .dockerignore
├── .env
├── vendor/
├── ...
```

This layout represents a typical Laravel project, with Docker configurations stored in a unified `docker` directory. You’ll find **two** Compose files — `compose.dev.yaml` (for development) and `compose.prod.yaml` (for production) — to keep your environments separate and manageable.

The environment includes a `workspace` service, a sidecar container for tasks like building front-end assets, running Artisan commands, and other CLI tools your project may require. While this extra container may seem unusual, it’s a familiar pattern in solutions like **Laravel Sail** and **Laradock**. It also includes **Xdebug** to aid in debugging.

### Create a Dockerfile for PHP-FPM

This Dockerfile **extends** the production image by installing Xdebug and adjusting user permissions to ease local development. That way, your development environment stays consistent with production while still offering extra debug features and improved file mounting.

```dockerfile
# Builds a dev-only layer on top of the production image
FROM production AS development

# Use ARGs to define environment variables passed from the Docker build command or Docker Compose.
ARG XDEBUG_ENABLED=true
ARG XDEBUG_MODE=develop,coverage,debug,profile
ARG XDEBUG_HOST=host.docker.internal
ARG XDEBUG_IDE_KEY=DOCKER
ARG XDEBUG_LOG=/dev/stdout
ARG XDEBUG_LOG_LEVEL=0

USER root

# Configure Xdebug if enabled
RUN if [ "${XDEBUG_ENABLED}" = "true" ]; then \
    pecl install xdebug && \
    docker-php-ext-enable xdebug && \
    echo "xdebug.mode=${XDEBUG_MODE}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.idekey=${XDEBUG_IDE_KEY}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.log=${XDEBUG_LOG}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.log_level=${XDEBUG_LOG_LEVEL}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.client_host=${XDEBUG_HOST}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini ; \
    echo "xdebug.start_with_request=yes" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini ; \
fi

# Add ARGs for syncing permissions
ARG UID=1000
ARG GID=1000

# Create a new user with the specified UID and GID, reusing an existing group if GID exists
RUN if getent group ${GID}; then \
      group_name=$(getent group ${GID} | cut -d: -f1); \
      useradd -m -u ${UID} -g ${GID} -s /bin/bash www; \
    else \
      groupadd -g ${GID} www && \
      useradd -m -u ${UID} -g www -s /bin/bash www; \
      group_name=www; \
    fi

# Dynamically update php-fpm to use the new user and group
RUN sed -i "s/user = www-data/user = www/g" /usr/local/etc/php-fpm.d/www.conf && \
    sed -i "s/group = www-data/group = $group_name/g" /usr/local/etc/php-fpm.d/www.conf


# Set the working directory
WORKDIR /var/www

# Copy the entrypoint script
COPY ./docker/development/php-fpm/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Switch back to the non-privileged user to run the application
USER www-data

# Change the default command to run the entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Expose port 9000 and start php-fpm server
EXPOSE 9000
CMD ["php-fpm"]
```

### Create a Dockerfile for Workspace

A workspace container provides a dedicated shell for asset compilation, Artisan/Composer commands, and other CLI tasks. This approach follows patterns from Laravel Sail and Laradock, consolidating all development tools into one container for convenience.

```dockerfile
# docker/development/workspace/Dockerfile
# Use the official PHP CLI image as the base
FROM php:8.5-cli

# Set environment variables for user and group ID
ARG UID=1000
ARG GID=1000
ARG NODE_VERSION=22.0.0

# Install system dependencies and build libraries
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    unzip \
    libpq-dev \
    libonig-dev \
    libssl-dev \
    libxml2-dev \
    libcurl4-openssl-dev \
    libicu-dev \
    libzip-dev \
    && docker-php-ext-install -j$(nproc) \
    pdo_mysql \
    pdo_pgsql \
    pgsql \
    intl \
    zip \
    bcmath \
    soap \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && apt-get autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Use ARG to define environment variables passed from the Docker build command or Docker Compose.
ARG XDEBUG_ENABLED
ARG XDEBUG_MODE
ARG XDEBUG_HOST
ARG XDEBUG_IDE_KEY
ARG XDEBUG_LOG
ARG XDEBUG_LOG_LEVEL

# Configure Xdebug if enabled
RUN if [ "${XDEBUG_ENABLED}" = "true" ]; then \
    pecl install xdebug && \
    docker-php-ext-enable xdebug && \
    echo "xdebug.mode=${XDEBUG_MODE}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.idekey=${XDEBUG_IDE_KEY}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.log=${XDEBUG_LOG}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.log_level=${XDEBUG_LOG_LEVEL}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.client_host=${XDEBUG_HOST}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini ; \
    echo "xdebug.start_with_request=yes" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini ; \
fi

# If the group already exists, use it; otherwise, create the 'www' group
RUN if getent group ${GID}; then \
      useradd -m -u ${UID} -g ${GID} -s /bin/bash www; \
    else \
      groupadd -g ${GID} www && \
      useradd -m -u ${UID} -g www -s /bin/bash www; \
    fi && \
    usermod -aG sudo www && \
    echo 'www ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Switch to the non-root user to install NVM and Node.js
USER www

# Install NVM (Node Version Manager) as the www user
RUN export NVM_DIR="$HOME/.nvm" && \
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.0/install.sh | bash && \
    [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh" && \
    nvm install ${NODE_VERSION} && \
    nvm alias default ${NODE_VERSION} && \
    nvm use default

# Ensure NVM is available for all future shells
RUN echo 'export NVM_DIR="$HOME/.nvm"' >> /home/www/.bashrc && \
    echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> /home/www/.bashrc && \
    echo '[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"' >> /home/www/.bashrc

# Set the working directory
WORKDIR /var/www

# Override the entrypoint to avoid the default php entrypoint
ENTRYPOINT []

# Default command to keep the container running
CMD ["bash"]
```

> [!NOTE]
> If you prefer a **one-service-per-container** approach, simply omit the workspace container and run separate containers for each task. For example, you could use a dedicated `php-cli` container for your PHP scripts, and a `node` container to handle the asset building.

### Create a Docker Compose configuration for development

Here's the `compose.yaml` file to set up the development environment:

```yaml
services:
  web:
    image: nginx:latest # Using the default Nginx image with custom configuration.
    volumes:
      # Mount the application code for live updates
      - ./:/var/www
      # Mount the Nginx configuration file
      - ./docker/development/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      # Map port 80 inside the container to the port specified by 'NGINX_PORT' on the host machine
      - "80:80"
    environment:
      - NGINX_HOST=localhost
    networks:
      - laravel-development
    depends_on:
      php-fpm:
        condition: service_started # Wait for php-fpm to start

  php-fpm:
    # For the php-fpm service, we will use our common PHP-FPM Dockerfile with the development target
    build:
      context: .
      dockerfile: ./docker/common/php-fpm/Dockerfile
      target: development
      args:
        UID: ${UID:-1000}
        GID: ${GID:-1000}
        XDEBUG_ENABLED: ${XDEBUG_ENABLED:-true}
        XDEBUG_MODE: develop,coverage,debug,profile
        XDEBUG_HOST: ${XDEBUG_HOST:-host.docker.internal}
        XDEBUG_IDE_KEY: ${XDEBUG_IDE_KEY:-DOCKER}
        XDEBUG_LOG: /dev/stdout
        XDEBUG_LOG_LEVEL: 0
    env_file:
      # Load the environment variables from the Laravel application
      - .env
    user: "${UID:-1000}:${GID:-1000}"
    volumes:
      # Mount the application code for live updates
      - ./:/var/www
    networks:
      - laravel-development
    depends_on:
      postgres:
        condition: service_started # Wait for postgres to start

  workspace:
    # For the workspace service, we will also create a custom image to install and setup all the necessary stuff.
    build:
      context: .
      dockerfile: ./docker/development/workspace/Dockerfile
      args:
        UID: ${UID:-1000}
        GID: ${GID:-1000}
        XDEBUG_ENABLED: ${XDEBUG_ENABLED:-true}
        XDEBUG_MODE: develop,coverage,debug,profile
        XDEBUG_HOST: ${XDEBUG_HOST:-host.docker.internal}
        XDEBUG_IDE_KEY: ${XDEBUG_IDE_KEY:-DOCKER}
        XDEBUG_LOG: /dev/stdout
        XDEBUG_LOG_LEVEL: 0
    tty: true # Enables an interactive terminal
    stdin_open: true # Keeps standard input open for 'docker exec'
    env_file:
      - .env
    volumes:
      - ./:/var/www
    networks:
      - laravel-development

  postgres:
    image: postgres:18
    ports:
      - "${POSTGRES_PORT:-5432}:5432"
    environment:
      - POSTGRES_DB=app
      - POSTGRES_USER=laravel
      - POSTGRES_PASSWORD=secret
    volumes:
      - postgres-data-development:/var/lib/postgresql
    networks:
      - laravel-development

  redis:
    image: redis:alpine
    networks:
      - laravel-development

networks:
  laravel-development:

volumes:
  postgres-data-development:
```

> [!NOTE]
> Ensure you have an `.env` file at the root of your Laravel project with the necessary configurations. You can use the `.env.example` file as a template.

### Run your development environment

To start the development environment, use:

```console
$ docker compose -f compose.dev.yaml up --build -d
```

Run this command to build and start the development environment in detached mode. When the containers finish initializing, visit [http://localhost/](http://localhost/) to see your Laravel app in action.

### Summary

By building on top of the production image and adding debug tools like Xdebug, you create a Laravel development workflow that closely mirrors production. The optional workspace container simplifies tasks like asset building and running Artisan commands. If you prefer a separate container for every service (e.g., a dedicated `php-cli` and `node` container), you can skip the workspace approach. Either way, Docker Compose provides an efficient, consistent way to develop your Laravel project.

## Common Questions on Using Laravel with Docker

<!-- vale Docker.HeadingLength = NO -->

### 1. Why should I use Docker Compose for Laravel?

Docker Compose is a powerful tool for managing multi-container environments, particularly in development due to its simplicity. With Docker Compose, you can define and connect all necessary services for Laravel, such as PHP, Nginx, and databases, in a single configuration (`compose.*.yaml`). This setup ensures consistency across development, testing, and production environments, streamlining onboarding and reducing discrepancies between local and server setups.

While Docker Compose is a great choice for development, tools like **Docker Swarm** or **Kubernetes** offer advanced scaling and orchestration features, which may be beneficial for complex production deployments.

### 2. How do I debug my Laravel application with Docker Compose?

To debug your Laravel application in a Docker environment, use **Xdebug**. In the development setup, Xdebug is installed in the `php-fpm` container to enable debugging. Ensure Xdebug is enabled in your `compose.dev.yaml` file by setting the environment variable `XDEBUG_ENABLED=true` and configuring your IDE (e.g., Visual Studio Code or PHPStorm) to connect to the remote container for debugging.

### 3. Can I use Docker Compose with databases other than PostgreSQL?

Yes, Docker Compose supports various database services for Laravel. While PostgreSQL is used in the examples, you can easily substitute **MySQL**, **MariaDB**, or even **SQLite**. Update the `compose.*.yaml` file to specify the required Docker image and adjust your `.env` file to reflect the new database configuration.

### 4. How can I persist data in development and production?

In both development and production, Docker volumes are used to persist data. For instance, in the `compose.*.yaml` file, the `postgres-data-*` volume stores PostgreSQL data, ensuring that data is retained even if the container restarts. You can also define named volumes for other services where data persistence is essential.

### 5. What is the difference between development and production Docker configurations?

In a development environment, Docker configurations include tools that streamline coding and debugging, such as Xdebug for debugging, and volume mounts to enable real-time code updates without requiring image rebuilds.

In production, the configuration is optimized for performance, security, and efficiency. This setup uses multi-stage builds to keep the image lightweight and includes only essential tools, packages, and libraries.

It’s recommended to use `alpine`-based images in production for smaller image sizes, enhancing deployment speed and security.

Additionally, consider using [Docker Scout](/manuals/scout/_index.md) to detect and analyze vulnerabilities, especially in production environments.

For additional information about using Docker Compose in production, see [this guide](/compose/how-tos/production/).
