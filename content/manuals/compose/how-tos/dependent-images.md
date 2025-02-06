---
description: Build images for services with shared definition
keywords: compose, build
title: Build dependent images
weight: 50
---

{{< summary-bar feature_name="Compose dependent images" >}}

To reduce push/pull time and image weight, a common practice for Compose applications is to have services
share base layers as much as possible. You will typically select the same operating system base image for
all services. But you also can get one step further sharing image layers when your images share the same
system packages. The challenge to address is then to avoid repeating the exact same Dockerfile instruction 
in all services.

For illustration, this page assumes you want all your services to be built with an `alpine` base
image and install system package `openssl`.

## Multi-stage Dockerfile

The recommended approach is to group the shared declaration in a single Dockerfile, and use multi-stage features
so that service images are built from this shared declaration.

Dockerfile:

```dockerfile
FROM alpine as base
RUN /bin/sh -c apk add --update --no-cache openssl

FROM base as service_a
# build service a
...

FROM base as service_b
# build service b
...
```

Compose file:

```yaml
services:
  a:
     build:
       target: service_a
  b:
     build:
       target: service_b
```

## Use another service's image as the base image

A popular pattern is to reuse a service image as a base image in another service.
As Compose does not parse the Dockerfile, it can't automatically detect this dependency 
between services to correctly order the build execution.

a.Dockerfile:

```dockerfile
FROM alpine
RUN /bin/sh -c apk add --update --no-cache openssl
```

b.Dockerfile:

```dockerfile
FROM service_a
# build service b
```

Compose file:

```yaml
services:
  a:
     image: service_a 
     build:
       dockerfile: a.Dockerfile
  b:
     image: service_b
     build:
       dockerfile: b.Dockerfile
```

Legacy Docker Compose v1 used to build images sequentially, which made this pattern usable
out of the box. Compose v2 uses BuildKit to optimise builds and build images in parallel 
and requires an explicit declaration.

The recommended approach is to declare the dependent base image as an additional build context:

Compose file:

```yaml
services:
  a:
     image: service_a
     build: 
       dockerfile: a.Dockerfile
  b:
     image: service_b
     build:
       context: b/
       dockerfile: b.Dockerfile
       additional_contexts:
         # `FROM service_a` will be resolved as a dependency on service a which has to be built first
         service_a: "service:a"  
```

## Build with Bake

Using [Bake](/manuals/build/bake/_index.md) let you pass the complete build definition for all services
and to orchestrate build execution in the most efficient way. 

To enable this feature, run Compose with the `COMPOSE_BAKE=true` variable set in your environment.

```console
$ COMPOSE_BAKE=true docker compose build
[+] Building 0.0s (0/1)                                                         
 => [internal] load local bake definitions                                 0.0s
...
[+] Building 2/2 manifest list sha256:4bd2e88a262a02ddef525c381a5bdb08c83  0.0s
 ✔ service_b  Built                                                        0.7s 
 ✔ service_a  Built    
```

Bake can also be selected as the default builder by editing your `$HOME/.docker/config.json` config file:
```json
{
  ...
  "plugins": {
    "compose": {
      "build": "bake"
    }
  }
  ...
}
```