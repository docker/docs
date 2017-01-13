---
datafolder: engine-cli
datafile: docker_image_tag
title: docker image tag
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Tagging an image referenced by ID

To tag a local image with ID "0e5574283393" into the "fedora" repository with
"version1.0":

    docker image tag 0e5574283393 fedora/httpd:version1.0

### Tagging an image referenced by Name

To tag a local image with name "httpd" into the "fedora" repository with
"version1.0":

    docker image tag httpd fedora/httpd:version1.0

Note that since the tag name is not specified, the alias is created for an
existing local version `httpd:latest`.

### Tagging an image referenced by Name and Tag

To tag a local image with name "httpd" and tag "test" into the "fedora"
repository with "version1.0.test":

    docker image tag httpd:test fedora/httpd:version1.0.test

### Tagging an image for a private repository

To push an image to a private registry and not the central Docker
registry you must tag it with the registry hostname and port (if needed).

    docker image tag 0e5574283393 myregistryhost:5000/fedora/httpd:version1.0
