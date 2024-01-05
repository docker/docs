---
title: Volume
id: volume
short_description: >
  A specially-designated directory within one or more containers
  that bypasses the Union File System. Volumes are designed to persist data,
  independent of the container's life cycle. 
---

A volume is a specially-designated directory within one or more containers
  that bypasses the Union File System. Volumes are designed to persist data,
  independent of the container's life cycle. Docker therefore never automatically deletes volumes when you remove a container, nor will it "garbage collect" volumes that are no longer referenced by a container.


There are three types of volumes: *host, anonymous, and named*:

  - A **host volume** lives on the Docker host's filesystem and can be accessed from within the container.

  - A **named volume** is a volume which Docker manages where on disk the volume is created, but it is given a name.

  - An **anonymous volume** is similar to a named volume, however, it can be difficult to refer to the same volume over time when it is an anonymous volume. Docker handles where the files are stored.
