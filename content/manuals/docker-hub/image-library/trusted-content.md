---
description: Learn about Docker Hub's trusted content.
keywords: Docker Hub, Hub, trusted content
title: Trusted content
weight: 15
aliases:
- /trusted-content/official-images/using/
- /trusted-content/official-images/
---

Docker Hub's trusted content provides a curated selection of high-quality,
secure images designed to give developers confidence in the reliability and
security of the resources they use. These images are stable, regularly updated,
and adhere to industry best practices, making them a strong foundation for
building and deploying applications. Docker Hub's trusted content includes,
Docker Official Images, Verified Publisher images, and Docker-Sponsored Open
Source Software images.

## Docker Official Images

The Docker Official Images are a curated set of Docker repositories hosted on
Docker Hub.

Docker recommends you use the Docker Official Images in your projects. These
images have clear documentation, promote best practices, and are regularly
updated. Docker Official Images support most common use cases, making them
perfect for new Docker users. Advanced users can benefit from more specialized
image variants as well as review Docker Official Images as part of your
`Dockerfile` learning process.

> [!NOTE]
>
> Use of Docker Official Images is subject to [Docker's Terms of Service](https://www.docker.com/legal/docker-terms-service/).

These images provide essential base repositories that serve as the starting
point for the majority of users.

These include operating systems such as
[Ubuntu](https://hub.docker.com/_/ubuntu/) and
[Alpine](https://hub.docker.com/_/alpine/), programming language runtimes such as
[Python](https://hub.docker.com/_/python) and
[Node](https://hub.docker.com/_/node), and other essential tools such as
[memcached](https://hub.docker.com/_/memcached) and
[MySQL](https://hub.docker.com/_/mysql).

The images are some of the [most secure images](https://www.docker.com/blog/enhancing-security-and-transparency-with-docker-official-images/)
on Docker Hub. This is particularly important as Docker Official Images are
some of the most popular on Docker Hub. Typically, Docker Official images have
few or no packages containing CVEs.

The images exemplify [Dockerfile best practices](/manuals/build/building/best-practices.md)
and provide clear documentation to serve as a reference for other Dockerfile authors.

Images that are part of this program have a special badge on Docker Hub making
it easier for you to identify projects that are part of Docker Official Images.

![Docker official image badge](../images/official-image-badge-iso.png)

### Supported tags and respective Dockerfile links

The repository description for each Docker Official Image contains a
**Supported tags and respective Dockerfile links** section that lists all the
current tags with links to the Dockerfiles that created the image with those
tags. The purpose of this section is to show what image variants are available.

![Example: supported tags for Ubuntu](../images/supported_tags.webp)

Tags listed on the same line all refer to the same underlying image. Multiple
tags can point to the same image. For example, in the previous screenshot taken
from the `ubuntu` Docker Official Images repository, the tags `24.04`,
`noble-20240225`, `noble`, and `devel` all refer to the same image.

The `latest` tag for a Docker Official Image is often optimized for ease of use
and includes a wide variety of useful software, such as developer and build tools.
By tagging an image as `latest`, the image maintainers are essentially suggesting
that image be used as the default. In other words, if you do not know what tag to
use or are unfamiliar with the underlying software, you should probably start with
the `latest` image. As your understanding of the software and image variants advances,
you may find other image variants better suit your needs.

### Slim images

A number of language stacks such as
[Node.js](https://hub.docker.com/_/node/),
[Python](https://hub.docker.com/_/python/), and
[Ruby](https://hub.docker.com/_/ruby/) have `slim` tag variants
designed to provide a lightweight, production-ready base image
with fewer packages.

A typical consumption pattern for `slim`
images is as the base image for the final stage of a
[multi-staged build](https://docs.docker.com/build/building/multi-stage/).
For example, you build your application in the first stage of the build
using the `latest` variant and then copy your application into the final
stage based upon the `slim` variant. Here is an example `Dockerfile`.

```dockerfile
FROM node:latest AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . ./
FROM node:slim
WORKDIR /app
COPY --from=build /app /app
CMD ["node", "app.js"]
```

### Alpine images

Many Docker Official Images repositories also offer `alpine` variants. These
images are built on top of the [Alpine Linux](https://www.alpinelinux.org/)
distribution rather than Debian or Ubuntu. Alpine Linux is focused on providing
a small, simple, and secure base for container images, and Docker Official
Images `alpine` variants typically aim to install only necessary packages. As a
result, Docker Official Images `alpine` variants are typically even smaller
than `slim` variants.

The main caveat to note is that Alpine Linux uses [musl libc](https://musl.libc.org/)
instead of [glibc](https://www.gnu.org/software/libc/). Additionally, to
minimize image size, it's uncommon for Alpine-based images to include tools
such as Git or Bash by default. Depending on the depth of libc requirements or
assumptions in your programs, you may find yourself running into issues due to
missing libraries or tools.

When you use Alpine images as a base, consider the following options in order
to make your program compatible with Alpine Linux and musl:

- Compile your program against musl libc
- Statically link glibc libraries into your program
- Avoid C dependencies altogether (for example, build Go programs without CGO)
- Add the software you need yourself in your Dockerfile.

Refer to the `alpine` image [description](https://hub.docker.com/_/alpine) on
Docker Hub for examples on how to install packages if you are unfamiliar.

### Codenames

Tags with words that look like Toy Story characters (for example, `bookworm`,
`bullseye`, and `trixie`) or adjectives (such as `focal`, `jammy`, and
`noble`), indicate the codename of the Linux distribution they use as a base
image. Debian release codenames are [based on Toy Story characters](https://en.wikipedia.org/wiki/Debian_version_history#Naming_convention),
and Ubuntu's take the form of "Adjective Animal". For example, the
codename for Ubuntu 24.04 is "Noble Numbat".

Linux distribution indicators are helpful because many Docker Official Images
provide variants built upon multiple underlying distribution versions (for
example, `postgres:bookworm` and `postgres:bullseye`).

### Other tags

Docker Official Images tags may contain other hints to the purpose of
their image variant in addition to those described here. Often these
tag variants are explained in the Docker Official Images repository
documentation. Reading through the "How to use this image" and
"Image Variants" sections will help you to understand how to use these
variants.

## Verified Publisher images

The Docker Verified Publisher program provides high-quality images from
commercial publishers verified by Docker.

These images help development teams build secure software supply chains,
minimizing exposure to malicious content early in the process to save time and
money later.

Images that are part of this program have a special badge on Docker Hub making
it easier for users to identify projects that Docker has verified as
high-quality commercial publishers.

![Docker-Sponsored Open Source badge](../images/verified-publisher-badge-iso.png)

## Docker-Sponsored Open Source Software images

The Docker-Sponsored Open Source Software (OSS) program provides images that are
published and maintained by open-source projects sponsored by Docker.

Images that are part of this program have a special badge on Docker Hub making
it easier for users to identify projects that Docker has verified as trusted,
secure, and active open-source projects.

![Docker-Sponsored Open Source badge](../images/sponsored-badge-iso.png)