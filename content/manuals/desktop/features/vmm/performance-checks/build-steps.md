---
title: Time the per-step cost of an image build
linkTitle: Build with 20 steps
description: Measure the fixed cost of each Dockerfile instruction by timing a one-line build and a 20-step build.
keywords: docker desktop, performance, benchmark, docker build, buildkit, dockerfile, layers
weight: 50
---

Each Dockerfile instruction makes the builder start a container, run a command,
and save the result as a layer. This check measures that fixed cost per
instruction by timing two builds and subtracting:

```text
(20-step build - one-line build) / 20
```

The first build is a single `FROM` line. Its duration is the fixed cost of
starting the builder, preparing the base image, and writing the image. The
second build adds 20 instructions that each write a tiny file. Subtracting the
first duration from the second and dividing by 20 removes the fixed costs.

Both builds start from an empty build cache. Each duration is the wall-clock time
of a `docker build` command on your machine. Each pair of builds gives one
result.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/).

> [!WARNING]
>
> This check runs `docker buildx prune --all --force`, which deletes your entire
> local build cache. Don't run it on a machine where you need that cache.

## Set up

```console
$ mkdir -p ~/perfcheck/twenty-steps && cd ~/perfcheck/twenty-steps

$ printf 'FROM alpine:3.24\n' > Dockerfile.baseline

$ printf 'FROM alpine:3.24\n' > Dockerfile.twentysteps
$ i=1
$ while [ "$i" -le 20 ]; do
  printf 'RUN mkdir -p /steps && printf "%02d" > /steps/%02d\n' "$i" "$i" >> Dockerfile.twentysteps
  i=$((i + 1))
done

$ printf 'Dockerfile.baseline\nDockerfile.twentysteps\n.dockerignore\n' > .dockerignore

$ docker pull alpine:3.24
```

Pulling the base image keeps its download out of both timed builds. Stay in this
directory for the runs.

Confirm that `Dockerfile.twentysteps` has exactly 20 `RUN` lines:

```console
$ grep -c '^RUN ' Dockerfile.twentysteps
20
```

## Run the check

1. Remove the previous images and empty the build cache. This step isn't timed.
   On the first run, the removal reports `No such image`. That is expected.

   ```console
   $ docker image rm --force perfcheck/twenty-baseline:current perfcheck/twenty-steps:current
   $ docker buildx prune --all --force
   ```

2. Time the one-line build:

   ```console
   $ time docker build --pull=false --progress=plain \
     -f Dockerfile.baseline -t perfcheck/twenty-baseline:current .
   ```

3. Remove that image and empty the build cache again, so the second build also
   starts with nothing. This step isn't timed.

   ```console
   $ docker image rm --force perfcheck/twenty-baseline:current
   $ docker buildx prune --all --force
   ```

4. Time the 20-step build:

   ```console
   $ time docker build --pull=false --progress=plain \
     -f Dockerfile.twentysteps -t perfcheck/twenty-steps:current .
   ```

5. Check that the run is valid, then clean up:

   ```console
   $ docker image inspect --format '{{.Id}}' perfcheck/twenty-steps:current
   $ docker image rm --force perfcheck/twenty-steps:current
   ```

6. Calculate the result in milliseconds:

   ```text
   (step 4 duration - step 2 duration) / 20
   ```

   A negative result means something went wrong, most often that the base image
   downloaded during a timed build.

Calculate the result once for each pair of builds. Don't average several builds
and divide once.
