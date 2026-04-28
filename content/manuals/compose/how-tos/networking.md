---
description: How Docker Compose sets up networking between containers
keywords: documentation, docs, docker, compose, orchestration, containers, networking
title: Networking in Compose
linkTitle: Networking
weight: 70
aliases:
  - /compose/networking/
---

By default Compose sets up a single
[network](/reference/cli/docker/network/create/) for your app. Each
container for a service joins the default network and is both reachable by
other containers on that network, and discoverable by the service's name.

For most development setups the default network is sufficient. You'll want to define custom networks when you need to isolate services from each other, or when connecting to infrastructure that's managed outside of Compose.

## Default network and service discovery

When you run `docker compose up`, Compose creates a network named `<project-name>_default` and attaches all services to it. Each service registers its name with an internal DNS server, so containers can reach each other using the service name directly. No IP addresses or manual configuration is needed.

For example, suppose your app is in a directory called `myapp`, and your `compose.yaml` looks like this:

```yaml
services:
  web:
    build: .
    ports:
      - "8000:8000"
  db:
    image: postgres:18
    ports:
      - "8001:5432"
```

When you run `docker compose up`, the following happens:

1.  A network called `myapp_default` is created.
2.  A container is created using `web`'s configuration. It joins the network
    `myapp_default` under the name `web`.
3.  A container is created using `db`'s configuration. It joins the network
    `myapp_default` under the name `db`.

Each container can now look up the service name `web` or `db` and
get back the appropriate container's IP address. For example, `web`'s
application code could connect to the URL `postgres://db:5432` and start
using the Postgres database.

> [!TIP]
>
> Reference containers by name, not IP, whenever possible. Otherwise you’ll need to constantly update the IP address you use.

It is important to note the distinction between `HOST_PORT` and `CONTAINER_PORT`.
In the above example, for `db`, the `HOST_PORT` is `8001` and the container port is
`5432` (Postgres default). Networked service-to-service
communication uses the `CONTAINER_PORT`. The host port only comes into play when accessing the service from outside the network.

Within the `web` container, your connection string to `db` would look like
`postgres://db:5432`, and from the host machine, the connection string would
look like `postgres://{DOCKER_IP}:8001` for example `postgres://localhost:8001` if your container is running locally.

> [!NOTE]
>
> Your app's network is given a name based on the "project name",
> which is based on the name of the directory it lives in. You can override the
> project name with either the [`--project-name` flag](/reference/cli/docker/compose/)
> or the [`COMPOSE_PROJECT_NAME` environment variable](environment-variables/envvars.md#compose_project_name).

### Updating containers on the network

If you make a configuration change to a service and run `docker compose up` to update it, the old container is removed and the new one joins the network under a different IP address but the same name. Running containers can look up that name and connect to the new address, but the old address stops working.

If any containers have connections open to the old container, they are closed. It is a container's responsibility to detect this condition, look up the name again and reconnect.

## Use an existing external network

If you've manually created a bridge network outside of Compose using the `docker network create` command, you can connect your Compose services to it by marking the network as [`external`](/reference/compose-file/networks.md#external).

```yaml
services:
  # ...
networks:
  network1:
    name: my-pre-existing-network
    external: true
```

Instead of attempting to create a network called `<project-name>_default`, Compose looks for a network called `my-pre-existing-network` and connects your app's containers to it.

External networks are particularly useful when services in separate Compose projects need to communicate. Create a shared network once, then reference it as external in each project. Services on the same external network can reach each other by service name, just like services within a single project.

Hybrid networking
A service can connect to both an external shared network and its own project-internal network. This is useful when you want a service to be reachable across projects while keeping other services — like a database — fully isolated:
yamlservices:
  api:
    image: myapp-api
    networks:
      - shared    # Reachable from other projects
      - internal  # Can reach the database

  database:
    image: postgres:18
    networks:
      - internal  # Not exposed on the shared network

networks:
  shared:
    name: shared-network
    external: true
  internal: {}    # Project-specific, isolated

Troubleshooting
If a service can't reach another service on an external network, verify that both containers are actually attached to it:
bashdocker network inspect shared-network

## Specify custom networks

Instead of just using the default app network, you can specify your own networks with the top-level `networks` key. This lets you create more complex topologies and specify [custom network drivers](/engine/extend/plugins_network/) and options. You can also use it to connect services to externally-created networks which aren't managed by Compose.

Each service can specify what networks to connect to with the service-level `networks` key, which is a list of names referencing entries under the top-level `networks` key.

The following example shows a Compose file which defines two custom networks. The `proxy` service is isolated from the `db` service, because they do not share a network in common. Only `app` can talk to both.

```yaml
services:
  proxy:
    build: ./proxy
    networks:
      - frontend
  app:
    build: ./app
    networks:
      - frontend
      - backend
  db:
    image: postgres:18
    networks:
      - backend

networks:
  frontend:
    # Specify driver options
    driver: bridge
    driver_opts:
      com.docker.network.bridge.host_binding_ipv4: "127.0.0.1"
  backend:
    # Use a custom driver
    driver: custom-driver
```

Networks can be configured with static IP addresses by setting the [ipv4_address and/or ipv6_address](/reference/compose-file/services.md#ipv4_address-ipv6_address) for each attached network.

Networks can also be given a [custom name](/reference/compose-file/networks.md#name):

```yaml
services:
  # ...
networks:
  frontend:
    name: custom_frontend
    driver: custom-driver-1
```

## Configure the default network

Instead of, or as well as, specifying your own networks, you can also change the settings of the app-wide default network by defining an entry under `networks` named `default`:

```yaml
services:
  web:
    build: .
    ports:
      - "8000:8000"
  db:
    image: postgres:18

networks:
  default:
    # Use a custom driver
    driver: custom-driver-1
```

## Custom DNS with `extra_hosts`

You can add custom hostname-to-IP mappings to a container's /etc/hosts file using extra_hosts. This is useful when a service needs to resolve a hostname that isn't registered in the default DNS — for example, a staging API endpoint or a service on a fixed IP:
yamlservices:
  app:
    image: myapp
    extra_hosts:
      - "api.staging:192.168.1.100"
      - "cache.internal:192.168.1.101"
To map a hostname to the host machine's IP dynamically, use the special host-gateway value:
yamlservices:
  app:
    image: myapp
    extra_hosts:
      - "host.docker.internal:host-gateway"

## Multi-host networking

When deploying a Compose application on a Docker Engine with [Swarm mode enabled](/manuals/engine/swarm/_index.md),
you can make use of the built-in `overlay` driver to enable multi-host communication.

Overlay networks are always created as `attachable`. You can optionally set the [`attachable`](/reference/compose-file/networks.md#attachable) property to `false`.

See the [overlay network driver documentation](/manuals/engine/network/drivers/overlay.md)
to learn about multi-host overlay networks.

## Link containers

Links allow you to define extra aliases by which a service is reachable from another service. They are not required to enable services to communicate. By default, any service can reach any other service at that service's name. In the following example, `db` is reachable from `web` at the hostnames `db` and `database`:

```yaml
services:
  web:
    build: .
    links:
      - "db:database"
  db:
    image: postgres:18
```

See the [links reference](/reference/compose-file/services.md#links) for more information.

## Debugging

To find out which host port maps to a container port, use docker compose port:

bash# Which host port maps to container port 5432 on db?
docker compose port db 5432
# Output: 0.0.0.0:8001
This is particularly useful when you use dynamic port mapping and the host port changes on every docker compose up:
yamlservices:
  web:
    image: nginx
    ports:
      - "80"  # Docker assigns the host port dynamically
bashdocker compose port web 80
# Output: 0.0.0.0:55432

## Further reference information

For full details of the network configuration options available, see the following references:

- [Top-level `networks` element](/reference/compose-file/networks.md)
- [Service-level `networks` attribute](/reference/compose-file/services.md#networks)
