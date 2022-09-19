---
title: Official Docker samples
description: Learn how to develop and ship containerized applications, by walking through samples that exhibit canonical practices.
redirect_from:
- /en/latest/examples/
- /engine/examples/
- /examples/
---
{%assign awesomeComposeRepo="https://github.com/docker/awesome-compose/tree/master"%}
{%assign awesomeComposeDevEnv="https://open.docker.com/dashboard/dev-envs?url=https://github.com/docker/awesome-compose/tree/master"%}

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
| [ASP.NET / MS-SQL]({{awesomeComposeRepo}}/aspnet-mssql){: target="_blank" rel="noopener" class="_"} | - |
| [ASP.NET / NGINX / MySQL]({{awesomeComposeRepo}}/nginx-aspnet-mysql){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/nginx-aspnet-mysql) |
| [Elasticsearch / Logstash / Kibana]({{awesomeComposeRepo}}/elasticsearch-logstash-kibana){: target="_blank" rel="noopener" class="_"} | - |
| [Go / NGINX / MySQL]({{awesomeComposeRepo}}/nginx-golang-mysql){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/nginx-golang-mysql) |
| [Go / NGINX / PostgreSQL]({{awesomeComposeRepo}}/nginx-golang-postgres){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/nginx-golang-postgres) |
| [Go / NGINX]({{awesomeComposeRepo}}/nginx-golang){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/nginx-golang) |
| [Flask / NGINX / MongoDB]({{awesomeComposeRepo}}/nginx-flask-mongo) | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/nginx-flask-mongo) |
| [Flask / NGINX / MySQL]({{awesomeComposeRepo}}/nginx-flask-mysql){: target="_blank" rel="noopener" class="_"} | - |
| [Flask / NGINX / WSGI ]({{awesomeComposeRepo}}/nginx-wsgi-flask){: target="_blank" rel="noopener" class="_"} | - |
| [Flask / Redis]({{awesomeComposeRepo}}/flask-redis){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/flask-redis) |
| [Node.js / NGINX / Redis]({{awesomeComposeRepo}}/nginx-nodejs-redis){: target="_blank" rel="noopener" class="_"} | - |
| [Java Spark / MySQL]({{awesomeComposeRepo}}/sparkjava-mysql){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/sparkjava-mysql) |
| [PostgreSQL / pgAdmin]({{awesomeComposeRepo}}/postgresql-pgadmin){: target="_blank" rel="noopener" class="_"} | - |
| [React / Spring / MySQL]({{awesomeComposeRepo}}/react-java-mysql){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/react-java-mysql) |
| [React / Express / MySQL]({{awesomeComposeRepo}}/react-express-mysql){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/react-express-mysql) |
| [React / Express / MongoDB]({{awesomeComposeRepo}}/react-express-mongodb){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/react-express-mongodb) |
| [React / Rust / PostgreSQL]({{awesomeComposeRepo}}/react-rust-postgres){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/react-rust-postgres) |
| [React / NGINX ]({{awesomeComposeRepo}}/react-nginx){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/react-nginx) |
| [Spring / PostgreSQL]({{awesomeComposeRepo}}/spring-postgres){: target="_blank" rel="noopener" class="_"}  | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/spring-postgres) |

### Compose applications with a single service

> **Note**
>
> Samples compatible with [Docker Dev Environments](../desktop/dev-environments/index.md) require Docker Desktop version 4.10 or later.

| Sample | Dev Environment (if compatible) |
| ------ | ------------------------------- |
| [Angular]({{awesomeComposeRepo}}/angular){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/angular) |
| [Django]({{awesomeComposeRepo}}/django){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/django) |
| [FastAPI]({{awesomeComposeRepo}}/fastapi){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/fastapi) |
| [Flask]({{awesomeComposeRepo}}/flask){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/flask) |
| [Minecraft server]({{awesomeComposeRepo}}/minecraft){: target="_blank" rel="noopener" class="_"} | - |
| [PHP]({{awesomeComposeRepo}}/apache-php){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/apache-php) |
| [Plex]({{awesomeComposeRepo}}/plex){: target="_blank" rel="noopener" class="_"} | - |
| [Portainer]({{awesomeComposeRepo}}/portainer){: target="_blank" rel="noopener" class="_"} | - |
| [Spark]({{awesomeComposeRepo}}/sparkjava){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/sparkjava) |
| [Traefik]({{awesomeComposeRepo}}/traefik-golang){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/traefik-golang) |
| [Wireguard]({{awesomeComposeRepo}}/wireguard){: target="_blank" rel="noopener" class="_"} | - |
| [VueJS]({{awesomeComposeRepo}}/vuejs){: target="_blank" rel="noopener" class="_"} | [Open in Docker Dev Environment]({{awesomeComposeDevEnv}}/vuejs) |

### Platform setups using Compose

- [Gitea / PostgreSQL]({{awesomeComposeRepo}}/gitea-postgres){: target="_blank" rel="noopener" class="_"}
- [Nextcloud / PostgreSQL]({{awesomeComposeRepo}}/nextcloud-postgres){: target="_blank" rel="noopener" class="_"}
- [Nextcloud / Redis / MariaDB]({{awesomeComposeRepo}}/nextcloud-redis-mariadb){: target="_blank" rel="noopener" class="_"}
- [Pi-hole / cloudflared]({{awesomeComposeRepo}}/pihole-cloudflared-DoH){: target="_blank" rel="noopener" class="_"}
- [Prometheus / Grafana]({{awesomeComposeRepo}}/prometheus-grafana){: target="_blank" rel="noopener" class="_"}
- [Wordpress / MySQL]({{awesomeComposeRepo}}/wordpress-mysql){: target="_blank" rel="noopener" class="_"}

## Docker samples

[https://github.com/dockersamples](https://github.com/dockersamples?q=&type=all&language=&sort=stargazers)

A collection of over 30 repositories that offer sample containerized applications, tutorials, and labs.

### Top three Docker samples

| Sample | Services | Description |
| ------ | -------- | ----------- |
| [atsea-sample-shop-app](https://github.com/dockersamples/atsea-sample-shop-app){: target="_blank" rel="noopener" class="_"} | React / Spring / PostgreSQL | A sample Java REST application. |
| [example-voting-app](https://github.com/dockersamples/example-voting-app){: target="_blank" rel="noopener" class="_"} | Python / Node.js / .NET / Java / Redis / PostgreSQL | A sample distributed application running across multiple Docker containers. |
| [k8s-wordsmith-demo](https://github.com/dockersamples/k8s-wordsmith-demo){: target="_blank" rel="noopener" class="_"} | Go / Java / PostgreSQL | A sample Wordsmith project that runs across three containers: a Postgres database, a Java REST API, and a Go web application. |
