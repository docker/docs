This branch contains Dockerfiles and configuration files which create base
images used by the Docker docs publication process.


> **Warning**: Each time a change is pushed to this branch, all the images built
from this branch will be automatically rebuilt on Docker Cloud. This will in
turn cause all the docs archives to be rebuilt.

## Overview of creating an archive image

1.  The archive's `Dockerfile` is invoked.

2.  It is based on the `docker.github.io/docs:docs-builder` image (built by the
    [Dockerfile.builder](Dockerfile.builder) Dockerfile in the `publish-tools`
    branch). That image in turn invokes the
    `docker.github.io/docs:docs-builder-onbuild` image (built by the
    [Dockerfile.builder.onbuild](Dockerfile.builder.onbuild) Dockerfile in the
    `publish-tools` branch). Post-processing scripts included in this image.

    At the end of step 2, all the static HTML has been built and post-processing
    has been done on it.

3.  The archive's `Dockerfile` resets to the
    `docker.github.io/docs:nginx-onbuild` image (built by the
    [Dockerfile.nginx](Dockerfile.nginx.onbuild) Dockerfile in the `publish-tools`
    branch). This image contains a Nginx environment, our custom Nginx
    configuration file, and some (tiny) scripts we use for post-processing HTML.

    At the end of step 3, the static HTML from step 2 has been copied into the
    much smaller layer created by the `docker.github.io/docs:nginx-onbuild`
    image, along with the Nginx configuration. The static HTML for the archive
    is now self-browseable.

The result of these three steps is the archive Dockerfile, which is tagged as
`docker.github.io/docs:v<VER>` as set in the Dockerfile in step 1. This image
has two uses:

- It can be deployed as a standalone docs archive for that version.
- It is also incorporated into the process which builds the
[`docker.github.io/docs:docs-base`](https://github.com/docker/docker.github.io/tree/docs-base)
image. That image holds all of the archives, one per directory, and is the base
image for the documentation published on https://docs.docker.com/).

## Build all of the required images locally

All of the images are built using the auto-builder function of Docker Cloud.
To test the entire process end-to-end on your local system, you need to build
each of the required images locally and tag it appropriately:

1.  Locally build and tag all tooling images:

    ```bash
    $ git checkout publish-tools
    $ docker build -t docs/docker.github.io:docs-builder -f Dockerfile.builder .
    $ docker build -t docs/docker.github.io:docs-builder-onbuild -f Dockerfile.builder.onbuild .
    $ docker build -t docs/docker.github.io:nginx-onbuild -f Dockerfile.nginx.onbuild .
    ```

2.  For each archive branch (`v1.4` through whatever is the newest archive
    (currently `v17.09`)), build that archive branch's image. This example does
    that for the `v1.4` archive branch:

    ```bash
    $ git checkout v1.4
    $ docker build -t docs/docker.github.io:v1.4 .
    ```

    > **Note**: The archive Dockerfile looks like this (comments have been
    > removed). Each of the two `FROM` lines will use the `VER` build-time
    > argument as a parameter.
    >
    > ```Dockerfile
    > ARG VER=v1.4
    > FROM docs/docker.github.io:docs-builder-onbuild AS builder
    > FROM docs/docker.github.io:nginx-onbuild
    > ```

3.  After repeating step 2 for each archive branch, build the image for `master`:

    ```bash
    $ git checkout master
    $ docker build -t docs/docker.github.io:latest -t docker.github.io/docs:livedocs .
    ```

    The resulting image has the static HTML for each archive and for the
    contents of `master`. To test it:

    ```bash
    $ docker run --rm -it -p 4000:4000 docs/docker.github.io:latest
    ```

## When to change each file in this branch

- `Dockerfile.builder`: to update the version of Jekyll or to add or modify
  tools needed by the Jekyll environment.
- `Dockerfile.builder.onbuild`: to change the logic for building archives using
  Jekyll or post-processing the static HTML.
- contents of the `scripts` directory: To change the behavior of any of the
  individual post-processing scripts which run against the static HTML.
- `Dockerfile.nginx.onbuild`: To change the base Nginx image or to change the
  command that starts Nginx for an archive.
- `nginx-overrides.conf`: To change the Nginx configuration used by all of the
  images which serve static HTML.


