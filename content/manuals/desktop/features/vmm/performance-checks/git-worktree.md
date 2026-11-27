---
title: Time git worktree add in a shared folder
linkTitle: Git worktree add
description: Measure how long git worktree add takes in a container when the repository and the new copy are in a folder shared from your machine.
keywords: docker desktop, performance, benchmark, git worktree, file sharing, virtiofs, bind mount
weight: 30
---

`git worktree add` creates a second checkout of a repository next to the first
one, writing tens of thousands of files. This check runs it inside a container,
with both the repository and the new copy in a folder shared in from your
machine, so every write goes through the engine's file sharing.

The repository is PostHog at a pinned commit, with about 47 000 files and 589 MB.
The result is the duration of the single `git worktree add` command in
milliseconds, measured inside the container.

For how durations are taken and how to repeat a check, see
[Performance checks](/manuals/desktop/features/vmm/performance-checks/_index.md).

## Get the source

```console
$ mkdir -p ~/perfcheck/worktree && cd ~/perfcheck/worktree

$ git init --bare -q posthog-mirror.git
$ git -C posthog-mirror.git remote add origin https://github.com/PostHog/posthog.git
$ git -C posthog-mirror.git fetch -q --depth=1 origin 4fa8e10ae2e988fd64ce93c8408856100af16c4f
$ git -C posthog-mirror.git update-ref refs/heads/pinned FETCH_HEAD

$ git -C posthog-mirror.git rev-parse refs/heads/pinned
4fa8e10ae2e988fd64ce93c8408856100af16c4f
```

## Run the check

1. Create a fresh checkout in an otherwise empty folder. This step isn't timed.
   `fixture/` must contain only `repo/`. The timed command creates
   `posthog-worktree/`.

   ```console
   $ rm -rf fixture && mkdir fixture
   $ git clone -q --no-local --no-checkout posthog-mirror.git fixture/repo
   $ git -C fixture/repo checkout -q --detach 4fa8e10ae2e988fd64ce93c8408856100af16c4f
   ```

2. Share the whole `fixture` folder into a container and install Git. Both the
   repository and the destination must be on the shared path. A destination in
   the container's own storage measures a different filesystem.

   ```console
   $ docker run -d --rm --name perfcheck-worktree \
     -v "$PWD/fixture:/fixture" -w /fixture/repo \
     alpine:3.24.1 sleep infinity

   $ docker exec perfcheck-worktree apk add --no-cache git=~2.54.0 coreutils

   $ docker exec perfcheck-worktree git --version
   git version 2.54.0
   ```

3. Define a helper that runs commands in the container. On a Linux engine that
   isn't rootless, the helper runs commands as your user, so the files the next
   step creates stay yours. On Mac, Windows, and rootless Linux engines, it
   runs commands unchanged.

   ```bash
   in_container() {
     if [ "$(uname -s)" = Linux ] &&
        ! docker info --format '{{.SecurityOptions}}' | grep -q name=rootless; then
       docker exec --user "$(id -u):$(id -g)" --env HOME=/tmp perfcheck-worktree "$@"
     else
       docker exec perfcheck-worktree "$@"
     fi
   }
   ```

4. Run the timed command. The `--detach` flag keeps branch handling out of the
   measurement, so the file writes are the work being timed.

   ```console
   $ in_container sh -c '
     set -e
     start=$(date +%s%3N)
     git -c safe.directory=/fixture/repo -C /fixture/repo \
         worktree add --quiet --detach /fixture/posthog-worktree \
         4fa8e10ae2e988fd64ce93c8408856100af16c4f
     end=$(date +%s%3N)
     printf "worktree_add_ms=%s\n" "$((end - start))"
   '
   ```

   The value of `worktree_add_ms` is the result.

5. Check that the run is valid. This step isn't timed. The new copy must be at
   the pinned commit and clean.

   ```console
   $ in_container sh -c '
     git -c safe.directory=/fixture/posthog-worktree -C /fixture/posthog-worktree rev-parse HEAD
     git -c safe.directory=/fixture/posthog-worktree -C /fixture/posthog-worktree \
         status --porcelain=v1 --untracked-files=all | head
   '
   ```

   The first command must print `4fa8e10ae2e988fd64ce93c8408856100af16c4f`. The
   second must print nothing.

6. Clean up:

   ```console
   $ docker rm -f perfcheck-worktree
   $ rm -rf fixture
   ```
