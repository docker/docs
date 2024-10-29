---
title: Laravel Production Setup with Docker Compose
description: Set up a production-ready environment for Laravel using Docker Compose.
weight: 30
---

## Production Environment Setup

In this section, you'll learn how to set up a production-ready environment for a Laravel application using Docker Compose. This setup will focus on creating an efficient, secure, and scalable environment for your Laravel application.

### Project Structure

The project structure for production is similar to the development setup, but with specific Dockerfiles and configurations for production:

```
laravel-dockerize
├── production
│   └── nginx-fpm
│       ├── php-fpm
│       │   └── Dockerfile
│       ├── nginx
│       │   └── Dockerfile
│       └── compose.yaml
├── example-app
│   ├── app, config, routes, tests, etc.
```

This structure contains:
- **production**: Docker configurations for the production environment.
- **example-app**: The Laravel application code.

### Writing the Dockerfile for PHP-FPM (Production)

For production, the `php-fpm` Dockerfile aims to create a smaller, optimized image with only the necessary extensions and dependencies:

```dockerfile
# production/nginx-fpm/php-fpm/Dockerfile
FROM php:8.3-fpm AS builder

# Install dependencies required for Laravel
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    unzip \
    libpq-dev \
    && docker-php-ext-install pdo_mysql pdo_pgsql zip

# Copy application files
COPY ./example-app /var/www

# Set working directory
WORKDIR /var/www

# Run optimizations
RUN php artisan config:cache && php artisan route:cache

# Expose port
EXPOSE 9000

CMD ["php-fpm"]
```

This Dockerfile focuses on creating a lightweight image, installing only the required dependencies, and running optimization commands such as `config:cache` and `route:cache` to improve performance.

### Writing the Dockerfile for Nginx (Production)

Nginx serves as the web server for the Laravel application. Here's the Dockerfile for Nginx:

```dockerfile
# production/nginx-fpm/nginx/Dockerfile
FROM nginx:alpine

# Copy the Nginx configuration file
COPY nginx.conf /etc/nginx/nginx.conf

# Set the working directory
WORKDIR /var/www

# Expose port 80 for HTTP traffic
EXPOSE 80
```

This Dockerfile uses the lightweight `nginx:alpine` image and copies an Nginx configuration file to the appropriate location.

### Docker Compose Configuration for Production

Here's the `compose.yaml` file for the production environment:

```yaml
version: '3.8'
services:
  web:
    build:
      context: ./production/nginx-fpm/nginx
    ports:
      - "80:80"
    depends_on:
      - php-fpm
    networks:
      - laravel

  php-fpm:
    build:
      context: ./production/nginx-fpm/php-fpm
    volumes:
      - ./example-app:/var/www:ro
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

### Running Your Production Environment

To start the production environment, use the following command:

```sh
$ docker compose -f production/nginx-fpm/compose.yaml up --build -d
```

This command will build and start all the services in detached mode, providing a scalable and production-ready setup for your Laravel application.

### Summary

By setting up a Docker Compose environment for Laravel in production, you ensure that your application is optimized for performance, scalable, and secure. This setup makes deployments consistent and easier to manage, reducing the likelihood of errors due to differences between environments.

<div id="compose-lp-survey-anchor"></div>
