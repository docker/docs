---
description: Learn about labels, a tool to manage metadata on Docker objects.
keywords: labels, metadata, docker, annotations
title: Docker object labels
aliases:
  - /engine/userguide/labels-custom-metadata/
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
strings which may contain periods (`.`), underscores (`_`), and hyphens (`-`). Most Docker users use
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
managing them and using them as they relate to that type of object. These links
provide a good place to start learning about how you can use labels in your
Docker deployments.

Labels on images, containers, local daemons, volumes, and networks are static for
the lifetime of the object. To change these labels you must recreate the object.
Labels on Swarm nodes and services can be updated dynamically.

- Images and containers

  - [Adding labels to images](../reference/dockerfile.md#label)
  - [Overriding a container's labels at runtime](../reference/cli/docker/container/run.md#label)
  - [Inspecting labels on images or containers](../reference/cli/docker/inspect.md)
  - [Filtering images by label](../reference/cli/docker/image/ls.md#filter)
  - [Filtering containers by label](../reference/cli/docker/container/ls.md#filter)

- Local Docker daemons

  - [Adding labels to a Docker daemon at runtime](../reference/cli/dockerd.md)
  - [Inspecting a Docker daemon's labels](../reference/cli/docker/system/info.md)

- Volumes

  - [Adding labels to volumes](../reference/cli/docker/volume/create.md)
  - [Inspecting a volume's labels](../reference/cli/docker/volume/inspect.md)
  - [Filtering volumes by label](../reference/cli/docker/volume/ls.md#filter)

- Networks

  - [Adding labels to a network](../reference/cli/docker/network/create.md)
  - [Inspecting a network's labels](../reference/cli/docker/network/inspect.md)
  - [Filtering networks by label](../reference/cli/docker/network/ls.md#filter)

- Swarm nodes

  - [Adding or updating a Swarm node's labels](../reference/cli/docker/node/update.md#label-add)
  - [Inspecting a Swarm node's labels](../reference/cli/docker/node/inspect.md)
  - [Filtering Swarm nodes by label](../reference/cli/docker/node/ls.md#filter)

- Swarm services
  - [Adding labels when creating a Swarm service](../reference/cli/docker/service/create.md#label)
  - [Updating a Swarm service's labels](../reference/cli/docker/service/update.md)
  - [Inspecting a Swarm service's labels](../reference/cli/docker/service/inspect.md)
  - [Filtering Swarm services by label](../reference/cli/docker/service/ls.md#filter)
