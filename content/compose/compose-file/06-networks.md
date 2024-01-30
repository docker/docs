---
title: Networks top-level elements
description: Explore all the attributes the networks top-level element can have.
keywords: compose, compose specification, networks, compose file reference
---

Networks are the layer that allow services to communicate with each other.

The top-level `networks` element lets you configure named networks that can be reused across multiple services. To use a network across multiple services, you must explicitly grant each service access by using the [networks](05-services.md) attribute within the `services` top-level element. The `networks` top-level element has additional syntax that provides more granular control.

## Examples

### Basic example

In the following example, at runtime, networks `front-tier` and `back-tier` are created and the `frontend` service
is connected to `front-tier` and `back-tier` networks.

```yml
services:
  frontend:
    image: example/webapp
    networks:
      - front-tier
      - back-tier

networks:
  front-tier:
  back-tier:
```

### Advanced example

```yml
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
    image: postgres
    networks:
      - backend

networks:
  frontend:
    # Use a custom driver
    driver: custom-driver-1
  backend:
    # Use a custom driver which takes special options
    driver: custom-driver-2
    driver_opts:
      foo: "1"
      bar: "2"
```

The advanced example shows a Compose file which defines two custom networks. The `proxy` service is isolated from the `db` service, because they do not share a network in common. Only `app` can talk to both.

## Attributes

### driver

`driver` specifies which driver should be used for this network. Compose returns an error if the
driver is not available on the platform.

```yml
networks:
  db-data:
    driver: overlay
```

Compose supports the following drivers:
`none` and `host`

- `host`: Use the host's networking stack.
- `none`: Turn off networking.

#### host or none

The syntax for using built-in networks such as `host` and `none` is different, as such networks implicitly exist outside
the scope of Compose. To use them, you must define an external network with the name `host` or `none` and
an alias that Compose can use (`hostnet` and `nonet` in the following example), then grant the service
access to that network using its alias.

```yml
services:
  web:
    networks:
      hostnet: {}

networks:
  hostnet:
    external: true
    name: host
```

```yml
services:
  web:
    ...
    networks:
      nonet: {}

networks:
  nonet:
    external: true
    name: none
```

### driver_opts

`driver_opts` specifies a list of options as key-value pairs to pass to the driver. These options are
driver-dependent. Consult the driver's documentation for more information. 

```yml
networks:
  db-data:
    driver_opts:
      foo: "bar"
      baz: 1
```

### attachable

If `attachable` is set to `true`, then standalone containers should be able to attach to this network, in addition to services.
If a standalone container attaches to the network, it can communicate with services and other standalone containers
that are also attached to the network.

```yml
networks:
  mynet1:
    driver: overlay
    attachable: true
```

### enable_ipv6

`enable_ipv6` enables IPv6 networking. For an example, see step four of [Create an IPv6 network](../../config/daemon/ipv6.md).

## external

If set to `true`:
 - `external` specifies that this network’s lifecycle is maintained outside of that of the application.
Compose doesn't attempt to create these networks, and returns an error if one doesn't exist.
 - All other attributes apart from name are irrelevant. If Compose detects any other attribute, it rejects the Compose file as invalid.

In the example below, `proxy` is the gateway to the outside world. Instead of attempting to create a network, Compose
queries the platform for an existing network simply called `outside` and connects the
`proxy` service's containers to it.

```yml

services:
  proxy:
    image: example/proxy
    networks:
      - outside
      - default
  app:
    image: example/app
    networks:
      - default

networks:
  outside:
    external: true
```

### ipam

`ipam` specifies a custom IPAM configuration. This is an object with several properties, each of which is optional:

- `driver`: Custom IPAM driver, instead of the default.
- `config`: A list with zero or more configuration elements, each containing a:
  - `subnet`: Subnet in CIDR format that represents a network segment
  - `ip_range`: Range of IPs from which to allocate container IPs
  - `gateway`: IPv4 or IPv6 gateway for the master subnet
  - `aux_addresses`: Auxiliary IPv4 or IPv6 addresses used by Network driver, as a mapping from hostname to IP
- `options`: Driver-specific options as a key-value mapping.

```yml
networks:
  mynet1:
    ipam:
      driver: default
      config:
        - subnet: 172.28.0.0/16
          ip_range: 172.28.5.0/24
          gateway: 172.28.5.254
          aux_addresses:
            host1: 172.28.1.5
            host2: 172.28.1.6
            host3: 172.28.1.7
      options:
        foo: bar
        baz: "0"
```

## internal

By default, Compose provides external connectivity to networks. `internal`, when set to `true`, allows you to
create an externally isolated network.

## labels

Add metadata to containers using `labels`. You can use either an array or a dictionary.

It is recommended that you use reverse-DNS notation to prevent labels from conflicting with those used by other software.

```yml
networks:
  mynet1:
    labels:
      com.example.description: "Financial transaction network"
      com.example.department: "Finance"
      com.example.label-with-empty-value: ""
```

```yml
networks:
  mynet1:
    labels:
      - "com.example.description=Financial transaction network"
      - "com.example.department=Finance"
      - "com.example.label-with-empty-value"
```

Compose sets `com.docker.compose.project` and `com.docker.compose.network` labels.

## name

`name` sets a custom name for the network. The name field can be used to reference networks which contain special characters.
The name is used as is and is not scoped with the project name.

```yml
networks:
  network1:
    name: my-app-net
```

It can also be used in conjunction with the `external` property to define the platform network that Compose
should retrieve, typically by using a parameter so the Compose file doesn't need to hard-code runtime specific values:

```yml
networks:
  network1:
    external: true
    name: "${NETWORK_ID}"
```

## Additional resources

For more examples, see [Networking in Compose](../networking.md).