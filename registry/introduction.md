---
description: Explains what the Registry is, basic use cases and requirements
keywords: registry, on-prem, images, tags, repository, distribution, use cases, requirements
title: About Registry
---

{% include registry.md %}

A registry is a storage and content delivery system, holding named Docker
images, available in different tagged versions. For example, the image `distribution/registry`, with tags `2.0` and `2.1`.

Users interact with a registry by using `docker push` and `docker pull` commands. See the [example](#example-registry-interaction) to try the commands yourself.

Storage itself is delegated to drivers. The default storage driver is the local
posix filesystem, which is suitable for development or small deployments.
Additional cloud-based storage drivers like S3, Microsoft Azure, OpenStack Swift,
and Aliyun OSS are also supported. If you're looking into using other storage
backends, you may do so by writing your own driver implementing the
[Storage API](storage-drivers/index.md).

Since securing access to your hosted images is paramount, the Registry natively
supports TLS and basic authentication.

The Registry is a stateless, highly scalable server side application that stores
and distributes Docker images. The Registry is open-source, under the
permissive [Apache license](https://en.wikipedia.org/wiki/Apache_License).
You can find the source code on
[GitHub](https://github.com/distribution/distribution){: target="blank" rel="noopener" class=""}.

The Registry GitHub repository includes additional information about advanced
authentication and authorization methods. Only very large or public deployments
are expected to extend the Registry in this way.

The Registry ships with a robust [notification system](notifications.md),
calling webhooks in response to activity, and both extensive logging and reporting,
mostly useful for large installations that want to collect metrics.

## Why host your own registry?

You should host your own registry if you want to:

 * control where you store your images
 * fully own the distribution pipeline of your images
 * integrate image storage and distribution into your in-house development workflow


## Understanding image naming

Image names as used in typical Docker commands reflect their origin:

 * `docker pull ubuntu` instructs Docker to pull an image named `ubuntu` from the Docker Hub. This is simply a shortcut for the longer `docker pull docker.io/library/ubuntu` command.
 * `docker pull myregistrydomain:port/foo/bar` instructs Docker to contact the registry located at `myregistrydomain:port` to find the image `foo/bar`.

You can find out more about the various Docker commands in
the [Command-line reference](../engine/reference/commandline/cli.md).

## Use cases

Running your own registry is a great solution to integrate with and complement
your CI/CD system. In a typical workflow, a commit to your source revision
control system would trigger a build on your CI system, which would then push a
new image to your registry if the build is successful. A notification from your
registry would then trigger a deployment on a staging environment, or notify
other systems that a new image is available.

It's also an essential component if you want to quickly deploy a new image over
a large cluster of machines.

Finally, it's the best way to distribute images inside an isolated network.

## Requirements

The Registry is compatible with Docker Engine version 1.6.0 or later.

To host your own registry, you need to be familiar with Docker, specifically with regard to
pushing and pulling images. You must understand the difference between the
daemon and the cli, and at least grasp basic concepts about networking.

Also, while just starting a registry is fairly easy, operating it in a
production environment requires operational skills, just like any other service.
You are expected to be familiar with systems availability and scalability,
logging and log processing, systems monitoring, and security 101. Strong
understanding of HTTP and overall network communications, plus familiarity with
Go are certainly useful as well for advanced operations or hacking.

## Example registry interaction

The following steps provide a simple example of interacting with a locally hosted registry as well as Docker Hub.

1. Start your registry.
   ```console
   $ docker run -d -p 5000:5000 --name registry registry:2
   ```

2. Pull the image from Docker Hub.
   ```console
   $ docker pull ubuntu
   ```

3. Tag the image so that it points to your registry.
   ```console
   $ docker image tag ubuntu localhost:5000/myfirstimage
   ```

4. Push the image to your registry
   ```console
   $ docker push localhost:5000/myfirstimage
```

5. Pull the image from your registry.
   ```console
   $ docker pull localhost:5000/myfirstimage
   ```

6. Stop your registry and remove all data.
   ```console
   $ docker container stop registry && docker container rm -v registry
   ```

## Next

Dive into [deploying your registry](deploying.md).
