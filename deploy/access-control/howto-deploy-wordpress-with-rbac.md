---
title: Deploy simple Wordpress application with RBAC
description: Learn how to deploy a simple application and customize access to resources.
keywords: rbac, authorize, authentication, users, teams, UCP, Docker
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
---

{% if include.ui %}
{% if include.version=="ucp-3.0" %}
{% elsif include.version=="ucp-2.2" %}
{% endif %}

This tutorial explains how to create a simple application with two services,
Worpress and MySQL, and use role-based access control (RBAC) to authorize access
across the organization.

## Build the organization

Acme company wants to start a blog to better communicate with its users.

```
Acme Datacenter
├── Dev
│   └── Alex Alutin
└── DBA
    └── Brad Bhatia
```

## Deploy Wordpress with MySQL

```
version: "3.1"

services:
  db:
    image: mysql:5.7
    deploy:
      replicas: 1
      labels:
        com.docker.ucp.access.label: "/Shared/acme-blog/mysql-collection"
      restart_policy:
        condition: on-failure
        max_attempts: 3
    volumes:
      - db_data:/var/lib/mysql
    networks:
      - wordpress-net
    environment:
      MYSQL_ROOT_PASSWORD: wordpress
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: wordpress
  wordpress:
    depends_on:
      - db
    image: wordpress:latest
    deploy:
      replicas: 1
      labels:
        com.docker.ucp.access.label: "/Shared/acme-blog/wordpress-collection"
      restart_policy:
        condition: on-failure
        max_attempts: 3
    volumes:
      - wordpress_data:/var/www/html
    networks:
      - wordpress-net
    ports:
      - "8000:80"
    environment:
      WORDPRESS_DB_HOST: db:3306
      WORDPRESS_DB_PASSWORD: wordpress

volumes:
  db_data:
  wordpress_data:

networks:
  wordpress-net:
    labels:
      com.docker.ucp.access.label: "/Shared/acme-blog"
```

{% endif %}
