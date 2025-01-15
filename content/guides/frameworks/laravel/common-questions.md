---
title: Common Questions on Using Laravel with Docker
description: Find answers to common questions about setting up and managing Laravel environments with Docker Compose, including troubleshooting and best practices.
weight: 40
---

<!-- vale Docker.HeadingLength = NO -->

## 1. Why should I use Docker Compose for Laravel?

Docker Compose is a powerful tool for managing multi-container environments, particularly in development due to its simplicity. With Docker Compose, you can define and connect all necessary services for Laravel, such as PHP, Nginx, and databases, in a single configuration (`compose.*.yaml`). This setup ensures consistency across development, testing, and production environments, streamlining onboarding and reducing discrepancies between local and server setups.

While Docker Compose is a great choice for development, tools like **Docker Swarm** or **Kubernetes** offer advanced scaling and orchestration features, which may be beneficial for complex production deployments.

## 2. How do I debug my Laravel application with Docker Compose?

To debug your Laravel application in a Docker environment, use **Xdebug**. In the development setup, Xdebug is installed in the `php-fpm` container to enable debugging. Ensure Xdebug is enabled in your `compose.dev.yaml` file by setting the environment variable `XDEBUG_ENABLED=true` and configuring your IDE (e.g., Visual Studio Code or PHPStorm) to connect to the remote container for debugging.

## 3. Can I use Docker Compose with databases other than PostgreSQL?

Yes, Docker Compose supports various database services for Laravel. While PostgreSQL is used in the examples, you can easily substitute **MySQL**, **MariaDB**, or even **SQLite**. Update the `compose.*.yaml` file to specify the required Docker image and adjust your `.env` file to reflect the new database configuration.

## 4. How can I persist data in development and production?

In both development and production, Docker volumes are used to persist data. For instance, in the `compose.*.yaml` file, the `postgres-data-*` volume stores PostgreSQL data, ensuring that data is retained even if the container restarts. You can also define named volumes for other services where data persistence is essential.

## 5. What is the difference between development and production Docker configurations?

In a development environment, Docker configurations include tools that streamline coding and debugging, such as Xdebug for debugging, and volume mounts to enable real-time code updates without requiring image rebuilds.

In production, the configuration is optimized for performance, security, and efficiency. This setup uses multi-stage builds to keep the image lightweight and includes only essential tools, packages, and libraries.

It’s recommended to use `alpine`-based images in production for smaller image sizes, enhancing deployment speed and security.

Additionally, consider using [Docker Scout](/manuals/scout/_index.md) to detect and analyze vulnerabilities, especially in production environments.

For additional information about using Docker Compose in production, see [this guide](/compose/how-tos/production/).
