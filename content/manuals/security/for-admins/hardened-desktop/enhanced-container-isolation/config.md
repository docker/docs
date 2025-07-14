---
description: Advanced Configuration for Enhanced Container Isolation
title: Advanced configuration options for ECI
linkTitle: Advanced configuration
keywords: enhanced container isolation, Docker Desktop, Docker socket, bind mount, configuration
aliases:
 - /desktop/hardened-desktop/enhanced-container-isolation/config/
weight: 30
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

## Docker socket mount permissions

By default, when Enhanced Container Isolation (ECI) is enabled, Docker Desktop does not allow bind-mounting the
Docker Engine socket into containers:

```console
$ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock docker:cli
docker: Error response from daemon: enhanced container isolation: docker socket mount denied for container with image "docker.io/library/docker"; image is not in the allowed list; if you wish to allow it, configure the docker socket image list in the Docker Desktop settings.
```
This prevents malicious containers from gaining access to the Docker Engine, as
such access could allow them to perform supply chain attacks. For example, build and
push malicious images into the organization's repositories or similar.

However, some legitimate use cases require containers to have access to the
Docker Engine socket. For example, the popular [Testcontainers](https://testcontainers.com/)
framework sometimes bind-mounts the Docker Engine socket into containers to
manage them or perform post-test cleanup. Similarly, some Buildpack frameworks,
for example [Paketo](https://paketo.io/), require Docker socket bind-mounts into
containers.

Administrators can optionally configure ECI to allow
bind mounting the Docker Engine socket into containers, but in a controlled way.

This can be done via the Docker Socket mount permissions section in the
[`admin-settings.json`](../settings-management/configure-json-file.md) file. For example:


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
        ],
        "allowDerivedImages": true
      },
      "commandList": {
        "type": "deny",
        "commands": ["push"]
      }
    }
  }
}
```

> [!TIP]
>
> You can now also configure these settings in the [Docker Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md).

As shown above, there are two configurations for bind-mounting the Docker
socket into containers: the `imageList` and the `commandList`. These are
described below.

### Image list

The `imageList` is a list of container images that are allowed to bind-mount the
Docker socket. By default the list is empty, no containers are allowed to
bind-mount the Docker socket when ECI is enabled. However, an administrator can add
images to the list, using either of these formats:

| Image Reference Format  | Description |
| :---------------------- | :---------- |
| `<image_name>[:<tag>]`  | Name of the image, with optional tag. If the tag is omitted, the `:latest` tag is used. If the tag is the wildcard `*`, then it means "any tag for that image." |
| `<image_name>@<digest>` | Name of the image, with a specific repository digest (e.g., as reported by `docker buildx imagetools inspect <image>`). This means only the image that matches that name and digest is allowed. |

The image name follows the standard convention, so it can point to any registry
and repository.

In the previous example, the image list was configured with three images:

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

> [!TIP]
>
> Be restrictive with the images you allow, as described in [Recommendations](#recommendations).

In general, it's easier to specify the image using the tag wildcard format, for example `<image-name>:*`, because then `imageList` doesn't need to be updated whenever a new version of the
image is used. Alternatively, you can use an immutable tag, for example `:latest`,
but it does not always work as well as the wildcard because, for example,
Testcontainers uses specific versions of the image, not necessarily the latest
one.

When ECI is enabled, Docker Desktop periodically downloads the image digests
for the allowed images from the appropriate registry and stores them in
memory. Then, when a container is started with a Docker socket bind-mount,
Docker Desktop checks if the container's image digest matches one of the allowed
digests. If so, the container is allowed to start, otherwise it's blocked.

Due to the digest comparison, it's not possible to bypass the Docker socket
mount permissions by re-tagging a disallowed image to the name of an allowed
one. In other words, if a user does:

```console
$ docker image rm <allowed_image>
$ docker tag <disallowed_image> <allowed_image>
$ docker run -v /var/run/docker.sock:/var/run/docker.sock <allowed_image>
```

then the tag operation succeeds, but the `docker run` command fails
because the image digest of the disallowed image won't match that of the allowed
ones in the repository.

### Docker Socket Mount Permissions for derived images

{{< summary-bar feature_name="Docker Scout Mount Permissions" >}}

As described in the prior section, administrators can configure the list of container
images that are allowed to mount the Docker socket via the `imageList`.

This works for most scenarios, but not always, because it requires knowing upfront
the name of the image(s) on which the Docker socket mounts should be allowed.
Some container tools such as [Paketo](https://paketo.io/) buildpacks,
build ephemeral local images that require Docker socket bind mounts. Since the name of
those ephemeral images is not known upfront, the `imageList` is not sufficient.

To overcome this, starting with Docker Desktop version 4.34, the Docker Socket
mount permissions not only apply to the images listed in the `imageList`; they
also apply to any local images derived (i.e., built from) an image in the
`imageList`.

That is, if a local image called "myLocalImage" is built from "myBaseImage"
(i.e., has a Dockerfile with a `FROM myBaseImage`), then if "myBaseImage" is in
the `imageList`, both "myBaseImage" and "myLocalImage" are allowed to mount the
Docker socket.

For example, to enable Paketo buildpacks to work with Docker Desktop and ECI,
simply add the following image to the `imageList`:

```json
"imageList": {
  "images": [
    "paketobuildpacks/builder:base"
  ],
  "allowDerivedImages": true
}
```

When the buildpack runs, it will create an ephemeral image derived from
`paketobuildpacks/builder:base` and mount the Docker socket to it. ECI will
allow this because it will notice that the ephemeral image is derived from an
allowed image.

The behavior is disabled by default and must be explicitly enabled by setting
`"allowDerivedImages": true` as shown above. In general it is recommended that
you disable this setting unless you know it's required.

A few caveats:

* Setting `"allowedDerivedImages" :true` will impact the startup time of
  containers by up to 1 extra second, as Docker Desktop needs to perform
  some more checks on the container image.

* The `allowDerivedImages` setting only applies to local-only images built from
  an allowed image. That is, the derived image must not be present in a remote
  repository because if it were, you would just list its name in the `imageList`.

* For derived image checking to work, the parent image (i.e., the image in the
  `imageList`) must be present locally (i.e., must have been explicitly pulled
  from a repository). This is usually not a problem as the tools that need this
  feature (e.g., Paketo buildpacks) will do the pre-pull of the parent image.

* For Docker Desktop versions 4.34 and 4.35 only: The `allowDerivedImages` setting
  applies to all images in the `imageList` specified with an explicit tag (e.g.,
  `<name>:<tag>`). It does not apply to images specified using the tag wildcard
  (e.g., `<name>:*`) described in the prior section. In Docker Desktop 4.36 and
  later, this caveat no longer applies, meaning that the `allowDerivedImages`
  settings applies to images specified with or without a wildcard tag. This
  makes it easier to manage the ECI Docker socket image list.

### Allowing all containers to mount the Docker socket

In Docker Desktop version 4.36 and later, it's possible to configure the image
list to allow any container to mount the Docker socket. You do this by adding
`"*"` to the `imageList`:

```json
"imageList": {
  "images": [
    "*"
  ]
}
```

This tells Docker Desktop to allow all containers to mount the Docker socket
which increases flexibility but reduces security. It also improves container
startup time when using Enhanced Container Isolation.

It is recommended that you use this only in scenarios where explicitly listing
allowed container images is not flexible enough.

### Command list

In addition to the `imageList` described in the prior sections, ECI can further
restrict the commands that a container can issue via a bind mounted Docker
socket. This is done via the Docker socket mount permission `commandList`, and
acts as a complementary security mechanism to the `imageList` (i.e., like a
second line of defense).

For example, say the `imageList` is configured to allow image `docker:cli` to
mount the Docker socket, and a container is started with it:

```console
$ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock sh
/ #
```

By default, this allows the container to issue any command via that Docker
socket (e.g., build and push images to the organization's repositories), which
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
  execute. For example, for local testing (e.g., Testcontainers), containers
  that bind-mount the Docker socket typically create / run / remove containers,
  volumes, and networks, but don't typically build images or push them into
  repositories (though some may legitimately do this). What commands to allow or
  block depends on the use case.

  - Note that all "docker" commands issued by the container via the bind-mounted
    Docker socket will also execute under enhanced container isolation (i.e.,
    the resulting container uses a the Linux user-namespace, sensitive system
    calls are vetted, etc.)

### Caveats and limitations

* When Docker Desktop is restarted, it's possible that an image that is allowed
  to mount the Docker socket is unexpectedly blocked from doing so. This can
  happen when the image digest changes in the remote repository (e.g., a
  ":latest" image was updated) and the local copy of that image (e.g., from a
  prior `docker pull`) no longer matches the digest in the remote repository. In
  this case, remove the local image and pull it again (e.g., `docker rm <image>`
  and `docker pull <image>`).

* It's not possible to allow Docker socket bind-mounts on containers using
  local-only images (i.e., images that are not on a registry) unless they are
  [derived from an allowed image](#docker-socket-mount-permissions-for-derived-images)
  or you've [allowed all containers to mount the Docker socket](#allowing-all-containers-to-mount-the-docker-socket).
  That is because Docker Desktop pulls the digests for the allowed images from
  the registry, and then uses that to compare against the local copy of the
  image.

* The `commandList` configuration applies to all containers that are allowed to
  bind-mount the Docker socket. Therefore it can't be configured differently per
  container.

* The following commands are not yet supported in the `commandList`:

| Unsupported command  | Description |
| :------------------- | :---------- |
| `compose`              | Docker Compose |
| `dev`                  | Dev environments |
| `extension`            | Manages Docker Extensions |
| `feedback`             | Send feedback to Docker |
| `init`                 | Creates Docker-related starter files |
| `manifest`             | Manages Docker image manifests |
| `plugin`              | Manages plugins |
| `sbom`                 | View Software Bill of Materials (SBOM) |
| `scout`                | Docker Scout |
| `trust`                | Manage trust on Docker images |

> [!NOTE]
>
> Docker socket mount permissions do not apply when running "true"
> Docker-in-Docker (i.e., when running the Docker Engine inside a container). In
> this case there's no bind-mount of the host's Docker socket into the
> container, and therefore no risk of the container leveraging the configuration
> and credentials of the host's Docker Engine to perform malicious activity.
> Enhanced Container Isolation is capable of running Docker-in-Docker securely,
> without giving the outer container true root permissions in the Docker Desktop
> VM.
