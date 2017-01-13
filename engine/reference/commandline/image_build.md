---
datafolder: engine-cli
datafile: docker_image_build
title: docker image build
---

<!--
Sorry, but the contents of this page are automatically generated from
Docker's source code. If you want to suggest a change to the text that appears
here, you'll need to find the string by searching this repo:

https://www.github.com/docker/docker
-->

{% include cli.md %}

## Examples

### Build with PATH

```bash
$ docker build .

Uploading context 10240 bytes
Step 1/3 : FROM busybox
Pulling repository busybox
 ---> e9aa60c60128MB/2.284 MB (100%) endpoint: https://cdn-registry-1.docker.io/v1/
Step 2/3 : RUN ls -lh /
 ---> Running in 9c9e81692ae9
total 24
drwxr-xr-x    2 root     root        4.0K Mar 12  2013 bin
drwxr-xr-x    5 root     root        4.0K Oct 19 00:19 dev
drwxr-xr-x    2 root     root        4.0K Oct 19 00:19 etc
drwxr-xr-x    2 root     root        4.0K Nov 15 23:34 lib
lrwxrwxrwx    1 root     root           3 Mar 12  2013 lib64 -> lib
dr-xr-xr-x  116 root     root           0 Nov 15 23:34 proc
lrwxrwxrwx    1 root     root           3 Mar 12  2013 sbin -> bin
dr-xr-xr-x   13 root     root           0 Nov 15 23:34 sys
drwxr-xr-x    2 root     root        4.0K Mar 12  2013 tmp
drwxr-xr-x    2 root     root        4.0K Nov 15 23:34 usr
 ---> b35f4035db3f
Step 3/3 : CMD echo Hello world
 ---> Running in 02071fceb21b
 ---> f52f38b7823e
Successfully built f52f38b7823e
Removing intermediate container 9c9e81692ae9
Removing intermediate container 02071fceb21b
```

This example specifies that the `PATH` is `.`, and so all the files in the
local directory get `tar`d and sent to the Docker daemon. The `PATH` specifies
where to find the files for the "context" of the build on the Docker daemon.
Remember that the daemon could be running on a remote machine and that no
parsing of the Dockerfile happens at the client side (where you're running
`docker build`). That means that *all* the files at `PATH` get sent, not just
the ones listed to [*ADD*](../builder.md#add) in the Dockerfile.

The transfer of context from the local machine to the Docker daemon is what the
`docker` client means when you see the "Sending build context" message.

If you wish to keep the intermediate containers after the build is complete,
you must use `--rm=false`. This does not affect the build cache.

### Build with URL

```bash
$ docker build github.com/creack/docker-firefox
```

This will clone the GitHub repository and use the cloned repository as context.
The Dockerfile at the root of the repository is used as Dockerfile. You can
specify an arbitrary Git repository by using the `git://` or `git@` scheme.

```bash
$ docker build -f ctx/Dockerfile http://server/ctx.tar.gz

Downloading context: http://server/ctx.tar.gz [===================>]    240 B/240 B
Step 1/3 : FROM busybox
 ---> 8c2e06607696
Step 2/3 : ADD ctx/container.cfg /
 ---> e7829950cee3
Removing intermediate container b35224abf821
Step 3/3 : CMD /bin/ls
 ---> Running in fbc63d321d73
 ---> 3286931702ad
Removing intermediate container fbc63d321d73
Successfully built 377c409b35e4
```

This sends the URL `http://server/ctx.tar.gz` to the Docker daemon, which
downloads and extracts the referenced tarball. The `-f ctx/Dockerfile`
parameter specifies a path inside `ctx.tar.gz` to the `Dockerfile` that is used
to build the image. Any `ADD` commands in that `Dockerfile` that refers to local
paths must be relative to the root of the contents inside `ctx.tar.gz`. In the
example above, the tarball contains a directory `ctx/`, so the `ADD
ctx/container.cfg /` operation works as expected.

### Build with -

```bash
$ docker build - < Dockerfile
```

This will read a Dockerfile from `STDIN` without context. Due to the lack of a
context, no contents of any local directory will be sent to the Docker daemon.
Since there is no context, a Dockerfile `ADD` only works if it refers to a
remote URL.

```bash
$ docker build - < context.tar.gz
```

This will build an image for a compressed context read from `STDIN`.  Supported
formats are: bzip2, gzip and xz.

### Usage of .dockerignore

```bash
$ docker build .

Uploading context 18.829 MB
Uploading context
Step 1/2 : FROM busybox
 ---> 769b9341d937
Step 2/2 : CMD echo Hello world
 ---> Using cache
 ---> 99cc1ad10469
Successfully built 99cc1ad10469
$ echo ".git" > .dockerignore
$ docker build .
Uploading context  6.76 MB
Uploading context
Step 1/2 : FROM busybox
 ---> 769b9341d937
Step 2/2 : CMD echo Hello world
 ---> Using cache
 ---> 99cc1ad10469
Successfully built 99cc1ad10469
```

This example shows the use of the `.dockerignore` file to exclude the `.git`
directory from the context. Its effect can be seen in the changed size of the
uploaded context. The builder reference contains detailed information on
[creating a .dockerignore file](../builder.md#dockerignore-file)

### Tag image (-t)

```bash
$ docker build -t vieux/apache:2.0 .
```

This will build like the previous example, but it will then tag the resulting
image. The repository name will be `vieux/apache` and the tag will be `2.0`.
[Read more about valid tags](tag.md).

You can apply multiple tags to an image. For example, you can apply the `latest`
tag to a newly built image and add another tag that references a specific
version.
For example, to tag an image both as `whenry/fedora-jboss:latest` and
`whenry/fedora-jboss:v2.1`, use the following:

```bash
$ docker build -t whenry/fedora-jboss:latest -t whenry/fedora-jboss:v2.1 .
```
### Specify Dockerfile (-f)

```bash
$ docker build -f Dockerfile.debug .
```

This will use a file called `Dockerfile.debug` for the build instructions
instead of `Dockerfile`.

```bash
$ docker build -f dockerfiles/Dockerfile.debug -t myapp_debug .
$ docker build -f dockerfiles/Dockerfile.prod  -t myapp_prod .
```

The above commands will build the current build context (as specified by the
`.`) twice, once using a debug version of a `Dockerfile` and once using a
production version.

```bash
$ cd /home/me/myapp/some/dir/really/deep
$ docker build -f /home/me/myapp/dockerfiles/debug /home/me/myapp
$ docker build -f ../../../../dockerfiles/debug /home/me/myapp
```

These two `docker build` commands do the exact same thing. They both use the
contents of the `debug` file instead of looking for a `Dockerfile` and will use
`/home/me/myapp` as the root of the build context. Note that `debug` is in the
directory structure of the build context, regardless of how you refer to it on
the command line.

> **Note:**
> `docker build` will return a `no such file or directory` error if the
> file or directory does not exist in the uploaded context. This may
> happen if there is no context, or if you specify a file that is
> elsewhere on the Host system. The context is limited to the current
> directory (and its children) for security reasons, and to ensure
> repeatable builds on remote Docker hosts. This is also the reason why
> `ADD ../file` will not work.

### Optional parent cgroup (--cgroup-parent)

When `docker build` is run with the `--cgroup-parent` option the containers
used in the build will be run with the [corresponding `docker run`
flag](../run.md#specifying-custom-cgroups).

### Set ulimits in container (--ulimit)

Using the `--ulimit` option with `docker build` will cause each build step's
container to be started using those [`--ulimit`
flag values](./run.md#set-ulimits-in-container-ulimit).

### Set build-time variables (--build-arg)

You can use `ENV` instructions in a Dockerfile to define variable
values. These values persist in the built image. However, often
persistence is not what you want. Users want to specify variables differently
depending on which host they build an image on.

A good example is `http_proxy` or source versions for pulling intermediate
files. The `ARG` instruction lets Dockerfile authors define values that users
can set at build-time using the  `--build-arg` flag:

```bash
$ docker build --build-arg HTTP_PROXY=http://10.20.30.2:1234 .
```

This flag allows you to pass the build-time variables that are
accessed like regular environment variables in the `RUN` instruction of the
Dockerfile. Also, these values don't persist in the intermediate or final images
like `ENV` values do.

Using this flag will not alter the output you see when the `ARG` lines from the
Dockerfile are echoed during the build process.

For detailed information on using `ARG` and `ENV` instructions, see the
[Dockerfile reference](../builder.md).

### Optional security options (--security-opt)

This flag is only supported on a daemon running on Windows, and only supports
the `credentialspec` option. The `credentialspec` must be in the format
`file://spec.txt` or `registry://keyname`.

### Specify isolation technology for container (--isolation)

This option is useful in situations where you are running Docker containers on
Windows. The `--isolation=<value>` option sets a container's isolation
technology. On Linux, the only supported is the `default` option which uses
Linux namespaces. On Microsoft Windows, you can specify these values:


| Value     | Description                                                                                                                                                   |
|-----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `default` | Use the value specified by the Docker daemon's `--exec-opt` . If the `daemon` does not specify an isolation technology, Microsoft Windows uses `process` as its default value.  |
| `process` | Namespace isolation only.                                                                                                                                     |
| `hyperv`  | Hyper-V hypervisor partition-based isolation.                                                                                                                 |

Specifying the `--isolation` flag without a value is the same as setting `--isolation="default"`.


### Squash an image's layers (--squash) **Experimental Only**

Once the image is built, squash the new layers into a new image with a single
new layer. Squashing does not destroy any existing image, rather it creates a new
image with the content of the squashed layers. This effectively makes it look
like all `Dockerfile` commands were created with a single layer. The build
cache is preserved with this method.

**Note**: using this option means the new image will not be able to take
advantage of layer sharing with other images and may use significantly more
space.

**Note**: using this option you may see significantly more space used due to
storing two copies of the image, one for the build cache with all the cache
layers in tact, and one for the squashed version.
