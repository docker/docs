---
datafolder: engine-cli
datafile: docker_container_commit
title: docker container commit
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Creating a new image from an existing container
An existing Fedora based container has had Apache installed while running
in interactive mode with the bash shell. Apache is also running. To
create a new image run `docker ps` to find the container's ID and then run:

```bash
$ docker commit -m="Added Apache to Fedora base image" \
  -a="A D Ministrator" 98bd7fc99854 fedora/fedora_httpd:20
```

Note that only `a-z0-9-_.` are allowed when naming images from an
existing container.

### Apply specified Dockerfile instructions while committing the image
If an existing container was created without the DEBUG environment
variable set to "true", you can create a new image based on that
container by first getting the container's ID with `docker ps` and
then running:

```bash
$ docker container commit -c="ENV DEBUG true" 98bd7fc99854 debug-image
```
