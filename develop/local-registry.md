---
title: Run a local registry
description: >
  Easily run a local image registry using the official registry image,
  maintained by CNCF.
keywords: registry, distribution, local
---

This guide shows you how to run a local registry on your development machine.

## What's a registry?

A registry is where you store the container images you've built. It's also a
platform for downloading and distributing images.
[Docker Hub](../docker-hub/index.md) is one example of a hosted registry
service.

It's possible to run your own, local registry as well. In most cases, you'll
want to use a hosted registry service for managing your images. But there are
scenarios where you may need to use a self-hosted registry instead. For example,
if you need to have control where you store your images, or you want to fully
own your images distribution pipeline.

## Step one: Start the registry

You can use the `registry` official image to run a container registry on your
local machine.

Run the following command to start the registry service:

```console
$ docker run -d -p 5000:5000 --restart always --name registry registry:2
```

If you now run `docker ps` you should see the registry container running:

```console
$ docker ps
CONTAINER ID   IMAGE        COMMAND                  CREATED        STATUS        PORTS                    NAMES
a71cdb4c020e   registry:2   "/entrypoint.sh /etc…"   1 second ago   Up 1 second   0.0.0.0:5000->5000/tcp   registry
```

To verify that everything works, try using the registry API. Since you haven't
uploaded any images yet, it returns an empty list of repositories.

```console
$ curl localhost:5000/v2/_catalog
{"repositories":[]}
```

## Step two: Push an image

Now push an image to the registry. You can pull any existing image from Docker
Hub, change the image tag, and upload it to your local registry.

The Docker Buildx CLI provides a shorthand command for doing this:
`docker buildx imagetools create`. Run the following command, which takes the
latest version of the Alpine image and re-uploads it to your local registry:

```console
docker buildx imagetools create alpine --tag localhost:5000/alpine
```

Now try invoking the API again, and you should see the `alpine` repository:

```console
$ curl localhost:5000/v2/_catalog
{"repositories":["alpine"]}
```

## Step three: Pull the image

Now that you've upload an image, try pulling it to your Docker client.

```console
$ docker pull localhost:5000/alpine
```

The `localhost:5000` part of the image name tells Docker to pull the `alpine`
image from your local registry. By default, if you leave out the hostname of an
image, Docker tries to pull the image from Docker Hub:

- `docker pull ubuntu` instructs Docker to pull an image named `ubuntu` from the
  official Docker Hub. This is a shortcut for the longer
  `docker pull docker.io/library/ubuntu` command.
- `docker pull example.com:1337/ubuntu` instructs Docker to pull the `ubuntu`
  image from the registry located at `example.com:1337`.

## Step four: Clean up

Stop the registry container:

```console
$ docker stop registry
```

Remove the data volume created by the registry:

```console
$ docker rm -v registry
```

## Next steps

You can find the full documentation and more information on the registry image
in the
[distribution/distrubution GitHub repository](https://github.com/distribution/distribution){:
target="blank" rel="noopener"}.
