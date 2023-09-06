---
description: Automated builds
keywords: automated, build, images
title: Advanced options for Autobuild and Autotest
aliases:
- /docker-cloud/builds/advanced/
---

> **Note**
>
> Automated builds require a
> [Docker Pro, Team, or Business subscription](../../subscription/index.md).

The following options allow you to customize your automated build and automated
test processes.

## Environment variables for building and testing

Several utility environment variables are set by the build process, and are
available during automated builds, automated tests, and while executing
hooks.

> **Note**
>
> These environment variables are only available to the build and test
processes and don't affect your service's run environment.

* `SOURCE_BRANCH`: the name of the branch or the tag that is currently being tested.
* `SOURCE_COMMIT`: the SHA1 hash of the commit being tested.
* `COMMIT_MSG`: the message from the commit being tested and built.
* `DOCKER_REPO`: the name of the Docker repository being built.
* `DOCKERFILE_PATH`: the dockerfile currently being built.
* `DOCKER_TAG`: the Docker repository tag being built.
* `IMAGE_NAME`: the name and tag of the Docker repository being built. (This variable is a combination of `DOCKER_REPO`:`DOCKER_TAG`.)

If you are using these build environment variables in a
`docker-compose.test.yml` file for automated testing, declare them in your `sut`
service's environment as shown below.

```yaml
services:
  sut:
    build: .
    command: run_tests.sh
    environment:
      - SOURCE_BRANCH
```


## Override build, test or push commands

Docker Hub allows you to override and customize the `build`, `test` and `push`
commands during automated build and test processes using hooks. For example, you
might use a build hook to set build arguments used only during the build
process. You can also set up [custom build phase hooks](#custom-build-phase-hooks)
to perform actions in between these commands.

> **Important**
>
>Use these hooks with caution. The contents of these hook files replace the
basic `docker` commands, so you must include a similar build, test or push
command in the hook or your automated process does not complete.
{ .important }

To override these phases, create a folder called `hooks` in your source code
repository at the same directory level as your Dockerfile. Create a file called
`hooks/build`, `hooks/test`, or `hooks/push` and include commands that the
builder process can execute, such as `docker` and `bash` commands (prefixed
appropriately with `#!/bin/bash`).

These hooks run on an instance of [Ubuntu](https://releases.ubuntu.com/),
which includes interpreters
such as Perl or Python, and utilities such as `git` or `curl`. Refer to the
[Ubuntu documentation](https://ubuntu.com/)
for the full list of available interpreters and utilities.

## Custom build phase hooks

You can run custom commands between phases of the build process by creating
hooks. Hooks allow you to provide extra instructions to the autobuild and
autotest processes.

Create a folder called `hooks` in your source code repository at the same
directory level as your Dockerfile. Place files that define the hooks in that
folder. Hook files can include both `docker` commands, and `bash` commands as
long as they are prefixed appropriately with `#!/bin/bash`. The builder executes
the commands in the files before and after each step.

The following hooks are available:

* `hooks/post_checkout`
* `hooks/pre_build`
* `hooks/post_build`
* `hooks/pre_test`
* `hooks/post_test`
* `hooks/pre_push` (only used when executing a build rule or [Automated build](index.md) )
* `hooks/post_push` (only used when executing a build rule or [Automated build](index.md) )

### Build hook examples

#### Override the "build" phase to set variables

Docker Hub allows you to define build environment variables either in the hook
files, or from the automated build interface, which you can then reference in hooks.

The following example defines a build hook that uses `docker build` arguments to
set the variable `CUSTOM` based on the value of variable defined using the
Docker Hub build settings. `$DOCKERFILE_PATH` is a variable that you provide
with the name of the Dockerfile you want to build, and `$IMAGE_NAME` is the name
of the image being built.

```console
$ docker build --build-arg CUSTOM=$VAR -f $DOCKERFILE_PATH -t $IMAGE_NAME .
```

> **Important**
>
> A `hooks/build` file overrides the basic [docker build](../../engine/reference/commandline/build.md) command used by the builder, so you must include a similar build command in the hook or
the automated build fails.
{ .important }

Refer to the [docker build documentation](../../engine/reference/commandline/build.md#build-arg)
to learn more about Docker build-time variables.

#### Push to multiple repositories

By default the build process pushes the image only to the repository where the
build settings are configured. If you need to push the same image to multiple
repositories, you can set up a `post_push` hook to add additional tags and push
to more repositories.

```console
$ docker tag $IMAGE_NAME $DOCKER_REPO:$SOURCE_COMMIT
$ docker push $DOCKER_REPO:$SOURCE_COMMIT
```

## Source repository or branch clones

When Docker Hub pulls a branch from a source code repository, it performs
a shallow clone, it clones only the tip of the specified branch. This has the advantage
of minimizing the amount of data transfer necessary from the repository and
speeding up the build because it pulls only the minimal code necessary.

As a result, if you need to perform a custom action that relies on a different
branch, such as a `post_push` hook, you can't checkout that branch unless
you do one of the following:

* You can get a shallow checkout of the target branch by doing the following:

    ```console
    $ git fetch origin branch:mytargetbranch --depth 1
    ```

* You can also "unshallow" the clone, which fetches the whole Git history (and
  potentially takes a long time / moves a lot of data) by using the `--unshallow`
  flag on the fetch:

    ```console
    $ git fetch --unshallow origin
    ```