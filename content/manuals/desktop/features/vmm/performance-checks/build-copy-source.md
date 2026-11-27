---
title: Time an image build that copies a large source tree
linkTitle: Build with large COPY
description: Measure how long docker build takes to send and copy a 589 MB source tree into an image.
keywords: docker desktop, performance, benchmark, docker build, buildkit, build context, copy
weight: 40
---

Most `docker build` runs start by sending your source directory to the builder
and copying it into the image. This check measures only that step: an image
whose entire content is a 589 MB source tree, with no compiler and no package
downloads.

The result is the wall-clock time of the single `docker build` command on your
machine, starting from an empty build cache.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/_index.md).

> [!WARNING]
>
> Step 3 runs `docker buildx prune --all --force`, which deletes your entire
> local build cache. Every later build on the machine starts from scratch. Don't
> run this check on a machine where you need that cache.

## Get the source

```console
$ mkdir -p ~/perfcheck/copy-build && cd ~/perfcheck/copy-build

$ git init --bare -q posthog-mirror.git
$ git -C posthog-mirror.git remote add origin https://github.com/PostHog/posthog.git
$ git -C posthog-mirror.git fetch -q --depth=1 origin 4fa8e10ae2e988fd64ce93c8408856100af16c4f
$ git -C posthog-mirror.git update-ref refs/heads/pinned FETCH_HEAD

$ git -C posthog-mirror.git rev-parse refs/heads/pinned
4fa8e10ae2e988fd64ce93c8408856100af16c4f
```

## Run the check

1. Create a fresh checkout. This step isn't timed.

   ```console
   $ rm -rf fixture && mkdir fixture
   $ git clone -q --no-local --no-checkout posthog-mirror.git fixture/repo
   $ git -C fixture/repo checkout -q --detach 4fa8e10ae2e988fd64ce93c8408856100af16c4f
   ```

2. Write the build files into the checkout. `FROM scratch` means the image has
   no base layer to fetch, so the build uses no network. The `.dockerignore`
   file keeps Git's data out of the build context.

   ```console
   $ cd fixture/repo
   $ printf 'FROM scratch\nCOPY . /src\n' > Dockerfile.copy
   $ printf '.git\nDockerfile.copy\n.dockerignore\n' > .dockerignore
   ```

3. Remove the previous image and empty the build cache. This step isn't timed.
   On the first run, the image removal reports `No such image`. That is
   expected.

   ```console
   $ docker image rm --force perfcheck/posthog-copy:current
   $ docker buildx prune --all --force
   ```

4. Run the timed command:

   ```console
   $ time docker build --pull=false --progress=plain \
     -f Dockerfile.copy -t perfcheck/posthog-copy:current .
   ```

   Record the elapsed time.

5. Check that the run is valid. The image must exist:

   ```console
   $ docker image inspect --format '{{.Id}}' perfcheck/posthog-copy:current
   ```

6. Clean up. Remove the 589 MB image so it doesn't put the disk under pressure
   and affect later measurements.

   ```console
   $ docker image rm --force perfcheck/posthog-copy:current
   $ cd ../..
   $ rm -rf fixture
   ```

Report whether the engine uses the containerd image store alongside your result.
The two image stores write images differently, which affects this check.
