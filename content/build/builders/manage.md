---
title: Manage builders
keywords: build, buildx, builders, buildkit, drivers, backend
description: null
---

You can create, inspect, and manage builders using `docker buildx` commands,
or [using Docker Desktop](#manage-builders-with-docker-desktop).

## Create a new builder

The default builder uses the [`docker` driver](drivers/docker.md).
You can't manually create new `docker` builders, but you can create builders
that use other drivers, such as the
[`docker-container` driver](drivers/docker-container.md),
which runs the BuildKit daemon in a container.

Use the
[`docker buildx create`](../../reference/cli/docker/buildx/create.md)
command to create a builder.

```console
$ docker buildx create --name=<builder-name>
```

Buildx uses the `docker-container` driver by default if you omit the `--driver`
flag. For more information about available drivers, see
[Build drivers](drivers/_index.md).

## List available builders

Use `docker buildx ls` to see builder instances available on your system, and
the drivers they're using.

```console
$ docker buildx ls
NAME/NODE       DRIVER/ENDPOINT      STATUS   BUILDKIT PLATFORMS
default *       docker
  default       default              running  v0.11.6  linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
my_builder      docker-container
  my_builder0   default              running  v0.11.6  linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
```

The asterisk (`*`) next to the builder name indicates the
[selected builder](_index.md#selected-builder).

## Inspect a builder

To inspect a builder with the CLI, use `docker buildx inspect <name>`.
You can only inspect a builder if the builder is active.
You can add the `--bootstrap` flag to the command to start the builder.

```console
$ docker buildx inspect --bootstrap my_builder
[+] Building 1.7s (1/1) FINISHED                                                                  
 => [internal] booting buildkit                                                              1.7s
 => => pulling image moby/buildkit:buildx-stable-1                                           1.3s
 => => creating container buildx_buildkit_my_builder0                                        0.4s
Name:          my_builder
Driver:        docker-container
Last Activity: 2023-06-21 18:28:37 +0000 UTC

Nodes:
Name:      my_builder0
Endpoint:  unix:///var/run/docker.sock
Status:    running
Buildkit:  v0.11.6
Platforms: linux/arm64, linux/amd64, linux/amd64/v2, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/mips64le, linux/mips64, linux/arm/v7, linux/arm/v6
```

If you want to see how much disk space a builder is using, use the
`docker buildx du` command. By default, this command shows the total disk usage
for all available builders. To see usage for a specific builder, use the
`--builder` flag.

```console
$ docker buildx du --builder my_builder
ID                                        RECLAIMABLE SIZE        LAST ACCESSED
olkri5gq6zsh8q2819i69aq6l                 true        797.2MB     37 seconds ago
6km4kasxgsywxkm6cxybdumbb*                true        438.5MB     36 seconds ago
qh3wwwda7gx2s5u4hsk0kp4w7                 true        213.8MB     37 seconds ago
54qq1egqem8max3lxq6180cj8                 true        200.2MB     37 seconds ago
ndlp969ku0950bmrw9muolw0c*                true        116.7MB     37 seconds ago
u52rcsnfd1brwc0chwsesb3io*                true        116.7MB     37 seconds ago
rzoeay0s4nmss8ub59z6lwj7d                 true        46.25MB     4 minutes ago
itk1iibhmv7awmidiwbef633q                 true        33.33MB     37 seconds ago
4p78yqnbmgt6xhcxqitdieeln                 true        19.46MB     4 minutes ago
dgkjvv4ay0szmr9bl7ynla7fy*                true        19.24MB     36 seconds ago
tuep198kmcw299qc9e4d1a8q2                 true        8.663MB     4 minutes ago
n1wzhauk9rpmt6ib1es7dktvj                 true        20.7kB      4 minutes ago
0a2xfhinvndki99y69157udlm                 true        16.56kB     37 seconds ago
gf0z1ypz54npfererqfeyhinn                 true        16.38kB     37 seconds ago
nz505f12cnsu739dw2pw0q78c                 true        8.192kB     37 seconds ago
hwpcyq5hdfvioltmkxu7fzwhb*                true        8.192kB     37 seconds ago
acekq89snc7j6im1rjdizvsg1*                true        8.192kB     37 seconds ago
Reclaimable:  2.01GB
Total:        2.01GB
```

## Remove a builder

Use the
[`docker buildx remove`](../../reference/cli/docker/buildx/create.md)
command to remove a builder.

```console
$ docker buildx rm <builder-name>
```

If you remove your currently selected builder,
the default `docker` builder is automatically selected.
You can't remove the default builder.

Local build cache for the builder is also removed.

### Removing remote builders

Removing a remote builder doesn't affect the remote build cache.
It also doesn't stop the remote BuildKit daemon.
It only removes your connection to the builder.

## Manage builders with Docker Desktop

If you have turned on the
[Docker Desktop Builds view](../../desktop/use-desktop/builds.md),
you can inspect builders in Docker Desktop settings. See:

- [Change settings, Windows](../../desktop/settings/windows.md#builders)
- [Change settings, Mac](../../desktop/settings/mac.md#builders)
- [Change settings, Linux](../../desktop/settings/linux.md#builders)
