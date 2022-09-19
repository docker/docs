---
title: Official Docker samples
description: Learn how to develop and ship containerized applications, by walking through samples that exhibit canonical practices.
redirect_from:
- /en/latest/examples/
- /engine/examples/
- /examples/
---
## Awesome Compose

[https://github.com/docker/awesome-compose](https://github.com/docker/awesome-compose)

A curated repository containing over 30 Docker Compose samples. These samples offer a starting point for how to integrate different services using a Compose file and to manage their deployment with Docker Compose.

> **Important**
>
> Use the following samples as a learning tool in local development environments only. Do not deploy these samples in production environments.
{: .important}

### Compose applications with multiple integrated services

> **Note**
>
> Samples compatible with [Docker Dev Environments](../desktop/dev-environments/index.md) require Docker Desktop version 4.10 or later.

| Sample | Dev Environment (if compatible) |
| ------ | ------------------------------- |
| [ASP.NET / MS-SQL](https://github.com/docker/awesome-compose/tree/master/aspnet-mssql){: target="_blank" rel="noopener" class="_"} | - |
| [ASP.NET / NGINX / MySQL](https://github.com/docker/awesome-compose/tree/master/nginx-aspnet-mysql){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/nginx-aspnet-mysql) |
| [Elasticsearch / Logstash / Kibana](https://github.com/docker/awesome-compose/tree/master/elasticsearch-logstash-kibana){: target="_blank" rel="noopener" class="_"} | - |
| [Go / NGINX / MySQL](https://github.com/docker/awesome-compose/tree/master/nginx-golang-mysql){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/nginx-golang-mysql) |
| [Go / NGINX / PostgreSQL](https://github.com/docker/awesome-compose/tree/master/nginx-golang-postgres){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/nginx-golang-postgres) |
| [Go / NGINX](https://github.com/docker/awesome-compose/tree/master/nginx-golang){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/nginx-golang) |
| [Flask / NGINX / MongoDB](https://github.com/docker/awesome-compose/tree/master/nginx-flask-mongo) | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/nginx-flask-mongo) |
| [Flask / NGINX / MySQL](https://github.com/docker/awesome-compose/tree/master/nginx-flask-mysql){: target="_blank" rel="noopener" class="_"} | - |
| [Flask / NGINX / WSGI ](https://github.com/docker/awesome-compose/tree/master/nginx-wsgi-flask){: target="_blank" rel="noopener" class="_"} | - |
| [Flask / Redis](https://github.com/docker/awesome-compose/tree/master/flask-redis){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/flask-redis) |
| [Node.js / NGINX / Redis](https://github.com/docker/awesome-compose/tree/master/nginx-nodejs-redis){: target="_blank" rel="noopener" class="_"} | - |
| [Java Spark / MySQL](https://github.com/docker/awesome-compose/tree/master/sparkjava-mysql){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/sparkjava-mysql) |
| [PostgreSQL / pgAdmin](https://github.com/docker/awesome-compose/tree/master/postgresql-pgadmin){: target="_blank" rel="noopener" class="_"} | - |
| [React / Spring / MySQL](https://github.com/docker/awesome-compose/tree/master/react-java-mysql){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/react-java-mysql) |
| [React / Express / MySQL](https://github.com/docker/awesome-compose/tree/master/react-express-mysql){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/react-express-mysql) |
| [React / Express / MongoDB](https://github.com/docker/awesome-compose/tree/master/react-express-mongodb){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/react-express-mongodb) |
| [React / Rust / PostgreSQL](https://github.com/docker/awesome-compose/tree/master/react-rust-postgres){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/react-rust-postgres) |
| [React / NGINX ](https://github.com/docker/awesome-compose/tree/master/react-nginx){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/react-nginx) |
| [Spring / PostgreSQL](https://github.com/docker/awesome-compose/tree/master/spring-postgres){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/spring-postgres) |

### Compose applications with a single service

> **Note**
>
> Samples compatible with [Docker Dev Environments](../desktop/dev-environments/index.md) require Docker Desktop version 4.10 or later.

| Sample | Dev Environment (if compatible) |
| ------ | ------------------------------- |
| [Angular](https://github.com/docker/awesome-compose/tree/master/angular){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/angular) |
| [Django](https://github.com/docker/awesome-compose/tree/master/django){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/django) |
| [FastAPI](https://github.com/docker/awesome-compose/tree/master/fastapi){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/fastapi) |
| [Flask](https://github.com/docker/awesome-compose/tree/master/flask){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/flask) |
| [Minecraft server](https://github.com/docker/awesome-compose/tree/master/minecraft){: target="_blank" rel="noopener" class="_"} | - |
| [PHP](https://github.com/docker/awesome-compose/tree/master/apache-php){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/apache-php) |
| [Plex](https://github.com/docker/awesome-compose/tree/master/plex){: target="_blank" rel="noopener" class="_"} | - |
| [Portainer](https://github.com/docker/awesome-compose/tree/master/portainer){: target="_blank" rel="noopener" class="_"} | - |
| [Spark](https://github.com/docker/awesome-compose/tree/master/sparkjava){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/sparkjava) |
| [Traefik](https://github.com/docker/awesome-compose/tree/master/traefik-golang){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/traefik-golang) |
| [Wireguard](https://github.com/docker/awesome-compose/tree/master/wireguard){: target="_blank" rel="noopener" class="_"} | - |
| [VueJS](https://github.com/docker/awesome-compose/tree/master/vuejs){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment](https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master/vuejs) |

### Platform setups using Compose

- [Gitea / PostgreSQL](https://github.com/docker/awesome-compose/tree/master/gitea-postgres){: target="_blank" rel="noopener" class="_"}
- [Nextcloud / PostgreSQL](https://github.com/docker/awesome-compose/tree/master/nextcloud-postgres){: target="_blank" rel="noopener" class="_"}
- [Nextcloud / Redis / MariaDB](https://github.com/docker/awesome-compose/tree/master/nextcloud-redis-mariadb){: target="_blank" rel="noopener" class="_"}
- [Pi-hole / cloudflared](https://github.com/docker/awesome-compose/tree/master/pihole-cloudflared-DoH){: target="_blank" rel="noopener" class="_"}
- [Prometheus / Grafana](https://github.com/docker/awesome-compose/tree/master/prometheus-grafana){: target="_blank" rel="noopener" class="_"}
- [Wordpress / MySQL](https://github.com/docker/awesome-compose/tree/master/wordpress-mysql){: target="_blank" rel="noopener" class="_"}

## Docker samples

[https://github.com/dockersamples](https://github.com/dockersamples?q=&type=all&language=&sort=stargazers)

A collection of over 30 repositories that offer sample containerized applications, tutorials, and labs.

### Top three Docker samples

| Sample | Services | Description |
| ------ | -------- | ----------- |
| [atsea-sample-shop-app](https://github.com/dockersamples/atsea-sample-shop-app){: target="_blank" rel="noopener" class="_"} | React / Spring / PostgreSQL | A sample Java REST application. |
| [example-voting-app](https://github.com/dockersamples/example-voting-app){: target="_blank" rel="noopener" class="_"} | Python / Node.js / .NET / Java / Redis / PostgreSQL | A sample distributed application running across multiple Docker containers. |
| [k8s-wordsmith-demo](https://github.com/dockersamples/k8s-wordsmith-demo){: target="_blank" rel="noopener" class="_"} | Go / Java / PostgreSQL | A sample Wordsmith project that runs across three containers: a Postgres database, a Java REST API, and a Go web application. |
