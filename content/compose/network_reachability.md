---
description: How to connect to containers in Docker Compose
keywords: documentation, docs, docker, compose, orchestration, containers, networking
title: Network reachability between containers in Compose
---

{{< include "compose-eol.md" >}}

> **Warning**
> 
> The `service name` is the recommended method of connection.

## Connection methods

There are multiple ways to connect to other containers within Docker Compose. Each approach has its own nuances. General users only need the service name.

| Name          | Naming rules           | Containers                                   |
|---------------|------------------------|----------------------------------------------|
| service       | Required field.        | Applies to first container only.  | 
| container     | Default name is `project_name-service_name-number`. Define custom name with `container_name:` key.  | Unique per container.  |
| alias         | Define custom alias with `aliases:` key which is nested under `networks:` key.  | Applies to first container only. | 
| links         | Define custom extra aliases. A cycle is not allowed i.e. A->B->A. | Applies to all containers of that service. |
| hostname      | Default name is `$CONTAINER_ID`. Define custom name with `hostname:` key. | Default name is unique per container. Custom name applies to all containers of that service. |
| extra_hosts   | Define custom domain names. `"host.docker.internal:host-gateway"` maps Docker DNS to Host PC. Useful for contacting URLs external to the docker network. | Applies to all containers of that service. | 
| IP address    | Default or custom IP address is possible. Useful for troubleshooting connectivity. | Unique per container. | 


## Commands to verify computer name
### Services, container and alias
```console
$ docker inspect -f='{{range .NetworkSettings.Networks}}{{.Aliases}}{{end}}' $CONTAINER_ID
```
### Links
```console
$ docker inspect -f='{{range .NetworkSettings.Networks}}{{.Links}}{{end}}' $CONTAINER_ID
```
### Hostname
```console
$ cat /etc/hosts
172.20.0.3 web_hostname
```
### Extra hosts
```console
$ cat /etc/hosts
172.17.0.1 host.docker.internal
172.19.0.3 my_domain
172.20.0.3 web_hostname
```
### IP address
```console
$ docker inspect -f='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINER_ID
```

> **Note**
>
> `"hostname:"` key should not to be confused with the general computer term "hostname" that means the name of the computer. 

This is the sample setup.

```yaml
version: "3.8"
services:
  web_service: # service name
    image: redis:6.2
    container_name: web_container # container name
      networks:
      my_network:
        aliases: # alias
          - web_alias
    links: # links
      - "db_service:db_link"
    hostname: web_hostname # hostname
    extra_hosts: # extra_hosts
      - "host.docker.internal:host-gateway"
      - "my_domain:172.19.0.3"
    restart: always
  db_service:
    image: redis:6.2
    container_name: db_container
    hostname: db_hostname
    networks:
      my_network:
        aliases:
          - db_alias
networks:
  my_network:
```