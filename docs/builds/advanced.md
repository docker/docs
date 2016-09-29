<!--[metadata]>
+++
title = "Advanced options for Autobuild and Autotest"
description = "Automated builds"
keywords = ["automated, build, images"]
[menu.main]
parent="builds"
weight=-40
+++
<![end-metadata]-->

# Advanced options for Autobuild and Autotest

The following options allow you to customize your automated build and automated test processes.

## Environment variables for building and testing

Several utility environment variables are set by the build process, and are
available during automated builds, automated tests, and while executing
hooks.

> **Note**: These environment variables are only available to the build and test
processes and do not affect your service's run environment.

* `SOURCE_BRANCH`: the name of the branch or the tag that is currently being tested.
* `SOURCE_COMMIT`: the SHA1 hash of the commit being tested.
* `COMMIT_MSG`: the message from the commit being tested and built.
* `DOCKER_REPO`: the name of the Docker repository being built.
* `DOCKER_TAG`: the Docker repository tag being built.
* `IMAGE_NAME`: the name and tag of the Docker repository being built. (This variable is a combination of `DOCKER_REPO`/`DOCKER_TAG`.)

If you are using these environment variables in a `docker-compose.test.yml` file
for automated testing, declare them in your `sut` service's environment as shows
below.

```yml
sut:
  build: .
  command: run_tests.sh
  environment:
    - SOURCE_BRANCH
```

## Custom build phase hooks

You can run custom commands between phases of the build process by creating
hooks. Hooks allow you to provide extra instructions to the autobuild and
autotest processes.

Create a folder called `hooks` in your source code repository at the same
directory level as your Dockerfile. Place files that define the hooks in that
folder. The builder executes them before and after each step.

The following hooks are available:

* `hooks/post_checkout`
* `hooks/pre_build`
* `hooks/post_build`
* `hooks/pre_test`
* `hooks/post_test`
* `hooks/pre_push` (only used when executing a build rule or [automated build](automated-build.md) )
* `hooks/post_push` (only used when executing a build rule or [automated build](automated-build.md) )


## Override build, test or push commands

In addition to the custom build phase hooks above, you can also use
`hooks/build`, `hooks/test`, and `hooks/push` to override and customize the
`build`, `test` and `push` commands during automated build and test processes.

**Use these hooks with caution.** The contents of these hook files replace the
basic `docker` commands, so you must include a similar build, test or push
command in the hook or your automated process will not complete.
