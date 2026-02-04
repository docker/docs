---
title: Build a custom Docker Hardened Image
linktitle: Build a custom image
description: Learn how to write a DHI definition file and build your own Docker Hardened Image from the declarative YAML schema.
keywords: hardened images, DHI, custom image, build, yaml, security, sbom, provenance, declarative, catalog
weight: 26
---

Docker Hardened Images (DHI) are built from declarative YAML definition files
instead of traditional Dockerfiles. A single YAML file describes exactly what
goes into an image: packages, users, environment variables, entrypoint, and
metadata. The DHI build system produces a minimal, signed image with a Software
Bill of Materials (SBOM) and SLSA Build Level 3 provenance.

This page explains how to write a DHI definition file, build images locally, and
use advanced patterns such as build stages, third-party repositories, file
paths, and dev variants.

> [!IMPORTANT]
>
> You must authenticate to the Docker Hardened Images registry (`dhi.io`) to
> pull base images and build tools. Use your Docker ID credentials (the same
> username and password you use for Docker Hub) when signing in.
>
> Run `docker login dhi.io` to authenticate.

## How DHI builds differ from Dockerfiles

A Dockerfile is a sequence of imperative instructions: `RUN`, `COPY`, `FROM`.
A DHI definition file is a declarative specification. You describe the desired
state of the image, and the build system figures out how to produce it.

Every DHI definition starts with a syntax directive that tells BuildKit to use
the DHI frontend instead of the Dockerfile parser:

```yaml
# syntax=dhi.io/build:2-alpine3.23
```

The frontend version corresponds to the base distribution:

| Distribution | Syntax directive |
|-------------|-----------------|
| Alpine 3.23 | `# syntax=dhi.io/build:2-alpine3.23` |
| Debian 13 (Trixie) | `# syntax=dhi.io/build:2-debian13` |

## Explore the catalog for reference

The [DHI catalog repository](https://github.com/docker-hardened-images/catalog)
is open source under Apache 2.0 and contains every official image definition.
Studying existing definitions is the best way to learn the YAML patterns for
different image types.

The catalog follows this directory structure:

```text
catalog/
├── image/
│   ├── alpine-base/
│   │   └── alpine-3.23/
│   │       ├── 3.23.yaml
│   │       └── 3.23-dev.yaml
│   ├── nginx/
│   │   ├── alpine-3.23/
│   │   │   ├── mainline.yaml
│   │   │   └── mainline-dev.yaml
│   │   └── debian-13/
│   └── redis/
│       └── debian-13/
│           ├── 8.0.yaml
│           └── 8.0-dev.yaml
├── chart/
└── package/
```

Each image organizes its variants by distribution. A `runtime` variant is
minimal and typically runs as a non-root user. A `dev` variant adds a shell,
package manager, and development tools.

## YAML schema reference

The following tables describe the fields available in a DHI definition file.

### Required fields

Every definition must include these top-level fields:

| Field | Description |
|-------|-------------|
| `name` | Human-readable name for the image. |
| `image` | Full registry path, such as `dhi.io/my-image`. |
| `variant` | Either `runtime` or `dev`. |
| `tags` | List of image tags. |
| `platforms` | Target architectures, such as `linux/amd64` and `linux/arm64`. |
| `contents` | Package repositories and packages to install. |

### Container configuration

| Field | Description |
|-------|-------------|
| `accounts` | Users, groups, and the `run-as` user. |
| `environment` | Environment variables. |
| `entrypoint` | Container entrypoint command. |
| `cmd` | Default command arguments. |
| `work-dir` | Working directory inside the container. |
| `volumes` | Volume mount points. |
| `ports` | Exposed network ports. |
| `paths` | Directories, files, and symlinks to create. |
| `os-release` | Customizes `/etc/os-release` inside the image. |
| `annotations` | OCI image annotations such as description and license. |

### Advanced fields

| Field | Description |
|-------|-------------|
| `contents.builds` | Build stages with shell pipelines. |
| `contents.keyring` | Signing keys for third-party package repositories. |
| `contents.artifacts` | Pre-built OCI artifacts to include. |
| `contents.mappings` | Package URL (purl) mappings for SBOM accuracy. |
| `vars` | Build-time variables for templating. |
| `dates` | Release date and end-of-life date. |

## Build a minimal Alpine image

The following example shows the simplest possible definition: an Alpine base
image with a non-root user.

Create a directory for your project and add a file called `base.yaml`:

```yaml
# syntax=dhi.io/build:2-alpine3.23

name: My Base Image
image: my-registry/my-base
variant: runtime
tags:
  - "1.0"
  - "latest"
platforms:
  - linux/amd64
  - linux/arm64

contents:
  repositories:
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/main
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/community
  packages:
    - alpine-baselayout-data
    - busybox
    - ca-certificates-bundle

accounts:
  run-as: nonroot
  users:
    - name: nonroot
      uid: 65532
      gid: 65532
  groups:
    - name: nonroot
      gid: 65532
      members:
        - nonroot

os-release:
  name: Docker Hardened Images (Alpine)
  id: alpine
  version-id: "3.23"
  pretty-name: My Custom Hardened Image
  home-url: https://docker.com/products/hardened-images/

environment:
  SSL_CERT_FILE: /etc/ssl/certs/ca-certificates.crt

cmd:
  - /bin/sh
```

In this definition:

- `contents.repositories` uses full URLs to Alpine package mirrors.
- `contents.packages` lists exact Alpine package names. The package
  `alpine-baselayout-data` provides essential filesystem structure and is
  required in most Alpine-based images.
- The `accounts` section creates a `nonroot` user with UID 65532 (a common
  convention for hardened images) and sets it as the default runtime user.
- The `os-release` block customizes what appears in `/etc/os-release`.

Build the image:

```console
$ docker buildx build . -f base.yaml \
    --sbom=generator=dhi.io/scout-sbom-indexer:1 \
    --provenance=1 \
    --tag my-base:latest \
    --load
```

Verify the image by checking the running user:

```console
$ docker run --rm my-base:latest id
```

This should show `uid=65532(nonroot)`.

## Add application packages

To create a useful image, add application-specific packages to the
`contents.packages` list. The following example adds Python:

```yaml
# syntax=dhi.io/build:2-alpine3.23

name: Python Runtime
image: my-registry/my-python
variant: runtime
tags:
  - "3.13"
platforms:
  - linux/amd64
  - linux/arm64

contents:
  repositories:
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/main
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/community
  packages:
    - alpine-baselayout-data
    - ca-certificates-bundle
    - bzip2
    - expat
    - libffi
    - mpdecimal
    - musl
    - ncurses
    - openssl
    - readline
    - sqlite-libs
    - python3
    - tzdata
    - zlib

accounts:
  run-as: nonroot
  users:
    - name: nonroot
      uid: 65532
      gid: 65532
  groups:
    - name: nonroot
      gid: 65532
      members:
        - nonroot

environment:
  SSL_CERT_FILE: /etc/ssl/certs/ca-certificates.crt
  PYTHON_VERSION: "3.13"

annotations:
  org.opencontainers.image.description: A minimal Python runtime image
  org.opencontainers.image.licenses: PSF-2.0

cmd:
  - python3
```

Only include the packages your application actually needs. Fewer packages means
a smaller attack surface.

For reproducible builds, pin packages to specific versions:

```yaml
  packages:
    - python3=3.13.2-r1
```

To look up available versions, run:

```console
$ docker run --rm alpine:3.23 apk search <package-name>
```

## Build a Debian-based image

To build a Debian-based image instead of Alpine, change the syntax directive
and adjust the repository format. The following example defines a Redis image
on Debian 13 with a third-party package repository:

```yaml
# syntax=dhi.io/build:2-debian13

name: Redis 8.0.x
image: my-registry/my-redis
variant: runtime
tags:
  - "8.0"
platforms:
  - linux/amd64
  - linux/arm64

contents:
  repositories:
    - deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb trixie main
  keyring:
    - https://packages.redis.io/gpg
  packages:
    - '!libelogind0'
    - '!mawk'
    - '!original-awk'
    - base-files
    - libpcre2-8-0
    - libssl3t64
    - libstdc++6
    - libsystemd0
    - redis=6:8.0.5-1rl1~trixie1
    - redis-server=6:8.0.5-1rl1~trixie1
    - redis-tools=6:8.0.5-1rl1~trixie1
    - tini

accounts:
  run-as: nonroot
  users:
    - name: nonroot
      uid: 65532
      gid: 65532
  groups:
    - name: nonroot
      gid: 65532
      members:
        - nonroot

os-release:
  name: Docker Hardened Images (Debian)
  id: debian
  version-id: "13"
  version-codename: trixie
  pretty-name: Docker Hardened Images/Debian GNU/Linux 13 (trixie)
  home-url: https://docker.com/products/hardened-images/

work-dir: /data

environment:
  REDIS_VERSION: 8.0.5

entrypoint:
  - /usr/bin/tini
  - --

cmd:
  - redis-server
  - /etc/redis/redis.conf
  - --include
  - /etc/redis/conf.d/*.conf
```

There are several Debian-specific differences from the Alpine examples:

| Feature | Alpine | Debian |
|---------|--------|--------|
| Syntax directive | `dhi.io/build:2-alpine3.23` | `dhi.io/build:2-debian13` |
| Repository format | Plain URL | `deb [signed-by=...] <url> <suite> <component>` |

### Exclude unwanted packages

Debian packages sometimes pull in dependencies you don't need. Prefix a package
name with `!` to exclude it:

```yaml
  packages:
    - '!libelogind0'
    - '!mawk'
    - '!original-awk'
```

Wrap exclusions in quotes because YAML treats `!` as a special character.

To look up available Debian packages, run:

```console
$ docker run --rm debian:trixie apt-cache search <package-name>
```

## Create files, directories, and symlinks

Use the `paths` field to create filesystem entries without running shell
commands. This is cleaner and more auditable than using build stages for simple
file operations.

```yaml
paths:
  - type: directory
    path: /var/lib/myapp
    uid: 65532
    gid: 65532
    mode: "0755"

  - type: file
    path: /etc/myapp/config.conf
    content: |
      daemonize no
      bind 0.0.0.0
      logfile ""
    uid: 0
    gid: 0
    mode: "0555"

  - type: symlink
    path: /usr/bin/myapp-alias
    uid: 0
    gid: 0
    source: /usr/bin/myapp
```

The `paths` field supports three types:

| Type | Description |
|------|-------------|
| `directory` | Creates a directory with the specified ownership and permissions. |
| `file` | Creates a file with inline content. |
| `symlink` | Creates a symbolic link pointing to `source`. |

## Use build stages for runtime configuration

Some images need shell commands during the build to configure the application.
Use `contents.builds` to define build stages with shell pipelines. This is
conceptually similar to multi-stage Dockerfiles.

The following example configures Nginx to run as a non-root user:

```yaml
contents:
  repositories:
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/main
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/community
    - http://nginx.org/packages/mainline/alpine/v3.23/main
  keyring:
    - https://nginx.org/keys/nginx_signing.rsa.pub
  packages:
    - alpine-baselayout-data
    - busybox
    - musl-utils
    - nginx=1.29.4-r1
  builds:
    - name: nginx
      contents:
        repositories:
          - https://dl-cdn.alpinelinux.org/alpine/v3.23/main
          - https://dl-cdn.alpinelinux.org/alpine/v3.23/community
          - http://nginx.org/packages/mainline/alpine/v3.23/main
        keyring:
          - https://nginx.org/keys/nginx_signing.rsa.pub
        packages:
          - alpine-baselayout-data
          - bash
          - musl-utils
          - nginx=1.29.4-r1
      pipeline:
        - name: install
          runs: |
            set -eux -o pipefail

            ln -sf /dev/stdout /var/log/nginx/access.log
            ln -sf /dev/stderr /var/log/nginx/error.log

            sed -i "s,listen       80;,listen       8080;," \
              /etc/nginx/conf.d/default.conf

            sed -i "/user  nginx;/d" /etc/nginx/nginx.conf
            sed -i "s,pid        /run/nginx.pid;,pid        /var/run/nginx.pid;," \
              /etc/nginx/nginx.conf
            sed -i '/^http {$/a\    server_tokens off;' \
              /etc/nginx/nginx.conf

            chown -R 65532:65532 /var/cache/nginx
            chmod -R g+w /var/cache/nginx
            chown -R 65532:65532 /etc/nginx
            chmod -R g+w /etc/nginx
            chown -R 65532:65532 /run
            chown -R 65532:65532 /run/lock
            chown -R 65532:65532 /var/run
            chown -R 65532:65532 /var/log/nginx
      outputs:
        - source: /
          target: /
          uid: 0
          gid: 0
          diff: true
```

Key concepts for build stages:

| Field | Description |
|-------|-------------|
| `contents` | Each build stage has its own `contents` section. Include packages needed only during the build, such as `bash`. |
| `pipeline` | Contains named steps that run shell commands. Always start scripts with `set -eux -o pipefail`. |
| `outputs` | Copies results from the build stage into the final image. Setting `diff: true` copies only files that changed, keeping the image minimal. |

## Create a dev variant

Production images should be minimal, but developers often need additional tools
for debugging and troubleshooting. A common pattern is to maintain both a
`runtime` and a `dev` variant of each image.

```yaml
# syntax=dhi.io/build:2-alpine3.23

name: My Base Image (dev)
image: my-registry/my-base
variant: dev
tags:
  - "1.0-dev"
platforms:
  - linux/amd64
  - linux/arm64

contents:
  repositories:
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/main
    - https://dl-cdn.alpinelinux.org/alpine/v3.23/community
  packages:
    - alpine-baselayout-data
    - apk-tools
    - busybox
    - ca-certificates-bundle

accounts:
  root: true
  run-as: root
  users:
    - name: nonroot
      uid: 65532
      gid: 65532
  groups:
    - name: nonroot
      gid: 65532
      members:
        - nonroot

environment:
  SSL_CERT_FILE: /etc/ssl/certs/ca-certificates.crt

cmd:
  - /bin/sh
```

The dev variant typically adds the following:

| Feature | Runtime variant | Dev variant |
|---------|----------------|-------------|
| Package manager | Not included | `apk-tools` (Alpine) or `apt` (Debian) |
| Shell | `busybox` (Alpine) or none (Debian) | `busybox` (Alpine) or `bash` (Debian) |
| User | `nonroot` (65532) | `root` |
| Debug tools | Not included | Install via package manager at runtime |

The runtime and dev variants share the same `image` field but use different
`variant` values and tag suffixes. This keeps them in the same repository while
making it clear which variant a consumer is pulling.

> [!NOTE]
>
> Using the dev variant increases the attack surface. It is not recommended as a
> runtime for production environments.

## Expose ports and volumes

For services that listen on network ports or persist data, use the `ports` and
`volumes` fields:

```yaml
ports:
  - 8080/tcp
  - 8443/tcp

volumes:
  - /data
  - /var/log/myapp
```

Always use unprivileged ports (higher than 1024) when the container runs as a
non-root user. This avoids the need for `NET_BIND_SERVICE` capabilities.

## Build and verify an image

### Build

Build a single-platform image for local testing:

```console
$ docker buildx build . -f <my-image>.yaml \
    --sbom=generator=dhi.io/scout-sbom-indexer:1 \
    --provenance=1 \
    --tag <my-image>:latest \
    --load
```

The flags break down as follows:

| Flag | Description |
|------|-------------|
| `-f <my-image>.yaml` | Points to the DHI definition file. |
| `--sbom=generator=dhi.io/scout-sbom-indexer:1` | Generates a Software Bill of Materials. |
| `--provenance=1` | Attaches SLSA provenance attestation. |
| `--load` | Loads the built image into the local Docker image store. |

### Inspect the SBOM

View the generated Software Bill of Materials:

```console
$ docker scout sbom <my-image>:latest
```

### Scan for vulnerabilities

Check the image against known CVE databases:

```console
$ docker scout cves <my-image>:latest
```

To learn more about scanning, see [Scan Docker Hardened Images](./scan.md).

### Compare with a non-hardened image

Measure the security improvement against an equivalent non-hardened image:

```console
$ docker scout compare <my-image>:latest \
    --to <non-hardened-image>:<tag> 
```

To learn more about comparing images, see
[Compare Docker Hardened Images](./compare.md).

### Inspect with Docker Debug

Verify the entrypoint configuration and explore the image contents:

```console
$ docker debug <my-image>:latest
```

To learn more about debugging, see
[Debug a Docker Hardened Image container](./debug.md).

## Push to a registry

Tag and push the image to your container registry:

```console
$ docker tag <my-image>:latest <your-namespace>/<my-image>:latest
```

```console
$ docker push <your-namespace>/<my-image>:latest
```

Replace `<your-namespace>` with your Docker Hub username or organization
namespace.
