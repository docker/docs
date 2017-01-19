---
datafolder: engine-cli
datafile: docker_image_import
title: docker image import
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Import from a remote location

    # docker image import http://example.com/exampleimage.tgz example/imagerepo

### Import from a local file

Import to docker via pipe and stdin:

    # cat exampleimage.tgz | docker image import - example/imagelocal

Import with a commit message.

    # cat exampleimage.tgz | docker image import --message "New image imported from tarball" - exampleimagelocal:new

Import to a Docker image from a local file.

    # docker image import /path/to/exampleimage.tgz


### Import from a local file and tag

Import to docker via pipe and stdin:

    # cat exampleimageV2.tgz | docker image import - example/imagelocal:V-2.0

### Import from a local directory

    # tar -c . | docker image import - exampleimagedir

### Apply specified Dockerfile instructions while importing the image
This example sets the docker image ENV variable DEBUG to true by default.

    # tar -c . | docker image import -c="ENV DEBUG true" - exampleimagedir
