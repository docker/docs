---
title: Architecture-specific images
description: Learn how to use images that are specific to particular hardware architectures in Docker Universal Control Plane.
keywords: UCP, Docker EE, image, IBM z
---

Docker Universal Control Plane deploys images for a number of different
hardware architectures, including IBM z systems. Some architectures
require pulling images that have specific tags indicating the target
architecture.

## Tag for IBM z Systems

Append the tag `-s390x` to a UCP system image to pull the appropriate image
for IBM z Systems. For example, you can modify the CLI command for getting a
[UCP support dump](..\..\get-support.md) to use an environment variable
that indicates the current architecture:

```bash
[[ $(docker info --format='{{.Architecture}}') == s390x ]] && export _ARCH='-s390x' || export _ARCH=''

docker container run --rm \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --log-driver none \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }}${_ARCH} \
  support > docker-support.tgz
```

In this case, the environment variable is named `_ARCH`, but you can use any 
valid shell name.




