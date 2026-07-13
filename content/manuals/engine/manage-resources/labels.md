---
description: Learn about labels, a tool to manage metadata on Docker objects.
keywords: labels, metadata, docker, annotations
title: Docker object labels
aliases:
  - /engine/userguide/labels-custom-metadata/
  - /config/labels-custom-metadata/
---

Labels are a mechanism for applying metadata to Docker objects, including:

- Images
- Containers
- Local daemons
- Volumes
- Networks
- Swarm nodes
- Swarm services

You can use labels to organize your images, record licensing information, annotate
relationships between containers, volumes, and networks, or in any way that makes
sense for your business or application.

## Label keys and values

A label is a key-value pair, stored as a string. You can specify multiple labels
for an object, but each key must be unique within an object. If the
same key is given multiple values, the most-recently-written value overwrites
all previous values.

### Key format recommendations

A label key is the left-hand side of the key-value pair. Keys are alphanumeric
strings which may contain periods (`.`), underscores (`_`), slashes (`/`), and hyphens (`-`). Most Docker users use
images created by other organizations, and the following guidelines help to
prevent inadvertent duplication of labels across objects, especially if you plan
to use labels as a mechanism for automation.

- Authors of third-party tools should prefix each label key with the
  reverse DNS notation of a domain they own, such as `com.example.some-label`.

- Don't use a domain in your label key without the domain owner's permission.

- The `com.docker.*`, `io.docker.*`, and `org.dockerproject.*` namespaces are
  reserved by Docker for internal use.

- Label keys should begin and end with a lower-case letter and should only
  contain lower-case alphanumeric characters, the period character (`.`), and
  the hyphen character (`-`). Consecutive periods or hyphens aren't allowed.

- The period character (`.`) separates namespace "fields". Label keys without
  namespaces are reserved for CLI use, allowing users of the CLI to interactively
  label Docker objects using shorter typing-friendly strings.

These guidelines aren't currently enforced and additional guidelines may apply
to specific use cases.

### Value guidelines

Label values can contain any data type that can be represented as a string,
including (but not limited to) JSON, XML, CSV, or YAML. The only requirement is
that the value be serialized to a string first, using a mechanism specific to
the type of structure. For instance, to serialize JSON into a string, you might
use the `JSON.stringify()` JavaScript method.

Since Docker doesn't deserialize the value, you can't treat a JSON or XML
document as a nested structure when querying or filtering by label value unless
you build this functionality into third-party tooling.

## Manage labels on objects

Each type of object with support for labels has mechanisms for adding and
managing them and using them as they relate to that type of object.

Labels on images, containers, local daemons, volumes, and networks are static
for the lifetime of the object. To change these labels you must recreate the
object. Labels on Swarm nodes and services can be updated dynamically.

### Images

Add labels to images using the [`LABEL` instruction](/reference/dockerfile.md#label) in a Dockerfile:

```dockerfile
LABEL com.example.version="1.0"
LABEL com.example.description="Web application"
```

You can also set labels at build time with the `--label` flag, without needing
a `LABEL` instruction in the Dockerfile:

```console
$ docker build --label "com.example.version=1.0" -t myapp .
```

Inspect labels on an image using `docker inspect`:

```console
$ docker inspect --format='{{json .Config.Labels}}' myapp
```

Filter images by label with [`docker image ls --filter`](/reference/cli/docker/image/ls/#filter):

```console
$ docker image ls --filter "label=com.example.version"
```

### Containers

Override or add labels when starting a container with
[`docker run --label`](/reference/cli/docker/container/run/#label):

```console
$ docker run --label "com.example.env=prod" myapp
```

Inspect labels on a container:

```console
$ docker inspect --format='{{json .Config.Labels}}' mycontainer
```

Filter containers by label with [`docker container ls --filter`](/reference/cli/docker/container/ls/#filter):

```console
$ docker container ls --filter "label=com.example.env=prod"
```

### Local Docker daemons

Add labels to the Docker daemon by passing `--label` flags when starting
`dockerd`, or by setting `"labels"` in the
[daemon configuration file](/reference/cli/dockerd.md#daemon-configuration-file):

```json
{
  "labels": ["com.example.environment=production"]
}
```

View daemon labels with `docker system info`.

### Volumes

Add labels when [creating a volume](/reference/cli/docker/volume/create/):

```console
$ docker volume create --label "com.example.purpose=database" myvolume
```

Inspect volume labels:

```console
$ docker volume inspect myvolume --format='{{json .Labels}}'
```

Filter volumes by label with [`docker volume ls --filter`](/reference/cli/docker/volume/ls/#filter):

```console
$ docker volume ls --filter "label=com.example.purpose"
```

### Networks

Add labels when [creating a network](/reference/cli/docker/network/create/):

```console
$ docker network create --label "com.example.purpose=frontend" mynetwork
```

Inspect network labels:

```console
$ docker network inspect mynetwork --format='{{json .Labels}}'
```

Filter networks by label with [`docker network ls --filter`](/reference/cli/docker/network/ls/#filter):

```console
$ docker network ls --filter "label=com.example.purpose"
```

### Swarm nodes

- [Adding or updating a Swarm node's labels](/reference/cli/docker/node/update/#label-add)
- [Filtering Swarm nodes by label](/reference/cli/docker/node/ls/#filter)

### Swarm services

- [Adding labels when creating a Swarm service](/reference/cli/docker/service/create/#label)
- [Updating a Swarm service's labels](/reference/cli/docker/service/update/)
- [Filtering Swarm services by label](/reference/cli/docker/service/ls/#filter)
