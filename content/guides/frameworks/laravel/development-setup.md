---
title: Laravel Development Setup with Docker Compose
description: Set up a Laravel development environment using Docker Compose.
weight: 20
---

## Development Environment Setup

This guide demonstrates how to set up a development environment for a Laravel application using Docker Compose. This setup includes essential services such as PHP-FPM, Nginx, and a database (we will use Postgres, but MySQL/MariaDB can be easily set), which enable you to develop in an isolated and consistent environment.

> [!NOTE]
> If you want to quickly test this setup without configuring everything manually, you can download the [Laravel Docker Examples](https://github.com/rw4lll/laravel-docker-examples) repository. It includes pre-configured examples for both development and production environments.

### Project Structure

To start, create a project structure that will include both the Laravel application and Docker-related files:

```
example-app
├── app, config, routes, tests, etc.
├── docker/
│   ├── php-fpm
│   │   └── Dockerfile
│   │   └── entrypoint.sh
│   ├── workspace
│   │   └── Dockerfile
│   └── web
│       └── nginx.conf
├── compose.yaml
├── .dockerignore
└── .env
└── other files
```


This structure includes a typical Laravel app, with a `docker` directory for Docker-related files like `php-fpm` and `workspace` Dockerfiles, as well as `nginx.conf` config file, and the `compose.yaml` file to define the services.


### Writing the Dockerfile for PHP-FPM

The PHP-FPM Dockerfile defines the environment in which PHP will run. Here is an example:

```dockerfile
# docker/php-fpm/Dockerfile
# For development environment we can use one-stage build for simplicity.
FROM php:8.3-fpm

# Install system dependencies and PHP extensions required for Laravel + MySQL/PostgreSQL support
# Some dependencies are required for PHP extensions only in the build stage
# We don't need to install Node.js or build assets, as it was done in the Nginx image
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
    opcache \
    intl \
    zip \
    bcmath \
    soap \
    && pecl install redis xdebug \
    && docker-php-ext-enable redis \
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
    docker-php-ext-enable xdebug && \
    echo "xdebug.mode=${XDEBUG_MODE}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.idekey=${XDEBUG_IDE_KEY}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.log=${XDEBUG_LOG}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.log_level=${XDEBUG_LOG_LEVEL}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.client_host=${XDEBUG_HOST}" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini ; \
    echo "xdebug.start_with_request=yes" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini ; \
fi

# Set environment variables for user and group ID
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
COPY ./docker/php-fpm/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Change the default command to run the entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Expose port 9000 and start php-fpm server
EXPOSE 9000
CMD ["php-fpm"]
```

This Dockerfile installs the necessary PHP extensions required by Laravel, including database drivers and the Xdebug extension for debugging.

### Writing the Dockerfile for Workspace

The workspace container is used to run Artisan commands, Composer, and NPM. Here's the Dockerfile for the workspace:

```dockerfile
# docker/workspace/Dockerfile
# Use the official PHP CLI image as the base
FROM php:8.3-cli

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
    opcache \
    intl \
    zip \
    bcmath \
    soap \
    && pecl install redis xdebug \
    && docker-php-ext-enable redis xdebug\
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

# Install NVM as the www user
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

### Docker Compose Configuration for Development

Here's the `compose.yaml` file to set up the development environment:

```yaml
services:
  web:
    image: nginx:latest # We don't need to customize the image. Just pass the configuration to the Dockerfile.
    volumes:
      # Mount the application code for live updates
      - ./:/var/www
      # Mount the Nginx configuration file
      - ./docker/web/nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      # Map port 80 inside the container to the port specified by 'NGINX_PORT' on the host machine
      - "${NGINX_PORT:-80}:80"
    environment:
      - NGINX_HOST=${NGINX_HOST}
    networks:
      - laravel
    depends_on:
      php-fpm:
        condition: service_started  # Wait for php-fpm to start

  php-fpm:
    # For the php-fpm service, we will create a custom image to install the necessary PHP extensions and setup proper permissions.
    build:
      context: .
      dockerfile: ./docker/php-fpm/Dockerfile
      args:
        UID: ${UID}
        GID: ${GID}
        XDEBUG_ENABLED: ${XDEBUG_ENABLED}
        XDEBUG_MODE: ${XDEBUG_MODE}
        XDEBUG_HOST: ${XDEBUG_HOST}
        XDEBUG_IDE_KEY: ${XDEBUG_IDE_KEY}
        XDEBUG_LOG: ${XDEBUG_LOG}
        XDEBUG_LOG_LEVEL: ${XDEBUG_LOG_LEVEL}
    env_file:
      # Load the environment variables from the Laravel application
      - .env
    user: "${UID}:${GID}"
    volumes:
      # Mount the application code for live updates
      - ./:/var/www
    networks:
      - laravel
    depends_on:
      postgres:
        condition: service_started  # Wait for postgres to start

  workspace:
   # For the workspace service, we will also create a custom image to install and setup all the necessary stuff.
    build:
      context: .
      dockerfile: ./docker/workspace/Dockerfile
      args:
        UID: ${UID}
        GID: ${GID}
        XDEBUG_ENABLED: ${XDEBUG_ENABLED}
        XDEBUG_MODE: ${XDEBUG_MODE}
        XDEBUG_HOST: ${XDEBUG_HOST}
        XDEBUG_IDE_KEY: ${XDEBUG_IDE_KEY}
        XDEBUG_LOG: ${XDEBUG_LOG}
        XDEBUG_LOG_LEVEL: ${XDEBUG_LOG_LEVEL}
    tty: true  # Enables an interactive terminal
    stdin_open: true  # Keeps standard input open for 'docker exec'
    env_file:
      - .env
    volumes:
      - ./:/var/www
    networks:
      - laravel

  postgres:
    image: postgres:16
    ports:
      - "${POSTGRES_PORT}:5432"
    environment:
      - POSTGRES_DB=${POSTGRES_DATABASE}
      - POSTGRES_USER=${POSTGRES_USERNAME}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - laravel

  redis:
    image: redis:alpine
    networks:
      - laravel

networks:
  laravel:

volumes:
  postgres-data:
```

> [!NOTE]
> Ensure you have an `.env` file at the root of your Laravel project with the necessary configurations (e.g., database and Xdebug settings) to match the Docker Compose setup.

### Running Your Development Environment

To start your Laravel development environment, run the following command in your terminal:

```sh
$ docker compose up -d
```

This command will build and start all the required services, including PHP, Nginx, and the PostgreSQL database. You can now access your Laravel application at `http://localhost/`.

### Summary

By setting up a Docker Compose environment for Laravel development, you ensure that your development setup is consistent and easily reproducible. This makes it easier for you and your team to collaborate on the same project, without worrying about differences in local environments.

<div id="compose-lp-survey-anchor"></div>
