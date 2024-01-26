---
description: Advanced Configuration for Enhanced Container Isolation
title: Advanced configuration options
keywords: enhanced container isolation, Docker Desktop, Docker socket, bind mount, configuration
---

> **Note**
>
> This feature is available with Docker Desktop version 4.27 and later. It's currently in
> [Beta](../../../release-lifecycle.md/#beta). 

This page describes optional, advanced configurations for ECI, once ECI is enabled.

## Docker socket mount permissions 

> **Important**
>
> It does not yet work on Windows hosts when Docker Desktop configured to use WSL, but does work with Hyper-V.

By default, when ECI is enabled, Docker Desktop does not allow bind-mounting the
Docker Engine socket into containers:

```console
$ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock docker:cli
docker: Error response from daemon: enhanced container isolation: docker socket mount denied for container with image "docker.io/library/docker"; image is not in the allowed list; if you wish to allow it, configure the docker socket image list in the Docker Desktop settings.
```
This prevents malicious containers from gaining access to the Docker Engine, as
such access could allow them to perform supply chain attacks (e.g., build and
push malicious images into the organization's repositories) or similar.

However, some legitimate use cases require containers to have access to the
Docker Engine socket. For example, the popular [TestContainers](https://testcontainers.com/)
framework sometimes bind-mounts the Docker Engine socket into containers to
manage them or perform post-test cleanup.

Starting with Docker Desktop 4.27, admins can optionally configure ECI to allow
bind mounting the Docker Engine socket into containers, but in a controlled way.

This can be done via the `admin-settings.json` file, as described in
[Settings Management](../settings-management/configure.md). For example:

```json
{
  "configurationFileVersion": 2,
  "enhancedContainerIsolation": {
    "locked": true,
    "value": true,
    "dockerSocketMount": {
      "imageList": {
        "images": [
          "docker.io/localstack/localstack:*",
          "docker.io/testcontainers/ryuk:*",
          "docker:cli"
        ]
      },
      "commandList": {
        "type": "deny",
        "commands": ["push"]
      }
    }
  }
}
```

As shown above, there are two configurations for bind-mounting the Docker
socket into containers: the `imageList` and the `commandList`. These are
described below.

### Image list

The `imageList` is a list of container images that are allowed to bind-mount the
Docker socket. By default the list is empty (i.e., no containers are allowed to
bind-mount the Docker socket when ECI is enabled). However, an admin can add
images to the list, using either of these formats:

| Image Reference Format  | Description |
| :---------------------- | :---------- |
| `<image_name>[:<tag>]`  | Name of the image, with optional tag. If the tag is omitted, the `:latest` tag is used. If the tag is the wildcard `*`, then it means "any tag for that image." |
| `<image_name>@<digest>` | Name of the image, with a specific repository digest (e.g., as reported by `docker buildx imagetools inspect <image>`). This means only the image that matches that name and digest is allowed. |

The image name follows the standard convention, so it can point to any registry
and repository.

In the example above, the image list was configured with three images:

```json
"imageList": {
  "images": [
    "docker.io/localstack/localstack:*",
    "docker.io/testcontainers/ryuk:*",
    "docker:cli"
  ]
}
```

This means that containers that use either the `docker.io/localstack/localstack`
or the `docker.io/testcontainers/ryuk` image (with any tag), or the `docker:cli`
image, are allowed to bind-mount the Docker socket when ECI is enabled. Thus,
the following works:

```console
$ docker run -it -v /var/run/docker.sock:/var/run/docker.sock docker:cli sh
/ #
```

> **Tip**
>
> Be restrictive on the images you allow, as described in [Recommendations](#recommendations) below.
{ .tip }

In general, it's easier to specify the image using the tag wildcard format
(e.g., `<image-name>:*`) because then `imageList` doesn't need to be updated whenever a new version of the
image is used. Alternatively, you can use an immutable tag (e.g., `:latest`),
but it does not always work as well as the wildcard because, for example,
TestContainers uses specific versions of the image, not necessarily the latest
one.

When ECI is enabled, Docker Desktop periodically downloads the image digests
for the allowed images from the appropriate registry and stores them in
memory. Then, when a container is started with a Docker socket bind-mount,
Docker Desktop checks if the container's image digest matches one of the allowed
digests. If so, the container is allowed to start, otherwise it's blocked.

Note that due to the digest comparison mentioned in the prior paragraph, it's
not possible to bypass the Docker socket mount permissions by retagging a
disallowed image to the name of an allowed one. In other words, if a user
does:

```console
$ docker image rm <allowed_image>
$ docker tag <disallowed_image> <allowed_image>
$ docker run -v /var/run/docker.sock:/var/run/docker.sock <allowed_image>
```

then the tag operation succeeds, but the `docker run` command fails
because the image digest of the disallowed image won't match that of the allowed
ones in the repository.

### Command List

The `commandList` restricts the Docker commands that a container can issue via a
bind-mounted Docker socket when ECI is enabled. It acts as a complementary
security mechanism to the `imageList` (i.e., like a second line of defense).

For example, say the `imageList` is configured to allow
image `docker:cli` to mount the Docker socket, and a container is started with
it:

```console
$ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock sh
/ #
```

By default, this allows the container to issue any command via that Docker
socket (e.g., build and push images to the organisation's repositories), which
is generally not desirable.

To improve security, the `commandList` can be configured to restrict the
commands that the processes inside the container can issue on the bind-mounted
Docker socket. The `commandList` can be configured as a "deny" list (default) or
an "allow" list, depending on your preference.

Each command in the list is specified by its name, as reported by `docker
--help` (e.g., "ps", "build", "pull", "push", etc.) In addition, the following
command wildcards are allowed to block an entire group of commands:

| Command wildcard  | Description |
| :---------------- | :---------- |
| "container\*"     | Refers to all "docker container ..." commands |
| "image\*"         | Refers to all "docker image ..." commands |
| "volume\*"        | Refers to all "docker volume ..." commands |
| "network\*"       | Refers to all "docker network ..." commands |
| "build\*"         | Refers to all "docker build ..." commands |
| "system\*"        | Refers to all "docker system ..." commands |

For example, the following configuration blocks the `build` and `push` commands
on the Docker socket:

```json
"commandList": {
  "type": "deny",
  "commands": ["build", "push"]
}
```

Thus, if inside the container, you issue either of those commands on the
bind-mounted Docker socket, they will be blocked:

```console
/ # docker push myimage
Error response from daemon: enhanced container isolation: docker command "/v1.43/images/myimage/push?tag=latest" is blocked; if you wish to allow it, configure the docker socket command list in the Docker Desktop settings or admin-settings.
```

Similarly:

```console
/ # curl --unix-socket /var/run/docker.sock -XPOST http://localhost/v1.43/images/myimage/push?tag=latest
Error response from daemon: enhanced container isolation: docker command "/v1.43/images/myimage/push?tag=latest" is blocked; if you wish to allow it, configure the docker socket command list in the Docker Desktop settings or admin-settings.
```

Note that if the `commandList` had been configured as an "allow" list, then the
effect would be the opposite: only the listed commands would have been allowed.
Whether to configure the list as an allow or deny list depends on the use case.

### Recommendations

* Be restrictive on the list of container images for which you allow bind-mounting
  of the Docker socket (i.e., the `imageList`). Generally, only allow this for
  images that are absolutely needed and that you trust.

* Use the tag wildcard format if possible in the `imageList`
  (e.g., `<image_name>:*`), as this eliminates the need to update the
  `admin-settings.json` file due to image tag changes.

* In the `commandList`, block commands that you don't expect the container to
  execute. For example, for local testing (e.g., TestContainers), containers that bind-mount the Docker
  socket typically create / run / remove containers, volumes, and networks, but
  don't typically build images or push them into repositories (though some may
  legitimately do this). What commands to allow or block depends on the use case.

  - Note that all "docker" commands issued by the container via the bind-mounted
    Docker socket will also execute under enhanced container isolation (i.e.,
    the resulting container uses a the Linux user-namespace, sensitive system
    calls are vetted, etc.)

### Caveats and limitations

* Docker Socket Mount permissions don't yet work on Docker Desktop on Windows
  hosts with WSL (but they work on Hyper-V). Support for WSL is expected to be
  added soon.

* When Docker Desktop is restarted, it's possible that an image that is allowed
  to mount the Docker socket is unexpectedly blocked from doing so. This can
  happen when the image digest changes in the remote repository (e.g., a
  ":latest" image was updated) and the local copy of that image (e.g., from a
  prior `docker pull`) no longer matches the digest in the remote repository. In
  this case, remove the local image and pull it again (e.g., `docker rm <image>`
  and `docker pull <image>`).

* It's not possible to allow Docker socket bind-mounts on images that are not on
  a registry (e.g., images built locally and not yet pushed to a
  registry). That's because Docker Desktop pulls the digests for the allowed
  images from the registry, and then uses that to compare against the local copy
  of the image.

* The `commandList` configuration applies to all containers that are allowed to
  bind-mount the Docker socket. Therefore it can't be configured differently per
  container.

* The following commands are not yet supported in the `commandList`:

| Unsupported command  | Description |
| :------------------- | :---------- |
| compose              | Docker compose |
| dev                  | Docker dev environments |
| extension            | Manages Docker extensions |
| feedback             | Send feedback to Docker |
| init                 | Creates Docker-related starter files |
| manifest             | Manages Docker image manifests |
| plugins              | Manages plugins |
| sbom                 | View Software Bill of Materials (SBOM) |
| scan                 | Docker Scan |
| scout                | Docker Scout |
| trust                | Manage trust on Docker images |

> **Note**
>
> Docker socket mount permissions do not apply when running "true"
> Docker-in-Docker (i.e., when running the Docker Engine inside a container). In
> this case there's no bind-mount of the host's Docker socket into the
> container, and therefore no risk of the container leveraging the configuration
> and credentials of the host's Docker Engine to perform malicious activity.
> Enhanced Container Isolation is capable of running Docker-in-Docker securely,
> without giving the outer container true root permissions in the Docker Desktop
> VM.
